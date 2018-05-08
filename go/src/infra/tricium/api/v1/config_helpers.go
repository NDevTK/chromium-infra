// Copyright 2017 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package tricium

import (
	"errors"
	"fmt"

	"golang.org/x/net/context"

	"go.chromium.org/luci/server/auth"
)

// ProjectIsKnown checks if the provided project is known to the Tricium service.
//
// DEPRECATED, will be removed with Projects is removed from ServiceConfig,
// see https://crbug.com/824558.
func ProjectIsKnown(sc *ServiceConfig, project string) bool {
	for _, p := range sc.Projects {
		if p.Name == project {
			return true
		}
	}
	return false
}

// LookupProjectDetails lookups up the project details entry for the provided project.
//
// Unknown projects results in nil.
//
// DEPRECATED, will be removed with Projects is removed from ServiceConfig,
// see https://crbug.com/824558.
func LookupProjectDetails(sc *ServiceConfig, project string) *ProjectDetails {
	for _, p := range sc.Projects {
		if p.Name == project {
			return p
		}
	}
	return nil
}

// LookupRepoDetails gets corresponding RepoDetails for the
// repo that matches the given AnalyzeRequest.
//
// Returns nil if no such repo is found.
func LookupRepoDetails(pc *ProjectConfig, request *AnalyzeRequest) *RepoDetails {
	requestURL := urlFromRequest(request)
	if requestURL == "" {
		return nil
	}
	for _, repo := range pc.Repos {
		repoURL := urlFromRepo(repo)
		if requestURL == repoURL {
			return repo
		}
	}
	return nil
}

func urlFromRequest(request *AnalyzeRequest) string {
	if revision := request.GetGerritRevision(); revision != nil {
		return revision.GitUrl
	}
	if commit := request.GetGitCommit(); commit != nil {
		return commit.Url
	}
	return ""
}

func urlFromRepo(repo *RepoDetails) string {
	if project := repo.GetGerritProject(); project != nil {
		return project.GitUrl
	}
	if repo := repo.GetGitRepo(); repo != nil {
		return repo.Url
	}
	return ""
}

// CanRequest checks the current user can make service requests for the project.
func CanRequest(c context.Context, pc *ProjectConfig) (bool, error) {
	return checkAcls(c, pc, Acl_REQUESTER)
}

// CanRead checks the current user can read project results.
func CanRead(c context.Context, pc *ProjectConfig) (bool, error) {
	return checkAcls(c, pc, Acl_READER)
}

func checkAcls(c context.Context, pc *ProjectConfig, role Acl_Role) (bool, error) {
	var groups []string
	for _, acl := range pc.Acls {
		if acl.Role != role {
			continue
		}
		if acl.Group != "" {
			groups = append(groups, acl.Group)
		}
		if acl.Identity == string(auth.CurrentIdentity(c)) {
			return true, nil
		}
	}
	ok, err := auth.IsMember(c, groups...)
	if err != nil {
		return false, fmt.Errorf("failed to check member in group(s) (%v): %v", groups, err)
	}
	return ok, nil
}

// LookupProjectFunction looks up the given function in the project config.
func LookupProjectFunction(pc *ProjectConfig, function string) *Function {
	return lookupFunction(pc.Functions, function)
}

// LookupServiceFunction looks up the given function in the service config.
func LookupServiceFunction(sc *ServiceConfig, function string) *Function {
	return lookupFunction(sc.Functions, function)
}

func lookupFunction(functions []*Function, function string) *Function {
	for _, f := range functions {
		if f.Name == function {
			return f
		}
	}
	return nil
}

// LookupImplForPlatform returns the first impl providing data for the provided platform.
func LookupImplForPlatform(f *Function, platform Platform_Name) *Impl {
	for _, i := range f.Impls {
		if i.ProvidesForPlatform == platform {
			return i
		}
	}
	if platform == Platform_ANY && len(f.Impls) != 0 {
		return f.Impls[0]
	}
	return nil
}

// LookupPlatform returns the first platform matching the provided platform name.
func LookupPlatform(sc *ServiceConfig, platform Platform_Name) *Platform_Details {
	for _, p := range sc.Platforms {
		if p.Name == platform {
			return p
		}
	}
	return nil
}

// SupportsPlatform checks if the provided function has an implementation providing
// data for the provided platform.
func SupportsPlatform(f *Function, platform Platform_Name) bool {
	for _, i := range f.Impls {
		if i.ProvidesForPlatform == platform {
			return true
		}
	}
	return false
}

