// Copyright 2015 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package analyzer

import (
	"errors"
	"expvar"
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"infra/monitoring/client"
	"infra/monitoring/messages"

	"github.com/luci/luci-go/common/data/stringset"
)

const (
	// StepCompletedRun is a synthetic step name used to indicate the build run is complete.
	StepCompletedRun = "completed run"

	// Order of severity, worst to least bad.
	treeCloserSev = iota
	staleMasterSev
	hungBuilderSev
	infraFailureSev
	reliableFailureSev
	newFailureSev
	staleBuilderSev
	idleBuilderSev
	offlineBuilderSev

	// Step result values.
	resOK           = float64(1)
	resInfraFailure = float64(4)
)

var (
	errLog  = log.New(os.Stderr, "", log.Lshortfile|log.Ltime)
	infoLog = log.New(os.Stdout, "", log.Lshortfile|log.Ltime)
	expvars = expvar.NewMap("analyzer")
	cpRE    = regexp.MustCompile("Cr-Commit-Position: (.*)@{#([0-9]+)}")
)

var (
	errNoBuildSteps      = errors.New("No build steps")
	errNoRecentBuilds    = errors.New("No recent builds")
	errNoCompletedBuilds = errors.New("No completed builds")
)

// StepAnalyzer reasons about a stepFailure and produces a set of reasons for the
// failure.  It also indicates whether or not it recognizes the stepFailure.
type StepAnalyzer interface {
	// Analyze returns a list of reasons for the step failure, and a boolean
	// value indicating whether or not the step failure was recognized by
	// the analyzer.
	Analyze(f stepFailure) (*StepAnalyzerResult, error)
}

// StepAnalyzerResult is the result of running analysis on a stepFailure.
type StepAnalyzerResult struct {
	// Recognized is true if the StepFailureAnalyzer recognized the stepFailure.
	Recognized bool

	// Reasons lists the reasons for the stepFailure determined by the StepFailureAnalyzer.
	Reasons []string
}

// Analyzer runs the process of checking masters, builders, test results and so on,
// in order to produce alerts.
type Analyzer struct {
	// MaxRecentBuilds is the maximum number of recent builds to check, per builder.
	MaxRecentBuilds int

	// MinRecentBuilds is the minimum number of recent builds to check, per builder.
	MinRecentBuilds int

	// StepAnalzers are the set of build step failure analyzers to be checked on
	// build step failures.
	StepAnalyzers []StepAnalyzer

	// Reader is the Reader implementation for fetching json from CBE, builds, etc.
	Reader client.Reader

	// HungBuilerThresh is the maxumum length of time a builder may be in state "building"
	// before triggering a "hung builder" alert.
	HungBuilderThresh time.Duration

	// OfflineBuilderThresh is the maximum length of time a builder may be in state "offline"
	//  before triggering an "offline builder" alert.
	OfflineBuilderThresh time.Duration

	// IdleBuilderCountThresh is the maximum number of builds a builder may have in queue
	// while in the "idle" state before triggering an "idle builder" alert.
	IdleBuilderCountThresh int64

	// StaleMasterThreshold is the maximum age that master data from CBE can be before
	// triggering a "stale master" alert.
	StaleMasterThreshold time.Duration

	// Gatekeeper is a the parsed gatekeeper.json config file.
	Gatekeeper *GatekeeperRules

	// These limit the scope analysis, useful for debugging.
	MasterOnly  string
	BuilderOnly string
	BuildOnly   int64

	// rslck protects revisionSummaries from concurrent access.
	rslck             *sync.Mutex
	revisionSummaries map[string]messages.RevisionSummary

	// Now is useful for mocking the system clock in testing and simulating time
	// during replay.
	Now func() time.Time
}

// New returns a new Analyzer. If client is nil, it assigns a default implementation.
// maxBuilds is the maximum number of builds to check, per builder.
func New(c client.Reader, minBuilds, maxBuilds int) *Analyzer {
	if c == nil {
		c = client.NewReader()
	}

	return &Analyzer{
		Reader:                 c,
		MaxRecentBuilds:        maxBuilds,
		MinRecentBuilds:        minBuilds,
		HungBuilderThresh:      3 * time.Hour,
		OfflineBuilderThresh:   90 * time.Minute,
		IdleBuilderCountThresh: 50,
		StaleMasterThreshold:   10 * time.Minute,
		StepAnalyzers: []StepAnalyzer{
			&TestFailureAnalyzer{Reader: c},
			&CompileFailureAnalyzer{Reader: c},
		},
		Gatekeeper:        &GatekeeperRules{},
		rslck:             &sync.Mutex{},
		revisionSummaries: map[string]messages.RevisionSummary{},
		Now: func() time.Time {
			return time.Now()
		},
	}
}

