// Code generated by protoc-gen-go. DO NOT EDIT.
// source: infra/appengine/arquebus/app/config/config.proto

package config

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	duration "github.com/golang/protobuf/ptypes/duration"
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

type Oncall_Position int32

const (
	Oncall_UNSET     Oncall_Position = 0
	Oncall_PRIMARY   Oncall_Position = 1
	Oncall_SECONDARY Oncall_Position = 2
)

var Oncall_Position_name = map[int32]string{
	0: "UNSET",
	1: "PRIMARY",
	2: "SECONDARY",
}

var Oncall_Position_value = map[string]int32{
	"UNSET":     0,
	"PRIMARY":   1,
	"SECONDARY": 2,
}

func (x Oncall_Position) String() string {
	return proto.EnumName(Oncall_Position_name, int32(x))
}

func (Oncall_Position) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_421d741a02045ab0, []int{2, 0}
}

// Config is the service-wide configuration data for Arquebus
type Config struct {
	// AccessGroup is the luci-auth group who has access to admin pages and
	// APIs.
	AccessGroup string `protobuf:"bytes,1,opt,name=access_group,json=accessGroup,proto3" json:"access_group,omitempty"`
	// The endpoint for Monorail APIs.
	MonorailHostname string `protobuf:"bytes,2,opt,name=monorail_hostname,json=monorailHostname,proto3" json:"monorail_hostname,omitempty"`
	// A list of Assigner config(s).
	Assigners []*Assigner `protobuf:"bytes,3,rep,name=assigners,proto3" json:"assigners,omitempty"`
	// The endpoint for RotaNG APIs.
	RotangHostname       string   `protobuf:"bytes,4,opt,name=rotang_hostname,json=rotangHostname,proto3" json:"rotang_hostname,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Config) Reset()         { *m = Config{} }
func (m *Config) String() string { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()    {}
func (*Config) Descriptor() ([]byte, []int) {
	return fileDescriptor_421d741a02045ab0, []int{0}
}

func (m *Config) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Config.Unmarshal(m, b)
}
func (m *Config) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Config.Marshal(b, m, deterministic)
}
func (m *Config) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Config.Merge(m, src)
}
func (m *Config) XXX_Size() int {
	return xxx_messageInfo_Config.Size(m)
}
func (m *Config) XXX_DiscardUnknown() {
	xxx_messageInfo_Config.DiscardUnknown(m)
}

var xxx_messageInfo_Config proto.InternalMessageInfo

func (m *Config) GetAccessGroup() string {
	if m != nil {
		return m.AccessGroup
	}
	return ""
}

func (m *Config) GetMonorailHostname() string {
	if m != nil {
		return m.MonorailHostname
	}
	return ""
}

func (m *Config) GetAssigners() []*Assigner {
	if m != nil {
		return m.Assigners
	}
	return nil
}

func (m *Config) GetRotangHostname() string {
	if m != nil {
		return m.RotangHostname
	}
	return ""
}

// IssueQuery describes the issue query to be used for searching unassigned
// issues in Monorail.
type IssueQuery struct {
	// Free-form text query.
	Q string `protobuf:"bytes,1,opt,name=q,proto3" json:"q,omitempty"`
	// String name of the projects to search issues for, e.g. "chromium".
	ProjectNames         []string `protobuf:"bytes,2,rep,name=project_names,json=projectNames,proto3" json:"project_names,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *IssueQuery) Reset()         { *m = IssueQuery{} }
func (m *IssueQuery) String() string { return proto.CompactTextString(m) }
func (*IssueQuery) ProtoMessage()    {}
func (*IssueQuery) Descriptor() ([]byte, []int) {
	return fileDescriptor_421d741a02045ab0, []int{1}
}

func (m *IssueQuery) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_IssueQuery.Unmarshal(m, b)
}
func (m *IssueQuery) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_IssueQuery.Marshal(b, m, deterministic)
}
func (m *IssueQuery) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IssueQuery.Merge(m, src)
}
func (m *IssueQuery) XXX_Size() int {
	return xxx_messageInfo_IssueQuery.Size(m)
}
func (m *IssueQuery) XXX_DiscardUnknown() {
	xxx_messageInfo_IssueQuery.DiscardUnknown(m)
}

