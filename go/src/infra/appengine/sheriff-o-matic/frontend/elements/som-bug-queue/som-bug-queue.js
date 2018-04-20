'use strict';

const UNSET_PRIORITY = Number.MAX_SAFE_INTEGER;

class SomBugQueue extends Polymer.Element {

  static get is() {
    return 'som-bug-queue';
  }

  static get properties() {
    return {
      bugQueueLabel: {
        type: String,
        observer: '_bugQueueLabelChanged',
      },
      bugs: {
        type: Array,
        notify: true,
        computed: '_computeBugs(_bugQueueJson, _uncachedBugsJson)',
      },
      treeDisplayName: String,
      _activeRequests: Object,
      _bugsByPriority: {
        type: Array,
        computed: '_computeBugsByPriority(bugs)',
      },
      _bugQueueJson: {
        type: Object,
        value: null,
      },
      _bugQueueJsonError: {
        type: Object,
        value: null,
      },
      _bugsLoaded: {
        type: Boolean,
        value: false,
      },
      _defaultOpenState: {
        type: Boolean,
        value: true,
      },
      _hideBugQueue: {
        type: Boolean,
        value: true,
        computed: '_computeHideBugQueue(bugQueueLabel)',
      },
      _showNoBugs: {
        type: Boolean,
        value: false,
        computed: '_computeShowNoBugs(bugs, _bugsLoaded, _bugQueueJsonError)',
      },
      _toggleSectionIcon: {
        type: String,
        computed: '_computeToggleSectionIcon(_opened)',
      },
      _opened: {
        type: Boolean,
        value: true,
      },
      _uncachedBugsJson: {
        type: Object,
        value: null,
      },
      _uncachedBugsJsonError: {
        type: Object,
        value: null,
      },
    };
  }

  ready() {
    super.ready();

    // This is to expose the UNSET_PRIORITY constant for use in unit testing.
    this.UNSET_PRIORITY = UNSET_PRIORITY;
  }

  refresh() {
    if (this._hideBugQueue) {
      return;
    }

    if (this._activeRequests) {
      this._activeRequests.forEach((req) => {
        req.abort();
      });
    }

    let requests = [this.$.bugQueueAjax.generateRequest()];

    let promises = requests.map((r) => {
      return r.completes;
    });

    this._activeRequests = requests;
    Promise.all(promises).then(() => {
      this._bugsLoaded = true;
    });
  }

  _bugQueueLabelChanged() {
    this._bugQueueJson = null;
    this._bugQueueJsonError = null;

    this._uncachedBugsJson = null;
    this._uncachedBugsJsonError = null;

    this._bugsLoaded = false;

    this.refresh();
  }

  _computeBugs(bugQueueJson, uncachedBugsJson) {
    let hasBugJson = bugQueueJson && bugQueueJson.items;
    let hasUncachedJson = uncachedBugsJson && uncachedBugsJson.items;
    if (!hasBugJson && !hasUncachedJson) {
      return [];
    } else if (!hasUncachedJson) {
      return bugQueueJson.items;
    }
    return uncachedBugsJson.items;
  }

  _computeBugsByPriority(bugs) {
    // update last updated time as relative time
    for(var i = 0; i < bugs.length; i++) {
      if (bugs[i].updated) {
        bugs[i].updated = moment.tz(bugs[i].updated, 'Atlantic/Reykjavik').fromNow();
      }
    }
    let buckets = bugs.reduce((function(obj, b) {
                                let p = this._computePriority(b);
                                if (!(p in obj)) {
                                  obj[p] = [b];
                                } else {
                                  obj[p].push(b);
                                }
                                return obj;
                              }).bind(this),
                              {});

    // Flatten the buckets into an array for use in dom-repeat.
    let result = Object.keys(buckets).sort().map(function(key) {
      return {'priority': key, 'bugs': buckets[key]};
    });
    return result;
  }

  _computeHideBugQueue(bugQueueLabel) {
    // No loading or empty message is shown unless a bug queue exists.
    return !bugQueueLabel || bugQueueLabel === '' ||
           bugQueueLabel === 'Performance-Sheriff-BotHealth';
  }

  _computePriority(bug) {
    if (!bug || !bug.labels) {
      return this.UNSET_PRIORITY;
    }
    for (let i in bug.labels) {
      let match = bug.labels[i].match(/^Pri-(\d)$/);
      if (match) {
        let result = parseInt(match[1]);
        return result !== NaN ? result : this.UNSET_PRIORITY;
      }
    }
    return this.UNSET_PRIORITY;
  }

  _computeShowNoBugs(bugs, bugsLoaded, error) {
    // Show the "No bugs" message only when the queue is done loading
    return bugsLoaded && this._haveNoBugs(bugs) && this._haveNoErrors(error);
  }

  _filterBugLabels(labels, bugQueueLabel) {
    if (!labels) {
      return [];
    }
    bugQueueLabel = bugQueueLabel || '';
    return labels.filter((label) => {
      return label.toLowerCase() != bugQueueLabel.toLowerCase() &&
             !label.match(/^Pri-(\d)$/);
    });
  }

  _haveNoBugs(bugs) {
    return !bugs || bugs.length == 0;
  }

  _haveNoErrors(error) {
    return !error;
  }

  _priorityText(pri) {
    if (this._validPriority(pri)) {
      return `Priority ${pri}`;
    }
    return 'No Priority';
  }

  _showBugsLoading(bugsLoaded, error) {
    return !bugsLoaded && this._haveNoErrors(error);
  }

  _validPriority(pri) {
    return pri != this.UNSET_PRIORITY;
  }

  ////////////////////// Collapsing by priority ///////////////////////////

  _computeCollapseId(pri) {
    return `collapsePri${pri}`;
  }

  _computeCollapseIcon(opened) {
    return opened ? 'remove' : 'add';
  }

  _collapseAll() {
    for (let i = 0; i < this._bugsByPriority.length; i++) {
      let pri = this._bugsByPriority[i].priority;
      let id = this._computeCollapseId(pri);
      let collapse = this.shadowRoot.querySelector('#' + id);

      collapse.opened = false;
      this.shadowRoot.querySelector('#toggleIconPri' + pri).icon =
          this._computeCollapseIcon(collapse.opened);
    }
  }

  _expandAll() {
    for (let i = 0; i < this._bugsByPriority.length; i++) {
      let pri = this._bugsByPriority[i].priority;
      let id = this._computeCollapseId(pri);
      let collapse = this.shadowRoot.querySelector('#' + id);

      collapse.opened = true;
      this.shadowRoot.querySelector('#toggleIconPri' + pri).icon =
          this._computeCollapseIcon(collapse.opened);
    }

    this._opened = true;
  }

  _togglePriorityCollapse(evt) {
    let i = evt.model.get('index');
    let pri = this._bugsByPriority[i].priority;
    let id = this._computeCollapseId(pri);
    let collapse = this.shadowRoot.querySelector('#' + id);
    if (!collapse) {
      console.error(id + ' is not a valid Id.');
    } else {
      collapse.toggle();

      this.shadowRoot.querySelector('#toggleIconPri' + pri).icon =
          this._computeCollapseIcon(collapse.opened);
    }
  }

  ////////////////////// Collapsing the section ///////////////////////////

  _computeToggleSectionIcon(opened) {
    return opened ? 'unfold-less' : 'unfold-more';
  }

  _toggleSection() {
    this._opened = !this._opened;
  }
}

customElements.define(SomBugQueue.is, SomBugQueue);
