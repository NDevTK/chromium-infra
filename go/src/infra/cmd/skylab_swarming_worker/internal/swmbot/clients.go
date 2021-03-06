// Copyright 2018 The Chromium OS Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package swmbot

import (
	"context"

	"go.chromium.org/luci/auth"
	"go.chromium.org/luci/common/errors"
	"go.chromium.org/luci/lucictx"
	"google.golang.org/grpc/metadata"

	invV2 "infra/appengine/cros/lab_inventory/api/v1"
	fleet "infra/appengine/crosskylabadmin/api/fleet/v1"
	ufsAPI "infra/unifiedfleet/api/v1/rpc"
	ufsUtil "infra/unifiedfleet/app/util"

	"infra/cmd/skylab_swarming_worker/internal/admin"
)

// WithTaskAccount returns a context using the Swarming task service
// account.
func WithTaskAccount(ctx context.Context) (context.Context, error) {
	return lucictx.SwitchLocalAccount(ctx, "task")
}

// WithSystemAccount returns acontext to using the Swarming bot system
// service account.
func WithSystemAccount(ctx context.Context) (context.Context, error) {
	return lucictx.SwitchLocalAccount(ctx, "system")
}

// InventoryClient returns an InventoryClient for the current Swarming
// bot task.  The context should use an explicit service account using
// WithTaskAccount or WithSystemAccount; otherwise the default service
// account is used.
func InventoryClient(ctx context.Context, b *Info) (fleet.InventoryClient, error) {
	o := auth.Options{
		Method: auth.LUCIContextMethod,
		Scopes: []string{
			auth.OAuthScopeEmail,
			"https://www.googleapis.com/auth/cloud-platform",
		},
	}
	c, err := admin.NewInventoryClient(ctx, b.AdminService, o)
	if err != nil {
		return nil, errors.Annotate(err, "create inventory client").Err()
	}
	return c, nil
}

// InventoryV2Client returns an InventoryClient for the current Swarming
// bot task. The context should use an explicit service account using
// WithTaskAccount or WithSystemAccount; otherwise the default service
// account is used.
func InventoryV2Client(ctx context.Context, b *Info) (invV2.InventoryClient, error) {
	o := auth.Options{
		Method: auth.LUCIContextMethod,
		Scopes: []string{
			auth.OAuthScopeEmail,
			"https://www.googleapis.com/auth/cloud-platform",
		},
	}
	pc, err := admin.NewPrpcClient(ctx, b.InventoryService, o)
	c := invV2.NewInventoryPRPCClient(pc)
	if err != nil {
		return nil, errors.Annotate(err, "create inventory V2 client").Err()
	}
	return c, nil
}

// UFSClient returns a FleetClient to communicate with UFS service.
// The context should use an explicit service account using
// WithTaskAccount or WithSystemAccount; otherwise the default service
// account is used.
func UFSClient(ctx context.Context, b *Info) (ufsAPI.FleetClient, error) {
	o := auth.Options{
		Method: auth.LUCIContextMethod,
		Scopes: []string{
			auth.OAuthScopeEmail,
			"https://www.googleapis.com/auth/cloud-platform",
		},
	}
	pc, err := admin.NewUFSClient(ctx, b.UFSService, o)
	if err != nil {
		return nil, errors.Annotate(err, "create UFS client").Err()
	}
	c := ufsAPI.NewFleetPRPCClient(pc)
	return c, nil
}

// SetupContext set up the outgoing context for API calls.
func SetupContext(ctx context.Context, namespace string) context.Context {
	md := metadata.Pairs(ufsUtil.Namespace, namespace)
	return metadata.NewOutgoingContext(ctx, md)
}
