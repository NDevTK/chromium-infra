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
	"math"
	"math/rand"
	"time"

	"infra/rts/presubmit/eval"
)

func main() {
	ctx := context.Background()
	rand.Seed(time.Now().Unix())
	eval.Main(ctx, func(ctx context.Context, in eval.Input) (eval.Output, error) {
		if len(in.TestVariants) > 32 {
			return eval.Output{ShouldRunAny: true}, nil
		}
		oneOf := int(math.Pow(2, float64(len(in.TestVariants))))
		return eval.Output{
			ShouldRunAny: rand.Intn(oneOf) == 0,
		}, nil
	})
}
