# Copyright 2016 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

import textwrap

from testing_utils import testing

from infra_api_clients.codereview.rietveld import Rietveld
from libs.http import retry_http_client


class DummyHttpClient(retry_http_client.RetryHttpClient):

  def __init__(self):
    super(DummyHttpClient, self).__init__()
    self.responses = {}
    self.requests = []

  def SetResponse(self, url, result):
    self.responses[url] = result

  def GetBackoff(self, *_):  # pragma: no cover
    """Override to avoid sleep."""
    return 0

  def _Get(self, *_):  # pragma: no cover
    pass

  def _Post(self, url, data, _, headers):
    self.requests.append((url, data, headers))
    return self.responses.get(url, (404, 'Not Found'))

  def _Put(self, *_):  # pragma: no cover
    pass


class RietveldTest(testing.AppengineTestCase):

  def setUp(self):
    super(RietveldTest, self).setUp()
    self.http_client = DummyHttpClient()
    self.rietveld = Rietveld()
    self.rietveld.HTTP_CLIENT = self.http_client

  def testGetXsrfTokenSuccess(self):
    rietveld_url = 'https://test'
    self.http_client.SetResponse('%s/xsrf_token' % rietveld_url, (200, 'abc'))
    self.assertEqual('abc', self.rietveld._GetXsrfToken(rietveld_url))
    self.assertEqual(1, len(self.http_client.requests))
    _, _, headers = self.http_client.requests[0]
    self.assertTrue('X-Requesting-XSRF-Token' in headers)

  def testGetXsrfTokenFailure(self):
    rietveld_url = 'https://test'
    self.http_client.SetResponse('%s/xsrf_token' % rietveld_url, (302, 'login'))
    self.assertIsNone(self.rietveld._GetXsrfToken(rietveld_url))

  def testGetRietveldUrlAndIssueNumber(self):
    cases = {
        'http://abc/123': '123',
        'https://abc/456/': '456',
        'http://abc/789/diff': '789',
    }
    for issue_url, expected_issue_number in cases.iteritems():
      rietveld_url, issue_number = self.rietveld._GetRietveldUrlAndIssueNumber(
          issue_url)
      self.assertEqual('https://abc', rietveld_url)
      self.assertEqual(expected_issue_number, issue_number)

  def testEncodeMultipartFormData(self):
    content_type, body = self.rietveld._EncodeMultipartFormData({'a':'b'})
    expected_content_type = (
        'multipart/form-data; boundary=-F-I-N-D-I-T-M-E-S-S-A-G-E-')
    expected_body = textwrap.dedent("""
    ---F-I-N-D-I-T-M-E-S-S-A-G-E-\r
    Content-Disposition: form-data; name="a"\r
    \r
    b\r
    ---F-I-N-D-I-T-M-E-S-S-A-G-E---\r
    """)[1:]
    self.assertEqual(expected_content_type, content_type)
    self.assertEqual(expected_body, body)

  def testPostMessageSuccess(self):
    rietveld_url = 'https://test'
    issue_url = '%s/123' % rietveld_url
    message_publish_url = '%s/publish' % issue_url
    self.http_client.SetResponse('%s/xsrf_token' % rietveld_url, (200, 'abc'))
    self.http_client.SetResponse(message_publish_url, (200, 'OK'))
    self.assertTrue(self.rietveld.PostMessage(issue_url, 'message'))
    self.assertEqual(2, len(self.http_client.requests))

  def testPostMessageFailOnXsrfToken(self):
    rietveld_url = 'https://test'
    issue_url = '%s/123' % rietveld_url
    message_publish_url = '%s/publish' % issue_url
    self.http_client.SetResponse('%s/xsrf_token' % rietveld_url, (302, 'login'))
    self.http_client.SetResponse(message_publish_url, (200, 'OK'))
    self.assertFalse(self.rietveld.PostMessage(issue_url, 'message'))
    self.assertEqual(1, len(self.http_client.requests))

  def testPostMessageFailOnPublish(self):
    rietveld_url = 'https://test'
    issue_url = '%s/123' % rietveld_url
    message_publish_url = '%s/publish' % issue_url
    self.http_client.SetResponse('%s/xsrf_token' % rietveld_url, (302, 'login'))
    self.http_client.SetResponse(message_publish_url, (429, 'Error'))
    self.assertFalse(self.rietveld.PostMessage(issue_url, 'message'))
    self.assertEqual(1, len(self.http_client.requests))
