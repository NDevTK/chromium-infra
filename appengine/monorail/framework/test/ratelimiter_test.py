# Copyright 2016 The Chromium Authors. All rights reserved.
# Use of this source code is govered by a BSD-style
# license that can be found in the LICENSE file or at
# https://developers.google.com/open-source/licenses/bsd

"""Unit tests for RateLimiter.
"""
import unittest

from google.appengine.api import memcache
from google.appengine.ext import testbed

import mox
import os
import settings

from framework import ratelimiter
from services import service_manager
from testing import fake
from testing import testing_helpers


class RateLimiterTest(unittest.TestCase):
  def setUp(self):
    settings.ratelimiting_enabled = True
    self.testbed = testbed.Testbed()
    self.testbed.activate()
    self.testbed.init_memcache_stub()
    self.testbed.init_user_stub()

    self.mox = mox.Mox()
    self.services = service_manager.Services(
      config=fake.ConfigService(),
      issue=fake.IssueService(),
      user=fake.UserService(),
      project=fake.ProjectService(),
    )
    self.project = self.services.project.TestAddProject('proj', project_id=987)

    self.ratelimiter = ratelimiter.RateLimiter()
    ratelimiter.COUNTRY_LIMITS = {}
    os.environ['USER_EMAIL'] = ''
    settings.ratelimiting_enabled = True
    settings.ratelimiting_cost_enabled = True
    ratelimiter.DEFAULT_LIMIT = 10

  def tearDown(self):
    self.testbed.deactivate()
    self.mox.UnsetStubs()
    self.mox.ResetAll()
    # settings.ratelimiting_enabled = True

  def testCheckStart_pass(self):
    request, _ = testing_helpers.GetRequestObjects(
      project=self.project)
    request.headers['X-AppEngine-Country'] = 'US'
    request.remote_addr = '192.168.1.0'
    self.ratelimiter.CheckStart(request)
    # Should not throw an exception.

  def testCheckStart_fail(self):
    request, _ = testing_helpers.GetRequestObjects(
      project=self.project)
    request.headers['X-AppEngine-Country'] = 'US'
    request.remote_addr = '192.168.1.0'
    now = 0.0
    cachekeysets, _, _, _ = ratelimiter._CacheKeys(request, now)
    values = [{key: ratelimiter.DEFAULT_LIMIT for key in cachekeys} for
              cachekeys in cachekeysets]
    for value in values:
      memcache.add_multi(value)
    with self.assertRaises(ratelimiter.RateLimitExceeded):
      self.ratelimiter.CheckStart(request, now)

  def testCheckStart_expiredEntries(self):
    request, _ = testing_helpers.GetRequestObjects(
      project=self.project)
    request.headers['X-AppEngine-Country'] = 'US'
    request.remote_addr = '192.168.1.0'
    now = 0.0
    cachekeysets, _, _, _ = ratelimiter._CacheKeys(request, now)
    values = [{key: ratelimiter.DEFAULT_LIMIT for key in cachekeys} for
              cachekeys in cachekeysets]
    for value in values:
      memcache.add_multi(value)

    now = now + 2 * ratelimiter.EXPIRE_AFTER_SECS
    self.ratelimiter.CheckStart(request, now)
    # Should not throw an exception.

  def testCheckStart_repeatedCalls(self):
    request, _ = testing_helpers.GetRequestObjects(
      project=self.project)
    request.headers['X-AppEngine-Country'] = 'US'
    request.remote_addr = '192.168.1.0'
    now = 0.0

    # Call CheckStart once every minute.  Should be ok.
    for _ in range(ratelimiter.N_MINUTES):
      self.ratelimiter.CheckStart(request, now)
      now = now + 120.0

    # Call CheckStart more than DEFAULT_LIMIT times in the same minute.
    with self.assertRaises(ratelimiter.RateLimitExceeded):
      for _ in range(ratelimiter.DEFAULT_LIMIT + 2):
        now = now + 0.001
        self.ratelimiter.CheckStart(request, now)

  def testCheckStart_differentIPs(self):
    now = 0.0

    ratelimiter.COUNTRY_LIMITS = {}
    # Exceed DEFAULT_LIMIT calls, but vary remote_addr so different
    # remote addresses aren't ratelimited together.
    for m in range(ratelimiter.DEFAULT_LIMIT * 2):
      request, _ = testing_helpers.GetRequestObjects(
        project=self.project)
      request.headers['X-AppEngine-Country'] = 'US'
      request.remote_addr = '192.168.1.%d' % (m % 16)
      ratelimiter._CacheKeys(request, now)
      self.ratelimiter.CheckStart(request, now)
      now = now + 0.001

    # Exceed the limit, but only for one IP address. The
    # others should be fine.
    with self.assertRaises(ratelimiter.RateLimitExceeded):
      for m in range(ratelimiter.DEFAULT_LIMIT):
        request, _ = testing_helpers.GetRequestObjects(
          project=self.project)
        request.headers['X-AppEngine-Country'] = 'US'
        request.remote_addr = '192.168.1.0'
        ratelimiter._CacheKeys(request, now)
        self.ratelimiter.CheckStart(request, now)
        now = now + 0.001

    # Now proceed to make requests for all of the other IP
    # addresses besides .0.
    for m in range(ratelimiter.DEFAULT_LIMIT * 2):
      request, _ = testing_helpers.GetRequestObjects(
        project=self.project)
      request.headers['X-AppEngine-Country'] = 'US'
      # Skip .0 since it's already exceeded the limit.
      request.remote_addr = '192.168.1.%d' % (m + 1)
      ratelimiter._CacheKeys(request, now)
      self.ratelimiter.CheckStart(request, now)
      now = now + 0.001

  def testCheckStart_sameIPDifferentUserIDs(self):
    # Behind a NAT, e.g.
    now = 0.0

    # Exceed DEFAULT_LIMIT calls, but vary user_id so different
    # users behind the same IP aren't ratelimited together.
    for m in range(ratelimiter.DEFAULT_LIMIT * 2):
      request, _ = testing_helpers.GetRequestObjects(
        project=self.project)
      request.remote_addr = '192.168.1.0'
      os.environ['USER_EMAIL'] = '%s@example.com' % m
      request.headers['X-AppEngine-Country'] = 'US'
      ratelimiter._CacheKeys(request, now)
      self.ratelimiter.CheckStart(request, now)
      now = now + 0.001

    # Exceed the limit, but only for one userID+IP address. The
    # others should be fine.
    with self.assertRaises(ratelimiter.RateLimitExceeded):
      for m in range(ratelimiter.DEFAULT_LIMIT + 2):
        request, _ = testing_helpers.GetRequestObjects(
          project=self.project)
        request.headers['X-AppEngine-Country'] = 'US'
        request.remote_addr = '192.168.1.0'
        os.environ['USER_EMAIL'] = '42@example.com'
        ratelimiter._CacheKeys(request, now)
        self.ratelimiter.CheckStart(request, now)
        now = now + 0.001

    # Now proceed to make requests for other user IDs
    # besides 42.
    for m in range(ratelimiter.DEFAULT_LIMIT * 2):
      request, _ = testing_helpers.GetRequestObjects(
        project=self.project)
      request.headers['X-AppEngine-Country'] = 'US'
      # Skip .0 since it's already exceeded the limit.
      request.remote_addr = '192.168.1.0'
      os.environ['USER_EMAIL'] = '%s@example.com' % (43 + m)
      ratelimiter._CacheKeys(request, now)
      self.ratelimiter.CheckStart(request, now)
      now = now + 0.001

  def testCheckStart_ratelimitingDisabled(self):
    settings.ratelimiting_enabled = False
    request, _ = testing_helpers.GetRequestObjects(
      project=self.project)
    request.headers['X-AppEngine-Country'] = 'US'
    request.remote_addr = '192.168.1.0'
    now = 0.0

    # Call CheckStart a lot.  Should be ok.
    for _ in range(ratelimiter.DEFAULT_LIMIT):
      self.ratelimiter.CheckStart(request, now)
      now = now + 0.001

  def testCheckStart_perCountryLoggedOutLimit(self):
    ratelimiter.COUNTRY_LIMITS['US'] = 10

    request, _ = testing_helpers.GetRequestObjects(
      project=self.project)
    request.headers[ratelimiter.COUNTRY_HEADER] = 'US'
    request.remote_addr = '192.168.1.1'
    now = 0.0

    with self.assertRaises(ratelimiter.RateLimitExceeded):
      for m in range(ratelimiter.DEFAULT_LIMIT + 2):
        self.ratelimiter.CheckStart(request, now)
        # Vary remote address to make sure the limit covers
        # the whole country, regardless of IP.
        request.remote_addr = '192.168.1.%d' % m
        now = now + 0.001

    # CheckStart for a country that isn't covered by a country-specific limit.
    request.headers['X-AppEngine-Country'] = 'UK'
    for m in range(11):
      self.ratelimiter.CheckStart(request, now)
      # Vary remote address to make sure the limit covers
      # the whole country, regardless of IP.
      request.remote_addr = '192.168.1.%d' % m
      now = now + 0.001

    # And regular rate limits work per-IP.
    request.remote_addr = '192.168.1.1'
    with self.assertRaises(ratelimiter.RateLimitExceeded):
      for m in range(ratelimiter.DEFAULT_LIMIT):
        self.ratelimiter.CheckStart(request, now)
        # Vary remote address to make sure the limit covers
        # the whole country, regardless of IP.
        now = now + 0.001

  def testCheckEnd_overCostThresh(self):
    request, _ = testing_helpers.GetRequestObjects(
      project=self.project)
    request.headers[ratelimiter.COUNTRY_HEADER] = 'US'
    request.remote_addr = '192.168.1.1'
    start_time = 0.0

    # Send some requests, all under the limit.
    for _ in range(ratelimiter.DEFAULT_LIMIT-1):
      start_time = start_time + 0.001
      self.ratelimiter.CheckStart(request, start_time)
      now = start_time + 0.010
      self.ratelimiter.CheckEnd(request, now, start_time)

    # Now issue some more request, this time taking long
    # enough to get the cost threshold penalty.
    # Fast forward enough to impact a later bucket than the
    # previous requests.
    start_time = now + 120.0
    self.ratelimiter.CheckStart(request, start_time)

    # Take longer than the threshold to process the request.
    now = start_time + (settings.ratelimiting_cost_thresh_ms + 1) / 1000

    # The request finished, taking longer than the cost
    # threshold.
    self.ratelimiter.CheckEnd(request, now, start_time)

    with self.assertRaises(ratelimiter.RateLimitExceeded):
      # One more request after the expensive query should
      # throw an excpetion.
      self.ratelimiter.CheckStart(request, start_time)

  def testCheckEnd_overCostThreshButDisabled(self):
    request, _ = testing_helpers.GetRequestObjects(
      project=self.project)
    request.headers[ratelimiter.COUNTRY_HEADER] = 'US'
    request.remote_addr = '192.168.1.1'
    start_time = 0.0
    settings.ratelimiting_cost_enabled = False

    # Send some requests, all under the limit.
    for _ in range(ratelimiter.DEFAULT_LIMIT-1):
      start_time = start_time + 0.001
      self.ratelimiter.CheckStart(request, start_time)
      now = start_time + 0.010
      self.ratelimiter.CheckEnd(request, now, start_time)

    # Now issue some more request, this time taking long
    # enough to get the cost threshold penalty.
    # Fast forward enough to impact a later bucket than the
    # previous requests.
    start_time = now + 120.0
    self.ratelimiter.CheckStart(request, start_time)

    # Take longer than the threshold to process the request.
    now = start_time + (settings.ratelimiting_cost_thresh_ms + 10)/1000

    # The request finished, taking longer than the cost
    # threshold.
    self.ratelimiter.CheckEnd(request, now, start_time)

    # One more request after the expensive query should
    # throw an excpetion, but cost thresholds are disabled.
    self.ratelimiter.CheckStart(request, start_time)

  def testChekcEnd_underCostThresh(self):
    request, _ = testing_helpers.GetRequestObjects(
      project=self.project)
    request.headers[ratelimiter.COUNTRY_HEADER] = 'asdasd'
    request.remote_addr = '192.168.1.1'
    start_time = 0.0

    # Send some requests, all under the limit.
    for _ in range(ratelimiter.DEFAULT_LIMIT):
      self.ratelimiter.CheckStart(request, start_time)
      now = start_time + 0.010
      self.ratelimiter.CheckEnd(request, now, start_time)
      start_time = now + 0.010

  def testChekcEnd_underCostThresh(self):
    request, _ = testing_helpers.GetRequestObjects(
      project=self.project)
    request.headers[ratelimiter.COUNTRY_HEADER] = 'asdasd'
    request.remote_addr = '192.168.1.1'
    start_time = 0.0

    # Send some requests, all under the limit.
    for _ in range(ratelimiter.DEFAULT_LIMIT):
      self.ratelimiter.CheckStart(request, start_time)
      now = start_time + 0.01
      self.ratelimiter.CheckEnd(request, now, start_time)
      start_time = now + 0.01
