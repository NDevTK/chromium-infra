// Copyright 2021 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: infra/appengine/weetbix/internal/config/project_config.proto

package config

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

// ProjectConfig is the project-specific configuration data for Weetbix.
type ProjectConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The monorail configuration to use when filing bugs.
	Monorail *MonorailProject `protobuf:"bytes,1,opt,name=monorail,proto3" json:"monorail,omitempty"`
	// The threshold at which to file bugs. If a cluster's impact exceeds
	// the given threshold, a bug will be filed for it.
	BugFilingThreshold *ImpactThreshold `protobuf:"bytes,2,opt,name=bug_filing_threshold,json=bugFilingThreshold,proto3" json:"bug_filing_threshold,omitempty"`
}

func (x *ProjectConfig) Reset() {
	*x = ProjectConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_infra_appengine_weetbix_internal_config_project_config_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProjectConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProjectConfig) ProtoMessage() {}

func (x *ProjectConfig) ProtoReflect() protoreflect.Message {
	mi := &file_infra_appengine_weetbix_internal_config_project_config_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProjectConfig.ProtoReflect.Descriptor instead.
func (*ProjectConfig) Descriptor() ([]byte, []int) {
	return file_infra_appengine_weetbix_internal_config_project_config_proto_rawDescGZIP(), []int{0}
}

func (x *ProjectConfig) GetMonorail() *MonorailProject {
	if x != nil {
		return x.Monorail
	}
	return nil
}

func (x *ProjectConfig) GetBugFilingThreshold() *ImpactThreshold {
	if x != nil {
		return x.BugFilingThreshold
	}
	return nil
}

// MonorailProject describes the configuration to use when filing bugs
// into a given monorail project.
type MonorailProject struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The monorail project being described.
	// E.g. "chromium".
	Project string `protobuf:"bytes,1,opt,name=project,proto3" json:"project,omitempty"`
	// The field values to use when creating new bugs.
	// For example, on chromium issue tracker, there is a manadatory
	// issue type field (field 10), which must be set to "Bug".
	DefaultFieldValues []*MonorailFieldValue `protobuf:"bytes,2,rep,name=default_field_values,json=defaultFieldValues,proto3" json:"default_field_values,omitempty"`
	// The ID of the issue's priority field. You can find this by visiting
	// https://monorail-prod.appspot.com/p/<project>/adminLabels, scrolling
	// down to Custom fields and finding the ID of the field you wish to set.
	PriorityFieldId int64 `protobuf:"varint,3,opt,name=priority_field_id,json=priorityFieldId,proto3" json:"priority_field_id,omitempty"`
	// The possible bug priorities and their associated impact thresholds.
	// Priorities must be listed from highest (i.e. P0) to lowest (i.e. P3).
	// Higher priorities can only be reached if the thresholds for all lower
	// priorities are also met.
	// The impact thresholds for setting the lowest priority implicitly
	// identifies the bug closure threshold -- if no priority can be
	// matched, the bug is closed. Satisfying the threshold for filing bugs MUST
	// at least imply the threshold for the lowest priority, and MAY imply
	// the thresholds of higher priorities.
	Priorities []*MonorailPriority `protobuf:"bytes,4,rep,name=priorities,proto3" json:"priorities,omitempty"`
	// Controls the amount of hysteresis used in setting bug priorities.
	// Once a bug is assigned a given priority, its priority will only be
	// increased if it exceeds the next priority's thresholds by the
	// specified percentage margin, and decreased if the current priority's
	// thresholds exceed the bug's impact by the given percentage margin.
	//
	// A value of 100 indicates impact may be double the threshold for
	// the next highest priority value, (or half the threshold of the
	// current priority value,) before a bug's priority is increased
	// (or decreased).
	//
	// Valid values are from 0 (no hystersis) to 1,000 (10x hysteresis).
	PriorityHysteresisPercent int64 `protobuf:"varint,5,opt,name=priority_hysteresis_percent,json=priorityHysteresisPercent,proto3" json:"priority_hysteresis_percent,omitempty"`
}

