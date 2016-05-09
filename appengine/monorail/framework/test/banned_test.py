# Copyright 2016 The Chromium Authors. All rights reserved.
# Use of this source code is govered by a BSD-style
# license that can be found in the LICENSE file or at
# https://developers.google.com/open-source/licenses/bsd

"""Unittests for monorail.framework.banned."""

import unittest

import webapp2

from framework import banned
from framework import monorailrequest
from services import service_manager
from testing import testing_helpers


class BannedTest(unittest.TestCase):

  def testAssertBasePermission(self):
    servlet = banned.Banned(
        'request', 'response', services=service_manager.Services())

    mr = monorailrequest.MonorailRequest()
    mr.auth.user_id = 0L  # Anon user cannot see banned page.
    with self.assertRaises(webapp2.HTTPException) as cm:
      servlet.AssertBasePermission(mr)
    self.assertEquals(404, cm.exception.code)

    mr.auth.user_id = 111L  # User who is not banned cannot view banned page.
    with self.assertRaises(webapp2.HTTPException) as cm:
      servlet.AssertBasePermission(mr)
    self.assertEquals(404, cm.exception.code)

    # This should not throw exception.
    mr.auth.user_pb.banned = 'spammer'
    servlet.AssertBasePermission(mr)

  def testGatherPageData(self):
    servlet = banned.Banned(
        'request', 'response', services=service_manager.Services())
    self.assertNotEquals(servlet.template, None)

    _request, mr = testing_helpers.GetRequestObjects()
    page_data = servlet.GatherPageData(mr)

    self.assertEquals(None, page_data['currentPageURLEncoded'])
