# Copyright 2015 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

import base64
from collections import defaultdict
import json
import urllib
import zlib

from google.appengine.ext import ndb

from model.wf_step import WfStep
from waterfall import auth_util
from waterfall import waterfall_config
from waterfall.swarming_task_request import SwarmingTaskRequest


STATES_RUNNING = ('RUNNING', 'PENDING')
STATE_COMPLETED = 'COMPLETED'
STATES_NOT_RUNNING = (
    'EXPIRED', 'TIMED_OUT', 'BOT_DIED', 'CANCELED', 'COMPLETED')


def _SendRequestToServer(url, http_client, post_data=None):
  """Sends GET/POST request to arbitrary url and returns response content."""
  headers = {'Authorization': 'Bearer ' + auth_util.GetAuthToken()}
  if post_data:
    post_data = json.dumps(post_data, sort_keys=True, separators=(',', ':'))
    headers['Content-Type'] = 'application/json; charset=UTF-8'
    headers['Content-Length'] = len(post_data)
    status_code, content = http_client.Post(url, post_data, headers=headers)
  else:
    status_code, content = http_client.Get(url, headers=headers)

  if status_code != 200:
    # The retry upon 50x (501 excluded) is automatically handled in the
    # underlying http_client.
    # By default, it retries 5 times with exponential backoff.
    return None
  return content


def GetSwarmingTaskRequest(task_id, http_client):
  """Returns an instance of SwarmingTaskRequest representing the given task."""
  swarming_server_host = waterfall_config.GetSwarmingSettings().get(
      'server_host')
  url = ('https://%s/_ah/api/swarming/v1/task/%s/request') % (
      swarming_server_host, task_id)
  json_data = json.loads(_SendRequestToServer(url, http_client))
  return SwarmingTaskRequest.Deserialize(json_data)


def TriggerSwarmingTask(request, http_client):
  """Triggers a new Swarming task for the given request.

  The Swarming task priority will be overwritten, and extra tags might be added.
  Args:
    request (SwarmingTaskRequest): A Swarming task request.
    http_client (RetryHttpClient): An http client with automatic retry.
  """
  # Use a priority much lower than CQ for now (CQ's priority is 30).
  # Later we might use a higher priority -- a lower value here.
  # Note: the smaller value, the higher priority.
  swarming_settings = waterfall_config.GetSwarmingSettings()
  request_expiration_hours = swarming_settings.get('request_expiration_hours')
  request.priority = max(100, swarming_settings.get('default_request_priority'))
  request.expiration_secs = request_expiration_hours * 60 * 60

  request.tags.extend(['findit:1', 'project:Chromium', 'purpose:post-commit'])

  url = 'https://%s/_ah/api/swarming/v1/tasks/new' % swarming_settings.get(
      'server_host')
  response_data = _SendRequestToServer(url, http_client, request.Serialize())
  return json.loads(response_data)['task_id']


def ListSwarmingTasksDataByTags(
    master_name, builder_name, build_number, http_client, step_name=None):
  """Downloads tasks data from swarming server."""
  base_url = ('https://%s/_ah/api/swarming/v1/tasks/'
              'list?tags=%s&tags=%s&tags=%s') % (
                  waterfall_config.GetSwarmingSettings().get('server_host'),
                  urllib.quote('master:%s' % master_name),
                  urllib.quote('buildername:%s' % builder_name),
                  urllib.quote('buildnumber:%d' % build_number))
  if step_name:
    base_url += '&tags=%s' % urllib.quote('stepname:%s' % step_name)

  items = []
  cursor = None

  while True:
    if not cursor:
      url = base_url
    else:
      url = base_url + '&cursor=%s' % urllib.quote(cursor)
    new_data = _SendRequestToServer(url, http_client)
    if not new_data:
      break

    new_data_json = json.loads(new_data)
    if new_data_json.get('items'):
      items.extend(new_data_json['items'])

    if new_data_json.get('cursor'):
      cursor = new_data_json['cursor']
    else:
      break

  return items


