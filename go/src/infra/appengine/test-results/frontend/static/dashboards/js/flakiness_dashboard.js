// Copyright (C) 2012 Google Inc. All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//     * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//     * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

//////////////////////////////////////////////////////////////////////////////
// CONSTANTS
//////////////////////////////////////////////////////////////////////////////
var FORWARD = 'forward';
var BACKWARD = 'backward';
var TEST_URL_BASE_PATH_FOR_BROWSING =
    'https://chromium.googlesource.com/' +
    'chromium/src/+/master/third_party/WebKit/LayoutTests/';
var TEST_URL_BASE_PATH_FOR_XHR = 'http://src.chromium.org/blink/trunk/LayoutTests/';
var TEST_RESULTS_BASE_PATH = 'https://storage.googleapis.com/chromium-layout-test-archives/';
var GPU_RESULTS_BASE_PATH = 'http://chromium-browser-gpu-tests.commondatastorage.googleapis.com/runs/'

var RELEASE_TIMEOUT = 6;
var DEBUG_TIMEOUT = 12;
var SLOW_MULTIPLIER = 5;

// FIXME: Figure out how to make this not be hard-coded.
// Probably just include in the results.json files and get it from there.
var VIRTUAL_SUITES = {
    'virtual/gpu/fast/canvas': 'fast/canvas',
    'virtual/gpu/canvas/philip': 'canvas/philip',
    'virtual/threaded/compositing/visibility': 'compositing/visibility',
    'virtual/threaded/compositing/webgl': 'compositing/webgl',
    'virtual/gpu/fast/hidpi': 'fast/hidpi',
    'virtual/deferred/fast/images': 'fast/images',
    'virtual/gpu/compositedscrolling/overflow': 'compositing/overflow',
    'virtual/gpu/compositedscrolling/scrollbars': 'scrollbars',
};

var ACTUAL_RESULT_SUFFIXES = ['expected.txt', 'expected.png', 'actual.txt', 'actual.png', 'diff.txt', 'diff.png', 'wdiff.html', 'crash-log.txt'];

var EXPECTATIONS_ORDER = ACTUAL_RESULT_SUFFIXES.filter(function(suffix) {
    return !string.endsWith(suffix, 'png');
}).map(function(suffix) {
    return suffix.split('.')[0]
});

var resourceLoader;

function generatePage(historyInstance)
{
    if (historyInstance.crossDashboardState.useTestData)
        return;

    resourceLoader.showErrors();

    var currentBuilders = _currentBuilders();

    // tests expands to all tests that match the CSV list.
    // result expands to all tests that ever have the given result
    if (historyInstance.dashboardSpecificState.tests || historyInstance.dashboardSpecificState.result) {
        generatePageForIndividualTests(individualTests());
    } else if (currentBuilders.length) {
        generatePageForBuilder(builders.builderFromKey(historyInstance.dashboardSpecificState.builder) || currentBuilders[0]);
        currentBuilders.forEach(function(builder) {
            processTestResultsForBuilderAsync(builder);
        });
    } else {
        appendHTML(htmlForNavBar() + '<br><strong>No builders found for the ' +
            'given test type. Please select another.</strong>');
        hideLoadingUI();
    }

    postHeightChangedMessage();
}

function handleValidHashParameter(historyInstance, key, value)
{
    switch(key) {
    case 'result':
    case 'tests':
        history.validateParameter(historyInstance.dashboardSpecificState, key, value,
            function() {
                return string.isValidName(value);
            });
        return true;

    case 'builder':
        history.validateParameter(historyInstance.dashboardSpecificState, key, value,
            function() {
                var currentBuilders = _currentBuilders();
                for (var c = 0; c < currentBuilders.length; c++) {
                    if (currentBuilders[c].key() == value)
                        return true;
                }
                return false;
            });

        return true;

    case 'sortColumn':
        history.validateParameter(historyInstance.dashboardSpecificState, key, value,
            function() {
                // Get all possible headers since the actual used set of headers
                // depends on the values in historyInstance.dashboardSpecificState, which are currently being set.
                var getAllTableHeaders = true;
                var headers = tableHeaders(getAllTableHeaders);
                for (var i = 0; i < headers.length; i++) {
                    if (value == sortColumnFromTableHeader(headers[i]))
                        return true;
                }
                return value == 'test' || value == 'builder';
            });
        return true;

    case 'sortOrder':
        history.validateParameter(historyInstance.dashboardSpecificState, key, value,
            function() {
                return value == FORWARD || value == BACKWARD;
            });
        return true;

    case 'resultsHeight':
    case 'revision':
        history.validateParameter(historyInstance.dashboardSpecificState, key, Number(value),
            function() {
                return value.match(/^\d+$/);
            });
        return true;

    case 'showChrome':
    case 'showExpectations':
    case 'showFlaky':
    case 'showLargeExpectations':
    case 'showNonFlaky':
    case 'showSlow':
    case 'showSkip':
    case 'showUnexpectedPasses':
    case 'showWontFix':
        historyInstance.dashboardSpecificState[key] = value == 'true';
        return true;

    default:
        return false;
    }
}

// @param {Object} params New or modified query parameters as key: value.
function handleQueryParameterChange(historyInstance, params)
{
    for (key in params) {
        if (key == 'tests') {
            // Entering cross-builder view, only keep valid keys for that view.
            for (var currentKey in historyInstance.dashboardSpecificState) {
              if (isInvalidKeyForCrossBuilderView(currentKey)) {
                delete historyInstance.dashboardSpecificState[currentKey];
              }
            }
        } else if (isInvalidKeyForCrossBuilderView(key)) {
            delete historyInstance.dashboardSpecificState.tests;
            delete historyInstance.dashboardSpecificState.result;
        }
    }

    return true;
}

var defaultDashboardSpecificStateValues = {
    sortOrder: BACKWARD,
    sortColumn: 'flakiness',
    showExpectations: false,
    // FIXME: Show flaky tests by default if you have a builder picked.
    // Ideally, we'd fix the dashboard to not pick a default builder and have
    // you pick one. In the interim, this is a good way to make the default
    // page load faster since we don't need to generate/layout a large table.
    showFlaky: false,
    showLargeExpectations: false,
    showChrome: true,
    showWontFix: false,
    showNonFlaky: false,
    showSkip: false,
    showUnexpectedPasses: false,
    resultsHeight: 300,
    revision: null,
    tests: '',
    result: '',
    builder: null
};

var DB_SPECIFIC_INVALIDATING_PARAMETERS = {
    'tests' : 'builder',
    'testType': 'builder'
};

var flakinessConfig = {
    defaultStateValues: defaultDashboardSpecificStateValues,
    generatePage: generatePage,
    handleValidHashParameter: handleValidHashParameter,
    handleQueryParameterChange: handleQueryParameterChange,
    invalidatingHashParameters: DB_SPECIFIC_INVALIDATING_PARAMETERS
};

// FIXME(jparent): Eventually remove all usage of global history object.
var g_history = new history.History(flakinessConfig);
g_history.parseCrossDashboardParameters();

