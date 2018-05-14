# Copyright 2014 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

import datetime

from google.appengine.ext import ndb

from model.build_run import BuildRun

# A particular occurrence of a single flake occuring.
class FlakeOccurrence(ndb.Model):
  # step name, i.e. browser_tests
  name = ndb.StringProperty(required=True)
  # failure, i.e. FooTest.Bar
  failure = ndb.StringProperty(required=True)
  # issue in which this occurrence was reported
  issue_id = ndb.IntegerProperty(default=0)


# Represents a patchset with a successful and failed try run for a particular
# builder. The flaky failed run can have multiple different flakes that cause
# it to turn red, each represented by a FlakeOccurrence.
class FlakyRun(ndb.Model):
  failure_run = ndb.KeyProperty(BuildRun, required=True)
  # A copy of failure_run.time_started to reduce lookups.
  failure_run_time_started = ndb.DateTimeProperty(default=datetime.datetime.max)
  # A copy of failure_run.time_finished to reduce lookups.
  failure_run_time_finished = ndb.DateTimeProperty(required=True)
  success_run = ndb.KeyProperty(BuildRun, required=True)
  flakes = ndb.StructuredProperty(FlakeOccurrence, repeated=True)
  comment = ndb.StringProperty()


# Represents a step that flakes. The name could be a test_suite:test name (i.e.
# unit_tests:FooTest), a ninja step in case of compile flake, etc... This entity
# groups together all the occurrences of this flake, with each particular
# instance of a flake being represented by a FlakyRun.
class Flake(ndb.Model):
  name = ndb.StringProperty(required=True)
  is_step = ndb.BooleanProperty(default=False)
  occurrences = ndb.KeyProperty(FlakyRun, repeated=True)
  comment = ndb.StringProperty(default='')

  # Used so we can quickly query and sort by number of occurrences per time
  # range.
  count_hour = ndb.IntegerProperty(default=0)
  count_day = ndb.IntegerProperty(default=0)
  count_week = ndb.IntegerProperty(default=0)
  count_month = ndb.IntegerProperty(default=0)
  count_all = ndb.IntegerProperty(default=0)
  last_time_seen = ndb.DateTimeProperty()

  # This is needed to allow the query in update_flake_date_flags to be fast.
  last_hour = ndb.BooleanProperty(default=False)
  last_day = ndb.BooleanProperty(default=False)
  last_week = ndb.BooleanProperty(default=False)
  last_month = ndb.BooleanProperty(default=False)

  # An issue filed on issue tracker to track flakiness of this step/test.
  issue_id = ndb.IntegerProperty(default=0)
  issue_last_updated = ndb.DateTimeProperty(default=datetime.datetime.min)

  # Stores previous issue ID when the issue need to be re-created.
  old_issue_id = ndb.IntegerProperty(default=0)

  # Number of occurences that were already reported to the issue tracker and
  # Findit respectively. They don't share the same value because Chromium Try
  # Flake can only report flaky runs to at most 10 issues per day, however,
  # there is no upper limit for reporting to Findit.
  # NOTE: num_reported_flaky_runs should be renamed to
  # num_reported_flaky_runs_to_issue, but it didn't happen due to backward
  # compatbility.
  num_reported_flaky_runs = ndb.IntegerProperty(default=0)
  num_reported_flaky_runs_to_findit = ndb.IntegerProperty(default=0)


# The following two entities are used to track updates posted to the issue
# tracker and prevent too many updates filed. The FlakeUpdateSingleton entity is
# a singleton and all FlakyUpdate entities should be child entities of this
# singleton entity, which allows us to query all of them within a single
# transaction.
class FlakeUpdateSingleton(ndb.Model):
  pass


class FlakeUpdate(ndb.Model):
  time_updated = ndb.DateTimeProperty(auto_now_add=True)


class FlakeType(ndb.Model):
  project = ndb.StringProperty(required=True)
  step_name = ndb.StringProperty(required=True)
  test_name = ndb.StringProperty()
  config = ndb.StringProperty()
  last_updated = ndb.DateTimeProperty(required=True)


class Issue(ndb.Model):
  project = ndb.StringProperty(required=True)
  issue_id = ndb.IntegerProperty(required=True)
  flake_type_keys = ndb.KeyProperty(kind=FlakeType, repeated=True)


class CacheAncestor(ndb.Model):
  pass


class Cache(ndb.Model):
  idx = ndb.IntegerProperty(required=True)
  data = ndb.BlobProperty()