// SupportsConfig checks if the function has a config def matching the provided config.
func SupportsConfig(f *Function, config *Config) bool {
	for _, c := range f.ConfigDefs {
		if c.Name == config.Name {
			return true
		}
	}
	return false
}

// GetRecipePackages returns the base service recipe packages for recipe-based implementations.
func GetRecipePackages(sc *ServiceConfig, platform Platform_Name) ([]*CipdPackage, error) {
	if len(sc.RecipePackages) == 0 {
		return nil, errors.New("service recipe packages missing")
	}
	// TODO(emso): adjust packages for platform.
	return sc.RecipePackages, nil
}

// GetRecipeCmd returns the base service command for recipe-based implementations.
func GetRecipeCmd(sc *ServiceConfig, platform Platform_Name) (*Cmd, error) {
	if sc.GetRecipeCmd() == nil {
		return nil, errors.New("service recipe command missing")
	}
	// TODO(emso): Adjust for platform?
	return &Cmd{
		Exec: sc.RecipeCmd.Exec,
		Args: sc.RecipeCmd.Args,
	}, nil
}

// ValidateFunction checks if the function config entry is valid.
//
// A valid function config entry has a name, valid deps and valid impl entries.
// Note that there are more requirements for a function config to be fully
// valid in a merged config; for instance, data dependencies are required.
func ValidateFunction(f *Function, sc *ServiceConfig) error {
	switch f.GetType() {
	case Function_NONE:
		return errors.New("missing type in function config")
	case Function_ANALYZER:
		if f.GetProvides() != Data_RESULTS {
			return errors.New("analyzer function must return results")
		}
	case Function_ISOLATOR:
		if f.GetProvides() == Data_RESULTS {
			return errors.New("isolator functions must not return results")
		}
	}
	if f.GetName() == "" {
		return errors.New("missing name in function config")
	}
	if f.GetNeeds() == Data_NONE {
		return errors.New("missing input type in function config")
	}
	if f.GetProvides() == Data_NONE {
		return errors.New("missing output type in function config")
	}
	pm := make(map[Platform_Name]bool)
	for _, i := range f.Impls {
		needs, err := LookupDataTypeDetails(sc, f.Needs)
		if err != nil {
			return errors.New("function has impl that needs unknown data type")
		}
		provides, err := LookupDataTypeDetails(sc, f.Provides)
		if err != nil {
			return errors.New("function has impl that provides unknown data type")
		}
		if err = validateImpl(i, sc, needs, provides); err != nil {
			return fmt.Errorf("invalid impl for function %s: %v", f.Name, err)
		}
		if i.ProvidesForPlatform == Platform_ANY {
			continue
		}
		if _, ok := pm[i.ProvidesForPlatform]; ok {
			return fmt.Errorf("multiple impl providing data for platform %v", i.ProvidesForPlatform)
		}
		pm[i.ProvidesForPlatform] = true
	}
	return nil
}

// validateImpl checks if the Function Impl entry is valid.
//
// A valid Impl entry has a valid runtime platform, specifies platforms for
// data dependencies if necessary, has a cmd or recipe, and has a deadline.
func validateImpl(i *Impl, sc *ServiceConfig, needs *Data_TypeDetails, provides *Data_TypeDetails) error {
	if i.GetRuntimePlatform() == Platform_ANY {
		return errors.New("must provide runtime platform for impl")
	}
	runtime := LookupPlatform(sc, i.GetRuntimePlatform())
	if runtime == nil {
		return fmt.Errorf("impl using unknown runtime platform: %v", i.RuntimePlatform)
	}
	if !runtime.GetHasRuntime() {
		return errors.New("must provide a runtime platform that has a runtime")
	}
	if needs.GetIsPlatformSpecific() && i.GetNeedsForPlatform() == Platform_ANY {
		return errors.New("must specify platform for needed platform-specific data type")
	}
	if provides.GetIsPlatformSpecific() && i.GetProvidesForPlatform() == Platform_ANY {
		return errors.New("must specify platform for provided platform-specific data type")
	}
	if i.GetCmd() == nil && i.GetRecipe() == nil {
		return errors.New("must include either command or recipe in impl")
	}
	if i.GetDeadline() == 0 {
		return errors.New("missing deadline")
	}
	return nil
}

// LookupDataTypeDetails looks up data type details for a given type from the
// provided service config.
func LookupDataTypeDetails(sc *ServiceConfig, dt Data_Type) (*Data_TypeDetails, error) {
	for _, d := range sc.GetDataDetails() {
		if d.Type == dt {
			return d, nil
		}
	}
	return nil, fmt.Errorf("data type undefined: %v", dt)
}
