// Copyright 2020 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package rules

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go.chromium.org/luci/common/proto/git"
)

func TestAccountRules(t *testing.T) {
	t.Parallel()
	Convey("AccountRules", t, func() {
		commit := &git.Commit{
			Author:    &git.Commit_User{Email: "dummy1@test1.com"},
			Committer: &git.Commit_User{Email: "dummy2@test2.com"},
		}
		So(AccountRules{Account: "dummy1@test1.com"}.MatchesCommit(commit), ShouldBeTrue)
		So(AccountRules{Account: "dummy2@test2.com"}.MatchesCommit(commit), ShouldBeTrue)
		So(AccountRules{Account: "dummy3@test3.com"}.MatchesCommit(commit), ShouldBeFalse)
		So(AccountRules{Account: "*"}.MatchesCommit(commit), ShouldBeTrue)
	})
}

func TestGetReleaseConfig(t *testing.T) {
	t.Parallel()
	Convey("Ensure starting commit, branch name, and milestone number are valid", t, func() {
		ctx := context.Background()
		cfg := RefConfig{}
		milestoneNumber := ""
		success := false
		branchRefsURLContents := []string{"\"commit\": \"e8b8df68cc0a4623567482825115b8a321d01eb9\""}
		branchInfos := []BranchInfo{{PdfiumBranch: "", SkiaBranch: "", WebrtcBranch: "", V8Branch: "", ChromiumBranch: "3440", Milestone: 68}}
		// Call the release config method
		concreteConfigs, err := GetReleaseConfig(ctx, cfg, branchRefsURLContents, branchInfos)
		So(err, ShouldBeNil)
		for i := range concreteConfigs {
			So(concreteConfigs[i].StartingCommit, ShouldEqual, "e8b8df68cc0a4623567482825115b8a321d01eb9")
			So(concreteConfigs[i].BranchName, ShouldEqual, "refs/branch-heads/3440")
			milestoneNumber, success = GetToken(ctx, "MilestoneNumber", concreteConfigs[i].Metadata)
			So(milestoneNumber, ShouldEqual, "68")
			So(success, ShouldEqual, true)
		}
	})
}
