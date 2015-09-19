// Copyright 2014 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

/*
Package common defines structures and functions used by all other cipd/ packages.
*/
package common

import (
	"fmt"
	"regexp"
	"strings"
)

// packageNameRe is a regular expression for a package name: <word>/<word/<word>
// Package names must be lower case.
var packageNameRe = regexp.MustCompile(`^([a-z0-9_\-]+/)*[a-z0-9_\-]+$`)

// instanceTagKeyRe is a regular expression for a tag key.
var instanceTagKeyRe = regexp.MustCompile(`^[a-z0-9_\-]+$`)

// packageRefRe is a regular expression for a ref.
var packageRefRe = regexp.MustCompile(`^[a-z0-9_\-]{1,100}$`)

// Pin uniquely identifies an instance of some package.
type Pin struct {
	PackageName string `json:"package"`
	InstanceID  string `json:"instance_id"`
}

// String converts pin to a human readable string.
func (pin Pin) String() string {
	return fmt.Sprintf("%s:%s", pin.PackageName, pin.InstanceID)
}

// ValidatePackageName returns error if a string doesn't look like a valid package name.
func ValidatePackageName(name string) error {
	if !packageNameRe.MatchString(name) {
		return fmt.Errorf("invalid package name: %s", name)
	}
	return nil
}

// ValidateInstanceID returns error if a string doesn't look like a valid package instance id.
func ValidateInstanceID(s string) error {
	// Instance id is SHA1 hex digest currently.
	if len(s) != 40 {
		return fmt.Errorf("not a valid package instance ID %q: not 40 bytes", s)
	}
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			return fmt.Errorf("not a valid package instance ID %q: wrong char %c", s, c)
		}
	}
	return nil
}

// ValidatePin returns error if package name of instance id of a pin are not valid.
func ValidatePin(pin Pin) error {
	if err := ValidatePackageName(pin.PackageName); err != nil {
		return err
	}
	if err := ValidateInstanceID(pin.InstanceID); err != nil {
		return err
	}
	return nil
}

// ValidatePackageRef returns error if a string doesn't look like a valid ref.
func ValidatePackageRef(r string) error {
	if ValidateInstanceID(r) == nil {
		return fmt.Errorf("invalid ref name (looks like an instance ID): %q", r)
	}
	if !packageRefRe.MatchString(r) {
		return fmt.Errorf("invalid ref name: %q", r)
	}
	return nil
}

// ValidateInstanceTag returns error if a string doesn't look like a valid tag.
func ValidateInstanceTag(t string) error {
	chunks := strings.SplitN(t, ":", 2)
	if len(chunks) != 2 {
		return fmt.Errorf("%q doesn't look like a tag (a key:value pair)", t)
	}
	if len(t) > 400 {
		return fmt.Errorf("the tag is too long: %q", t)
	}
	if !instanceTagKeyRe.MatchString(chunks[0]) {
		return fmt.Errorf("invalid tag key in %q (should be a lowercase word)", t)
	}
	return nil
}

// ValidateInstanceVersion return error if a string doesn't look like
// an instance ID or a package ref or an instance tag.
func ValidateInstanceVersion(v string) error {
	if ValidateInstanceID(v) == nil || ValidatePackageRef(v) == nil || ValidateInstanceTag(v) == nil {
		return nil
	}
	return fmt.Errorf("bad version (not an instance ID, a ref or a tag): %q", v)
}
