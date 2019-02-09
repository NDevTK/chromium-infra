// Copyright 2019 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

import {assert} from 'chai';
import {MrAccountDropdown} from './mr-account-dropdown.js';

let element;

suite('mr-account-dropdown', () => {
  setup(() => {
    element = document.createElement('mr-account-dropdown');
    document.body.appendChild(element);
  });

  teardown(() => {
    document.body.removeChild(element);
  });

  test('initializes', () => {
    assert.instanceOf(element, MrAccountDropdown);
  });
});
