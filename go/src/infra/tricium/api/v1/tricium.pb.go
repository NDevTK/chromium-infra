// Code generated by protoc-gen-go. DO NOT EDIT.
// source: infra/tricium/api/v1/tricium.proto

package tricium

import prpc "github.com/luci/luci-go/grpc/prpc"

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type State int32

const (
	// Pending is for when an analysis request has been received but the corresponding
	// workflow, analyzer, or workers are not running yet.
	State_PENDING State = 0
	// Running is for when the workflow, analyzer, or workers of a request have been launched
	// but have not finished.
	State_RUNNING State = 1
	// Success is for a workflow, analyzer, or worker that successfully completed.
	//
	// Success of workflows and analyzers, is aggregated from underlying analyzers and workers,
	// where full success means success is aggregated.
	State_SUCCESS State = 2
	// Failure is for a workflow, analyzer, or worker that completed with failure.
	//
	// Failure of workflows and analyzers, is aggregated from underlying analyzers and workers,
	// where any occurence of failure means failure is aggregated.
	State_FAILURE State = 3
	// Canceled is for user canceled workflows, analyzers, and workers.
	// NB! Not supported yet.
	State_CANCELED State = 4
	// Timed out is for workers where the triggered swarming task timed out.
	// NB! Not supported yet.
	State_TIMED_OUT State = 5
)

var State_name = map[int32]string{
	0: "PENDING",
	1: "RUNNING",
	2: "SUCCESS",
	3: "FAILURE",
	4: "CANCELED",
	5: "TIMED_OUT",
}
var State_value = map[string]int32{
	"PENDING":   0,
	"RUNNING":   1,
	"SUCCESS":   2,
	"FAILURE":   3,
	"CANCELED":  4,
	"TIMED_OUT": 5,
}

func (x State) String() string {
	return proto.EnumName(State_name, int32(x))
}
func (State) EnumDescriptor() ([]byte, []int) { return fileDescriptor3, []int{0} }

type Consumer int32

const (
	// No consumer means progress and results can be pulled from the Tricium service API.
	//
	// This is the default value used when no consumer is configured for a project.
	Consumer_NONE Consumer = 0
	// A Gerrit consumer means the Tricium service polls Gerrit for changes and reports
	// progress updates and results to Gerrit.
	//
	// Gerrit details need to be configured for a project before the Gerrit consumer
	// option is enabled.
	Consumer_GERRIT Consumer = 1
)

var Consumer_name = map[int32]string{
	0: "NONE",
	1: "GERRIT",
}
var Consumer_value = map[string]int32{
	"NONE":   0,
	"GERRIT": 1,
}

func (x Consumer) String() string {
	return proto.EnumName(Consumer_name, int32(x))
}
func (Consumer) EnumDescriptor() ([]byte, []int) { return fileDescriptor3, []int{1} }

// AnalyzeRequest contains the details needed for an analysis request.
type AnalyzeRequest struct {
	// Name of the project hosting the paths listed in the request. The name
	// should map to the project name as it is connected to Tricium.
	Project string `protobuf:"bytes,1,opt,name=project" json:"project,omitempty"`
	GitRef  string `protobuf:"bytes,2,opt,name=git_ref,json=gitRef" json:"git_ref,omitempty"`
	// Paths to analyze in the project. Listed from the root of the Git
	// repository.
	// TODO(emso): document path separators or add listing of path segments.
	Paths []string `protobuf:"bytes,3,rep,name=paths" json:"paths,omitempty"`
	// Consumer to send progress updates and results to.
	//
	// This field is optional. If included it will push progress and result
	// updates to the provided consumer. The selected consumer must be
	// configured for the project of the request.
	//
	// Note that progress and results can be accessed via the Tricium
	// API regardless of whether a consumer has been included in the request.
	Consumer Consumer `protobuf:"varint,4,opt,name=consumer,enum=tricium.Consumer" json:"consumer,omitempty"`
	// Gerrit details for a Gerrit consumer; change and revision.
	//
	// These fields are required if Gerrit is selected as consumer.
	GerritChange   string `protobuf:"bytes,5,opt,name=gerrit_change,json=gerritChange" json:"gerrit_change,omitempty"`
	GerritRevision string `protobuf:"bytes,6,opt,name=gerrit_revision,json=gerritRevision" json:"gerrit_revision,omitempty"`
}

