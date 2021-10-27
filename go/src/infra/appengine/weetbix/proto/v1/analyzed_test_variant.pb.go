// Copyright 2021 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: infra/appengine/weetbix/proto/v1/analyzed_test_variant.proto

package weetbixpb

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

// Status of a test variant.
type AnalyzedTestVariantStatus int32

const (
	// Status was not specified.
	// Not to be used in actual test variants; serves as a default value for an unset field.
	AnalyzedTestVariantStatus_STATUS_UNSPECIFIED AnalyzedTestVariantStatus = 0
	// The test variant has unexpected results, but Weetbix cannot determine
	// If it is FLAKY or CONSISTENTLY_UNEXPECTED.
	// This status can be used when
	// * in in-build flakiness cases, a test variant with flaky results in a build
	//   is newly detected but the service has not been notified if the build
	//   contributes to a CL's submission or not.
	//   *  Note that this does not apply to Chromium flaky analysis because for
	//      Chromium Weetbix only ingests test results from builds contribute to
	//      CL submissions.
	// * in cross-build flakiness cases, a test variant is newly detected in a build
	//   where all of its results are unexpected.
	AnalyzedTestVariantStatus_HAS_UNEXPECTED_RESULTS AnalyzedTestVariantStatus = 5
	// The test variant is currently flaky.
	AnalyzedTestVariantStatus_FLAKY AnalyzedTestVariantStatus = 10
	// Results of the test variant have been consistently unexpected for
	// a period of time.
	AnalyzedTestVariantStatus_CONSISTENTLY_UNEXPECTED AnalyzedTestVariantStatus = 20
	// Results of the test variant have been consistently expected for
	// a period of time.
	// TODO(chanli@): mention the configuration that specifies the time range.
	AnalyzedTestVariantStatus_CONSISTENTLY_EXPECTED AnalyzedTestVariantStatus = 30
	// There are no new results of the test variant for a period of time.
	// It's likely that this test variant has been disabled or removed.
	AnalyzedTestVariantStatus_NO_NEW_RESULTS AnalyzedTestVariantStatus = 40
)

// Enum value maps for AnalyzedTestVariantStatus.
var (
	AnalyzedTestVariantStatus_name = map[int32]string{
		0:  "STATUS_UNSPECIFIED",
		5:  "HAS_UNEXPECTED_RESULTS",
		10: "FLAKY",
		20: "CONSISTENTLY_UNEXPECTED",
		30: "CONSISTENTLY_EXPECTED",
		40: "NO_NEW_RESULTS",
	}
	AnalyzedTestVariantStatus_value = map[string]int32{
		"STATUS_UNSPECIFIED":      0,
		"HAS_UNEXPECTED_RESULTS":  5,
		"FLAKY":                   10,
		"CONSISTENTLY_UNEXPECTED": 20,
		"CONSISTENTLY_EXPECTED":   30,
		"NO_NEW_RESULTS":          40,
	}
)

func (x AnalyzedTestVariantStatus) Enum() *AnalyzedTestVariantStatus {
	p := new(AnalyzedTestVariantStatus)
	*p = x
	return p
}

