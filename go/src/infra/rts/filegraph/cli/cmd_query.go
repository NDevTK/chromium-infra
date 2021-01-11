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

package cli

import (
	"errors"
	"fmt"

	"github.com/maruel/subcommands"

	"go.chromium.org/luci/common/cli"
	"go.chromium.org/luci/common/data/text"

	"infra/rts/filegraph"
)

var cmdQuery = &subcommands.Command{
	UsageLine: `query [flags] SOURCE_FILE [SOURCE_FILE...]`,
	ShortDesc: "print graph files in the distance-ascending order",
	LongDesc: text.Doc(`
		Print graph files in the distance-ascending order from SOURCE_FILEs.

		Each output line has format "<distance> <filename>",
		where the filename is forward-slash-separated and has "//" prefix.
		Example: "0.4 //foo/bar.cpp".

		All SOURCE_FILEs must be in the same git repository.
		Does not print unreachable files.
	`),
	CommandRun: func() subcommands.CommandRun {
		r := &queryRun{}
		r.git.RegisterFlags(&r.Flags)
		return r
	},
}

type queryRun struct {
	baseCommandRun
	git gitGraph
}

func (r *queryRun) Run(a subcommands.Application, args []string, env subcommands.Env) int {
	ctx := cli.GetContext(a, r, env)
	if err := r.git.Validate(); err != nil {
		return r.done(err)
	}
	if len(args) == 0 {
		return r.done(errors.New("expected filenames as positional arguments"))
	}

	var err error
	if r.git.q.Sources, err = r.git.loadSyncedNodes(ctx, args...); err != nil {
		return r.done(err)
	}

	r.git.q.Run(func(sp *filegraph.ShortestPath) bool {
		fmt.Printf("%.2f %s\n", sp.Distance, sp.Node.Name())
		return ctx.Err() == nil
	})
	return 0
}
