// Copyright 2021 The Chromium OS Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package cros

import (
	"context"
	"strings"

	"go.chromium.org/luci/common/errors"

	"infra/cros/recovery/internal/execs"
	"infra/cros/recovery/internal/log"
)

const (
	readHWIDCommand = "crossystem hwid"
)

// updateHWIDToInvExec read HWID from the resource and update DUT info.
func updateHWIDToInvExec(ctx context.Context, args *execs.RunArgs, actionArgs []string) error {
	r := args.Access.Run(ctx, args.ResourceName, readHWIDCommand)
	if r.ExitCode != 0 {
		return errors.Reason("update HWID in DUT-info: failed with code: %d, %q", r.ExitCode, r.Stderr).Err()
	}
	hwid := strings.TrimSpace(r.Stdout)
	if hwid == "" {
		return errors.Reason("update HWID in DUT-info: is empty").Err()
	}
	log.Debug(ctx, "Update HWID %q in DUT-info.", hwid)
	args.DUT.Hwid = hwid
	return nil
}

// matchHWIDToInvExec matches HWID from the resource to value in the Inventory.
func matchHWIDToInvExec(ctx context.Context, args *execs.RunArgs, actionArgs []string) error {
	r := args.Access.Run(ctx, args.ResourceName, readHWIDCommand)
	if r.ExitCode != 0 {
		return errors.Reason("match HWID to inventory: failed with code: %d, %q", r.ExitCode, r.Stderr).Err()
	}
	expectedHWID := args.DUT.Hwid
	actualHWID := strings.TrimSpace(r.Stdout)
	if actualHWID != expectedHWID {
		return errors.Reason("match HWID to inventory: failed, expected: %q, but got %q", expectedHWID, actualHWID).Err()
	}
	return nil
}

func init() {
	execs.Register("cros_update_hwid_to_inventory", updateHWIDToInvExec)
	execs.Register("cros_match_hwid_to_inventory", matchHWIDToInvExec)
}
