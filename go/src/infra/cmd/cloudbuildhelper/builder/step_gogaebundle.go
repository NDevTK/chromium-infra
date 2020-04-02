// Copyright 2020 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package builder

import (
	"context"
	"encoding/json"
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	"strings"

	"go.chromium.org/luci/common/data/stringset"
	"go.chromium.org/luci/common/errors"
	"go.chromium.org/luci/common/logging"

	"infra/cmd/cloudbuildhelper/fileset"
)

// runRunBuildStep executes manifest.RunBuildStep.
func runGoGAEBundleBuildStep(ctx context.Context, inv *stepRunnerInv) error {
	logging.Infof(ctx, "Bundling %q", inv.BuildStep.GoGAEBundle)

	mainDir, err := filepath.Abs(inv.BuildStep.GoGAEBundle)
	if err != nil {
		return errors.Annotate(err, "failed to convert the path %q to absolute", inv.BuildStep.GoGAEBundle).Err()
	}

	// Get a build.Context as if we are building for linux amd64.
	bc := build.Default
	bc.GOARCH = "amd64"
	bc.GOOS = "linux"
	bc.Dir = mainDir

	// Find where main package is actually located.
	mainPkg, err := bc.ImportDir(mainDir, 0)
	if err != nil {
		return errors.Annotate(err, "failed to locate the go code").Err()
	}
	logging.Infof(ctx, "Import path is %q", mainPkg.ImportPath)
	if mainPkg.Name != "main" {
		return errors.Annotate(err, "only \"main\" package can be bundled, got %q", mainPkg.Name).Err()
	}

	// Drop a hint to gaedeploy that helps it to reconstruct the GOPATH.
	var info struct {
		Go struct {
			ImportPath string `json:"import_path"`
		} `json:"go"`
	}
	info.Go.ImportPath = mainPkg.ImportPath
	blob, err := json.MarshalIndent(&info, "", "  ")
	if err != nil {
		panic(fmt.Sprintf("impossible: %s", err))
	}
	err = inv.addBlobToOutput(ctx, filepath.Join(inv.BuildStep.Dest, ".gaedeploy.json"), blob)
	if err != nil {
		return err
	}

	// Copy all files that make up "main" package (they can be only at the root
	// of `mainDir`), and copy all non-go files recursively (they can potentially
	// be referenced by static_files in app.yaml). We'll deal with Go dependencies
	// separately.
	err = inv.addFilesToOutput(ctx, mainDir, inv.BuildStep.Dest, func(absPath string, isDir bool) bool {
		if isDir {
			return false // do not exclude directories, may have contain static files
		}
		rel, err := filepath.Rel(mainDir, absPath)
		if err != nil {
			panic(fmt.Sprintf("impossible: %s", err))
		}
		switch {
		// Do not exclude non-code files regardless of where they are.
		case !isGoSourceFile(rel):
			return false
		// Exclude code files not in the mainDir. If they are needed, they'll be
		// discovered by the next step.
		case rel != filepath.Base(rel):
			return true
		// For code files in the mainDir, pick up only ones matching the build
		// context (linux amd64).
		default:
			matches, err := bc.MatchFile(mainDir, rel)
			if err != nil {
				logging.Warningf(ctx, "Failed to check whether %q matches the build context, skipping it: %s", absPath, err)
				return true
			}
			return !matches
		}
	})
	if err != nil {
		return err
	}

	// Find all packages that mainPkg transitively depends on.
	if inv.State.goStdlib == nil {
		logging.Infof(ctx, "Enumerating stdlib packages to know when to skip them...")
		inv.State.goStdlib, err = findStdlib(bc.GOROOT)
		if err != nil {
			return err
		}
	}
	logging.Infof(ctx, "Discovering transitive dependencies...")
	deps, err := findDeps(&bc, mainPkg, inv.State.goStdlib)
	if err != nil {
		return err
	}

	if inv.State.goDeps == nil {
		inv.State.goDeps = stringset.New(len(deps))
	}

	// Add them all to the tarball.
	logging.Infof(ctx, "Found %d dependencies. Copying them to the output...", len(deps))
	for _, pkg := range deps {
		if !inv.State.goDeps.Add(pkg.ImportPath) {
			continue // added it already in some previous build step
		}

		srcDir := filepath.Join(pkg.SrcRoot, pkg.ImportPath)
		dstDir := filepath.Join("_gopath", "src", pkg.ImportPath)

		// All non-test files Go compiler ever touches.
		var names []string
		names = append(names, pkg.GoFiles...)
		names = append(names, pkg.CgoFiles...)
		names = append(names, pkg.CFiles...)
		names = append(names, pkg.CXXFiles...)
		names = append(names, pkg.MFiles...)
		names = append(names, pkg.HFiles...)
		names = append(names, pkg.FFiles...)
		names = append(names, pkg.SFiles...)

		// Add them all to the tarball.
		for _, name := range names {
			err := inv.Output.AddFromDisk(filepath.Join(srcDir, name), filepath.Join(dstDir, name), nil)
			if err != nil {
				return errors.Annotate(err, "failed to copy %q from %q", name, srcDir).Err()
			}
		}
	}

	// Drop a script that can be used to test sanity of this tarball.
	envPath := filepath.Join("_gopath", "env")
	return inv.Output.AddFromMemory(envPath, []byte(envScript), &fileset.File{
		Executable: true,
	})
}

