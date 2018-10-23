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

// TODO(akeshet): The tests in this file make use of a lot of unexported
// methods and fields. It would be better if they were rewritten to use
// only the exported API of the scheduler. That would also entail building
// an exported API for getting job prioritization.

package scheduler

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/kylelemons/godebug/pretty"

	"infra/qscheduler/qslib/tutils"
	"infra/qscheduler/qslib/types/account"
	"infra/qscheduler/qslib/types/vector"

	. "github.com/smartystreets/goconvey/convey"
)

// epoch is an arbitrary time for testing purposes, corresponds to
// 01/01/2018 @ 1:00 am UTC
var epoch = time.Unix(1514768400, 0)

// TestPrioritizeOne tests that PrioritizeRequests behaves correctly
// for a single request.
func TestPrioritizeOne(t *testing.T) {

	aid := "a1"
	rid := "r1"

	Convey("Given a scheduler, with a request for account", t, func() {
		tm := time.Unix(0, 0)
		s := New()
		s.AddRequest(rid, &TaskRequest{AccountId: aid}, tm)

		accountCases := []struct {
			accountBalance   *vector.Vector
			expectedPriority int32
		}{
			{
				vector.New(1),
				0,
			},
			{
				vector.New(0, 1),
				1,
			},
			{
				vector.New(),
				account.FreeBucket,
			},
		}

		for _, c := range accountCases {
			Convey(fmt.Sprintf("given account balance is %v", c.accountBalance.Values), func() {
				s.AddAccount(aid, &account.Config{}, c.accountBalance)
				Convey("when prioritizing", func() {
					got := s.prioritizeRequests()
					Convey(fmt.Sprintf("then the request gets priority %d.", c.expectedPriority), func() {
						So(got[0].Priority, ShouldEqual, c.expectedPriority)
						So(got[0].RequestId, ShouldEqual, rid)
					})
				})
			})
		}

		Convey("when no such account exists", func() {
			Convey("when prioritizing", func() {
				got := s.prioritizeRequests()
				Convey("then the request is put in free bucket priority.", func() {
					So(got[0].Priority, ShouldEqual, account.FreeBucket)
					So(got[0].RequestId, ShouldEqual, rid)
				})
			})
		})
	})
}

// TestPrioritizeMany tests that PrioritizeRequests behaves correctly
// for a number of requests.
func TestPrioritizeMany(t *testing.T) {
	t.Parallel()
	nReqs := 10
	aid := "a1"
	Convey("Given requests with different enqueue times, but inserted in random order", t, func() {
		s := New()
		// Use a fixed seed, so the test is reproducible and the request order is
		// pseudo-random.
		rand.Seed(10)
		perm := rand.Perm(nReqs)
		for _, i := range perm {
			tm := time.Unix(int64(i), 0)
			s.AddRequest(strconv.Itoa(i), &TaskRequest{AccountId: aid, EnqueueTime: tutils.TimestampProto(tm)}, tm)
		}

		Convey("given no matching account", func() {
			Convey("when prioritizing", func() {
				got := s.prioritizeRequests()
				Convey("then requests are prioritized by enqueue time.", func() {
					times := make([]time.Time, nReqs)
					for i, g := range got {
						times[i] = tutils.Timestamp(g.Request.EnqueueTime)
					}
					So(times, ShouldBeChronological)
					So(got, ShouldHaveLength, nReqs)
				})
			})
		})

		Convey("given an account with no maximum fanout", func() {
			s.AddAccount(aid, account.NewConfig(), nil)
			Convey("when prioritizing", func() {
				got := s.prioritizeRequests()
				Convey("then requests are prioritized by enqueue time.", func() {
					times := make([]time.Time, nReqs)
					for i, g := range got {
						times[i] = tutils.Timestamp(g.Request.EnqueueTime)
					}
					So(times, ShouldBeChronological)
					So(got, ShouldHaveLength, nReqs)
				})
			})
		})

		Convey("given the account specifices a maximum fanout and some requests for that account are already running", func() {
			maxFanout := 5
			s.AddAccount(aid, &account.Config{MaxFanout: int32(maxFanout)}, vector.New(0, 1))

			// Two requests are already running.
			s.AddRequest("11", &TaskRequest{AccountId: aid}, time.Unix(0, 0))
			s.AddRequest("12", &TaskRequest{AccountId: aid}, time.Unix(0, 0))
			s.MarkIdle("11", []string{}, time.Unix(0, 0))
			s.MarkIdle("12", []string{}, time.Unix(0, 0))
			s.state.applyAssignment(&Assignment{Priority: 1, RequestId: "11", WorkerId: "11", Type: Assignment_IDLE_WORKER})
			s.state.applyAssignment(&Assignment{Priority: 1, RequestId: "12", WorkerId: "12", Type: Assignment_IDLE_WORKER})

			Convey("when prioritizing", func() {
				got := s.prioritizeRequests()
				Convey("then requests beyond the maximum fanout are put in the free bucket priority.", func() {
					// First three get the account's available bucket (P1), the remaining
					// ones get the free bucket.
					fb := account.FreeBucket
					expectedPriorities := []int32{1, 1, 1, fb, fb, fb, fb, fb, fb, fb}

					actualPriorities := make([]int32, nReqs)
					for i := 0; i < nReqs; i++ {
						actualPriorities[i] = got[i].Priority
					}
					So(actualPriorities, ShouldResemble, expectedPriorities)
					So(got, ShouldHaveLength, nReqs)
				})
			})
		})
	})
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