//////////////////////////////////////////////////////////////////////////////
// GLOBALS
//////////////////////////////////////////////////////////////////////////////

var g_perBuilderFailures = {};
// Maps test path to an array of {builder, testResults} objects.
var g_testToResultsMap = {};

function createResultsObjectForTest(test, builder)
{
    return {
        test: test,
        builder: builder,
        // HTML for display of the results in the flakiness column
        html: '',
        flipCount: 0,
        slowestTime: 0,
        isFlaky: false,
        bugs: [],
        expectations : '',
        rawResults: '',
        // List of all the results the test actually has.
        actualResults: []
    };
}

var TestTrie = function(builders, resultsByBuilder)
{
    this._trie = {};
    builders.forEach(function(builder) {
        if (!resultsByBuilder[builder.key()]) {
            console.warn("No results for builder: ", builder.key());
            return;
        }
        var testsForBuilder = resultsByBuilder[builder.key()].tests;
        for (var test in testsForBuilder)
            this._addTest(test.split('/'), this._trie);
    }.bind(this));
}

TestTrie.prototype.forEach = function(callback, startingTriePath)
{
    var testsTrie = this._trie;
    if (startingTriePath) {
        var splitPath = startingTriePath.split('/');
        while (splitPath.length && testsTrie)
            testsTrie = testsTrie[splitPath.shift()];
    }

    if (!testsTrie)
        return;

    function traverse(trie, triePath) {
        if (trie == true)
            callback(triePath);
        else {
            for (var member in trie)
                traverse(trie[member], triePath ? triePath + '/' + member : member);
        }
    }
    traverse(testsTrie, startingTriePath);
}

TestTrie.prototype._addTest = function(test, trie)
{
    var rootComponent = test.shift();
    if (!test.length) {
        if (!trie[rootComponent])
            trie[rootComponent] = true;
        return;
    }

    if (!trie[rootComponent] || trie[rootComponent] == true)
        trie[rootComponent] = {};
    this._addTest(test, trie[rootComponent]);
}

// Map of all tests to true values. This is just so we can have the list of
// all tests across all the builders.
var g_allTestsTrie;

function getAllTestsTrie()
{
    if (!g_allTestsTrie)
        g_allTestsTrie = new TestTrie(_currentBuilders(), g_resultsByBuilder);

    return g_allTestsTrie;
}

// Returns an array of tests to be displayed in the individual tests view.
// Note that a directory can be listed as a test, so we expand that into all
// tests in the directory.
function individualTests()
{
    if (g_history.dashboardSpecificState.result)
        return allTestsWithResult(g_history.dashboardSpecificState.result);

    if (!g_history.dashboardSpecificState.tests)
        return [];

    return individualTestsForSubstringList();
}

function splitTestList()
{
    // Convert windows slashes to unix slashes and spaces/newlines to commas.
    var tests = g_history.dashboardSpecificState.tests.replace(/\\/g, '/').replace('\n', ' ').replace(/\s+/g, ',');
    return tests.split(',');
}

function individualTestsForSubstringList()
{
    var testList = splitTestList();
    // If listing a lot of tests, assume you've passed in an explicit list of tests
    // instead of patterns to match against. The matching code below is super slow.
    //
    // Also, when showChrome is false, we're embedding the dashboard elsewhere and
    // an explicit test list is passed in. In that case, we don't want
    // a search for fast/canvas/foo.html to also show virtual/gpu/fast/canvas/foo.html.
    if (testList.length > 10 || !g_history.dashboardSpecificState.showChrome)
        return testList;

    // Put the tests into an object first and then move them into an array
    // as a way of deduping.
    var testsMap = {};
    for (var i = 0; i < testList.length; i++) {
        var path = testList[i];

        // Ignore whitespace entries as they'd match every test.
        if (path.match(/^\s*$/))
            continue;

        var hasAnyMatches = false;
        getAllTestsTrie().forEach(function(triePath) {
            if (string.caseInsensitiveContains(triePath, path)) {
                testsMap[triePath] = 1;
                hasAnyMatches = true;
            }
        });

        // If a path doesn't match any tests, then assume it's a full path
        // to a test that passes on all builders.
        if (!hasAnyMatches)
            testsMap[path] = 1;
    }

    var testsArray = [];
    for (var test in testsMap)
        testsArray.push(test);

    return testsArray;
}

function allTestsWithResult(result)
{
    processTestRunsForAllBuilders();
    var retVal = [];

    getAllTestsTrie().forEach(function(triePath) {
        for (var i = 0; i < g_testToResultsMap[triePath].length; i++) {
            if (g_testToResultsMap[triePath][i].actualResults.indexOf(result.toUpperCase()) != -1) {
                retVal.push(triePath);
                break;
            }
        }
    });

    return retVal;
}

function processTestResultsForBuilderAsync(builder)
{
    setTimeout(function() { processTestRunsForBuilder(builder); }, 0);
}

function processTestRunsForAllBuilders()
{
    _currentBuilders().forEach(function(builder) {
        processTestRunsForBuilder(builder);
    });
}

function processTestRunsForBuilder(builder)
{
    if (g_perBuilderFailures[builder.key()])
      return;

    if (!g_resultsByBuilder[builder.key()]) {
        console.error('No tests found for ' + builder.key());
        g_perBuilderFailures[builder.key()] = [];
        return;
    }

    var failures = [];
    var allTestsForThisBuilder = g_resultsByBuilder[builder.key()].tests;

    for (var test in allTestsForThisBuilder) {
        var resultsForTest = createResultsObjectForTest(test, builder);

        var rawTest = g_resultsByBuilder[builder.key()].tests[test];
        resultsForTest.rawTimes = rawTest.times;
        var rawResults = rawTest.results;
        resultsForTest.rawResults = rawResults;

        if (rawTest.expected)
            resultsForTest.expectations = rawTest.expected;

        if (rawTest.bugs)
            resultsForTest.bugs = rawTest.bugs;

        var failureMap = g_resultsByBuilder[builder.key()][results.FAILURE_MAP];
        // FIXME: Switch to resultsByBuild
        var times = resultsForTest.rawTimes;
        var numTimesSeen = 0;
        var numResultsSeen = 0;
        var resultsIndex = 0;
        var resultsMap = {}

        for (var i = 0; i < times.length; i++) {
            numTimesSeen += times[i][results.RLE.LENGTH];

            while (rawResults[resultsIndex] && numTimesSeen > (numResultsSeen + rawResults[resultsIndex][results.RLE.LENGTH])) {
                numResultsSeen += rawResults[resultsIndex][results.RLE.LENGTH];
                resultsIndex++;
            }

            if (rawResults && rawResults[resultsIndex]) {
                var allResults = rawResults[resultsIndex][results.RLE.VALUE];
                for (var result = 0; result < allResults.length; result++) {
                  resultsMap[failureMap[allResults[result]]] = true;
                }
            }

            resultsForTest.slowestTime = Math.max(resultsForTest.slowestTime, times[i][results.RLE.VALUE]);
        }

        resultsForTest.actualResults = Object.keys(resultsMap);

        results.determineFlakiness(failureMap, rawResults, resultsForTest);
        failures.push(resultsForTest);

        if (!g_testToResultsMap[test])
            g_testToResultsMap[test] = [];
        g_testToResultsMap[test].push(resultsForTest);
    }

    g_perBuilderFailures[builder.key()] = failures;
}

