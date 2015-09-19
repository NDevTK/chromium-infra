// Copyright 2014 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

/*
Client side of for Chrome Infra Package Deployer.

TODO: write more.

Subcommand starting with 'pkg-' are low level commands operating on package
files on disk.
*/
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/luci/luci-go/client/authcli"
	"github.com/luci/luci-go/common/auth"
	"github.com/luci/luci-go/common/logging/gologger"
	gol "github.com/op/go-logging"

	"github.com/maruel/subcommands"

	cipd_lib "infra/libs/cipd"

	"infra/tools/cipd"
	"infra/tools/cipd/common"
	"infra/tools/cipd/local"
)

// log is also overwritten in Subcommand's init.
var log = gologger.Get()

////////////////////////////////////////////////////////////////////////////////
// Common subcommand functions.

// Subcommand is a base of all CIPD subcommands. It defines some common flags,
// such as logging and JSON output parameters.
type Subcommand struct {
	subcommands.CommandRunBase

	jsonOutput string
	verbose    bool
}

// registerBaseFlags registers common flags used by all subcommands.
func (c *Subcommand) registerBaseFlags() {
	c.Flags.StringVar(&c.jsonOutput, "json-output", "", "Path to write operation results to.")
	c.Flags.BoolVar(&c.verbose, "verbose", false, "Enable more logging.")
}

// init initializes subcommand state (including global logger). It ensures all
// required positional and flag-like parameters are set. Returns true if they
// are, or false (and prints to stderr) if not.
func (c *Subcommand) init(args []string, minPosCount, maxPosCount int) bool {
	// Check number of expected positional arguments.
	if maxPosCount == 0 && len(args) != 0 {
		c.printError(makeCLIError("unexpected arguments %v", args))
		return false
	}
	if len(args) < minPosCount || len(args) > maxPosCount {
		var err error
		if minPosCount == maxPosCount {
			err = makeCLIError("expecting %d positional argument, got %d instead", minPosCount, len(args))
		} else {
			err = makeCLIError(
				"expecting from %d to %d positional arguments, got %d instead",
				minPosCount, maxPosCount, len(args))
		}
		c.printError(err)
		return false
	}

	// Check required unset flags.
	unset := []*flag.Flag{}
	c.Flags.VisitAll(func(f *flag.Flag) {
		if strings.HasPrefix(f.DefValue, "<") && f.Value.String() == f.DefValue {
			unset = append(unset, f)
		}
	})
	if len(unset) != 0 {
		missing := []string{}
		for _, f := range unset {
			missing = append(missing, f.Name)
		}
		c.printError(makeCLIError("missing required flags: %v", missing))
		return false
	}

	// Setup logger.
	loggerConfig := gologger.LoggerConfig{
		Format: `[P%{pid} %{time:15:04:05.000} %{shortfile} %{level:.1s}] %{message}`,
		Out:    os.Stderr,
		Level:  gol.INFO,
	}
	if c.verbose {
		loggerConfig.Level = gol.DEBUG
	}
	log = loggerConfig.Get()
	return true
}

// printError prints error to stderr (recognizing commandLineError).
func (c *Subcommand) printError(err error) {
	if _, ok := err.(commandLineError); ok {
		fmt.Fprintf(os.Stderr, "Bad command line: %s.\n\n", err)
		c.Flags.Usage()
	} else {
		fmt.Fprintln(os.Stderr, err)
	}
}

// writeJSONOutput writes result to JSON output file. It returns original error
// if it is non-nil.
func (c *Subcommand) writeJSONOutput(result interface{}, err error) error {
	// -json-output flag wasn't specified.
	if c.jsonOutput == "" {
		return err
	}

	// Prepare the body of the output file.
	var body struct {
		Error  string      `json:"error,omitempty"`
		Result interface{} `json:"result,omitempty"`
	}
	if err != nil {
		body.Error = err.Error()
	}
	body.Result = result
	out, e := json.MarshalIndent(&body, "", "  ")
	if e != nil {
		if err == nil {
			err = e
		} else {
			fmt.Fprintf(os.Stderr, "Failed to serialize JSON output: %s\n", e)
		}
		return err
	}

	e = ioutil.WriteFile(c.jsonOutput, out, 0600)
	if e != nil {
		if err == nil {
			err = e
		} else {
			fmt.Fprintf(os.Stderr, "Failed write JSON output to %s: %s\n", c.jsonOutput, e)
		}
		return err
	}

	return err
}

// done is called as a last step of processing a subcommand. It dumps command
// result (or error) to JSON output file, prints error message and generates
// process exit code.
func (c *Subcommand) done(result interface{}, err error) int {
	err = c.writeJSONOutput(result, err)
	if err != nil {
		c.printError(err)
		return 1
	}
	return 0
}

// commandLineError is used to tag errors related to CLI.
type commandLineError struct {
	error
}

// makeCLIError returns new commandLineError.
func makeCLIError(msg string, args ...interface{}) error {
	return commandLineError{fmt.Errorf(msg, args...)}
}