func (x AnalyzedTestVariantStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (AnalyzedTestVariantStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_infra_appengine_weetbix_proto_v1_analyzed_test_variant_proto_enumTypes[0].Descriptor()
}

func (AnalyzedTestVariantStatus) Type() protoreflect.EnumType {
	return &file_infra_appengine_weetbix_proto_v1_analyzed_test_variant_proto_enumTypes[0]
}

func (x AnalyzedTestVariantStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use AnalyzedTestVariantStatus.Descriptor instead.
func (AnalyzedTestVariantStatus) EnumDescriptor() ([]byte, []int) {
	return file_infra_appengine_weetbix_proto_v1_analyzed_test_variant_proto_rawDescGZIP(), []int{0}
}

// Information about a test.
//
// As of Oct 2021, it's an exact copy of luci.resultdb.v1.TestMetadata, but
// we'd like to keep a local definition of the proto to keep the possibility that
// we need to diverge down the track.
type TestMetadata struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The original test name.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// Where the test is defined, e.g. the file name.
	// location.repo MUST be specified.
	Location *TestLocation `protobuf:"bytes,2,opt,name=location,proto3" json:"location,omitempty"`
}

func (x *TestMetadata) Reset() {
	*x = TestMetadata{}
	if protoimpl.UnsafeEnabled {
		mi := &file_infra_appengine_weetbix_proto_v1_analyzed_test_variant_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TestMetadata) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TestMetadata) ProtoMessage() {}

func (x *TestMetadata) ProtoReflect() protoreflect.Message {
	mi := &file_infra_appengine_weetbix_proto_v1_analyzed_test_variant_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TestMetadata.ProtoReflect.Descriptor instead.
func (*TestMetadata) Descriptor() ([]byte, []int) {
	return file_infra_appengine_weetbix_proto_v1_analyzed_test_variant_proto_rawDescGZIP(), []int{0}
}

func (x *TestMetadata) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *TestMetadata) GetLocation() *TestLocation {
	if x != nil {
		return x.Location
	}
	return nil
}

// Location of the test definition.
type TestLocation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Gitiles URL as the identifier for a repo.
	// Format for Gitiles URL: https://<host>/<project>
	// For example "https://chromium.googlesource.com/chromium/src"
	// Must not end with ".git".
	// SHOULD be specified.
	Repo string `protobuf:"bytes,1,opt,name=repo,proto3" json:"repo,omitempty"`
	// Name of the file where the test is defined.
	// For files in a repository, must start with "//"
	// Example: "//components/payments/core/payment_request_data_util_unittest.cc"
	// Max length: 512.
	// MUST not use backslashes.
	// Required.
	FileName string `protobuf:"bytes,2,opt,name=file_name,json=fileName,proto3" json:"file_name,omitempty"`
	// One-based line number where the test is defined.
	Line int32 `protobuf:"varint,3,opt,name=line,proto3" json:"line,omitempty"`
}

