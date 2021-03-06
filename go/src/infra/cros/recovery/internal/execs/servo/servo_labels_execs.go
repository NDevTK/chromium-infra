// Copyright 2021 The Chromium OS Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package servo

import (
	"context"

	"go.chromium.org/luci/common/errors"

	"infra/cros/recovery/internal/execs"
	"infra/cros/recovery/internal/log"
)

// updateServoTypeLabelExec updates DUT's servo type to the correct servo type string.
func updateServoTypeLabelExec(ctx context.Context, args *execs.RunArgs, actionArgs []string) error {
	servoType, err := GetServoType(ctx, args)
	if err != nil {
		return errors.Annotate(err, "update servo type label").Err()
	}
	args.DUT.ServoHost.Servo.Type = servoType
	log.Info(ctx, "Set DUT's servo type to be: %s", servoType)
	return nil
}

func init() {
	execs.Register("servo_update_servo_type_label", updateServoTypeLabelExec)
}
