# Copyright 2016 The Chromium Authors. All rights reserved.
# Use of this source code is govered by a BSD-style
# license that can be found in the LICENSE file or at
# https://developers.google.com/open-source/licenses/bsd

"""Unit tests for the monorailrequest module."""

import re
import unittest

import mox

from google.appengine.api import users

import webapp2

from framework import monorailrequest
from framework import permissions
from framework import profiler
from proto import project_pb2
from proto import tracker_pb2
from services import service_manager
from testing import fake
from testing import testing_helpers
from tracker import tracker_constants


class HostportReTest(unittest.TestCase):

  def testGood(self):
    test_data = [
      'localhost:8080',
      'app.appspot.com',
      'bugs-staging.chromium.org',
      'vers10n-h3x-dot-app-id.appspot.com',
      ]
    for hostport in test_data:
      self.assertTrue(monorailrequest._HOSTPORT_RE.match(hostport),
                      msg='Incorrectly rejected %r' % hostport)

  def testBad(self):
    test_data = [
      '',
      ' ',
      '\t',
      '\n',
      '\'',
      '"',
      'version"cruft-dot-app-id.appspot.com',
      '\nother header',
      'version&cruft-dot-app-id.appspot.com',
      ]
    for hostport in test_data:
      self.assertFalse(monorailrequest._HOSTPORT_RE.match(hostport),
                       msg='Incorrectly accepted %r' % hostport)

class AuthDataTest(unittest.TestCase):

  def setUp(self):
    self.mox = mox.Mox()

  def tearDown(self):
    self.mox.UnsetStubs()

  def testGetUserID(self):
    pass  # TODO(jrobbins): re-impement

  def testExamineRequestUserID(self):
    pass  # TODO(jrobbins): re-implement


