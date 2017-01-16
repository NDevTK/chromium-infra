// Copyright 2016 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Package common implements common functionality for the Tricium service modules.
package common

import (
	"errors"
	"net/http"
	"strings"

	"golang.org/x/net/context"

	"google.golang.org/appengine"

	"github.com/golang/protobuf/proto"
	ds "github.com/luci/gae/service/datastore"
	"github.com/luci/gae/service/info"
	"github.com/luci/luci-go/appengine/gaeauth/server"
	"github.com/luci/luci-go/appengine/gaemiddleware"
	"github.com/luci/luci-go/common/logging"
	"github.com/luci/luci-go/grpc/discovery"
	"github.com/luci/luci-go/grpc/prpc"
	"github.com/luci/luci-go/server/auth"
	"github.com/luci/luci-go/server/auth/xsrf"
	"github.com/luci/luci-go/server/router"
	"github.com/luci/luci-go/server/templates"

	"infra/tricium/api/admin/v1"
	"infra/tricium/api/v1"
)

const (
	// AnalyzeQueue specifies the name of the analyze task queue.
	AnalyzeQueue = "analyze-queue"

	// LauncherQueue specifies the name of the launcher task queue.
	LauncherQueue = "launcher-queue"

	// DriverQueue specified the name of the driver task queue.
	DriverQueue = "driver-queue"

	// TrackerQueue specified the name of the tracker task queue.
	TrackerQueue = "tracker-queue"
)

// TODO(emso): Use string IDs every where and use a ID translation scheme between
// the key visible to users and the ID used to store in datastore. This removes the
// temptation for users to try things like ID+1.

// Workflow config entry for storing in datastore.
type Workflow struct {
	ID int64 `gae:"$id"`
	// Serialized workflow config proto.
	SerializedWorkflow []byte `gae:",noindex"`
}

// ReportServerError reports back a server error (http code 500).
func ReportServerError(c *router.Context, err error) {
	logging.WithError(err).Errorf(c.Context, "HTTP 500")
	http.Error(c.Writer, "An internal server error occured. We are working on it ;)",
		http.StatusInternalServerError)
}

// MiddlewareBase returns a middleware chain applied to ALL routes.
func MiddlewareBase() router.MiddlewareChain {
	return gaemiddleware.BaseProd()
}

// MiddlewareForInternal returns a middleware chain applied to internal routes.
//
// It assumes internal routes are protected by specifying 'login: admin'
// app.yaml.
func MiddlewareForInternal() router.MiddlewareChain {
	// TODO(vadimsh): Figure out how to assert that the handler is called by GAE
	// itself or by PubSub. That's how internal routes are supposed to be called.
	return MiddlewareBase()
}

// MiddlewareForUI returns a middleware chain intended for Web UI routes.
//
// It's same as MiddlewareBase, with addition of authentication based on
// cookies. It is supposed to be used for all routes that result in HTML pages.
//
// Anonymous access is still allowed. The handlers should do authorization
// checks by examining auth.CurrentIdentity(ctx) (it can be identity.Anonymous).
func MiddlewareForUI() router.MiddlewareChain {
	// Configure auth system to use cookies and actually attempt to do the
	// authentication. Finally, configure templating system.
	return MiddlewareBase().Extend(
		auth.Use(auth.Authenticator{server.CookieAuth}),
		auth.Authenticate,
		templates.WithTemplates(prepareTemplates()),
	)
}

// MiddlewareForREST returns a middleware chain intended for REST API routes.
//
// It's same as MiddlewareBase, with addition of authentication based on
// OAuth2 access tokens. It is supposed to be used for all REST API routes.
//
// Anonymous access is still allowed. The handlers should do authorization
// checks by examining auth.CurrentIdentity(ctx) (it can be identity.Anonymous).
func MiddlewareForREST() router.MiddlewareChain {
	return MiddlewareBase().Extend(
		auth.Use(auth.Authenticator{
			&server.OAuth2Method{Scopes: []string{server.EmailScope}},
		}),
		auth.Authenticate,
	)
}

// MiddlewareForRPC returns a middleware chain intended for pRPC routes.
//
// It is identical to MiddlewareBase currently, since pRPC does its
// authentication (based on OAuth2 access tokens) internally.
func MiddlewareForRPC() router.MiddlewareChain {
	return MiddlewareBase()
}