////////////////////////////////////////////////////////////////////////////////
// ServiceOptions mixin.

// ServiceOptions defines command line arguments related to communication
// with the remote service. Subcommands that interact with the network embed it.
type ServiceOptions struct {
	authFlags  authcli.Flags
	serviceURL string
}

func (opts *ServiceOptions) registerFlags(f *flag.FlagSet) {
	f.StringVar(&opts.serviceURL, "service-url", "", "URL of a backend to use instead of the default one.")
	opts.authFlags.Register(f, auth.Options{})
}

func (opts *ServiceOptions) makeCipdClient(root string) (cipd.Client, error) {
	authOpts, err := opts.authFlags.Options()
	if err != nil {
		return nil, err
	}
	return cipd.NewClient(cipd.ClientOptions{
		ServiceURL: opts.serviceURL,
		Root:       root,
		Logger:     log,
		AuthenticatedClientFactory: func() (*http.Client, error) {
			return auth.AuthenticatedClient(auth.OptionalLogin, auth.NewAuthenticator(authOpts))
		},
	}), nil
}

////////////////////////////////////////////////////////////////////////////////
// InputOptions mixin.

// PackageVars holds array of '-pkg-var' command line options.
type PackageVars map[string]string

func (vars *PackageVars) String() string {
	// String() for empty vars used in -help output.
	if len(*vars) == 0 {
		return "key:value"
	}
	chunks := make([]string, 0, len(*vars))
	for k, v := range *vars {
		chunks = append(chunks, fmt.Sprintf("%s:%s", k, v))
	}
	return strings.Join(chunks, " ")
}

// Set is called by 'flag' package when parsing command line options.
func (vars *PackageVars) Set(value string) error {
	// <key>:<value> pair.
	chunks := strings.Split(value, ":")
	if len(chunks) != 2 {
		return makeCLIError("expecting <key>:<value> pair, got %q", value)
	}
	(*vars)[chunks[0]] = chunks[1]
	return nil
}

// InputOptions defines command line arguments that specify where to get data
// for a new package. Subcommands that build packages embed it.
type InputOptions struct {
	// Path to *.yaml file with package definition.
	packageDef string
	vars       PackageVars

	// Alternative to 'pkg-def'.
	packageName string
	inputDir    string
	installMode local.InstallMode
}

func (opts *InputOptions) registerFlags(f *flag.FlagSet) {
	opts.vars = PackageVars{}

	// Interface to accept package definition file.
	f.StringVar(&opts.packageDef, "pkg-def", "", "*.yaml file that defines what to put into the package.")
	f.Var(&opts.vars, "pkg-var", "Variables accessible from package definition file.")

	// Interface to accept a single directory (alternative to -pkg-def).
	f.StringVar(&opts.packageName, "name", "", "Package name (unused with -pkg-def).")
	f.StringVar(&opts.inputDir, "in", "", "Path to a directory with files to package (unused with -pkg-def).")
	f.Var(&opts.installMode, "install-mode",
		"How the package should be installed: \"copy\" or \"symlink\" (unused with -pkg-def).")
}

// prepareInput processes InputOptions by collecting all files to be added to
// a package and populating BuildInstanceOptions. Caller is still responsible to
// fill out Output field of BuildInstanceOptions.
func (opts *InputOptions) prepareInput() (local.BuildInstanceOptions, error) {
	out := local.BuildInstanceOptions{Logger: log}

	// Handle -name and -in if defined. Do not allow -pkg-def and -pkg-var in that case.
	if opts.inputDir != "" {
		if opts.packageName == "" {
			return out, makeCLIError("missing required flag: -name")
		}
		if opts.packageDef != "" {
			return out, makeCLIError("-pkg-def and -in can not be used together")
		}
		if len(opts.vars) != 0 {
			return out, makeCLIError("-pkg-var and -in can not be used together")
		}

		// Simply enumerate files in the directory.
		var files []local.File
		files, err := local.ScanFileSystem(opts.inputDir, opts.inputDir, nil)
		if err != nil {
			return out, err
		}
		out = local.BuildInstanceOptions{
			Input:       files,
			PackageName: opts.packageName,
			InstallMode: opts.installMode,
			Logger:      log,
		}
		return out, nil
	}

	// Handle -pkg-def case. -in is "" (already checked), reject -name.
	if opts.packageDef != "" {
		if opts.packageName != "" {
			return out, makeCLIError("-pkg-def and -name can not be used together")
		}
		if opts.installMode != "" {
			return out, makeCLIError("-install-mode is ignored if -pkd-def is used")
		}

		// Parse the file, perform variable substitution.
		f, err := os.Open(opts.packageDef)
		if err != nil {
			return out, err
		}
		defer f.Close()
		pkgDef, err := local.LoadPackageDef(f, opts.vars)
		if err != nil {
			return out, err
		}

		// Scan the file system. Package definition may use path relative to the
		// package definition file itself, so pass its location.
		fmt.Println("Enumerating files to zip...")
		files, err := pkgDef.FindFiles(filepath.Dir(opts.packageDef))
		if err != nil {
			return out, err
		}
		out = local.BuildInstanceOptions{
			Input:       files,
			PackageName: pkgDef.Package,
			VersionFile: pkgDef.VersionFile(),
			InstallMode: pkgDef.InstallMode,
			Logger:      log,
		}
		return out, nil
	}

	// All command line options are missing.
	return out, makeCLIError("-pkg-def or -name/-in are required")
}

