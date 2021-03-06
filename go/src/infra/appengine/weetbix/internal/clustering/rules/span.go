// Copyright 2021 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package rules

import (
	"context"
	"regexp"
	"time"

	"cloud.google.com/go/spanner"

	"go.chromium.org/luci/common/errors"
	"go.chromium.org/luci/server/span"

	"infra/appengine/weetbix/internal/bugs"
	"infra/appengine/weetbix/internal/clustering"
	"infra/appengine/weetbix/internal/clustering/rules/lang"
	"infra/appengine/weetbix/internal/config"
	spanutil "infra/appengine/weetbix/internal/span"
)

// Identifiers is the set of identifiers that are permitted in failure
// association rules.
var Identifiers = []string{"test", "reason"}

// RuleIDRe matches validly formed rule IDs.
var RuleIDRe = regexp.MustCompile(`^[0-9a-f]{32}$`)

// UserRe matches valid users. These are email addresses or the special
// value "weetbix".
var UserRe = regexp.MustCompile(`^weetbix|([a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+)$`)

// WeetbixSystem is the special user that identifies changes made by the
// Weetbix system itself in audit fields.
const WeetbixSystem = "weetbix"

// StartingEpoch is the rule version used for projects that have no rules
// (even inactive rules). It is deliberately different from the timestamp
// zero value to be discernible from "timestamp not populated" programming
// errors.
var StartingEpoch = time.Date(1900, time.January, 1, 0, 0, 0, 0, time.UTC)

// FailureAssociationRule associates failures with a bug. When the rule
// is used to match incoming test failures, the resultant cluster is
// known as a 'bug cluster' because the cluster is associated with a bug
// (via the failure association rule).
type FailureAssociationRule struct {
	// The LUCI Project for which this rule is defined.
	Project string `json:"project"`
	// The unique identifier for the failure association rule,
	// as 32 lowercase hexadecimal characters.
	RuleID string `json:"ruleId"`
	// The rule predicate, defining which failures are being associated.
	RuleDefinition string `json:"ruleDefinition"`
	// The time the rule was created. Output only.
	CreationTime time.Time `json:"creationTime"`
	// The user which created the rule. Output only.
	CreationUser string `json:"creationUser"`
	// The time the rule was last updated. Output only.
	LastUpdated time.Time `json:"lastUpdated"`
	// The user which last updated the rule. Output only.
	LastUpdatedUser string `json:"lastUpdatedUser"`
	// Bug is the identifier of the bug that the failures are
	// associated with.
	Bug bugs.BugID `json:"bug"`
	// Whether the bug should be updated by Weetbix, and whether failures
	// should still be matched against the rule.
	IsActive bool `json:"isActive"`
	// The suggested cluster this rule was created from (if any).
	// Until re-clustering is complete and has reduced the residual impact
	// of the source cluster, this cluster ID tells bug filing to ignore
	// the source cluster when determining whether new bugs need to be filed.
	SourceCluster clustering.ClusterID `json:"sourceCluster"`
}

// Read reads the failure association rule with the given rule ID.
func Read(ctx context.Context, project string, id string) (*FailureAssociationRule, error) {
	whereClause := `RuleId = @ruleId`
	params := map[string]interface{}{
		"ruleId": id,
	}
	rs, err := readWhere(ctx, project, whereClause, params)
	if err != nil {
		return nil, errors.Annotate(err, "query rule by id").Err()
	}
	if len(rs) == 0 {
		return nil, errors.New("no matching rule exists")
	}
	return rs[0], nil
}

// ReadActive reads all active Weetbix failure association rules in the given LUCI project.
func ReadActive(ctx context.Context, projectID string) ([]*FailureAssociationRule, error) {
	whereClause := `IsActive`
	rs, err := readWhere(ctx, projectID, whereClause, nil)
	if err != nil {
		return nil, errors.Annotate(err, "query active rules").Err()
	}
	return rs, nil
}

// ReadDelta reads the changed failure association rules since the given
// timestamp, in the given LUCI project.
func ReadDelta(ctx context.Context, projectID string, sinceTime time.Time) ([]*FailureAssociationRule, error) {
	if sinceTime.Before(StartingEpoch) {
		return nil, errors.New("cannot query rule deltas from before project inception")
	}
	whereClause := `LastUpdated > @sinceTime`
	params := map[string]interface{}{
		"sinceTime": sinceTime,
	}
	rs, err := readWhere(ctx, projectID, whereClause, params)
	if err != nil {
		return nil, errors.Annotate(err, "query rules since").Err()
	}
	return rs, nil
}

