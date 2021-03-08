// Copyright 2021 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package dut

import (
	"github.com/maruel/subcommands"
	"go.chromium.org/luci/auth/client/authcli"
	"go.chromium.org/luci/common/cli"
	"google.golang.org/protobuf/types/known/timestamppb"
	"infra/cmd/crosfleet/internal/buildbucket"
	"infra/cmd/crosfleet/internal/common"
	"infra/cmd/crosfleet/internal/site"
	"infra/cmdsupport/cmdlib"
	"time"
)

var abandon = &subcommands.Command{
	UsageLine: "abandon HOST [HOST...]",
	ShortDesc: "abandon DUTs which were previously leased via 'dut lease'",
	LongDesc: `Abandon DUTs which were previously leased via 'dut lease'.

This command's behavior is subject to change without notice.
Do not build automation around this subcommand.`,
	CommandRun: func() subcommands.CommandRun {
		c := &abandonRun{}
		c.authFlags.Register(&c.Flags, site.DefaultAuthOptions)
		c.envFlags.Register(&c.Flags)
		c.Flags.StringVar(&c.reason, "reason", "", "Optional reason for abandoning.")
		return c
	},
}

type abandonRun struct {
	subcommands.CommandRunBase
	reason    string
	authFlags authcli.Flags
	envFlags  common.EnvFlags
}

func (c *abandonRun) Run(a subcommands.Application, args []string, env subcommands.Env) int {
	if err := c.innerRun(a, args, env); err != nil {
		cmdlib.PrintError(a, err)
		return 1
	}
	return 0
}

func (c *abandonRun) innerRun(a subcommands.Application, args []string, env subcommands.Env) error {
	if len(args) == 0 {
		return cmdlib.NewUsageError(c.Flags, "must specify at least one DUT hostname")
	}

	ctx := cli.GetContext(a, c, env)
	userEmail, err := common.GetUserEmail(ctx, &c.authFlags)
	if err != nil {
		return err
	}
	swarmingService, err := newSwarmingService(ctx, c.envFlags.Env().SwarmingService, &c.authFlags)
	if err != nil {
		return err
	}
	bbClient, err := buildbucket.NewClient(ctx, c.envFlags.Env().DUTLeaserBuilderInfo, c.authFlags)
	if err != nil {
		return err
	}
	earliestCreateTime := earliestActiveLeaseTimestamp()
	for _, hostname := range args {
		correctedHostname := correctedHostname(hostname)
		id, err := hostnameToBotID(ctx, swarmingService, correctedHostname)
		if err != nil {
			return err
		}
		err = bbClient.CancelBuildWithBotID(ctx, id, earliestCreateTime, userEmail, c.reason, a.GetOut())
		if err != nil {
			return err
		}
	}
	return nil
}

// earliestActiveLeaseTimestamp returns the earliest possible creation
// timestamp for leases that haven't timed out.
func earliestActiveLeaseTimestamp() *timestamppb.Timestamp {
	return timestamppb.New(time.Now().Add(-1 * time.Minute * maxLeaseLengthMinutes))
}