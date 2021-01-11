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

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/api/iterator"

	evalpb "infra/rts/presubmit/eval/proto"
)

// rejections calls f for each found CQ rejection.
func (r *presubmitHistoryRun) rejections(ctx context.Context, f func(*evalpb.Rejection) error) error {
	q, err := r.bqQuery(ctx, rejectedPatchSetsSQL)
	if err != nil {
		return err
	}
	it, err := q.Read(ctx)
	if err != nil {
		return err
	}

	// Each result row represents a failed test variant.
	// The rows are sorted by change and patchset, and one rejection
	// may span multiple rows.
	// Accumulate failed tests to the current rejection until we encounter
	// a different patchset.
	var curRej *evalpb.Rejection
	for {
		var row rejectionRow
		switch err := it.Next(&row); {
		case err == iterator.Done:
			// Report the last one.
			if curRej != nil {
				return f(curRej)
			}
			return nil

		case err != nil:
			return err
		}

		// Report the current rejection if we have one and if `row` is a different
		// patchset.
		if curRej != nil && (curRej.Patchsets[0].Change.Number != int64(row.Change) || curRej.Patchsets[0].Patchset != int64(row.Patchset)) {
			if err := f(curRej); err != nil {
				return err
			}
			curRej = nil
		}

		// Initialize a new rejection if we don't have one.
		if curRej == nil {
			curRej = &evalpb.Rejection{
				Patchsets: []*evalpb.GerritPatchset{
					{
						Change: &evalpb.GerritChange{
							// Assume all Chromium source code is in
							// https://chromium.googlesource.com/chromium/src
							// TOOD(nodir): make it fail if it is not.
							Host:    "chromium-review.googlesource.com",
							Project: "chromium/src",
							Number:  int64(row.Change),
						},
						Patchset: int64(row.Patchset),
					},
				},
			}
			curRej.Timestamp, _ = ptypes.TimestampProto(row.Timestamp)
		}

		curRej.FailedTestVariants = append(curRej.FailedTestVariants, row.FailedTestVariant.proto())
	}
}

type rejectionRow struct {
	Change            int
	Patchset          int
	Timestamp         time.Time
	FailedTestVariant testVariantRow
}

type testVariantRow struct {
	ID       string
	FileName string
	Variant  []struct {
		Key   string
		Value string
	}
}

func (t *testVariantRow) proto() *evalpb.TestVariant {
	ret := &evalpb.TestVariant{
		Id:       t.ID,
		FileName: t.FileName,
		Variant:  make(map[string]string, len(t.Variant)),
	}
	for _, kv := range t.Variant {
		ret.Variant[kv.Key] = kv.Value
	}
	return ret
}

// rejectedPatchSetsSQL is a BigQuery query that returns patchsets with test
// failures. Ignores tests that don't have a test location.
//
// This query intentionally does not use ARRAY_AGG for failed tests because
// otherwise a result row size may exceed the API limitation of 10MB.
const rejectedPatchSetsSQL = `
	WITH
		tryjobs AS (
			SELECT
				b.id,
				ps,
				partition_time as ps_approx_timestamp,
			FROM commit-queue.chromium.attempts a, a.gerrit_changes ps, a.builds b
			-- Read next-day attempts too in case the CQ attempt finished soon after 12am.
			-- Note that for test results, the cutoff is still @endTime.
			WHERE partition_time BETWEEN @startTime AND TIMESTAMP_ADD(@endTime, INTERVAL 1 DAY)
		),
		failed_test_variants AS (
			SELECT
				CAST(REGEXP_EXTRACT(exported.id, r'build-(\d+)') as INT64) as build_id,
				ANY_VALUE(IFNULL(test_location.file_name, '')) as file_name,
				test_id,
				variant_hash,
				ANY_VALUE(variant) as variant,
			FROM luci-resultdb.chromium.try_test_results tr
			WHERE partition_time BETWEEN @startTime and @endTime
				AND (@test_id_regexp = '' OR REGEXP_CONTAINS(test_id, @test_id_regexp))
				AND (@builder_regexp = '' OR EXISTS (SELECT 0 FROM tr.variant WHERE key='builder' AND REGEXP_CONTAINS(value, @builder_regexp)))
				-- Exclude broken test locations.
				-- TODO(nodir): remove this after crbug.com/1130425 is fixed.
				AND REGEXP_CONTAINS(IFNULL(test_location.file_name, ''), r'(?i)^(|.*\.(cc|html|m|c|cpp))$')
				-- Exclude broken prefixes.
				-- TODO(nodir): remove after crbug.com/1017288 is fixed.
				AND (test_id NOT LIKE 'ninja://:blink_web_tests/%' OR test_location.file_name LIKE '//third_party/%')
			GROUP BY build_id, test_id, variant_hash
			-- TODO(crbug.com/1112125): consider doing this filtering after joining with patchsets.
			-- This will increase the query cost, but might improve
			-- data quality. Currenly the query is vulnerable to situations where
			-- a test failed in one CQ attempt for the patchset, and succeeded
			-- in another for the same patchset.
			HAVING LOGICAL_AND(NOT expected AND NOT exonerated)
		)
	SELECT
		ps.change as Change,
		ps.patchset as Patchset,
		MIN(ps_approx_timestamp) as Timestamp,
		STRUCT(
			test_id as ID,
			ANY_VALUE(variant) as Variant,
			ANY_VALUE(file_name) as FileName
		) as FailedTestVariant
	FROM tryjobs t
	JOIN failed_test_variants f ON t.id = f.build_id
	GROUP BY ps.change, ps.patchset, test_id, variant_hash
	ORDER BY ps.change, ps.patchset
`
