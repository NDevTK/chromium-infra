# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: project_config.proto

import sys
_b=sys.version_info[0]<3 and (lambda x:x) or (lambda x:x.encode('latin1'))
from google.protobuf.internal import enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.protobuf import wrappers_pb2 as google_dot_protobuf_dot_wrappers__pb2
import common_pb2 as common__pb2


DESCRIPTOR = _descriptor.FileDescriptor(
  name='project_config.proto',
  package='buildbucket',
  syntax='proto3',
  serialized_options=_b('Z4go.chromium.org/luci/buildbucket/proto;buildbucketpb'),
  serialized_pb=_b('\n\x14project_config.proto\x12\x0b\x62uildbucket\x1a\x1egoogle/protobuf/wrappers.proto\x1a\x0c\x63ommon.proto\"z\n\x03\x41\x63l\x12#\n\x04role\x18\x01 \x01(\x0e\x32\x15.buildbucket.Acl.Role\x12\r\n\x05group\x18\x02 \x01(\t\x12\x10\n\x08identity\x18\x03 \x01(\t\"-\n\x04Role\x12\n\n\x06READER\x10\x00\x12\r\n\tSCHEDULER\x10\x01\x12\n\n\x06WRITER\x10\x02\"6\n\x06\x41\x63lSet\x12\x0c\n\x04name\x18\x01 \x01(\t\x12\x1e\n\x04\x61\x63ls\x18\x02 \x03(\x0b\x32\x10.buildbucket.Acl\"\xc2\x06\n\x07\x42uilder\x12\x0c\n\x04name\x18\x01 \x01(\t\x12\x15\n\rswarming_host\x18\x15 \x01(\t\x12\x0e\n\x06mixins\x18\n \x03(\t\x12\x10\n\x08\x63\x61tegory\x18\x06 \x01(\t\x12\x15\n\rswarming_tags\x18\x02 \x03(\t\x12\x12\n\ndimensions\x18\x03 \x03(\t\x12+\n\x06recipe\x18\x04 \x01(\x0b\x32\x1b.buildbucket.Builder.Recipe\x12\'\n\x03\x65xe\x18\x17 \x01(\x0b\x32\x1a.buildbucket.v2.Executable\x12\x12\n\nproperties\x18\x18 \x01(\t\x12\x10\n\x08priority\x18\x05 \x01(\r\x12\x1e\n\x16\x65xecution_timeout_secs\x18\x07 \x01(\r\x12\x17\n\x0f\x65xpiration_secs\x18\x14 \x01(\r\x12/\n\x06\x63\x61\x63hes\x18\t \x03(\x0b\x32\x1f.buildbucket.Builder.CacheEntry\x12*\n\rbuild_numbers\x18\x10 \x01(\x0e\x32\x13.buildbucket.Toggle\x12\x17\n\x0fservice_account\x18\x0c \x01(\t\x12\x33\n\x16\x61uto_builder_dimension\x18\x11 \x01(\x0e\x32\x13.buildbucket.Toggle\x12)\n\x0c\x65xperimental\x18\x12 \x01(\x0e\x32\x13.buildbucket.Toggle\x12\x1b\n\x13luci_migration_host\x18\x13 \x01(\t\x12\x45\n\x1ftask_template_canary_percentage\x18\x16 \x01(\x0b\x32\x1c.google.protobuf.UInt32Value\x1aJ\n\nCacheEntry\x12\x0c\n\x04name\x18\x01 \x01(\t\x12\x0c\n\x04path\x18\x02 \x01(\t\x12 \n\x18wait_for_warm_cache_secs\x18\x03 \x01(\x05\x1ar\n\x06Recipe\x12\x0c\n\x04name\x18\x02 \x01(\t\x12\x14\n\x0c\x63ipd_package\x18\x06 \x01(\t\x12\x14\n\x0c\x63ipd_version\x18\x05 \x01(\t\x12\x12\n\nproperties\x18\x03 \x03(\t\x12\x14\n\x0cproperties_j\x18\x04 \x03(\tJ\x04\x08\x01\x10\x02J\x04\x08\x08\x10\tJ\x04\x08\x0b\x10\x0cJ\x04\x08\r\x10\x0eJ\x04\x08\x0f\x10\x10\"\xcf\x01\n\x08Swarming\x12\x10\n\x08hostname\x18\x01 \x01(\t\x12\x12\n\nurl_format\x18\x02 \x01(\t\x12.\n\x10\x62uilder_defaults\x18\x03 \x01(\x0b\x32\x14.buildbucket.Builder\x12&\n\x08\x62uilders\x18\x04 \x03(\x0b\x32\x14.buildbucket.Builder\x12\x45\n\x1ftask_template_canary_percentage\x18\x05 \x01(\x0b\x32\x1c.google.protobuf.UInt32Value\"q\n\x06\x42ucket\x12\x0c\n\x04name\x18\x01 \x01(\t\x12\x1e\n\x04\x61\x63ls\x18\x02 \x03(\x0b\x32\x10.buildbucket.Acl\x12\x10\n\x08\x61\x63l_sets\x18\x04 \x03(\t\x12\'\n\x08swarming\x18\x03 \x01(\x0b\x32\x15.buildbucket.Swarming\"\x8b\x01\n\x0e\x42uildbucketCfg\x12$\n\x07\x62uckets\x18\x01 \x03(\x0b\x32\x13.buildbucket.Bucket\x12%\n\x08\x61\x63l_sets\x18\x02 \x03(\x0b\x32\x13.buildbucket.AclSet\x12,\n\x0e\x62uilder_mixins\x18\x03 \x03(\x0b\x32\x14.buildbucket.Builder*$\n\x06Toggle\x12\t\n\x05UNSET\x10\x00\x12\x07\n\x03YES\x10\x01\x12\x06\n\x02NO\x10\x02\x42\x36Z4go.chromium.org/luci/buildbucket/proto;buildbucketpbb\x06proto3')
  ,
  dependencies=[google_dot_protobuf_dot_wrappers__pb2.DESCRIPTOR,common__pb2.DESCRIPTOR,])