func (x *MonorailProject) Reset() {
	*x = MonorailProject{}
	if protoimpl.UnsafeEnabled {
		mi := &file_infra_appengine_weetbix_internal_config_project_config_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MonorailProject) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MonorailProject) ProtoMessage() {}

func (x *MonorailProject) ProtoReflect() protoreflect.Message {
	mi := &file_infra_appengine_weetbix_internal_config_project_config_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MonorailProject.ProtoReflect.Descriptor instead.
func (*MonorailProject) Descriptor() ([]byte, []int) {
	return file_infra_appengine_weetbix_internal_config_project_config_proto_rawDescGZIP(), []int{1}
}

func (x *MonorailProject) GetProject() string {
	if x != nil {
		return x.Project
	}
	return ""
}

func (x *MonorailProject) GetDefaultFieldValues() []*MonorailFieldValue {
	if x != nil {
		return x.DefaultFieldValues
	}
	return nil
}

func (x *MonorailProject) GetPriorityFieldId() int64 {
	if x != nil {
		return x.PriorityFieldId
	}
	return 0
}

func (x *MonorailProject) GetPriorities() []*MonorailPriority {
	if x != nil {
		return x.Priorities
	}
	return nil
}

func (x *MonorailProject) GetPriorityHysteresisPercent() int64 {
	if x != nil {
		return x.PriorityHysteresisPercent
	}
	return 0
}

// MonorailFieldValue describes a monorail field/value pair.
type MonorailFieldValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The ID of the field to set. You can find this by visiting
	// https://monorail-prod.appspot.com/p/<project>/adminLabels, scrolling
	// down to Custom fields and finding the ID of the field you wish to set.
	FieldId int64 `protobuf:"varint,1,opt,name=field_id,json=fieldId,proto3" json:"field_id,omitempty"`
	// The field value. Values are encoded according to the field type:
	// - Enumeration types: the string enumeration value (e.g. "Bug").
	// - Integer types: the integer, converted to a string (e.g. "1052").
	// - String types: the value, included verbatim.
	// - User types: the user's resource name (e.g. "users/2627516260").
	//   User IDs can be identified by looking at the people listing for a
	//   project:  https://monorail-prod.appspot.com/p/<project>/people/list.
	//   The User ID is included in the URL as u=<number> when clicking into
	//   the page for a particular user. For example, "user/3816576959" is
	//   https://monorail-prod.appspot.com/p/chromium/people/detail?u=3816576959.
	// - Date types: the number of seconds since epoch, as a string
	//   (e.g. "1609459200" for 1 January 2021).
	// - URL type: the URL value, as a string (e.g. "https://www.google.com/").
	//
	// The source of truth for mapping of field types to values is as
	// defined in the Monorail v3 API, found here:
	// https://source.chromium.org/chromium/infra/infra/+/main:appengine/monorail/api/v3/api_proto/issue_objects.proto?q=%22message%20FieldValue%22
	Value string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *MonorailFieldValue) Reset() {
	*x = MonorailFieldValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_infra_appengine_weetbix_internal_config_project_config_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MonorailFieldValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MonorailFieldValue) ProtoMessage() {}

func (x *MonorailFieldValue) ProtoReflect() protoreflect.Message {
	mi := &file_infra_appengine_weetbix_internal_config_project_config_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MonorailFieldValue.ProtoReflect.Descriptor instead.
func (*MonorailFieldValue) Descriptor() ([]byte, []int) {
	return file_infra_appengine_weetbix_internal_config_project_config_proto_rawDescGZIP(), []int{2}
}

func (x *MonorailFieldValue) GetFieldId() int64 {
	if x != nil {
		return x.FieldId
	}
	return 0
}

func (x *MonorailFieldValue) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

// MonorailPriority represents configuration for when to use a given
// priority value in a bug.
type MonorailPriority struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The monorail priority value. For example, "0". This depends on the
	// valid priority field values you have defined in your monorail project.
	Priority string `protobuf:"bytes,1,opt,name=priority,proto3" json:"priority,omitempty"`
	// The threshold at which to apply the priority.
	Threshold *ImpactThreshold `protobuf:"bytes,2,opt,name=threshold,proto3" json:"threshold,omitempty"`
}

