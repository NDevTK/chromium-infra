// Code generated by protoc-gen-go. DO NOT EDIT.
// source: infra/tricium/api/v1/config.proto

/*
Package tricium is a generated protocol buffer package.

It is generated from these files:
	infra/tricium/api/v1/config.proto
	infra/tricium/api/v1/data.proto
	infra/tricium/api/v1/function.proto
	infra/tricium/api/v1/platform.proto
	infra/tricium/api/v1/tricium.proto

It has these top-level messages:
	ServiceConfig
	ProjectConfig
	RepoDetails
	GerritProject
	GitRepo
	Acl
	Selection
	Config
	Data
	Function
	ConfigDef
	Impl
	Recipe
	Property
	Cmd
	CipdPackage
	Platform
	AnalyzeRequest
	GerritRevision
	GitCommit
	AnalyzeResponse
	ProgressRequest
	ProgressResponse
	FunctionProgress
	ProjectProgressRequest
	ProjectProgressResponse
	RunProgress
	ResultsRequest
	ResultsResponse
	FeedbackRequest
	FeedbackResponse
	ReportNotUsefulRequest
	ReportNotUsefulResponse
*/
package tricium

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Roles relevant to Tricium.
type Acl_Role int32

const (
	// Can read progress/results.
	Acl_READER Acl_Role = 0
	// Can request analysis.
	Acl_REQUESTER Acl_Role = 1
)

var Acl_Role_name = map[int32]string{
	0: "READER",
	1: "REQUESTER",
}
var Acl_Role_value = map[string]int32{
	"READER":    0,
	"REQUESTER": 1,
}

func (x Acl_Role) String() string {
	return proto.EnumName(Acl_Role_name, int32(x))
}
func (Acl_Role) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{5, 0} }

// Tricium service configuration.
//
// Listing supported platforms and analyzers shared between projects connected
// to Tricium.
type ServiceConfig struct {
	// Supported platforms.
	Platforms []*Platform_Details `protobuf:"bytes,1,rep,name=platforms" json:"platforms,omitempty"`
	// Supported data types.
	DataDetails []*Data_TypeDetails `protobuf:"bytes,2,rep,name=data_details,json=dataDetails" json:"data_details,omitempty"`
	// List of shared functions.
	Functions []*Function `protobuf:"bytes,3,rep,name=functions" json:"functions,omitempty"`
	// Base recipe command used for workers implemented as recipes.
	//
	// Specific recipe details for the worker will be added as flags at the
	// end of the argument list.
	RecipeCmd *Cmd `protobuf:"bytes,5,opt,name=recipe_cmd,json=recipeCmd" json:"recipe_cmd,omitempty"`
	// Base recipe packages used for workers implemented as recipes.
	//
	// These packages will be adjusted for the platform in question, by appending
	// platform name details to the end of the package name.
	RecipePackages []*CipdPackage `protobuf:"bytes,6,rep,name=recipe_packages,json=recipePackages" json:"recipe_packages,omitempty"`
	// Swarming server to use for this service instance.
	//
	// This should be a full URL with no trailing slash.
	SwarmingServer string `protobuf:"bytes,7,opt,name=swarming_server,json=swarmingServer" json:"swarming_server,omitempty"`
	// Isolate server to use for this service instance.
	//
	// This should be a full URL with no trailing slash.
	IsolateServer string `protobuf:"bytes,8,opt,name=isolate_server,json=isolateServer" json:"isolate_server,omitempty"`
}

func (m *ServiceConfig) Reset()                    { *m = ServiceConfig{} }
func (m *ServiceConfig) String() string            { return proto.CompactTextString(m) }
func (*ServiceConfig) ProtoMessage()               {}
func (*ServiceConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *ServiceConfig) GetPlatforms() []*Platform_Details {
	if m != nil {
		return m.Platforms
	}
	return nil
}

func (m *ServiceConfig) GetDataDetails() []*Data_TypeDetails {
	if m != nil {
		return m.DataDetails
	}
	return nil
}

func (m *ServiceConfig) GetFunctions() []*Function {
	if m != nil {
		return m.Functions
	}
	return nil
}

func (m *ServiceConfig) GetRecipeCmd() *Cmd {
	if m != nil {
		return m.RecipeCmd
	}
	return nil
}

func (m *ServiceConfig) GetRecipePackages() []*CipdPackage {
	if m != nil {
		return m.RecipePackages
	}
	return nil
}

