// copyright 2021 the chromium os authors. all rights reserved.
// use of this source code is governed by a bsd-style license that can be
// found in the license file.

package engine

import (
	"context"
	"testing"
	"time"

	"infra/cros/recovery/internal/execs"
	"infra/cros/recovery/internal/planpb"
	"infra/cros/recovery/logger/metrics"

	"github.com/google/go-cmp/cmp"
)

// TODO(otabek@) Add cases with verification the cache.

// Predefined exec functions.
const (
	exec_pass = "sample_pass"
	exec_fail = "sample_fail"
)

var planTestCases = []struct {
	name       string
	got        *planpb.Plan
	expSuccess bool
}{
	{
		"simple",
		&planpb.Plan{},
		true,
	},
	{
		"critical action fail",
		&planpb.Plan{
			CriticalActions: []string{
				"a1",
				"a2",
			},
			Actions: map[string]*planpb.Action{
				"a1": {
					ExecName: exec_pass,
				},
				"a2": {
					ExecName: exec_fail,
				},
			},
		},
		false,
	},
	{
		"allowed critical action fail",
		&planpb.Plan{
			AllowFail: true,
			CriticalActions: []string{
				"a1",
				"a2",
			},
			Actions: map[string]*planpb.Action{
				"a1": {
					ExecName: exec_pass,
				},
				"a2": {
					ExecName: exec_fail,
				},
			},
		},
		true,
	},
	{
		"skip fail action as not applicable",
		&planpb.Plan{
			CriticalActions: []string{
				"a1",
			},
			Actions: map[string]*planpb.Action{
				"a1": {
					ExecName:   exec_fail,
					Conditions: []string{"c1"},
				},
				"c1": {
					ExecName: exec_fail,
				},
			},
		},
		true,
	},
	{
		"skip fail dependency as not applicable",
		&planpb.Plan{
			CriticalActions: []string{
				"a1",
			},
			Actions: map[string]*planpb.Action{
				"a1": {
					ExecName:     exec_pass,
					Dependencies: []string{"d1"},
				},
				"d1": {
					ExecName:   exec_fail,
					Conditions: []string{"c1"},
				},
				"c1": {
					ExecName: exec_fail,
				},
			},
		},
		true,
	},
	{
		"fail action by dependencies",
		&planpb.Plan{
			CriticalActions: []string{
				"a1",
			},
			Actions: map[string]*planpb.Action{
				"a1": {
					ExecName:     exec_pass,
					Dependencies: []string{"d1"},
				},
				"d1": {
					ExecName: exec_fail,
				},
			},
		},
		false,
	},
	{
		"success run",
		&planpb.Plan{
			CriticalActions: []string{
				"a1",
			},
			Actions: map[string]*planpb.Action{
				"a1": {
					ExecName:     exec_pass,
					Conditions:   []string{"c1"},
					Dependencies: []string{"d1"},
				},
				"c1": {
					ExecName:     exec_pass,
					Dependencies: []string{"d2"},
				},
				"d1": {
					ExecName:     exec_pass,
					Dependencies: []string{"d2"},
				},
				"d2": {
					ExecName: exec_pass,
				},
			},
		},
		true,
	},
}

func TestRun(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	for _, c := range planTestCases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			// t.Parallel() -- TODO(gregorynisbet): Consider parallelizing.
			args := &execs.RunArgs{
				EnableRecovery: true,
			}
			err := Run(ctx, c.name, c.got, args)
			if c.expSuccess {
				if err != nil {
					t.Errorf("Case %q fail but expected to pass. Received error: %s", c.name, err)
				}
			} else {
				if err == nil {
					t.Errorf("Case %q expected to fail but pass", c.name)
				}
			}
		})
	}
}

func TestRunPlanDoNotRunActionAsResultInCache(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	r := recoveryEngine{
		plan: &planpb.Plan{
			CriticalActions: []string{"a"},
			Actions: map[string]*planpb.Action{
				"a": {},
			},
		},
		args: &execs.RunArgs{},
	}
	r.initCache()
	r.cacheActionResult("a", nil)
	err := r.runPlan(ctx)
	if err != nil {
		t.Errorf("Expected plan pass as single action cached with result=nil. Received error: %s", err)
	}
}

