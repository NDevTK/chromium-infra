// Copyright 2021 The Chromium OS Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Package engine provides struts and functionality of recovery engine.
// For more details please read go/paris-recovery-engine.
package engine

import (
	"context"
	"fmt"
	"time"

	"go.chromium.org/luci/common/errors"
	"go.chromium.org/luci/luciexe/build"

	"infra/cros/recovery/internal/execs"
	"infra/cros/recovery/internal/log"
	"infra/cros/recovery/internal/planpb"
	"infra/cros/recovery/logger/metrics"
)

// recoveryEngine holds info required for running a recovery plan.
type recoveryEngine struct {
	planName string
	plan     *planpb.Plan
	args     *execs.RunArgs
	// Caches
	actionResultsCache map[string]error
	recoveryUsageCache map[recoveryUsageKey]error
}

// Error tag to track error with request to start critical actions over.
var startOverTag = errors.BoolTag{Key: errors.NewTagKey("start-over")}

// Run runs the recovery plan.
func Run(ctx context.Context, planName string, plan *planpb.Plan, args *execs.RunArgs) error {
	r := &recoveryEngine{
		planName: planName,
		plan:     plan,
		args:     args,
	}
	r.initCache()
	defer func() { r.close() }()
	log.Debug(ctx, "Received plan %s for %s \n%s", r.planName, r.args.ResourceName, r.describe())
	return r.runPlan(ctx)
}

// close free up used resources.
func (r *recoveryEngine) close() {
	if r.actionResultsCache != nil {
		r.actionResultsCache = nil
	}
	// TODO(otabek@): Close the caches.
}

// runPlan executes recovery plan with critical-actions.
func (r *recoveryEngine) runPlan(ctx context.Context) (rErr error) {
	newCtx := ctx
	log.Info(newCtx, "Plan %q: started", r.planName)

	var metricAction *metrics.Action
	var restartTally int64
	var forgivenFailureTally int64

	if r.args != nil && r.args.Metrics != nil {
		var closer execs.CloserFunc
		metricAction, closer = r.args.NewMetric(
			newCtx,
			fmt.Sprintf("plan:%s", r.planName),
		)
		defer func() {
			if closer != nil {
				metricAction.Observations = append(
					metricAction.Observations,
					metrics.NewInt64Observation("restarts", restartTally),
					metrics.NewInt64Observation("forgiven_failures", forgivenFailureTally),
				)
				closer(ctx, rErr)
			}
		}()
	}

	for {
		if err := r.runActions(ctx, r.plan.GetCriticalActions(), r.args.EnableRecovery); err != nil {
			if startOverTag.In(err) {
				log.Info(ctx, "Plan %q for %s: received request to start over.", r.planName, r.args.ResourceName)
				r.resetCacheAfterSuccessfulRecoveryAction()
				restartTally++
				continue
			}
			if r.plan.GetAllowFail() {
				log.Info(ctx, "Plan %q for %s: failed with error: %s.", r.planName, r.args.ResourceName, err)
				log.Info(ctx, "Plan %q for %s: is allowed to fail, continue.", r.planName, r.args.ResourceName)
				forgivenFailureTally++
				return nil
			}
			return errors.Annotate(err, "run plan %q", r.planName).Err()
		}
		break
	}
	log.Info(ctx, "Plan %q: finished successfully.", r.planName)
	log.Info(ctx, "Plan %q: recorded %d restarts during execution.", r.planName, restartTally)
	log.Info(ctx, "Plan %q: recorded %d forgiven failures during execution.", r.planName, forgivenFailureTally)
	return nil
}

// runActions runs actions in order.
func (r *recoveryEngine) runActions(ctx context.Context, actions []string, enableRecovery bool) error {
	for _, actionName := range actions {
		if err := r.runAction(ctx, actionName, enableRecovery); err != nil {
			return errors.Annotate(err, "run actions").Err()
		}
	}
	return nil
}

