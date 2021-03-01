// Copyright 2021 The Chromium OS Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package controller

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	. "go.chromium.org/luci/common/testing/assertions"
	"google.golang.org/genproto/protobuf/field_mask"

	ufspb "infra/unifiedfleet/api/v1/models"
	"infra/unifiedfleet/app/model/inventory"
)

func mockMachineLSEDeployment(serialNumber string) *ufspb.MachineLSEDeployment {
	return &ufspb.MachineLSEDeployment{
		SerialNumber: serialNumber,
	}
}

func TestUpdateMachineLSEDeployment(t *testing.T) {
	t.Parallel()
	ctx := testingContext()
	Convey("UpdateMachineLSEDeployment", t, func() {
		Convey("Update MachineLSEDeployment for non-existing MachineLSEDeployment - happy path", func() {
			dr1 := mockMachineLSEDeployment("serial-1")
			resp, err := UpdateMachineLSEDeployment(ctx, dr1, nil)
			So(err, ShouldBeNil)
			So(resp.GetHostname(), ShouldEqual, "no-host-yet-serial-1")
			So(resp.GetSerialNumber(), ShouldEqual, "serial-1")

			resGet, err := inventory.GetMachineLSEDeployment(ctx, "serial-1")
			So(err, ShouldBeNil)
			So(resGet, ShouldResembleProto, resp)
		})

		Convey("Update MachineLSEDeployment for existing MachineLSEDeployment - happy path", func() {
			dr2 := mockMachineLSEDeployment("serial-2")
			_, err := UpdateMachineLSEDeployment(ctx, dr2, nil)
			So(err, ShouldBeNil)

			dr2.Hostname = "hostname-2"
			resp, err := UpdateMachineLSEDeployment(ctx, dr2, nil)
			So(err, ShouldBeNil)
			So(resp.GetHostname(), ShouldEqual, "hostname-2")
			So(resp.GetSerialNumber(), ShouldEqual, "serial-2")

			resGet, err := inventory.GetMachineLSEDeployment(ctx, "serial-2")
			So(err, ShouldBeNil)
			So(resGet, ShouldResembleProto, dr2)
		})

		Convey("Update MachineLSEDeployment for existing MachineLSEDeployment - partial update hostname", func() {
			dr3 := mockMachineLSEDeployment("serial-3")
			_, err := UpdateMachineLSEDeployment(ctx, dr3, nil)
			So(err, ShouldBeNil)

			dr3.Hostname = "hostname-3"
			dr3.DeploymentIdentifier = "identifier-3"
			resp, err := UpdateMachineLSEDeployment(ctx, dr3, &field_mask.FieldMask{Paths: []string{"hostname"}})
			So(err, ShouldBeNil)
			So(resp.GetSerialNumber(), ShouldEqual, "serial-3")
			So(resp.GetHostname(), ShouldEqual, "hostname-3")
			So(resp.GetDeploymentIdentifier(), ShouldBeEmpty)

			resGet, err := inventory.GetMachineLSEDeployment(ctx, "serial-3")
			So(err, ShouldBeNil)
			dr3.DeploymentIdentifier = ""
			So(resGet.GetHostname(), ShouldEqual, "hostname-3")
			So(resGet.GetDeploymentIdentifier(), ShouldBeEmpty)
		})
	})
}