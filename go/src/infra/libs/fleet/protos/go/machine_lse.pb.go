// Code generated by protoc-gen-go. DO NOT EDIT.
// source: machine_lse.proto

package _go

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	lab "go.chromium.org/chromiumos/infra/proto/go/lab"
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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type MachineLSE struct {
	Id *LabSetupEnvID `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// The prototype that this machine LSE should follow. System will use this
	// prototype to detect if the LSE is completed or valid.
	PrototypeId *MachineLSEPrototypeID `protobuf:"bytes,2,opt,name=prototype_id,json=prototypeId,proto3" json:"prototype_id,omitempty"`
	// Types that are valid to be assigned to Lse:
	//	*MachineLSE_ChromeMachineLse
	//	*MachineLSE_ChromeosMachineLse
	Lse isMachineLSE_Lse `protobuf_oneof:"lse"`
	// The machines that this LSE is linked to. No machine is linked if it's NULL.
	Machines             []*MachineID `protobuf:"bytes,5,rep,name=machines,proto3" json:"machines,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *MachineLSE) Reset()         { *m = MachineLSE{} }
func (m *MachineLSE) String() string { return proto.CompactTextString(m) }
func (*MachineLSE) ProtoMessage()    {}
func (*MachineLSE) Descriptor() ([]byte, []int) {
	return fileDescriptor_6f342dc1a43117d0, []int{0}
}

func (m *MachineLSE) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MachineLSE.Unmarshal(m, b)
}
func (m *MachineLSE) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MachineLSE.Marshal(b, m, deterministic)
}
func (m *MachineLSE) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MachineLSE.Merge(m, src)
}
func (m *MachineLSE) XXX_Size() int {
	return xxx_messageInfo_MachineLSE.Size(m)
}
func (m *MachineLSE) XXX_DiscardUnknown() {
	xxx_messageInfo_MachineLSE.DiscardUnknown(m)
}

var xxx_messageInfo_MachineLSE proto.InternalMessageInfo

func (m *MachineLSE) GetId() *LabSetupEnvID {
	if m != nil {
		return m.Id
	}
	return nil
}

func (m *MachineLSE) GetPrototypeId() *MachineLSEPrototypeID {
	if m != nil {
		return m.PrototypeId
	}
	return nil
}

type isMachineLSE_Lse interface {
	isMachineLSE_Lse()
}

type MachineLSE_ChromeMachineLse struct {
	ChromeMachineLse *ChromeMachineLSE `protobuf:"bytes,3,opt,name=chrome_machine_lse,json=chromeMachineLse,proto3,oneof"`
}

type MachineLSE_ChromeosMachineLse struct {
	ChromeosMachineLse *ChromeOSMachineLSE `protobuf:"bytes,4,opt,name=chromeos_machine_lse,json=chromeosMachineLse,proto3,oneof"`
}

func (*MachineLSE_ChromeMachineLse) isMachineLSE_Lse() {}

func (*MachineLSE_ChromeosMachineLse) isMachineLSE_Lse() {}

func (m *MachineLSE) GetLse() isMachineLSE_Lse {
	if m != nil {
		return m.Lse
	}
	return nil
}

func (m *MachineLSE) GetChromeMachineLse() *ChromeMachineLSE {
	if x, ok := m.GetLse().(*MachineLSE_ChromeMachineLse); ok {
		return x.ChromeMachineLse
	}
	return nil
}

func (m *MachineLSE) GetChromeosMachineLse() *ChromeOSMachineLSE {
	if x, ok := m.GetLse().(*MachineLSE_ChromeosMachineLse); ok {
		return x.ChromeosMachineLse
	}
	return nil
}

func (m *MachineLSE) GetMachines() []*MachineID {
	if m != nil {
		return m.Machines
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*MachineLSE) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*MachineLSE_ChromeMachineLse)(nil),
		(*MachineLSE_ChromeosMachineLse)(nil),
	}
}