func (m *AnalyzeRequest) Reset()                    { *m = AnalyzeRequest{} }
func (m *AnalyzeRequest) String() string            { return proto.CompactTextString(m) }
func (*AnalyzeRequest) ProtoMessage()               {}
func (*AnalyzeRequest) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{0} }

func (m *AnalyzeRequest) GetProject() string {
	if m != nil {
		return m.Project
	}
	return ""
}

func (m *AnalyzeRequest) GetGitRef() string {
	if m != nil {
		return m.GitRef
	}
	return ""
}

func (m *AnalyzeRequest) GetPaths() []string {
	if m != nil {
		return m.Paths
	}
	return nil
}

func (m *AnalyzeRequest) GetConsumer() Consumer {
	if m != nil {
		return m.Consumer
	}
	return Consumer_NONE
}

func (m *AnalyzeRequest) GetGerritChange() string {
	if m != nil {
		return m.GerritChange
	}
	return ""
}

func (m *AnalyzeRequest) GetGerritRevision() string {
	if m != nil {
		return m.GerritRevision
	}
	return ""
}

type AnalyzeResponse struct {
	// ID of the run started for this request.
	//
	// This ID can be used to track progress and request results.
	RunId string `protobuf:"bytes,1,opt,name=run_id,json=runId" json:"run_id,omitempty"`
}

func (m *AnalyzeResponse) Reset()                    { *m = AnalyzeResponse{} }
func (m *AnalyzeResponse) String() string            { return proto.CompactTextString(m) }
func (*AnalyzeResponse) ProtoMessage()               {}
func (*AnalyzeResponse) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{1} }

func (m *AnalyzeResponse) GetRunId() string {
	if m != nil {
		return m.RunId
	}
	return ""
}

type ProgressRequest struct {
	// Run ID returned by an analyze request.
	//
	// This field must be provided. If nothing else is provided, then
	// all known progress for the run is returned.
	RunId string `protobuf:"bytes,1,opt,name=run_id,json=runId" json:"run_id,omitempty"`
	// An optional analyzer name.
	//
	// If provided, only progress for the provided analyzer will be returned.
	// The analyzer name should match the name of the analyzer in the Tricium
	// configuration.
	//
	// NB! Currently not supported.
	Analyzer string `protobuf:"bytes,2,opt,name=analyzer" json:"analyzer,omitempty"`
	// Optional platform that may be provided together with an analyzer name.
	//
	// If provided, only progress for the provided analyzer and platform will be provided.
	//
	// NB! Currently not supported.
	Platform *Platform `protobuf:"bytes,3,opt,name=platform" json:"platform,omitempty"`
}

func (m *ProgressRequest) Reset()                    { *m = ProgressRequest{} }
func (m *ProgressRequest) String() string            { return proto.CompactTextString(m) }
func (*ProgressRequest) ProtoMessage()               {}
func (*ProgressRequest) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{2} }

func (m *ProgressRequest) GetRunId() string {
	if m != nil {
		return m.RunId
	}
	return ""
}

func (m *ProgressRequest) GetAnalyzer() string {
	if m != nil {
		return m.Analyzer
	}
	return ""
}

func (m *ProgressRequest) GetPlatform() *Platform {
	if m != nil {
		return m.Platform
	}
	return nil
}

type ProgressResponse struct {
	// Overall state for the run provided in the progress request.
	State State `protobuf:"varint,1,opt,name=state,enum=tricium.State" json:"state,omitempty"`
	// Analyzer progress matching the requested progress report.
	//
	// For a provided run ID this corresponds to all analyzers and platforms, and
	// for any selection of these, a subset is returned.
	//
	// NB! Selection of a subset is currently not supported.
	AnalyzerProgress []*AnalyzerProgress `protobuf:"bytes,2,rep,name=analyzer_progress,json=analyzerProgress" json:"analyzer_progress,omitempty"`
}

