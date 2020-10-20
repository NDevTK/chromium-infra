// Copyright 2020 The Chromium OS Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package controller

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	. "go.chromium.org/luci/common/testing/assertions"
	"google.golang.org/genproto/protobuf/field_mask"

	ufspb "infra/unifiedfleet/api/v1/proto"
	"infra/unifiedfleet/app/model/history"
	"infra/unifiedfleet/app/model/registration"
	"infra/unifiedfleet/app/util"
)

func mockRack(name, row string, zone ufspb.Zone) *ufspb.Rack {
	return &ufspb.Rack{
		Name: name,
		Location: &ufspb.Location{
			Zone:        zone,
			Row:         row,
			BarcodeName: name,
		},
	}
}

func mockAsset(name, model, row, rack, pos, host string, aType ufspb.AssetType, zone ufspb.Zone) *ufspb.Asset {
	return &ufspb.Asset{
		Name:  name,
		Type:  aType,
		Model: model,
		Location: &ufspb.Location{
			Zone:        zone,
			Row:         row,
			Rack:        rack,
			Position:    pos,
			BarcodeName: host,
		},
	}
}

func mockAssetInfo(serialNumber, costCenter, googleCodeName, model, buildTarget, referenceBoard, ethernetMacAddress, sku, phase string) *ufspb.AssetInfo {
	return &ufspb.AssetInfo{
		SerialNumber:       serialNumber,
		CostCenter:         costCenter,
		GoogleCodeName:     googleCodeName,
		Model:              model,
		BuildTarget:        buildTarget,
		ReferenceBoard:     referenceBoard,
		EthernetMacAddress: ethernetMacAddress,
		Sku:                sku,
		Phase:              phase,
	}
}

func TestAssetRegistration(t *testing.T) {
	t.Parallel()
	ctx := testingContext()
	Convey("Testing AssetRegistration", t, func() {
		Convey("Register asset with existing rack", func() {
			r := mockRack("chromeos6-row2-rack3", "2", ufspb.Zone_ZONE_CHROMEOS6)
			_, err := registration.CreateRack(ctx, r)
			So(err, ShouldBeNil)
			a := mockAsset("C001001", "eve", "2", "chromeos6-row2-rack3", "1", "chromeos6-row2-rack3-host1", ufspb.AssetType_DUT, ufspb.Zone_ZONE_CHROMEOS6)
			_, err = AssetRegistration(ctx, a)
			So(err, ShouldBeNil)
			changes, err := history.QueryChangesByPropertyName(ctx, "name", "assets/C001001")
			So(err, ShouldBeNil)
			So(changes, ShouldHaveLength, 1)
			So(changes[0].GetEventLabel(), ShouldEqual, "asset")
			So(changes[0].GetOldValue(), ShouldEqual, LifeCycleRegistration)
			So(changes[0].GetNewValue(), ShouldEqual, LifeCycleRegistration)
		})
		Convey("Register asset with non-existent rack", func() {
			a := mockAsset("C001002", "eve", "2", "chromeos6-row3-rack3", "1", "chromeos6-row2-rack3-host1", ufspb.AssetType_DUT, ufspb.Zone_ZONE_CHROMEOS6)
			_, err := AssetRegistration(ctx, a)
			So(err, ShouldNotBeNil)
			changes, err := history.QueryChangesByPropertyName(ctx, "name", "assets/C001002")
			So(err, ShouldBeNil)
			So(changes, ShouldHaveLength, 0)
		})
		Convey("Register asset with invalid name", func() {
			r := mockRack("chromeos6-row4-rack3", "2", ufspb.Zone_ZONE_CHROMEOS6)
			_, err := registration.CreateRack(ctx, r)
			So(err, ShouldBeNil)
			a := mockAsset("", "eve", "4", "chromeos6-row4-rack3", "1", "chromeos6-row4-rack3-host1", ufspb.AssetType_DUT, ufspb.Zone_ZONE_CHROMEOS6)
			_, err = AssetRegistration(ctx, a)
			So(err, ShouldNotBeNil)
		})
		Convey("Register existing asset", func() {
			r := mockRack("chromeos6-row2-rack4", "2", ufspb.Zone_ZONE_CHROMEOS6)
			_, err := registration.CreateRack(ctx, r)
			So(err, ShouldBeNil)
			a := mockAsset("C001004", "eve", "2", "chromeos6-row2-rack4", "1", "chromeos6-row2-rack4-host1", ufspb.AssetType_DUT, ufspb.Zone_ZONE_CHROMEOS6)
			_, err = AssetRegistration(ctx, a)
			So(err, ShouldBeNil)
			changes, err := history.QueryChangesByPropertyName(ctx, "name", "assets/C001001")
			So(err, ShouldBeNil)
			So(changes, ShouldHaveLength, 1)
			So(changes[0].GetEventLabel(), ShouldEqual, "asset")
			So(changes[0].GetOldValue(), ShouldEqual, LifeCycleRegistration)
			So(changes[0].GetNewValue(), ShouldEqual, LifeCycleRegistration)
			_, err = AssetRegistration(ctx, a)
			So(err, ShouldNotBeNil)
		})
	})
}