var xxx_messageInfo_IssueQuery proto.InternalMessageInfo

func (m *IssueQuery) GetQ() string {
	if m != nil {
		return m.Q
	}
	return ""
}

func (m *IssueQuery) GetProjectNames() []string {
	if m != nil {
		return m.ProjectNames
	}
	return nil
}

// Oncall represents a rotation shift modelled in RotaNG.
type Oncall struct {
	// The name of a rotation.
	Rotation string `protobuf:"bytes,1,opt,name=rotation,proto3" json:"rotation,omitempty"`
	// The oncall position in the shift.
	Position             Oncall_Position `protobuf:"varint,2,opt,name=position,proto3,enum=arquebus.config.Oncall_Position" json:"position,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *Oncall) Reset()         { *m = Oncall{} }
func (m *Oncall) String() string { return proto.CompactTextString(m) }
func (*Oncall) ProtoMessage()    {}
func (*Oncall) Descriptor() ([]byte, []int) {
	return fileDescriptor_421d741a02045ab0, []int{2}
}

func (m *Oncall) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Oncall.Unmarshal(m, b)
}
func (m *Oncall) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Oncall.Marshal(b, m, deterministic)
}
func (m *Oncall) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Oncall.Merge(m, src)
}
func (m *Oncall) XXX_Size() int {
	return xxx_messageInfo_Oncall.Size(m)
}
func (m *Oncall) XXX_DiscardUnknown() {
	xxx_messageInfo_Oncall.DiscardUnknown(m)
}

var xxx_messageInfo_Oncall proto.InternalMessageInfo

func (m *Oncall) GetRotation() string {
	if m != nil {
		return m.Rotation
	}
	return ""
}

func (m *Oncall) GetPosition() Oncall_Position {
	if m != nil {
		return m.Position
	}
	return Oncall_UNSET
}

// UserSource represents a single source to find a valid Monorail user to whom
// Arquebus will assign or cc issues found.
type UserSource struct {
	// Types that are valid to be assigned to From:
	//	*UserSource_Oncall
	//	*UserSource_Email
	From                 isUserSource_From `protobuf_oneof:"from"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *UserSource) Reset()         { *m = UserSource{} }
func (m *UserSource) String() string { return proto.CompactTextString(m) }
func (*UserSource) ProtoMessage()    {}
func (*UserSource) Descriptor() ([]byte, []int) {
	return fileDescriptor_421d741a02045ab0, []int{3}
}

func (m *UserSource) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UserSource.Unmarshal(m, b)
}
func (m *UserSource) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UserSource.Marshal(b, m, deterministic)
}
func (m *UserSource) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserSource.Merge(m, src)
}
func (m *UserSource) XXX_Size() int {
	return xxx_messageInfo_UserSource.Size(m)
}
func (m *UserSource) XXX_DiscardUnknown() {
	xxx_messageInfo_UserSource.DiscardUnknown(m)
}

var xxx_messageInfo_UserSource proto.InternalMessageInfo

type isUserSource_From interface {
	isUserSource_From()
}

type UserSource_Oncall struct {
	Oncall *Oncall `protobuf:"bytes,1,opt,name=oncall,proto3,oneof"`
}

type UserSource_Email struct {
	Email string `protobuf:"bytes,2,opt,name=email,proto3,oneof"`
}

func (*UserSource_Oncall) isUserSource_From() {}

func (*UserSource_Email) isUserSource_From() {}

func (m *UserSource) GetFrom() isUserSource_From {
	if m != nil {
		return m.From
	}
	return nil
}

func (m *UserSource) GetOncall() *Oncall {
	if x, ok := m.GetFrom().(*UserSource_Oncall); ok {
		return x.Oncall
	}
	return nil
}

func (m *UserSource) GetEmail() string {
	if x, ok := m.GetFrom().(*UserSource_Email); ok {
		return x.Email
	}
	return ""
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*UserSource) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*UserSource_Oncall)(nil),
		(*UserSource_Email)(nil),
	}
}

