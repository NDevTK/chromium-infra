// Copyright 2019 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

import {LitElement, html, css} from 'lit-element';

import 'elements/issue-detail/mr-flipper/mr-flipper.js';
import 'elements/chops/chops-dialog/chops-dialog.js';
import 'elements/chops/chops-timestamp/chops-timestamp.js';
import {store, connectStore} from 'elements/reducers/base.js';
import * as issue from 'elements/reducers/issue.js';
import * as project from 'elements/reducers/project.js';
import {SHARED_STYLES} from 'elements/shared/shared-styles.js';
import 'elements/framework/links/mr-user-link/mr-user-link.js';
import 'elements/framework/links/mr-crbug-link/mr-crbug-link.js';
import 'elements/framework/mr-code-font-toggle/mr-code-font-toggle.js';
import 'elements/framework/mr-dropdown/mr-dropdown.js';
import {ISSUE_EDIT_PERMISSION, ISSUE_DELETE_PERMISSION,
  ISSUE_FLAGSPAM_PERMISSION} from 'elements/shared/permissions.js';
import {prpcClient} from 'prpc-client-instance.js';

const DELETE_ISSUE_CONFIRMATION_NOTICE = `\
Normally, you would just close issues by setting their status to a closed value.
Are you sure you want to delete this issue?`;


/**
 * `<mr-issue-header>`
 *
 * The header for a given launch issue.
 *
 */
export class MrIssueHeader extends connectStore(LitElement) {
  static get styles() {
    return [
      SHARED_STYLES,
      css`
        :host {
          width: 100%;
          margin-top: 0;
          font-size: var(--chops-large-font-size);
          background-color: var(--monorail-metadata-toggled-bg);
          border-bottom: var(--chops-normal-border);
          padding: 0.25em 8px;
          box-sizing: border-box;
          display: flex;
          flex-direction: row;
          justify-content: space-between;
          align-items: center;
        }
        h1 {
          font-size: 100%;
          line-height: 140%;
          font-weight: bolder;
          padding: 0;
          margin: 0;
        }
        mr-flipper {
          border-left: var(--chops-normal-border);
          padding-left: 8px;
          margin-left: 4px;
          font-size: var(--chops-main-font-size);
        }
        .issue-actions {
          min-width: fit-content;
          display: flex;
          flex-direction: row;
          align-items: center;
        }
        .issue-actions a {
          color: var(--chops-link-color);
          cursor: pointer;
        }
        .issue-actions a:hover {
          text-decoration: underline;
        }
        .code-font-and-description-edit {
          min-width: fit-content;
          display: flex;
          flex-direction: column;
          align-items: flex-end;
          font-size: var(--chops-main-font-size);
        }
        .code-font-and-description-edit div {
          min-width: 70px;
          display: flex;
          justify-content: space-between;
        }
        .code-font-and-description-edit a {
          white-space: nowrap;
        }
        .spam-notice {
          display: inline-flex;
          align-items: center;
          justify-content: center;
          padding: 1px 6px;
          border-radius: 3px;
          background: #F44336;
          color: white;
          font-weight: bold;
          font-size: var(--chops-main-font-size);
          margin-right: 4px;
        }
        .byline {
          display: block;
          font-size: var(--chops-main-font-size);
          width: 100%;
          line-height: 140%;
          color: var(--chops-text-color);
        }
        .main-text-outer {
          flex-basis: 100%;
          display: flex;
          justify-content: flex-start;
          flex-direction: row;
          align-items: center;
        }
        .main-text {
          flex-basis: 100%;
        }
        @media (max-width: 840px) {
          :host {
            flex-wrap: wrap;
            justify-content: center;
          }
          .main-text {
            width: 100%;
            margin-bottom: 0.5em;
          }
        }
      `,
    ];
  }

  render() {
    return html`
      <div class="main-text-outer">
        <div class="main-text">
          <h1>
            ${this.issue.isSpam ? html`
              <span class="spam-notice">Spam</span>
            `: ''}
            Issue ${this.issue.localId}: ${this.issue.summary}
          </h1>
          <small class="byline">
            Reported by
            <mr-user-link .userRef=${this.issue.reporterRef}></mr-user-link>
            on <chops-timestamp .timestamp=${this.issue.openedTimestamp}></chops-timestamp>
          </small>
        </div>
      </div>
      <div class="issue-actions">
        <div class="code-font-and-description-edit">
          <div>
            <mr-crbug-link .issue=${this.issue}></mr-crbug-link>
            <mr-code-font-toggle
              .userDisplayName=${this.userDisplayName}
            ></mr-code-font-toggle>
          </div>
          ${this._canEditIssue ? html`
            <a @click=${this._openEditDescription}>Edit description</a>
          `: ''}
        </div>
        ${this._issueOptions.length ? html`
          <mr-dropdown
            .items=${this._issueOptions}
            icon="more_vert"
          ></mr-dropdown>
        ` : ''}
        <mr-flipper></mr-flipper>
      </div>
    `;
  }

