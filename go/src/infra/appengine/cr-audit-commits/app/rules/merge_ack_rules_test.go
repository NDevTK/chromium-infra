// Copyright 2018 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package rules

import (
	"context"
	"testing"

	"infra/monorail"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMergeAckRules(t *testing.T) {
	t.Parallel()
	Convey("Merge Acknowledgement rules work", t, func() {
		ctx := context.Background()
		rc := &RelevantCommit{
			CommitHash:    "b07c0de",
			Status:        AuditScheduled,
			CommitMessage: "Acknowledging merges into a release branch",
		}
		cfg := &RefConfig{
			BaseRepoURL: "https://a.googlesource.com/a.git",
			GerritURL:   "https://a-review.googlesource.com/",
			BranchName:  "3325",
			Metadata:    "MilestoneNumber:65",
		}
		ap := &AuditParams{
			TriggeringAccount: "releasebot@sample.com",
			RepoCfg:           cfg,
		}
		testClients := &Clients{}
		testClients.Monorail = MockMonorailClient{
			Gi: &monorail.Issue{},
			Ii: &monorail.InsertIssueResponse{
				Issue: &monorail.Issue{},
			},
		}
		Convey("Change to commit has a valid bug", func() {
			testClients.Monorail = MockMonorailClient{
				Gi: &monorail.Issue{
					Id: 123456,
				},
			}
			rc.CommitMessage = "This change has a valid bug ID \nBug:123456"
			// Run rule
			rr, _ := AcknowledgeMerge{}.Run(ctx, ap, rc, testClients)
			So(rr.RuleResultStatus, ShouldEqual, RulePassed)
		})
		Convey("Change to commit has no bug", func() {
			testClients.Monorail = MockMonorailClient{
				Gi: &monorail.Issue{
					Id: 123456,
				},
			}
			rc.CommitMessage = "This change has no bug attached"
			// Run rule
			rr, _ := AcknowledgeMerge{}.Run(ctx, ap, rc, testClients)
			So(rr.RuleResultStatus, ShouldEqual, RulePassed)
		})
	})
}
