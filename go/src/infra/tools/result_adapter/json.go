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

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"go.chromium.org/luci/common/data/stringset"
	"go.chromium.org/luci/common/errors"
	"go.chromium.org/luci/common/logging"

	"go.chromium.org/luci/resultdb/pbutil"
	pb "go.chromium.org/luci/resultdb/proto/v1"
	sinkpb "go.chromium.org/luci/resultdb/sink/proto/v1"
)

const testNamePrefixKey = "test_name_prefix"

var (
	// Subdirectory of outputs root identifying run index.
	testRunSubdirRe = regexp.MustCompile("retry_([0-9]+)/")

	// Regexp for matching about:blank reference artifact.
	aboutBlankRe = regexp.MustCompile(`\babout:blank$`)

	// Known subdirectories of the outputs root where artifacts might be located.
	artifactDirectories = []string{
		"artifacts",
		"layout-test-results",
	}

	// These artifacts contain HTML code that expects to be rendered on the test
	// result server. This is an XSS vulnerability and assumes server
	// implementation details. Milo should implement their functionality itself.
	// Context: crbug.com/1078883
	artifactIgnorelist = stringset.NewFromSlice("pretty_image_diff", "pretty_text_diff")
)

// JSONTestResults represents the structure in
// https://chromium.googlesource.com/chromium/src/+/master/docs/testing/json_test_results_format.md
//
// Deprecated fields and fields not used by Test Results are omitted.
type JSONTestResults struct {
	Interrupted bool `json:"interrupted"`

	PathDelimiter string `json:"path_delimiter"`

	TestsRaw json.RawMessage `json:"tests"`
	Tests    map[string]*TestFields

	Version int32 `json:"version"`

	ArtifactTypes map[string]string `json:"artifact_types"`

	BuildNumber string `json:"build_number"`
	BuilderName string `json:"builder_name"`

	// Metadata associated with results, which may include a list of expectation_files, or
	// test_name_prefix e.g. in GPU tests (distinct from test_id_prefix passed in the recorder API
	// request).
	Metadata map[string]json.RawMessage `json:"metadata"`
}

// TestFields represents the test fields structure in
// https://chromium.googlesource.com/chromium/src/+/master/docs/testing/json_test_results_format.md
//
// Deprecated fields and fields not used by Test Results are omitted.
type TestFields struct {
	Actual   string `json:"actual"`
	Expected string `json:"expected"`

	// TODO(crbug/1034021): Support only map[string][]string, as spec'd by the JSON Test Results
	// Format.
	ArtifactsRaw map[string]json.RawMessage `json:"artifacts"`
	Artifacts    map[string][]string

	Time  float64   `json:"time"`
	Times []float64 `json:"times"`
}

// ConvertFromJSON converts a JSON of test results in the JSON Test Results
// format to the internal struct format.
//
// The receiver is cleared and its fields overwritten.
func (r *JSONTestResults) ConvertFromJSON(reader io.Reader) error {
	*r = JSONTestResults{}
	if err := json.NewDecoder(reader).Decode(r); err != nil {
		return err
	}

	// Convert Tests and return.
	if err := r.convertTests("", r.TestsRaw); err != nil {
		return err
	}
	return nil
}

