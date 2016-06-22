#!/usr/bin/env python
# Copyright 2014 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

"""Prepares a local hermetic Go installation.

- Downloads and unpacks the Go toolset in ../../golang.
- Downloads and installs Glide (used by deps.py).
- Fetches code dependencies via deps.py.
"""

import argparse
import contextlib
import json
import logging
import os
import platform
import shutil
import stat
import subprocess
import sys
import tarfile
import tempfile
import urllib
import zipfile


LOGGER = logging.getLogger(__name__)


# /path/to/infra
ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))

# The current overarching Infra version. If this changes, everything will be
# updated regardless of its version.
INFRA_VERSION = 1

# Where to install Go toolset to. GOROOT would be <TOOLSET_ROOT>/go.
TOOLSET_ROOT = os.path.join(os.path.dirname(ROOT), 'golang')

# Where Go GAE SDK is supposedly installed.
GO_APPENGINE = os.path.join(os.path.dirname(ROOT), 'go_appengine')

# Default workspace with infra go code.
WORKSPACE = os.path.join(ROOT, 'go')

# Platform depended suffix for executable files.
EXE_SFX = '.exe' if sys.platform == 'win32' else ''

# On Windows we use git from depot_tools.
GIT_EXE = 'git.bat' if sys.platform == 'win32' else 'git'

# Pinned version of Go toolset to download. See http://golang.org/dl/.
TOOLSET_VERSION = 'go1.6.2'

# Platform and toolset dependent portion of a download URL.
# See NATIVE_DOWNLOAD_URL_PREFIX and CROSS_DOWNLOAD_URL_PREFIX.
TOOLSET_VARIANTS = {
  'darwin-amd64': 'darwin-amd64.tar.gz',
  'linux2-386': 'linux-386.tar.gz',
  'linux2-amd64': 'linux-amd64.tar.gz',
  'win32-386': 'windows-386.zip',
  'win32-amd64': 'windows-amd64.zip',
  # Cross-compilation toolsets must support all native ones above.
  #'linux2-amd64-cross': 'linux2-amd64-cross.tar.gz',
}


# Google Storage utility wrapper in depot_tools.
GSUTIL = os.path.join(os.path.dirname(ROOT), 'depot_tools', 'gsutil.py')
# Google Storage bucket for go cross-compilation toolsets .
BUCKET = 'chrome-infra-public'
# A path inside bucket where toolset archives are stored.
BUCKET_PATH = 'go/toolset'

# Download URL root for official native toolsets.
NATIVE_DOWNLOAD_URL_PREFIX = 'https://storage.googleapis.com/golang'
# Download URL for cross-compiled toolsets.
CROSS_DOWNLOAD_URL_PREFIX = (
    'https://storage-download.googleapis.com/%s/%s' % (BUCKET, BUCKET_PATH))

# Describes how to fetch 'glide'.
GLIDE_SOURCE = {
  'src/github.com/Masterminds/glide': {
    'url': (
        'https://chromium.googlesource.com/external/github.com/'
        'Masterminds/glide.git'),
    'rev': 'refs/tags/0.10.2',
  },
}

# The set of directories to copy from the go_appengine SDK's goroot into our
# toolchain's root. See patch_appengine_sdk for more details.
GAE_PKGS = frozenset(('appengine', 'appengine_internal', 'github.com'))

# Counter-hacks against GAE's hacks :). The init module writes a gae-added
# boolean into the os package to disable file writes, and overwrites the
# http.DefaultTransport to prevent usage. Since we're still using the real
# GAE SDK for running apps locally and doing app uploads, remove their hacks.
GAE_HACK_PKGS = frozenset(('appengine_internal/init',))


class Failure(Exception):
  """Bootstrap failed."""


def get_default_toolset():
  """Name of a toolset native for this platform."""
  machine = platform.machine().lower()
  if (machine == 'x86_64' and platform.system() == 'Linux' and
      sys.maxsize == (2 ** 31) - 1):
    # This is 32bit python on 64bit CPU on linux, which probably means the
    # entire userland is 32bit and thus we should play along and install 32bit
    # packages.
    machine = 'x86'
  # Name arch the same way Go does it.
  arch = {
    'amd64': 'amd64',
    'x86_64': 'amd64',
    'i386': '386',
    'i686': '386',
    'x86': '386',
  }.get(machine)
  return '%s-%s' % (sys.platform, arch)


