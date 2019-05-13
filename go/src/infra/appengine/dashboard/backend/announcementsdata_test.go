// Copyright 2019 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package backend

import (
	dashpb "infra/appengine/dashboard/api/dashboard"
	"reflect"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"go.chromium.org/gae/service/datastore"

	. "github.com/smartystreets/goconvey/convey"
)

var chickenAnn = &Announcement{Message: "chicken is missing", Creator: "farmer1"}
var cowAnn = &Announcement{Message: "cow is missing", Creator: "farmer2"}
var retiredAnn = &Announcement{Message: "fox is missing", Creator: "farmer3", Retired: true}

var chickBarnPlat = &Platform{Name: "barn"}
var chickHousePlat = &Platform{Name: "house", URLPaths: []string{"kitchen/*"}}
var chickenPlats = []*Platform{chickBarnPlat, chickHousePlat}

var cowBarnPlat = &Platform{Name: "barn"}
var cowFieldPlat = &Platform{Name: "field"}
var cowPlats = []*Platform{cowBarnPlat, cowFieldPlat}

func TestConvertAnnouncement(t *testing.T) {
	startTS := int64(764797594)
	endTS := int64(764883994)
	testCases := []struct {
		ann       Announcement
		platforms []*Platform
		expected  *dashpb.Announcement
	}{
		{
			ann: Announcement{
				ID:        1234,
				Retired:   true,
				Message:   "CQ issues",
				Creator:   "trooper",
				StartTime: time.Unix(startTS, 0),
				EndTime:   time.Unix(endTS, 0),
			},
			platforms: []*Platform{
				{Name: "monorail"},
				{
					Name:            "gerrit",
					URLPaths:        []string{"c/infra/infra/*", "src/*"},
					AnnouncementKey: &datastore.Key{},
				},
			},
			expected: &dashpb.Announcement{
				Id:             1234,
				MessageContent: "CQ issues",
				Creator:        "trooper",
				Retired:        true,
				StartTime:      &timestamp.Timestamp{Seconds: startTS},
				EndTime:        &timestamp.Timestamp{Seconds: endTS},
				Platforms: []*dashpb.Platform{
					{Name: "monorail"},
					{
						Name:     "gerrit",
						UrlPaths: []string{"c/infra/infra/*", "src/*"},
					},
				},
			},
		},
		{
			ann: Announcement{
				ID:        13,
				StartTime: time.Unix(startTS, 0),
				EndTime:   time.Unix(endTS, 0),
			},
			expected: &dashpb.Announcement{
				Id:             13,
				MessageContent: "",
				Creator:        "",
				Retired:        false,
				StartTime:      &timestamp.Timestamp{Seconds: startTS},
				EndTime:        &timestamp.Timestamp{Seconds: endTS},
				Platforms:      []*dashpb.Platform{},
			},
		},
	}
	for i, tc := range testCases {
		actual, err := tc.ann.ToProto(tc.platforms)
		if err != nil {
			t.Errorf("%d: unexpected error - %s", i, err)
		}
		if !reflect.DeepEqual(tc.expected, actual) {
			t.Errorf("%d: expected %+v, found %+v", i, tc.expected, actual)
		}
	}
}

func TestCreateLiveAnnouncement(t *testing.T) {
	Convey("CreateLiveAnnouncement", t, func() {
		ctx := newTestContext()
		Convey("successful Announcement creator", func() {
			platforms := []*Platform{
				{
					Name:     "monorail",
					URLPaths: []string{"p/chromium/*"},
				},
				{
					Name:     "som",
					URLPaths: []string{"c/infra/infra/*"},
				},
			}
			ann, err := CreateLiveAnnouncement(
				ctx, "Cow cow cow", "cowman", platforms)
			So(err, ShouldBeNil)
			So(ann.Platforms, ShouldHaveLength, 2)
			// Test getting platforms and announcement does not result
			// in error and they were saved correctly in datastore.
			annKey := datastore.NewKey(ctx, "Announcement", "", ann.Id, nil)
			for _, platform := range platforms {
				pKey := datastore.NewKey(ctx, "Platform", platform.Name, 0, annKey)
				existsR, _ := datastore.Exists(ctx, pKey)
				So(existsR.All(), ShouldBeTrue)
			}
			existsR, _ := datastore.Exists(ctx, annKey)
			So(existsR.All(), ShouldBeTrue)
		})
	})
}

func TestGetLiveAnnouncements(t *testing.T) {
	Convey("GetLiveAnnouncements", t, func() {
		ctx := newTestContext()
		datastore.Put(ctx, retiredAnn)
		cowProto, _ := CreateLiveAnnouncement(ctx, cowAnn.Message, cowAnn.Creator, cowPlats)
		chickenProto, _ := CreateLiveAnnouncement(ctx, chickenAnn.Message, chickenAnn.Creator, chickenPlats)
		Convey("get all live announcements", func() {
			anns, err := GetLiveAnnouncements(ctx, "")
			So(err, ShouldBeNil)
			So(anns, ShouldResemble, []*dashpb.Announcement{cowProto, chickenProto})
		})
		Convey("get live announcements for house", func() {
			anns, err := GetLiveAnnouncements(ctx, "house")
			So(err, ShouldBeNil)
			So(anns, ShouldResemble, []*dashpb.Announcement{chickenProto})
		})
		Convey("get live announcements for barn", func() {
			anns, err := GetLiveAnnouncements(ctx, "barn")
			So(err, ShouldBeNil)
			So(anns, ShouldResemble, []*dashpb.Announcement{cowProto, chickenProto})
		})
	})
}