  static get properties() {
    return {
      userDisplayName: {type: String},
      issue: {type: Object},
      issuePermissions: {type: Object},
      isRestricted: {type: Boolean},
      projectTemplates: {type: Array},
      _action: {type: String},
      _targetProjectError: {type: String},
    };
  }

  constructor() {
    super();
    this.issuePermissions = [];
    this.projectTemplates = [];
    this.issue = {};
    this.isRestricted = false;
  }

  stateChanged(state) {
    this.issue = issue.issue(state);
    this.issuePermissions = issue.permissions(state);
    this.projectTemplates = project.project(state).templates;

    const restrictions = issue.restrictions(state);
    this.isRestricted = restrictions && Object.keys(restrictions).length;
  }

  get _canEditIssue() {
    return this.issuePermissions.includes(ISSUE_EDIT_PERMISSION);
  }

  get _issueOptions() {
    // We create two edit Arrays for the top and bottom half of the menu,
    // to be separated by a separator in the UI.
    const editOptions = [];
    const riskyOptions = [];
    const isSpam = this.issue.isSpam;
    const isRestricted = this.isRestricted;

    const permissions = this.issuePermissions;
    const templates = this.projectTemplates;


    if (permissions.includes(ISSUE_FLAGSPAM_PERMISSION)) {
      const text = (isSpam ? 'Un-flag' : 'Flag') + ' issue as spam';
      riskyOptions.push({
        text,
        handler: this._markIssue.bind(this),
      });
    }
    if (permissions.includes(ISSUE_DELETE_PERMISSION)) {
      riskyOptions.push({
        text: 'Delete issue',
        handler: this._deleteIssue.bind(this),
      });
      if (!isRestricted) {
        editOptions.push({
          text: 'Move issue',
          handler: this._openMoveCopyIssue.bind(this, 'Move'),
        });
        editOptions.push({
          text: 'Copy issue',
          handler: this._openMoveCopyIssue.bind(this, 'Copy'),
        });
      }
    }

    if (permissions.includes(ISSUE_EDIT_PERMISSION) && templates.length) {
      editOptions.push({
        text: 'Convert issue template',
        handler: this._openConvertIssue.bind(this),
      });
    }

    if (editOptions.length && riskyOptions.length) {
      editOptions.push({separator: true});
    }
    return editOptions.concat(riskyOptions);
  }

  _markIssue() {
    prpcClient.call('monorail.Issues', 'FlagIssues', {
      issueRefs: [{
        projectName: this.issue.projectName,
        localId: this.issue.localId,
      }],
      flag: !this.issue.isSpam,
    }).then(() => {
      const message = {
        issueRef: {
          projectName: this.issue.projectName,
          localId: this.issue.localId,
        },
      };
      store.dispatch(issue.fetch(message));
    });
  }

  _deleteIssue() {
    const ok = confirm(DELETE_ISSUE_CONFIRMATION_NOTICE);
    if (ok) {
      prpcClient.call('monorail.Issues', 'DeleteIssue', {
        issueRef: {
          projectName: this.issue.projectName,
          localId: this.issue.localId,
        },
        delete: true,
      }).then(() => {
        const message = {
          issueRef: {
            projectName: this.issue.projectName,
            localId: this.issue.localId,
          },
        };
        store.dispatch(issue.fetch(message));
      });
    }
  }

  _openEditDescription() {
    this.dispatchEvent(new CustomEvent('open-dialog', {
      bubbles: true,
      composed: true,
      detail: {
        dialogId: 'edit-description',
        fieldName: '',
      },
    }));
  }

  _openMoveCopyIssue(action) {
    this.dispatchEvent(new CustomEvent('open-dialog', {
      bubbles: true,
      composed: true,
      detail: {
        dialogId: 'move-copy-issue',
        action,
      },
    }));
  }

  _openConvertIssue() {
    this.dispatchEvent(new CustomEvent('open-dialog', {
      bubbles: true,
      composed: true,
      detail: {
        dialogId: 'convert-issue',
      },
    }));
  }
}

customElements.define('mr-issue-header', MrIssueHeader);
