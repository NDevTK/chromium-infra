# Copyright 2014 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

import re

from common import dependency
from common import deps_parser
from common import git_repository
from common import http_client_appengine


_CHROMIUM_ROOT_DIR = 'src/'
_CHROMIUM_REPO_MASTER = 'https://chromium.googlesource.com/chromium/src.git'

_CHROME_VERSION_PATTERN = re.compile(r'^\d+\.\d+\.\d+\.\d+$')

_BUILDSPEC_REPO = ('https://chrome-internal.googlesource.com/chrome/tools/'
                   'buildspec.git/')


def IsChromeVersion(revision):
  """Determines if a revision is a chrome version."""
  if _CHROME_VERSION_PATTERN.match(revision):
    return True

  return False


class DEPSDownloader(deps_parser.DEPSLoader):
  """Downloads DEPS from remote Git repo."""

  def __init__(self, check_deps_git_first=False):
    """
    Args:
      check_deps_git_first (bool): If True, use .DEPS.git instead of DEPS.
    """
    self.check_deps_git_first = check_deps_git_first

  def Load(self, repo_url, revision, deps_file):
    http_client = http_client_appengine.HttpClientAppengine()
    repo = git_repository.GitRepository(repo_url, http_client)

    content = None

    if self.check_deps_git_first and deps_file == 'DEPS':
      # When the given deps_file is "DEPS" and .DEPS.git should be checked
      # first, it's because before migration from SVN to Git, .DEPS.git contains
      # dependencies hosted in Git while DEPS contains those in SVN.
      # If .DEPS.git is not found, fallback to the given deps_file. Assume it is
      # a commit after migration from SVN to Git.
      content = repo.GetSource('.DEPS.git', revision)

    if content is None:
      content = repo.GetSource(deps_file, revision)

    if content is None:
      raise Exception(
          'Failed to pull %s file from %s, at revision %s.' % (
              deps_file, repo_url, revision))

    return content


def GetChromeDependency(revision, platform, check_deps_git_first=False):
  """Returns all dependencies of Chrome as a dict for the given revision and OS.

  Args:
    revision (str): The revision of a Chrome build, it can be a githash or a
      chrome version for a official build.
    platform (str): The target platform of the Chrome build, should be one of
      'win', 'ios', 'mac', 'unix', 'android', or 'all'.
    check_deps_git_first (bool): If True, use .DEPS.git instead of DEPS.

  Returns:
    A map from dependency path to the dependency info.
  """
  deps_repo_info = {'deps_file': 'DEPS'}

  if IsChromeVersion(revision):
    # For chrome version, get the DEPS file from internal buildspec/ repo
    # instead of chromium trunk.
    deps_repo_info['deps_repo_url'] = _BUILDSPEC_REPO
    deps_repo_info['deps_repo_revision'] = 'master'
    deps_repo_info['deps_file'] = 'release/%s/DEPS' % revision

  root_dep = dependency.Dependency(
      _CHROMIUM_ROOT_DIR, _CHROMIUM_REPO_MASTER, revision, **deps_repo_info)

  deps_parser.UpdateDependencyTree(
      root_dep, [platform], DEPSDownloader(check_deps_git_first))

  dependencies = {}

  # Flatten the dependency tree into a one-level dict.
  def FlattenDepTree(dep):
    dependencies[dep.path] = dep
    for child in dep.children.values():
      FlattenDepTree(child)

  FlattenDepTree(root_dep)

  return dependencies


def GetChromiumDEPSRolls(old_cr_revision, new_cr_revision, platform,
                         check_deps_git_first=False, skip_chromium_roll=True):
  """Returns a list of dependency rolls between the given Chromium revisions.

  Args:
    old_cr_revision (str): The old Chromium revision, it can be a githash or a
      chrome version for a official build.
    new_cr_revision (str): The new Chromium revision, it can be a githash or a
      chrome version for a official build.
    platform (str): The target OS platform of the Chrome or test binary.
    check_deps_git_first (bool): If True, use .DEPS.git instead of DEPS.
    skip_chromium_roll (bool): If False, chromium roll will be contained in
      the return.

  Returns:
    A list of DependencyRoll objects in the revision range.
  """
  old_deps = GetChromeDependency(
      old_cr_revision, platform, check_deps_git_first)
  new_deps = GetChromeDependency(
      new_cr_revision, platform, check_deps_git_first)

  rolls = []

  for path, new_dep in new_deps.iteritems():
    if skip_chromium_roll and path == _CHROMIUM_ROOT_DIR:
      continue

    old_revision = None
    if path in old_deps:
      old_revision = old_deps[path].revision

    if old_revision != new_dep.revision:
      rolls.append(
          dependency.DependencyRoll(
              path, new_dep.repo_url, old_revision, new_dep.revision))

  for path, old_dep in old_deps.iteritems():
    if path not in new_deps:
      rolls.append(
          dependency.DependencyRoll(
              path, old_dep.repo_url, old_dep.revision, None))

  return rolls


def GetDEPSRollsDict(old_cr_revision, new_cr_revision, platform):
  """Gets dep_path to DependencyRoll dictionary for deps between revisions.

  Args:
    old_cr_revision (str): The old Chromium revision, it can be a githash or a
      chrome version for a official build.
    new_cr_revision (str): The new Chromium revision, it can be a githash or a
      chrome version for a official build.
    platform (str): The target OS platform of the Chrome or test binary.

  Returns:
    A dict, mapping dep path to its DependencyRoll.
  """
  deps_rolls = GetChromiumDEPSRolls(old_cr_revision, new_cr_revision, platform,
                                    skip_chromium_roll=False)

  deps_rolls_dict = {}
  for dep_roll in deps_rolls:
    deps_rolls_dict[dep_roll.path] = dep_roll

  return deps_rolls_dict
