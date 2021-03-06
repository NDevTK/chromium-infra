// Copyright 2019 The Chromium OS Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package stableversion

import (
	"fmt"
	"sort"
	"strings"

	"github.com/golang/protobuf/jsonpb"
	proto "github.com/golang/protobuf/proto"

	sv "go.chromium.org/chromiumos/infra/proto/go/lab_platform"
)

// AddUpdatedCros add and update the new cros stable version to old.
func AddUpdatedCros(old, updated []*sv.StableCrosVersion) []*sv.StableCrosVersion {
	oldM := make(map[string]*sv.StableCrosVersion, len(old))
	for _, osv := range old {
		oldM[crosSVKey(osv)] = osv
	}

	for _, u := range updated {
		k := crosSVKey(u)
		osv, ok := oldM[k]
		if ok {
			osv.Version = u.GetVersion()
		} else {
			old = append(old, u)
		}
	}
	return old
}

// AddUpdatedFirmware add and update the new firmware stable version to old.
func AddUpdatedFirmware(old, updated []*sv.StableFirmwareVersion) []*sv.StableFirmwareVersion {
	oldM := make(map[string]*sv.StableFirmwareVersion, len(old))
	for _, osv := range old {
		oldM[firmwareSVKey(osv)] = osv
	}

	for _, u := range updated {
		k := firmwareSVKey(u)
		osv, ok := oldM[k]
		if ok {
			osv.Version = u.GetVersion()
		} else {
			old = append(old, u)
		}
	}
	return old
}

func crosSVKey(c *sv.StableCrosVersion) string {
	return c.GetKey().GetBuildTarget().GetName()
}

func firmwareSVKey(f *sv.StableFirmwareVersion) string {
	return fmt.Sprintf("%s:%s", f.GetKey().GetBuildTarget().GetName(), f.GetKey().GetModelId().GetValue())
}

func faftSVKey(f *sv.StableFaftVersion) string {
	return fmt.Sprintf("%s:%s", f.GetKey().GetBuildTarget().GetName(), f.GetKey().GetModelId().GetValue())
}

// WriteSVToString marshals stable version information into a string.
func WriteSVToString(s *sv.StableVersions) (string, error) {
	all := proto.Clone(s).(*sv.StableVersions)
	SortSV(all)
	return (&jsonpb.Marshaler{Indent: "\t"}).MarshalToString(all)
}

// SortSV sorts all the individual entries in a stable version config file.
func SortSV(s *sv.StableVersions) {
	if s == nil {
		return
	}

	c := s.Cros
	sort.SliceStable(c, func(i, j int) bool {
		return strings.ToLower(crosSVKey(c[i])) < strings.ToLower(crosSVKey(c[j]))
	})
	faft := s.Faft
	sort.SliceStable(faft, func(i, j int) bool {
		return strings.ToLower(faftSVKey(faft[i])) < strings.ToLower(faftSVKey(faft[j]))
	})
	fi := s.Firmware
	sort.SliceStable(fi, func(i, j int) bool {
		return strings.ToLower(firmwareSVKey(fi[i])) < strings.ToLower(firmwareSVKey(fi[j]))
	})
}

const separator = ";"

// JoinBuildTargetModel -- join a buildTarget string and a model string to produce a combined key
func JoinBuildTargetModel(buildTarget string, model string) (string, error) {
	b := strings.ToLower(buildTarget)
	m := strings.ToLower(model)
	if err := ValidateJoinBuildTargetModel(b, m); err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%s%s", b, separator, m), nil
}

// FallbackBuildTargetKey creates the key based on the given build target
// This kind of key should only ever be used as a fallback when looking up a stable version.
func FallbackBuildTargetKey(buildTarget string) string {
	return strings.ToLower(buildTarget)
}

// ValidateJoinBuildTargetModel -- checks that a buildTarget and model are valid
// The model is explicitly allowed to be empty.
func ValidateJoinBuildTargetModel(buildTarget string, model string) error {
	if buildTarget == "" {
		return fmt.Errorf("ValidateJoinBuildTargetModel: buildTarget cannot be \"\"")
	}
	if strings.Contains(buildTarget, separator) {
		return fmt.Errorf("ValidateJoinBuildTargetModel: buildTarget cannot contain separator")
	}
	if strings.Contains(model, separator) {
		return fmt.Errorf("ValidateJoinBuildTargetModel: model cannot contain separator")
	}
	return nil
}
