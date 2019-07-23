// Copyright 2019 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"infra/cmd/skylab/internal/site"
	"io"
	"strconv"
	"time"

	"github.com/maruel/subcommands"
	"go.chromium.org/chromiumos/infra/proto/go/test_platform"
	"go.chromium.org/chromiumos/infra/proto/go/test_platform/steps"
	"go.chromium.org/luci/auth/client/authcli"
	swarming_api "go.chromium.org/luci/common/api/swarming/swarming/v1"
	"go.chromium.org/luci/common/cli"
	"go.chromium.org/luci/common/errors"

	"infra/libs/skylab/swarming"
)

// WaitTask subcommand: wait for a task to finish.
var WaitTask = &subcommands.Command{
	UsageLine: "wait-task [FLAGS...] TASK_ID",
	ShortDesc: "wait for a task to complete",
	LongDesc:  `Wait for the task with the given swarming task id to complete, and summarize its results.`,
	CommandRun: func() subcommands.CommandRun {
		c := &waitTaskRun{}
		c.authFlags.Register(&c.Flags, site.DefaultAuthOptions)
		c.envFlags.Register(&c.Flags)

		c.Flags.IntVar(&c.timeoutMins, "timeout-mins", -1, "The maxinum number of minutes to wait for the task to finish. Default: no timeout.")
		c.Flags.BoolVar(&c.buildBucket, "bb", false, "(Expert use only, not a stable API) use buildbucket recipe backend. If specified, TASK_ID is interpreted as a buildbucket task id.")
		return c
	},
}

type waitTaskRun struct {
	subcommands.CommandRunBase
	authFlags   authcli.Flags
	envFlags    envFlags
	timeoutMins int
	buildBucket bool
}

type taskResult struct {
	Name  string `json:"name"`
	State string `json:"state"`
	// TODO(crbug.com/964573): Deprecate this field.
	Failure bool `json:"failure"`
	Success bool `json:"success"`
	// TODO(akeshet): Deprecate this field; the run ID is not independently
	// meaningful to callers, it depends on the namespace of the ID (e.g.
	// swarming? buildbucket? dev? prod?)
	TaskRunID string `json:"task-run-id"`
	// Note: these URL fields are only populated for -bb runs; eventually,
	// non-bb runs will be deprecated.
	TaskRunURL  string `json:"task-run-url"`
	TaskLogsURL string `json:"task-logs-url"`
}

type waitTaskResult struct {
	TaskResult   *taskResult   `json:"task-result"`
	Stdout       string        `json:"stdout"`
	ChildResults []*taskResult `json:"child-results"`
}

func (c *waitTaskRun) Run(a subcommands.Application, args []string, env subcommands.Env) int {
	if err := c.innerRun(a, args, env); err != nil {
		PrintError(a.GetErr(), err)
		return 1
	}
	return 0
}

func (c *waitTaskRun) innerRun(a subcommands.Application, args []string, env subcommands.Env) error {
	if c.buildBucket {
		return c.innerRunBuildbucket(a, env, a.GetOut())
	}

	taskID := c.Flags.Arg(0)
	if taskID == "" {
		return NewUsageError(c.Flags, "missing swarming task ID")
	}

	siteEnv := c.envFlags.Env()
	ctx := cli.GetContext(a, c, env)
	h, err := httpClient(ctx, &c.authFlags)
	if err != nil {
		return errors.Annotate(err, "failed to create http client").Err()
	}
	client, err := swarming.New(ctx, h, siteEnv.SwarmingService)
	if err != nil {
		return err
	}

	var taskWaitCtx context.Context
	var taskWaitCancel context.CancelFunc
	if c.timeoutMins >= 0 {
		taskWaitCtx, taskWaitCancel = context.WithTimeout(ctx, time.Duration(c.timeoutMins)*time.Minute)
	} else {
		taskWaitCtx, taskWaitCancel = context.WithCancel(ctx)
	}
	defer taskWaitCancel()

	if err = waitTask(taskWaitCtx, taskID, client); err != nil {
		if err == context.DeadlineExceeded {
			return errors.New("timed out waiting for task to complete")
		}
		return err
	}

	return postWaitTask(ctx, taskID, client, a.GetOut())
}