// NewRPCServer returns preconfigured pRPC server that can host gRPC APIs.
//
// Usage:
//   srv := NewRPCServer()
//   someapi.RegisterSomeAPIServer(&srv)
//   ...
//   srv.InstallHandlers(router, MiddlewareForRPC())
func NewRPCServer() *prpc.Server {
	srv := &prpc.Server{
		Authenticator: auth.Authenticator{
			&server.OAuth2Method{Scopes: []string{server.EmailScope}},
		},
		// TODO(vadimsh): Enable monitoring interceptor.
		// UnaryServerInterceptor: grpcmon.NewUnaryServerInterceptor(nil),
	}
	discovery.Enable(srv)
	return srv
}

// NewGAEContext constructs a context compatible with standard appengine lib.
//
// The returned context is compatible with both LUCI libs and GAE std lib.
//
// TODO(emso): Get rid of it once everything is converted to use luci/gae.
func NewGAEContext(c *router.Context) context.Context {
	return appengine.WithContext(c.Context, c.Request)
}

// prepareTemplates returns templates.Bundle used by all UI handlers.
//
// In particular it includes a set of default arguments passed to all templates.
func prepareTemplates() *templates.Bundle {
	return &templates.Bundle{
		Loader:          templates.FileSystemLoader("templates"),
		DebugMode:       info.IsDevAppServer,
		DefaultTemplate: "base", // defined in includes/base.hml
		DefaultArgs: func(c context.Context) (templates.Args, error) {
			loginURL, err := auth.LoginURL(c, "/")
			if err != nil {
				return nil, err
			}
			logoutURL, err := auth.LogoutURL(c, "/")
			if err != nil {
				return nil, err
			}
			token, err := xsrf.Token(c)
			if err != nil {
				return nil, err
			}
			return templates.Args{
				"AppVersion":  strings.Split(info.VersionID(c), ".")[0],
				"IsAnonymous": auth.CurrentIdentity(c) == "anonymous:anonymous",
				"User":        auth.CurrentUser(c),
				"LoginURL":    loginURL,
				"LogoutURL":   logoutURL,
				"XsrfToken":   token,
			}, nil
		},
	}
}

// WorkflowProvider provides a workflow config from a project or a run ID.
type WorkflowProvider interface {
	ReadConfigForProject(context.Context, string) (*admin.Workflow, error)
	ReadConfigForRun(context.Context, int64) (*admin.Workflow, error)
}

// LuciConfigProvider provides workflow configurations from the Luci-config service.
type LuciConfigProvider struct {
}

// ReadConfigForProject reads a workflow config for a project from Luci-config.
func (*LuciConfigProvider) ReadConfigForProject(c context.Context, project string) (*admin.Workflow, error) {
	// TODO(emso): Replace this dummy config with one read from luci-config.
	return &admin.Workflow{
		WorkerTopic:    "projects/tricium-dev/topics/worker-completion",
		ServiceAccount: "emso@chromium.org",
		Workers: []*admin.Worker{
			{
				Name:     "Hello_Ubuntu14.04_x86-64",
				Needs:    tricium.Data_GIT_FILE_DETAILS,
				Provides: tricium.Data_FILES,
				Platform: "Ubuntu14.04_x86-64",
				Dimensions: []string{
					"pool:Chrome",
					"os:Ubuntu-14.04",
					"cpu:x84-64",
				},
				Cmd: &tricium.Cmd{
					Exec: "echo",
					Args: []string{
						"'hello'",
					},
				},
				Deadline: 30,
			},
		},
	}, nil
}

// ReadConfigForRun is not supported by this workflow provider, included to match the interface.
func (*LuciConfigProvider) ReadConfigForRun(c context.Context, runID int64) (*admin.Workflow, error) {
	return nil, errors.New("Luci-config workflow provider cannot provide config for run ID")
}

// DatastoreConfigProvider provides workflow configurations from Datastore.
type DatastoreConfigProvider struct {
}

// ReadConfigForProject is not supported by this workflow provider, included to match the interface.
func (*DatastoreConfigProvider) ReadConfigForProject(c context.Context, project string) (*admin.Workflow, error) {
	return nil, errors.New("Datastore workflow provider cannot provide config for project name")
}

// ReadConfigForRun provides workflow configurations for a run ID from Datastore.
func (*DatastoreConfigProvider) ReadConfigForRun(c context.Context, runID int64) (*admin.Workflow, error) {
	wfb := &Workflow{ID: runID}
	if err := ds.Get(c, wfb); err != nil {
		return nil, err
	}
	wf := &admin.Workflow{}
	if err := proto.Unmarshal(wfb.SerializedWorkflow, wf); err != nil {
		return nil, err
	}
	return wf, nil
}