func (m *ServiceConfig) GetSwarmingServer() string {
	if m != nil {
		return m.SwarmingServer
	}
	return ""
}

func (m *ServiceConfig) GetIsolateServer() string {
	if m != nil {
		return m.IsolateServer
	}
	return ""
}

// Tricium project configuration.
//
// Specifies details needed to connect a project to Tricium.
// Adds project-specific functions and selects shared function
// implementations.
type ProjectConfig struct {
	// Access control rules for the project.
	Acls []*Acl `protobuf:"bytes,2,rep,name=acls" json:"acls,omitempty"`
	// Project-specific function details.
	//
	// This includes project-specific analyzer implementations and full
	// project-specific analyzer specifications.
	Functions []*Function `protobuf:"bytes,3,rep,name=functions" json:"functions,omitempty"`
	// Selection of function implementations to run for this project.
	Selections []*Selection `protobuf:"bytes,4,rep,name=selections" json:"selections,omitempty"`
	// Repositories, including Git and Gerrit details.
	Repos []*RepoDetails `protobuf:"bytes,5,rep,name=repos" json:"repos,omitempty"`
	// General service account for this project.
	// Used for any service interaction, with the exception of swarming.
	ServiceAccount string `protobuf:"bytes,6,opt,name=service_account,json=serviceAccount" json:"service_account,omitempty"`
	// Project-specific swarming service account.
	SwarmingServiceAccount string `protobuf:"bytes,7,opt,name=swarming_service_account,json=swarmingServiceAccount" json:"swarming_service_account,omitempty"`
}

func (m *ProjectConfig) Reset()                    { *m = ProjectConfig{} }
func (m *ProjectConfig) String() string            { return proto.CompactTextString(m) }
func (*ProjectConfig) ProtoMessage()               {}
func (*ProjectConfig) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *ProjectConfig) GetAcls() []*Acl {
	if m != nil {
		return m.Acls
	}
	return nil
}

func (m *ProjectConfig) GetFunctions() []*Function {
	if m != nil {
		return m.Functions
	}
	return nil
}

func (m *ProjectConfig) GetSelections() []*Selection {
	if m != nil {
		return m.Selections
	}
	return nil
}

func (m *ProjectConfig) GetRepos() []*RepoDetails {
	if m != nil {
		return m.Repos
	}
	return nil
}

func (m *ProjectConfig) GetServiceAccount() string {
	if m != nil {
		return m.ServiceAccount
	}
	return ""
}

func (m *ProjectConfig) GetSwarmingServiceAccount() string {
	if m != nil {
		return m.SwarmingServiceAccount
	}
	return ""
}

// Repository details for a project.
// DEPRECATED, see https://crbug.com/824558
type RepoDetails struct {
	// Could be renamed to kind when the above kind is removed.
	//
	// Types that are valid to be assigned to Source:
	//	*RepoDetails_GerritProject
	//	*RepoDetails_GitRepo
	Source isRepoDetails_Source `protobuf_oneof:"source"`
	// Whether to disable reporting results back.
	DisableReporting bool `protobuf:"varint,6,opt,name=disable_reporting,json=disableReporting" json:"disable_reporting,omitempty"`
	// Whitelisted groups.
	//
	// The owner of a change will be checked for membership of a whitelisted
	// group. Absence of this field means all groups are whitelisted.
	//
	// Group names must be known to the Chrome infra auth service,
	// https://chrome-infra-auth.appspot.com. Contact a Chromium trooper
	// if you need to add or modify a group: g.co/bugatrooper.
	WhitelistedGroup []string `protobuf:"bytes,7,rep,name=whitelisted_group,json=whitelistedGroup" json:"whitelisted_group,omitempty"`
}

func (m *RepoDetails) Reset()                    { *m = RepoDetails{} }
func (m *RepoDetails) String() string            { return proto.CompactTextString(m) }
func (*RepoDetails) ProtoMessage()               {}
func (*RepoDetails) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type isRepoDetails_Source interface {
	isRepoDetails_Source()
}

type RepoDetails_GerritProject struct {
	GerritProject *GerritProject `protobuf:"bytes,4,opt,name=gerrit_project,json=gerritProject,oneof"`
}
type RepoDetails_GitRepo struct {
	GitRepo *GitRepo `protobuf:"bytes,5,opt,name=git_repo,json=gitRepo,oneof"`
}