func (x *MonorailPriority) Reset() {
	*x = MonorailPriority{}
	if protoimpl.UnsafeEnabled {
		mi := &file_infra_appengine_weetbix_internal_config_project_config_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MonorailPriority) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MonorailPriority) ProtoMessage() {}

func (x *MonorailPriority) ProtoReflect() protoreflect.Message {
	mi := &file_infra_appengine_weetbix_internal_config_project_config_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MonorailPriority.ProtoReflect.Descriptor instead.
func (*MonorailPriority) Descriptor() ([]byte, []int) {
	return file_infra_appengine_weetbix_internal_config_project_config_proto_rawDescGZIP(), []int{3}
}

func (x *MonorailPriority) GetPriority() string {
	if x != nil {
		return x.Priority
	}
	return ""
}

func (x *MonorailPriority) GetThreshold() *ImpactThreshold {
	if x != nil {
		return x.Threshold
	}
	return nil
}

// ImpactThreshold specifies a condition on a cluster's impact metrics.
// The threshold is considered satisfied if any of the individual metric
// thresholds is satisfied (i.e. if multiple thresholds are set, they are
// combined using an OR-semantic). If no threshold is set on any individual
// metric, the threshold as a whole is unsatisfiable.
type ImpactThreshold struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The minimum number of unexpected failures that should occur in 1 day.
	UnexpectedFailures_1D *int64 `protobuf:"varint,1,opt,name=unexpected_failures_1d,json=unexpectedFailures1d,proto3,oneof" json:"unexpected_failures_1d,omitempty"`
	// The minimum number of unexpected failures that should occur in 3 days.
	UnexpectedFailures_3D *int64 `protobuf:"varint,2,opt,name=unexpected_failures_3d,json=unexpectedFailures3d,proto3,oneof" json:"unexpected_failures_3d,omitempty"`
	// The minimum number of unexpected failures that should occur in 7 days.
	UnexpectedFailures_7D *int64 `protobuf:"varint,3,opt,name=unexpected_failures_7d,json=unexpectedFailures7d,proto3,oneof" json:"unexpected_failures_7d,omitempty"`
}

func (x *ImpactThreshold) Reset() {
	*x = ImpactThreshold{}
	if protoimpl.UnsafeEnabled {
		mi := &file_infra_appengine_weetbix_internal_config_project_config_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ImpactThreshold) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ImpactThreshold) ProtoMessage() {}