class MonorailRequestUnitTest(unittest.TestCase):

  def setUp(self):
    self.services = service_manager.Services(
        project=fake.ProjectService(),
        user=fake.UserService(),
        usergroup=fake.UserGroupService())
    self.project = self.services.project.TestAddProject('proj')
    self.services.user.TestAddUser('jrobbins@example.com', 111)

    self.profiler = profiler.Profiler()
    self.mox = mox.Mox()
    self.mox.StubOutWithMock(users, 'get_current_user')
    users.get_current_user().AndReturn(None)
    self.mox.ReplayAll()

  def tearDown(self):
    self.mox.UnsetStubs()

  def testGetIntParamConvertsQueryParamToInt(self):
    notice_id = 12345
    mr = testing_helpers.MakeMonorailRequest(
        path='/foo?notice=%s' % notice_id)

    value = mr.GetIntParam('notice')
    self.assert_(isinstance(value, int))
    self.assertEqual(notice_id, value)

  def testGetIntParamConvertsQueryParamToLong(self):
    notice_id = 12345678901234567890
    mr = testing_helpers.MakeMonorailRequest(
        path='/foo?notice=%s' % notice_id)

    value = mr.GetIntParam('notice')
    self.assertTrue(isinstance(value, long))
    self.assertEqual(notice_id, value)

  def testGetIntListParamNoParam(self):
    mr = monorailrequest.MonorailRequest()
    mr.ParseRequest(
        webapp2.Request.blank('servlet'), self.services, self.profiler)
    self.assertEquals(mr.GetIntListParam('ids'), None)
    self.assertEquals(mr.GetIntListParam('ids', default_value=['test']),
                      ['test'])

  def testGetIntListParamOneValue(self):
    mr = monorailrequest.MonorailRequest()
    mr.ParseRequest(
        webapp2.Request.blank('servlet?ids=11'), self.services, self.profiler)
    self.assertEquals(mr.GetIntListParam('ids'), [11])
    self.assertEquals(mr.GetIntListParam('ids', default_value=['test']),
                      [11])

  def testGetIntListParamMultiValue(self):
    mr = monorailrequest.MonorailRequest()
    mr.ParseRequest(
        webapp2.Request.blank('servlet?ids=21,22,23'), self.services,
        self.profiler)
    self.assertEquals(mr.GetIntListParam('ids'), [21, 22, 23])
    self.assertEquals(mr.GetIntListParam('ids', default_value=['test']),
                      [21, 22, 23])

  def testGetIntListParamBogusValue(self):
    mr = monorailrequest.MonorailRequest()
    mr.ParseRequest(
        webapp2.Request.blank('servlet?ids=not_an_int'), self.services,
        self.profiler)
    self.assertEquals(mr.GetIntListParam('ids'), None)
    self.assertEquals(mr.GetIntListParam('ids', default_value=['test']),
                      ['test'])

  def testGetIntListParamMalformed(self):
    mr = monorailrequest.MonorailRequest()
    mr.ParseRequest(
        webapp2.Request.blank('servlet?ids=31,32,,'), self.services,
        self.profiler)
    self.assertEquals(mr.GetIntListParam('ids'), None)
    self.assertEquals(mr.GetIntListParam('ids', default_value=['test']),
                      ['test'])

  def testDefaultValuesNoUrl(self):
    """If request has no param, default param values should be used."""
    mr = monorailrequest.MonorailRequest()
    mr.ParseRequest(
        webapp2.Request.blank('servlet'), self.services, self.profiler)
    self.assertEquals(mr.GetParam('r', 3), 3)
    self.assertEquals(mr.GetIntParam('r', 3), 3)
    self.assertEquals(mr.GetPositiveIntParam('r', 3), 3)
    self.assertEquals(mr.GetIntListParam('r', [3, 4]), [3, 4])

  def _MRWithMockRequest(
      self, path, headers=None, *mr_args, **mr_kwargs):
    request = webapp2.Request.blank(path, headers=headers)
    mr = monorailrequest.MonorailRequest(*mr_args, **mr_kwargs)
    mr.ParseRequest(request, self.services, self.profiler)
    return mr

  def testParseQueryParameters(self):
    mr = self._MRWithMockRequest(
        '/p/proj/issues/list?q=foo+OR+bar&num=50')
    self.assertEquals('foo OR bar', mr.query)
    self.assertEquals(50, mr.num)

  def testParseRequest_Scheme(self):
    mr = self._MRWithMockRequest('/p/proj/')
    self.assertEquals('http', mr.request.scheme)

  def testParseRequest_HostportAndCurrentPageURL(self):
    mr = self._MRWithMockRequest('/p/proj/', headers={
        'Host': 'example.com',
        'Cookie': 'asdf',
        })
    self.assertEquals('http', mr.request.scheme)
    self.assertEquals('example.com', mr.request.host)
    self.assertEquals('http://example.com/p/proj/', mr.current_page_url)

  def testViewedUser_WithEmail(self):
    mr = self._MRWithMockRequest('/u/jrobbins@example.com/')
    self.assertEquals('jrobbins@example.com', mr.viewed_username)
    self.assertEquals(111, mr.viewed_user_auth.user_id)
    self.assertEquals(
        self.services.user.GetUser('fake cnxn', 111),
        mr.viewed_user_auth.user_pb)

  def testViewedUser_WithUserID(self):
    mr = self._MRWithMockRequest('/u/111/')
    self.assertEquals('jrobbins@example.com', mr.viewed_username)
    self.assertEquals(111, mr.viewed_user_auth.user_id)
    self.assertEquals(
        self.services.user.GetUser('fake cnxn', 111),
        mr.viewed_user_auth.user_pb)

  def testViewedUser_NoSuchEmail(self):
    with self.assertRaises(webapp2.HTTPException) as cm:
      self._MRWithMockRequest('/u/unknownuser@example.com/')
    self.assertEquals(404, cm.exception.code)

  def testViewedUser_NoSuchUserID(self):
    with self.assertRaises(webapp2.HTTPException) as cm:
      self._MRWithMockRequest('/u/234521111/')
    self.assertEquals(404, cm.exception.code)

  def testGetParam(self):
    mr = testing_helpers.MakeMonorailRequest(
        path='/foo?syn=error!&a=a&empty=',
        params=dict(over1='over_value1', over2='over_value2'))

    # test tampering
    self.assertRaises(monorailrequest.InputException, mr.GetParam, 'a',
                      antitamper_re=re.compile(r'^$'))
    self.assertRaises(monorailrequest.InputException, mr.GetParam,
                      'undefined', default_value='default',
                      antitamper_re=re.compile(r'^$'))

    # test empty value
    self.assertEquals('', mr.GetParam(
        'empty', default_value='default', antitamper_re=re.compile(r'^$')))

    # test default
    self.assertEquals('default', mr.GetParam(
        'undefined', default_value='default'))

  def testComputeColSpec(self):
    # No config passed, and nothing in URL
    mr = testing_helpers.MakeMonorailRequest(
        path='/p/proj/issues/detail?id=123')
    mr.ComputeColSpec(None)
    self.assertEquals(tracker_constants.DEFAULT_COL_SPEC, mr.col_spec)

    # No config passed, but set in URL
    mr = testing_helpers.MakeMonorailRequest(
        path='/p/proj/issues/detail?id=123&colspec=a b C')
    mr.ComputeColSpec(None)
    self.assertEquals('a b C', mr.col_spec)

    config = tracker_pb2.ProjectIssueConfig()

    # No default in the config, and nothing in URL
    mr = testing_helpers.MakeMonorailRequest(
        path='/p/proj/issues/detail?id=123')
    mr.ComputeColSpec(config)
    self.assertEquals(tracker_constants.DEFAULT_COL_SPEC, mr.col_spec)

    # No default in the config, but set in URL
    mr = testing_helpers.MakeMonorailRequest(
        path='/p/proj/issues/detail?id=123&colspec=a b C')
    mr.ComputeColSpec(config)
    self.assertEquals('a b C', mr.col_spec)

    config.default_col_spec = 'd e f'

    # Default in the config, and nothing in URL
    mr = testing_helpers.MakeMonorailRequest(
        path='/p/proj/issues/detail?id=123')
    mr.ComputeColSpec(config)
    self.assertEquals('d e f', mr.col_spec)

    # Default in the config, but overrided via URL
    mr = testing_helpers.MakeMonorailRequest(
        path='/p/proj/issues/detail?id=123&colspec=a b C')
    mr.ComputeColSpec(config)
    self.assertEquals('a b C', mr.col_spec)

  def testComputeColSpec_XSS(self):
    config_1 = tracker_pb2.ProjectIssueConfig()
    config_2 = tracker_pb2.ProjectIssueConfig()
    config_2.default_col_spec = "id '+alert(1)+'"
    mr_1 = testing_helpers.MakeMonorailRequest(
        path='/p/proj/issues/detail?id=123')
    mr_2 = testing_helpers.MakeMonorailRequest(
        path="/p/proj/issues/detail?id=123&colspec=id '+alert(1)+'")

    # Normal colspec in config but malicious request
    self.assertRaises(
        monorailrequest.InputException,
        mr_2.ComputeColSpec, config_1)

    # Malicious colspec in config but normal request
    self.assertRaises(
        monorailrequest.InputException,
        mr_1.ComputeColSpec, config_2)

    # Malicious colspec in config and malicious request
    self.assertRaises(
        monorailrequest.InputException,
        mr_2.ComputeColSpec, config_2)