// isGoSourceFile returns true if rel may be read by Go compiler.
//
// See https://golang.org/src/go/build/build.go.
func isGoSourceFile(rel string) bool {
	switch filepath.Ext(rel) {
	case ".go", ".c", ".cc", ".cxx", ".cpp", ".m", ".s", ".h", ".hh", ".hpp", ".hxx", ".f", ".F", ".f90", ".S", ".sx", ".swig", ".swigcxx":
		return true
	default:
		return false
	}
}

// findStdlib examines GOROOT to find names of most of stdlib packages.
func findStdlib(goRoot string) (s stringset.Set, err error) {
	s = stringset.New(100)

	dir := filepath.Join(goRoot, "src")
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || !info.IsDir() {
			return err
		}

		// Convert to an import path.
		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)

		// Skip some not really importable magic directories.
		if rel == "cmd" || rel == "vendor" || rel == "internal" || strings.HasSuffix(rel, "/internal") {
			return filepath.SkipDir
		}

		s.Add(rel)
		return err
	})
	return
}

// findDeps finds all non-stdlib dependencies that `pkg` depends on.
func findDeps(bc *build.Context, pkg *build.Package, stdlib stringset.Set) (deps []*build.Package, err error) {
	goRootSrc := filepath.Join(bc.GOROOT, "src")
	visitedDeps := stringset.New(0)

	// In go `import "a/b/c"` can import physically different packages depending
	// on what package contains the import, due to existence of magical "vendor"
	// folder. We account for that by using importFrom struct as a map key instead
	// of just ImportPath. It records both what package is imported and from
	// where.
	type importFrom struct {
		path string
		from string
	}
	queue := make([]importFrom, 0, len(pkg.Imports))
	visited := make(map[importFrom]bool, len(pkg.Imports))

	enqueue := func(i importFrom) {
		// Skip stdlib package. This check is fast, but may not be 100% reliable
		// since `stdlib` is constructed in not very rigorous way. We'll do a
		// separate strict check later. The check here is just an optimization.
		if !stdlib.Has(i.path) && !visited[i] {
			queue = append(queue, i)
			visited[i] = true
		}
	}

	for _, importPath := range pkg.Imports {
		enqueue(importFrom{
			path: importPath,
			from: pkg.Dir,
		})
	}

	for len(queue) != 0 {
		cur := queue[0]
		queue = queue[1:]

		pkg, err := bc.Import(cur.path, cur.from, 0)
		if err != nil {
			return nil, err
		}

		// Skip stdlib packages (in case our simplistic first check failed). Note
		// that relying on this strong check exclusively makes the code very slow,
		// since it keeps revisiting stdlib packages (via bc.Import) a lot.
		if pkg.SrcRoot == goRootSrc {
			continue
		}

		// Get the actual physical location of the package on disk. And make sure
		// we didn't pick it up already.
		if visitedDeps.Add(filepath.Join(pkg.SrcRoot, pkg.ImportPath)) {
			deps = append(deps, pkg)

			// Visit its imports.
			for _, importPath := range pkg.Imports {
				enqueue(importFrom{
					path: importPath,
					from: pkg.Dir,
				})
			}
		}
	}

	return deps, nil
}

// envScript spits out a script that modifies Go env vars to point to files
// in the tarball. Can be used to manually test the tarball's soundness.
const envScript = `#!/usr/bin/env bash
cd $(dirname "${BASH_SOURCE[0]}")

echo "export GOARCH=amd64"
echo "export GOOS=linux"
echo "export GO111MODULE=off"
echo "export GOPATH=$(pwd)"
`
