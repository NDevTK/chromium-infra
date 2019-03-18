import {assert} from 'chai';
import sinon from 'sinon';

import ClientLogger from './client-logger.js';
import MonorailTSMon from './monorail-ts-mon.js';

suite('ClientLogger', () => {
  const startedKey = 'ClientLogger.rutabaga.started';
  let c;

  setup(() => {
    window.CS_env = {
      token: 'rutabaga-token',
      tokenExpiresSec: 1234,
      app_version: 'rutabaga-version',
    };
    window.chops = {rpc: {PrpcClient: sinon.spy()}};
    window.ga = sinon.spy();
    MonorailTSMon.prototype.disableAfterNextFlush = sinon.spy();
    c = new ClientLogger('rutabaga');
  });

  teardown(() => {
    sessionStorage.clear();
  });

  suite('constructor', () => {
    test('assigns this.category', () => {
      assert.equal(c.category, 'rutabaga');
    });

    test('gets started events from sessionStorage', () => {
      const startedEvents = {
        event1: {
          time: 12345678,
          labels: ['label1', 'label2'],
        },
        event2: {
          time: 87654321,
          labels: ['label2'],
        },
      };
      sessionStorage[startedKey] = JSON.stringify(startedEvents);

      c = new ClientLogger('rutabaga');
      assert.deepEqual(startedEvents, c.startedEvents);
    });
  });

  suite('records ts_mon metrics', () => {
    let issueCreateMetric;
    let issueUpdateMetric;
    let autocompleteMetric;
    let c;

    setup(() => {
      window.ga = sinon.spy();
      c = new ClientLogger('issues');
      issueCreateMetric = c.tsMon._userTimingMetrics[0].metric;
      issueCreateMetric.add = sinon.spy();

      issueUpdateMetric = c.tsMon._userTimingMetrics[1].metric;
      issueUpdateMetric.add = sinon.spy();

      autocompleteMetric = c.tsMon._userTimingMetrics[2].metric;
      autocompleteMetric.add = sinon.spy();
    });

    test('bogus', () => {
      c.logStart('rutabaga');
      c.logEnd('rutabaga');
      sinon.assert.notCalled(issueCreateMetric.add);
      sinon.assert.notCalled(issueUpdateMetric.add);
      sinon.assert.notCalled(autocompleteMetric.add);
    });

    test('new-issue', () => {
      c.logStart('new-issue', 'server-time');
      c.logEnd('new-issue', 'server-time');
      sinon.assert.notCalled(issueUpdateMetric.add);
      sinon.assert.notCalled(autocompleteMetric.add);

      sinon.assert.calledOnce(issueCreateMetric.add);
      assert.isNumber(issueCreateMetric.add.getCall(0).args[0]);
      assert.isString(issueCreateMetric.add.getCall(0).args[1].get('client_id'));
      assert.equal(issueCreateMetric.add.getCall(0).args[1].get('host_name'),
        'rutabaga-version');
    });

    test('issue-update', () => {
      c.logStart('issue-update', 'computer-time');
      c.logEnd('issue-update', 'computer-time');
      sinon.assert.notCalled(issueCreateMetric.add);
      sinon.assert.notCalled(autocompleteMetric.add);

      sinon.assert.calledOnce(issueUpdateMetric.add);
      assert.isNumber(issueUpdateMetric.add.getCall(0).args[0]);
      assert.isString(issueUpdateMetric.add.getCall(0).args[1].get('client_id'));
      assert.equal(issueUpdateMetric.add.getCall(0).args[1].get('host_name'),
        'rutabaga-version');
    });

    test('populate-options', () => {
      c.logStart('populate-options');
      c.logEnd('populate-options');
      sinon.assert.notCalled(issueCreateMetric.add);
      sinon.assert.notCalled(issueUpdateMetric.add);
      // Autocomplete is not called in issues category.
      sinon.assert.notCalled(autocompleteMetric.add);

      c = new ClientLogger('autocomplete');
      autocompleteMetric = c.tsMon._userTimingMetrics[2].metric;
      autocompleteMetric.add = sinon.spy();

      c.logStart('populate-options', 'user-time');
      c.logEnd('populate-options', 'user-time');
      sinon.assert.notCalled(issueCreateMetric.add);
      sinon.assert.notCalled(issueUpdateMetric.add);

      sinon.assert.calledOnce(autocompleteMetric.add);
      assert.isNumber(autocompleteMetric.add.getCall(0).args[0]);
      assert.isString(autocompleteMetric.add.getCall(0).args[1].get('client_id'));
      assert.equal(autocompleteMetric.add.getCall(0).args[1].get('host_name'),
        'rutabaga-version');
    });
  });

  suite('logStart', () => {
    let c, clock;
    const currentTime = 5000;

    setup(() => {
      c = new ClientLogger('rutabaga');
      clock = sinon.useFakeTimers(currentTime);
    });

    teardown(() => {
      clock.restore();
      sessionStorage.clear();
    });

    test('creates a new startedEvent if none', () => {
      c.logStart('event-name', 'event-label');

      sinon.assert.calledOnce(ga);
      sinon.assert.calledWith(ga, 'send', 'event', 'rutabaga',
        'event-name-start', 'event-label');

      const expectedStartedEvents = {
        'event-name': {
          time: currentTime,
          labels: {
            'event-label': currentTime,
          },
        },
      };
      assert.deepEqual(c.startedEvents, expectedStartedEvents);
      assert.deepEqual(JSON.parse(sessionStorage[startedKey]),
        expectedStartedEvents);
    });

    test('uses an existing startedEvent', () => {
      c.startedEvents['event-name'] = {
        time: 1234,
        labels: {
          'event-label': 1000,
        },
      };
      c.logStart('event-name', 'event-label');

      sinon.assert.calledOnce(ga);
      sinon.assert.calledWith(ga, 'send', 'event', 'rutabaga',
        'event-name-start', 'event-label');

      // TODO(jeffcarp): Audit is this wanted behavior? Replacing event time
      // but not label time?
      const expectedStartedEvents = {
        'event-name': {
          time: 1234,
          labels: {
            'event-label': currentTime,
          },
        },
      };
      assert.deepEqual(c.startedEvents, expectedStartedEvents);
      assert.deepEqual(JSON.parse(sessionStorage[startedKey]),
        expectedStartedEvents);
    });
  });

  suite('logPause', () => {
    const startTime = 1234;
    const currentTime = 5000;
    let c, clock;

    setup(() => {
      clock = sinon.useFakeTimers(currentTime);
      c = new ClientLogger('rutabaga');
      c.startedEvents['event-name'] = {
        time: startTime,
        labels: {
          'event-label': startTime,
        },
      };
    });

    teardown(() => {
      clock.restore();
      sessionStorage.clear();
    });

    test('throws if no label given', () => {
      assert.throws(() => {
        c.logPause('bogus');
      }, 'event with no label');
    });

    test('exits early if no start event exists', () => {
      const originalStartedEvents = Object.assign(c.startedEvents, {});
      c.logPause('bogus', 'fogus');
      assert.deepEqual(c.startedEvents, originalStartedEvents);
    });

    test('exits early if no label exists', () => {
      const originalStartedEvents = Object.assign(c.startedEvents, {});
      c.logPause('event-name', 'fogus');
      assert.deepEqual(c.startedEvents, originalStartedEvents);
    });

    test('adds elapsed time to start event', () => {
      c.logPause('event-name', 'event-label');

      const expectedStartedEvents = {
        'event-name': {
          time: startTime,
          labels: {
            'event-label': startTime,
          },
          elapsed: {
            'event-label': currentTime - startTime,
          },
        },
      };
      assert.deepEqual(c.startedEvents, expectedStartedEvents);
      assert.deepEqual(
        JSON.parse(sessionStorage['ClientLogger.rutabaga.started']),
        expectedStartedEvents);
    });
  });

  suite('logResume', () => {
    let c, clock;
    const startTimeEvent = 1234;
    const startTimeLabel = 2345;
    const labelElapsed = 4321;
    const currentTime = 6000;

    setup(() => {
      clock = sinon.useFakeTimers(currentTime);
      c = new ClientLogger('rutabaga');
      c.startedEvents['event-name'] = {
        time: startTimeEvent,
        labels: {
          'event-label': startTimeLabel,
        },
        elapsed: {
          'event-label': labelElapsed,
        },
      };
    });

    teardown(() => {
      clock.restore();
      sessionStorage.clear();
    });

    test('throws if no label given', () => {
      assert.throws(() => {
        c.logResume('bogus');
      }, 'no label');
    });

    test('exits early if no start event exists', () => {
      const originalStartedEvents = Object.assign(c.startedEvents, {});
      c.logResume('bogus', 'fogus');
      assert.deepEqual(c.startedEvents, originalStartedEvents);
    });

    test('exits early if the label was never paused', () => {
      c.startedEvents['event-name'] = {
        time: startTimeEvent,
        labels: {
          'event-label': startTimeLabel,
        },
        elapsed: {},
      };

      const originalStartedEvents = Object.assign(c.startedEvents, {});
      c.logResume('event-name', 'event-label');
      assert.deepEqual(c.startedEvents, originalStartedEvents);
    });

    test('sets start event time to current time', () => {
      c.logResume('event-name', 'event-label');

      const expectedStartedEvents = {
        'event-name': {
          time: startTimeEvent,
          labels: {
            'event-label': currentTime,
          },
          elapsed: {
            'event-label': labelElapsed,
          },
        },
      };
      assert.deepEqual(c.startedEvents, expectedStartedEvents);
      assert.deepEqual(
        JSON.parse(sessionStorage['ClientLogger.rutabaga.started']),
        expectedStartedEvents);
    });
  });

  suite('logEnd', () => {
    let c, clock;
    const startTimeEvent = 1234;
    const startTimeLabel1 = 2345;
    const startTimeLabel2 = 3456;
    const currentTime = 10000;

    setup(() => {
      c = new ClientLogger('rutabaga');
      clock = sinon.useFakeTimers(currentTime);
      c.tsMon.recordUserTiming = sinon.spy();
      c.startedEvents = {
        someEvent: {
          time: startTimeEvent,
          labels: {
            label1: startTimeLabel1,
            label2: startTimeLabel2,
          },
        },
      };
    });

    teardown(() => {
      clock.restore();
    });

    test('returns early if no event was started', () => {
      c.startedEvents = {someEvent: {}};
      const originalStartedEvents = Object.assign(c.startedEvents, {});
      c.logEnd('bogus');
      sinon.assert.notCalled(window.ga);
      assert.isNull(sessionStorage.getItem(startedKey));
      assert.deepEqual(c.startedEvents, originalStartedEvents);
    });

    test('returns early if label was not started', () => {
      c.startedEvents = {someEvent: {labels: {}}};
      const originalStartedEvents = Object.assign(c.startedEvents, {});
      c.logEnd('someEvent', 'bogus');
      sinon.assert.notCalled(window.ga);
      assert.isNull(sessionStorage.getItem(startedKey));
      assert.deepEqual(c.startedEvents, originalStartedEvents);
    });

    test('does not log non-labeled events over threshold', () => {
      c.startedEvents = {someEvent: {time: currentTime - 1000}};
      c.logEnd('someEvent', null, 999);

      sinon.assert.calledOnce(window.ga);
      sinon.assert.calledWith(window.ga, 'send', 'event', 'rutabaga',
        'someEvent-end', null, undefined);
      sinon.assert.notCalled(c.tsMon.recordUserTiming);
      assert.equal(sessionStorage.getItem(startedKey), '{}');
    });

    test('does not log labeled events over threshold', () => {
      const elapsedLabel2 = 2000;
      c.startedEvents.someEvent.elapsed = {
        label1: currentTime - 1000,
        label2: elapsedLabel2,
      };
      c.logEnd('someEvent', 'label1', 999);

      sinon.assert.calledOnce(window.ga);
      sinon.assert.calledWith(window.ga, 'send', 'event', 'rutabaga',
        'someEvent-end', 'label1', undefined);
      // TODO(jeffcarp): Feature: add GA event if over threshold.
      sinon.assert.notCalled(c.tsMon.recordUserTiming);

      const expectedStartedEvents = {
        someEvent: {
          time: startTimeEvent,
          labels: {
            label2: startTimeLabel2,
          },
          elapsed: {
            label2: elapsedLabel2,
          },
        },
      };
      assert.deepEqual(c.startedEvents, expectedStartedEvents);
      assert.deepEqual(JSON.parse(sessionStorage[startedKey]),
        expectedStartedEvents);
    });

    test('calls ga() with timing and event info for all labels', () => {
      const label1Elapsed = 1000;
      const label2Elapsed = 2500;
      c.startedEvents.someEvent.elapsed = {
        label1: label1Elapsed,
        label2: label2Elapsed,
      };
      c.logEnd('someEvent');

      assert.deepEqual(ga.getCall(0).args, [
        'send', 'timing', {
          timingCategory: 'rutabaga',
          timingValue: currentTime - startTimeEvent,
          timingVar: 'someEvent',
        }]);

      assert.deepEqual(ga.getCall(1).args, [
        'send', 'timing', {
          timingCategory: 'rutabaga',
          timingValue: (currentTime - startTimeLabel1) + label1Elapsed,
          timingVar: 'someEvent',
          timingLabel: 'label1',
        }]);
      assert.deepEqual(ga.getCall(2).args, [
        'send', 'timing', {
          timingCategory: 'rutabaga',
          timingValue: (currentTime - startTimeLabel2) + label2Elapsed,
          timingVar: 'someEvent',
          timingLabel: 'label2',
        }]);
      assert.deepEqual(ga.getCall(3).args, [
        'send', 'event', 'rutabaga', 'someEvent-end', undefined, undefined,
      ]);
      assert.deepEqual(c.tsMon.recordUserTiming.getCall(0).args, [
        'rutabaga', 'someEvent', null, currentTime - startTimeEvent,
      ]);
      assert.deepEqual(c.tsMon.recordUserTiming.getCall(1).args, [
        'rutabaga', 'someEvent', 'label1',
        (currentTime - startTimeLabel1) + label1Elapsed,
      ]);
      assert.deepEqual(c.tsMon.recordUserTiming.getCall(2).args, [
        'rutabaga', 'someEvent', 'label2',
        (currentTime - startTimeLabel2) + label2Elapsed,
      ]);

      assert.deepEqual(c.startedEvents, {});
      assert.equal(sessionStorage.getItem(startedKey), '{}');
    });

    test('calling with a label calls ga() only for that label', () => {
      const label1Elapsed = 1000;
      const label2Elapsed = 2500;
      c.startedEvents.someEvent.elapsed = {
        label1: label1Elapsed,
        label2: label2Elapsed,
      };
      c.logEnd('someEvent', 'label2');

      assert.deepEqual(ga.getCall(0).args, [
        'send', 'timing', {
          timingCategory: 'rutabaga',
          timingValue: (currentTime - startTimeLabel2) + label2Elapsed,
          timingVar: 'someEvent',
          timingLabel: 'label2',
        }]);
      assert.deepEqual(window.ga.getCall(1).args, [
        'send', 'event', 'rutabaga', 'someEvent-end', 'label2', undefined,
      ]);
      sinon.assert.calledOnce(c.tsMon.recordUserTiming);
      sinon.assert.calledWith(c.tsMon.recordUserTiming, 'rutabaga',
        'someEvent', 'label2', (currentTime - startTimeLabel2) + label2Elapsed);

      const expectedStartedEvents = {
        someEvent: {
          time: startTimeEvent,
          labels: {
            label1: startTimeLabel1,
          },
          elapsed: {
            label1: label1Elapsed,
          },
        },
      };
      assert.deepEqual(c.startedEvents, expectedStartedEvents);
      assert.deepEqual(JSON.parse(sessionStorage[startedKey]),
        expectedStartedEvents);
    });

    test('calling logStart, logPause, logResume, and logEnd works for labels', () => {
      let countedElapsedTime = 0;
      c.logStart('someEvent', 'label1');
      clock.tick(1000);
      countedElapsedTime += 1000;
      c.logPause('someEvent', 'label1');
      clock.tick(1000);
      c.logResume('someEvent', 'label1');
      clock.tick(1000);
      countedElapsedTime += 1000;
      c.logEnd('someEvent', 'label1');

      assert.deepEqual(ga.getCall(0).args, [
        'send', 'event', 'rutabaga', 'someEvent-start', 'label1', undefined,
      ]);
      assert.deepEqual(ga.getCall(1).args, [
        'send', 'timing', {
          timingCategory: 'rutabaga',
          timingValue: countedElapsedTime,
          timingVar: 'someEvent',
          timingLabel: 'label1',
        }]);
      assert.deepEqual(window.ga.getCall(2).args, [
        'send', 'event', 'rutabaga', 'someEvent-end', 'label1', undefined,
      ]);
      sinon.assert.calledOnce(c.tsMon.recordUserTiming);
      sinon.assert.calledWith(c.tsMon.recordUserTiming, 'rutabaga',
        'someEvent', 'label1', countedElapsedTime);

      const expectedStartedEvents = {
        someEvent: {
          time: startTimeEvent,
          labels: {
            label2: startTimeLabel2,
          },
          elapsed: {},
        },
      };
      assert.deepEqual(c.startedEvents, expectedStartedEvents);
      assert.deepEqual(JSON.parse(sessionStorage[startedKey]),
        expectedStartedEvents);
    });

    test('logs some events when others are above threshold', () => {
      c.startedEvents = {
        someEvent: {
          time: 9500,
          labels: {
            overThresholdWithoutElapsed: 8000,
            overThresholdWithElapsed: 9500,
            underThresholdWithoutElapsed: 9750,
            underThresholdWithElapsed: 9650,
            exactlyOnThresholdWithoutElapsed: 9001,
            exactlyOnThresholdWithElapsed: 9002,
          },
          elapsed: {
            overThresholdWithElapsed: 1000,
            underThresholdWithElapsed: 100,
            exactlyOnThresholdWithElapsed: 1,
          },
        },
      };
      c.logEnd('someEvent', null, 999);

      // Verify ga() calls.
      assert.equal(window.ga.getCalls().length, 6);
      assert.deepEqual(ga.getCall(0).args, [
        'send', 'timing', {
          timingCategory: 'rutabaga',
          timingValue: 500,
          timingVar: 'someEvent',
        }]);
      assert.deepEqual(ga.getCall(1).args, [
        'send', 'timing', {
          timingCategory: 'rutabaga',
          timingValue: 250,
          timingVar: 'someEvent',
          timingLabel: 'underThresholdWithoutElapsed',
        }]);
      assert.deepEqual(ga.getCall(2).args, [
        'send', 'timing', {
          timingCategory: 'rutabaga',
          timingValue: 450,
          timingVar: 'someEvent',
          timingLabel: 'underThresholdWithElapsed',
        }]);
      assert.deepEqual(ga.getCall(3).args, [
        'send', 'timing', {
          timingCategory: 'rutabaga',
          timingValue: 999,
          timingVar: 'someEvent',
          timingLabel: 'exactlyOnThresholdWithoutElapsed',
        }]);
      assert.deepEqual(ga.getCall(4).args, [
        'send', 'timing', {
          timingCategory: 'rutabaga',
          timingValue: 999,
          timingVar: 'someEvent',
          timingLabel: 'exactlyOnThresholdWithElapsed',
        }]);
      assert.deepEqual(ga.getCall(5).args, [
        'send', 'event', 'rutabaga', 'someEvent-end', null, undefined,
      ]);

      // Verify ts_mon.recordUserTiming() calls.
      assert.equal(c.tsMon.recordUserTiming.getCalls().length, 5);
      assert.deepEqual(c.tsMon.recordUserTiming.getCall(0).args, [
        'rutabaga', 'someEvent', null, 500,
      ]);
      assert.deepEqual(c.tsMon.recordUserTiming.getCall(1).args, [
        'rutabaga', 'someEvent', 'underThresholdWithoutElapsed', 250,
      ]);
      assert.deepEqual(c.tsMon.recordUserTiming.getCall(2).args, [
        'rutabaga', 'someEvent', 'underThresholdWithElapsed', 450,
      ]);
      assert.deepEqual(c.tsMon.recordUserTiming.getCall(3).args, [
        'rutabaga', 'someEvent', 'exactlyOnThresholdWithoutElapsed', 999,
      ]);
      assert.deepEqual(c.tsMon.recordUserTiming.getCall(4).args, [
        'rutabaga', 'someEvent', 'exactlyOnThresholdWithElapsed', 999,
      ]);
      assert.deepEqual(c.startedEvents, {});
      assert.deepEqual(JSON.parse(sessionStorage[startedKey]), {});
    });
  });
});