// runAction runs single action.
// Execution steps:
// 1) Check action's result in cache.
// 2) Check if the action is applicable based on conditions. Skip if any fail.
// 3) Run dependencies of the action. Fail if any fails.
// 4) Run action exec function. Fail if any fail.
func (r *recoveryEngine) runAction(ctx context.Context, actionName string, enableRecovery bool) (rErr error) {
	newCtx := ctx
	if r.args != nil {
		if r.args.Metrics != nil {
			log.Debug(ctx, "Recording metrics for action %q", actionName)
			_, closer := r.args.NewMetric(
				newCtx,
				// TODO(gregorynisbet): Consider adding a new field to Karte to explicitly track the name
				//                      assigned to an action by recoverylib.
				fmt.Sprintf("action:%s", actionName),
			)
			defer func() {
				closer(ctx, rErr)
			}()
		} else {
			log.Debug(ctx, "Skipping metrics for action %q", actionName)
		}
		if r.args.ShowSteps {
			var step *build.Step
			step, ctx = build.StartStep(ctx, fmt.Sprintf("Run %s", actionName))
			defer func() { step.End(rErr) }()
		}
		if r.args.Logger != nil {
			r.args.Logger.IndentLogging()
			defer func() { r.args.Logger.DedentLogging() }()
		}
	}
	log.Info(ctx, "Action %q: started.", actionName)
	if aErr, ok := r.actionResultFromCache(actionName); ok {
		if aErr == nil {
			log.Info(ctx, "Action %q: pass (cached).", actionName)
			// Return nil error so we can continue execution of next actions...
			return nil
		}
		a := r.getAction(actionName)
		if a.GetAllowFailAfterRecovery() {
			log.Info(ctx, "Action %q: fail (cached). Error: %s", actionName, aErr)
			log.Debug(ctx, "Action %q: error ignored as action is allowed to fail.", actionName)
			// Return nil error so we can continue execution of next actions...
			return nil
		}
		return errors.Annotate(aErr, "run action %q: (cached)", actionName).Err()
	}
	conditionName, err := r.runActionConditions(ctx, actionName)
	if err != nil {
		log.Info(ctx, "Action %q: one of conditions %q failed, skipping...", actionName, conditionName)
		log.Debug(ctx, "Action %q: conditions fail with %s", actionName, err)
		// Return nil error so we can continue execution of next actions...
		return nil
	}
	if err := r.runDependencies(ctx, actionName, enableRecovery); err != nil {
		return errors.Annotate(err, "run action %q", actionName).Err()
	}
	if err := r.runActionExec(ctx, actionName, enableRecovery); err != nil {
		if startOverTag.In(err) {
			return errors.Annotate(err, "run action %q", actionName).Err()
		}
		a := r.getAction(actionName)
		if a.GetAllowFailAfterRecovery() {
			log.Info(ctx, "Action %q: fail. Error: %s", actionName, err)
			log.Debug(ctx, "Action %q: error ignored as action is allowed to fail.", actionName)
		} else {
			return errors.Annotate(err, "run action %q", actionName).Err()
		}
	} else {
		log.Info(ctx, "Action %q: finished successfully.", actionName)
	}
	// Return nil error so we can continue execution of next actions...
	return nil
}

// runActionExec runs action's exec function and initiates recovery flow if exec fails.
// The recover flow start only recoveries is enabled.
func (r *recoveryEngine) runActionExec(ctx context.Context, actionName string, enableRecovery bool) error {
	a := r.getAction(actionName)
	if err := r.runActionExecWithTimeout(ctx, a); err != nil {
		if enableRecovery && len(a.GetRecoveryActions()) > 0 {
			log.Info(ctx, "Action %q: starting recovery actions.", actionName)
			log.Debug(ctx, "Action %q: fail. Error: %s", actionName, err)
			if rErr := r.runRecoveries(ctx, actionName); rErr != nil {
				return errors.Annotate(rErr, "run action %q exec", actionName).Err()
			}
			log.Info(ctx, "Run action %q exec: no recoveries left to try", actionName)
		}
		// Cache the action error only after running recoveries.
		// If no recoveries were run, we still cache the action.
		r.cacheActionResult(actionName, err)
		return errors.Annotate(err, "run action %q exec", actionName).Err()
	}
	r.cacheActionResult(actionName, nil)
	return nil
}

// Default time limit per action exec function.
const defaultExecTimeout = 60 * time.Second

func actionExecTimeout(a *planpb.Action) time.Duration {
	if a.ExecTimeout != nil {
		return a.ExecTimeout.AsDuration()
	}
	return defaultExecTimeout
}

