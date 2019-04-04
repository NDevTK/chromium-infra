// Copyright 2019 The Chromium OS Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Package worker implements a constructor for skylab_swarming_worker
// commands.  This package is intended to be used by package that need
// to construct a command line for running skylab_swarming_worker.
package worker

import (
	"encoding/json"
	"fmt"
	"strings"
)

// DefaultPath is the default path for skylab_swarming_worker.
const DefaultPath = "/opt/infra-tools/skylab_swarming_worker"

// Command is a constructor for skylab_swarming_worker commands.
type Command struct {
	// Path to skylab_swarming_worker.  The default is DefaultPath.
	Path string

	// TaskName is required.
	TaskName string
	// LogDogAnnotationURL can be set automatically with Env.
	LogDogAnnotationURL string
	ForceFresh          bool
	ClientTest          bool
	ProvisionLabels     []string
	Keyvals             map[string]string
	TestArgs            string
}

// Args returns the arg strings for running the command.
func (c *Command) Args() []string {
	var args []string
	if c.Path != "" {
		args = append(args, c.Path)
	} else {
		args = append(args, DefaultPath)
	}
	if c.TaskName != "" {
		args = append(args, "-task-name", c.TaskName)
	}
	if c.LogDogAnnotationURL != "" {
		args = append(args, "-logdog-annotation-url", c.LogDogAnnotationURL)
	}
	if c.ForceFresh {
		args = append(args, "-force-fresh")
	}
	if c.ClientTest {
		args = append(args, "-client-test")
	}
	if len(c.ProvisionLabels) > 0 {
		args = append(args, "-provision-labels", strings.Join(c.ProvisionLabels, ","))
	}
	if c.Keyvals != nil {
		b, err := json.Marshal(c.Keyvals)
		// Marshal map[string]string should never error.
		if err != nil {
			panic(err)
		}
		args = append(args, "-keyvals", string(b))
	}
	if c.TestArgs != "" {
		args = append(args, "-test-args", c.TestArgs)
	}
	return args
}

// Config configures the command with the given options.  This method
// returns the command receiver to allow method chaining.
func (c *Command) Config(o ...Option) *Command {
	for _, o := range o {
		o(c)
	}
	return c
}

// Option is used to help construct a worker command.
type Option func(*Command)

// Environment defines a Skylab environment (e.g., dev vs prod) for
// configuring a worker command.
type Environment interface {
	LUCIProject() string
	LogDogHost() string
	GenerateLogPrefix() string
}

// Env returns an option used to configure a worker command for a
// Skylab environment (e.g., dev vs prod).
func Env(e Environment) Option {
	return func(c *Command) {
		c.LogDogAnnotationURL = GenerateLogDogURL(e)
	}
}

// GenerateLogDogURL generates a LogDog annotation URL that is
// suitable for a worker command.
func GenerateLogDogURL(e Environment) string {
	u := logDogURL{
		Host:    e.LogDogHost(),
		Project: e.LUCIProject(),
		Prefix:  e.GenerateLogPrefix(),
	}
	return u.String()
}

// logDogURL is a constructor for LogDog annotation URLs.
type logDogURL struct {
	Host    string
	Project string
	Prefix  string
}

func (u logDogURL) String() string {
	return fmt.Sprintf("logdog://%s/%s/%s/+/annotations", u.Host, u.Project, u.Prefix)
}
