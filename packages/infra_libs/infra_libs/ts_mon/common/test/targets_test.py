# Copyright 2015 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

import unittest

from infra_libs.ts_mon.common import targets
from infra_libs.ts_mon.protos import metrics_pb2


class TargetTest(unittest.TestCase):

  def setUp(self):
    self.task0 = targets.TaskTarget('serv', 'job', 'reg', 'host', 0)
    self.task1 = targets.TaskTarget('serv', 'job', 'reg', 'host', 0)
    self.task2 = targets.TaskTarget('serv', 'job', 'reg', 'host', 1)
    self.device0 = targets.DeviceTarget('reg', 'role', 'net', 'host0')
    self.device1 = targets.DeviceTarget('reg', 'role', 'net', 'host0')
    self.device2 = targets.DeviceTarget('reg', 'role', 'net', 'host1')

  def test_eq(self):
    self.assertTrue(self.task0 == self.task1)
    self.assertTrue(self.device0 == self.device1)

    self.assertFalse(self.task0 == self.task2)
    self.assertFalse(self.device0 == self.device2)

    self.assertFalse(self.task0 == self.device0)

  def test_hash(self):
    d = {}
    d[self.task0] = 1
    d[self.task1] = 2
    d[self.task2] = 3
    d[self.device0] = 4
    d[self.device1] = 5
    d[self.device2] = 6

    self.assertEqual(2, d[self.task0])
    self.assertEqual(2, d[self.task1])
    self.assertEqual(3, d[self.task2])
    self.assertEqual(5, d[self.device0])
    self.assertEqual(5, d[self.device1])
    self.assertEqual(6, d[self.device2])


class DeviceTargetTest(unittest.TestCase):

  def test_populate_target_pb(self):
    pb = metrics_pb2.MetricsCollection()
    target = targets.DeviceTarget('reg', 'role', 'net', 'host')
    target.populate_target_pb(pb)

    expected = metrics_pb2.MetricsCollection()
    expected.root_labels.add(key='metro', string_value='reg')
    expected.root_labels.add(key='role', string_value='role')
    expected.root_labels.add(key='hostgroup', string_value='net')
    expected.root_labels.add(key='hostname', string_value='host')
    expected.root_labels.add(key='realm', string_value='ACQ_CHROME')
    expected.root_labels.add(key='alertable', bool_value=True)

    self.assertEqual(pb, expected)

  def test_update_to_dict(self):
    target = targets.DeviceTarget('reg', 'role', 'net', 'host')
    self.assertEqual({
      'region': 'reg',
      'role': 'role',
      'network': 'net',
      'hostname': 'host'}, target.to_dict())
    target.update({'region': 'other', 'hostname': 'guest'})
    self.assertEqual({
      'region': 'other',
      'role': 'role',
      'network': 'net',
      'hostname': 'guest'}, target.to_dict())

  def test_update_private_field(self):
    target = targets.DeviceTarget('reg', 'role', 'net', 'host')
    with self.assertRaises(AttributeError):
      target.update({'realm': 'boo'})

  def test_update_nonexistent_field(self):
    target = targets.DeviceTarget('reg', 'role', 'net', 'host')
    # Simulate a bug: exporting a non-existent field.
    target._fields += ('bad',)
    with self.assertRaises(AttributeError):
      target.update({'bad': 'boo'})


class TaskTargetTest(unittest.TestCase):

  def test_populate_target_pb(self):
    pb = metrics_pb2.MetricsCollection()
    target = targets.TaskTarget('serv', 'job', 'reg', 'host')
    target.populate_target_pb(pb)

    expected = metrics_pb2.MetricsCollection()
    expected.root_labels.add(key='service_name', string_value='serv')
    expected.root_labels.add(key='job_name', string_value='job')
    expected.root_labels.add(key='data_center', string_value='reg')
    expected.root_labels.add(key='host_name', string_value='host')
    expected.root_labels.add(key='task_num', int64_value=0)

    self.assertEqual(pb, expected)

  def test_update_to_dict(self):
    target = targets.TaskTarget('serv', 'job', 'reg', 'host', 5)
    self.assertEqual({
      'service_name': 'serv',
      'job_name': 'job',
      'region': 'reg',
      'hostname': 'host',
      'task_num': 5}, target.to_dict())
    target.update({'region': 'other', 'hostname': 'guest'})
    self.assertEqual({
      'service_name': 'serv',
      'job_name': 'job',
      'region': 'other',
      'hostname': 'guest',
      'task_num': 5}, target.to_dict())

  def test_update_private_field(self):
    target = targets.TaskTarget('serv', 'job', 'reg', 'host')
    with self.assertRaises(AttributeError):
      target.update({'realm': 'boo'})

  def test_update_nonexistent_field(self):
    target = targets.TaskTarget('serv', 'job', 'reg', 'host')
    # Simulate a bug: exporting a non-existent field.
    target._fields += ('bad',)
    with self.assertRaises(AttributeError):
      target.update({'bad': 'boo'})
