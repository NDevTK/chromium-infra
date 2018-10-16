// Copyright 2018 The LUCI Authors.
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

package scheduler

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/kylelemons/godebug/pretty"

	"infra/qscheduler/qslib/tutils"
	"infra/qscheduler/qslib/types/account"
	"infra/qscheduler/qslib/types/vector"
)

// epoch is an arbitrary time for testing purposes, corresponds to
// 01/01/2018 @ 1:00 am UTC
var epoch = time.Unix(1514768400, 0)

// TestBasicPrioritization tests that PrioritizeRequests behaves correctly
// for simple cases.
func TestBasicPrioritization(t *testing.T) {
	t.Parallel()

	cases := []struct {
		s        *Scheduler
		expected orderedRequests
	}{
		// One request with quota, should be given appropriate priority.
		{
			&Scheduler{
				&State{
					Balances: map[string]*vector.Vector{"a1": vector.New(1, 0, 0)},
					QueuedRequests: map[string]*TaskRequest{
						"t1": &TaskRequest{AccountId: "a1"},
					},
				},
				NewConfig(),
			},
			[]prioritizedRequest{{RequestId: "t1", Priority: 0, Request: &TaskRequest{AccountId: "a1"}}},
		},

		// One request without quota, should be in the FreeBucket.
		{
			&Scheduler{
				&State{
					QueuedRequests: map[string]*TaskRequest{"t1": &TaskRequest{}},
				},
				NewConfig(),
			},
			[]prioritizedRequest{{RequestId: "t1", Priority: account.FreeBucket, Request: &TaskRequest{}}},
		},
	}

	for i, test := range cases {
		actual := test.s.prioritizeRequests()

		if diff := pretty.Compare(actual, test.expected); diff != "" {
			t.Errorf(fmt.Sprintf("Case %d got unexpected slice diff (-got +want): %s", i, diff))
		}
	}
}

// TestPrioritizeWithEnqueueTimeTieBreaker tests that PrioritizeRequests
// behaves correctly in circumstances where the enqueue time of a request
// is used as tiebreaker.
func TestPrioritizeWithEnqueueTimeTieBreaker(t *testing.T) {
	t.Parallel()
	e := time.Unix(100, 100)
	l := e.Add(10 * time.Second)

	eT := tutils.TimestampProto(e)
	lT := tutils.TimestampProto(l)

	eR := TaskRequest{AccountId: "a1", EnqueueTime: eT}
	lR := TaskRequest{AccountId: "a1", EnqueueTime: lT}

	state := &State{
		Balances: map[string]*vector.Vector{"a1": vector.New(1, 0, 0)},
		QueuedRequests: map[string]*TaskRequest{
			"t2": &lR,
			"t1": &eR,
		},
	}
	actual := (&Scheduler{state, NewConfig()}).prioritizeRequests()
	expected := orderedRequests([]prioritizedRequest{
		{RequestId: "t1", Priority: 0, Request: &eR},
		{RequestId: "t2", Priority: 0, Request: &lR},
	})

	if diff := pretty.Compare(actual, expected); diff != "" {
		t.Errorf(fmt.Sprintf("Got unexpected diff (-got +want): %s", diff))
	}
}

// TestDemoteBeyondFanout tests that for a given account, once the number
// of running or requested tasks exceeds that account's MaxFanout, further
// requests are assigned to the FreeBucket.
func TestDemoteBeyondFanout(t *testing.T) {
	t.Parallel()
	config := &Config{
		AccountConfigs: map[string]*account.Config{
			"a1": {MaxFanout: 3},
			"a2": {},
		},
	}
	running := []*TaskRun{
		{RequestId: "1", Priority: 0, Request: &TaskRequest{AccountId: "a1"}},
		{RequestId: "2", Priority: 0, Request: &TaskRequest{AccountId: "a1"}},
		{RequestId: "3", Priority: 0, Request: &TaskRequest{AccountId: "a2"}},
		{RequestId: "4", Priority: account.FreeBucket, Request: &TaskRequest{AccountId: "a3"}},
	}
	workers := getWorkers(running)

	r1 := TaskRequest{AccountId: "a1"}
	r2 := TaskRequest{AccountId: "a1"}
	r3 := TaskRequest{AccountId: "a2"}
	r4 := TaskRequest{AccountId: "a3"}
	reqs := map[string]*TaskRequest{
		"5": &r1,
		"6": &r2,
		"7": &r3,
		"8": &r4,
	}
	state := &State{
		Balances: map[string]*vector.Vector{
			"a1": {},
			"a2": {},
		},
		QueuedRequests: reqs,
		Workers:        workers,
	}

	priList := []prioritizedRequest{
		{RequestId: "5", Priority: 0, Request: &r1},
		{RequestId: "6", Priority: 0, Request: &r2},
		{RequestId: "7", Priority: 0, Request: &r3},
		{RequestId: "8", Priority: account.FreeBucket, Request: &r4},
	}

	expected := []prioritizedRequest{
		{RequestId: "5", Priority: 0, Request: &r1},
		// This request got demoted from P0 to FreeBucket because it
		// exceeded the account's max fanout.
		{RequestId: "6", Priority: account.FreeBucket, Request: &r2},
		{RequestId: "7", Priority: 0, Request: &r3},
		{RequestId: "8", Priority: account.FreeBucket, Request: &r4},
	}

	demoteTasksBeyondFanout(priList, state, config)

	actual := priList
	if diff := pretty.Compare(actual, expected); diff != "" {
		t.Errorf(fmt.Sprintf("Got unexpected diff (-got +want): %s", diff))
	}
}