type ChromeMachineLSE struct {
	// The hostname is also recorded in DHCP configs
	Hostname string `protobuf:"bytes,1,opt,name=hostname,proto3" json:"hostname,omitempty"`
	// Indicate if VM is needed to set up
	Vms                  []*VM    `protobuf:"bytes,2,rep,name=vms,proto3" json:"vms,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ChromeMachineLSE) Reset()         { *m = ChromeMachineLSE{} }
func (m *ChromeMachineLSE) String() string { return proto.CompactTextString(m) }
func (*ChromeMachineLSE) ProtoMessage()    {}
func (*ChromeMachineLSE) Descriptor() ([]byte, []int) {
	return fileDescriptor_6f342dc1a43117d0, []int{1}
}

func (m *ChromeMachineLSE) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChromeMachineLSE.Unmarshal(m, b)
}
func (m *ChromeMachineLSE) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChromeMachineLSE.Marshal(b, m, deterministic)
}
func (m *ChromeMachineLSE) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChromeMachineLSE.Merge(m, src)
}
func (m *ChromeMachineLSE) XXX_Size() int {
	return xxx_messageInfo_ChromeMachineLSE.Size(m)
}
func (m *ChromeMachineLSE) XXX_DiscardUnknown() {
	xxx_messageInfo_ChromeMachineLSE.DiscardUnknown(m)
}

var xxx_messageInfo_ChromeMachineLSE proto.InternalMessageInfo

func (m *ChromeMachineLSE) GetHostname() string {
	if m != nil {
		return m.Hostname
	}
	return ""
}

func (m *ChromeMachineLSE) GetVms() []*VM {
	if m != nil {
		return m.Vms
	}
	return nil
}

type VM struct {
	// A unique vm name
	Id                   *VMID      `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	OsVersion            *OSVersion `protobuf:"bytes,2,opt,name=os_version,json=osVersion,proto3" json:"os_version,omitempty"`
	MacAddress           string     `protobuf:"bytes,3,opt,name=mac_address,json=macAddress,proto3" json:"mac_address,omitempty"`
	Hostname             string     `protobuf:"bytes,4,opt,name=hostname,proto3" json:"hostname,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *VM) Reset()         { *m = VM{} }
func (m *VM) String() string { return proto.CompactTextString(m) }
func (*VM) ProtoMessage()    {}
func (*VM) Descriptor() ([]byte, []int) {
	return fileDescriptor_6f342dc1a43117d0, []int{2}
}

func (m *VM) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VM.Unmarshal(m, b)
}
func (m *VM) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VM.Marshal(b, m, deterministic)
}
func (m *VM) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VM.Merge(m, src)
}
func (m *VM) XXX_Size() int {
	return xxx_messageInfo_VM.Size(m)
}
func (m *VM) XXX_DiscardUnknown() {
	xxx_messageInfo_VM.DiscardUnknown(m)
}

var xxx_messageInfo_VM proto.InternalMessageInfo

func (m *VM) GetId() *VMID {
	if m != nil {
		return m.Id
	}
	return nil
}

func (m *VM) GetOsVersion() *OSVersion {
	if m != nil {
		return m.OsVersion
	}
	return nil
}

func (m *VM) GetMacAddress() string {
	if m != nil {
		return m.MacAddress
	}
	return ""
}

func (m *VM) GetHostname() string {
	if m != nil {
		return m.Hostname
	}
	return ""
}

type VMID struct {
	Value                string   `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *VMID) Reset()         { *m = VMID{} }
func (m *VMID) String() string { return proto.CompactTextString(m) }
func (*VMID) ProtoMessage()    {}
func (*VMID) Descriptor() ([]byte, []int) {
	return fileDescriptor_6f342dc1a43117d0, []int{3}
}

func (m *VMID) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VMID.Unmarshal(m, b)
}
func (m *VMID) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VMID.Marshal(b, m, deterministic)
}
func (m *VMID) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VMID.Merge(m, src)
}
func (m *VMID) XXX_Size() int {
	return xxx_messageInfo_VMID.Size(m)
}
func (m *VMID) XXX_DiscardUnknown() {
	xxx_messageInfo_VMID.DiscardUnknown(m)
}

var xxx_messageInfo_VMID proto.InternalMessageInfo

func (m *VMID) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

