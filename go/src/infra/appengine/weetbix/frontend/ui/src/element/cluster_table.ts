
// Copyright 2021 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

import { LitElement, html, customElement, property, css, state } from 'lit-element';
import '@material/mwc-checkbox';
import '@material/mwc-formfield';
import "@material/mwc-icon/mwc-icon";
import "@material/mwc-list/mwc-list-item";
import '@material/mwc-select';

// ClusterTable lists the clusters tracked by Weetbix.
@customElement('cluster-table')
export class ClusterTable extends LitElement {
    @property()
    project: string;

    @property({ type: Number })
    days: number = 1;

    @property({ type: Boolean })
    preexoneration: boolean;

    @property({ type: Boolean })
    residual: boolean = true;

    @property()
    sortMetric: MetricName = 'presubmitRejects';

    @property({ type: Boolean })
    ascending: boolean = false;

    @state()
    clusters: Cluster[] | undefined;

    connectedCallback() {
        super.connectedCallback()
        this.project = "chromium";
        fetch(`/api/projects/${encodeURIComponent(this.project)}/clusters`).then(r => r.json()).then(clusters => this.clusters = clusters);
    }

    onDaysChanged() {
        const item = this.shadowRoot.querySelector('#days [selected]');
        if (item) {
            this.days = parseInt(item.getAttribute('value'));
        }
    }

    sort(metric: MetricName) {
        if (metric === this.sortMetric) {
            this.ascending = !this.ascending;
        } else {
            this.sortMetric = metric;
            this.ascending = false;
        }
    }

