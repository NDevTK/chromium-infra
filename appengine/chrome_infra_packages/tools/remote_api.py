#!/usr/bin/env python
# Copyright 2014 The Swarming Authors. All rights reserved.
# Use of this source code is governed by the Apache v2.0 license that can be
# found in the LICENSE file.

"""Accesses an GAE instance via remote_api."""

import os
import sys

APP_DIR = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
COMPONENTS_DIR = os.path.abspath(os.path.join(
    APP_DIR, '..', 'swarming', 'appengine', 'components'))

sys.path.insert(0, COMPONENTS_DIR)
sys.path.insert(0, os.path.join(COMPONENTS_DIR, 'third_party'))

from tools import remote_api  # pylint: disable=E0611


def setup_context():
  """Symbols to import into interactive console."""
  sys.path.insert(0, APP_DIR)

  # Unused variable 'XXX'; they are accessed via locals().
  # pylint: disable=W0612
  # TODO(vadimsh): Import important modules.
  return locals().copy()


if __name__ == '__main__':
  sys.exit(remote_api.main(sys.argv[1:], APP_DIR, setup_context))