////////////////////////////////////////////////////////////////////////////////
// RefsOptions mixin.

// Refs holds an array of '-ref' command line options.
type Refs []string

func (refs *Refs) String() string {
	// String() for empty vars used in -help output.
	if len(*refs) == 0 {
		return "ref"
	}
	return strings.Join(*refs, " ")
}

// Set is called by 'flag' package when parsing command line options.
func (refs *Refs) Set(value string) error {
	err := common.ValidatePackageRef(value)
	if err != nil {
		return commandLineError{err}
	}
	*refs = append(*refs, value)
	return nil
}

// RefsOptions defines command line arguments for commands that accept a set
// of refs.
type RefsOptions struct {
	refs Refs
}

func (opts *RefsOptions) registerFlags(f *flag.FlagSet) {
	opts.refs = []string{}
	f.Var(&opts.refs, "ref", "A ref to point to the package instance (can be used multiple times).")
}

////////////////////////////////////////////////////////////////////////////////
// TagsOptions mixin.

// Tags holds an array of '-tag' command line options.
type Tags []string

func (tags *Tags) String() string {
	// String() for empty vars used in -help output.
	if len(*tags) == 0 {
		return "key:value"
	}
	return strings.Join(*tags, " ")
}

// Set is called by 'flag' package when parsing command line options.
func (tags *Tags) Set(value string) error {
	err := common.ValidateInstanceTag(value)
	if err != nil {
		return commandLineError{err}
	}
	*tags = append(*tags, value)
	return nil
}

// TagsOptions defines command line arguments for commands that accept a set
// of tags.
type TagsOptions struct {
	tags Tags
}

func (opts *TagsOptions) registerFlags(f *flag.FlagSet) {
	opts.tags = []string{}
	f.Var(&opts.tags, "tag", "A tag to attach to the package instance (can be used multiple times).")
}

////////////////////////////////////////////////////////////////////////////////
// Support for running operations concurrently.

// batchOperation defines what to do with a packages matching a prefix.
type batchOperation struct {
	client        cipd.Client
	packagePrefix string   // a package name or a prefix
	packages      []string // packages to operate on, overrides packagePrefix
	callback      func(pkg string) (common.Pin, error)
}

// pinOrError is passed through channels and also dumped to JSON results file.
type pinOrError struct {
	Pkg string      `json:"package"`
	Pin *common.Pin `json:"pin,omitempty"`
	Err string      `json:"error,omitempty"`
}

// expandPkgDir takes a package name or '<prefix>/' and returns a list
// of matching packages (asking backend if necessary). Doesn't recurse, returns
// only direct children.
func expandPkgDir(c cipd.Client, packagePrefix string) ([]string, error) {
	if !strings.HasSuffix(packagePrefix, "/") {
		return []string{packagePrefix}, nil
	}
	pkgs, err := c.ListPackages(packagePrefix, false)
	if err != nil {
		return nil, err
	}
	// Skip directories.
	out := []string{}
	for _, p := range pkgs {
		if !strings.HasSuffix(p, "/") {
			out = append(out, p)
		}
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no packages under %s", packagePrefix)
	}
	return out, nil
}

// performBatchOperation expands a package prefix into a list of packages and
// calls callback for each of them (concurrently) gathering the results.
func performBatchOperation(op batchOperation) ([]pinOrError, error) {
	pkgs := op.packages
	if len(pkgs) == 0 {
		var err error
		pkgs, err = expandPkgDir(op.client, op.packagePrefix)
		if err != nil {
			return nil, err
		}
	}
	return callConcurrently(pkgs, func(pkg string) pinOrError {
		pin, err := op.callback(pkg)
		if err != nil {
			return pinOrError{pkg, nil, err.Error()}
		}
		return pinOrError{pkg, &pin, ""}
	}), nil
}

func callConcurrently(pkgs []string, callback func(pkg string) pinOrError) []pinOrError {
	// Push index through channel to make results ordered as 'pkgs'.
	ch := make(chan struct {
		int
		pinOrError
	})
	for idx, pkg := range pkgs {
		go func(idx int, pkg string) {
			ch <- struct {
				int
				pinOrError
			}{idx, callback(pkg)}
		}(idx, pkg)
	}
	pins := make([]pinOrError, len(pkgs))
	for i := 0; i < len(pkgs); i++ {
		res := <-ch
		pins[res.int] = res.pinOrError
	}
	return pins
}

