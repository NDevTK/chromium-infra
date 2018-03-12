# Copyright 2017 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.
import datetime
import mock

from common import constants
from gae_libs.gitiles.cached_gitiles_repository import CachedGitilesRepository
from gae_libs.pipelines import pipeline
from gae_libs.pipelines import pipeline_handlers
from libs import analysis_status
from model.flake.flake_swarming_task import FlakeSwarmingTask
from model.flake.master_flake_analysis import DataPoint
from model.flake.master_flake_analysis import MasterFlakeAnalysis
from waterfall import build_util
from waterfall.build_info import BuildInfo
from waterfall.flake import flake_constants
from waterfall.flake import update_flake_analysis_data_points_pipeline
from waterfall.flake.update_flake_analysis_data_points_pipeline import (
    UpdateFlakeAnalysisDataPointsPipeline)
from waterfall.test import wf_testcase


class UpdateFlakeAnalysisDataPointsPipelineTest(wf_testcase.WaterfallTestCase):

  app_module = pipeline_handlers._APP

  @mock.patch.object(
      CachedGitilesRepository,
      'GetCommitsBetweenRevisions',
      return_value=['r4', 'r3', 'r2', 'r1'])
  def testGetCommitsBetweenRevisions(self, _):
    self.assertEqual(
        update_flake_analysis_data_points_pipeline._GetCommitsBetweenRevisions(
            'r0', 'r4'), ['r1', 'r2', 'r3', 'r4'])

  def testGetPassRateTestDoesNotExist(self):
    task = FlakeSwarmingTask.Create('m', 'b', 123, 's', 't')
    task.tries = 0
    self.assertEqual(
        flake_constants.PASS_RATE_TEST_NOT_FOUND,
        update_flake_analysis_data_points_pipeline._GetPassRate(task))

  def testGetPassRate(self):
    task = FlakeSwarmingTask.Create('m', 'b', 123, 's', 't')
    task.successes = 50
    task.tries = 100
    self.assertEqual(
        0.5, update_flake_analysis_data_points_pipeline._GetPassRate(task))

  @mock.patch.object(update_flake_analysis_data_points_pipeline,
                     '_GetCommitsBetweenRevisions')
  @mock.patch.object(build_util, 'GetBuildInfo')
  def testUpdateAnalysisDataPointsForSwarmingTaskWithPrevious(
      self, mocked_build_info, mocked_commits):
    master_name = 'm'
    builder_name = 'b'
    build_number = 123
    step_name = 's'
    test_name = 't'
    task_id = 'task_id'
    has_valid_artifact = True
    tries = 100
    successes = 50
    chromium_revision = 'r1000'
    commit_position = 1000
    blame_list = [
        'r1000', 'r999', 'r998', 'r997', 'r996', 'r995', 'r994', 'r993', 'r992',
        'r991'
    ]

    analysis = MasterFlakeAnalysis.Create(master_name, builder_name,
                                          build_number, step_name, test_name)
    analysis.put()

    task = FlakeSwarmingTask.Create(master_name, builder_name, build_number,
                                    step_name, test_name)
    task.task_id = task_id
    task.has_valid_artifact = has_valid_artifact
    task.tries = tries
    task.successes = successes
    task.started_time = datetime.datetime(1, 1, 1, 0, 0)
    task.completed_time = datetime.datetime(1, 1, 1, 0, 1)

    build_info_122 = BuildInfo(master_name, builder_name, build_number)
    build_info_122.commit_position = commit_position - 10
    build_info_122.chromium_revision = 'r990'
    build_info_123 = BuildInfo(master_name, builder_name, build_number)
    build_info_123.commit_position = commit_position
    build_info_123.chromium_revision = chromium_revision

    mocked_build_info.side_effect = [(200, build_info_123), (200,
                                                             build_info_122)]
    mocked_commits.return_value = blame_list

    (update_flake_analysis_data_points_pipeline.
     _UpdateAnalysisDataPointsWithSwarmingTask(task, analysis))
    data_point = analysis.FindMatchingDataPointWithBuildNumber(build_number)

    self.assertEqual(build_number, data_point.build_number)
    self.assertEqual(0.5, data_point.pass_rate)
    self.assertEqual(commit_position, data_point.commit_position)
    self.assertEqual(chromium_revision, data_point.git_hash)
    self.assertEqual('r990', data_point.previous_build_git_hash)
    self.assertEqual(990, data_point.previous_build_commit_position)
    self.assertEqual(blame_list, data_point.blame_list)
    self.assertEqual(60, data_point.elapsed_seconds)

  @mock.patch.object(build_util, 'GetBuildInfo')
  def testUpdateAnalysisDataPointsForSwarmingTaskExistingDataPoint(
      self, mocked_build_info):
    master_name = 'm'
    builder_name = 'b'
    build_number = 0
    step_name = 's'
    test_name = 't'
    task_id = 'task_id'
    has_valid_artifact = True
    tries = 100
    successes = 50
    chromium_revision = 'r1000'
    commit_position = 1000
    blame_list = [
        'r1000', 'r999', 'r998', 'r997', 'r996', 'r995', 'r994', 'r993', 'r992',
        'r991'
    ]

    analysis = MasterFlakeAnalysis.Create(master_name, builder_name,
                                          build_number, step_name, test_name)
    analysis.data_points = [
        DataPoint.Create(
            build_number=build_number,
            pass_rate=1.0,
            task_ids=['t1', 't2'],
            git_hash=chromium_revision,
            commit_position=commit_position,
            blame_list=blame_list,
            iterations=100,
            elapsed_seconds=100)
    ]
    analysis.put()

    task = FlakeSwarmingTask.Create(master_name, builder_name, build_number,
                                    step_name, test_name)
    task.task_id = task_id
    task.has_valid_artifact = has_valid_artifact
    task.tries = tries
    task.successes = successes
    task.started_time = datetime.datetime(1, 1, 1, 0, 0)
    task.completed_time = datetime.datetime(1, 1, 1, 0, 2)

    build_info = BuildInfo(master_name, builder_name, build_number)
    build_info.commit_position = commit_position
    build_info.chromium_revision = chromium_revision
    build_info.blame_list = blame_list

    mocked_build_info.return_value = (200, build_info)

    (update_flake_analysis_data_points_pipeline.
     _UpdateAnalysisDataPointsWithSwarmingTask(task, analysis))
    data_point = analysis.FindMatchingDataPointWithCommitPosition(
        commit_position)

    self.assertEqual(build_number, data_point.build_number)
    self.assertEqual(0.75, data_point.pass_rate)
    self.assertEqual(commit_position, data_point.commit_position)
    self.assertEqual(chromium_revision, data_point.git_hash)
    self.assertIsNone(data_point.previous_build_git_hash)
    self.assertIsNone(data_point.previous_build_commit_position)
    self.assertEqual(blame_list, data_point.blame_list)
    self.assertEqual(220, data_point.elapsed_seconds)

  @mock.patch.object(build_util, 'GetBuildInfo')
  def testUpdateAnalysisDataPointsForSwarmingTaskNoPreviousBuild(
      self, mocked_build_info):
    master_name = 'm'
    builder_name = 'b'
    build_number = 0
    step_name = 's'
    test_name = 't'
    task_id = 'task_id'
    has_valid_artifact = True
    tries = 100
    successes = 50
    chromium_revision = 'r1000'
    commit_position = 1000
    blame_list = [
        'r1000', 'r999', 'r998', 'r997', 'r996', 'r995', 'r994', 'r993', 'r992',
        'r991'
    ]

    analysis = MasterFlakeAnalysis.Create(master_name, builder_name,
                                          build_number, step_name, test_name)
    analysis.put()

    task = FlakeSwarmingTask.Create(master_name, builder_name, build_number,
                                    step_name, test_name)
    task.task_id = task_id
    task.has_valid_artifact = has_valid_artifact
    task.tries = tries
    task.successes = successes
    task.started_time = datetime.datetime(1, 1, 1, 0, 0)
    task.completed_time = datetime.datetime(1, 1, 1, 0, 2)

    build_info = BuildInfo(master_name, builder_name, build_number)
    build_info.commit_position = commit_position
    build_info.chromium_revision = chromium_revision
    build_info.blame_list = blame_list

    mocked_build_info.return_value = (200, build_info)

    (update_flake_analysis_data_points_pipeline.
     _UpdateAnalysisDataPointsWithSwarmingTask(task, analysis))
    data_point = analysis.FindMatchingDataPointWithCommitPosition(
        commit_position)

    self.assertEqual(build_number, data_point.build_number)
    self.assertEqual(0.5, data_point.pass_rate)
    self.assertEqual(commit_position, data_point.commit_position)
    self.assertEqual(chromium_revision, data_point.git_hash)
    self.assertIsNone(data_point.previous_build_git_hash)
    self.assertIsNone(data_point.previous_build_commit_position)
    self.assertEqual(blame_list, data_point.blame_list)
    self.assertEqual(120, data_point.elapsed_seconds)

  @mock.patch.object(build_util, 'GetBuildInfo', return_value=(404, None))
  def testCreateDataPointNoBuildInfo(self, _):
    master_name = 'm'
    builder_name = 'b'
    build_number = 0
    step_name = 's'
    test_name = 't'
    task_id = 'task_id'
    has_valid_artifact = True
    tries = 100
    successes = 50

    analysis = MasterFlakeAnalysis.Create(master_name, builder_name,
                                          build_number, step_name, test_name)
    analysis.put()

    task = FlakeSwarmingTask.Create(master_name, builder_name, build_number,
                                    step_name, test_name)
    task.task_id = task_id
    task.has_valid_artifact = has_valid_artifact
    task.tries = tries
    task.successes = successes
    task.started_time = datetime.datetime(1, 1, 1, 0, 0)
    task.completed_time = datetime.datetime(1, 1, 1, 0, 1)

    with self.assertRaises(pipeline.Retry):
      (update_flake_analysis_data_points_pipeline.
       _UpdateAnalysisDataPointsWithSwarmingTask(task, analysis))

  @mock.patch.object(build_util, 'GetBuildInfo')
  def testCreateDataPointNoPreviousBuildInfo(self, mocked_build_info):
    master_name = 'm'
    builder_name = 'b'
    build_number = 1
    step_name = 's'
    test_name = 't'
    task_id = 'task_id'
    has_valid_artifact = True
    tries = 100
    successes = 50
    chromium_revision = 'r1000'
    commit_position = 1000
    blame_list = [
        'r1000', 'r999', 'r998', 'r997', 'r996', 'r995', 'r994', 'r993', 'r992',
        'r991'
    ]

    analysis = MasterFlakeAnalysis.Create(master_name, builder_name,
                                          build_number, step_name, test_name)
    analysis.put()

    task = FlakeSwarmingTask.Create(master_name, builder_name, build_number,
                                    step_name, test_name)
    task.task_id = task_id
    task.has_valid_artifact = has_valid_artifact
    task.tries = tries
    task.successes = successes
    task.started_time = datetime.datetime(1, 1, 1, 0, 0)
    task.completed_time = datetime.datetime(1, 1, 1, 0, 1)

    build_info = BuildInfo(master_name, builder_name, build_number)
    build_info.commit_position = commit_position
    build_info.chromium_revision = chromium_revision
    build_info.blame_list = blame_list

    mocked_build_info.side_effect = [(200, build_info), (404, None)]

    with self.assertRaises(pipeline.Retry):
      (update_flake_analysis_data_points_pipeline.
       _UpdateAnalysisDataPointsWithSwarmingTask(task, analysis))

  @mock.patch.object(update_flake_analysis_data_points_pipeline,
                     '_UpdateAnalysisDataPointsWithSwarmingTask')
  def testUpdateFlakeAnalysisDataPointsPipeline(self, mocked_create_data_point):
    master_name = 'm'
    builder_name = 'b'
    build_number = 123
    step_name = 's'
    test_name = 't'

    tries = 100
    successes = 50
    task_id = 'task_id'
    has_valid_artifact = True
    commit_position = 1000
    git_hash = 'r1000'
    previous_build_commit_position = 990
    previous_build_git_hash = 'r990'

    analysis = MasterFlakeAnalysis.Create(master_name, builder_name,
                                          build_number, step_name, test_name)
    analysis.put()

    flake_swarming_task = FlakeSwarmingTask.Create(
        master_name, builder_name, build_number, step_name, test_name)
    flake_swarming_task.tries = tries
    flake_swarming_task.successes = successes
    flake_swarming_task.put()

    expected_data_point = DataPoint.Create(
        pass_rate=0.5,
        build_number=build_number,
        task_ids=[task_id],
        commit_position=commit_position,
        git_hash=git_hash,
        previous_build_git_hash=previous_build_git_hash,
        previous_build_commit_position=previous_build_commit_position,
        has_valid_artifact=has_valid_artifact)

    def create_data_point_fn(_, analysis):
      analysis.data_points.append(expected_data_point)
      analysis.put()

    mocked_create_data_point.side_effect = create_data_point_fn

    pipeline_job = UpdateFlakeAnalysisDataPointsPipeline(
        analysis.key.urlsafe(), build_number)

    pipeline_job.start(queue_name=constants.DEFAULT_QUEUE)
    self.execute_queued_tasks()

    self.assertEqual(len(analysis.data_points), 1)
    self.assertIn(expected_data_point, analysis.data_points)

  def testUpdateFlakeAnalysisDataPointsPipelineWithFailedSwarmingTask(self):
    master_name = 'm'
    builder_name = 'b'
    build_number = 123
    step_name = 's'
    test_name = 't'

    analysis = MasterFlakeAnalysis.Create(master_name, builder_name,
                                          build_number, step_name, test_name)
    analysis.put()

    flake_swarming_task = FlakeSwarmingTask.Create(
        master_name, builder_name, build_number, step_name, test_name)
    flake_swarming_task.status = analysis_status.ERROR
    flake_swarming_task.put()

    pipeline_job = UpdateFlakeAnalysisDataPointsPipeline(
        analysis.key.urlsafe(), build_number)

    pipeline_job.start(queue_name=constants.DEFAULT_QUEUE)
    self.execute_queued_tasks()

    self.assertEqual([], analysis.data_points)

  @mock.patch.object(update_flake_analysis_data_points_pipeline,
                     '_UpdateAnalysisDataPointsWithSwarmingTask')
  def testUpdateFlakeAnalysisDataPointsPipelineWithSalvagedSwarmingTask(
      self, mock_update_fn):
    master_name = 'm'
    builder_name = 'b'
    build_number = 0
    step_name = 's'
    test_name = 't'
    task_id = 'task_id'
    has_valid_artifact = True
    tries = 100
    successes = 50

    analysis = MasterFlakeAnalysis.Create(master_name, builder_name,
                                          build_number, step_name, test_name)
    analysis.put()

    task = FlakeSwarmingTask.Create(master_name, builder_name, build_number,
                                    step_name, test_name)
    task.status = analysis_status.ERROR
    task.task_id = task_id
    task.has_valid_artifact = has_valid_artifact
    task.tries = tries
    task.successes = successes
    task.started_time = datetime.datetime(1, 1, 1, 0, 0)
    task.completed_time = datetime.datetime(1, 1, 1, 0, 1)
    task.put()

    pipeline_job = UpdateFlakeAnalysisDataPointsPipeline(
        analysis.key.urlsafe(), build_number)

    pipeline_job.start(queue_name=constants.DEFAULT_QUEUE)
    self.execute_queued_tasks()

    mock_update_fn.assert_called_with(task, analysis)