// TestPrioritize tests the full behavior set of PrioritizeRequests.
func TestPrioritize(t *testing.T) {
	t.Parallel()
	// Setup common variables.
	a1 := "a1"
	a2 := "a2"
	a3 := "a3"
	a4 := "a4"
	// a1: Account with P0 quota, fanout limit 3.
	// a2: Account with P1 quota, no fanout limit.
	// a3: Account with no quota.
	// a4: Invalid / nonexistent account.
	balances := map[string]*vector.Vector{
		a1: vector.New(1, 0, 0),
		a2: vector.New(0, 1, 0),
		a3: vector.New(),
	}
	config := &Config{
		AccountConfigs: map[string]*account.Config{
			a1: &account.Config{MaxFanout: 3},
			a2: &account.Config{},
			a3: &account.Config{},
		},
	}

	// 6 Jobs are already running. 2 for A1, 2 for A2, 1 for each of A3, A4
	run1 := TaskRun{Priority: 0, Request: &TaskRequest{AccountId: a1}}
	run2 := TaskRun{Priority: 0, Request: &TaskRequest{AccountId: a1}}
	run3 := TaskRun{Priority: 1, Request: &TaskRequest{AccountId: a2}}
	run4 := TaskRun{Priority: 1, Request: &TaskRequest{AccountId: a2}}
	run5 := TaskRun{Priority: 3, Request: &TaskRequest{AccountId: a3}}
	run6 := TaskRun{Priority: 3, Request: &TaskRequest{AccountId: a4}}
	running := []*TaskRun{
		&run1,
		&run2,
		&run3,
		&run4,
		&run5,
		&run6,
	}

	// 6 Jobs are requested. 3 for A1, 1 for each of the remaining
	// A3's requests are the earliest, and 1 second apart.
	req1 := TaskRequest{AccountId: a1, EnqueueTime: atTime(0)}
	req2 := TaskRequest{AccountId: a1, EnqueueTime: atTime(1)}
	req3 := TaskRequest{AccountId: a1, EnqueueTime: atTime(2)}
	// The remaining requests are later by 1 second each.
	req4 := TaskRequest{AccountId: a2, EnqueueTime: atTime(3)}
	req5 := TaskRequest{AccountId: a3, EnqueueTime: atTime(4)}
	req6 := TaskRequest{AccountId: a4, EnqueueTime: atTime(5)}

	reqs := map[string]*TaskRequest{
		"1": &req1,
		"2": &req2,
		"3": &req3,
		"4": &req4,
		"5": &req5,
		"6": &req6,
	}

	state := &State{
		Balances:       balances,
		QueuedRequests: reqs,
		Workers:        getWorkers(running),
	}

	expected := orderedRequests([]prioritizedRequest{
		// A1 gets one additional request at P0, prior to overflowing fanout.
		{RequestId: "1", Priority: 0, Request: &req1},
		// A2 gets a P1 request.
		{RequestId: "4", Priority: 1, Request: &req4},
		// Remaining requests are all in the free bucket, ordered by enqueue time.
		{RequestId: "2", Priority: account.FreeBucket, Request: &req2},
		{RequestId: "3", Priority: account.FreeBucket, Request: &req3},
		{RequestId: "5", Priority: account.FreeBucket, Request: &req5},
		{RequestId: "6", Priority: account.FreeBucket, Request: &req6},
	})

	actual := (&Scheduler{state, config}).prioritizeRequests()

	if diff := pretty.Compare(actual, expected); diff != "" {
		t.Errorf(fmt.Sprintf("Got unexpected diff (-got +want): %s", diff))
	}
}

// TestForPriority tests that ForPriority method returns the correct
// sub-slices of a prioritized list.
func TestForPriority(t *testing.T) {
	t.Parallel()
	pRequests := orderedRequests([]prioritizedRequest{
		prioritizedRequest{Priority: 0},
		prioritizedRequest{Priority: 0},
		prioritizedRequest{Priority: 1},
		prioritizedRequest{Priority: 3},
		prioritizedRequest{Priority: 3},
		prioritizedRequest{Priority: 4},
	})

	expecteds := []orderedRequests{
		pRequests[0:2],
		pRequests[2:3],
		pRequests[3:3],
		pRequests[3:5],
		pRequests[5:6],
		pRequests[6:6],
	}

	for p := int32(0); p < 6; p++ {
		actual := pRequests.forPriority(p)
		expected := expecteds[p]
		if diff := pretty.Compare(actual, expected); diff != "" {
			t.Errorf(fmt.Sprintf("P%d slice got unexpected diff (-got +want): %s", p, diff))
		}
	}
}

// atTime is a helper method to create proto.Timestamp objects at various
// times relative to a fixed "0" time.
func atTime(seconds time.Duration) *timestamp.Timestamp {
	timeAfter := epoch.Add(seconds * time.Second)
	return tutils.TimestampProto(timeAfter)
}

// getWorkers is a helper function to turn a slice of running tasks
// into a workers map.
func getWorkers(running []*TaskRun) map[string]*Worker {
	workers := make(map[string]*Worker)
	for i, r := range running {
		wid := fmt.Sprintf("w%d", i)
		workers[wid] = &Worker{Labels: []string{}, RunningTask: r}
	}
	return workers
}
