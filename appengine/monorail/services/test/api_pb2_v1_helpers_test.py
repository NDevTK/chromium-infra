# Copyright 2016 The Chromium Authors. All rights reserved.
# Use of this source code is govered by a BSD-style
# license that can be found in the LICENSE file or at
# https://developers.google.com/open-source/licenses/bsd

"""Tests for the API v1 helpers."""

import datetime
import mock
import unittest

from framework import framework_constants
from framework import permissions
from services import api_pb2_v1_helpers
from services import service_manager
from proto import api_pb2_v1
from proto import project_pb2
from proto import tracker_pb2
from proto import usergroup_pb2
from testing import fake


def MakeTemplate(prefix):
  return tracker_pb2.TemplateDef(
      name='%s-template' % prefix,
      content='%s-content' % prefix,
      summary='%s-summary' % prefix,
      summary_must_be_edited=True,
      status='New',
      labels=['%s-label1' % prefix, '%s-label2' % prefix],
      members_only=True,
      owner_defaults_to_member=True,
      component_required=True,
  )


def MakeLabel(prefix):
  return tracker_pb2.LabelDef(
      label='%s-label' % prefix,
      label_docstring='%s-description' % prefix
  )


def MakeStatus(prefix):
  return tracker_pb2.StatusDef(
      status='%s-New' % prefix,
      means_open=True,
      status_docstring='%s-status' % prefix
  )


def MakeProjectIssueConfig(prefix):
  return tracker_pb2.ProjectIssueConfig(
      restrict_to_known=True,
      default_col_spec='ID Type Priority Summary',
      default_sort_spec='ID Priority',
      well_known_statuses=[
          MakeStatus('%s-status1' % prefix),
          MakeStatus('%s-status2' % prefix),
      ],
      well_known_labels=[
          MakeLabel('%s-label1' % prefix),
          MakeLabel('%s-label2' % prefix),
      ],
      templates=[
          MakeTemplate('%s-template1' % prefix),
          MakeTemplate('%s-template2' % prefix),
      ],
      default_template_for_developers=1,
      default_template_for_users=2
  )


def MakeProject(prefix):
  return project_pb2.MakeProject(
      project_name='%s-project' % prefix,
      summary='%s-summary' % prefix,
      description='%s-description' % prefix,
  )


