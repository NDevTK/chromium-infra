// Copyright 2019 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Package metrics stores the reported JSON metrics from depot_tools into a
// BigQuery table.
package metrics

import (
	"go.chromium.org/luci/common/data/stringset"
	"go.chromium.org/luci/common/tsmon/distribution"
	"go.chromium.org/luci/common/tsmon/field"
	"go.chromium.org/luci/common/tsmon/metric"
	"go.chromium.org/luci/common/tsmon/types"
	"golang.org/x/net/context"
	"infra/appengine/depot_tools_metrics/schema"
)

// chromiumSrc is the URL of the chromium/src repo. It is counted apart from all
// the other repos.
const chromiumSrc = "https://chromium.googlesource.com/chromium/src"

var (
	// GitLatency keeps track of how long it takes to run a git command per repo and exit code.
	// We only keep track of the git commands executed by a depot_tools command.
	GitLatency = metric.NewCumulativeDistribution(
		"depot_tools_metrics/git/latency",
		"Time it takes to run a git command.",
		&types.MetricMetadata{Units: types.Seconds},
		// A growth factor if 1.045 with 200 buckets covers up to about 100m,
		// which is the interval we're interested about.
		distribution.GeometricBucketer(1.045, 200),
		field.String("command"),
		field.Int("exit_code"),
		field.String("repo"),
	)
)

// reportGitPushMetrics reports git push metrics to ts_mon.
func reportGitPushMetrics(ctx context.Context, m schema.Metrics) {
	if len(m.ProjectUrls) == 0 {
		return
	}
	if len(m.SubCommands) == 0 {
		return
	}
	repo := "everything_else"
	if stringset.NewFromSlice(m.ProjectUrls...).Has(chromiumSrc) {
		repo = chromiumSrc
	}
	for _, sc := range m.SubCommands {
		if sc.Command != "git push" {
			continue
		}
		GitLatency.Add(ctx, sc.ExecutionTime, "git push", sc.ExitCode, repo)
	}
}
