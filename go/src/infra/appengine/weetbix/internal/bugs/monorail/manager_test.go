// Copyright 2021 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package monorail

import (
	"context"
	"infra/appengine/weetbix/internal/bugs"
	"infra/appengine/weetbix/internal/clustering"
	mpb "infra/monorailv2/api/v3/api_proto"
	"testing"

	"cloud.google.com/go/bigquery"
	. "github.com/smartystreets/goconvey/convey"
	. "go.chromium.org/luci/common/testing/assertions"
	"google.golang.org/genproto/protobuf/field_mask"
)

func NewCluster() *clustering.Cluster {
	return &clustering.Cluster{
		ClusterID:              "ClusterID",
		UnexpectedFailures1d:   1300,
		UnexpectedFailures3d:   3300,
		UnexpectedFailures7d:   7300,
		UnexoneratedFailures1d: 120,
		UnexoneratedFailures3d: 320,
		UnexoneratedFailures7d: 720,
		AffectedRuns1d:         11,
		AffectedRuns3d:         31,
		AffectedRuns7d:         71,
		ExampleFailureReason:   bigquery.NullString{StringVal: "Some failure reason.", Valid: true},
	}
}

func TestManager(t *testing.T) {
	t.Parallel()

	Convey("With Bug Manager", t, func() {
		ctx := context.Background()
		f := &FakeIssuesStore{
			NextID: 100,
		}
		user := AutomationUsers[0]
		cl, err := NewClient(UseFakeIssuesClient(ctx, f, user), "myhost")
		So(err, ShouldBeNil)
		bm := NewBugManager(cl)

		Convey("Create", func() {
			c := NewCluster()
			Convey("With reason-based failure cluster", func() {
				reason := `Expected equality of these values:
					"Expected_Value"
					my_expr.evaluate(123)
						Which is: "Unexpected_Value"`
				c.ClusterID = "ClusterIDShouldNotAppearInOutput"
				c.ExampleFailureReason = bigquery.NullString{StringVal: reason, Valid: true}

				bug, err := bm.Create(ctx, c)
				So(err, ShouldBeNil)
				So(bug, ShouldEqual, "chromium/100")
				So(len(f.Issues), ShouldEqual, 1)
				issue := f.Issues[0]

				So(issue.Issue, ShouldResembleProto, &mpb.Issue{
					Name:     "projects/chromium/issues/100",
					Summary:  "Tests are failing: Expected equality of these values: \"Expected_Value\" my_expr.evaluate(123) Which is: \"Unexpected_Value\"",
					Reporter: AutomationUsers[0],
					State:    mpb.IssueContentState_ACTIVE,
					Status:   &mpb.Issue_StatusValue{Status: "Untriaged"},
					FieldValues: []*mpb.FieldValue{
						{
							// Type field.
							Field: "projects/chromium/fieldDefs/10",
							Value: "Bug",
						},
						{
							// Priority field.
							Field: "projects/chromium/fieldDefs/11",
							Value: "0",
						},
					},
					Labels: []*mpb.Issue_LabelValue{{
						Label: "Restrict-View-Google",
					}, {
						Label: "Weetbix-Managed",
					}},
				})
				So(len(issue.Comments), ShouldEqual, 1)
				So(issue.Comments[0].Content, ShouldContainSubstring, reason)
				So(issue.Comments[0].Content, ShouldNotContainSubstring, "ClusterIDShouldNotAppearInOutput")
			})
			Convey("With test name failure cluster", func() {
				c.ClusterID = "ninja://:blink_web_tests/media/my-suite/my-test.html"
				c.ExampleFailureReason = bigquery.NullString{Valid: false}

				bug, err := bm.Create(ctx, c)
				So(err, ShouldBeNil)
				So(bug, ShouldEqual, "chromium/100")
				So(len(f.Issues), ShouldEqual, 1)
				issue := f.Issues[0]

				So(issue.Issue, ShouldResembleProto, &mpb.Issue{
					Name:     "projects/chromium/issues/100",
					Summary:  "Tests are failing: ninja://:blink_web_tests/media/my-suite/my-test.html",
					Reporter: AutomationUsers[0],
					State:    mpb.IssueContentState_ACTIVE,
					Status:   &mpb.Issue_StatusValue{Status: "Untriaged"},
					FieldValues: []*mpb.FieldValue{
						{
							// Type field.
							Field: "projects/chromium/fieldDefs/10",
							Value: "Bug",
						},
						{
							// Priority field.
							Field: "projects/chromium/fieldDefs/11",
							Value: "0",
						},
					},
					Labels: []*mpb.Issue_LabelValue{{
						Label: "Restrict-View-Google",
					}, {
						Label: "Weetbix-Managed",
					}},
				})
				So(len(issue.Comments), ShouldEqual, 1)
				So(issue.Comments[0].Content, ShouldContainSubstring, "ninja://:blink_web_tests/media/my-suite/my-test.html")
			})
			Convey("Does nothing if in simulation mode", func() {
				bm.Simulate = true
				_, err := bm.Create(ctx, c)
				So(err, ShouldEqual, bugs.ErrCreateSimulated)
				So(len(f.Issues), ShouldEqual, 0)
			})
		})
		Convey("Update", func() {
			c := NewCluster()
			bug, err := bm.Create(ctx, c)
			So(err, ShouldBeNil)
			So(bug, ShouldEqual, "chromium/100")
			So(len(f.Issues), ShouldEqual, 1)
			So(IssuePriority(f.Issues[0].Issue), ShouldEqual, "0")

			bugsToUpdate := []*bugs.BugToUpdate{
				{
					BugName: bug,
					Cluster: c,
				},
			}
			updateDoesNothing := func() {
				originalIssues := CopyIssuesStore(f)
				err := bm.Update(ctx, bugsToUpdate)
				So(err, ShouldBeNil)
				So(f, ShouldResembleIssuesStore, originalIssues)
			}

			Convey("If impact unchanged, does nothing", func() {
				updateDoesNothing()
			})
			Convey("If impact changed", func() {
				c.UnexpectedFailures1d = 1
				bugsToUpdate := []*bugs.BugToUpdate{
					{
						BugName: bug,
						Cluster: c,
					},
				}
				Convey("Adjusts priority in response to changed impact", func() {
					err := bm.Update(ctx, bugsToUpdate)
					So(err, ShouldBeNil)
					So(IssuePriority(f.Issues[0].Issue), ShouldEqual, "3")

					// Verify repeated update has no effect.
					updateDoesNothing()
				})
				Convey("Does not adjust priority if priority manually set", func() {
					// Create a monorail client that interacts with monorail
					// as an end-user. This is needed as we distinguish user
					// updates to the bug from system updates.
					user := "users/100"
					usercl, err := NewClient(UseFakeIssuesClient(ctx, f, user), "myhost")
					So(err, ShouldBeNil)

					updateReq := updateBugPriorityRequest(f.Issues[0].Issue.Name, "1")
					err = usercl.ModifyIssues(ctx, updateReq)
					So(err, ShouldBeNil)
					So(IssuePriority(f.Issues[0].Issue), ShouldEqual, "1")

					// Check the update sets the label.
					expectedIssue := CopyIssue(f.Issues[0].Issue)
					expectedIssue.Labels = append(expectedIssue.Labels, &mpb.Issue_LabelValue{
						Label: manualPriorityLabel,
					})
					SortLabels(expectedIssue.Labels)

					err = bm.Update(ctx, bugsToUpdate)
					So(err, ShouldBeNil)
					So(f.Issues[0].Issue, ShouldResembleProto, expectedIssue)

					// Check repeated update does nothing more.
					updateDoesNothing()

					Convey("Unless manual priority cleared", func() {
						updateReq := clearManualPriorityRequest(f.Issues[0].Issue.Name)
						err = usercl.ModifyIssues(ctx, updateReq)
						So(err, ShouldBeNil)
						So(hasLabel(f.Issues[0].Issue, manualPriorityLabel), ShouldBeFalse)

						err := bm.Update(ctx, bugsToUpdate)
						So(err, ShouldBeNil)
						So(IssuePriority(f.Issues[0].Issue), ShouldEqual, "3")

						// Verify repeated update has no effect.
						updateDoesNothing()
					})
				})
				Convey("Does nothing if in simulation mode", func() {
					bm.Simulate = true
					updateDoesNothing()
				})
			})
		})
	})
}

func updateBugPriorityRequest(name string, priority string) *mpb.ModifyIssuesRequest {
	return &mpb.ModifyIssuesRequest{
		Deltas: []*mpb.IssueDelta{
			{
				Issue: &mpb.Issue{
					Name: name,
					FieldValues: []*mpb.FieldValue{
						{
							Field: priorityFieldName,
							Value: priority,
						},
					},
				},
				UpdateMask: &field_mask.FieldMask{
					Paths: []string{"field_values"},
				},
			},
		},
		CommentContent: "User comment.",
	}
}

func clearManualPriorityRequest(name string) *mpb.ModifyIssuesRequest {
	return &mpb.ModifyIssuesRequest{
		Deltas: []*mpb.IssueDelta{
			{
				Issue: &mpb.Issue{
					Name: name,
				},
				UpdateMask:   &field_mask.FieldMask{},
				LabelsRemove: []string{manualPriorityLabel},
			},
		},
		CommentContent: "User comment.",
	}
}