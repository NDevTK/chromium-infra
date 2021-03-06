// Copyright 2021 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"google.golang.org/genproto/protobuf/field_mask"

	"go.chromium.org/luci/common/logging"
	"go.chromium.org/luci/server/auth"
	"go.chromium.org/luci/server/router"
	"go.chromium.org/luci/server/span"

	"infra/appengine/weetbix/internal/clustering/rules"
	"infra/appengine/weetbix/internal/clustering/rules/lang"
)

// ListRules serves a GET request for
// /api/projects/:project/rules.
func (h *Handlers) ListRules(ctx *router.Context) {
	transctx, cancel := span.ReadOnlyTransaction(ctx.Context)
	defer cancel()

	projectID, ok := obtainProjectOrError(ctx)
	if !ok {
		return
	}
	rs, err := rules.ReadActive(transctx, projectID)
	if err != nil {
		logging.Errorf(ctx.Context, "Reading rules: %s", err)
		http.Error(ctx.Writer, "Internal server error.", http.StatusInternalServerError)
		return
	}

	respondWithJSON(ctx, rs)
}

// GetRule serves a GET request for
// /api/projects/:project/rules/:id.
func (h *Handlers) GetRule(ctx *router.Context) {
	projectID, ok := obtainProjectOrError(ctx)
	if !ok {
		return
	}
	ruleID := ctx.Params.ByName("id")
	if !rules.RuleIDRe.MatchString(ruleID) {
		http.Error(ctx.Writer, "Please supply a valid rule ID.", http.StatusBadRequest)
		return
	}
	txn, cancel := span.ReadOnlyTransaction(ctx.Context)
	defer cancel()

	rule, err := rules.Read(txn, projectID, ruleID)
	if err != nil {
		logging.Errorf(ctx.Context, "Reading rule %s: %s", ruleID, err)
		http.Error(ctx.Writer, "Internal server error.", http.StatusInternalServerError)
		return
	}

	ctx.Writer.Header().Add("ETag", ruleETag(rule))
	respondWithJSON(ctx, rule)
}

// Designed to conform to https://google.aip.dev/134.
type ruleUpdateRequest struct {
	Rule       *rules.FailureAssociationRule `json:"rule"`
	UpdateMask field_mask.FieldMask          `json:"updateMask"`
}

func ruleETag(rule *rules.FailureAssociationRule) string {
	// While this ETag is strong, GAE's NGINX proxy will sometimes
	// remove them or modify them to be Weak as it compresses content.
	// Marking the ETag as weak to start with addresses this.
	return fmt.Sprintf(`W/"%s"`, rule.LastUpdated.UTC().Format(time.RFC3339Nano))
}

// PatchRule serves a PATCH request for
// /api/projects/:project/rules/:id.
func (h *Handlers) PatchRule(ctx *router.Context) {
	projectID, ok := obtainProjectOrError(ctx)
	if !ok {
		return
	}

	ruleID := ctx.Params.ByName("id")
	if !rules.RuleIDRe.MatchString(ruleID) {
		http.Error(ctx.Writer, "Please supply a valid rule ID.", http.StatusBadRequest)
		return
	}

	blob, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		http.Error(ctx.Writer, "Failed to read request body.", http.StatusBadRequest)
		return
	}
	var u *ruleUpdateRequest
	if err := json.Unmarshal(blob, &u); err != nil {
		logging.Errorf(ctx.Context, "Failed to umarshal rule update request: %s", err)
		http.Error(ctx.Writer, "Incorrectly formed request: invalid json.", http.StatusBadRequest)
		return
	}
	if msg := validateUpdate(u); msg != "" {
		http.Error(ctx.Writer, msg, http.StatusBadRequest)
		return
	}

	requestedETag := ctx.Request.Header.Get("If-Match")
	user := auth.CurrentUser(ctx.Context).Email

	var concurrentModification bool
	var updatedRule *rules.FailureAssociationRule
	commitTime, err := span.ReadWriteTransaction(ctx.Context, func(ctx context.Context) error {
		rule, err := rules.Read(ctx, projectID, ruleID)
		if err != nil {
			return err
		}
		if requestedETag != "" && ruleETag(rule) != requestedETag {
			concurrentModification = true
			return nil
		}
		for _, path := range u.UpdateMask.Paths {
			// Only limited fields may be modified by the client.
			switch path {
			case "ruleDefinition":
				rule.RuleDefinition = u.Rule.RuleDefinition
			case "bug":
				rule.Bug = u.Rule.Bug
			case "isActive":
				rule.IsActive = u.Rule.IsActive
			default:
				return fmt.Errorf("unsupported field update: %s", path)
			}
		}
		if err := rules.Update(ctx, rule, user); err != nil {
			return err
		}
		concurrentModification = false
		updatedRule = rule
		return nil
	})
	if err != nil {
		logging.Errorf(ctx.Context, "Updating rule %s: %s", ruleID, err)
		http.Error(ctx.Writer, "Internal server error.", http.StatusInternalServerError)
		return
	}
	if concurrentModification {
		http.Error(ctx.Writer, "The rule was modified since it was last read; the update was not applied.", http.StatusConflict)
		return
	}
	updatedRule.LastUpdated = commitTime.In(time.UTC)
	updatedRule.LastUpdatedUser = user
	ctx.Writer.Header().Add("ETag", ruleETag(updatedRule))
	respondWithJSON(ctx, updatedRule)
}

func validateUpdate(update *ruleUpdateRequest) string {
	for _, path := range update.UpdateMask.Paths {
		switch path {
		case "ruleDefinition":
			_, err := lang.Parse(update.Rule.RuleDefinition, rules.Identifiers...)
			if err != nil {
				return fmt.Sprintf("Rule definition is not valid: %s", err)
			}
		case "bug":
			if err := update.Rule.Bug.Validate(); err != nil {
				return "Bug is not valid."
			}
		case "isActive":
			continue
		default:
			return fmt.Sprintf("Unsupported update mask - %s", path)
		}
	}
	return ""
}
