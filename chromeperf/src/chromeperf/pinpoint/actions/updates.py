# Copyright 2020 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

import functools
import itertools
import logging

from chromeperf.pinpoint.models import task as task_module

VALID_TRANSITIONS = {
    'pending': {'ongoing', 'completed', 'failed', 'cancelled'},
    'ongoing': {'completed', 'failed', 'cancelled'},
    'cancelled': {'pending'},
    'completed': {'pending'},
    'failed': {'pending'},
}


class Error(Exception):
  pass


class InvalidAmendment(Error):
  pass


class TaskNotFound(Error):
  pass


class InvalidTransition(Error):
  pass


def update_task(client, job, task_id, new_state=None, payload=None):
  """Update a task.

  This enforces that the status transitions are semantically correct, where only
  the transitions defined in the VALID_TRANSITIONS map are allowed.

  When either new_state or payload are not None, this function performs the
  update transactionally. At least one of `new_state` or `payload` must be
  provided in calls to this function.
  """
  if new_state is None and payload is None:
    raise ValueError('Set one of `new_state` or `payload`.')

  if new_state and new_state not in VALID_TRANSITIONS:
    raise InvalidTransition('Unknown state: %s' % (new_state,))

  with client.transaction():
    task = client.get(client.key('Task', task_id, parent=job.key))
    if not task:
      raise TaskNotFound('Task with id "%s" not found for job "%s".' %
                         (task_id, job.job_id))

    if new_state:
      valid_transitions = VALID_TRANSITIONS.get(task['status'])
      if new_state not in valid_transitions:
        raise InvalidTransition(
            'Attempting transition from "%s" to "%s" not in %s; task = %s' %
            (task['status'], new_state, valid_transitions, task))
      task['status'] = new_state

    if payload:
      task['payload'] = payload

    client.put(task)


def extend_task_graph(client, job, vertices, dependencies):
  """Add new vertices and dependency links to the graph.

  Args:
    job: a dashboard.pinpoint.model.job.Job instance.
    vertices: an iterable of TaskVertex instances.
    dependencies: an iterable of Dependency instances.
  """
  if job is None:
    raise ValueError('job must not be None.')
  if not vertices and not dependencies:
    return

  job_key = job.key
  amendment_task_graph = {
      v.id: task_module.Task(
          key=client.key('Task', v.id, parent=job_key),
          task_type=v.vertex_type,
          status='pending',
          payload=v.payload) for v in vertices
  }

  with client.transaction():
    # Ensure that the keys we're adding are not in the graph yet.
    current_tasks = client.query(kind='Task', ancestor=job_key).fetch()
    current_task_keys = set(t.key for t in current_tasks)
    new_task_keys = set(t.key for t in amendment_task_graph.values())
    overlap = new_task_keys & current_task_keys
    if overlap:
      raise InvalidAmendment('vertices (%r) already in task graph.' %
                             (overlap,))

    # Then we add the dependencies.
    current_task_graph = {t.key.id(): t for t in current_tasks}
    handled_dependencies = set()
    update_filter = set(amendment_task_graph)
    for dependency in dependencies:
      dependency_key = client.key('Task', dependency.to, parent=job_key)
      if dependency not in handled_dependencies:
        current_task = current_task_graph.get(dependency.from_)
        amendment_task = amendment_task_graph.get(dependency.from_)
        if current_task is None and amendment_task is None:
          raise InvalidAmendment(
              'dependency `from` (%s) not in amended graph.' %
              (dependency.from_,))
        if current_task:
          current_task_graph[dependency.from_].dependencies.append(
              dependency_key)
        if amendment_task:
          amendment_task_graph[dependency.from_].dependencies.append(
              dependency_key)
        handled_dependencies.add(dependency)
        update_filter.add(dependency.from_)

    client.put_multi(
        itertools.chain(amendment_task_graph.values(), [
            t for id_, t in current_task_graph.items() if id_ in update_filter
        ]),
        use_cache=True)


def log_transition_failures(wrapped_action):
  """Decorator to log state transition failures.

  This is a convenience decorator to handle state transition failures, and
  suppress further exception propagation of the transition failure.
  """

  @functools.wraps(wrapped_action)
  def ActionWrapper(*args, **kwargs):
    try:
      return wrapped_action(*args, **kwargs)
    except InvalidTransition as e:
      logging.error('State transition failed: %s', e)
      return None

  return ActionWrapper