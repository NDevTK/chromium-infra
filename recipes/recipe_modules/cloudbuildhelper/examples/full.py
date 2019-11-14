# Copyright 2019 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

DEPS = [
  'cloudbuildhelper',
  'recipe_engine/json',
]


def RunSteps(api):
  api.cloudbuildhelper.report_version()

  # Updating pins.
  assert api.cloudbuildhelper.update_pins('some/pins.yaml')

  # Building with all args.
  img = api.cloudbuildhelper.build(
      manifest='some/dir/target.yaml',
      canonical_tag='123_456',
      build_id='bid',
      infra='dev',
      labels={'l1': 'v1', 'l2': 'v2'},
      tags=['latest', 'another'],
  )

  expected = api.cloudbuildhelper.Image(
      image='example.com/fake-registry/target',
      digest='sha256:34a04005bcaf206e...',
      tag='123_456',
      view_image_url=None,
      view_build_url=None)
  assert img == expected, img

  # With minimal args and custom emulated output.
  custom = api.cloudbuildhelper.Image(
      image='xxx',
      digest='yyy',
      tag=None,
      view_image_url=None,
      view_build_url=None)
  img = api.cloudbuildhelper.build('another.yaml', step_test_image=custom)
  assert img == custom, img

  # Using non-canonical tag.
  api.cloudbuildhelper.build('a.yaml', tags=['something'])

  # Using :inputs-hash canonical tag.
  api.cloudbuildhelper.build('b.yaml', canonical_tag=':inputs-hash')

  # Use custom binary from this point onward, for test coverage.
  api.cloudbuildhelper.command = 'custom_cloudbuildhelper'

  # Image that wasn't uploaded anywhere.
  img = api.cloudbuildhelper.build(
      'third.yaml', step_test_image=api.cloudbuildhelper.NotUploadedImage)
  assert img == api.cloudbuildhelper.NotUploadedImage, img

  # Possibly failing build.
  api.cloudbuildhelper.build('fail_maybe.yaml')


def GenTests(api):
  yield api.test('simple')

  yield (
      api.test('failing') +
      api.step_data(
          'cloudbuildhelper build fail_maybe',
          api.cloudbuildhelper.build_error_output('Boom'),
          retcode=1)
  )
