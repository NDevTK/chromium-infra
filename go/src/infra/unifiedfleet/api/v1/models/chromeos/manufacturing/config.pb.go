// Copyright 2021 The Chromium OS Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.12.1
// source: infra/unifiedfleet/api/v1/models/chromeos/manufacturing/config.proto

package ufspb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// phase for the device. Next Tag: 9
type Config_Phase int32

const (
	Config_PHASE_INVALID Config_Phase = 0
	Config_PHASE_EVT     Config_Phase = 1
	Config_PHASE_EVT2    Config_Phase = 2
	Config_PHASE_DVT     Config_Phase = 3
	Config_PHASE_DVT2    Config_Phase = 4
	Config_PHASE_PVT     Config_Phase = 5
	Config_PHASE_PVT2    Config_Phase = 6
	Config_PHASE_PVT3    Config_Phase = 7
	Config_PHASE_MP      Config_Phase = 8
)

// Enum value maps for Config_Phase.
var (
	Config_Phase_name = map[int32]string{
		0: "PHASE_INVALID",
		1: "PHASE_EVT",
		2: "PHASE_EVT2",
		3: "PHASE_DVT",
		4: "PHASE_DVT2",
		5: "PHASE_PVT",
		6: "PHASE_PVT2",
		7: "PHASE_PVT3",
		8: "PHASE_MP",
	}
	Config_Phase_value = map[string]int32{
		"PHASE_INVALID": 0,
		"PHASE_EVT":     1,
		"PHASE_EVT2":    2,
		"PHASE_DVT":     3,
		"PHASE_DVT2":    4,
		"PHASE_PVT":     5,
		"PHASE_PVT2":    6,
		"PHASE_PVT3":    7,
		"PHASE_MP":      8,
	}
)

func (x Config_Phase) Enum() *Config_Phase {
	p := new(Config_Phase)
	*p = x
	return p
}

func (x Config_Phase) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Config_Phase) Descriptor() protoreflect.EnumDescriptor {
	return file_infra_unifiedfleet_api_v1_models_chromeos_manufacturing_config_proto_enumTypes[0].Descriptor()
}

func (Config_Phase) Type() protoreflect.EnumType {
	return &file_infra_unifiedfleet_api_v1_models_chromeos_manufacturing_config_proto_enumTypes[0]
}

func (x Config_Phase) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Config_Phase.Descriptor instead.
func (Config_Phase) EnumDescriptor() ([]byte, []int) {
	return file_infra_unifiedfleet_api_v1_models_chromeos_manufacturing_config_proto_rawDescGZIP(), []int{0, 0}
}

// phases for cr50 module. Next Tag: 3
type Config_CR50Phase int32

const (
	Config_CR50_PHASE_INVALID Config_CR50Phase = 0
	Config_CR50_PHASE_PREPVT  Config_CR50Phase = 1
	Config_CR50_PHASE_PVT     Config_CR50Phase = 2
)

// Enum value maps for Config_CR50Phase.
var (
	Config_CR50Phase_name = map[int32]string{
		0: "CR50_PHASE_INVALID",
		1: "CR50_PHASE_PREPVT",
		2: "CR50_PHASE_PVT",
	}
	Config_CR50Phase_value = map[string]int32{
		"CR50_PHASE_INVALID": 0,
		"CR50_PHASE_PREPVT":  1,
		"CR50_PHASE_PVT":     2,
	}
)

func (x Config_CR50Phase) Enum() *Config_CR50Phase {
	p := new(Config_CR50Phase)
	*p = x
	return p
}

func (x Config_CR50Phase) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Config_CR50Phase) Descriptor() protoreflect.EnumDescriptor {
	return file_infra_unifiedfleet_api_v1_models_chromeos_manufacturing_config_proto_enumTypes[1].Descriptor()
}

func (Config_CR50Phase) Type() protoreflect.EnumType {
	return &file_infra_unifiedfleet_api_v1_models_chromeos_manufacturing_config_proto_enumTypes[1]
}

func (x Config_CR50Phase) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Config_CR50Phase.Descriptor instead.
func (Config_CR50Phase) EnumDescriptor() ([]byte, []int) {
	return file_infra_unifiedfleet_api_v1_models_chromeos_manufacturing_config_proto_rawDescGZIP(), []int{0, 1}
}

// key env for cr50 RW version. Next Tag: 3
type Config_CR50KeyEnv int32

