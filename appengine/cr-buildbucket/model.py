# Copyright 2014 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

import collections
import contextlib
import datetime
import itertools
import random
import zlib

from components import auth
from components import datastore_utils
from components import utils
from google.appengine.api import app_identity
from google.appengine.ext import ndb
from google.appengine.ext.ndb import msgprop
from google.protobuf import struct_pb2
from protorpc import messages

from go.chromium.org.luci.buildbucket.proto import build_pb2
from go.chromium.org.luci.buildbucket.proto import common_pb2
import buildtags
import config
import experiments

BEGINING_OF_THE_WORLD = datetime.datetime(2010, 1, 1, 0, 0, 0, 0)
BUILD_TIMEOUT = datetime.timedelta(days=2)

# For how long to store builds.
BUILD_STORAGE_DURATION = datetime.timedelta(days=30 * 18)  # ~18mo

# If builds weren't scheduled for this duration on a given builder, the
# Builder entity is deleted.
BUILDER_EXPIRATION_DURATION = datetime.timedelta(weeks=4)


class BuildStatus(messages.Enum):
  # Status is not set.
  UNSET = 0
  # A build is created, can be leased by someone and started.
  SCHEDULED = 1
  # Someone has leased the build and marked it as started.
  STARTED = 2
  # A build is completed. See BuildResult for more details.
  COMPLETED = 3


class BuildResult(messages.Enum):
  # Result is not set.
  UNSET = 0
  # A build has completed successfully.
  SUCCESS = 1
  # A build has completed unsuccessfully.
  FAILURE = 2
  # A build was canceled.
  CANCELED = 3


class FailureReason(messages.Enum):
  # Reason is not set.
  UNSET = 0
  # Build failed
  BUILD_FAILURE = 1
  # Something happened within buildbucket.
  BUILDBUCKET_FAILURE = 2
  # Something happened with build infrastructure, but not buildbucket.
  INFRA_FAILURE = 3
  # A build-system rejected a build because its definition is invalid.
  INVALID_BUILD_DEFINITION = 4


class CancelationReason(messages.Enum):
  # Reason is not set.
  UNSET = 0
  # A build was canceled explicitly, probably by an API call.
  CANCELED_EXPLICITLY = 1
  # A build was canceled by buildbucket due to timeout.
  TIMEOUT = 2


class PubSubCallback(ndb.Model):
  """Parameters for a callack push task."""
  topic = ndb.StringProperty(required=True, indexed=False)
  auth_token = ndb.StringProperty(indexed=False)
  user_data = ndb.StringProperty(indexed=False)


class BucketState(ndb.Model):
  """Persistent state of a single bucket."""
  # If True, no new builds may be leased for this bucket.
  is_paused = ndb.BooleanProperty()


def is_terminal_status(status):  # pragma: no cover
  return status not in (
      common_pb2.STATUS_UNSPECIFIED, common_pb2.SCHEDULED, common_pb2.STARTED
  )