func TestUpdateAsset(t *testing.T) {
	t.Parallel()
	ctx := testingContext()
	Convey("Testing UpdateAsset", t, func() {
		Convey("Update non existent asset", func() {
			r := mockRack("chromeos6-row2-rack3", "2", ufspb.Zone_ZONE_CHROMEOS6)
			_, err := registration.CreateRack(ctx, r)
			So(err, ShouldBeNil)
			b := mockAsset("C001001", "eve", "2", "chromeos6-row2-rack3", "1", "chromeos6-row2-rack3-host1", ufspb.AssetType_DUT, ufspb.Zone_ZONE_CHROMEOS6)
			_, err = UpdateAsset(ctx, b, &field_mask.FieldMask{Paths: []string{"type", "model"}})
			fmt.Println(err)
			So(err.Error(), ShouldContainSubstring, "unable to update asset C001001")
			changes, err := history.QueryChangesByPropertyName(ctx, "name", "assets/C001001")
			So(err, ShouldBeNil)
			So(changes, ShouldHaveLength, 0)
		})

		Convey("Move asset to non-existent rack", func() {
			a := mockAsset("C001001", "eve", "2", "chromeos6-row2-rack3", "1", "chromeos6-row2-rack3-host1", ufspb.AssetType_DUT, ufspb.Zone_ZONE_CHROMEOS6)
			_, err := AssetRegistration(ctx, a)
			So(err, ShouldBeNil)
			b := mockAsset("C001001", "eve", "2", "chromeos6-row2-rack4", "1", "chromeos6-row2-rack3-host1", ufspb.AssetType_DUT, ufspb.Zone_ZONE_CHROMEOS6)
			_, err = UpdateAsset(ctx, b, &field_mask.FieldMask{Paths: []string{"location.rack"}})
			So(err.Error(), ShouldContainSubstring, "There is no Rack with RackID chromeos6-row2-rack4")
			changes, err := history.QueryChangesByPropertyName(ctx, "name", "assets/C001001")
			So(err, ShouldBeNil)
			So(changes, ShouldHaveLength, 1)
		})

		Convey("Move Asset to existing rack", func() {
			r := mockRack("chromeos6-row2-rack5", "2", ufspb.Zone_ZONE_CHROMEOS6)
			_, err := registration.CreateRack(ctx, r)
			So(err, ShouldBeNil)
			r = mockRack("chromeos6-row2-rack6", "2", ufspb.Zone_ZONE_CHROMEOS6)
			_, err = registration.CreateRack(ctx, r)
			So(err, ShouldBeNil)
			a := mockAsset("C001002", "eve", "2", "chromeos6-row2-rack5", "1", "chromeos6-row2-rack3-host1", ufspb.AssetType_DUT, ufspb.Zone_ZONE_CHROMEOS6)
			_, err = AssetRegistration(ctx, a)
			So(err, ShouldBeNil)
			b := mockAsset("C001002", "eve", "2", "chromeos6-row2-rack6", "2", "chromeos6-row2-rack4-host2", ufspb.AssetType_DUT, ufspb.Zone_ZONE_CHROMEOS6)
			_, err = UpdateAsset(ctx, b, &field_mask.FieldMask{Paths: []string{"location.rack"}})
			So(err, ShouldBeNil)
			changes, err := history.QueryChangesByPropertyName(ctx, "name", "assets/C001002")
			So(err, ShouldBeNil)
			So(changes, ShouldHaveLength, 2)
			msgs, err := history.QuerySnapshotMsgByPropertyName(ctx, "name", "assets/C001002")
			So(err, ShouldBeNil)
			So(msgs, ShouldHaveLength, 0)
		})

		Convey("Update Asset info of an asset", func() {
			r := mockRack("chromeos6-row2-rack7", "2", ufspb.Zone_ZONE_CHROMEOS6)
			_, err := registration.CreateRack(ctx, r)
			So(err, ShouldBeNil)
			a := mockAsset("C001003", "eve", "2", "chromeos6-row2-rack5", "3", "chromeos6-row2-rack3-host1", ufspb.AssetType_DUT, ufspb.Zone_ZONE_CHROMEOS6)
			_, err = AssetRegistration(ctx, a)
			So(err, ShouldBeNil)
			ai := mockAssetInfo("MTV100212", "", "", "", "", "", "", "", "")
			a.Info = ai
			_, err = UpdateAsset(ctx, a, &field_mask.FieldMask{Paths: []string{"info.serial_number"}})
			So(err, ShouldBeNil)
			changes, err := history.QueryChangesByPropertyName(ctx, "name", "assets/C001003")
			So(err, ShouldBeNil)
			So(changes, ShouldHaveLength, 2)
			So(changes[0].GetEventLabel(), ShouldEqual, "asset")
			So(changes[0].GetOldValue(), ShouldEqual, LifeCycleRegistration)
			So(changes[0].GetNewValue(), ShouldEqual, LifeCycleRegistration)
			So(changes[1].GetEventLabel(), ShouldEqual, "asset.info.serial_number")
			So(changes[1].GetNewValue(), ShouldEqual, "MTV100212")
			So(changes[1].GetOldValue(), ShouldEqual, "")
		})

		Convey("Update Asset with invalid mask", func() {
			a := mockAsset("C001004", "eve", "2", "chromeos6-row2-rack3", "1", "chromeos6-row2-rack3-host1", ufspb.AssetType_DUT, ufspb.Zone_ZONE_CHROMEOS6)
			_, err := AssetRegistration(ctx, a)
			So(err, ShouldBeNil)
			b := mockAsset("C001004", "eve", "2", "chromeos6-row2-rack3", "1", "chromeos6-row2-rack3-host1", ufspb.AssetType_DUT, ufspb.Zone_ZONE_CHROMEOS6)
			// Attempt to update name of the asset
			_, err = UpdateAsset(ctx, b, &field_mask.FieldMask{Paths: []string{"name"}})
			So(err, ShouldNotBeNil)
			changes, err := history.QueryChangesByPropertyName(ctx, "name", "assets/C001004")
			So(err, ShouldBeNil)
			So(changes, ShouldHaveLength, 1)
			// Attempt to update name of the asset in asset info
			_, err = UpdateAsset(ctx, b, &field_mask.FieldMask{Paths: []string{"info.asset_tag"}})
			So(err, ShouldNotBeNil)
			changes, err = history.QueryChangesByPropertyName(ctx, "name", "assets/C001004")
			So(err, ShouldBeNil)
			So(changes, ShouldHaveLength, 1)
			// Attempt to update timestamp of the asset
			_, err = UpdateAsset(ctx, b, &field_mask.FieldMask{Paths: []string{"update_time"}})
			So(err, ShouldNotBeNil)
			changes, err = history.QueryChangesByPropertyName(ctx, "name", "assets/C001004")
			So(err, ShouldBeNil)
			So(changes, ShouldHaveLength, 1)
			// Attempt to clear zone of the asset
			b.Location.Zone = ufspb.Zone_ZONE_UNSPECIFIED
			_, err = UpdateAsset(ctx, b, &field_mask.FieldMask{Paths: []string{"location.zone"}})
			So(err, ShouldNotBeNil)
			changes, err = history.QueryChangesByPropertyName(ctx, "name", "assets/C001004")
			So(err, ShouldBeNil)
			So(changes, ShouldHaveLength, 1)
			b.Location.Zone = ufspb.Zone_ZONE_CHROMEOS6
			// Attempt to clear rack of the asset
			b.Location.Rack = ""
			_, err = UpdateAsset(ctx, b, &field_mask.FieldMask{Paths: []string{"location.rack"}})
			So(err, ShouldNotBeNil)
			changes, err = history.QueryChangesByPropertyName(ctx, "name", "assets/C001004")
			So(err, ShouldBeNil)
			So(changes, ShouldHaveLength, 1)
		})
	})
}

