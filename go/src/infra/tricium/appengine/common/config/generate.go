// Copyright 2017 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package config

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"go.chromium.org/luci/common/errors"

	admin "infra/tricium/api/admin/v1"
	"infra/tricium/api/v1"
)

// Generate generates a Tricium workflow based on the provided configs and
// paths to analyze.
//
// The workflow will be computed from the validated and merged config for the
// project in question, and filtered to only include workers relevant to the
// files to be analyzed.
func Generate(sc *tricium.ServiceConfig, pc *tricium.ProjectConfig, files []*tricium.Data_File, gitRef, gitURL string) (*admin.Workflow, error) {
	vpc, err := Validate(sc, pc)
	if err != nil {
		return nil, errors.Annotate(err, "failed to validate project config").Err()
	}
	var workers []*admin.Worker
	functions := map[string]*tricium.Function{}
	for _, s := range vpc.Selections {
		if _, ok := functions[s.Function]; !ok {
			f := tricium.LookupFunction(vpc.Functions, s.Function)
			if f == nil {
				return nil, errors.Annotate(err, "failed to lookup project function").Err()
			}
			functions[s.Function] = f
		}
		shouldInclude, err := includeFunction(functions[s.Function], files)
		if err != nil {
			return nil, errors.Annotate(err, "failed include function check").Err()
		}
		if shouldInclude {
			w, err := createWorker(s, sc, functions[s.Function], gitRef, gitURL)
			if err != nil {
				return nil, errors.Annotate(err, "failed to create worker").Err()
			}
			workers = append(workers, w)
		}
	}
	if err := resolveSuccessorWorkers(sc, workers); err != nil {
		return nil, errors.Annotate(err, "workflow is not sane").Err()
	}

	return &admin.Workflow{
		ServiceAccount:        pc.SwarmingServiceAccount,
		Workers:               workers,
		SwarmingServer:        sc.SwarmingServer,
		BuildbucketServerHost: sc.BuildbucketServerHost,
		IsolateServer:         sc.IsolateServer,
		Functions:             vpc.Functions,
	}, nil
}

// resolveSuccessorWorkers computes successor workers based on data
// dependencies.
//
// The resulting list of successors are added to the Next fields of the
// provided workers. Platform-specific data types add an additional platform
// check to make successors of workers providing a platform-specific type only
// include successors running on that platform.
//
// The resulting workflow is sanity checked and returns an error on failure.
func resolveSuccessorWorkers(sc *tricium.ServiceConfig, workers []*admin.Worker) error {
	specific := map[tricium.Data_Type]bool{}
	for _, d := range sc.GetDataDetails() {
		if _, ok := specific[d.Type]; ok {
			return errors.Reason("multiple declarations of the same data type in the service config, type: %s", d).Err()
		}
		specific[d.Type] = d.IsPlatformSpecific
	}
	needs := map[tricium.Data_Type][]*admin.Worker{}
	for _, w := range workers {
		needs[w.Needs] = append(needs[w.Needs], w)
	}
	for _, w := range workers {
		for _, ws := range needs[w.Provides] {
			if !specific[w.Provides] || specific[w.Provides] && w.ProvidesForPlatform == ws.NeedsForPlatform {
				w.Next = append(w.Next, ws.Name)
			}
		}
	}
	return checkWorkflowSanity(workers)
}

// checkWorkflowSanity checks if the workflow is a tree.
//
// A sane workflow has one path to each worker and includes all workers.
// Multiple paths could mean multiple predecessors to a worker, or could be a
// circularity.
func checkWorkflowSanity(workers []*admin.Worker) error {
	var roots []*admin.Worker
	m := map[string]*admin.Worker{}
	for _, w := range workers {
		if w.Needs == tricium.Data_GIT_FILE_DETAILS {
			roots = append(roots, w)
		}
		m[w.Name] = w
	}
	visited := map[string]*admin.Worker{}
	for _, w := range roots {
		if err := checkWorkerDeps(w, m, visited); err != nil {
			return err
		}
	}
	if len(visited) < len(workers) {
		return errors.Reason("non-accessible workers in workflow").Err()
	}
	return nil
}

// checkWorkerDeps detects joined/circular deps and unknown successors.
//
// Deps are recursively followed via Next pointers for the provided worker.
// The provided visited map is used to track already visited workers to detect
// multiple paths to a worker.
func checkWorkerDeps(w *admin.Worker, m map[string]*admin.Worker, visited map[string]*admin.Worker) error {
	if _, ok := visited[w.Name]; ok {
		return errors.Reason("multiple paths to worker %s", w.Name).Err()
	}
	visited[w.Name] = w
	for _, n := range w.Next {
		wn, ok := m[n]
		if !ok {
			return errors.Reason("unknown next worker %s", n).Err()
		}
		if err := checkWorkerDeps(wn, m, visited); err != nil {
			return err
		}
	}
	return nil
}