class Build(ndb.Model):
  """Describes a build.

  Build key:
    Build keys are autogenerated, monotonically decreasing integers.
    That is, when sorted by key, new builds are first.
    Build has no parent.

    Build id is a 64 bits integer represented as a string to the user.
    - 1 highest order bit is set to 0 to keep value positive.
    - 43 bits are 43 lower bits of bitwise-inverted time since
      BEGINING_OF_THE_WORLD at 1ms resolution.
      It is good for 2**43 / 365.3 / 24 / 60 / 60 / 1000 = 278 years
      or 2010 + 278 = year 2288.
    - 16 bits are set to a random value. Assuming an instance is internally
      consistent with itself, it can ensure to not reuse the same 16 bits in two
      consecutive requests and/or throttle itself to one request per
      millisecond. Using random value reduces to 2**-15 the probability of
      collision on exact same timestamp at 1ms resolution, so a maximum
      theoretical rate of 65536000 requests/sec but an effective rate in the
      range of ~64k qps without much transaction conflicts. We should be fine.
    - 4 bits are 0. This is to represent the 'version' of the entity
      schema.

    The idea is taken from Swarming TaskRequest entity:
    https://source.chromium.org/chromium/_/chromium/infra/luci/luci-py/+/a4a91d5e1e14b8b866b68b68bc1055b0b8ffef3b:appengine/swarming/server/task_request.py;l=1380-1404
  """

  # ndb library sometimes silently ignores memcache errors
  # => memcache is not synchronized with datastore
  # => a build never finishes from the app code perspective
  # => builder is stuck for days.
  # We workaround this problem by setting a timeout.
  _memcache_timeout = 600  # 10m

  @classmethod
  def _use_memcache(cls, _):
    # See main.py for reasons why memcache is being disabled.
    return False

  # Stores the build proto. The primary property of this entity.
  # Majority of the other properties are either derivatives of this field or
  # legacy.
  #
  # Does not include:
  #   output.properties: see BuildOutputProperties
  #   steps: see BuildSteps.
  #   tags: stored in tags attribute, because we have to index them anyway.
  #   input.properties: see BuildInputProperties.
  #     CAVEAT: field input.properties does exist during build creation, and
  #     moved into BuildInputProperties right before initial datastore.put.
  #   infra: see BuildInfra.
  #     CAVEAT: field infra does exist during build creation, and moved into
  #     BuildInfra right before initial datastore.put.
  #
  # Transition period: proto is either None or complete, i.e. created by
  # creation.py or fix_builds.py.
  proto = datastore_utils.ProtobufProperty(build_pb2.Build)

  # Update token required to finalize the invocation.
  # Set at invocation creation time.
  resultdb_update_token = ndb.StringProperty(indexed=False)

  # Build token required to update the build.
  # Set at build creation time.
  update_token = ndb.StringProperty(indexed=False)

  # == proto-derived properties ================================================
  #
  # These properties are derived from "proto" properties.
  # They are used to index builds.

  status = ndb.ComputedProperty(
      lambda self: self.proto.status, name='status_v2'
  )

  @property
  def is_ended(self):  # pragma: no cover
    return is_terminal_status(self.proto.status)

  incomplete = ndb.ComputedProperty(lambda self: not self.is_ended)

  # ID of the LUCI project to which this build belongs.
  project = ndb.ComputedProperty(lambda self: self.proto.builder.project)

  # Indexed string "<project_id>/<bucket_name>".
  # Example: "chromium/try".
  # Prefix "luci.<project_id>." is stripped from bucket name.
  bucket_id = ndb.ComputedProperty(
      lambda self: config.format_bucket_id(
          self.proto.builder.project, self.proto.builder.bucket))

  # Indexed string "<project_id>/<bucket_name>/<builder_name>".
  # Example: "chromium/try/linux-rel".
  # Prefix "luci.<project_id>." is stripped from bucket name.
  builder_id = ndb.ComputedProperty(
      lambda self: config.builder_id_string(self.proto.builder)
  )

  canary = ndb.ComputedProperty(lambda self: self.proto.canary)

  # Value of proto.create_time.
  # Making this property computed is not-entirely trivial because
  # ComputedProperty saves it as int, as opposed to datetime.datetime.
  # TODO(nodir): remove usages of create_time indices, rely on build id ordering
  # instead.
  create_time = ndb.DateTimeProperty()

  # A list of colon-separated key-value pairs. Indexed.
  # Used to populate tags in builds_to_protos_async, if requested.
  tags = ndb.StringProperty(repeated=True)

  # If True, the build won't affect monitoring and won't be surfaced in
  # search results unless explicitly requested.
  experimental = ndb.ComputedProperty(
      lambda self: self.proto.input.experimental
  )

  # A list of experiments enabled or disabled on this build.
  # Each entry should look like "[-+]$experiment_name".
  experiments = ndb.StringProperty(repeated=True)

  # Value of proto.created_by.
  # Making this property computed is not-entirely trivial because
  # ComputedProperty saves it as string, but IdentityProperty stores it
  # as a blob property.
  created_by = auth.IdentityProperty()

  is_luci = ndb.BooleanProperty()

  @property
  def realm(self):  # pragma: no cover
    return '%s:%s' % (self.proto.builder.project, self.proto.builder.bucket)

  # == Legacy properties =======================================================

  status_legacy = msgprop.EnumProperty(
      BuildStatus, default=BuildStatus.SCHEDULED, name='status'
  )

  status_changed_time = ndb.DateTimeProperty(auto_now_add=True)

  # immutable arbitrary build parameters.
  parameters = datastore_utils.DeterministicJsonProperty(json_type=dict)

  # PubSub message parameters for build status change notifications.
  # TODO(nodir): replace with notification_pb2.NotificationConfig.
  pubsub_callback = ndb.StructuredProperty(PubSubCallback, indexed=False)

  # id of the original build that this build was derived from.
  retry_of = ndb.IntegerProperty()

  # a URL to a build-system-specific build, viewable by a human.
  url = ndb.StringProperty(indexed=False)

  # V1 status properties. Computed by _pre_put_hook and _post_get_hook.
  # TODO(crbug/1090540): Remove _pre_put_hook computation of legacy properties.
  result = msgprop.EnumProperty(BuildResult)
  result_details = datastore_utils.DeterministicJsonProperty(json_type=dict)
  cancelation_reason = msgprop.EnumProperty(CancelationReason)
  failure_reason = msgprop.EnumProperty(FailureReason)

  # Lease-time properties.

  # TODO(nodir): move lease to a separate entity under Build.
  # It would be more efficient.
  # current lease expiration date.
  # The moment the build is leased, |lease_expiration_date| is set to
  # (utcnow + lease_duration).
  lease_expiration_date = ndb.DateTimeProperty()
  # None if build is not leased, otherwise a random value.
  # Changes every time a build is leased. Can be used to verify that a client
  # is the leaseholder.
  lease_key = ndb.IntegerProperty(indexed=False)
  # True if the build is currently leased. Otherwise False
  is_leased = ndb.ComputedProperty(lambda self: bool(self.lease_key))
  leasee = auth.IdentityProperty()
  never_leased = ndb.BooleanProperty()

  # ============================================================================

  def _pre_put_hook(self):
    """Checks Build invariants before putting."""
    super(Build, self)._pre_put_hook()

    config.validate_project_id(self.proto.builder.project)
    config.validate_bucket_name(self.proto.builder.bucket)

    self.update_v1_status_fields()
    self.proto.update_time.FromDatetime(utils.utcnow())

    is_started = self.proto.status == common_pb2.STARTED
    is_ended = self.is_ended

    # See note in _post_get_hook
    if self.lease_key == 0:  # pragma: no cover
      self.lease_key = None
      self.lease_expiration_date = None
      self.leasee = None

    is_leased = self.lease_key is not None

    assert not (is_ended and is_leased)
    assert (self.lease_expiration_date is not None) == is_leased
    assert (self.leasee is not None) == is_leased

    tag_delm = buildtags.DELIMITER
    assert not self.tags or all(tag_delm in t for t in self.tags)

    assert not self.experiments or (
        all(exp.startswith((
            '+',
            '-',
        )) for exp in self.experiments)
    )

    assert self.proto.HasField('create_time')
    assert self.proto.HasField('end_time') == is_ended
    assert not is_started or self.proto.HasField('start_time')

    def _ts_less(ts1, ts2):
      return ts1.seconds and ts2.seconds and ts1.ToDatetime() < ts2.ToDatetime()

    assert not _ts_less(self.proto.start_time, self.proto.create_time)
    assert not _ts_less(self.proto.end_time, self.proto.create_time)
    assert not _ts_less(self.proto.end_time, self.proto.start_time)

    self.tags = sorted(set(self.tags))
    self.experiments = sorted(set(self.experiments))

  @classmethod
  def _post_get_hook(cls, key, future):
    """Computes v1 legacy fields."""
    build = future.get_result()
    if build:
      build.update_v1_status_fields()

      # NOTE: When Go writes Build, it assigns these fields 0-values, i.e. 0 and
      # empty-string, rather than the ndb None value. Since these fields are
      # legacy related to BBv1 we just have a bit of a hack here so that if we
      # have to write these from python, we write them with the expected values.
      #
      # We expect that any Go code which processes these will assume zero-value
      # equals unset, meaning that the None values here are just to placate
      # Python assumptions.
      if build.lease_key == 0:  # pragma: no cover
        build.lease_key = None
        build.lease_expiration_date = None
        build.leasee = None
    else:  # pragma: no cover
      pass

  def update_v1_status_fields(self):
    """Updates V1 status fields."""
    # Reset these to None instead of UNSET for backwards compatibility.
    # UNSET is only required to load values written by the Go service,
    # the Python service expects these to be None when they aren't set.
    self.status_legacy = None
    self.result = None
    self.failure_reason = None
    self.cancelation_reason = None

    status = self.proto.status
    if status == common_pb2.SCHEDULED:
      self.status_legacy = BuildStatus.SCHEDULED
    elif status == common_pb2.STARTED:
      self.status_legacy = BuildStatus.STARTED
    elif status == common_pb2.SUCCESS:
      self.status_legacy = BuildStatus.COMPLETED
      self.result = BuildResult.SUCCESS
    elif status == common_pb2.FAILURE:
      self.status_legacy = BuildStatus.COMPLETED
      self.result = BuildResult.FAILURE
      self.failure_reason = FailureReason.BUILD_FAILURE
    elif status == common_pb2.INFRA_FAILURE:
      self.status_legacy = BuildStatus.COMPLETED
      if self.proto.status_details.HasField('timeout'):
        self.result = BuildResult.CANCELED
        self.cancelation_reason = CancelationReason.TIMEOUT
      else:
        self.result = BuildResult.FAILURE
        self.failure_reason = FailureReason.INFRA_FAILURE
    elif status == common_pb2.CANCELED:
      self.status_legacy = BuildStatus.COMPLETED
      self.result = BuildResult.CANCELED
      self.cancelation_reason = CancelationReason.CANCELED_EXPLICITLY
    else:  # pragma: no cover
      assert False, status

  def regenerate_lease_key(self):
    """Changes lease key to a different random int."""
    while True:
      new_key = random.randint(1, 1 << 31)
      if new_key != self.lease_key:  # pragma: no branch
        self.lease_key = new_key
        break

  def clear_lease(self):  # pragma: no cover
    """Clears build's lease attributes."""
    self.lease_key = None
    self.lease_expiration_date = None
    self.leasee = None

  def tags_to_protos(self, dest):
    """Adds non-hidden self.tags to a repeated StringPair container."""
    for t in self.tags:
      k, v = buildtags.parse(t)
      if k not in buildtags.HIDDEN_TAG_KEYS:
        dest.add(key=k, value=v)

  @property
  def uses_realms(self):  # pragma: no cover
    """True if the build opted-in into using LUCI Realms."""
    return '+%s' % (experiments.USE_REALMS,) in self.experiments