func TestGetAsset(t *testing.T) {
	t.Parallel()
	ctx := testingContext()
	Convey("Testing GetAsset", t, func() {
		Convey("Get existing assets", func() {
			r := mockRack("chromeos6-row2-rack3", "2", ufspb.Zone_ZONE_CHROMEOS6)
			_, err := RackRegistration(ctx, r)
			So(err, ShouldBeNil)
			a := mockAsset("C001001", "eve", "2", "chromeos6-row2-rack3", "1", "chromeos6-row2-rack3-host1", ufspb.AssetType_DUT, ufspb.Zone_ZONE_CHROMEOS6)
			_, err = AssetRegistration(ctx, a)
			So(err, ShouldBeNil)
			respA, err := GetAsset(ctx, "C001001")
			So(err, ShouldBeNil)
			So(respA, ShouldResembleProto, a)
		})
		Convey("Get non existing assets", func() {
			respA, err := GetAsset(ctx, "C001004")
			So(err, ShouldNotBeNil)
			So(respA, ShouldBeNil)
		})
		Convey("Get invalid assets", func() {
			respB, err := GetAsset(ctx, "")
			So(err, ShouldNotBeNil)
			So(respB, ShouldBeNil)
		})
	})
}

func createArrayOfMockAssets(n int, prefix, zone, assetType, model string) []*ufspb.Asset {
	var assets []*ufspb.Asset
	for i := 0; i < n; i++ {
		aType := ufspb.AssetType_UNDEFINED
		if assetType == "dut" {
			aType = ufspb.AssetType_DUT
		} else if assetType == "labstation" {
			aType = ufspb.AssetType_LABSTATION
		}
		a := mockAsset(fmt.Sprintf("%s00%d", prefix, i), model, "3", fmt.Sprintf("%s-row3-rack3", zone), fmt.Sprintf("%d", i), fmt.Sprintf("%s-row3-rack3-host%d", zone, i), aType, util.ToUFSZone(zone))
		assets = append(assets, a)
	}
	return assets
}

