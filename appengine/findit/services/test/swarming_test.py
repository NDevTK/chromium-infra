# Copyright 2018 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

from datetime import datetime
import json
import mock

from google.appengine.api import app_identity

from gae_libs import token
from infra_api_clients.swarming import swarming_util
from infra_api_clients.swarming.swarming_task_data import SwarmingTaskData
from infra_api_clients.swarming.swarming_task_request import SwarmingTaskRequest
from libs import time_util
from services import swarming
from services.parameters import TestFailedSteps
from waterfall import waterfall_config
from waterfall.test import wf_testcase

_SAMPLE_BUILD_STEP_DATA = [{
    'failure': True,
    'task_id': '2944afa502297110',
    'internal_failure': False,
    'name': 'unit_tests/Windows-XP-SP3/123c706dde/XP Tests (1)/39454',
    'state': 112,
    'tags': [
        'buildername:XP Tests (1)', 'buildnumber:39442', 'cpu:x86-32',
        'data:123c706ddeeadcc761d612f732ddde9322400446', 'gpu:none',
        'master:chromium.win', 'name:unit_tests', 'os:Windows-XP-SP3',
        'pool:Chrome', 'priority:25', 'project:chromium', 'purpose:CI',
        'purpose:post-commit', 'stepname:unit_tests', 'user:'
    ],
    'try_number': 0,
    'user': '',
    'outputs_ref': {
        'isolatedserver': 'https://isolateserver.appspot.com',
        'namespace': 'default-gzip',
        'isolated': 'isolatedhashunittests'
    }
}, {
    'failure': False,
    'task_id': '2944afa502297111',
    'internal_failure': False,
    'tags': [
        'buildername:XP Tests (1)', 'buildnumber:39442', 'cpu:x86-32',
        'data:123c706ddeeadcc761d612f732ddde9322400446', 'gpu:none',
        'master:chromium.win', 'name:unit_tests', 'os:Windows-XP-SP3',
        'pool:Chrome', 'priority:25', 'project:chromium', 'purpose:CI',
        'purpose:post-commit', 'stepname:unit_tests', 'user:'
    ],
    'outputs_ref': {
        'isolatedserver': 'https://isolateserver.appspot.com',
        'namespace': 'default-gzip',
        'isolated': 'isolatedhashunittests1'
    }
}]

_REF_REQUEST_WITH_COMMAND = {
    'expiration_secs': '3600',
    'name': 'ref_task_request',
    'parent_task_id': 'pti',
    'priority': '25',
    'properties': {
        'command': [
            'cmd',
            '--flag=value',
            '--gtest_filter=d.f',
            '--test-launcher-filter-file=path/to/filter/file',
            '--isolated-script-test-filter',
        ],
        'dimensions': [{
            'key': 'k',
            'value': 'v'
        }],
        'env': [
            {
                'key': 'a',
                'value': '1'
            },
            {
                'key': 'GTEST_SHARD_INDEX',
                'value': '1'
            },
            {
                'key': 'GTEST_TOTAL_SHARDS',
                'value': '5'
            },
        ],
        'execution_timeout_secs': '3600',
        'grace_period_secs': '30',
        'idempotent': True,
        'inputs_ref': {
            'isolatedserver': 'isolatedserver',
            'isolated': 'isolatedsha',
        },
        'io_timeout_secs': '1200',
    },
    'tags': ['master:%s' % 'b',
             'buildername:%s' % 'b', 'name:a_tests'],
    'user': 'user',
    'pubsub_topic': None,
    'pubsub_auth_token': None,
    'pubsub_userdata': None,
}