function linkHTMLToOpenWindow(url, text)
{
    return '<a href="' + url + '" target="_blank">' + text + '</a>';
}

// Returns whether the result for index'th result for testName on builder was
// a failure.
function isFailure(builder, testName, index)
{
    var currentIndex = 0;
    var rawResults = g_resultsByBuilder[builder.key()].tests[testName].results;
    var failureMap = g_resultsByBuilder[builder.key()][results.FAILURE_MAP];
    for (var i = 0; i < rawResults.length; i++) {
        currentIndex += rawResults[i][results.RLE.LENGTH];
        if (currentIndex > index)
            return results.isFailingResult(failureMap, rawResults[i][results.RLE.VALUE]);
    }
    console.error('Index exceeds number of results: ' + index);
}

// Returns an array of indexes for all builds where this test failed.
function indexesForFailures(builder, testName)
{
    var rawResults = g_resultsByBuilder[builder.key()].tests[testName].results;
    var buildNumbers = g_resultsByBuilder[builder.key()].buildNumbers;
    var failureMap = g_resultsByBuilder[builder.key()][results.FAILURE_MAP];
    var index = 0;
    var failures = [];
    for (var i = 0; i < rawResults.length; i++) {
        var numResults = rawResults[i][results.RLE.LENGTH];
        if (results.isFailingResult(failureMap, rawResults[i][results.RLE.VALUE])) {
            for (var j = 0; j < numResults; j++)
                failures.push(index + j);
        }
        index += numResults;
    }
    return failures;
}

// Returns the path to the failure log for this non-webkit test.
function pathToFailureLog(testName)
{
    return '/steps/' + g_history.crossDashboardState.testType + '/logs/' + testName.split('.')[1]
}

function showPopup(event)
{
    if (event.target.classList.contains("interpolatedResult")) {
        showPopupForInterpolatedResult(event, event.target.attributes.rev.value);
    } else if (event.target.classList.contains("results")) {
        var tableData = event.target.parentElement;
        var tableRow = tableData.parentElement;
        var builder = tableRow.attributes.builder.value;
        var cells = tableData.children;
        var buildIndex  = -1;
        for (var i = 0; i < cells.length; ++i) {
            if (cells[i].classList.contains("interpolatedResult"))
                continue;
            ++buildIndex;
            if (cells[i] === event.target)
                break;
        }
        var test = tableRow.attributes.test.value;
        showPopupForBuild(event, builder, buildIndex, test);
    }
}

function htmlForPopupForBuild(builderKey, index, opt_testName)
{
    var html = '';
    var builder = builders.builderFromKey(builderKey);

    var time = g_resultsByBuilder[builder.key()].secondsSinceEpoch[index];
    if (time) {
        var date = new Date(time * 1000);
        html += date.toLocaleDateString() + ' ' +
                date.toLocaleTimeString(undefined, {timeZoneName: 'short'});
    }

    var buildNumber = g_resultsByBuilder[builder.key()].buildNumbers[index];
    var master = builder.master();
    var buildBasePath = master.logPath(builder.builderName, buildNumber);

    html += '<ul><li>' + linkHTMLToOpenWindow(buildBasePath, 'Build log') + '</li>';

    var chromiumLink = ui.html.chromiumRevisionLink(g_resultsByBuilder[builder.key()], index);
    if (chromiumLink)
        html += '<li>Chromium: ' + chromiumLink + '</li>';

    var chromeRevision = g_resultsByBuilder[builder.key()].chromeRevision[index];
    if (chromeRevision && g_history.isLayoutTestResults()) {
        html += '<li><a href="' + TEST_RESULTS_BASE_PATH + builder.builderNameForPath() +
            '/' + buildNumber + '/layout-test-results.zip">layout-test-results.zip</a></li>';
    }

    if (!g_history.isLayoutTestResults() && opt_testName && isFailure(builder, opt_testName, index))
        html += '<li>' + linkHTMLToOpenWindow(buildBasePath + pathToFailureLog(opt_testName), 'Failure log') + '</li>';

    html += '</ul>';
    return html;
}

function showPopupForBuild(e, builderKey, index, opt_testName)
{
    ui.popup.show(e.target, htmlForPopupForBuild(builderKey, index, opt_testName));
}

function showPopupForInterpolatedResult(e, revision)
{
    var html = 'This bot did not run the test for r' + revision + '.';
    ui.popup.show(e.target, html);
}

function classNameForFailureString(failure)
{
    return failure.replace(/(\+|\ )/, '');
}

function htmlForTestResults(test, revisions)
{
    var html = '<td class="results-container">';
    var testResults = test.rawResults.concat();
    var times = test.rawTimes.concat();
    var builder = test.builder;
    var buildNumbers = g_resultsByBuilder[builder.key()].buildNumbers;

    var cells = [];
    if (revisions) {
        for (var index = 0; index < revisions.length; index++)
            cells.push({revision: revisions[index]});
    }

    var indexToReplaceCurrentResult = -1;
    var indexToReplaceCurrentTime = -1;
    for (var i = 0; i < buildNumbers.length; i++) {
        var currentResultArray, currentTimeArray, currentTime;
        var resultString, detailedResultString;

        if (i > indexToReplaceCurrentResult) {
            currentResultArray = testResults.shift();
            if (currentResultArray) {
                var encodedResults = currentResultArray[results.RLE.VALUE];
                if (encodedResults.length > 1) {
                  resultString = 'FLAKY';
                  var failureMap = g_resultsByBuilder[builder.key()][results.FAILURE_MAP];
                  detailedResultString = '';
                  for (var j = 0; j < encodedResults.length; j++) {
                    detailedResultString += failureMap[encodedResults[j]] + ' ';
                  }
                  // We treat three attempts as a special case since many test
                  // runners run the first attempt with other tests in parallel,
                  // resulting in occassional races. Tests *should* be resilient
                  // to that, so we label them as FLAKY above. If tests still
                  // fail once when run sequentially afterward it's noteworthy
                  // and we label them VERYFLAKY.
                  if (encodedResults.length > 2) {
                    resultString = 'VERYFLAKY';
                  }
                } else {
                  resultString = g_resultsByBuilder[builder.key()][results.FAILURE_MAP][encodedResults];
                  detailedResultString = resultString;
                }
                indexToReplaceCurrentResult += currentResultArray[results.RLE.LENGTH];
            } else {
                resultString = results.NO_DATA;
                detailedResultString = resultString;
                indexToReplaceCurrentResult += buildNumbers.length;
            }
        }

        if (i > indexToReplaceCurrentTime) {
            currentTimeArray = times.shift();
            currentTime = 0;
            if (currentResultArray) {
              currentTime = currentTimeArray[results.RLE.VALUE];
              indexToReplaceCurrentTime += currentTimeArray[results.RLE.LENGTH];
            } else
              indexToReplaceCurrentTime += buildNumbers.length;
        }

        var rawRevision = g_resultsByBuilder[builder.key()][results.CHROME_REVISIONS][i];
        var revision = Number(rawRevision) || rawRevision;

        // Locate the empty cell corresponding to this revision.
        var cell;

        if (revisions) {
            for (var index = 0; index < cells.length; index++) {
                if (cells[index].revision == revision && !cells[index].html) {
                    cell = cells[index];
                    break;
                }
            }
        } else {
            cell = {revision: revision};
            cells.push(cell);
        }

        if (!cell) {
            console.error(
                'Could not determine test results column for r' + revision);
            continue;
        }

        cell.className = classNameForFailureString(resultString);
        cell.hasResult =
            resultString !== results.NO_DATA && resultString !== results.NOTRUN;
        cell.html = '<div'
            + ' title="' + detailedResultString + '. Click for more info."'
            + ' class="results ' + cell.className + '"'
            + '>'
            + (currentTime || (cell.hasResult ? '' : '?')) + '</div>';
    }

    populateEmptyCells(cells);

    for (var index = 0; index < cells.length; index++)
        html += cells[index].html;

    return html + "</td>";
}