func printPinsAndError(pins []pinOrError) {
	hasPins := false
	hasErrors := false
	for _, p := range pins {
		if p.Err == "" {
			hasPins = true
		} else {
			hasErrors = true
		}
	}
	if hasPins {
		fmt.Println("Instances:")
		for _, p := range pins {
			if p.Err == "" {
				fmt.Printf("  %s\n", p.Pin)
			}
		}
	}
	if hasErrors {
		fmt.Fprintln(os.Stderr, "Errors:")
		for _, p := range pins {
			if p.Err != "" {
				fmt.Fprintf(os.Stderr, "  %s\n", p.Err)
			}
		}
	}
}

func hasErrors(pins []pinOrError) bool {
	for _, p := range pins {
		if p.Err != "" {
			return true
		}
	}
	return false
}

////////////////////////////////////////////////////////////////////////////////
// 'create' subcommand.

var cmdCreate = &subcommands.Command{
	UsageLine: "create [options]",
	ShortDesc: "builds and uploads a package instance file",
	LongDesc:  "Builds and uploads a package instance file.",
	CommandRun: func() subcommands.CommandRun {
		c := &createRun{}
		c.registerBaseFlags()
		c.InputOptions.registerFlags(&c.Flags)
		c.RefsOptions.registerFlags(&c.Flags)
		c.TagsOptions.registerFlags(&c.Flags)
		c.ServiceOptions.registerFlags(&c.Flags)
		return c
	},
}

type createRun struct {
	Subcommand
	InputOptions
	RefsOptions
	TagsOptions
	ServiceOptions
}

func (c *createRun) Run(a subcommands.Application, args []string) int {
	if !c.init(args, 0, 0) {
		return 1
	}
	return c.done(buildAndUploadInstance(c.InputOptions, c.RefsOptions, c.TagsOptions, c.ServiceOptions))
}

func buildAndUploadInstance(inputOpts InputOptions, refsOpts RefsOptions, tagsOpts TagsOptions, serviceOpts ServiceOptions) (common.Pin, error) {
	f, err := ioutil.TempFile("", "cipd_pkg")
	if err != nil {
		return common.Pin{}, err
	}
	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()
	err = buildInstanceFile(f.Name(), inputOpts)
	if err != nil {
		return common.Pin{}, err
	}
	return registerInstanceFile(f.Name(), refsOpts, tagsOpts, serviceOpts)
}

////////////////////////////////////////////////////////////////////////////////
// 'ensure' subcommand.

var cmdEnsure = &subcommands.Command{
	UsageLine: "ensure [options]",
	ShortDesc: "installs, removes and updates packages in one go",
	LongDesc:  "Installs, removes and updates packages in one go.",
	CommandRun: func() subcommands.CommandRun {
		c := &ensureRun{}
		c.registerBaseFlags()
		c.ServiceOptions.registerFlags(&c.Flags)
		c.Flags.StringVar(&c.rootDir, "root", "<path>", "Path to a installation site root directory.")
		c.Flags.StringVar(&c.listFile, "list", "<path>", "A file with a list of '<package name> <version>' pairs.")
		return c
	},
}

type ensureRun struct {
	Subcommand
	ServiceOptions

	rootDir  string
	listFile string
}

func (c *ensureRun) Run(a subcommands.Application, args []string) int {
	if !c.init(args, 0, 0) {
		return 1
	}
	return c.done(ensurePackages(c.rootDir, c.listFile, c.ServiceOptions))
}

