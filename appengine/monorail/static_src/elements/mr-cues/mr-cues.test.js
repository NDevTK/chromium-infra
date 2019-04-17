// Copyright 2019 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

import {assert} from 'chai';
import {MrCues} from './mr-cues.js';
import page from 'page';

let element;

suite('mr-cues', () => {
  setup(() => {
    element = document.createElement('mr-cues');
    document.body.appendChild(element);

    sinon.stub(page, 'call');
  });

  teardown(() => {
    document.body.removeChild(element);

    page.call.restore();
  });

  test('initializes', () => {
    assert.instanceOf(element, MrCues);
  });

  test('map state', () => {
    const state = {user: {currentUser:
      {prefs: new Map(), prefsLoaded: false}}};
    assert.deepEqual(
      {prefs: new Map(), prefsLoaded: false},
      MrCues.mapStateToProps(state, element));
  });


  test('anon does not see privacy dialog', () => {
    assert.isFalse(element._computeShowPrivacyDialog(
      'unused', 'unused', 'unused'));
  });

  test('signed in user sees no privacy dialog before prefs load', () => {
    element.userDisplayName = 'user@example.com';
    assert.isFalse(element._computeShowPrivacyDialog(
      'unused', false, 'unused'));
  });

  test('signed in user sees no privacy dialog if already dismissed', () => {
    element.userDisplayName = 'user@example.com';
    assert.isFalse(element._computeShowPrivacyDialog(
      'unused', true, true));
  });

  test('signed in user sees no privacy dialog if dismissal pref set', () => {
    element.userDisplayName = 'user@example.com';
    assert.isFalse(element._computeShowPrivacyDialog(
      new Map([['privacy_click_through', 'true']]), true, false));
  });

  test('signed in user sees privacy dialog if dismissal pref missing', () => {
    element.userDisplayName = 'user@example.com';
    assert.isTrue(element._computeShowPrivacyDialog(
      new Map(), true, false));
  });

  test('show or close privacy dialog', () => {
    sinon.spy(HTMLDialogElement.prototype, 'showModal');
    sinon.spy(HTMLDialogElement.prototype, 'close');
    element._showPrivacyDialogChanged(true);
    assert.isTrue(HTMLDialogElement.prototype.showModal.calledOnce);
    element._showPrivacyDialogChanged(false);
    assert.isTrue(HTMLDialogElement.prototype.close.calledOnce);
  });

});
