# Copyright 2018 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

import logging
import os

from google.appengine.ext import ndb

from gae_libs import appengine_util
from model.flake.flake import Flake
from model.flake.detection.flake_occurrence import FlakeOccurrence
from model.flake.detection.flake_occurrence import FlakeType
from services import bigquery_helper
from services import monitoring

# Path to the query used to detect flaky tests that caused cq false rejections.
PATH_TO_FLAKY_TESTS_QUERY = os.path.realpath(
    os.path.join(__file__, os.path.pardir,
                 'flaky_tests.cq_false_rejection.sql'))


def _CreateFlakeFromRow(row):
  """Creates a Flake entity from a row fetched from BigQuery."""
  luci_project = row['luci_project']
  step_ui_name = row['step_ui_name']
  test_name = row['test_name']
  luci_builder = row['luci_builder']
  legacy_master_name = row['legacy_master_name']
  legacy_build_number = row['legacy_build_number']

  normalized_step_name = Flake.NormalizeStepName(
      step_name=step_ui_name,
      master_name=legacy_master_name,
      builder_name=luci_builder,
      build_number=legacy_build_number)
  normalized_test_name = Flake.NormalizeTestName(test_name, step_ui_name)
  test_label_name = Flake.GetTestLabelName(test_name, step_ui_name)

  return Flake.Create(
      luci_project=luci_project,
      normalized_step_name=normalized_step_name,
      normalized_test_name=normalized_test_name,
      test_label_name=test_label_name)


def _CreateFlakeOccurrenceFromRow(row):
  """Creates a FlakeOccurrence from a row fetched from BigQuery."""
  luci_project = row['luci_project']
  step_ui_name = row['step_ui_name']
  test_name = row['test_name']
  luci_builder = row['luci_builder']
  legacy_master_name = row['legacy_master_name']
  legacy_build_number = row['legacy_build_number']

  normalized_step_name = Flake.NormalizeStepName(
      step_name=step_ui_name,
      master_name=legacy_master_name,
      builder_name=luci_builder,
      build_number=legacy_build_number)
  normalized_test_name = Flake.NormalizeTestName(test_name)

  flake_id = Flake.GetId(
      luci_project=luci_project,
      normalized_step_name=normalized_step_name,
      normalized_test_name=normalized_test_name)
  flake_key = ndb.Key(Flake, flake_id)

  build_id = row['build_id']
  luci_bucket = row['luci_bucket']
  time_happened = row['test_start_msec']
  gerrit_cl_id = row['gerrit_cl_id']
  flake_occurrence = FlakeOccurrence.Create(
      flake_type=FlakeType.CQ_FALSE_REJECTION,
      build_id=build_id,
      step_ui_name=step_ui_name,
      test_name=test_name,
      luci_project=luci_project,
      luci_bucket=luci_bucket,
      luci_builder=luci_builder,
      legacy_master_name=legacy_master_name,
      legacy_build_number=legacy_build_number,
      time_happened=time_happened,
      gerrit_cl_id=gerrit_cl_id,
      parent_flake_key=flake_key)

  return flake_occurrence


def _StoreMultipleLocalEntities(local_entities):
  """Stores multiple ndb Model entities.

  NOTE: This method doesn't overwrite existing entities.

  Args:
    local_entities: A list of Model entities in local memory. It is OK for
                    local_entities to have duplicates, this method will
                    automatically de-duplicate them.

  Returns:
    Distinct new entities that were written to the ndb.
  """
  key_to_local_entities = {}
  for entity in local_entities:
    key_to_local_entities[entity.key] = entity

  # |local_entities| may have duplicates, need to de-duplicate them.
  unique_entity_keys = key_to_local_entities.keys()

  # get_multi returns a list, and a list item is None if the key wasn't found.
  remote_entities = ndb.get_multi(unique_entity_keys)
  non_existent_entity_keys = [
      unique_entity_keys[i]
      for i in range(len(remote_entities))
      if not remote_entities[i]
  ]
  non_existent_local_entities = [
      key_to_local_entities[key] for key in non_existent_entity_keys
  ]
  ndb.put_multi(non_existent_local_entities)
  return non_existent_local_entities


def _UpdateLastFlakeHappenedTimeForFlakes(occurrences):
  """Updates flakes' last_occurred_time.

  Args:
    occurrences(list): A list of FlakeOccurrence entities.
  """

  flake_key_to_latest_false_rejection_time = {}
  for occurrence in occurrences:
    flake_key = occurrence.key.parent()
    if (not flake_key_to_latest_false_rejection_time.get(flake_key) or
        flake_key_to_latest_false_rejection_time[flake_key] <
        occurrence.time_happened):
      flake_key_to_latest_false_rejection_time[
          flake_key] = occurrence.time_happened

  for flake_key, latest in flake_key_to_latest_false_rejection_time.iteritems():
    flake = flake_key.get()
    if (not flake.last_occurred_time or flake.last_occurred_time < latest):
      flake.last_occurred_time = latest
    flake.put()


def QueryAndStoreFlakes():
  """Runs the query to fetch flake occurrences and store them."""
  with open(PATH_TO_FLAKY_TESTS_QUERY) as f:
    query = f.read()

  success, rows = bigquery_helper.ExecuteQuery(
      appengine_util.GetApplicationId(), query)

  if not success:
    logging.error(
        'Failed executing the query to detect cq false rejection flakes.')
    monitoring.OnFlakeDetectionQueryFailed(flake_type='cq false rejection')
    return

  logging.info('Fetched %d cq false rejection flake occurrences from BigQuery.',
               len(rows))

  local_flakes = [_CreateFlakeFromRow(row) for row in rows]
  _StoreMultipleLocalEntities(local_flakes)

  local_flake_occurrences = [_CreateFlakeOccurrenceFromRow(row) for row in rows]
  new_occurrences = _StoreMultipleLocalEntities(local_flake_occurrences)
  _UpdateLastFlakeHappenedTimeForFlakes(new_occurrences)
  monitoring.OnFlakeDetectionDetectNewOccurrences(
      flake_type='cq false rejection', num_occurrences=len(new_occurrences))
