# Copyright 2014 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

"""Gory implementation details of CAS service.

Append only currently. Once object is added, it can't be removed.

Upload protocol:
  1) Client asks the server to initiate a new upload session (for hash X).
  2) Server starts Resumable Upload protocol to a temporary file in GS.
  3) Client uploads data to this temporary file (using resumable upload_id for
     authentication).
  4) Client finalizes the upload (thus making the temp file visible).
  5) Client notifies server that upload has finished.
  6) Server starts hash verification task.
  7) On successful verification, server copies the file to the final location.
  8) Meanwhile client polls server for verification operation status.
  9) Once verification finishes, client polls 'PUBLISHED' status.
    a) Concurrent uploads of a same file are fine, upload session moves to
       PUBLISHED state whenever corresponding hash becomes available in
       the store, regardless of who exactly uploaded it.

Features of Google Storage used or taken into consideration:
  * upload_id is enough to authenticate the request (no access_token needed).
  * upload_id is NOT consumed when upload is finalized and may be reused.
  * Each object has ETag that identified its content.
  * There's copy-object-if-etag-matches atomic operation.
  * Lifecycle management for temporary files, to cleanup old garbage.

Also this module is sensitive to implementation details of 'cloudstorage'
library since it uses its non-public APIs:
  * StreamingBuffer._api.api_url and StreamingBuffer._path_with_token.
  * ReadBuffer._etag.
  * storage_api._get_storage_api(...) and _StorageApi it returns.
"""

import base64
import hashlib
import logging
import random
import re
import string
import urllib
import webapp2

from google.appengine import runtime
from google.appengine.api import app_identity
from google.appengine.api import datastore_errors
from google.appengine.ext import ndb

# We use cloud storage guts to implement "copy if ETag matches".
import cloudstorage
from cloudstorage import api_utils
from cloudstorage import errors
from cloudstorage import storage_api

from components import auth
from components import decorators
from components import utils

import config

# TODO(vadimsh): Garbage collect expired UploadSession. Right know only public
# upload_session_id expires, rendering sessions unreachable by clients. But the
# entities themselves unnecessarily stay in the datastore.

# How long to keep signed fetch URL alive.
FETCH_URL_EXPIRATION_SEC = 60 * 60

# How long to keep pending upload session alive.
SESSION_EXPIRATION_TIME_SEC = 6 * 60 * 60

# Chunks to read when verifying the hash.
READ_BUFFER_SIZE = 1024 * 1024

# Hash algorithms we are willing to accept: name -> (factory, hex digest len).
SUPPORTED_HASH_ALGOS = {
  'SHA1': (hashlib.sha1, 40),
}

# Return values of task queue task handling function.
TASK_DONE = 1
TASK_RETRY = 2


def is_supported_hash_algo(hash_algo):
  """True if given algorithm is supported by CAS service."""
  return hash_algo in SUPPORTED_HASH_ALGOS


def is_valid_hash_digest(hash_algo, hash_digest):
  """True if given hex digest looks like a valid hex digest for given algo."""
  assert is_supported_hash_algo(hash_algo)
  _, digest_size = SUPPORTED_HASH_ALGOS[hash_algo]
  return re.match('^[0-9a-f]{%d}$' % digest_size, hash_digest)


def get_cas_service():
  """Factory method that returns configured CASService instance.

  If the service is not configured, returns None. Also acts as a mocking point
  for unit tests.
  """
  conf = config.cached()
  if not conf.cas_gs_path or not conf.cas_gs_temp:
    return None
  try:
    cloudstorage.validate_file_path(conf.cas_gs_path.rstrip('/'))
    cloudstorage.validate_file_path(conf.cas_gs_temp.rstrip('/'))
  except ValueError as err:
    logging.error("Invalid CAS config: %s", err)
    return None
  return CASService(
      conf.cas_gs_path.rstrip('/'),
      conf.cas_gs_temp.rstrip('/'))


