# Copyright 2014 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

import datetime
import re
import string

from components import utils

BUCKET_NAME_REGEX = re.compile(r'^[0-9a-z_\.\-]{1,100}$')
MAX_LEASE_DURATION = datetime.timedelta(hours=2)

BUILDER_NAME_VALID_CHARS = string.ascii_letters + string.digits + '()-_. '
_BUILDER_NAME_VALID_CHAR_SET = frozenset(BUILDER_NAME_VALID_CHARS)


class Error(Exception):

  def __init__(self, message=''):
    # passing None instead of empty docstring so that
    # Exception class applies its own default.
    super(Error, self).__init__(message or self.__doc__ or None)


class NotFoundError(Error):
  """Requested resource not found."""


class BuildNotFoundError(NotFoundError):
  """Requested build was not found."""


class BuilderNotFoundError(NotFoundError):
  """Requested builder was not found."""


class BuildIsCompletedError(Error):
  """Build is complete and cannot be changed."""


class InvalidInputError(Error):
  """Raised when service method argument value is invalid."""


class LeaseExpiredError(Error):
  """Raised when provided lease_key does not match the current one."""


class TagIndexIncomplete(Error):
  """Raised when a tag index is permanently incomplete and cannot be used."""


# TODO(nodir): move to config.py. Cannot be done because config.py depends on
# swarmingcfg.py which needs this function.
def validate_builder_name(name):
  if not name:
    raise InvalidInputError('unspecified')

  if len(name) > 128:
    raise InvalidInputError('length is > 128')

  invalid_chars = ''.join(
      sorted(set(c for c in name if c not in _BUILDER_NAME_VALID_CHAR_SET))
  )
  if invalid_chars:
    raise InvalidInputError(
        'invalid char(s) %r. Alphabet: "%s"' %
        (invalid_chars, BUILDER_NAME_VALID_CHARS)
    )


# TODO(crbug.com/851036): move to config.py
def validate_bucket_name(bucket, project_id=None):
  """Raises InvalidInputError if bucket name is invalid."""
  if not bucket:
    raise InvalidInputError('Bucket not specified')
  if (project_id and bucket.startswith('luci.') and
      not bucket.startswith('luci.%s.' % project_id)):
    raise InvalidInputError(
        'Bucket must start with "luci.%s." because it starts with "luci." '
        'and is defined in the %s project' % (project_id, project_id)
    )

  if not isinstance(bucket, basestring):
    raise InvalidInputError(
        'Bucket must be a string. It is %s.' % type(bucket).__name__
    )
  if not BUCKET_NAME_REGEX.match(bucket):
    raise InvalidInputError(
        'Bucket name "%s" does not match regular expression %s' %
        (bucket, BUCKET_NAME_REGEX.pattern)
    )


def validate_lease_expiration_date(expiration_date):
  """Raises errors.InvalidInputError if |expiration_date| is invalid."""
  if expiration_date is None:
    return
  if not isinstance(expiration_date, datetime.datetime):
    raise InvalidInputError('Lease expiration date must be datetime.datetime')
  duration = expiration_date - utils.utcnow()
  if duration <= datetime.timedelta(0):
    raise InvalidInputError('Lease expiration date cannot be in the past')
  if duration > MAX_LEASE_DURATION:
    raise InvalidInputError(
        'Lease duration cannot exceed %s' % MAX_LEASE_DURATION
    )