// runActionExecWithTimeout runs action's exec function with timeout.
func (r *recoveryEngine) runActionExecWithTimeout(ctx context.Context, a *planpb.Action) error {
	timeout := actionExecTimeout(a)
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer func() { cancel() }()
	cw := make(chan error, 1)
	go func() {
		err := execs.Run(ctx, a.ExecName, r.args, a.ExecExtraArgs)
		cw <- err
	}()
	select {
	case err := <-cw:
		return errors.Annotate(err, "run exec %q with timeout %s", a.ExecName, timeout).Err()
	case <-ctx.Done():
		log.Info(ctx, "Run exec %q with timeout %s: excited timeout", a.ExecName, timeout)
		return errors.Reason("run exec %q with timeout %s: excited timeout", a.ExecName, timeout).Err()
	}
}

// runActionConditions checks if action is applicable based on condition actions.
// If return err then not applicable, if nil then applicable.
func (r *recoveryEngine) runActionConditions(ctx context.Context, actionName string) (conditionName string, err error) {
	a := r.getAction(actionName)
	if len(a.GetConditions()) == 0 {
		return "", nil
	}
	if r.args != nil {
		if r.args.ShowSteps {
			var step *build.Step
			step, ctx = build.StartStep(ctx, "Run continions")
			defer func() { step.End(err) }()
		}
		if r.args.Logger != nil {
			r.args.Logger.IndentLogging()
			defer func() { r.args.Logger.DedentLogging() }()
		}
	}
	log.Debug(ctx, "Action %q: running conditions...", actionName)
	enableRecovery := false
	for _, condition := range a.GetConditions() {
		if err := r.runAction(ctx, condition, enableRecovery); err != nil {
			log.Debug(ctx, "Action %q: condition %q fails. Error: %s", actionName, condition, err)
			return condition, errors.Annotate(err, "run conditions").Err()
		}
	}
	log.Debug(ctx, "Action %q: all conditions passed.", actionName)
	return "", nil
}

// runDependencies runs action's dependencies.
func (r *recoveryEngine) runDependencies(ctx context.Context, actionName string, enableRecovery bool) (rErr error) {
	a := r.getAction(actionName)
	if len(a.GetDependencies()) == 0 {
		return nil
	}
	if r.args != nil {
		if r.args.ShowSteps {
			var step *build.Step
			step, ctx = build.StartStep(ctx, "Run dependencies")
			defer func() { step.End(rErr) }()
		}
		if r.args.Logger != nil {
			r.args.Logger.IndentLogging()
			defer func() { r.args.Logger.DedentLogging() }()
		}
	}
	err := r.runActions(ctx, a.GetDependencies(), enableRecovery)
	return errors.Annotate(err, "run dependencies").Err()
}

// runRecoveries runs action's recoveries.
// Recovery actions are expected to fail. If recovery action fails then next will be attempted.
// Finishes with nil if no recovery action provided or nether succeeded.
// Finishes with start-over request if any recovery succeeded.
// Recovery action will skip if used before.
func (r *recoveryEngine) runRecoveries(ctx context.Context, actionName string) (rErr error) {
	a := r.getAction(actionName)
	if len(a.GetRecoveryActions()) == 0 {
		return nil
	}
	if r.args != nil {
		if r.args.ShowSteps {
			var step *build.Step
			step, ctx = build.StartStep(ctx, "Run recoveries")
			defer func() { step.End(rErr) }()
		}
		if r.args.Logger != nil {
			r.args.Logger.IndentLogging()
			defer func() { r.args.Logger.DedentLogging() }()
		}
	}
	for _, recoveryName := range a.GetRecoveryActions() {
		if r.isRecoveryUsed(actionName, recoveryName) {
			// Engine allows to use each recovery action only once in scope of the action.
			continue
		}
		if err := r.runActions(ctx, []string{recoveryName}, false); err != nil {
			log.Debug(ctx, "Recovery action %q: fail. Error: %s ", recoveryName, err)
			r.registerRecoveryUsage(actionName, recoveryName, err)
			continue
		}
		r.registerRecoveryUsage(actionName, recoveryName, nil)
		log.Info(ctx, "Recovery action %q: request to start-over.", recoveryName)
		return errors.Reason("recovery action %q: request to start over", recoveryName).Tag(startOverTag).Err()
	}
	return nil
}

// getAction finds and provides action from the plan collection.
func (r *recoveryEngine) getAction(name string) *planpb.Action {
	if a, ok := r.plan.Actions[name]; ok {
		return a
	}
	// If we reach this place then we have issues with plan validation logic.
	panic(fmt.Sprintf("action %q not found in the plan", name))
}