class NotFoundError(Exception):
  """Raised by 'open' when the file is not in CAS."""


class UploadIdSignature(auth.TokenKind):
  """Token to use to generate and validate signed upload_session_id."""
  expiration_sec = SESSION_EXPIRATION_TIME_SEC
  secret_key = auth.SecretKey('upload_session_id_signing')
  version = 1


class CASService(object):
  """CAS implementation on top of Google Storage.

  It uses GAE App Identity API to sign Google Storage URLs.
  """

  def __init__(self, gs_path, gs_temp):
    self._gs_path = gs_path.rstrip('/')
    self._gs_temp = gs_temp.rstrip('/')
    self._retry_params = api_utils.RetryParams()
    self._app_identity_sign_blob = app_identity.sign_blob  # mocked in tests
    cloudstorage.validate_file_path(self._gs_path)
    cloudstorage.validate_file_path(self._gs_temp)

  def is_object_present(self, hash_algo, hash_digest):
    """True if the given object is in the store."""
    assert is_valid_hash_digest(hash_algo, hash_digest)
    return self._is_gs_file_present(
        self._verified_gs_path(hash_algo, hash_digest))

  def generate_fetch_url(self, hash_algo, hash_digest, filename=None):
    """Returns a signed URL that can be used to fetch an object.

    See https://developers.google.com/storage/docs/accesscontrol#Signed-URLs
    for more info about signed URLs.
    """
    assert is_valid_hash_digest(hash_algo, hash_digest)

    # Generate the signature for GET operation.
    gs_path = self._verified_gs_path(hash_algo, hash_digest)
    signature, expires = self._sign_gs_fetch(gs_path)

    # Query parameters for the final URL.
    params = [
      ('GoogleAccessId', utils.get_service_account_name()),
      ('Expires', str(expires)),
      ('Signature', base64.b64encode(signature)),
    ]

    # Oddly, response-content-disposition is not signed and can be slapped onto
    # existing signed URL.
    if filename:
      assert '"' not in filename, filename
      params.append((
          'response-content-disposition',
          'attachment; filename="%s"' % filename))

    # Generate the final URL.
    query_params = urllib.urlencode(params)
    assert gs_path.startswith('/'), gs_path
    return 'https://storage.googleapis.com%s?%s' % (gs_path, query_params)

  def open(self, hash_algo, hash_digest, read_buffer_size=None):
    """Opens a file in CAS for reading.

    Args:
      hash_algo: valid supported hash algorithm to use for verification.
      hash_digest: hex hash digest of the content to be uploaded.
      read_buffer_size: length of chunk of data to read with each RPC.

    Returns:
      File-like object, caller takes ownership and should close it.

    Raises:
      NotFoundError if file is missing.
    """
    read_buffer_size = read_buffer_size or READ_BUFFER_SIZE
    try:
      return cloudstorage.open(
          filename=self._verified_gs_path(hash_algo, hash_digest),
          mode='r',
          read_buffer_size=read_buffer_size,
          retry_params=self._retry_params)
    except cloudstorage.NotFoundError:
      raise NotFoundError()

  def start_direct_upload(self, hash_algo):
    """Can be used to upload data to CAS directly from an Appengine handler.

    Opens a temp file for writing (and returns wrapper around it). Hashes the
    data while it is being written, and moves the temp file to an appropriate
    location in CAS once it is closed.

    Args:
      hash_algo: algorithm to use to calculate data hash.

    Returns:
      DirectUpload object to write data to.
    """
    assert is_supported_hash_algo(hash_algo)
    ts_sec = utils.datetime_to_timestamp(utils.utcnow()) / 1000000.
    temp_path = self._temp_direct_upload_gs_path(ts_sec)
    temp_file = cloudstorage.open(
        filename=temp_path,
        mode='w',
        retry_params=self._retry_params)
    def commit_callback(hash_digest, commit):
      if commit:
        self._gs_copy(temp_path, self._verified_gs_path(hash_algo, hash_digest))
      self._gs_delete(temp_path)
    return DirectUpload(
        file_obj=temp_file,
        hasher=SUPPORTED_HASH_ALGOS[hash_algo][0](),
        callback=commit_callback)

  def create_upload_session(self, hash_algo, hash_digest, caller):
    """Starts a new upload operation.

    Starts an upload regardless of whether the object is already stored or not.
    Generates upload_url for GS resumable upload protocol.

    Args:
      hash_algo: valid supported hash algorithm to use for verification.
      hash_digest: hex hash digest of the content to be uploaded.
      caller: auth.Identity of whoever makes the request.

    Returns:
      tuple (UploadSession object, signed upload session ID).
    """
    assert is_valid_hash_digest(hash_algo, hash_digest)

    # TODO(vadimsh): Check that number of pending uploads opened by |caller|
    # is low enough. To prevent malicious client from creating tons of uploads.

    # New unique ID (long int).
    upload_id = UploadSession.allocate_ids(size=1)[0]

    # Opening a GCS file and not closing it keeps upload session active.
    ts_sec = utils.datetime_to_timestamp(utils.utcnow()) / 1000000.
    temp_gs_location = self._temp_upload_session_gs_path(upload_id, ts_sec)
    temp_file = cloudstorage.open(
        filename=temp_gs_location,
        mode='w',
        retry_params=self._retry_params)

    # See cloudstorage/storage_api.py, StreamingBuffer for _path_with_token.
    upload_url = '%s%s' % (temp_file._api.api_url, temp_file._path_with_token)

    # New session.
    upload_session = UploadSession(
        id=upload_id,
        hash_algo=hash_algo,
        hash_digest=hash_digest,
        temp_gs_location=temp_gs_location,
        final_gs_location=self._verified_gs_path(hash_algo, hash_digest),
        upload_url=upload_url,
        status=UploadSession.STATUS_UPLOADING,
        created_by=caller)
    upload_session.put()

    # Generate signed ID. It will be usable only by |caller|.
    upload_session_id = UploadIdSignature.generate(
        message=[caller.to_bytes()],
        embedded={'id': '%s' % upload_session.key.id()})
    return upload_session, upload_session_id

  def fetch_upload_session(self, upload_session_id, caller):
    """Returns an existing non-expired upload session given its signed ID.

    Args:
      upload_session_id: signed upload session ID, see create_upload_session.
      caller: auth.Identity of whoever makes the request.

    Returns:
      UploadSession object, or None if session is expired, missing or signature
      is not valid.
    """
    try:
      # Verify the signature, extract upload_id.
      embedded = UploadIdSignature.validate(
          upload_session_id, [caller.to_bytes()])
      upload_id = long(embedded['id'])
    except (auth.InvalidTokenError, KeyError, ValueError):
      logging.error('Using invalid or expired upload_session_id')
      return None
    return UploadSession.get_by_id(upload_id)

  def maybe_finish_upload(self, upload_session):
    """Called whenever a client checks the status of the upload session.

    Args:
      upload_session: UploadSession object.

    Returns:
      Updated UploadSession object.
    """
    # Fast check before starting the transaction.
    if upload_session.status != UploadSession.STATUS_UPLOADING:
      return upload_session

    # Move to VERIFYING state, adding the verification task.
    @ndb.transactional
    def run():
      refreshed = upload_session.key.get()
      if refreshed.status != UploadSession.STATUS_UPLOADING:  # pragma: no cover
        return refreshed
      success = utils.enqueue_task(
          url='/internal/taskqueue/cas-verify/%d' % refreshed.key.id(),
          queue_name='cas-verify',
          transactional=True)
      if not success:  # pragma: no cover
        raise datastore_errors.TransactionFailedError()
      refreshed.status = UploadSession.STATUS_VERIFYING
      refreshed.put()
      return refreshed

    return run()

  def verify_pending_upload(self, unsigned_upload_id):
    """Task queue task that checks the hash of a pending upload, finalizes it.

    Args:
      unsigned_upload_id: long int ID of upload session to check.

    Returns:
      TASK_RETRY if task should be retried, TASK_DONE if not.
    """
    upload_session = UploadSession.get_by_id(unsigned_upload_id)
    if upload_session is None:
      logging.error('Verifying missing upload session:\n%d', unsigned_upload_id)
      return TASK_DONE
    if upload_session.status != UploadSession.STATUS_VERIFYING:
      return TASK_DONE

    # Moves upload_session to some state if it's still in VERIFYING state.
    @ndb.transactional
    def set_status(status, error_message=None):
      refreshed = upload_session.key.get()
      if refreshed.status != UploadSession.STATUS_VERIFYING:  # pragma: no cover
        return False, refreshed
      refreshed.status = status
      refreshed.error_message = error_message
      refreshed.put()
      return True, refreshed

    # Moves upload_session to ERROR state, logs the error, cleans temp files.
    def set_error(error_message):
      changed, _ = set_status(UploadSession.STATUS_ERROR, error_message)
      if not changed:  # pragma: no cover
        return False
      logging.error('\n'.join([
        'CAS upload verification failed.',
        error_message.strip(),
        'upload_id=%d' % unsigned_upload_id,
        'uploader=%s' % upload_session.created_by.to_bytes(),
      ]))
      self._cleanup_temp(upload_session)
      return True

    # Maybe someone else uploaded (and verified) the resulting file already?
    if self._is_gs_file_present(upload_session.final_gs_location):
      self._cleanup_temp(upload_session)
      set_status(UploadSession.STATUS_PUBLISHED)
      return TASK_DONE

    # Client MUST finalize GS upload before invoking verification. If client
    # fails to do so, abort the protocol. Also 'cloudstorage.open' verifies
    # that file is not modified midway by checking ETag with each request. We
    # then perform conditional copy to the final destination using this ETag.
    try:
      temp_file = cloudstorage.open(
          filename=upload_session.temp_gs_location,
          mode='r',
          read_buffer_size=READ_BUFFER_SIZE,
          retry_params=self._retry_params)
    except errors.NotFoundError:
      set_error('Google Storage upload wasn\'t finalized.')
      return TASK_DONE

    # ETag MUST be available by this moment, since cloudstorage.open() performs
    # blocking HEAD call. Also for some weird reason _etag is wrapped in "".
    etag = temp_file._etag.strip('"')
    assert etag

    # TODO(vadimsh): This will always timeout for large files (gigabytes).
    # Verification task can be split into multiple subsequent tasks with the
    # hasher internal state (and file offset and ETag) transported between them.
    try:
      hasher = SUPPORTED_HASH_ALGOS[upload_session.hash_algo][0]()
      while True:
        buf = temp_file.read(READ_BUFFER_SIZE)
        if not buf:
          break
        hasher.update(buf)
        # Help GC to collect this buffer before new one is allocated. Appengine
        # is very memory constrained environment.
        del buf
      digest = hasher.hexdigest()
    finally:
      temp_file.close()
      del temp_file

    # Moment of truth.
    if upload_session.hash_digest != digest:
      set_error(
          'Invalid %s hash: expected %s, got %s.' %
          (upload_session.hash_algo, upload_session.hash_digest, digest))
      return TASK_DONE

    # Copy to the final destination verifying the ETag is the same.
    try:
      self._gs_copy(
          src=upload_session.temp_gs_location,
          dst=upload_session.final_gs_location,
          src_etag=etag)
    except errors.NotFoundError:  # pragma: no cover
      # Probably some concurrent finalization removed temp_gs_location already.
      # Retry the task to check this.
      return TASK_RETRY
    except errors.FatalError as err:  # pragma: no cover
      # Precondition failed, i.e. temp file was modified after verification.
      set_error(str(err))
      return TASK_DONE

    # Everything is in place. Cleanup temp garbage.
    set_status(upload_session.STATUS_PUBLISHED)
    self._cleanup_temp(upload_session)
    return TASK_DONE

  def _verified_gs_path(self, hash_algo, hash_digest):
    """Google Storage path to a verified file."""
    return str('%s/%s/%s' % (self._gs_path, hash_algo, hash_digest))

  def _temp_upload_session_gs_path(self, upload_id, timestamp_sec):
    """Path to a temporary drop location for upload sessions."""
    # Use timestamp prefix to enable cheap and dirty "time range" queries when
    # listing bucket with a prefix scan.
    return str('%s/%d_%s' % (self._gs_temp, timestamp_sec, upload_id))

  def _temp_direct_upload_gs_path(self, timestamp_sec):
    """Randomized path to a temporary drop location for direct uploads."""
    rnd_str = ''.join(random.choice(string.ascii_letters) for i in range(20))
    return str('%s/%d_direct_%s' % (self._gs_temp, timestamp_sec, rnd_str))

  def _is_gs_file_present(self, gs_path):
    """True if given GS file exists."""
    try:
      cloudstorage.stat(
          filename=gs_path,
          retry_params=self._retry_params)
    except cloudstorage.NotFoundError:
      return False
    return True

  def _gs_copy(self, src, dst, src_etag=None):  # pragma: no cover
    """Copy |src| file to |dst| optionally checking src ETag.

    Raises cloudstorage.FatalError on precondition error.
    """
    # See cloudstorage.cloudstorage_api._copy2.
    cloudstorage.validate_file_path(src)
    cloudstorage.validate_file_path(dst)
    headers = {
      'x-goog-copy-source': src,
      'x-goog-metadata-directive': 'COPY',
    }
    if src_etag is not None:
      headers['x-goog-copy-source-if-match'] = src_etag
    api = storage_api._get_storage_api(retry_params=self._retry_params)
    status, resp_headers, content = api.put_object(
        api_utils._quote_filename(dst), headers=headers)
    errors.check_status(status, [200], src, headers, resp_headers, body=content)

  def _gs_delete(self, gs_path):
    """Wrapper around cloudstorage.delete that catches NotFoundError."""
    try:
      cloudstorage.delete(filename=gs_path, retry_params=self._retry_params)
    except cloudstorage.NotFoundError:  # pragma: no cover
      pass

  def _cleanup_temp(self, upload_session):
    """Removes temporary drop file, best effort.

    Temp bucket's life cycle management will ensure old files are deleted even
    if this call fails.
    """
    # TODO(vadimsh): Finalize pending upload using upload_session.upload_url.
    self._gs_delete(upload_session.temp_gs_location)

  def _sign_gs_fetch(self, gs_path):
    """Returns signature (and its properties) for Google Storage GET operation.

    Callers can use it to construct signed Google Storage URL. The returned
    signature is guaranteed to have at least 30 min lifetime (but may have
    more).

    Returns tuple (signature, expires), where:
      * signature is byte blob with the signature
      * expires is Unix timestamp (integer seconds) when the signature expires
    """
    # TODO(vadimsh): Popular packages are fetched very frequently when they are
    # updated. Cache the signature in memcache to avoid hammering 'sign_blob'.

    # See https://cloud.google.com/storage/docs/access-control/signed-urls.
    #
    # Basically, we should sign a specially crafted multi-line string that
    # encodes expected parameters of the request. During the actual request,
    # Google Storage backend will construct the same string and verify that
    # provided signature matches it.
    expires = int(utils.time_time() + FETCH_URL_EXPIRATION_SEC)
    to_sign = '\n'.join([
      'GET',
      '',  # expected value of 'Content-MD5' header, not used
      '',  # expected value of 'Content-Type' header, not used
      str(expires),
      gs_path,
    ])
    _, signature = self._app_identity_sign_blob(to_sign)
    return signature, expires