// MasterAlerts returns alerts generated from the master.
func (a *Analyzer) MasterAlerts(master *messages.MasterLocation, be *messages.BuildExtract) []messages.Alert {
	ret := []messages.Alert{}

	// Copied logic from builder_messages.
	// No created_timestamp should be a warning sign, no?
	if be.CreatedTimestamp == messages.EpochTime(0) {
		return ret
	}
	expvars.Add("MasterAlerts", 1)
	defer expvars.Add("MasterAlerts", -1)
	elapsed := a.Now().Sub(be.CreatedTimestamp.Time())
	if elapsed > a.StaleMasterThreshold {
		ret = append(ret, messages.Alert{
			Key:       fmt.Sprintf("stale master: %v", master),
			Title:     fmt.Sprintf("Stale %s master data", master),
			Body:      fmt.Sprintf("%dh %2dm elapsed since last update.", int(elapsed.Hours()), int(elapsed.Minutes())),
			StartTime: messages.TimeToEpochTime(be.CreatedTimestamp.Time()),
			Severity:  staleMasterSev,
			Time:      messages.TimeToEpochTime(a.Now()),
			Links:     []messages.Link{{"Master", master.URL.String()}},
			Type:      messages.AlertStaleMaster,
			// No extension for now.
		})
	}
	if elapsed < 0 {
		// Add this to the alerts returned, rather than just log it?
		errLog.Printf("Master %s timestamp is newer than current time (%s): %s old.", master, a.Now(), elapsed)
	}

	return ret
}

// BuilderAlerts returns alerts generated from builders connected to the master.
func (a *Analyzer) BuilderAlerts(tree string, master *messages.MasterLocation, be *messages.BuildExtract) []messages.Alert {

	// TODO: Collect activeBuilds from be.Slaves.RunningBuilds
	type r struct {
		builderName string
		b           messages.Builder
		alerts      []messages.Alert
		err         []error
	}
	c := make(chan r, len(be.Builders))

	// TODO: get a list of all the running builds from be.Slaves? It
	// appears to be used later on in the original py.
	scannedBuilders := []string{}
	for builderName, builder := range be.Builders {
		if a.BuilderOnly != "" && builderName != a.BuilderOnly {
			continue
		}
		scannedBuilders = append(scannedBuilders, builderName)
		go func(builderName string, b messages.Builder) {
			out := r{builderName: builderName, b: b}
			defer func() {
				c <- out
			}()

			expvars.Add("BuilderAlerts", 1)
			defer expvars.Add("BuilderAlerts", -1)
			// Each call to builderAlerts may trigger blocking json fetches,
			// but it has a data dependency on the above cache-warming call, so
			// the logic remains serial.
			out.alerts, out.err = a.builderAlerts(tree, master, builderName, &b)
		}(builderName, builder)
	}

	ret := []messages.Alert{}
	for _, builderName := range scannedBuilders {
		r := <-c
		if len(r.err) != 0 {
			// TODO: add a special alert for this too?
			errLog.Printf("Error getting alerts for builder %s: %v", builderName, r.err)
		} else {
			ret = append(ret, r.alerts...)
		}
	}

	ret = a.mergeAlertsByReason(ret)

	return ret
}

// TODO: actually write the on-disk cache.
func filenameForCacheKey(cc string) string {
	cc = strings.Replace(cc, "/", "_", -1)
	return fmt.Sprintf("/tmp/dispatcher.cache/%s", cc)
}

func alertKey(master, builder, step, testName string) string {
	return fmt.Sprintf("%s.%s.%s.%s", master, builder, step, testName)
}

// This type is used for sorting build IDs.
type buildNums []int64

func (a buildNums) Len() int           { return len(a) }
func (a buildNums) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a buildNums) Less(i, j int) bool { return a[i] > a[j] }

// latestBuildStep returns the latest build step name and update time, and an error
// if there were any errors.
func (a *Analyzer) latestBuildStep(b *messages.Build) (lastStep string, lastUpdate messages.EpochTime, err error) {
	if len(b.Steps) == 0 {
		return "", messages.TimeToEpochTime(a.Now()), errNoBuildSteps
	}
	if len(b.Times) > 1 && b.Times[1] != 0 {
		return StepCompletedRun, b.Times[1], nil
	}

	for _, step := range b.Steps {
		if len(step.Times) > 1 && step.Times[1] > lastUpdate {
			// Step is done.
			lastUpdate = step.Times[1]
			lastStep = step.Name
		} else if len(step.Times) > 0 && step.Times[0] > lastUpdate {
			// Step has started.
			lastUpdate = step.Times[0]
			lastStep = step.Name
		}
	}
	return
}