_REF_REQUEST_WITH_EXTRA_ARGS = {
    'expiration_secs': '3600',
    'name': 'ref_task_request',
    'parent_task_id': 'pti',
    'priority': '25',
    'properties': {
        'caches': [
            {
                'name': 'a',
                'path': '.foo'
            },
            {
                'name': 'task_template_vpython_cache',
                'path': '.task_template_vpython_cache/vpython'
            },
        ],
        'dimensions': [{
            'key': 'k',
            'value': 'v'
        }],
        'env': [
            {
                'key': 'a',
                'value': '1'
            },
            {
                'key': 'GTEST_SHARD_INDEX',
                'value': '1'
            },
            {
                'key': 'GTEST_TOTAL_SHARDS',
                'value': '5'
            },
        ],
        'env_prefixes': [
            {
                'key': 'a',
                'value': ['1']
            },
            {
                'key':
                    'PATH',
                'value': [
                    '.task_template_packages',
                    '.task_template_packages/cpython/bin',
                    '.task_template_packages/cpython3/bin'
                ]
            },
            {
                'key': 'VPYTHON_VIRTUALENV_ROOT',
                'value': ['.task_template_vpython_cache/vpython']
            },
        ],
        'execution_timeout_secs': '3600',
        'extra_args': [
            '--flag=value',
            '--gtest_filter=d.f',
            '--test-launcher-filter-file=path/to/filter/file',
            '--isolated-script-test-filter',
        ],
        'grace_period_secs': '30',
        'idempotent': True,
        'cas_input_root': {
            'cas_instance': 'cas_instance',
            'digest': {
                'size_bytes': '100',
                'hash': 'hash',
            },
        },
        'io_timeout_secs': '1200',
        'cipd_input': {
            'packages': [{
                'path': '.path',
                'version': '1',
                'package_name': 'foo',
            }, {
                'path': '.task_template_packages',
                'version': 'version',
                'package_name': 'package_name',
            }],
            'client_package': {
                'version': 'version',
                'package_name': 'package_name',
            },
            'server': 'server',
        },
    },
    'tags': ['master:%s' % 'b',
             'buildername:%s' % 'b', 'name:a_tests'],
    'user': 'user',
    'pubsub_topic': None,
    'pubsub_auth_token': None,
    'pubsub_userdata': None,
}


