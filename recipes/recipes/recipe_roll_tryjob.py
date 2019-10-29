# Copyright 2019 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

from recipe_engine.recipe_api import Property


DEPS = [
  'recipe_engine/buildbucket',
  'recipe_engine/context',
  'recipe_engine/file',
  'recipe_engine/json',
  'recipe_engine/path',
  'recipe_engine/properties',
  'recipe_engine/python',
  'recipe_engine/raw_io',
  'recipe_engine/runtime',
  'recipe_engine/step',

  'depot_tools/bot_update',
  'depot_tools/gclient',
  'depot_tools/git',
  'depot_tools/tryserver',
]


PROPERTIES = {
  'upstream_id': Property(
      kind=str,
      help='ID of the project to patch'),
  'upstream_url': Property(
      kind=str,
      help='URL of git repo of the upstream project'),

  'downstream_id': Property(
      kind=str,
      help=('ID of the project that includes |upstream_id| in its recipes.cfg '
            'to be tested with upstream patch')),
  'downstream_url': Property(
      kind=str,
      help='URL of the git repo of the downstream project'),
}


NONTRIVIAL_ROLL_FOOTER = 'Recipe-Nontrivial-Roll'
MANUAL_CHANGE_FOOTER = 'Recipe-Manual-Change'
BYPASS_FOOTER = 'Recipe-Tryjob-Bypass-Reason'
KNOWN_FOOTERS = [NONTRIVIAL_ROLL_FOOTER, MANUAL_CHANGE_FOOTER, BYPASS_FOOTER]

FOOTER_ADD_TEMPLATE = '''

Add

    {footer}: {down_id}

To your CL message.

'''

MANUAL_CHANGE_MSG = '''
This means that your upstream CL (this one) will require MANUAL CODE CHANGES
in the downstream repo {down_id!r}. Best practice is to prepare all downstream
changes before landing the upstream CL, using:

    {down_id}/{down_recipes} -O {up_id}=/path/to/local/{up_id} test train

When that CL has been reviewed, you can land this upstream change. Once the
upstream change lands, roll it into your downstream CL:

    {down_id}/recipes.py manual_roll   # may require running multiple times.

Re-train expectations and upload the expectations plus the roll to your
downstream CL. It's customary to copy the outputs of manual_roll to create
a changelog to attach to the downstream CL as well to help reviewers understand
what the roll contains.
'''.strip()

NONTRIVIAL_CHANGE_MSG = '''
This means that your upstream CL (this one) will change the EXPECTATION FILES
in the downstream repo {down_id!r}.

The recipe roller will automatically prepare the non-trivial CL and will upload
it with `git cl upload --r-owners` to the downstream repo. Best practice is to
review this non-trivial roll CL to ensure that the expectations you see there
are expected.
'''

EXTRA_MSG = {
  NONTRIVIAL_ROLL_FOOTER: NONTRIVIAL_CHANGE_MSG,
  MANUAL_CHANGE_FOOTER: MANUAL_CHANGE_MSG,
}


class RecipeTrainingFailure(Exception):
  pass