// Fill in cells where no test data is available with a question mark.
// If the next and previous runs were equal, an "interpolated" result is shown.
function populateEmptyCells(cells)
{
    for (var index = 1; index < cells.length; index++) {
        var prevCell = cells[index - 1];
        cells[index].nextClassName =
            prevCell.hasResult ? prevCell.className : prevCell.nextClassName;
    }

    for (var index = cells.length - 2; index >= 0; --index) {
        var nextCell = cells[index + 1];
        cells[index].prevClassName =
            nextCell.hasResult ? nextCell.className : nextCell.prevClassName;
    }

    cells.forEach(function(cell) {
        if (!cell.html) {
            var interpolatedResult =
                cell.nextClassName && cell.nextClassName == cell.prevClassName;
            cell.className = interpolatedResult ? cell.nextClassName : "NODATA";
            cell.html = '<div'
                + ' title="Unknown result. Did not run tests."'
                + ' rev="' + cell.revision + '"'
                + ' class="results interpolatedResult ' + cell.className + '"'
                + '>?</div>';
        }
    });
}

function shouldShowTest(testResult)
{
    if (!g_history.isLayoutTestResults())
        return true;

    if (testResult.expectations == 'WONTFIX')
        return g_history.dashboardSpecificState.showWontFix;

    if (testResult.expectations == results.SKIP)
        return g_history.dashboardSpecificState.showSkip;

    if (testResult.isFlaky)
        return g_history.dashboardSpecificState.showFlaky;

    return g_history.dashboardSpecificState.showNonFlaky;
}

function createBugHTML(test)
{
    var symptom = test.isFlaky ? 'flaky' : 'failing';
    var title = encodeURIComponent('Layout Test ' + test.test + ' is ' + symptom);
    var description = encodeURIComponent('The following layout test is ' + symptom + ' on ' +
        '[insert platform]\n\n' + test.test + '\n\nProbable cause:\n\n' +
        '[insert probable cause]');

    url = 'https://code.google.com/p/chromium/issues/entry?template=Layout%20Test%20Failure&summary=' + title + '&comment=' + description;
    return '<a class="file-new-bug" href="' + url + '">File new bug</a>';
}

function isCrossBuilderView()
{
    return g_history.dashboardSpecificState.tests || g_history.dashboardSpecificState.result;
}

function tableHeaders(opt_getAll)
{
    var headers = [];
    if (isCrossBuilderView() || opt_getAll)
        headers.push('builder');

    if (!isCrossBuilderView() || opt_getAll)
        headers.push('test');

    if (g_history.isLayoutTestResults() || opt_getAll)
        headers.push('bugs', 'expectations');

    headers.push('slowest run', 'flakiness (numbers are runtimes in seconds)');
    return headers;
}

function linkifyBugs(bugs)
{
    var html = '';
    bugs.forEach(function(bug) {
        var bugHtml;
        if (string.startsWith(bug, 'Bug('))
            bugHtml = bug;
        else
            bugHtml = '<a href="http://' + bug + '">' + bug + '</a>';
        html += '<div>' + bugHtml + '</div>'
    });
    return html;
}

function htmlForSingleTestRow(test, showBuilderNames, revisions, isFyi)
{
    var headers = tableHeaders();
    var html = '';
    for (var i = 0; i < headers.length; i++) {
        var header = headers[i];
        var rowClass = isFyi ? ' class="fyi"' : '';
        if (string.startsWith(header, 'test') || string.startsWith(header, 'builder')) {
            var testCellClassName = 'test-link' + (showBuilderNames ? ' builder-name' : '');
            var testCellHTML = showBuilderNames ? test.builder.builderName : '<span class="link" onclick="g_history.setQueryParameter(\'tests\',\'' + test.test +'\');">' + test.test + '</span>';
            html += '<tr builder="' + test.builder.key() + '" test="' + test.test + '"' + rowClass + '>';
            html += '<td class="' + testCellClassName + '">' + testCellHTML;
        } else if (string.startsWith(header, 'bugs'))
            // FIXME: linkify bugs.
            html += '<td class=options-container>' + (linkifyBugs(test.bugs) || createBugHTML(test));
        else if (string.startsWith(header, 'expectations'))
            html += '<td class=options-container>' + test.expectations;
        else if (string.startsWith(header, 'slowest'))
            html += '<td>' + (test.slowestTime ? test.slowestTime + 's' : '');
        else if (string.startsWith(header, 'flakiness'))
            html += htmlForTestResults(test, revisions);
    }
    return html;
}

function sortColumnFromTableHeader(headerText)
{
    return headerText.split(' ', 1)[0];
}

function htmlForTableColumnHeader(headerName, opt_fillColSpan)
{
    // Use the first word of the header title as the sortkey
    var thisSortValue = sortColumnFromTableHeader(headerName);
    var arrowHTML = thisSortValue == g_history.dashboardSpecificState.sortColumn ?
        '<span class=' + g_history.dashboardSpecificState.sortOrder + '>' + (g_history.dashboardSpecificState.sortOrder == FORWARD ? '&uarr;' : '&darr;' ) + '</span>' : '';
    return '<th sortValue=' + thisSortValue +
        // Extend last th through all the rest of the columns.
        (opt_fillColSpan ? ' colspan=10000' : '') +
        // Extra span here is so flex boxing actually centers.
        // There's probably a better way to do this with CSS only though.
        '><div class=table-header-content><span></span>' + arrowHTML +
        '<span class=header-text>' + headerName + '</span>' + arrowHTML + '</div></th>';
}