func (m *ProgressResponse) Reset()                    { *m = ProgressResponse{} }
func (m *ProgressResponse) String() string            { return proto.CompactTextString(m) }
func (*ProgressResponse) ProtoMessage()               {}
func (*ProgressResponse) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{3} }

func (m *ProgressResponse) GetState() State {
	if m != nil {
		return m.State
	}
	return State_PENDING
}

func (m *ProgressResponse) GetAnalyzerProgress() []*AnalyzerProgress {
	if m != nil {
		return m.AnalyzerProgress
	}
	return nil
}

type AnalyzerProgress struct {
	// The analyzer name.
	Analyzer string `protobuf:"bytes,1,opt,name=analyzer" json:"analyzer,omitempty"`
	// The platform for which the analyzer progress is reported.
	Platform Platform_Name `protobuf:"varint,2,opt,name=platform,enum=tricium.Platform_Name" json:"platform,omitempty"`
	// The state of the analyzer.
	//
	// For an analyzer on a specific platform this state corresponds to the state
	// of the worker, else it is the aggregated state of all workers for the analyzer.
	State State `protobuf:"varint,3,opt,name=state,enum=tricium.State" json:"state,omitempty"`
	// The Swarming task ID running the analyzer worker.
	SwarmingTaskId string `protobuf:"bytes,4,opt,name=swarming_task_id,json=swarmingTaskId" json:"swarming_task_id,omitempty"`
	// Number of comments.
	//
	// For analyzers that are done and produce comments.
	NumComments int32 `protobuf:"varint,5,opt,name=num_comments,json=numComments" json:"num_comments,omitempty"`
}

func (m *AnalyzerProgress) Reset()                    { *m = AnalyzerProgress{} }
func (m *AnalyzerProgress) String() string            { return proto.CompactTextString(m) }
func (*AnalyzerProgress) ProtoMessage()               {}
func (*AnalyzerProgress) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{4} }

func (m *AnalyzerProgress) GetAnalyzer() string {
	if m != nil {
		return m.Analyzer
	}
	return ""
}

func (m *AnalyzerProgress) GetPlatform() Platform_Name {
	if m != nil {
		return m.Platform
	}
	return Platform_ANY
}

func (m *AnalyzerProgress) GetState() State {
	if m != nil {
		return m.State
	}
	return State_PENDING
}

func (m *AnalyzerProgress) GetSwarmingTaskId() string {
	if m != nil {
		return m.SwarmingTaskId
	}
	return ""
}

func (m *AnalyzerProgress) GetNumComments() int32 {
	if m != nil {
		return m.NumComments
	}
	return 0
}

type ResultsRequest struct {
	// Run ID returned by an analyze request.
	RunId string `protobuf:"bytes,1,opt,name=run_id,json=runId" json:"run_id,omitempty"`
	// An optional analyzer name.
	//
	// If provided, only results for the provided analyzer are returned.
	// If an analyzer is being run on more than one platform then the merged
	// results of the analyzer can be returned by exclusion of a specific platform.
	//
	// NB! Currently not supported.
	Analyzer string `protobuf:"bytes,2,opt,name=analyzer" json:"analyzer,omitempty"`
	// Optional platform that can be provided together with an analyzer name.
	//
	// If provided, only results for the provided platform and analyzer are returned.
	//
	// NB! Currently not supported.
	Platform Platform_Name `protobuf:"varint,3,opt,name=platform,enum=tricium.Platform_Name" json:"platform,omitempty"`
}

func (m *ResultsRequest) Reset()                    { *m = ResultsRequest{} }
func (m *ResultsRequest) String() string            { return proto.CompactTextString(m) }
func (*ResultsRequest) ProtoMessage()               {}
func (*ResultsRequest) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{5} }

func (m *ResultsRequest) GetRunId() string {
	if m != nil {
		return m.RunId
	}
	return ""
}

