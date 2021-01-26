// Copyright 2020 The Chromium OS Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package hostinfo

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	grpc "google.golang.org/grpc"

	fleet "infra/appengine/crosskylabadmin/api/fleet/v1"
	"infra/libs/skylab/inventory"
)

const fullResponse = `{
	"labels": [
		"arc",
		"board:FAKE-BOARD",
		"model:FAKE-MODEL",
		"device-sku:FAKE-DEVICE-SKU",
		"sku:FAKE-SKU",
		"platform:FAKE-PLATFORM",
		"wifi_chip:FAKE-CHIP",
		"ec:cros",
		"os:cros",
		"phase:DVT",
		"variant:FAKE-VARIANT",
		"bluetooth",
		"fingerprint",
		"internal_display",
		"touchpad",
		"power:battery",
		"storage:nvme",
		"hw_video_acc_h264",
		"cr50:prepvt",
		"cr50-ro-keyid:prod",
		"cts_abi_x86",
		"cts_cpu_x86",
		"hwid_component:FAKE-BATTERY",
		"hwid_component:FAKE-DISPLAY",
		"audio_loopback_dongle",
		"servo",
		"servo_state:WORKING",
		"servo_type:servo_v4_with_ccd_cr50",
		"servo_topology:eyJtYWluIjp7InR5cGUiOiJzZXJ2b192NCIsInN5c2ZzX3Byb2R1Y3QiOiJTZXJ2byBWNCIsInNlcmlhbCI6IkZBS0UtU0VSVk8tU0VSSUFMIiwidXNiX2h1Yl9wb3J0IjoiRkFLRS1TRVJWTy1VU0ItSFVCLVBPUlQifSwiY2hpbGRyZW4iOlt7InR5cGUiOiJjY2RfY3I1MCIsInN5c2ZzX3Byb2R1Y3QiOiJDcjUwIiwic2VyaWFsIjoiRkFLRS1UT1BPTE9HWS1JVEVNIiwidXNiX2h1Yl9wb3J0IjoiRkFLRS1VU0ItSFVCLVBPUlQifV19",
		"servo_usb_state:NORMAL",
		"storage_state:NORMAL",
		"pool:faft-test"
	],
	"attributes": {
		"HWID": "FAKE-HWID",
		"powerunit_hostname": "FAKE-POWERUNIT-HOSTNAME",
		"powerunit_outlet": "AA6",
		"serial_number": "FAKE-SERIAL-NUMBER",
		"servo_host": "FAKE-SERVO-HOST",
		"servo_port": "FAKE-SERVO-PORT",
		"servo_serial": "FAKE-SERVO-SERIAL",
		"servo_setup": "REGULAR",
		"servo_type": "servo_v4_with_ccd_cr50"
	},
	"stable_versions": {
		"cros": "FAKE-CROS-VERSION",
		"faft": "FAKE-FAFT-VERSION",
		"firmware": "FAKE-FIRMWARE-VERSION",
		"servo-cros": "FAKE-SERVO-CROS-VERSION"
	},
	"serializer_version": 1
}`

type FakeGetDutInfo struct {
	hostname string
	response *inventory.DeviceUnderTest
}

func (f *FakeGetDutInfo) GetDutInfo(ctx context.Context, id string, byHostname bool) (*inventory.DeviceUnderTest, error) {
	if !byHostname {
		return nil, fmt.Errorf("by hostname not provided")
	}
	if id != f.hostname {
		return nil, fmt.Errorf("bad hostname")
	}
	return f.response, nil
}

type FakeGetStableVersion struct {
	cros      string
	faft      string
	firmware  string
	servoCros string
}

func (f *FakeGetStableVersion) GetStableVersion(ctx context.Context, in *fleet.GetStableVersionRequest, opts ...grpc.CallOption) (*fleet.GetStableVersionResponse, error) {
	resp := &fleet.GetStableVersionResponse{}
	resp.CrosVersion = f.cros
	resp.FaftVersion = f.faft
	resp.FirmwareVersion = f.firmware
	resp.ServoCrosVersion = f.servoCros
	return resp, nil
}