def is_cross_compile(toolset):
  return toolset.endswith('-cross')


def is_compatible_toolset(installed_toolset, toolset):
  return (installed_toolset == toolset or
      toolset + '-cross' == installed_toolset)


def get_toolset_url(toolset):
  """URL of a specific Go toolset archive."""
  variant = TOOLSET_VARIANTS.get(toolset)
  if not variant:
    # TODO(vadimsh): Compile go lang from source.
    raise Failure('Unrecognized toolset')
  if is_cross_compile(toolset):
    return '%s/%s.%s' % (CROSS_DOWNLOAD_URL_PREFIX, TOOLSET_VERSION, variant)
  else:
    return '%s/%s.%s' % (NATIVE_DOWNLOAD_URL_PREFIX, TOOLSET_VERSION, variant)


def read_file(path):
  """Returns contents of a given file or None if not readable."""
  assert isinstance(path, (list, tuple))
  try:
    with open(os.path.join(*path), 'r') as f:
      return f.read()
  except IOError:
    return None


def write_file(path, data):
  """Writes |data| to a file."""
  assert isinstance(path, (list, tuple))
  with open(os.path.join(*path), 'w') as f:
    f.write(data)


def remove_directory(path):
  """Recursively removes a directory."""
  assert isinstance(path, (list, tuple))
  p = os.path.join(*path)
  if not os.path.exists(p):
    return
  LOGGER.info('Removing %s', p)
  # Crutch to remove read-only file (.git/* in particular) on Windows.
  def onerror(func, path, _exc_info):
    if not os.access(path, os.W_OK):
      os.chmod(path, stat.S_IWUSR)
      func(path)
    else:
      raise
  shutil.rmtree(p, onerror=onerror if sys.platform == 'win32' else None)


def install_toolset(toolset_root, toolset):
  """Downloads and installs Go toolset.

  GOROOT would be <toolset_root>/go/.
  """
  if not os.path.exists(toolset_root):
    os.makedirs(toolset_root)
  url = get_toolset_url(toolset)
  pkg_path = os.path.join(toolset_root, url[url.rfind('/')+1:])

  LOGGER.info('Downloading %s...', url)
  download_file(url, pkg_path)

  LOGGER.info('Extracting...')
  if pkg_path.endswith('.zip'):
    with zipfile.ZipFile(pkg_path, 'r') as f:
      f.extractall(toolset_root)
  elif pkg_path.endswith('.tar.gz'):
    with tarfile.open(pkg_path, 'r:gz') as f:
      f.extractall(toolset_root)
  else:
    raise Failure('Unrecognized archive format')

  LOGGER.info('Validating...')
  if not check_hello_world(toolset_root):
    raise Failure('Something is not right, test program doesn\'t work')


def download_file(url, path):
  """Fetches |url| to |path|."""
  last_progress = [0]
  def report(a, b, c):
    progress = int(a * b * 100.0 / c)
    if progress != last_progress[0]:
      print >> sys.stderr, 'Downloading... %d%%' % progress
      last_progress[0] = progress
  # TODO(vadimsh): Use something less crippled, something that validates SSL.
  urllib.urlretrieve(url, path, reporthook=report)


@contextlib.contextmanager
def temp_dir(path):
  """Creates a temporary directory, then deletes it."""
  tmp = tempfile.mkdtemp(dir=path)
  try:
    yield tmp
  finally:
    remove_directory([tmp])


def check_hello_world(toolset_root):
  """Compiles and runs 'hello world' program to verify that toolset works."""
  with temp_dir(toolset_root) as tmp:
    path = os.path.join(tmp, 'hello.go')
    write_file([path], r"""
        package main
        func main() { println("hello, world\n") }
    """)
    out = subprocess.check_output(
        [get_go_exe(toolset_root), 'run', path],
        env=get_go_environ(toolset_root, tmp),
        stderr=subprocess.STDOUT)
    if out.strip() != 'hello, world':
      LOGGER.error('Failed to run sample program:\n%s', out)
      return False
    return True


def discover_targets(src_dir, pkg):
  ret = []
  for root, _, fnames in os.walk(os.path.join(src_dir, pkg)):
    if any(n.endswith('.go') for n in fnames):
      ret.append(root[len(src_dir)+1:])
  return ret


