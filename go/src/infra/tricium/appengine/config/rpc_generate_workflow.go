// Copyright 2017 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package config

import (
	"go.chromium.org/luci/common/logging"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	admin "infra/tricium/api/admin/v1"
	"infra/tricium/appengine/common/config"
)

func (*configServer) GenerateWorkflow(c context.Context, req *admin.GenerateWorkflowRequest) (*admin.GenerateWorkflowResponse, error) {
	if req.Project == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "missing project name")
	}
	pc, err := config.LuciConfigServer.GetProjectConfig(c, req.Project)
	if err != nil {
		logging.WithError(err).Errorf(c, "failed to get project config")
		return nil, grpc.Errorf(codes.InvalidArgument, "failed to get project config")
	}
	sc, err := config.LuciConfigServer.GetServiceConfig(c)
	if err != nil {
		logging.WithError(err).Errorf(c, "failed to get service config")
		return nil, grpc.Errorf(codes.InvalidArgument, "failed to get service config")
	}
	wf, err := config.Generate(sc, pc, req.Files)
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "failed to validate config")
	}
	return &admin.GenerateWorkflowResponse{Workflow: wf}, nil
}
