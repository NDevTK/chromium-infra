// Copyright 2020 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package reviewer

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
	"go.chromium.org/luci/common/proto"
	gerritpb "go.chromium.org/luci/common/proto/gerrit"
	. "go.chromium.org/luci/common/testing/assertions"
	"go.chromium.org/luci/gae/impl/memory"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"

	"infra/appengine/rubber-stamper/config"
	"infra/appengine/rubber-stamper/tasks/taskspb"
)

func TestReviewCleanCherryPick(t *testing.T) {
	Convey("review clean cherry pick", t, func() {
		ctx := memory.Use(context.Background())

		ctl := gomock.NewController(t)
		defer ctl.Finish()
		gerritMock := gerritpb.NewMockGerritClient(ctl)

		cfg := &config.Config{
			DefaultTimeWindow: "7d",
			HostConfigs: map[string]*config.HostConfig{
				"test-host": {
					RepoConfigs: map[string]*config.RepoConfig{},
				},
			},
		}

		t := &taskspb.ChangeReviewTask{
			Host:               "test-host",
			Number:             12345,
			Revision:           "123abc",
			Repo:               "dummy",
			AutoSubmit:         false,
			CherryPickOfChange: 12121,
		}

		Convey("decline when cherry pick has more than 1 revisions", func() {
			t.RevisionsCount = 2
			msg, err := reviewCleanCherryPick(ctx, cfg, gerritMock, t)
			So(err, ShouldBeNil)
			So(msg, ShouldEqual, "The change cannot be reviewed. There are more than one revision uploaded.")
		})
		Convey("decline when out of configured time window", func() {
			Convey("global time window works", func() {
				gerritMock.EXPECT().GetChange(gomock.Any(), proto.MatcherEqual(&gerritpb.GetChangeRequest{
					Number:  t.CherryPickOfChange,
					Options: []gerritpb.QueryOption{gerritpb.QueryOption_CURRENT_REVISION},
				})).Return(&gerritpb.ChangeInfo{
					CurrentRevision: "456def",
					Revisions: map[string]*gerritpb.RevisionInfo{
						"456def": {
							Created: timestamppb.New(time.Now().Add(-8 * 24 * time.Hour)),
						},
						"789aaa": {
							Created: timestamppb.New(time.Now().Add(-9 * 24 * time.Hour)),
						},
					},
					Labels: map[string]*gerritpb.LabelInfo{
						"Bot-Commit": {
							All: []*gerritpb.ApprovalInfo{
								{
									User: &gerritpb.AccountInfo{
										Name:  "a cute bot",
										Email: "bot@example.com",
									},
									Value: 1,
								},
							},
						},
					},
				}, nil)
				msg, err := reviewCleanCherryPick(ctx, cfg, gerritMock, t)
				So(err, ShouldBeNil)
				So(msg, ShouldEqual, "The change is not in the configured time window. Rubber Stamper is only allowed to review cherry-picks within 7 day(s).")
			})
			Convey("host-level time window works", func() {
				cfg.HostConfigs["test-host"].CleanCherryPickTimeWindow = "5d"
				gerritMock.EXPECT().GetChange(gomock.Any(), proto.MatcherEqual(&gerritpb.GetChangeRequest{
					Number:  t.CherryPickOfChange,
					Options: []gerritpb.QueryOption{gerritpb.QueryOption_CURRENT_REVISION},
				})).Return(&gerritpb.ChangeInfo{
					CurrentRevision: "456def",
					Revisions: map[string]*gerritpb.RevisionInfo{
						"456def": {
							Created: timestamppb.New(time.Now().Add(-6 * 24 * time.Hour)),
						},
						"789aaa": {
							Created: timestamppb.New(time.Now().Add(-9 * 24 * time.Hour)),
						},
					},
					Labels: map[string]*gerritpb.LabelInfo{
						"Bot-Commit": {
							All: []*gerritpb.ApprovalInfo{
								{
									User: &gerritpb.AccountInfo{
										Name:  "a cute bot",
										Email: "bot@example.com",
									},
									Value: 1,
								},
							},
						},
					},
				}, nil)
				msg, err := reviewCleanCherryPick(ctx, cfg, gerritMock, t)
				So(err, ShouldBeNil)
				So(msg, ShouldEqual, "The change is not in the configured time window. Rubber Stamper is only allowed to review cherry-picks within 5 day(s).")
			})
			Convey("repo-level time window works", func() {
				cfg.HostConfigs["test-host"].CleanCherryPickTimeWindow = "5d"
				cfg.HostConfigs["test-host"].RepoConfigs["dummy"] = &config.RepoConfig{
					CleanCherryPickPattern: &config.CleanCherryPickPattern{
						TimeWindow: "58m",
					},
				}
				gerritMock.EXPECT().GetChange(gomock.Any(), proto.MatcherEqual(&gerritpb.GetChangeRequest{
					Number:  t.CherryPickOfChange,
					Options: []gerritpb.QueryOption{gerritpb.QueryOption_CURRENT_REVISION},
				})).Return(&gerritpb.ChangeInfo{
					CurrentRevision: "456def",
					Revisions: map[string]*gerritpb.RevisionInfo{
						"456def": {
							Created: timestamppb.New(time.Now().Add(-time.Hour)),
						},
						"789aaa": {
							Created: timestamppb.New(time.Now().Add(-9 * 24 * time.Hour)),
						},
					},
					Labels: map[string]*gerritpb.LabelInfo{
						"Bot-Commit": {
							All: []*gerritpb.ApprovalInfo{
								{
									User: &gerritpb.AccountInfo{
										Name:  "a cute bot",
										Email: "bot@example.com",
									},
									Value: 1,
								},
							},
						},
					},
				}, nil)
				msg, err := reviewCleanCherryPick(ctx, cfg, gerritMock, t)
				So(err, ShouldBeNil)
				So(msg, ShouldEqual, "The change is not in the configured time window. Rubber Stamper is only allowed to review cherry-picks within 58 minute(s).")
			})
		})
		Convey("decline when alters any excluded file", func() {
			cfg.HostConfigs["test-host"].RepoConfigs["dummy"] = &config.RepoConfig{
				CleanCherryPickPattern: &config.CleanCherryPickPattern{
					ExcludedPaths: []string{"p/q/**", "**.c"},
				},
			}
			gerritMock.EXPECT().GetChange(gomock.Any(), proto.MatcherEqual(&gerritpb.GetChangeRequest{
				Number:  t.CherryPickOfChange,
				Options: []gerritpb.QueryOption{gerritpb.QueryOption_CURRENT_REVISION},
			})).Return(&gerritpb.ChangeInfo{
				CurrentRevision: "456def",
				Revisions: map[string]*gerritpb.RevisionInfo{
					"456def": {
						Created: timestamppb.New(time.Now().Add(-5 * 24 * time.Hour)),
					},
					"789aaa": {
						Created: timestamppb.New(time.Now().Add(-9 * 24 * time.Hour)),
					},
				},
				Labels: map[string]*gerritpb.LabelInfo{
					"Bot-Commit": {
						All: []*gerritpb.ApprovalInfo{
							{
								User: &gerritpb.AccountInfo{
									Name:  "a cute bot",
									Email: "bot@example.com",
								},
								Value: 1,
							},
						},
					},
				},
			}, nil)
			gerritMock.EXPECT().ListFiles(gomock.Any(), proto.MatcherEqual(&gerritpb.ListFilesRequest{
				Number:     t.Number,
				RevisionId: t.Revision,
			})).Return(&gerritpb.ListFilesResponse{
				Files: map[string]*gerritpb.FileInfo{
					"p/q/o/0.txt": nil,
					"valid.md":    nil,
					"a/invalid.c": nil,
				},
			}, nil)
			msg, err := reviewCleanCherryPick(ctx, cfg, gerritMock, t)
			So(err, ShouldBeNil)
			So(msg, ShouldEqual, "The change contains the following files which require a human reviewer: a/invalid.c, p/q/o/0.txt.")
		})
		Convey("decline when not mergeable", func() {
			gerritMock.EXPECT().GetChange(gomock.Any(), proto.MatcherEqual(&gerritpb.GetChangeRequest{
				Number:  t.CherryPickOfChange,
				Options: []gerritpb.QueryOption{gerritpb.QueryOption_CURRENT_REVISION},
			})).Return(&gerritpb.ChangeInfo{
				CurrentRevision: "456def",
				Revisions: map[string]*gerritpb.RevisionInfo{
					"456def": {
						Created: timestamppb.New(time.Now().Add(-5 * 24 * time.Hour)),
					},
					"789aaa": {
						Created: timestamppb.New(time.Now().Add(-9 * 24 * time.Hour)),
					},
				},
				Labels: map[string]*gerritpb.LabelInfo{
					"Bot-Commit": {
						All: []*gerritpb.ApprovalInfo{
							{
								User: &gerritpb.AccountInfo{
									Name:  "a cute bot",
									Email: "bot@example.com",
								},
								Value: 1,
							},
						},
					},
				},
			}, nil)
			gerritMock.EXPECT().GetMergeable(gomock.Any(), proto.MatcherEqual(&gerritpb.GetMergeableRequest{
				Number:     t.Number,
				Project:    t.Repo,
				RevisionId: t.Revision,
			})).Return(&gerritpb.MergeableInfo{
				Mergeable: false,
			}, nil)
			msg, err := reviewCleanCherryPick(ctx, cfg, gerritMock, t)
			So(err, ShouldBeNil)
			So(msg, ShouldEqual, "The change is not mergeable.")
		})
		Convey("return error works", func() {
			Convey("Gerrit GetChange API returns error", func() {
				gerritMock.EXPECT().GetChange(gomock.Any(), proto.MatcherEqual(&gerritpb.GetChangeRequest{
					Number:  t.CherryPickOfChange,
					Options: []gerritpb.QueryOption{gerritpb.QueryOption_CURRENT_REVISION},
				})).Return(nil, grpc.Errorf(codes.NotFound, "not found"))
				msg, err := reviewCleanCherryPick(ctx, cfg, gerritMock, t)
				So(msg, ShouldEqual, "")
				So(err, ShouldErrLike, "gerrit GetChange rpc call failed with error")
			})
			Convey("Gerrit ListFiles API returns error", func() {
				cfg.HostConfigs["test-host"].RepoConfigs["dummy"] = &config.RepoConfig{
					CleanCherryPickPattern: &config.CleanCherryPickPattern{
						ExcludedPaths: []string{"p/q/**", "**.c"},
					},
				}
				gerritMock.EXPECT().GetChange(gomock.Any(), proto.MatcherEqual(&gerritpb.GetChangeRequest{
					Number:  t.CherryPickOfChange,
					Options: []gerritpb.QueryOption{gerritpb.QueryOption_CURRENT_REVISION},
				})).Return(&gerritpb.ChangeInfo{
					CurrentRevision: "456def",
					Revisions: map[string]*gerritpb.RevisionInfo{
						"456def": {
							Created: timestamppb.New(time.Now().Add(-5 * 24 * time.Hour)),
						},
						"789aaa": {
							Created: timestamppb.New(time.Now().Add(-9 * 24 * time.Hour)),
						},
					},
					Labels: map[string]*gerritpb.LabelInfo{
						"Bot-Commit": {
							All: []*gerritpb.ApprovalInfo{
								{
									User: &gerritpb.AccountInfo{
										Name:  "a cute bot",
										Email: "bot@example.com",
									},
									Value: 1,
								},
							},
						},
					},
				}, nil)
				gerritMock.EXPECT().ListFiles(gomock.Any(), proto.MatcherEqual(&gerritpb.ListFilesRequest{
					Number:     t.Number,
					RevisionId: t.Revision,
				})).Return(nil, grpc.Errorf(codes.NotFound, "not found"))
				msg, err := reviewCleanCherryPick(ctx, cfg, gerritMock, t)
				So(msg, ShouldEqual, "")
				So(err, ShouldErrLike, "gerrit ListFiles rpc call failed with error")
			})
			Convey("Gerrit GetMergeable API returns error", func() {
				gerritMock.EXPECT().GetChange(gomock.Any(), proto.MatcherEqual(&gerritpb.GetChangeRequest{
					Number:  t.CherryPickOfChange,
					Options: []gerritpb.QueryOption{gerritpb.QueryOption_CURRENT_REVISION},
				})).Return(&gerritpb.ChangeInfo{
					CurrentRevision: "456def",
					Revisions: map[string]*gerritpb.RevisionInfo{
						"456def": {
							Created: timestamppb.New(time.Now().Add(-5 * 24 * time.Hour)),
						},
						"789aaa": {
							Created: timestamppb.New(time.Now().Add(-9 * 24 * time.Hour)),
						},
					},
					Labels: map[string]*gerritpb.LabelInfo{
						"Bot-Commit": {
							All: []*gerritpb.ApprovalInfo{
								{
									User: &gerritpb.AccountInfo{
										Name:  "a cute bot",
										Email: "bot@example.com",
									},
									Value: 1,
								},
							},
						},
					},
				}, nil)
				gerritMock.EXPECT().GetMergeable(gomock.Any(), proto.MatcherEqual(&gerritpb.GetMergeableRequest{
					Number:     t.Number,
					Project:    t.Repo,
					RevisionId: t.Revision,
				})).Return(nil, grpc.Errorf(codes.NotFound, "not found"))
				msg, err := reviewCleanCherryPick(ctx, cfg, gerritMock, t)
				So(msg, ShouldEqual, "")
				So(err, ShouldErrLike, "gerrit GetMergeable rpc call failed with error")
			})
			Convey("time window config error", func() {
				cfg.HostConfigs["test-host"].CleanCherryPickTimeWindow = "112-1d"
				gerritMock.EXPECT().GetChange(gomock.Any(), proto.MatcherEqual(&gerritpb.GetChangeRequest{
					Number:  t.CherryPickOfChange,
					Options: []gerritpb.QueryOption{gerritpb.QueryOption_CURRENT_REVISION},
				})).Return(&gerritpb.ChangeInfo{
					CurrentRevision: "456def",
					Revisions: map[string]*gerritpb.RevisionInfo{
						"456def": {
							Created: timestamppb.New(time.Now().Add(-5 * 24 * time.Hour)),
						},
						"789aaa": {
							Created: timestamppb.New(time.Now().Add(-9 * 24 * time.Hour)),
						},
					},
					Labels: map[string]*gerritpb.LabelInfo{
						"Bot-Commit": {
							All: []*gerritpb.ApprovalInfo{
								{
									User: &gerritpb.AccountInfo{
										Name:  "a cute bot",
										Email: "bot@example.com",
									},
									Value: 1,
								},
							},
						},
					},
				}, nil)
				msg, err := reviewCleanCherryPick(ctx, cfg, gerritMock, t)
				So(msg, ShouldEqual, "")
				So(err, ShouldErrLike, "invalid time_window config 112-1d")
			})
		})
	})
}
