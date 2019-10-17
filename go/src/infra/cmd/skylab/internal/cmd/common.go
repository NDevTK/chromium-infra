// Copyright 2018 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package cmd

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/maruel/subcommands"

	"go.chromium.org/luci/auth"
	"go.chromium.org/luci/auth/client/authcli"
	"go.chromium.org/luci/common/errors"
	lflag "go.chromium.org/luci/common/flag"

	"infra/cmd/skylab/internal/site"
	"infra/libs/skylab/common/errctx"
)

const progName = "skylab"

var defaultTaskPriority = 140

var jsonPBMarshaller = &jsonpb.Marshaler{
	EmitDefaults: true,
}

var jsonPBUnmarshaller = jsonpb.Unmarshaler{
	AllowUnknownFields: true,
}

type commonFlags struct {
	debug bool
}

func (f *commonFlags) Register(fl *flag.FlagSet) {
	fl.BoolVar(&f.debug, "debug", false, "Enable debug output.")
}

func (f commonFlags) DebugLogger(a subcommands.Application) *log.Logger {
	out := ioutil.Discard
	if f.debug {
		out = a.GetErr()
	}
	return log.New(out, progName, log.LstdFlags|log.Lshortfile)
}

type envFlags struct {
	dev bool
}

func (f *envFlags) Register(fl *flag.FlagSet) {
	fl.BoolVar(&f.dev, "dev", false, "Run in dev environment.")
}

func (f envFlags) Env() site.Environment {
	if f.dev {
		return site.Dev
	}
	return site.Prod
}

type removalReason struct {
	bug     string
	comment string
	expire  time.Time
}

func (rr *removalReason) Register(f *flag.FlagSet) {
	f.StringVar(&rr.bug, "bug", "", "Bug link for why DUT is being removed.  Required.")
	f.StringVar(&rr.comment, "comment", "", "Short comment about why DUT is being removed.")
	f.Var(lflag.RelativeTime{T: &rr.expire}, "expires-in", "Expire removal reason in `days`.")
}

// newHTTPClient returns an HTTP client with authentication set up.
func newHTTPClient(ctx context.Context, f *authcli.Flags) (*http.Client, error) {
	o, err := f.Options()
	if err != nil {
		return nil, errors.Annotate(err, "failed to get auth options").Err()
	}
	a := auth.NewAuthenticator(ctx, auth.OptionalLogin, o)
	c, err := a.Client()
	if err != nil {
		return nil, errors.Annotate(err, "failed to create HTTP client").Err()
	}
	return c, nil
}

// UserErrorReporter reports a detailed error message to the user.
//
// PrintError() uses a UserErrorReporter to print multi-line user error details
// along with the actual error.
type UserErrorReporter interface {
	// Report a user-friendly error through w.
	ReportUserError(w io.Writer)
}

// PrintError reports errors back to the user.
//
// Detailed error information is printed if err is a UserErrorReporter.
func PrintError(w io.Writer, err error) {
	if u, ok := err.(UserErrorReporter); ok {
		u.ReportUserError(w)
	} else {
		fmt.Fprintf(w, "%s: %s\n", progName, err)
	}
}

// NewUsageError creates a new error that also reports flags usage error
// details.
func NewUsageError(flags flag.FlagSet, format string, a ...interface{}) error {
	return &usageError{
		error: fmt.Errorf(format, a...),
		flags: flags,
	}
}

type usageError struct {
	error
	flags flag.FlagSet
}

func (e *usageError) ReportUserError(w io.Writer) {
	fmt.Fprintf(w, "%s\n\nUsage:\n\n", e.error)
	e.flags.Usage()
}

func maybeWithTimeout(ctx context.Context, timeoutMins int) (context.Context, func(error)) {
	if timeoutMins >= 0 {
		return errctx.WithTimeout(ctx, time.Duration(timeoutMins)*time.Minute,
			fmt.Errorf("timed out after %d minutes while waiting for task(s) to complete", timeoutMins))
	}
	return errctx.WithCancel(ctx)
}
