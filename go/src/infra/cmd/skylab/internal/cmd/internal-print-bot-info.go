// Copyright 2019 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/maruel/subcommands"
	"go.chromium.org/luci/auth/client/authcli"
	"go.chromium.org/luci/common/cli"
	"go.chromium.org/luci/grpc/prpc"

	fleet "infra/appengine/crosskylabadmin/api/fleet/v1"
	"infra/cmd/skylab/internal/site"
	"infra/libs/skylab/inventory"
	"infra/libs/skylab/inventory/swarming"
)

// InternalPrintBotInfo subcommand: Print Swarming dimensions for a DUT.
var InternalPrintBotInfo = &subcommands.Command{
	UsageLine: "internal-print-bot-info DUT_ID",
	ShortDesc: "print Swarming bot info for a DUT",
	LongDesc: `Print Swarming bot info for a DUT.

For internal use only.`,
	CommandRun: func() subcommands.CommandRun {
		c := &internalPrintBotInfoRun{}
		c.authFlags.Register(&c.Flags, site.DefaultAuthOptions)
		c.envFlags.Register(&c.Flags)
		return c
	},
}

type internalPrintBotInfoRun struct {
	subcommands.CommandRunBase
	authFlags authcli.Flags
	envFlags  envFlags
}

func (c *internalPrintBotInfoRun) Run(a subcommands.Application, args []string, env subcommands.Env) int {
	if err := c.innerRun(a, args, env); err != nil {
		fmt.Fprintf(a.GetErr(), "%s: %s\n", progName, err)
		return 1
	}
	return 0
}

func (c *internalPrintBotInfoRun) innerRun(a subcommands.Application, args []string, env subcommands.Env) error {
	if len(args) != 1 {
		return NewUsageError(c.Flags, "exactly one DUT_ID must be provided")
	}
	dutID := args[0]
	ctx := cli.GetContext(a, c, env)
	hc, err := httpClient(ctx, &c.authFlags)
	if err != nil {
		return err
	}
	siteEnv := c.envFlags.Env()
	ic := fleet.NewInventoryPRPCClient(&prpc.Client{
		C:       hc,
		Host:    siteEnv.AdminService,
		Options: site.DefaultPRPCOptions,
	})
	req := fleet.GetDutInfoRequest{Id: dutID}
	res, err := ic.GetDutInfo(ctx, &req)
	if err != nil {
		return err
	}
	var d inventory.DeviceUnderTest
	if err := proto.Unmarshal(res.GetSpec(), &d); err != nil {
		return err
	}
	stderr := a.GetErr()
	r := func(e error) { fmt.Fprintf(stderr, "sanitize dimensions: %s\n", err) }
	bi := botInfoForDUT(&d, r)
	enc, err := json.Marshal(bi)
	if err != nil {
		return err
	}
	a.GetOut().Write(enc)
	return nil
}

type botInfo struct {
	Dimensions swarming.Dimensions
	State      botState
}

type botState map[string][]string

func botInfoForDUT(d *inventory.DeviceUnderTest, r swarming.ReportFunc) botInfo {
	var bi botInfo
	bi.Dimensions = swarming.Convert(d.GetCommon().GetLabels())
	swarming.Sanitize(bi.Dimensions, r)
	bi.State = botStateForDUT(d)
	return bi
}

func botStateForDUT(d *inventory.DeviceUnderTest) botState {
	s := make(botState)
	for _, kv := range d.GetCommon().GetAttributes() {
		k, v := kv.GetKey(), kv.GetValue()
		s[k] = append(s[k], v)
	}
	return s
}
