// Copyright 2020 The Chromium OS Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.12.1
// source: infra/unifiedfleet/api/v1/proto/state.proto

package ufspb

import (
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

// Next tag: 10
type State int32

const (
	State_STATE_UNSPECIFIED State = 0
	// Equlavant to the concept in ChromeOS lab: needs_deploy
	State_STATE_REGISTERED State = 1
	// Deployed but not placed in prod. It's only a temporarily state for browser machine
	// as there's no service to push a deployed machine to prod automatically yet.
	State_STATE_DEPLOYED_PRE_SERVING State = 9
	// Deployed to the prod infrastructure, but for testing.
	State_STATE_DEPLOYED_TESTING State = 2
	// Deployed to the prod infrastructure, serving.
	State_STATE_SERVING State = 3
	// Deployed to the prod infrastructure, but needs repair.
	State_STATE_NEEDS_REPAIR State = 5
	// Deployed to the prod infrastructure, but get disabled.
	State_STATE_DISABLED State = 6
	// Deployed to the prod infrastructure, but get reserved (e.g. locked).
	State_STATE_RESERVED State = 7
	// Decommissioned from the prod infrastructure, but still leave in UFS record.
	State_STATE_DECOMMISSIONED State = 8
)

// Enum value maps for State.
var (
	State_name = map[int32]string{
		0: "STATE_UNSPECIFIED",
		1: "STATE_REGISTERED",
		9: "STATE_DEPLOYED_PRE_SERVING",
		2: "STATE_DEPLOYED_TESTING",
		3: "STATE_SERVING",
		5: "STATE_NEEDS_REPAIR",
		6: "STATE_DISABLED",
		7: "STATE_RESERVED",
		8: "STATE_DECOMMISSIONED",
	}
	State_value = map[string]int32{
		"STATE_UNSPECIFIED":          0,
		"STATE_REGISTERED":           1,
		"STATE_DEPLOYED_PRE_SERVING": 9,
		"STATE_DEPLOYED_TESTING":     2,
		"STATE_SERVING":              3,
		"STATE_NEEDS_REPAIR":         5,
		"STATE_DISABLED":             6,
		"STATE_RESERVED":             7,
		"STATE_DECOMMISSIONED":       8,
	}
)

func (x State) Enum() *State {
	p := new(State)
	*p = x
	return p
}

func (x State) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (State) Descriptor() protoreflect.EnumDescriptor {
	return file_infra_unifiedfleet_api_v1_proto_state_proto_enumTypes[0].Descriptor()
}

func (State) Type() protoreflect.EnumType {
	return &file_infra_unifiedfleet_api_v1_proto_state_proto_enumTypes[0]
}

func (x State) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use State.Descriptor instead.
func (State) EnumDescriptor() ([]byte, []int) {
	return file_infra_unifiedfleet_api_v1_proto_state_proto_rawDescGZIP(), []int{0}
}

// There's no exposed API for users to directly retrieve a state record.
//
// Ideally, state record can only be modified internally by UFS after some essential
// preconditions are fulfilled.
//
// Users will focus on the tasks triggered by any state change instead of the state
// itself, e.g. once the state of a machine is changed to registered, lab admins will
// know it by founding more machines are listed for waiting for further configurations,
// instead of actively monitoring it by any tooling.
type StateRecord struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The string resource_name could be an ID of a rack, machine, RPM and switches.
	// It can also be the ID of virtual concepts, e.g. LSE and vlan.
	// The format of the resource name will be “racks/XXX” or “rpms/XXX” to help to
	// distinguish the type of the resource.
	ResourceName string `protobuf:"bytes,1,opt,name=resource_name,json=resourceName,proto3" json:"resource_name,omitempty"`
	State        State  `protobuf:"varint,2,opt,name=state,proto3,enum=unifiedfleet.api.v1.proto.State" json:"state,omitempty"`
	User         string `protobuf:"bytes,3,opt,name=user,proto3" json:"user,omitempty"`
	Ticket       string `protobuf:"bytes,4,opt,name=ticket,proto3" json:"ticket,omitempty"`
	Description  string `protobuf:"bytes,5,opt,name=description,proto3" json:"description,omitempty"`
	// Record the last update timestamp of this machine (In UTC timezone)
	UpdateTime *timestamp.Timestamp `protobuf:"bytes,6,opt,name=update_time,json=updateTime,proto3" json:"update_time,omitempty"`
}