def is_suitable_toolset_installed(installed, toolset):
  if not installed:
    return False
  try:
    installed_toolset, installed_version = installed.strip().split(' ')
    return (TOOLSET_VERSION == installed_version and
        is_compatible_toolset(installed_toolset, toolset))
  except ValueError:
    return False


def infra_version_outdated(root):
  infra = read_file([root, 'INFRA_VERSION'])
  if not infra:
    return True
  return int(infra.strip()) < INFRA_VERSION


def write_infra_version(root):
  write_file([root, 'INFRA_VERSION'], str(INFRA_VERSION))


def ensure_toolset_installed(toolset_root, toolset=None):
  """Installs or updates Go toolset if necessary.

  Returns True if new toolset was installed.
  """
  if toolset is None:
    toolset = get_default_toolset()
  infra_outdated = infra_version_outdated(toolset_root)
  installed = read_file([toolset_root, 'INSTALLED_TOOLSET'])

  if infra_outdated:
    LOGGER.info('Infra version is out of date.')
  elif is_suitable_toolset_installed(installed, toolset):
    LOGGER.debug('Go toolset is up-to-date: %s', installed)
    return False

  available = '%s %s' % (toolset,  TOOLSET_VERSION)
  LOGGER.info('Installing Go toolset.')
  LOGGER.info('  Old toolset is %s', installed)
  LOGGER.info('  New toolset is %s', available)
  remove_directory([toolset_root])
  install_toolset(toolset_root, toolset)
  LOGGER.info('Go toolset installed: %s', available)
  write_file([toolset_root, 'INSTALLED_TOOLSET'], available)
  write_infra_version(toolset_root)
  return True


def ensure_glide_installed(toolset_root):
  """Installs or updates 'glide' tool."""
  installed_tools = read_file([toolset_root, 'INSTALLED_TOOLS'])
  available_tools = json.dumps(GLIDE_SOURCE, sort_keys=True)
  if installed_tools == available_tools:
    LOGGER.debug('Glide is up-to-date')
    return

  def install(workspace, pkg):
    subprocess.check_call(
        [get_go_exe(toolset_root), 'install', pkg],
        cwd=tmp,
        env=get_go_environ(toolset_root, workspace),
        stdout=sys.stderr)
    # Windows os.rename doesn't support overwrites.
    name = pkg[pkg.rfind('/')+1:]
    dest = os.path.join(toolset_root, 'go', 'bin', name + EXE_SFX)
    if os.path.exists(dest):
      os.remove(dest)
    os.rename(os.path.join(workspace, 'bin', name + EXE_SFX), dest)
    LOGGER.info('Installed %s', dest)

  LOGGER.info('Installing Glide...')
  with temp_dir(toolset_root) as tmp:
    fetch_glide_code(tmp, GLIDE_SOURCE)
    install(tmp, 'github.com/Masterminds/glide')

  LOGGER.info('Glide is installed')
  write_file([toolset_root, 'INSTALLED_TOOLS'], available_tools)


def fetch_glide_code(workspace, spec):
  """Fetches glide source code."""
  def git(cmd, cwd):
    subprocess.check_call([GIT_EXE] + cmd, cwd=cwd, stdout=sys.stderr)
  for path, repo in sorted(spec.iteritems()):
    path = os.path.join(workspace, path.replace('/', os.sep))
    os.makedirs(path)
    git(['clone', repo['url'], '.'], cwd=path)
    git(['checkout', repo['rev']], cwd=path)


def update_vendor_packages(workspace, force=False):
  """Runs deps.py to fetch and install pinned packages."""
  if not os.path.isfile(os.path.join(workspace, 'deps.lock')):
    return
  cmd = [
    sys.executable, '-u', os.path.join(ROOT, 'go', 'deps.py'),
    '--workspace', workspace, 'install',
  ]
  if force:
    cmd.append('--force')
  subprocess.check_call(cmd, stdout=sys.stderr)


