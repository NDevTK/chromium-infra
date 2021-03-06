// Copyright 2021 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/golang/protobuf/jsonpb"
	"github.com/maruel/subcommands"
	build_api "go.chromium.org/chromiumos/config/go/build/api"
	"go.chromium.org/chromiumos/config/go/test/api"
	lab_api "go.chromium.org/chromiumos/config/go/test/lab/api"
	"go.chromium.org/luci/auth"
	"go.chromium.org/luci/auth/client/authcli"
	"go.chromium.org/luci/common/cli"
	"go.chromium.org/luci/common/errors"
	"go.chromium.org/luci/lucictx"
	"golang.org/x/sync/errgroup"

	"infra/cmdsupport/cmdlib"
	"infra/cros/cmd/cros-tool-runner/internal/provision"
)

// Provision executes the provisioning for requested devices.
func Provision(authOpts auth.Options) *subcommands.Command {
	return &subcommands.Command{
		UsageLine: "provision -input input.json -output output.json",
		ShortDesc: "Run provisioning for ChromeOS devices",
		LongDesc: `Run provisioning for ChromeOS devices

Tool used to perfrom provisioning OS, components and FW to ChromeOS device specified by ProvisionState.

Example:
cros-tool-runner provision -images docker-images.json -input provision_request.json -output provision_result.json
`,
		CommandRun: func() subcommands.CommandRun {
			c := &runCmd{}
			c.authFlags.Register(&c.Flags, authOpts)
			// Used to provide input by files.
			c.Flags.StringVar(&c.inputPath, "input", "", "The input file contains a jsonproto representation of provision requests (CrosToolRunnerProvisionRequest)")
			c.Flags.StringVar(&c.outputPath, "output", "", "The output file contains a jsonproto representation of provision responses (CrosToolRunnerProvisionResponse)")
			c.Flags.StringVar(&c.imagesPath, "images", "", "The input file contains a jsonproto representation of containers metadata (ContainerMetadata)")
			return c
		},
	}
}

type runCmd struct {
	subcommands.CommandRunBase
	authFlags authcli.Flags

	inputPath  string
	outputPath string
	imagesPath string

	in                 *api.CrosToolRunnerProvisionRequest
	crosDutImage       *build_api.ContainerImageInfo
	crosProvisionImage *build_api.ContainerImageInfo
}

// Run executes the tool.
func (c *runCmd) Run(a subcommands.Application, args []string, env subcommands.Env) int {
	ctx := cli.GetContext(a, c, env)
	out, err := c.innerRun(ctx, a, args, env)
	// Unexpected error will counted as incorrect request data.
	// all expected cases has to generate responses.
	if err != nil && len(out.GetResponses()) == 0 {
		log.Printf("Run: add error to output, %s", err)
		out.Responses = []*api.CrosProvisionResponse{
			{
				Outcome: &api.CrosProvisionResponse_Failure{
					Failure: &api.InstallFailure{
						Reason: api.InstallFailure_REASON_INVALID_REQUEST,
					},
				},
			},
		}
	}
	if err := saveOutput(out, c.outputPath); err != nil {
		log.Printf("Run: %s", err)
	}
	printOutput(out, a)
	if err != nil {
		cmdlib.PrintError(a, err)
		return 1
	}
	return 0
}

func (c *runCmd) innerRun(ctx context.Context, a subcommands.Application, args []string, env subcommands.Env) (*api.CrosToolRunnerProvisionResponse, error) {
	out := &api.CrosToolRunnerProvisionResponse{}
	ctx, err := useSystemAuth(ctx, &c.authFlags)
	if err != nil {
		return out, errors.Annotate(err, "inner run: read system auth").Err()
	}
	req, err := readProvisionRequest(c.inputPath)
	if err != nil {
		return out, errors.Annotate(err, "inner run").Err()
	}

	cm, err := readContainersMetadata(c.imagesPath)
	if err != nil {
		return out, errors.Annotate(err, "inner run").Err()
	}

	if isEmptyEndPoint(req.GetInventoryServer()) {
		return out, errors.Annotate(err, "inner run: inventory service is not provided").Err()
	}

	// TODO(otabek): Listen signal to cancel execution by client.

	// errgroup returns the first error but doesn't stop execution of other goroutines.
	g, ctx := errgroup.WithContext(ctx)
	provisionResults := make([]*api.CrosProvisionResponse, len(req.GetDevices()))
	// Each DUT will run in parallel execution.
	for i, device := range req.GetDevices() {
		i, device := i, device
		g.Go(func() error {
			result := provision.Run(ctx, device, req.GetInventoryServer(), findContainer(cm, device, "cros-dut"), findContainer(cm, device, "cros-provision"))
			provisionResults[i] = result.Out
			return result.Err
		})
	}
	err = g.Wait()
	// Read all generated results for the output.
	for _, result := range provisionResults {
		out.Responses = append(out.Responses, result)
	}
	return out, errors.Annotate(err, "inner run").Err()
}

