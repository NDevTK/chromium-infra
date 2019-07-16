// Copyright 2019 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package backend

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"infra/appengine/arquebus/app/config"
	"infra/appengine/rotang/proto/rotangapi"
	"infra/monorailv2/api/api_proto"
)

func TestSearchAndUpdateIssues(t *testing.T) {
	t.Parallel()
	assignerID := "test-assigner"

	Convey("searchAndUpdateIssues", t, func() {
		c := createTestContextWithTQ()

		// create a sample assigner with tasks.
		assigner := createAssigner(c, assignerID)
		assigner.AssigneesRaw = createRawUserSources(
			oncallUserSource("Rotation 1", config.Oncall_PRIMARY),
		)
		assigner.CCsRaw = createRawUserSources()
		tasks := triggerScheduleTaskHandler(c, assignerID)
		So(tasks, ShouldNotBeNil)
		task := tasks[0]

		mockListIssues(
			c, &monorail.Issue{ProjectName: "test", LocalId: 123},
			&monorail.Issue{ProjectName: "test", LocalId: 456},
		)

		Convey("tickets with opt-out label are filtered in search", func() {
			countOptOptLabel := func(query string) int {
				assigner.IssueQuery.Q = query
				_, err := searchAndUpdateIssues(c, assigner, task)
				So(err, ShouldBeNil)
				req := getListIssuesRequest(c)
				So(req, ShouldNotBeNil)
				return strings.Count(req.Query, fmt.Sprintf("-label:%s", OptOutLabel))
			}
			So(countOptOptLabel("ABC"), ShouldEqual, 1)
			So(countOptOptLabel("ABC OR "), ShouldEqual, 1)
			So(countOptOptLabel("ABC OR"), ShouldEqual, 1)
			So(countOptOptLabel("ABC DEF"), ShouldEqual, 1)
			So(countOptOptLabel(" OR ABC"), ShouldEqual, 1)
			So(countOptOptLabel("OR ABC DEF"), ShouldEqual, 1)
			So(countOptOptLabel("ABC OR DEF"), ShouldEqual, 2)
			So(countOptOptLabel("ABC OR DEF OR FOO"), ShouldEqual, 3)
		})

		Convey("issues are updated", func() {
			nUpdated, err := searchAndUpdateIssues(c, assigner, task)
			So(err, ShouldBeNil)
			So(nUpdated, ShouldEqual, 2)

			req := getIssueUpdateRequest(c, "test", 123)
			So(req, ShouldNotBeNil)
			So(
				req.Delta.OwnerRef.DisplayName, ShouldEqual,
				findPrimaryOncall(sampleOncallShifts["Rotation 1"]).DisplayName,
			)
			req = getIssueUpdateRequest(c, "test", 456)
			So(req, ShouldNotBeNil)
			So(
				req.Delta.OwnerRef.DisplayName, ShouldEqual,
				findPrimaryOncall(sampleOncallShifts["Rotation 1"]).DisplayName,
			)
		})

		Convey("no issues are updated", func() {
			mockListIssues(
				c, &monorail.Issue{ProjectName: "test", LocalId: 123},
			)

			Convey("if no oncaller is available", func() {
				// simulate an oncall with empty shifts.
				mockOncall(c, "Rotation 1", &rotangapi.ShiftEntry{})

				// nUpdated should be 0
				nUpdated, err := searchAndUpdateIssues(c, assigner, task)
				So(err, ShouldBeNil)
				So(nUpdated, ShouldEqual, 0)
			})

			Convey("if no assignees and ccs set in config", func() {
				assigner.AssigneesRaw = createRawUserSources()
				assigner.CCsRaw = createRawUserSources()

				// nUpdated should be 0
				nUpdated, err := searchAndUpdateIssues(c, assigner, task)
				So(err, ShouldBeNil)
				So(nUpdated, ShouldEqual, 0)
			})

			Convey("if no delta was found", func() {
				assigner.AssigneesRaw = createRawUserSources(
					emailUserSource("foo@example.org"),
				)
				assigner.CCsRaw = createRawUserSources(
					emailUserSource("bar@example.net"),
				)
				mockListIssues(
					c, &monorail.Issue{
						ProjectName: "test", LocalId: 123,
						OwnerRef: monorailUser("foo@example.org"),
						CcRefs: []*monorail.UserRef{
							monorailUser("bar@example.net"),
						},
					},
				)
				nUpdated, err := searchAndUpdateIssues(c, assigner, task)
				So(err, ShouldBeNil)
				So(nUpdated, ShouldEqual, 0)
			})

			Convey("if dry-run is set", func() {
				assigner.IsDryRun = true
				assigner.AssigneesRaw = createRawUserSources(
					emailUserSource("foo@example.org"),
				)
				assigner.CCsRaw = createRawUserSources(
					emailUserSource("bar@example.net"),
				)
				mockListIssues(
					c, &monorail.Issue{ProjectName: "test", LocalId: 123},
				)
				nUpdated, err := searchAndUpdateIssues(c, assigner, task)
				So(err, ShouldBeNil)
				So(nUpdated, ShouldEqual, 0)
			})
		})
	})
}
