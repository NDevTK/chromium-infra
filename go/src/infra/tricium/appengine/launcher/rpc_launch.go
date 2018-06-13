// Copyright 2017 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package launcher

import (
	"encoding/json"
	"fmt"

	"github.com/golang/protobuf/proto"
	ds "go.chromium.org/gae/service/datastore"
	tq "go.chromium.org/gae/service/taskqueue"
	"go.chromium.org/luci/common/logging"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"infra/tricium/api/admin/v1"
	"infra/tricium/api/v1"
	"infra/tricium/appengine/common"
	"infra/tricium/appengine/common/config"
)

// LauncherServer represents the Tricium pRPC Launcher server.
type launcherServer struct{}

var server = &launcherServer{}

// Launch processes one launch request to the Tricium Launcher.
func (r *launcherServer) Launch(c context.Context, req *admin.LaunchRequest) (*admin.LaunchResponse, error) {
	logging.Infof(c, "[launcher] Launch request (run ID: %d)", req.RunId)
	if req.RunId == 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "missing run ID")
	}
	if req.Project == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "missing project")
	}
	if req.GitUrl == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "missing git URL")
	}
	if req.GitRef == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "missing git ref")
	}
	if len(req.Files) == 0 {
		return nil, grpc.Errorf(codes.InvalidArgument, "missing paths to analyze")
	}
	if err := launch(c, req, config.LuciConfigServer, common.IsolateServer,
		common.SwarmingServer, common.PubsubServer); err != nil {
		return nil, grpc.Errorf(codes.Internal, "failed to launch workflow: %v", err)
	}
	return &admin.LaunchResponse{}, nil
}

func launch(c context.Context, req *admin.LaunchRequest, cp config.ProviderAPI, isolator common.IsolateAPI,
	swarming common.SwarmingAPI, pubsub common.PubSubAPI) error {
	// Guard checking if there is already a stored workflow for the run ID
	// in the request; if so stop here.
	w := &config.Workflow{ID: req.RunId}
	if err := ds.Get(c, w); err != ds.ErrNoSuchEntity {
		logging.Infof(c, "[launcher] Launch request for launched workflow, run ID: %s", req.RunId)
		return nil
	}

	// Generate workflow and convert to string.
	sc, err := cp.GetServiceConfig(c)
	if err != nil {
		logging.WithError(err).Errorf(c, "failed to get service config")
		return err
	}
	pc, err := cp.GetProjectConfig(c, req.Project)
	if err != nil {
		logging.WithError(err).Errorf(c, "failed to get project config")
		return err
	}
	configJSON, _ := json.Marshal(pc)
	logging.Fields{
		"project": req.Project,
		"config":  string(configJSON),
	}.Debugf(c, "Got project config")
	wf, err := config.Generate(sc, pc, req.Files)
	if err != nil {
		return fmt.Errorf("failed to generate workflow config for project %s: %v", req.Project, err)
	}

	// Set up pubsub for worker completion notification.
	err = pubsub.Setup(c)
	if err != nil {
		return fmt.Errorf("failed to setup pubsub for workflow: %v", err)
	}
	wfb, err := proto.Marshal(wf)
	if err != nil {
		return fmt.Errorf("failed to marshal workflow proto: %v", err)
	}

	// Prepare workflow config entry to store.
	wfConfig := &config.Workflow{
		ID:                 req.RunId,
		SerializedWorkflow: wfb,
	}

	// Prepare workflow launched request.
	b, err := proto.Marshal(&admin.WorkflowLaunchedRequest{RunId: req.RunId})
	if err != nil {
		return fmt.Errorf("failed to marshal trigger request proto: %v", err)
	}
	wfTask := tq.NewPOSTTask("/tracker/internal/workflow-launched", nil)
	wfTask.Payload = b

	// Isolate initial input.
	inputHash, err := isolator.IsolateGitFileDetails(c, wf.IsolateServer, &tricium.Data_GitFileDetails{
		Repository: req.GitUrl,
		Ref:        req.GitRef,
		Files:      req.Files,
	})
	if err != nil {
		return fmt.Errorf("failed to isolate git file details: %v", err)
	}
	logging.Infof(c, "[launcher] Isolated git file details, hash: %q", inputHash)

	// TODO(emso): select dev/prod swarming/isolate serve based on
	// devserver and dev/prod tricium instance. Prepare trigger requests
	// for root workers.
	wTasks := []*tq.Task{}
	for _, worker := range wf.RootWorkers() {
		b, err := proto.Marshal(&admin.TriggerRequest{
			RunId:             req.RunId,
			IsolatedInputHash: inputHash,
			Worker:            worker,
		})
		if err != nil {
			return fmt.Errorf("failed to encode driver request: %v", err)
		}
		t := tq.NewPOSTTask("/driver/internal/trigger", nil)
		t.Payload = b
		wTasks = append(wTasks, t)
	}
	return ds.RunInTransaction(c, func(c context.Context) (err error) {
		// Store workflow config.
		if err := ds.Put(c, wfConfig); err != nil {
			return fmt.Errorf("failed to store workflow: %v", err)
		}
		// Run the below two operations in parallel.
		done := make(chan error)
		defer func() {
			if err2 := <-done; err2 != nil {
				err = err2
			}
		}()
		go func() {
			// Mark workflow as launched. Processing of this
			// request needs the stored workflow config.
			if err := tq.Add(c, common.TrackerQueue, wfTask); err != nil {
				done <- fmt.Errorf("failed to enqueue workflow launched track request: %v", err)
			}
			done <- nil
		}()
		// Trigger root workers. Processing of this request needs the
		// stored workflow config.
		if err := tq.Add(c, common.DriverQueue, wTasks...); err != nil {
			return fmt.Errorf("failed to enqueue trigger request(s) for root worker(s): %v", err)
		}
		return nil
	}, nil)
}