// lastBuild returns the last build (which may or may not have completed) and the last completed build (which may be
// the same as the last build), and any error that occurred while finding them.
func (a *Analyzer) lastBuilds(master *messages.MasterLocation, builderName string, recentBuildIDs []int64) (lastBuild, lastCompletedBuild *messages.Build, err error) {
	// Check for stale/idle/offline builders.  Latest build is the first in the list.

	for i, buildNum := range recentBuildIDs {
		infoLog.Printf("Checking last %s/%s build ID: %d", master.Name(), builderName, buildNum)

		var build *messages.Build
		build, err = a.Reader.Build(master, builderName, buildNum)
		if err != nil {
			return
		}

		if i == 0 {
			lastBuild = build
		}

		lastStep, _, _ := a.latestBuildStep(build)
		if lastStep == StepCompletedRun {
			lastCompletedBuild = build
			return
		}
	}

	return nil, nil, errNoCompletedBuilds
}

// TODO: also check the build slaves to see if there are alerts for currently running builds that
// haven't shown up in CBE yet.
func (a *Analyzer) builderAlerts(tree string, master *messages.MasterLocation, builderName string, b *messages.Builder) ([]messages.Alert, []error) {
	if len(b.CachedBuilds) == 0 {
		// TODO: Make an alert for this?
		return nil, []error{errNoRecentBuilds}
	}

	recentBuildIDs := b.CachedBuilds
	// Should be a *reverse* sort.
	sort.Sort(buildNums(recentBuildIDs))
	if len(recentBuildIDs) > a.MaxRecentBuilds {
		recentBuildIDs = recentBuildIDs[:a.MaxRecentBuilds]
	}

	alerts, errs := []messages.Alert{}, []error{}

	lastBuild, lastCompletedBuild, err := a.lastBuilds(master, builderName, recentBuildIDs)
	if err != nil {
		errs = append(errs, err)
		return nil, errs
	}

	if a.Gatekeeper.ExcludeBuilder(tree, master, builderName) {
		return nil, nil
	}

	// Examining only the latest build is probably suboptimal since if it's still in progress it might
	// not have hit a step that is going to fail and has failed repeatedly for the last few builds.
	// AKA "Reliable failures".  TODO: Identify "Reliable failures"
	lastStep, lastUpdated, err := a.latestBuildStep(lastBuild)
	if err != nil {
		errs = append(errs, fmt.Errorf("Couldn't get latest build step for %s.%s: %v", master.Name(), builderName, err))
		return alerts, errs
	}
	elapsed := a.Now().Sub(lastUpdated.Time())
	links := []messages.Link{
		{"Builder", client.BuilderURL(master, builderName).String()},
		{"Last build", client.BuildURL(master, builderName, lastBuild.Number).String()},
		{"Last build step", client.StepURL(master, builderName, lastStep, lastBuild.Number).String()},
	}

	switch b.State {
	case messages.StateBuilding:
		if elapsed > a.HungBuilderThresh && lastStep != StepCompletedRun {
			alerts = append(alerts, messages.Alert{
				Key:       fmt.Sprintf("%s.%s.hung", master.Name(), builderName),
				Title:     fmt.Sprintf("%s.%s is hung in step %s.", master.Name(), builderName, lastStep),
				Body:      fmt.Sprintf("%s.%s has been building for %v (last step update %s), past the alerting threshold of %v", master.Name(), builderName, elapsed, lastUpdated.Time(), a.HungBuilderThresh),
				Severity:  hungBuilderSev,
				Time:      messages.TimeToEpochTime(a.Now()),
				StartTime: messages.TimeToEpochTime(lastUpdated.Time()),
				Links:     links,
				Type:      messages.AlertHungBuilder,
			})
			// Note, just because it's building doesn't mean it's in a good state. If the last N builds
			// all failed (for some large N) then this might still be alertable.
		}
	case messages.StateOffline:
		if elapsed > a.OfflineBuilderThresh {
			alerts = append(alerts, messages.Alert{
				Key:       fmt.Sprintf("%s.%s.offline", master.Name(), builderName),
				Title:     fmt.Sprintf("%s.%s is offline.", master.Name(), builderName),
				Body:      fmt.Sprintf("%s.%s has been offline for %v (last step update %s %v), past the alerting threshold of %v", master.Name(), builderName, elapsed, lastUpdated.Time(), float64(lastUpdated), a.OfflineBuilderThresh),
				Severity:  offlineBuilderSev,
				Time:      messages.TimeToEpochTime(a.Now()),
				StartTime: messages.TimeToEpochTime(lastUpdated.Time()),
				Links:     links,
				Type:      messages.AlertOfflineBuilder,
			})
		}
	case messages.StateIdle:
		if b.PendingBuilds > a.IdleBuilderCountThresh {
			alerts = append(alerts, messages.Alert{
				Key:       fmt.Sprintf("%s.%s.idle", master.Name(), builderName),
				Title:     fmt.Sprintf("%s.%s is idle with too many pending builds.", master.Name(), builderName),
				Body:      fmt.Sprintf("%s.%s is idle with %d pending builds, past the alerting threshold of %d", master.Name(), builderName, b.PendingBuilds, a.IdleBuilderCountThresh),
				Severity:  idleBuilderSev,
				Time:      messages.TimeToEpochTime(a.Now()),
				StartTime: messages.TimeToEpochTime(lastUpdated.Time()),
				Links:     links,
				Type:      messages.AlertIdleBuilder,
			})
		}
	default:
		errLog.Printf("Unknown %s.%s builder state: %s", master.Name(), builderName, b.State)
	}

	// Check for alerts on the most recent complete build
	infoLog.Printf("Checking %d most recent builds for alertable step failures: %s/%s", len(recentBuildIDs), master.Name(), builderName)

	mostRecentComplete := 0
	for i, id := range recentBuildIDs {
		if id == lastCompletedBuild.Number {
			mostRecentComplete = i
		}
	}
	as, es := a.builderStepAlerts(tree, master, builderName, recentBuildIDs[mostRecentComplete:])
	alerts = append(alerts, as...)
	errs = append(errs, es...)

	return alerts, errs
}

