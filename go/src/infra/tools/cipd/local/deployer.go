// Copyright 2014 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package local

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"

	"github.com/luci/luci-go/common/logging"

	"infra/tools/cipd/common"
)

// TODO(vadimsh): How to handle path conflicts between two packages? Currently
// the last one installed wins.

// TODO(vadimsh): Use some sort of file lock to reduce a chance of corruption.

// File system layout of a site directory <root> for "symlink" install method:
// <root>/.cipd/pkgs/
//   <package name digest>/
//     _current -> symlink to fea3ab83440e9dfb813785e16d4101f331ed44f4
//     fea3ab83440e9dfb813785e16d4101f331ed44f4/
//       bin/
//         tool
//         ...
//       ...
// bin/
//    tool -> symlink to ../.cipd/pkgs/<package name digest>/_current/bin/tool
//    ...
//
// Where <package name digest> is derived from a package name. It doesn't have
// to be reversible though, since the package name is still stored in the
// installed package manifest and can be read from there.
//
// Some efforts are made to make sure that during the deployment a window of
// inconsistency in the file system is as small as possible.
//
// For "copy" install method everything is much simpler: files are directly
// copied to the site root directory and .cipd/pkgs/* contains only metadata,
// such as manifest file with a list of extracted files (to know what to
// uninstall).

// Deployer knows how to unzip and place packages into site root directory.
type Deployer interface {
	// DeployInstance installs a specific instance of a package into a site root
	// directory. It unpacks the package into <root>/.cipd/pkgs/*, and rearranges
	// symlinks to point to unpacked files. It tries to make it as "atomic" as
	// possible. Returns information about the deployed instance.
	DeployInstance(PackageInstance) (common.Pin, error)

	// CheckDeployed checks whether a given package is deployed and returns
	// information about installed version if it is or error if not.
	CheckDeployed(packageName string) (common.Pin, error)

	// FindDeployed returns a list of packages deployed to a site root.
	FindDeployed() (out []common.Pin, err error)

	// RemoveDeployed deletes a package given its name.
	RemoveDeployed(packageName string) error

	// TempFile returns os.File located in <root>/tmp/*.
	TempFile(prefix string) (*os.File, error)
}

// NewDeployer return default Deployer implementation.
func NewDeployer(root string, logger logging.Logger) Deployer {
	var err error
	if root == "" {
		err = fmt.Errorf("Site root path is not provided")
	} else {
		root, err = filepath.Abs(filepath.Clean(root))
	}
	if err != nil {
		return errDeployer{err}
	}
	if logger == nil {
		logger = logging.Null()
	}
	return &deployerImpl{NewFileSystem(root, logger), logger}
}

////////////////////////////////////////////////////////////////////////////////
// Implementation that returns error on all requests.

type errDeployer struct{ err error }

func (d errDeployer) DeployInstance(PackageInstance) (common.Pin, error)   { return common.Pin{}, d.err }
func (d errDeployer) CheckDeployed(packageName string) (common.Pin, error) { return common.Pin{}, d.err }
func (d errDeployer) FindDeployed() (out []common.Pin, err error)          { return nil, d.err }
func (d errDeployer) RemoveDeployed(packageName string) error              { return d.err }
func (d errDeployer) TempFile(prefix string) (*os.File, error)             { return nil, d.err }

////////////////////////////////////////////////////////////////////////////////
// Real deployer implementation.

// packagesDir is a subdirectory of site root to extract packages to.
const packagesDir = siteServiceDir + "/pkgs"

// currentSymlink is a name of a symlink that points to latest deployed version.
// Used on Linux and Mac.
const currentSymlink = "_current"

// currentTxt is a name of a text file with instance ID of latest deployed
// version. Used on Windows.
const currentTxt = "_current.txt"

// deployerImpl implements Deployer interface.
type deployerImpl struct {
	fs     FileSystem
	logger logging.Logger
}

