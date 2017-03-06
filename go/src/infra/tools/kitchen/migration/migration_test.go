// Copyright 2016 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package migration

import (
	"testing"

	. "github.com/luci/luci-go/common/testing/assertions"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMigration(t *testing.T) {
	t.Parallel()

	Convey("transformProperties", t, func() {
		Convey("works", func() {
			props := map[string]interface{}{
				"mastername":  "luci.chromium.try",
				"buildername": "LUCI linux_chromium_rel_ng",
				"foo":         "bar",
			}
			So(TransformProperties(props), ShouldBeNil)
			So(props["mastername"], ShouldEqual, "tryserver.chromium.linux")
			So(props["buildername"], ShouldEqual, "linux_chromium_rel_ng")
			So(props["foo"], ShouldEqual, "bar")
		})
		Convey("noop if no master", func() {
			props := map[string]interface{}{
				"buildername": "linux_chromium_rel_ng",
			}
			So(TransformProperties(props), ShouldBeNil)
			So(props["buildername"], ShouldEqual, "linux_chromium_rel_ng")
		})
		Convey("noop if builder does not start with LUCI", func() {
			props := map[string]interface{}{
				"mastername":  "luci.chromium.try",
				"buildername": "linux_chromium_rel_ng",
			}
			So(TransformProperties(props), ShouldBeNil)
			So(props["mastername"], ShouldEqual, "luci.chromium.try")
			So(props["buildername"], ShouldEqual, "linux_chromium_rel_ng")
		})
		Convey("noop if builder isn't known", func() {
			props := map[string]interface{}{
				"mastername":  "luci.chromium.try",
				"buildername": "LUCI my builder",
			}
			So(TransformProperties(props), ShouldBeNil)
			So(props["mastername"], ShouldEqual, "luci.chromium.try")
			So(props["buildername"], ShouldEqual, "LUCI my builder")
		})
		Convey("noop if master isn't known", func() {
			props := map[string]interface{}{
				"mastername":  "luci.fuchsia.try",
				"buildername": "linux_chromium_rel_ng",
			}
			So(TransformProperties(props), ShouldBeNil)
			So(props["mastername"], ShouldEqual, "luci.fuchsia.try")
			So(props["buildername"], ShouldEqual, "linux_chromium_rel_ng")
		})
		Convey("fails if buildername is not present", func() {
			props := map[string]interface{}{
				"mastername": "luci.chromium.try",
				"foo":        "bar",
			}
			So(TransformProperties(props), ShouldErrLike, "buildername property is not set")
		})
	})
}
