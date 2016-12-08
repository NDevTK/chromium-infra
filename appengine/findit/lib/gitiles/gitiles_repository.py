# Copyright 2014 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

import base64
from datetime import datetime
from datetime import timedelta
import json
import re

from lib.gitiles import commit_util
from lib.gitiles import diff
from lib.gitiles.blame import Blame
from lib.gitiles.blame import Region
from lib.gitiles.change_log import ChangeLog
from lib.gitiles.change_log import FileChangeInfo
from lib.gitiles.git_repository import GitRepository
from lib.time_util import TimeZoneInfo

COMMIT_POSITION_PATTERN = re.compile(
    '^Cr-Commit-Position: refs/heads/master@{#(\d+)}$', re.IGNORECASE)
CODE_REVIEW_URL_PATTERN = re.compile(
    '^(?:Review URL|Review-Url): (.*\d+).*$', re.IGNORECASE)
REVERTED_REVISION_PATTERN = re.compile(
    '^> Committed: https://.+/([0-9a-fA-F]{40})$', re.IGNORECASE)
TIMEZONE_PATTERN = re.compile('[-+]\d{4}$')
CACHE_EXPIRE_TIME_SECONDS = 24 * 60 * 60


class GitilesRepository(GitRepository):
  """Use Gitiles to access a repository on https://chromium.googlesource.com."""

  def __init__(self, http_client, repo_url=None):
    super(GitilesRepository, self).__init__()
    if repo_url and repo_url.endswith('/'):
      self._repo_url = repo_url[:-1]
    else:
      self._repo_url = repo_url

    self._http_client = http_client

  @property
  def repo_url(self):
    return self._repo_url

  @repo_url.setter
  def repo_url(self, repo_url):
    self._repo_url = repo_url

  @property
  def http_client(self):
    return self._http_client

  @property
  def identifier(self):
    return self.repo_url

  def _SendRequestForJsonResponse(self, url, params=None):
    if params is None:  # pragma: no cover
      params = {}
    params['format'] = 'json'

    # Gerrit prepends )]}' to json-formatted response.
    prefix = ')]}\'\n'

    status_code, content = self.http_client.Get(url, params)
    if status_code != 200:
      return None
    elif not content or not content.startswith(prefix):
      raise Exception('Response does not begin with %s' % prefix)

    return json.loads(content[len(prefix):])

  def _SendRequestForTextResponse(self, url):
    status_code, content = self.http_client.Get(url, {'format': 'text'})
    if status_code != 200:
      return None
    return base64.b64decode(content)

  def _GetDateTimeFromString(self, datetime_string,
                             date_format='%a %b %d %H:%M:%S %Y'):
    if TIMEZONE_PATTERN.findall(datetime_string):
      # Need to handle timezone conversion.
      naive_datetime_str, _, offset_str = datetime_string.rpartition(' ')
      naive_datetime = datetime.strptime(naive_datetime_str, date_format)
      return TimeZoneInfo(offset_str).LocalToUTC(naive_datetime)

    return datetime.strptime(datetime_string, date_format)

  def _DownloadChangeLogData(self, revision):
    url = '%s/+/%s' % (self.repo_url, revision)
    return url, self._SendRequestForJsonResponse(url)

  def _ParseChangeLogFromLogData(self, data):
    commit_position, code_review_url = (
        commit_util.ExtractCommitPositionAndCodeReviewUrl(data['message']))

    touched_files = []
    for file_diff in data['tree_diff']:
      change_type = file_diff['type'].lower()
      if not diff.IsKnownChangeType(change_type):
        raise Exception('Unknown change type "%s"' % change_type)
      touched_files.append(
          FileChangeInfo(
              change_type, file_diff['old_path'], file_diff['new_path']))

    author_time = self._GetDateTimeFromString(data['author']['time'])
    committer_time = self._GetDateTimeFromString(data['committer']['time'])
    reverted_revision = commit_util.GetRevertedRevision(data['message'])
    url = '%s/+/%s' % (self.repo_url, data['commit'])

    return ChangeLog(
        data['author']['name'],
        commit_util.NormalizeEmail(data['author']['email']),
        author_time,
        data['committer']['name'],
        commit_util.NormalizeEmail(data['committer']['email']),
        committer_time, data['commit'], commit_position,
        data['message'], touched_files, url, code_review_url,
        reverted_revision)

  def GetChangeLog(self, revision):
    """Returns the change log of the given revision."""
    _, data = self._DownloadChangeLogData(revision)
    if not data:
      return None

    return self._ParseChangeLogFromLogData(data)

  def GetCommitsBetweenRevisions(self, start_revision, end_revision, n=1000):
    """Gets a list of commit hashes between start_revision and end_revision.

    Args:
      start_revision: The oldest revision in the range.
      end_revision: The latest revision in the range.
      n: The maximum number of revisions to request at a time.

    Returns:
      A list of commit hashes made since start_revision through and including
      end_revision in order from most-recent to least-recent. This includes
      end_revision, but not start_revision.
    """
    params = {'n': n}
    next_end_revision = end_revision
    commits = []

    while next_end_revision:
      url = '%s/+log/%s..%s' % (
          self.repo_url, start_revision, next_end_revision)
      data = self._SendRequestForJsonResponse(url, params)

      if not data:
        break

      for log in data.get('log', []):
        commit = log.get('commit')
        if commit:
          commits.append(commit)

      next_end_revision = data.get('next')

    return commits

  def GetChangeDiff(self, revision):
    """Returns the raw diff of the given revision."""
    url = '%s/+/%s%%5E%%21/' % (self.repo_url, revision)
    return self._SendRequestForTextResponse(url)

  def GetBlame(self, path, revision):
    """Returns blame of the file at ``path`` of the given revision."""
    url = '%s/+blame/%s/%s' % (self.repo_url, revision, path)

    data = self._SendRequestForJsonResponse(url)
    if not data:
      return None

    blame = Blame(revision, path)
    for region in data['regions']:
      author_time = self._GetDateTimeFromString(
          region['author']['time'], '%Y-%m-%d %H:%M:%S')

      blame.AddRegion(
          Region(region['start'], region['count'], region['commit'],
                 region['author']['name'],
                 commit_util.NormalizeEmail(region['author']['email']),
                 author_time))

    return blame

  def GetSource(self, path, revision):
    """Returns source code of the file at ``path`` of the given revision."""
    url = '%s/+/%s/%s' % (self.repo_url, revision, path)
    return self._SendRequestForTextResponse(url)

  def GetChangeLogs(self, start_revision, end_revision, n=1000):
    """Gets a list of ChangeLogs in revision range by batch.

    Args:
      start_revision (str): The oldest revision in the range.
      end_revision (str): The latest revision in the range.
      n (int): The maximum number of revisions to request at a time (default
        to 1000).

    Returns:
      A list of changelogs in (start_revision, end_revision].
    """
    next_end_revision = end_revision
    changelogs = []

    while next_end_revision:
      url = '%s/+log/%s..%s' % (self.repo_url,
                                start_revision, next_end_revision)
      data = self._SendRequestForJsonResponse(url, params={'n': str(n),
                                                           'name-status': '1'})

      for log in data['log']:
        changelogs.append(self._ParseChangeLogFromLogData(log))

      if 'next' in data:
        next_end_revision = data['next']
      else:
        next_end_revision = None

    return changelogs