_TOGGLE = _descriptor.EnumDescriptor(
  name='Toggle',
  full_name='buildbucket.Toggle',
  filename=None,
  file=DESCRIPTOR,
  values=[
    _descriptor.EnumValueDescriptor(
      name='UNSET', index=0, number=0,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='YES', index=1, number=1,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='NO', index=2, number=2,
      serialized_options=None,
      type=None),
  ],
  containing_type=None,
  serialized_options=None,
  serialized_start=1567,
  serialized_end=1603,
)
_sym_db.RegisterEnumDescriptor(_TOGGLE)

Toggle = enum_type_wrapper.EnumTypeWrapper(_TOGGLE)
UNSET = 0
YES = 1
NO = 2


_ACL_ROLE = _descriptor.EnumDescriptor(
  name='Role',
  full_name='buildbucket.Acl.Role',
  filename=None,
  file=DESCRIPTOR,
  values=[
    _descriptor.EnumValueDescriptor(
      name='READER', index=0, number=0,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='SCHEDULER', index=1, number=1,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='WRITER', index=2, number=2,
      serialized_options=None,
      type=None),
  ],
  containing_type=None,
  serialized_options=None,
  serialized_start=160,
  serialized_end=205,
)
_sym_db.RegisterEnumDescriptor(_ACL_ROLE)


_ACL = _descriptor.Descriptor(
  name='Acl',
  full_name='buildbucket.Acl',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='role', full_name='buildbucket.Acl.role', index=0,
      number=1, type=14, cpp_type=8, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='group', full_name='buildbucket.Acl.group', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='identity', full_name='buildbucket.Acl.identity', index=2,
      number=3, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
    _ACL_ROLE,
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=83,
  serialized_end=205,
)


_ACLSET = _descriptor.Descriptor(
  name='AclSet',
  full_name='buildbucket.AclSet',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='name', full_name='buildbucket.AclSet.name', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='acls', full_name='buildbucket.AclSet.acls', index=1,
      number=2, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
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
  serialized_start=207,
  serialized_end=261,
)


_BUILDER_CACHEENTRY = _descriptor.Descriptor(
  name='CacheEntry',
  full_name='buildbucket.Builder.CacheEntry',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='name', full_name='buildbucket.Builder.CacheEntry.name', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='path', full_name='buildbucket.Builder.CacheEntry.path', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='wait_for_warm_cache_secs', full_name='buildbucket.Builder.CacheEntry.wait_for_warm_cache_secs', index=2,
      number=3, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
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
  serialized_start=884,
  serialized_end=958,
)

_BUILDER_RECIPE = _descriptor.Descriptor(
  name='Recipe',
  full_name='buildbucket.Builder.Recipe',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='name', full_name='buildbucket.Builder.Recipe.name', index=0,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='cipd_package', full_name='buildbucket.Builder.Recipe.cipd_package', index=1,
      number=6, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='cipd_version', full_name='buildbucket.Builder.Recipe.cipd_version', index=2,
      number=5, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='properties', full_name='buildbucket.Builder.Recipe.properties', index=3,
      number=3, type=9, cpp_type=9, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='properties_j', full_name='buildbucket.Builder.Recipe.properties_j', index=4,
      number=4, type=9, cpp_type=9, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
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
  serialized_start=960,
  serialized_end=1074,
)

