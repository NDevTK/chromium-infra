// Copyright 2020 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package operations

import (
	"github.com/maruel/subcommands"
	"go.chromium.org/luci/common/cli"

	"infra/cmd/shivas/ufs/subcmds/asset"
	"infra/cmd/shivas/ufs/subcmds/chromeplatform"
	"infra/cmd/shivas/ufs/subcmds/drac"
	"infra/cmd/shivas/ufs/subcmds/dut"
	"infra/cmd/shivas/ufs/subcmds/host"
	"infra/cmd/shivas/ufs/subcmds/kvm"
	"infra/cmd/shivas/ufs/subcmds/machine"
	"infra/cmd/shivas/ufs/subcmds/machineprototype"
	"infra/cmd/shivas/ufs/subcmds/nic"
	"infra/cmd/shivas/ufs/subcmds/rack"
	"infra/cmd/shivas/ufs/subcmds/rackprototype"
	"infra/cmd/shivas/ufs/subcmds/rpm"
	"infra/cmd/shivas/ufs/subcmds/static"
	"infra/cmd/shivas/ufs/subcmds/switches"
	"infra/cmd/shivas/ufs/subcmds/vlan"
	"infra/cmd/shivas/ufs/subcmds/vm"
)

type get struct {
	subcommands.CommandRunBase
}

// GetCmd contains get command specification
var GetCmd = &subcommands.Command{
	UsageLine: "get <sub-command>",
	ShortDesc: "Get details of a resource/entity",
	LongDesc: `Get details for
	machine/rack/kvm/rpm/switch/drac/nic
	dut
	host/vm/vm-slots
	machine-prototype/rack-prototype/platform/vlan`,
	CommandRun: func() subcommands.CommandRun {
		c := &get{}
		return c
	},
}

type getApp struct {
	cli.Application
}

// Run implementing subcommands.CommandRun interface
func (c *get) Run(a subcommands.Application, args []string, env subcommands.Env) int {
	d := a.(*cli.Application)
	return subcommands.Run(&getApp{*d}, args)
}

// GetCommands lists all the subcommands under get
func (c getApp) GetCommands() []*subcommands.Command {
	return []*subcommands.Command{
		subcommands.CmdHelp,
		asset.GetAssetCmd,
		machine.GetMachineCmd,
		dut.GetDutCmd,
		host.GetHostCmd,
		kvm.GetKVMCmd,
		rpm.GetRPMCmd,
		switches.GetSwitchCmd,
		drac.GetDracCmd,
		nic.GetNicCmd,
		vm.GetVMCmd,
		vm.ListVMSlotCmd,
		rack.GetRackCmd,
		machineprototype.GetMachineLSEPrototypeCmd,
		rackprototype.GetRackLSEPrototypeCmd,
		chromeplatform.GetChromePlatformCmd,
		vlan.GetVlanCmd,
		static.GetStatesCmd,
		static.GetZonesCmd,
	}
}

// GetName is cli.Application interface implementation
func (c getApp) GetName() string {
	return "get"
}
