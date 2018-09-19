// Copyright 2017 The LUCI Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"net/url"

	"go.chromium.org/luci/common/system/environ"

	"infra/tools/kitchen/cookflags"
)

var cookModeSelector = map[cookflags.CookMode]cookMode{
	cookflags.CookSwarming: swarmingCookMode{},
	cookflags.CookBuildBot: buildBotCookMode{},
}

const (
	// PropertyBotID must be added by cookMode.addProperties.
	PropertyBotID = "bot_id"
)

// cookMode integrates environment-specific behaviors into Kitchen's "cook"
// command.
type cookMode interface {
	needsIOKeepAlive() bool
	alwaysForwardAnnotations() bool

	// addProperties adds builtin properties. Must add PropertyBotId.
	addProperties(props map[string]interface{}, env environ.Env) error

	shouldEmitLogDogLinks() bool
	addLogDogGlobalTags(tags map[string]string, props map[string]interface{}, env environ.Env) error

	allowCustomGitAuth() bool
	allowDevShell() bool
	allowDockerAuth() bool
	allowFirebaseAuth() bool
}

// readSwarmingEnv reads relevent data out of the environment.
func readSwarmingEnv(env environ.Env) (botID string, err error) {
	var ok bool
	botID, ok = env.Get("SWARMING_BOT_ID")
	if !ok {
		err = inputError("no Swarming bot ID in $SWARMING_BOT_ID")
		return
	}

	return
}

type swarmingCookMode struct{}

func (m swarmingCookMode) needsIOKeepAlive() bool         { return false }
func (m swarmingCookMode) alwaysForwardAnnotations() bool { return false }

func (m swarmingCookMode) addProperties(props map[string]interface{}, env environ.Env) error {
	botID, err := readSwarmingEnv(env)
	if err != nil {
		return err
	}
	props[PropertyBotID] = botID
	return nil
}
func (m swarmingCookMode) shouldEmitLogDogLinks() bool { return false }
func (m swarmingCookMode) addLogDogGlobalTags(tags map[string]string, props map[string]interface{},
	env environ.Env) error {

	// SWARMING_SERVER is the full URL: https://example.com
	// We want just the hostname.
	if v, ok := env.Get("SWARMING_SERVER"); ok {
		if u, err := url.Parse(v); err == nil && u.Host != "" {
			tags["swarming.host"] = u.Host
		}
	}
	if v, ok := env.Get("SWARMING_TASK_ID"); ok {
		tags["swarming.run_id"] = v
	}
	if v, ok := env.Get("SWARMING_BOT_ID"); ok {
		tags["bot_id"] = v
	}

	return nil
}

func (m swarmingCookMode) allowCustomGitAuth() bool { return true }
func (m swarmingCookMode) allowDevShell() bool      { return true }
func (m swarmingCookMode) allowDockerAuth() bool    { return true }
func (m swarmingCookMode) allowFirebaseAuth() bool  { return true }

type buildBotCookMode struct{}

func (m buildBotCookMode) needsIOKeepAlive() bool         { return true }
func (m buildBotCookMode) alwaysForwardAnnotations() bool { return true }

func (m buildBotCookMode) addProperties(props map[string]interface{}, env environ.Env) error {
	botID, ok := env.Get("BUILDBOT_SLAVENAME")
	if !ok {
		return inputError("no slave name in $BUILDBOT_SLAVENAME")
	}
	props[PropertyBotID] = botID
	return nil
}

func (m buildBotCookMode) shouldEmitLogDogLinks() bool { return true }
func (m buildBotCookMode) addLogDogGlobalTags(tags map[string]string, props map[string]interface{},
	env environ.Env) error {

	if v, ok := props["mastername"].(string); ok && v != "" {
		tags["buildbot.master"] = v
	}
	if v, ok := props["buildername"].(string); ok && v != "" {
		tags["buildbot.builder"] = v
	}
	if v, ok := props["buildnumber"].(json.Number); ok && v != "" {
		tags["buildbot.buildnumber"] = string(v)
	}

	if v, ok := props["slavename"].(string); ok && v != "" {
		tags["bot_id"] = v
	}

	return nil
}

func (m buildBotCookMode) allowCustomGitAuth() bool { return false }
func (m buildBotCookMode) allowDevShell() bool      { return false }
func (m buildBotCookMode) allowDockerAuth() bool    { return false }
func (m buildBotCookMode) allowFirebaseAuth() bool  { return false }