var recoveryTestCases = []struct {
	name         string
	got          map[string]*planpb.Action
	expStartOver bool
}{
	{
		"no recoveries",
		map[string]*planpb.Action{
			"a": {
				RecoveryActions: nil,
			},
		},
		false,
	},
	{
		"recoveries stopped on passed r2 and create request to start over",
		map[string]*planpb.Action{
			"a": {
				RecoveryActions: []string{"r1", "r2", "r3"},
			},
			"r1": {
				ExecName: exec_fail,
			},
			"r2": {
				ExecName: exec_pass,
			},
			"r3": {}, //Should not reached
		},
		true,
	},
	{
		"recoveries fail but the process still pass",
		map[string]*planpb.Action{
			"a": {
				RecoveryActions: []string{"r1", "r2", "r3"},
			},
			"r1": {
				ExecName: exec_fail,
			},
			"r2": {
				ExecName: exec_fail,
			},
			"r3": {
				ExecName: exec_fail,
			},
		},
		false,
	},
}

func TestRunRecovery(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	for _, c := range recoveryTestCases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			r := recoveryEngine{
				plan: &planpb.Plan{
					Actions: c.got,
				},
			}
			r.initCache()
			err := r.runRecoveries(ctx, "a")
			if c.expStartOver {
				if !startOverTag.In(err) {
					t.Errorf("Case %q expected to get request to start over. Received: %s", c.name, err)
				}
			} else {
				if err != nil {
					t.Errorf("Case %q expected to receive nil. Received error: %s", c.name, err)
				}
			}
		})
	}
}

var runExecTestCases = []struct {
	name           string
	enableRecovery bool
	got            map[string]*planpb.Action
	expError       bool
	expStartOver   bool
}{
	{
		"do not start recovery flow if action passed",
		true,
		map[string]*planpb.Action{
			"a": {
				ExecName: exec_pass,
				// Will fail if reached any recovery actions.
				RecoveryActions: []string{"r11"},
			},
		},
		false,
		false,
	},
	{
		"do not start recovery flow if it is not allowed",
		false,
		map[string]*planpb.Action{
			"a": {
				ExecName: exec_fail,
				// Will fail if reached any recovery actions.
				RecoveryActions: []string{"r21"},
			},
		},
		true,
		false,
	},
	{
		"receive start over request after run successful recovery action",
		true,
		map[string]*planpb.Action{
			"a": {
				ExecName:        exec_fail,
				RecoveryActions: []string{"r31"},
			},
			"r31": {
				ExecName: exec_pass,
			},
		},
		true,
		true,
	},
	{
		"receive error after try recovery action",
		true,
		map[string]*planpb.Action{
			"a": {
				ExecName:        exec_fail,
				RecoveryActions: []string{"r41"},
			},
			"r41": {
				ExecName: exec_fail,
			},
		},
		true,
		false,
	},
}

func TestActionExec(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	for _, c := range runExecTestCases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			r := recoveryEngine{
				plan: &planpb.Plan{
					Actions: c.got,
				},
			}
			r.initCache()
			err := r.runActionExec(ctx, "a", c.enableRecovery)
			if c.expError && c.expStartOver {
				if !startOverTag.In(err) {
					t.Errorf("Case %q expected to get request to start over. Received error: %s", c.name, err)
				}
			} else if c.expError {
				if err == nil {
					t.Errorf("Case %q expected to fail by passed", c.name)
				}
			} else {
				if err != nil {
					t.Errorf("Case %q expected to receive nil. Received error: %s", c.name, err)
				}
			}
		})
	}
}