// describe describes the plan details with critical actions.
func (r *recoveryEngine) describe() string {
	d := fmt.Sprintf("Plan %q, AllowFail: %v ", r.planName, r.plan.AllowFail)
	if len(r.plan.GetCriticalActions()) > 0 {
		prefix := "\n "
		d += fmt.Sprintf("%sCritical-actions:", prefix)
		for i, a := range r.plan.GetCriticalActions() {
			d += fmt.Sprintf("%s %d: %s", prefix, i, r.describeAction(a, prefix+"  "))
		}
	} else {
		d += "\n No critical-actions"
	}
	return d
}

// describeAction describes the action structure recursively.
func (r *recoveryEngine) describeAction(actionName string, prefix string) string {
	a := r.getAction(actionName)
	ap := fmt.Sprintf("Action %q, AllowFailAfterRecovery: %v, RunControl: %v",
		actionName, a.GetAllowFailAfterRecovery(), a.GetRunControl())
	if len(a.GetConditions()) > 0 {
		ap += fmt.Sprintf("%sConditions:", prefix)
		for i, d := range a.GetConditions() {
			ap += fmt.Sprintf("%s%d: %s", prefix, i, r.describeAction(d, prefix+"  "))
		}
	}
	ap += fmt.Sprintf("%sExec: %s", prefix, r.describeActionExec(actionName))
	if len(a.GetDependencies()) > 0 {
		ap += fmt.Sprintf("%sDependencies:", prefix)
		for i, d := range a.GetDependencies() {
			ap += fmt.Sprintf("%s%d: %s", prefix, i, r.describeAction(d, prefix+"  "))
		}
	}
	if len(a.GetRecoveryActions()) > 0 {
		ap += fmt.Sprintf("%sRecoveryActions:", prefix)
		for i, d := range a.GetRecoveryActions() {
			ap += fmt.Sprintf("%s%d: %s", prefix, i, r.describeAction(d, prefix+"  "))
		}
	}
	return ap
}

// describeActionExec describes the action's exec function with details.
func (r *recoveryEngine) describeActionExec(actionName string) string {
	a := r.getAction(actionName)
	er := a.GetExecName()
	if len(a.GetExecExtraArgs()) > 0 {
		er += fmt.Sprintf(", Args: %s", a.GetExecExtraArgs())
	}
	return er
}

// initCache initializes cache on engine.
// The function extracted to supported testing.
func (r *recoveryEngine) initCache() {
	r.actionResultsCache = make(map[string]error, len(r.plan.GetActions()))
	r.recoveryUsageCache = make(map[recoveryUsageKey]error)
}

// actionResultFromCache reads action's result from cache.
func (r *recoveryEngine) actionResultFromCache(actionName string) (err error, ok bool) {
	err, ok = r.actionResultsCache[actionName]
	return err, ok
}

// cacheActionResult sets action's result to the cache.
func (r *recoveryEngine) cacheActionResult(actionName string, err error) {
	switch r.getAction(actionName).GetRunControl() {
	case planpb.RunControl_RERUN_AFTER_RECOVERY, planpb.RunControl_RUN_ONCE:
		r.actionResultsCache[actionName] = err
	case planpb.RunControl_ALWAYS_RUN:
		// Do not cache the value
	}
}

// resetCacheAfterSuccessfulRecoveryAction resets cache for actions
// with run-control=RERUN_AFTER_RECOVERY.
func (r *recoveryEngine) resetCacheAfterSuccessfulRecoveryAction() {
	for name, a := range r.plan.GetActions() {
		if a.GetRunControl() == planpb.RunControl_RERUN_AFTER_RECOVERY {
			delete(r.actionResultsCache, name)
		}
	}
}

// isRecoveryUsed checks if recovery action is used in plan or action level scope.
func (r *recoveryEngine) isRecoveryUsed(actionName, recoveryName string) bool {
	k := recoveryUsageKey{
		action:   actionName,
		recovery: recoveryName,
	}
	// If the recovery has been used in previous actions then it can be in
	// the action result cache.
	if err, ok := r.actionResultsCache[recoveryName]; ok {
		r.recoveryUsageCache[k] = err
	}
	_, ok := r.recoveryUsageCache[k]
	return ok
}

// registerRecoveryUsage sets recovery action usage to the cache.
func (r *recoveryEngine) registerRecoveryUsage(actionName, recoveryName string, err error) {
	r.recoveryUsageCache[recoveryUsageKey{
		action:   actionName,
		recovery: recoveryName,
	}] = err
}

// recoveryUsageKey holds action and action's recovery name as key for recovery-usage cache.
type recoveryUsageKey struct {
	action   string
	recovery string
}