// readWhere failure association rules matching the given where clause,
// substituting params for any SQL parameters used in that calsue.
func readWhere(ctx context.Context, projectID string, whereClause string, params map[string]interface{}) ([]*FailureAssociationRule, error) {
	stmt := spanner.NewStatement(`
		SELECT RuleId, RuleDefinition, BugSystem, BugId,
		  CreationTime, LastUpdated,
		  CreationUser, LastUpdatedUser,
		  IsActive,
		  SourceClusterAlgorithm, SourceClusterId
		FROM FailureAssociationRules
		WHERE Project = @projectID AND (` + whereClause + `)
		ORDER BY BugSystem, BugId
	`)
	stmt.Params = make(map[string]interface{})
	for k, v := range params {
		stmt.Params[k] = v
	}
	stmt.Params["projectID"] = projectID

	it := span.Query(ctx, stmt)
	rs := []*FailureAssociationRule{}
	err := it.Do(func(r *spanner.Row) error {
		var ruleID, ruleDefinition, bugSystem, bugID string
		var creationTime, lastUpdated time.Time
		var creationUser, lastUpdatedUser string
		var isActive spanner.NullBool
		var sourceClusterAlgorithm, sourceClusterID string
		err := r.Columns(
			&ruleID, &ruleDefinition, &bugSystem, &bugID,
			&creationTime, &lastUpdated,
			&creationUser, &lastUpdatedUser,
			&isActive,
			&sourceClusterAlgorithm, &sourceClusterID,
		)
		if err != nil {
			return errors.Annotate(err, "read rule row").Err()
		}

		rule := &FailureAssociationRule{
			Project:         projectID,
			RuleID:          ruleID,
			RuleDefinition:  ruleDefinition,
			CreationTime:    creationTime,
			CreationUser:    creationUser,
			LastUpdated:     lastUpdated,
			LastUpdatedUser: lastUpdatedUser,
			Bug:             bugs.BugID{System: bugSystem, ID: bugID},
			IsActive:        isActive.Valid && isActive.Bool,
			SourceCluster: clustering.ClusterID{
				Algorithm: sourceClusterAlgorithm,
				ID:        sourceClusterID,
			},
		}
		rs = append(rs, rule)
		return nil
	})
	return rs, err
}

// ReadLastUpdated reads the last timestamp a rule in the given project was
// updated. This is used to version the set of rules retrieved by ReadActive
// and is typically called in the same transaction.
//
// Simply reading the last LastUpdated time of the rules read by ReadActive
// is not sufficient to version the set of rules read, as the most recent
// update may have been to mark a rule inactive (removing it from the set
// that is read).
//
// If the project has no failure association rules, the timestamp
// StartingEpoch is returned.
func ReadLastUpdated(ctx context.Context, projectID string) (time.Time, error) {
	stmt := spanner.NewStatement(`
		SELECT MAX(LastUpdated) as LastUpdated
		FROM FailureAssociationRules
		WHERE Project = @projectID
	`)
	stmt.Params = map[string]interface{}{
		"projectID": projectID,
	}
	var lastUpdated spanner.NullTime
	it := span.Query(ctx, stmt)
	err := it.Do(func(r *spanner.Row) error {
		err := r.Columns(&lastUpdated)
		if err != nil {
			return errors.Annotate(err, "read last updated row").Err()
		}
		return nil
	})
	if err != nil {
		return time.Time{}, errors.Annotate(err, "query last updated").Err()
	}
	// No failure association rules.
	if !lastUpdated.Valid {
		return StartingEpoch, nil
	}
	return lastUpdated.Time, nil
}

// Create inserts a new failure association rule with the specified details.
func Create(ctx context.Context, r *FailureAssociationRule, user string) error {
	if err := validateRule(r); err != nil {
		return err
	}
	if err := validateUser(user); err != nil {
		return err
	}
	ms := spanutil.InsertMap("FailureAssociationRules", map[string]interface{}{
		"Project":         r.Project,
		"RuleId":          r.RuleID,
		"RuleDefinition":  r.RuleDefinition,
		"CreationTime":    spanner.CommitTimestamp,
		"CreationUser":    user,
		"LastUpdated":     spanner.CommitTimestamp,
		"LastUpdatedUser": user,
		"BugSystem":       r.Bug.System,
		"BugId":           r.Bug.ID,
		// IsActive uses the value 'NULL' to indicate false, and true to indicate true.
		"IsActive":               spanner.NullBool{Bool: r.IsActive, Valid: r.IsActive},
		"SourceClusterAlgorithm": r.SourceCluster.Algorithm,
		"SourceClusterId":        r.SourceCluster.ID,
	})
	span.BufferWrite(ctx, ms)
	return nil
}

// Update updates an existing failure association rule to have the specified
// details.
func Update(ctx context.Context, r *FailureAssociationRule, user string) error {
	if err := validateRule(r); err != nil {
		return err
	}
	if err := validateUser(user); err != nil {
		return err
	}
	ms := spanutil.UpdateMap("FailureAssociationRules", map[string]interface{}{
		"Project":         r.Project,
		"RuleId":          r.RuleID,
		"RuleDefinition":  r.RuleDefinition,
		"LastUpdated":     spanner.CommitTimestamp,
		"LastUpdatedUser": user,
		"BugSystem":       r.Bug.System,
		"BugId":           r.Bug.ID,
		// IsActive uses the value 'NULL' to indicate false, and true to indicate true.
		"IsActive":               spanner.NullBool{Bool: r.IsActive, Valid: r.IsActive},
		"SourceClusterAlgorithm": r.SourceCluster.Algorithm,
		"SourceClusterId":        r.SourceCluster.ID,
	})
	span.BufferWrite(ctx, ms)
	return nil
}

func validateRule(r *FailureAssociationRule) error {
	switch {
	case !config.ProjectRe.MatchString(r.Project):
		return errors.New("project must be valid")
	case !RuleIDRe.MatchString(r.RuleID):
		return errors.New("rule ID must be valid")
	case r.Bug.Validate() != nil:
		return errors.Annotate(r.Bug.Validate(), "bug is not valid").Err()
	case r.SourceCluster.Validate() != nil && !r.SourceCluster.IsEmpty():
		return errors.Annotate(r.SourceCluster.Validate(), "source cluster ID is not valid").Err()
	}
	_, err := lang.Parse(r.RuleDefinition, Identifiers...)
	if err != nil {
		return errors.Annotate(err, "rule definition is not valid").Err()
	}
	return nil
}

func validateUser(u string) error {
	if !UserRe.MatchString(u) {
		return errors.New("user must be valid")
	}
	return nil
}
