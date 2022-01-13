// Copyright 2022 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

import { LitElement, html, customElement, property, css, state, TemplateResult } from 'lit-element';
import { styleMap } from 'lit-html/directives/style-map';
import { DateTime } from 'luxon';
import '@material/mwc-button';
import '@material/mwc-icon';
import "@material/mwc-list/mwc-list-item";
import '@material/mwc-select';

// Indent of each level of grouping in the table in pixels.
const levelIndent = 10;

// FailureTable lists the failures in a cluster tracked by Weetbix.
@customElement('failure-table')
export class FailureTable extends LitElement {
    @property()
    project: string = '';

    @property()
    clusterAlgorithm: string = '';

    @property()
    clusterID: string = '';

    @state()
    failures: ClusterFailure[] | undefined;

    @state()
    groups: FailureGroup[] = [];

    @state()
    variants: FailureVariant[] = [];

    @state()
    impactFilter: ImpactFilter = impactFilters[0];

    @property()
    sortMetric: MetricName = 'latestFailureTime';

    @property({ type: Boolean })
    ascending: boolean = false;

    connectedCallback() {
        super.connectedCallback()

        fetch(`/api/projects/${encodeURIComponent(this.project)}/clusters/${encodeURIComponent(this.clusterAlgorithm)}/${encodeURIComponent(this.clusterID)}/failures`)
            .then(r => r.json())
            .then((failures: ClusterFailure[]) => {
                this.failures = failures
                this.countDistictVariantValues();
                this.groupAndSortFailures();
            });
    }

    countDistictVariantValues() {
        if (!this.failures) {
            return;
        }
        this.variants = [];
        this.failures.forEach(f => {
            f.variant.forEach(v => {
                if (!v.key) {
                    return;
                }
                const variant = this.variants.filter(e => e.key === v.key)?.[0];
                if (!variant) {
                    this.variants.push({ key: v.key, values: [v.value || ''], isSelected: false });
                } else {
                    if (variant.values.indexOf(v.value || '') === -1) {
                        variant.values.push(v.value || '');
                    }
                }
            });
        });
    }

    groupAndSortFailures() {
        if (this.failures) {
            this.groups = groupFailures(this.failures, f => {
                const variantValues = this.variants.filter(v => v.isSelected)
                    .map(v => f.variant.filter(fv => fv.key === v.key)?.[0]?.value || '');
                return [...variantValues, f.testId || ''];
            });
            this.groups.forEach(group => {
                treeDistinctValues(group, failureIdExtractor(), (g, values) => g.failures = values.size);
                treeDistinctValues(group, rejectedTestRunIdExtractor(this.impactFilter), (g, values) => g.testRunFailures = values.size);
                treeDistinctValues(group, rejectedIngestedInvocationIdExtractor(this.impactFilter), (g, values) => g.invocationFailures = values.size);
                treeDistinctValues(group, rejectedPresubmitRunIdExtractor(this.impactFilter), (g, values) => g.presubmitRejects = values.size);
            });
        }
        this.sortFailures();
    }

    sortFailures() {
        sortFailureGroups(this.groups, this.sortMetric, this.ascending);
        this.requestUpdate();
    }

    toggleSort(metric: MetricName) {
        if (metric === this.sortMetric) {
            this.ascending = !this.ascending;
        } else {
            this.sortMetric = metric;
            this.ascending = false;
        }
        this.sortFailures();
    }

    onImpactFilterChanged() {
        const item = this.shadowRoot!.querySelector('#impact-filter [selected]');
        if (item) {
            const selected = item.getAttribute('value')
            this.impactFilter = impactFilters.filter(f => f.name == selected)?.[0] || impactFilters[0];
        }
        this.groupAndSortFailures();
    }

    toggleVariant(variant: FailureVariant) {
        const index = this.variants.indexOf(variant);
        this.variants.splice(index, 1);
        variant.isSelected = !variant.isSelected;
        const numSelected = this.variants.filter(v => v.isSelected).length;
        this.variants.splice(numSelected, 0, variant);
        this.groupAndSortFailures();
    }

    toggleExpand(group: FailureGroup) {
        group.isExpanded = !group.isExpanded;
        this.requestUpdate();
    }