function htmlForTestTable(rowsHTML, opt_excludeHeaders)
{
    var html = '<table class=test-table onclick="showPopup(event)">';
    if (!opt_excludeHeaders) {
        html += '<thead><tr>';
        var headers = tableHeaders();
        for (var i = 0; i < headers.length; i++)
            html += htmlForTableColumnHeader(headers[i], i == headers.length - 1);
        html += '</tr></thead>';
    }
    return html + '<tbody>' + rowsHTML + '</tbody></table>';
}

function appendHTML(html)
{
    // InnerHTML to a div that's not in the document. This is
    // ~300ms faster in Safari 4 and Chrome 4 on mac.
    var div = document.createElement('div');
    div.innerHTML = html;
    document.body.appendChild(div);
    postHeightChangedMessage();
}

function alphanumericCompare(column, reverse)
{
    return reversibleCompareFunction(function(a, b) {
        // Put null entries at the bottom
        var a = a[column] ? String(a[column]) : 'z';
        var b = b[column] ? String(b[column]) : 'z';

        if (a < b)
            return -1;
        else if (a == b)
            return 0;
        else
            return 1;
    }, reverse);
}

function numericSort(column, reverse)
{
    return reversibleCompareFunction(function(a, b) {
        a = parseFloat(a[column]);
        b = parseFloat(b[column]);
        return a - b;
    }, reverse);
}

function reversibleCompareFunction(compare, reverse)
{
    return function(a, b) {
        return compare(reverse ? b : a, reverse ? a : b);
    };
}

function changeSort(e)
{
    var target = e.currentTarget;
    e.preventDefault();

    var sortValue = target.getAttribute('sortValue');
    while (target && target.tagName != 'TABLE')
        target = target.parentNode;

    var sort = 'sortColumn';
    var orderKey = 'sortOrder';
    if (sortValue == g_history.dashboardSpecificState[sort] && g_history.dashboardSpecificState[orderKey] == FORWARD)
        order = BACKWARD;
    else
        order = FORWARD;

    g_history.setQueryParameter(sort, sortValue, orderKey, order);
}

function sortTests(tests, column, order)
{
    var resultsProperty, sortFunctionGetter;
    if (column == 'flakiness') {
        sortFunctionGetter = numericSort;
        resultsProperty = 'flipCount';
    } else if (column == 'slowest') {
        sortFunctionGetter = numericSort;
        resultsProperty = 'slowestTime';
    } else {
        sortFunctionGetter = alphanumericCompare;
        resultsProperty = column;
    }

    tests.sort(sortFunctionGetter(resultsProperty, order == BACKWARD));
}

// Return an array of revisions across all builders such that all of a builder's
// builds have a unique entry in the list.
function collapsedRevisionList(testResults)
{
    var revisionsCountedSet = {};
    for (var resultIndex = 0; resultIndex < testResults.length; resultIndex++) {
        var builder = testResults[resultIndex].builder;
        var build = g_resultsByBuilder[builder.key()];
        var buildNumbers = build.buildNumbers;
        var builderRevisionsCountedSet = {};
        for (var i = 0; i < buildNumbers.length; i++) {
            var revision = build[results.CHROME_REVISIONS][i];

            if (!Number(revision))
                return null;

            builderRevisionsCountedSet[revision] =
                (builderRevisionsCountedSet[revision] || 0) + 1;
        }

        // Join builder's revisions with the total revisions for all builders.
        for (var revision in builderRevisionsCountedSet)
            revisionsCountedSet[revision] = Math.max(revisionsCountedSet[revision] || 0, builderRevisionsCountedSet[revision]);
    }

    var revisionsArray = [];
    for (var revision in revisionsCountedSet) {
        for (var i = revisionsCountedSet[revision] - 1; i >= 0; --i)
            revisionsArray.push(revision);
    }

    revisionsArray.sort(function(a, b) {
        return (b - a);
    });

    return revisionsArray;
}

function htmlForIndividualTestOnAllBuilders(test)
{
    processTestRunsForAllBuilders();

    var testResults = g_testToResultsMap[test];
    if (!testResults)
        return '<div class="not-found">Test not found. Either it does not exist, is skipped or passes on all recorded runs.</div>';

    var masters = {};
    testResults.forEach(function(testResult) {
        var masterName = testResult.builder.masterName;
        var masterResults = masters[masterName];
        if (!masterResults) {
            masterResults = [];
            masters[masterName] = masterResults;
        }
        masterResults.push(testResult);
    });

    var html = '';
    var shownBuilderKeys = [];
    var revisions = collapsedRevisionList(testResults);
    var sortedMasters = Object.keys(masters).sort(function(a, b) {
      // FYI masters sort after other masters, sort alphanumerically otherwise.
      var aFyi = (a.indexOf('.fyi') != -1);
      var bFyi = (b.indexOf('.fyi') != -1);
      if (aFyi == bFyi) {
        return (a < b) ?  -1 : (a > b ? 1 : 0);
      }
      return aFyi ? 1 : -1;
    })
    sortedMasters.forEach(function(masterName) {
        var isFyi = (masterName.indexOf('.fyi') != -1);
        var rowClass = isFyi ? ' class="fyi"' : '';
        html += '<tr' + rowClass + '><td class="master-name" colspan=5>' +
            masterName + '</td></tr>';
        var masterResults = masters[masterName];
        for (var j = 0; j < masterResults.length; j++) {
            shownBuilderKeys.push(masterResults[j].builder.key());
            var showBuilderNames = true;
            html += htmlForSingleTestRow(
                masterResults[j], showBuilderNames, revisions, isFyi);
        }
    });

    var skippedBuilderKeys = []
    _currentBuilders().forEach(function(builder) {
        if (shownBuilderKeys.indexOf(builder.key()) == -1)
            skippedBuilderKeys.push(builder.key());
    });

    var skippedBuildersHtml = '';
    if (skippedBuilderKeys.length) {
        skippedBuildersHtml = '<div>The following builders either don\'t run this test (e.g. it\'s skipped) or all recorded runs passed:</div>' +
            '<div class=skipped-builder-list><div class=skipped-builder>' + skippedBuilderKeys.join('</div><div class=skipped-builder>') + '</div></div>';
    }

    return htmlForTestTable(html) + skippedBuildersHtml;
}

function htmlForIndividualTestOnAllBuildersWithResultsLinks(test)
{
    processTestRunsForAllBuilders();

    var testResults = g_testToResultsMap[test];
    var html = '';
    html += htmlForIndividualTestOnAllBuilders(test);

    if (!g_history.dashboardSpecificState.showChrome) {
        var newWindowURL = window.location.hash.replace(/showChrome=false/, 'showChrome=true');
        html += linkHTMLToOpenWindow(newWindowURL, 'Pop out in a new tab');
    } else {
        html += '<div class=expectations test=' + test + '><div>';
        html += linkHTMLToToggleState('showExpectations', 'results');

        if (g_history.isLayoutTestResults() || g_history.isGPUTestResults()) {
            if (g_history.isLayoutTestResults())
                html += ' | ' + linkHTMLToToggleState('showLargeExpectations', 'large thumbnails');
            html += ' | <b>Only shows actual results/diffs from the most recent *failure* on each bot.</b>';
        } else {
          html += ' | <span>Results height:<input ' +
              'onchange="g_history.setQueryParameter(\'resultsHeight\',this.value)" value="' +
              g_history.dashboardSpecificState.resultsHeight + '" style="width:2.5em">px</span>';
        }
        html += '</div></div>';
    }
    return html;
}

