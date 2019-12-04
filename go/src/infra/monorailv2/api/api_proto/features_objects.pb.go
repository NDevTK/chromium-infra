// Code generated by protoc-gen-go. DO NOT EDIT.
// source: api/api_proto/features_objects.proto

package monorail

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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// TODO(jojwang): add editors and followers
// Next available tag: 8
type Hotlist struct {
	OwnerRef             *UserRef   `protobuf:"bytes,1,opt,name=owner_ref,json=ownerRef,proto3" json:"owner_ref,omitempty"`
	EditorRefs           []*UserRef `protobuf:"bytes,5,rep,name=editor_refs,json=editorRefs,proto3" json:"editor_refs,omitempty"`
	FollowerRefs         []*UserRef `protobuf:"bytes,6,rep,name=follower_refs,json=followerRefs,proto3" json:"follower_refs,omitempty"`
	Name                 string     `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Summary              string     `protobuf:"bytes,3,opt,name=summary,proto3" json:"summary,omitempty"`
	Description          string     `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	DefaultColSpec       string     `protobuf:"bytes,7,opt,name=default_col_spec,json=defaultColSpec,proto3" json:"default_col_spec,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *Hotlist) Reset()         { *m = Hotlist{} }
func (m *Hotlist) String() string { return proto.CompactTextString(m) }
func (*Hotlist) ProtoMessage()    {}
func (*Hotlist) Descriptor() ([]byte, []int) {
	return fileDescriptor_806b6b78af767289, []int{0}
}

func (m *Hotlist) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Hotlist.Unmarshal(m, b)
}
func (m *Hotlist) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Hotlist.Marshal(b, m, deterministic)
}
func (m *Hotlist) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Hotlist.Merge(m, src)
}
func (m *Hotlist) XXX_Size() int {
	return xxx_messageInfo_Hotlist.Size(m)
}
func (m *Hotlist) XXX_DiscardUnknown() {
	xxx_messageInfo_Hotlist.DiscardUnknown(m)
}

var xxx_messageInfo_Hotlist proto.InternalMessageInfo

func (m *Hotlist) GetOwnerRef() *UserRef {
	if m != nil {
		return m.OwnerRef
	}
	return nil
}

func (m *Hotlist) GetEditorRefs() []*UserRef {
	if m != nil {
		return m.EditorRefs
	}
	return nil
}

func (m *Hotlist) GetFollowerRefs() []*UserRef {
	if m != nil {
		return m.FollowerRefs
	}
	return nil
}

func (m *Hotlist) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Hotlist) GetSummary() string {
	if m != nil {
		return m.Summary
	}
	return ""
}

func (m *Hotlist) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *Hotlist) GetDefaultColSpec() string {
	if m != nil {
		return m.DefaultColSpec
	}
	return ""
}

// Next available tag: 6
type HotlistItem struct {
	Issue                *Issue   `protobuf:"bytes,1,opt,name=issue,proto3" json:"issue,omitempty"`
	Rank                 uint32   `protobuf:"varint,2,opt,name=rank,proto3" json:"rank,omitempty"`
	AdderRef             *UserRef `protobuf:"bytes,3,opt,name=adder_ref,json=adderRef,proto3" json:"adder_ref,omitempty"`
	AddedTimestamp       uint32   `protobuf:"varint,4,opt,name=added_timestamp,json=addedTimestamp,proto3" json:"added_timestamp,omitempty"`
	Note                 string   `protobuf:"bytes,5,opt,name=note,proto3" json:"note,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HotlistItem) Reset()         { *m = HotlistItem{} }
func (m *HotlistItem) String() string { return proto.CompactTextString(m) }
func (*HotlistItem) ProtoMessage()    {}
func (*HotlistItem) Descriptor() ([]byte, []int) {
	return fileDescriptor_806b6b78af767289, []int{1}
}

func (m *HotlistItem) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HotlistItem.Unmarshal(m, b)
}
func (m *HotlistItem) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HotlistItem.Marshal(b, m, deterministic)
}
func (m *HotlistItem) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HotlistItem.Merge(m, src)
}
func (m *HotlistItem) XXX_Size() int {
	return xxx_messageInfo_HotlistItem.Size(m)
}
func (m *HotlistItem) XXX_DiscardUnknown() {
	xxx_messageInfo_HotlistItem.DiscardUnknown(m)
}

var xxx_messageInfo_HotlistItem proto.InternalMessageInfo

func (m *HotlistItem) GetIssue() *Issue {
	if m != nil {
		return m.Issue
	}
	return nil
}

func (m *HotlistItem) GetRank() uint32 {
	if m != nil {
		return m.Rank
	}
	return 0
}