// mergeAlertsByReason merges alerts for step failures occurring across multiple builders into
// one alert with multiple builders indicated.
func (a *Analyzer) mergeAlertsByReason(alerts []messages.Alert) []messages.Alert {
	mergedAlerts := []messages.Alert{}
	byReason := map[string][]messages.Alert{}
	for _, alert := range alerts {
		bf, ok := alert.Extension.(messages.BuildFailure)
		if !ok {
			infoLog.Printf("%s failed, but isn't a builder-failure: %s", alert.Key, alert.Type)
			// Not a builder failure, so don't bother trying to group it by step name.
			mergedAlerts = append(mergedAlerts, alert)
			continue
		}
		for _, r := range bf.Reasons {
			testNames := r.TestNames
			sort.Strings(testNames)
			reason := strings.Join(append([]string{r.Step}, testNames...), ",")
			byReason[reason] = append(byReason[reason], alert)
		}
	}

	sortedReasons := []string{}
	for reason := range byReason {
		sortedReasons = append(sortedReasons, reason)
	}

	sort.Strings(sortedReasons)

	// Merge build failures by step, then make a single alert listing all of the builders
	// failing for that step.
	for _, reason := range sortedReasons {
		stepAlerts := byReason[reason]
		if len(stepAlerts) == 1 {
			mergedAlerts = append(mergedAlerts, stepAlerts[0])
			continue
		}

		sort.Sort(messages.Alerts(stepAlerts))
		merged := stepAlerts[0]

		mergedBF := merged.Extension.(messages.BuildFailure)
		if len(mergedBF.Builders) > 1 {
			errLog.Printf("Alert shouldn't have multiple builders before merging by reason: %+v", reason)
		}

		// Clear out the list of builders because we're going to reconstruct it.
		mergedBF.Builders = []messages.AlertedBuilder{}
		mergedBF.RegressionRanges = []messages.RegressionRange{}

		builders := map[string]messages.AlertedBuilder{}
		regressionRanges := map[string][]messages.RegressionRange{}

		for _, alert := range stepAlerts { // stepAlerts[1:]? already have [0] in mergedBf
			bf := alert.Extension.(messages.BuildFailure)
			if len(bf.Builders) > 1 {
				errLog.Printf("Alert shouldn't have multiple builders before merging by reason: %+v", reason)
			}
			if bf.TreeCloser {
				mergedBF.TreeCloser = true
			}

			builder := bf.Builders[0]
			// If any of the builders would call it a tree closer,
			// mark the merged alert as one.
			mergedBF.TreeCloser = bf.TreeCloser || mergedBF.TreeCloser
			if ab, ok := builders[builder.Name]; ok {
				if ab.FirstFailure < builder.FirstFailure {
					builder.FirstFailure = ab.FirstFailure
				}
				if ab.LatestFailure > builder.LatestFailure {
					builder.LatestFailure = ab.LatestFailure
				}
				if ab.StartTime < builder.StartTime || builder.StartTime == 0 {
					builder.StartTime = ab.StartTime
				}
			}
			builders[builder.Name] = builder
			regressionRanges[builder.Name] = bf.RegressionRanges
		}

		builderNames := []string{}
		for name := range builders {
			builderNames = append(builderNames, name)
		}
		sort.Strings(builderNames)

		for _, name := range builderNames {
			builder := builders[name]
			mergedBF.Builders = append(mergedBF.Builders, builder)
			// Fix this so it de-dupes regression ranges, or at least dedupes the Revisions in
			// in each repo.
			mergedBF.RegressionRanges = append(mergedBF.RegressionRanges, regressionRanges[builder.Name]...)
		}

		// De-dupe regression ranges by repo.
		posByRepo := map[string][]string{}
		for _, regRange := range mergedBF.RegressionRanges {
			posByRepo[regRange.Repo] = append(posByRepo[regRange.Repo], regRange.Positions...)
		}

		mergedBF.RegressionRanges = []messages.RegressionRange{}
		for repo, pos := range posByRepo {
			mergedBF.RegressionRanges = append(mergedBF.RegressionRanges, messages.RegressionRange{
				Repo:      repo,
				Positions: uniques(pos),
			})
		}

		sort.Sort(byRepo(mergedBF.RegressionRanges))

		if len(mergedBF.Builders) > 1 {
			merged.Title = fmt.Sprintf("%s failing on %d builders", mergedBF.Reasons[0].Step, len(mergedBF.Builders))
			builderNames := []string{}
			for _, b := range mergedBF.Builders {
				builderNames = append(builderNames, b.Name)
				if b.StartTime < merged.StartTime || merged.StartTime == 0 {
					merged.StartTime = b.StartTime
				}
			}
			merged.Body = strings.Join(builderNames, ", ")
		}

		shrunkRegressionRanges := []messages.RegressionRange{}

		// Save space for long commit position lists by just keeping the first and last.
		for _, r := range mergedBF.RegressionRanges {
			if len(r.Positions) > 2 {
				r.Positions = []string{r.Positions[0], r.Positions[len(r.Positions)-1]}
			}
			shrunkRegressionRanges = append(shrunkRegressionRanges, r)
		}
		mergedBF.RegressionRanges = shrunkRegressionRanges

		merged.Extension = mergedBF
		mergedAlerts = append(mergedAlerts, merged)
	}

	return mergedAlerts
}

