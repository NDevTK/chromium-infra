# Copyright 2016 The Chromium Authors. All rights reserved.
# Use of this source code is govered by a BSD-style
# license that can be found in the LICENSE file or at
# https://developers.google.com/open-source/licenses/bsd

"""Tests for the user profile page."""

import unittest

import mox

from framework import framework_helpers
from framework import framework_views
from proto import project_pb2
from proto import user_pb2
from services import service_manager
from sitewide import userprofile
from testing import fake


REGULAR_USER_ID = 111L
ADMIN_USER_ID = 222L
OTHER_USER_ID = 333L
STATES = {
    'live': project_pb2.ProjectState.LIVE,
    'archived': project_pb2.ProjectState.ARCHIVED,
}


def MakeReqInfo(
    user_pb, user_id, viewed_user_pb, viewed_user_id, viewed_user_name,
    _reveal_email=False, _params=None):
  mr = fake.MonorailRequest()
  mr.auth.user_pb = user_pb
  mr.auth.user_id = user_id
  mr.auth.effective_ids = {user_id}
  mr.viewed_user_auth.email = viewed_user_name
  mr.viewed_user_auth.user_pb = viewed_user_pb
  mr.viewed_user_auth.user_id = viewed_user_id
  mr.viewed_user_auth.effective_ids = {viewed_user_id}
  mr.viewed_user_auth.user_view = framework_views.UserView(
      viewed_user_id, viewed_user_pb.email, viewed_user_pb.obscure_email)
  mr.viewed_user_name = viewed_user_name
  return mr