type ChromeOSMachineLSE struct {
	Hostname string `protobuf:"bytes,1,opt,name=hostname,proto3" json:"hostname,omitempty"`
	// Types that are valid to be assigned to ChromeosLse:
	//	*ChromeOSMachineLSE_Dut
	//	*ChromeOSMachineLSE_Server
	ChromeosLse          isChromeOSMachineLSE_ChromeosLse `protobuf_oneof:"chromeos_lse"`
	XXX_NoUnkeyedLiteral struct{}                         `json:"-"`
	XXX_unrecognized     []byte                           `json:"-"`
	XXX_sizecache        int32                            `json:"-"`
}

func (m *ChromeOSMachineLSE) Reset()         { *m = ChromeOSMachineLSE{} }
func (m *ChromeOSMachineLSE) String() string { return proto.CompactTextString(m) }
func (*ChromeOSMachineLSE) ProtoMessage()    {}
func (*ChromeOSMachineLSE) Descriptor() ([]byte, []int) {
	return fileDescriptor_6f342dc1a43117d0, []int{4}
}

func (m *ChromeOSMachineLSE) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChromeOSMachineLSE.Unmarshal(m, b)
}
func (m *ChromeOSMachineLSE) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChromeOSMachineLSE.Marshal(b, m, deterministic)
}
func (m *ChromeOSMachineLSE) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChromeOSMachineLSE.Merge(m, src)
}
func (m *ChromeOSMachineLSE) XXX_Size() int {
	return xxx_messageInfo_ChromeOSMachineLSE.Size(m)
}
func (m *ChromeOSMachineLSE) XXX_DiscardUnknown() {
	xxx_messageInfo_ChromeOSMachineLSE.DiscardUnknown(m)
}

var xxx_messageInfo_ChromeOSMachineLSE proto.InternalMessageInfo

func (m *ChromeOSMachineLSE) GetHostname() string {
	if m != nil {
		return m.Hostname
	}
	return ""
}

type isChromeOSMachineLSE_ChromeosLse interface {
	isChromeOSMachineLSE_ChromeosLse()
}

type ChromeOSMachineLSE_Dut struct {
	Dut *ChromeOSDeviceLSE `protobuf:"bytes,2,opt,name=dut,proto3,oneof"`
}

type ChromeOSMachineLSE_Server struct {
	Server *ChromeOSServerLSE `protobuf:"bytes,3,opt,name=server,proto3,oneof"`
}

func (*ChromeOSMachineLSE_Dut) isChromeOSMachineLSE_ChromeosLse() {}

func (*ChromeOSMachineLSE_Server) isChromeOSMachineLSE_ChromeosLse() {}

func (m *ChromeOSMachineLSE) GetChromeosLse() isChromeOSMachineLSE_ChromeosLse {
	if m != nil {
		return m.ChromeosLse
	}
	return nil
}

func (m *ChromeOSMachineLSE) GetDut() *ChromeOSDeviceLSE {
	if x, ok := m.GetChromeosLse().(*ChromeOSMachineLSE_Dut); ok {
		return x.Dut
	}
	return nil
}

func (m *ChromeOSMachineLSE) GetServer() *ChromeOSServerLSE {
	if x, ok := m.GetChromeosLse().(*ChromeOSMachineLSE_Server); ok {
		return x.Server
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*ChromeOSMachineLSE) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*ChromeOSMachineLSE_Dut)(nil),
		(*ChromeOSMachineLSE_Server)(nil),
	}
}

type ChromeOSDeviceLSE struct {
	Config                 *lab.DeviceUnderTest `protobuf:"bytes,1,opt,name=config,proto3" json:"config,omitempty"`
	RpmInterface           *RPMInterface        `protobuf:"bytes,2,opt,name=rpm_interface,json=rpmInterface,proto3" json:"rpm_interface,omitempty"`
	NetworkDeviceInterface *SwitchInterface     `protobuf:"bytes,3,opt,name=network_device_interface,json=networkDeviceInterface,proto3" json:"network_device_interface,omitempty"`
	XXX_NoUnkeyedLiteral   struct{}             `json:"-"`
	XXX_unrecognized       []byte               `json:"-"`
	XXX_sizecache          int32                `json:"-"`
}

