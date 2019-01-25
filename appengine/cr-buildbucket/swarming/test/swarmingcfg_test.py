# Copyright 2015 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

import json

from components import utils
utils.fix_protobuf_package()

from google import protobuf
from parameterized import parameterized

from components import config as config_component
from testing_utils import testing

from proto.config import project_config_pb2
from proto.config import service_config_pb2
from test import config_test
from swarming import swarmingcfg
import errors


class ProjectCfgTest(testing.AppengineTestCase):

  def cfg_test(self, swarming_text, mixins_text, expected_errors):
    swarming_cfg = project_config_pb2.Swarming()
    protobuf.text_format.Merge(swarming_text, swarming_cfg)

    buildbucket_cfg = project_config_pb2.BuildbucketCfg()
    protobuf.text_format.Merge(mixins_text, buildbucket_cfg)
    mixins = {m.name: m for m in buildbucket_cfg.builder_mixins}

    ctx = config_component.validation.Context()
    swarmingcfg.validate_project_cfg(swarming_cfg, mixins, True, ctx)
    self.assert_errors(ctx, expected_errors)

  def assert_errors(self, ctx, expected_errors):
    self.assertEqual(
        map(config_test.errmsg, expected_errors),
        ctx.result().messages
    )

  def test_valid(self):
    self.cfg_test(
        '''
          hostname: "example.com"
          builder_defaults {
            swarming_tags: "master:master.a"
            dimensions: "cores:8"
            dimensions: "60:cores:64"
            dimensions: "pool:default"
            dimensions: "cpu:x86-64"
            service_account: "bot"
          }
          builders {
            name: "release"
            swarming_tags: "a:b'"
            dimensions: "os:Linux"
            dimensions: "cpu:"
            service_account: "robot@example.com"
            caches {
              name: "git_chromium"
              path: "git_cache"
            }
            recipe {
              repository: "https://x.com"
              name: "foo"
              properties: "a:b'"
              properties_j: "x:true"
            }
          }
          builders {
            name: "release cipd"
            recipe {
              cipd_package: "some/package"
              name: "foo"
            }
          }
        ''', '', []
    )

  def test_validate_recipe_properties(self):

    def test(properties, properties_j, expected_errors):
      ctx = config_component.validation.Context()
      swarmingcfg.validate_recipe_properties(properties, properties_j, ctx)
      self.assertEqual(
          map(config_test.errmsg, expected_errors),
          ctx.result().messages
      )

    test([], [], [])

    runtime = '$recipe_engine/runtime:' + json.dumps({
        'is_luci': False,
        'is_experimental': True,
    })
    test(
        properties=[
            '',
            ':',
            'buildbucket:foobar',
            'x:y',
        ],
        properties_j=[
            'x:"y"',
            'y:b',
            'z',
            runtime,
        ],
        expected_errors=[
            'properties \'\': does not have a colon',
            'properties \':\': key not specified',
            'properties \'buildbucket:foobar\': reserved property',
            'properties_j \'x:"y"\': duplicate property',
            'properties_j \'y:b\': No JSON object could be decoded',
            'properties_j \'z\': does not have a colon',
            'properties_j %r: key \'is_luci\': reserved key' % runtime,
            'properties_j %r: key \'is_experimental\': reserved key' % runtime,
        ]
    )

    test([], ['$recipe_engine/runtime:1'], [
        ('properties_j \'$recipe_engine/runtime:1\': '
         'not a JSON object'),
    ])

    test([], ['$recipe_engine/runtime:{"unrecognized_is_fine": 1}'], [])

  def test_bad(self):
    self.cfg_test(
        '''
          builders {}
        ''',
        '',
        [
            'hostname: unspecified',
            'builder #1: name: unspecified',
            'builder #1: recipe: name: unspecified',
            'builder #1: recipe: specify either cipd_package or repository',
        ],
    )

    self.cfg_test(
        '''
          hostname: "example.com"
          builder_defaults {
            recipe {
              name: "meeper"
              repository: "https://example.com"
            }
          }
          builders {
            name: "meep"
          }
          builders {
            name: "meep"
          }
        ''',
        '',
        [
            'builder meep: name: duplicate',
        ],
    )

    self.cfg_test(
        '''
          hostname: "example.com"
          builders {
            name: ":/:"
          }
        ''',
        '',
        [
            ('builder :/:: name: invalid char(s) u\'/:\'. '
             'Alphabet: "%s"') % errors.BUILDER_NAME_VALID_CHARS,
        ],
    )

    self.cfg_test(
        '''
          hostname: "example.com"
          builders {
            name: "veeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeery"
                  "looooooooooooooooooooooooooooooooooooooooooooooooooooooooong"
                  "naaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaame"
          }
        ''',
        '',
        [(
            'builder '
            'veeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeery'
            'looooooooooooooooooooooooooooooooooooooooooooooooooooooooong'
            'naaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaame: '
            'name: length is > 128'
        )],
    )

    self.cfg_test(
        '''
          hostname: "example.com"
          builder_defaults {name: "x"}
          builders {
            name: "release"
            dimensions: "pool:a"
            recipe {
              repository: "https://x.com"
              name: "foo"
            }
          }
        ''',
        '',
        [
            'builder_defaults: name: not allowed',
        ],
    )

    self.cfg_test(
        '''
          hostname: "https://example.com"
          task_template_canary_percentage { value: 102 }
          builder_defaults {
            swarming_tags: "wrong"
          }
          builders {
            swarming_tags: "wrong2"
          }
          builders {
            name: "b2"
            swarming_tags: "builder:b2"
            caches {}
            caches { name: "a/b" path: "a" }
            caches { name: "b" path: "a\\c" }
            caches { name: "c" path: "a/.." }
            caches { name: "d" path: "/a" }
            priority: 300
          }
        ''',
        '',
        [
            'hostname: must not contain "://"',
            'task_template_canary_percentage.value must must be in [0, 100]',
            'builder_defaults: tag #1: does not have ":": wrong',
            'builder #1: tag #1: does not have ":": wrong2',
            (
                'builder b2: tag #1: do not specify builder tag; '
                'it is added by swarmbucket automatically'
            ),
            'builder b2: cache #1: name: required',
            'builder b2: cache #1: path: required',
            (
                'builder b2: cache #2: '
                'name: "a/b" does not match ^[a-z0-9_]{1,4096}$'
            ),
            (
                'builder b2: cache #3: path: cannot contain \\. '
                'On Windows forward-slashes will be replaced with back-slashes.'
            ),
            'builder b2: cache #4: path: cannot contain ".."',
            'builder b2: cache #5: path: cannot start with "/"',
            'builder b2: priority: must be in [0, 200] range; got 300',
        ],
    )

    self.cfg_test(
        '''
          hostname: "example.com"
          builders {
            name: "rel"
            caches { path: "a" name: "a" }
            caches { path: "a" name: "a" }
          }
        ''',
        '',
        [
            'builder rel: cache #2: duplicate name',
            'builder rel: cache #2: duplicate path',
        ],
    )

    self.cfg_test(
        '''
          hostname: "example.com"
          builders {
            name: "rel"
            caches { path: "a" name: "a" wait_for_warm_cache_secs: 61 }
          }
        ''',
        '',
        [
            'builder rel: cache #1: wait_for_warm_cache_secs: must be rounded '
            'on 60 seconds',
        ],
    )

    self.cfg_test(
        '''
          hostname: "example.com"
          builders {
            name: "rel"
            caches { path: "a" name: "a" wait_for_warm_cache_secs: 59 }
          }
        ''',
        '',
        [
            'builder rel: cache #1: wait_for_warm_cache_secs: must be at least '
            '60 seconds'
        ],
    )

    self.cfg_test(
        '''
          hostname: "example.com"
          builders {
            name: "rel"
            caches { path: "a" name: "a" wait_for_warm_cache_secs: 60 }
            caches { path: "b" name: "b" wait_for_warm_cache_secs: 120 }
            caches { path: "c" name: "c" wait_for_warm_cache_secs: 180 }
            caches { path: "d" name: "d" wait_for_warm_cache_secs: 240 }
            caches { path: "e" name: "e" wait_for_warm_cache_secs: 300 }
            caches { path: "f" name: "f" wait_for_warm_cache_secs: 360 }
            caches { path: "g" name: "g" wait_for_warm_cache_secs: 420 }
            caches { path: "h" name: "h" wait_for_warm_cache_secs: 480 }
          }
        ''',
        '',
        [
            'builder rel: too many different (8) wait_for_warm_cache_secs '
            'values; max 7',
        ],
    )

    self.cfg_test(
        '''
          hostname: "example.com"
          builders {
            name: "b"
            service_account: "not an email"
          }
        ''',
        '',
        [
            'builder b: service_account: value "not an email" does not match '
            '^[0-9a-zA-Z_\\-\\.\\+\\%]+@[0-9a-zA-Z_\\-\\.]+$',
        ],
    )

    self.cfg_test(
        '''
          hostname: "example.com"
          builders {
            name: "b"
            expiration_secs: 158400  # 44h
            execution_timeout_secs: 14400  # 4h
          }
        ''',
        '',
        [
            'builder b: expiration_secs + execution_timeout_secs '
            'must be at most 47h'
        ],
    )

  @parameterized.expand([
      (['a:b'], ''),
      ([''], 'dimension "": does not have ":"'),
      (
          ['caches:a'],
          (
              'dimension "caches:a": dimension key must not be "caches"; '
              'caches must be declared via caches field'
          ),
      ),
      (
          ['a:b', 'a:c'],
          (
              'dimension "a:c": '
              'multiple values for dimension key "a" and expiration 0s'
          ),
      ),
      ([':'], 'dimension ":": no key'),
      (
          ['a.b:c'],
          (
              'dimension "a.b:c": '
              r'key "a.b" does not match pattern "^[a-zA-Z\_\-]+$"'
          ),
      ),
      (['0:'], 'dimension "0:": has expiration_secs but missing value'),
      (['a:', '60:a:b'], 'dimension "60:a:b": mutually exclusive with "a:"'),
      (
          ['-1:a:1'],
          (
              'dimension "-1:a:1": '
              'expiration_secs is outside valid range; up to 21 days'
          ),
      ),
      (
          ['1:a:b'],
          'dimension "1:a:b": expiration_secs must be a multiple of 60 seconds',
      ),
      (
          ['1814400:a:1'],  # 21*24*60*6
          '',
      ),
      (
          ['1814401:a:1'],  # 21*24*60*60+
          (
              'dimension "1814401:a:1": '
              'expiration_secs is outside valid range; up to 21 days'
          ),
      ),
      (
          [
              '60:a:1',
              '120:a:1',
              '180:a:1',
              '240:a:1',
              '300:a:1',
              '360:a:1',
              '420:a:1',
          ],
          'at most 6 different expiration_secs values can be used',
      ),
  ])
  def test_validate_dimensions(self, dimensions, expected_error):
    ctx = config_component.validation.Context()
    swarmingcfg._validate_dimensions('dimension', dimensions, ctx)
    self.assert_errors(ctx, [expected_error] if expected_error else [])

  def test_default_recipe(self):
    self.cfg_test(
        '''
          hostname: "chromium-swarm.appspot.com"
          builder_defaults {
            dimensions: "pool:default"
            recipe {
              repository: "https://x.com"
              name: "foo"
              properties: "a:b"
              properties: "x:y"
           }
          }
          builders { name: "debug" }
          builders {
            name: "release"
            recipe {
              properties: "a:c"
              properties_j: "x:null"
            }
          }
        ''', '', []
    )

  def test_default_recipe_bad(self):
    self.cfg_test(
        '''
          hostname: "chromium-swarm.appspot.com"
          builder_defaults {
            dimensions: "pool:default"
            recipe {
              name: "foo"
              properties: "a"
            }
          }
          builders { name: "debug" }
        ''',
        '',
        ['builder_defaults: recipe: properties u\'a\': does not have a colon'],
    )

  def test_cipd_and_repository_bad(self):
    self.cfg_test(
        '''
          hostname: "chromium-swarm.appspot.com"
          builders {
            name: "debug"
            dimensions: "pool:default"
            recipe {
              name: "foo"
              repository: "https://example.com"
              cipd_package: "some/package"
            }
          }
        ''', '', [
            (
                'builder debug: recipe: specify either cipd_package '
                'or repository, not both'
            ),
        ]
    )

  def test_clear_recipe_repository(self):
    self.cfg_test(
        '''
          hostname: "chromium-swarm.appspot.com"
          builder_defaults {
            recipe {
              name: "foo"
              repository: "https://example.com"
            }
          }
          builders {
            name: "debug"
            recipe {
              repository: "-"
              cipd_package: "some/package"
            }
          }
        ''', '', []
    )

  def test_validate_builder_mixins(self):

    def test(cfg_text, expected_errors):
      ctx = config_component.validation.Context()
      cfg = project_config_pb2.BuildbucketCfg()
      protobuf.text_format.Merge(cfg_text, cfg)
      swarmingcfg.validate_builder_mixins(cfg.builder_mixins, ctx)
      self.assertEqual(
          map(config_test.errmsg, expected_errors),
          ctx.result().messages
      )

    test(
        '''
          builder_mixins {
            name: "a"
            dimensions: "a:b"
            dimensions: "60:a:c"
          }
          builder_mixins {
            name: "b"
            mixins: "a"
            dimensions: "a:b"
          }
        ''', []
    )

    test(
        '''
          builder_mixins {
            name: "b"
            mixins: "a"
          }
          builder_mixins {
            name: "a"
          }
        ''', []
    )

    test(
        '''
          builder_mixins {}
        ''', ['builder_mixin #1: name: unspecified']
    )

    test(
        '''
          builder_mixins { name: "a" }
          builder_mixins { name: "a" }
        ''', ['builder_mixin a: name: duplicate']
    )

    test(
        '''
          builder_mixins {
            name: "a"
            mixins: ""
          }
        ''', ['builder_mixin a: referenced mixin name is empty']
    )

    test(
        '''
          builder_mixins {
            name: "a"
            mixins: "b"
          }
        ''', ['builder_mixin a: mixin "b" is not defined']
    )

    test(
        '''
          builder_mixins {
            name: "a"
            mixins: "a"
          }
        ''', [
            'circular mixin chain: a -> a',
        ]
    )

    test(
        '''
          builder_mixins {
            name: "a"
            mixins: "b"
          }
          builder_mixins {
            name: "b"
            mixins: "c"
          }
          builder_mixins {
            name: "c"
            mixins: "a"
          }
        ''', [
            'circular mixin chain: a -> b -> c -> a',
        ]
    )

  def test_builder_with_mixins(self):

    def test(cfg_text, expected_errors):
      ctx = config_component.validation.Context()
      cfg = project_config_pb2.BuildbucketCfg()
      protobuf.text_format.Merge(cfg_text, cfg)
      swarmingcfg.validate_builder_mixins(cfg.builder_mixins, ctx)
      self.assertEqual([], ctx.result().messages)
      mixins = {m.name: m for m in cfg.builder_mixins}
      swarmingcfg.validate_project_cfg(
          cfg.buckets[0].swarming, mixins, True, ctx
      )
      self.assertEqual(
          map(config_test.errmsg, expected_errors),
          ctx.result().messages
      )

    test(
        '''
          builder_mixins {
            name: "a"

            dimensions: "cores:8"
            dimensions: "cpu:x86-64"
            dimensions: "os:Linux"
            dimensions: "pool:default"
            caches {
              name: "git"
              path: "git"
            }
            recipe {
              repository: "https://x.com"
              name: "foo"
              properties: "a:b'"
              properties_j: "x:true"
            }
          }
          builder_mixins {
            name: "b"
            mixins: "a"
          }
          builder_mixins {
            name: "c"
            mixins: "a"
            mixins: "b"
          }
          buckets {
            name: "a"
            swarming {
              hostname: "chromium-swarm.appspot.com"
              builders {
                name: "release"
                mixins: "b"
                mixins: "c"
              }
            }
          }
        ''', []
    )


class ServiceCfgTest(testing.AppengineTestCase):

  def cfg_test(self, swarming_text, expected_errors):
    ctx = config_component.validation.Context()

    settings = service_config_pb2.SwarmingSettings()
    protobuf.text_format.Merge(swarming_text, settings)

    swarmingcfg.validate_service_cfg(settings, ctx)
    self.assertEqual(
        map(config_test.errmsg, expected_errors),
        ctx.result().messages
    )

  def test_valid(self):
    self.cfg_test('', [])

  def test_schema_in_hostname(self):
    self.cfg_test(
        '''
          milo_hostname: "https://milo.example.com"
        ''', [
            'milo_hostname: must not contain "://"',
        ]
    )