function maybeAddPngChecksum(expectationDiv, pngUrl)
{
    // pngUrl gets served from the browser cache since we just loaded it in an
    // <img> tag.
    loader.request(pngUrl,
        function(xhr) {
            // Convert the first 2k of the response to a byte string.
            var bytes = xhr.responseText.substring(0, 2048);
            for (var position = 0; position < bytes.length; ++position)
                bytes[position] = bytes[position] & 0xff;

            // Look for the comment.
            var commentKey = 'tEXtchecksum\x00';
            var checksumPosition = bytes.indexOf(commentKey);
            if (checksumPosition == -1)
                return;

            var checksum = bytes.substring(checksumPosition + commentKey.length, checksumPosition + commentKey.length + 32);
            var checksumContainer = document.createElement('span');
            checksumContainer.innerText = 'Embedded checksum: ' + checksum;
            checksumContainer.setAttribute('class', 'pngchecksum');
            expectationDiv.parentNode.appendChild(checksumContainer);
        },
        function(xhr) {},
        true);
}

function getOrCreate(className, parent)
{
    var element = parent.querySelector('.' + className);
    if (!element) {
        element = document.createElement('div');
        element.className = className;
        parent.appendChild(element);
    }
    return element;
}

function handleExpectationsItemLoad(title, item, itemType, parent)
{
    item.className = 'expectation';
    if (g_history.dashboardSpecificState.showLargeExpectations)
        item.className += ' large';

    var titleContainer = document.createElement('h3');
    titleContainer.className = 'expectations-title';
    titleContainer.textContent = title;

    var itemContainer = document.createElement('span');
    itemContainer.appendChild(titleContainer);
    itemContainer.className = 'expectations-item ' + title;
    itemContainer.appendChild(item);

    // Separate text and image results into separate divs..
    var typeContainer = getOrCreate(itemType, parent);

    // Insert results in a consistent order.
    var index = EXPECTATIONS_ORDER.indexOf(title);
    while (index < EXPECTATIONS_ORDER.length) {
        index++;
        var elementAfter = typeContainer.querySelector('.' + EXPECTATIONS_ORDER[index]);
        if (elementAfter) {
            typeContainer.insertBefore(itemContainer, elementAfter);
            break;
        }
    }
    if (!itemContainer.parentNode)
        typeContainer.appendChild(itemContainer);

    handleFinishedLoadingExpectations(parent);
}

function addExpectationItem(expectationsContainers, parentContainer, url, opt_builder)
{
    // Group expectations by builder, putting test and reference files first.
    var builderName = opt_builder && opt_builder.builderName || "Test and reference files";
    var container = expectationsContainers[builderName];

    if (!container) {
        container = document.createElement('div');
        container.className = 'expectations-container';
        container.setAttribute('data-builder-name', builderName);
        parentContainer.appendChild(container);
        expectationsContainers[builderName] = container;
    }

    var numUnloaded = container.getAttribute('data-unloaded') || 0;
    container.setAttribute('data-unloaded', ++numUnloaded);

    var isImage = url.match(/\.png$/);

    var appendExpectationsItem = function(item) {
        var itemType = isImage ? 'image' : 'text';
        handleExpectationsItemLoad(expectationsTitle(url), item, itemType, container);
    };

    var handleLoadError = function() {
        handleFinishedLoadingExpectations(container);
    };

    if (isImage) {
        var dummyNode = document.createElement('img');
        dummyNode.onload = function() {
            var item = dummyNode;
            maybeAddPngChecksum(item, url);
            appendExpectationsItem(item);
        }
        dummyNode.onerror = handleLoadError;
        dummyNode.src = url;
    } else {
        loader.request(url,
            function(xhr) {
                var item = document.createElement('pre');
                if (string.endsWith(url, '-wdiff.html'))
                    item.innerHTML = xhr.responseText;
                else
                    item.textContent = xhr.responseText;
                appendExpectationsItem(item);
            },
            handleLoadError);
    }
}

function handleFinishedLoadingExpectations(container)
{
    var numUnloaded = container.getAttribute('data-unloaded') - 1;
    container.setAttribute('data-unloaded', numUnloaded);
    if (numUnloaded)
        return;

    if (!container.firstChild) {
        container.remove();
        return;
    }

    var builderName = container.getAttribute('data-builder-name');
    if (!builderName)
        return;

    var header = document.createElement('h2');
    header.textContent = builderName;
    container.insertBefore(header, container.firstChild);
}

function expectationsTitle(url)
{
    var matchingSuffixes = ACTUAL_RESULT_SUFFIXES.filter(function(suffix) {
        return string.endsWith(url, suffix);
    });

    if (matchingSuffixes.length)
        return matchingSuffixes[0].split('.')[0];

    var parts = url.split('/');
    return parts[parts.length - 1];
}

function loadExpectations(expectationsContainer)
{
    var test = expectationsContainer.getAttribute('test');
    if (g_history.isLayoutTestResults())
        loadExpectationsLayoutTests(test, expectationsContainer);
    else {
        var testResults = g_testToResultsMap[test];
        for (var i = 0; i < testResults.length; i++)
            if (g_history.isGPUTestResults())
                loadGPUResultsForBuilder(testResults[i].builder, test, expectationsContainer);
            else
                loadNonWebKitResultsForBuilder(testResults[i].builder, test, expectationsContainer);
    }
}

function gpuResultsPath(chromeRevision, builderName)
{
  return chromeRevision + '_' + builderName.replace(/[^A-Za-z0-9]+/g, '_');
}

function loadGPUResultsForBuilder(builder, test, expectationsContainer)
{
    var container = document.createElement('div');
    container.className = 'expectations-container';
    container.innerHTML = '<div><b>' + builder.key() + '</b></div>';
    expectationsContainer.appendChild(container);

    var failureIndex = indexesForFailures(builder, test)[0];

    var buildNumber = g_resultsByBuilder[builder.key()].buildNumbers[failureIndex];
    var pathToLog = builder.master().logPath(builder.builderName, buildNumber) + pathToFailureLog(test);

    var chromeRevision = g_resultsByBuilder[builder.key()].chromeRevision[failureIndex];
    var resultsUrl = GPU_RESULTS_BASE_PATH + gpuResultsPath(chromeRevision, builder.builderName);
    var filename = test.split(/\./)[1] + '.png';

    appendNonWebKitResults(container, pathToLog, 'non-webkit-results');
    appendNonWebKitResults(container, resultsUrl + '/gen/' + filename, 'gpu-test-results', 'Generated');
    appendNonWebKitResults(container, resultsUrl + '/ref/' + filename, 'gpu-test-results', 'Reference');
    appendNonWebKitResults(container, resultsUrl + '/diff/' + filename, 'gpu-test-results', 'Diff');
}

