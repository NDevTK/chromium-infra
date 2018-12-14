// Code generated by protoc-gen-go. DO NOT EDIT.
// source: infra/qscheduler/qslib/scheduler/worker.proto

package scheduler

import (
	fmt "fmt"
	_ "infra/qscheduler/qslib/types/account"
	_ "infra/qscheduler/qslib/types/vector"
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

// Worker represents a resource that can run 1 task at a time. This corresponds
// to the swarming concept of a Bot. This representation considers only the
// subset of Labels that are Provisionable (can be changed by running a task),
// because the quota scheduler algorithm is expected to run against a pool of
// otherwise homogenous workers.
type Worker struct {
	// Labels represents the set of provisionable labels that this worker
	// possesses.
	Labels []string `protobuf:"bytes,1,rep,name=labels,proto3" json:"labels,omitempty"`
	// RunningTask is, if non-nil, the task that is currently running on the
	// worker.
	RunningTask *TaskRun `protobuf:"bytes,2,opt,name=running_task,json=runningTask,proto3" json:"running_task,omitempty"`
	// ConfirmedIdleTime is the most recent time at which the Worker state was
	// directly confirmed as idle by external authority (via a call to MarkIdle or
	// NotifyRequest).
	ConfirmedTime        *timestamp.Timestamp `protobuf:"bytes,3,opt,name=confirmed_time,json=confirmedTime,proto3" json:"confirmed_time,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Worker) Reset()         { *m = Worker{} }
func (m *Worker) String() string { return proto.CompactTextString(m) }
func (*Worker) ProtoMessage()    {}
func (*Worker) Descriptor() ([]byte, []int) {
	return fileDescriptor_9853fc992bf091d9, []int{0}
}

func (m *Worker) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Worker.Unmarshal(m, b)
}
func (m *Worker) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Worker.Marshal(b, m, deterministic)
}
func (m *Worker) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Worker.Merge(m, src)
}
func (m *Worker) XXX_Size() int {
	return xxx_messageInfo_Worker.Size(m)
}
func (m *Worker) XXX_DiscardUnknown() {
	xxx_messageInfo_Worker.DiscardUnknown(m)
}

var xxx_messageInfo_Worker proto.InternalMessageInfo

func (m *Worker) GetLabels() []string {
	if m != nil {
		return m.Labels
	}
	return nil
}

func (m *Worker) GetRunningTask() *TaskRun {
	if m != nil {
		return m.RunningTask
	}
	return nil
}

func (m *Worker) GetConfirmedTime() *timestamp.Timestamp {
	if m != nil {
		return m.ConfirmedTime
	}
	return nil
}

func init() {
	proto.RegisterType((*Worker)(nil), "scheduler.Worker")
}

func init() {
	proto.RegisterFile("infra/qscheduler/qslib/scheduler/worker.proto", fileDescriptor_9853fc992bf091d9)
}

var fileDescriptor_9853fc992bf091d9 = []byte{
	// 242 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x90, 0x41, 0x4a, 0xc4, 0x30,
	0x14, 0x86, 0xa9, 0x03, 0x85, 0xc9, 0xa8, 0x8b, 0x2c, 0xa4, 0x74, 0xe3, 0xe0, 0x6a, 0x40, 0x4c,
	0x64, 0xc4, 0x03, 0x78, 0x85, 0x32, 0xe0, 0x72, 0x48, 0x33, 0xaf, 0x35, 0x34, 0x4d, 0x3a, 0x2f,
	0x89, 0xe2, 0x55, 0x3c, 0xad, 0x34, 0x4d, 0xeb, 0x4a, 0x5d, 0x3d, 0xfe, 0xd7, 0xef, 0x6b, 0x92,
	0x9f, 0x3c, 0x28, 0xd3, 0xa0, 0xe0, 0x67, 0x27, 0xdf, 0xe0, 0x14, 0x34, 0x20, 0x3f, 0x3b, 0xad,
	0x6a, 0xfe, 0x93, 0x3f, 0x2c, 0x76, 0x80, 0x6c, 0x40, 0xeb, 0x2d, 0x5d, 0x2f, 0xfb, 0xf2, 0xb6,
	0xb5, 0xb6, 0xd5, 0xc0, 0xe3, 0x87, 0x3a, 0x34, 0xdc, 0xab, 0x1e, 0x9c, 0x17, 0xfd, 0x30, 0xb1,
	0xe5, 0xfd, 0xbf, 0xbf, 0xf6, 0xc2, 0x75, 0x09, 0x7e, 0xfc, 0x05, 0xf6, 0x9f, 0x03, 0x38, 0xfe,
	0x0e, 0xd2, 0x5b, 0x4c, 0x23, 0x19, 0xfb, 0x3f, 0x0d, 0x21, 0xa5, 0x0d, 0xc6, 0xcf, 0x73, 0x72,
	0xee, 0xbe, 0x32, 0x92, 0xbf, 0xc6, 0xf7, 0xd0, 0x1b, 0x92, 0x6b, 0x51, 0x83, 0x76, 0x45, 0xb6,
	0x5d, 0xed, 0xd6, 0x55, 0x4a, 0xf4, 0x99, 0x5c, 0x62, 0x30, 0x46, 0x99, 0xf6, 0x38, 0x5e, 0xaf,
	0xb8, 0xd8, 0x66, 0xbb, 0xcd, 0x9e, 0xb2, 0xe5, 0x18, 0x76, 0x10, 0xae, 0xab, 0x82, 0xa9, 0x36,
	0x89, 0x1b, 0x33, 0x7d, 0x21, 0xd7, 0xd2, 0x9a, 0x46, 0x61, 0x0f, 0xa7, 0xe3, 0xd8, 0x44, 0xb1,
	0x8a, 0x62, 0xc9, 0xa6, 0x9a, 0xd8, 0x5c, 0x13, 0x3b, 0xcc, 0x35, 0x55, 0x57, 0x8b, 0x31, 0xee,
	0xea, 0x3c, 0x22, 0x4f, 0xdf, 0x01, 0x00, 0x00, 0xff, 0xff, 0xf0, 0xe4, 0xc6, 0x70, 0x93, 0x01,
	0x00, 0x00,
}
