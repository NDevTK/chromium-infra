// Copyright 2016 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var train = flag.Bool("test.train", false, "regenerate expectations")

func TestMain(t *testing.T) {
	t.Parallel()

	Convey("Main", t, func() {
		testBuildersPyl, err := os.Open("testdata/builders.pyl")
		So(err, ShouldBeNil)
		defer testBuildersPyl.Close()

		buf := &bytes.Buffer{}
		run(testBuildersPyl, buf)
		actual := buf.String()

		const expectedPath = "testdata/expected.cfg"
		if *train {
			err := ioutil.WriteFile(expectedPath, []byte(actual), 0777)
			So(err, ShouldBeNil)
		}

		expected, err := ioutil.ReadFile(expectedPath)
		So(err, ShouldBeNil)
		So(actual, ShouldEqual, string(expected))
	})
}
