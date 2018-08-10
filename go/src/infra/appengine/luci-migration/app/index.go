// Copyright 2017 The LUCI Authors.
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

package app

import (
	"sort"

	"golang.org/x/net/context"

	"go.chromium.org/gae/service/datastore"
	"go.chromium.org/luci/common/errors"
	"go.chromium.org/luci/server/auth"
	"go.chromium.org/luci/server/router"
	"go.chromium.org/luci/server/templates"

	"infra/appengine/luci-migration/config"
	"infra/appengine/luci-migration/storage"
)

type indexViewModel struct {
	Masters []*indexMasterViewModel
}

type indexMasterViewModel struct {
	Name string

	// Counts of builders in combinations:
	// i: 0 -> not prod on LUCI; 1 -> prod on LUCI
	// j: 0 -> not found on Buildbot; 1 -> found on Buildbot
	LUCIBuildbotCounts     [2][2]int
	WAIBuilderCount        int
	WAIBuilderPercent      int
	MigratedBuilderCount   int
	MigratedBuilderPercent int
	TotalBuilderCount      int
}

func handleIndexPage(c *router.Context) error {
	viewModel, err := indexPage(c.Context)
	if err != nil {
		return err
	}
	templates.MustRender(c.Context, c.Writer, "pages/index.html", templates.Args{"Model": viewModel})
	return nil
}

func indexPage(c context.Context) (*indexViewModel, error) {
	cfg := config.Get(c)
	cfgMap := make(map[string]*config.Master, len(cfg.Masters))
	for _, m := range cfg.Masters {
		cfgMap[m.Name] = m
	}

	hasInternalAccess, err := auth.IsMember(c, accessGroup)
	if err != nil {
		return nil, errors.Annotate(err, "could not check group membership").Err()
	}

	masters := map[string]*indexMasterViewModel{}
	masterNames := []string{}
	// Note: may have to cache this if we have a lot of builders.
	q := datastore.NewQuery(storage.BuilderKind)
	err = datastore.Run(c, q, func(b *storage.Builder) {
		m := masters[b.ID.Master]
		if m == nil {
			switch mCfg := cfgMap[b.ID.Master]; {
			case mCfg == nil:
				return
			case !mCfg.Public && !hasInternalAccess:
				return
			}

			m = &indexMasterViewModel{Name: b.ID.Master}
			masters[b.ID.Master] = m
			masterNames = append(masterNames, m.Name)
		}

		m.TotalBuilderCount++

		i, j := 0, 0
		if b.LUCIIsProd {
			i = 1
		}
		if !b.NotOnBuildbot {
			j = 1
		}
		m.LUCIBuildbotCounts[i][j]++

		switch b.Migration.Status {
		case storage.StatusMigrated:
			m.MigratedBuilderCount++
			fallthrough
		case storage.StatusLUCIWAI:
			m.WAIBuilderCount++
		}
	})
	if err != nil {
		return nil, err
	}

	sort.Strings(masterNames)
	model := &indexViewModel{Masters: make([]*indexMasterViewModel, len(masterNames))}
	for i, name := range masterNames {
		m := masters[name]
		if m.TotalBuilderCount > 0 {
			m.WAIBuilderPercent = 100 * m.WAIBuilderCount / m.TotalBuilderCount
			m.MigratedBuilderPercent = 100 * m.MigratedBuilderCount / m.TotalBuilderCount
		}
		model.Masters[i] = m
	}

	return model, nil
}