class UploadSession(ndb.Model):
  """Some pending upload operation.

  Entity id is autogenerated by the datastore. No parent entity.
  """
  # Upload session never existed or already expired.
  STATUS_MISSING = 0
  # Client is still uploading the file.
  STATUS_UPLOADING = 1
  # Server is verifying the hash of the uploaded file.
  STATUS_VERIFYING = 2
  # The file is in the store and visible by all clients. Final state.
  STATUS_PUBLISHED = 3
  # Some other unexpected fatal error happened.
  STATUS_ERROR = 4

  # Hash algorithm to use to verify the content.
  hash_algo = ndb.StringProperty(required=True, indexed=False)
  # Expected hex digest of the file.
  hash_digest = ndb.StringProperty(required=True, indexed=False)

  # Full path in the GS to the temporary drop file that the client upload to.
  temp_gs_location = ndb.TextProperty(required=True)
  # Full path in the GS where to store the verified file.
  final_gs_location = ndb.TextProperty(required=True)

  # URL to put file content too.
  upload_url = ndb.TextProperty(required=True)

  # Status of the upload operation. See STATUS_* constants.
  status = ndb.IntegerProperty(required=True, choices=[
    STATUS_ERROR,
    STATUS_MISSING,
    STATUS_PUBLISHED,
    STATUS_UPLOADING,
    STATUS_VERIFYING,
  ])
  # For STATUS_ERROR may contain an error message.
  error_message = ndb.TextProperty(required=False)

  # Who started the upload.
  created_by = auth.IdentityProperty(required=True)
  # When the entity was created.
  created_ts = ndb.DateTimeProperty(required=True, auto_now_add=True)