_BUILDER = _descriptor.Descriptor(
  name='Builder',
  full_name='buildbucket.Builder',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='name', full_name='buildbucket.Builder.name', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='swarming_host', full_name='buildbucket.Builder.swarming_host', index=1,
      number=21, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='mixins', full_name='buildbucket.Builder.mixins', index=2,
      number=10, type=9, cpp_type=9, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='category', full_name='buildbucket.Builder.category', index=3,
      number=6, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='swarming_tags', full_name='buildbucket.Builder.swarming_tags', index=4,
      number=2, type=9, cpp_type=9, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='dimensions', full_name='buildbucket.Builder.dimensions', index=5,
      number=3, type=9, cpp_type=9, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='recipe', full_name='buildbucket.Builder.recipe', index=6,
      number=4, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='exe', full_name='buildbucket.Builder.exe', index=7,
      number=23, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='properties', full_name='buildbucket.Builder.properties', index=8,
      number=24, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='priority', full_name='buildbucket.Builder.priority', index=9,
      number=5, type=13, cpp_type=3, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='execution_timeout_secs', full_name='buildbucket.Builder.execution_timeout_secs', index=10,
      number=7, type=13, cpp_type=3, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='expiration_secs', full_name='buildbucket.Builder.expiration_secs', index=11,
      number=20, type=13, cpp_type=3, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='caches', full_name='buildbucket.Builder.caches', index=12,
      number=9, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='build_numbers', full_name='buildbucket.Builder.build_numbers', index=13,
      number=16, type=14, cpp_type=8, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='service_account', full_name='buildbucket.Builder.service_account', index=14,
      number=12, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='auto_builder_dimension', full_name='buildbucket.Builder.auto_builder_dimension', index=15,
      number=17, type=14, cpp_type=8, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='experimental', full_name='buildbucket.Builder.experimental', index=16,
      number=18, type=14, cpp_type=8, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='luci_migration_host', full_name='buildbucket.Builder.luci_migration_host', index=17,
      number=19, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='task_template_canary_percentage', full_name='buildbucket.Builder.task_template_canary_percentage', index=18,
      number=22, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[_BUILDER_CACHEENTRY, _BUILDER_RECIPE, ],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=264,
  serialized_end=1098,
)


_SWARMING = _descriptor.Descriptor(
  name='Swarming',
  full_name='buildbucket.Swarming',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='hostname', full_name='buildbucket.Swarming.hostname', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='url_format', full_name='buildbucket.Swarming.url_format', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='builder_defaults', full_name='buildbucket.Swarming.builder_defaults', index=2,
      number=3, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='builders', full_name='buildbucket.Swarming.builders', index=3,
      number=4, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='task_template_canary_percentage', full_name='buildbucket.Swarming.task_template_canary_percentage', index=4,
      number=5, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
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
  serialized_start=1101,
  serialized_end=1308,
)


_BUCKET = _descriptor.Descriptor(
  name='Bucket',
  full_name='buildbucket.Bucket',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='name', full_name='buildbucket.Bucket.name', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='acls', full_name='buildbucket.Bucket.acls', index=1,
      number=2, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='acl_sets', full_name='buildbucket.Bucket.acl_sets', index=2,
      number=4, type=9, cpp_type=9, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='swarming', full_name='buildbucket.Bucket.swarming', index=3,
      number=3, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
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
  serialized_start=1310,
  serialized_end=1423,
)


_BUILDBUCKETCFG = _descriptor.Descriptor(
  name='BuildbucketCfg',
  full_name='buildbucket.BuildbucketCfg',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='buckets', full_name='buildbucket.BuildbucketCfg.buckets', index=0,
      number=1, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='acl_sets', full_name='buildbucket.BuildbucketCfg.acl_sets', index=1,
      number=2, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='builder_mixins', full_name='buildbucket.BuildbucketCfg.builder_mixins', index=2,
      number=3, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
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
  serialized_start=1426,
  serialized_end=1565,
)

