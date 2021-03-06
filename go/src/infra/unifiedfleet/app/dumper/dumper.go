// Copyright 2020 The Chromium OS Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package dumper

import (
	"context"
	"time"

	"cloud.google.com/go/bigquery"
	"go.chromium.org/luci/common/errors"
	"go.chromium.org/luci/common/logging"
	"go.chromium.org/luci/server"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	bqlib "infra/cros/lab_inventory/bq"
	"infra/unifiedfleet/app/cron"
	"infra/unifiedfleet/app/model/configuration"
	"infra/unifiedfleet/app/util"
)

// Jobs is a list of all the cron jobs that are currently available for running
var Jobs = []*cron.CronTab{
	{
		// Dump configs, registrations, inventory and states to BQ
		Name:     "ufs.dumper",
		Time:     20 * time.Minute,
		TrigType: cron.DAILY,
		Job:      dump,
	},
	{
		// Dump change events to BQ
		Name:     "ufs.change_event.BqDump",
		Time:     10 * time.Minute,
		TrigType: cron.EVERY,
		Job:      dumpChangeEvent,
	},
	{
		// Dump snapshots to BQ
		Name:     "ufs.snapshot_msg.BqDump",
		Time:     10 * time.Minute,
		TrigType: cron.EVERY,
		Job:      dumpChangeSnapshots,
	},
	{
		// Dump network configs to BQ
		Name:     "ufs.cros_network.dump",
		Time:     60 * time.Minute,
		TrigType: cron.EVERY,
		Job:      dumpCrosNetwork,
	},
	{
		// Sync asset info from HaRT
		Name:     "ufs.sync_devices.sync",
		TrigType: cron.HOURLY,
		Job:      SyncAssetInfoFromHaRT,
	},
	{
		// Push changes to dron queen
		Name:     "ufs.push_to_drone_queen",
		Time:     10 * time.Minute,
		TrigType: cron.EVERY,
		Job:      pushToDroneQueen,
	},
	{
		// Dump dut states to IV2
		Name:     "ufs.dump_to_invv2_dutstates",
		Time:     15 * time.Minute,
		TrigType: cron.DAILY,
		Job:      DumpToInventoryDutStateSnapshot,
	},
	{
		// Report UFS metrics
		Name:     "ufs.report_inventory",
		Time:     5 * time.Minute,
		TrigType: cron.EVERY,
		Job:      reportUFSInventoryCronHandler,
	},
}

// InitServer initializes a cron server.
func InitServer(srv *server.Server) {
	for _, job := range Jobs {
		// make a copy of the job to avoid race condition.
		t := job
		// Start all the cron jobs in background.
		srv.RunInBackground(job.Name, func(ctx context.Context) {
			cron.Run(ctx, t)
		})
	}
}

// Triggers a job by name. Returns error if the job is not found.
func TriggerJob(name string) error {
	for _, job := range Jobs {
		if job.Name == name {
			return cron.Trigger(job)
		}
	}
	return status.Errorf(codes.NotFound, "Invalid cron job %s. Not found", name)
}

func dump(ctx context.Context) error {
	ctx = logging.SetLevel(ctx, logging.Info)
	// Execute importing before dumping
	err1 := importCrimson(ctx)
	err2 := exportToBQ(ctx, dumpToBQ)
	if err1 == nil && err2 == nil {
		return nil
	}
	return errors.NewMultiError(err1, err2)
}

func dumpToBQ(ctx context.Context, bqClient *bigquery.Client) (err error) {
	defer func() {
		dumpToBQTick.Add(ctx, 1, err == nil)
	}()
	logging.Infof(ctx, "Dumping to BQ")
	curTime := time.Now()
	curTimeStr := bqlib.GetPSTTimeStamp(curTime)
	if err := configuration.SaveProjectConfig(ctx, &configuration.ProjectConfigEntity{
		Name:             getProject(ctx),
		DailyDumpTimeStr: curTimeStr,
	}); err != nil {
		return err
	}
	if err := dumpConfigurations(ctx, bqClient, curTimeStr); err != nil {
		return errors.Annotate(err, "dump configurations").Err()
	}
	if err := dumpRegistration(ctx, bqClient, curTimeStr); err != nil {
		return errors.Annotate(err, "dump registrations").Err()
	}
	if err := dumpInventory(ctx, bqClient, curTimeStr); err != nil {
		return errors.Annotate(err, "dump inventories").Err()
	}
	if err := dumpState(ctx, bqClient, curTimeStr); err != nil {
		return errors.Annotate(err, "dump states").Err()
	}
	logging.Debugf(ctx, "Dump is successfully finished")
	return nil
}

func dumpChangeEvent(ctx context.Context) (err error) {
	defer func() {
		dumpChangeEventTick.Add(ctx, 1, err == nil)
	}()
	ctx = logging.SetLevel(ctx, logging.Info)
	logging.Debugf(ctx, "Dumping change event to BQ")
	return exportToBQ(ctx, dumpChangeEventHelper)
}

func dumpChangeSnapshots(ctx context.Context) (err error) {
	defer func() {
		dumpChangeSnapshotTick.Add(ctx, 1, err == nil)
	}()
	ctx = logging.SetLevel(ctx, logging.Info)
	logging.Debugf(ctx, "Dumping change snapshots to BQ")
	return exportToBQ(ctx, dumpChangeSnapshotHelper)
}

func dumpCrosNetwork(ctx context.Context) (err error) {
	defer func() {
		dumpCrosNetworkTick.Add(ctx, 1, err == nil)
	}()
	// In UFS write to 'os' namespace
	ctx, err = util.SetupDatastoreNamespace(ctx, util.OSNamespace)
	if err != nil {
		return err
	}
	return importCrosNetwork(ctx)
}

// unique key used to store and retrieve context.
var contextKey = util.Key("ufs bigquery-client key")
var projectKey = util.Key("ufs project key")

// Use installs bigquery client to context.
func Use(ctx context.Context, bqClient *bigquery.Client) context.Context {
	return context.WithValue(ctx, contextKey, bqClient)
}

func get(ctx context.Context) *bigquery.Client {
	return ctx.Value(contextKey).(*bigquery.Client)
}

// UseProject installs project name to context.
func UseProject(ctx context.Context, project string) context.Context {
	return context.WithValue(ctx, projectKey, project)
}

func getProject(ctx context.Context) string {
	return ctx.Value(projectKey).(string)
}

func exportToBQ(ctx context.Context, f func(ctx context.Context, bqClient *bigquery.Client) error) (err error) {
	for _, ns := range util.ClientToDatastoreNamespace {
		newCtx, err1 := util.SetupDatastoreNamespace(ctx, ns)
		if ns == "" {
			// This is only for printing error message for default namespace.
			ns = "default"
		}
		logging.Debugf(newCtx, "Exporting to BQ for namespace %q", ns)
		if err1 != nil {
			err1 = errors.Annotate(err, "Setting namespace %q failed. BQ export skipped for the namespace %q", ns, ns).Err()
			logging.Errorf(ctx, err.Error())
			err = errors.NewMultiError(err, err1)
			continue
		}
		err1 = f(newCtx, get(newCtx))
		if err1 != nil {
			err1 = errors.Annotate(err, "BQ export failed for the namespace %q", ns).Err()
			logging.Errorf(ctx, err.Error())
			err = errors.NewMultiError(err, err1)
		}
	}
	return err
}