func TestListAssets(t *testing.T) {
	t.Parallel()
	ctx := testingContext()
	r := mockRack("chromeos6-row3-rack3", "3", ufspb.Zone_ZONE_CHROMEOS6)
	RackRegistration(ctx, r)
	r = mockRack("chromeos2-row3-rack3", "3", ufspb.Zone_ZONE_CHROMEOS2)
	RackRegistration(ctx, r)
	dutChromeos6 := createArrayOfMockAssets(4, "EVE6", "chromeos6", "dut", "eve")
	labstationsChromeos6 := createArrayOfMockAssets(4, "FIZ6", "chromeos6", "labstation", "fizz")
	guadoChromeos2 := createArrayOfMockAssets(4, "GUA2", "chromeos2", "labstation", "guado")
	fizzChromeos2 := createArrayOfMockAssets(4, "FIZ2", "chromeos2", "labstation", "fizz")
	assets := append(dutChromeos6, labstationsChromeos6...)
	assets = append(assets, guadoChromeos2...)
	assets = append(assets, fizzChromeos2...)
	chromeos2Assets := append(fizzChromeos2, guadoChromeos2...)
	chromeos6Assets := append(dutChromeos6, labstationsChromeos6...)
	labstationAssets := append(fizzChromeos2, labstationsChromeos6...)
	labstationAssets = append(labstationAssets, guadoChromeos2...)
	for _, asset := range assets {
		AssetRegistration(ctx, asset)
	}
	Convey("Testing ListAssets", t, func() {
		Convey("List all existing assets", func() {
			respAssets, _, err := ListAssets(ctx, 16, "", "", false)
			So(err, ShouldBeNil)
			So(respAssets, ShouldHaveLength, 16)
		})
		Convey("List assets by zone", func() {
			respAssets, _, err := ListAssets(ctx, 10, "", "zone=chromeos2", false)
			So(err, ShouldBeNil)
			So(respAssets, ShouldHaveLength, 8)
			So(respAssets, ShouldResembleProto, chromeos2Assets)
			respAssets, _, err = ListAssets(ctx, 10, "", "zone=chromeos6", false)
			So(err, ShouldBeNil)
			So(respAssets, ShouldHaveLength, 8)
			So(respAssets, ShouldResembleProto, chromeos6Assets)
		})
		Convey("List assets by model", func() {
			respAssets, _, err := ListAssets(ctx, 10, "", "model=guado", false)
			So(err, ShouldBeNil)
			So(respAssets, ShouldHaveLength, 4)
			So(respAssets, ShouldResembleProto, guadoChromeos2)
			respAssets, _, err = ListAssets(ctx, 10, "", "model=eve", false)
			So(err, ShouldBeNil)
			So(respAssets, ShouldHaveLength, 4)
			So(respAssets, ShouldResembleProto, dutChromeos6)
		})
		Convey("List assets by type", func() {
			respAssets, _, err := ListAssets(ctx, 10, "", "type=dut", false)
			So(err, ShouldBeNil)
			So(respAssets, ShouldHaveLength, 4)
			So(respAssets, ShouldResembleProto, dutChromeos6)
			respAssets, _, err = ListAssets(ctx, 12, "", "type=labstation", false)
			So(err, ShouldBeNil)
			So(respAssets, ShouldHaveLength, 12)
			So(respAssets, ShouldResembleProto, labstationAssets)
		})
		Convey("List assets by combination of filters", func() {
			respAssets, _, err := ListAssets(ctx, 10, "", "type=dut&model=eve", false)
			So(err, ShouldBeNil)
			So(respAssets, ShouldHaveLength, 4)
			So(respAssets, ShouldResembleProto, dutChromeos6)
			respAssets, _, err = ListAssets(ctx, 10, "", "type=labstation&zone=chromeos2", false)
			So(err, ShouldBeNil)
			So(respAssets, ShouldHaveLength, 8)
			So(respAssets, ShouldResembleProto, chromeos2Assets)
			respAssets, _, err = ListAssets(ctx, 10, "", "type=labstation&zone=chromeos2&model=guado", false)
			So(err, ShouldBeNil)
			So(respAssets, ShouldHaveLength, 4)
			So(respAssets, ShouldResembleProto, guadoChromeos2)
		})
	})
}