func (m *ResultsRequest) GetAnalyzer() string {
	if m != nil {
		return m.Analyzer
	}
	return ""
}

func (m *ResultsRequest) GetPlatform() Platform_Name {
	if m != nil {
		return m.Platform
	}
	return Platform_ANY
}

type ResultsResponse struct {
	// TODO(emso): Support paging of results to deal with large number of results.
	Results *Data_Results `protobuf:"bytes,1,opt,name=results" json:"results,omitempty"`
	// Whether the returned results are merged.
	//
	// Results may be merged if a result request for an analyzer running on multiple
	// platforms was made and the request did not include a specific platform.
	// Results for a run with no specific analyzer selected will be marked as merged
	// if any included analyzer results were merged.
	IsMerged bool `protobuf:"varint,2,opt,name=is_merged,json=isMerged" json:"is_merged,omitempty"`
}

func (m *ResultsResponse) Reset()                    { *m = ResultsResponse{} }
func (m *ResultsResponse) String() string            { return proto.CompactTextString(m) }
func (*ResultsResponse) ProtoMessage()               {}
func (*ResultsResponse) Descriptor() ([]byte, []int) { return fileDescriptor3, []int{6} }

func (m *ResultsResponse) GetResults() *Data_Results {
	if m != nil {
		return m.Results
	}
	return nil
}

func (m *ResultsResponse) GetIsMerged() bool {
	if m != nil {
		return m.IsMerged
	}
	return false
}