func (x *StateRecord) Reset() {
	*x = StateRecord{}
	if protoimpl.UnsafeEnabled {
		mi := &file_infra_unifiedfleet_api_v1_proto_state_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StateRecord) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StateRecord) ProtoMessage() {}

func (x *StateRecord) ProtoReflect() protoreflect.Message {
	mi := &file_infra_unifiedfleet_api_v1_proto_state_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StateRecord.ProtoReflect.Descriptor instead.
func (*StateRecord) Descriptor() ([]byte, []int) {
	return file_infra_unifiedfleet_api_v1_proto_state_proto_rawDescGZIP(), []int{0}
}

func (x *StateRecord) GetResourceName() string {
	if x != nil {
		return x.ResourceName
	}
	return ""
}

func (x *StateRecord) GetState() State {
	if x != nil {
		return x.State
	}
	return State_STATE_UNSPECIFIED
}

func (x *StateRecord) GetUser() string {
	if x != nil {
		return x.User
	}
	return ""
}

func (x *StateRecord) GetTicket() string {
	if x != nil {
		return x.Ticket
	}
	return ""
}

func (x *StateRecord) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *StateRecord) GetUpdateTime() *timestamp.Timestamp {
	if x != nil {
		return x.UpdateTime
	}
	return nil
}

var File_infra_unifiedfleet_api_v1_proto_state_proto protoreflect.FileDescriptor

var file_infra_unifiedfleet_api_v1_proto_state_proto_rawDesc = []byte{
	0x0a, 0x2b, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x2f, 0x75, 0x6e, 0x69, 0x66, 0x69, 0x65, 0x64, 0x66,
	0x6c, 0x65, 0x65, 0x74, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x19, 0x75,
	0x6e, 0x69, 0x66, 0x69, 0x65, 0x64, 0x66, 0x6c, 0x65, 0x65, 0x74, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x76, 0x31, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xf5, 0x01, 0x0a, 0x0b, 0x53, 0x74,
	0x61, 0x74, 0x65, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x12, 0x23, 0x0a, 0x0d, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0c, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x36,
	0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x20, 0x2e,
	0x75, 0x6e, 0x69, 0x66, 0x69, 0x65, 0x64, 0x66, 0x6c, 0x65, 0x65, 0x74, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x76, 0x31, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52,
	0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x75, 0x73, 0x65, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x74, 0x69,
	0x63, 0x6b, 0x65, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x74, 0x69, 0x63, 0x6b,
	0x65, 0x74, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x3b, 0x0a, 0x0b, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x5f, 0x74,
	0x69, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d,
	0x65, 0x2a, 0xdd, 0x01, 0x0a, 0x05, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x15, 0x0a, 0x11, 0x53,
	0x54, 0x41, 0x54, 0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44,
	0x10, 0x00, 0x12, 0x14, 0x0a, 0x10, 0x53, 0x54, 0x41, 0x54, 0x45, 0x5f, 0x52, 0x45, 0x47, 0x49,
	0x53, 0x54, 0x45, 0x52, 0x45, 0x44, 0x10, 0x01, 0x12, 0x1e, 0x0a, 0x1a, 0x53, 0x54, 0x41, 0x54,
	0x45, 0x5f, 0x44, 0x45, 0x50, 0x4c, 0x4f, 0x59, 0x45, 0x44, 0x5f, 0x50, 0x52, 0x45, 0x5f, 0x53,
	0x45, 0x52, 0x56, 0x49, 0x4e, 0x47, 0x10, 0x09, 0x12, 0x1a, 0x0a, 0x16, 0x53, 0x54, 0x41, 0x54,
	0x45, 0x5f, 0x44, 0x45, 0x50, 0x4c, 0x4f, 0x59, 0x45, 0x44, 0x5f, 0x54, 0x45, 0x53, 0x54, 0x49,
	0x4e, 0x47, 0x10, 0x02, 0x12, 0x11, 0x0a, 0x0d, 0x53, 0x54, 0x41, 0x54, 0x45, 0x5f, 0x53, 0x45,
	0x52, 0x56, 0x49, 0x4e, 0x47, 0x10, 0x03, 0x12, 0x16, 0x0a, 0x12, 0x53, 0x54, 0x41, 0x54, 0x45,
	0x5f, 0x4e, 0x45, 0x45, 0x44, 0x53, 0x5f, 0x52, 0x45, 0x50, 0x41, 0x49, 0x52, 0x10, 0x05, 0x12,
	0x12, 0x0a, 0x0e, 0x53, 0x54, 0x41, 0x54, 0x45, 0x5f, 0x44, 0x49, 0x53, 0x41, 0x42, 0x4c, 0x45,
	0x44, 0x10, 0x06, 0x12, 0x12, 0x0a, 0x0e, 0x53, 0x54, 0x41, 0x54, 0x45, 0x5f, 0x52, 0x45, 0x53,
	0x45, 0x52, 0x56, 0x45, 0x44, 0x10, 0x07, 0x12, 0x18, 0x0a, 0x14, 0x53, 0x54, 0x41, 0x54, 0x45,
	0x5f, 0x44, 0x45, 0x43, 0x4f, 0x4d, 0x4d, 0x49, 0x53, 0x53, 0x49, 0x4f, 0x4e, 0x45, 0x44, 0x10,
	0x08, 0x42, 0x27, 0x5a, 0x25, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x2f, 0x75, 0x6e, 0x69, 0x66, 0x69,
	0x65, 0x64, 0x66, 0x6c, 0x65, 0x65, 0x74, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x3b, 0x75, 0x66, 0x73, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_infra_unifiedfleet_api_v1_proto_state_proto_rawDescOnce sync.Once
	file_infra_unifiedfleet_api_v1_proto_state_proto_rawDescData = file_infra_unifiedfleet_api_v1_proto_state_proto_rawDesc
)

