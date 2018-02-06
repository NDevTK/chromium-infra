// Copyright 2017 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package gerrit

import (
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"infra/tricium/api/v1"
)

func TestComposeChangesQueryURL(t *testing.T) {
	Convey("Test Environment", t, func() {
		instance := "https://chromium-review.googlesource.com"
		project := "playground/gerrit-tricium"
		formattedProject := "playground%2Fgerrit-tricium"
		const form = "2006-01-02 15:04:05.000000000"
		time, err := time.Parse(form, "2016-10-01 10:00:05.640000000")
		So(err, ShouldBeNil)
		formattedTime := "2016-10-01+10%3A00%3A05.640000000"
		Convey("First page of poll", func() {
			So(composeChangesQueryURL(instance, project, time, 0), ShouldEqual,
				fmt.Sprintf("%s/a/changes/?o=CURRENT_REVISION&o=CURRENT_FILES&o=DETAILED_ACCOUNTS&q=project%%3A%s+after%%3A%%22%s%%22&start=0",
					instance, formattedProject, formattedTime))
		})
	})
}

func TestCreateRobotComment(t *testing.T) {
	Convey("Test Enviroinment", t, func() {
		runID := int64(1234567)
		uuid := "7ae6f43d-22e9-4350-ace4-1fee9014509a"
		Convey("File comment has no position info", func() {
			roco := createRobotComment(runID, tricium.Data_Comment{
				Id:       uuid,
				Path:     "README.md",
				Message:  "Message",
				Category: "Hello",
			})
			So(roco.Line, ShouldEqual, 0)
			So(roco.Range, ShouldBeNil)
		})
		Convey("Line comment has no range info", func() {
			line := int32(10)
			roco := createRobotComment(runID, tricium.Data_Comment{
				Id:        uuid,
				Path:      "README.md",
				Message:   "Message",
				Category:  "Hello",
				StartLine: line,
			})
			So(roco.Line, ShouldEqual, line)
			So(roco.Range, ShouldBeNil)
		})
		Convey("Range comment has range", func() {
			startLine := int32(10)
			endLine := int32(20)
			startChar := int32(2)
			endChar := int32(18)
			roco := createRobotComment(runID, tricium.Data_Comment{
				Id:        uuid,
				Path:      "README.md",
				Message:   "Message",
				Category:  "Hello",
				StartLine: startLine,
				EndLine:   endLine,
				StartChar: startChar,
				EndChar:   endChar,
			})
			So(roco.Message, ShouldEqual, "Message")
			So(roco.RobotID, ShouldEqual, "Hello")
			So(roco.Properties["tricium_comment_uuid"], ShouldEqual, uuid)
			So(roco.Line, ShouldEqual, startLine)
			So(roco.Range, ShouldNotBeNil)
			So(roco.Range.StartLine, ShouldEqual, startLine)
			So(roco.Range.EndLine, ShouldEqual, endLine)
			So(roco.Range.StartCharacter, ShouldEqual, startChar)
			So(roco.Range.EndCharacter, ShouldEqual, endChar)
		})

	})
}
