// Copyright 2021 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Package testname contains the test name-based clustering algorithm for Weetbix.
package testname

import (
	"crypto/sha256"
	"fmt"
	"infra/appengine/weetbix/internal/clustering"
	"strconv"
)

// AlgorithmVersion is the version of the clustering algorithm. The algorithm
// version should be incremented whenever existing test results may be
// clustered differently (i.e. Cluster(f) returns a different value for some
// f that may have been already ingested).
const AlgorithmVersion = 1

// AlgorithmName is the identifier for the clustering algorithm.
// Weetbix requires all clustering algorithms to have a unique identifier.
// Must match the pattern ^[a-z0-9-.]{1,32}$.
//
// The AlgorithmName must encode the algorithm version, so that each version
// of an algorithm has a different name.
var AlgorithmName = fmt.Sprintf("testname-v%v", AlgorithmVersion)

// Algorithm represents an instance of the test name-based clustering
// algorithm.
type Algorithm struct{}

// Name returns the identifier of the clustering algorithm.
func (a *Algorithm) Name() string {
	return AlgorithmName
}

// Cluster clusters the given test failure and returns its cluster ID (if it
// can be clustered) or nil otherwise.
func (a *Algorithm) Cluster(failure *clustering.Failure) []byte {
	id := failure.TestID
	// Hash test ID to generate a unique fingerprint.
	h := sha256.Sum256([]byte(id))
	// Take first 16 bytes as the ID. (Risk of collision is
	// so low as to not warrant full 32 bytes.)
	return h[0:16]
}

const bugDescriptionTemplate = `This bug is for all test failures with the test name: %s`

// ClusterDescription returns a description of the cluster, for use when
// filing bugs, with the help of the given example failure.
func (a *Algorithm) ClusterDescription(example *clustering.Failure) *clustering.ClusterDescription {
	return &clustering.ClusterDescription{
		Title:       example.TestID,
		Description: fmt.Sprintf(bugDescriptionTemplate, example.TestID),
	}
}

// FailureAssociationRule returns a failure association rule that
// captures the definition of cluster containing the given example.
func (a *Algorithm) FailureAssociationRule(example *clustering.Failure) string {
	stringLiteral := strconv.QuoteToGraphic(example.TestID)
	return fmt.Sprintf("test = %s", stringLiteral)
}