class RecipesRepo(object):
  """An abstraction of a recipes project to encapsulate common interactions."""

  def __init__(self, api, workdir_base, name, url, manifest_name):
    """
    Args:
      api (RecipeApi): The recipe api for this build.
      workdir_base (Path): The global directory for all recipe repo checkouts.
      name (str): See `name` property.
      url (str): The remote URL for this repo.
      manifest_name (str): The name of the manifest to upload to Logdog (must
        be unique per-build).
    """
    self._api = api
    self._workdir = workdir_base.join(name)
    self._name = name
    self._url = url
    self._manifest_name = manifest_name

    self._root = None
    self._cl_revision = None
    self._recipes_py = None

  @property
  def name(self):
    """The name of this recipes project, e.g. 'recipe_engine'."""
    return self._name

  @property
  def root(self):
    """The absolute path to the root of the checkout for this repo.

    Will be None until `init_checkout()` is called.
    """
    return self._root

  @property
  def recipes_py(self):
    """The path to the recipes.py file for this repo, relative to the root."""
    if self._recipes_py is None:
      recipes_cfg = self._api.file.read_json(
          'parse recipes.cfg',
          self._root.join('infra', 'config', 'recipes.cfg'),
          test_data={
            'recipes_path': 'some/path',
          })
      self._recipes_py = self._api.path.join(
          recipes_cfg.get('recipes_path', ''), 'recipes.py')
    return self._recipes_py

  def checkout_cl(self):
    """Sync the repo the CL that triggered this build.

    Assumes this repo is the repo for the CL.
    """
    assert self._cl_revision
    return self._checkout(self._cl_revision, 'sync %s to CL' % self.name)

  def checkout_master(self):
    """Sync the repo to master."""
    return self._checkout('refs/heads/master', 'sync %s to master' % self.name)

  def _checkout(self, checkout_ref, step_name):
    with self._api.context(cwd=self.root):
      # Clean out those stale pyc's!
      self._api.git('clean', '-xf')
      return self._api.git('checkout', '-f', checkout_ref, name=step_name)

  def init_checkout(self):
    """Checks out the repo into a subdirectory of _workdir.

    Sets `root` to the root of the checkout, and `_cl_revision` to the
    """
    assert not self._root, 'checkout already initialized'

    is_triggering_repo = self._url == self._api.tryserver.gerrit_change_repo_url

    self._api.file.ensure_directory(
        '%s checkout' % self._name, self._workdir)

    gclient_config = self._api.gclient.make_config()
    gclient_config.got_revision_reverse_mapping['got_revision'] = self._name
    soln = gclient_config.solutions.add()
    soln.name = self._name
    soln.url = self._url

    with self._api.context(cwd=self._workdir):
      ret = self._api.bot_update.ensure_checkout(
          gclient_config=gclient_config,
          # Only try to checkout the CL if this repo is the one that triggered
          # the current build.
          patch=is_triggering_repo,
          manifest_name=self._manifest_name)
      self._root = self._workdir.join(ret.json.output['root'])

      if is_triggering_repo:
        with self._api.context(cwd=self._root):
          rev_parse_step = self._api.git(
              'rev-parse', 'HEAD',
              name='read CL revision',
              stdout=self._api.raw_io.output(),
              step_test_data=lambda:
                  self._api.raw_io.test_api.stream_output('deadbeef'))
          self._cl_revision = rev_parse_step.stdout.strip()

  def is_dirty(self, name):
    """Check whether the repo has any unstaged changes.

    Specifically, this can be used after calling `train()` to determine
    whether the training caused a change in the expectation files.
    """
    with self._api.context(cwd=self._root):
      # This has the benefit of showing the expectation diff to the user.
      diff_step = self._api.git('diff', '--exit-code', name=name, ok_ret='any')
      dirty = diff_step.retcode != 0
      if dirty:
        diff_step.presentation.status = 'FAILURE'
      return dirty

  def train(self, upstream_repo, step_name):
    """Re-trains the expectation files for this repo.

    Args:
      upstream_repo (RecipeRepo): A locally checked-out recipes repo that's
        among the dependencies of `self`. The training will be run using the
        local version of the dependency rather than the version pinned in
        recipes.cfg.
      step_name (str): The name to use for the training step.

    Raises:
      `RecipeTrainingFailure` if the training produces an uncaught exception.
    """
    try:
      return self._api.python(
          step_name, self._root.join(self.recipes_py),
          ['-O', '%s=%s' % (upstream_repo.name, upstream_repo.root),
          'test', 'train', '--no-docs'])
    except self._api.step.StepFailure:
      raise RecipeTrainingFailure('failed to train recipes')


def _find_footer(api, repo_id):
  all_footers = api.tryserver.get_footers()

  if BYPASS_FOOTER in all_footers:
    api.python.succeeding_step(
        'BYPASS ENABLED',
        'Roll tryjob bypassed for %r' % (
          # It's unlikely that there's more than one value, but just in case...
          ', '.join(all_footers[BYPASS_FOOTER]),))
    return None, True

  found_set = set()
  for footer in KNOWN_FOOTERS:
    values = all_footers.get(footer, ())
    if repo_id in values:
      found_set.add(footer)

  if len(found_set) > 1:
    api.python.failing_step(
        'Too many footers for %r' % (repo_id,),
        'Found too many footers in CL message:\n' + (
          '\n'.join(' * '+f for f in sorted(found_set)))
    )

  return found_set.pop() if found_set else None, False


