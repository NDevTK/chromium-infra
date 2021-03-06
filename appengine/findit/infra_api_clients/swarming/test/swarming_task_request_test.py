# Copyright 2018 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

from gae_libs.testcase import TestCase
from infra_api_clients.swarming import swarming_task_request
from infra_api_clients.swarming.swarming_task_request import (
    SwarmingTaskInputsRef, SwarmingTaskProperties, SwarmingTaskRequest,
    SwarmingResultDBCfg)
from libs.list_of_basestring import ListOfBasestring


class SwarmingTaskRequestTest(TestCase):

  def testGetSwarmingTaskRequestTemplate(self):
    expected_request = SwarmingTaskRequest(
        created_ts=None,
        expiration_secs='3600',
        name='',
        parent_task_id='',
        priority='150',
        properties=SwarmingTaskProperties(
            caches=[],
            command=ListOfBasestring(),
            relative_cwd=None,
            dimensions=[],
            env=[],
            env_prefixes=[],
            execution_timeout_secs='3600',
            extra_args=ListOfBasestring(),
            grace_period_secs='30',
            io_timeout_secs='1200',
            idempotent=True,
            inputs_ref=SwarmingTaskInputsRef(
                isolated=None, isolatedserver=None, namespace=None),
            cipd_input=swarming_task_request.CIPDInput(
                packages=swarming_task_request.CIPDPackages(),
                client_package=swarming_task_request.CIPDClientPackage(
                    version=None,
                    package_name=None,
                ),
                server=None),
            cas_input_root=swarming_task_request.CASInputRoot(
                cas_instance=None,
                digest=swarming_task_request.CASInputRootDigest(
                    size_bytes=None,
                    hash=None,
                ),
            ),
        ),
        pubsub_auth_token=None,
        pubsub_topic=None,
        pubsub_userdata=None,
        service_account=None,
        tags=ListOfBasestring(),
        user='',
        realm=None,
        resultdb=SwarmingResultDBCfg(enable=False))

    self.assertEqual(expected_request,
                     SwarmingTaskRequest.GetSwarmingTaskRequestTemplate())

  def testFromSerializable(self):
    cipd_packages = [{
        'path': 'path',
        'version': 'version',
        'package_name': 'package_name',
    }]
    data = {
        'expiration_secs': '50',
        'name': 'a swarming task',
        'parent_task_id': 'parent task id',
        'priority': '150',
        'tags': ['a'],
        'user': 'someone',
        'some_unused_field': 'blabla',
        'pubsub_topic': 'topic',
        'pubsub_auth_token': 'token',
        'pubsub_userdata': 'data',
        'realm': 'foo:realm',
        'resultdb': {
            'enable': True
        },
        'properties': {
            'caches': [{
                'name': 'a',
                'path': '.foo'
            },],
            'command': ['path/to/binary'],
            'relative_cwd': 'out/Release',
            'unused_property': 'blabla',
            'dimensions': [{
                'key': 'cpu',
                'value': 'x86-64',
            },],
            'env': [{
                'key': 'name',
                'value': '1',
            },],
            'execution_timeout_secs': 10,
            'grace_period_secs': 5,
            'extra_args': ['--arg=value'],
            'idempotent': True,
            'inputs_ref': {
                'namespace': 'default-gzip',
                'isolated': 'a-hash',
                'random_field': 'blabla'
            },
            'io_timeout_secs': 10,
            'cipd_input': {
                'packages': cipd_packages,
                'client_package': {
                    'version': 'version',
                    'package_name': 'package_name',
                },
                'server': 'server',
            },
            'cas_input_root': {
                'cas_instance': 'projects/chromium-swarm/instances/default',
                'digest': {
                    'size_bytes': '1135',
                    'hash': 'bec400ea1f34229085710fcc7ff174147f4f4af4db',
                }
            }
        },
    }

    expected_request = SwarmingTaskRequest(
        created_ts=None,
        expiration_secs='50',
        name='a swarming task',
        parent_task_id='parent task id',
        priority='150',
        properties=SwarmingTaskProperties(
            caches=[{
                'name': 'a',
                'path': '.foo'
            }],
            command=ListOfBasestring.FromSerializable(['path/to/binary']),
            relative_cwd='out/Release',
            dimensions=[
                {
                    'key': 'cpu',
                    'value': 'x86-64',
                },
            ],
            env=[
                {
                    'key': 'name',
                    'value': '1',
                },
            ],
            env_prefixes=[],
            execution_timeout_secs='10',
            extra_args=ListOfBasestring.FromSerializable(['--arg=value']),
            grace_period_secs='5',
            io_timeout_secs='10',
            idempotent=True,
            inputs_ref=SwarmingTaskInputsRef(
                isolated='a-hash',
                isolatedserver=None,
                namespace='default-gzip'),
            cipd_input=swarming_task_request.CIPDInput(
                packages=swarming_task_request.CIPDPackages.FromSerializable(
                    cipd_packages),
                client_package=swarming_task_request.CIPDClientPackage(
                    version='version', package_name='package_name'),
                server='server'),
            cas_input_root=swarming_task_request.CASInputRoot(
                cas_instance='projects/chromium-swarm/instances/default',
                digest=swarming_task_request.CASInputRootDigest(
                    size_bytes='1135',
                    hash='bec400ea1f34229085710fcc7ff174147f4f4af4db',
                ),
            )),
        pubsub_auth_token='token',
        pubsub_topic='topic',
        pubsub_userdata='data',
        service_account=None,
        tags=ListOfBasestring.FromSerializable(['a']),
        user='someone',
        realm='foo:realm',
        resultdb=SwarmingResultDBCfg(enable=True))

    self.assertEqual(expected_request,
                     SwarmingTaskRequest.FromSerializable(data))
