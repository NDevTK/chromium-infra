// Copyright 2019 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package cmd

import (
	"context"
	"fmt"
	"infra/cmd/skylab/internal/site"
	"io"
	"time"

	"github.com/maruel/subcommands"

	"go.chromium.org/chromiumos/infra/proto/go/test_platform"
	"go.chromium.org/chromiumos/infra/proto/go/test_platform/skylab_tool"
	"go.chromium.org/chromiumos/infra/proto/go/test_platform/steps"
	"go.chromium.org/luci/auth/client/authcli"
	swarming_api "go.chromium.org/luci/common/api/swarming/swarming/v1"
	"go.chromium.org/luci/common/cli"

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

func (c *waitTaskRun) Run(a subcommands.Application, args []string, env subcommands.Env) int {
	if err := c.innerRun(a, env); err != nil {
		PrintError(a.GetErr(), err)
		return 1
	}
	return 0
}

func (c *waitTaskRun) innerRun(a subcommands.Application, env subcommands.Env) error {
	var result *skylab_tool.WaitTaskResult
	var err error
	switch c.buildBucket {
	case true:
		result, err = c.innerRunBuildbucket(a, env)
	case false:
		result, err = c.innerRunSwarming(a, env)
	}

	if err != nil {
		return err
	}

	printJSONResults(a.GetOut(), result)
	return nil
}

func (c *waitTaskRun) innerRunSwarming(a subcommands.Application, env subcommands.Env) (*skylab_tool.WaitTaskResult, error) {
	taskID := c.Flags.Arg(0)
	if taskID == "" {
		return nil, NewUsageError(c.Flags, "missing swarming task ID")
	}

	ctx := cli.GetContext(a, c, env)
	ctx, cancel := maybeWithTimeout(ctx, c.timeoutMins)
	defer cancel(context.Canceled)

	client, err := newSwarmingClient(ctx, c.authFlags, c.envFlags.Env())
	if err != nil {
		return nil, err
	}

	if err = waitSwarmingTask(ctx, taskID, client); err != nil {
		return nil, err
	}

	return extractSwarmingResult(ctx, taskID, client)
}

func (c *waitTaskRun) innerRunBuildbucket(a subcommands.Application, env subcommands.Env) (*skylab_tool.WaitTaskResult, error) {
	taskIDString := c.Flags.Arg(0)
	if taskIDString == "" {
		return nil, NewUsageError(c.Flags, "missing buildbucket task id")
	}

	ctx := cli.GetContext(a, c, env)
	ctx, cancel := maybeWithTimeout(ctx, c.timeoutMins)
	defer cancel(context.Canceled)

	bClient, err := bbClient(ctx, c.envFlags.Env(), c.authFlags)
	if err != nil {
		return nil, err
	}

	return waitBuildbucketTask(ctx, taskIDString, bClient, c.envFlags.Env())
}

func responseToTaskResult(e site.Environment, buildID int64, response *steps.ExecuteResponse) *skylab_tool.WaitTaskResult {
	u := bbURL(e, buildID)
	verdict := response.GetState().GetVerdict()
	failure := verdict == test_platform.TaskState_VERDICT_FAILED
	success := verdict == test_platform.TaskState_VERDICT_PASSED
	tr := &skylab_tool.WaitTaskResult_Task{
		Name:          "Test Platform Invocation",
		TaskRunUrl:    u,
		TaskRunId:     fmt.Sprintf("%d", buildID),
		TaskRequestId: fmt.Sprintf("%d", buildID),
		Failure:       failure,
		Success:       success,
	}
	var childResults []*skylab_tool.WaitTaskResult_Task
	for _, child := range response.TaskResults {
		verdict := child.GetState().GetVerdict()
		failure := verdict == test_platform.TaskState_VERDICT_FAILED
		success := verdict == test_platform.TaskState_VERDICT_PASSED
		childResult := &skylab_tool.WaitTaskResult_Task{
			Name:        child.Name,
			TaskLogsUrl: child.LogUrl,
			TaskRunUrl:  child.TaskUrl,
			// Note: TaskRunID is deprecated and excluded here.
			Failure: failure,
			Success: success,
		}
		childResults = append(childResults, childResult)
	}
	return &skylab_tool.WaitTaskResult{
		ChildResults: childResults,
		Result:       tr,
		// Note: Stdout it not set.
	}
}

// waitSwarmingTask waits until the task with the given ID has completed.
//
// It returns an error if the given context was cancelled or in case of swarming
// rpc failures (after transient retry).
func waitSwarmingTask(ctx context.Context, taskID string, t *swarming.Client) error {
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

func asTaskResult(s *swarming_api.SwarmingRpcsTaskResult) *skylab_tool.WaitTaskResult_Task {
	return &skylab_tool.WaitTaskResult_Task{
		Name:  s.Name,
		State: s.State,
		// TODO(crbug.com/964573): Deprecate this field.
		Failure:       s.Failure,
		Success:       !s.Failure && (s.State == "COMPLETED" || s.State == "COMPLETED_SUCCESS"),
		TaskRunId:     s.RunId,
		TaskRequestId: s.TaskId,
	}
}

func extractSwarmingResult(ctx context.Context, taskID string, t *swarming.Client) (*skylab_tool.WaitTaskResult, error) {
	results, err := t.GetResults(ctx, []string{taskID})
	if err != nil {
		return nil, err
	}
	stdouts, err := t.GetTaskOutputs(ctx, []string{taskID})
	if err != nil {
		return nil, err
	}
	childs, err := t.GetResults(ctx, results[0].ChildrenTaskIds)
	if err != nil {
		return nil, err
	}
	childResults := make([]*skylab_tool.WaitTaskResult_Task, len(childs))
	for i, c := range childs {
		childResults[i] = asTaskResult(c)
	}
	tr := asTaskResult(results[0])
	result := &skylab_tool.WaitTaskResult{
		Result:       tr,
		Stdout:       stdouts[0].Output,
		ChildResults: childResults,
	}
	return result, nil
}

func printJSONResults(w io.Writer, m *skylab_tool.WaitTaskResult) {
	err := jsonPBMarshaller.Marshal(w, m)
	if err != nil {
		panic(err)
	}
}