func (*RepoDetails_GerritProject) isRepoDetails_Source() {}
func (*RepoDetails_GitRepo) isRepoDetails_Source()       {}

func (m *RepoDetails) GetSource() isRepoDetails_Source {
	if m != nil {
		return m.Source
	}
	return nil
}

func (m *RepoDetails) GetGerritProject() *GerritProject {
	if x, ok := m.GetSource().(*RepoDetails_GerritProject); ok {
		return x.GerritProject
	}
	return nil
}

func (m *RepoDetails) GetGitRepo() *GitRepo {
	if x, ok := m.GetSource().(*RepoDetails_GitRepo); ok {
		return x.GitRepo
	}
	return nil
}

func (m *RepoDetails) GetDisableReporting() bool {
	if m != nil {
		return m.DisableReporting
	}
	return false
}

func (m *RepoDetails) GetWhitelistedGroup() []string {
	if m != nil {
		return m.WhitelistedGroup
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*RepoDetails) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _RepoDetails_OneofMarshaler, _RepoDetails_OneofUnmarshaler, _RepoDetails_OneofSizer, []interface{}{
		(*RepoDetails_GerritProject)(nil),
		(*RepoDetails_GitRepo)(nil),
	}
}

func _RepoDetails_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*RepoDetails)
	// source
	switch x := m.Source.(type) {
	case *RepoDetails_GerritProject:
		b.EncodeVarint(4<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.GerritProject); err != nil {
			return err
		}
	case *RepoDetails_GitRepo:
		b.EncodeVarint(5<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.GitRepo); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("RepoDetails.Source has unexpected type %T", x)
	}
	return nil
}

func _RepoDetails_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*RepoDetails)
	switch tag {
	case 4: // source.gerrit_project
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(GerritProject)
		err := b.DecodeMessage(msg)
		m.Source = &RepoDetails_GerritProject{msg}
		return true, err
	case 5: // source.git_repo
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(GitRepo)
		err := b.DecodeMessage(msg)
		m.Source = &RepoDetails_GitRepo{msg}
		return true, err
	default:
		return false, nil
	}
}

func _RepoDetails_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*RepoDetails)
	// source
	switch x := m.Source.(type) {
	case *RepoDetails_GerritProject:
		s := proto.Size(x.GerritProject)
		n += proto.SizeVarint(4<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *RepoDetails_GitRepo:
		s := proto.Size(x.GitRepo)
		n += proto.SizeVarint(5<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

// Specifies a Gerrit project and its corresponding git repo.
type GerritProject struct {
	// The Gerrit host to connect to.
	//
	// Value must not include the schema part; it will be assumed to be "https".
	Host string `protobuf:"bytes,1,opt,name=host" json:"host,omitempty"`
	// Gerrit project name.
	Project string `protobuf:"bytes,2,opt,name=project" json:"project,omitempty"`
	// Full URL for the corresponding git repo.
	GitUrl string `protobuf:"bytes,3,opt,name=git_url,json=gitUrl" json:"git_url,omitempty"`
}

func (m *GerritProject) Reset()                    { *m = GerritProject{} }
func (m *GerritProject) String() string            { return proto.CompactTextString(m) }
func (*GerritProject) ProtoMessage()               {}
func (*GerritProject) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *GerritProject) GetHost() string {
	if m != nil {
		return m.Host
	}
	return ""
}

func (m *GerritProject) GetProject() string {
	if m != nil {
		return m.Project
	}
	return ""
}

func (m *GerritProject) GetGitUrl() string {
	if m != nil {
		return m.GitUrl
	}
	return ""
}

type GitRepo struct {
	// Full repository url, including schema.
	Url string `protobuf:"bytes,3,opt,name=url" json:"url,omitempty"`
}

func (m *GitRepo) Reset()                    { *m = GitRepo{} }
func (m *GitRepo) String() string            { return proto.CompactTextString(m) }
func (*GitRepo) ProtoMessage()               {}
func (*GitRepo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *GitRepo) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

// Access control rules.
type Acl struct {
	// Role of a group or identity.
	Role Acl_Role `protobuf:"varint,1,opt,name=role,enum=tricium.Acl_Role" json:"role,omitempty"`
	// Name of group, as defined in the auth service. Specify either group or
	// identity, not both.
	Group string `protobuf:"bytes,2,opt,name=group" json:"group,omitempty"`
	// Identity, as defined by the auth service. Can be either an email address
	// or an identity string, for instance, "anonymous:anonymous" for anonymous
	// users. Specify either group or identity, not both.
	Identity string `protobuf:"bytes,3,opt,name=identity" json:"identity,omitempty"`
}

func (m *Acl) Reset()                    { *m = Acl{} }
func (m *Acl) String() string            { return proto.CompactTextString(m) }
func (*Acl) ProtoMessage()               {}
func (*Acl) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *Acl) GetRole() Acl_Role {
	if m != nil {
		return m.Role
	}
	return Acl_READER
}

func (m *Acl) GetGroup() string {
	if m != nil {
		return m.Group
	}
	return ""
}

func (m *Acl) GetIdentity() string {
	if m != nil {
		return m.Identity
	}
	return ""
}

// Selection of function implementations to run for a project.
type Selection struct {
	// Name of function to run.
	Function string `protobuf:"bytes,1,opt,name=function" json:"function,omitempty"`
	// Name of platform to retrieve results from.
	Platform Platform_Name `protobuf:"varint,2,opt,name=platform,enum=tricium.Platform_Name" json:"platform,omitempty"`
	// Function configuration to use on this platform.
	Configs []*Config `protobuf:"bytes,3,rep,name=configs" json:"configs,omitempty"`
}

func (m *Selection) Reset()                    { *m = Selection{} }
func (m *Selection) String() string            { return proto.CompactTextString(m) }
func (*Selection) ProtoMessage()               {}
func (*Selection) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *Selection) GetFunction() string {
	if m != nil {
		return m.Function
	}
	return ""
}