function loadNonWebKitResultsForBuilder(builder, test, expectationsContainer)
{
    var failureIndexes = indexesForFailures(builder, test);
    var container = document.createElement('div');
    container.innerHTML = '<div><b>' + builder.key() + '</b></div>';
    expectationsContainer.appendChild(container);
    for (var i = 0; i < failureIndexes.length; i++) {
        // FIXME: This doesn't seem to work anymore. Did the paths change?
        // Once that's resolved, see if we need to try each gtest modifier prefix as well.
        var buildNumber = g_resultsByBuilder[builder.key()].buildNumbers[failureIndexes[i]];
        var pathToLog = builder.master().logPath(builder.builderName, buildNumber) + pathToFailureLog(test);
        appendNonWebKitResults(container, pathToLog, 'non-webkit-results');
    }
}

function appendNonWebKitResults(container, url, itemClassName, opt_title)
{
    // Use a script tag to detect whether the URL 404s.
    // Need to use a script tag since the URL is cross-domain.
    var dummyNode = document.createElement('script');
    dummyNode.src = url;

    dummyNode.onload = function() {
        var item = document.createElement('iframe');
        item.src = dummyNode.src;
        item.className = itemClassName;
        item.style.height = g_history.dashboardSpecificState.resultsHeight + 'px';

        if (opt_title) {
            var childContainer = document.createElement('div');
            childContainer.style.display = 'inline-block';
            var title = document.createElement('div');
            title.textContent = opt_title;
            childContainer.appendChild(title);
            childContainer.appendChild(item);
            container.replaceChild(childContainer, dummyNode);
        } else
            container.replaceChild(item, dummyNode);
    }
    dummyNode.onerror = function() {
        container.removeChild(dummyNode);
    }

    container.appendChild(dummyNode);
}

function lookupVirtualTestSuite(test) {
    for (var suite in VIRTUAL_SUITES) {
        if (test.indexOf(suite) != -1)
            return suite;
    }
    return '';
}

function baseTest(test, suite) {
    base = VIRTUAL_SUITES[suite];
    return base ? test.replace(suite, base) : test;
}

function loadTestAndReferenceFiles(expectationsContainers, expectationsContainer, test) {
    var testWithoutSuffix = test.substring(0, test.lastIndexOf('.'));
    var reftest_html_file = testWithoutSuffix + "-expected.html";
    var reftest_mismatch_html_file = testWithoutSuffix + "-expected-mismatch.html";

    var suite = lookupVirtualTestSuite(test);
    if (suite) {
        loadTestAndReferenceFiles(expectationsContainers, expectationsContainer, baseTest(test, suite));
        return;
    }

    addExpectationItem(expectationsContainers, expectationsContainer, TEST_URL_BASE_PATH_FOR_XHR + test);
    addExpectationItem(expectationsContainers, expectationsContainer, TEST_URL_BASE_PATH_FOR_XHR + reftest_html_file);
    addExpectationItem(expectationsContainers, expectationsContainer, TEST_URL_BASE_PATH_FOR_XHR + reftest_mismatch_html_file);
}

function loadExpectationsLayoutTests(test, expectationsContainer)
{
    // Map from file extension to container div for expectations of that type.
    var expectationsContainers = {};
    loadTestAndReferenceFiles(expectationsContainers, expectationsContainer, test);

    var testWithoutSuffix = test.substring(0, test.lastIndexOf('.'));

    _currentBuilders().forEach(function(builder) {
        var actualResultsBase = TEST_RESULTS_BASE_PATH + builder.builderNameForPath() + '/results/layout-test-results/';
        ACTUAL_RESULT_SUFFIXES.forEach(function(suffix) {{
            addExpectationItem(expectationsContainers, expectationsContainer, actualResultsBase + testWithoutSuffix + '-' + suffix, builder);
        }})
    });

    // Add a clearing element so floated elements don't bleed out of their
    // containing block.
    var br = document.createElement('br');
    br.style.clear = 'both';
    expectationsContainer.appendChild(br);
}

function appendExpectations()
{
    var expectations = g_history.dashboardSpecificState.showExpectations ? document.getElementsByClassName('expectations') : [];
    g_chunkedActionState = {
        items: expectations,
        index: 0
    }
    performChunkedAction(function(expectation) {
            loadExpectations(expectation);
            postHeightChangedMessage();
        },
        hideLoadingUI,
        expectations);
}

function hideLoadingUI()
{
    var loadingDiv = $('loading-ui');
    if (loadingDiv)
        loadingDiv.style.display = 'none';
    postHeightChangedMessage();
}

function generatePageForIndividualTests(tests)
{
    if (g_history.dashboardSpecificState.showChrome)
        appendHTML(htmlForNavBar());
    performChunkedAction(function(test) {
            appendHTML(htmlForIndividualTest(test));
        },
        appendExpectations,
        tests);
    if (g_history.dashboardSpecificState.showChrome) {
        $('tests-input').value = g_history.dashboardSpecificState.tests;
        $('result-input').value = g_history.dashboardSpecificState.result;
    }
}

var g_chunkedActionRequestId;
function performChunkedAction(action, onComplete, items, opt_index) {
    if (g_chunkedActionRequestId)
        cancelAnimationFrame(g_chunkedActionRequestId);

    var index = opt_index || 0;
    g_chunkedActionRequestId = requestAnimationFrame(function() {
        if (index < items.length) {
            action(items[index]);
            performChunkedAction(action, onComplete, items, ++index);
        } else {
            onComplete();
        }
    });
}

function htmlForIndividualTest(test)
{
    var testNameHtml = '';
    if (g_history.dashboardSpecificState.showChrome) {
        if (g_history.isLayoutTestResults()) {
            var suite = lookupVirtualTestSuite(test);
            var base = suite ? baseTest(test, suite) : test;
            var versionControlUrl = TEST_URL_BASE_PATH_FOR_BROWSING + base;
            testNameHtml += '<h2>' + linkHTMLToOpenWindow(versionControlUrl, test) + '</h2>';
        } else
            testNameHtml += '<h2>' + test + '</h2>';
    }

    return testNameHtml + htmlForIndividualTestOnAllBuildersWithResultsLinks(test);
}

function setTestsParameter(input)
{
    g_history.setQueryParameter('tests', input.value);
}