func (x *TestLocation) Reset() {
	*x = TestLocation{}
	if protoimpl.UnsafeEnabled {
		mi := &file_infra_appengine_weetbix_proto_v1_analyzed_test_variant_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TestLocation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TestLocation) ProtoMessage() {}

func (x *TestLocation) ProtoReflect() protoreflect.Message {
	mi := &file_infra_appengine_weetbix_proto_v1_analyzed_test_variant_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TestLocation.ProtoReflect.Descriptor instead.
func (*TestLocation) Descriptor() ([]byte, []int) {
	return file_infra_appengine_weetbix_proto_v1_analyzed_test_variant_proto_rawDescGZIP(), []int{1}
}

func (x *TestLocation) GetRepo() string {
	if x != nil {
		return x.Repo
	}
	return ""
}

func (x *TestLocation) GetFileName() string {
	if x != nil {
		return x.FileName
	}
	return ""
}

func (x *TestLocation) GetLine() int32 {
	if x != nil {
		return x.Line
	}
	return 0
}

// Flake statistics of a test variant.
type FlakeStatistics struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Flake verdict rate calculated by the verdicts within the time range.
	FlakyVerdictRate float32 `protobuf:"fixed32,1,opt,name=flaky_verdict_rate,json=flakyVerdictRate,proto3" json:"flaky_verdict_rate,omitempty"`
	// Count of verdicts with flaky status.
	FlakyVerdictCount int32 `protobuf:"varint,2,opt,name=flaky_verdict_count,json=flakyVerdictCount,proto3" json:"flaky_verdict_count,omitempty"`
	// Count of total verdicts.
	TotalVerdictCount int32 `protobuf:"varint,3,opt,name=total_verdict_count,json=totalVerdictCount,proto3" json:"total_verdict_count,omitempty"`
	// Unexpected result rate calculated by the test results within the time range.
	UnexpectedResultRate float32 `protobuf:"fixed32,4,opt,name=unexpected_result_rate,json=unexpectedResultRate,proto3" json:"unexpected_result_rate,omitempty"`
	// Count of unexpected results.
	UnexpectedResultCount int32 `protobuf:"varint,5,opt,name=unexpected_result_count,json=unexpectedResultCount,proto3" json:"unexpected_result_count,omitempty"`
	// Count of total results.
	TotalResultCount int32 `protobuf:"varint,6,opt,name=total_result_count,json=totalResultCount,proto3" json:"total_result_count,omitempty"`
}

func (x *FlakeStatistics) Reset() {
	*x = FlakeStatistics{}
	if protoimpl.UnsafeEnabled {
		mi := &file_infra_appengine_weetbix_proto_v1_analyzed_test_variant_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FlakeStatistics) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FlakeStatistics) ProtoMessage() {}

func (x *FlakeStatistics) ProtoReflect() protoreflect.Message {
	mi := &file_infra_appengine_weetbix_proto_v1_analyzed_test_variant_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FlakeStatistics.ProtoReflect.Descriptor instead.
func (*FlakeStatistics) Descriptor() ([]byte, []int) {
	return file_infra_appengine_weetbix_proto_v1_analyzed_test_variant_proto_rawDescGZIP(), []int{2}
}

func (x *FlakeStatistics) GetFlakyVerdictRate() float32 {
	if x != nil {
		return x.FlakyVerdictRate
	}
	return 0
}

func (x *FlakeStatistics) GetFlakyVerdictCount() int32 {
	if x != nil {
		return x.FlakyVerdictCount
	}
	return 0
}

func (x *FlakeStatistics) GetTotalVerdictCount() int32 {
	if x != nil {
		return x.TotalVerdictCount
	}
	return 0
}

func (x *FlakeStatistics) GetUnexpectedResultRate() float32 {
	if x != nil {
		return x.UnexpectedResultRate
	}
	return 0
}

func (x *FlakeStatistics) GetUnexpectedResultCount() int32 {
	if x != nil {
		return x.UnexpectedResultCount
	}
	return 0
}

func (x *FlakeStatistics) GetTotalResultCount() int32 {
	if x != nil {
		return x.TotalResultCount
	}
	return 0
}

type AnalyzedTestVariant struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Can be used to refer to this test variant.
	// Format:
	// "realms/{REALM}/tests/{URL_ESCAPED_TEST_ID}/variants/{VARIANT_HASH}"
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// Realm that the test variant exists under.
	// See https://source.chromium.org/chromium/infra/infra/+/main:go/src/go.chromium.org/luci/common/proto/realms/realms_config.proto
	Realm string `protobuf:"bytes,2,opt,name=realm,proto3" json:"realm,omitempty"`
	// Test id, identifier of the test. Unique in a LUCI realm.
	TestId string `protobuf:"bytes,3,opt,name=test_id,json=testId,proto3" json:"test_id,omitempty"`
	// Hash of the variant.
	VariantHash string `protobuf:"bytes,4,opt,name=variant_hash,json=variantHash,proto3" json:"variant_hash,omitempty"`
	// Description of one specific way of running the test,
	// e.g. a specific bucket, builder and a test suite.
	Variant *Variant `protobuf:"bytes,5,opt,name=variant,proto3" json:"variant,omitempty"`
	// Information about the test at the time of its execution.
	TestMetadata *TestMetadata `protobuf:"bytes,6,opt,name=test_metadata,json=testMetadata,proto3" json:"test_metadata,omitempty"`
	// Metadata for the test variant.
	// See luci.resultdb.v1.Tags for details.
	Tags []*StringPair `protobuf:"bytes,7,rep,name=tags,proto3" json:"tags,omitempty"`
	// A range of time. Flake statistics are calculated using test results
	// within that range.
	TimeRange *TimeRange `protobuf:"bytes,8,opt,name=time_range,json=timeRange,proto3" json:"time_range,omitempty"`
	// Status of the test valiant.
	Status AnalyzedTestVariantStatus `protobuf:"varint,9,opt,name=status,proto3,enum=weetbix.v1.AnalyzedTestVariantStatus" json:"status,omitempty"`
	// Flakiness statistics of the test variant.
	FlakeStatistics *FlakeStatistics `protobuf:"bytes,10,opt,name=flake_statistics,json=flakeStatistics,proto3" json:"flake_statistics,omitempty"`
}

func (x *AnalyzedTestVariant) Reset() {
	*x = AnalyzedTestVariant{}
	if protoimpl.UnsafeEnabled {
		mi := &file_infra_appengine_weetbix_proto_v1_analyzed_test_variant_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AnalyzedTestVariant) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AnalyzedTestVariant) ProtoMessage() {}

func (x *AnalyzedTestVariant) ProtoReflect() protoreflect.Message {
	mi := &file_infra_appengine_weetbix_proto_v1_analyzed_test_variant_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AnalyzedTestVariant.ProtoReflect.Descriptor instead.
func (*AnalyzedTestVariant) Descriptor() ([]byte, []int) {
	return file_infra_appengine_weetbix_proto_v1_analyzed_test_variant_proto_rawDescGZIP(), []int{3}
}

func (x *AnalyzedTestVariant) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *AnalyzedTestVariant) GetRealm() string {
	if x != nil {
		return x.Realm
	}
	return ""
}

func (x *AnalyzedTestVariant) GetTestId() string {
	if x != nil {
		return x.TestId
	}
	return ""
}

func (x *AnalyzedTestVariant) GetVariantHash() string {
	if x != nil {
		return x.VariantHash
	}
	return ""
}

func (x *AnalyzedTestVariant) GetVariant() *Variant {
	if x != nil {
		return x.Variant
	}
	return nil
}

func (x *AnalyzedTestVariant) GetTestMetadata() *TestMetadata {
	if x != nil {
		return x.TestMetadata
	}
	return nil
}

func (x *AnalyzedTestVariant) GetTags() []*StringPair {
	if x != nil {
		return x.Tags
	}
	return nil
}

func (x *AnalyzedTestVariant) GetTimeRange() *TimeRange {
	if x != nil {
		return x.TimeRange
	}
	return nil
}

func (x *AnalyzedTestVariant) GetStatus() AnalyzedTestVariantStatus {
	if x != nil {
		return x.Status
	}
	return AnalyzedTestVariantStatus_STATUS_UNSPECIFIED
}

func (x *AnalyzedTestVariant) GetFlakeStatistics() *FlakeStatistics {
	if x != nil {
		return x.FlakeStatistics
	}
	return nil
}

var File_infra_appengine_weetbix_proto_v1_analyzed_test_variant_proto protoreflect.FileDescriptor

var file_infra_appengine_weetbix_proto_v1_analyzed_test_variant_proto_rawDesc = []byte{
	0x0a, 0x3c, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x2f, 0x61, 0x70, 0x70, 0x65, 0x6e, 0x67, 0x69, 0x6e,
	0x65, 0x2f, 0x77, 0x65, 0x65, 0x74, 0x62, 0x69, 0x78, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x76, 0x31, 0x2f, 0x61, 0x6e, 0x61, 0x6c, 0x79, 0x7a, 0x65, 0x64, 0x5f, 0x74, 0x65, 0x73, 0x74,
	0x5f, 0x76, 0x61, 0x72, 0x69, 0x61, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a,
	0x77, 0x65, 0x65, 0x74, 0x62, 0x69, 0x78, 0x2e, 0x76, 0x31, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x5f, 0x62, 0x65, 0x68,
	0x61, 0x76, 0x69, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2d, 0x69, 0x6e, 0x66,
	0x72, 0x61, 0x2f, 0x61, 0x70, 0x70, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x2f, 0x77, 0x65, 0x65,
	0x74, 0x62, 0x69, 0x78, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x58, 0x0a, 0x0c, 0x54, 0x65,
	0x73, 0x74, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x34,
	0x0a, 0x08, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x18, 0x2e, 0x77, 0x65, 0x65, 0x74, 0x62, 0x69, 0x78, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x65,
	0x73, 0x74, 0x4c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x08, 0x6c, 0x6f, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x22, 0x53, 0x0a, 0x0c, 0x54, 0x65, 0x73, 0x74, 0x4c, 0x6f, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x72, 0x65, 0x70, 0x6f, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x72, 0x65, 0x70, 0x6f, 0x12, 0x1b, 0x0a, 0x09, 0x66, 0x69, 0x6c, 0x65,
	0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c,
	0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6c, 0x69, 0x6e, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x04, 0x6c, 0x69, 0x6e, 0x65, 0x22, 0xbb, 0x02, 0x0a, 0x0f, 0x46, 0x6c,
	0x61, 0x6b, 0x65, 0x53, 0x74, 0x61, 0x74, 0x69, 0x73, 0x74, 0x69, 0x63, 0x73, 0x12, 0x2c, 0x0a,
	0x12, 0x66, 0x6c, 0x61, 0x6b, 0x79, 0x5f, 0x76, 0x65, 0x72, 0x64, 0x69, 0x63, 0x74, 0x5f, 0x72,
	0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x02, 0x52, 0x10, 0x66, 0x6c, 0x61, 0x6b, 0x79,
	0x56, 0x65, 0x72, 0x64, 0x69, 0x63, 0x74, 0x52, 0x61, 0x74, 0x65, 0x12, 0x2e, 0x0a, 0x13, 0x66,
	0x6c, 0x61, 0x6b, 0x79, 0x5f, 0x76, 0x65, 0x72, 0x64, 0x69, 0x63, 0x74, 0x5f, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x11, 0x66, 0x6c, 0x61, 0x6b, 0x79, 0x56,
	0x65, 0x72, 0x64, 0x69, 0x63, 0x74, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x2e, 0x0a, 0x13, 0x74,
	0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x76, 0x65, 0x72, 0x64, 0x69, 0x63, 0x74, 0x5f, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x11, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x56,
	0x65, 0x72, 0x64, 0x69, 0x63, 0x74, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x34, 0x0a, 0x16, 0x75,
	0x6e, 0x65, 0x78, 0x70, 0x65, 0x63, 0x74, 0x65, 0x64, 0x5f, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74,
	0x5f, 0x72, 0x61, 0x74, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x02, 0x52, 0x14, 0x75, 0x6e, 0x65,
	0x78, 0x70, 0x65, 0x63, 0x74, 0x65, 0x64, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x52, 0x61, 0x74,
	0x65, 0x12, 0x36, 0x0a, 0x17, 0x75, 0x6e, 0x65, 0x78, 0x70, 0x65, 0x63, 0x74, 0x65, 0x64, 0x5f,
	0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x15, 0x75, 0x6e, 0x65, 0x78, 0x70, 0x65, 0x63, 0x74, 0x65, 0x64, 0x52, 0x65,
	0x73, 0x75, 0x6c, 0x74, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x2c, 0x0a, 0x12, 0x74, 0x6f, 0x74,
	0x61, 0x6c, 0x5f, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x05, 0x52, 0x10, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x52, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0xda, 0x03, 0x0a, 0x13, 0x41, 0x6e, 0x61, 0x6c,
	0x79, 0x7a, 0x65, 0x64, 0x54, 0x65, 0x73, 0x74, 0x56, 0x61, 0x72, 0x69, 0x61, 0x6e, 0x74, 0x12,
	0x1a, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x06, 0xe0,
	0x41, 0x03, 0xe0, 0x41, 0x05, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x72,
	0x65, 0x61, 0x6c, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x72, 0x65, 0x61, 0x6c,
	0x6d, 0x12, 0x17, 0x0a, 0x07, 0x74, 0x65, 0x73, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x74, 0x65, 0x73, 0x74, 0x49, 0x64, 0x12, 0x21, 0x0a, 0x0c, 0x76, 0x61,
	0x72, 0x69, 0x61, 0x6e, 0x74, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0b, 0x76, 0x61, 0x72, 0x69, 0x61, 0x6e, 0x74, 0x48, 0x61, 0x73, 0x68, 0x12, 0x2d, 0x0a,
	0x07, 0x76, 0x61, 0x72, 0x69, 0x61, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13,
	0x2e, 0x77, 0x65, 0x65, 0x74, 0x62, 0x69, 0x78, 0x2e, 0x76, 0x31, 0x2e, 0x56, 0x61, 0x72, 0x69,
	0x61, 0x6e, 0x74, 0x52, 0x07, 0x76, 0x61, 0x72, 0x69, 0x61, 0x6e, 0x74, 0x12, 0x3d, 0x0a, 0x0d,
	0x74, 0x65, 0x73, 0x74, 0x5f, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x77, 0x65, 0x65, 0x74, 0x62, 0x69, 0x78, 0x2e, 0x76, 0x31,
	0x2e, 0x54, 0x65, 0x73, 0x74, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x0c, 0x74,
	0x65, 0x73, 0x74, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x2a, 0x0a, 0x04, 0x74,
	0x61, 0x67, 0x73, 0x18, 0x07, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x77, 0x65, 0x65, 0x74,
	0x62, 0x69, 0x78, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x50, 0x61, 0x69,
	0x72, 0x52, 0x04, 0x74, 0x61, 0x67, 0x73, 0x12, 0x34, 0x0a, 0x0a, 0x74, 0x69, 0x6d, 0x65, 0x5f,
	0x72, 0x61, 0x6e, 0x67, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x77, 0x65,
	0x65, 0x74, 0x62, 0x69, 0x78, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x52, 0x61, 0x6e,
	0x67, 0x65, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x12, 0x3d, 0x0a,
	0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x25, 0x2e,
	0x77, 0x65, 0x65, 0x74, 0x62, 0x69, 0x78, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x6e, 0x61, 0x6c, 0x79,
	0x7a, 0x65, 0x64, 0x54, 0x65, 0x73, 0x74, 0x56, 0x61, 0x72, 0x69, 0x61, 0x6e, 0x74, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x46, 0x0a, 0x10,
	0x66, 0x6c, 0x61, 0x6b, 0x65, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x69, 0x73, 0x74, 0x69, 0x63, 0x73,
	0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x77, 0x65, 0x65, 0x74, 0x62, 0x69, 0x78,
	0x2e, 0x76, 0x31, 0x2e, 0x46, 0x6c, 0x61, 0x6b, 0x65, 0x53, 0x74, 0x61, 0x74, 0x69, 0x73, 0x74,
	0x69, 0x63, 0x73, 0x52, 0x0f, 0x66, 0x6c, 0x61, 0x6b, 0x65, 0x53, 0x74, 0x61, 0x74, 0x69, 0x73,
	0x74, 0x69, 0x63, 0x73, 0x2a, 0xa6, 0x01, 0x0a, 0x19, 0x41, 0x6e, 0x61, 0x6c, 0x79, 0x7a, 0x65,
	0x64, 0x54, 0x65, 0x73, 0x74, 0x56, 0x61, 0x72, 0x69, 0x61, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x12, 0x16, 0x0a, 0x12, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x55, 0x4e, 0x53,
	0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x1a, 0x0a, 0x16, 0x48, 0x41,
	0x53, 0x5f, 0x55, 0x4e, 0x45, 0x58, 0x50, 0x45, 0x43, 0x54, 0x45, 0x44, 0x5f, 0x52, 0x45, 0x53,
	0x55, 0x4c, 0x54, 0x53, 0x10, 0x05, 0x12, 0x09, 0x0a, 0x05, 0x46, 0x4c, 0x41, 0x4b, 0x59, 0x10,
	0x0a, 0x12, 0x1b, 0x0a, 0x17, 0x43, 0x4f, 0x4e, 0x53, 0x49, 0x53, 0x54, 0x45, 0x4e, 0x54, 0x4c,
	0x59, 0x5f, 0x55, 0x4e, 0x45, 0x58, 0x50, 0x45, 0x43, 0x54, 0x45, 0x44, 0x10, 0x14, 0x12, 0x19,
	0x0a, 0x15, 0x43, 0x4f, 0x4e, 0x53, 0x49, 0x53, 0x54, 0x45, 0x4e, 0x54, 0x4c, 0x59, 0x5f, 0x45,
	0x58, 0x50, 0x45, 0x43, 0x54, 0x45, 0x44, 0x10, 0x1e, 0x12, 0x12, 0x0a, 0x0e, 0x4e, 0x4f, 0x5f,
	0x4e, 0x45, 0x57, 0x5f, 0x52, 0x45, 0x53, 0x55, 0x4c, 0x54, 0x53, 0x10, 0x28, 0x42, 0x2c, 0x5a,
	0x2a, 0x69, 0x6e, 0x66, 0x72, 0x61, 0x2f, 0x61, 0x70, 0x70, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65,
	0x2f, 0x77, 0x65, 0x65, 0x74, 0x62, 0x69, 0x78, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x76,
	0x31, 0x3b, 0x77, 0x65, 0x65, 0x74, 0x62, 0x69, 0x78, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_infra_appengine_weetbix_proto_v1_analyzed_test_variant_proto_rawDescOnce sync.Once
	file_infra_appengine_weetbix_proto_v1_analyzed_test_variant_proto_rawDescData = file_infra_appengine_weetbix_proto_v1_analyzed_test_variant_proto_rawDesc
)

func file_infra_appengine_weetbix_proto_v1_analyzed_test_variant_proto_rawDescGZIP() []byte {
	file_infra_appengine_weetbix_proto_v1_analyzed_test_variant_proto_rawDescOnce.Do(func() {
		file_infra_appengine_weetbix_proto_v1_analyzed_test_variant_proto_rawDescData = protoimpl.X.CompressGZIP(file_infra_appengine_weetbix_proto_v1_analyzed_test_variant_proto_rawDescData)
	})
	return file_infra_appengine_weetbix_proto_v1_analyzed_test_variant_proto_rawDescData
}

var file_infra_appengine_weetbix_proto_v1_analyzed_test_variant_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_infra_appengine_weetbix_proto_v1_analyzed_test_variant_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_infra_appengine_weetbix_proto_v1_analyzed_test_variant_proto_goTypes = []interface{}{
	(AnalyzedTestVariantStatus)(0), // 0: weetbix.v1.AnalyzedTestVariantStatus
	(*TestMetadata)(nil),           // 1: weetbix.v1.TestMetadata
	(*TestLocation)(nil),           // 2: weetbix.v1.TestLocation
	(*FlakeStatistics)(nil),        // 3: weetbix.v1.FlakeStatistics
	(*AnalyzedTestVariant)(nil),    // 4: weetbix.v1.AnalyzedTestVariant
	(*Variant)(nil),                // 5: weetbix.v1.Variant
	(*StringPair)(nil),             // 6: weetbix.v1.StringPair
	(*TimeRange)(nil),              // 7: weetbix.v1.TimeRange
}
var file_infra_appengine_weetbix_proto_v1_analyzed_test_variant_proto_depIdxs = []int32{
	2, // 0: weetbix.v1.TestMetadata.location:type_name -> weetbix.v1.TestLocation
	5, // 1: weetbix.v1.AnalyzedTestVariant.variant:type_name -> weetbix.v1.Variant
	1, // 2: weetbix.v1.AnalyzedTestVariant.test_metadata:type_name -> weetbix.v1.TestMetadata
	6, // 3: weetbix.v1.AnalyzedTestVariant.tags:type_name -> weetbix.v1.StringPair
	7, // 4: weetbix.v1.AnalyzedTestVariant.time_range:type_name -> weetbix.v1.TimeRange
	0, // 5: weetbix.v1.AnalyzedTestVariant.status:type_name -> weetbix.v1.AnalyzedTestVariantStatus
	3, // 6: weetbix.v1.AnalyzedTestVariant.flake_statistics:type_name -> weetbix.v1.FlakeStatistics
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_infra_appengine_weetbix_proto_v1_analyzed_test_variant_proto_init() }
func file_infra_appengine_weetbix_proto_v1_analyzed_test_variant_proto_init() {
	if File_infra_appengine_weetbix_proto_v1_analyzed_test_variant_proto != nil {
		return
	}
	file_infra_appengine_weetbix_proto_v1_common_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_infra_appengine_weetbix_proto_v1_analyzed_test_variant_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TestMetadata); i {
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
		file_infra_appengine_weetbix_proto_v1_analyzed_test_variant_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TestLocation); i {
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
		file_infra_appengine_weetbix_proto_v1_analyzed_test_variant_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FlakeStatistics); i {
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
		file_infra_appengine_weetbix_proto_v1_analyzed_test_variant_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AnalyzedTestVariant); i {
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
			RawDescriptor: file_infra_appengine_weetbix_proto_v1_analyzed_test_variant_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_infra_appengine_weetbix_proto_v1_analyzed_test_variant_proto_goTypes,
		DependencyIndexes: file_infra_appengine_weetbix_proto_v1_analyzed_test_variant_proto_depIdxs,
		EnumInfos:         file_infra_appengine_weetbix_proto_v1_analyzed_test_variant_proto_enumTypes,
		MessageInfos:      file_infra_appengine_weetbix_proto_v1_analyzed_test_variant_proto_msgTypes,
	}.Build()
	File_infra_appengine_weetbix_proto_v1_analyzed_test_variant_proto = out.File
	file_infra_appengine_weetbix_proto_v1_analyzed_test_variant_proto_rawDesc = nil
	file_infra_appengine_weetbix_proto_v1_analyzed_test_variant_proto_goTypes = nil
	file_infra_appengine_weetbix_proto_v1_analyzed_test_variant_proto_depIdxs = nil
}
