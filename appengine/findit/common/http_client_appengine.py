# Copyright 2014 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

import json
import logging

from google.appengine.api import urlfetch

from common import auth_util
from common.retry_http_client import RetryHttpClient


#TODO(katesonia): Move this to config.
_INTERNAL_HOSTS = ['https://chrome-internal.googlesource.com/']


class HttpClientAppengine(RetryHttpClient):  # pragma: no cover
  """A http client for running on appengine."""

  def _ExpandAuthorizationHeaders(self, headers):
    headers['Authorization'] = 'Bearer ' + auth_util.GetAuthToken()

  def _ShouldLogError(self, status_code):
    if not self.no_error_logging_statuses:
      return True
    return status_code not in self.no_error_logging_statuses

  def _SendRequest(self, url, method, data, timeout, headers=None):
    # We wanted to validate certificate to avoid the man in the middle.
    if not headers:
      headers = {}

    # For google internal hosts, expand Oauth2.0 token to headers to authorize
    # the requests.
    for host in _INTERNAL_HOSTS:
      if url.startswith(host):
        self._ExpandAuthorizationHeaders(headers)
        break

    if method in (urlfetch.POST, urlfetch.PUT):
      result = urlfetch.fetch(
          url, payload=data, method=method,
          headers=headers, deadline=timeout, validate_certificate=True)
    else:
      result = urlfetch.fetch(
          url, headers=headers, deadline=timeout, validate_certificate=True)

    if (result.status_code != 200 and self._ShouldLogError(result.status_code)):
      logging.error('Request to %s resulted in %d, headers:%s', url,
                    result.status_code, json.dumps(result.headers.items()))

    return result.status_code, result.content

  def _Get(self, url, timeout, headers):
    return self._SendRequest(url, urlfetch.GET, None, timeout, headers)

  def _Post(self, url, data, timeout, headers):
    return self._SendRequest(url, urlfetch.POST, data, timeout, headers)

  def _Put(self, url, data, timeout, headers):
    return self._SendRequest(url, urlfetch.PUT, data, timeout, headers)