func (d *deployerImpl) DeployInstance(inst PackageInstance) (common.Pin, error) {
	pin := inst.Pin()
	d.logger.Infof("Deploying %s into %s", pin, d.fs.Root())

	// Be paranoid.
	if err := common.ValidatePin(pin); err != nil {
		return common.Pin{}, err
	}
	if _, err := d.fs.EnsureDirectory(d.fs.Root()); err != nil {
		return common.Pin{}, err
	}

	// Extract new version to the .cipd/pkgs/* guts. For "symlink" install mode it
	// is the final destination. For "copy" install mode it's a temp destination
	// and files will be moved to the site root later (in addToSiteRoot call).
	// ExtractPackageInstance knows how to build full paths and how to atomically
	// extract a package. No need to delete garbage if it fails.
	pkgPath := d.packagePath(pin.PackageName)
	destPath := filepath.Join(pkgPath, pin.InstanceID)
	if err := ExtractInstance(inst, NewFileSystemDestination(destPath, d.fs)); err != nil {
		return common.Pin{}, err
	}
	newManifest, err := d.readManifest(destPath)
	if err != nil {
		return common.Pin{}, err
	}

	// Remember currently deployed version (to remove it later). Do not freak out
	// if it's not there (prevInstanceID == "") or broken (err != nil).
	prevInstanceID, err := d.getCurrentInstanceID(pkgPath)
	prevManifest := Manifest{}
	if err == nil && prevInstanceID != "" {
		prevManifest, err = d.readManifest(filepath.Join(pkgPath, prevInstanceID))
	}
	if err != nil {
		d.logger.Warningf("Previous version of the package is broken: %s", err)
		prevManifest = Manifest{} // to make sure prevManifest.Files == nil.
	}

	// Install all new files to the site root.
	err = d.addToSiteRoot(newManifest.Files, newManifest.InstallMode, pkgPath, destPath)
	if err != nil {
		d.fs.EnsureDirectoryGone(destPath)
		return common.Pin{}, err
	}

	// Mark installed instance as a current one. After this call the package is
	// considered installed and the function must not fail. All cleanup below is
	// best effort.
	if err = d.setCurrentInstanceID(pkgPath, pin.InstanceID); err != nil {
		d.fs.EnsureDirectoryGone(destPath)
		return common.Pin{}, err
	}

	// Wait for async cleanup to finish.
	wg := sync.WaitGroup{}
	defer wg.Wait()

	// Remove old instance directory completely.
	if prevInstanceID != "" && prevInstanceID != pin.InstanceID {
		wg.Add(1)
		go func() {
			defer wg.Done()
			d.fs.EnsureDirectoryGone(filepath.Join(pkgPath, prevInstanceID))
		}()
	}

	// Remove no longer present files from the site root directory.
	if len(prevManifest.Files) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			toKeep := map[string]bool{}
			for _, f := range newManifest.Files {
				toKeep[f.Name] = true
			}
			toKill := []FileInfo{}
			for _, f := range prevManifest.Files {
				if !toKeep[f.Name] {
					toKill = append(toKill, f)
				}
			}
			d.removeFromSiteRoot(toKill)
		}()
	}

	// Verify it's all right.
	newPin, err := d.CheckDeployed(pin.PackageName)
	if err == nil && newPin.InstanceID != pin.InstanceID {
		err = fmt.Errorf("Other instance (%s) was deployed concurrently", newPin.InstanceID)
	}
	if err == nil {
		d.logger.Infof("Successfully deployed %s", pin)
	} else {
		d.logger.Errorf("Failed to deploy %s: %s", pin, err)
	}
	return newPin, err
}

func (d *deployerImpl) CheckDeployed(pkg string) (common.Pin, error) {
	current, err := d.getCurrentInstanceID(d.packagePath(pkg))
	if err != nil {
		return common.Pin{}, err
	}
	if current == "" {
		return common.Pin{}, fmt.Errorf("package %s is not installed", pkg)
	}
	return common.Pin{
		PackageName: pkg,
		InstanceID:  current,
	}, nil
}

func (d *deployerImpl) FindDeployed() ([]common.Pin, error) {
	// Directories with packages are direct children of .cipd/pkgs/.
	pkgs := filepath.Join(d.fs.Root(), filepath.FromSlash(packagesDir))
	infos, err := ioutil.ReadDir(pkgs)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	found := map[string]common.Pin{}
	keys := []string{}
	for _, info := range infos {
		if !info.IsDir() {
			continue
		}
		// Attempt to read the manifest. If it is there -> valid package is found.
		pkgPath := filepath.Join(pkgs, info.Name())
		currentID, err := d.getCurrentInstanceID(pkgPath)
		if err != nil || currentID == "" {
			continue
		}
		manifest, err := d.readManifest(filepath.Join(pkgPath, currentID))
		if err != nil {
			continue
		}
		// Ignore duplicate entries, they can appear if someone messes with pkgs/*
		// structure manually.
		if _, ok := found[manifest.PackageName]; !ok {
			keys = append(keys, manifest.PackageName)
			found[manifest.PackageName] = common.Pin{
				PackageName: manifest.PackageName,
				InstanceID:  currentID,
			}
		}
	}

	// Sort by package name.
	sort.Strings(keys)
	out := make([]common.Pin, len(found))
	for i, k := range keys {
		out[i] = found[k]
	}
	return out, nil
}

