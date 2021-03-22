// Copyright 2020 The Chromium OS Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.12.4
// source: infra/unifiedfleet/api/v1/models/chopsasset.proto

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

// Next Tag: 3
type ChopsAsset struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Asset's state and location
	Id       string    `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Location *Location `protobuf:"bytes,2,opt,name=location,proto3" json:"location,omitempty"`
}

func (x *ChopsAsset) Reset() {
	*x = ChopsAsset{}
	if protoimpl.UnsafeEnabled {
		mi := &file_infra_unifiedfleet_api_v1_models_chopsasset_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChopsAsset) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChopsAsset) ProtoMessage() {}

func (x *ChopsAsset) ProtoReflect() protoreflect.Message {
	mi := &file_infra_unifiedfleet_api_v1_models_chopsasset_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChopsAsset.ProtoReflect.Descriptor instead.
func (*ChopsAsset) Descriptor() ([]byte, []int) {
	return file_infra_unifiedfleet_api_v1_models_chopsasset_proto_rawDescGZIP(), []int{0}
}

func (x *ChopsAsset) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *ChopsAsset) GetLocation() *Location {
	if x != nil {
		return x.Location
	}
	return nil
}

var File_infra_unifiedfleet_api_v1_models_chopsasset_proto protoreflect.FileDescriptor

var file_infra_unifiedfleet_api_v1_models_chopsasset_proto_rawDesc = []byte{
	0x0a, 0x31, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x2f, 0x75, 0x6e, 0x69, 0x66, 0x69, 0x65, 0x64, 0x66,
	0x6c, 0x65, 0x65, 0x74, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x6f, 0x64, 0x65,
	0x6c, 0x73, 0x2f, 0x63, 0x68, 0x6f, 0x70, 0x73, 0x61, 0x73, 0x73, 0x65, 0x74, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x1a, 0x75, 0x6e, 0x69, 0x66, 0x69, 0x65, 0x64, 0x66, 0x6c, 0x65, 0x65,
	0x74, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x1a,
	0x2f, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x2f, 0x75, 0x6e, 0x69, 0x66, 0x69, 0x65, 0x64, 0x66, 0x6c,
	0x65, 0x65, 0x74, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c,
	0x73, 0x2f, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x5e, 0x0a, 0x0a, 0x43, 0x68, 0x6f, 0x70, 0x73, 0x41, 0x73, 0x73, 0x65, 0x74, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x40,
	0x0a, 0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x24, 0x2e, 0x75, 0x6e, 0x69, 0x66, 0x69, 0x65, 0x64, 0x66, 0x6c, 0x65, 0x65, 0x74, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x2e, 0x4c, 0x6f,
	0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x42, 0x28, 0x5a, 0x26, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x2f, 0x75, 0x6e, 0x69, 0x66, 0x69, 0x65,
	0x64, 0x66, 0x6c, 0x65, 0x65, 0x74, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x6f,
	0x64, 0x65, 0x6c, 0x73, 0x3b, 0x75, 0x66, 0x73, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_infra_unifiedfleet_api_v1_models_chopsasset_proto_rawDescOnce sync.Once
	file_infra_unifiedfleet_api_v1_models_chopsasset_proto_rawDescData = file_infra_unifiedfleet_api_v1_models_chopsasset_proto_rawDesc
)

func file_infra_unifiedfleet_api_v1_models_chopsasset_proto_rawDescGZIP() []byte {
	file_infra_unifiedfleet_api_v1_models_chopsasset_proto_rawDescOnce.Do(func() {
		file_infra_unifiedfleet_api_v1_models_chopsasset_proto_rawDescData = protoimpl.X.CompressGZIP(file_infra_unifiedfleet_api_v1_models_chopsasset_proto_rawDescData)
	})
	return file_infra_unifiedfleet_api_v1_models_chopsasset_proto_rawDescData
}

var file_infra_unifiedfleet_api_v1_models_chopsasset_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_infra_unifiedfleet_api_v1_models_chopsasset_proto_goTypes = []interface{}{
	(*ChopsAsset)(nil), // 0: unifiedfleet.api.v1.models.ChopsAsset
	(*Location)(nil),   // 1: unifiedfleet.api.v1.models.Location
}
var file_infra_unifiedfleet_api_v1_models_chopsasset_proto_depIdxs = []int32{
	1, // 0: unifiedfleet.api.v1.models.ChopsAsset.location:type_name -> unifiedfleet.api.v1.models.Location
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_infra_unifiedfleet_api_v1_models_chopsasset_proto_init() }
func file_infra_unifiedfleet_api_v1_models_chopsasset_proto_init() {
	if File_infra_unifiedfleet_api_v1_models_chopsasset_proto != nil {
		return
	}
	file_infra_unifiedfleet_api_v1_models_location_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_infra_unifiedfleet_api_v1_models_chopsasset_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChopsAsset); i {
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
			RawDescriptor: file_infra_unifiedfleet_api_v1_models_chopsasset_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_infra_unifiedfleet_api_v1_models_chopsasset_proto_goTypes,
		DependencyIndexes: file_infra_unifiedfleet_api_v1_models_chopsasset_proto_depIdxs,
		MessageInfos:      file_infra_unifiedfleet_api_v1_models_chopsasset_proto_msgTypes,
	}.Build()
	File_infra_unifiedfleet_api_v1_models_chopsasset_proto = out.File
	file_infra_unifiedfleet_api_v1_models_chopsasset_proto_rawDesc = nil
	file_infra_unifiedfleet_api_v1_models_chopsasset_proto_goTypes = nil
	file_infra_unifiedfleet_api_v1_models_chopsasset_proto_depIdxs = nil
}
