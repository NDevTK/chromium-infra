# Copyright 2015 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

"""Task queue endpoints for creating and updating issues on issue tracker."""

import datetime
import json
import logging
import urllib2
import webapp2

from google.appengine.api import app_identity
from google.appengine.api import taskqueue
from google.appengine.api import urlfetch
from google.appengine.ext import ndb

from infra_libs import ts_mon
from issue_tracker import issue_tracker_api, issue
from model.flake import (
    Flake, FlakeOccurrence, FlakeUpdate, FlakeUpdateSingleton, FlakyRun)
from status import build_result, util
from test_results.util import normalize_test_type


MAX_UPDATED_ISSUES_PER_DAY = 50
MAX_TIME_DIFFERENCE_SECONDS = 12 * 60 * 60
MIN_REQUIRED_FLAKY_RUNS = 3
DAYS_TILL_STALE = 30
USE_MONORAIL = True
DAYS_TO_REOPEN_ISSUE = 3
FLAKY_RUNS_TEMPLATE = (
    'Detected %(new_flakes_count)d new flakes for test/step "%(name)s". To see '
    'the actual flakes, please visit %(flakes_url)s. This message was posted '
    'automatically by the chromium-try-flakes app. Since flakiness is ongoing, '
    'the issue was moved back into %(queue_name)s (unless already there).')
SUMMARY_TEMPLATE = '"%(name)s" is flaky'
DESCRIPTION_TEMPLATE = (
    '%(summary)s.\n\n'
    'This issue was created automatically by the chromium-try-flakes app. '
    'Please find the right owner to fix the respective test/step and assign '
    'this issue to them. %(other_queue_msg)s\n\n'
    'We have detected %(flakes_count)d recent flakes. List of all flakes can '
    'be found at %(flakes_url)s.\n\n'
    'Flaky tests should be disabled within 30 minutes unless culprit CL is '
    'found and reverted. Please see more details here: '
    'https://sites.google.com/a/chromium.org/dev/developers/tree-sheriffs/'
    'sheriffing-bug-queues#triaging-auto-filed-flakiness-bugs')
SHERIFF_QUEUE_MSG = (
    'If the step/test is infrastructure-related, please add Infra-Troopers '
    'label and change issue status to Untriaged. When done, please remove the '
    'issue from Sheriff Bug Queue by removing the Sheriff-Chromium label.')
TROOPER_QUEUE_MSG = (
    'If the step/test is not infrastructure-related (e.g. flaky test), please '
    'add Sheriff-Chromium label and change issue status to Untriaged. When '
    'done, please remove the issue from Trooper Bug Queue by removing the '
    'Infra-Troopers label.')
REOPENED_DESCRIPTION_TEMPLATE = (
    '%(description)s\n\n'
    'This flaky test/step was previously tracked in issue %(old_issue)d.')
FLAKES_URL_TEMPLATE = (
    'https://chromium-try-flakes.appspot.com/all_flake_occurrences?key=%s')
TEST_RESULTS_URL_TEMPLATE = (
    'http://test-results.appspot.com/testfile?builder=%(buildername)s&name='
    'full_results.json&master=%(mastername)s&testtype=%(stepname)s&buildnumber='
    '%(buildnumber)s')
VERY_STALE_FLAKES_MESSAGE = (
    'Reporting to stale-flakes-reports@google.com to investigate why this '
    'issue is not being processed despite being in an appropriate queue.')
STALE_FLAKES_ML = 'stale-flakes-reports@google.com'
MAX_GAP_FOR_FLAKINESS_PERIOD = datetime.timedelta(days=3)
KNOWN_TROOPER_FAILURES = [
    'analyze', 'bot_update', 'compile (with patch)', 'compile',
    'device_status_check', 'gclient runhooks (with patch)', 'Patch failure',
    'process_dumps', 'provision_devices', 'update_scripts']


def is_trooper_flake(flake_name):
  return flake_name in KNOWN_TROOPER_FAILURES


def get_queue_details(flake_name):
  if is_trooper_flake(flake_name):
    return 'Trooper Bug Queue', 'Infra-Troopers'
  else:
    return 'Sheriff Bug Queue', 'Sheriff-Chromium'