func (m *Selection) GetPlatform() Platform_Name {
	if m != nil {
		return m.Platform
	}
	return Platform_ANY
}

func (m *Selection) GetConfigs() []*Config {
	if m != nil {
		return m.Configs
	}
	return nil
}

// Function configuration used when selecting a function implementation.
type Config struct {
	// Name of the configuration option.
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	// Value of the configuration.
	Value string `protobuf:"bytes,2,opt,name=value" json:"value,omitempty"`
}

func (m *Config) Reset()                    { *m = Config{} }
func (m *Config) String() string            { return proto.CompactTextString(m) }
func (*Config) ProtoMessage()               {}
func (*Config) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *Config) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Config) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

func init() {
	proto.RegisterType((*ServiceConfig)(nil), "tricium.ServiceConfig")
	proto.RegisterType((*ProjectConfig)(nil), "tricium.ProjectConfig")
	proto.RegisterType((*RepoDetails)(nil), "tricium.RepoDetails")
	proto.RegisterType((*GerritProject)(nil), "tricium.GerritProject")
	proto.RegisterType((*GitRepo)(nil), "tricium.GitRepo")
	proto.RegisterType((*Acl)(nil), "tricium.Acl")
	proto.RegisterType((*Selection)(nil), "tricium.Selection")
	proto.RegisterType((*Config)(nil), "tricium.Config")
	proto.RegisterEnum("tricium.Acl_Role", Acl_Role_name, Acl_Role_value)
}