class TestMonorailRequestFunctions(unittest.TestCase):

  def testExtractPathIdenifiers_ProjectOnly(self):
    username, project_name = monorailrequest._ParsePathIdentifiers(
        '/p/proj/issues/list?q=foo+OR+bar&ts=1234')
    self.assertIsNone(username)
    self.assertEquals('proj', project_name)

  def testExtractPathIdenifiers_ViewedUserOnly(self):
    username, project_name = monorailrequest._ParsePathIdentifiers(
        '/u/jrobbins@example.com/')
    self.assertEquals('jrobbins@example.com', username)
    self.assertIsNone(project_name)

  def testExtractPathIdenifiers_ViewedUserURLSpace(self):
    username, project_name = monorailrequest._ParsePathIdentifiers(
        '/u/jrobbins@example.com/updates')
    self.assertEquals('jrobbins@example.com', username)
    self.assertIsNone(project_name)

  def testExtractPathIdenifiers_ViewedGroupURLSpace(self):
    username, project_name = monorailrequest._ParsePathIdentifiers(
        '/g/user-group@example.com/updates')
    self.assertEquals('user-group@example.com', username)
    self.assertIsNone(project_name)

  def testParseColSpec(self):
    parse = monorailrequest.ParseColSpec
    self.assertEqual(['PageName', 'Summary', 'Changed', 'ChangedBy'],
                     parse(u'PageName Summary Changed ChangedBy'))
    self.assertEqual(['Foo-Bar', 'Foo-Bar-Baz', 'Release-1.2', 'Hey', 'There'],
                     parse('Foo-Bar Foo-Bar-Baz Release-1.2 Hey!There'))
    self.assertEqual(
        ['\xe7\xaa\xbf\xe8\x8b\xa5\xe7\xb9\xb9'.decode('utf-8'),
         '\xe5\x9f\xba\xe5\x9c\xb0\xe3\x81\xaf'.decode('utf-8')],
        parse('\xe7\xaa\xbf\xe8\x8b\xa5\xe7\xb9\xb9 '
              '\xe5\x9f\xba\xe5\x9c\xb0\xe3\x81\xaf'.decode('utf-8')))


