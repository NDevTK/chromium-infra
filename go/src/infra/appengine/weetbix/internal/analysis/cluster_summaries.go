// Copyright 2021 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package analysis

import (
	"context"
	"fmt"
	"math"
	"time"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"

	"go.chromium.org/luci/common/errors"

	"infra/appengine/weetbix/internal/bqutil"
	"infra/appengine/weetbix/internal/clustering"
	"infra/appengine/weetbix/internal/config"
)

// ImpactfulClusterReadOptions specifies options for ReadImpactfulClusters().
type ImpactfulClusterReadOptions struct {
	// Project is the LUCI Project for which analysis is being performed.
	Project string
	// Thresholds is the set of thresholds, which if any are met
	// or exceeded, should result in the cluster being returned.
	// Thresholds are applied based on cluster residual impact.
	Thresholds *config.ImpactThreshold
	// AlwaysInclude is the set of clusters to always include.
	AlwaysInclude []clustering.ClusterID
}

// ClusterSummary represents a statistical summary of a cluster's failures,
// and their impact.
type ClusterSummary struct {
	ClusterID            clustering.ClusterID `json:"clusterId"`
	PresubmitRejects1d   Counts               `json:"presubmitRejects1d"`
	PresubmitRejects3d   Counts               `json:"presubmitRejects3d"`
	PresubmitRejects7d   Counts               `json:"presubmitRejects7d"`
	TestRunFails1d       Counts               `json:"testRunFailures1d"`
	TestRunFails3d       Counts               `json:"testRunFailures3d"`
	TestRunFails7d       Counts               `json:"testRunFailures7d"`
	Failures1d           Counts               `json:"failures1d"`
	Failures3d           Counts               `json:"failures3d"`
	Failures7d           Counts               `json:"failures7d"`
	AffectedTests1d      []SubCluster         `json:"affectedTests1d"`
	AffectedTests3d      []SubCluster         `json:"affectedTests3d"`
	AffectedTests7d      []SubCluster         `json:"affectedTests7d"`
	ExampleFailureReason bigquery.NullString  `json:"exampleFailureReason"`
	ExampleTestID        string               `json:"exampleTestId"`
}

// SubCluster represents the name of a test and the number of times
// a failure has impacted it.
type SubCluster struct {
	Value     string `json:"value"`
	Num_Fails int    `json:"numFails"`
}

// Counts captures the values of an integer-valued metric in different
// calculation bases.
type Counts struct {
	// The statistic value after impact has been reduced by exoneration.
	Nominal int64 `json:"nominal"`
	// The statistic value before impact has been reduced by exoneration.
	PreExoneration int64 `json:"preExoneration"`
	// The statistic value:
	// - excluding impact already counted under other higher-priority clusters
	//   (I.E. bug clusters.)
	// - after impact has been reduced by exoneration.
	Residual int64 `json:"residual"`
	// The statistic value:
	// - excluding impact already counted under other higher-priority clusters
	//   (I.E. bug clusters.)
	// - before impact has been reduced by exoneration.
	ResidualPreExoneration int64 `json:"residualPreExoneration"`
}

// NewClient creates a new client for reading clusters. Close() MUST
// be called after you have finished using this client.
func NewClient(ctx context.Context, gcpProject string) (*Client, error) {
	client, err := bqutil.Client(ctx, gcpProject)
	if err != nil {
		return nil, err
	}
	return &Client{client: client}, nil
}

// Client may be used to read Weetbix clusters.
type Client struct {
	client *bigquery.Client
}

// Close releases any resources held by the client.
func (c *Client) Close() error {
	return c.client.Close()
}

// RebuildAnalysis re-builds the cluster summaries analysis from
// clustered test results.
func (c *Client) RebuildAnalysis(ctx context.Context, project string) error {
	datasetID, err := bqutil.DatasetForProject(project)
	if err != nil {
		return errors.Annotate(err, "getting dataset").Err()
	}
	dataset := c.client.Dataset(datasetID)

	dstTable := dataset.Table("cluster_summaries")

	q := c.client.Query(clusterSummariesAnalysis)
	q.DefaultDatasetID = dataset.DatasetID
	q.Dst = dstTable
	q.CreateDisposition = bigquery.CreateIfNeeded
	q.WriteDisposition = bigquery.WriteTruncate
	job, err := q.Run(ctx)
	if err != nil {
		return errors.Annotate(err, "starting cluster summary analysis").Err()
	}

	waitCtx, cancel := context.WithTimeout(ctx, time.Minute*5)
	defer cancel()

	js, err := job.Wait(waitCtx)
	if err != nil {
		return errors.Annotate(err, "waiting for cluster summary analysis to complete").Err()
	}
	if js.Err() != nil {
		return errors.Annotate(err, "cluster summary analysis failed").Err()
	}
	return nil
}