// GetRevisionSummaries returns a slice of RevisionSummaries for the list of hashes.
func (a *Analyzer) GetRevisionSummaries(hashes []string) ([]messages.RevisionSummary, error) {
	ret := []messages.RevisionSummary{}
	for _, h := range hashes {
		a.rslck.Lock()
		s, ok := a.revisionSummaries[h]
		a.rslck.Unlock()
		if !ok {
			return nil, fmt.Errorf("Unrecognized hash:  %s", h)
		}
		ret = append(ret, s)
	}

	return ret, nil
}

// builderStepAlerts scans the steps of recent builds done on a particular builder,
// generating an Alert for each step output that warrants one.  Alerts are then
// merged by key so that failures that occur across a range of builds produce a single
// alert instead of one for each build.
//
// We assume the first build in the list of recent build IDs is the most recent;
// this build is used to select the steps which we think are still failing and
// should show up as alerts.
func (a *Analyzer) builderStepAlerts(tree string, master *messages.MasterLocation, builderName string, recentBuildIDs []int64) (alerts []messages.Alert, errs []error) {
	if len(recentBuildIDs) == 0 {
		return nil, nil
	}

	sort.Sort(buildNums(recentBuildIDs))
	// Check for alertable step failures.  We group them by key to de-duplicate and merge values
	// once we've scanned everything.
	stepAlertsByKey := map[string][]messages.Alert{}

	latestBuild := recentBuildIDs[0]
	importantFailures, err := a.stepFailures(master, builderName, latestBuild)
	if err != nil {
		return nil, []error{err}
	}
	if len(importantFailures) == 0 {
		return nil, errs
	}

	// Get findit results
	stepNames := make([]string, len(importantFailures))
	for i, f := range importantFailures {
		stepNames[i] = f.step.Name
	}
	finditResults, err := a.Reader.Findit(master, builderName, latestBuild, stepNames)
	if err != nil {
		return nil, []error{fmt.Errorf("while getting findit results: %s", err)}
	}

	importantAlerts, err := a.stepFailureAlerts(tree, importantFailures)
	if err != nil {
		return nil, []error{err}
	}
	importantKeys := stringset.New(0)
	for _, alr := range importantAlerts {
		importantKeys.Add(alr.Key)
	}

	for _, buildNum := range recentBuildIDs {
		failures, err := a.stepFailures(master, builderName, buildNum)
		if err != nil {
			errs = append(errs, err)
		}
		if len(failures) == 0 {
			// Bail as soon as we find a successful build.
			break
		}

		as, err := a.stepFailureAlerts(tree, failures)
		if err != nil {
			errs = append(errs, err)
		}

		// Group alerts by key so they can be merged across builds/regression ranges.
		for _, alr := range as {
			if importantKeys.Has(alr.Key) {
				stepAlertsByKey[alr.Key] = append(stepAlertsByKey[alr.Key], alr)
			}
		}
	}

	// Now coalesce alerts with the same key into single alerts with merged properties.
	for key, keyedAlerts := range stepAlertsByKey {
		mergedAlert := keyedAlerts[0] // Merge everything into the first one
		mergedBF, ok := mergedAlert.Extension.(messages.BuildFailure)
		if !ok {
			errLog.Printf("Couldn't cast extension as BuildFailure: %s", mergedAlert.Type)
		}

		for _, alr := range keyedAlerts[1:] {
			if alr.Title != mergedAlert.Title {
				// Sanity checking.
				errLog.Printf("Merging alerts with same key (%q), different title: (%q vs %q)", key, alr.Title, mergedAlert.Title)
				continue
			}
			bf, ok := alr.Extension.(messages.BuildFailure)
			if !ok {
				errLog.Printf("Couldn't cast a %q extension as BuildFailure", alr.Type)
				continue
			}
			// At this point, there should only be one builder per failure because
			// alert keys include the builder name.  We merge builders by step failure
			// in another pass, after this function is called.
			if len(bf.Builders) != 1 {
				errLog.Printf("bf.Builders len is not 1: %d", len(bf.Builders))
			}
			firstBuilder := bf.Builders[0]
			mergedBuilder := mergedBF.Builders[0]

			// These failure numbers are build numbers. The UI should probably indicate
			// gnumbd sequence numbers instead or in addition to build numbers.

			if firstBuilder.FirstFailure < mergedBuilder.FirstFailure {
				mergedBuilder.FirstFailure = firstBuilder.FirstFailure
			}
			if firstBuilder.LatestFailure > mergedBuilder.LatestFailure {
				mergedBuilder.LatestFailure = firstBuilder.LatestFailure
			}
			if firstBuilder.StartTime < mergedBuilder.StartTime || mergedBuilder.StartTime == 0 {
				mergedBuilder.StartTime = firstBuilder.StartTime
			}
			mergedBF.Builders[0] = mergedBuilder
			mergedBF.RegressionRanges = append(mergedBF.RegressionRanges, bf.RegressionRanges...)
		}

		// Now merge regression ranges by repo.
		positionsByRepo := map[string][]string{}
		for _, rr := range mergedBF.RegressionRanges {
			positionsByRepo[rr.Repo] = append(positionsByRepo[rr.Repo], rr.Positions...)
		}

		mergedBF.RegressionRanges = []messages.RegressionRange{}

		for repo, pos := range positionsByRepo {
			mergedBF.RegressionRanges = append(mergedBF.RegressionRanges,
				messages.RegressionRange{
					Repo:      repo,
					Positions: uniques(pos),
				})
		}

		// Necessary for test cases to be repeatable.
		sort.Sort(byRepo(mergedBF.RegressionRanges))

		// FIXME: This is a very simplistic model, and we're throwing away a lot of the findit data.
		// This data should really be a part of the regression range data.
		for _, result := range finditResults {
			mergedBF.SuspectedCLs = append(mergedBF.SuspectedCLs, result.SuspectedCLs...)
		}

		mergedAlert.Extension = mergedBF

		for _, failingBuilder := range mergedBF.Builders {
			if failingBuilder.LatestFailure-failingBuilder.FirstFailure > 0 {
				mergedAlert.Severity = reliableFailureSev
			}
			if failingBuilder.StartTime < mergedAlert.StartTime || mergedAlert.StartTime == 0 {
				mergedAlert.StartTime = failingBuilder.StartTime
			}
		}

		alerts = append(alerts, mergedAlert)
	}

	return alerts, errs
}