const (
	Config_CR50_KEYENV_INVALID Config_CR50KeyEnv = 0
	Config_CR50_KEYENV_PROD    Config_CR50KeyEnv = 1
	Config_CR50_KEYENV_DEV     Config_CR50KeyEnv = 2
)

// Enum value maps for Config_CR50KeyEnv.
var (
	Config_CR50KeyEnv_name = map[int32]string{
		0: "CR50_KEYENV_INVALID",
		1: "CR50_KEYENV_PROD",
		2: "CR50_KEYENV_DEV",
	}
	Config_CR50KeyEnv_value = map[string]int32{
		"CR50_KEYENV_INVALID": 0,
		"CR50_KEYENV_PROD":    1,
		"CR50_KEYENV_DEV":     2,
	}
)

func (x Config_CR50KeyEnv) Enum() *Config_CR50KeyEnv {
	p := new(Config_CR50KeyEnv)
	*p = x
	return p
}

func (x Config_CR50KeyEnv) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Config_CR50KeyEnv) Descriptor() protoreflect.EnumDescriptor {
	return file_infra_unifiedfleet_api_v1_models_chromeos_manufacturing_config_proto_enumTypes[2].Descriptor()
}

func (Config_CR50KeyEnv) Type() protoreflect.EnumType {
	return &file_infra_unifiedfleet_api_v1_models_chromeos_manufacturing_config_proto_enumTypes[2]
}

func (x Config_CR50KeyEnv) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Config_CR50KeyEnv.Descriptor instead.
func (Config_CR50KeyEnv) EnumDescriptor() ([]byte, []int) {
	return file_infra_unifiedfleet_api_v1_models_chromeos_manufacturing_config_proto_rawDescGZIP(), []int{0, 2}
}

// These are the configs that's provided in manufacture stage of a ChromeOS device.
// Next Tag: 7
type Config struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ManufacturingId *ConfigID        `protobuf:"bytes,1,opt,name=manufacturing_id,json=manufacturingId,proto3" json:"manufacturing_id,omitempty"`
	DevicePhase     Config_Phase     `protobuf:"varint,2,opt,name=device_phase,json=devicePhase,proto3,enum=unifiedfleet.api.v1.models.chromeos.manufacturing.Config_Phase" json:"device_phase,omitempty"`
	Cr50Phase       Config_CR50Phase `protobuf:"varint,3,opt,name=cr50_phase,json=cr50Phase,proto3,enum=unifiedfleet.api.v1.models.chromeos.manufacturing.Config_CR50Phase" json:"cr50_phase,omitempty"`
	// Detected based on the cr50 RW version that the DUT is running on.
	Cr50KeyEnv Config_CR50KeyEnv `protobuf:"varint,4,opt,name=cr50_key_env,json=cr50KeyEnv,proto3,enum=unifiedfleet.api.v1.models.chromeos.manufacturing.Config_CR50KeyEnv" json:"cr50_key_env,omitempty"`
	// wifi chip that is installed on the DUT in manufacturing stage.
	WifiChip string `protobuf:"bytes,5,opt,name=wifi_chip,json=wifiChip,proto3" json:"wifi_chip,omitempty"`
	// Save repeated hwid components obtained from hwid service
	HwidComponent []string `protobuf:"bytes,6,rep,name=hwid_component,json=hwidComponent,proto3" json:"hwid_component,omitempty"`
}

func (x *Config) Reset() {
	*x = Config{}
	if protoimpl.UnsafeEnabled {
		mi := &file_infra_unifiedfleet_api_v1_models_chromeos_manufacturing_config_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Config) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Config) ProtoMessage() {}

func (x *Config) ProtoReflect() protoreflect.Message {
	mi := &file_infra_unifiedfleet_api_v1_models_chromeos_manufacturing_config_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Config.ProtoReflect.Descriptor instead.
func (*Config) Descriptor() ([]byte, []int) {
	return file_infra_unifiedfleet_api_v1_models_chromeos_manufacturing_config_proto_rawDescGZIP(), []int{0}
}

func (x *Config) GetManufacturingId() *ConfigID {
	if x != nil {
		return x.ManufacturingId
	}
	return nil
}

func (x *Config) GetDevicePhase() Config_Phase {
	if x != nil {
		return x.DevicePhase
	}
	return Config_PHASE_INVALID
}

func (x *Config) GetCr50Phase() Config_CR50Phase {
	if x != nil {
		return x.Cr50Phase
	}
	return Config_CR50_PHASE_INVALID
}

func (x *Config) GetCr50KeyEnv() Config_CR50KeyEnv {
	if x != nil {
		return x.Cr50KeyEnv
	}
	return Config_CR50_KEYENV_INVALID
}

func (x *Config) GetWifiChip() string {
	if x != nil {
		return x.WifiChip
	}
	return ""
}

func (x *Config) GetHwidComponent() []string {
	if x != nil {
		return x.HwidComponent
	}
	return nil
}

// Message contains all ChromeOS manufacturing configs.
type ConfigList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value []*Config `protobuf:"bytes,1,rep,name=value,proto3" json:"value,omitempty"`
}