func ensurePackages(root string, desiredStateFile string, serviceOpts ServiceOptions) ([]common.Pin, error) {
	f, err := os.Open(desiredStateFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	client, err := serviceOpts.makeCipdClient(root)
	if err != nil {
		return nil, err
	}
	desiredState, err := client.ProcessEnsureFile(f)
	if err != nil {
		return nil, err
	}
	if err = client.EnsurePackages(desiredState); err != nil {
		return nil, err
	}
	return desiredState, nil
}

////////////////////////////////////////////////////////////////////////////////
// 'resolve' subcommand.

var cmdResolve = &subcommands.Command{
	UsageLine: "resolve <package or package prefix> [options]",
	ShortDesc: "returns concrete package instance ID given a version",
	LongDesc:  "Returns concrete package instance ID given a version.",
	CommandRun: func() subcommands.CommandRun {
		c := &resolveRun{}
		c.registerBaseFlags()
		c.ServiceOptions.registerFlags(&c.Flags)
		c.Flags.StringVar(&c.version, "version", "<version>", "Package version to resolve.")
		return c
	},
}

type resolveRun struct {
	Subcommand
	ServiceOptions

	version string
}

func (c *resolveRun) Run(a subcommands.Application, args []string) int {
	if !c.init(args, 1, 1) {
		return 1
	}
	pins, err := resolveVersion(args[0], c.version, c.ServiceOptions)
	printPinsAndError(pins)
	ret := c.done(pins, err)
	if hasErrors(pins) && ret == 0 {
		return 1
	}
	return ret
}

func resolveVersion(packagePrefix, version string, serviceOpts ServiceOptions) ([]pinOrError, error) {
	client, err := serviceOpts.makeCipdClient("")
	if err != nil {
		return nil, err
	}
	return performBatchOperation(batchOperation{
		client:        client,
		packagePrefix: packagePrefix,
		callback: func(pkg string) (common.Pin, error) {
			return client.ResolveVersion(pkg, version)
		},
	})
}

////////////////////////////////////////////////////////////////////////////////
// 'set-ref' subcommand.

var cmdSetRef = &subcommands.Command{
	UsageLine: "set-ref <package or package prefix> [options]",
	ShortDesc: "moves a ref to point to a given version",
	LongDesc:  "Moves a ref to point to a given version.",
	CommandRun: func() subcommands.CommandRun {
		c := &setRefRun{}
		c.registerBaseFlags()
		c.RefsOptions.registerFlags(&c.Flags)
		c.ServiceOptions.registerFlags(&c.Flags)
		c.Flags.StringVar(&c.version, "version", "<version>", "Package version to point the ref to.")
		return c
	},
}

type setRefRun struct {
	Subcommand
	RefsOptions
	ServiceOptions

	version string
}

func (c *setRefRun) Run(a subcommands.Application, args []string) int {
	if !c.init(args, 1, 1) {
		return 1
	}
	if len(c.refs) == 0 {
		return c.done(nil, makeCLIError("at least one -ref must be provided"))
	}
	pins, err := setRef(args[0], c.version, c.RefsOptions, c.ServiceOptions)
	printPinsAndError(pins)
	ret := c.done(pins, err)
	if hasErrors(pins) && ret == 0 {
		return 1
	}
	return ret
}

func setRef(packagePrefix, version string, refsOpts RefsOptions, serviceOpts ServiceOptions) ([]pinOrError, error) {
	client, err := serviceOpts.makeCipdClient("")
	if err != nil {
		return nil, err
	}

	// Do not touch anything if some packages do not have requested version. So
	// resolve versions first and only then move refs.
	pins, err := performBatchOperation(batchOperation{
		client:        client,
		packagePrefix: packagePrefix,
		callback: func(pkg string) (common.Pin, error) {
			return client.ResolveVersion(pkg, version)
		},
	})
	if err != nil {
		return nil, err
	}
	if hasErrors(pins) {
		printPinsAndError(pins)
		return nil, fmt.Errorf("can't find %q version in all packages, aborting.", version)
	}

	// Prepare for the next batch call.
	packages := []string{}
	pinsToUse := map[string]common.Pin{}
	for _, p := range pins {
		packages = append(packages, p.Pkg)
		pinsToUse[p.Pkg] = *p.Pin
	}

	// Update all refs.
	return performBatchOperation(batchOperation{
		client:   client,
		packages: packages,
		callback: func(pkg string) (common.Pin, error) {
			pin := pinsToUse[pkg]
			for _, ref := range refsOpts.refs {
				if err := client.SetRefWhenReady(ref, pin); err != nil {
					return common.Pin{}, err
				}
			}
			return pin, nil
		},
	})
}

////////////////////////////////////////////////////////////////////////////////
// 'ls' subcommand.

var cmdListPackages = &subcommands.Command{
	UsageLine: "ls [-r] [<prefix string>]",
	ShortDesc: "lists matching packages",
	LongDesc:  "List packages in the given path to which the user has access, optionally recursively.",
	CommandRun: func() subcommands.CommandRun {
		c := &listPackagesRun{}
		c.registerBaseFlags()
		c.ServiceOptions.registerFlags(&c.Flags)
		c.Flags.BoolVar(&c.recursive, "r", false, "Whether to list packages in subdirectories.")
		return c
	},
}

type listPackagesRun struct {
	Subcommand
	ServiceOptions

	recursive bool
}

func (c *listPackagesRun) Run(a subcommands.Application, args []string) int {
	if !c.init(args, 0, 1) {
		return 1
	}
	path := ""
	if len(args) == 1 {
		path = args[0]
	}
	return c.done(listPackages(path, c.recursive, c.ServiceOptions))
}

func listPackages(path string, recursive bool, serviceOpts ServiceOptions) ([]string, error) {
	client, err := serviceOpts.makeCipdClient("")
	if err != nil {
		return nil, err
	}
	packages, err := client.ListPackages(path, recursive)
	if err != nil {
		return nil, err
	}
	if len(packages) == 0 {
		fmt.Println("No matching packages.")
	} else {
		for _, p := range packages {
			fmt.Println(p)
		}
	}
	return packages, nil
}

////////////////////////////////////////////////////////////////////////////////
// 'acl-list' subcommand.

var cmdListACL = &subcommands.Command{
	UsageLine: "acl-list <package subpath>",
	ShortDesc: "lists package path Access Control List",
	LongDesc:  "Lists package path Access Control List.",
	CommandRun: func() subcommands.CommandRun {
		c := &listACLRun{}
		c.registerBaseFlags()
		c.ServiceOptions.registerFlags(&c.Flags)
		return c
	},
}

type listACLRun struct {
	Subcommand
	ServiceOptions
}

func (c *listACLRun) Run(a subcommands.Application, args []string) int {
	if !c.init(args, 1, 1) {
		return 1
	}
	return c.done(nil, listACL(args[0], c.ServiceOptions))
}

func listACL(packagePath string, serviceOpts ServiceOptions) error {
	client, err := serviceOpts.makeCipdClient("")
	if err != nil {
		return err
	}
	acls, err := client.FetchACL(packagePath)
	if err != nil {
		return err
	}

	// Split by role, drop empty ACLs.
	byRole := map[string][]cipd.PackageACL{}
	for _, a := range acls {
		if len(a.Principals) != 0 {
			byRole[a.Role] = append(byRole[a.Role], a)
		}
	}

	listRoleACL := func(title string, acls []cipd.PackageACL) {
		fmt.Printf("%s:\n", title)
		if len(acls) == 0 {
			fmt.Printf("  none\n")
			return
		}
		for _, a := range acls {
			fmt.Printf("  via %q:\n", a.PackagePath)
			for _, u := range a.Principals {
				fmt.Printf("    %s\n", u)
			}
		}
	}

	listRoleACL("Owners", byRole["OWNER"])
	listRoleACL("Writers", byRole["WRITER"])
	listRoleACL("Readers", byRole["READER"])

	return nil
}

////////////////////////////////////////////////////////////////////////////////
// 'acl-edit' subcommand.

// principalsList is used as custom flag value. It implements flag.Value.
type principalsList []string

func (l *principalsList) String() string {
	return fmt.Sprintf("%v", *l)
}

func (l *principalsList) Set(value string) error {
	// Ensure <type>:<id> syntax is used. Let the backend to validate the rest.
	chunks := strings.Split(value, ":")
	if len(chunks) != 2 {
		return makeCLIError("%q doesn't look like principal id (<type>:<id>)", value)
	}
	*l = append(*l, value)
	return nil
}

var cmdEditACL = &subcommands.Command{
	UsageLine: "acl-edit <package subpath> [options]",
	ShortDesc: "modifies package path Access Control List",
	LongDesc:  "Modifies package path Access Control List.",
	CommandRun: func() subcommands.CommandRun {
		c := &editACLRun{}
		c.registerBaseFlags()
		c.ServiceOptions.registerFlags(&c.Flags)
		c.Flags.Var(&c.owner, "owner", "Users or groups to grant OWNER role.")
		c.Flags.Var(&c.writer, "writer", "Users or groups to grant WRITER role.")
		c.Flags.Var(&c.reader, "reader", "Users or groups to grant READER role.")
		c.Flags.Var(&c.revoke, "revoke", "Users or groups to remove from all roles.")
		return c
	},
}

type editACLRun struct {
	Subcommand
	ServiceOptions

	owner  principalsList
	writer principalsList
	reader principalsList
	revoke principalsList
}

func (c *editACLRun) Run(a subcommands.Application, args []string) int {
	if !c.init(args, 1, 1) {
		return 1
	}
	return c.done(nil, editACL(args[0], c.owner, c.writer, c.reader, c.revoke, c.ServiceOptions))
}

func editACL(packagePath string, owners, writers, readers, revoke principalsList, serviceOpts ServiceOptions) error {
	changes := []cipd.PackageACLChange{}

	makeChanges := func(action cipd.PackageACLChangeAction, role string, list principalsList) {
		for _, p := range list {
			changes = append(changes, cipd.PackageACLChange{
				Action:    action,
				Role:      role,
				Principal: p,
			})
		}
	}

	makeChanges(cipd.GrantRole, "OWNER", owners)
	makeChanges(cipd.GrantRole, "WRITER", writers)
	makeChanges(cipd.GrantRole, "READER", readers)

	makeChanges(cipd.RevokeRole, "OWNER", revoke)
	makeChanges(cipd.RevokeRole, "WRITER", revoke)
	makeChanges(cipd.RevokeRole, "READER", revoke)

	if len(changes) == 0 {
		return nil
	}

	client, err := serviceOpts.makeCipdClient("")
	if err != nil {
		return err
	}
	err = client.ModifyACL(packagePath, changes)
	if err != nil {
		return err
	}
	fmt.Println("ACL changes applied.")
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// 'pkg-build' subcommand.

var cmdBuild = &subcommands.Command{
	UsageLine: "pkg-build [options]",
	ShortDesc: "builds a package instance file",
	LongDesc:  "Builds a package instance producing *.cipd file.",
	CommandRun: func() subcommands.CommandRun {
		c := &buildRun{}
		c.registerBaseFlags()
		c.InputOptions.registerFlags(&c.Flags)
		c.Flags.StringVar(&c.outputFile, "out", "<path>", "Path to a file to write the final package to.")
		return c
	},
}

type buildRun struct {
	Subcommand
	InputOptions

	outputFile string
}

func (c *buildRun) Run(a subcommands.Application, args []string) int {
	if !c.init(args, 0, 0) {
		return 1
	}
	err := buildInstanceFile(c.outputFile, c.InputOptions)
	if err != nil {
		return c.done(nil, err)
	}
	return c.done(inspectInstanceFile(c.outputFile, false))
}

func buildInstanceFile(instanceFile string, inputOpts InputOptions) error {
	// Read the list of files to add to the package.
	buildOpts, err := inputOpts.prepareInput()
	if err != nil {
		return err
	}

	// Prepare the destination, update build options with io.Writer to it.
	out, err := os.OpenFile(instanceFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	buildOpts.Output = out

	// Build the package.
	err = local.BuildInstance(buildOpts)
	out.Close()
	if err != nil {
		os.Remove(instanceFile)
		return err
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// 'pkg-deploy' subcommand.

var cmdDeploy = &subcommands.Command{
	UsageLine: "pkg-deploy <package instance file> [options]",
	ShortDesc: "deploys a package instance file",
	LongDesc:  "Deploys a *.cipd package instance into a site root.",
	CommandRun: func() subcommands.CommandRun {
		c := &deployRun{}
		c.registerBaseFlags()
		c.Flags.StringVar(&c.rootDir, "root", "<path>", "Path to a installation site root directory.")
		return c
	},
}

type deployRun struct {
	Subcommand

	rootDir string
}

func (c *deployRun) Run(a subcommands.Application, args []string) int {
	if !c.init(args, 1, 1) {
		return 1
	}
	return c.done(deployInstanceFile(c.rootDir, args[0]))
}

func deployInstanceFile(root string, instanceFile string) (common.Pin, error) {
	inst, err := local.OpenInstanceFile(instanceFile, "")
	if err != nil {
		return common.Pin{}, err
	}
	defer inst.Close()
	inspectInstance(inst, false)
	return local.NewDeployer(root, log).DeployInstance(inst)
}

////////////////////////////////////////////////////////////////////////////////
// 'pkg-fetch' subcommand.

var cmdFetch = &subcommands.Command{
	UsageLine: "pkg-fetch <package> [options]",
	ShortDesc: "fetches a package instance file from the repository",
	LongDesc:  "Fetches a package instance file from the repository.",
	CommandRun: func() subcommands.CommandRun {
		c := &fetchRun{}
		c.registerBaseFlags()
		c.ServiceOptions.registerFlags(&c.Flags)
		c.Flags.StringVar(&c.version, "version", "<version>", "Package version to fetch.")
		c.Flags.StringVar(&c.outputPath, "out", "<path>", "Path to a file to write fetch to.")
		return c
	},
}

type fetchRun struct {
	Subcommand
	ServiceOptions

	version    string
	outputPath string
}

func (c *fetchRun) Run(a subcommands.Application, args []string) int {
	if !c.init(args, 1, 1) {
		return 1
	}
	return c.done(fetchInstanceFile(args[0], c.version, c.outputPath, c.ServiceOptions))
}

func fetchInstanceFile(packageName, version, instanceFile string, serviceOpts ServiceOptions) (common.Pin, error) {
	client, err := serviceOpts.makeCipdClient("")
	if err != nil {
		return common.Pin{}, err
	}
	pin, err := client.ResolveVersion(packageName, version)
	if err != nil {
		return common.Pin{}, err
	}

	out, err := os.OpenFile(instanceFile, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return common.Pin{}, err
	}
	ok := false
	defer func() {
		if !ok {
			out.Close()
			os.Remove(instanceFile)
		}
	}()

	err = client.FetchInstance(pin, out)
	if err != nil {
		return common.Pin{}, err
	}

	// Verify it (by checking that instanceID matches the file content).
	out.Close()
	ok = true
	inst, err := local.OpenInstanceFile(instanceFile, pin.InstanceID)
	if err != nil {
		os.Remove(instanceFile)
		return common.Pin{}, err
	}
	defer inst.Close()
	inspectInstance(inst, false)
	return inst.Pin(), nil
}

////////////////////////////////////////////////////////////////////////////////
// 'pkg-inspect' subcommand.

var cmdInspect = &subcommands.Command{
	UsageLine: "pkg-inspect <package instance file>",
	ShortDesc: "inspects contents of a package instance file",
	LongDesc:  "Reads contents *.cipd file and prints information about it.",
	CommandRun: func() subcommands.CommandRun {
		c := &inspectRun{}
		c.registerBaseFlags()
		return c
	},
}

type inspectRun struct {
	Subcommand
}

func (c *inspectRun) Run(a subcommands.Application, args []string) int {
	if !c.init(args, 1, 1) {
		return 1
	}
	return c.done(inspectInstanceFile(args[0], true))
}

func inspectInstanceFile(instanceFile string, listFiles bool) (common.Pin, error) {
	inst, err := local.OpenInstanceFile(instanceFile, "")
	if err != nil {
		return common.Pin{}, err
	}
	defer inst.Close()
	inspectInstance(inst, listFiles)
	return inst.Pin(), nil
}

func inspectInstance(inst local.PackageInstance, listFiles bool) {
	fmt.Printf("Instance: %s\n", inst.Pin())
	if listFiles {
		fmt.Println("Package files:")
		for _, f := range inst.Files() {
			if f.Symlink() {
				target, err := f.SymlinkTarget()
				if err != nil {
					fmt.Printf(" E %s (%s)\n", f.Name(), err)
				} else {
					fmt.Printf(" S %s -> %s\n", f.Name(), target)
				}
			} else {
				fmt.Printf(" F %s\n", f.Name())
			}
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// 'pkg-register' subcommand.

var cmdRegister = &subcommands.Command{
	UsageLine: "pkg-register <package instance file>",
	ShortDesc: "uploads and registers package instance in the package repository",
	LongDesc:  "Uploads and registers package instance in the package repository.",
	CommandRun: func() subcommands.CommandRun {
		c := &registerRun{}
		c.registerBaseFlags()
		c.RefsOptions.registerFlags(&c.Flags)
		c.TagsOptions.registerFlags(&c.Flags)
		c.ServiceOptions.registerFlags(&c.Flags)
		return c
	},
}

type registerRun struct {
	Subcommand
	RefsOptions
	TagsOptions
	ServiceOptions
}

func (c *registerRun) Run(a subcommands.Application, args []string) int {
	if !c.init(args, 1, 1) {
		return 1
	}
	return c.done(registerInstanceFile(args[0], c.RefsOptions, c.TagsOptions, c.ServiceOptions))
}

func registerInstanceFile(instanceFile string, refsOpts RefsOptions, tagsOpts TagsOptions, serviceOpts ServiceOptions) (common.Pin, error) {
	inst, err := local.OpenInstanceFile(instanceFile, "")
	if err != nil {
		return common.Pin{}, err
	}
	defer inst.Close()
	client, err := serviceOpts.makeCipdClient("")
	if err != nil {
		return common.Pin{}, err
	}
	inspectInstance(inst, false)
	err = client.RegisterInstance(inst)
	if err != nil {
		return common.Pin{}, err
	}
	err = client.AttachTagsWhenReady(inst.Pin(), tagsOpts.tags)
	if err != nil {
		return common.Pin{}, err
	}
	for _, ref := range refsOpts.refs {
		err = client.SetRefWhenReady(ref, inst.Pin())
		if err != nil {
			return common.Pin{}, err
		}
	}
	return inst.Pin(), nil
}

////////////////////////////////////////////////////////////////////////////////
// Main.

var application = &subcommands.DefaultApplication{
	Name:  "cipd",
	Title: "Chrome Infra Package Deployer",
	Commands: []*subcommands.Command{
		subcommands.CmdHelp,
		cipd_lib.SubcommandVersion,

		// Authentication related commands.
		authcli.SubcommandInfo(auth.Options{Logger: log}, "auth-info"),
		authcli.SubcommandLogin(auth.Options{Logger: log}, "auth-login"),
		authcli.SubcommandLogout(auth.Options{Logger: log}, "auth-logout"),

		// High level commands.
		cmdListPackages,
		cmdCreate,
		cmdEnsure,
		cmdResolve,
		cmdSetRef,

		// ACLs.
		cmdListACL,
		cmdEditACL,

		// Low level pkg-* commands.
		cmdBuild,
		cmdDeploy,
		cmdFetch,
		cmdInspect,
		cmdRegister,
	},
}

func splitCmdLine(args []string) (cmd string, flags []string, pos []string) {
	// No subcomand, just flags.
	if len(args) == 0 || strings.HasPrefix(args[0], "-") {
		return "", args, nil
	}
	// Pick subcommand, than collect all positional args up to a first flag.
	cmd = args[0]
	firstFlagIdx := -1
	for i := 1; i < len(args); i++ {
		if strings.HasPrefix(args[i], "-") {
			firstFlagIdx = i
			break
		}
	}
	// No flags at all.
	if firstFlagIdx == -1 {
		return cmd, nil, args[1:]
	}
	return cmd, args[firstFlagIdx:], args[1:firstFlagIdx]
}

func fixFlagsPosition(args []string) []string {
	// 'flags' package requires positional arguments to be after flags. This is
	// very inconvenient choice, it makes commands like "set-ref" look awkward:
	// Compare "set-ref -ref=abc -version=def package/name" to more natural
	// "set-ref package/name -ref=abc -version=def". Reshuffle arguments to put
	// all positional args at the end of the command line.
	cmd, flags, positional := splitCmdLine(args)
	newArgs := []string{}
	if cmd != "" {
		newArgs = append(newArgs, cmd)
	}
	newArgs = append(newArgs, flags...)
	newArgs = append(newArgs, positional...)
	return newArgs
}

func main() {
	os.Exit(subcommands.Run(application, fixFlagsPosition(os.Args[1:])))
}