func (m *ChromeOSDeviceLSE) Reset()         { *m = ChromeOSDeviceLSE{} }
func (m *ChromeOSDeviceLSE) String() string { return proto.CompactTextString(m) }
func (*ChromeOSDeviceLSE) ProtoMessage()    {}
func (*ChromeOSDeviceLSE) Descriptor() ([]byte, []int) {
	return fileDescriptor_6f342dc1a43117d0, []int{5}
}

func (m *ChromeOSDeviceLSE) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChromeOSDeviceLSE.Unmarshal(m, b)
}
func (m *ChromeOSDeviceLSE) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChromeOSDeviceLSE.Marshal(b, m, deterministic)
}
func (m *ChromeOSDeviceLSE) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChromeOSDeviceLSE.Merge(m, src)
}
func (m *ChromeOSDeviceLSE) XXX_Size() int {
	return xxx_messageInfo_ChromeOSDeviceLSE.Size(m)
}
func (m *ChromeOSDeviceLSE) XXX_DiscardUnknown() {
	xxx_messageInfo_ChromeOSDeviceLSE.DiscardUnknown(m)
}

var xxx_messageInfo_ChromeOSDeviceLSE proto.InternalMessageInfo

func (m *ChromeOSDeviceLSE) GetConfig() *lab.DeviceUnderTest {
	if m != nil {
		return m.Config
	}
	return nil
}

func (m *ChromeOSDeviceLSE) GetRpmInterface() *RPMInterface {
	if m != nil {
		return m.RpmInterface
	}
	return nil
}

func (m *ChromeOSDeviceLSE) GetNetworkDeviceInterface() *SwitchInterface {
	if m != nil {
		return m.NetworkDeviceInterface
	}
	return nil
}

