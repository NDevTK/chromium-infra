// Copyright 2019 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

import {LitElement, html, css} from 'lit-element';

import page from 'page';
import {store, connectStore} from 'reducers/base.js';
import * as hotlist from 'reducers/hotlist.js';
import * as sitewide from 'reducers/sitewide.js';
import * as userV0 from 'reducers/userV0.js';

import 'elements/chops/chops-button/chops-button.js';
import 'elements/hotlist/mr-hotlist-header/mr-hotlist-header.js';

/** Hotlist Settings page */
class _MrHotlistSettingsPage extends LitElement {
  /** @override */
  static get styles() {
    return css`
      :host {
        display: block;
      }
      section {
        margin: 16px 24px;
      }
      h2 {
        font-weight: normal;
      }
      dt {
        font-weight: bold;
      }
      dd {
        margin: 0;
      }
      div {
        margin: 16px 24px;
      }
    `;
  }

  /** @override */
  render() {
    return html`
      <mr-hotlist-header selected=2></mr-hotlist-header>
      ${this._hotlist ? this._renderPage() : 'Loading...'}
    `;
  }

  /**
   * @return {TemplateResult}
   */
  _renderPage() {
    const defaultColumns = this._hotlist.defaultColumns
        .map((col) => col.column).join(' ');
    return html`
      <section>
        <h2>Hotlist Settings</h2>
        <dl>
          <dt>Name</dt>
          <dd>${this._hotlist.displayName}</dd>
          <dt>Summary</dt>
          <dd>${this._hotlist.summary}</dd>
          <dt>Description</dt>
          <dd>${this._hotlist.description}</dd>
        </dl>
      </section>

      <section>
        <h2>Hotlist Defaults</h2>
        <dl>
          <dt>Default columns shown in list view</dt>
          <dd>${defaultColumns}</dd>
        </dl>
      </section>

      <section>
        <h2>Hotlist Access</h2>
        <dl>
          <dt>Who can view this hotlist</dt>
          <dd>
            ${this._hotlist.hotlistPrivacy ?
              'Anyone on the internet' : 'Members only'}
          </dd>
        </dl>
        <p>
          Individual issues in the list can only be seen by users who can
          normally see them. The privacy status of an issue is considered
          when it is being displayed (or not displayed) in a hotlist.
      </section>

      <div>
        <chops-button @click=${this._delete} id="delete-hotlist">
          Delete hotlist
        </chops-button>
      </div>
    `;
  }

  /** @override */
  static get properties() {
    return {
      _currentUser: {type: Object},
      _hotlist: {type: Object},
    };
  }

  /** @override */
  constructor() {
    super();
    /** @type {?HotlistV1} */
    this._hotlist = null;

    // Expose page.js for test stubbing.
    this.page = page;
  }

  /** Deletes the hotlist, dispatching an action to Redux. */
  async _delete() {}
};

/** Redux-connected version of _MrHotlistSettingsPage. */
export class MrHotlistSettingsPage
  extends connectStore(_MrHotlistSettingsPage) {
  /** @override */
  stateChanged(state) {
    this._hotlist = hotlist.viewedHotlist(state);
    this._currentUser = userV0.currentUser(state);
  }

  /** @override */
  updated(changedProperties) {
    super.updated(changedProperties);

    if (changedProperties.has('_hotlist') && this._hotlist) {
      const pageTitle = 'Settings - ' + this._hotlist.displayName;
      store.dispatch(sitewide.setPageTitle(pageTitle));
      const headerTitle = 'Hotlist ' + this._hotlist.displayName;
      store.dispatch(sitewide.setHeaderTitle(headerTitle));
    }
  }

  /** @override */
  async _delete() {
    if (confirm(
        'Are you sure you want to delete this hotlist? This cannot be undone.')
    ) {
      const action = hotlist.deleteHotlist(this._hotlist.name);
      await store.dispatch(action);

      // TODO(crbug/monorail/7430): Handle an error and add <chops-snackbar>.
      // Note that this will redirect regardless of an error.
      this.page(`/u/${this._currentUser.displayName}/hotlists`);
    }
  }
}

customElements.define('mr-hotlist-settings-page-base', _MrHotlistSettingsPage);
customElements.define('mr-hotlist-settings-page', MrHotlistSettingsPage);
