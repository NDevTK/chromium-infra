// Copyright 2017 The LUCI Authors.
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

package analysis

import (
	"fmt"
	"time"

	"go.chromium.org/luci/buildbucket"
)

type groupKey struct {
	// do not put interfaces in this struct,
	// because it is used as a map key.

	buildbucket.GerritChange
	GotRevision string
}

func (k *groupKey) String() string {
	return fmt.Sprintf("%s @ %q", &k.GerritChange, k.GotRevision)
}

// group is two sets of builds, for Buildbot and LUCI, that should have the same
// results.
type group struct {
	Key            groupKey
	LUCI, Buildbot groupSide
}

// trustworthy returns true if g can be used for correctness analysis.
func (g *group) trustworthy() bool {
	return g.Buildbot.trustworthy() && g.LUCI.trustworthy()
}

// build contains minimal information needed for analysis.
type build struct {
	Status         buildbucket.Status
	CreationTime   time.Time
	CompletionTime time.Time
	RunDuration    time.Duration
	URL            string
}

// groupSide is a list of builds ordered from oldest to newest
type groupSide []*build

func (s groupSide) avgRunDuration() time.Duration {
	avg := time.Duration(0)
	count := 0
	for _, b := range s {
		avg += b.RunDuration
		count++
	}
	if count == 0 {
		return 0
	}
	return avg / time.Duration(count)
}

// MostRecentlyCompleted returns completion time of the most recently created
// build.
func (s groupSide) MostRecentlyCompleted() time.Time {
	if len(s) == 0 {
		return time.Time{}
	}
	return s[len(s)-1].CompletionTime
}

// success returns true if at least one build succeeded, otherwise false.
func (s groupSide) success() bool {
	for _, b := range s {
		if b.Status == buildbucket.StatusSuccess {
			return true
		}
	}
	return false
}

// trustworthy returns true if s can be used for correctness analysis.
func (s groupSide) trustworthy() bool {
	if s.success() {
		return true
	}

	// If there are no successful builds and fewer than 3 trustworthy failures,
	// consider this result too vulnerable to flakes.
	failures := 0
	for _, b := range s {
		if b.Status == buildbucket.StatusFailure {
			failures++
			if failures >= 3 {
				return true
			}
		}
	}
	return false
}