class ProcessIssue(webapp2.RequestHandler):
  reporting_delay = ts_mon.FloatMetric(
      'flakiness_pipeline/reporting_delay',
      description='The delay in seconds from the moment first flake occurrence '
                  'in this flakiness period happens and until the time an '
                  'issue is created to track it.')
  issue_updates = ts_mon.CounterMetric(
      'flakiness_pipeline/issue_updates',
      description='Issues updated/created')

  @ndb.transactional
  def _get_flake_update_singleton_key(self):
    singleton_key = ndb.Key('FlakeUpdateSingleton', 'singleton')
    if not singleton_key.get():
      FlakeUpdateSingleton(key=singleton_key).put()
    return singleton_key

  @ndb.transactional
  def _increment_update_counter(self):
    FlakeUpdate(parent=self._get_flake_update_singleton_key()).put()

  @ndb.non_transactional
  def _time_difference(self, flaky_run):
    return (flaky_run.success_run.get().time_finished -
            flaky_run.failure_run_time_finished).total_seconds()

  @ndb.non_transactional
  def _is_same_day(self, flaky_run):
    time_since_finishing = (
        datetime.datetime.utcnow() - flaky_run.failure_run_time_finished)
    return time_since_finishing <= datetime.timedelta(days=1)

  @ndb.non_transactional
  def _get_new_flakes(self, flake):
    num_runs = len(flake.occurrences) - flake.num_reported_flaky_runs
    flaky_runs = ndb.get_multi(flake.occurrences[-num_runs:])
    flaky_runs = [run for run in flaky_runs if run is not None]
    return [
      flaky_run for flaky_run in flaky_runs
      if self._is_same_day(flaky_run) and
         self._time_difference(flaky_run) <= MAX_TIME_DIFFERENCE_SECONDS]

  @staticmethod
  @ndb.non_transactional
  def _get_first_flake_occurrence_time(flake):
    assert flake.occurrences, 'Flake entity has no occurrences'
    flaky_runs = ndb.get_multi(flake.occurrences)
    flaky_runs = [run for run in flaky_runs if run is not None]
    rev_occ = sorted(flaky_runs, key=lambda run: run.failure_run_time_finished,
                     reverse=True)
    last_occ_time = rev_occ[0].failure_run_time_finished
    for occ in rev_occ[1:]:
      occ_time = occ.failure_run_time_finished
      if last_occ_time - occ_time > MAX_GAP_FOR_FLAKINESS_PERIOD:
        break
      last_occ_time = occ_time
    return last_occ_time

  @ndb.transactional
  def _recreate_issue_for_flake(self, flake):
    """Updates a flake to re-create an issue and creates a respective task."""
    flake.old_issue_id = flake.issue_id
    flake.issue_id = 0
    taskqueue.add(url='/issues/process/%s' % flake.key.urlsafe(),
                  queue_name='issue-updates', transactional=True)

  @staticmethod
  @ndb.non_transactional
  def _update_new_occurrences_with_issue_id(name, new_flaky_runs, issue_id):
    # TODO(sergiyb): Find a way to do this asynchronously to avoid block
    # transaction-bound method calling this. Possible solutions are to use
    # put_multi_sync (need to find a way to test this) or to use deferred
    # execution.
    for fr in new_flaky_runs:
      for occ in fr.flakes:
        if occ.failure == name:
          occ.issue_id = issue_id
    ndb.put_multi(new_flaky_runs)

  @ndb.transactional
  def _update_issue(self, api, flake, new_flakes, now):
    """Updates an issue on the issue tracker."""
    flake_issue = api.getIssue(flake.issue_id)

    # Handle cases when an issue has been closed. We need to do this in a loop
    # because we might move onto another issue.
    seen_issues = set()
    while not flake_issue.open:
      if flake_issue.status == 'Duplicate':
        # If the issue was marked as duplicate, we update the issue ID stored in
        # datastore to the one it was merged into and continue working with the
        # new issue.
        seen_issues.add(flake_issue.id)
        if flake_issue.merged_into not in seen_issues:
          flake.issue_id = flake_issue.merged_into
          flake_issue = api.getIssue(flake.issue_id)
        else:
          logging.info('Detected issue duplication loop: %s. Re-creating an '
                       'issue for the flake %s.', seen_issues, flake.name)
          self._recreate_issue_for_flake(flake)
          return
      else:  # Fixed, WontFix, Verified, Archived, custom status
        # If the issue was closed, we do not update it. This allows changes made
        # to reduce flakiness to propagate and take effect. If after
        # DAYS_TO_REOPEN_ISSUE days we still detect flakiness, we will create a
        # new issue.
        recent_cutoff = now - datetime.timedelta(days=DAYS_TO_REOPEN_ISSUE)
        if flake_issue.updated < recent_cutoff:
          self._recreate_issue_for_flake(flake)
        return

    # Make sure issue is in the appropriate bug queue as flakiness is ongoing.
    queue_name, expected_label = get_queue_details(flake.name)
    if expected_label not in flake_issue.labels:
      flake_issue.labels.append(expected_label)

    new_flaky_runs_msg = FLAKY_RUNS_TEMPLATE % {
        'name': flake.name,
        'new_flakes_count': len(new_flakes),
        'flakes_url': FLAKES_URL_TEMPLATE % flake.key.urlsafe(),
        'queue_name': queue_name}
    api.update(flake_issue, comment=new_flaky_runs_msg)
    self.issue_updates.increment_by(1, {'operation': 'update'})
    logging.info('Updated issue %d for flake %s with %d flake runs',
                 flake.issue_id, flake.name, len(new_flakes))
    self._update_new_occurrences_with_issue_id(
        flake.name, new_flakes, flake_issue.id)
    flake.num_reported_flaky_runs = len(flake.occurrences)
    flake.issue_last_updated = now

  @ndb.transactional
  def _create_issue(self, api, flake, new_flakes, now):
    _, qlabel = get_queue_details(flake.name)
    labels = ['Type-Bug', 'Pri-1', 'Via-TryFlakes', qlabel]
    if is_trooper_flake(flake.name):
      other_queue_msg = TROOPER_QUEUE_MSG
    else:
      other_queue_msg = SHERIFF_QUEUE_MSG

    summary = SUMMARY_TEMPLATE % {'name': flake.name}
    description = DESCRIPTION_TEMPLATE % {
        'summary': summary,
        'flakes_url': FLAKES_URL_TEMPLATE % flake.key.urlsafe(),
        'flakes_count': len(new_flakes),
        'other_queue_msg': other_queue_msg}
    if flake.old_issue_id:
      description = REOPENED_DESCRIPTION_TEMPLATE % {
          'description': description, 'old_issue': flake.old_issue_id}

    new_issue = issue.Issue({'summary': summary,
                             'description': description,
                             'status': 'Untriaged',
                             'labels': labels,
                             'components': ['Tests>Flaky']})
    flake_issue = api.create(new_issue)
    flake.issue_id = flake_issue.id
    self._update_new_occurrences_with_issue_id(
        flake.name, new_flakes, flake_issue.id)
    flake.num_reported_flaky_runs = len(flake.occurrences)
    flake.issue_last_updated = now
    self.issue_updates.increment_by(1, {'operation': 'create'})
    logging.info('Created a new issue %d for flake %s', flake.issue_id,
                 flake.name)

    delay = (now - self._get_first_flake_occurrence_time(flake)).total_seconds()
    self.reporting_delay.set(delay)
    logging.info('Reported reporting_delay %d for flake %s', delay, flake.name)

  @ndb.transactional(xg=True)  # pylint: disable=E1120
  def post(self, urlsafe_key):
    api = issue_tracker_api.IssueTrackerAPI(
        'chromium', use_monorail=USE_MONORAIL)

    # Check if we should stop processing this issue because we've posted too
    # many updates to issue tracker today already.
    day_ago = datetime.datetime.utcnow() - datetime.timedelta(days=1)
    num_updates_last_day = FlakeUpdate.query(
        FlakeUpdate.time_updated > day_ago,
        ancestor=self._get_flake_update_singleton_key()).count()
    if num_updates_last_day >= MAX_UPDATED_ISSUES_PER_DAY:
      return

    now = datetime.datetime.utcnow()
    flake = ndb.Key(urlsafe=urlsafe_key).get()
    # Only update/file issues if there are new flaky runs.
    if flake.num_reported_flaky_runs == len(flake.occurrences):
      return

    # Retrieve flaky runs outside of the transaction, because we are not
    # planning to modify them and because there could be more of them than the
    # number of groups supported by cross-group transactions on AppEngine.
    new_flakes = self._get_new_flakes(flake)

    if len(new_flakes) < MIN_REQUIRED_FLAKY_RUNS:
      return

    if flake.issue_id > 0:
      # Update issues at most once a day.
      if flake.issue_last_updated > now - datetime.timedelta(days=1):
        return

      self._update_issue(api, flake, new_flakes, now)
      self._increment_update_counter()
    else:
      self._create_issue(api, flake, new_flakes, now)
      # Don't update the issue just yet, this may fail, and we need the
      # transaction to succeed in order to avoid filing duplicate bugs.
      self._increment_update_counter()

    # Note that if transaction fails for some reason at this point, we may post
    # updates or create issues multiple times. On the other hand, this should be
    # extremely rare because we set the number of concurrently running tasks to
    # 1, therefore there should be no contention for updating this issue's
    # entity.
    flake.put()


