// Copyright 2021 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package bootstrap

import (
	"fmt"

	buildbucketpb "go.chromium.org/luci/buildbucket/proto"
	"go.chromium.org/luci/luciexe/exe"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

func setPropertiesFromJson(build *buildbucketpb.Build, propsJson map[string]string) {
	props := make(map[string]interface{}, len(propsJson))
	for key, p := range propsJson {
		s := &structpb.Value{}
		if err := protojson.Unmarshal([]byte(p), s); err != nil {
			panic(err)
		}
		props[key] = s
	}
	if err := exe.WriteProperties(build.Input.Properties, props); err != nil {
		panic(err)
	}
}

func setBootstrapProperties(build *buildbucketpb.Build, propsJson string) {
	setPropertiesFromJson(build, map[string]string{
		"$bootstrap": propsJson,
	})
}

func strPtr(s string) *string {
	return &s
}

func getInput(build *buildbucketpb.Build) *Input {
	input, err := NewInput(build)
	if err != nil {
		panic(err)
	}
	return input
}

func getValueAtPath(s *structpb.Struct, path ...string) *structpb.Value {
	if len(path) < 1 {
		panic("at least one path element must be provided")
	}
	original := s
	for i, p := range path[:len(path)-1] {
		value, ok := s.Fields[p]
		if !ok {
			panic(fmt.Sprintf("path %s is not present in struct %v", path[:i+1], original))
		}
		s = value.GetStructValue()
		if s == nil {
			panic(fmt.Sprintf("path %s is not present in struct %v", path[:i+2], original))
		}
	}
	value, ok := s.Fields[path[len(path)-1]]
	if !ok {
		panic(fmt.Sprintf("path %s is not present in struct %v", path, original))
	}
	return value
}