func TestGetter(t *testing.T) {
	bg := context.Background()

	const hostname = "FAKE-HOSTNAME"
	const expected = fullResponse
	const expectedErr = ""

	g := NewGetter(
		&FakeGetDutInfo{
			hostname,
			&inventory.DeviceUnderTest{
				Common: &inventory.CommonDeviceSpecs{
					Attributes: []*inventory.KeyValue{
						{
							Key:   s("HWID"),
							Value: s("FAKE-HWID"),
						},
						{
							Key:   s("powerunit_hostname"),
							Value: s("FAKE-POWERUNIT-HOSTNAME"),
						},
						{
							Key:   s("powerunit_outlet"),
							Value: s("AA6"),
						},
						{
							Key:   s("serial_number"),
							Value: s("FAKE-SERIAL-NUMBER"),
						},
						{
							Key:   s("servo_host"),
							Value: s("FAKE-SERVO-HOST"),
						},
						{
							Key:   s("servo_port"),
							Value: s("FAKE-SERVO-PORT"),
						},
						{
							Key:   s("servo_setup"),
							Value: s("REGULAR"),
						},
						{
							Key:   s("servo_serial"),
							Value: s("FAKE-SERVO-SERIAL"),
						},
						{
							Key:   s("servo_type"),
							Value: s("servo_v4_with_ccd_cr50"),
						},
					},
					Labels: &inventory.SchedulableLabels{
						Arc:   b(true),
						Board: s("FAKE-BOARD"),
						Brand: nil,
						Capabilities: &inventory.HardwareCapabilities{
							Atrus:           nil,
							Bluetooth:       b(true),
							Carrier:         nil,
							Fingerprint:     b(true),
							GpuFamily:       nil,
							Graphics:        nil,
							InternalDisplay: b(true),
							Power:           s("battery"),
							Storage:         s("nvme"),
							Touchpad:        b(true),
							VideoAcceleration: []inventory.HardwareCapabilities_VideoAcceleration{
								inventory.HardwareCapabilities_VIDEO_ACCELERATION_H264,
							},
						},
						Cr50Phase:   cr50(inventory.SchedulableLabels_CR50_PHASE_PREPVT),
						Cr50RoKeyid: s("prod"),
						CtsAbi: []inventory.SchedulableLabels_CTSABI{
							inventory.SchedulableLabels_CTS_ABI_X86,
						},
						CtsCpu: []inventory.SchedulableLabels_CTSCPU{
							inventory.SchedulableLabels_CTS_CPU_X86,
						},
						EcType: ectype(inventory.SchedulableLabels_EC_TYPE_CHROME_OS),
						HwidComponent: []string{
							"FAKE-BATTERY",
							"FAKE-DISPLAY",
						},
						HwidSku: s("FAKE-SKU"),
						Model:   s("FAKE-MODEL"),
						Sku:     s("FAKE-DEVICE-SKU"),
						OsType:  ostype(inventory.SchedulableLabels_OS_TYPE_CROS),
						Peripherals: &inventory.Peripherals{
							AudioBoard:          nil,
							AudioBox:            nil,
							AudioCable:          nil,
							AudioLoopbackDongle: b(true),
							Chameleon:           nil,
							Conductive:          nil,
							Huddly:              nil,
							Mimo:                nil,
							Servo:               b(true),
							ServoState:          peripheralState(inventory.PeripheralState_WORKING),
							ServoType:           s("servo_v4_with_ccd_cr50"),
							ServoTopology: &inventory.ServoTopology{
								Main: &inventory.ServoTopologyItem{
									Type:         s("servo_v4"),
									SysfsProduct: s("Servo V4"),
									Serial:       s("FAKE-SERVO-SERIAL"),
									UsbHubPort:   s("FAKE-SERVO-USB-HUB-PORT"),
								},
								Children: []*inventory.ServoTopologyItem{
									{
										Type:         s("ccd_cr50"),
										SysfsProduct: s("Cr50"),
										Serial:       s("FAKE-TOPOLOGY-ITEM"),
										UsbHubPort:   s("FAKE-USB-HUB-PORT"),
									},
								},
							},
							SmartUsbhub:     nil,
							Camerabox:       nil,
							Wificell:        nil,
							Router_802_11Ax: nil,
							StorageState:    hardwarestate(inventory.HardwareState_HARDWARE_NORMAL),
							ServoUsbState:   hardwarestate(inventory.HardwareState_HARDWARE_NORMAL),
						},
						Phase:    phase(inventory.SchedulableLabels_PHASE_DVT),
						Platform: s("FAKE-PLATFORM"),
						SelfServePools: []string{
							"faft-test",
						},
						TestCoverageHints: &inventory.TestCoverageHints{
							ChaosDut: nil,
						},
						Variant: []string{
							"FAKE-VARIANT",
						},
						WifiChip: s("FAKE-CHIP"),
					},
				},
			},
		},
		&FakeGetStableVersion{
			cros:      "FAKE-CROS-VERSION",
			faft:      "FAKE-FAFT-VERSION",
			firmware:  "FAKE-FIRMWARE-VERSION",
			servoCros: "FAKE-SERVO-CROS-VERSION",
		},
	)

	out, e := g.GetContentsForHostname(bg, hostname)
	eMsg := errToString(e)

	if diff := cmp.Diff(expected, out); diff != "" {
		t.Errorf("wanted: (%s) got: (%s)\n(%s)", expected, out, diff)
	}

	if diff := cmp.Diff(expectedErr, eMsg); diff != "" {
		t.Errorf("wanted: (%s) got: (%s)\n(%s)", expectedErr, eMsg, diff)
	}
}

func errToString(e error) string {
	if e == nil {
		return ""
	}
	return fmt.Sprintf("<%s>", e.Error())
}

// Pointer to string, for building protos.
func s(s string) *string {
	return &s
}

// Pointer to bool, for building protos.
func b(b bool) *bool {
	return &b
}

// Pointer to Cr50 Phase, for building protos.
func cr50(cr50 inventory.SchedulableLabels_CR50_Phase) *inventory.SchedulableLabels_CR50_Phase {
	return &cr50
}

// Pointer to ECType, for building protos.
func ectype(ectype inventory.SchedulableLabels_ECType) *inventory.SchedulableLabels_ECType {
	return &ectype
}

// Pointer to OSType, for building protos.
func ostype(ostype inventory.SchedulableLabels_OSType) *inventory.SchedulableLabels_OSType {
	return &ostype
}

// Pointer to peripheral state, for building protos.
func peripheralState(peripheralState inventory.PeripheralState) *inventory.PeripheralState {
	return &peripheralState
}

// Pointer to hardware state, for building protos.
func hardwarestate(hardwareState inventory.HardwareState) *inventory.HardwareState {
	return &hardwareState
}

// Pointer to phase, for building protos.
func phase(phase inventory.SchedulableLabels_Phase) *inventory.SchedulableLabels_Phase {
	return &phase
}