def _get_expected_footer(api, upstream_repo, downstream_repo):
  # Run a 'train' on the downstream repo, first using the upstream repo at the
  # CL revision and then at HEAD. We compare these two runs to avoid taking
  # unrolled CLs into account in the resulting diff.
  #
  # If the CL train fails, we require a Manual-Change footer
  # If there's a diff between the two trains, we require a Nontrivial-Roll
  # footer
  # If there's no diff between the two trains, we require no footers

  with api.step.nest('initialize checkouts'):
    upstream_repo.init_checkout()
    downstream_repo.init_checkout()

  upstream_repo.checkout_cl()

  try:
    downstream_repo.train(upstream_repo, 'train recipes at upstream CL')
  except RecipeTrainingFailure:
    return MANUAL_CHANGE_FOOTER

  # If there is no diff introduced by training against this CL along with all
  # the other unrolled changes, then there is a trivial CL that the roller will
  # use rather than trying to roll in any nontrivial ancestor commits. So this
  # change can be considered trivial.
  trivial = not downstream_repo.is_dirty('post-train diff at upstream CL')
  if trivial:
    return None

  with api.context(cwd=downstream_repo.root):
    api.git('add', '--all', name='save post-train diff')

  upstream_repo.checkout_master()

  try:
    downstream_repo.train(upstream_repo, 'train recipes at upstream master')
  except RecipeTrainingFailure:
    # If this CL fixes an older unrolled CL that requires a manual change
    # downstream for it to be rolled, then we can let this CL in without
    # requiring a manual change footer.
    return NONTRIVIAL_ROLL_FOOTER

  # If the training passes when run against upstream at HEAD, we can use the
  # diff between training against upstream HEAD and training against the
  # upstream CL to more accurately determine if this CL introduces diffs,
  # while disregarding diffs introduced by older unrolled changes.
  trivial = not downstream_repo.is_dirty('post-train diff at upstream master')
  return NONTRIVIAL_ROLL_FOOTER if not trivial else None


def RunSteps(api, upstream_id, upstream_url, downstream_id, downstream_url):
  # NOTE: this recipe is only useful as a tryjob with patch applied against
  # upstream repo, which means upstream_url must always match that specified in
  # api.buildbucket.build.input.gerrit_changes[0]. upstream_url remains as a
  # required parameter for symmetric input for upstream/downstream.
  # TODO: figure out upstream_id from downstream's repo recipes.cfg file using
  # patch and deprecated both upstream_id and upstream_url parameters.
  workdir_base = api.path['cache'].join('builder')

  upstream_repo = RecipesRepo(
    api, workdir_base, upstream_id, upstream_url, manifest_name='upstream')
  downstream_repo = RecipesRepo(
    api, workdir_base, downstream_id, downstream_url,
    manifest_name='downstream')

  # First, check to see if the user has bypassed this tryjob's analysis
  # entirely.
  actual_footer, bypass = _find_footer(api, downstream_id)
  if bypass:
    return

  expected_footer = _get_expected_footer(api, upstream_repo, downstream_repo)

  # Either expected_footer and actual_footer are both None or both matching
  # footers.
  if expected_footer == actual_footer:
    if expected_footer:
      msg = (
        'CL message contains correct footer (%r) for this repo.'
      ) % expected_footer
    else:
      msg = 'CL is trivial and message contains no footers for this repo.'
    api.python.succeeding_step('Roll OK', msg)
    return

  # trivial roll, but user has footer in CL message.
  if expected_footer is None and actual_footer is not None:
    api.python.failing_step(
        'UNEXPECTED FOOTER IN CL MESSAGE',
        'Change is trivial, but found %r footer' % (actual_footer,))

  # nontrivial/manual roll, but user has wrong footer in CL message.
  if expected_footer is not None and actual_footer is not None:
    api.python.failing_step(
        'WRONG FOOTER IN CL MESSAGE',
        'Change reqires %r, but found %r footer' % (
          expected_footer, actual_footer,))

  # expected != None at this point, so actual_footer must be None
  msg = FOOTER_ADD_TEMPLATE + EXTRA_MSG[expected_footer]
  api.python.failing_step(
      'MISSING FOOTER IN CL MESSAGE',
      msg.format(
          footer=expected_footer,
          up_id=upstream_id,
          down_id=downstream_id,
          down_recipes=downstream_repo.recipes_py,
      ))


