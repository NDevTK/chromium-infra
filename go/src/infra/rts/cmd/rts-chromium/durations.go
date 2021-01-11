// Copyright 2020 The LUCI Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/api/iterator"

	"go.chromium.org/luci/common/errors"

	evalpb "infra/rts/presubmit/eval/proto"
)

// durations calls f for each found test duration.
func (r *presubmitHistoryRun) durations(ctx context.Context, f func(*evalpb.TestDuration) error) error {
	bq, err := r.bqClient(ctx)
	if err != nil {
		return errors.Annotate(err, "failed to init BigQuery client").Err()
	}

	q := bq.Query(testDurationsSQL)
	q.Parameters = []bigquery.QueryParameter{
		{Name: "startTime", Value: r.startTime},
		{Name: "endTime", Value: r.endTime},
		{Name: "frac", Value: r.durationDataFrac},
	}
	it, err := q.Read(ctx)
	if err != nil {
		return err
	}

	for {
		var row durationRow
		switch err := it.Next(&row); {
		case err == iterator.Done:
			return nil
		case err != nil:
			return err
		default:
			if err := f(row.proto()); err != nil {
				return err
			}
		}
	}
}

type durationRow struct {
	Change       int
	Patchset     int
	TestID       string
	TestFileName string
	Duration     float64
}

func (r *durationRow) proto() *evalpb.TestDuration {
	return &evalpb.TestDuration{
		Patchsets: []*evalpb.GerritPatchset{
			{
				Change: &evalpb.GerritChange{
					// Assume all Chromium source code is in
					// https://chromium.googlesource.com/chromium/src
					// TOOD(nodir): make it fail if it is not.
					Host:    "chromium-review.googlesource.com",
					Project: "chromium/src",
					Number:  int64(r.Change),
				},
				Patchset: int64(r.Patchset),
			},
		},
		Test:     &evalpb.Test{Id: r.TestID, FileName: r.TestFileName},
		Duration: ptypes.DurationProto(time.Duration(r.Duration * 1e9)),
	}
}

const testDurationsSQL = `
WITH
	tryjobs AS (
		SELECT
			b.id,
			ps,
		FROM commit-queue.chromium.attempts a, a.gerrit_changes ps, a.builds b
		WHERE partition_time BETWEEN @startTime AND TIMESTAMP_ADD(@endTime, INTERVAL 1 DAY)
	),
	test_results AS (
		SELECT
			CAST(REGEXP_EXTRACT(exported.id, r'build-(\d+)') as INT64) build_id,
			test_location.file_name,
			test_id,
			duration,
		FROM luci-resultdb.chromium.try_test_results
		WHERE partition_time BETWEEN @startTime and @endTime
			AND RAND() <= @frac
			AND duration > 0

			-- Exclude broken test locations.
			-- TODO(nodir): remove this after crbug.com/1130425 is fixed.
			AND REGEXP_CONTAINS(test_location.file_name, r'(?i)\.(cc|html|m|c|cpp)$')
			-- Exclude broken prefixes.
			-- TODO(nodir): remove after crbug.com/1017288 is fixed.
			AND (test_id NOT LIKE 'ninja://:blink_web_tests/%' OR test_location.file_name LIKE '//third_party/%')
	)
SELECT
	ps.change as Change,
	ps.patchset as Patchset,
	test_id as TestID,
	file_name as TestFileName,
	duration as Duration
FROM tryjobs t
JOIN test_results tr ON t.id = tr.build_id
`
