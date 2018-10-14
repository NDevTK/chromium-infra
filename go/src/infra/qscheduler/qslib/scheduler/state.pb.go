// Code generated by protoc-gen-go. DO NOT EDIT.
// source: infra/qscheduler/qslib/scheduler/state.proto

package scheduler

import (
	fmt "fmt"
	task "infra/qscheduler/qslib/types/task"
	vector "infra/qscheduler/qslib/types/vector"
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
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// State represents the overall state of a quota scheduler worker pool,
// account set, and task queue. This is represented separately from
// configuration information. The state is expected to be updated frequently,
// on each scheduler tick.
type State struct {
	// QueuedRequests is the set of Requests that are waiting to be assigned to a
	// worker, keyed by request id.
	QueuedRequests map[string]*task.Request `protobuf:"bytes,1,rep,name=queued_requests,json=queuedRequests,proto3" json:"queued_requests,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// Balance of all quota accounts for this pool, keyed by account id.
	Balances map[string]*vector.Vector `protobuf:"bytes,2,rep,name=balances,proto3" json:"balances,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// Workers that may run tasks, and their states, keyed by worker id.
	Workers map[string]*Worker `protobuf:"bytes,3,rep,name=workers,proto3" json:"workers,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// LastUpdateTime is the last time at which UpdateTime was called on a scheduler,
	// and corresponds to the when the quota account balances were updated.
	LastUpdateTime       *timestamp.Timestamp `protobuf:"bytes,4,opt,name=last_update_time,json=lastUpdateTime,proto3" json:"last_update_time,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *State) Reset()         { *m = State{} }
func (m *State) String() string { return proto.CompactTextString(m) }
func (*State) ProtoMessage()    {}
func (*State) Descriptor() ([]byte, []int) {
	return fileDescriptor_d6b2685915185dd1, []int{0}
}

func (m *State) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_State.Unmarshal(m, b)
}
func (m *State) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_State.Marshal(b, m, deterministic)
}
func (m *State) XXX_Merge(src proto.Message) {
	xxx_messageInfo_State.Merge(m, src)
}
func (m *State) XXX_Size() int {
	return xxx_messageInfo_State.Size(m)
}
func (m *State) XXX_DiscardUnknown() {
	xxx_messageInfo_State.DiscardUnknown(m)
}

var xxx_messageInfo_State proto.InternalMessageInfo

func (m *State) GetQueuedRequests() map[string]*task.Request {
	if m != nil {
		return m.QueuedRequests
	}
	return nil
}

func (m *State) GetBalances() map[string]*vector.Vector {
	if m != nil {
		return m.Balances
	}
	return nil
}

func (m *State) GetWorkers() map[string]*Worker {
	if m != nil {
		return m.Workers
	}
	return nil
}

func (m *State) GetLastUpdateTime() *timestamp.Timestamp {
	if m != nil {
		return m.LastUpdateTime
	}
	return nil
}

func init() {
	proto.RegisterType((*State)(nil), "scheduler.State")
	proto.RegisterMapType((map[string]*vector.Vector)(nil), "scheduler.State.BalancesEntry")
	proto.RegisterMapType((map[string]*task.Request)(nil), "scheduler.State.QueuedRequestsEntry")
	proto.RegisterMapType((map[string]*Worker)(nil), "scheduler.State.WorkersEntry")
}

func init() {
	proto.RegisterFile("infra/qscheduler/qslib/scheduler/state.proto", fileDescriptor_d6b2685915185dd1)
}

