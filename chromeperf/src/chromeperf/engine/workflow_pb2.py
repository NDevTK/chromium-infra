# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: chromeperf/engine/workflow.proto

from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.protobuf import struct_pb2 as google_dot_protobuf_dot_struct__pb2
from google.protobuf import timestamp_pb2 as google_dot_protobuf_dot_timestamp__pb2
from chromeperf.engine import graph_pb2 as chromeperf_dot_engine_dot_graph__pb2


DESCRIPTOR = _descriptor.FileDescriptor(
  name='chromeperf/engine/workflow.proto',
  package='chromeperf.engine',
  syntax='proto3',
  serialized_options=None,
  create_key=_descriptor._internal_create_key,
  serialized_pb=b'\n chromeperf/engine/workflow.proto\x12\x11\x63hromeperf.engine\x1a\x1cgoogle/protobuf/struct.proto\x1a\x1fgoogle/protobuf/timestamp.proto\x1a\x1d\x63hromeperf/engine/graph.proto\"C\n\x11GraphTemplateSpec\x12\x19\n\x11graph_template_id\x18\x01 \x01(\t\x12\x13\n\x0brevision_id\x18\x02 \x01(\t\"\x9b\x03\n\x08Workflow\x12\n\n\x02id\x18\x01 \x01(\t\x12-\n\ttimestamp\x18\x02 \x01(\x0b\x32\x1a.google.protobuf.Timestamp\x12\'\n\x06inputs\x18\x03 \x01(\x0b\x32\x17.google.protobuf.Struct\x12:\n\x06status\x18\x04 \x01(\x0e\x32*.chromeperf.engine.Workflow.WorkflowStatus\x12\x17\n\x0frequesting_user\x18\x05 \x01(\t\x12\x17\n\x0fservice_account\x18\x06 \x01(\t\x12\x41\n\x13graph_template_spec\x18\x07 \x01(\x0b\x32$.chromeperf.engine.GraphTemplateSpec\x12\'\n\x05graph\x18\x08 \x01(\x0b\x32\x18.chromeperf.engine.Graph\x12\x14\n\x0c\x61rchive_link\x18\t \x01(\t\";\n\x0eWorkflowStatus\x12\x0f\n\x0bUNSPECIFIED\x10\x00\x12\n\n\x06\x41\x43TIVE\x10\x01\x12\x0c\n\x08\x41RCHIVED\x10\x02\x62\x06proto3'
  ,
  dependencies=[google_dot_protobuf_dot_struct__pb2.DESCRIPTOR,google_dot_protobuf_dot_timestamp__pb2.DESCRIPTOR,chromeperf_dot_engine_dot_graph__pb2.DESCRIPTOR,])



_WORKFLOW_WORKFLOWSTATUS = _descriptor.EnumDescriptor(
  name='WorkflowStatus',
  full_name='chromeperf.engine.Workflow.WorkflowStatus',
  filename=None,
  file=DESCRIPTOR,
  create_key=_descriptor._internal_create_key,
  values=[
    _descriptor.EnumValueDescriptor(
      name='UNSPECIFIED', index=0, number=0,
      serialized_options=None,
      type=None,
      create_key=_descriptor._internal_create_key),
    _descriptor.EnumValueDescriptor(
      name='ACTIVE', index=1, number=1,
      serialized_options=None,
      type=None,
      create_key=_descriptor._internal_create_key),
    _descriptor.EnumValueDescriptor(
      name='ARCHIVED', index=2, number=2,
      serialized_options=None,
      type=None,
      create_key=_descriptor._internal_create_key),
  ],
  containing_type=None,
  serialized_options=None,
  serialized_start=571,
  serialized_end=630,
)
_sym_db.RegisterEnumDescriptor(_WORKFLOW_WORKFLOWSTATUS)


_GRAPHTEMPLATESPEC = _descriptor.Descriptor(
  name='GraphTemplateSpec',
  full_name='chromeperf.engine.GraphTemplateSpec',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='graph_template_id', full_name='chromeperf.engine.GraphTemplateSpec.graph_template_id', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='revision_id', full_name='chromeperf.engine.GraphTemplateSpec.revision_id', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=149,
  serialized_end=216,
)


_WORKFLOW = _descriptor.Descriptor(
  name='Workflow',
  full_name='chromeperf.engine.Workflow',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='id', full_name='chromeperf.engine.Workflow.id', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='timestamp', full_name='chromeperf.engine.Workflow.timestamp', index=1,
      number=2, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='inputs', full_name='chromeperf.engine.Workflow.inputs', index=2,
      number=3, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='status', full_name='chromeperf.engine.Workflow.status', index=3,
      number=4, type=14, cpp_type=8, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='requesting_user', full_name='chromeperf.engine.Workflow.requesting_user', index=4,
      number=5, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='service_account', full_name='chromeperf.engine.Workflow.service_account', index=5,
      number=6, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='graph_template_spec', full_name='chromeperf.engine.Workflow.graph_template_spec', index=6,
      number=7, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='graph', full_name='chromeperf.engine.Workflow.graph', index=7,
      number=8, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='archive_link', full_name='chromeperf.engine.Workflow.archive_link', index=8,
      number=9, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
    _WORKFLOW_WORKFLOWSTATUS,
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=219,
  serialized_end=630,
)

_WORKFLOW.fields_by_name['timestamp'].message_type = google_dot_protobuf_dot_timestamp__pb2._TIMESTAMP
_WORKFLOW.fields_by_name['inputs'].message_type = google_dot_protobuf_dot_struct__pb2._STRUCT
_WORKFLOW.fields_by_name['status'].enum_type = _WORKFLOW_WORKFLOWSTATUS
_WORKFLOW.fields_by_name['graph_template_spec'].message_type = _GRAPHTEMPLATESPEC
_WORKFLOW.fields_by_name['graph'].message_type = chromeperf_dot_engine_dot_graph__pb2._GRAPH
_WORKFLOW_WORKFLOWSTATUS.containing_type = _WORKFLOW
DESCRIPTOR.message_types_by_name['GraphTemplateSpec'] = _GRAPHTEMPLATESPEC
DESCRIPTOR.message_types_by_name['Workflow'] = _WORKFLOW
_sym_db.RegisterFileDescriptor(DESCRIPTOR)

GraphTemplateSpec = _reflection.GeneratedProtocolMessageType('GraphTemplateSpec', (_message.Message,), {
  'DESCRIPTOR' : _GRAPHTEMPLATESPEC,
  '__module__' : 'chromeperf.engine.workflow_pb2'
  # @@protoc_insertion_point(class_scope:chromeperf.engine.GraphTemplateSpec)
  })
_sym_db.RegisterMessage(GraphTemplateSpec)

Workflow = _reflection.GeneratedProtocolMessageType('Workflow', (_message.Message,), {
  'DESCRIPTOR' : _WORKFLOW,
  '__module__' : 'chromeperf.engine.workflow_pb2'
  # @@protoc_insertion_point(class_scope:chromeperf.engine.Workflow)
  })
_sym_db.RegisterMessage(Workflow)


# @@protoc_insertion_point(module_scope)
