// Copyright 2020 The LUCI Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/maruel/subcommands"
	"golang.org/x/sync/errgroup"
	"golang.org/x/time/rate"
	"google.golang.org/api/option"

	"go.chromium.org/luci/auth"
	"go.chromium.org/luci/common/api/gerrit"
	"go.chromium.org/luci/common/cli"
	"go.chromium.org/luci/common/data/caching/lru"
	"go.chromium.org/luci/common/data/text"
	"go.chromium.org/luci/common/errors"
	luciflag "go.chromium.org/luci/common/flag"
	gerritpb "go.chromium.org/luci/common/proto/gerrit"

	"infra/rts/presubmit/eval/history"
	evalpb "infra/rts/presubmit/eval/proto"
)

func cmdPresubmitHistory(authOpt *auth.Options) *subcommands.Command {
	return &subcommands.Command{
		UsageLine: `presubmit-history`,
		ShortDesc: "fetch presubmit history",
		LongDesc: text.Doc(`
			Fetch presubmit history, suitable for RTS evaluation.
		`),
		CommandRun: func() subcommands.CommandRun {
			r := &presubmitHistoryRun{authOpt: authOpt}
			r.Flags.StringVar(&r.out, "out", "", "Path to the output file")
			r.Flags.Var(luciflag.Date(&r.startTime), "from", "Fetch results starting from this date; format: 2020-01-02")
			r.Flags.Var(luciflag.Date(&r.endTime), "to", "Fetch results until this date; format: 2020-01-02")
			r.Flags.Float64Var(&r.durationDataFrac, "duration-data-frac", 0.001, "Fraction of duration data to fetch")
			r.Flags.DurationVar(&r.minDuration, "min-duration", time.Second, "Minimum duration to fetch")
			r.Flags.StringVar(&r.builderRegex, "builder", ".*", "A regular expression for builder. Implicitly wrapped with ^ and $.")
			r.Flags.StringVar(&r.testIDRegex, "test", ".*", "A regular expression for test. Implicitly wrapped with ^ and $.")
			r.Flags.IntVar(&r.minCLFlakes, "min-cl-flakes", 5, text.Doc(`
				In order to conlude that a test variant is flaky and exclude it from analysis,
				it must have mixed results in <min-cl-flakes> unique CLs.
			`))
			return r
		},
	}
}

type presubmitHistoryRun struct {
	baseCommandRun
	out              string
	startTime        time.Time
	endTime          time.Time
	durationDataFrac float64
	builderRegex     string
	testIDRegex      string
	minDuration      time.Duration
	minCLFlakes      int

	gerrit *gerritClient

	authenticator *auth.Authenticator
	authOpt       *auth.Options

	mu                    sync.Mutex
	w                     *history.Writer
	recordsWrote          int
	recordCountNextReport time.Time
}

func (r *presubmitHistoryRun) validate() error {
	switch {
	case r.out == "":
		return errors.New("-out is required")
	case r.startTime.IsZero():
		return errors.New("-from is required")
	case r.endTime.IsZero():
		return errors.New("-to is required")
	case r.endTime.Before(r.startTime):
		return errors.New("the -to date must not be before the -from date")
	default:
		return nil
	}
}

func (r *presubmitHistoryRun) Run(a subcommands.Application, args []string, env subcommands.Env) int {
	ctx := cli.GetContext(a, r, env)
	if len(args) != 0 {
		return r.done(errors.New("unexpected positional arguments"))
	}

	if err := r.validate(); err != nil {
		return r.done(err)
	}

	r.authenticator = auth.NewAuthenticator(ctx, auth.InteractiveLogin, *r.authOpt)

	var err error
	if r.gerrit, err = r.newGerritClient(); err != nil {
		return r.done(errors.Annotate(err, "failed to init Gerrit client").Err())
	}

	// Create the history file.
	if r.w, err = history.CreateFile(r.out); err != nil {
		return r.done(errors.Annotate(err, "failed to create the output file").Err())
	}
	defer r.w.Close()

	fmt.Printf("starting BigQuery queries...\n")

	eg, ctx := errgroup.WithContext(ctx)

	// Fetch the rejections.
	var rejections int
	eg.Go(func() error {
		err := r.rejections(ctx, func(frag *evalpb.RejectionFragment) error {
			if frag.Terminal {
				rejections++
			}
			return r.write(&evalpb.Record{
				Data: &evalpb.Record_RejectionFragment{RejectionFragment: frag},
			})
		})
		return errors.Annotate(err, "failed to fetch rejections").Err()
	})

	// Fetch test durations.
	var durations int
	if r.durationDataFrac > 0 {
		eg.Go(func() error {
			err := r.durations(ctx, func(td *evalpb.TestDuration) error {
				durations++
				return r.write(&evalpb.Record{
					Data: &evalpb.Record_TestDuration{TestDuration: td},
				})
			})
			return errors.Annotate(err, "failed to fetch test durations").Err()
		})
	}

	if err = eg.Wait(); err != nil {
		return r.done(err)
	}

	r.mu.Lock()
	fmt.Printf("total: %d rejections, %d durations, %d records\n", rejections, durations, r.recordsWrote)
	r.mu.Unlock()

	return r.done(r.w.Close())
}