class BuildDetailEntity(ndb.Model):
  """A base class for a Datastore entity that stores some details of one Build.

  Entity key: Parent is Build entity key. ID is 1.
  """

  @classmethod
  def _use_memcache(cls, _):
    # Fallback in case memcache is still enabled (observed in v1 API handlers).
    # See main.py for reasons why memcache is being disabled.
    return not app_identity.get_application_id().endswith('-dev')

  @classmethod
  def key_for(cls, build_key):  # pragma: no cover
    return ndb.Key(cls, 1, parent=build_key)


class BuildProperties(BuildDetailEntity):
  """Base class for storing build input/output properties."""

  # google.protobuf.Struct message in binary form.
  properties = ndb.BlobProperty()

  def parse(self):  # pragma: no cover
    s = struct_pb2.Struct()
    s.ParseFromString(self.properties or '')
    return s

  def serialize(self, struct):  # pragma: no cover
    assert isinstance(struct, struct_pb2.Struct)
    self.properties = struct.SerializeToString()


class BuildInputProperties(BuildProperties):
  """Stores buildbucket.v2.Build.input.properties."""


class BuildOutputProperties(BuildProperties):
  """Stores buildbucket.v2.Build.output.properties."""

  @classmethod
  def _use_memcache(cls, _):
    # See main.py for reasons why memcache is being disabled.
    return False