class UpdateIfStaleIssue(webapp2.RequestHandler):
  def _remove_issue_from_flakes(self, issue_id):
    for flake in Flake.query(Flake.issue_id == issue_id):
      logging.info('Removing issue_id %s from %s', issue_id, flake.key)
      flake.old_issue_id = issue_id
      flake.issue_id = 0
      flake.put()

  def post(self, issue_id):
    """Check if an issue is stale and report it back to appropriate queue.

    Check if the issue is stale, i.e. has not been updated by anyone else than
    this app in the last DAYS_TILL_STALE days, and if this is the case, then
    move it back to the appropriate queue. Also if the issue is stale for 7
    days, report it to stale-flakes-reports@google.com to investigate why it is
    not being processed despite being in the appropriate queue.
    """
    issue_id = int(issue_id)
    api = issue_tracker_api.IssueTrackerAPI(
        'chromium', use_monorail=USE_MONORAIL)
    flake_issue = api.getIssue(issue_id)
    now = datetime.datetime.utcnow()

    if not flake_issue.open:
      # Remove issue_id from all flakes unless it has recently been updated. We
      # should not remove issue_id too soon, otherwise issues will get reopened
      # before changes made will propogate and reduce flakiness.
      recent_cutoff = now - datetime.timedelta(days=DAYS_TO_REOPEN_ISSUE)
      if flake_issue.updated < recent_cutoff:
        self._remove_issue_from_flakes(issue_id)
      return

    # Find the last update, which defaults to when issue was created if no third
    # party updates were posted.
    comments = api.getComments(issue_id)
    last_third_party_update = flake_issue.created
    for comment in sorted(comments, key=lambda c: c.created, reverse=True):
      if comment.author != app_identity.get_service_account_name():
        last_third_party_update = comment.created
        break

    # Parse the flake name from the first comment (which we post ourselves).
    original_summary = comments[0].comment.splitlines()[0]
    flake_name = original_summary[len('"'):-len('" is flaky')]
    _, expected_label = get_queue_details(flake_name)

    # Report to stale-flakes-reports@ if the issue has been in appropriate queue
    # without any updates for 7 days.
    week_ago = now - datetime.timedelta(days=7)
    if (last_third_party_update < week_ago and
        expected_label in flake_issue.labels and
        STALE_FLAKES_ML not in flake_issue.cc):
      flake_issue.cc.append(STALE_FLAKES_ML)
      logging.info('Reporting issue %s to %s', flake_issue.id, STALE_FLAKES_ML)
      api.update(flake_issue, comment=VERY_STALE_FLAKES_MESSAGE)