// ReadImpactfulClusters reads clusters exceeding specified impact metrics, or are otherwise
// nominated to be read.
func (c *Client) ReadImpactfulClusters(ctx context.Context, opts ImpactfulClusterReadOptions) ([]*ClusterSummary, error) {
	if opts.Thresholds == nil {
		return nil, errors.New("thresholds must be specified")
	}

	dataset, err := bqutil.DatasetForProject(opts.Project)
	if err != nil {
		return nil, errors.Annotate(err, "getting dataset").Err()
	}

	whereFailures, failuresParams := whereThresholdsExceeded("failures", opts.Thresholds.TestResultsFailed)
	whereTestRuns, testRunsParams := whereThresholdsExceeded("test_run_fails", opts.Thresholds.TestRunsFailed)
	wherePresubmits, presubmitParams := whereThresholdsExceeded("presubmit_rejects", opts.Thresholds.PresubmitRunsFailed)

	q := c.client.Query(`
		SELECT
			STRUCT(cluster_algorithm AS Algorithm, cluster_id as ID) as ClusterID,` +
		selectCounts("presubmit_rejects", "PresubmitRejects", "1d") +
		selectCounts("presubmit_rejects", "PresubmitRejects", "3d") +
		selectCounts("presubmit_rejects", "PresubmitRejects", "7d") +
		selectCounts("test_run_fails", "TestRunFails", "1d") +
		selectCounts("test_run_fails", "TestRunFails", "3d") +
		selectCounts("test_run_fails", "TestRunFails", "7d") +
		selectCounts("failures", "Failures", "1d") +
		selectCounts("failures", "Failures", "3d") +
		selectCounts("failures", "Failures", "7d") + `
			example_failure_reason.primary_error_message as ExampleFailureReason,
			example_test_id as ExampleTestID
		FROM ` + dataset + `.cluster_summaries
		WHERE (` + whereFailures + `) OR (` + whereTestRuns + `) OR (` + wherePresubmits + `)
			OR STRUCT(cluster_algorithm AS Algorithm, cluster_id as ID) IN UNNEST(@alwaysInclude)
		ORDER BY
			presubmit_rejects_residual_1d DESC,
			test_run_fails_residual_1d DESC,
			failures_residual_1d DESC
	`)

	params := []bigquery.QueryParameter{
		{
			Name:  "alwaysInclude",
			Value: opts.AlwaysInclude,
		},
	}
	params = append(params, failuresParams...)
	params = append(params, testRunsParams...)
	params = append(params, presubmitParams...)
	q.Parameters = params

	job, err := q.Run(ctx)
	if err != nil {
		return nil, errors.Annotate(err, "querying cluster summaries").Err()
	}
	it, err := job.Read(ctx)
	if err != nil {
		return nil, errors.Annotate(err, "obtain result iterator").Err()
	}
	clusters := []*ClusterSummary{}
	for {
		row := &ClusterSummary{}
		err := it.Next(row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, errors.Annotate(err, "obtain next cluster summary row").Err()
		}
		clusters = append(clusters, row)
	}
	return clusters, nil
}

func valueOrDefault(value *int64, defaultValue int64) int64 {
	if value != nil {
		return *value
	}
	return defaultValue
}

// selectCounts generates SQL to select a set of Counts.
func selectCounts(sqlPrefix, fieldPrefix, suffix string) string {
	return `STRUCT(` +
		sqlPrefix + `_` + suffix + ` AS Nominal,` +
		sqlPrefix + `_pre_exon_` + suffix + ` AS PreExoneration,` +
		sqlPrefix + `_residual_` + suffix + ` AS Residual,` +
		sqlPrefix + `_residual_pre_exon_` + suffix + ` AS ResidualPreExoneration` +
		`) AS ` + fieldPrefix + suffix + `,`
}

// whereThresholdsExceeded generates a SQL Where clause to query
// where a particular metric exceeds a given threshold.
func whereThresholdsExceeded(sqlPrefix string, threshold *config.MetricThreshold) (string, []bigquery.QueryParameter) {
	if threshold == nil {
		threshold = &config.MetricThreshold{}
	}
	sql := sqlPrefix + "_residual_1d > @" + sqlPrefix + "_1d OR " +
		sqlPrefix + "_residual_3d > @" + sqlPrefix + "_3d OR " +
		sqlPrefix + "_residual_7d > @" + sqlPrefix + "_7d"
	parameters := []bigquery.QueryParameter{
		{
			Name:  sqlPrefix + "_1d",
			Value: valueOrDefault(threshold.OneDay, math.MaxInt64),
		},
		{
			Name:  sqlPrefix + "_3d",
			Value: valueOrDefault(threshold.ThreeDay, math.MaxInt64),
		},
		{
			Name:  sqlPrefix + "_7d",
			Value: valueOrDefault(threshold.SevenDay, math.MaxInt64),
		},
	}
	return sql, parameters
}