class BuildSteps(BuildDetailEntity):
  """Stores buildbucket.v2.Build.steps."""

  # max length of step_container_bytes attribute.
  MAX_STEPS_LEN = 1e6

  # buildbucket.v2.Build binary protobuf message with only "steps" field set.
  # zlib-compressed if step_container_bytes_zipped is True.
  step_container_bytes = ndb.BlobProperty(name='steps')

  # Whether step_container_bytes are zlib-compressed.
  # We don't reuse ndb compression because we want to enforce the size limit
  # at the API level after compression.
  step_container_bytes_zipped = ndb.BooleanProperty(indexed=False)

  def _pre_put_hook(self):
    """Checks BuildSteps invariants before putting."""
    super(BuildSteps, self)._pre_put_hook()
    assert self.step_container_bytes is not None
    assert len(self.step_container_bytes) <= self.MAX_STEPS_LEN

  @classmethod
  def make(cls, build_proto):
    """Creates BuildSteps for the build_proto.

    Does not verify step size.
    """
    assert build_proto.id
    build_key = ndb.Key(Build, build_proto.id)
    ret = cls(key=cls.key_for(build_key))
    ret.write_steps(build_proto)
    return ret

  def write_steps(self, build_proto):
    """Serializes build_proto.steps into self."""
    container = build_pb2.Build(steps=build_proto.steps)
    container_bytes = container.SerializeToString()

    # Compress only if necessary.
    zipped = len(container_bytes) > self.MAX_STEPS_LEN
    if zipped:
      container_bytes = zlib.compress(container_bytes)

    self.step_container_bytes = container_bytes
    self.step_container_bytes_zipped = zipped

  def read_steps(self, build_proto):
    """Deserializes steps into build_proto.steps."""
    container_bytes = self.step_container_bytes
    if self.step_container_bytes_zipped:
      container_bytes = zlib.decompress(container_bytes)
    build_proto.ClearField('steps')
    build_proto.MergeFromString(container_bytes)

  @classmethod
  @ndb.tasklet
  def cancel_incomplete_steps_async(cls, build_id, end_ts):
    """Marks incomplete steps as canceled in the Datastore, if any."""
    assert end_ts.seconds
    assert ndb.in_transaction()
    entity = yield cls.key_for(ndb.Key(Build, build_id)).get_async()
    if not entity:
      return

    container = build_pb2.Build()
    entity.read_steps(container)

    changed = False
    for s in container.steps:
      if not is_terminal_status(s.status):
        s.status = common_pb2.CANCELED
        s.end_time.CopyFrom(end_ts)
        changed = True

    if changed:  # pragma: no branch
      entity.write_steps(container)
      yield entity.put_async()

  @classmethod
  def _use_memcache(cls, _):
    # See main.py for reasons why memcache is being disabled.
    return False