func TestDeleteAsset(t *testing.T) {
	t.Parallel()
	ctx := testingContext()
	Convey("Testing DeleteAsset", t, func() {
		Convey("Delete existing assets", func() {
			r := mockRack("chromeos6-row2-rack3", "2", ufspb.Zone_ZONE_CHROMEOS6)
			_, err := RackRegistration(ctx, r)
			So(err, ShouldBeNil)
			a := mockAsset("C001001", "eve", "2", "chromeos6-row2-rack3", "1", "chromeos6-row2-rack3-host1", ufspb.AssetType_DUT, ufspb.Zone_ZONE_CHROMEOS6)
			_, err = AssetRegistration(ctx, a)
			So(err, ShouldBeNil)
			err = DeleteAsset(ctx, "C001001")
			So(err, ShouldBeNil)
			changes, err := history.QueryChangesByPropertyName(ctx, "name", "assets/C001001")
			So(err, ShouldBeNil)
			So(changes, ShouldHaveLength, 2)
			So(changes[0].GetEventLabel(), ShouldEqual, "asset")
			So(changes[0].GetOldValue(), ShouldEqual, LifeCycleRegistration)
			So(changes[0].GetNewValue(), ShouldEqual, LifeCycleRegistration)
			So(changes[1].GetEventLabel(), ShouldEqual, "asset")
			So(changes[1].GetOldValue(), ShouldEqual, LifeCycleRetire)
			So(changes[1].GetNewValue(), ShouldEqual, LifeCycleRetire)
		})
		Convey("Delete non existing assets", func() {
			err := DeleteAsset(ctx, "C001004")
			So(err, ShouldNotBeNil)
			changes, err := history.QueryChangesByPropertyName(ctx, "name", "assets/C001004")
			So(err, ShouldBeNil)
			So(changes, ShouldHaveLength, 0)
		})
		Convey("Delete invalid assets", func() {
			err := DeleteAsset(ctx, "")
			So(err, ShouldNotBeNil)
		})
	})
}