// ToProtos converts test results in r to []*sinkpb.TestResult.
//
// Does not populate TestResult.Name, TestResult.ResultId or TestResult.TestLocation.
func (r *JSONTestResults) ToProtos(ctx context.Context, normPathToFullPath map[string]string, testLocations bool) ([]*sinkpb.TestResult, error) {
	if r.Version != 3 {
		return nil, errors.Reason("unknown JSON Test Results version %d", r.Version).Err()
	}

	// Sort the test name to make the output deterministic.
	testNames := make([]string, 0, len(r.Tests))
	for name := range r.Tests {
		testNames = append(testNames, name)
	}
	sort.Strings(testNames)

	// Get tags from metadata if any.
	tags, err := r.extractTags()
	if err != nil {
		return nil, err
	}

	// The code below does not return errors, so it is safe to make in-place
	// modifications of inv.
	globalTags := make([]*pb.StringPair, 0, len(tags)+2)
	for _, tag := range tags {
		globalTags = append(globalTags, pbutil.StringPair("json_format_tag", tag))
	}
	globalTags = append(globalTags, pbutil.StringPair(originalFormatTagKey, formatJTR))
	if r.BuildNumber != "" {
		globalTags = append(globalTags, pbutil.StringPair("build_number", r.BuildNumber))
	}

	ret := make([]*sinkpb.TestResult, 0, len(r.Tests))
	buf := &bytes.Buffer{}
	for _, name := range testNames {
		// Populate protos.
		if err := r.Tests[name].toProtos(ctx, &ret, buf, name, globalTags, normPathToFullPath, testLocations); err != nil {
			return nil, errors.Annotate(err, "test %q failed to convert run fields", name).Err()
		}
	}

	return ret, nil
}

// convertTests converts the trie of tests.
func (r *JSONTestResults) convertTests(curPath string, curNode json.RawMessage) error {
	// curNode should certainly be a map.
	var maybeNode map[string]json.RawMessage
	if err := json.Unmarshal(curNode, &maybeNode); err != nil {
		return errors.Annotate(err, "%q not map[string]json.RawMessage", curNode).Err()
	}

	// Convert the tree.
	for key, value := range maybeNode {
		// Set up test id.
		delim := "/"
		testID := key
		if r.PathDelimiter != "" {
			delim = r.PathDelimiter
		}

		if curPath != "" {
			testID = fmt.Sprintf("%s%s%s", curPath, delim, key)
		} else {
			if prefixJSON, ok := r.Metadata[testNamePrefixKey]; ok {
				var prefix string
				if err := json.Unmarshal(prefixJSON, &prefix); err != nil {
					return errors.Annotate(err, "%s not string, got %q", testNamePrefixKey, prefixJSON).Err()
				}
				testID = prefix + key
			}
		}

		// Try to unmarshal value to TestFields. We check success by checking fields we expect to
		// be populated.
		maybeFields := &TestFields{}
		json.Unmarshal(value, maybeFields)
		// TODO(crbug/1034025): Check maybeFields.Expected is populated.
		if maybeFields.Actual != "" {
			if err := maybeFields.convertArtifacts(); err != nil {
				return errors.Annotate(err,
					"artifacts in %q neither map[string][]string nor map[string]string", value).Err()
			}

			if r.Tests == nil {
				r.Tests = make(map[string]*TestFields)
			}
			r.Tests[testID] = maybeFields
			continue
		}

		// Otherwise, try to process it as an intermediate node.
		if err := r.convertTests(testID, value); err != nil {
			return errors.Annotate(err, "error attempting conversion of %q as intermediate node", value).Err()
		}
	}
	return nil
}

// convertArtifacts converts the raw artifacts into the supported
// map[string][]string representation.
// It also accepts but does not store artifacts in the map[string]string
// representation used by WPT results.
// TODO(crbug/1034021): Support only map[string][]string and remove (unmarshal directly).
func (f *TestFields) convertArtifacts() error {
	for name, arts := range f.ArtifactsRaw {
		if artifactIgnorelist.Has(name) {
			continue
		}
		var asPathsErr, asStringErr error

		// Try interpreting the artifacts as both formats in turn.
		// Store only the ones that are in the expected map[string][]string format.
		var maybePaths []string
		if asPathsErr = json.Unmarshal(arts, &maybePaths); asPathsErr == nil {
			if f.Artifacts == nil {
				f.Artifacts = map[string][]string{}
			}
			f.Artifacts[name] = maybePaths
			continue
		}

		// Accept string formats, but do not store the converted artifacts, because
		// they are not actually supported by JSON Test Results.
		var maybeString string
		if asStringErr = json.Unmarshal(arts, &maybeString); asStringErr == nil {
			continue
		}

		if asPathsErr != nil || asStringErr != nil {
			return errors.Annotate(errors.NewMultiError(asStringErr, asPathsErr),
				"converting artifacts for %s", name).Err()
		}
	}

	return nil
}

