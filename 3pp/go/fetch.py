#!/usr/bin/env python
# Copyright 2018 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

import argparse
import json
import os
import sys
import urllib

from pkg_resources import parse_version


def do_latest():
  versions = []
  for release in json.load(urllib.urlopen('https://golang.org/dl/?mode=json')):
    versions.append(parse_version(release['version'].replace('go', '')))
  versions.sort()
  print versions[-1]


# TODO(akashmukherjee): Remove
def do_checkout(version, platform, kind, checkout_path):
  if kind == 'prebuilt':
    platform = platform.replace('mac', 'darwin')
    ext = 'zip' if platform.startswith('windows') else 'tar.gz'
    download_url = (
      'https://storage.googleapis.com/golang/go%(version)s.%(platform)s.%(ext)s'
      % {
        'version': version,
        'platform': platform,
        'ext': ext
      })
  else:
    ext = 'tar.gz'
    download_url = (
      'https://storage.googleapis.com/golang/go%s.src.tar.gz' % (version,))

  print >>sys.stderr, 'fetching', download_url
  urllib.urlretrieve(download_url,
                     os.path.join(checkout_path, 'archive.' + ext))


def get_download_url(version, platform, kind):
  if kind == 'prebuilt':
    platform = platform.replace('mac', 'darwin')
    ext = '.zip' if platform.startswith('windows') else '.tar.gz'
    download_url = (
      'https://storage.googleapis.com/golang/go%(version)s.%(platform)s%(ext)s'
      % {
        'version': version,
        'platform': platform,
        'ext': ext
      })
  else:
    ext = '.tar.gz'
    download_url = (
      'https://storage.googleapis.com/golang/go%s.src.tar.gz' % (version,))

  partial_manifest = {
    'url': download_url,
    'ext': ext,
  }
  print(json.dumps(partial_manifest))


def main():
  ap = argparse.ArgumentParser()
  ap.add_argument('kind', choices=('prebuilt', 'source'))

  sub = ap.add_subparsers()

  latest = sub.add_parser('latest')
  latest.set_defaults(func=lambda _opts: do_latest())

  # TODO(akashmukherjee): Remove
  checkout = sub.add_parser('checkout')
  checkout.add_argument('checkout_path')
  checkout.set_defaults(
    func=lambda opts: do_checkout(
      os.environ['_3PP_VERSION'], os.environ['_3PP_PLATFORM'],
      opts.kind, opts.checkout_path))

  download = sub.add_parser('get_url')
  download.set_defaults(
    func=lambda opts: get_download_url(
      os.environ['_3PP_VERSION'], os.environ['_3PP_PLATFORM'],
      opts.kind))

  opts = ap.parse_args()
  return opts.func(opts)

if __name__ == '__main__':
  sys.exit(main())
