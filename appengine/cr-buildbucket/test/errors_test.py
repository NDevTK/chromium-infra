# Copyright 2014 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

from testing_utils import testing
import errors


class ErrorsTest(testing.AppengineTestCase):
  def test_validate_bucket_name(self):
    with self.assertRaises(errors.InvalidInputError):
      errors.validate_bucket_name(1)
    with self.assertRaises(errors.InvalidInputError):
      errors.validate_bucket_name('luci.x', project_id='y')
    with self.assertRaises(errors.InvalidInputError):
      errors.validate_bucket_name('')
    with self.assertRaises(errors.InvalidInputError):
      errors.validate_bucket_name('no spaces')
    with self.assertRaises(errors.InvalidInputError):
      errors.validate_bucket_name('no spaces')
    errors.validate_bucket_name('good-name')
    errors.validate_bucket_name('luci.infra.try', project_id='infra')

  def test_default_message(self):
    ex = errors.BuildIsCompletedError()
    self.assertEqual(ex.message, 'Build is complete and cannot be changed.')