// extractTags tries to read the optional "tags" field in "metadata" as a slice of strings.
func (r *JSONTestResults) extractTags() ([]string, error) {
	maybeTags, ok := r.Metadata["tags"]
	if !ok {
		return nil, nil
	}

	var tags []string
	if err := json.Unmarshal(maybeTags, &tags); err != nil {
		return nil, errors.Annotate(err, "tags not []string, got %q", maybeTags).Err()
	}

	return tags, nil
}

func fromJSONStatus(s string) (pb.TestStatus, error) {
	switch s {
	case "CRASH":
		return pb.TestStatus_CRASH, nil
	case "FAIL":
		return pb.TestStatus_FAIL, nil
	case "PASS":
		return pb.TestStatus_PASS, nil
	case "SKIP":
		return pb.TestStatus_SKIP, nil
	case "TIMEOUT":
		return pb.TestStatus_ABORT, nil

	// The below are web test-specific statuses. They are officially deprecated, but in practice
	// still generated by the tests and should be converted.
	case "IMAGE", "TEXT", "IMAGE+TEXT", "AUDIO", "LEAK", "MISSING":
		return pb.TestStatus_FAIL, nil

	default:
		// There are a number of web test-specific statuses not handled here as they are deprecated.
		return 0, errors.Reason("unknown or unexpected JSON Test Format status %s", s).Err()
	}
}

// toProtos converts the TestFields into zero or more sinkpb.TestResult and
// appends them to dest.
//
// Logs unresolved artifacts.
func (f *TestFields) toProtos(ctx context.Context, dest *[]*sinkpb.TestResult, buf *bytes.Buffer, testName string, globalTags []*pb.StringPair, normPathToFullPath map[string]string, testLocations bool) error {
	// Process statuses.
	actualStatuses := strings.Split(f.Actual, " ")

	expectedSlice := strings.Split(f.Expected, " ")
	// TODO(crbug/1034025): Remove.
	if len(expectedSlice) == 1 && expectedSlice[0] == "" {
		expectedSlice = []string{"PASS"}
	}
	expectedSet := stringset.NewFromSlice(expectedSlice...)
	expectedSet.Add("SKIP")

	// Process times.
	// Time and Times are both optional, but if Times is present, its length should match the number
	// of runs. Otherwise we have only Time as the duration of the first run.
	if len(f.Times) > 0 && len(f.Times) != len(actualStatuses) {
		return errors.Reason(
			"%d durations populated but has %d test statuses; should match",
			len(f.Times), len(actualStatuses)).Err()
	}

	var durations []float64
	if len(f.Times) > 0 {
		durations = f.Times
	} else if f.Time != 0 { // Do not set duration if it is unknown.
		durations = []float64{f.Time}
	}

	// Get artifacts.
	// We expect that if we have any artifacts, the number of runs from deriving the artifacts
	// should match the number of actual runs. Because the arts are a map from run index to
	// *sinkpb.Artifacts slice, we will not error if artifacts are missing for a run, but log a warning
	// in case the number of runs do not match each other for further investigation.
	arts := f.parseArtifacts(ctx, testName, normPathToFullPath)
	if len(arts) > 0 && len(actualStatuses) != len(arts) {
		logging.Infof(ctx,
			"Test %s generated %d statuses (%v); does not match number of runs generated from artifacts (%d)",
			testName, len(actualStatuses), actualStatuses, len(arts))
	}

	// Populate protos.
	for i, runStatus := range actualStatuses {
		status, err := fromJSONStatus(runStatus)
		if err != nil {
			return err
		}

		tags := pbutil.StringPairs(
			"json_format_status", runStatus,
			"test_name", testName,
		)
		tags = append(tags, globalTags...)

		tr := &sinkpb.TestResult{
			TestId:       testName,
			Expected:     expectedSet.Has(runStatus),
			Status:       status,
			Tags:         tags,
			TestMetadata: &pb.TestMetadata{Name: testName},
		}
		if testLocations {
			tr.TestMetadata.Location = &pb.TestLocation{
				Repo:     chromiumSrcRepo,
				FileName: testName,
			}
		}

		if container, ok := arts[i]; ok {
			buf.Reset()
			err := summaryTmpl.ExecuteTemplate(buf, "jtr", map[string]interface{}{
				"links": container.links,
			})
			if err != nil {
				return err
			}
			tr.SummaryHtml = buf.String()
			tr.Artifacts = container.artifacts
		}

		if i < len(durations) {
			tr.Duration = msToDuration(durations[i] * 1000)
		}

		*dest = append(*dest, tr)
	}

	return nil
}

