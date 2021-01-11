// Copyright 2020 The LUCI Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package eval

import (
	"testing"

	evalpb "infra/rts/presubmit/eval/proto"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPSURL(t *testing.T) {
	t.Parallel()
	Convey(`psURL`, t, func() {
		patchSet := &evalpb.GerritPatchset{
			Change: &evalpb.GerritChange{
				Host:   "example.googlesource.com",
				Number: 123,
			},
			Patchset: 4,
		}
		So(psURL(patchSet), ShouldEqual, "https://example.googlesource.com/c/123/4")
	})
}
