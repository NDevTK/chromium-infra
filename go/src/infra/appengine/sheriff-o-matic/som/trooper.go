// Copyright 2016 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Package som implements HTTP server that handles requests to default module.
package som

import (
	"encoding/json"
	"time"

	"infra/monitoring/messages"

	"golang.org/x/net/context"

	"github.com/luci/gae/service/datastore"
	"github.com/luci/luci-go/common/api/swarming/swarming/v1"
	"github.com/luci/luci-go/common/logging"
)

var (
	swarmingBasePath = "https://chromium-swarm.appspot.com/_ah/api/swarming/v1/"
)

// TrooperAlert ... Extended alert struct type for use in the trooper tab.
type TrooperAlert struct {
	messages.Alert
	Tree string `json:"tree"`
}

type swarmingAlerts struct {
	Dead        []*swarming.SwarmingRpcsBotInfo `json:"dead"`
	Quarantined []*swarming.SwarmingRpcsBotInfo `json:"quarantined"`
	Error       []string                        `json:"errors"`
}

func getTrooperAlerts(c context.Context) ([]byte, error) {
	swarmAlerts := make(chan *swarmingAlerts)
	go func() {
		swarmAlerts <- getSwarmingAlerts(c)
	}()

	q := datastore.NewQuery("Tree")
	trees := []*Tree{}
	datastore.GetAll(c, q, &trees)

	result := make(map[string]interface{})
	alerts := []*TrooperAlert{}

	// Assume that none of the timestamps will be from after right now.
	timestamp := messages.EpochTime(time.Now().Unix())

	for _, t := range trees {
		q := datastore.NewQuery("AlertsJSON")
		q = q.Ancestor(datastore.MakeKey(c, "Tree", t.Name))
		q = q.Order("-Date")
		q = q.Limit(1)

		alertsJSON := []*AlertsJSON{}
		err := datastore.GetAll(c, q, &alertsJSON)
		if err != nil {
			return nil, err
		}

		if len(alertsJSON) > 0 {
			data := alertsJSON[0].Contents

			result["date"] = alertsJSON[0].Date

			alertsSummary := &messages.AlertsSummary{}

			err = json.Unmarshal(data, alertsSummary)
			if err != nil {
				return nil, err
			}

			newTime := alertsSummary.Timestamp
			if newTime > 0 && newTime < timestamp {
				timestamp = newTime
			}
			result["revision_summaries"] = alertsSummary.RevisionSummaries

			for _, a := range alertsSummary.Alerts {
				if a.Type == messages.AlertInfraFailure ||
					a.Type == messages.AlertOfflineBuilder {
					newAlert := &TrooperAlert{a, t.Name}
					alerts = append(alerts, newAlert)
				}
			}
		}
	}

	result["timestamp"] = timestamp
	result["alerts"] = alerts
	result["swarming"] = <-swarmAlerts

	out, err := json.Marshal(result)

	if err != nil {
		return nil, err
	}

	return out, nil
}

func getSwarmingAlerts(c context.Context) *swarmingAlerts {
	// TODO(seanmmccullough): Cache these responses for a few minutes.
	ret := &swarmingAlerts{}
	c, _ = context.WithTimeout(c, 60*time.Second)
	oauthClient, err := getOAuthClient(c)
	if err != nil {
		logging.Errorf(c, "getting oauth client: %v", err)
		ret.Error = append(ret.Error, err.Error())
		return ret
	}

	swarmingService, err := swarming.New(oauthClient)
	if err != nil {
		logging.Errorf(c, "getting swarming client: %v", err)
		ret.Error = append(ret.Error, err.Error())
		return ret
	}
	swarmingService.BasePath = swarmingBasePath

	botCh := make(chan *swarmingAlerts)

	go func() {
		deadBots, err := swarmingService.Bots.List().IsDead("TRUE").Do()
		sa := &swarmingAlerts{}
		if err != nil {
			logging.Errorf(c, "getting dead bots: %v", err)
			sa.Error = append(sa.Error, err.Error())
		}
		if deadBots != nil {
			sa.Dead = filterBots(deadBots.Items)
		}

		botCh <- sa
	}()

	go func() {
		quarantinedBots, err := swarmingService.Bots.List().Quarantined("TRUE").Do()
		sa := &swarmingAlerts{}
		if err != nil {
			logging.Errorf(c, "getting quarantined bots: %v", err.Error())
			sa.Error = append(sa.Error, err.Error())
		}
		if quarantinedBots != nil {
			sa.Quarantined = filterBots(quarantinedBots.Items)
		}
		botCh <- sa
	}()

	for i := 0; i < 2; i++ {
		bots := <-botCh
		if bots.Dead != nil {
			ret.Dead = bots.Dead
		}
		if bots.Quarantined != nil {
			ret.Quarantined = bots.Quarantined
		}
		ret.Error = append(ret.Error, bots.Error...)
	}

	return ret
}

// filterBots removes the State field, which is escaped JSON which we do not
// parse and can get arbitrarily large.
func filterBots(bots []*swarming.SwarmingRpcsBotInfo) []*swarming.SwarmingRpcsBotInfo {
	for _, bot := range bots {
		bot.State = ""
	}
	return bots
}
