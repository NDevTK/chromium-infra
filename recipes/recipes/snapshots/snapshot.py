# Copyright 2018 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

DEPS = [
  'recipe_engine/buildbucket',
  'recipe_engine/cipd',
  'recipe_engine/path',
  'recipe_engine/platform',
  'recipe_engine/properties',
  'recipe_engine/step',
  'recipe_engine/time',
  'recipe_engine/url',
]


def RunSteps(api):
  packages_dir = api.path['start_dir'].join('packages')
  ensure_file = api.cipd.EnsureFile()
  ensure_file.add_package(
      'infra/machine-provider/snapshot/gce/${platform}', 'canary')
  api.cipd.ensure(packages_dir, ensure_file)

  result = api.url.get_text(
      'http://metadata/computeMetadata/v1/instance/zone',
      default_test_data='projects/p/zones/us-zone1-a',
      headers={'Metadata-Flavor': 'Google'},
      step_name='get zone')
  zone = result.output.rsplit('/', 1)[-1]

  snapshot = packages_dir.join('snapshot')
  api.step('snapshot', [
      snapshot,
      '-disk', '%s-disk' % api.properties['bot_id'],
      # Name doesn't matter since the snapshot will be identified by labels.
      # TODO(smut): Add labels.
      '-name', 'snapshot-%s' % int(api.time.time() * 100),
      '-project', 'google.com:chromecompute',
      '-service-account-json', ':gce',
      '-zone', zone,
  ])


def GenTests(api):
  yield (
    api.test('snapshot') +
    api.platform('linux', 64) +
    api.properties(bot_id='bot-id') +
    api.buildbucket.ci_build()
  )