    render() {
        const unselectedVariants = this.variants.filter(v => !v.isSelected).map(v => v.key);
        if (this.failures === undefined) {
            return html`Loading cluster failures...`;
        }
        const ungroupedVariants = (failure: ClusterFailure) => {
            return unselectedVariants.map(key => failure.variant.filter(v => v.key == key)?.[0]).filter(v => v);
        }
        const failureLink = (failure: ClusterFailure) => {
            const query = `ID:${failure.testId} `
            return `https://ci.chromium.org/ui/inv/${failure.ingestedInvocationId}/test-results?q=${encodeURIComponent(query)}`;
        }
        const indentStyle = (level: number) => {
            return styleMap({ paddingLeft: (levelIndent * level) + 'px' });
        }
        const groupRow = (group: FailureGroup): TemplateResult => {
            return html`
            <tr>
                ${group.failure ?
                    html`<td style=${indentStyle(group.level)}>
                        <a href=${failureLink(group.failure)} target="_blank">${group.failure.ingestedInvocationId}</a>
                        <span class="variant-info">${ungroupedVariants(group.failure).map(v => `${v.key}: ${v.value}`).join(', ')}</span>
                    </td>` :
                    html`<td class="group" style=${indentStyle(group.level)} @click=${() => this.toggleExpand(group)}>
                        <mwc-icon>${group.isExpanded ? 'keyboard_arrow_down' : 'keyboard_arrow_right'}</mwc-icon>
                        ${group.name || 'none'}
                    </td>`}
                <td class="number">${group.presubmitRejects}</td>
                <td class="number">${group.invocationFailures}</td>
                <td class="number">${group.testRunFailures}</td>
                <td class="number">${group.failures}</td>
                <td>${group.latestFailureTime.toRelative()}</td>
            </tr>
            ${group.isExpanded ? group.children.map(child => groupRow(child)) : null}`
        }
        const filterButton = (variant: FailureVariant) => {
            return html`
                <mwc-button
                    label=${`${variant.key} (${variant.values.length})`}
                    ?unelevated=${variant.isSelected}
                    ?outlined=${!variant.isSelected}
                    @click=${() => this.toggleVariant(variant)}></mwc-button>`;
        }
        return html`
            <div class="controls">
                <div class="select-offset">
                    <mwc-select id="impact-filter" outlined label="Impact" @change=${() => this.onImpactFilterChanged()}>
                        ${impactFilters.map((filter) => html`<mwc-list-item ?selected=${filter == this.impactFilter} value="${filter.name}">${filter.name}</mwc-list-item>`)}
                    </mwc-select>
                </div>
                <div>
                    <div class="label">
                        Group By
                    </div>
                    ${this.variants.map(v => filterButton(v))}
                </div>
            </div>
            <table>
                <thead>
                    <tr>
                        <th></th>
                        <th class="sortable" @click=${() => this.toggleSort('presubmitRejects')}>
                            Presubmit Runs Failed
                            ${this.sortMetric === 'presubmitRejects' ? html`<mwc-icon>${this.ascending ? 'expand_less' : 'expand_more'}</mwc-icon>` : null}
                        </th>
                        <th class="sortable" @click=${() => this.toggleSort('invocationFailures')}>
                            Invocations Failed
                            ${this.sortMetric === 'invocationFailures' ? html`<mwc-icon>${this.ascending ? 'expand_less' : 'expand_more'}</mwc-icon>` : null}
                        </th>
                        <th class="sortable" @click=${() => this.toggleSort('testRunFailures')}>
                            Test Runs Failed
                            ${this.sortMetric === 'testRunFailures' ? html`<mwc-icon>${this.ascending ? 'expand_less' : 'expand_more'}</mwc-icon>` : null}
                        </th>
                        <th class="sortable" @click=${() => this.toggleSort('failures')}>
                            Unexpected Failures
                            ${this.sortMetric === 'failures' ? html`<mwc-icon>${this.ascending ? 'expand_less' : 'expand_more'}</mwc-icon>` : null}
                        </th>
                        <th class="sortable" @click=${() => this.toggleSort('latestFailureTime')}>
                            Latest Failure Time
                            ${this.sortMetric === 'latestFailureTime' ? html`<mwc-icon>${this.ascending ? 'expand_less' : 'expand_more'}</mwc-icon>` : null}
                        </th>
                    </tr>
                </thead>
                <tbody>
                    ${this.groups.map(group => groupRow(group))}
                </tbody>
            </table>
        `;
    }
    static styles = [css`
        .controls {
            display: flex;            
            gap: 30px;
        }
        .label {
            color: var(--greyed-out-text-color);
            font-size: var(--font-size-small);
        }
        .select-offset {
            padding-top: 7px
        }
        #impact-filter {
            width: 230px;
        }
        table {
            border-collapse: collapse;
            width: 100%;
            table-layout: fixed;
        }
        th {
            font-weight: normal;
            color: var(--greyed-out-text-color);
            font-size: var(--font-size-small);
            text-align: left;
        }
        td,th {
            padding: 4px;
            max-width: 80%;
        }
        td.number {
            text-align: right;
        }
        td.group {
            word-break: break-all;
        }
        th.sortable {
            cursor: pointer;
            width:120px;
        }
        tbody tr:hover {
            background-color: var(--light-active-color);
        }
        .group {
            cursor: pointer;
            --mdc-icon-size: var(--font-size-default);
        }
        .variant-info {
            color: var(--greyed-out-text-color);
            font-size: var(--font-size-small);
        }
    `];
}