_ACL.fields_by_name['role'].enum_type = _ACL_ROLE
_ACL_ROLE.containing_type = _ACL
_ACLSET.fields_by_name['acls'].message_type = _ACL
_BUILDER_CACHEENTRY.containing_type = _BUILDER
_BUILDER_RECIPE.containing_type = _BUILDER
_BUILDER.fields_by_name['recipe'].message_type = _BUILDER_RECIPE
_BUILDER.fields_by_name['exe'].message_type = common__pb2._EXECUTABLE
_BUILDER.fields_by_name['caches'].message_type = _BUILDER_CACHEENTRY
_BUILDER.fields_by_name['build_numbers'].enum_type = _TOGGLE
_BUILDER.fields_by_name['auto_builder_dimension'].enum_type = _TOGGLE
_BUILDER.fields_by_name['experimental'].enum_type = _TOGGLE
_BUILDER.fields_by_name['task_template_canary_percentage'].message_type = google_dot_protobuf_dot_wrappers__pb2._UINT32VALUE
_SWARMING.fields_by_name['builder_defaults'].message_type = _BUILDER
_SWARMING.fields_by_name['builders'].message_type = _BUILDER
_SWARMING.fields_by_name['task_template_canary_percentage'].message_type = google_dot_protobuf_dot_wrappers__pb2._UINT32VALUE
_BUCKET.fields_by_name['acls'].message_type = _ACL
_BUCKET.fields_by_name['swarming'].message_type = _SWARMING
_BUILDBUCKETCFG.fields_by_name['buckets'].message_type = _BUCKET
_BUILDBUCKETCFG.fields_by_name['acl_sets'].message_type = _ACLSET
_BUILDBUCKETCFG.fields_by_name['builder_mixins'].message_type = _BUILDER
DESCRIPTOR.message_types_by_name['Acl'] = _ACL
DESCRIPTOR.message_types_by_name['AclSet'] = _ACLSET
DESCRIPTOR.message_types_by_name['Builder'] = _BUILDER
DESCRIPTOR.message_types_by_name['Swarming'] = _SWARMING
DESCRIPTOR.message_types_by_name['Bucket'] = _BUCKET
DESCRIPTOR.message_types_by_name['BuildbucketCfg'] = _BUILDBUCKETCFG
DESCRIPTOR.enum_types_by_name['Toggle'] = _TOGGLE
_sym_db.RegisterFileDescriptor(DESCRIPTOR)

Acl = _reflection.GeneratedProtocolMessageType('Acl', (_message.Message,), dict(
  DESCRIPTOR = _ACL,
  __module__ = 'project_config_pb2'
  # @@protoc_insertion_point(class_scope:buildbucket.Acl)
  ))
_sym_db.RegisterMessage(Acl)

AclSet = _reflection.GeneratedProtocolMessageType('AclSet', (_message.Message,), dict(
  DESCRIPTOR = _ACLSET,
  __module__ = 'project_config_pb2'
  # @@protoc_insertion_point(class_scope:buildbucket.AclSet)
  ))
_sym_db.RegisterMessage(AclSet)

Builder = _reflection.GeneratedProtocolMessageType('Builder', (_message.Message,), dict(

  CacheEntry = _reflection.GeneratedProtocolMessageType('CacheEntry', (_message.Message,), dict(
    DESCRIPTOR = _BUILDER_CACHEENTRY,
    __module__ = 'project_config_pb2'
    # @@protoc_insertion_point(class_scope:buildbucket.Builder.CacheEntry)
    ))
  ,

  Recipe = _reflection.GeneratedProtocolMessageType('Recipe', (_message.Message,), dict(
    DESCRIPTOR = _BUILDER_RECIPE,
    __module__ = 'project_config_pb2'
    # @@protoc_insertion_point(class_scope:buildbucket.Builder.Recipe)
    ))
  ,
  DESCRIPTOR = _BUILDER,
  __module__ = 'project_config_pb2'
  # @@protoc_insertion_point(class_scope:buildbucket.Builder)
  ))
_sym_db.RegisterMessage(Builder)
_sym_db.RegisterMessage(Builder.CacheEntry)
_sym_db.RegisterMessage(Builder.Recipe)

Swarming = _reflection.GeneratedProtocolMessageType('Swarming', (_message.Message,), dict(
  DESCRIPTOR = _SWARMING,
  __module__ = 'project_config_pb2'
  # @@protoc_insertion_point(class_scope:buildbucket.Swarming)
  ))
_sym_db.RegisterMessage(Swarming)

Bucket = _reflection.GeneratedProtocolMessageType('Bucket', (_message.Message,), dict(
  DESCRIPTOR = _BUCKET,
  __module__ = 'project_config_pb2'
  # @@protoc_insertion_point(class_scope:buildbucket.Bucket)
  ))
_sym_db.RegisterMessage(Bucket)

BuildbucketCfg = _reflection.GeneratedProtocolMessageType('BuildbucketCfg', (_message.Message,), dict(
  DESCRIPTOR = _BUILDBUCKETCFG,
  __module__ = 'project_config_pb2'
  # @@protoc_insertion_point(class_scope:buildbucket.BuildbucketCfg)
  ))
_sym_db.RegisterMessage(BuildbucketCfg)


DESCRIPTOR._options = None
# @@protoc_insertion_point(module_scope)