class CreateFlakyRun(webapp2.RequestHandler):
  flaky_runs = ts_mon.CounterMetric(
      'flakiness_pipeline/flake_occurrences_recorded',
      description='Recorded flake occurrences.')

  # We execute below method in an indepedent transaction since otherwise we
  # would exceed the maximum number of entities allowed within a single
  # transaction.
  @staticmethod
  # pylint: disable=E1120
  @ndb.transactional(xg=True, propagation=ndb.TransactionOptions.INDEPENDENT)
  def add_failure_to_flake(name, flaky_run_key, failure_time):
    flake = Flake.get_by_id(name)
    if not flake:
      flake = Flake(name=name, id=name, last_time_seen=datetime.datetime.min)
      flake.put()

    flake.occurrences.append(flaky_run_key)
    util.add_occurrence_time_to_flake(flake, failure_time)
    flake.put()

  @classmethod
  def _flatten_tests(cls, tests, delimiter='/', prefix=None):
    """Flattens hierarchical GTest JSON test structure.

    The hiearchical GTest JSON test structure is described in
    https://www.chromium.org/developers/the-json-test-results-format (see
    top-level 'tests' key).

    We only return 3 test types:
     - passed, i.e. expected is "PASS" and last actual run is "PASS"
     - failed, i.e. expected is "PASS" and last actual run is "FAIL", "TIMEOUT"
       or "CRASH"
     - skipped, i.e. expected and actual are both "SKIP"

    We do not classify or return any other tests, in particular:
     - known flaky, i.e. expected to have varying results, e.g. "PASS FAIL".
     - known failing, i.e. expected is "FAIL", "TIMEOUT" or "CRASH".
     - unexpected flakiness, i.e. failures than hapeneed before last PASS.

    Args:
      prefix: Prefix to be added before test names if this is a parent node.
      delimiter: Delimiter to use for concatenating parts of test name.
      tests: Any non-leaf node of the hierarchical GTest JSON test structure.

    Returns:
      A tuple (passed, failed, skpped), where each is a list of test names.
    """
    passed = []
    failed = []
    skipped = []
    for node_name, node in tests.iteritems():
      # Compute current node name, which would be a test name for leaf nodes or
      # new prefix for parent nodes.
      test_name = prefix + delimiter + node_name if prefix else node_name

      # Check if it is a leaf node first.
      if 'actual' in node and 'expected' in node:
        if node['expected'] == 'PASS':
          actual_results = node['actual'].split(' ')
          if actual_results[-1] == 'PASS':
            passed.append(test_name)
          elif actual_results[-1] in ('FAIL', 'TIMEOUT', 'CRASH'):
            failed.append(test_name)
        elif node['expected'] == 'SKIP' and node['actual'] == 'SKIP':
          skipped.append(test_name)
      else:
        node_passed, node_failed, node_skipped = cls._flatten_tests(
            node, delimiter, test_name)
        passed.extend(node_passed)
        failed.extend(node_failed)
        skipped.extend(node_skipped)

    return passed, failed, skipped

  # see examples:
  # compile http://build.chromium.org/p/tryserver.chromium.mac/json/builders/
  #         mac_chromium_compile_dbg/builds/11167?as_text=1
  # gtest http://build.chromium.org/p/tryserver.chromium.win/json/builders/
  #       win_chromium_x64_rel_swarming/builds/4357?as_text=1
  # TODO(jam): get specific problem with compile so we can use that as name
  @classmethod
  def get_flakes(cls, mastername, buildername, buildnumber, step):
    # If test results were invalid, report whole step as flaky.
    steptext = ' '.join(step['text'])
    stepname = normalize_test_type(step['name'])
    if 'TEST RESULTS WERE INVALID' in steptext:
      return [stepname]

    url = TEST_RESULTS_URL_TEMPLATE % {
      'mastername': urllib2.quote(mastername),
      'buildername': urllib2.quote(buildername),
      'buildnumber': urllib2.quote(str(buildnumber)),
      'stepname': urllib2.quote(stepname),
    }

    try:
      result = urlfetch.fetch(url)

      if result.status_code >= 200 and result.status_code < 400:
        json_result = json.loads(result.content)

        _, failed, _ = cls._flatten_tests(
            json_result.get('tests', {}),
            json_result.get('path_delimiter', '/'))
        return failed

      if result.status_code == 404:
        # This is quite a common case (only some failing steps are actually
        # running tests and reporting results to flakiness dashboard).
        logging.info('Failed to retrieve JSON from %s', url)
      else:
        logging.exception('Failed to retrieve JSON from %s', url)
    except Exception:
      logging.exception('Failed to retrieve or parse JSON from %s', url)

    return [stepname]

  @ndb.transactional(xg=True)  # pylint: disable=E1120
  def post(self):
    if (not self.request.get('failure_run_key') or
        not self.request.get('success_run_key')):
      self.response.set_status(400, 'Invalid request parameters')
      return

    failure_run = ndb.Key(urlsafe=self.request.get('failure_run_key')).get()
    success_run = ndb.Key(urlsafe=self.request.get('success_run_key')).get()

    flaky_run = FlakyRun(
        failure_run=failure_run.key,
        failure_run_time_started=failure_run.time_started,
        failure_run_time_finished=failure_run.time_finished,
        success_run=success_run.key)

    success_time = success_run.time_finished
    failure_time = failure_run.time_finished
    patchset_builder_runs = failure_run.key.parent().get()

    # TODO(sergiyb): The parsing logic below is very fragile and will break with
    # any changes to step names and step text. We should move away from parsing
    # buildbot to tools like flakiness dashboard (test-results.appspot.com),
    # which uses a standartized JSON format.
    url = ('http://build.chromium.org/p/' + patchset_builder_runs.master +
           '/json/builders/' + patchset_builder_runs.builder +'/builds/' +
           str(failure_run.buildnumber))
    urlfetch.set_default_fetch_deadline(60)
    logging.info('get_flaky_run_reason ' + url)
    result = urlfetch.fetch(url).content
    try:
      json_result = json.loads(result)
    except ValueError:
      logging.exception('couldnt decode json for %s', url)
      return
    steps = json_result['steps']

    failed_steps = []
    passed_steps = []
    for step in steps:
      result = step['results'][0]
      if build_result.isResultSuccess(result):
        passed_steps.append(step)
        continue
      if not build_result.isResultFailure(result):
        continue
      step_name = step['name']
      step_text = ' '.join(step['text'])
      # The following step failures are ignored:
      #  - steps: always red when any other step is red (not a failure)
      #  - [swarming] ...: summary step would also be red (do not double count)
      #  - presubmit: typically red due to missing OWNERs LGTM, not a flake
      #  - recipe failure reason: always red when build fails (not a failure)
      #  - Patch failure: if success run was before failure run, it is
      #    likely a legitimate failure. For example it often happens that
      #    developers use CQ dry run and then wait for a review. Once getting
      #    LGTM they check CQ checkbox, but the patch does not cleanly apply
      #    anymore.
      #  - bot_update PATCH FAILED: Corresponds to 'Patch failure' step.
      #  - test results: always red when another step is red (not a failure)
      #  - Uncaught Exception: summary step referring to an exception in another
      #    step (e.g. bot_update)
      #  - ... (retry summary): this is an artificial step to fail the build due
      #    to another step that has failed earlier (do not double count).
      if (step_name == 'steps' or step_name.startswith('[swarming]') or
          step_name == 'presubmit' or step_name == 'recipe failure reason' or
          (step_name == 'Patch failure' and success_time < failure_time) or
          (step_name == 'bot_update' and 'PATCH FAILED' in step_text) or
          step_name == 'test results' or step_name == 'Uncaught Exception' or
          step_name.endswith(' (retry summary)')):
        continue
      failed_steps.append(step)

    steps_to_ignore = []
    for step in failed_steps:
      step_name = step['name']
      if ' (with patch)' in step_name:
        # Android instrumentation tests add a prefix before the step name, which
        # doesn't appear on the summary step (without suffixes). To make sure we
        # correctly ignore duplicate failures, we remove the prefix.
        step_name = step_name.replace('Instrumentation test ', '')

        # If a step fails without the patch, then the tree is busted. Don't
        # count as flake.
        step_name_with_no_modifier = step_name.replace(' (with patch)', '')
        step_name_without_patch = (
            '%s (without patch)' % step_name_with_no_modifier)
        for other_step in failed_steps:
          if other_step['name'] == step_name_without_patch:
            steps_to_ignore.append(step['name'])
            steps_to_ignore.append(other_step['name'])

    flakes_to_update = []
    for step in failed_steps:
      step_name = step['name']
      if step_name in steps_to_ignore:
        continue
      flakes = self.get_flakes(
          patchset_builder_runs.master, patchset_builder_runs.builder,
          failure_run.buildnumber, step)
      for flake in flakes:
        flake_occurrence = FlakeOccurrence(name=step_name, failure=flake)
        flaky_run.flakes.append(flake_occurrence)
        flakes_to_update.append(flake)

    flaky_run_key = flaky_run.put()
    for flake in flakes_to_update:
      self.add_failure_to_flake(flake, flaky_run_key, failure_time)
    self.flaky_runs.increment_by(1)
