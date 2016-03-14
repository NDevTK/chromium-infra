# Copyright 2014 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

import gae_ts_mon
import gae_event_mon
import os
import sys
import webapp2

THIRD_PARTY_DIR = os.path.join(os.path.dirname(__file__), 'third_party')
sys.path.insert(0, os.path.join(THIRD_PARTY_DIR, 'dateutil'))
sys.path.insert(0, os.path.join(THIRD_PARTY_DIR, 'gae-pytz'))
sys.path.insert(0, os.path.join(THIRD_PARTY_DIR, 'google-api-python-client'))
sys.path.insert(0, os.path.join(THIRD_PARTY_DIR, 'httplib2', 'python2'))
sys.path.insert(0, os.path.join(THIRD_PARTY_DIR, 'oauth2client'))
sys.path.insert(0, os.path.join(THIRD_PARTY_DIR, 'six'))
sys.path.insert(0, os.path.join(THIRD_PARTY_DIR, 'test_results'))
sys.path.insert(0, os.path.join(THIRD_PARTY_DIR, 'time_functions'))
sys.path.insert(0, os.path.join(THIRD_PARTY_DIR, 'uritemplate'))

from handlers.cron_dispatch import CronDispatch
from handlers.index import Index
from handlers.post_comment import PostComment
from handlers.all_flake_occurrences import AllFlakeOccurrences
from handlers.search import Search
from handlers import flake_issues

handlers = [
  (r'/', Index),
  (r'/post_comment', PostComment),
  (r'/all_flake_occurrences', AllFlakeOccurrences),
  (r'/search', Search),
  (r'/cron/(.*)', CronDispatch),
  (r'/issues/process/(.*)', flake_issues.ProcessIssue),
  (r'/issues/update-if-stale/(.*)', flake_issues.UpdateIfStaleIssue),
  (r'/issues/create_flaky_run', flake_issues.CreateFlakyRun),
]

app = webapp2.WSGIApplication(handlers, debug=True)
gae_ts_mon.initialize(app)
gae_event_mon.initialize('flakiness_pipeline')
