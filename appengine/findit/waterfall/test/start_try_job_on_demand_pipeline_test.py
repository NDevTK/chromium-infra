# Copyright 2016 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

import mock

from common.waterfall import failure_type
from model.wf_try_job import WfTryJob
from waterfall import start_try_job_on_demand_pipeline
from waterfall.start_try_job_on_demand_pipeline import (
    StartTryJobOnDemandPipeline)
from waterfall.test import wf_testcase

class StartTryJobOnDemandPipelineTest(wf_testcase.WaterfallTestCase):

  def testGetLastPassCurrentBuildIsNotFirstFailure(self):
    failure_info = {
        'failed_steps': {
          'a': {
              'first_failure': 1,
              'last_pass': 0
          }
        }
    }
    last_pass = start_try_job_on_demand_pipeline._GetLastPass(
      2, failure_info, failure_type.COMPILE)
    self.assertIsNone(last_pass)

  def testGetLastPassUnknownType(self):
    last_pass = start_try_job_on_demand_pipeline._GetLastPass(
      1, {}, failure_type.UNKNOWN)
    self.assertIsNone(last_pass)

  def testGetLastPassTestNoLastPass(self):
    try_job_type = failure_type.TEST
    failure_info = {
        'failure_type': try_job_type,
        'builds': {
            '0': {
                'blame_list': ['r0', 'r1'],
                'chromium_revision': 'r1'
            },
            '1': {
                'blame_list': ['r2'],
                'chromium_revision': 'r2'
            }
        },
        'failed_steps': {
            'a': {
                'first_failure': 1,
                'last_pass': 0,
                'tests': {
                    'test1': {
                        'first_failure': 1
                    }
                }
            }
        }
    }
    last_pass = start_try_job_on_demand_pipeline._GetLastPass(
      1, failure_info, try_job_type)
    self.assertIsNone(last_pass)

  def testGetSuspectsFromHeuristicResult(self):
    heuristic_result = {
        'failures': [
            {
                'step_name': 'compile',
                'suspected_cls': [
                    {
                        'revision': 'r1',
                    },
                    {
                        'revision': 'r2',
                    },
                ],
            },
        ]
    }
    expected_suspected_revisions = ['r1', 'r2']
    self.assertEqual(
        expected_suspected_revisions,
        start_try_job_on_demand_pipeline._GetSuspectsFromHeuristicResult(
            heuristic_result))

  def testNotScheduleTryJobIfBuildNotCompleted(self):
    pipeline = start_try_job_on_demand_pipeline.StartTryJobOnDemandPipeline()
    result = pipeline.run(
        'm', 'b', 1, {}, {}, {}, False, False)
    self.assertEqual(list(result), [])

  @mock.patch.object(start_try_job_on_demand_pipeline, 'try_job_util')
  def testNotScheduleTryJobIfDontNeedTryJob(self, mock_module):
    mock_module.NeedANewTryJob.return_value = False
    pipeline = start_try_job_on_demand_pipeline.StartTryJobOnDemandPipeline()
    result = pipeline.run(
        'm', 'b', 1, {}, {}, {}, True, False)
    self.assertEqual(list(result), [])

  @mock.patch.object(start_try_job_on_demand_pipeline, 'try_job_util')
  def testNotScheduleTryJobIfUnsupportedFailureType(self, mock_module):
    mock_module.NeedANewTryJob.return_value = True
    try_job_type = failure_type.UNKNOWN
    failure_info = {
        'failure_type': try_job_type,
        'builds': {
            '0': {
                'blame_list': ['r0', 'r1'],
                'chromium_revision': 'r1'
            },
            '1': {
                'blame_list': ['r2'],
                'chromium_revision': 'r2'
            }
        },
        'failed_steps': {
            'a': {
                'first_failure': 1,
                'last_pass': 0
            }
        }
    }
    pipeline = start_try_job_on_demand_pipeline.StartTryJobOnDemandPipeline()
    result = pipeline.run(
        'm', 'b', 1, failure_info, {}, {}, True, False)
    self.assertEqual(list(result), [])

  @mock.patch.object(start_try_job_on_demand_pipeline, 'try_job_util')
  def testCompileTryJob(self, mock_module):

    master_name = 'm'
    builder_name = 'b'
    build_number = 1
    try_job_type = failure_type.COMPILE
    failure_info = {
        'failure_type': try_job_type,
        'builds': {
            '0': {
                'blame_list': ['r0', 'r1'],
                'chromium_revision': 'r1'
            },
            '1': {
                'blame_list': ['r2'],
                'chromium_revision': 'r2'
            }
        },
        'failed_steps': {
          'compile': {
              'first_failure': 1,
              'last_pass': 0
          }
        }
    }
    good_revision = 'r1'
    bad_revision = 'r2'
    WfTryJob.Create('m', 'b', 1).put()

    mock_module.NeedANewTryJob.return_value = True
    mock_module.GetFailedTargetsFromSignals.return_value = {}

    self.MockPipeline(
        start_try_job_on_demand_pipeline.ScheduleCompileTryJobPipeline,
        'try_job_id',
        expected_args=[
            master_name, builder_name, build_number, good_revision,
           bad_revision, try_job_type, {}, []],
        expected_kwargs={})
    self.MockPipeline(
        start_try_job_on_demand_pipeline.MonitorTryJobPipeline,
        'try_job_result',
        expected_args=[
            master_name, builder_name, build_number, try_job_type, 'try_job_id'
        ],
        expected_kwargs={})
    self.MockPipeline(
        start_try_job_on_demand_pipeline.IdentifyTryJobCulpritPipeline,
        'final_result',
        expected_args=[
            master_name, builder_name, build_number, ['r2'], try_job_type,
            'try_job_id', 'try_job_result'],
        expected_kwargs={})

    pipeline = StartTryJobOnDemandPipeline()
    result = pipeline.run(
        'm', 'b', 1, failure_info, {}, {}, True, False)
    self.assertNotEqual(list(result), [])

  @mock.patch.object(start_try_job_on_demand_pipeline, 'try_job_util')
  def testTestTryJob(self, mock_module):

    master_name = 'm'
    builder_name = 'b'
    build_number = 1
    try_job_type = failure_type.TEST
    failure_info = {
        'failure_type': try_job_type,
        'builds': {
            '0': {
                'blame_list': ['r0', 'r1'],
                'chromium_revision': 'r1'
            },
            '1': {
                'blame_list': ['r2'],
                'chromium_revision': 'r2'
            }
        },
        'failed_steps': {
            'a': {
                'first_failure': 1,
                'last_pass': 0,
                'tests': {
                    'test1': {
                        'first_failure': 1,
                        'last_pass': 0
                    },
                    'test2': {
                        'first_failure': 0
                    }
                }
            },
            'b': {
                'first_failure': 0,
                'tests': {
                    'b_test1': {
                        'first_failure': 0
                    }
                }
            }
        }
    }
    good_revision = 'r1'
    bad_revision = 'r2'

    mock_module.NeedANewTryJob.return_value = True
    mock_module.GetFailedTargetsFromSignals.return_value = {}

    self.MockPipeline(
        start_try_job_on_demand_pipeline.ScheduleTestTryJobPipeline,
        'try_job_id',
        expected_args=[
            master_name, builder_name, build_number, good_revision,
           bad_revision, try_job_type, 'targeted_tests', []],
        expected_kwargs={})
    self.MockPipeline(
        start_try_job_on_demand_pipeline.MonitorTryJobPipeline,
        'try_job_result',
        expected_args=[
            master_name, builder_name, build_number, try_job_type, 'try_job_id'
        ],
        expected_kwargs={})
    self.MockPipeline(
        start_try_job_on_demand_pipeline.IdentifyTryJobCulpritPipeline,
        'final_result',
        expected_args=[
            master_name, builder_name, build_number, ['r2'], try_job_type,
            'try_job_id', 'try_job_result'],
        expected_kwargs={})

    pipeline = StartTryJobOnDemandPipeline()
    result = pipeline.run(
        'm', 'b', 1, failure_info, {}, {}, True, False)
    WfTryJob.Create('m', 'b', 1).put()
    self.assertNotEqual(list(result), [])