var actionResultsCacheTestCases = []struct {
	name       string
	got        map[string]*planpb.Action
	expInCashe bool
	expError   bool
}{
	{
		"set pass to the cache",
		map[string]*planpb.Action{
			"a": {
				ExecName: exec_pass,
			},
		},
		true,
		false,
	},
	{
		"do not set pass to the cache when run_control:run_always",
		map[string]*planpb.Action{
			"a": {
				ExecName:   exec_pass,
				RunControl: planpb.RunControl_ALWAYS_RUN,
			},
		},
		false,
		false,
	},
	{
		"set fail to the cache",
		map[string]*planpb.Action{
			"a": {
				ExecName: exec_fail,
			},
		},
		true,
		true,
	},
	{
		"do not set if recovery finished with success",
		map[string]*planpb.Action{
			"a": {
				ExecName:        exec_fail,
				RecoveryActions: []string{"r"},
			},
			"r": {
				ExecName: exec_pass,
			},
		},
		false,
		false,
	},
	{
		"set fail when all recoveries failed",
		map[string]*planpb.Action{
			"a": {
				ExecName:        exec_fail,
				RecoveryActions: []string{"r"},
			},
			"r": {
				ExecName: exec_fail,
			},
		},
		true,
		true,
	},
	{
		"do not set pass to cache when all recoveries failed and action has run_control:run_always",
		map[string]*planpb.Action{
			"a": {
				ExecName:        exec_fail,
				RecoveryActions: []string{"r"},
				RunControl:      planpb.RunControl_ALWAYS_RUN,
			},
			"r": {
				ExecName: exec_fail,
			},
		},
		false,
		false,
	},
}

func TestActionExecCache(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	for _, c := range actionResultsCacheTestCases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			r := recoveryEngine{
				plan: &planpb.Plan{
					Actions: c.got,
				},
			}
			r.initCache()
			r.runActionExec(ctx, "a", true)
			err, ok := r.actionResultFromCache("a")
			if c.expInCashe {
				if !ok {
					t.Errorf("Case %q: action result in not in the cache", c.name)
				}
				if c.expError && err == nil {
					t.Errorf("Case %q: expected has error as action result but got nil", c.name)
				} else if !c.expError && err != nil {
					t.Errorf("Case %q: expected do not have error as action result but got it: %s", c.name, err)
				}
			} else {
				if ok {
					t.Errorf("Case %q: does not expected result in the cache", c.name)
				}
			}
		})
	}
}

var resetCacheTestCases = []struct {
	name    string
	got     map[string]planpb.RunControl
	present []string
	removed []string
}{
	{
		"clean all",
		map[string]planpb.RunControl{
			"a1": planpb.RunControl_RERUN_AFTER_RECOVERY,
			"a2": planpb.RunControl_RERUN_AFTER_RECOVERY,
			"a3": planpb.RunControl_RERUN_AFTER_RECOVERY,
			"a4": planpb.RunControl_RERUN_AFTER_RECOVERY,
		},
		nil,
		[]string{"a1", "a2", "a3", "a4"},
	},
	{
		"partially clean up",
		map[string]planpb.RunControl{
			"a1": planpb.RunControl_RUN_ONCE,
			"a2": planpb.RunControl_RUN_ONCE,
			"a3": planpb.RunControl_RERUN_AFTER_RECOVERY,
			"a4": planpb.RunControl_RERUN_AFTER_RECOVERY,
		},
		[]string{"a1", "a2"},
		[]string{"a3", "a4"},
	},
}

func TestResetCacheAfterSuccessfulRecoveryAction(t *testing.T) {
	t.Parallel()
	for _, c := range resetCacheTestCases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			actions := make(map[string]*planpb.Action)
			for name, rc := range c.got {
				actions[name] = &planpb.Action{
					RunControl: rc,
				}
			}
			r := recoveryEngine{
				plan: &planpb.Plan{
					Actions: actions,
				},
			}
			r.initCache()
			for name := range c.got {
				r.cacheActionResult(name, nil)
			}
			r.resetCacheAfterSuccessfulRecoveryAction()
			for _, name := range c.present {
				if _, ok := r.actionResultFromCache(name); !ok {
					t.Errorf("Case %q: expected to have result for action %q in the cache", c.name, name)
				}
			}
			for _, name := range c.removed {
				if _, ok := r.actionResultFromCache(name); ok {
					t.Errorf("Case %q: not expected to have result for action %q in the cache", c.name, name)
				}
			}
		})
	}
}

var setCacheTestCases = []struct {
	name string
	got  planpb.RunControl
	exp  bool
}{
	{
		"run once",
		planpb.RunControl_RUN_ONCE,
		true,
	},
	{
		"rerun after recovery",
		planpb.RunControl_RERUN_AFTER_RECOVERY,
		true,
	},
	{
		"always run",
		planpb.RunControl_ALWAYS_RUN,
		false,
	},
}