// ImpactFilter represents what kind of impact should be counted or ignored in
// calculating impact for failures.
export interface ImpactFilter {
    name: string;
    ignoreExoneration: boolean;
    ignoreIngestedInvocationBlocked: boolean;
    ignoreTestRunBlocked: boolean;
}
export const impactFilters: ImpactFilter[] = [
    {
        name: 'Actual Impact',
        ignoreExoneration: false,
        ignoreIngestedInvocationBlocked: false,
        ignoreTestRunBlocked: false,
    }, {
        name: 'Without Exoneration',
        ignoreExoneration: true,
        ignoreIngestedInvocationBlocked: false,
        ignoreTestRunBlocked: false,
    }, {
        name: 'Without Invocation Retries',
        ignoreExoneration: true,
        ignoreIngestedInvocationBlocked: true,
        ignoreTestRunBlocked: false,
    }, {
        name: 'Without Any Retries',
        ignoreExoneration: true,
        ignoreIngestedInvocationBlocked: true,
        ignoreTestRunBlocked: true,
    }
];

// group a number of failures into a tree of failure groups.
// grouper is a function that returns a list of keys, one corresponding to each level of the grouping tree.
// impactFilter controls how metric counts are aggregated from failures into parent groups (see treeCounts and rejected... functions).
const groupFailures = (failures: ClusterFailure[], grouper: (f: ClusterFailure) => string[]): FailureGroup[] => {
    const topGroups: FailureGroup[] = [];
    failures.forEach(f => {
        const keys = grouper(f);
        let groups = topGroups;
        let failureTime = DateTime.fromISO(f.partitionTime || '');
        let level = 0;
        for (const key of keys) {
            const group = getOrCreateGroup(groups, key, failureTime);
            group.level = level;
            level += 1;
            groups = group.children;
        }
        const failureGroup = newGroup('', failureTime);
        failureGroup.failure = f;
        failureGroup.level = level;
        groups.push(failureGroup);
    });
    return topGroups;
}

// Create a new group.
const newGroup = (name: string, failureTime: DateTime): FailureGroup => {
    return {
        name: name,
        failures: 0,
        invocationFailures: 0,
        testRunFailures: 0,
        presubmitRejects: 0,
        children: [],
        isExpanded: false,
        latestFailureTime: failureTime,
        level: 0
    };
}

// Find a group by name in the given list of groups, create a new one and insert it if it is not found.
// failureTime is only used when creating a new group.
const getOrCreateGroup = (groups: FailureGroup[], name: string, failureTime: DateTime): FailureGroup => {
    let group = groups.filter(g => g.name == name)?.[0];
    if (group) {
        return group;
    }
    group = newGroup(name, failureTime);
    groups.push(group);
    return group;
}

// Returns the distinct values returned by featureExtractor for all children of the group.
// If featureExtractor returns undefined, the failure will be ignored.
// The distinct values for each group in the tree are also reported to `visitor` as the tree is traversed.
// A typical `visitor` function will store the count of distinct values in a property of the group.
const treeDistinctValues = (group: FailureGroup,
    featureExtractor: FeatureExtractor,
    visitor: (group: FailureGroup, distinctValues: Set<string>) => void): Set<string> => {
    const values: Set<string> = new Set();
    if (group.failure) {
        const f = featureExtractor(group.failure);
        if (f !== undefined) {
            values.add(f)
        }
    } else {
        for (const child of group.children) {
            for (let value of treeDistinctValues(child, featureExtractor, visitor)) {
                values.add(value);
            }
        }
    }
    visitor(group, values);
    return values;
}

// A FeatureExtractor returns a string representing some feature of a ClusterFailure.
// Returns undefined if there is no such feature for this failure.
type FeatureExtractor = (failure: ClusterFailure) => string | undefined;

// failureIdExtractor returns an extractor that returns a unique failure id for each failure.
// As failures don't actually have ids, it just returns an incrementing integer.
const failureIdExtractor = (): FeatureExtractor => {
    let unique = 0;
    return _f => {
        unique += 1;
        return '' + unique;
    }
}