def get_go_environ(
    toolset_root,
    workspace=None,
    go_paths=(),
    vendor_paths=(),
    go_appengine_path=None):
  """Returns a copy of os.environ with added GO* environment variables.

  Overrides GOROOT, GOPATH and GOBIN. Keeps everything else. Idempotent.

  Args:
    toolset_root: GOROOT would be <toolset_root>/go.
    workspace: main workspace directory or None if compiling in GOROOT.
    go_paths: additional paths to add to GOPATH (used by bootstrap_internal.py).
    vendor_paths: directories with .vendor directories.
    go_appengine_path: path to GAE Go SDK to add to PATH.
  """
  env = os.environ.copy()
  env['GOROOT'] = os.path.join(toolset_root, 'go')
  if workspace:
    env['GOBIN'] = os.path.join(workspace, 'bin')
  else:
    env.pop('GOBIN', None)

  all_go_paths = [os.path.join(p, '.vendor') for p in vendor_paths]
  all_go_paths.extend(go_paths)
  if workspace:
    all_go_paths.append(workspace)
  env['GOPATH'] = os.pathsep.join(all_go_paths)

  # Remove preexisting bin/ paths (including .vendor/bin) pointing to infra
  # or infra_internal Go workspaces. It's important when switching from
  # infra_internal to infra environments: infra_internal bin paths should
  # be removed.
  path = env['PATH'].split(os.pathsep)
  def should_keep(p):
    # TODO(vadimsh): This code knows about gclient checkout layout.
    gclient_root = os.path.dirname(ROOT)
    for d in ['infra', 'infra_internal']:
      if p.startswith(os.path.join(gclient_root, d, 'go')):
        return False
    return True
  path = filter(should_keep, path)

  # New PATH entries.
  paths_to_add = [
    os.path.join(env['GOROOT'], 'bin'),
    env.get('GOBIN'),
    os.path.join(ROOT, 'cipd'),
    os.path.join(ROOT, 'luci', 'appengine', 'components', 'tools'),
  ]
  paths_to_add.extend(os.path.join(p, '.vendor', 'bin') for p in vendor_paths)
  if go_appengine_path:
    paths_to_add.append(go_appengine_path)

  # Make sure not to add duplicates entries to PATH over and over again when
  # get_go_environ is invoked multiple times.
  paths_to_add = [p for p in paths_to_add if p and p not in path]
  env['PATH'] = os.pathsep.join(paths_to_add + path)

  # APPENGINE_DEV_APPSERVER is used by "appengine/aetest" package. If it's
  # missing, aetest will scan PATH looking for dev_appserver.py, possibly
  # finding it in some other place.
  if go_appengine_path:
    env['APPENGINE_DEV_APPSERVER'] = os.path.join(
        go_appengine_path, 'dev_appserver.py')
  else:
    env.pop('APPENGINE_DEV_APPSERVER', None)

  # Add a tag to the prompt
  prompt = env.get('PS1')
  if prompt:
    prompttag = env.get('INFRA_PROMPT_TAG', '[cr go] ')
    if prompttag not in prompt:
      env['PS1'] = prompttag + prompt

  return env


def get_go_exe(toolset_root):
  """Returns path to go executable."""
  return os.path.join(toolset_root, 'go', 'bin', 'go' + EXE_SFX)


def bootstrap(go_paths, logging_level, toolset=None):
  """Installs all dependencies in default locations.

  Supposed to be called at the beginning of some script (it modifies logger).

  Args:
    go_paths: list of paths to search for deps.lock, for each path deps.py
        will install all dependencies in <path>/.vendor/src/*.
    logging_level: logging level of bootstrap process.
    toolset: custom toolset if given, otherwise native.
  """
  logging.basicConfig()
  LOGGER.setLevel(logging_level)
  updated = ensure_toolset_installed(TOOLSET_ROOT, toolset)
  ensure_glide_installed(TOOLSET_ROOT)
  for p in go_paths:
    update_vendor_packages(p, force=updated)
  if updated:
    # GOPATH/pkg may have binaries generated with previous version of toolset,
    # they may not be compatible and "go build" isn't smart enough to rebuild
    # them.
    for p in go_paths:
      remove_directory([p, 'pkg'])


def prepare_go_environ(toolset=None):
  """Returns dict with environment variables to set to use Go toolset.

  Installs or updates the toolset and vendored dependencies if necessary.
  """
  bootstrap([WORKSPACE], logging.INFO, toolset)
  return get_go_environ(
    toolset_root=TOOLSET_ROOT,
    workspace=WORKSPACE,       # primary GOPATH with source code
    vendor_paths=[WORKSPACE],  # where to look for deps.yaml and .vendor dirs
    go_appengine_path=GO_APPENGINE)