func TestCacheActionResult(t *testing.T) {
	t.Parallel()
	for _, c := range setCacheTestCases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			r := recoveryEngine{
				plan: &planpb.Plan{
					Actions: map[string]*planpb.Action{
						"a": {
							RunControl: c.got,
						},
					},
				},
			}
			r.initCache()
			r.cacheActionResult("a", nil)
			_, ok := r.actionResultFromCache("a")
			if c.exp {
				if !ok {
					t.Errorf("Case %q: expected to have result but not present in cache", c.name)
				}
			} else {
				if ok {
					t.Errorf("Case %q: not expected to have result but present in cache", c.name)
				}
			}
		})
	}
}

var isRecoveryUsageTestCases = []struct {
	name          string
	actionCache   []string
	recoveryCache []recoveryUsageKey
	used          bool
}{
	{
		"not used",
		[]string{"a", "b"},
		[]recoveryUsageKey{
			{
				action:   "a",
				recovery: "a",
			},
			{
				action:   "a",
				recovery: "b",
			},
			{
				action:   "b",
				recovery: "a",
			},
			{
				action:   "b",
				recovery: "r",
			},
		},
		false,
	},
	{
		"used by action result",
		[]string{"r"},
		nil,
		true,
	},
	{
		"used by recovery result from other action",
		nil,
		[]recoveryUsageKey{
			{
				action:   "a",
				recovery: "r",
			},
		},
		true,
	},
}

func TestRecoveryCachePersistence(t *testing.T) {
	t.Parallel()
	for _, c := range isRecoveryUsageTestCases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			r := recoveryEngine{
				plan: &planpb.Plan{
					Actions: map[string]*planpb.Action{
						"a": {},
						"b": {},
						"r": {},
					},
				},
			}
			r.initCache()
			for _, name := range c.actionCache {
				r.cacheActionResult(name, nil)
			}
			for _, k := range c.recoveryCache {
				r.registerRecoveryUsage(k.action, k.recovery, nil)
			}
			if r.isRecoveryUsed("a", "r") != c.used {
				t.Errorf("Case %q before rest: expectaton did not matche expectations: Expected: %v, Got: %v", c.name, c.used, !c.used)
			}
			r.resetCacheAfterSuccessfulRecoveryAction()
			if r.isRecoveryUsed("a", "r") != c.used {
				t.Errorf("Case %q after reset: expectaton did not matche expectations: Expected: %v, Got: %v", c.name, c.used, !c.used)
			}
		})
	}
}

// TestCallMetricsInSimplePlan tests that calling a simple plan with a fake implementation of a metrics interface calls the metrics implementation.
func TestCallMetricsInSimplePlan(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	m := newFakeMetrics()
	r := &recoveryEngine{
		planName: "2e9aa66a-5fa1-4eaa-933c-eee8e4337823",
	}
	// NOTE: This is a bit subtle, but there really should be TWO records here. The fake implementation service always appends new records,
	// regardless of whether Karte would create a new record or update one in place. This is good for unit tests because it means that every
	// intermediate state is visible, so we really are testing the entire interaction.
	expected := []*metrics.Action{
		{
			ActionKind: "plan:2e9aa66a-5fa1-4eaa-933c-eee8e4337823",
			Observations: []*metrics.Observation{
				{MetricKind: "restarts", ValueType: "number", Value: "0"},
				{MetricKind: "forgiven_failures", ValueType: "number", Value: "0"},
			},
		},
		{
			ActionKind: "plan:2e9aa66a-5fa1-4eaa-933c-eee8e4337823",
			Observations: []*metrics.Observation{
				{MetricKind: "restarts", ValueType: "number", Value: "0"},
				{MetricKind: "forgiven_failures", ValueType: "number", Value: "0"},
			},
		},
	}
	r.plan = &planpb.Plan{
		Actions: map[string]*planpb.Action{
			"a": {},
			"b": {},
			"r": {},
		},
	}
	r.args = &execs.RunArgs{
		Metrics: m,
	}
	err := r.runPlan(ctx)
	// TODO(gregorynisbet): Mock the time.Now() function everywhere instead of removing times
	// from test cases.
	for i := range m.actions {
		var zero time.Time
		m.actions[i].StartTime = zero
	}
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if diff := cmp.Diff(expected, m.actions); diff != "" {
		t.Errorf("unexpected diff (-want +got): %s", diff)
	}
}