class TestPermissionLookup(unittest.TestCase):
  OWNER_ID = 1
  OTHER_USER_ID = 2

  def setUp(self):
    self.services = service_manager.Services(
        project=fake.ProjectService(),
        user=fake.UserService(),
        usergroup=fake.UserGroupService())
    self.services.user.TestAddUser('owner@gmail.com', self.OWNER_ID)
    self.services.user.TestAddUser('user@gmail.com', self.OTHER_USER_ID)
    self.live_project = self.services.project.TestAddProject(
        'live', owner_ids=[self.OWNER_ID])
    self.archived_project = self.services.project.TestAddProject(
        'archived', owner_ids=[self.OWNER_ID],
        state=project_pb2.ProjectState.ARCHIVED)
    self.members_only_project = self.services.project.TestAddProject(
        'members-only', owner_ids=[self.OWNER_ID],
        access=project_pb2.ProjectAccess.MEMBERS_ONLY)

    self.mox = mox.Mox()

  def tearDown(self):
    self.mox.UnsetStubs()

  def CheckPermissions(self, perms, expect_view, expect_commit, expect_edit):
    may_view = perms.HasPerm(permissions.VIEW, None, None)
    self.assertEqual(expect_view, may_view)
    may_commit = perms.HasPerm(permissions.COMMIT, None, None)
    self.assertEqual(expect_commit, may_commit)
    may_edit = perms.HasPerm(permissions.EDIT_PROJECT, None, None)
    self.assertEqual(expect_edit, may_edit)

  def MakeRequestAsUser(self, project_name, email):
    self.mox.StubOutWithMock(users, 'get_current_user')
    users.get_current_user().AndReturn(testing_helpers.Blank(
        email=lambda: email))
    self.mox.ReplayAll()

    request = webapp2.Request.blank('/p/' + project_name)
    mr = monorailrequest.MonorailRequest()
    prof = profiler.Profiler()
    with prof.Phase('parse user info'):
      mr.ParseRequest(request, self.services, prof)
    return mr

  def testOwnerPermissions_Live(self):
    mr = self.MakeRequestAsUser('live', 'owner@gmail.com')
    self.CheckPermissions(mr.perms, True, True, True)

  def testOwnerPermissions_Archived(self):
    mr = self.MakeRequestAsUser('archived', 'owner@gmail.com')
    self.CheckPermissions(mr.perms, True, False, True)

  def testOwnerPermissions_MembersOnly(self):
    mr = self.MakeRequestAsUser('members-only', 'owner@gmail.com')
    self.CheckPermissions(mr.perms, True, True, True)

  def testExternalUserPermissions_Live(self):
    mr = self.MakeRequestAsUser('live', 'user@gmail.com')
    self.CheckPermissions(mr.perms, True, False, False)

  def testExternalUserPermissions_Archived(self):
    mr = self.MakeRequestAsUser('archived', 'user@gmail.com')
    self.CheckPermissions(mr.perms, False, False, False)

  def testExternalUserPermissions_MembersOnly(self):
    mr = self.MakeRequestAsUser('members-only', 'user@gmail.com')
    self.CheckPermissions(mr.perms, False, False, False)
