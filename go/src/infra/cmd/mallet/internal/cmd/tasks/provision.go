// Copyright 2021 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package tasks

import (
	b64 "encoding/base64"
	"fmt"

	"github.com/maruel/subcommands"
	"go.chromium.org/luci/auth/client/authcli"
	"go.chromium.org/luci/common/cli"
	"go.chromium.org/luci/common/errors"

	"infra/cmd/mallet/internal/site"
	"infra/cmdsupport/cmdlib"
	"infra/cros/recovery/tasknames"
	"infra/libs/skylab/buildbucket"
	"infra/libs/skylab/buildbucket/labpack"
)

// Recovery subcommand: Recovering the devices.
var CustomProvision = &subcommands.Command{
	UsageLine: "provision",
	ShortDesc: "Quick provision Cros DUT.",
	LongDesc:  "Quick provision Cros DUT. Tool allows provide custom values for provisioning.",
	CommandRun: func() subcommands.CommandRun {
		c := &customProvisionRun{}
		c.authFlags.Register(&c.Flags, site.DefaultAuthOptions)
		c.envFlags.Register(&c.Flags)
		c.Flags.StringVar(&c.osName, "os", "", "ChromeOS version name like eve-release/R86-13380.0.0")
		c.Flags.StringVar(&c.osPath, "os-path", "", "GS path to where the payloads are located. Example: gs://chromeos-image-archive/eve-release/R86-13380.0.0")
		return c
	},
}

type customProvisionRun struct {
	subcommands.CommandRunBase
	authFlags authcli.Flags
	envFlags  site.EnvFlags

	osName string
	osPath string
}

func (c *customProvisionRun) Run(a subcommands.Application, args []string, env subcommands.Env) int {
	if err := c.innerRun(a, args, env); err != nil {
		cmdlib.PrintError(a, err)
		return 1
	}
	return 0
}

func (c *customProvisionRun) innerRun(a subcommands.Application, args []string, env subcommands.Env) error {
	ctx := cli.GetContext(a, c, env)

	bc, err := buildbucket.NewClient(ctx, c.authFlags, site.DefaultPRPCOptions, site.BBProject, site.MalletBucket, site.MalletBuilder)
	if err != nil {
		return err
	}
	if len(args) == 0 {
		return errors.Reason("create recovery task: unit is not specified").Err()
	}
	unit := args[0]
	e := c.envFlags.Env()
	configuration := b64.StdEncoding.EncodeToString([]byte(c.createPlan()))
	taskID, err := labpack.ScheduleTask(
		ctx,
		bc,
		&labpack.Params{
			UnitName:         unit,
			TaskName:         string(tasknames.Custom),
			AdminService:     e.AdminService,
			InventoryService: e.UFSService,
			NoMetrics:        true,
			Configuration:    configuration,
		},
	)
	if err != nil {
		return errors.Annotate(err, "create recovery task").Err()
	}
	fmt.Fprintf(a.GetOut(), "Created recovery task for %s: %s\n", unit, bc.BuildURL(taskID))
	return nil
}

// Custom plan to execute provision
// TODO(otabek): Replace by build plan on fly.
const customProvisionPlanStart = `
{
	"plan_names": [
		"cros"
	],
	"plans": {
		"cros": {
			"critical_actions": [
				"cros_ssh",
				"cros_default_boot",
				"Custom provision"
			],
			"actions": {
				"cros_default_boot": {
					"dependencies": [
						"cros_storage_writing"
					],
					"exec_name": "cros_is_default_boot_from_disk"
				},
				"cros_ping": {
					"exec_name": "cros_ping"
				},
				"cros_ssh": {
					"dependencies": [
						"has_dut_name",
						"has_dut_board_name",
						"has_dut_model_name",
						"cros_ping"
					],
					"exec_name": "cros_ssh"
				},
				"cros_storage_writing": {
					"dependencies": [
						"cros_ssh"
					],
					"exec_name": "cros_is_file_system_writable"
				},
				"Custom provision":{
					"docs":[
						"Provision device to the custom os version"
					],
					"exec_name": "cros_provision",
					"exec_extra_args":[`
const customProvisionPlanTail = `
					],
					"exec_timeout": {
						"seconds": 3600
					}
				}
			}
		}
	}
}`

func (c *customProvisionRun) createPlan() string {
	var customArg string
	if c.osPath != "" {
		customArg = fmt.Sprintf("os_image_path:%s", c.osPath)
	} else if c.osName != "" {
		customArg = fmt.Sprintf("os_name:%s", c.osName)
	}
	return customProvisionPlanStart + customArg + customProvisionPlanTail
}
