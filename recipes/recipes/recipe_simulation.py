# Copyright 2016 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.
"""A continuous builder which runs recipe tests."""

from recipe_engine.recipe_api import Property

DEPS = [
  'depot_tools/bot_update',
  'depot_tools/gclient',

  'recipe_engine/buildbucket',
  'recipe_engine/context',
  'recipe_engine/file',
  'recipe_engine/json',
  'recipe_engine/path',
  'recipe_engine/properties',
  'recipe_engine/python',
  'recipe_engine/step',

  'build/luci_config',
  'build/puppet_service_account',
]

PROPERTIES = {
  'project_under_test': Property(
      default='build', kind=str, help='luci-config project to run tests for'),
  'auth_with_account': Property(
      default=None, kind=str,
      help="Try to authenticate with given service account."),
}


def RunSteps(api, project_under_test, auth_with_account):
  if auth_with_account:
    api.luci_config.c.auth_token = api.puppet_service_account.get_access_token(
        auth_with_account)

  safe_project_name = ''.join(
      c if c.isalnum() else '_' for c in project_under_test)
  root_dir = api.path['cache'].join('builder', safe_project_name)
  api.file.ensure_directory('ensure cache dir', root_dir)
  c = api.gclient.make_config()
  soln = c.solutions.add()
  soln.name = project_under_test
  soln.url = api.luci_config.get_project_metadata(
      project_under_test)['repo_url']
  soln.revision = 'HEAD'

  with api.context(cwd=root_dir):
    api.bot_update.ensure_checkout(gclient_config=c)

  # TODO(martiniss): allow recipes.cfg patches to take affect
  # This requires getting the refs.cfg from luci_config, reading the local
  # patched version, etc.
  result = api.luci_config.get_ref_config(
      project_under_test, 'refs/heads/master', 'recipes.cfg')
  cfg_path = result['url'].split(result['revision'])[-1]
  recipes_path = api.json.loads(result['content']).get('recipes_path', '')

  api.step(
      'recipe simulation test', [
          root_dir.join(
              project_under_test,
              *(recipes_path.split('/') + ['recipes.py'])),
          '--package', root_dir.join(
              project_under_test,
              *cfg_path.split('/')),
          'test', 'run',
      ])


def GenTests(api):
  yield (
      api.test('normal') +
      api.buildbucket.ci_build(
          project='infra',
          builder='recipe simulation tester',
          # This hardcodes test repo URL in luci_config/test_api.py
          git_repo='https://repo.repo/build',
      ) +
      api.properties(
          project_under_test='build',
      ) +
      api.luci_config.get_projects(('build',)) +
      api.luci_config.get_ref_config(
          'build', 'refs/heads/master', 'recipes.cfg',
          content='{"recipes_path": "foobar"}',
          found_at_path='infra/config/')
  )

  yield (
      api.test('with_auth') +
      api.buildbucket.ci_build(
          project='infra',
          builder='recipe simulation tester',
          # This hardcodes test repo URL in luci_config/test_api.py
          git_repo='https://repo.repo/build',
      ) +
      api.properties(
          project_under_test='build',
          auth_with_account='build_limited',
      ) +
      api.luci_config.get_projects(('build',)) +
      api.luci_config.get_ref_config(
          'build', 'refs/heads/master', 'recipes.cfg',
          content='{}',
          found_at_path='custom/config/dir/')
  )