class BuildInfra(BuildDetailEntity):
  """Stores buildbucket.v2.Build.infra."""

  # buildbucket.v2.Build.infra serialized to bytes.
  infra = ndb.BlobProperty()

  def parse(self):  # pragma: no cover
    """Deserializes infra."""
    ret = build_pb2.BuildInfra()
    ret.ParseFromString(self.infra)
    return ret

  @contextlib.contextmanager
  def mutate(self):  # pragma: no cover
    """Returns a context manager that provides a mutable BuildInfra proto.

    Deserializes infra, yields it, and serializes back.
    """
    proto = self.parse()
    yield proto
    self.infra = proto.SerializeToString()


# Tuple of classes representing entity kinds that living under Build entity.
# Such entities must be deleted if Build entity is deleted.
BUILD_CHILD_CLASSES = (
    BuildInfra,
    BuildInputProperties,
    BuildOutputProperties,
    BuildSteps,
)

BuildBundleBase = collections.namedtuple(
    'BuildBundleBase',
    [
        'build',  # instance of Build
        'infra',  # instance of BuildInfra
        'input_properties',  # instance of BuildInputProperties
        'output_properties',  # instance of BuildOutputProperties
        'steps',  # instance of BuildSteps
    ]
)


class BuildBundle(BuildBundleBase):
  """A tuple of entities describing one build."""

  def __new__(
      cls,
      build,
      infra=None,
      input_properties=None,
      output_properties=None,
      steps=None
  ):
    assert isinstance(build, Build), build
    assert not infra or isinstance(infra, BuildInfra)
    assert (
        not input_properties or
        isinstance(input_properties, BuildInputProperties)
    )
    assert (
        not output_properties or
        isinstance(output_properties, BuildOutputProperties)
    )
    assert not steps or isinstance(steps, BuildSteps)
    return super(BuildBundle, cls).__new__(
        cls,
        build=build,
        infra=infra,
        input_properties=input_properties,
        output_properties=output_properties,
        steps=steps,
    )

  @classmethod
  @ndb.tasklet
  def get_async(
      cls,
      build,
      infra=False,
      output_properties=False,
      input_properties=False,
      steps=False
  ):
    """Fetches a BuildBundle.

    build must be either an int/long build id or pre-fetched Build.
    If it is an id, the build will be fetched.
    If not found, this function returns None future.
    """
    assert isinstance(build, (int, long, Build)), build

    if isinstance(build, Build):
      build_key = build.key
      build_fut = ndb.Future()
      build_fut.set_result(build)
    else:
      build_key = ndb.Key(Build, build)
      build_fut = build_key.get_async()

    def get_child_async(do_load, clazz):
      if not do_load:
        f = ndb.Future()
        f.set_result(None)
        return f
      return clazz.key_for(build_key).get_async()

    kwarg_futs = dict(
        build=build_fut,
        infra=get_child_async(infra, BuildInfra),
        input_properties=get_child_async(
            input_properties, BuildInputProperties
        ),
        output_properties=get_child_async(
            output_properties, BuildOutputProperties
        ),
        steps=get_child_async(steps, BuildSteps),
    )
    build = yield build_fut
    if not build:  # pragma: no cover
      raise ndb.Return(None)

    kwargs = {}
    for k, f in kwarg_futs.iteritems():
      kwargs[k] = yield f
    raise ndb.Return(cls(**kwargs))

  @classmethod
  def get(cls, *args, **kwargs):
    return cls.get_async(*args, **kwargs).get_result()

  @ndb.tasklet
  def put_async(self):
    """Puts all non-None entities."""
    yield ndb.put_multi_async(filter(None, self))

  def put(self):
    return self.put_async().get_result()

  def to_proto(self, dest, load_tags):
    """Writes build to the dest Build proto. Returns dest."""
    if dest is not self.build.proto:  # pragma: no branch
      dest.CopyFrom(self.build.proto)
    dest.id = self.build.key.id()  # old builds do not have id field

    if load_tags:
      self.build.tags_to_protos(dest.tags)

    if self.infra:
      dest.infra.ParseFromString(self.infra.infra)

    if self.steps:
      self.steps.read_steps(dest)

    if self.input_properties:
      dest.input.properties.ParseFromString(self.input_properties.properties)
    if self.output_properties:
      dest.output.properties.ParseFromString(self.output_properties.properties)
    return dest