func uniques(s []string) []string {
	m := map[string]bool{}
	for _, s := range s {
		m[s] = true
	}

	ret := []string{}
	for k := range m {
		ret = append(ret, k)
	}
	sort.Strings(ret)
	return ret
}

type byRepo []messages.RegressionRange

func (a byRepo) Len() int           { return len(a) }
func (a byRepo) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byRepo) Less(i, j int) bool { return a[i].Repo < a[j].Repo }

// stepFailures returns the steps that have failed recently on builder builderName.
func (a *Analyzer) stepFailures(master *messages.MasterLocation, builderName string, bID int64) ([]stepFailure, error) {

	var err error // To avoid re-scoping b in the nested conditional below with a :=.
	b, err := a.Reader.Build(master, builderName, bID)
	if err != nil || b == nil {
		errLog.Printf("Error fetching build: %v", err)
		return nil, err
	}

	ret := []stepFailure{}

	for _, s := range b.Steps {
		if !s.IsFinished || len(s.Results) == 0 {
			continue
		}
		// Because Results in the json data is a homogeneous array, the unmarshaler
		// doesn't have any type information to assert about it. We have to do
		// some ugly runtime type assertion ourselves.
		if r, ok := s.Results[0].(float64); ok {
			if r <= resOK {
				// This 0/1 check seems to be a convention or heuristic. A 0 or 1
				// result is apparently "ok", accoring to the original python code.
				continue
			}
		} else {
			errLog.Printf("Couldn't unmarshal first step result into a float64: %v", s.Results[0])
		}

		// We have a failure of some kind, so queue it up to check later.
		ret = append(ret, stepFailure{
			master:      master,
			builderName: builderName,
			build:       *b,
			step:        s,
		})
	}

	return ret, nil
}