// includeFunction checks if an analyzer should be included based on paths.
//
// The paths are checked against the path filters included for the function. If
// there are no path filters or no paths, then the function is included
// without further checking. With both paths and path filters, there needs to
// be at least one path match for the function to be included.
//
// The path filter only applies to the last part of the path.
//
// Also, path filters are only provided for analyzers; analyzer functions are
// always included regardless of path matching.
func includeFunction(f *tricium.Function, files []*tricium.Data_File) (bool, error) {
	if f.Type == tricium.Function_ISOLATOR || len(files) == 0 || len(f.PathFilters) == 0 {
		return true, nil
	}
	for _, file := range files {
		p := file.Path
		for _, filter := range f.PathFilters {
			ok, err := filepath.Match(filter, filepath.Base(p))
			if err != nil {
				return false, errors.Reason("failed to check path filter %s for path %s", filter, p).Err()
			}
			if ok {
				return true, nil
			}
		}
	}
	return false, nil
}

// createWorker creates a worker from the provided function, selection and
// service config.
//
// The provided function is assumed to be verified.
func createWorker(s *tricium.Selection, sc *tricium.ServiceConfig, f *tricium.Function, gitRef, gitURL string) (*admin.Worker, error) {
	i := tricium.LookupImplForPlatform(f, s.Platform) // If verified, there should be an Impl.
	p := tricium.LookupPlatform(sc, s.Platform)       // If verified, the platform should be known.
	// TODO(qyearsley): The character that's used as a separator in worker
	// names should be explicitly disallowed from function and platform
	// names. Currently the character is "_"; a check could be added that
	// the function and worker do not contain "_". If this is not feasible,
	// the separator character could be changed.
	w := &admin.Worker{
		Name:                fmt.Sprintf("%s_%s", s.Function, s.Platform),
		Needs:               f.Needs,
		Provides:            f.Provides,
		NeedsForPlatform:    i.NeedsForPlatform,
		ProvidesForPlatform: i.ProvidesForPlatform,
		RuntimePlatform:     i.RuntimePlatform,
		Dimensions:          p.Dimensions,
		CipdPackages:        i.CipdPackages,
		Deadline:            i.Deadline,
	}
	switch ii := i.Impl.(type) {
	case *tricium.Impl_Recipe:
		recipe := ii.Recipe
		properties := make(map[string]interface{})
		if recipe.Properties != "" {
			err := json.Unmarshal([]byte(recipe.Properties), &properties)
			if err != nil {
				return nil, errors.Annotate(err, "failed to unmarshal").Err()
			}
		}
		for _, c := range s.Configs {
			switch v := c.ValueType.(type) {
			case *tricium.Config_Value:
				properties[c.Name] = v.Value
			case *tricium.Config_ValueJ:
				var value interface{}
				err := json.Unmarshal([]byte(v.ValueJ), &value)
				if err != nil {
					return nil, errors.Annotate(err, "failed to unmarshal value_j").Err()
				}
				properties[c.Name] = value
			default:
				return nil, errors.Reason("please specify value or value_j").Err()
			}
		}
		properties["ref"] = gitRef
		properties["repository"] = gitURL
		properties_bytes, err := json.Marshal(properties)
		if err != nil {
			return nil, errors.Annotate(err, "failed to marshal").Err()
		}
		recipe.Properties = string(properties_bytes)
		w.Impl = &admin.Worker_Recipe{Recipe: ii.Recipe}
	case *tricium.Impl_Cmd:
		cmd := ii.Cmd
		for _, c := range s.Configs {
			switch v := c.ValueType.(type) {
			case *tricium.Config_Value:
				cmd.Args = append(cmd.Args, "--"+c.Name, v.Value)
			case *tricium.Config_ValueJ:
				cmd.Args = append(cmd.Args, "--"+c.Name, v.ValueJ)
			default:
				return nil, errors.Reason("please specify value or value_j").Err()
			}
		}
		w.Impl = &admin.Worker_Cmd{Cmd: ii.Cmd}
	case nil:
		return nil, errors.Reason("missing Impl when constructing worker %s", w.Name).Err()
	default:
		return nil, errors.Reason("Impl.Impl has unexpected type %T", ii).Err()
	}
	return w, nil
}
