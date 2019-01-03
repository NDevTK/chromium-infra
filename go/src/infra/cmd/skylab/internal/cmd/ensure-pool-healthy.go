// Copyright 2018 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package cmd

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"text/tabwriter"

	fleet "infra/appengine/crosskylabadmin/api/fleet/v1"
	"infra/cmd/skylab/internal/site"

	"github.com/maruel/subcommands"

	"go.chromium.org/luci/auth/client/authcli"
	"go.chromium.org/luci/common/cli"
	"go.chromium.org/luci/common/errors"
	"go.chromium.org/luci/grpc/prpc"
)

// EnsurePoolHealthy subcommand: Balance DUT pools
var EnsurePoolHealthy = &subcommands.Command{
	UsageLine: "ensure-pool-healthy [-dryrun] [-all-models] [-s spare] target [model...]",
	ShortDesc: "Ensure DUT pool is healthy",
	LongDesc: `
Ensure that given models' target pool contains healthy DUTs.
If needed, swap in healthy DUTs from spare pool.`,
	CommandRun: func() subcommands.CommandRun {
		c := &ensurePoolHealthyRun{}
		c.authFlags.Register(&c.Flags, site.DefaultAuthOptions)
		c.envFlags.Register(&c.Flags)

		c.Flags.BoolVar(&c.dryrun, "dryrun", false, "Dryrun mode -- do not commit inventory changes")
		c.Flags.BoolVar(&c.allModels, "all-models", false, "Ensure pool health for all known models")
		c.Flags.StringVar(&c.spare, "spare", "DUT_POOL_SUITES", "DUT pool to swap in healthy DUTs from")
		return c
	},
}

type ensurePoolHealthyRun struct {
	subcommands.CommandRunBase
	authFlags authcli.Flags
	envFlags  envFlags

	dryrun    bool
	allModels bool
	spare     string
}

type userError struct {
	error
}

func (c *ensurePoolHealthyRun) printError(w io.Writer, err error) {
	fmt.Fprintf(w, "%s\n\n", err)
	switch err.(type) {
	case userError:
		c.Flags.Usage()
	default:
		// Nothing more to say
	}
}

func (c *ensurePoolHealthyRun) Run(a subcommands.Application, args []string, env subcommands.Env) int {
	if err := c.innerRun(a, args, env); err != nil {
		c.printError(a.GetErr(), err)
		return 1
	}
	return 0
}

func (c *ensurePoolHealthyRun) innerRun(a subcommands.Application, args []string, env subcommands.Env) error {
	ctx := cli.GetContext(a, c, env)
	hc, err := httpClient(ctx, &c.authFlags)
	if err != nil {
		return err
	}
	e := c.envFlags.Env()

	target, err := c.getTargetPool(args)
	if err != nil {
		return err
	}
	models, err := c.getModels(ctx, hc, args)
	if err != nil {
		return err
	}

	ic := fleet.NewInventoryPRPCClient(&prpc.Client{
		C:       hc,
		Host:    e.AdminService,
		Options: site.DefaultPRPCOptions,
	})

	if c.dryrun {
		fmt.Fprintf(a.GetOut(), "DRYRUN: These changes are recommendations. Rerun without dryrun to apply changes.\n")
	}
	out := a.GetOut()
	for _, m := range models {
		if err := c.ensurePoolForModel(ctx, out, ic, target, m); err != nil {
			return err
		}
	}
	return nil
}

func (c *ensurePoolHealthyRun) ensurePoolForModel(ctx context.Context, w io.Writer, ic fleet.InventoryClient, target, model string) error {
	res, err := ic.EnsurePoolHealthy(
		ctx,
		&fleet.EnsurePoolHealthyRequest{
			DutSelector: &fleet.DutSelector{Model: model},
			TargetPool:  target,
			SparePool:   c.spare,
			Options: &fleet.EnsurePoolHealthyRequest_Options{
				Dryrun: c.dryrun,
			},
		},
	)
	if err != nil {
		return errors.Annotate(err, "ensure pool for %s", model).Err()
	}
	c.printEnsurePoolHealthyResult(w, model, target, res)
	return nil
}

func (c *ensurePoolHealthyRun) printEnsurePoolHealthyResult(w io.Writer, model, target string, res *fleet.EnsurePoolHealthyResponse) {
	bw := bufio.NewWriter(w)
	defer bw.Flush()

	// Align summary output
	tw := tabwriter.NewWriter(bw, 0, 2, 2, ' ', 0)
	defer tw.Flush()
	fmt.Fprintf(tw, "Model:\t%s\t\n", model)
	fmt.Fprintf(tw, "Target:\t%s\t\n", target)
	fmt.Fprintf(tw, "Spare:\t%s\t\n", c.spare)
	tp := res.GetTargetPoolStatus()
	fmt.Fprintf(tw, "New target pool status:\t%d/%d healthy\t\n", tp.GetHealthyCount(), tp.GetSize())
	sp := res.GetSparePoolStatus()
	fmt.Fprintf(tw, "New spare pool status:\t%d/%d healthy\t\n", sp.GetHealthyCount(), sp.GetSize())
	if len(res.GetFailures()) > 0 {
		fs := make([]string, 0, len(res.Failures))
		for _, f := range res.Failures {
			fs = append(fs, f.String())
		}
		fmt.Fprintf(tw, "Failures encountered:\t%s\t\n", strings.Join(fs, ", "))
	}
	if res.GetUrl() != "" {
		fmt.Fprintf(tw, "Inventory changes commited at:\t%s\t\n", res.GetUrl())
	}

	// Do not align inventory changes with the summary output.
	if len(res.GetChanges()) > 0 {
		fmt.Fprintf(bw, "Inventory changes:\n")
		for _, c := range res.GetChanges() {
			fmt.Fprintf(bw, "\t%s: %s\t->\t%s\n", c.DutId, c.OldPool, c.NewPool)
		}
	}
	fmt.Fprintf(bw, "\n")
}

func (*ensurePoolHealthyRun) getTargetPool(args []string) (string, error) {
	if len(args) < 1 {
		return "", userError{errors.New("want at least 1 arguments, have none")}
	}
	return args[0], nil
}

func (c *ensurePoolHealthyRun) getModels(ctx context.Context, hc *http.Client, args []string) ([]string, error) {
	numModelPosArgs := len(args) - 1
	if c.allModels {
		if numModelPosArgs > 0 {
			return []string{}, userError{fmt.Errorf("want no model postional arguments with -all-models, got %d", numModelPosArgs)}
		}
		return c.getAllModels(ctx, hc)
	}

	if numModelPosArgs < 1 {
		return []string{}, userError{fmt.Errorf("want at least 1 model positional argument, have %d", numModelPosArgs)}
	}
	return args[1:], nil
}

func (c *ensurePoolHealthyRun) getAllModels(ctx context.Context, hc *http.Client) ([]string, error) {
	// TODO(pprabhu) Consider implementing an RPC directly to ensure pool health
	// for all models.
	e := c.envFlags.Env()
	tc := fleet.NewTrackerPRPCClient(&prpc.Client{
		C:       hc,
		Host:    e.AdminService,
		Options: site.DefaultPRPCOptions,
	})
	res, err := tc.SummarizeBots(ctx, &fleet.SummarizeBotsRequest{})
	if err != nil {
		return []string{}, err
	}
	r := compileInventoryReport(res.GetBots())
	return modelsFromInventory(r.models), nil
}

func modelsFromInventory(ics []*inventoryCount) []string {
	ms := make([]string, 0, len(ics))
	for _, ic := range ics {
		ms = append(ms, ic.name)
	}
	return ms
}
