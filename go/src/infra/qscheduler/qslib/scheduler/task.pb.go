// Code generated by protoc-gen-go. DO NOT EDIT.
// source: infra/qscheduler/qslib/scheduler/task.proto

package scheduler

import (
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// TaskRequest represents a requested task in the queue, and refers to the
// quota account to run it against. This representation intentionally
// excludes most of the details of a Swarming task request.
type TaskRequest struct {
	// AccountId is the id of the account that this request charges to.
	AccountId string `protobuf:"bytes,1,opt,name=account_id,json=accountId,proto3" json:"account_id,omitempty"`
	// EnqueueTime is the time at which the request was enqueued.
	EnqueueTime *timestamp.Timestamp `protobuf:"bytes,2,opt,name=enqueue_time,json=enqueueTime,proto3" json:"enqueue_time,omitempty"`
	// The set of Provisionable Labels for this task.
	Labels []string `protobuf:"bytes,3,rep,name=labels,proto3" json:"labels,omitempty"`
	// ConfirmedTime is the most recent time at which the Request state was
	// provided or confirmed by external authority (via a call to Enforce or
	// AddRequest).
	ConfirmedTime        *timestamp.Timestamp `protobuf:"bytes,4,opt,name=confirmed_time,json=confirmedTime,proto3" json:"confirmed_time,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *TaskRequest) Reset()         { *m = TaskRequest{} }
func (m *TaskRequest) String() string { return proto.CompactTextString(m) }
func (*TaskRequest) ProtoMessage()    {}
func (*TaskRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_2be5499420917688, []int{0}
}

func (m *TaskRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TaskRequest.Unmarshal(m, b)
}
func (m *TaskRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TaskRequest.Marshal(b, m, deterministic)
}
func (m *TaskRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TaskRequest.Merge(m, src)
}
func (m *TaskRequest) XXX_Size() int {
	return xxx_messageInfo_TaskRequest.Size(m)
}
func (m *TaskRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_TaskRequest.DiscardUnknown(m)
}

var xxx_messageInfo_TaskRequest proto.InternalMessageInfo

func (m *TaskRequest) GetAccountId() string {
	if m != nil {
		return m.AccountId
	}
	return ""
}

func (m *TaskRequest) GetEnqueueTime() *timestamp.Timestamp {
	if m != nil {
		return m.EnqueueTime
	}
	return nil
}

func (m *TaskRequest) GetLabels() []string {
	if m != nil {
		return m.Labels
	}
	return nil
}

func (m *TaskRequest) GetConfirmedTime() *timestamp.Timestamp {
	if m != nil {
		return m.ConfirmedTime
	}
	return nil
}

// TaskRun represents a task that has been assigned to a worker and is
// now running.
type TaskRun struct {
	// Cost is the total cost that has been incurred on this task while running.
	Cost []float64 `protobuf:"fixed64,1,rep,packed,name=cost,proto3" json:"cost,omitempty"`
	// Request is the request that this running task corresponds to.
	Request *TaskRequest `protobuf:"bytes,2,opt,name=request,proto3" json:"request,omitempty"`
	// RequestId is the request id of the request that this running task
	// corresponds to.
	RequestId string `protobuf:"bytes,3,opt,name=request_id,json=requestId,proto3" json:"request_id,omitempty"`
	// Priority is the current priority level of the running task.
	Priority             int32    `protobuf:"varint,4,opt,name=priority,proto3" json:"priority,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TaskRun) Reset()         { *m = TaskRun{} }
func (m *TaskRun) String() string { return proto.CompactTextString(m) }
func (*TaskRun) ProtoMessage()    {}
func (*TaskRun) Descriptor() ([]byte, []int) {
	return fileDescriptor_2be5499420917688, []int{1}
}

func (m *TaskRun) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TaskRun.Unmarshal(m, b)
}
func (m *TaskRun) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TaskRun.Marshal(b, m, deterministic)
}
func (m *TaskRun) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TaskRun.Merge(m, src)
}
func (m *TaskRun) XXX_Size() int {
	return xxx_messageInfo_TaskRun.Size(m)
}
func (m *TaskRun) XXX_DiscardUnknown() {
	xxx_messageInfo_TaskRun.DiscardUnknown(m)
}

var xxx_messageInfo_TaskRun proto.InternalMessageInfo

func (m *TaskRun) GetCost() []float64 {
	if m != nil {
		return m.Cost
	}
	return nil
}

func (m *TaskRun) GetRequest() *TaskRequest {
	if m != nil {
		return m.Request
	}
	return nil
}

