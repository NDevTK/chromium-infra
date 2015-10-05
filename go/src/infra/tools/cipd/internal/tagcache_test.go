// Copyright 2015 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package internal

import (
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"infra/tools/cipd/common"
)

func TestTagCacheWorks(t *testing.T) {
	Convey("Works", t, func(c C) {
		tc := TagCache{}
		So(tc.ResolveTag("pkg", "tag:1"), ShouldResemble, common.Pin{})
		So(tc.Dirty(), ShouldBeFalse)

		// Add new.
		tc.AddTag(common.Pin{
			PackageName: "pkg",
			InstanceID:  strings.Repeat("a", 40),
		}, "tag:1")
		So(tc.Dirty(), ShouldBeTrue)
		So(tc.ResolveTag("pkg", "tag:1"), ShouldResemble, common.Pin{
			PackageName: "pkg",
			InstanceID:  strings.Repeat("a", 40),
		})

		// Replace existing.
		tc.AddTag(common.Pin{
			PackageName: "pkg",
			InstanceID:  strings.Repeat("b", 40),
		}, "tag:1")
		So(tc.Dirty(), ShouldBeTrue)
		So(tc.ResolveTag("pkg", "tag:1"), ShouldResemble, common.Pin{
			PackageName: "pkg",
			InstanceID:  strings.Repeat("b", 40),
		})

		// Save\load.
		blob, err := tc.Save()
		So(err, ShouldBeNil)
		So(tc.Dirty(), ShouldBeFalse)
		another := TagCache{}
		So(another.Load(blob), ShouldBeNil)
		So(another.ResolveTag("pkg", "tag:1"), ShouldResemble, common.Pin{
			PackageName: "pkg",
			InstanceID:  strings.Repeat("b", 40),
		})
	})
}