// Assigner contains specifications for an Assigner job.
type Assigner struct {
	// The unique ID of the Assigner.
	//
	// This value will be used in URLs of UI, so keep it short. Note that
	// only lowercase alphabet letters and numbers are allowed. A hyphen may
	// be placed between letters and numbers.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// An email list of the owners of the Assigner.
	Owners []string `protobuf:"bytes,2,rep,name=owners,proto3" json:"owners,omitempty"`
	// The duration between the start of an Assigner run and the next one.
	//
	// This value should be at least a minute long.
	Interval *duration.Duration `protobuf:"bytes,3,opt,name=interval,proto3" json:"interval,omitempty"`
	// IssueQuery describes the search criteria to look for issues to assign.
	IssueQuery *IssueQuery `protobuf:"bytes,4,opt,name=issue_query,json=issueQuery,proto3" json:"issue_query,omitempty"`
	// If multiple values are specified in assignees, Arquebus iterates the list
	// in the order until it finds a currently available assignee. Note that
	// Monorail users are always assumed to be available.
	Assignees []*UserSource `protobuf:"bytes,6,rep,name=assignees,proto3" json:"assignees,omitempty"`
	// If multiple values are specified in ccs, all the available roations and
	// users are added to the CC of searched issues.
	Ccs []*UserSource `protobuf:"bytes,7,rep,name=ccs,proto3" json:"ccs,omitempty"`
	// If DryRun is set, Assigner doesn't update the found issues.
	DryRun bool `protobuf:"varint,8,opt,name=dry_run,json=dryRun,proto3" json:"dry_run,omitempty"`
	// The description shown on UI.
	Description string `protobuf:"bytes,9,opt,name=description,proto3" json:"description,omitempty"`
	// Comment is an additional message that is added to the body of the issue
	// comment that is posted when an issue gets updated.
	Comment              string   `protobuf:"bytes,10,opt,name=comment,proto3" json:"comment,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Assigner) Reset()         { *m = Assigner{} }
func (m *Assigner) String() string { return proto.CompactTextString(m) }
func (*Assigner) ProtoMessage()    {}
func (*Assigner) Descriptor() ([]byte, []int) {
	return fileDescriptor_421d741a02045ab0, []int{4}
}

func (m *Assigner) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Assigner.Unmarshal(m, b)
}
func (m *Assigner) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Assigner.Marshal(b, m, deterministic)
}
func (m *Assigner) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Assigner.Merge(m, src)
}
func (m *Assigner) XXX_Size() int {
	return xxx_messageInfo_Assigner.Size(m)
}
func (m *Assigner) XXX_DiscardUnknown() {
	xxx_messageInfo_Assigner.DiscardUnknown(m)
}

var xxx_messageInfo_Assigner proto.InternalMessageInfo

func (m *Assigner) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Assigner) GetOwners() []string {
	if m != nil {
		return m.Owners
	}
	return nil
}

func (m *Assigner) GetInterval() *duration.Duration {
	if m != nil {
		return m.Interval
	}
	return nil
}

func (m *Assigner) GetIssueQuery() *IssueQuery {
	if m != nil {
		return m.IssueQuery
	}
	return nil
}

func (m *Assigner) GetAssignees() []*UserSource {
	if m != nil {
		return m.Assignees
	}
	return nil
}

func (m *Assigner) GetCcs() []*UserSource {
	if m != nil {
		return m.Ccs
	}
	return nil
}

func (m *Assigner) GetDryRun() bool {
	if m != nil {
		return m.DryRun
	}
	return false
}

func (m *Assigner) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *Assigner) GetComment() string {
	if m != nil {
		return m.Comment
	}
	return ""
}

func init() {
	proto.RegisterEnum("arquebus.config.Oncall_Position", Oncall_Position_name, Oncall_Position_value)
	proto.RegisterType((*Config)(nil), "arquebus.config.Config")
	proto.RegisterType((*IssueQuery)(nil), "arquebus.config.IssueQuery")
	proto.RegisterType((*Oncall)(nil), "arquebus.config.Oncall")
	proto.RegisterType((*UserSource)(nil), "arquebus.config.UserSource")
	proto.RegisterType((*Assigner)(nil), "arquebus.config.Assigner")
}

func init() {
	proto.RegisterFile("infra/appengine/arquebus/app/config/config.proto", fileDescriptor_421d741a02045ab0)
}

var fileDescriptor_421d741a02045ab0 = []byte{
	// 553 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x53, 0xdd, 0x6e, 0x94, 0x40,
	0x18, 0x2d, 0x6c, 0xcb, 0xc2, 0x47, 0xff, 0x9c, 0x8b, 0x96, 0xd6, 0xc4, 0x20, 0x5e, 0xb8, 0x89,
	0x91, 0xb5, 0x6b, 0x8c, 0x31, 0x69, 0x62, 0xfa, 0x17, 0xdb, 0x0b, 0xdb, 0x3a, 0xb5, 0x17, 0x7a,
	0x43, 0x66, 0x61, 0x16, 0xc7, 0xc0, 0x0c, 0x3b, 0x03, 0x9a, 0x7d, 0x10, 0x5f, 0xc6, 0x07, 0xf2,
	0x39, 0x0c, 0x03, 0x2c, 0xc6, 0x8d, 0xf1, 0x8a, 0x7c, 0x67, 0xce, 0xf7, 0x7b, 0x0e, 0xf0, 0x82,
	0xf1, 0x99, 0x24, 0x63, 0x52, 0x14, 0x94, 0xa7, 0x8c, 0xd3, 0x31, 0x91, 0xf3, 0x8a, 0x4e, 0x2b,
	0x55, 0x43, 0xe3, 0x58, 0xf0, 0x19, 0x4b, 0xdb, 0x4f, 0x58, 0x48, 0x51, 0x0a, 0xb4, 0xd3, 0x31,
	0xc2, 0x06, 0x3e, 0x7c, 0x94, 0x0a, 0x91, 0x66, 0x74, 0xac, 0x9f, 0xa7, 0xd5, 0x6c, 0x9c, 0x54,
	0x92, 0x94, 0x4c, 0xf0, 0x26, 0x21, 0xf8, 0x69, 0x80, 0x75, 0xa6, 0xa9, 0xe8, 0x31, 0x6c, 0x92,
	0x38, 0xa6, 0x4a, 0x45, 0xa9, 0x14, 0x55, 0xe1, 0x19, 0xbe, 0x31, 0x72, 0xb0, 0xdb, 0x60, 0xef,
	0x6a, 0x08, 0x3d, 0x83, 0x07, 0xb9, 0xe0, 0x42, 0x12, 0x96, 0x45, 0x5f, 0x84, 0x2a, 0x39, 0xc9,
	0xa9, 0x67, 0x6a, 0xde, 0x6e, 0xf7, 0x70, 0xd9, 0xe2, 0xe8, 0x35, 0x38, 0x44, 0x29, 0x96, 0x72,
	0x2a, 0x95, 0x37, 0xf0, 0x07, 0x23, 0x77, 0x72, 0x10, 0xfe, 0x35, 0x5f, 0x78, 0xd2, 0x32, 0x70,
	0xcf, 0x45, 0x4f, 0x61, 0x47, 0x8a, 0x92, 0xf0, 0xb4, 0xef, 0xb1, 0xae, 0x7b, 0x6c, 0x37, 0x70,
	0xd7, 0x21, 0x78, 0x0b, 0x70, 0xa5, 0x54, 0x45, 0x3f, 0x54, 0x54, 0x2e, 0xd0, 0x26, 0x18, 0xf3,
	0x76, 0x68, 0x63, 0x8e, 0x9e, 0xc0, 0x56, 0x21, 0xc5, 0x57, 0x1a, 0x97, 0x51, 0xcd, 0x55, 0x9e,
	0xe9, 0x0f, 0x46, 0x0e, 0xde, 0x6c, 0xc1, 0xeb, 0x1a, 0x0b, 0x7e, 0x18, 0x60, 0xdd, 0xf0, 0x98,
	0x64, 0x19, 0x3a, 0x04, 0xbb, 0xae, 0x5e, 0x9f, 0xa6, 0x2d, 0xb2, 0x8c, 0xd1, 0x31, 0xd8, 0x85,
	0x50, 0x4c, 0xbf, 0xd5, 0xdb, 0x6e, 0x4f, 0xfc, 0x95, 0x45, 0x9a, 0x32, 0xe1, 0x6d, 0xcb, 0xc3,
	0xcb, 0x8c, 0xe0, 0x08, 0xec, 0x0e, 0x45, 0x0e, 0x6c, 0xdc, 0x5f, 0xdf, 0x5d, 0x7c, 0xdc, 0x5d,
	0x43, 0x2e, 0x0c, 0x6f, 0xf1, 0xd5, 0xfb, 0x13, 0xfc, 0x69, 0xd7, 0x40, 0x5b, 0xe0, 0xdc, 0x5d,
	0x9c, 0xdd, 0x5c, 0x9f, 0xd7, 0xa1, 0x19, 0x44, 0x00, 0xf7, 0x8a, 0xca, 0x3b, 0x51, 0xc9, 0x98,
	0xa2, 0x23, 0xb0, 0x84, 0xae, 0xae, 0x07, 0x73, 0x27, 0xfb, 0xff, 0x68, 0x7e, 0xb9, 0x86, 0x5b,
	0x22, 0xda, 0x83, 0x0d, 0x9a, 0x13, 0x96, 0x35, 0xe2, 0x5c, 0xae, 0xe1, 0x26, 0x3c, 0xb5, 0x60,
	0x7d, 0x26, 0x45, 0x1e, 0xfc, 0x32, 0xc1, 0xee, 0x4e, 0x8f, 0xb6, 0xc1, 0x64, 0x49, 0xbb, 0xb4,
	0xc9, 0x12, 0xb4, 0x07, 0x96, 0xf8, 0xae, 0x55, 0x6b, 0x6e, 0xd6, 0x46, 0xe8, 0x15, 0xd8, 0x8c,
	0x97, 0x54, 0x7e, 0x23, 0x99, 0x37, 0xd0, 0x93, 0x1c, 0x84, 0x8d, 0xbd, 0xc2, 0xce, 0x5e, 0xe1,
	0x79, 0x6b, 0x2f, 0xbc, 0xa4, 0xa2, 0x63, 0x70, 0x59, 0xad, 0x52, 0x34, 0xaf, 0x65, 0xd2, 0x52,
	0xba, 0x93, 0x87, 0x2b, 0x3b, 0xf4, 0x4a, 0x62, 0x60, 0xbd, 0xaa, 0x6f, 0x96, 0x2e, 0xa2, 0xca,
	0xb3, 0xb4, 0x8b, 0x56, 0x73, 0xfb, 0x63, 0xe1, 0x9e, 0x8d, 0x9e, 0xc3, 0x20, 0x8e, 0x95, 0x37,
	0xfc, 0x7f, 0x52, 0xcd, 0x43, 0xfb, 0x30, 0x4c, 0xe4, 0x22, 0x92, 0x15, 0xf7, 0x6c, 0xdf, 0x18,
	0xd9, 0xd8, 0x4a, 0xe4, 0x02, 0x57, 0x1c, 0xf9, 0xe0, 0x26, 0x54, 0xc5, 0x92, 0x15, 0xda, 0x01,
	0x4e, 0xf3, 0x5f, 0xfc, 0x01, 0x21, 0x0f, 0x86, 0xb1, 0xc8, 0x73, 0xca, 0x4b, 0x0f, 0xf4, 0x6b,
	0x17, 0x9e, 0xda, 0x9f, 0xad, 0xa6, 0xdd, 0xd4, 0xd2, 0x37, 0x7a, 0xf9, 0x3b, 0x00, 0x00, 0xff,
	0xff, 0x63, 0x29, 0xf4, 0xef, 0xd5, 0x03, 0x00, 0x00,
}
