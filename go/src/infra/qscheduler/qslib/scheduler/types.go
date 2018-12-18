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
	"time"

	"infra/qscheduler/qslib/tutils"
)

// NewConfig creates an returns a new Config instance with all maps initialized.
func NewConfig() *Config {
	return &Config{
		AccountConfigs: map[string]*AccountConfig{},
	}
}

// newState creates an returns a new State instance with all maps initialized.
func newState(t time.Time) *state {
	return &state{
		balances:       map[string]balance{},
		queuedRequests: map[string]*request{},
		workers:        map[string]*worker{},
		lastUpdateTime: t,
	}
}

func newStateFromProto(sp *StateProto) *state {
	s := &state{}
	s.lastUpdateTime = tutils.Timestamp(sp.LastUpdateTime)
	s.queuedRequests = make(map[string]*request, len(sp.QueuedRequests))
	for rid, req := range sp.QueuedRequests {
		s.queuedRequests[rid] = &request{
			accountID:     req.AccountId,
			confirmedTime: tutils.Timestamp(req.ConfirmedTime),
			enqueueTime:   tutils.Timestamp(req.EnqueueTime),
			labels:        req.Labels,
		}
	}

	s.runningRequestsCache = make(map[string]string, len(sp.Workers))
	s.workers = make(map[string]*worker, len(sp.Workers))
	for wid, w := range sp.Workers {
		var tr *taskRun
		if w.RunningTask != nil {
			cost := balance{}
			copy(cost[:], w.RunningTask.Cost)
			tr = &taskRun{
				cost:     cost,
				priority: int(w.RunningTask.Priority),
				request: &request{
					accountID:     w.RunningTask.Request.AccountId,
					confirmedTime: tutils.Timestamp(w.RunningTask.Request.ConfirmedTime),
					enqueueTime:   tutils.Timestamp(w.RunningTask.Request.EnqueueTime),
					labels:        w.RunningTask.Request.Labels,
				},
				requestID: w.RunningTask.RequestId,
			}
			s.runningRequestsCache[w.RunningTask.RequestId] = wid
		}
		s.workers[wid] = &worker{
			confirmedTime: tutils.Timestamp(w.ConfirmedTime),
			labels:        w.Labels,
			runningTask:   tr,
		}
	}

	s.balances = make(map[string]balance, len(sp.Balances))
	for aid, bal := range sp.Balances {
		newBal := balance{}
		copy(newBal[:], bal.Value)
		s.balances[aid] = newBal
	}

	return s
}

func (s *state) toProto() *StateProto {
	balances := make(map[string]*StateProto_Balance, len(s.balances))
	for aid, bal := range s.balances {
		bCopy := bal
		balances[aid] = &StateProto_Balance{Value: bCopy[:]}
	}

	queuedRequests := make(map[string]*TaskRequest, len(s.queuedRequests))
	for rid, rq := range s.queuedRequests {
		queuedRequests[rid] = newTaskRequest(rq)
	}

	workers := make(map[string]*Worker, len(s.workers))
	for wid, w := range s.workers {
		var rt *TaskRun
		if w.runningTask != nil {
			costCopy := w.runningTask.cost
			rt = &TaskRun{
				Cost:      costCopy[:],
				Priority:  int32(w.runningTask.priority),
				Request:   newTaskRequest(w.runningTask.request),
				RequestId: w.runningTask.requestID,
			}
		}
		workers[wid] = &Worker{
			ConfirmedTime: tutils.TimestampProto(w.confirmedTime),
			Labels:        w.labels,
			RunningTask:   rt,
		}
	}

	return &StateProto{
		Balances:       balances,
		LastUpdateTime: tutils.TimestampProto(s.lastUpdateTime),
		QueuedRequests: queuedRequests,
		Workers:        workers,
	}
}

// Clone returns a deep copy of state, by doing a round-trip proto serialization.
func (s *state) Clone() *state {
	return newStateFromProto(s.toProto())
}
