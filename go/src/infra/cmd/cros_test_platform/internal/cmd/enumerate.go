// Copyright 2098 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package cmd

import (
	"context"
	"fmt"

	"io/ioutil"
	"os"
	"path/filepath"

	"infra/cmd/cros_test_platform/internal/autotest/artifacts"
	"infra/cmd/cros_test_platform/internal/autotest/testspec"
	"infra/cmd/cros_test_platform/internal/enumeration"
	"infra/cmd/cros_test_platform/internal/site"

	"github.com/maruel/subcommands"
	"go.chromium.org/chromiumos/infra/proto/go/chromite/api"
	"go.chromium.org/chromiumos/infra/proto/go/test_platform/steps"
	"go.chromium.org/luci/auth/client/authcli"
	"go.chromium.org/luci/common/cli"
	"go.chromium.org/luci/common/errors"
	"go.chromium.org/luci/common/gcloud/gs"
)

// Enumerate is the `enumerate` subcommand implementation.
var Enumerate = &subcommands.Command{
	UsageLine: "enumerate -input_json /path/to/input.json -output_json /path/to/output.json",
	ShortDesc: "Enumerate tasks to execute for a request.",
	LongDesc: `Enumerate tasks to execute for a request.

Step input and output is JSON encoded protobuf defined at
https://chromium.googlesource.com/chromiumos/infra/proto/+/master/src/test_platform/steps/enumeration.proto`,
	CommandRun: func() subcommands.CommandRun {
		c := &enumerateRun{}
		c.authFlags.Register(&c.Flags, site.DefaultAuthOptions)
		c.Flags.StringVar(&c.inputPath, "input_json", "", "Path that contains JSON encoded test_platform.steps.EnumerationRequest")
		c.Flags.StringVar(&c.outputPath, "output_json", "", "Path where JSON encoded test_platform.steps.EnumerationResponse should be written.")
		return c
	},
}

type enumerateRun struct {
	subcommands.CommandRunBase
	authFlags authcli.Flags

	inputPath  string
	outputPath string
}

func (c *enumerateRun) Run(a subcommands.Application, args []string, env subcommands.Env) int {
	if err := c.innerRun(a, args, env); err != nil {
		fmt.Fprintf(a.GetErr(), "%s\n", err)
		return 1
	}
	return 0
}

func (c *enumerateRun) innerRun(a subcommands.Application, args []string, env subcommands.Env) error {
	if err := c.processCLIArgs(args); err != nil {
		return err
	}
	request, err := c.readRequest()
	if err != nil {
		return err
	}

	workspace, err := ioutil.TempDir("", "enumerate")
	if err != nil {
		return err
	}
	defer func() {
		os.RemoveAll(workspace)
	}()

	ctx := cli.GetContext(a, c, env)
	gsDir := gs.Path(request.GetMetadata().GetTestMetadataUrl())
	lp, err := c.downloadArtifacts(ctx, gsDir, workspace)
	if err != nil {
		return err
	}
	tm, err := computeMetadata(lp, workspace)
	if err != nil {
		return err
	}
	if err := c.writeResponse(c.enumerate(tm, request)); err != nil {
		return err
	}
	return nil
}

func (c *enumerateRun) processCLIArgs(args []string) error {
	if len(args) > 0 {
		return errors.Reason("have %d positional args, want 0", len(args)).Err()
	}
	if c.inputPath == "" {
		return errors.Reason("-input_json not specified").Err()
	}
	if c.outputPath == "" {
		return errors.Reason("-output_json not specified").Err()
	}
	return nil
}

func (c *enumerateRun) readRequest() (*steps.EnumerationRequest, error) {
	r, err := os.Open(c.inputPath)
	if err != nil {
		return nil, errors.Annotate(err, "read request").Err()
	}
	defer r.Close()

	request := steps.EnumerationRequest{}
	if err := unmarshaller.Unmarshal(r, &request); err != nil {
		return nil, errors.Annotate(err, "read request").Err()
	}
	return &request, nil
}

func (c *enumerateRun) downloadArtifacts(ctx context.Context, gsDir gs.Path, workspace string) (artifacts.LocalPaths, error) {
	outDir := filepath.Join(workspace, "artifacts")
	if err := os.Mkdir(outDir, 0750); err != nil {
		return artifacts.LocalPaths{}, errors.Annotate(err, "download artifacts").Err()
	}
	client, err := c.newGSClient(ctx)
	if err != nil {
		return artifacts.LocalPaths{}, errors.Annotate(err, "download artifacts").Err()
	}
	lp, err := artifacts.DownloadFromGoogleStorage(client, gsDir, outDir)
	if err != nil {
		return artifacts.LocalPaths{}, errors.Annotate(err, "download artifacts").Err()
	}
	return lp, err
}

func (c *enumerateRun) newGSClient(ctx context.Context) (gs.Client, error) {
	t, err := newAuthenticatedTransport(ctx, &c.authFlags)
	if err != nil {
		return nil, errors.Annotate(err, "create GS client").Err()
	}
	return gs.NewProdClient(ctx, t)
}

func (c *enumerateRun) enumerate(tm *api.TestMetadataResponse, request *steps.EnumerationRequest) []*api.AutotestTest {
	var ts []*api.AutotestTest
	ts = append(ts, enumeration.GetForTests(tm.GetAutotest(), request.GetTests())...)
	ts = append(ts, enumeration.GetForSuites(tm.GetAutotest(), request.GetSuites())...)
	return ts
}

func (c *enumerateRun) writeResponse(tests []*api.AutotestTest) error {
	w, err := os.Create(c.outputPath)
	if err != nil {
		return errors.Annotate(err, "write response").Err()
	}
	defer w.Close()
	resp := steps.EnumerationResponse{AutotestTests: tests}
	if err := marshaller.Marshal(w, &resp); err != nil {
		return errors.Annotate(err, "write response").Err()
	}
	return nil
}

func computeMetadata(localPaths artifacts.LocalPaths, workspace string) (*api.TestMetadataResponse, error) {
	extracted := filepath.Join(workspace, "extracted")
	if err := os.Mkdir(extracted, 0750); err != nil {
		return nil, errors.Annotate(err, "compute metadata").Err()
	}
	if err := artifacts.ExtractControlFiles(localPaths, extracted); err != nil {
		return nil, errors.Annotate(err, "compute metadata").Err()
	}
	return testspec.Get(extracted)
}
