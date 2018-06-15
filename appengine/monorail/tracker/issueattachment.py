# Copyright 2016 The Chromium Authors. All rights reserved.
# Use of this source code is govered by a BSD-style
# license that can be found in the LICENSE file or at
# https://developers.google.com/open-source/licenses/bsd

"""Issue Tracker code to serve out issue attachments.

Summary of page classes:
  AttachmentPage: Serve the content of an attachment w/ the appropriate
                  MIME type.
  IssueAttachmentDeletion: Form handler for deleting attachments.
"""

import base64
import logging
import os
import re
import urllib

import webapp2

from google.appengine.api import app_identity
from google.appengine.api import images

from framework import exceptions
from framework import framework_constants
from framework import framework_helpers
from framework import gcs_helpers
from framework import permissions
from framework import servlet
from framework import urls
from tracker import attachment_helpers
from tracker import tracker_helpers
from tracker import tracker_views


# This will likely appear blank or as a broken image icon in the browser.
NO_PREVIEW_ICON = ''
NO_PREVIEW_MIME_TYPE = 'image/png'


class AttachmentPage(servlet.Servlet):
  """AttachmentPage serves issue attachments."""

  def GatherPageData(self, mr):
    """Parse the attachment ID from the request and serve its content.

    Args:
      mr: commonly used info parsed from the request.

    Returns: dict of values used by EZT for rendering the page.
    """
    if mr.signed_aid != attachment_helpers.SignAttachmentID(mr.aid):
      webapp2.abort(400, 'Please reload the issue page')

    try:
      attachment, _issue = tracker_helpers.GetAttachmentIfAllowed(
          mr, self.services)
    except exceptions.NoSuchIssueException:
      webapp2.abort(404, 'issue not found')
    except exceptions.NoSuchAttachmentException:
      webapp2.abort(404, 'attachment not found')
    except exceptions.NoSuchCommentException:
      webapp2.abort(404, 'comment not found')

    if not attachment.gcs_object_id:
      webapp2.abort(404, 'attachment data not found')

    bucket_name = app_identity.get_default_gcs_bucket_name()

    gcs_object_id = attachment.gcs_object_id

    logging.info('attachment id %d is %s', mr.aid, gcs_object_id)

    # By default GCS will return images and attachments displayable inline.
    if mr.thumb:
      # Thumbnails are stored in a separate obj always displayed inline.
      gcs_object_id = gcs_object_id + '-thumbnail'
    elif not mr.inline:
      # Downloads are stored in a separate obj with disposiiton set.
      filename = attachment.filename
      if not framework_constants.FILENAME_RE.match(filename):
        logging.info('bad file name: %s' % attachment.attachment_id)
        filename = 'attachment-%d.dat' % attachment.attachment_id
      if gcs_helpers.MaybeCreateDownload(
          bucket_name, gcs_object_id, filename):
        gcs_object_id = gcs_object_id + '-download'

    url = gcs_helpers.SignUrl(bucket_name, gcs_object_id)
    self.redirect(url, abort=True)


class IssueAttachmentDeletion(servlet.Servlet):
  """Form handler that allows user to hard-delete attachments."""

  def ProcessFormData(self, mr, post_data):
    """Process the form that soft-deletes an issue attachment.

    Args:
      mr: commonly used info parsed from the request.
      post_data: HTML form data from the request.

    Returns:
      String URL to redirect the user to after processing.
    """
    local_id = int(post_data['id'])
    sequence_num = int(post_data['sequence_num'])
    attachment_id = int(post_data['aid'])
    delete = 'delete' in post_data

    issue = self.services.issue.GetIssueByLocalID(
        mr.cnxn, mr.project_id, local_id)

    all_comments = self.services.issue.GetCommentsForIssue(
        mr.cnxn, issue.issue_id)
    logging.info('comments on %s are: %s', local_id, all_comments)
    comment = all_comments[sequence_num]

    if not permissions.CanDelete(
        mr.auth.user_id, mr.auth.effective_ids, mr.perms,
        comment.deleted_by, comment.user_id, mr.project,
        permissions.GetRestrictions(issue)):
      raise permissions.PermissionException(
          'Cannot un/delete attachment')

    self.services.issue.SoftDeleteAttachment(
        mr.cnxn, mr.project_id, local_id, sequence_num,
        attachment_id, self.services.user, delete=delete)

    return framework_helpers.FormatAbsoluteURL(
        mr, urls.ISSUE_DETAIL, id=local_id)