func file_infra_unifiedfleet_api_v1_proto_state_proto_rawDescGZIP() []byte {
	file_infra_unifiedfleet_api_v1_proto_state_proto_rawDescOnce.Do(func() {
		file_infra_unifiedfleet_api_v1_proto_state_proto_rawDescData = protoimpl.X.CompressGZIP(file_infra_unifiedfleet_api_v1_proto_state_proto_rawDescData)
	})
	return file_infra_unifiedfleet_api_v1_proto_state_proto_rawDescData
}

var file_infra_unifiedfleet_api_v1_proto_state_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_infra_unifiedfleet_api_v1_proto_state_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_infra_unifiedfleet_api_v1_proto_state_proto_goTypes = []interface{}{
	(State)(0),                  // 0: unifiedfleet.api.v1.proto.State
	(*StateRecord)(nil),         // 1: unifiedfleet.api.v1.proto.StateRecord
	(*timestamp.Timestamp)(nil), // 2: google.protobuf.Timestamp
}
var file_infra_unifiedfleet_api_v1_proto_state_proto_depIdxs = []int32{
	0, // 0: unifiedfleet.api.v1.proto.StateRecord.state:type_name -> unifiedfleet.api.v1.proto.State
	2, // 1: unifiedfleet.api.v1.proto.StateRecord.update_time:type_name -> google.protobuf.Timestamp
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_infra_unifiedfleet_api_v1_proto_state_proto_init() }
func file_infra_unifiedfleet_api_v1_proto_state_proto_init() {
	if File_infra_unifiedfleet_api_v1_proto_state_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_infra_unifiedfleet_api_v1_proto_state_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StateRecord); i {
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
			RawDescriptor: file_infra_unifiedfleet_api_v1_proto_state_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_infra_unifiedfleet_api_v1_proto_state_proto_goTypes,
		DependencyIndexes: file_infra_unifiedfleet_api_v1_proto_state_proto_depIdxs,
		EnumInfos:         file_infra_unifiedfleet_api_v1_proto_state_proto_enumTypes,
		MessageInfos:      file_infra_unifiedfleet_api_v1_proto_state_proto_msgTypes,
	}.Build()
	File_infra_unifiedfleet_api_v1_proto_state_proto = out.File
	file_infra_unifiedfleet_api_v1_proto_state_proto_rawDesc = nil
	file_infra_unifiedfleet_api_v1_proto_state_proto_goTypes = nil
	file_infra_unifiedfleet_api_v1_proto_state_proto_depIdxs = nil
}