func (x *ImpactThreshold) ProtoReflect() protoreflect.Message {
	mi := &file_infra_appengine_weetbix_internal_config_project_config_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ImpactThreshold.ProtoReflect.Descriptor instead.
func (*ImpactThreshold) Descriptor() ([]byte, []int) {
	return file_infra_appengine_weetbix_internal_config_project_config_proto_rawDescGZIP(), []int{4}
}

func (x *ImpactThreshold) GetUnexpectedFailures_1D() int64 {
	if x != nil && x.UnexpectedFailures_1D != nil {
		return *x.UnexpectedFailures_1D
	}
	return 0
}

func (x *ImpactThreshold) GetUnexpectedFailures_3D() int64 {
	if x != nil && x.UnexpectedFailures_3D != nil {
		return *x.UnexpectedFailures_3D
	}
	return 0
}

func (x *ImpactThreshold) GetUnexpectedFailures_7D() int64 {
	if x != nil && x.UnexpectedFailures_7D != nil {
		return *x.UnexpectedFailures_7D
	}
	return 0
}

var File_infra_appengine_weetbix_internal_config_project_config_proto protoreflect.FileDescriptor

var file_infra_appengine_weetbix_internal_config_project_config_proto_rawDesc = []byte{
	0x0a, 0x3c, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x2f, 0x61, 0x70, 0x70, 0x65, 0x6e, 0x67, 0x69, 0x6e,
	0x65, 0x2f, 0x77, 0x65, 0x65, 0x74, 0x62, 0x69, 0x78, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e,
	0x61, 0x6c, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63,
	0x74, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a,
	0x77, 0x65, 0x65, 0x74, 0x62, 0x69, 0x78, 0x2e, 0x76, 0x31, 0x22, 0x97, 0x01, 0x0a, 0x0d, 0x50,
	0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x37, 0x0a, 0x08,
	0x6d, 0x6f, 0x6e, 0x6f, 0x72, 0x61, 0x69, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b,
	0x2e, 0x77, 0x65, 0x65, 0x74, 0x62, 0x69, 0x78, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x6f, 0x6e, 0x6f,
	0x72, 0x61, 0x69, 0x6c, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x52, 0x08, 0x6d, 0x6f, 0x6e,
	0x6f, 0x72, 0x61, 0x69, 0x6c, 0x12, 0x4d, 0x0a, 0x14, 0x62, 0x75, 0x67, 0x5f, 0x66, 0x69, 0x6c,
	0x69, 0x6e, 0x67, 0x5f, 0x74, 0x68, 0x72, 0x65, 0x73, 0x68, 0x6f, 0x6c, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x77, 0x65, 0x65, 0x74, 0x62, 0x69, 0x78, 0x2e, 0x76, 0x31,
	0x2e, 0x49, 0x6d, 0x70, 0x61, 0x63, 0x74, 0x54, 0x68, 0x72, 0x65, 0x73, 0x68, 0x6f, 0x6c, 0x64,
	0x52, 0x12, 0x62, 0x75, 0x67, 0x46, 0x69, 0x6c, 0x69, 0x6e, 0x67, 0x54, 0x68, 0x72, 0x65, 0x73,
	0x68, 0x6f, 0x6c, 0x64, 0x22, 0xa7, 0x02, 0x0a, 0x0f, 0x4d, 0x6f, 0x6e, 0x6f, 0x72, 0x61, 0x69,
	0x6c, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x72, 0x6f, 0x6a,
	0x65, 0x63, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x70, 0x72, 0x6f, 0x6a, 0x65,
	0x63, 0x74, 0x12, 0x50, 0x0a, 0x14, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x5f, 0x66, 0x69,
	0x65, 0x6c, 0x64, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x1e, 0x2e, 0x77, 0x65, 0x65, 0x74, 0x62, 0x69, 0x78, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x6f,
	0x6e, 0x6f, 0x72, 0x61, 0x69, 0x6c, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x52, 0x12, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x73, 0x12, 0x2a, 0x0a, 0x11, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79,
	0x5f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x0f, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x49, 0x64,
	0x12, 0x3c, 0x0a, 0x0a, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x69, 0x65, 0x73, 0x18, 0x04,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x77, 0x65, 0x65, 0x74, 0x62, 0x69, 0x78, 0x2e, 0x76,
	0x31, 0x2e, 0x4d, 0x6f, 0x6e, 0x6f, 0x72, 0x61, 0x69, 0x6c, 0x50, 0x72, 0x69, 0x6f, 0x72, 0x69,
	0x74, 0x79, 0x52, 0x0a, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x69, 0x65, 0x73, 0x12, 0x3e,
	0x0a, 0x1b, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x5f, 0x68, 0x79, 0x73, 0x74, 0x65,
	0x72, 0x65, 0x73, 0x69, 0x73, 0x5f, 0x70, 0x65, 0x72, 0x63, 0x65, 0x6e, 0x74, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x19, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x48, 0x79, 0x73,
	0x74, 0x65, 0x72, 0x65, 0x73, 0x69, 0x73, 0x50, 0x65, 0x72, 0x63, 0x65, 0x6e, 0x74, 0x22, 0x45,
	0x0a, 0x12, 0x4d, 0x6f, 0x6e, 0x6f, 0x72, 0x61, 0x69, 0x6c, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x49, 0x64, 0x12,
	0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x69, 0x0a, 0x10, 0x4d, 0x6f, 0x6e, 0x6f, 0x72, 0x61, 0x69,
	0x6c, 0x50, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x72, 0x69,
	0x6f, 0x72, 0x69, 0x74, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x72, 0x69,
	0x6f, 0x72, 0x69, 0x74, 0x79, 0x12, 0x39, 0x0a, 0x09, 0x74, 0x68, 0x72, 0x65, 0x73, 0x68, 0x6f,
	0x6c, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x77, 0x65, 0x65, 0x74, 0x62,
	0x69, 0x78, 0x2e, 0x76, 0x31, 0x2e, 0x49, 0x6d, 0x70, 0x61, 0x63, 0x74, 0x54, 0x68, 0x72, 0x65,
	0x73, 0x68, 0x6f, 0x6c, 0x64, 0x52, 0x09, 0x74, 0x68, 0x72, 0x65, 0x73, 0x68, 0x6f, 0x6c, 0x64,
	0x22, 0x93, 0x02, 0x0a, 0x0f, 0x49, 0x6d, 0x70, 0x61, 0x63, 0x74, 0x54, 0x68, 0x72, 0x65, 0x73,
	0x68, 0x6f, 0x6c, 0x64, 0x12, 0x39, 0x0a, 0x16, 0x75, 0x6e, 0x65, 0x78, 0x70, 0x65, 0x63, 0x74,
	0x65, 0x64, 0x5f, 0x66, 0x61, 0x69, 0x6c, 0x75, 0x72, 0x65, 0x73, 0x5f, 0x31, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x03, 0x48, 0x00, 0x52, 0x14, 0x75, 0x6e, 0x65, 0x78, 0x70, 0x65, 0x63, 0x74,
	0x65, 0x64, 0x46, 0x61, 0x69, 0x6c, 0x75, 0x72, 0x65, 0x73, 0x31, 0x64, 0x88, 0x01, 0x01, 0x12,
	0x39, 0x0a, 0x16, 0x75, 0x6e, 0x65, 0x78, 0x70, 0x65, 0x63, 0x74, 0x65, 0x64, 0x5f, 0x66, 0x61,
	0x69, 0x6c, 0x75, 0x72, 0x65, 0x73, 0x5f, 0x33, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x48,
	0x01, 0x52, 0x14, 0x75, 0x6e, 0x65, 0x78, 0x70, 0x65, 0x63, 0x74, 0x65, 0x64, 0x46, 0x61, 0x69,
	0x6c, 0x75, 0x72, 0x65, 0x73, 0x33, 0x64, 0x88, 0x01, 0x01, 0x12, 0x39, 0x0a, 0x16, 0x75, 0x6e,
	0x65, 0x78, 0x70, 0x65, 0x63, 0x74, 0x65, 0x64, 0x5f, 0x66, 0x61, 0x69, 0x6c, 0x75, 0x72, 0x65,
	0x73, 0x5f, 0x37, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x48, 0x02, 0x52, 0x14, 0x75, 0x6e,
	0x65, 0x78, 0x70, 0x65, 0x63, 0x74, 0x65, 0x64, 0x46, 0x61, 0x69, 0x6c, 0x75, 0x72, 0x65, 0x73,
	0x37, 0x64, 0x88, 0x01, 0x01, 0x42, 0x19, 0x0a, 0x17, 0x5f, 0x75, 0x6e, 0x65, 0x78, 0x70, 0x65,
	0x63, 0x74, 0x65, 0x64, 0x5f, 0x66, 0x61, 0x69, 0x6c, 0x75, 0x72, 0x65, 0x73, 0x5f, 0x31, 0x64,
	0x42, 0x19, 0x0a, 0x17, 0x5f, 0x75, 0x6e, 0x65, 0x78, 0x70, 0x65, 0x63, 0x74, 0x65, 0x64, 0x5f,
	0x66, 0x61, 0x69, 0x6c, 0x75, 0x72, 0x65, 0x73, 0x5f, 0x33, 0x64, 0x42, 0x19, 0x0a, 0x17, 0x5f,
	0x75, 0x6e, 0x65, 0x78, 0x70, 0x65, 0x63, 0x74, 0x65, 0x64, 0x5f, 0x66, 0x61, 0x69, 0x6c, 0x75,
	0x72, 0x65, 0x73, 0x5f, 0x37, 0x64, 0x42, 0x30, 0x5a, 0x2e, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x2f,
	0x61, 0x70, 0x70, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2f, 0x77, 0x65, 0x65, 0x74, 0x62, 0x69,
	0x78, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x3b, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_infra_appengine_weetbix_internal_config_project_config_proto_rawDescOnce sync.Once
	file_infra_appengine_weetbix_internal_config_project_config_proto_rawDescData = file_infra_appengine_weetbix_internal_config_project_config_proto_rawDesc
)

func file_infra_appengine_weetbix_internal_config_project_config_proto_rawDescGZIP() []byte {
	file_infra_appengine_weetbix_internal_config_project_config_proto_rawDescOnce.Do(func() {
		file_infra_appengine_weetbix_internal_config_project_config_proto_rawDescData = protoimpl.X.CompressGZIP(file_infra_appengine_weetbix_internal_config_project_config_proto_rawDescData)
	})
	return file_infra_appengine_weetbix_internal_config_project_config_proto_rawDescData
}

var file_infra_appengine_weetbix_internal_config_project_config_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_infra_appengine_weetbix_internal_config_project_config_proto_goTypes = []interface{}{
	(*ProjectConfig)(nil),      // 0: weetbix.v1.ProjectConfig
	(*MonorailProject)(nil),    // 1: weetbix.v1.MonorailProject
	(*MonorailFieldValue)(nil), // 2: weetbix.v1.MonorailFieldValue
	(*MonorailPriority)(nil),   // 3: weetbix.v1.MonorailPriority
	(*ImpactThreshold)(nil),    // 4: weetbix.v1.ImpactThreshold
}
var file_infra_appengine_weetbix_internal_config_project_config_proto_depIdxs = []int32{
	1, // 0: weetbix.v1.ProjectConfig.monorail:type_name -> weetbix.v1.MonorailProject
	4, // 1: weetbix.v1.ProjectConfig.bug_filing_threshold:type_name -> weetbix.v1.ImpactThreshold
	2, // 2: weetbix.v1.MonorailProject.default_field_values:type_name -> weetbix.v1.MonorailFieldValue
	3, // 3: weetbix.v1.MonorailProject.priorities:type_name -> weetbix.v1.MonorailPriority
	4, // 4: weetbix.v1.MonorailPriority.threshold:type_name -> weetbix.v1.ImpactThreshold
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_infra_appengine_weetbix_internal_config_project_config_proto_init() }
func file_infra_appengine_weetbix_internal_config_project_config_proto_init() {
	if File_infra_appengine_weetbix_internal_config_project_config_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_infra_appengine_weetbix_internal_config_project_config_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProjectConfig); i {
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
		file_infra_appengine_weetbix_internal_config_project_config_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MonorailProject); i {
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
		file_infra_appengine_weetbix_internal_config_project_config_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MonorailFieldValue); i {
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
		file_infra_appengine_weetbix_internal_config_project_config_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MonorailPriority); i {
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
		file_infra_appengine_weetbix_internal_config_project_config_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ImpactThreshold); i {
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
	file_infra_appengine_weetbix_internal_config_project_config_proto_msgTypes[4].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_infra_appengine_weetbix_internal_config_project_config_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_infra_appengine_weetbix_internal_config_project_config_proto_goTypes,
		DependencyIndexes: file_infra_appengine_weetbix_internal_config_project_config_proto_depIdxs,
		MessageInfos:      file_infra_appengine_weetbix_internal_config_project_config_proto_msgTypes,
	}.Build()
	File_infra_appengine_weetbix_internal_config_project_config_proto = out.File
	file_infra_appengine_weetbix_internal_config_project_config_proto_rawDesc = nil
	file_infra_appengine_weetbix_internal_config_project_config_proto_goTypes = nil
	file_infra_appengine_weetbix_internal_config_project_config_proto_depIdxs = nil
}