// ReadCluster reads information about a single cluster.
func (c *Client) ReadCluster(ctx context.Context, luciProject string, clusterID clustering.ClusterID) (*ClusterSummary, error) {
	dataset, err := bqutil.DatasetForProject(luciProject)
	if err != nil {
		return nil, errors.Annotate(err, "getting dataset").Err()
	}

	q := c.client.Query(`
		SELECT
			STRUCT(cluster_algorithm AS Algorithm, cluster_id as ID) as ClusterID,` +
		selectCounts("presubmit_rejects", "PresubmitRejects", "1d") +
		selectCounts("presubmit_rejects", "PresubmitRejects", "3d") +
		selectCounts("presubmit_rejects", "PresubmitRejects", "7d") +
		selectCounts("test_run_fails", "TestRunFails", "1d") +
		selectCounts("test_run_fails", "TestRunFails", "3d") +
		selectCounts("test_run_fails", "TestRunFails", "7d") +
		selectCounts("failures", "Failures", "1d") +
		selectCounts("failures", "Failures", "3d") +
		selectCounts("failures", "Failures", "7d") + `
			affected_tests_1d as AffectedTests1d,
			affected_tests_3d as AffectedTests3d,
			affected_tests_7d as AffectedTests7d,
			example_failure_reason.primary_error_message as ExampleFailureReason,
			example_test_id as ExampleTestID
		FROM ` + dataset + `.cluster_summaries
		WHERE cluster_algorithm = @clusterAlgorithm
		  AND cluster_id = @clusterID
	`)
	q.Parameters = []bigquery.QueryParameter{
		{Name: "clusterAlgorithm", Value: clusterID.Algorithm},
		{Name: "clusterID", Value: clusterID.ID},
	}
	job, err := q.Run(ctx)
	if err != nil {
		return nil, errors.Annotate(err, "querying cluster summary").Err()
	}
	it, err := job.Read(ctx)
	if err != nil {
		return nil, errors.Annotate(err, "obtain result iterator").Err()
	}
	clusters := []*ClusterSummary{}
	for {
		row := &ClusterSummary{}
		err := it.Next(row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, errors.Annotate(err, "obtain next cluster summary row").Err()
		}
		clusters = append(clusters, row)
	}
	if len(clusters) == 0 {
		return nil, fmt.Errorf("cluster %s not found", clusterID)
	}
	return clusters[0], nil
}

type ClusterFailure struct {
	Realm                      bigquery.NullString    `json:"realm"`
	IngestedInvocationID       bigquery.NullString    `json:"ingestedInvocationID"`
	TestID                     bigquery.NullString    `json:"testID"`
	Variant                    []*Variant             `json:"variant"`
	PresubmitRunID             *PresubmitRunID        `json:"presubmitRunID"`
	PartitionTime              bigquery.NullTimestamp `json:"partitionTime"`
	IsIncluded                 bigquery.NullBool      `json:"isIncluded"`
	IsIncludedWithHighPriority bigquery.NullBool      `json:"isIncludedWithHighPriority"`
	IsExonerated               bigquery.NullBool      `json:"isExonerated"`
}

type Variant struct {
	Key   bigquery.NullString `json:"key"`
	Value bigquery.NullString `json:"value"`
}

type PresubmitRunID struct {
	System bigquery.NullString `json:"system"`
	ID     bigquery.NullString `json:"id"`
}

// ReadClusterFailures reads the latest 2000 failures for a single cluster for the last 7 days.
func (c *Client) ReadClusterFailures(ctx context.Context, luciProject string, clusterID clustering.ClusterID) ([]*ClusterFailure, error) {
	dataset, err := bqutil.DatasetForProject(luciProject)
	if err != nil {
		return nil, errors.Annotate(err, "getting dataset").Err()
	}
	// TODO(mwarton): change this to read from materialized
	// table instead of directly from the view for better performance
	// (from 30 seconds -> 2 seconds).
	q := c.client.Query(`
		SELECT
			realm as Realm,
			ingested_invocation_id as IngestedInvocationID,
			test_id as TestID,
			variant as Variant,
			presubmit_run_id as PresubmitRunID,
			partition_time as PartitionTime,
			is_included as IsIncluded,
			is_included_with_high_priority as IsIncludedWithHighPriority,
			is_exonerated as IsExonerated
		FROM
			` + dataset + `.clustered_failures_latest_7d
		WHERE cluster_algorithm = @clusterAlgorithm
		  AND cluster_id = @clusterID
		ORDER BY partition_time DESC
		LIMIT 2000
	`)
	q.Parameters = []bigquery.QueryParameter{
		{Name: "clusterAlgorithm", Value: clusterID.Algorithm},
		{Name: "clusterID", Value: clusterID.ID},
	}
	job, err := q.Run(ctx)
	if err != nil {
		return nil, errors.Annotate(err, "querying cluster failures").Err()
	}
	it, err := job.Read(ctx)
	if err != nil {
		return nil, errors.Annotate(err, "obtain result iterator").Err()
	}
	failures := []*ClusterFailure{}
	for {
		row := &ClusterFailure{}
		err := it.Next(row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, errors.Annotate(err, "obtain next cluster failure row").Err()
		}
		failures = append(failures, row)
	}
	return failures, nil
}