function htmlForNavBar()
{
    var extraHTML = '';
    var html = ui.html.testTypeSwitcher(false, extraHTML, isCrossBuilderView());
    html += '<div class=forms><form id=result-form ' +
        'onsubmit="g_history.setQueryParameter(\'result\', result.value);' +
        'return false;">Show all tests with result: ' +
        '<input name=result placeholder="e.g. CRASH" id=result-input>' +
        '</form><span>Show tests on all platforms: </span>' +
        // Use a textarea to avoid the 32k limit on the length of inputs.
        '<textarea name=tests ' +
        'placeholder="Comma or space-separated list of tests or partial ' +
        'paths to show test results across all builders, e.g., ' +
        'foo/bar.html,foo/baz,domstorage" id=tests-input onchange="setTestsParameter(this)" ' +
        'onkeydown="if (event.keyCode == 13) { setTestsParameter(this); return false; }"></textarea>' +
        '<span class=link onclick="showLegend()">Show legend [type ?]</span></div>';
    return html;
}

function checkBoxToToggleState(key, text)
{
    var stateEnabled = g_history.dashboardSpecificState[key];
    return '<label><input type=checkbox ' + (stateEnabled ? 'checked ' : '') + 'onclick="g_history.setQueryParameter(\'' + key + '\', ' + !stateEnabled + ')">' + text + '</label> ';
}

function linkHTMLToToggleState(key, linkText)
{
    var stateEnabled = g_history.dashboardSpecificState[key];
    return '<span class=link onclick="g_history.setQueryParameter(\'' + key + '\', ' + !stateEnabled + ')">' + (stateEnabled ? 'Hide' : 'Show') + ' ' + linkText + '</span>';
}

function headerForTestTableHtml()
{
    return '<h2 style="display:inline-block">Failing tests</h2>' +
        checkBoxToToggleState('showFlaky', 'Show flaky') +
        checkBoxToToggleState('showNonFlaky', 'Show non-flaky') +
        checkBoxToToggleState('showSkip', 'Show Skip') +
        checkBoxToToggleState('showWontFix', 'Show WontFix');
}

function generatePageForBuilder(builder)
{
    processTestRunsForBuilder(builder);

    var filteredResults = g_perBuilderFailures[builder.key()].filter(shouldShowTest);
    sortTests(filteredResults, g_history.dashboardSpecificState.sortColumn, g_history.dashboardSpecificState.sortOrder);

    var testsHTML = '';
    if (filteredResults.length) {
        var tableRowsHTML = '';
        var showBuilderNames = false;
        var revisions = collapsedRevisionList(filteredResults);
        for (var i = 0; i < filteredResults.length; i++)
            tableRowsHTML += htmlForSingleTestRow(filteredResults[i], showBuilderNames, revisions, false);
        testsHTML = htmlForTestTable(tableRowsHTML);
    } else {
        if (g_history.isLayoutTestResults())
            testsHTML += '<div>Fill in one of the text inputs or checkboxes above to show failures.</div>';
        else
            testsHTML += '<div>No tests have failed!</div>';
    }

    var html = htmlForNavBar();

    if (g_history.isLayoutTestResults())
        html += headerForTestTableHtml();

    html += '<br>' + testsHTML;
    appendHTML(html);

    var ths = document.getElementsByTagName('th');
    for (var i = 0; i < ths.length; i++) {
        ths[i].addEventListener('click', changeSort, false);
        ths[i].className = "sortable";
    }

    hideLoadingUI();
}

var VALID_KEYS_FOR_CROSS_BUILDER_VIEW = {
    tests: 1,
    result: 1,
    showChrome: 1,
    showExpectations: 1,
    showLargeExpectations: 1,
    resultsHeight: 1,
    revision: 1
};

function isInvalidKeyForCrossBuilderView(key)
{
    return !(key in VALID_KEYS_FOR_CROSS_BUILDER_VIEW) && !(key in history.DEFAULT_CROSS_DASHBOARD_STATE_VALUES);
}

function hideLegend()
{
    var legend = $('legend');
    if (legend)
        legend.parentNode.removeChild(legend);
}

function showLegend()
{
    var legend = $('legend');
    if (!legend) {
        legend = document.createElement('div');
        legend.id = 'legend';
        document.body.appendChild(legend);
    }

    var html = '<div id=legend-toggle onclick="hideLegend()">Hide ' +
        'legend [type esc]</div><div id=legend-contents>';

    // Just grab the first failureMap. Technically, different builders can have different maps if they
    // haven't all cycled after the map was changed, but meh.
    var failureMap = g_resultsByBuilder[Object.keys(g_resultsByBuilder)[0]][results.FAILURE_MAP];
    var seen = {};
    for (var expectation in failureMap) {
        var failureString = failureMap[expectation];
        if (failureString in seen) {
          continue;
        }
        seen[failureString] = true;
        html += '<div class=' + classNameForFailureString(failureString) + '>' + failureString + '</div>';
    }

    if (g_history.isLayoutTestResults()) {
      html += '</div><br style="clear:both">' +
          '</div>';

      html += '<div>RELEASE TIMEOUTS:</div>' +
          htmlForSlowTimes(RELEASE_TIMEOUT) +
          '<div>DEBUG TIMEOUTS:</div>' +
          htmlForSlowTimes(DEBUG_TIMEOUT);
    }

    legend.innerHTML = html;
}

function htmlForSlowTimes(minTime)
{
    return '<ul><li>' + minTime + ' seconds</li><li>' +
        SLOW_MULTIPLIER * minTime + ' seconds if marked Slow in TestExpectations</li></ul>';
}

function postHeightChangedMessage()
{
    if (window == parent)
        return;

    var root = document.documentElement;
    var height = root.offsetHeight;
    if (root.offsetWidth < root.scrollWidth) {
        // We have a horizontal scrollbar. Include it in the height.
        var dummyNode = document.createElement('div');
        dummyNode.style.overflow = 'scroll';
        document.body.appendChild(dummyNode);
        var scrollbarWidth = dummyNode.offsetHeight - dummyNode.clientHeight;
        document.body.removeChild(dummyNode);
        height += scrollbarWidth;
    }
    parent.postMessage({command: 'heightChanged', height: height}, '*')
}

function _currentBuilders()
{
    var bldrs = builders.getBuilders(g_history.crossDashboardState.testType);
    // TODO(sergiyb): Remove this when all layout tests are using swarming (crbug.com/524758).
    if (g_history.crossDashboardState.testType == "webkit_layout_tests")
      bldrs = bldrs.concat(builders.getBuilders("webkit_tests"));
    if (g_history.crossDashboardState.testType == "webkit_tests")
      bldrs = bldrs.concat(builders.getBuilders("webkit_layout_tests"));
    return bldrs
}

if (window != parent)
    window.addEventListener('blur', ui.popup.hide);

document.addEventListener('keydown', function(e) {
    if (e.keyIdentifier == 'U+003F' || e.keyIdentifier == 'U+00BF') {
        // WebKit MAC retursn 3F. WebKit WIN returns BF. This is a bug!
        // ? key
        showLegend();
    } else if (e.keyIdentifier == 'U+001B') {
        // escape key
        hideLegend();
        ui.popup.hide();
    }
}, false);

window.addEventListener('load', function() {
    resourceLoader = new loader.Loader();
    resourceLoader.load();
}, false);