def GenTests(api):
  def test(name, *footers):
    upstream_id = 'recipe_engine'
    downstream_id = 'depot_tools'
    repo_urls = {
      'build':
      'https://chromium.googlesource.com/chromium/tools/build',
      'depot_tools':
      'https://chromium.googlesource.com/chromium/tools/depot_tools',
      'recipe_engine':
      'https://chromium.googlesource.com/infra/luci/recipes-py',
    }
    return (
      api.test(name)
      + api.runtime(is_luci=True, is_experimental=False)
      + api.properties(
          upstream_id=upstream_id,
          upstream_url=repo_urls[upstream_id],
          downstream_id=downstream_id,
          downstream_url=repo_urls[downstream_id])
      + api.buildbucket.try_build(
          git_repo=repo_urls[upstream_id],
          change_number=456789,
          patch_set=12)
      + api.override_step_data('gerrit changes', api.json.output([{
        'revisions': {
          'deadbeef': {'_number': 12, 'commit': {'message': ''}},
        }
      }]))
      + api.step_data(
          'parse description', api.json.output({
            k: ['Reasons' if k == BYPASS_FOOTER else downstream_id]
            for k in footers
          }))
    )

  yield (
    test('find_trivial_roll')
  )

  yield (
    test('bypass', BYPASS_FOOTER)
    + api.post_check(lambda check, steps: check('BYPASS ENABLED' in steps))
  )

  yield (
    test('too_many_footers', MANUAL_CHANGE_FOOTER, NONTRIVIAL_ROLL_FOOTER)
    + api.post_check(lambda check, steps: check(
        "Too many footers for 'depot_tools'" in steps
    ))
  )

  yield (
    test('find_trivial_roll_unexpected', MANUAL_CHANGE_FOOTER)
    + api.post_check(lambda check, steps: check(
        'UNEXPECTED FOOTER IN CL MESSAGE' in steps
    ))
  )

  yield (
    test('find_manual_roll_missing')
    + api.step_data('train recipes at upstream CL', retcode=1)
    + api.post_check(lambda check, steps: check(
        MANUAL_CHANGE_FOOTER in steps['MISSING FOOTER IN CL MESSAGE'].step_text
    ))
  )

  yield (
    test('find_manual_roll_wrong', NONTRIVIAL_ROLL_FOOTER)
    + api.step_data('train recipes at upstream CL', retcode=1)
    + api.post_check(lambda check, steps: check(
        MANUAL_CHANGE_FOOTER in steps['WRONG FOOTER IN CL MESSAGE'].step_text
    ))
  )

  yield (
    test('find_non_trivial_roll')
    + api.step_data('post-train diff at upstream CL', retcode=1)
    + api.step_data('post-train diff at upstream master', retcode=1)
    + api.post_check(lambda check, steps: check(
      NONTRIVIAL_ROLL_FOOTER in steps['MISSING FOOTER IN CL MESSAGE'].step_text
    ))
  )

  yield (
    test('trivial_roll_unrolled_changes')
    + api.step_data('post-train diff at upstream CL', retcode=1)
  )

  yield (
    test('nontrivial_roll_unrolled_changes')
    + api.step_data('post-train diff at upstream CL', retcode=1)
    + api.step_data('post-train diff at upstream master', retcode=1)
    + api.post_check(lambda check, steps: check(
      NONTRIVIAL_ROLL_FOOTER in steps['MISSING FOOTER IN CL MESSAGE'].step_text
    ))
  )

  yield (
    test('nontrivial_roll_match', NONTRIVIAL_ROLL_FOOTER)
    + api.step_data('post-train diff at upstream CL', retcode=1)
    + api.step_data('post-train diff at upstream master', retcode=1)
  )

  # The upstream HEAD commit cannot be rolled downstream without a manual
  # change, but this nontrivial CL fixes that.
  yield (
    test('nontrivial_fix_for_broken_head', NONTRIVIAL_ROLL_FOOTER)
    + api.step_data('post-train diff at upstream CL', retcode=1)
    + api.step_data('train recipes at upstream master', retcode=1)
  )
