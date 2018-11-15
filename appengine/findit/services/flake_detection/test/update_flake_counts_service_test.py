# Copyright 2018 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

from datetime import datetime
import mock

from libs import time_util
from model.flake.detection.flake_occurrence import FlakeOccurrence
from model.flake.flake import Flake
from model.flake.flake import FlakeCountsByType
from model.flake.flake_type import FlakeType
from services.flake_detection.update_flake_counts_service import (
    UpdateFlakeCounts)
from waterfall.test.wf_testcase import WaterfallTestCase


class UpdateFlakeCountsTest(WaterfallTestCase):

  @mock.patch.object(
      time_util, 'GetDateDaysBeforeNow', return_value=datetime(2018, 8, 23))
  def testFlakeUpdates(self, _):
    luci_project = 'chromium'
    step_ui_name = 'step'
    normalized_step_name = 'normalized_step_name'
    flake1 = Flake.Create(
        luci_project=luci_project,
        normalized_step_name=normalized_step_name,
        normalized_test_name='normalized_test_name_1',
        test_label_name='test_label1')
    flake1.last_occurred_time = datetime(2018, 9, 1)
    flake1.put()
    flake1_key = flake1.key

    flake2 = Flake.Create(
        luci_project=luci_project,
        normalized_step_name=normalized_step_name,
        normalized_test_name='normalized_test_name_2',
        test_label_name='test_label2')
    flake2.last_occurred_time = datetime(2017, 9, 1)
    flake2.false_rejection_count_last_week = 5
    flake2.impacted_cl_count_last_week = 3
    flake2.put()
    flake2_key = flake2.key

    flake3 = Flake.Create(
        luci_project=luci_project,
        normalized_step_name=normalized_step_name,
        normalized_test_name='normalized_test_name_3',
        test_label_name='test_label3')
    flake3.last_occurred_time = datetime(2018, 9, 1)
    flake3.false_rejection_count_last_week = 5
    flake3.impacted_cl_count_last_week = 3
    flake3.put()
    flake3_key = flake3.key

    luci_bucket = 'try'
    luci_builder = 'luci builder'
    legacy_master_name = 'buildbot master'
    legacy_build_number = 999

    occurrence1 = FlakeOccurrence.Create(
        flake_type=FlakeType.CQ_FALSE_REJECTION,
        build_id=1,
        step_ui_name=step_ui_name,
        test_name='t1',
        luci_project=luci_project,
        luci_bucket=luci_bucket,
        luci_builder=luci_builder,
        legacy_master_name=legacy_master_name,
        legacy_build_number=legacy_build_number,
        time_happened=datetime(2018, 9, 1),
        gerrit_cl_id=98761,
        parent_flake_key=flake1_key)
    occurrence1.put()

    occurrence2 = FlakeOccurrence.Create(
        flake_type=FlakeType.CQ_FALSE_REJECTION,
        build_id=2,
        step_ui_name=step_ui_name,
        test_name='t1',
        luci_project=luci_project,
        luci_bucket=luci_bucket,
        luci_builder=luci_builder,
        legacy_master_name=legacy_master_name,
        legacy_build_number=legacy_build_number,
        time_happened=datetime(2018, 8, 31),
        gerrit_cl_id=98761,
        parent_flake_key=flake1_key)
    occurrence2.put()

    occurrence3 = FlakeOccurrence.Create(
        flake_type=FlakeType.CQ_FALSE_REJECTION,
        build_id=3,
        step_ui_name=step_ui_name,
        test_name='t2',
        luci_project=luci_project,
        luci_bucket=luci_bucket,
        luci_builder=luci_builder,
        legacy_master_name=legacy_master_name,
        legacy_build_number=legacy_build_number,
        time_happened=datetime(2018, 7, 31),
        gerrit_cl_id=98763,
        parent_flake_key=flake1_key)
    occurrence3.put()

    occurrence4 = FlakeOccurrence.Create(
        flake_type=FlakeType.CQ_FALSE_REJECTION,
        build_id=4,
        step_ui_name=step_ui_name,
        test_name='t1',
        luci_project=luci_project,
        luci_bucket=luci_bucket,
        luci_builder=luci_builder,
        legacy_master_name=legacy_master_name,
        legacy_build_number=legacy_build_number,
        time_happened=datetime(2018, 8, 31),
        gerrit_cl_id=98764,
        parent_flake_key=flake3_key)
    occurrence4.put()

    occurrence5 = FlakeOccurrence.Create(
        flake_type=FlakeType.RETRY_WITH_PATCH,
        build_id=5,
        step_ui_name=step_ui_name,
        test_name='t2',
        luci_project=luci_project,
        luci_bucket=luci_bucket,
        luci_builder=luci_builder,
        legacy_master_name=legacy_master_name,
        legacy_build_number=legacy_build_number,
        time_happened=datetime(2018, 8, 31),
        gerrit_cl_id=98761,
        parent_flake_key=flake1_key)
    occurrence5.put()

    occurrence6 = FlakeOccurrence.Create(
        flake_type=FlakeType.RETRY_WITH_PATCH,
        build_id=6,
        step_ui_name=step_ui_name,
        test_name='t2',
        luci_project=luci_project,
        luci_bucket=luci_bucket,
        luci_builder=luci_builder,
        legacy_master_name=legacy_master_name,
        legacy_build_number=legacy_build_number,
        time_happened=datetime(2018, 8, 31),
        gerrit_cl_id=98766,
        parent_flake_key=flake1_key)
    occurrence6.put()

    occurrence7 = FlakeOccurrence.Create(
        flake_type=FlakeType.RETRY_WITH_PATCH,
        build_id=7,
        step_ui_name=step_ui_name,
        test_name='t2',
        luci_project=luci_project,
        luci_bucket=luci_bucket,
        luci_builder=luci_builder,
        legacy_master_name=legacy_master_name,
        legacy_build_number=legacy_build_number,
        time_happened=datetime(2018, 8, 31),
        gerrit_cl_id=98767,
        parent_flake_key=flake1_key)
    occurrence7.put()

    UpdateFlakeCounts()

    flake1 = flake1_key.get()
    self.assertEqual(5, flake1.false_rejection_count_last_week)
    self.assertEqual(3, flake1.impacted_cl_count_last_week)
    self.assertEqual([
        FlakeCountsByType(
            flake_type=FlakeType.CQ_FALSE_REJECTION,
            impacted_cl_count=1,
            occurrence_count=2),
        FlakeCountsByType(
            flake_type=FlakeType.RETRY_WITH_PATCH,
            impacted_cl_count=2,
            occurrence_count=3)
    ], flake1.flake_counts_last_week)
    self.assertEqual(120, flake1.flake_score_last_week)

    flake2 = flake2_key.get()
    self.assertEqual([], flake2.flake_counts_last_week)
    self.assertEqual(0, flake2.flake_score_last_week)

    flake3 = flake3_key.get()
    self.assertEqual([
        FlakeCountsByType(
            flake_type=FlakeType.CQ_FALSE_REJECTION,
            impacted_cl_count=1,
            occurrence_count=1)
    ], flake3.flake_counts_last_week)
    self.assertEqual(0, flake3.flake_score_last_week)