def _GenerateIsolatedData(outputs_ref):
  return {
      'digest': outputs_ref['isolated'],
      'namespace': outputs_ref['namespace'],
      'isolatedserver': outputs_ref['isolatedserver']
  }


def GetSwarmingTaskResultById(task_id, http_client):
  """Gets swarming result, checks state and returns outputs ref if needed."""
  base_url = ('https://%s/_ah/api/swarming/v1/task/%s/result') % (
      waterfall_config.GetSwarmingSettings().get('server_host'), task_id)
  data = _SendRequestToServer(base_url, http_client)
  json_data = json.loads(data)

  return json_data


def GetSwarmingTaskFailureLog(outputs_ref, http_client):
  """Downloads failure log from isolated server."""
  isolated_data = _GenerateIsolatedData(outputs_ref)
  return _DownloadTestResults(isolated_data, http_client)


def GetTagValue(tags, tag_name):
  """Returns the content for a specific tag."""
  tag_prefix = tag_name + ':'
  content = None
  for tag in tags:
    if tag.startswith(tag_prefix):
      content = tag[len(tag_prefix):]
      break
  return content


def GetIsolatedDataForFailedBuild(
    master_name, builder_name, build_number, failed_steps, http_client):
  """Checks failed step_names in swarming log for the build.

  Searches each failed step_name to identify swarming/non-swarming tests
  and keeps track of isolated data for each failed swarming steps.
  """
  data = ListSwarmingTasksDataByTags(
      master_name, builder_name, build_number, http_client)
  if not data:
    return False

  tag_name = 'stepname'
  build_isolated_data = defaultdict(list)
  for item in data:
    if item['failure'] and not item['internal_failure']:
      # Only retrieves test results from tasks which have failures and
      # the failure should not be internal infrastructure failure.
      swarming_step_name = GetTagValue(item['tags'], tag_name)
      if swarming_step_name in failed_steps:
        isolated_data = _GenerateIsolatedData(item['outputs_ref'])
        build_isolated_data[swarming_step_name].append(isolated_data)

  new_steps = []
  for step_name in build_isolated_data:
    failed_steps[step_name]['list_isolated_data'] = (
        build_isolated_data[step_name])

    # Create WfStep object for all the failed steps.
    step = WfStep.Create(master_name, builder_name, build_number, step_name)
    step.isolated = True
    new_steps.append(step)

  ndb.put_multi(new_steps)
  return True


def GetIsolatedDataForStep(
    master_name, builder_name, build_number, step_name,
    http_client):
  """Returns the isolated data for a specific step."""
  step_isolated_data = []
  data = ListSwarmingTasksDataByTags(master_name, builder_name, build_number,
                                     http_client, step_name)
  if not data:
    return step_isolated_data

  for item in data:
    if item['failure'] and not item['internal_failure']:
      # Only retrieves test results from tasks which have failures and
      # the failure should not be internal infrastructure failure.
      isolated_data = _GenerateIsolatedData(item['outputs_ref'])
      step_isolated_data.append(isolated_data)

  return step_isolated_data


def _FetchOutputJsonInfoFromIsolatedServer(isolated_data, http_client):
  """Sends POST request to isolated server and returns response content.

  This function is used for fetching
    1. hash code for the output.json file,
    2. the redirect url.
  """
  post_data = {
      'digest': isolated_data['digest'],
      'namespace': isolated_data['namespace']
  }
  url = '%s/_ah/api/isolateservice/v1/retrieve' % (
      isolated_data['isolatedserver'])
  content = _SendRequestToServer(url, http_client, post_data)
  return content


def _GetOutputJsonHash(content):
  """Gets hash for output.json.

  Parses response content of the request using hash for .isolated file and
  returns the hash for output.json file.

  Args:
    content (string): Content returned by the POST request to isolated server
        for hash to output.json.
  """
  content_json = json.loads(content)
  content_string = zlib.decompress(base64.b64decode(content_json['content']))
  json_result = json.loads(content_string)
  return json_result.get('files', {}).get('output.json', {}).get('h')