class ApiV1HelpersTest(unittest.TestCase):

  def testConvertTemplate(self):
    """Test convert_template."""
    template = MakeTemplate('test')
    prompt = api_pb2_v1_helpers.convert_template(template)
    self.assertEquals(template.name, prompt.name)
    self.assertEquals(template.summary, prompt.title)
    self.assertEquals(template.content, prompt.description)
    self.assertEquals(
        template.summary_must_be_edited, prompt.titleMustBeEdited)
    self.assertEquals(template.status, prompt.status)
    self.assertEquals(template.labels, prompt.labels)
    self.assertEquals(template.members_only, prompt.membersOnly)
    self.assertEquals(
        template.owner_defaults_to_member, prompt.defaultToMember)
    self.assertEquals(template.component_required, prompt.componentRequired)

  def testConvertLabel(self):
    """Test convert_label."""
    labeldef = MakeLabel('test')
    label = api_pb2_v1_helpers.convert_label(labeldef)
    self.assertEquals(labeldef.label, label.label)
    self.assertEquals(labeldef.label_docstring, label.description)

  def testConvertStatus(self):
    """Test convert_status."""
    statusdef = MakeStatus('test')
    status = api_pb2_v1_helpers.convert_status(statusdef)
    self.assertEquals(statusdef.status, status.status)
    self.assertEquals(statusdef.means_open, status.meansOpen)
    self.assertEquals(statusdef.status_docstring, status.description)

  def testConvertProjectIssueConfig(self):
    """Test convert_project_config."""
    config = MakeProjectIssueConfig('test')
    config_api = api_pb2_v1_helpers.convert_project_config(config)
    self.assertEquals(config.restrict_to_known, config_api.restrictToKnown)
    self.assertEquals(
        config.default_col_spec.split(), config_api.defaultColumns)
    self.assertEquals(
        config.default_sort_spec.split(), config_api.defaultSorting)
    self.assertEquals(2, len(config_api.statuses))
    self.assertEquals(2, len(config_api.labels))
    self.assertEquals(2, len(config_api.prompts))
    self.assertEquals(
        config.default_template_for_developers,
        config_api.defaultPromptForMembers)
    self.assertEquals(
        config.default_template_for_users,
        config_api.defaultPromptForNonMembers)

  def testConvertProject(self):
    """Test convert_project."""
    project = MakeProject('testprj')
    config = MakeProjectIssueConfig('testconfig')
    role = api_pb2_v1.Role.owner
    project_api = api_pb2_v1_helpers.convert_project(project, config, role)
    self.assertEquals(project.project_name, project_api.name)
    self.assertEquals(project.project_name, project_api.externalId)
    self.assertEquals('/p/%s/' % project.project_name, project_api.htmlLink)
    self.assertEquals(project.summary, project_api.summary)
    self.assertEquals(project.description, project_api.description)
    self.assertEquals(role, project_api.role)
    self.assertIsInstance(
        project_api.issuesConfig, api_pb2_v1.ProjectIssueConfig)

  def testConvertPerson(self):
    """Test convert_person."""
    svcs = service_manager.Services()
    svcs.user = fake.UserService()
    svcs.user.TestAddUser('user@example.com', 1)
    result = api_pb2_v1_helpers.convert_person(1, None, svcs)
    self.assertIsInstance(result, api_pb2_v1.AtomPerson)
    self.assertEquals('user@example.com', result.name)

  def testConvertIssueIDs(self):
    """Test convert_issue_ids."""
    svcs = service_manager.Services()
    svcs.issue = fake.IssueService()
    issue1 = fake.MakeTestIssue(789, 1, 'one', 'New', 111L)
    svcs.issue.TestAddIssue(issue1)
    issue_ids = [100001]
    mar = mock.Mock()
    mar.cnxn = None
    mar.project_name = 'test-project'
    result = api_pb2_v1_helpers.convert_issue_ids(issue_ids, mar, svcs)
    self.assertEquals(1, len(result))
    self.assertEquals(1, result[0].issueId)

  def testConvertIssueRef(self):
    """Test convert_issueref_pbs."""
    svcs = service_manager.Services()
    svcs.issue = fake.IssueService()
    issue1 = fake.MakeTestIssue(12345, 1, 'one', 'New', 111L)
    svcs.issue.TestAddIssue(issue1)
    svcs.project = fake.ProjectService()
    svcs.project.TestAddProject(
        'test-project', owner_ids=[2],
        project_id=12345)
    mar = mock.Mock()
    mar.cnxn = None
    mar.project_name = 'test-project'
    mar.project_id = 12345
    ir = api_pb2_v1.IssueRef(
        issueId=1,
        projectId='test-project'
    )
    result = api_pb2_v1_helpers.convert_issueref_pbs([ir], mar, svcs)
    self.assertEquals(1, len(result))
    self.assertEquals(100001, result[0])

  def testConvertAttachment(self):
    """Test convert_attachment."""

    attachment = tracker_pb2.Attachment(
        attachment_id=1,
        filename='stats.txt',
        filesize=12345,
        mimetype='text/plain',
        deleted=False)

    result = api_pb2_v1_helpers.convert_attachment(attachment)
    self.assertEquals(attachment.attachment_id, result.attachmentId)
    self.assertEquals(attachment.filename, result.fileName)
    self.assertEquals(attachment.filesize, result.fileSize)
    self.assertEquals(attachment.mimetype, result.mimetype)
    self.assertEquals(attachment.deleted, result.isDeleted)

  def testConvertAmendments(self):
    """Test convert_amendments."""

    svcs = service_manager.Services()
    svcs.user = fake.UserService()
    svcs.user.TestAddUser('user@example.com', 1)
    svcs.user.TestAddUser('user2@example.com', 2)
    mar = mock.Mock()
    mar.cnxn = None
    issue = mock.Mock()
    issue.project_name = 'test-project'

    amendment_summary = tracker_pb2.Amendment(
        field=tracker_pb2.FieldID.SUMMARY,
        newvalue='new summary')
    amendment_status = tracker_pb2.Amendment(
        field=tracker_pb2.FieldID.STATUS,
        newvalue='new status')
    amendment_owner = tracker_pb2.Amendment(
        field=tracker_pb2.FieldID.OWNER,
        added_user_ids=[1])
    amendment_labels = tracker_pb2.Amendment(
        field=tracker_pb2.FieldID.LABELS,
        newvalue='label1 -label2')
    amendment_cc_add = tracker_pb2.Amendment(
        field=tracker_pb2.FieldID.CC,
        added_user_ids=[1])
    amendment_cc_remove = tracker_pb2.Amendment(
        field=tracker_pb2.FieldID.CC,
        removed_user_ids=[2])
    amendment_blockedon = tracker_pb2.Amendment(
        field=tracker_pb2.FieldID.BLOCKEDON,
        newvalue='1')
    amendment_blocking = tracker_pb2.Amendment(
        field=tracker_pb2.FieldID.BLOCKING,
        newvalue='other:2 -3')
    amendment_mergedinto = tracker_pb2.Amendment(
        field=tracker_pb2.FieldID.MERGEDINTO,
        newvalue='4')
    amendments = [
        amendment_summary, amendment_status, amendment_owner,
        amendment_labels, amendment_cc_add, amendment_cc_remove,
        amendment_blockedon, amendment_blocking, amendment_mergedinto]

    result = api_pb2_v1_helpers.convert_amendments(
        issue, amendments, mar, svcs)
    self.assertEquals(amendment_summary.newvalue, result.summary)
    self.assertEquals(amendment_status.newvalue, result.status)
    self.assertEquals('user@example.com', result.owner)
    self.assertEquals(['label1', '-label2'], result.labels)
    self.assertEquals(['user@example.com', '-user2@example.com'], result.cc)
    self.assertEquals(['test-project:1'], result.blockedOn)
    self.assertEquals(['other:2', '-test-project:3'], result.blocking)
    self.assertEquals(amendment_mergedinto.newvalue, result.mergedInto)

  def testConvertComment(self):
    """Test convert_comment."""

    svcs = service_manager.Services()
    svcs.user = fake.UserService()
    svcs.user.TestAddUser('user@example.com', 1)
    mar = mock.Mock()
    mar.cnxn = None
    issue = fake.MakeTestIssue(project_id=12345, local_id=1, summary='sum',
                               status='New', owner_id=1001)

    comment = tracker_pb2.IssueComment(
        user_id=1,
        content='test content',
        sequence=1,
        deleted_by=1,
        timestamp=1437700000,
    )
    result = api_pb2_v1_helpers.convert_comment(
        issue, comment, mar, svcs, None)
    self.assertEquals('user@example.com', result.author.name)
    self.assertEquals(comment.content, result.content)
    self.assertEquals('user@example.com', result.deletedBy.name)
    self.assertEquals(1, result.id)
    # Ensure that the published timestamp falls in a timestamp range to account
    # for the test being run in different timezones.
    # Using "Fri, 23 Jul 2015 00:00:00" and "Fri, 25 Jul 2015 00:00:00".
    self.assertTrue(
        datetime.datetime(2015, 7, 23, 0, 0, 0) <= result.published <=
        datetime.datetime(2015, 7, 25, 0, 0, 0))

  def testGetUserEmail(self):
    svcs = service_manager.Services()
    svcs.user = fake.UserService()
    svcs.user.TestAddUser('user@example.com', 1)

    email = api_pb2_v1_helpers._get_user_email(svcs.user, '', 1)
    self.assertEquals('user@example.com', email)

    no_email = api_pb2_v1_helpers._get_user_email(svcs.user, '', 2)
    self.assertEquals(framework_constants.DELETED_USER_NAME, no_email)

  def testSplitRemoveAdd(self):
    """Test split_remove_add."""

    items = ['1', '-2', '-3', '4']
    list_to_add, list_to_remove = api_pb2_v1_helpers.split_remove_add(items)

    self.assertEquals(['1', '4'], list_to_add)
    self.assertEquals(['2', '3'], list_to_remove)

  def testIssueGlobalIDs(self):
    """Test issue_global_ids."""

    svcs = service_manager.Services()
    svcs.issue = fake.IssueService()
    issue1 = fake.MakeTestIssue(12345, 1, 'one', 'New', 111L)
    svcs.issue.TestAddIssue(issue1)
    svcs.project = fake.ProjectService()
    svcs.project.TestAddProject(
        'test-project', owner_ids=[2],
        project_id=12345)
    mar = mock.Mock()
    mar.cnxn = None
    mar.project_name = 'test-project'
    mar.project_id = 12345
    pairs = ['test-project:1']
    result = api_pb2_v1_helpers.issue_global_ids(pairs, 12345, mar, svcs)
    self.assertEquals(100001, result[0])

  def testConvertGroupSettings(self):
    """Test convert_group_settings."""

    setting = usergroup_pb2.MakeSettings('owners', 'mdb', 0)
    result = api_pb2_v1_helpers.convert_group_settings('test-group', setting)
    self.assertEquals('test-group', result.groupName)
    self.assertEquals(
        setting.who_can_view_members, result.who_can_view_members)
    self.assertEquals(setting.ext_group_type, result.ext_group_type)
    self.assertEquals(setting.last_sync_time, result.last_sync_time)