class UserProfileTest(unittest.TestCase):

  def setUp(self):
    self.mox = mox.Mox()
    self.mox.StubOutWithMock(
        framework_helpers.UserSettings, 'GatherUnifiedSettingsPageData')

    services = service_manager.Services(
        project=fake.ProjectService(),
        user=fake.UserService(),
        usergroup=fake.UserGroupService(),
        project_star=fake.ProjectStarService(),
        user_star=fake.UserStarService())
    self.servlet = userprofile.UserProfile('req', 'res', services=services)

    for user_id in (
        REGULAR_USER_ID, ADMIN_USER_ID, OTHER_USER_ID):
      services.user.TestAddUser('%s@gmail.com' % user_id, user_id)

    for user in ['regular', 'other']:
      for relation in ['owner', 'member']:
        for state_name, state in STATES.iteritems():
          services.project.TestAddProject(
              '%s-%s-%s' % (user, relation, state_name), state=state)

    # Add projects
    for state_name, state in STATES.iteritems():
      services.project.TestAddProject(
          'regular-owner-%s' % state_name, state=state,
          owner_ids=[REGULAR_USER_ID])
      services.project.TestAddProject(
          'regular-member-%s' % state_name, state=state,
          committer_ids=[REGULAR_USER_ID])
      services.project.TestAddProject(
          'other-owner-%s' % state_name, state=state,
          owner_ids=[OTHER_USER_ID])
      services.project.TestAddProject(
          'other-member-%s' % state_name, state=state,
          committer_ids=[OTHER_USER_ID])

    self.regular_user = services.user.GetUser('fake cnxn', REGULAR_USER_ID)
    self.admin_user = services.user.GetUser('fake cnxn', ADMIN_USER_ID)
    self.admin_user.is_site_admin = True
    self.other_user = services.user.GetUser('fake cnxn', OTHER_USER_ID)

  def tearDown(self):
    self.mox.UnsetStubs()

  def assertProjectsAnyOrder(self, value_to_test, *expected_project_names):
    actual_project_names = [project_view.project_name
                            for project_view in value_to_test]
    self.assertItemsEqual(expected_project_names, actual_project_names)

    # TODO(jrobbins): re-implement captchas to reveal full
    # email address and add tests for that.

  def testGatherPageData_RegularUserViewingOtherUserProjects(self):
    """A user can see the other users' live projects, but not archived ones."""
    mr = MakeReqInfo(
        self.regular_user, REGULAR_USER_ID, self.other_user,
        OTHER_USER_ID, 'other@xyz.com')

    framework_helpers.UserSettings.GatherUnifiedSettingsPageData(
        111L, mr.viewed_user_auth.user_view,
        mr.viewed_user_auth.user_pb).AndReturn({'unified': None})
    self.mox.ReplayAll()

    page_data = self.servlet.GatherPageData(mr)
    self.assertProjectsAnyOrder(page_data['owner_of_projects'],
                                'other-owner-live')
    self.assertProjectsAnyOrder(page_data['committer_of_projects'],
                                'other-member-live')
    self.assertFalse(page_data['owner_of_archived_projects'])
    self.assertEquals('ot...@xyz.com', page_data['viewed_user_display_name'])

    self.mox.VerifyAll()

  def testGatherPageData_RegularUserViewingOwnProjects(self):
    """A user can see all their own projects: live or archived."""
    mr = MakeReqInfo(
        self.regular_user, REGULAR_USER_ID, self.regular_user,
        REGULAR_USER_ID, 'self@xyz.com')

    framework_helpers.UserSettings.GatherUnifiedSettingsPageData(
        111L, mr.viewed_user_auth.user_view,
        mr.viewed_user_auth.user_pb).AndReturn({'unified': None})
    self.mox.ReplayAll()

    page_data = self.servlet.GatherPageData(mr)
    self.assertEquals('self@xyz.com', page_data['viewed_user_display_name'])
    self.assertProjectsAnyOrder(page_data['owner_of_projects'],
                                'regular-owner-live')
    self.assertProjectsAnyOrder(page_data['committer_of_projects'],
                                'regular-member-live')
    self.assertProjectsAnyOrder(
        page_data['owner_of_archived_projects'],
        'regular-owner-archived')

    self.mox.VerifyAll()

  def testGatherPageData_AdminViewingOtherUserAddress(self):
    """Site admins always see full email addresses of other users."""
    mr = MakeReqInfo(
        self.admin_user, ADMIN_USER_ID, self.other_user,
        OTHER_USER_ID, 'other@xyz.com')

    framework_helpers.UserSettings.GatherUnifiedSettingsPageData(
        222L, mr.viewed_user_auth.user_view,
        mr.viewed_user_auth.user_pb).AndReturn({'unified': None})
    self.mox.ReplayAll()

    page_data = self.servlet.GatherPageData(mr)
    self.assertEquals('other@xyz.com', page_data['viewed_user_display_name'])

    self.mox.VerifyAll()

  def testGatherPageData_RegularUserViewingOtherUserAddress(self):
    """Email should be revealed to others depending on obscure_email."""
    mr = MakeReqInfo(
        self.regular_user, REGULAR_USER_ID, self.other_user,
        OTHER_USER_ID, 'other@xyz.com')

    framework_helpers.UserSettings.GatherUnifiedSettingsPageData(
        111L, mr.viewed_user_auth.user_view,
        mr.viewed_user_auth.user_pb).AndReturn({'unified': None})
    framework_helpers.UserSettings.GatherUnifiedSettingsPageData(
        111L, mr.viewed_user_auth.user_view,
        mr.viewed_user_auth.user_pb).AndReturn({'unified': None})
    self.mox.ReplayAll()

    mr.viewed_user_auth.user_view.obscure_email = False
    page_data = self.servlet.GatherPageData(mr)
    self.assertEquals('other@xyz.com', page_data['viewed_user_display_name'])

    mr.viewed_user_auth.user_view.obscure_email = True
    page_data = self.servlet.GatherPageData(mr)
    self.assertEquals('ot...@xyz.com', page_data['viewed_user_display_name'])

    self.mox.VerifyAll()


if __name__ == '__main__':
  unittest.main()