class Builder(ndb.Model):
  """A builder in a bucket.

  Used internally for metrics.
  Registered automatically by scheduling a build.
  Unregistered automatically by not scheduling builds for
  BUILDER_EXPIRATION_DURATION.

  Entity key:
    No parent. ID is a string with format "{project}:{bucket}:{builder}".
  """

  # Last time we received a valid build scheduling request for this builder.
  # Probabilistically updated by services.py, see its _should_update_builder.
  last_scheduled = ndb.DateTimeProperty()

  @classmethod
  def make_key(cls, builder_id):  # pragma: no cover
    bid = builder_id
    return ndb.Key(cls, '%s:%s:%s' % (bid.project, bid.bucket, bid.builder))


_TIME_RESOLUTION = datetime.timedelta(milliseconds=1)
_BUILD_ID_SUFFIX_LEN = 20
# Size of a build id segment covering one millisecond.
ONE_MS_BUILD_ID_RANGE = 1 << _BUILD_ID_SUFFIX_LEN


def _id_time_segment(dtime):
  assert dtime
  assert dtime >= BEGINING_OF_THE_WORLD
  delta = dtime - BEGINING_OF_THE_WORLD
  now = int(delta.total_seconds() * 1000.)
  return (~now & ((1 << 43) - 1)) << 20


def create_build_ids(dtime, count, randomness=True):
  """Returns a range of valid build ids, as integers and based on a datetime.

  See Build's docstring, "Build key" section.
  """
  # Build ID bits: "0N{43}R{16}V{4}"
  # where N is now bits, R is random bits and V is version bits.
  build_id = int(_id_time_segment(dtime))
  build_id = build_id | ((random.getrandbits(16) << 4) if randomness else 0)
  # Subtract so that ids are descending.
  return [build_id - i * (1 << 4) for i in xrange(count)]


def build_id_range(create_time_low, create_time_high):
  """Converts a creation time range to build id range.

  Low/high bounds are inclusive/exclusive respectively, for both time and id
  ranges.
  """
  id_low = None
  id_high = None
  if create_time_low is not None:  # pragma: no branch
    # convert inclusive to exclusive
    id_high = _id_time_segment(create_time_low - _TIME_RESOLUTION)
  if create_time_high is not None:  # pragma: no branch
    # convert exclusive to inclusive
    id_low = _id_time_segment(create_time_high - _TIME_RESOLUTION)
  return id_low, id_high


@ndb.tasklet
def builds_to_protos_async(
    builds,
    load_tags,
    load_input_properties,
    load_output_properties,
    load_steps,
    load_infra,
):
  """Converts Build objects to build_pb2.Build messages.

  builds must be a list of (model.Build, build_pb2.Build) tuples,
  where the build_pb2.Build is the destination.
  """

  bundle_futs = [(
      dest,
      BuildBundle.get_async(
          b,
          infra=load_infra,
          input_properties=load_input_properties,
          output_properties=load_output_properties,
          steps=load_steps
      )
  ) for b, dest in builds]

  for dest, bundle_fut in bundle_futs:
    bundle = yield bundle_fut
    bundle.to_proto(dest, load_tags=load_tags)
