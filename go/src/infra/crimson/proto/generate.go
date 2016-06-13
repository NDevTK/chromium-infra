// Copyright 2016 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package crimson

//go:generate go install github.com/luci/luci-go/tools/cmd/cproto
//go:generate cproto
//go:generate go install github.com/luci/luci-go/tools/cmd/svcdec
//go:generate svcdec -type CrimsonServer
