// Copyright 2020 The Chromium OS Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.12.1
// source: infra/unifiedfleet/api/v1/proto/chromeos/device/config_id.proto

package ufspb

import (
	proto "github.com/golang/protobuf/proto"
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

// These are the globally unique identifiers that determine what set of
// configuration data is used for a given device.
type ConfigId struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Required.
	PlatformId *PlatformId `protobuf:"bytes,1,opt,name=platform_id,json=platformId,proto3" json:"platform_id,omitempty"`
	// Required.
	ModelId *ModelId `protobuf:"bytes,2,opt,name=model_id,json=modelId,proto3" json:"model_id,omitempty"`
	// Required.
	VariantId *VariantId `protobuf:"bytes,3,opt,name=variant_id,json=variantId,proto3" json:"variant_id,omitempty"`
	// Required.
	BrandId *BrandId `protobuf:"bytes,4,opt,name=brand_id,json=brandId,proto3" json:"brand_id,omitempty"`
}

func (x *ConfigId) Reset() {
	*x = ConfigId{}
	if protoimpl.UnsafeEnabled {
		mi := &file_infra_unifiedfleet_api_v1_proto_chromeos_device_config_id_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConfigId) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConfigId) ProtoMessage() {}

func (x *ConfigId) ProtoReflect() protoreflect.Message {
	mi := &file_infra_unifiedfleet_api_v1_proto_chromeos_device_config_id_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConfigId.ProtoReflect.Descriptor instead.
func (*ConfigId) Descriptor() ([]byte, []int) {
	return file_infra_unifiedfleet_api_v1_proto_chromeos_device_config_id_proto_rawDescGZIP(), []int{0}
}

func (x *ConfigId) GetPlatformId() *PlatformId {
	if x != nil {
		return x.PlatformId
	}
	return nil
}

func (x *ConfigId) GetModelId() *ModelId {
	if x != nil {
		return x.ModelId
	}
	return nil
}

func (x *ConfigId) GetVariantId() *VariantId {
	if x != nil {
		return x.VariantId
	}
	return nil
}

func (x *ConfigId) GetBrandId() *BrandId {
	if x != nil {
		return x.BrandId
	}
	return nil
}

var File_infra_unifiedfleet_api_v1_proto_chromeos_device_config_id_proto protoreflect.FileDescriptor

var file_infra_unifiedfleet_api_v1_proto_chromeos_device_config_id_proto_rawDesc = []byte{
	0x0a, 0x3f, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x2f, 0x75, 0x6e, 0x69, 0x66, 0x69, 0x65, 0x64, 0x66,
	0x6c, 0x65, 0x65, 0x74, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2f, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65, 0x6f, 0x73, 0x2f, 0x64, 0x65, 0x76, 0x69, 0x63,
	0x65, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x5f, 0x69, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x29, 0x75, 0x6e, 0x69, 0x66, 0x69, 0x65, 0x64, 0x66, 0x6c, 0x65, 0x65, 0x74, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x68, 0x72,
	0x6f, 0x6d, 0x65, 0x6f, 0x73, 0x2e, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x1a, 0x41, 0x69, 0x6e,
	0x66, 0x72, 0x61, 0x2f, 0x75, 0x6e, 0x69, 0x66, 0x69, 0x65, 0x64, 0x66, 0x6c, 0x65, 0x65, 0x74,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x68,
	0x72, 0x6f, 0x6d, 0x65, 0x6f, 0x73, 0x2f, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x70, 0x6c,
	0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x5f, 0x69, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x3e, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x2f, 0x75, 0x6e, 0x69, 0x66, 0x69, 0x65, 0x64, 0x66, 0x6c,
	0x65, 0x65, 0x74, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65, 0x6f, 0x73, 0x2f, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65,
	0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x5f, 0x69, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x3e, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x2f, 0x75, 0x6e, 0x69, 0x66, 0x69, 0x65, 0x64, 0x66, 0x6c,
	0x65, 0x65, 0x74, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65, 0x6f, 0x73, 0x2f, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65,
	0x2f, 0x62, 0x72, 0x61, 0x6e, 0x64, 0x5f, 0x69, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x40, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x2f, 0x75, 0x6e, 0x69, 0x66, 0x69, 0x65, 0x64, 0x66, 0x6c,
	0x65, 0x65, 0x74, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65, 0x6f, 0x73, 0x2f, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65,
	0x2f, 0x76, 0x61, 0x72, 0x69, 0x61, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0xd5, 0x02, 0x0a, 0x08, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x49, 0x64, 0x12, 0x56,
	0x0a, 0x0b, 0x70, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x35, 0x2e, 0x75, 0x6e, 0x69, 0x66, 0x69, 0x65, 0x64, 0x66, 0x6c, 0x65,
	0x65, 0x74, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65, 0x6f, 0x73, 0x2e, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x2e,
	0x50, 0x6c, 0x61, 0x74, 0x66, 0x6f, 0x72, 0x6d, 0x49, 0x64, 0x52, 0x0a, 0x70, 0x6c, 0x61, 0x74,
	0x66, 0x6f, 0x72, 0x6d, 0x49, 0x64, 0x12, 0x4d, 0x0a, 0x08, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x5f,
	0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x32, 0x2e, 0x75, 0x6e, 0x69, 0x66, 0x69,
	0x65, 0x64, 0x66, 0x6c, 0x65, 0x65, 0x74, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65, 0x6f, 0x73, 0x2e, 0x64, 0x65,
	0x76, 0x69, 0x63, 0x65, 0x2e, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x49, 0x64, 0x52, 0x07, 0x6d, 0x6f,
	0x64, 0x65, 0x6c, 0x49, 0x64, 0x12, 0x53, 0x0a, 0x0a, 0x76, 0x61, 0x72, 0x69, 0x61, 0x6e, 0x74,
	0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x34, 0x2e, 0x75, 0x6e, 0x69, 0x66,
	0x69, 0x65, 0x64, 0x66, 0x6c, 0x65, 0x65, 0x74, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65, 0x6f, 0x73, 0x2e, 0x64,
	0x65, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x56, 0x61, 0x72, 0x69, 0x61, 0x6e, 0x74, 0x49, 0x64, 0x52,
	0x09, 0x76, 0x61, 0x72, 0x69, 0x61, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x4d, 0x0a, 0x08, 0x62, 0x72,
	0x61, 0x6e, 0x64, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x32, 0x2e, 0x75,
	0x6e, 0x69, 0x66, 0x69, 0x65, 0x64, 0x66, 0x6c, 0x65, 0x65, 0x74, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x76, 0x31, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65, 0x6f,
	0x73, 0x2e, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x42, 0x72, 0x61, 0x6e, 0x64, 0x49, 0x64,
	0x52, 0x07, 0x62, 0x72, 0x61, 0x6e, 0x64, 0x49, 0x64, 0x42, 0x37, 0x5a, 0x35, 0x69, 0x6e, 0x66,
	0x72, 0x61, 0x2f, 0x75, 0x6e, 0x69, 0x66, 0x69, 0x65, 0x64, 0x66, 0x6c, 0x65, 0x65, 0x74, 0x2f,
	0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x68, 0x72,
	0x6f, 0x6d, 0x65, 0x6f, 0x73, 0x2f, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x3b, 0x75, 0x66, 0x73,
	0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_infra_unifiedfleet_api_v1_proto_chromeos_device_config_id_proto_rawDescOnce sync.Once
	file_infra_unifiedfleet_api_v1_proto_chromeos_device_config_id_proto_rawDescData = file_infra_unifiedfleet_api_v1_proto_chromeos_device_config_id_proto_rawDesc
)

func file_infra_unifiedfleet_api_v1_proto_chromeos_device_config_id_proto_rawDescGZIP() []byte {
	file_infra_unifiedfleet_api_v1_proto_chromeos_device_config_id_proto_rawDescOnce.Do(func() {
		file_infra_unifiedfleet_api_v1_proto_chromeos_device_config_id_proto_rawDescData = protoimpl.X.CompressGZIP(file_infra_unifiedfleet_api_v1_proto_chromeos_device_config_id_proto_rawDescData)
	})
	return file_infra_unifiedfleet_api_v1_proto_chromeos_device_config_id_proto_rawDescData
}

var file_infra_unifiedfleet_api_v1_proto_chromeos_device_config_id_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_infra_unifiedfleet_api_v1_proto_chromeos_device_config_id_proto_goTypes = []interface{}{
	(*ConfigId)(nil),   // 0: unifiedfleet.api.v1.proto.chromeos.device.ConfigId
	(*PlatformId)(nil), // 1: unifiedfleet.api.v1.proto.chromeos.device.PlatformId
	(*ModelId)(nil),    // 2: unifiedfleet.api.v1.proto.chromeos.device.ModelId
	(*VariantId)(nil),  // 3: unifiedfleet.api.v1.proto.chromeos.device.VariantId
	(*BrandId)(nil),    // 4: unifiedfleet.api.v1.proto.chromeos.device.BrandId
}
var file_infra_unifiedfleet_api_v1_proto_chromeos_device_config_id_proto_depIdxs = []int32{
	1, // 0: unifiedfleet.api.v1.proto.chromeos.device.ConfigId.platform_id:type_name -> unifiedfleet.api.v1.proto.chromeos.device.PlatformId
	2, // 1: unifiedfleet.api.v1.proto.chromeos.device.ConfigId.model_id:type_name -> unifiedfleet.api.v1.proto.chromeos.device.ModelId
	3, // 2: unifiedfleet.api.v1.proto.chromeos.device.ConfigId.variant_id:type_name -> unifiedfleet.api.v1.proto.chromeos.device.VariantId
	4, // 3: unifiedfleet.api.v1.proto.chromeos.device.ConfigId.brand_id:type_name -> unifiedfleet.api.v1.proto.chromeos.device.BrandId
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_infra_unifiedfleet_api_v1_proto_chromeos_device_config_id_proto_init() }
func file_infra_unifiedfleet_api_v1_proto_chromeos_device_config_id_proto_init() {
	if File_infra_unifiedfleet_api_v1_proto_chromeos_device_config_id_proto != nil {
		return
	}
	file_infra_unifiedfleet_api_v1_proto_chromeos_device_platform_id_proto_init()
	file_infra_unifiedfleet_api_v1_proto_chromeos_device_model_id_proto_init()
	file_infra_unifiedfleet_api_v1_proto_chromeos_device_brand_id_proto_init()
	file_infra_unifiedfleet_api_v1_proto_chromeos_device_variant_id_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_infra_unifiedfleet_api_v1_proto_chromeos_device_config_id_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConfigId); i {
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
			RawDescriptor: file_infra_unifiedfleet_api_v1_proto_chromeos_device_config_id_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_infra_unifiedfleet_api_v1_proto_chromeos_device_config_id_proto_goTypes,
		DependencyIndexes: file_infra_unifiedfleet_api_v1_proto_chromeos_device_config_id_proto_depIdxs,
		MessageInfos:      file_infra_unifiedfleet_api_v1_proto_chromeos_device_config_id_proto_msgTypes,
	}.Build()
	File_infra_unifiedfleet_api_v1_proto_chromeos_device_config_id_proto = out.File
	file_infra_unifiedfleet_api_v1_proto_chromeos_device_config_id_proto_rawDesc = nil
	file_infra_unifiedfleet_api_v1_proto_chromeos_device_config_id_proto_goTypes = nil
	file_infra_unifiedfleet_api_v1_proto_chromeos_device_config_id_proto_depIdxs = nil
}
