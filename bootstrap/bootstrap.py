#!/usr/bin/env python
# Copyright 2014 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

import argparse
import contextlib
import glob
import logging
import os
import re
import shutil
import subprocess
import sys
import tempfile

from util import STORAGE_URL, OBJECT_URL, LOCAL_STORAGE_PATH, LOCAL_OBJECT_URL
from util import read_deps, merge_deps, print_deps, platform_tag

LOGGER = logging.getLogger(__name__)

# /path/to/infra
ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))

PYTHON_BAT_WIN = '@%~dp0\\..\\Scripts\\python.exe %*'


class NoWheelException(Exception):
  def __init__(self, name, version, build, source_sha):
    super(NoWheelException, self).__init__(
        'No matching wheel found for (%s==%s (build %s_%s))' %
        (name, version, build, source_sha))


def check_pydistutils():
  if os.path.exists(os.path.expanduser('~/.pydistutils.cfg')):
    print >> sys.stderr, '\n'.join([
      '',
      '',
      '=========== ERROR ===========',
      'You have a ~/.pydistutils.cfg file, which interferes with the ',
      'infra virtualenv environment. Please move it to the side and bootstrap ',
      'again. Once infra has bootstrapped, you may move it back.',
      '',
      'Upstream bug: https://github.com/pypa/virtualenv/issues/88/',
      ''
    ])
    sys.exit(1)


def ls(prefix):
  from pip._vendor import requests  # pylint: disable=E0611
  data = requests.get(STORAGE_URL, params=dict(
      prefix=prefix,
      fields='items(name,md5Hash)'
  )).json()
  entries = data.get('items', [])
  for entry in entries:
    entry['md5Hash'] = entry['md5Hash'].decode('base64').encode('hex')
    entry['local'] = False
  # Also look in the local cache
  entries.extend([
    {'name': fname, 'md5Hash': None, 'local': True}
    for fname in glob.glob(os.path.join(LOCAL_STORAGE_PATH,
                                        prefix.split('/')[-1] + '*'))])
  return entries


def sha_for(deps_entry):
  if 'rev' in deps_entry:
    return deps_entry['rev']
  else:
    return deps_entry['gs'].split('.')[0]


def get_links(deps):
  import pip.wheel  # pylint: disable=E0611
  plat_tag = platform_tag()

  links = []

  for name, dep in deps.iteritems():
    version, source_sha = dep['version'] , sha_for(dep)
    prefix = 'wheels/{}-{}-{}_{}'.format(name, version, dep['build'],
                                         source_sha)
    generic_link = None
    binary_link = None
    local_link = None

    for entry in ls(prefix):
      fname = entry['name'].split('/')[-1]
      md5hash = entry['md5Hash']
      wheel_info = pip.wheel.Wheel.wheel_file_re.match(fname)
      if not wheel_info:
        LOGGER.warn('Skipping invalid wheel: %r', fname)
        continue

      if pip.wheel.Wheel(fname).supported():
        if entry['local']:
          link = LOCAL_OBJECT_URL.format(entry['name'])
          local_link = link
          continue
        else:
          link = OBJECT_URL.format(entry['name'], md5hash)
        if fname.endswith('none-any.whl'):
          if generic_link:
            LOGGER.warning(
              'Found more than one generic matching wheel for %r: %r',
              prefix, dep)
            continue
          generic_link = link
        elif plat_tag in fname:
          if binary_link:
            LOGGER.warning(
              'Found more than one binary matching wheel for %r: %r\n'
              '  Picking:        %s\n'
              '  Also available: %s\n',
              prefix, dep, binary_link, link)
            continue
          binary_link = link

    if not binary_link and not generic_link and not local_link:
      raise NoWheelException(name, version, dep['build'], source_sha)

    links.append(local_link or binary_link or generic_link)

  return links


@contextlib.contextmanager
def html_index(links):
  tf = tempfile.mktemp('.html')
  try:
    with open(tf, 'w') as f:
      print >> f, '<html><body>'
      for link in links:
        print >> f, '<a href="%s">wat</a>' % link
      print >> f, '</body></html>'
    yield tf
  finally:
    os.unlink(tf)


def install(deps):
  if sys.platform.startswith('win'):
    # On Windows, "pip" is installed as a standalone binary called "pip.exe".
    pip = [os.path.join(sys.prefix, 'Scripts', 'pip')]
  else:
    # On Linux, "pip" is a "#!/...python"-bootstrapped wrapper. Because of
    # shebang length limitations, we will manually run this through the
    # Python interpreter rather than relying on shebang interpretation.
    pip = [
        os.path.join(sys.prefix, 'bin', 'python'),
        os.path.join(sys.prefix, 'bin', 'pip'),
    ]

  links = get_links(deps)
  with html_index(links) as ipath:
    requirements = []
    # TODO(iannucci): Do this as a requirements.txt
    for name, deps_entry in deps.iteritems():
      if not deps_entry.get('implicit'):
        requirements.append('%s==%s' % (name, deps_entry['version']))
    subprocess.check_call(
        pip + ['install', '--no-index', '-f', ipath] + requirements)