class SwarmingTest(wf_testcase.WaterfallTestCase):

  @mock.patch.object(swarming_util, 'ListTasks', return_value={})
  def testListSwarmingTasksDataByTags(self, _):
    self.assertEqual({}, swarming.ListSwarmingTasksDataByTags(None, 'b', 123))
    self.assertEqual({},
                     swarming.ListSwarmingTasksDataByTags(
                         None, 'b', 123, 'step'))

  @mock.patch.object(swarming, 'ListSwarmingTasksDataByTags')
  def testGetIsolatedShaForStep(self, mocked_list_swarming_tasks_data):
    builder_name = 'b'
    build_number = 123
    step_name = 's'
    isolated_sha = 'a1b2c3d4'

    mocked_list_swarming_tasks_data.return_value = [
        SwarmingTaskData({'tags': ['data:a1b2c3d4']})
    ]

    self.assertEqual(
        isolated_sha,
        swarming.GetIsolatedShaForStep(builder_name, build_number, step_name,
                                       None))

  @mock.patch.object(swarming, 'ListSwarmingTasksDataByTags', return_value=None)
  def testGetIsolatedShaForStepNoData(self, _):
    self.assertIsNone(swarming.GetIsolatedShaForStep('b', 123, 's', None))

  @mock.patch.object(swarming, 'ListSwarmingTasksDataByTags')
  def testGetIsolatedShaForStepNoShaFound(self,
                                          mocked_list_swarming_tasks_data):
    builder_name = 'b'
    build_number = 123
    step_name = 's'

    mocked_list_swarming_tasks_data.return_value = [
        SwarmingTaskData({'tags': []})
    ]

    self.assertIsNone(
        swarming.GetIsolatedShaForStep(builder_name, build_number, step_name,
                                       None))

  @mock.patch.object(swarming, 'ListSwarmingTasksDataByTags')
  def testGetIsolatedDataForStep(self, mock_fn):
    builder_name = 'b'
    build_number = 223
    step_name = 'unit_tests'

    mock_fn.return_value = [
        SwarmingTaskData(item) for item in _SAMPLE_BUILD_STEP_DATA
    ]

    data = swarming.GetIsolatedDataForStep(builder_name, build_number,
                                           step_name, None)
    expected_data = [{
        'digest':
            'isolatedhashunittests',
        'namespace':
            'default-gzip',
        'isolatedserver':
            waterfall_config.GetSwarmingSettings().get('isolated_server')
    }]
    self.assertEqual(expected_data, data)

  @mock.patch.object(swarming, 'ListSwarmingTasksDataByTags')
  def testGetIsolatedDataForStepNotOnlyFailure(self, mock_fn):
    builder_name = 'b'
    build_number = 223
    step_name = 'unit_tests'

    mock_fn.return_value = [
        SwarmingTaskData(item) for item in _SAMPLE_BUILD_STEP_DATA
    ]

    data = swarming.GetIsolatedDataForStep(
        builder_name, build_number, step_name, None, only_failure=False)
    expected_data = [{
        'digest':
            'isolatedhashunittests',
        'namespace':
            'default-gzip',
        'isolatedserver':
            waterfall_config.GetSwarmingSettings().get('isolated_server')
    }, {
        'digest':
            'isolatedhashunittests1',
        'namespace':
            'default-gzip',
        'isolatedserver':
            waterfall_config.GetSwarmingSettings().get('isolated_server')
    }]
    self.assertEqual(sorted(expected_data), sorted(data))

  @mock.patch.object(swarming, 'ListSwarmingTasksDataByTags', return_value=[])
  def testGetIsolatedDataForStepFailed(self, _):
    builder_name = 'download_failed'
    build_number = 223
    step_name = 's1'

    data = swarming.GetIsolatedDataForStep(builder_name, build_number,
                                           step_name, None)

    self.assertEqual([], data)

  @mock.patch.object(swarming, 'ListSwarmingTasksDataByTags')
  def testGetIsolatedDataForStepNoOutputsRef(self, mock_data):
    builder_name = 'download_failed'
    build_number = 223
    step_name = 's1'

    mock_data.return_value = [
        SwarmingTaskData({'failure': True}),
        SwarmingTaskData({'failure': False})
    ]

    data = swarming.GetIsolatedDataForStep(builder_name, build_number,
                                           step_name, None)
    expected_data = []

    self.assertEqual(expected_data, data)

  @mock.patch.object(swarming, 'ListSwarmingTasksDataByTags')
  def testUpdateSwarmingSteps(self, mock_data):
    builder_name = 'b'
    build_number = 223
    failed_steps = {
        'a_tests': {
            'current_failure': 2,
            'first_failure': 0,
            'supported': True
        },
        'unit_tests': {
            'current_failure': 2,
            'first_failure': 0,
            'supported': True
        },
        'compile': {
            'current_failure': 2,
            'first_failure': 0,
            'supported': True
        }
    }
    failed_steps = TestFailedSteps.FromSerializable(failed_steps)

    mock_data.return_value = [
        SwarmingTaskData({
            'failure': True,
            'internal_failure': False,
            'tags': [],
            'outputs_ref': {
                'isolatedserver': 'https://isolateserver.appspot.com',
                'namespace': 'default-gzip',
                'isolated': 'isolatedhashnetunittests'
            }
        }),
        SwarmingTaskData({
            'failure': False,
            'internal_failure': False,
            'tags': ['stepname:unit_tests'],
            'outputs_ref': {
                'isolatedserver': 'https://isolateserver.appspot.com',
                'namespace': 'default-gzip',
                'isolated': 'isolatedhashunittests'
            }
        }),
        SwarmingTaskData({
            'failure': True,
            'internal_failure': False,
            'tags': ['stepname:unit_tests'],
            'outputs_ref': {
                'isolatedserver': 'https://isolateserver.appspot.com',
                'namespace': 'default-gzip',
                'isolated': 'isolatedhashunittests1'
            }
        }),
        SwarmingTaskData({
            'failure': True,
            'internal_failure': False,
            'tags': ['stepname:a'],
            'outputs_ref': {
                'isolatedserver': 'https://isolateserver.appspot.com',
                'namespace': 'default-gzip',
                'isolated': 'isolatedhasha'
            }
        }),
        SwarmingTaskData({
            'failure': True,
            'internal_failure': False,
            'tags': ['stepname:a_tests'],
            'outputs_ref': {
                'isolatedserver': 'https://isolateserver.appspot.com',
                'namespace': 'default-gzip',
                'isolated': 'isolatedhashatests'
            }
        }),
        SwarmingTaskData({
            'failure': True,
            'internal_failure': False,
            'tags': ['stepname:abc_test'],
            'outputs_ref': {
                'isolatedserver': 'https://isolateserver.appspot.com',
                'namespace': 'default-gzip',
                'isolated': 'isolatedhashabctest-223'
            }
        }),
        SwarmingTaskData({
            'failure': True,
            'internal_failure': True
        })
    ]

    expected_isoloated_data = {
        'a_tests': [{
            'isolatedserver': 'https://isolateserver.appspot.com',
            'namespace': 'default-gzip',
            'digest': 'isolatedhashatests'
        }],
        'unit_tests': [{
            'isolatedserver': 'https://isolateserver.appspot.com',
            'namespace': 'default-gzip',
            'digest': 'isolatedhashunittests1'
        }]
    }

    isolated_data = swarming.GetIsolatedDataForFailedStepsInABuild(
        builder_name, build_number, failed_steps, None)

    self.assertEqual(expected_isoloated_data, isolated_data)

  @mock.patch.object(swarming, 'ListSwarmingTasksDataByTags', return_value=[])
  def testGetIsolatedDataForFailedStepsInABuildNoData(self, _):
    self.assertEqual({},
                     swarming.GetIsolatedDataForFailedStepsInABuild(
                         'b', 1, [], None))

  @mock.patch.object(swarming, 'ListSwarmingTasksDataByTags', return_value=[])
  def testGetReferredSwarmingTaskRequestInfoNoTaskFound(self, _):
    with self.assertRaises(Exception):
      swarming.GetReferredSwarmingTaskRequestInfo('m', 'b', 123, 's', None)

  @mock.patch.object(
      swarming,
      'ListSwarmingTasksDataByTags',
      return_value=[SwarmingTaskData({'task_id': 'task_id'})])
  @mock.patch.object(swarming_util, 'GetSwarmingTaskRequest')
  def testGetReferredSwarmingTaskRequestInfo(self, mock_get, _):
    request = SwarmingTaskRequest.FromSerializable(_REF_REQUEST_WITH_COMMAND)
    mock_get.return_value = request
    task_id, ref_request = swarming.GetReferredSwarmingTaskRequestInfo(
        'm', 'b', 123, 's', None)
    self.assertEqual('task_id', task_id)
    self.assertEqual(request, ref_request)

  @mock.patch.object(app_identity, 'get_application_id', return_value='app-id')
  @mock.patch.object(token, 'GenerateAuthToken', return_value='auth_token')
  @mock.patch.object(
      time_util, 'GetUTCNow', return_value=datetime(2018, 03, 15, 0, 0, 0))
  def testCreateNewSwarmingTaskRequestTemplateWithCommand(self, *_):
    ref_task_id = 'ref_task_id'
    master_name = 'm'
    builder_name = 'b'
    step_name = 'a_tests on platform'
    tests = ['a.b', 'a.c']
    iterations = 100

    new_request = swarming.CreateNewSwarmingTaskRequestTemplate(
        'runner_id', ref_task_id,
        SwarmingTaskRequest.FromSerializable(_REF_REQUEST_WITH_COMMAND),
        master_name, builder_name, step_name, tests, iterations)

    expected_new_request_json = {
        'expiration_secs': '72000',
        'name': 'findit/ref_task_id/ref_task_id/2018-03-15 00:00:00 000000',
        'parent_task_id': '',
        'priority': '150',
        'properties': {
            'command': [
                'cmd', '--flag=value', '--isolated-script-test-filter=a.b::a.c',
                '--isolated-script-test-repeat=%d' % iterations,
                '--isolated-script-test-launcher-retry-limit=0',
                '--isolated-script-test-also-run-disabled-tests'
            ],
            'dimensions': [{
                'key': 'k',
                'value': 'v'
            }],
            'env': [{
                'key': 'a',
                'value': '1'
            },],
            'execution_timeout_secs': '3600',
            'grace_period_secs': '30',
            'idempotent': False,
            'inputs_ref': {
                'isolatedserver': 'isolatedserver',
                'isolated': 'isolatedsha',
            },
            'io_timeout_secs': '1200',
        },
        'tags': [
            'ref_master:%s' % master_name,
            'ref_buildername:%s' % builder_name,
            'ref_stepname:%s' % step_name, 'ref_name:a_tests', 'findit:1',
            'project:Chromium', 'purpose:post-commit'
        ],
        'user': '',
        'pubsub_auth_token': 'auth_token',
        'pubsub_topic': 'projects/app-id/topics/swarming',
        'pubsub_userdata': json.dumps({'runner_id': 'runner_id'}),
    }

    self.assertEqual(
        SwarmingTaskRequest.FromSerializable(expected_new_request_json),
        new_request)

  @mock.patch.object(app_identity, 'get_application_id', return_value='app-id')
  @mock.patch.object(token, 'GenerateAuthToken', return_value='auth_token')
  @mock.patch.object(
      time_util, 'GetUTCNow', return_value=datetime(2018, 03, 15, 0, 0, 0))
  def testCreateNewSwarmingTaskRequestTemplateWithExtraArgs(self, *_):
    ref_task_id = 'ref_task_id'
    master_name = 'm'
    builder_name = 'b'
    step_name = 'a_tests on platform'
    tests = ['a.b', 'a.c']
    iterations = 100

    new_request = swarming.CreateNewSwarmingTaskRequestTemplate(
        'runner_id', ref_task_id,
        SwarmingTaskRequest.FromSerializable(_REF_REQUEST_WITH_EXTRA_ARGS),
        master_name, builder_name, step_name, tests, iterations)

    expected_new_request_json = {
        'expiration_secs': '72000',
        'name': 'findit/ref_task_id/ref_task_id/2018-03-15 00:00:00 000000',
        'parent_task_id': '',
        'priority': '150',
        'properties': {
            'caches': [{
                'name': 'a',
                'path': '.foo'
            }],
            'dimensions': [{
                'key': 'k',
                'value': 'v'
            }],
            'env': [{
                'key': 'a',
                'value': '1'
            },],
            'env_prefixes': [{
                'key': 'a',
                'value': ['1']
            },],
            'execution_timeout_secs': '3600',
            'extra_args': [
                '--flag=value',
                '--isolated-script-test-filter=a.b::a.c',
                '--isolated-script-test-repeat=%d' % iterations,
                '--isolated-script-test-launcher-retry-limit=0',
                '--isolated-script-test-also-run-disabled-tests',
            ],
            'grace_period_secs': '30',
            'idempotent': False,
            'cas_input_root': {
                'cas_instance': 'cas_instance',
                'digest': {
                    'size_bytes': '100',
                    'hash': 'hash',
                },
            },
            'io_timeout_secs': '1200',
            'cipd_input': {
                'packages': [{
                    'path': '.path',
                    'version': '1',
                    'package_name': 'foo',
                }],
                'client_package': {
                    'version': 'version',
                    'package_name': 'package_name',
                },
                'server': 'server',
            },
        },
        'tags': [
            'ref_master:%s' % master_name,
            'ref_buildername:%s' % builder_name,
            'ref_stepname:%s' % step_name, 'ref_name:a_tests', 'findit:1',
            'project:Chromium', 'purpose:post-commit'
        ],
        'user': '',
        'pubsub_auth_token': 'auth_token',
        'pubsub_topic': 'projects/app-id/topics/swarming',
        'pubsub_userdata': json.dumps({'runner_id': 'runner_id'}),
    }

    self.assertEqual(
        SwarmingTaskRequest.FromSerializable(expected_new_request_json),
        new_request)

  @mock.patch.object(swarming, 'ListSwarmingTasksDataByTags', return_value=[])
  def testCanFindSwarmingTaskFromBuildForAStep(self, _):
    self.assertFalse(
        swarming.CanFindSwarmingTaskFromBuildForAStep(None, 'b', 1, 's'))

  def testGetSwarmingTaskUrl(self):
    task_id = 'task_id'
    self.assertIn(task_id, swarming.GetSwarmingTaskUrl(task_id))

  @mock.patch.object(swarming, 'ListSwarmingTasksDataByTags')
  def testGetSwarmingTaskIdsForFailedSteps(self, mock_fn):
    builder_name = 'b'
    build_number = 223
    failed_steps = {'unit_tests': ''}

    mock_fn.return_value = [
        SwarmingTaskData(item) for item in _SAMPLE_BUILD_STEP_DATA
    ]

    data = swarming.GetSwarmingTaskIdsForFailedSteps(builder_name, build_number,
                                                     failed_steps, None)
    expected_data = {"unit_tests": ["2944afa502297110"]}
    self.assertEqual(expected_data, data)
