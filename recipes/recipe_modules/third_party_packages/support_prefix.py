# Copyright 2017 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

import collections
import itertools

from . import util


class Source(collections.namedtuple('_SourceBase', (
    # The Path to this Source's installation PREFIX.
    'prefix',
    # List of library names (e.g., "z", "curl") that this Source exports.
    'libs',
    # List of other Source entries that this Source depends on.
    'deps',
    # List of other shared libraries that this library requires when dynamically
    # linking.
    'shared_deps',
    ))):

  def __new__(cls, *args, **kwargs):
    include_prefix = kwargs.pop('include_prefix', None)
    src = super(Source, cls).__new__(cls, *args, **kwargs)
    src.include_prefix = include_prefix or src.prefix
    return src

  def _expand(self):
    exp = [self]
    for dep in self.deps:
      exp += dep._expand()
    return exp

  @property
  def bin_dir(self):
    return self.prefix.join('bin')

  @property
  def include_dirs(self):
    return [d.include_prefix.join('include') for d in self._expand()]

  @property
  def cppflags(self):
    return ['-I%s' % (d,) for d in self.include_dirs]

  @property
  def lib_dirs(self):
    return [d.include_prefix.join('lib') for d in self._expand()]

  @property
  def ldflags(self):
    return ['-L%s' % (d,) for d in self.lib_dirs]

  @property
  def full_static(self):
    full = []
    for s in self._expand():
      full += [str(s.include_prefix.join('lib', 'lib%s.a' % (lib,)))
               for lib in s.libs]
    return full

  @property
  def shared(self):
    link = []
    for s in self._expand():
      link += ['-l%s' % (lib,)
               for lib in itertools.chain(s.libs, s.shared_deps)]
    return link