func (m *TaskRun) GetRequestId() string {
	if m != nil {
		return m.RequestId
	}
	return ""
}

func (m *TaskRun) GetPriority() int32 {
	if m != nil {
		return m.Priority
	}
	return 0
}

func init() {
	proto.RegisterType((*TaskRequest)(nil), "scheduler.TaskRequest")
	proto.RegisterType((*TaskRun)(nil), "scheduler.TaskRun")
}

func init() {
	proto.RegisterFile("infra/qscheduler/qslib/scheduler/task.proto", fileDescriptor_2be5499420917688)
}

var fileDescriptor_2be5499420917688 = []byte{
	// 281 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x50, 0x3d, 0x4f, 0xc3, 0x30,
	0x14, 0x94, 0x49, 0x69, 0xc9, 0x2b, 0x30, 0x78, 0xa8, 0xa2, 0x48, 0x88, 0xa8, 0x53, 0x24, 0x24,
	0x07, 0xc1, 0xcc, 0xc0, 0xd8, 0xd5, 0xea, 0x5e, 0x39, 0xc9, 0x4b, 0xb1, 0x9a, 0xc4, 0x89, 0x3f,
	0x06, 0xfe, 0x02, 0x3f, 0x8c, 0xdf, 0x85, 0x92, 0x38, 0x81, 0x8d, 0xed, 0xbd, 0xf3, 0xf9, 0xee,
	0xde, 0xc1, 0x93, 0x6c, 0x2b, 0x2d, 0xb2, 0xde, 0x14, 0x1f, 0x58, 0xba, 0x1a, 0x75, 0xd6, 0x9b,
	0x5a, 0xe6, 0xd9, 0xef, 0x6e, 0x85, 0xb9, 0xb0, 0x4e, 0x2b, 0xab, 0x68, 0xb8, 0xa0, 0xf1, 0xe3,
	0x59, 0xa9, 0x73, 0x8d, 0xd9, 0xf8, 0x90, 0xbb, 0x2a, 0xb3, 0xb2, 0x41, 0x63, 0x45, 0xd3, 0x4d,
	0xdc, 0xfd, 0x37, 0x81, 0xed, 0x51, 0x98, 0x0b, 0xc7, 0xde, 0xa1, 0xb1, 0xf4, 0x01, 0x40, 0x14,
	0x85, 0x72, 0xad, 0x3d, 0xc9, 0x32, 0x22, 0x09, 0x49, 0x43, 0x1e, 0x7a, 0xe4, 0x50, 0xd2, 0x37,
	0xb8, 0xc5, 0xb6, 0x77, 0xe8, 0xf0, 0x34, 0x28, 0x45, 0x57, 0x09, 0x49, 0xb7, 0x2f, 0x31, 0x9b,
	0x6c, 0xd8, 0x6c, 0xc3, 0x8e, 0xb3, 0x0d, 0xdf, 0x7a, 0xfe, 0x80, 0xd0, 0x1d, 0xac, 0x6b, 0x91,
	0x63, 0x6d, 0xa2, 0x20, 0x09, 0xd2, 0x90, 0xfb, 0x8d, 0xbe, 0xc3, 0x7d, 0xa1, 0xda, 0x4a, 0xea,
	0x06, 0xcb, 0x49, 0x78, 0xf5, 0xaf, 0xf0, 0xdd, 0xf2, 0x63, 0xc0, 0xf6, 0x5f, 0x04, 0x36, 0xe3,
	0x21, 0xae, 0xa5, 0x14, 0x56, 0x85, 0x32, 0x36, 0x22, 0x49, 0x90, 0x12, 0x3e, 0xce, 0xf4, 0x19,
	0x36, 0x7a, 0xba, 0xd1, 0x87, 0xde, 0xb1, 0xa5, 0x26, 0xf6, 0xa7, 0x01, 0x3e, 0xd3, 0x86, 0x2a,
	0xfc, 0x38, 0x54, 0x11, 0x4c, 0x55, 0x78, 0xe4, 0x50, 0xd2, 0x18, 0x6e, 0x3a, 0x2d, 0x95, 0x96,
	0xf6, 0x73, 0x4c, 0x7b, 0xcd, 0x97, 0x3d, 0x5f, 0x8f, 0x79, 0x5f, 0x7f, 0x02, 0x00, 0x00, 0xff,
	0xff, 0xd1, 0xd7, 0xdb, 0x5c, 0xb7, 0x01, 0x00, 0x00,
}