type ChromeOSServerLSE struct {
	// The vlan that this server is going to serve
	ServedNetwork        *VlanID  `protobuf:"bytes,1,opt,name=served_network,json=servedNetwork,proto3" json:"served_network,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ChromeOSServerLSE) Reset()         { *m = ChromeOSServerLSE{} }
func (m *ChromeOSServerLSE) String() string { return proto.CompactTextString(m) }
func (*ChromeOSServerLSE) ProtoMessage()    {}
func (*ChromeOSServerLSE) Descriptor() ([]byte, []int) {
	return fileDescriptor_6f342dc1a43117d0, []int{6}
}

func (m *ChromeOSServerLSE) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChromeOSServerLSE.Unmarshal(m, b)
}
func (m *ChromeOSServerLSE) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChromeOSServerLSE.Marshal(b, m, deterministic)
}
func (m *ChromeOSServerLSE) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChromeOSServerLSE.Merge(m, src)
}
func (m *ChromeOSServerLSE) XXX_Size() int {
	return xxx_messageInfo_ChromeOSServerLSE.Size(m)
}
func (m *ChromeOSServerLSE) XXX_DiscardUnknown() {
	xxx_messageInfo_ChromeOSServerLSE.DiscardUnknown(m)
}

var xxx_messageInfo_ChromeOSServerLSE proto.InternalMessageInfo

func (m *ChromeOSServerLSE) GetServedNetwork() *VlanID {
	if m != nil {
		return m.ServedNetwork
	}
	return nil
}

func init() {
	proto.RegisterType((*MachineLSE)(nil), "fleet.MachineLSE")
	proto.RegisterType((*ChromeMachineLSE)(nil), "fleet.ChromeMachineLSE")
	proto.RegisterType((*VM)(nil), "fleet.VM")
	proto.RegisterType((*VMID)(nil), "fleet.VMID")
	proto.RegisterType((*ChromeOSMachineLSE)(nil), "fleet.ChromeOSMachineLSE")
	proto.RegisterType((*ChromeOSDeviceLSE)(nil), "fleet.ChromeOSDeviceLSE")
	proto.RegisterType((*ChromeOSServerLSE)(nil), "fleet.ChromeOSServerLSE")
}

func init() { proto.RegisterFile("machine_lse.proto", fileDescriptor_6f342dc1a43117d0) }

var fileDescriptor_6f342dc1a43117d0 = []byte{
	// 617 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x94, 0xdf, 0x6e, 0xd3, 0x3e,
	0x14, 0xc7, 0x7f, 0x6d, 0xd7, 0x69, 0x3d, 0xed, 0xa6, 0xce, 0x9b, 0xf6, 0x0b, 0xdd, 0x24, 0xa6,
	0x88, 0x8b, 0x5d, 0x4c, 0x89, 0x34, 0xb8, 0x40, 0xe2, 0x02, 0x31, 0x3a, 0x41, 0xc4, 0xca, 0x26,
	0x07, 0x76, 0xc1, 0x4d, 0xe4, 0x26, 0x6e, 0x6b, 0x91, 0xc4, 0x91, 0x9d, 0x66, 0xe2, 0x29, 0x78,
	0x07, 0x9e, 0x84, 0x17, 0xe0, 0x9d, 0x50, 0xfd, 0x27, 0x6d, 0x3a, 0x89, 0xbb, 0x9c, 0xe3, 0xcf,
	0xf9, 0xfa, 0x7c, 0x8f, 0xed, 0xc0, 0x61, 0x46, 0xe2, 0x05, 0xcb, 0x69, 0x94, 0x4a, 0xea, 0x15,
	0x82, 0x97, 0x1c, 0x75, 0x67, 0x29, 0xa5, 0xe5, 0xe8, 0xcd, 0x9c, 0x7b, 0xf1, 0x42, 0xf0, 0x8c,
	0x2d, 0x33, 0x8f, 0x8b, 0xb9, 0x6f, 0x03, 0x2e, 0x7d, 0x96, 0xcf, 0x04, 0xf1, 0x15, 0xee, 0x4b,
	0x11, 0xfb, 0x29, 0x99, 0xfa, 0x09, 0xad, 0x58, 0x6c, 0x34, 0x46, 0x47, 0x0a, 0xa6, 0x51, 0xcc,
	0xf3, 0x19, 0x9b, 0x9b, 0xe4, 0x20, 0x95, 0x34, 0x62, 0x89, 0x45, 0x56, 0x91, 0xfa, 0x2c, 0x7f,
	0x14, 0xb6, 0x6e, 0x98, 0xd3, 0xf2, 0x91, 0x8b, 0xef, 0x6b, 0xec, 0xb0, 0xa0, 0x82, 0x15, 0x0b,
	0x2a, 0x48, 0x2a, 0x2d, 0x64, 0x7b, 0xb6, 0x90, 0xfb, 0xbb, 0x0d, 0x30, 0xd1, 0xc9, 0xdb, 0xf0,
	0x06, 0xbd, 0x80, 0x36, 0x4b, 0x9c, 0xd6, 0x79, 0xeb, 0xa2, 0x7f, 0x75, 0xec, 0x29, 0x3b, 0xde,
	0x2d, 0x99, 0x86, 0xb4, 0x5c, 0x16, 0x37, 0x79, 0x15, 0x8c, 0x71, 0x9b, 0x25, 0xe8, 0x2d, 0x0c,
	0xea, 0xed, 0x23, 0x96, 0x38, 0x6d, 0xc5, 0x9f, 0x19, 0x7e, 0x2d, 0x77, 0x6f, 0xa1, 0x60, 0x8c,
	0xfb, 0x75, 0x45, 0x90, 0xa0, 0x0f, 0x80, 0x8c, 0xcd, 0x8d, 0x21, 0x3a, 0x1d, 0x25, 0xf3, 0xbf,
	0x91, 0x79, 0xaf, 0x80, 0xb5, 0xd8, 0xc7, 0xff, 0xf0, 0x30, 0x6e, 0xe4, 0x24, 0x45, 0x13, 0x38,
	0xd6, 0x39, 0x2e, 0x1b, 0x52, 0x3b, 0x4a, 0xea, 0x59, 0x43, 0xea, 0x2e, 0x6c, 0x88, 0x21, 0x5b,
	0xb8, 0x21, 0x77, 0x09, 0x7b, 0x46, 0x45, 0x3a, 0xdd, 0xf3, 0xce, 0x45, 0xff, 0x6a, 0xd8, 0x34,
	0x15, 0x8c, 0x71, 0x4d, 0x5c, 0x77, 0xa1, 0x93, 0x4a, 0xea, 0x7e, 0x82, 0xe1, 0x76, 0xaf, 0x68,
	0x04, 0x7b, 0x0b, 0x2e, 0xcb, 0x9c, 0x64, 0x54, 0x4d, 0xb3, 0x87, 0xeb, 0x18, 0x9d, 0x42, 0xa7,
	0xca, 0xa4, 0xd3, 0x56, 0xfa, 0x3d, 0xa3, 0xff, 0x30, 0xc1, 0xab, 0xac, 0xfb, 0xb3, 0x05, 0xed,
	0x87, 0x09, 0x3a, 0xdd, 0x38, 0x87, 0x7e, 0x8d, 0x98, 0xf1, 0xfb, 0x00, 0x5c, 0x46, 0x15, 0x15,
	0x92, 0xf1, 0xdc, 0x0c, 0xdf, 0xf6, 0x79, 0x17, 0x3e, 0xe8, 0x3c, 0xee, 0x71, 0x69, 0x3e, 0xd1,
	0x73, 0xe8, 0x67, 0x24, 0x8e, 0x48, 0x92, 0x08, 0x2a, 0xa5, 0x9a, 0x73, 0x0f, 0x43, 0x46, 0xe2,
	0x77, 0x3a, 0xd3, 0x68, 0x77, 0xa7, 0xd9, 0xae, 0x7b, 0x06, 0x3b, 0xab, 0x9d, 0xd1, 0x31, 0x74,
	0x2b, 0x92, 0x2e, 0xad, 0x1f, 0x1d, 0xb8, 0xbf, 0x5a, 0x80, 0x9e, 0x8e, 0xf7, 0x9f, 0xfe, 0x2f,
	0xa1, 0x93, 0x2c, 0x4b, 0xd3, 0xb7, 0xb3, 0x75, 0x44, 0x63, 0xf5, 0x16, 0xf4, 0x09, 0xad, 0x30,
	0x74, 0x05, 0xbb, 0x92, 0x8a, 0x8a, 0x0a, 0x73, 0x3d, 0xb6, 0x0b, 0x42, 0xb5, 0xa8, 0x0b, 0x0c,
	0x79, 0x7d, 0x00, 0x83, 0xfa, 0x56, 0xac, 0x4e, 0xe8, 0x4f, 0x0b, 0x0e, 0x9f, 0x6c, 0x80, 0x2e,
	0x61, 0x57, 0x3f, 0xb2, 0xfa, 0xbe, 0xa7, 0x64, 0xea, 0xe9, 0xf5, 0xaf, 0x79, 0x42, 0xc5, 0x17,
	0x2a, 0x4b, 0x6c, 0x18, 0xf4, 0x1a, 0xf6, 0x45, 0x91, 0x45, 0x2c, 0x2f, 0xa9, 0x98, 0x91, 0x98,
	0x9a, 0xfe, 0x8f, 0x4c, 0x3b, 0xf8, 0x7e, 0x12, 0xd8, 0x25, 0x3c, 0x10, 0x45, 0x56, 0x47, 0xe8,
	0x1e, 0x1c, 0xfb, 0x36, 0xf5, 0x4b, 0xdf, 0x10, 0xd1, 0x9e, 0x4e, 0x8c, 0x48, 0xf8, 0xc8, 0xca,
	0x78, 0xb1, 0xd6, 0x39, 0x31, 0x75, 0xba, 0xa7, 0x3a, 0xef, 0x06, 0x6b, 0x3b, 0xb5, 0x7d, 0xf4,
	0x0a, 0x0e, 0x94, 0xfd, 0x24, 0x32, 0x55, 0xc6, 0xd6, 0xbe, 0xbd, 0x3e, 0x29, 0xc9, 0x83, 0x31,
	0xde, 0xd7, 0xd0, 0x67, 0xcd, 0x5c, 0x9f, 0x7d, 0x1b, 0xe9, 0x5f, 0x52, 0xca, 0xa6, 0xd2, 0x57,
	0xa4, 0xfe, 0x3b, 0x49, 0x7f, 0xce, 0xa7, 0xbb, 0xea, 0xf3, 0xe5, 0xdf, 0x00, 0x00, 0x00, 0xff,
	0xff, 0xd1, 0x64, 0xd8, 0x9f, 0xec, 0x04, 0x00, 0x00,
}