func (c *waitTaskRun) innerRunBuildbucket(a subcommands.Application, env subcommands.Env, w io.Writer) error {
	taskIDString := c.Flags.Arg(0)
	if taskIDString == "" {
		return NewUsageError(c.Flags, "missing buildbucket task id")
	}

	buildID, err := strconv.ParseInt(taskIDString, 10, 64)
	if err != nil {
		return err
	}

	ctx := cli.GetContext(a, c, env)
	// TODO(akeshet): Respect wait timeout.
	build, err := bbWaitBuild(ctx, c.envFlags.Env(), c.authFlags, buildID)
	if err != nil {
		return err
	}

	response, err := bbExtractResponse(build)
	if err != nil {
		return err
	}

	waitResult := responseToTaskResult(c.envFlags.Env(), buildID, response)

	printJSONResults(w, waitResult)
	return nil
}

func responseToTaskResult(e site.Environment, buildID int64, response *steps.ExecuteResponse) *waitTaskResult {
	u := bbURL(e, buildID)
	verdict := response.GetState().GetVerdict()
	failure := verdict == test_platform.TaskState_VERDICT_FAILED
	success := verdict == test_platform.TaskState_VERDICT_PASSED
	tr := &taskResult{
		Name:       "Test Platform Invocation",
		TaskRunURL: u,
		// TODO(akeshet): Deprecate this field.
		TaskRunID: fmt.Sprintf("%d", buildID),
		Failure:   failure,
		Success:   success,
	}
	var childResults []*taskResult
	for _, child := range response.TaskResults {
		verdict := child.GetState().GetVerdict()
		failure := verdict == test_platform.TaskState_VERDICT_FAILED
		success := verdict == test_platform.TaskState_VERDICT_PASSED
		childResult := &taskResult{
			Name:        child.Name,
			TaskLogsURL: child.LogUrl,
			TaskRunURL:  child.TaskUrl,
			// Note: TaskRunID is deprecated and excluded here.
			Failure: failure,
			Success: success,
		}
		childResults = append(childResults, childResult)
	}
	return &waitTaskResult{
		ChildResults: childResults,
		TaskResult:   tr,
		// Note: Stdout it not set.
	}
}

// waitTask waits until the task with the given ID has completed.
//
// It returns an error if the given context was cancelled or in case of swarming
// rpc failures (after transient retry).
func waitTask(ctx context.Context, taskID string, t *swarming.Client) error {
	sleepInterval := time.Duration(15 * time.Second)
	for {
		results, err := t.GetResults(ctx, []string{taskID})
		if err != nil {
			return err
		}
		// Possible values:
		//   "BOT_DIED"
		//   "CANCELED"
		//   "COMPLETED"
		//   "EXPIRED"
		//   "INVALID"
		//   "KILLED"
		//   "NO_RESOURCE"
		//   "PENDING"
		//   "RUNNING"
		//   "TIMED_OUT"
		// Keep waiting if task state is RUNNING or PENDING
		if s := results[0].State; s != "RUNNING" && s != "PENDING" {
			return nil
		}
		if err = sleepOrCancel(ctx, sleepInterval); err != nil {
			return err
		}
	}
}

func sleepOrCancel(ctx context.Context, duration time.Duration) error {
	sleepTimer := time.NewTimer(duration)
	select {
	case <-sleepTimer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func asTaskResult(s *swarming_api.SwarmingRpcsTaskResult) *taskResult {
	return &taskResult{
		Name:  s.Name,
		State: s.State,
		// TODO(crbug.com/964573): Deprecate this field.
		Failure:   s.Failure,
		Success:   !s.Failure && (s.State == "COMPLETED" || s.State == "COMPLETED_SUCCESS"),
		TaskRunID: s.RunId,
	}
}

func postWaitTask(ctx context.Context, taskID string, t *swarming.Client, w io.Writer) error {
	results, err := t.GetResults(ctx, []string{taskID})
	if err != nil {
		return err
	}
	stdouts, err := t.GetTaskOutputs(ctx, []string{taskID})
	if err != nil {
		return err
	}
	childs, err := t.GetResults(ctx, results[0].ChildrenTaskIds)
	if err != nil {
		return err
	}
	childResults := make([]*taskResult, len(childs))
	for i, c := range childs {
		childResults[i] = asTaskResult(c)
	}
	tr := asTaskResult(results[0])
	result := &waitTaskResult{
		TaskResult:   tr,
		Stdout:       stdouts[0].Output,
		ChildResults: childResults,
	}
	printJSONResults(w, result)
	return nil
}

func printJSONResults(w io.Writer, m *waitTaskResult) {
	outputJSON, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, string(outputJSON))
}
