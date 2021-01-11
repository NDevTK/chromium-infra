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

package eval

import (
	"context"
	"math"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	evalpb "infra/rts/presubmit/eval/proto"
)

// Efficiency is result of evaluation how much compute time the RTS algorithm
// could save.
type Efficiency struct {
	// SampleDuration is the sum of test durations in the analyzed data sample.
	SampleDuration time.Duration

	// ForecastDuration is the sum of test durations for tests selected by the RTS
	// algorithm. It is a value between 0 and SampleDuration.
	// The lower the number the better.
	ForecastDuration time.Duration
}

// Score returns the efficiency score.
// May return NaN.
func (e *Efficiency) Score() float64 {
	if e.SampleDuration == 0 {
		return math.NaN()
	}
	saved := e.SampleDuration - e.ForecastDuration
	return float64(100*saved) / float64(e.SampleDuration)
}

func (r *evalRun) evaluateEfficiency(ctx context.Context, durationC <-chan *evalpb.TestDuration) (*Efficiency, error) {
	eg, ctx := errgroup.WithContext(ctx)

	// Run the algorithm in r.Concurrency goroutines.
	var ret Efficiency
	var mu sync.Mutex
	for i := 0; i < r.Concurrency; i++ {
		eg.Go(func() error {
			in := Input{TestVariants: make([]*evalpb.TestVariant, 1)}
			for td := range durationC {
				changedFiles, err := r.changedFiles(ctx, td.Patchsets...)
				switch {
				case err != nil:
					return err
				case len(changedFiles) == 0:
					continue // Ineligible.
				}

				// Run the algorithm.
				in.ChangedFiles = changedFiles
				in.TestVariants[0] = td.TestVariant
				out, err := r.Algorithm(ctx, in)
				if err != nil {
					return err
				}

				// Record results.
				mu.Lock()
				dur := td.Duration.AsDuration()
				ret.SampleDuration += dur
				if out.ShouldRunAny {
					ret.ForecastDuration += dur
				}
				mu.Unlock()

				r.progress.UpdateCurrentEfficiency(ctx, ret)
			}
			return ctx.Err()
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return &ret, nil
}
