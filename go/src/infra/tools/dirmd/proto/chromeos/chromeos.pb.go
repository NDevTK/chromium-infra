// Copyright 2021 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: infra/tools/dirmd/proto/chromeos/chromeos.proto

package chromeos

import (
	plan "go.chromium.org/chromiumos/config/go/test/plan"
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

// ChromeOS specific metadata.
type ChromeOS struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cq *ChromeOS_CQ `protobuf:"bytes,1,opt,name=cq,proto3" json:"cq,omitempty"`
}

func (x *ChromeOS) Reset() {
	*x = ChromeOS{}
	if protoimpl.UnsafeEnabled {
		mi := &file_infra_tools_dirmd_proto_chromeos_chromeos_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChromeOS) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChromeOS) ProtoMessage() {}

func (x *ChromeOS) ProtoReflect() protoreflect.Message {
	mi := &file_infra_tools_dirmd_proto_chromeos_chromeos_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChromeOS.ProtoReflect.Descriptor instead.
func (*ChromeOS) Descriptor() ([]byte, []int) {
	return file_infra_tools_dirmd_proto_chromeos_chromeos_proto_rawDescGZIP(), []int{0}
}

func (x *ChromeOS) GetCq() *ChromeOS_CQ {
	if x != nil {
		return x.Cq
	}
	return nil
}

// CQ specific metadata.
type ChromeOS_CQ struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// SourceTestPlans to configure testing specific to source paths.
	SourceTestPlans []*plan.SourceTestPlan `protobuf:"bytes,1,rep,name=source_test_plans,json=sourceTestPlans,proto3" json:"source_test_plans,omitempty"`
}

func (x *ChromeOS_CQ) Reset() {
	*x = ChromeOS_CQ{}
	if protoimpl.UnsafeEnabled {
		mi := &file_infra_tools_dirmd_proto_chromeos_chromeos_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChromeOS_CQ) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChromeOS_CQ) ProtoMessage() {}