func init() {
	proto.RegisterType((*AnalyzeRequest)(nil), "tricium.AnalyzeRequest")
	proto.RegisterType((*AnalyzeResponse)(nil), "tricium.AnalyzeResponse")
	proto.RegisterType((*ProgressRequest)(nil), "tricium.ProgressRequest")
	proto.RegisterType((*ProgressResponse)(nil), "tricium.ProgressResponse")
	proto.RegisterType((*AnalyzerProgress)(nil), "tricium.AnalyzerProgress")
	proto.RegisterType((*ResultsRequest)(nil), "tricium.ResultsRequest")
	proto.RegisterType((*ResultsResponse)(nil), "tricium.ResultsResponse")
	proto.RegisterEnum("tricium.State", State_name, State_value)
	proto.RegisterEnum("tricium.Consumer", Consumer_name, Consumer_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Tricium service

type TriciumClient interface {
	// Analyze requests analysis of a list of paths.
	//
	// An analysis request for a list of paths in a project connected to Tricium
	// via the Tricium configuration. On success, the ID of the resulting run is
	// returned.
	Analyze(ctx context.Context, in *AnalyzeRequest, opts ...grpc.CallOption) (*AnalyzeResponse, error)
	// Progress requests progress information for a run.
	//
	// A run corresponds to an analyze request and is identified with a run ID.
	Progress(ctx context.Context, in *ProgressRequest, opts ...grpc.CallOption) (*ProgressResponse, error)
	// Results requests analysis results from a run.
	//
	// A run corresponds to an analyze request and is identified with a run ID.
	Results(ctx context.Context, in *ResultsRequest, opts ...grpc.CallOption) (*ResultsResponse, error)
}
type triciumPRPCClient struct {
	client *prpc.Client
}

func NewTriciumPRPCClient(client *prpc.Client) TriciumClient {
	return &triciumPRPCClient{client}
}

func (c *triciumPRPCClient) Analyze(ctx context.Context, in *AnalyzeRequest, opts ...grpc.CallOption) (*AnalyzeResponse, error) {
	out := new(AnalyzeResponse)
	err := c.client.Call(ctx, "tricium.Tricium", "Analyze", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *triciumPRPCClient) Progress(ctx context.Context, in *ProgressRequest, opts ...grpc.CallOption) (*ProgressResponse, error) {
	out := new(ProgressResponse)
	err := c.client.Call(ctx, "tricium.Tricium", "Progress", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *triciumPRPCClient) Results(ctx context.Context, in *ResultsRequest, opts ...grpc.CallOption) (*ResultsResponse, error) {
	out := new(ResultsResponse)
	err := c.client.Call(ctx, "tricium.Tricium", "Results", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type triciumClient struct {
	cc *grpc.ClientConn
}

func NewTriciumClient(cc *grpc.ClientConn) TriciumClient {
	return &triciumClient{cc}
}

func (c *triciumClient) Analyze(ctx context.Context, in *AnalyzeRequest, opts ...grpc.CallOption) (*AnalyzeResponse, error) {
	out := new(AnalyzeResponse)
	err := grpc.Invoke(ctx, "/tricium.Tricium/Analyze", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *triciumClient) Progress(ctx context.Context, in *ProgressRequest, opts ...grpc.CallOption) (*ProgressResponse, error) {
	out := new(ProgressResponse)
	err := grpc.Invoke(ctx, "/tricium.Tricium/Progress", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *triciumClient) Results(ctx context.Context, in *ResultsRequest, opts ...grpc.CallOption) (*ResultsResponse, error) {
	out := new(ResultsResponse)
	err := grpc.Invoke(ctx, "/tricium.Tricium/Results", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Tricium service

type TriciumServer interface {
	// Analyze requests analysis of a list of paths.
	//
	// An analysis request for a list of paths in a project connected to Tricium
	// via the Tricium configuration. On success, the ID of the resulting run is
	// returned.
	Analyze(context.Context, *AnalyzeRequest) (*AnalyzeResponse, error)
	// Progress requests progress information for a run.
	//
	// A run corresponds to an analyze request and is identified with a run ID.
	Progress(context.Context, *ProgressRequest) (*ProgressResponse, error)
	// Results requests analysis results from a run.
	//
	// A run corresponds to an analyze request and is identified with a run ID.
	Results(context.Context, *ResultsRequest) (*ResultsResponse, error)
}

func RegisterTriciumServer(s prpc.Registrar, srv TriciumServer) {
	s.RegisterService(&_Tricium_serviceDesc, srv)
}

func _Tricium_Analyze_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AnalyzeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TriciumServer).Analyze(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tricium.Tricium/Analyze",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TriciumServer).Analyze(ctx, req.(*AnalyzeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Tricium_Progress_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProgressRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TriciumServer).Progress(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tricium.Tricium/Progress",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TriciumServer).Progress(ctx, req.(*ProgressRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Tricium_Results_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResultsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TriciumServer).Results(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tricium.Tricium/Results",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TriciumServer).Results(ctx, req.(*ResultsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Tricium_serviceDesc = grpc.ServiceDesc{
	ServiceName: "tricium.Tricium",
	HandlerType: (*TriciumServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Analyze",
			Handler:    _Tricium_Analyze_Handler,
		},
		{
			MethodName: "Progress",
			Handler:    _Tricium_Progress_Handler,
		},
		{
			MethodName: "Results",
			Handler:    _Tricium_Results_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "infra/tricium/api/v1/tricium.proto",
}

func init() { proto.RegisterFile("infra/tricium/api/v1/tricium.proto", fileDescriptor3) }

var fileDescriptor3 = []byte{
	// 650 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x54, 0x5d, 0x4f, 0xdb, 0x48,
	0x14, 0xc5, 0xf9, 0xb2, 0x73, 0x03, 0x89, 0x19, 0x2d, 0x8b, 0xc9, 0x3e, 0x6c, 0xd6, 0xac, 0xd4,
	0x08, 0xa9, 0x44, 0x4d, 0x5f, 0x91, 0x2a, 0x94, 0x18, 0x14, 0x09, 0x0c, 0x9a, 0x24, 0x52, 0xdf,
	0xac, 0x69, 0x32, 0x31, 0x2e, 0xf8, 0xa3, 0x33, 0x63, 0x50, 0xfb, 0xd4, 0x5f, 0xd7, 0x7f, 0xd0,
	0xfe, 0x9e, 0xca, 0xf6, 0xd8, 0x09, 0x26, 0xaa, 0x2a, 0xf5, 0xf1, 0x9c, 0x7b, 0x7c, 0xcf, 0x9c,
	0x7b, 0x67, 0x0c, 0xa6, 0x17, 0xac, 0x18, 0x19, 0x08, 0xe6, 0x2d, 0xbc, 0xd8, 0x1f, 0x90, 0xc8,
	0x1b, 0x3c, 0xbe, 0xc9, 0xe1, 0x69, 0xc4, 0x42, 0x11, 0x22, 0x55, 0xc2, 0xee, 0xbf, 0x5b, 0xc5,
	0x4b, 0x22, 0x48, 0xa6, 0xec, 0x1e, 0x6f, 0x15, 0x44, 0x0f, 0x44, 0xac, 0x42, 0x26, 0xdb, 0x99,
	0x3f, 0x14, 0x68, 0x9f, 0x07, 0xe4, 0xe1, 0xf3, 0x17, 0x8a, 0xe9, 0xa7, 0x98, 0x72, 0x81, 0x0c,
	0x50, 0x23, 0x16, 0x7e, 0xa4, 0x0b, 0x61, 0x28, 0x3d, 0xa5, 0xdf, 0xc4, 0x39, 0x44, 0x87, 0xa0,
	0xba, 0x9e, 0x70, 0x18, 0x5d, 0x19, 0x95, 0xb4, 0xd2, 0x70, 0x3d, 0x81, 0xe9, 0x0a, 0xfd, 0x05,
	0xf5, 0x88, 0x88, 0x3b, 0x6e, 0x54, 0x7b, 0xd5, 0x7e, 0x13, 0x67, 0x00, 0xbd, 0x06, 0x6d, 0x11,
	0x06, 0x3c, 0xf6, 0x29, 0x33, 0x6a, 0x3d, 0xa5, 0xdf, 0x1e, 0xee, 0x9f, 0xe6, 0x61, 0x46, 0xb2,
	0x80, 0x0b, 0x09, 0x3a, 0x86, 0x3d, 0x97, 0x32, 0xe6, 0x09, 0x67, 0x71, 0x47, 0x02, 0x97, 0x1a,
	0xf5, 0xd4, 0x63, 0x37, 0x23, 0x47, 0x29, 0x87, 0x5e, 0x41, 0x47, 0x8a, 0x18, 0x7d, 0xf4, 0xb8,
	0x17, 0x06, 0x46, 0x23, 0x95, 0xb5, 0x33, 0x1a, 0x4b, 0xd6, 0xec, 0x43, 0xa7, 0xc8, 0xc5, 0xa3,
	0x30, 0xe0, 0x14, 0x1d, 0x40, 0x83, 0xc5, 0x81, 0xe3, 0x2d, 0x65, 0xae, 0x3a, 0x8b, 0x83, 0xc9,
	0xd2, 0xe4, 0xd0, 0xb9, 0x65, 0xa1, 0xcb, 0x28, 0xe7, 0xf9, 0x08, 0xb6, 0x2b, 0x51, 0x17, 0x34,
	0x92, 0xf5, 0x64, 0x72, 0x00, 0x05, 0x4e, 0xc2, 0xe6, 0xa3, 0x35, 0xaa, 0x3d, 0xa5, 0xdf, 0xda,
	0x08, 0x7b, 0x2b, 0x0b, 0xb8, 0x90, 0x98, 0x5f, 0x15, 0xd0, 0xd7, 0xae, 0xf2, 0x80, 0xff, 0x43,
	0x9d, 0x0b, 0x22, 0x68, 0xea, 0xda, 0x1e, 0xb6, 0x8b, 0x06, 0xd3, 0x84, 0xc5, 0x59, 0x11, 0x5d,
	0xc0, 0x7e, 0xee, 0xea, 0x44, 0xb2, 0x85, 0x51, 0xe9, 0x55, 0xfb, 0xad, 0xe1, 0x51, 0xf1, 0x85,
	0xcc, 0xce, 0x0a, 0x0f, 0x9d, 0x94, 0x18, 0xf3, 0xbb, 0x02, 0x7a, 0x59, 0xf6, 0x2c, 0xa2, 0x52,
	0x8a, 0x38, 0xdc, 0x88, 0x58, 0x49, 0x4f, 0xf8, 0xf7, 0x8b, 0x88, 0xa7, 0x36, 0xf1, 0xe9, 0x3a,
	0xe7, 0x3a, 0x52, 0xf5, 0x57, 0x91, 0xfa, 0xa0, 0xf3, 0x27, 0xc2, 0x7c, 0x2f, 0x70, 0x1d, 0x41,
	0xf8, 0x7d, 0x32, 0xf9, 0x5a, 0xb6, 0xd6, 0x9c, 0x9f, 0x11, 0x7e, 0x3f, 0x59, 0xa2, 0xff, 0x60,
	0x37, 0x88, 0x7d, 0x67, 0x11, 0xfa, 0x3e, 0x0d, 0x04, 0x4f, 0xef, 0x48, 0x1d, 0xb7, 0x82, 0xd8,
	0x1f, 0x49, 0xca, 0x7c, 0x82, 0x36, 0xa6, 0x3c, 0x7e, 0x10, 0x7f, 0xb2, 0xce, 0x61, 0x69, 0x9d,
	0xbf, 0x91, 0xd5, 0x74, 0xa0, 0x53, 0x18, 0xcb, 0x8d, 0x0e, 0x40, 0x65, 0x19, 0x95, 0x5a, 0xb7,
	0x86, 0x07, 0x45, 0x97, 0x71, 0xf2, 0x52, 0x73, 0x7d, 0xae, 0x42, 0xff, 0x40, 0xd3, 0xe3, 0x8e,
	0x4f, 0x99, 0x4b, 0x97, 0xe9, 0xa1, 0x34, 0xac, 0x79, 0xfc, 0x3a, 0xc5, 0x27, 0xef, 0xa1, 0x9e,
	0x8e, 0x0d, 0xb5, 0x40, 0xbd, 0xb5, 0xec, 0xf1, 0xc4, 0xbe, 0xd4, 0x77, 0x12, 0x80, 0xe7, 0xb6,
	0x9d, 0x00, 0x25, 0x01, 0xd3, 0xf9, 0x68, 0x64, 0x4d, 0xa7, 0x7a, 0x25, 0x01, 0x17, 0xe7, 0x93,
	0xab, 0x39, 0xb6, 0xf4, 0x2a, 0xda, 0x05, 0x6d, 0x74, 0x6e, 0x8f, 0xac, 0x2b, 0x6b, 0xac, 0xd7,
	0xd0, 0x1e, 0x34, 0x67, 0x93, 0x6b, 0x6b, 0xec, 0xdc, 0xcc, 0x67, 0x7a, 0xfd, 0xa4, 0x07, 0x5a,
	0xfe, 0x22, 0x91, 0x06, 0x35, 0xfb, 0xc6, 0xb6, 0xf4, 0x1d, 0x04, 0xd0, 0xb8, 0xb4, 0x30, 0x9e,
	0xcc, 0x74, 0x65, 0xf8, 0x4d, 0x01, 0x75, 0x96, 0x1d, 0x1d, 0x9d, 0x81, 0x2a, 0x2f, 0x0e, 0x3a,
	0x2c, 0xdf, 0x38, 0x39, 0xf3, 0xae, 0xf1, 0xb2, 0x20, 0x67, 0xf2, 0x0e, 0xb4, 0xe2, 0xba, 0xad,
	0x55, 0xa5, 0x27, 0xd8, 0x3d, 0xda, 0x52, 0x91, 0x0d, 0xce, 0x40, 0x95, 0x73, 0xdb, 0xb0, 0x7f,
	0xbe, 0xf2, 0x0d, 0xfb, 0xd2, 0x4a, 0x3e, 0x34, 0xd2, 0x1f, 0xdf, 0xdb, 0x9f, 0x01, 0x00, 0x00,
	0xff, 0xff, 0x50, 0x1e, 0xf2, 0xb5, 0x6d, 0x05, 0x00, 0x00,
}
