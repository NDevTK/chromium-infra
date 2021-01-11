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

package git

import (
	"testing"

	"infra/rts/filegraph"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGraph(t *testing.T) {
	t.Parallel()

	Convey(`Graph`, t, func() {
		Convey(`Root of zero value`, func() {
			g := &Graph{}
			root := g.Node("//")
			So(root, ShouldNotBeNil)
			So(root.Name(), ShouldEqual, "//")
		})

		Convey(`node()`, func() {
			g := &Graph{
				root: node{
					children: map[string]*node{
						"dir": {
							children: map[string]*node{
								"foo": {},
							},
						},
					},
				},
			}

			Convey(`//`, func() {
				So(g.node("//"), ShouldEqual, &g.root)
			})

			Convey(`//dir`, func() {
				So(g.node("//dir"), ShouldEqual, g.root.children["dir"])
			})
			Convey(`//dir/foo`, func() {
				So(g.node("//dir/foo"), ShouldEqual, g.root.children["dir"].children["foo"])
			})
			Convey(`//dir/bar`, func() {
				So(g.node("//dir/bar"), ShouldBeNil)
			})
		})

		Convey(`ensureNode`, func() {
			g := &Graph{}
			Convey("//foo/bar", func() {
				bar := g.ensureNode("//foo/bar")
				So(bar, ShouldNotBeNil)
				So(bar.name, ShouldEqual, "//foo/bar")
				So(g.node("//foo/bar"), ShouldEqual, bar)

				foo := g.node("//foo")
				So(foo, ShouldNotBeNil)
				So(foo.name, ShouldEqual, "//foo")
				So(foo.children["bar"], ShouldEqual, bar)
			})

			Convey("already exists", func() {
				So(g.ensureNode("//foo/bar"), ShouldEqual, g.ensureNode("//foo/bar"))
			})

			Convey("//", func() {
				root := g.ensureNode("//")
				So(root, ShouldEqual, &g.root)
			})
		})

		Convey(`sortedChildKeys()`, func() {
			node := &node{
				children: map[string]*node{
					"foo": {},
					"bar": {},
				},
			}
			So(node.sortedChildKeys(), ShouldResemble, []string{"bar", "foo"})
		})

		Convey(`Node(non-existent) returns nil`, func() {
			g := &Graph{}
			// Do not use ShouldBeNil - it checks for interface{} with nil inside,
			// and we need exact nil.
			So(g.Node("//a/b") == nil, ShouldBeTrue)
		})

		Convey(`EdgeReader`, func() {
			bar := &node{commits: 4}
			foo := &node{commits: 2}
			foo.edges = []edge{{to: bar, commonCommits: 1}}
			bar.edges = []edge{{to: foo, commonCommits: 1}}

			type outgoingEdge struct {
				other    filegraph.Node
				distance float64
			}
			var actual []outgoingEdge
			callback := func(other filegraph.Node, distance float64) bool {
				actual = append(actual, outgoingEdge{other: other, distance: distance})
				return true
			}

			r := &EdgeReader{}
			Convey(`Works`, func() {
				r.ReadEdges(foo, callback)
				So(actual, ShouldResemble, []outgoingEdge{{
					other:    bar,
					distance: 1,
				}})
			})
			Convey(`Reversed`, func() {
				r.Reversed = true
				r.ReadEdges(foo, callback)
				So(actual, ShouldResemble, []outgoingEdge{{
					other:    bar,
					distance: 2,
				}})
			})
		})

		Convey(`splitName`, func() {
			Convey("//foo/bar.cc", func() {
				So(splitName("//foo/bar.cc"), ShouldResemble, []string{"foo", "bar.cc"})
			})
			Convey("//", func() {
				So(splitName("//"), ShouldResemble, []string(nil))
			})
		})
	})
}