var fileDescriptor_d6b2685915185dd1 = []byte{
	// 361 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x91, 0xcd, 0x4f, 0xf2, 0x40,
	0x10, 0xc6, 0x53, 0xfa, 0xf2, 0x2a, 0x8b, 0x20, 0xae, 0x97, 0xa6, 0x89, 0x4a, 0x94, 0x44, 0x0e,
	0xba, 0x35, 0x78, 0xd0, 0x70, 0x34, 0x7a, 0x32, 0x24, 0x5a, 0xbf, 0x8e, 0x64, 0x0b, 0x03, 0x12,
	0x0a, 0x6d, 0xf7, 0x03, 0xc3, 0x5f, 0xaf, 0xe9, 0xee, 0x82, 0x6d, 0xa8, 0x7a, 0xe9, 0xb6, 0x33,
	0xcf, 0xf3, 0xeb, 0xce, 0x3c, 0xe8, 0x6c, 0x32, 0x1f, 0x31, 0xea, 0x25, 0x7c, 0xf0, 0x0e, 0x43,
	0x19, 0x02, 0xf3, 0x12, 0x1e, 0x4e, 0x02, 0xef, 0xfb, 0x9b, 0x0b, 0x2a, 0x80, 0xc4, 0x2c, 0x12,
	0x11, 0xae, 0xac, 0xcb, 0xee, 0xd1, 0x38, 0x8a, 0xc6, 0x21, 0x78, 0xaa, 0x11, 0xc8, 0x91, 0x27,
	0x26, 0x33, 0xe0, 0x82, 0xce, 0x62, 0xad, 0x75, 0x7f, 0x22, 0x8b, 0x65, 0x0c, 0xdc, 0x13, 0x94,
	0x4f, 0xd5, 0xc3, 0xa8, 0x2f, 0x7e, 0x55, 0x2f, 0x60, 0x20, 0x22, 0x66, 0x0e, 0xe3, 0x38, 0xff,
	0xf3, 0xe6, 0x1f, 0x11, 0x9b, 0x82, 0x91, 0x1f, 0x7f, 0xda, 0xa8, 0xfc, 0x94, 0x8e, 0x82, 0x7b,
	0x68, 0x37, 0x91, 0x20, 0x61, 0xd8, 0x67, 0x90, 0x48, 0xe0, 0x82, 0x3b, 0x56, 0xd3, 0x6e, 0x57,
	0x3b, 0x2d, 0xb2, 0xf6, 0x12, 0x25, 0x25, 0x8f, 0x4a, 0xe7, 0x1b, 0xd9, 0xdd, 0x5c, 0xb0, 0xa5,
	0x5f, 0x4f, 0x72, 0x45, 0xdc, 0x45, 0xdb, 0x01, 0x0d, 0xe9, 0x7c, 0x00, 0xdc, 0x29, 0x29, 0xce,
	0xe1, 0x06, 0xe7, 0xc6, 0x08, 0x34, 0x61, 0xad, 0xc7, 0x57, 0x68, 0x4b, 0x5f, 0x92, 0x3b, 0xb6,
	0xb2, 0x1e, 0x6c, 0x58, 0xdf, 0x74, 0x5f, 0x3b, 0x57, 0x6a, 0x7c, 0x8b, 0x1a, 0x21, 0xe5, 0xa2,
	0x2f, 0xe3, 0x21, 0x15, 0xd0, 0x4f, 0x77, 0xef, 0xfc, 0x6b, 0x5a, 0xed, 0x6a, 0xc7, 0x25, 0x3a,
	0x18, 0xb2, 0x0a, 0x86, 0x3c, 0xaf, 0x82, 0xf1, 0xeb, 0xa9, 0xe7, 0x45, 0x59, 0xd2, 0xa2, 0xfb,
	0x80, 0xf6, 0x0b, 0x26, 0xc4, 0x0d, 0x64, 0x4f, 0x61, 0xe9, 0x58, 0x4d, 0xab, 0x5d, 0xf1, 0xd3,
	0x57, 0x7c, 0x82, 0xca, 0x0b, 0x1a, 0x4a, 0x70, 0x4a, 0xea, 0x1f, 0x35, 0xa2, 0x92, 0x33, 0x2e,
	0x5f, 0xf7, 0xba, 0xa5, 0x6b, 0xcb, 0xbd, 0x47, 0xb5, 0xdc, 0xac, 0x05, 0xac, 0x56, 0x9e, 0x55,
	0x27, 0x26, 0xd5, 0x57, 0x75, 0x64, 0x61, 0x3d, 0xb4, 0x93, 0x9d, 0xbe, 0x80, 0x75, 0x9a, 0x67,
	0xed, 0x65, 0xb6, 0xa7, 0x9d, 0x19, 0x5c, 0xf0, 0x5f, 0x6d, 0xe4, 0xf2, 0x2b, 0x00, 0x00, 0xff,
	0xff, 0x8a, 0xda, 0x0d, 0x73, 0xf3, 0x02, 0x00, 0x00,
}
