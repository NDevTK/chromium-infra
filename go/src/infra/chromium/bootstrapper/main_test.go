// Copyright 2021 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package main

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"infra/chromium/bootstrapper/cipd"
	fakecipd "infra/chromium/bootstrapper/fakes/cipd"
	fakegitiles "infra/chromium/bootstrapper/fakes/gitiles"
	"infra/chromium/bootstrapper/gitiles"
	. "infra/chromium/bootstrapper/util"

	. "github.com/smartystreets/goconvey/convey"
	buildbucketpb "go.chromium.org/luci/buildbucket/proto"
	. "go.chromium.org/luci/common/testing/assertions"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func strPtr(s string) *string {
	return &s
}

func createInput(buildJson string) io.Reader {
	build := &buildbucketpb.Build{}
	PanicOnError(protojson.Unmarshal([]byte(buildJson), build))
	buildProtoBytes, err := proto.Marshal(build)
	PanicOnError(err)
	return bytes.NewReader(buildProtoBytes)
}

type reader struct {
	readFn func([]byte) (int, error)
}

func (r reader) Read(p []byte) (n int, err error) {
	return r.readFn(p)
}

func TestPerformBootstrap(t *testing.T) {
	t.Parallel()

	project := &fakegitiles.Project{
		Refs:      map[string]string{},
		Revisions: map[string]*fakegitiles.Revision{},
	}

	pkg := &fakecipd.Package{
		Refs:      map[string]string{},
		Instances: map[string]*fakecipd.PackageInstance{},
	}

	opts := options{}

	ctx := context.Background()

	ctx = gitiles.UseGitilesClientFactory(ctx, fakegitiles.Factory(map[string]*fakegitiles.Host{
		"fake-host": {
			Projects: map[string]*fakegitiles.Project{
				"fake-project": project,
			},
		},
	}))

	ctx = cipd.UseCipdClientFactory(ctx, fakecipd.Factory(map[string]*fakecipd.Package{
		"fake-package": pkg,
	}))

	Convey("performBootstrap", t, func() {

		cipdRoot := t.TempDir()

		Convey("fails if reading input fails", func() {
			input := reader{func(p []byte) (int, error) {
				return 0, errors.New("test read failure")
			}}

			cmd, err := opts.performBootstrap(ctx, input, cipdRoot, "")

			So(err, ShouldNotBeNil)
			So(cmd, ShouldBeNil)
		})

		Convey("fails if unmarshalling build fails", func() {
			input := strings.NewReader("invalid-proto")

			cmd, err := opts.performBootstrap(ctx, input, cipdRoot, "")

			So(err, ShouldNotBeNil)
			So(cmd, ShouldBeNil)
		})

		Convey("fails if bootstrap fails", func() {
			input := createInput(`{}`)

			cmd, err := opts.performBootstrap(ctx, input, cipdRoot, "")

			So(err, ShouldNotBeNil)
			So(cmd, ShouldBeNil)
		})

		input := createInput(`{
			"input": {
				"properties": {
					"$bootstrap/properties": {
						"top_level_project": {
							"repo": {
								"host": "fake-host",
								"project": "fake-project"
							},
							"ref": "fake-ref"
						},
						"properties_file": "fake-properties-file"
					},
					"$bootstrap/exe": {
						"exe": {
							"cipd_package": "fake-package",
							"cipd_version": "fake-version",
							"cmd": ["fake-exe"]
						}
					}
				}
			}
		}`)

		Convey("fails if determining executable fails", func() {
			project.Refs["fake-ref"] = "fake-revision"
			project.Revisions["fake-revision"] = &fakegitiles.Revision{
				Files: map[string]*string{
					"fake-properties-file": strPtr(`{
						"foo": "bar"
					}`),
				},
			}
			pkg.Refs["fake-version"] = ""

			cmd, err := opts.performBootstrap(ctx, input, cipdRoot, "fake-output-path")

			So(err, ShouldNotBeNil)
			So(cmd, ShouldBeNil)
		})

		Convey("succeeds for valid input", func() {
			project.Refs["fake-ref"] = "fake-revision"
			project.Revisions["fake-revision"] = &fakegitiles.Revision{
				Files: map[string]*string{
					"fake-properties-file": strPtr(`{
						"foo": "bar"
					}`),
				},
			}
			pkg.Refs["fake-version"] = "fake-instance-id"

			cmd, err := opts.performBootstrap(ctx, input, cipdRoot, "fake-output-path")

			So(err, ShouldBeNil)
			So(cmd, ShouldNotBeNil)
			So(cmd.Args, ShouldResemble, []string{
				filepath.Join(cipdRoot, "fake-exe"),
				"--output",
				"fake-output-path",
			})
			contents, _ := ioutil.ReadAll(cmd.Stdin)
			build := &buildbucketpb.Build{}
			proto.Unmarshal(contents, build)
			So(build, ShouldResembleProtoJSON, `{
				"input": {
					"gitiles_commit": {
						"host": "fake-host",
						"project": "fake-project",
						"ref": "fake-ref",
						"id": "fake-revision"
					},
					"properties": {
						"$build/chromium_bootstrap": {
							"commits": [
								{
									"host": "fake-host",
									"project": "fake-project",
									"ref": "fake-ref",
									"id": "fake-revision"
								}
							],
							"exe": {
								"cipd": {
									"server": "https://chrome-infra-packages.appspot.com",
									"package": "fake-package",
									"requested_version": "fake-version",
									"actual_version": "fake-instance-id"
								},
								"cmd": ["fake-exe"]
							}
						},
						"foo": "bar"
					}
				}
			}`)
		})

		Convey("succeeds for properties-optional without $bootstrap/properties", func() {
			input := createInput(`{
				"input": {
					"properties": {
						"$bootstrap/exe": {
							"exe": {
								"cipd_package": "fake-package",
								"cipd_version": "fake-version",
								"cmd": ["fake-exe"]
							}
						}
					}
				}
			}`)
			opts := options{propertiesOptional: true}

			cmd, err := opts.performBootstrap(ctx, input, cipdRoot, "fake-output-path")

			So(err, ShouldBeNil)
			So(cmd, ShouldNotBeNil)
			So(cmd.Args, ShouldResemble, []string{
				filepath.Join(cipdRoot, "fake-exe"),
				"--output",
				"fake-output-path",
			})
			contents, _ := ioutil.ReadAll(cmd.Stdin)
			build := &buildbucketpb.Build{}
			proto.Unmarshal(contents, build)
			So(build, ShouldResembleProtoJSON, `{
				"input": {
					"properties": {
						"$build/chromium_bootstrap": {
							"exe": {
								"cipd": {
									"server": "https://chrome-infra-packages.appspot.com",
									"package": "fake-package",
									"requested_version": "fake-version",
									"actual_version": "fake-instance-id"
								},
								"cmd": ["fake-exe"]
							}
						}
					}
				}
			}`)
		})

	})
}
