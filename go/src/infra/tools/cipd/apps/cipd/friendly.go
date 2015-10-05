// Copyright 2015 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/luci/luci-go/client/authcli"
	"github.com/luci/luci-go/common/auth"
	"github.com/luci/luci-go/common/logging"

	"github.com/maruel/subcommands"

	"infra/tools/cipd"
	"infra/tools/cipd/local"
)

////////////////////////////////////////////////////////////////////////////////
// Site root path resolution.

// findSiteRoot returns a directory R such as R/.cipd exists and p is inside
// R or p is R. Returns empty string if no such directory.
func findSiteRoot(p string) string {
	for {
		if isSiteRoot(p) {
			return p
		}
		// Dir returns "/" (or C:\\) when it encounters the root directory. It is
		// the only case when Dir(...) return value ends with separator.
		parent := filepath.Dir(p)
		if parent[len(parent)-1] == filepath.Separator {
			// It is possible disk root has .cipd directory, check it.
			if isSiteRoot(parent) {
				return parent
			}
			return ""
		}
		p = parent
	}
}

// optionalSiteRoot takes a path to a site root or an empty string. If some
// path is given, it normalizes it and ensures that it is indeed a site root
// directory. If empty string is given, it discovers a site root for current
// directory.
func optionalSiteRoot(siteRoot string) (string, error) {
	if siteRoot == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		siteRoot = findSiteRoot(cwd)
		if siteRoot == "" {
			return "", fmt.Errorf("directory %s is not in a site root, use 'init' to create one", cwd)
		}
		return siteRoot, nil
	}
	siteRoot, err := filepath.Abs(siteRoot)
	if err != nil {
		return "", err
	}
	if !isSiteRoot(siteRoot) {
		return "", fmt.Errorf("directory %s doesn't look like a site root, use 'init' to create one", siteRoot)
	}
	return siteRoot, nil
}

// isSiteRoot returns true if <p>/.cipd exists.
func isSiteRoot(p string) bool {
	fi, err := os.Stat(filepath.Join(p, local.SiteServiceDir))
	return err == nil && fi.IsDir()
}

////////////////////////////////////////////////////////////////////////////////
// Config file parsing.

// installationSiteConfig is stored in .cipd/config.json.
type installationSiteConfig struct {
	// ServiceURL is https://<hostname> of a backend to use by default.
	ServiceURL string
	// DefaultVersion is what version to install if not specified.
	DefaultVersion string
	// TrackedVersions is mapping package name -> version to use in 'update'.
	TrackedVersions map[string]string
}

// read loads JSON from given path.
func (c *installationSiteConfig) read(path string) error {
	*c = installationSiteConfig{}
	r, err := os.Open(path)
	if err != nil {
		return err
	}
	defer r.Close()
	return json.NewDecoder(r).Decode(c)
}

// write dumps JSON to given path.
func (c *installationSiteConfig) write(path string) error {
	blob, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, blob, 0666)
}

// readConfig reads config, returning default one if missing.
func readConfig(siteRoot string) (installationSiteConfig, error) {
	path := filepath.Join(siteRoot, local.SiteServiceDir, "config.json")
	c := installationSiteConfig{}
	if err := c.read(path); err != nil && !os.IsNotExist(err) {
		return c, err
	}
	return c, nil
}

////////////////////////////////////////////////////////////////////////////////
// High level wrapper around site root.

// installationSite represents a site root directory with config and optional
// cipd.Client instance configured to install packages into that root.
type installationSite struct {
	siteRoot string                  // path to a site root directory
	cfg      *installationSiteConfig // parsed .cipd/config.json file
	client   cipd.Client             // initialized by initClient()
	log      logging.Logger
}

// getInstallationSite finds site root directory, reads config and constructs
// installationSite object. If siteRoot is "", will find a site root based on
// the current directory, otherwise will use siteRoot. Doesn't create any new
// files or directories, just reads what's on disk.
func getInstallationSite(siteRoot string) (*installationSite, error) {
	siteRoot, err := optionalSiteRoot(siteRoot)
	if err != nil {
		return nil, err
	}
	cfg, err := readConfig(siteRoot)
	if err != nil {
		return nil, err
	}
	return &installationSite{siteRoot, &cfg, nil, log}, nil
}

// initInstallationSite creates new site root directory on disk and returns
// corresponding *installationSite object. It does a bunch of sanity checks
// (like whether rootDir is empty) that are skipped if 'force' is set to true.
func initInstallationSite(rootDir string, force bool) (*installationSite, error) {
	rootDir, err := filepath.Abs(rootDir)
	if err != nil {
		return nil, err
	}

	// rootDir is inside an existing site root?
	existing := findSiteRoot(rootDir)
	if existing != "" {
		msg := fmt.Sprintf("directory %s is already inside a site root (%s)", rootDir, existing)
		if !force {
			return nil, errors.New(msg)
		}
		fmt.Fprintf(os.Stderr, "Warning: %s.\n", msg)
	}

	// Attempting to use in a non empty directory?
	items, err := ioutil.ReadDir(rootDir)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	if len(items) != 0 {
		msg := fmt.Sprintf("directory %s is not empty", rootDir)
		if !force {
			return nil, errors.New(msg)
		}
		fmt.Fprintf(os.Stderr, "Warning: %s.\n", msg)
	}

	// Good to go.
	if err = os.MkdirAll(filepath.Join(rootDir, local.SiteServiceDir), 0777); err != nil {
		return nil, err
	}
	site, err := getInstallationSite(rootDir)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Site root initialized at %s.\n", rootDir)
	return site, nil
}