func (m *HotlistItem) GetAdderRef() *UserRef {
	if m != nil {
		return m.AdderRef
	}
	return nil
}

func (m *HotlistItem) GetAddedTimestamp() uint32 {
	if m != nil {
		return m.AddedTimestamp
	}
	return 0
}

func (m *HotlistItem) GetNote() string {
	if m != nil {
		return m.Note
	}
	return ""
}

// Next available tag: 5
type HotlistPeopleDelta struct {
	NewOwnerRef          *UserRef   `protobuf:"bytes,1,opt,name=new_owner_ref,json=newOwnerRef,proto3" json:"new_owner_ref,omitempty"`
	AddEditorRefs        []*UserRef `protobuf:"bytes,2,rep,name=add_editor_refs,json=addEditorRefs,proto3" json:"add_editor_refs,omitempty"`
	AddFollowerRefs      []*UserRef `protobuf:"bytes,3,rep,name=add_follower_refs,json=addFollowerRefs,proto3" json:"add_follower_refs,omitempty"`
	RemoveUserRefs       []*UserRef `protobuf:"bytes,4,rep,name=remove_user_refs,json=removeUserRefs,proto3" json:"remove_user_refs,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *HotlistPeopleDelta) Reset()         { *m = HotlistPeopleDelta{} }
func (m *HotlistPeopleDelta) String() string { return proto.CompactTextString(m) }
func (*HotlistPeopleDelta) ProtoMessage()    {}
func (*HotlistPeopleDelta) Descriptor() ([]byte, []int) {
	return fileDescriptor_806b6b78af767289, []int{2}
}

func (m *HotlistPeopleDelta) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HotlistPeopleDelta.Unmarshal(m, b)
}
func (m *HotlistPeopleDelta) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HotlistPeopleDelta.Marshal(b, m, deterministic)
}
func (m *HotlistPeopleDelta) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HotlistPeopleDelta.Merge(m, src)
}
func (m *HotlistPeopleDelta) XXX_Size() int {
	return xxx_messageInfo_HotlistPeopleDelta.Size(m)
}
func (m *HotlistPeopleDelta) XXX_DiscardUnknown() {
	xxx_messageInfo_HotlistPeopleDelta.DiscardUnknown(m)
}

var xxx_messageInfo_HotlistPeopleDelta proto.InternalMessageInfo

func (m *HotlistPeopleDelta) GetNewOwnerRef() *UserRef {
	if m != nil {
		return m.NewOwnerRef
	}
	return nil
}

func (m *HotlistPeopleDelta) GetAddEditorRefs() []*UserRef {
	if m != nil {
		return m.AddEditorRefs
	}
	return nil
}

func (m *HotlistPeopleDelta) GetAddFollowerRefs() []*UserRef {
	if m != nil {
		return m.AddFollowerRefs
	}
	return nil
}

func (m *HotlistPeopleDelta) GetRemoveUserRefs() []*UserRef {
	if m != nil {
		return m.RemoveUserRefs
	}
	return nil
}

func init() {
	proto.RegisterType((*Hotlist)(nil), "monorail.Hotlist")
	proto.RegisterType((*HotlistItem)(nil), "monorail.HotlistItem")
	proto.RegisterType((*HotlistPeopleDelta)(nil), "monorail.HotlistPeopleDelta")
}

func init() {
	proto.RegisterFile("api/api_proto/features_objects.proto", fileDescriptor_806b6b78af767289)
}

var fileDescriptor_806b6b78af767289 = []byte{
	// 426 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x52, 0x4d, 0x8f, 0xd3, 0x30,
	0x10, 0x55, 0xbf, 0xb6, 0xbb, 0x13, 0xd2, 0xdd, 0xf5, 0x29, 0xea, 0xa9, 0x54, 0x20, 0x7a, 0xea,
	0x4a, 0x8b, 0x40, 0x42, 0x88, 0x13, 0x1f, 0x62, 0x4f, 0xa0, 0x00, 0x67, 0xcb, 0x9b, 0x4c, 0x24,
	0x83, 0x9d, 0x89, 0x6c, 0x87, 0x8a, 0x2b, 0x7f, 0x84, 0x5f, 0xc1, 0xff, 0x43, 0x99, 0x3a, 0x6a,
	0x2b, 0x11, 0x69, 0x6f, 0xe3, 0x37, 0xef, 0xc5, 0x2f, 0xef, 0x19, 0x9e, 0xa8, 0x46, 0xdf, 0xa8,
	0x46, 0xcb, 0xc6, 0x51, 0xa0, 0x9b, 0x0a, 0x55, 0x68, 0x1d, 0x7a, 0x49, 0xf7, 0xdf, 0xb1, 0x08,
	0x7e, 0xcb, 0xb0, 0x38, 0xb7, 0x54, 0x93, 0x53, 0xda, 0x2c, 0x97, 0xa7, 0xfc, 0x82, 0xac, 0xa5,
	0x7a, 0xcf, 0x5a, 0x3e, 0x3e, 0xdd, 0x69, 0xef, 0x5b, 0x3c, 0xfd, 0xd0, 0xfa, 0xcf, 0x18, 0xe6,
	0x1f, 0x29, 0x18, 0xed, 0x83, 0xd8, 0xc2, 0x05, 0xed, 0x6a, 0x74, 0xd2, 0x61, 0x95, 0x8d, 0x56,
	0xa3, 0x4d, 0x72, 0x7b, 0xbd, 0xed, 0x2f, 0xda, 0x7e, 0xf3, 0xe8, 0x72, 0xac, 0xf2, 0x73, 0xe6,
	0xe4, 0x58, 0x89, 0x5b, 0x48, 0xb0, 0xd4, 0x81, 0x58, 0xe0, 0xb3, 0xd9, 0x6a, 0xf2, 0x7f, 0x05,
	0xec, 0x59, 0x39, 0x56, 0x5e, 0xbc, 0x84, 0xb4, 0x22, 0x63, 0x68, 0x87, 0x51, 0x75, 0x36, 0xa4,
	0x7a, 0xd4, 0xf3, 0x58, 0x27, 0x60, 0x5a, 0x2b, 0x8b, 0xd9, 0x78, 0x35, 0xda, 0x5c, 0xe4, 0x3c,
	0x8b, 0x0c, 0xe6, 0xbe, 0xb5, 0x56, 0xb9, 0x5f, 0xd9, 0x84, 0xe1, 0xfe, 0x28, 0x56, 0x90, 0x94,
	0xe8, 0x0b, 0xa7, 0x9b, 0xa0, 0xa9, 0xce, 0xa6, 0xbc, 0x3d, 0x86, 0xc4, 0x06, 0xae, 0x4a, 0xac,
	0x54, 0x6b, 0x82, 0x2c, 0xc8, 0x48, 0xdf, 0x60, 0x91, 0xcd, 0x99, 0xb6, 0x88, 0xf8, 0x5b, 0x32,
	0x5f, 0x1a, 0x2c, 0xd6, 0x7f, 0x47, 0x90, 0xc4, 0x84, 0xee, 0x02, 0x5a, 0xf1, 0x14, 0x66, 0x1c,
	0x64, 0x4c, 0xe8, 0xf2, 0xe0, 0xfc, 0xae, 0x83, 0xf3, 0xfd, 0xb6, 0x33, 0xec, 0x54, 0xfd, 0x83,
	0x0d, 0xa7, 0x39, 0xcf, 0x5d, 0xc0, 0xaa, 0x2c, 0x63, 0xc0, 0x93, 0xc1, 0x80, 0x99, 0xd3, 0x05,
	0xfc, 0x0c, 0x2e, 0xbb, 0xb9, 0x94, 0x41, 0x5b, 0xf4, 0x41, 0xd9, 0x86, 0x7f, 0x25, 0xcd, 0x17,
	0x0c, 0x7f, 0xed, 0x51, 0x4e, 0x87, 0x02, 0x66, 0xb3, 0x98, 0x0e, 0x05, 0x5c, 0xff, 0x1e, 0x83,
	0x88, 0xbe, 0x3f, 0x23, 0x35, 0x06, 0xdf, 0xa1, 0x09, 0x4a, 0xbc, 0x80, 0xb4, 0xc6, 0x9d, 0x7c,
	0x40, 0xd1, 0x49, 0x8d, 0xbb, 0x4f, 0x7d, 0xd7, 0xaf, 0xd8, 0x8a, 0x3c, 0xee, 0x7b, 0x3c, 0xd4,
	0x5c, 0xaa, 0xca, 0xf2, 0xfd, 0xa1, 0xf2, 0x37, 0x70, 0xdd, 0x49, 0x4f, 0x6b, 0x9f, 0x0c, 0x89,
	0xbb, 0x6b, 0x3e, 0x1c, 0x37, 0xff, 0x1a, 0xae, 0x1c, 0x5a, 0xfa, 0x89, 0xb2, 0xf5, 0xbd, 0x7a,
	0x3a, 0xa4, 0x5e, 0xec, 0xa9, 0xf1, 0xe8, 0xef, 0xcf, 0xf8, 0x95, 0x3f, 0xff, 0x17, 0x00, 0x00,
	0xff, 0xff, 0xd8, 0x52, 0xaf, 0x0d, 0x56, 0x03, 0x00, 0x00,
}