def activate_env(env, deps, quiet=False, run_within_virtualenv=False):
  if hasattr(sys, 'real_prefix'):
    if not run_within_virtualenv:
      LOGGER.error('Already activated environment!')
      return
    LOGGER.info('Discarding current VirtualEnv (--run-within-virtualenv)')
    sys.prefix = sys.real_prefix

  if not quiet:
    print 'Activating environment: %r' % env
  assert isinstance(deps, dict)

  manifest_path = os.path.join(env, 'manifest.pyl')
  cur_deps = read_deps(manifest_path)
  if cur_deps != deps:
    if not quiet:
      print '  Removing old environment: %r' % cur_deps
    shutil.rmtree(env, ignore_errors=True)
    cur_deps = None

  if cur_deps is None:
    check_pydistutils()

    if not quiet:
      print '  Building new environment'
    # Add in bundled virtualenv lib
    sys.path.insert(0, os.path.join(os.path.dirname(__file__), 'virtualenv'))
    import virtualenv  # pylint: disable=F0401
    virtualenv.create_environment(
        env, search_dirs=virtualenv.file_search_dirs())

    # Hack: On windows orig-prefix.txt contains the hardcoded path
    # "E:\b\depot_tools\python276_bin", but some systems have depot_tools
    # installed on C:\ instead, so fiddle site.py to try loading it from there
    # as well.
    if sys.platform.startswith('win'):
      site_py_path = os.path.join(env, 'Lib\\site.py')
      with open(site_py_path) as fh:
        site_py = fh.read()

      m = re.search(r'( +)sys\.real_prefix = .*', site_py)
      replacement = ('%(indent)sif (sys.real_prefix.startswith("E:\\\\") and\n'
                     '%(indent)s    not os.path.exists(sys.real_prefix)):\n'
                     '%(indent)s  cand = "C:\\\\setup" + sys.real_prefix[4:]\n'
                     '%(indent)s  if os.path.exists(cand):\n'
                     '%(indent)s    sys.real_prefix = cand\n'
                     '%(indent)s  else:\n'
                     '%(indent)s    sys.real_prefix = "C" + sys.real_prefix'
                        '[1:]\n'
                     % {'indent': m.group(1)})

      site_py = site_py[:m.end(0)] + '\n' + replacement + site_py[m.end(0):]
      with open(site_py_path, 'w') as fh:
        fh.write(site_py)

  if not quiet:
    print '  Activating environment'
  # Ensure hermeticity during activation.
  os.environ.pop('PYTHONPATH', None)
  bin_dir = 'Scripts' if sys.platform.startswith('win') else 'bin'
  activate_this = os.path.join(env, bin_dir, 'activate_this.py')
  execfile(activate_this, dict(__file__=activate_this))

  if cur_deps is None:
    if not quiet:
      print '  Installing deps'
      print_deps(deps, indent=2, with_implicit=False)
    install(deps)
    virtualenv.make_environment_relocatable(env)
    with open(manifest_path, 'wb') as f:
      f.write(repr(deps) + '\n')

  # Create bin\python.bat on Windows to unify path where Python is found.
  if sys.platform.startswith('win'):
    bin_path = os.path.join(env, 'bin')
    if not os.path.isdir(bin_path):
      os.makedirs(bin_path)
    python_bat_path = os.path.join(bin_path, 'python.bat')
    if not os.path.isfile(python_bat_path):
      with open(python_bat_path, 'w') as python_bat_file:
        python_bat_file.write(PYTHON_BAT_WIN)

  if not quiet:
    print 'Done creating environment'


def main(args):
  parser = argparse.ArgumentParser()
  parser.add_argument('--deps-file', '--deps_file', action='append',
                      required=True,
                      help='Path to deps.pyl file (may be used multiple times)')
  parser.add_argument('-q', '--quiet', action='store_true', default=False,
                      help='Supress all output')
  parser.add_argument('-r', '--run-within-virtualenv', action='store_true',
                      help='Run even if the script is being run within a '
                           'VirtualEnv.')
  parser.add_argument('env_path',
                      help='Path to place environment (default: %(default)s)',
                      default='ENV')
  opts = parser.parse_args(args)

  deps = merge_deps(opts.deps_file)
  activate_env(opts.env_path, deps, opts.quiet, opts.run_within_virtualenv)


if __name__ == '__main__':
  logging.basicConfig()
  LOGGER.setLevel(logging.DEBUG)
  sys.exit(main(sys.argv[1:]))
