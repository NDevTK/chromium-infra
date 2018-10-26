// Code generated by protoc-gen-go. DO NOT EDIT.
// source: infra/qscheduler/qslib/scheduler/scheduler.proto

package scheduler

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Scheduler encapsulates the state and configuration of a running
// quotascheduler for a single pool, and its methods provide an implementation
// of the quotascheduler algorithm.
type Scheduler struct {
	// State is the state of the scheduler.
	State *State `protobuf:"bytes,1,opt,name=state,proto3" json:"state,omitempty"`
	// Config is the config of the scheduler.
	Config               *Config  `protobuf:"bytes,2,opt,name=config,proto3" json:"config,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Scheduler) Reset()         { *m = Scheduler{} }
func (m *Scheduler) String() string { return proto.CompactTextString(m) }
func (*Scheduler) ProtoMessage()    {}
func (*Scheduler) Descriptor() ([]byte, []int) {
	return fileDescriptor_9c9a368e1d01bf8a, []int{0}
}

func (m *Scheduler) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Scheduler.Unmarshal(m, b)
}
func (m *Scheduler) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Scheduler.Marshal(b, m, deterministic)
}
func (m *Scheduler) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Scheduler.Merge(m, src)
}
func (m *Scheduler) XXX_Size() int {
	return xxx_messageInfo_Scheduler.Size(m)
}
func (m *Scheduler) XXX_DiscardUnknown() {
	xxx_messageInfo_Scheduler.DiscardUnknown(m)
}

var xxx_messageInfo_Scheduler proto.InternalMessageInfo

func (m *Scheduler) GetState() *State {
	if m != nil {
		return m.State
	}
	return nil
}

func (m *Scheduler) GetConfig() *Config {
	if m != nil {
		return m.Config
	}
	return nil
}

func init() {
	proto.RegisterType((*Scheduler)(nil), "scheduler.Scheduler")
}

func init() {
	proto.RegisterFile("infra/qscheduler/qslib/scheduler/scheduler.proto", fileDescriptor_9c9a368e1d01bf8a)
}

var fileDescriptor_9c9a368e1d01bf8a = []byte{
	// 139 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x32, 0xc8, 0xcc, 0x4b, 0x2b,
	0x4a, 0xd4, 0x2f, 0x2c, 0x4e, 0xce, 0x48, 0x4d, 0x29, 0xcd, 0x49, 0x2d, 0xd2, 0x2f, 0x2c, 0xce,
	0xc9, 0x4c, 0xd2, 0x47, 0xf0, 0xe1, 0x2c, 0xbd, 0x82, 0xa2, 0xfc, 0x92, 0x7c, 0x21, 0x4e, 0xb8,
	0x80, 0x94, 0x0e, 0x61, 0xcd, 0x25, 0x89, 0x25, 0xa9, 0x10, 0x8d, 0x52, 0xba, 0x04, 0x55, 0x27,
	0xe7, 0xe7, 0xa5, 0x65, 0xa6, 0x43, 0x94, 0x2b, 0xc5, 0x71, 0x71, 0x06, 0xc3, 0x64, 0x84, 0xd4,
	0xb8, 0x58, 0xc1, 0x46, 0x49, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x1b, 0x09, 0xe8, 0x21, 0x5c, 0x15,
	0x0c, 0x12, 0x0f, 0x82, 0x48, 0x0b, 0x69, 0x72, 0xb1, 0x41, 0x0c, 0x91, 0x60, 0x02, 0x2b, 0x14,
	0x44, 0x52, 0xe8, 0x0c, 0x96, 0x08, 0x82, 0x2a, 0x48, 0x62, 0x03, 0x5b, 0x63, 0x0c, 0x08, 0x00,
	0x00, 0xff, 0xff, 0xa1, 0xf1, 0xa7, 0xca, 0x02, 0x01, 0x00, 0x00,
}
