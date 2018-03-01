# Copyright 2017 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

import copy
from datetime import datetime
import mock

from dto.test_location import TestLocation

from gae_libs import pipelines
from gae_libs.gitiles.cached_gitiles_repository import CachedGitilesRepository
from gae_libs.pipelines import pipeline_handlers
from libs import analysis_status
from libs import time_util
from libs.gitiles.blame import Blame
from model.flake.flake_swarming_task import FlakeSwarmingTask
from model.flake.master_flake_analysis import DataPoint
from model.flake.master_flake_analysis import MasterFlakeAnalysis
from pipelines import report_event_pipeline
from services import swarmed_test_util
from services.flake_failure import heuristic_analysis
from waterfall.flake import confidence
from waterfall.flake import finish_build_analysis_pipeline
from waterfall.flake import flake_constants
from waterfall.flake.finish_build_analysis_pipeline import (
    FinishBuildAnalysisPipeline)
from waterfall.flake.initialize_flake_try_job_pipeline import (
    InitializeFlakeTryJobPipeline)
from waterfall.flake.update_flake_bug_pipeline import UpdateFlakeBugPipeline
from waterfall.test import wf_testcase
from waterfall.test.wf_testcase import DEFAULT_CONFIG_DATA


class FinishBuildAnalysisPipelineTest(wf_testcase.WaterfallTestCase):
  app_module = pipeline_handlers._APP

  @mock.patch.object(finish_build_analysis_pipeline,
                     '_IdentifySuspectedRevisions')
  @mock.patch.object(heuristic_analysis, 'GenerateSuspectedRanges')
  @mock.patch.object(heuristic_analysis,
                     'SaveFlakeCulpritsForSuspectedRevisions')
  def testFinishBuildAnalysisPipeline(self, mock_save, mock_ranges,
                                      mock_revisions):
    master_name = 'm'
    builder_name = 'b'
    build_number = 100
    step_name = 's'
    test_name = 't'
    lower_bound = 1
    upper_bound = 10
    user_range = True
    iterations = 100
    suspected_revision = 'r1'
    suspected_revisions = [suspected_revision]
    suspected_ranges = [[None, suspected_revision]]
    mock_revisions.return_value = suspected_revisions
    mock_ranges.return_value = suspected_ranges

    analysis = MasterFlakeAnalysis.Create(master_name, builder_name,
                                          build_number, step_name, test_name)
    analysis.status = analysis_status.COMPLETED
    analysis.algorithm_parameters = copy.deepcopy(
        DEFAULT_CONFIG_DATA['check_flake_settings'])
    analysis.suspected_flake_build_number = build_number
    analysis.data_points = [
        DataPoint.Create(
            build_number=build_number, blame_list=[suspected_revision])
    ]
    analysis.Save()

    task = FlakeSwarmingTask.Create(master_name, builder_name, build_number,
                                    step_name, test_name)
    task.status = analysis_status.COMPLETED
    task.put()

    self.MockPipeline(
        InitializeFlakeTryJobPipeline,
        '',
        expected_args=[
            analysis.key.urlsafe(), suspected_ranges, iterations, user_range,
            False
        ],
        expected_kwargs={})

    self.MockPipeline(
        UpdateFlakeBugPipeline,
        '',
        expected_args=[analysis.key.urlsafe()],
        expected_kwargs={})

    report_event_input = pipelines.CreateInputObjectInstance(
        report_event_pipeline.ReportEventInput,
        analysis_urlsafe_key=analysis.key.urlsafe())
    self.MockGeneratorPipeline(
        report_event_pipeline.ReportAnalysisEventPipeline, report_event_input,
        None)

    pipeline = FinishBuildAnalysisPipeline(analysis.key.urlsafe(), lower_bound,
                                           upper_bound, iterations, False)
    pipeline.start()
    self.execute_queued_tasks()

    self.assertTrue(mock_save.called)

  @mock.patch.object(finish_build_analysis_pipeline,
                     '_IdentifySuspectedRevisions')
  @mock.patch.object(heuristic_analysis, 'GenerateSuspectedRanges')
  @mock.patch.object(heuristic_analysis,
                     'SaveFlakeCulpritsForSuspectedRevisions')
  def testFinishBuildAnalysisPipelineNoSuspectedDataPoint(
      self, mock_save, mock_ranges, mock_revisions):
    master_name = 'm'
    builder_name = 'b'
    build_number = 100
    step_name = 's'
    test_name = 't'
    lower_bound = 1
    upper_bound = 10
    iterations = 100

    analysis = MasterFlakeAnalysis.Create(master_name, builder_name,
                                          build_number, step_name, test_name)
    analysis.status = analysis_status.COMPLETED
    analysis.algorithm_parameters = copy.deepcopy(
        DEFAULT_CONFIG_DATA['check_flake_settings'])
    analysis.data_points = [DataPoint.Create(build_number=build_number)]
    analysis.Save()

    report_event_input = pipelines.CreateInputObjectInstance(
        report_event_pipeline.ReportEventInput,
        analysis_urlsafe_key=analysis.key.urlsafe())
    self.MockGeneratorPipeline(
        report_event_pipeline.ReportAnalysisEventPipeline, report_event_input,
        None)

    pipeline = FinishBuildAnalysisPipeline(analysis.key.urlsafe(), lower_bound,
                                           upper_bound, iterations, False)
    pipeline.start()
    self.execute_queued_tasks()

    mock_save.assert_not_called()
    mock_ranges.assert_not_called()
    mock_revisions.assert_not_called()

  @mock.patch.object(confidence, 'SteppinessForBuild', return_value=0.6)
  def testGetBuildConfidenceScore(self, _):
    analysis = MasterFlakeAnalysis.Create('m', 'b', 124, 's', 't')
    analysis.data_points = [DataPoint.Create(pass_rate=0.7, build_number=123)]
    self.assertIsNone(
        finish_build_analysis_pipeline._GetBuildConfidenceScore(
            analysis, None, []))
    self.assertEqual(0.6,
                     finish_build_analysis_pipeline._GetBuildConfidenceScore(
                         analysis, 123, analysis.data_points))

  def testGetBuildConfidenceScoreIntroducedNewFlakyTest(self):
    analysis = MasterFlakeAnalysis.Create('m', 'b', 124, 's', 't')
    analysis.data_points = [
        DataPoint.Create(pass_rate=0.7, build_number=123),
        DataPoint.Create(
            pass_rate=flake_constants.PASS_RATE_TEST_NOT_FOUND,
            build_number=122)
    ]
    self.assertEqual(1.0,
                     finish_build_analysis_pipeline._GetBuildConfidenceScore(
                         analysis, 123, analysis.data_points))

  def testUserSpecifiedRange(self):
    self.assertTrue(
        finish_build_analysis_pipeline._UserSpecifiedRange(123, 125))
    self.assertFalse(
        finish_build_analysis_pipeline._UserSpecifiedRange(None, 123))
    self.assertFalse(
        finish_build_analysis_pipeline._UserSpecifiedRange(123, None))
    self.assertFalse(
        finish_build_analysis_pipeline._UserSpecifiedRange(None, None))

  def testUpdateAnalysisResults(self):
    analysis = MasterFlakeAnalysis.Create('m', 'b', 123, 's', 't')
    analysis.last_attempted_swarming_task_id = '12345'
    finish_build_analysis_pipeline._UpdateAnalysisResults(
        analysis, 100, analysis_status.COMPLETED, None)
    self.assertEqual(analysis.suspected_flake_build_number, 100)
    self.assertIsNone(analysis.last_attempted_swarming_task_id)

  @mock.patch.object(time_util, 'GetUTCNow', return_value=datetime(2017, 6, 27))
  def testUpdateAnalysisResultsWithError(self, _):
    analysis = MasterFlakeAnalysis.Create('m', 'b', 123, 's', 't')
    analysis.start_time = datetime(2017, 6, 26, 23)
    analysis.last_attempted_swarming_task_id = '12345'
    finish_build_analysis_pipeline._UpdateAnalysisResults(
        analysis, None, analysis_status.ERROR, {
            'error': 'error'
        })
    self.assertEqual(analysis_status.ERROR, analysis.status)
    self.assertEqual(datetime(2017, 6, 27), analysis.end_time)

  @mock.patch.object(swarmed_test_util, 'GetTestLocation')
  @mock.patch.object(CachedGitilesRepository, 'GetBlame')
  def testIdentifySuspectedRanges(self, mock_blame, mock_test_location):
    mock_blame.return_value = [Blame('r1000', 'a/b.cc')]
    mock_test_location.return_value = TestLocation(file='a/b.cc', line=1)
    suspected_revision = 'r1000'
    analysis = MasterFlakeAnalysis.Create('m', 'b', 123, 's', 't')
    analysis.data_points = [
        DataPoint.Create(
            build_number=123, git_hash='r1000', blame_list=[suspected_revision])
    ]
    analysis.suspected_flake_build_number = 123
    analysis.Save()

    self.assertEqual([suspected_revision],
                     finish_build_analysis_pipeline._IdentifySuspectedRevisions(
                         analysis, None))

  @mock.patch.object(swarmed_test_util, 'GetTestLocation', return_value=None)
  def testIdentifysuspectedRangesNoTestLocation(self, _):
    analysis = MasterFlakeAnalysis.Create('m', 'b', 123, 's', 't')
    analysis.data_points = [
        DataPoint.Create(build_number=123, git_hash='r1000', blame_list=['r1'])
    ]
    analysis.suspected_flake_build_number = 123
    analysis.Save()

    self.assertEqual([],
                     finish_build_analysis_pipeline._IdentifySuspectedRevisions(
                         analysis, None))