// stepFailureAlerts returns alerts generated from step failures. It applies filtering
// logic specified in the gatekeeper config to ignore some failures.
func (a *Analyzer) stepFailureAlerts(tree string, failures []stepFailure) ([]messages.Alert, error) {
	ret := []messages.Alert{}
	type res struct {
		f   stepFailure
		a   *messages.Alert
		err error
	}

	// Might not need full capacity buffer, since some failures are ignored below.
	rs := make(chan res, len(failures))

	scannedFailures := []stepFailure{}
	for _, failure := range failures {
		// goroutine/channel because the reasonsForFailure call potentially
		// blocks on IO.
		if failure.step.Name == "steps" {
			infoLog.Printf("steps results failed, skipping ahead to actual failure: %s %s", failure.master.Name(), failure.builderName)
			continue
			// The actual breaking step will appear later.
		}

		// Check the gatekeeper configs to see if this is ignorable.
		if a.Gatekeeper.ExcludeFailure(tree, failure.master, failure.builderName, failure.step.Name) {
			continue
		}

		if len(failure.step.Results) > 0 {
			// Check results to see if it's an array of [4]
			// That's a purple failure, which should go to infra/trooper.
			if r, ok := failure.step.Results[0].(float64); ok && r == resInfraFailure {
				infoLog.Printf("INFRA FAILURE: %s/%s/%s", failure.master.Name(), failure.builderName, failure.step.Name)
				alr := messages.Alert{
					Title:     fmt.Sprintf("%s failing on %s/%s", failure.step.Name, failure.master.Name(), failure.builderName),
					Body:      "infrastructure failure",
					Type:      messages.AlertInfraFailure,
					StartTime: failure.build.Times[0],
					Time:      failure.build.Times[0],
					Severity:  infraFailureSev,
					Key:       alertKey(failure.master.Name(), failure.builderName, failure.step.Name, fmt.Sprintf("%v", failure.step.Results[1])),
					Extension: messages.BuildFailure{
						Reasons: []messages.Reason{{
							Step: failure.step.Name,
							URL:  failure.URL().String(),
						}},
						Builders: []messages.AlertedBuilder{
							{
								Name:          failure.builderName,
								URL:           client.BuilderURL(failure.master, failure.builderName).String(),
								StartTime:     failure.build.Times[0],
								FirstFailure:  failure.build.Number,
								LatestFailure: failure.build.Number,
							},
						},
						TreeCloser: a.Gatekeeper.WouldCloseTree(failure.master, failure.builderName, failure.step.Name),
					},
				}
				scannedFailures = append(scannedFailures, failure)
				rs <- res{
					f:   failure,
					a:   &alr,
					err: nil,
				}
				continue
			}
		}

		// Gets the named revision number from gnumbd metadata.
		getCommitPos := func(b messages.Build, name string) (string, bool) {
			for _, p := range b.Properties {
				if p[0] == name {
					s, ok := p[1].(string)
					return s, ok
				}
			}
			return "", false
		}

		scannedFailures = append(scannedFailures, failure)
		go func(f stepFailure) {
			expvars.Add("StepFailures", 1)
			defer expvars.Add("StepFailures", -1)
			alr := messages.Alert{
				Title:     fmt.Sprintf("%s failing on %s/%s", f.step.Name, f.master.Name(), f.builderName),
				Body:      "",
				Time:      f.build.Times[0],
				StartTime: f.build.Times[0],
				Type:      messages.AlertBuildFailure,
				Severity:  newFailureSev,
			}

			regRanges := []messages.RegressionRange{}
			revisionsByRepo := map[string][]string{}

			// Get gnumbd sequence numbers for whatever this build pulled in.
			chromiumPos, ok := getCommitPos(f.build, "got_revision_cp")
			if ok {
				regRanges = append(regRanges, messages.RegressionRange{
					Repo:      "chromium",
					Positions: []string{chromiumPos},
				})
			}

			blinkPos, ok := getCommitPos(f.build, "got_webkit_revision_cp")
			if ok {
				regRanges = append(regRanges, messages.RegressionRange{
					Repo:      "blink",
					Positions: []string{blinkPos},
				})
			}

			v8Pos, ok := getCommitPos(f.build, "got_v8_revision_cp")
			if ok {
				regRanges = append(regRanges, messages.RegressionRange{
					Repo:      "v8",
					Positions: []string{v8Pos},
				})
			}

			naclPos, ok := getCommitPos(f.build, "got_nacl_revision_cp")
			if ok {
				regRanges = append(regRanges, messages.RegressionRange{
					Repo:      "nacl",
					Positions: []string{naclPos},
				})
			}

			// TODO(seanmccullough):  Nix this? It adds a lot to the alerts json size.
			// Consider posting this information to a separate endpoint.
			for _, change := range f.build.SourceStamp.Changes {
				revisionsByRepo[change.Repository] = append(revisionsByRepo[change.Repository], change.Revision)
				// change.Revision is *not* always a git hash. Sometimes it is a position from gnumbd.
				// change.Revision is git hash or gnumbd depending on what exactly? Not obvious at this time.
				// A potential problem here is when multiple repos have overlapping gnumbd ranges.
				parts := cpRE.FindAllStringSubmatch(change.Comments, -1)
				pos, branch := "", ""
				if len(parts) > 0 {
					branch = parts[0][1]
					pos = parts[0][2]
				}
				a.rslck.Lock()
				a.revisionSummaries[change.Revision] = messages.RevisionSummary{
					GitHash:     change.Revision,
					Link:        change.Revlink,
					Description: trunc(change.Comments),
					Author:      change.Who,
					When:        change.When,
					Position:    pos,
					Branch:      branch,
				}
				a.rslck.Unlock()
			}

			for repo, revisions := range revisionsByRepo {
				regRanges = append(regRanges, messages.RegressionRange{
					Repo:      repo,
					Revisions: revisions,
				})
			}

			// If the builder has been failing on the same step for multiple builds in a row,
			// we should have only one alert but indicate the range of builds affected.
			// These are set in FirstFailure and LastFailure.
			bf := messages.BuildFailure{
				Builders: []messages.AlertedBuilder{
					{
						Name:          f.builderName,
						URL:           client.BuilderURL(f.master, f.builderName).String(),
						StartTime:     f.build.Times[0],
						FirstFailure:  f.build.Number,
						LatestFailure: f.build.Number,
					},
				},
				TreeCloser:       a.Gatekeeper.WouldCloseTree(f.master, f.builderName, f.step.Name),
				RegressionRanges: regRanges,
			}

			if bf.TreeCloser {
				alr.Severity = treeCloserSev
			}

			reasons := a.reasonsForFailure(f)

			bf.Reasons = append(bf.Reasons, messages.Reason{
				TestNames: reasons,
				Step:      f.step.Name,
				URL:       f.URL().String(),
			})

			alr.Key = alertKey(f.master.Name(), f.builderName, f.step.Name, "")
			alr.Extension = bf

			rs <- res{
				f:   f,
				a:   &alr,
				err: nil,
			}
		}(failure)
	}

	for range scannedFailures {
		r := <-rs
		if r.a != nil {
			ret = append(ret, *r.a)
		}
	}

	return ret, nil
}

