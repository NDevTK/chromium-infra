// Copyright 2017 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package frontend

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/golang/protobuf/proto"
	ds "github.com/luci/gae/service/datastore"
	tq "github.com/luci/gae/service/taskqueue"
	"github.com/luci/luci-go/common/clock"
	"github.com/luci/luci-go/common/logging"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	admin "infra/tricium/api/admin/v1"
	"infra/tricium/api/v1"
	"infra/tricium/appengine/common"
	"infra/tricium/appengine/common/track"
)

// TriciumServer represents the Tricium pRPC server.
type TriciumServer struct{}

// Server instance to use within this module/package.
var server = &TriciumServer{}

const repo = "https://chromium-review.googlesource.com/playground/gerrit-tricium"

// Analyze processes one analysis request to Tricium.
//
// Launched a workflow customized to the project and listed paths.  The run ID
// in the response can be used to track the progress and results of the request
// via the Tricium UI.
func (r *TriciumServer) Analyze(c context.Context, req *tricium.AnalyzeRequest) (*tricium.AnalyzeResponse, error) {
	if req.Project == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "missing project")
	}
	if req.GitRef == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "missing git ref")
	}
	if len(req.Paths) == 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "missing paths to analyze")
	}
	runID, err := analyze(c, req, &common.LuciConfigProvider{})
	if err != nil {
		return nil, grpc.Errorf(codes.Internal, "failed to execute analyze request")
	}
	logging.Infof(c, "[frontend] Run ID: %s", runID)
	return &tricium.AnalyzeResponse{runID}, nil
}

func analyze(c context.Context, req *tricium.AnalyzeRequest, cp common.ConfigProvider) (string, error) {
	// TODO(emso): let a failed project config request mean project is unknown and skip getting the service config?
	sc, err := cp.GetServiceConfig(c)
	if err != nil {
		return "", fmt.Errorf("failed to get service config: %v", err)
	}
	if !sc.ProjectIsKnown(req.Project) {
		return "", errors.New("unknown project")
	}
	pc, err := cp.GetProjectConfig(c, req.Project)
	if err != nil {
		return "", fmt.Errorf("failed to get project config: %v", err)
	}
	ok, err := pc.CanRequest(c)
	if err != nil {
		return "", fmt.Errorf("failed to authorize: %v", err)
	}
	if !ok {
		// TODO(emso): make this bubble up as a permission denied error
		return "", fmt.Errorf("no request access for project %s, context: %v", req.Project, c)
	}
	// TODO(emso): Verify that there is no current run for this request (map hashed requests to run IDs).
	// TODO(emso): Read Git repo info from the configuration projects/ endpoint.
	run := &track.Run{
		Received: clock.Now(c).UTC(),
		State:    track.Pending,
		Project:  req.Project,
	}
	sr := &track.ServiceRequest{
		Project: req.Project,
		Paths:   req.Paths,
		GitRepo: repo,
		GitRef:  req.GitRef,
	}
	lr := &admin.LaunchRequest{
		Project: sr.Project,
		Paths:   sr.Paths,
		GitRepo: repo,
		GitRef:  sr.GitRef,
	}
	// This is a cross-group transaction because first Run is stored to get the ID,
	// and then ServiceRequest is stored, with Run key as parent.
	err = ds.RunInTransaction(c, func(c context.Context) (err error) {
		// Add tracking entries for run.
		if err := ds.Put(c, run); err != nil {
			return fmt.Errorf("failed to store run entry: %v", err)
		}
		// Run the below operations in parallel.
		done := make(chan error)
		defer func() {
			if err2 := <-done; err == nil {
				err = err2
			}
		}()
		go func() {
			// Add tracking entry for service request.
			sr.Parent = ds.KeyForObj(c, run)
			if err := ds.Put(c, sr); err != nil {
				done <- fmt.Errorf("failed to store service request: %v", err)
				return
			}
			done <- nil
		}()
		// Launch workflow, enqueue launch request.
		lr.RunId = run.ID
		t := tq.NewPOSTTask("/launcher/internal/launch", nil)
		b, err := proto.Marshal(lr)
		if err != nil {
			return fmt.Errorf("failed to enqueue launch request: %v", err)
		}
		t.Payload = b
		return tq.Add(c, common.LauncherQueue, t)
	}, &ds.TransactionOptions{XG: true})
	if err != nil {
		return "", fmt.Errorf("failed to track and launch request: %v", err)
	}
	return strconv.FormatInt(run.ID, 10), nil
}
