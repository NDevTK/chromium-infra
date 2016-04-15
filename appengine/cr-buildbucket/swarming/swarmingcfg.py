# Copyright 2014 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

import re


DIMENSION_KEY_RGX = re.compile(r'^[a-zA-Z\_\-]+$')


def validate_tag(tag, ctx):
  # a valid swarming tag is a string that contains ":"
  if ':' not in tag:
    ctx.error('does not have ":": %s', tag)


def validate_dimension(dimension, ctx):
  components = dimension.split(':', 1)
  if len(components) != 2:
    ctx.error('does not have ":"')
  else:
    key, value = components
    if not key:
      ctx.error('no key')
    elif not DIMENSION_KEY_RGX.match(key):
      ctx.error(
        'key "%s" does not match pattern "%s"',
        key, DIMENSION_KEY_RGX.pattern)
    if not value:
      ctx.error('no value')


def validate_recipe_cfg(recipe, ctx):
  if not recipe.name:
    ctx.error('name unspecified')
  if not recipe.repository:
    ctx.error('repository unspecified')
  for i, p in enumerate(recipe.properties):
    with ctx.prefix('property #%d: ', i + 1):
      if ':' not in p:
        ctx.error('does not have colon')
      else:
        key, _ = p.split(':', 1)
        if not key:
          ctx.error('key not specified')


def validate_builder_cfg(builder, ctx):
  if not builder.name:
    ctx.error('name unspecified')

  for i, t in enumerate(builder.swarming_tags):
    with ctx.prefix('tag #%d: ', i + 1):
      validate_tag(t, ctx)

  for i, d in enumerate(builder.dimensions):
    with ctx.prefix('dimension #%d: ', i + 1):
      validate_dimension(d, ctx)

  if not builder.HasField('recipe'):
    ctx.error('recipe unspecified')
  else:
    with ctx.prefix('recipe: '):
      validate_recipe_cfg(builder.recipe, ctx)

  if builder.priority < 0 or builder.priority > 200:
    ctx.error('priority must be in [0, 200] range; got %d', builder.priority)


def validate_cfg(swarming, ctx):
  if not swarming.hostname:
    ctx.error('hostname unspecified')
  for i, t in enumerate(swarming.common_swarming_tags):
    with ctx.prefix('common tag #%d: ', i + 1):
      validate_tag(t, ctx)
  for i, d in enumerate(swarming.common_dimensions):
    with ctx.prefix('common dimension #%d: ', i + 1):
      validate_dimension(d, ctx)
  for i, b in enumerate(swarming.builders):
    with ctx.prefix('builder %s: ' % (b.name or '#%s' % (i + 1))):
      validate_builder_cfg(b, ctx)