func (x *ConfigList) Reset() {
	*x = ConfigList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_infra_unifiedfleet_api_v1_models_chromeos_manufacturing_config_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConfigList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConfigList) ProtoMessage() {}

func (x *ConfigList) ProtoReflect() protoreflect.Message {
	mi := &file_infra_unifiedfleet_api_v1_models_chromeos_manufacturing_config_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConfigList.ProtoReflect.Descriptor instead.
func (*ConfigList) Descriptor() ([]byte, []int) {
	return file_infra_unifiedfleet_api_v1_models_chromeos_manufacturing_config_proto_rawDescGZIP(), []int{1}
}

func (x *ConfigList) GetValue() []*Config {
	if x != nil {
		return x.Value
	}
	return nil
}

var File_infra_unifiedfleet_api_v1_models_chromeos_manufacturing_config_proto protoreflect.FileDescriptor

var file_infra_unifiedfleet_api_v1_models_chromeos_manufacturing_config_proto_rawDesc = []byte{
	0x0a, 0x44, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x2f, 0x75, 0x6e, 0x69, 0x66, 0x69, 0x65, 0x64, 0x66,
	0x6c, 0x65, 0x65, 0x74, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x6f, 0x64, 0x65,
	0x6c, 0x73, 0x2f, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65, 0x6f, 0x73, 0x2f, 0x6d, 0x61, 0x6e, 0x75,
	0x66, 0x61, 0x63, 0x74, 0x75, 0x72, 0x69, 0x6e, 0x67, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x31, 0x75, 0x6e, 0x69, 0x66, 0x69, 0x65, 0x64, 0x66,
	0x6c, 0x65, 0x65, 0x74, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x6d, 0x6f, 0x64, 0x65,
	0x6c, 0x73, 0x2e, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65, 0x6f, 0x73, 0x2e, 0x6d, 0x61, 0x6e, 0x75,
	0x66, 0x61, 0x63, 0x74, 0x75, 0x72, 0x69, 0x6e, 0x67, 0x1a, 0x47, 0x69, 0x6e, 0x66, 0x72, 0x61,
	0x2f, 0x75, 0x6e, 0x69, 0x66, 0x69, 0x65, 0x64, 0x66, 0x6c, 0x65, 0x65, 0x74, 0x2f, 0x61, 0x70,
	0x69, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2f, 0x63, 0x68, 0x72, 0x6f,
	0x6d, 0x65, 0x6f, 0x73, 0x2f, 0x6d, 0x61, 0x6e, 0x75, 0x66, 0x61, 0x63, 0x74, 0x75, 0x72, 0x69,
	0x6e, 0x67, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x5f, 0x69, 0x64, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x9e, 0x06, 0x0a, 0x06, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x66, 0x0a,
	0x10, 0x6d, 0x61, 0x6e, 0x75, 0x66, 0x61, 0x63, 0x74, 0x75, 0x72, 0x69, 0x6e, 0x67, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x3b, 0x2e, 0x75, 0x6e, 0x69, 0x66, 0x69, 0x65,
	0x64, 0x66, 0x6c, 0x65, 0x65, 0x74, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x6d, 0x6f,
	0x64, 0x65, 0x6c, 0x73, 0x2e, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65, 0x6f, 0x73, 0x2e, 0x6d, 0x61,
	0x6e, 0x75, 0x66, 0x61, 0x63, 0x74, 0x75, 0x72, 0x69, 0x6e, 0x67, 0x2e, 0x43, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x49, 0x44, 0x52, 0x0f, 0x6d, 0x61, 0x6e, 0x75, 0x66, 0x61, 0x63, 0x74, 0x75, 0x72,
	0x69, 0x6e, 0x67, 0x49, 0x64, 0x12, 0x62, 0x0a, 0x0c, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x5f,
	0x70, 0x68, 0x61, 0x73, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x3f, 0x2e, 0x75, 0x6e,
	0x69, 0x66, 0x69, 0x65, 0x64, 0x66, 0x6c, 0x65, 0x65, 0x74, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76,
	0x31, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65, 0x6f,
	0x73, 0x2e, 0x6d, 0x61, 0x6e, 0x75, 0x66, 0x61, 0x63, 0x74, 0x75, 0x72, 0x69, 0x6e, 0x67, 0x2e,
	0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x50, 0x68, 0x61, 0x73, 0x65, 0x52, 0x0b, 0x64, 0x65,
	0x76, 0x69, 0x63, 0x65, 0x50, 0x68, 0x61, 0x73, 0x65, 0x12, 0x62, 0x0a, 0x0a, 0x63, 0x72, 0x35,
	0x30, 0x5f, 0x70, 0x68, 0x61, 0x73, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x43, 0x2e,
	0x75, 0x6e, 0x69, 0x66, 0x69, 0x65, 0x64, 0x66, 0x6c, 0x65, 0x65, 0x74, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x76, 0x31, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x63, 0x68, 0x72, 0x6f, 0x6d,
	0x65, 0x6f, 0x73, 0x2e, 0x6d, 0x61, 0x6e, 0x75, 0x66, 0x61, 0x63, 0x74, 0x75, 0x72, 0x69, 0x6e,
	0x67, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x43, 0x52, 0x35, 0x30, 0x50, 0x68, 0x61,
	0x73, 0x65, 0x52, 0x09, 0x63, 0x72, 0x35, 0x30, 0x50, 0x68, 0x61, 0x73, 0x65, 0x12, 0x66, 0x0a,
	0x0c, 0x63, 0x72, 0x35, 0x30, 0x5f, 0x6b, 0x65, 0x79, 0x5f, 0x65, 0x6e, 0x76, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x44, 0x2e, 0x75, 0x6e, 0x69, 0x66, 0x69, 0x65, 0x64, 0x66, 0x6c, 0x65,
	0x65, 0x74, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73,
	0x2e, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65, 0x6f, 0x73, 0x2e, 0x6d, 0x61, 0x6e, 0x75, 0x66, 0x61,
	0x63, 0x74, 0x75, 0x72, 0x69, 0x6e, 0x67, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x43,
	0x52, 0x35, 0x30, 0x4b, 0x65, 0x79, 0x45, 0x6e, 0x76, 0x52, 0x0a, 0x63, 0x72, 0x35, 0x30, 0x4b,
	0x65, 0x79, 0x45, 0x6e, 0x76, 0x12, 0x1b, 0x0a, 0x09, 0x77, 0x69, 0x66, 0x69, 0x5f, 0x63, 0x68,
	0x69, 0x70, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x77, 0x69, 0x66, 0x69, 0x43, 0x68,
	0x69, 0x70, 0x12, 0x25, 0x0a, 0x0e, 0x68, 0x77, 0x69, 0x64, 0x5f, 0x63, 0x6f, 0x6d, 0x70, 0x6f,
	0x6e, 0x65, 0x6e, 0x74, 0x18, 0x06, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0d, 0x68, 0x77, 0x69, 0x64,
	0x43, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x22, 0x95, 0x01, 0x0a, 0x05, 0x50, 0x68,
	0x61, 0x73, 0x65, 0x12, 0x11, 0x0a, 0x0d, 0x50, 0x48, 0x41, 0x53, 0x45, 0x5f, 0x49, 0x4e, 0x56,
	0x41, 0x4c, 0x49, 0x44, 0x10, 0x00, 0x12, 0x0d, 0x0a, 0x09, 0x50, 0x48, 0x41, 0x53, 0x45, 0x5f,
	0x45, 0x56, 0x54, 0x10, 0x01, 0x12, 0x0e, 0x0a, 0x0a, 0x50, 0x48, 0x41, 0x53, 0x45, 0x5f, 0x45,
	0x56, 0x54, 0x32, 0x10, 0x02, 0x12, 0x0d, 0x0a, 0x09, 0x50, 0x48, 0x41, 0x53, 0x45, 0x5f, 0x44,
	0x56, 0x54, 0x10, 0x03, 0x12, 0x0e, 0x0a, 0x0a, 0x50, 0x48, 0x41, 0x53, 0x45, 0x5f, 0x44, 0x56,
	0x54, 0x32, 0x10, 0x04, 0x12, 0x0d, 0x0a, 0x09, 0x50, 0x48, 0x41, 0x53, 0x45, 0x5f, 0x50, 0x56,
	0x54, 0x10, 0x05, 0x12, 0x0e, 0x0a, 0x0a, 0x50, 0x48, 0x41, 0x53, 0x45, 0x5f, 0x50, 0x56, 0x54,
	0x32, 0x10, 0x06, 0x12, 0x0e, 0x0a, 0x0a, 0x50, 0x48, 0x41, 0x53, 0x45, 0x5f, 0x50, 0x56, 0x54,
	0x33, 0x10, 0x07, 0x12, 0x0c, 0x0a, 0x08, 0x50, 0x48, 0x41, 0x53, 0x45, 0x5f, 0x4d, 0x50, 0x10,
	0x08, 0x22, 0x4e, 0x0a, 0x09, 0x43, 0x52, 0x35, 0x30, 0x50, 0x68, 0x61, 0x73, 0x65, 0x12, 0x16,
	0x0a, 0x12, 0x43, 0x52, 0x35, 0x30, 0x5f, 0x50, 0x48, 0x41, 0x53, 0x45, 0x5f, 0x49, 0x4e, 0x56,
	0x41, 0x4c, 0x49, 0x44, 0x10, 0x00, 0x12, 0x15, 0x0a, 0x11, 0x43, 0x52, 0x35, 0x30, 0x5f, 0x50,
	0x48, 0x41, 0x53, 0x45, 0x5f, 0x50, 0x52, 0x45, 0x50, 0x56, 0x54, 0x10, 0x01, 0x12, 0x12, 0x0a,
	0x0e, 0x43, 0x52, 0x35, 0x30, 0x5f, 0x50, 0x48, 0x41, 0x53, 0x45, 0x5f, 0x50, 0x56, 0x54, 0x10,
	0x02, 0x22, 0x50, 0x0a, 0x0a, 0x43, 0x52, 0x35, 0x30, 0x4b, 0x65, 0x79, 0x45, 0x6e, 0x76, 0x12,
	0x17, 0x0a, 0x13, 0x43, 0x52, 0x35, 0x30, 0x5f, 0x4b, 0x45, 0x59, 0x45, 0x4e, 0x56, 0x5f, 0x49,
	0x4e, 0x56, 0x41, 0x4c, 0x49, 0x44, 0x10, 0x00, 0x12, 0x14, 0x0a, 0x10, 0x43, 0x52, 0x35, 0x30,
	0x5f, 0x4b, 0x45, 0x59, 0x45, 0x4e, 0x56, 0x5f, 0x50, 0x52, 0x4f, 0x44, 0x10, 0x01, 0x12, 0x13,
	0x0a, 0x0f, 0x43, 0x52, 0x35, 0x30, 0x5f, 0x4b, 0x45, 0x59, 0x45, 0x4e, 0x56, 0x5f, 0x44, 0x45,
	0x56, 0x10, 0x02, 0x22, 0x5d, 0x0a, 0x0a, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x4c, 0x69, 0x73,
	0x74, 0x12, 0x4f, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x39, 0x2e, 0x75, 0x6e, 0x69, 0x66, 0x69, 0x65, 0x64, 0x66, 0x6c, 0x65, 0x65, 0x74, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x63, 0x68,
	0x72, 0x6f, 0x6d, 0x65, 0x6f, 0x73, 0x2e, 0x6d, 0x61, 0x6e, 0x75, 0x66, 0x61, 0x63, 0x74, 0x75,
	0x72, 0x69, 0x6e, 0x67, 0x2e, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x42, 0x3f, 0x5a, 0x3d, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x2f, 0x75, 0x6e, 0x69, 0x66,
	0x69, 0x65, 0x64, 0x66, 0x6c, 0x65, 0x65, 0x74, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f,
	0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2f, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65, 0x6f, 0x73, 0x2f,
	0x6d, 0x61, 0x6e, 0x75, 0x66, 0x61, 0x63, 0x74, 0x75, 0x72, 0x69, 0x6e, 0x67, 0x3b, 0x75, 0x66,
	0x73, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_infra_unifiedfleet_api_v1_models_chromeos_manufacturing_config_proto_rawDescOnce sync.Once
	file_infra_unifiedfleet_api_v1_models_chromeos_manufacturing_config_proto_rawDescData = file_infra_unifiedfleet_api_v1_models_chromeos_manufacturing_config_proto_rawDesc
)

func file_infra_unifiedfleet_api_v1_models_chromeos_manufacturing_config_proto_rawDescGZIP() []byte {
	file_infra_unifiedfleet_api_v1_models_chromeos_manufacturing_config_proto_rawDescOnce.Do(func() {
		file_infra_unifiedfleet_api_v1_models_chromeos_manufacturing_config_proto_rawDescData = protoimpl.X.CompressGZIP(file_infra_unifiedfleet_api_v1_models_chromeos_manufacturing_config_proto_rawDescData)
	})
	return file_infra_unifiedfleet_api_v1_models_chromeos_manufacturing_config_proto_rawDescData
}

var file_infra_unifiedfleet_api_v1_models_chromeos_manufacturing_config_proto_enumTypes = make([]protoimpl.EnumInfo, 3)
var file_infra_unifiedfleet_api_v1_models_chromeos_manufacturing_config_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_infra_unifiedfleet_api_v1_models_chromeos_manufacturing_config_proto_goTypes = []interface{}{
	(Config_Phase)(0),      // 0: unifiedfleet.api.v1.models.chromeos.manufacturing.Config.Phase
	(Config_CR50Phase)(0),  // 1: unifiedfleet.api.v1.models.chromeos.manufacturing.Config.CR50Phase
	(Config_CR50KeyEnv)(0), // 2: unifiedfleet.api.v1.models.chromeos.manufacturing.Config.CR50KeyEnv
	(*Config)(nil),         // 3: unifiedfleet.api.v1.models.chromeos.manufacturing.Config
	(*ConfigList)(nil),     // 4: unifiedfleet.api.v1.models.chromeos.manufacturing.ConfigList
	(*ConfigID)(nil),       // 5: unifiedfleet.api.v1.models.chromeos.manufacturing.ConfigID
}
var file_infra_unifiedfleet_api_v1_models_chromeos_manufacturing_config_proto_depIdxs = []int32{
	5, // 0: unifiedfleet.api.v1.models.chromeos.manufacturing.Config.manufacturing_id:type_name -> unifiedfleet.api.v1.models.chromeos.manufacturing.ConfigID
	0, // 1: unifiedfleet.api.v1.models.chromeos.manufacturing.Config.device_phase:type_name -> unifiedfleet.api.v1.models.chromeos.manufacturing.Config.Phase
	1, // 2: unifiedfleet.api.v1.models.chromeos.manufacturing.Config.cr50_phase:type_name -> unifiedfleet.api.v1.models.chromeos.manufacturing.Config.CR50Phase
	2, // 3: unifiedfleet.api.v1.models.chromeos.manufacturing.Config.cr50_key_env:type_name -> unifiedfleet.api.v1.models.chromeos.manufacturing.Config.CR50KeyEnv
	3, // 4: unifiedfleet.api.v1.models.chromeos.manufacturing.ConfigList.value:type_name -> unifiedfleet.api.v1.models.chromeos.manufacturing.Config
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_infra_unifiedfleet_api_v1_models_chromeos_manufacturing_config_proto_init() }
func file_infra_unifiedfleet_api_v1_models_chromeos_manufacturing_config_proto_init() {
	if File_infra_unifiedfleet_api_v1_models_chromeos_manufacturing_config_proto != nil {
		return
	}
	file_infra_unifiedfleet_api_v1_models_chromeos_manufacturing_config_id_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_infra_unifiedfleet_api_v1_models_chromeos_manufacturing_config_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Config); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_infra_unifiedfleet_api_v1_models_chromeos_manufacturing_config_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConfigList); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_infra_unifiedfleet_api_v1_models_chromeos_manufacturing_config_proto_rawDesc,
			NumEnums:      3,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_infra_unifiedfleet_api_v1_models_chromeos_manufacturing_config_proto_goTypes,
		DependencyIndexes: file_infra_unifiedfleet_api_v1_models_chromeos_manufacturing_config_proto_depIdxs,
		EnumInfos:         file_infra_unifiedfleet_api_v1_models_chromeos_manufacturing_config_proto_enumTypes,
		MessageInfos:      file_infra_unifiedfleet_api_v1_models_chromeos_manufacturing_config_proto_msgTypes,
	}.Build()
	File_infra_unifiedfleet_api_v1_models_chromeos_manufacturing_config_proto = out.File
	file_infra_unifiedfleet_api_v1_models_chromeos_manufacturing_config_proto_rawDesc = nil
	file_infra_unifiedfleet_api_v1_models_chromeos_manufacturing_config_proto_goTypes = nil
	file_infra_unifiedfleet_api_v1_models_chromeos_manufacturing_config_proto_depIdxs = nil
}
