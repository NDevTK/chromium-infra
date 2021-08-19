// Copyright 2021 The Chromium Authors.
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
package cli

import (
	"context"
	"infra/chromeperf/pinpoint/proto"
	"path/filepath"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAnalyzeTelemetryExperiment(t *testing.T) {
	t.Parallel()
	// TODO: add tests for unhappy paths with errors
	// TODO: add more fine-grained unit tests for processing in-memory data
	// structures without requiring files
	Convey("Given a telemetry experiment manifest with known significant differences", t, func() {
		m, err := loadManifestFromPath("testdata/11ac8128320000/manifest.yaml")
		So(err, ShouldBeNil)
		Convey("When we analyze the artifacts", func() {
			rootDir, err := filepath.Abs("testdata/11ac8128320000")
			So(err, ShouldBeNil)
			r, err := analyzeExperiment(m, rootDir)
			So(err, ShouldBeNil)
			So(r, ShouldNotBeNil)
			Convey("Then we find the overall p-value", func() {
				So(r.OverallPValue, ShouldAlmostEqual, 0.0)
			})
			Convey("And we see that some p-values are less than 0.05", func() {
				less := 0
				for _, s := range r.Reports {
					if s.PValue != nil && *s.PValue < 0.05 {
						less += 1
					}
				}
				So(less, ShouldBeGreaterThan, 0)
			})
			Convey("And we see the stats for all metrics", func() {
				nan := 0
				withStats := 0
				for _, s := range r.Reports {
					if s.PValue == nil {
						nan += 1
						continue
					}
					if len(s.Measurements) > 1 {
						withStats += 1
					}
				}
				So(nan, ShouldNotEqual, 0)
				So(withStats, ShouldNotEqual, 0)
				So(nan+withStats, ShouldEqual, len(r.Reports))
			})
		})
		Convey("When we use the mixin to analyze the artifacts", func() {
			m := &analyzeExperimentMixin{analyzeExperiment: true}
			ctx := context.Background()
			// This is the minimal Job definition that's associated with the
			// manifest/artifacts in the test data.
			j := &proto.Job{
				Name: "jobs/legacy-11ac8128320000",
				JobSpec: &proto.JobSpec{
					JobKind: &proto.JobSpec_Experiment{
						Experiment: &proto.Experiment{},
					},
					Arguments: &proto.JobSpec_TelemetryBenchmark{
						TelemetryBenchmark: &proto.TelemetryBenchmark{},
					},
				},
			}
			wd, err := filepath.Abs("testdata")
			So(err, ShouldBeNil)
			r, err := m.doAnalyzeExperiment(ctx, wd, j)
			So(err, ShouldBeNil)
			So(r, ShouldNotBeNil)
			So(r.OverallPValue, ShouldNotEqual, 0)
			So(r.Reports, ShouldNotBeNil)
		})
	})

}