func findContainer(cm *build_api.ContainerMetadata, device *api.CrosToolRunnerProvisionRequest_Device, name string) *build_api.ContainerImageInfo {
	// TODO: need update logic base on real examples.
	for _, c := range cm.GetContainers() {
		for n, i := range c.GetImages() {
			if n == name {
				log.Printf("Found %q image %s", name, i)
				return i
			}
		}
	}
	log.Printf("Image %q not found", name)
	return nil
}

func isEmptyEndPoint(i *lab_api.IpEndpoint) bool {
	return i == nil || i.GetAddress() == "" || i.GetPort() <= 0
}

// readProvisionRequest reads the jsonproto at path input request data.
func readProvisionRequest(p string) (*api.CrosToolRunnerProvisionRequest, error) {
	in := &api.CrosToolRunnerProvisionRequest{}
	r, err := os.Open(p)
	if err != nil {
		return nil, errors.Annotate(err, "read provision request %q", p).Err()
	}
	err = jsonpb.Unmarshal(r, in)
	return in, errors.Annotate(err, "read provision request %q", p).Err()
}

// readContainersMetadata reads the jsonproto at path containers metadata file.
func readContainersMetadata(p string) (*build_api.ContainerMetadata, error) {
	in := &build_api.ContainerMetadata{}
	r, err := os.Open(p)
	if err != nil {
		return nil, errors.Annotate(err, "read container metadata %q", p).Err()
	}
	err = jsonpb.Unmarshal(r, in)
	return in, errors.Annotate(err, "read container metadata %q", p).Err()
}

// saveOutput saves output data to the file.
func saveOutput(out *api.CrosToolRunnerProvisionResponse, outputPath string) error {
	if outputPath != "" && out != nil {
		dir := filepath.Dir(outputPath)
		// Create the directory if it doesn't exist.
		if err := os.MkdirAll(dir, 0777); err != nil {
			return errors.Annotate(err, "save output").Err()
		}
		f, err := os.Create(outputPath)
		if err != nil {
			return errors.Annotate(err, "save output").Err()
		}
		defer f.Close()
		marshaler := jsonpb.Marshaler{}
		if err := marshaler.Marshal(f, out); err != nil {
			return errors.Annotate(err, "save output").Err()
		}
		if err := f.Close(); err != nil {
			return errors.Annotate(err, "save output").Err()
		}
	}
	return nil
}

func printOutput(out *api.CrosToolRunnerProvisionResponse, a subcommands.Application) {
	if out != nil {
		s, err := json.MarshalIndent(out, "", "\t")
		if err != nil {
			log.Printf("Output: fail to print info. Error: %s", err)
		} else {
			log.Println("Output:")
			fmt.Fprintf(a.GetOut(), "%s\n", s)
		}
	}
}

// readLocalAddress read local IP of the host.
func readLocalAddress() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", errors.Annotate(err, "read local address").Err()
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return "", errors.Annotate(err, "read local address").Err()
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			// TODO(otabek): Add option to work with IPv6 if we switched to it.
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.Reason("read local address: fail to find").Err()
}

func useSystemAuth(ctx context.Context, authFlags *authcli.Flags) (context.Context, error) {
	authOpts, err := authFlags.Options()
	if err != nil {
		return nil, errors.Annotate(err, "switching to system auth").Err()
	}

	authCtx, err := lucictx.SwitchLocalAccount(ctx, "system")
	if err == nil {
		// If there's a system account use it (the case of running on Swarming).
		// Otherwise default to user credentials (the local development case).
		authOpts.Method = auth.LUCIContextMethod
		return authCtx, nil
	}
	log.Printf("System account not found, err %s.\nFalling back to user credentials for auth.\n", err)
	return ctx, nil
}
