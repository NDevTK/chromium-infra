# Copyright 2016 The Chromium Authors. All rights reserved.
# Use of this source code is govered by a BSD-style
# license that can be found in the LICENSE file or at
# https://developers.google.com/open-source/licenses/bsd

"""Unit tests for the component_helpers module."""

import unittest

from proto import tracker_pb2
from services import service_manager
from testing import fake
from tracker import component_helpers
from tracker import tracker_bizobj


class ComponentHelpersTest(unittest.TestCase):

  def setUp(self):
    self.config = tracker_bizobj.MakeDefaultProjectIssueConfig(789)
    self.cd1 = tracker_bizobj.MakeComponentDef(
        1, 789, 'FrontEnd', 'doc', False, [], [111L], 0, 0)
    self.cd2 = tracker_bizobj.MakeComponentDef(
        2, 789, 'FrontEnd>Splash', 'doc', False, [], [222L], 0, 0)
    self.cd3 = tracker_bizobj.MakeComponentDef(
        3, 789, 'BackEnd', 'doc', True, [], [111L, 333L], 0, 0)
    self.config.component_defs = [self.cd1, self.cd2, self.cd3]
    self.services = service_manager.Services(
        user=fake.UserService())
    self.services.user.TestAddUser('a@example.com', 111L)
    self.services.user.TestAddUser('b@example.com', 222L)
    self.services.user.TestAddUser('c@example.com', 333L)
    self.mr = fake.MonorailRequest()
    self.mr.cnxn = fake.MonorailConnection()

  def testParseComponentRequest_Empty(self):
    post_data = fake.PostData(admins=[''], cc=[''])
    parsed = component_helpers.ParseComponentRequest(
        self.mr, post_data, self.services.user)
    self.assertEqual('', parsed.leaf_name)
    self.assertEqual('', parsed.docstring)
    self.assertEqual([], parsed.admin_usernames)
    self.assertEqual([], parsed.cc_usernames)
    self.assertEqual([], parsed.admin_ids)
    self.assertEqual([], parsed.cc_ids)
    self.assertFalse(self.mr.errors.AnyErrors())

  def testParseComponentRequest_Normal(self):
    post_data = fake.PostData(
        leaf_name=['FrontEnd'],
        docstring=['The server-side app that serves pages'],
        deprecated=[False],
        admins=['a@example.com'],
        cc=['b@example.com, c@example.com'])
    parsed = component_helpers.ParseComponentRequest(
        self.mr, post_data, self.services.user)
    self.assertEqual('FrontEnd', parsed.leaf_name)
    self.assertEqual('The server-side app that serves pages', parsed.docstring)
    self.assertEqual(['a@example.com'], parsed.admin_usernames)
    self.assertEqual(['b@example.com', 'c@example.com'], parsed.cc_usernames)
    self.assertEqual([111L], parsed.admin_ids)
    self.assertEqual([222L, 333L], parsed.cc_ids)
    self.assertFalse(self.mr.errors.AnyErrors())

  def testParseComponentRequest_InvalidUser(self):
    post_data = fake.PostData(
        leaf_name=['FrontEnd'],
        docstring=['The server-side app that serves pages'],
        deprecated=[False],
        admins=['a@example.com, invalid_user'],
        cc=['b@example.com, c@example.com'])
    parsed = component_helpers.ParseComponentRequest(
        self.mr, post_data, self.services.user)
    self.assertEqual('FrontEnd', parsed.leaf_name)
    self.assertEqual('The server-side app that serves pages', parsed.docstring)
    self.assertEqual(['a@example.com', 'invalid_user'], parsed.admin_usernames)
    self.assertEqual(['b@example.com', 'c@example.com'], parsed.cc_usernames)
    self.assertEqual([111L], parsed.admin_ids)
    self.assertEqual([222L, 333L], parsed.cc_ids)
    self.assertTrue(self.mr.errors.AnyErrors())
    self.assertEqual('invalid_user unrecognized', self.mr.errors.member_admins)

  def testGetComponentCcIDs(self):
    issue = tracker_pb2.Issue()
    issues_components_cc_ids = component_helpers.GetComponentCcIDs(
        issue, self.config)
    self.assertEqual(set(), issues_components_cc_ids)

    issue.component_ids = [1, 2]
    issues_components_cc_ids = component_helpers.GetComponentCcIDs(
        issue, self.config)
    self.assertEqual({111L, 222L}, issues_components_cc_ids)

  def testGetCcIDsForComponentAndAncestors(self):
    components_cc_ids = component_helpers.GetCcIDsForComponentAndAncestors(
        self.config, self.cd1)
    self.assertEqual({111L}, components_cc_ids)

    components_cc_ids = component_helpers.GetCcIDsForComponentAndAncestors(
        self.config, self.cd2)
    self.assertEqual({111L, 222L}, components_cc_ids)

    components_cc_ids = component_helpers.GetCcIDsForComponentAndAncestors(
        self.config, self.cd3)
    self.assertEqual({111L, 333L}, components_cc_ids)