def _RetrieveOutputJsonFile(output_json_content, http_client):
  """Downloads output.json file from isolated server or process it directly.

  If there is a url provided, send get request to that url to download log;
  else the log would be in content so use it directly.
  """
  json_content = json.loads(output_json_content)
  output_json_url = json_content.get('url')
  if output_json_url:
    get_content = _SendRequestToServer(output_json_url, http_client)
  elif json_content.get('content'):
    get_content = base64.b64decode(json_content['content'])
  else:  # pragma: no cover
    get_content = None  # Just for precausion.
  return json.loads(zlib.decompress(get_content)) if get_content else None


def _DownloadTestResults(isolated_data, http_client):
  """Downloads the output.json file and returns the json object."""
  # First POST request to get hash for the output.json file.
  content = _FetchOutputJsonInfoFromIsolatedServer(
      isolated_data, http_client)
  if not content:
    return None
  output_json_hash = _GetOutputJsonHash(content)
  if not output_json_hash:
    return None

  # Second POST request to get the redirect url for the output.json file.
  data_for_output_json = {
      'digest': output_json_hash,
      'namespace': isolated_data['namespace'],
      'isolatedserver': isolated_data['isolatedserver']
  }
  output_json_content = _FetchOutputJsonInfoFromIsolatedServer(
      data_for_output_json, http_client)
  if not output_json_content:
    return None

  # GET Request to get output.json file.
  return _RetrieveOutputJsonFile(
      output_json_content, http_client)


def _MergeListsOfDicts(merged, shard):
  output = []
  for i in xrange(max(len(merged), len(shard))):
    merged_dict = merged[i] if i < len(merged) else {}
    shard_dict = shard[i] if i < len(shard) else {}
    output_dict = merged_dict.copy()
    output_dict.update(shard_dict)
    output.append(output_dict)
  return output


def _MergeSwarmingTestShards(shard_results):
  """Merges the shards into one.

  Args:
    shard_results (list): A list of dicts with individual shard results.

  Returns:
    A dict with the following form:
    {
      'all_tests':[
        'AllForms/FormStructureBrowserTest.DataDrivenHeuristics/0',
        'AllForms/FormStructureBrowserTest.DataDrivenHeuristics/1',
        'AllForms/FormStructureBrowserTest.DataDrivenHeuristics/10',
        ...
      ]
      'per_iteration_data':[
        {
          'AllForms/FormStructureBrowserTest.DataDrivenHeuristics/109': [
            {
              'elapsed_time_ms': 4719,
              'losless_snippet': true,
              'output_snippet': '[ RUN      ] run outputs\\n',
              'output_snippet_base64': 'WyBSVU4gICAgICBdIEFsbEZvcm1zL0Zvcm1T'
              'status': 'SUCCESS'
            }
          ],
        },
        ...
      ]
    }
  """
  merged_results = {
      'all_tests': set(),
      'per_iteration_data': []
  }
  for shard_result in shard_results:
    merged_results['all_tests'].update(shard_result.get('all_tests', []))
    merged_results['per_iteration_data'] = _MergeListsOfDicts(
        merged_results['per_iteration_data'],
        shard_result.get('per_iteration_data', []))
  merged_results['all_tests'] = sorted(merged_results['all_tests'])
  return merged_results


def RetrieveShardedTestResultsFromIsolatedServer(
    list_isolated_data, http_client):
  """Gets test results from isolated server and merge the results."""
  shard_results = []
  for isolated_data in list_isolated_data:
    output_json = _DownloadTestResults(isolated_data, http_client)
    if not output_json:
      return None
    shard_results.append(output_json)

  if len(list_isolated_data) == 1:
    return shard_results[0]
  return _MergeSwarmingTestShards(shard_results)