// initClient initializes cipd.Client to use to talk to backend. Can be called
// only once. Use it directly via site.client.
func (site *installationSite) initClient(authFlags authcli.Flags) (err error) {
	if site.client != nil {
		return errors.New("client is already initialized")
	}
	serviceOpts := ServiceOptions{
		authFlags:  authFlags,
		serviceURL: site.cfg.ServiceURL,
	}
	site.client, err = serviceOpts.makeCipdClient(site.siteRoot)
	return
}

// closeClient closes the underlying cipd.Client if necessary.
func (site *installationSite) closeClient() {
	if site.client != nil {
		site.client.Close()
		site.client = nil
	}
}

// modifyConfig reads config file, calls callback to mutate it, then writes
// it back.
func (site *installationSite) modifyConfig(cb func(cfg *installationSiteConfig) error) error {
	path := filepath.Join(site.siteRoot, local.SiteServiceDir, "config.json")
	c := installationSiteConfig{}
	if err := c.read(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := cb(&c); err != nil {
		return err
	}
	return c.write(path)
}

// installedPackages discovers versions of packages installed in the site. If
// pkgs is empty array, it returns list of all installed packages.
func (site *installationSite) installedPackages(pkgs []string) ([]pinInfo, error) {
	d := local.NewDeployer(site.siteRoot, site.log)

	// List all?
	if len(pkgs) == 0 {
		pins, err := d.FindDeployed()
		if err != nil {
			return nil, err
		}
		output := make([]pinInfo, len(pins))
		for i, pin := range pins {
			cpy := pin
			output[i] = pinInfo{
				Pkg:      pin.PackageName,
				Pin:      &cpy,
				Tracking: site.cfg.TrackedVersions[pin.PackageName],
			}
		}
		return output, nil
	}

	// List specific packages only.
	output := make([]pinInfo, len(pkgs))
	for i, pkgName := range pkgs {
		pin, err := d.CheckDeployed(pkgName)
		if err == nil {
			output[i] = pinInfo{
				Pkg:      pkgName,
				Pin:      &pin,
				Tracking: site.cfg.TrackedVersions[pkgName],
			}
		} else {
			output[i] = pinInfo{
				Pkg:      pkgName,
				Tracking: site.cfg.TrackedVersions[pkgName],
				Err:      err.Error(),
			}
		}
	}
	return output, nil
}

// installPackage installs (or updates) a package. If 'force' is true, it will
// reinstall the package even if it is already marked as installed at requested
// version. On errors returns (nil, error).
func (site *installationSite) installPackage(pkgName, version string, force bool) (*pinInfo, error) {
	if site.client == nil {
		return nil, errors.New("client is not initialized")
	}

	// Figure out what exactly (what instance ID) to install.
	if version == "" {
		version = site.cfg.DefaultVersion
	}
	if version == "" {
		version = "latest"
	}
	resolved, err := site.client.ResolveVersion(pkgName, version)
	if err != nil {
		return nil, err
	}

	// Already installed?
	doInstall := true
	if !force {
		d := local.NewDeployer(site.siteRoot, site.log)
		existing, err := d.CheckDeployed(pkgName)
		if err == nil && existing == resolved {
			fmt.Printf("Package %s is up-to-date.\n", pkgName)
			doInstall = false
		}
	}

	// Go for it.
	if doInstall {
		fmt.Printf("Installing %s (version %q)...\n", pkgName, version)
		if err := site.client.FetchAndDeployInstance(resolved); err != nil {
			return nil, err
		}
	}

	// Update config saying what version to track. Remove tracking if an exact
	// instance ID was requested.
	trackedVersion := ""
	if version != resolved.InstanceID {
		trackedVersion = version
	}
	err = site.modifyConfig(func(cfg *installationSiteConfig) error {
		if cfg.TrackedVersions == nil {
			cfg.TrackedVersions = map[string]string{}
		}
		if cfg.TrackedVersions[pkgName] != trackedVersion {
			if trackedVersion == "" {
				fmt.Printf("Package %s is now pinned to %q.\n", pkgName, resolved.InstanceID)
			} else {
				fmt.Printf("Package %s is now tracking %q.\n", pkgName, trackedVersion)
			}
		}
		if trackedVersion == "" {
			delete(cfg.TrackedVersions, pkgName)
		} else {
			cfg.TrackedVersions[pkgName] = trackedVersion
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Success.
	return &pinInfo{
		Pkg:      pkgName,
		Pin:      &resolved,
		Tracking: trackedVersion,
	}, nil
}

////////////////////////////////////////////////////////////////////////////////
// Common command line flags.

// siteRootOptions defines command line flag for specifying existing site root
// directory. 'init' subcommand is NOT using it, since it creates a new site
// root, not reusing an existing one.
type siteRootOptions struct {
	rootDir string
}

func (opts *siteRootOptions) registerFlags(f *flag.FlagSet) {
	f.StringVar(
		&opts.rootDir, "root", "", "Path to an installation site root directory. "+
			"If omitted will try to discovery it by examining parent directories.")
}

////////////////////////////////////////////////////////////////////////////////
// 'init' subcommand.

var cmdInit = &subcommands.Command{
	UsageLine: "init [root dir] [options]",
	ShortDesc: "sets up a new site root directory to install packages into",
	LongDesc: "Sets up a new site root directory to install packages into.\n\n" +
		"Uses current working directory by default.\n" +
		"Unless -force is given, the new site root directory should be empty (or " +
		"do not exist at all) and not be under some other existing site root. " +
		"The command will create <root>/.cipd subdirectory with some " +
		"configuration files. This directory is used by CIPD client to keep " +
		"track of what is installed in the site root.",
	CommandRun: func() subcommands.CommandRun {
		c := &initRun{}
		c.registerBaseFlags()
		c.Flags.BoolVar(&c.force, "force", false, "Create the site root even if the directory is not empty or already under another site root directory.")
		c.Flags.StringVar(&c.serviceURL, "service-url", "", "URL of a backend to use instead of the default one.")
		return c
	},
}

type initRun struct {
	Subcommand

	force      bool
	serviceURL string
}

func (c *initRun) Run(a subcommands.Application, args []string) int {
	if !c.init(args, 0, 1) {
		return 1
	}
	rootDir := "."
	if len(args) == 1 {
		rootDir = args[0]
	}
	site, err := initInstallationSite(rootDir, c.force)
	if err != nil {
		return c.done(nil, err)
	}
	err = site.modifyConfig(func(cfg *installationSiteConfig) error {
		cfg.ServiceURL = c.serviceURL
		return nil
	})
	return c.done(site.siteRoot, err)
}

////////////////////////////////////////////////////////////////////////////////
// 'install' subcommand.

var cmdInstall = &subcommands.Command{
	UsageLine: "install <package> [<version>] [options]",
	ShortDesc: "installs or updates a package",
	LongDesc:  "Installs or updates a package.",
	CommandRun: func() subcommands.CommandRun {
		c := &installRun{}
		c.registerBaseFlags()
		c.authFlags.Register(&c.Flags, auth.Options{})
		c.siteRootOptions.registerFlags(&c.Flags)
		c.Flags.BoolVar(&c.force, "force", false, "Refetch and reinstall the package even if already installed.")
		return c
	},
}

type installRun struct {
	Subcommand
	authFlags authcli.Flags
	siteRootOptions

	force bool
}

func (c *installRun) Run(a subcommands.Application, args []string) int {
	if !c.init(args, 1, 2) {
		return 1
	}

	// Pkg and version to install.
	pkgName := args[0]
	version := ""
	if len(args) == 2 {
		version = args[1]
	}

	// Auto initialize site root directory if necessary. Don't be too aggressive
	// about it though (do not use force=true). Will do anything only if
	// c.rootDir points to an empty directory.
	var site *installationSite
	rootDir, err := optionalSiteRoot(c.rootDir)
	if err == nil {
		site, err = getInstallationSite(rootDir)
	} else {
		site, err = initInstallationSite(c.rootDir, false)
		if err != nil {
			err = fmt.Errorf("can't auto initialize cipd site root (%s), use 'init'", err)
		}
	}
	if err != nil {
		return c.done(nil, err)
	}
	if err = site.initClient(c.authFlags); err != nil {
		return c.done(nil, err)
	}
	defer site.closeClient()

	return c.done(site.installPackage(pkgName, version, c.force))
}

////////////////////////////////////////////////////////////////////////////////
// 'installed' subcommand.

var cmdInstalled = &subcommands.Command{
	UsageLine: "installed [<package> <package> ...] [options]",
	ShortDesc: "lists packages installed in the site root",
	LongDesc:  "Lists packages installed in the site root.",
	CommandRun: func() subcommands.CommandRun {
		c := &installedRun{}
		c.registerBaseFlags()
		c.siteRootOptions.registerFlags(&c.Flags)
		return c
	},
}

type installedRun struct {
	Subcommand
	siteRootOptions
}

func (c *installedRun) Run(a subcommands.Application, args []string) int {
	if !c.init(args, 0, -1) {
		return 1
	}
	site, err := getInstallationSite(c.rootDir)
	if err != nil {
		return c.done(nil, err)
	}
	return c.doneWithPins(site.installedPackages(args))
}