func (d *deployerImpl) RemoveDeployed(packageName string) error {
	d.logger.Infof("Removing %s from %s", packageName, d.fs.Root())
	if err := common.ValidatePackageName(packageName); err != nil {
		return err
	}
	pkgPath := d.packagePath(packageName)

	// Read the manifest of the currently installed version.
	manifest := Manifest{}
	currentID, err := d.getCurrentInstanceID(pkgPath)
	if err == nil && currentID != "" {
		manifest, err = d.readManifest(filepath.Join(pkgPath, currentID))
	}

	// Warn, but continue with removal anyway. EnsureDirectoryGone call below
	// will nuke everything (even if it's half broken).
	if err != nil {
		d.logger.Warningf("Package %s is in a broken state: %s", packageName, err)
	} else {
		d.removeFromSiteRoot(manifest.Files)
	}
	return d.fs.EnsureDirectoryGone(pkgPath)
}

func (d *deployerImpl) TempFile(prefix string) (*os.File, error) {
	dir, err := d.fs.EnsureDirectory(filepath.Join(d.fs.Root(), siteServiceDir, "tmp"))
	if err != nil {
		return nil, err
	}
	return ioutil.TempFile(dir, prefix)
}

////////////////////////////////////////////////////////////////////////////////
// Utility methods.

// packagePath returns a path to a package directory in .cipd/pkgs/.
func (d *deployerImpl) packagePath(pkg string) string {
	rel := filepath.Join(filepath.FromSlash(packagesDir), packageNameDigest(pkg))
	abs, err := d.fs.RootRelToAbs(rel)
	if err != nil {
		msg := fmt.Sprintf("Can't get absolute path of '%s'", rel)
		d.logger.Errorf("%s", msg)
		panic(msg)
	}
	return abs
}

// getCurrentInstanceID returns instance ID of currently installed instance
// given a path to a package directory (.cipd/pkgs/<name>). It returns ("", nil)
// if no package is installed there.
func (d *deployerImpl) getCurrentInstanceID(packageDir string) (string, error) {
	var current string
	var err error
	if runtime.GOOS == "windows" {
		var bytes []byte
		bytes, err = ioutil.ReadFile(filepath.Join(packageDir, currentTxt))
		if err == nil {
			current = strings.TrimSpace(string(bytes))
		}
	} else {
		current, err = os.Readlink(filepath.Join(packageDir, currentSymlink))
	}
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	if err = common.ValidateInstanceID(current); err != nil {
		return "", fmt.Errorf(
			"pointer to currently installed instance doesn't look like a valid instance id: %s", err)
	}
	return current, nil
}

// setCurrentInstanceID changes a pointer to currently installed instance ID. It
// takes a path to a package directory (.cipd/pkgs/<name>) as input.
func (d *deployerImpl) setCurrentInstanceID(packageDir string, instanceID string) error {
	if err := common.ValidateInstanceID(instanceID); err != nil {
		return err
	}
	if runtime.GOOS == "windows" {
		return d.fs.EnsureFile(filepath.Join(packageDir, currentTxt), []byte(instanceID), 0666)
	}
	return d.fs.EnsureSymlink(filepath.Join(packageDir, currentSymlink), instanceID)
}

// readManifest reads package manifest given a path to a package instance
// (.cipd/pkgs/<name>/<instance id>).
func (d *deployerImpl) readManifest(instanceDir string) (Manifest, error) {
	manifestPath := filepath.Join(instanceDir, filepath.FromSlash(manifestName))
	r, err := os.Open(manifestPath)
	if err != nil {
		return Manifest{}, err
	}
	defer r.Close()
	manifest, err := readManifest(r)
	if err != nil {
		return Manifest{}, err
	}
	// Older packages do not have Files section in the manifest, so reconstruct it
	// from actual files on disk.
	if len(manifest.Files) == 0 {
		if manifest.Files, err = scanPackageDir(instanceDir, d.logger); err != nil {
			return Manifest{}, err
		}
	}
	return manifest, nil
}

