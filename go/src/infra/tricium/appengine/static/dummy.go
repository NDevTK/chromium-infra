// Copyright 2017 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// A Go GAE module requires some .go files to be present, even if it is
// a pure static module, and the gae.py tool does not support different
// runtimes in the same deployment.

package static
