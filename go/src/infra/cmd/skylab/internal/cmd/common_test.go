// Copyright 2019 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package cmd

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"go.chromium.org/luci/common/clock"
	"go.chromium.org/luci/common/clock/testclock"

	"google.golang.org/api/googleapi"
)

func TestSwarmingCallWithRetries_TransientFailure(t *testing.T) {
	ctx, testClock := testclock.UseTime(context.Background(), time.Now())
	testClock.SetTimerCallback(func(time.Duration, clock.Timer) {
		testClock.Add(1 * time.Second)
	})
	count := 0
	f := func() error {
		defer func() { count++ }()
		if count == 0 {
			return &googleapi.Error{
				Code: http.StatusInternalServerError, // 500
			}
		}
		return nil
	}
	err := swarmingCallWithRetries(ctx, f)
	if err != nil {
		t.Fatalf("call error actual != expected, %v != %v", err, nil)
	}
	if count != 2 {
		t.Fatalf("try count actual != expected, %d != %d", count, 1)
	}
}

func TestSwarmingCallWithRetries_NontransientFailure(t *testing.T) {
	count := 0
	f := func() error {
		count++
		return errors.New("foo")
	}
	err := swarmingCallWithRetries(context.Background(), f)
	if err == nil {
		t.Fatalf("call error unexpectedly nil")
	}
	if count != 1 {
		t.Fatalf("try count actual != expected, %d != %d", count, 1)
	}
}