func init() { proto.RegisterFile("infra/tricium/api/v1/config.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 721 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x54, 0xcb, 0x6e, 0xdb, 0x38,
	0x14, 0x8d, 0xfc, 0xf6, 0x75, 0xec, 0x38, 0x44, 0x90, 0xd1, 0x64, 0x16, 0xe3, 0x68, 0x10, 0x8c,
	0xdb, 0xa0, 0x36, 0xea, 0x2e, 0xda, 0x45, 0x8b, 0xc2, 0x4d, 0xdc, 0x64, 0x55, 0xa4, 0x74, 0xd2,
	0xad, 0xc0, 0x48, 0xb4, 0xc2, 0x56, 0x2f, 0x50, 0xb4, 0x83, 0x2c, 0xbb, 0xe9, 0x9f, 0xf4, 0xc3,
	0xfa, 0x0d, 0xfd, 0x81, 0x42, 0x7c, 0x48, 0x4a, 0x10, 0x14, 0xe8, 0x8e, 0xf7, 0x9e, 0x73, 0xc8,
	0x7b, 0x0f, 0x2f, 0x09, 0x87, 0x2c, 0x5e, 0x71, 0x32, 0x15, 0x9c, 0x79, 0x6c, 0x1d, 0x4d, 0x49,
	0xca, 0xa6, 0x9b, 0xe7, 0x53, 0x2f, 0x89, 0x57, 0x2c, 0x98, 0xa4, 0x3c, 0x11, 0x09, 0x6a, 0x6b,
	0xf0, 0xe0, 0xdf, 0x47, 0xb9, 0x3e, 0x11, 0x44, 0x31, 0x0f, 0xfe, 0x7b, 0x94, 0xb0, 0x5a, 0xc7,
	0x9e, 0x60, 0x49, 0xfc, 0x5b, 0x52, 0x1a, 0x12, 0xb1, 0x4a, 0x78, 0xa4, 0x48, 0xce, 0xcf, 0x1a,
	0xf4, 0x97, 0x94, 0x6f, 0x98, 0x47, 0x4f, 0x64, 0x2d, 0xe8, 0x25, 0x74, 0x0d, 0x27, 0xb3, 0xad,
	0x51, 0x7d, 0xdc, 0x9b, 0xfd, 0x3d, 0xd1, 0x9b, 0x4c, 0x2e, 0x8c, 0xfa, 0x94, 0x0a, 0xc2, 0xc2,
	0x0c, 0x97, 0x5c, 0xf4, 0x1a, 0xb6, 0xf3, 0x12, 0x5d, 0x5f, 0x41, 0x76, 0xed, 0x81, 0xf6, 0x34,
	0xaf, 0xff, 0xf2, 0x2e, 0xa5, 0x46, 0xdb, 0xcb, 0xe9, 0x3a, 0x40, 0x53, 0xe8, 0x9a, 0xfa, 0x33,
	0xbb, 0x2e, 0xa5, 0xbb, 0x85, 0xf4, 0xbd, 0x46, 0x70, 0xc9, 0x41, 0xc7, 0x00, 0x9c, 0x7a, 0x2c,
	0xa5, 0xae, 0x17, 0xf9, 0x76, 0x73, 0x64, 0x8d, 0x7b, 0xb3, 0xed, 0x42, 0x71, 0x12, 0xf9, 0xb8,
	0xab, 0xf0, 0x93, 0xc8, 0x47, 0x6f, 0x60, 0x47, 0x93, 0x53, 0xe2, 0x7d, 0x21, 0x01, 0xcd, 0xec,
	0x96, 0x3c, 0x63, 0xaf, 0x54, 0xb0, 0xd4, 0xbf, 0x50, 0x20, 0x1e, 0x28, 0xb2, 0x0e, 0x33, 0xf4,
	0x3f, 0xec, 0x64, 0xb7, 0x84, 0x47, 0x2c, 0x0e, 0xdc, 0x8c, 0xf2, 0x0d, 0xe5, 0x76, 0x7b, 0x64,
	0x8d, 0xbb, 0x78, 0x60, 0xd2, 0x4b, 0x99, 0x45, 0x47, 0x30, 0x60, 0x59, 0x12, 0x12, 0x41, 0x0d,
	0xaf, 0x23, 0x79, 0x7d, 0x9d, 0x55, 0x34, 0xe7, 0x7b, 0x0d, 0xfa, 0x17, 0x3c, 0xf9, 0x4c, 0x3d,
	0xa1, 0x5d, 0x1f, 0x41, 0x83, 0x78, 0x85, 0x69, 0x65, 0x1f, 0x73, 0x2f, 0xc4, 0x12, 0xf9, 0x73,
	0x83, 0x66, 0x00, 0x19, 0x0d, 0xa9, 0x56, 0x34, 0xa4, 0x02, 0x15, 0x8a, 0xa5, 0x81, 0x70, 0x85,
	0x85, 0x9e, 0x42, 0x93, 0xd3, 0x34, 0xc9, 0xec, 0xe6, 0x03, 0x77, 0x30, 0x4d, 0x13, 0x73, 0x6f,
	0x8a, 0x22, 0x4d, 0x51, 0x93, 0xe3, 0x12, 0xcf, 0x4b, 0xd6, 0xb1, 0xb0, 0x5b, 0xda, 0x14, 0x95,
	0x9e, 0xab, 0x2c, 0x7a, 0x05, 0xf6, 0x3d, 0xf7, 0xaa, 0x0a, 0x65, 0xe3, 0x7e, 0xd5, 0xc6, 0x52,
	0xe9, 0xfc, 0xb0, 0xa0, 0x57, 0x39, 0x19, 0xbd, 0x85, 0x41, 0x40, 0x39, 0x67, 0xc2, 0x4d, 0x95,
	0x7b, 0x76, 0x43, 0xde, 0xfb, 0x7e, 0x51, 0xe7, 0x99, 0x84, 0xb5, 0xb7, 0xe7, 0x5b, 0xb8, 0x1f,
	0x54, 0x13, 0xe8, 0x19, 0x74, 0x02, 0x26, 0xdc, 0xbc, 0x01, 0x3d, 0x32, 0xc3, 0x52, 0xca, 0x44,
	0x7e, 0xd6, 0xf9, 0x16, 0x6e, 0x07, 0x6a, 0x89, 0x8e, 0x61, 0xd7, 0x67, 0x19, 0xb9, 0x0e, 0xa9,
	0x94, 0x70, 0xc1, 0xe2, 0x40, 0x36, 0xd9, 0xc1, 0x43, 0x0d, 0x60, 0x93, 0xcf, 0xc9, 0xb7, 0x37,
	0x4c, 0xd0, 0x90, 0x65, 0x82, 0xfa, 0x6e, 0xc0, 0x93, 0x75, 0x6a, 0xb7, 0x47, 0xf5, 0x71, 0x17,
	0x0f, 0x2b, 0xc0, 0x59, 0x9e, 0x7f, 0xd7, 0x81, 0x56, 0x96, 0xac, 0xb9, 0x47, 0x9d, 0x4f, 0xd0,
	0xbf, 0x57, 0x34, 0x42, 0xd0, 0xb8, 0x49, 0x32, 0x61, 0x5b, 0xd2, 0x1a, 0xb9, 0x46, 0x36, 0xb4,
	0x4d, 0xc7, 0x35, 0x99, 0x36, 0x21, 0xfa, 0x0b, 0xf2, 0x6a, 0xdd, 0x35, 0x0f, 0xed, 0xba, 0x44,
	0x5a, 0x01, 0x13, 0x57, 0x3c, 0x74, 0xfe, 0x81, 0xb6, 0xee, 0x08, 0x0d, 0xa1, 0x5e, 0xe2, 0xf9,
	0xd2, 0xf9, 0x6a, 0x41, 0x7d, 0xee, 0x85, 0xe8, 0x08, 0x1a, 0x3c, 0x09, 0xa9, 0x3c, 0x6b, 0x50,
	0x99, 0xa7, 0xb9, 0x17, 0x4e, 0x70, 0x12, 0x52, 0x2c, 0x61, 0xb4, 0x07, 0x4d, 0xd5, 0x8e, 0x3a,
	0x5c, 0x05, 0xe8, 0x00, 0x3a, 0xcc, 0xa7, 0xb1, 0x60, 0xe2, 0x4e, 0xef, 0x5d, 0xc4, 0xce, 0x21,
	0x34, 0x72, 0x3d, 0x02, 0x68, 0xe1, 0xc5, 0xfc, 0x74, 0x81, 0x87, 0x5b, 0xa8, 0x0f, 0x5d, 0xbc,
	0xf8, 0x78, 0xb5, 0x58, 0x5e, 0x2e, 0xf0, 0xd0, 0x72, 0xbe, 0x59, 0xd0, 0x2d, 0xa6, 0x30, 0xdf,
	0xcc, 0x8c, 0xae, 0xee, 0xbc, 0x88, 0xd1, 0x0c, 0x3a, 0xe6, 0x9b, 0x91, 0x15, 0x0c, 0x2a, 0x17,
	0x5e, 0xfc, 0x48, 0x1f, 0x48, 0x44, 0x71, 0xc1, 0x43, 0x4f, 0xa0, 0xad, 0x3e, 0x57, 0xf3, 0x58,
	0x76, 0xca, 0x97, 0x2e, 0xf3, 0xd8, 0xe0, 0xce, 0x0c, 0x5a, 0xfa, 0x15, 0x22, 0x68, 0xc4, 0x24,
	0xa2, 0xc6, 0xfa, 0x7c, 0x9d, 0xf7, 0xbe, 0x21, 0xe1, 0x9a, 0x9a, 0xde, 0x65, 0x70, 0xdd, 0x92,
	0xdf, 0xe7, 0x8b, 0x5f, 0x01, 0x00, 0x00, 0xff, 0xff, 0xce, 0xa0, 0xd3, 0xbe, 0xd7, 0x05, 0x00,
	0x00,
}
