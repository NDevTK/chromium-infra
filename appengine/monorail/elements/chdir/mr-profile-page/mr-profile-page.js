/* Copyright 2019 The Chromium Authors. All Rights Reserved.
 *
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

import '../../../node_modules/@polymer/polymer/polymer-legacy.js';
import {html} from '../../../node_modules/@polymer/polymer/lib/utils/html-tag.js';
import {PolymerElement} from '../../../node_modules/@polymer/polymer/polymer-element.js';

import '../../mr-header/mr-header.js';
import '../mr-activity-table/mr-activity-table.js';
import '../mr-comments-list/mr-comments-list.js';
import '../mr-commit-table/mr-commit-table.js';

/**
 * `<mr-profile-page>`
 *
 * The main entry point for a Monorail Polymer profile.
 *
 */
export class MrProfilePage extends PolymerElement {
  static get template() {
    return html`
      <style>
        .history-container {
          width: 70%;
          padding: 1em 16px;
          padding-top: 35px;
          padding-bottom: 75px;
          padding-left: 30px;
          padding-right: 30px;
          display: flex;
          flex-direction: column;
          flex-wrap: no-wrap;
          min-height: 100%;
        }
        .dataTable {
          width: 80%;
          overflow-x: auto;
          margin-bottom: 55px;
          margin-left: 30px;
          margin-right: 30px;
          max-height: 300px;
        }
        .activityGraph {
          width: 80%;
          overflow-x: auto;
          margin-bottom: 55px;
          margin-left: 30px;
          margin-right: 30px;
          max-height: 300px;
        }
        .metadata-container {
          font-size: 85%;
          background: hsl(120, 35%, 95%);
          border: 1px solid hsl(120, 15%, 90%);
          width: 15%;
          min-width: 256px;
          flex-grow: 0;
          flex-shrink: 0;
          margin-right: 16px;
          box-sizing: border-box;
          min-height: 100%;
        }
        .container-outside {
          box-sizing: border-box;
          width: 100%;
          max-width: 100%;
          margin: auto;
          padding: 0.75em 8px;
          display: flex;
          align-items: stretch;
          justify-content: space-between;
          flex-direction: row;
          flex-wrap: no-wrap;
          flex-grow: 0;
          min-height: 100%;
        }
        .profile-data {
          text-align: center;
          padding-top: 40%;
          font-size: 110%;
        }
        .commitTable {
          width: 80%;
          overflow-x: auto;
          height: 400px;
        }
      </style>
      <app-location query-params="{{queryParams}}" url-space-regex="^/u/(.*)/polymer$"></app-location>

      <mr-header
        user-display-name="[[user]]"
        login-url="[[loginUrl]]"
        logout-url="[[logoutUrl]]"
      >
        <span slot="subheader">
          &gt; Viewing Profile: [[viewedUser]]
        </span>
      </mr-header>
      <div class="container-outside">
        <div class="metadata-container">
          <div class="profile-data">
            [[viewedUser]] <br>
            <b>Last visit:</b> [[lastVisitStr]] <br>
            <b>Starred Developers:</b> [[_checkStarredUsers(starredUsers)]]
          </div>
        </div>
        <div class="history-container">
          <template is="dom-if" if="[[!_hideActivityTracker]]">
            <mr-activity-table
              class="activityGraph"
              user="[[viewedUser]]"
              viewed-user-id="[[viewedUserId]]"
              commits="[[commits]]"
              comments="[[comments]]"
              selected-date="{{selectedDate}}"
            ></mr-activity-table>
          </template>
          <div class="dataTable">
            <mr-commit-table
              user="[[viewedUser]]"
              commits="[[commits]]"
              selected-date="{{selectedDate}}">
            </mr-commit-table>
          </div>
          <div class="dataTable">
            <mr-comments-list
              user="[[viewedUser]]"
              viewed-user-id="[[viewedUserId]]"
              comments="[[comments]]"
              selected-date="{{selectedDate}}">
            </mr-comments-list>
          </div>
        </div>
      </div>
    `;
  }

  static get is() {
    return 'mr-profile-page';
  }

  static get properties() {
    return {
      user: {
        type: String,
        observer: '_getUserData',
      },
      logoutUrl: String,
      loginUrl: String,
      viewedUser: String,
      viewedUserId: Number,
      lastVisitStr: String,
      starredUsers: Array,
      commits: {
        type: Array,
      },
      comments: {
        type: Array,
      },
      selectedDate: {
        type: Number,
      },
      _hideActivityTracker: {
        type: Boolean,
        computed: '_computeHideActivityTracker(user, viewedUser)',
      },
    };
  }

  _checkStarredUsers(list) {
    if (list.length != 0) {
      return list;
    } else {
      return ['None'];
    }
  }

  _getUserData() {
    const d = new Date();
    const n = d.getTime();
    let currentTime = n / 1000;
    currentTime = Math.ceil(currentTime);
    let fromTime = currentTime - 31536000;

    const commitMessage = {
      email: this.viewedUser,
      fromTimestamp: fromTime,
      untilTimestamp: currentTime,
    };

    const getCommits = window.prpcCall(
      'monorail.Users', 'GetUserCommits', commitMessage
    );

    getCommits.then((resp) => {
      this.commits = resp.userCommits;
    }, (error) => {
    });

    const commentMessage = {
      userRef: {
        userId: this.viewedUserId,
      },
    };

    const listActivities = window.prpcCall(
      'monorail.Issues', 'ListActivities', commentMessage
    );

    listActivities.then(
      (resp) => {
        this.comments = resp.comments;
      },
      (error) => {}
    );
  }

  _computeHideActivityTracker(user, viewedUser) {
    return user !== viewedUser;
  }
}

customElements.define(MrProfilePage.is, MrProfilePage);