def find_executable(name, workspaces):
  """Returns full path to an executable in some bin/ (in GOROOT or GOBIN)."""
  basename = name
  if EXE_SFX and basename.endswith(EXE_SFX):
    basename = basename[:-len(EXE_SFX)]
  roots = [os.path.join(TOOLSET_ROOT, 'go', 'bin')]
  for path in workspaces:
    roots.extend([
      os.path.join(path, '.vendor', 'bin'),
      os.path.join(path, 'bin'),
    ])
  for root in roots:
    full_path = os.path.join(root, basename + EXE_SFX)
    if os.path.exists(full_path):
      return full_path
  return name


def build_cross_toolset(toolset):
  env = prepare_go_environ(toolset)

  LOGGER.info('Building toolset for cross-compilation. Grab a coffee or tea...')
  tmp_env = env.copy()
  goroot = env['GOROOT']
  gosrc = os.path.join(goroot, 'src')
  for GOOS in ['darwin', 'windows', 'linux']:
    for GOARCH in ['amd64', '386']:
      tmp_env['GOOS'] = GOOS
      tmp_env['GOARCH'] = GOARCH
      LOGGER.info('adding %s %s support...', GOOS, GOARCH)
      with open(os.devnull, 'wb') as devnull:
        subprocess.check_call(
            ['./make.bash'], shell=True, cwd=gosrc, env=tmp_env,
            stdout=devnull, stderr=subprocess.STDOUT)

  target_name = '%s.%s-cross.tar.gz' % (TOOLSET_VERSION, toolset)

  # Use temp directory to clean up the garbage in case of failures.
  with temp_dir(TOOLSET_ROOT) as tmp:
    tmp_tgz = os.path.join(tmp, target_name)

    def save_target_file():
      local_file = os.path.join(TOOLSET_ROOT, target_name)
      os.rename(tmp_tgz, local_file)
      LOGGER.info('resulting target file saved in %s', local_file)

    LOGGER.info('archiving into %s...', tmp_tgz)
    with tarfile.open(tmp_tgz, 'w:gz') as f:
      f.add(goroot, 'go')

    action = raw_input('Upload to Google Storage? (y)es / (n)o ')
    if not action[0] in 'yY':
      LOGGER.warn('skipping upload to Google Storage.')
      save_target_file()
      return

    gstorage_url = 'gs://%s/%s/' % (BUCKET, BUCKET_PATH)
    LOGGER.info('Uploading to Google Storage %s', gstorage_url)
    status = subprocess.call([
        sys.executable, GSUTIL, '--force-version', '4.7', 'cp',
        tmp_tgz, gstorage_url])
    if status:
      LOGGER.error('Failed to upload to Google Storage: %s', status)
      save_target_file()
      return
    LOGGER.info('Uploaded to Google Storage %s%s', gstorage_url, target_name)
  print 'Ensure that following item is TOOLSET_VARIANTS dict of go/boostrap.py'
  print '\'%s-cross\': \'%s-cross.tar.gz\'' % (toolset, toolset)


def main(args):
  parser = argparse.ArgumentParser(
      description=sys.modules[__name__].__doc__,
      formatter_class=argparse.RawDescriptionHelpFormatter)
  parser.add_argument(
      '--toolset', action='store', default=get_default_toolset(),
      help='Specific toolset to install, native by default')
  parser.add_argument(
      '--build-cross', dest='build_cross', action='store_true', default=False,
      help='Create cross-compiling toolset based on platform native toolset.')

  opts = parser.parse_args(args)
  if opts.build_cross and opts.toolset != get_default_toolset():
    parser.error(
        'building toolset for cross compilation is supported only based on '
        'native toolset "%s", but "%s" toolset was selected' % (
        get_default_toolset(), opts.toolset))
  if opts.build_cross and sys.platform != 'linux2':
    parser.error(
        'building toolset for cross compilation is currently supported only on '
        'Linux')

  bootstrap([WORKSPACE], logging.DEBUG, toolset=opts.toolset)
  if opts.build_cross:
    build_cross_toolset(opts.toolset)
  return 0


if __name__ == '__main__':
  sys.exit(main(sys.argv[1:]))