func trunc(s string) string {
	if len(s) < 100 {
		return s
	}
	return s[:100]
}

// reasonsForFailure examines the step failure and applies some heuristics to
// to find the cause. It may make blocking IO calls in the process.
func (a *Analyzer) reasonsForFailure(f stepFailure) []string {
	ret := []string(nil)
	recognized := false
	infoLog.Printf("Checking for reasons for failure step: %v", f.step.Name)
	for _, sfa := range a.StepAnalyzers {
		res, err := sfa.Analyze(f)
		if err != nil {
			// TODO: return something that contains errors *and* reasons.
			errLog.Printf("Error getting reasons from StepAnalyzer: %v", err)
			continue
		}
		if res.Recognized {
			recognized = true
			ret = append(ret, res.Reasons...)
		}
	}

	if !recognized {
		// TODO: log and report frequently encountered unrecognized builder step
		// failure names.
		errLog.Printf("Unrecognized step step failure type, unable to find reasons: %s", f.step.Name)
	}

	return ret
}

// unexpected returns the set of expected xor actual.
func unexpected(expected, actual []string) []string {
	e, a := make(map[string]bool), make(map[string]bool)
	for _, s := range expected {
		e[s] = true
	}
	for _, s := range actual {
		a[s] = true
	}

	ret := []string{}
	for k := range e {
		if !a[k] {
			ret = append(ret, k)
		}
	}

	for k := range a {
		if !e[k] {
			ret = append(ret, k)
		}
	}

	return ret
}

type stepFailure struct {
	master      *messages.MasterLocation
	builderName string
	build       messages.Build
	step        messages.Step
}

// URL returns a url to builder step failure page.
func (f stepFailure) URL() *url.URL {
	return client.StepURL(f.master, f.builderName, f.step.Name, f.build.Number)
}