// write write rec to the output file.
// Occasionally prints out progress.
func (r *presubmitHistoryRun) write(rec *evalpb.Record) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := r.w.Write(rec); err != nil {
		return err
	}

	// Occasionally print progress.
	r.recordsWrote++
	now := time.Now()
	if r.recordCountNextReport.Before(now) {
		if !r.recordCountNextReport.IsZero() {
			fmt.Printf("wrote %d records\n", r.recordsWrote)
		}
		r.recordCountNextReport = now.Add(time.Second)
	}
	return nil
}

func (r *presubmitHistoryRun) bqQuery(ctx context.Context, sql string) (*bigquery.Query, error) {
	http, err := r.authenticator.Client()
	if err != nil {
		return nil, err
	}
	client, err := bigquery.NewClient(ctx, "chrome-trooper-analytics", option.WithHTTPClient(http))
	if err != nil {
		return nil, err
	}

	prepRe := func(rgx string) string {
		if rgx == "" || rgx == ".*" {
			return ""
		}
		return fmt.Sprintf("^(%s)$", rgx)
	}

	q := client.Query(sql)
	q.Parameters = []bigquery.QueryParameter{
		{Name: "startTime", Value: r.startTime},
		{Name: "endTime", Value: r.endTime},
		{Name: "builder_regexp", Value: prepRe(r.builderRegex)},
		{Name: "test_id_regexp", Value: prepRe(r.testIDRegex)},
	}
	return q, nil
}

// populateChangedFiles populates ps.ChangedFiles.
// TODO(nodir): delete this function, gerrit.go and cache.go in February 2021,
// when we have enough BigQuery data with gerrit info.
func (r *presubmitHistoryRun) populateChangedFiles(ctx context.Context, ps *evalpb.GerritPatchset) error {
	changedFiles, err := r.gerrit.ChangedFiles(ctx, ps)
	if err != nil {
		return err
	}

	repo := fmt.Sprintf("https://%s/%s", ps.Change.Host, strings.TrimSuffix(ps.Change.Project, ".git"))

	ps.ChangedFiles = make([]*evalpb.SourceFile, len(changedFiles))
	for i, path := range changedFiles {
		ps.ChangedFiles[i] = &evalpb.SourceFile{
			Repo: repo,
			Path: "//" + path,
		}
	}
	return nil
}

// newGerritClient creates a new gitiles client. Does not mutate r.
func (r *presubmitHistoryRun) newGerritClient() (*gerritClient, error) {
	transport, err := r.authenticator.Transport()
	if err != nil {
		return nil, err
	}

	ucd, err := os.UserCacheDir()
	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{Transport: transport}
	return &gerritClient{
		listFilesRPC: func(ctx context.Context, host string, req *gerritpb.ListFilesRequest) (*gerritpb.ListFilesResponse, error) {
			client, err := gerrit.NewRESTClient(httpClient, host, true)
			if err != nil {
				return nil, errors.Annotate(err, "failed to create a Gerrit client").Err()
			}
			return client.ListFiles(ctx, req)
		},
		fileListCache: cache{
			dir:       filepath.Join(ucd, "chrome-rts", "gerrit-changed-files"),
			memory:    lru.New(1024),
			valueType: reflect.TypeOf(changedFiles{}),
		},
		limiter: rate.NewLimiter(10, 1),
	}, nil
}