// Returns an extractor that returns the id of the test run that was rejected by this failure, if any.
// The impact filter is taken into account in determining if the run was rejected by this failure.
const rejectedTestRunIdExtractor = (impactFilter: ImpactFilter): FeatureExtractor => {
    return f => {
        if (f.isExonerated && !impactFilter.ignoreExoneration) {
            return undefined;
        }
        if (!impactFilter.ignoreTestRunBlocked && !f.isTestRunBlocked) {
            return undefined;
        }
        return f.testRunId || undefined;
    }
}

// Returns an extractor that returns the id of the ingested invocation that was rejected by this failure, if any.
// The impact filter is taken into account in determining if the invocation was rejected by this failure.
const rejectedIngestedInvocationIdExtractor = (impactFilter: ImpactFilter): FeatureExtractor => {
    return f => {
        if (f.isExonerated && !impactFilter.ignoreExoneration) {
            return undefined;
        }
        if (!f.isIngestedInvocationBlocked && !impactFilter.ignoreIngestedInvocationBlocked) {
            return undefined;
        }
        if (!impactFilter.ignoreTestRunBlocked && !f.isTestRunBlocked) {
            return undefined;
        }
        return f.ingestedInvocationId || undefined;
    }
}

// Returns an extractor that returns the id of the presubmit run that was rejected by this failure, if any.
// The impact filter is taken into account in determining if the presubmit run was rejected by this failure.
const rejectedPresubmitRunIdExtractor = (impactFilter: ImpactFilter): FeatureExtractor => {
    return f => {
        if (f.isExonerated && !impactFilter.ignoreExoneration) {
            return undefined;
        }
        if (!f.isIngestedInvocationBlocked && !impactFilter.ignoreIngestedInvocationBlocked) {
            return undefined;
        }
        if (!impactFilter.ignoreTestRunBlocked && !f.isTestRunBlocked) {
            return undefined;
        }
        return f.presubmitRunId?.id || undefined;
    }
}

// Sorts child failure groups at each node of the tree by the given metric.
const sortFailureGroups = (groups: FailureGroup[], metric: MetricName, ascending: boolean) => {
    const getMetric = (group: FailureGroup): number => {
        switch (metric) {
            case 'failures':
                return group.failures;
            case 'presubmitRejects':
                return group.presubmitRejects;
            case 'invocationFailures':
                return group.invocationFailures;
            case 'testRunFailures':
                return group.testRunFailures;
            case 'latestFailureTime':
                return group.latestFailureTime.toSeconds();
            default:
                throw new Error('unknown metric: ' + metric);
        }
    }
    groups.sort((a, b) => ascending ? (getMetric(a) - getMetric(b)) : (getMetric(b) - getMetric(a)));;
    for (const group of groups) {
        if (group.children) {
            sortFailureGroups(group.children, metric, ascending);
        }
    }
}

// Flattens a group tree into a list of visible rows based on which groups are expanded or not.
const flattenGroupRows = (groups: FailureGroup[], flattened: FailureGroup[]) => {
    for (const group of groups) {
        flattened.push(group);
        if (group.isExpanded) {
            flattenGroupRows(group.children, flattened);
        }
    }
}

// The failure grouping code is complex, so export the parts for unit testing.
export const exportedForTesting = {
    groupFailures,
    impactFilters,
    rejectedIngestedInvocationIdExtractor,
    rejectedPresubmitRunIdExtractor,
    rejectedTestRunIdExtractor,
    sortFailureGroups,
    treeDistinctCounts: treeDistinctValues,
}

// ClusterFailure is the data returned by the server for each failure.
export interface ClusterFailure {
    realm: string | null;
    testId: string | null;
    variant: Variant[];
    presubmitRunId: PresubmitRunId | null;
    partitionTime: string | null;
    isExonerated: boolean | null;
    ingestedInvocationId: string | null;
    isIngestedInvocationBlocked: boolean | null;
    testRunId: string | null;
    isTestRunBlocked: boolean | null;
}

// Key/Value Variant pairs for failures.
interface Variant {
    key: string | null;
    value: string | null;
}

// Presubmit Run Ids of failures returned from the server.
interface PresubmitRunId {
    system: string | null;
    id: string | null;
}

// Metrics that can be used for sorting FailureGroups.
// Each value is a property of FailureGroup.
type MetricName = 'presubmitRejects' | 'invocationFailures' | 'testRunFailures' | 'failures' | 'latestFailureTime';

// FailureGroups are nodes in the failure tree hierarchy.
export interface FailureGroup {
    name: string;
    presubmitRejects: number;
    invocationFailures: number;
    testRunFailures: number;
    failures: number;
    latestFailureTime: DateTime;
    level: number;
    children: FailureGroup[];
    isExpanded: boolean;
    failure?: ClusterFailure;
}

// FailureVariant represents variant keys that appear on at least one failure.
interface FailureVariant {
    key: string;
    values: string[];
    isSelected: boolean;
}