func (x *ChromeOS_CQ) ProtoReflect() protoreflect.Message {
	mi := &file_infra_tools_dirmd_proto_chromeos_chromeos_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChromeOS_CQ.ProtoReflect.Descriptor instead.
func (*ChromeOS_CQ) Descriptor() ([]byte, []int) {
	return file_infra_tools_dirmd_proto_chromeos_chromeos_proto_rawDescGZIP(), []int{0, 0}
}

func (x *ChromeOS_CQ) GetSourceTestPlans() []*plan.SourceTestPlan {
	if x != nil {
		return x.SourceTestPlans
	}
	return nil
}

var File_infra_tools_dirmd_proto_chromeos_chromeos_proto protoreflect.FileDescriptor

var file_infra_tools_dirmd_proto_chromeos_chromeos_proto_rawDesc = []byte{
	0x0a, 0x2f, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x2f, 0x74, 0x6f, 0x6f, 0x6c, 0x73, 0x2f, 0x64, 0x69,
	0x72, 0x6d, 0x64, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65,
	0x6f, 0x73, 0x2f, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65, 0x6f, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x1c, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65, 0x2e, 0x64, 0x69, 0x72, 0x5f, 0x6d, 0x65,
	0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65, 0x6f, 0x73, 0x1a,
	0x53, 0x67, 0x6f, 0x2e, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x69, 0x75, 0x6d, 0x2e, 0x6f, 0x72, 0x67,
	0x2f, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x69, 0x75, 0x6d, 0x6f, 0x73, 0x2f, 0x63, 0x6f, 0x6e, 0x66,
	0x69, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x69, 0x75,
	0x6d, 0x6f, 0x73, 0x2f, 0x74, 0x65, 0x73, 0x74, 0x2f, 0x70, 0x6c, 0x61, 0x6e, 0x2f, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x5f, 0x74, 0x65, 0x73, 0x74, 0x5f, 0x70, 0x6c, 0x61, 0x6e, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x9d, 0x01, 0x0a, 0x08, 0x43, 0x68, 0x72, 0x6f, 0x6d, 0x65, 0x4f,
	0x53, 0x12, 0x39, 0x0a, 0x02, 0x63, 0x71, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x29, 0x2e,
	0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65, 0x2e, 0x64, 0x69, 0x72, 0x5f, 0x6d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0x2e, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65, 0x6f, 0x73, 0x2e, 0x43, 0x68, 0x72,
	0x6f, 0x6d, 0x65, 0x4f, 0x53, 0x2e, 0x43, 0x51, 0x52, 0x02, 0x63, 0x71, 0x1a, 0x56, 0x0a, 0x02,
	0x43, 0x51, 0x12, 0x50, 0x0a, 0x11, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x74, 0x65, 0x73,
	0x74, 0x5f, 0x70, 0x6c, 0x61, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x24, 0x2e,
	0x63, 0x68, 0x72, 0x6f, 0x6d, 0x69, 0x75, 0x6d, 0x6f, 0x73, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e,
	0x70, 0x6c, 0x61, 0x6e, 0x2e, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x54, 0x65, 0x73, 0x74, 0x50,
	0x6c, 0x61, 0x6e, 0x52, 0x0f, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x54, 0x65, 0x73, 0x74, 0x50,
	0x6c, 0x61, 0x6e, 0x73, 0x42, 0x22, 0x5a, 0x20, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x2f, 0x74, 0x6f,
	0x6f, 0x6c, 0x73, 0x2f, 0x64, 0x69, 0x72, 0x6d, 0x64, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65, 0x6f, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_infra_tools_dirmd_proto_chromeos_chromeos_proto_rawDescOnce sync.Once
	file_infra_tools_dirmd_proto_chromeos_chromeos_proto_rawDescData = file_infra_tools_dirmd_proto_chromeos_chromeos_proto_rawDesc
)

func file_infra_tools_dirmd_proto_chromeos_chromeos_proto_rawDescGZIP() []byte {
	file_infra_tools_dirmd_proto_chromeos_chromeos_proto_rawDescOnce.Do(func() {
		file_infra_tools_dirmd_proto_chromeos_chromeos_proto_rawDescData = protoimpl.X.CompressGZIP(file_infra_tools_dirmd_proto_chromeos_chromeos_proto_rawDescData)
	})
	return file_infra_tools_dirmd_proto_chromeos_chromeos_proto_rawDescData
}

var file_infra_tools_dirmd_proto_chromeos_chromeos_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_infra_tools_dirmd_proto_chromeos_chromeos_proto_goTypes = []interface{}{
	(*ChromeOS)(nil),            // 0: chrome.dir_metadata.chromeos.ChromeOS
	(*ChromeOS_CQ)(nil),         // 1: chrome.dir_metadata.chromeos.ChromeOS.CQ
	(*plan.SourceTestPlan)(nil), // 2: chromiumos.test.plan.SourceTestPlan
}
var file_infra_tools_dirmd_proto_chromeos_chromeos_proto_depIdxs = []int32{
	1, // 0: chrome.dir_metadata.chromeos.ChromeOS.cq:type_name -> chrome.dir_metadata.chromeos.ChromeOS.CQ
	2, // 1: chrome.dir_metadata.chromeos.ChromeOS.CQ.source_test_plans:type_name -> chromiumos.test.plan.SourceTestPlan
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_infra_tools_dirmd_proto_chromeos_chromeos_proto_init() }
func file_infra_tools_dirmd_proto_chromeos_chromeos_proto_init() {
	if File_infra_tools_dirmd_proto_chromeos_chromeos_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_infra_tools_dirmd_proto_chromeos_chromeos_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChromeOS); i {
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
		file_infra_tools_dirmd_proto_chromeos_chromeos_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChromeOS_CQ); i {
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
			RawDescriptor: file_infra_tools_dirmd_proto_chromeos_chromeos_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_infra_tools_dirmd_proto_chromeos_chromeos_proto_goTypes,
		DependencyIndexes: file_infra_tools_dirmd_proto_chromeos_chromeos_proto_depIdxs,
		MessageInfos:      file_infra_tools_dirmd_proto_chromeos_chromeos_proto_msgTypes,
	}.Build()
	File_infra_tools_dirmd_proto_chromeos_chromeos_proto = out.File
	file_infra_tools_dirmd_proto_chromeos_chromeos_proto_rawDesc = nil
	file_infra_tools_dirmd_proto_chromeos_chromeos_proto_goTypes = nil
	file_infra_tools_dirmd_proto_chromeos_chromeos_proto_depIdxs = nil
}