class DirectUpload(object):
  """A wrapper around temp GCS file, used to upload data directly to CAS.

  Returned by CASService.start_direct_upload. Must not be instantiated directly.
  """

  def __init__(self, file_obj, hasher, callback):
    self._file_obj = file_obj
    self._hasher = hasher
    self._callback = callback
    self._len = 0

  def __enter__(self):
    return self

  def __exit__(self, exc_type, exc_value, traceback):
    self.close(commit=exc_value is None)

  @property
  def hash_digest(self):
    """Hash hex digest of the data written so far."""
    return self._hasher.hexdigest()

  @property
  def length(self):
    """Total length of the data written so far."""
    return self._len

  def write(self, data):
    """Writes a chunk of data to the temp file."""
    assert self._file_obj is not None
    self._file_obj.write(data)
    self._hasher.update(data)
    self._len += len(data)

  def close(self, commit=True):
    """Closes the temp file.

    Args:
      commit: if True, moves the data to an appropriate location in CAS,
          otherwise just deletes the temp file.
    """
    if self._file_obj is None:
      return
    assert self._callback
    self._file_obj.close()
    self._file_obj = None
    self._callback(self.hash_digest, commit)
    self._callback = None


################################################################################
## Task queues and cron jobs.


class VerifyTaskQueueHandler(webapp2.RequestHandler):  # pragma: no cover
  """Verifies the hash of the pending upload, finalizes it."""
  # pylint: disable=R0201
  @decorators.silence(
      cloudstorage.TransientError,
      datastore_errors.InternalError,
      datastore_errors.Timeout,
      datastore_errors.TransactionFailedError,
      runtime.DeadlineExceededError)
  @decorators.require_taskqueue('cas-verify')
  def post(self, unsigned_upload_id):
    service = get_cas_service()
    if service is None:
      self.abort(500, detail='CAS service is misconfigured')
    result = service.verify_pending_upload(int(unsigned_upload_id))
    if result == TASK_RETRY:
      self.abort(500, detail='Retry')
    elif result != TASK_DONE:
      logging.error('Unexpected return value from task function')


def get_backend_routes():  # pragma: no cover
  """Returns a list of webapp2.Route to add to backend WSGI app.

  Task queues, cron jobs, etc.
  """
  return [
    webapp2.Route(
        r'/internal/taskqueue/cas-verify/<unsigned_upload_id:[0-9]+>',
        VerifyTaskQueueHandler),
  ]