    render() {
        if (this.clusters === undefined) {
            return html`Loading...`;
        }
        const clusterLink = (cluster: Cluster): string => {
            return `/projects/${encodeURIComponent(this.project)}/clusters/${encodeURIComponent(cluster.clusterId.algorithm)}/${encodeURIComponent(cluster.clusterId.id)}`;
        }
        const clusterDescription = (cluster: Cluster): string => {
            if (cluster.clusterId.algorithm.startsWith("testname-")) {
                return cluster.exampleTestId;
            } else if (cluster.clusterId.algorithm.startsWith("failurereason-")) {
                return cluster.exampleFailureReason;
            }
            return cluster.exampleFailureReason || cluster.exampleTestId || `${cluster.clusterId.algorithm}/${cluster.clusterId.id}`;
        }
        const metric = (c: Cluster, metric: MetricName): number => {
            let counts: Counts;
            switch (metric) {
                case 'presubmitRejects':
                    counts = this.days === 1 ? c.presubmitRejects1d : (this.days === 3 ? c.presubmitRejects3d : c.presubmitRejects7d);
                    break;
                case 'failures':
                    counts = this.days === 1 ? c.failures1d : (this.days === 3 ? c.failures3d : c.failures7d);
                    break;
                case 'testRunFailures':
                    counts = this.days === 1 ? c.testRunFailures1d : (this.days === 3 ? c.testRunFailures3d : c.testRunFailures7d);
                    break;
                default:
                    throw new Error('no such metric: ' + metric);
            }
            if (this.residual) {
                return this.preexoneration ? counts.residualPreExoneration : counts.residual;
            } else {
                return this.preexoneration ? counts.preExoneration : counts.nominal;
            }
        }
        const sortedClusters = [...this.clusters];
        sortedClusters.sort((c1, c2) => {
            const m1 = metric(c1, this.sortMetric);
            const m2 = metric(c2, this.sortMetric);
            return this.ascending ? m1 - m2 : m2 - m1;
        });
        return html`
        <div id="container">
            <h1>Clusters in project ${this.project}</h1>
            <mwc-select id="days" outlined label="Time Scale" @change=${() => this.onDaysChanged()}>
                <mwc-list-item selected value="1">1 Day</mwc-list-item>
                <mwc-list-item value="3">3 Days</mwc-list-item>
                <mwc-list-item value="7">7 Days</mwc-list-item>
            </mwc-select>
            <mwc-formfield label="Pre-Exoneration">
                <mwc-checkbox class="child" @change=${() => this.preexoneration = !this.preexoneration} ?checked=${this.preexoneration}></mwc-checkbox>
            </mwc-formfield>
            <mwc-formfield label="Residual">
                <mwc-checkbox checked class="child" @change=${() => this.residual = !this.residual} ?checked=${this.residual}></mwc-checkbox>
            </mwc-formfield>
            <table>
                <thead>
                    <tr>
                        <th>Cluster</th>
                        <th class="sortable" @click=${() => this.sort('presubmitRejects')}>
                            Presubmit Runs Failed
                            ${this.sortMetric === 'presubmitRejects' ? html`<mwc-icon>${this.ascending ? 'expand_less' : 'expand_more'}</mwc-icon>` : null}
                        </th>
                        <th class="sortable" @click=${() => this.sort('testRunFailures')}>
                            Test Runs Failed
                            ${this.sortMetric === 'testRunFailures' ? html`<mwc-icon>${this.ascending ? 'expand_less' : 'expand_more'}</mwc-icon>` : null}
                        </th>
                        <th class="sortable" @click=${() => this.sort('failures')}>
                            Unexpected Failures
                            ${this.sortMetric === 'failures' ? html`<mwc-icon>${this.ascending ? 'expand_less' : 'expand_more'}</mwc-icon>` : null}
                        </th>
                    </tr>
                </thead>
                <tbody>
                    ${sortedClusters.map(c => html`
                    <tr>
                        <td class="failure-reason">
                            <a ?data-suggested=${!c.clusterId.algorithm.startsWith('rules')} href=${clusterLink(c)}>
                                ${clusterDescription(c)}
                            </a>
                        </td>
                        <td class="number">
                            <a href=${clusterLink(c)}>
                                ${metric(c, 'presubmitRejects')}
                            </a>
                        </td>
                        <td class="number">
                            <a href=${clusterLink(c)}>
                                ${metric(c, 'testRunFailures')}
                            </a>
                        </td>
                        <td class="number">
                            <a href=${clusterLink(c)}>
                                ${metric(c, 'failures')}
                            </a>
                        </td>
                    </tr>`)}
                </tbody>
            </table>
        </div>`;
    }
    static styles = [css`
        #container {
            margin: 20px 14px;
        }
        h1 {
            font-size: 18px;
            font-weight: normal;
        }
        table {
            border-collapse: collapse;
            max-width: 100%;
        }
        th {
            font-weight: normal;
            color: var(--greyed-out-text-color);
            text-align: left;
            font-size: var(--font-size-small);
        }
        th.sortable {
            cursor: pointer;
        }
        td, th {
            padding: 4px;
            max-width: 80%;
        }
        td.number {
            text-align: right;
        }
        td a {
            display: block;
            text-decoration: none;
            color: var(--default-text-color);
        }
        tbody tr:hover {
            background-color: var(--light-active-color);
        }
        .failure-reason {
            word-break: break-all;
            font-size: var(--font-size-small);
        }
        a[data-suggested] {
            font-style: italic;
        }
        `];
}

type MetricName = 'presubmitRejects' | 'testRunFailures' | 'failures';

// Cluster is the cluster information sent by the server.
interface Cluster {
    clusterId: ClusterId;
    presubmitRejects1d: Counts;
    presubmitRejects3d: Counts;
    presubmitRejects7d: Counts;
    testRunFailures1d: Counts;
    testRunFailures3d: Counts;
    testRunFailures7d: Counts;
    failures1d: Counts;
    failures3d: Counts;
    failures7d: Counts;
    exampleFailureReason: string;
    exampleTestId: string;
}

interface ClusterId {
    algorithm: string;
    id: string;
}

interface Counts {
    nominal: number;
    preExoneration: number;
    residual: number;
    residualPreExoneration: number;
}