class SupportPrefix(util.ModuleShim):
  """Provides a shared compilation and external library support context.

  Using SupportPrefix allows for coordination between packages (Git, Python)
  and inter-package dependencies (curl -> libz) to ensure that any given
  support library or function is built consistently and on-demand (at most once)
  for any given run.
  """

  _SOURCES = {
    'infra/third_party/source/autoconf': 'version:2.69',
    'infra/third_party/source/automake': 'version:1.15',
    'infra/third_party/source/gnu_sed': 'version:4.2.2',
    'infra/third_party/source/bzip2': 'version:1.0.6',
    'infra/third_party/source/libidn2': 'version:2.0.4',
    'infra/third_party/source/openssl': 'version:1.1.0f',
    'infra/third_party/source/mac_openssl_headers': 'version:0.9.8zh',
    'infra/third_party/source/pcre': 'version:8.41',
    'infra/third_party/source/pcre2': 'version:10.23',
    'infra/third_party/source/readline': 'version:7.0',
    'infra/third_party/source/zlib': 'version:1.2.11',
    'infra/third_party/source/curl': 'version:7.59.0',
    'infra/third_party/source/ncurses': 'version:6.0',
    'infra/third_party/source/nsl': 'version:1.0.4',
    'infra/third_party/source/sqlite-autoconf': 'version:3.19.3',
    'infra/third_party/pip-packages': 'version:9.0.1',
  }

  # The name and versions of the universal wheels in "pip-packages" that should
  # be installed alongside "pip".
  #
  # This will be versioned with _SOURCES's "pip-packages" entry.
  _PIP_PACKAGES_WHEELS = {
    'pip': '9.0.1',
    'setuptools': '36.0.1',
    'wheel': '0.30.0a0',
  }

  def __init__(self, api, base):
    super(SupportPrefix, self).__init__(api)
    self._api = api
    self._base = base
    self._sources_installed = False
    self._built = {}

  @staticmethod
  def update_mac_autoconf(env):
    # Several functions are declared in OSX headers that aren't actually
    # present in its standard libraries. Autoconf will succeed at detecting
    # them, only to fail later due to a linker error. Override these autoconf
    # variables via env to prevent this.
    env.update({
        'ac_cv_func_getentropy': 'n',
        'ac_cv_func_clock_gettime': 'n',
    })

  def ensure_sources(self):
    sources = self._base.join('sources')
    if not self._sources_installed:
      self.m.cipd.ensure(sources, self._SOURCES)
      self._sources_installed = True
    return sources

  def _build_once(self, key, build_fn):
    result = self._built.get(key)
    if result:
      return result

    build_name = '-'.join(e for e in key if e)
    workdir = self._base.join(build_name)

    with self.m.step.nest(build_name):
      self.m.file.ensure_directory('makedirs workdir', workdir)

      with self.m.context(cwd=workdir):
        self._built[key] = build_fn()
    return self._built[key]

  def _ensure_and_build_archive(self, name, tag, build_fn, archive_name=None,
                                variant=None):
    sources = self.ensure_sources()
    archive_name = archive_name or '%s-%s.tar.gz' % (
        name, tag.lstrip('version:'))

    base = archive_name
    for ext in ('.tar.gz', '.tgz', '.zip'):
      if base.endswith(ext):
        base = base[:-len(ext)]
        break

    def build_archive():
      archive = sources.join(archive_name)
      build_root = self.m.context.cwd
      prefix = build_root.join('prefix')

      self.m.python(
          'extract',
          self.resource('archive_util.py'),
          [
            archive,
            self.m.context.cwd,
          ])
      build = build_root.join(base) # Archive is extracted here.

      try:
        with self.m.context(cwd=build):
          build_fn(prefix, build_root)
      finally:
        pass
      return prefix

    key = (base, variant)
    return self._build_once(key, build_archive)

  def _build_openssl(self, tag, shell=False):
    def build_fn(prefix, _build_root):
      target = {
        ('mac', 'intel', 64): 'darwin64-x86_64-cc',
        ('linux', 'intel', 32): 'linux-x86',
        ('linux', 'intel', 64): 'linux-x86_64',
      }[(
        self.m.platform.name,
        self.m.platform.arch,
        self.m.platform.bits)]

      configure_cmd = [
        './Configure',
        '--prefix=%s' % (prefix,),
        'no-shared',
        target,
      ]
      if shell:
        configure_cmd = ['bash'] + configure_cmd

      self.m.step('configure', configure_cmd)
      self.m.step('make', ['make', '-j', str(self.m.platform.cpu_count)])

      # Install OpenSSL. Note that "install_sw" is an OpenSSL-specific
      # sub-target that only installs headers and library, saving time.
      self.m.step('install', ['make', 'install_sw'])

    return Source(
        prefix=self._ensure_and_build_archive('openssl', tag, build_fn),
        libs=['ssl', 'crypto'],
        deps=[],
        shared_deps=[])

  def ensure_openssl(self):
    return self._build_openssl('version:1.1.0f')

  def ensure_mac_native_openssl(self):
    return self._build_openssl('version:0.9.8zh', shell=True)

  def _generic_build(self, name, tag, archive_name=None, configure_args=None,
                     libs=None, deps=None, shared_deps=None):
    def build_fn(prefix, _build_root):
      self.m.step('configure', [
        './configure',
        '--prefix=%s' % (prefix,),
      ] + (configure_args or []))
      self.m.step('make', [
        'make', 'install',
        '-j', str(self.m.platform.cpu_count),
        ])
    return Source(
        prefix=self._ensure_and_build_archive(
          name, tag, build_fn, archive_name=archive_name),
        deps=deps or [],
        libs=libs or [name],
        shared_deps=shared_deps or [])

  def ensure_curl(self):
    zlib = self.ensure_zlib()
    libidn2 = self.ensure_libidn2()

    env = {}
    configure_args = [
      '--disable-ldap',
      '--disable-shared',
      '--without-librtmp',
      '--with-zlib=%s' % (str(zlib.prefix,)),
      '--with-libidn2=%s' % (str(libidn2.prefix,)),
    ]
    deps = []
    shared_deps = []
    if self.m.platform.is_mac:
      configure_args += ['--with-darwinssl']
    elif self.m.platform.is_linux:
      ssl = self.ensure_openssl()
      env['LIBS'] = ' '.join(['-ldl', '-lpthread'])
      configure_args += ['--with-ssl=%s' % (str(ssl.prefix),)]
      deps += [ssl]
      shared_deps += ['dl', 'pthread']

    with self.m.context(env=env):
      return self._generic_build('curl', 'version:7.59.0',
                                 configure_args=configure_args, deps=deps,
                                 shared_deps=shared_deps)

  def ensure_pcre(self):
    return self._generic_build(
        'pcre', 'version:8.41',
        configure_args=[
          '--enable-static',
          '--disable-shared',
        ])

  def ensure_pcre2(self):
    return self._generic_build(
        'pcre2',
        'version:10.23',
        libs=['pcre2-8'],
        configure_args=[
          '--enable-static',
          '--disable-shared',
        ])

  def ensure_nsl(self):
    return self._generic_build('nsl', 'version:1.0.4',
        archive_name='libnsl-1.0.4.tar.gz',
        configure_args=['--disable-shared'])

  def ensure_ncurses(self):
    # The "ncurses" package, by default, uses a fixed-path location for terminal
    # information. This is not portable, so we need to disable it. Instead, we
    # will compile ncurses with a set of hand-picked custom terminal information
    # data baked in, as well as the ability to probe terminal via termcap if
    # needed.
    #
    # To do this, we need to bulid in multiple stages:
    # 1) Generic configure / make so thast the "tic" (terminfo compiler) and
    #    "toe" (table of entries) commands are built.
    # 2) Use "toe" tool to dump the set of available profiles and groom it.
    # 3) Build library with no database support using "tic" from (1), and
    #    configure it to statically embed all of the profiles from (2).
    def build_fn(prefix, build_root):
      src = self.m.context.cwd

      tic_build = build_root.join('tic_build')
      tic_prefix = build_root.join('tic_prefix')
      tic_bin = tic_prefix.join('bin')

      self.m.file.ensure_directory('makedirs tic build', tic_build)
      with self.m.context(cwd=tic_build):
        self.m.step('configure tic', [
          src.join('configure'),
          '--prefix=%s' % (tic_prefix,),
        ])
        self.m.step('make tic', [
          'make', 'install',
          '-j', str(self.m.platform.cpu_count),
          ])

      # Determine the list of all supported profiles. The "toe" command (table
      # of entries) will dump a list.
      toe = self.m.step(
          'get profiles',
          [tic_bin.join('toe')],
          stdout=self.m.raw_io.output_text(),
          step_test_data=lambda: self.m.raw_io.test_api.stream_output(
            '\n'.join([
              'foo        The foo profile.',
              'bar        The bar profile.',
              '9term      Should be pruned.',
              'guru+fake  Should be pruned.',
            ])),
      )
      fallbacks = [l.split()[0].strip() for l in toe.stdout.splitlines()]

      # Strip out fallbacks with bugs.
      #
      # This currently leaves 1591 profiles behind, which will be statically
      # compiled into the library.
      fallbacks = [f for f in fallbacks if (
        f and not any(f.startswith(x) for x in (
          # Some profiles do not generate valid C, either because:
          # - They begin with a number, which is not valid in C.
          # - They are flattened to a duplicate symbol as another profile. This
          #   usually happens when there are "+" and "-" variants; we choose
          #   "-".
          # - They include quotes in the description name.
          #
          # None of these identified terminals are really important, so we will
          # just avoid processing them.
          '9term', 'guru+', 'hp+', 'tvi912b+', 'tvi912b-vb', 'tvi920b-vb',
          'att4415+', 'nsterm+', 'xnuppc+', 'xterm+', 'wyse-vp',
        ))
      )]
      toe.presentation.step_text = 'Embedding %d profile(s)' % (
          len(fallbacks),)

      # Run the remainder of our build with our generated "tic" on PATH.
      #
      # Note that we only run "install.libs". Standard "install" expects the
      # full database to exist, and this will not be the case since we are
      # explicitly disabling it.
      with self.m.context(env_prefixes={'PATH': [tic_bin]}):
        self.m.step('configure', [
          './configure',
          '--prefix=%s' % (prefix,),
          '--disable-database',
          '--disable-db-install',
          '--enable-termcap',
          '--with-fallbacks=%s' % (','.join(fallbacks),),
        ])
        self.m.step('make', [
          'make', 'install.libs',
          '-j', str(self.m.platform.cpu_count),
          ])

    return Source(
        prefix=self._ensure_and_build_archive(
          'ncurses', 'version:6.0', build_fn),
        deps=[],
        libs=['panel', 'ncurses'],
        shared_deps=[])

  def ensure_zlib(self):
    return self._generic_build('zlib', 'version:1.2.11', libs=['z'],
                               configure_args=['--static'])

  def ensure_libidn2(self):
    return self._generic_build('libidn2', 'version:2.0.4', libs=['idn2'],
                               configure_args=['--enable-static=yes',
                                               '--enable-shared=no'])

  def ensure_sqlite(self):
    return self._generic_build('sqlite', 'version:3.19.3', libs=['sqlite3'],
        configure_args=[
          '--enable-static',
          '--disable-shared',
          '--with-pic',
          '--enable-fts5',
          '--enable-json1',
          '--enable-session',
        ],
        archive_name='sqlite-autoconf-3190300.tar.gz')

  def ensure_bzip2(self):
    def build_fn(prefix, _build_root):
      self.m.step('make', [
        'make',
        'install',
        'PREFIX=%s' % (prefix,),
      ])
    return Source(
        prefix=self._ensure_and_build_archive(
          'bzip2', 'version:1.0.6', build_fn),
        deps=[],
        libs=['bz2'],
        shared_deps=[])

  def ensure_readline(self):
    ncurses = self.ensure_ncurses()
    return self._generic_build('readline', 'version:7.0',
        deps=[ncurses])

  def ensure_autoconf(self):
    return self._generic_build('autoconf', 'version:2.69')

  def ensure_automake(self):
    autoconf = self.ensure_autoconf()
    with self.m.context(env_prefixes={'PATH': [autoconf.bin_dir]}):
      return self._generic_build('automake', 'version:1.15')

  def ensure_gnu_sed(self):
    return self._generic_build('gnu_sed', 'version:4.2.2',
        archive_name='sed-4.2.2.tar.gz')

  def ensure_pip_installer(self):
    """Returns information about the pip installation.

    Returns: (get_pip, links_path, wheels)
      get_pip (Path): Path to the "get-pip.py" script.
      links (Path): Path to the links directory containing all installation
          wheels.
      wheels (dict): key/value mapping of "pip" installation packages names
          and their verisons.
    """
    sources = self.ensure_sources()
    return (
        sources.join('get-pip.py'),
        sources,
        self._PIP_PACKAGES_WHEELS,
    )