// addToSiteRoot moves or symlinks files into the site root directory (depending
// on passed installMode).
func (d *deployerImpl) addToSiteRoot(files []FileInfo, installMode InstallMode, pkgDir, srcDir string) error {
	// On Windows only InstallModeCopy is supported.
	if runtime.GOOS == "windows" {
		installMode = InstallModeCopy
	} else if installMode == "" {
		installMode = InstallModeSymlink // default on non-Windows
	}
	if err := ValidateInstallMode(installMode); err != nil {
		return err
	}

	for _, f := range files {
		// E.g. bin/tool.
		relPath := filepath.FromSlash(f.Name)
		destAbs, err := d.fs.RootRelToAbs(relPath)
		if err != nil {
			d.logger.Warningf("Invalid relative path %q: %s", relPath, err)
			return err
		}
		if installMode == InstallModeSymlink {
			// E.g. <root>/.cipd/pkgs/name/_current/bin/tool.
			targetAbs := filepath.Join(pkgDir, currentSymlink, relPath)
			// E.g. ../.cipd/pkgs/name/_current/bin/tool.
			targetRel, err := filepath.Rel(filepath.Dir(destAbs), targetAbs)
			if err != nil {
				d.logger.Warningf("Can't get relative path from %s to %s", filepath.Dir(destAbs), targetAbs)
				return err
			}
			if err = d.fs.EnsureSymlink(destAbs, targetRel); err != nil {
				d.logger.Warningf("Failed to create symlink for %s", relPath)
				return err
			}
		} else if installMode == InstallModeCopy {
			// E.g. <root>/.cipd/pkgs/name/<id>/bin/tool.
			srcAbs := filepath.Join(srcDir, relPath)
			if err := d.fs.Replace(srcAbs, destAbs); err != nil {
				d.logger.Warningf("Failed to move %s to %s: %s", srcAbs, destAbs, err)
				return err
			}
		} else {
			// Should not happen. ValidateInstallMode checks this.
			return fmt.Errorf("impossible state")
		}
	}
	// Best effort cleanup of empty directories after all files has been moved.
	if installMode == InstallModeCopy {
		d.removeEmptyDirs(srcDir)
	}
	return nil
}

// removeFromSiteRoot deletes files from the site root directory. Best effort.
// Logs errors and carries on.
func (d *deployerImpl) removeFromSiteRoot(files []FileInfo) {
	for _, f := range files {
		absPath, err := d.fs.RootRelToAbs(filepath.FromSlash(f.Name))
		if err != nil {
			d.logger.Warningf("Refusing to remove %q: %s", f.Name, err)
			continue
		}
		if err := d.fs.EnsureFileGone(absPath); err != nil {
			d.logger.Warningf("Failed to remove a file from the site root: %s", err)
		}
	}
	d.removeEmptyDirs(d.fs.Root())
}

// removeEmptyDirs recursive removes empty directory subtrees (e.g. subtrees
// that do not have any files). The root directory itself will not be removed
// (even if it's empty). Best effort, logs errors.
func (d *deployerImpl) removeEmptyDirs(root string) {
	// TODO(vadimsh): Implement.
}

////////////////////////////////////////////////////////////////////////////////
// Utility functions.

// packageNameDigest returns a filename to use for naming a package directory in
// the file system. Using package names as is can introduce problems on file
// systems with path length limits (on Windows in particular). Returns last two
// components of the package name + stripped SHA1 of the whole package name.
func packageNameDigest(pkg string) string {
	// Be paranoid.
	err := common.ValidatePackageName(pkg)
	if err != nil {
		panic(err.Error())
	}

	// Grab stripped SHA1 of the full package name.
	digest := sha1.Sum([]byte(pkg))
	hash := base64.URLEncoding.EncodeToString(digest[:])[:10]

	// Grab last <= 2 components of the package path.
	chunks := strings.Split(pkg, "/")
	if len(chunks) > 2 {
		chunks = chunks[len(chunks)-2:]
	}

	// Join together with '_' as separator.
	chunks = append(chunks, hash)
	return strings.Join(chunks, "_")
}

// scanPackageDir finds a set of regular files (and symlinks) in a package
// instance directory and returns them as FileInfo structs (with slash-separated
// paths relative to dir directory). Skips package service directories (.cipdpkg
// and .cipd) since they contain package deployer gut files, not something that
// needs to be deployed.
func scanPackageDir(dir string, l logging.Logger) ([]FileInfo, error) {
	out := []FileInfo{}
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		if rel == packageServiceDir || rel == siteServiceDir {
			return filepath.SkipDir
		}
		if info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 {
			symlink := ""
			ok := true
			if info.Mode()&os.ModeSymlink != 0 {
				symlink, err = os.Readlink(path)
				if err != nil {
					l.Warningf("Can't readlink %q, skipping: %s", path, err)
					ok = false
				}
			}
			if ok {
				out = append(out, FileInfo{
					Name:       filepath.ToSlash(rel),
					Size:       uint64(info.Size()),
					Executable: (info.Mode().Perm() & 0111) != 0,
					Symlink:    symlink,
				})
			}
		}
		return nil
	})
	return out, err
}
