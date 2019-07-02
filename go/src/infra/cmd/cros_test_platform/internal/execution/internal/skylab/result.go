// Copyright 2019 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package skylab

import (
	"go.chromium.org/chromiumos/infra/proto/go/test_platform"
	"go.chromium.org/chromiumos/infra/proto/go/test_platform/skylab_test_runner"
	"go.chromium.org/chromiumos/infra/proto/go/test_platform/steps"

	"infra/cmd/cros_test_platform/internal/execution/swarming"
)

func toTaskResults(testRuns []*testRun, urler swarming.URLer) []*steps.ExecuteResponse_TaskResult {
	var results []*steps.ExecuteResponse_TaskResult
	for _, test := range testRuns {
		for _, attempt := range test.attempts {
			results = append(results, toTaskResult(test.test.Name, &attempt, urler))
		}
	}
	return results
}

func toTaskResult(testName string, attempt *attempt, urler swarming.URLer) *steps.ExecuteResponse_TaskResult {
	var verdict test_platform.TaskState_Verdict

	switch {
	case attempt.autotestResult == nil:
		verdict = test_platform.TaskState_VERDICT_UNSPECIFIED
	case attempt.autotestResult.Incomplete:
		verdict = test_platform.TaskState_VERDICT_FAILED
	default:
		verdict = flattenToVerdict(attempt.autotestResult.TestCases)
	}

	return &steps.ExecuteResponse_TaskResult{
		Name: testName,
		State: &test_platform.TaskState{
			LifeCycle: taskStateToLifeCycle[attempt.state],
			Verdict:   verdict,
		},
		TaskId:  attempt.taskID,
		TaskUrl: urler.GetTaskURL(attempt.taskID),
	}
}

func flattenToVerdict(tests []*skylab_test_runner.Result_Autotest_TestCase) test_platform.TaskState_Verdict {
	// By default (if no test cases ran), then there is no verdict.
	verdict := test_platform.TaskState_VERDICT_NO_VERDICT
	for _, c := range tests {
		switch c.Verdict {
		case skylab_test_runner.Result_Autotest_TestCase_VERDICT_FAIL:
			// Any case failing means the flat verdict is a failure.
			return test_platform.TaskState_VERDICT_FAILED
		case skylab_test_runner.Result_Autotest_TestCase_VERDICT_PASS:
			// Otherwise, at least 1 passing verdict means a pass.
			verdict = test_platform.TaskState_VERDICT_PASSED
		case skylab_test_runner.Result_Autotest_TestCase_VERDICT_UNDEFINED:
			// Undefined verdicts do not affect flat verdict.
		}
	}
	return verdict
}
