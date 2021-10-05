# Copyright 2021 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

import re

# Names of active well-known experiments
BBAGENT_GET_BUILD = 'luci.buildbucket.bbagent_getbuild'
CANARY = 'luci.buildbucket.canary_software'
NON_PROD = 'luci.non_production'
RECIPE_PY3 = 'luci.recipes.use_python3'
USE_BBAGENT = 'luci.buildbucket.use_bbagent'

# Decommissioned experiments which still may linger in configs but have
# no effect. Can be removed once configs no longer reference them and all
# Buildbucket GAE versions that use them are deleted.
USE_REALMS = 'luci.use_realms'

WELL_KNOWN = frozenset([
    BBAGENT_GET_BUILD,
    CANARY,
    NON_PROD,
    RECIPE_PY3,
    USE_BBAGENT,
    USE_REALMS,
])

_VALID_NAME = re.compile(r'^[a-z][a-z0-9_]*(?:\.[a-z][a-z0-9_]*)*$')


def check_invalid_name(exp_name):
  """Returns an error message string if this is an invalid expirement.

  Returns None if `exp_name` is a valid experiment name."""
  if not _VALID_NAME.match(exp_name):
    return 'does not match %r' % _VALID_NAME.pattern

  if exp_name.startswith('luci.') and exp_name not in WELL_KNOWN:
    return 'unknown experiment has reserved prefix "luci."'

  return