type parsedArtifacts struct {
	// artifacts maps from a short name to an artifact key
	artifacts map[string]*sinkpb.Artifact
	links     map[string]string
}

// parseArtifacts parses f.Artifacts field.
//
// It tries to derive the artifacts in the following order:
//   - look for them in the outputs represented as sinkpb.Artifacts
//   - check if they're a known special case
//   - fail to process and mark them as `unresolvedArtifacts`
func (f *TestFields) parseArtifacts(ctx context.Context, testID string, normPathToFullPath map[string]string) map[int]*parsedArtifacts {
	artifacts := map[int]*parsedArtifacts{}

	for name, paths := range f.Artifacts {
		for i, path := range paths {
			// We don't use the normalized path everywhere, e.g. if "path" is actually a URL.
			normPath := normalizePath(path)

			// Get the run ID of the artifact. Defaults to 0 (i.e. assumes there is only one run).
			runID, err := artifactRunID(normPath)
			if err != nil {
				logging.Warningf(ctx, "Test %q: failed to extract artifact run id from %q", testID, normPath)
				continue
			}
			container := artifacts[runID]
			if container == nil {
				container = &parsedArtifacts{
					links:     map[string]string{},
					artifacts: map[string]*sinkpb.Artifact{},
				}
				artifacts[runID] = container
			}

			// Look for the path in outputs.
			// TODO(crbug/1032779): Track outputs that were processed.
			if fullPath := findArtifactFullPath(normPathToFullPath, normPath); fullPath != "" {
				container.artifacts[name] = &sinkpb.Artifact{
					Body: &sinkpb.Artifact_FilePath{FilePath: fullPath},
				}
				continue
			}

			// If the name is otherwise understood by ResultDB, process it.

			// Accept Gold triage links.
			if name == "gold_triage_link" || name == "triage_link_for_entire_cl" {
				// We don't expect more than one triage link per test run, but if there is more than one,
				// suffix the name with index to ensure we retain it too.
				linkName := name
				if i > 0 {
					linkName = fmt.Sprintf("%s_%d", name, i)
				}
				container.links[linkName] = path
				continue
			}

			// Otherwise, could not populate artifact, log it as unresolved.
			logging.Warningf(ctx, "Test %q: failed to resolve artifact %q:%q", testID, name, path)
		}
	}

	return artifacts
}

// artifactRunID extracts a run ID, defaulting to 0, or error if it doesn't recognize the format.
func artifactRunID(path string) (int, error) {
	if m := testRunSubdirRe.FindStringSubmatch(path); m != nil {
		return strconv.Atoi(m[1])
	}

	// No retry_<i> subdirectory, so assume it's the first/0th run.
	return 0, nil
}

// findArtifactFullPath looks for an artifact's full path.
// Checks the root directory as well as known possible subdirectories.
func findArtifactFullPath(normToFull map[string]string, normPath string) string {
	// Check root.
	if fullPath, ok := normToFull[normPath]; ok {
		return fullPath
	}

	// Check known candidate subdirectories.
	// TODO(1027708,1031296): Remove.
	for _, dir := range artifactDirectories {
		key := path.Join(dir, normPath)
		if fullPath, ok := normToFull[key]; ok {
			return fullPath
		}
	}

	return ""
}
