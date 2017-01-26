(function() {
  'use strict';

  const DEFAULT_SNOOZE_TIME_MIN = 60;
  const ONE_MIN_MS = 1000 * 60;

  Polymer({
    is: 'som-annotations',
    behaviors: [AnnotationManagerBehavior],
    properties: {
      // All alert annotations. Includes values from localState.
      annotations: {
        notify: true,
        type: Object,
        value: function() {
          return {};
        },
        computed: '_computeAnnotations(_annotationsResp, localState)'
      },
      annotationError: {
        type: Object,
        value: function() {
          return {};
        },
      },
      // The raw response from the server of annotations.
      _annotationsResp: {
        type: Array,
        value: function() {
          return [];
        },
      },
      _bugErrorMessage: String,
      _bugInput: {
        type: Object,
        value: function() {
          return this.$.bug;
        }
      },
      _bugModel: Object,
      bugQueueLabel: String,
      _commentsErrorMessage: String,
      _commentsModel: Object,
      _commentsModelAnnotation: {
        type: Object,
        computed:
            '_computeCommentsModelAnnotation(annotations, _commentsModel)',
      },
      _commentsHidden: {
        type: Boolean,
        computed: '_computeCommentsHidden(_commentsModelAnnotation)',
      },
      _commentTextInput: {
        type: Object,
        value: function() {
          return this.$.commentText;
        }
      },
      _filedBug: {
        type: Boolean,
        value: false,
      },
      _snoozeErrorMessage: String,
      _snoozeModel: Object,
      _snoozeTimeInput: {
        type: Object,
        value: function() {
          return this.$.snoozeTime;
        }
      },
      user: String,
      xsrfToken: String,
    },

    ready: function() {
      this.fetchAnnotations();
    },

    fetch: function() {
      this.annotationError.action = 'Fetching all annotations';
      this.fetchAnnotations().catch((error) => {
        let old = this.annotationError;
        this.annotationError.message = error;
        this.notifyPath('annotationError.message');
      });
    },

    // Fetches new annotations from the server.
    fetchAnnotations: function() {
      return window.fetch('/api/v1/annotations', {credentials: 'include'})
          .then(jsonParsePromise)
          .then((req) => {
            this._annotationsResp = [];
            this._annotationsResp = req;
          });
    },

    // Send an annotation change. Also updates the in memory annotation
    // database.
    // Returns a promise of the POST request to the server to carry out the
    // annotation change.
    sendAnnotation: function(key, type, change) {
      let data = {
        xsrf_token: this.xsrfToken,
        data: change,
      };

      return this
          .postJSON(
              '/api/v1/annotations/' + encodeURIComponent(key) + '/' + type,
              change)
          .then(jsonParsePromise)
          .then(this._postResponse.bind(this));
    },

    // FIXME: Move to common behavior if other code does POST requests.
    postJSON: function(url, data, options) {
      options = options || {};
      options.body = JSON.stringify({
        xsrf_token: this.xsrfToken,
        data: data,
      });
      options.method = 'POST';
      options.credentials = 'include';
      return new Promise((resolve, reject) => {
        window.fetch(url, options).then((value) => {
          if (!value.ok) {
            value.text().then((txt) => {
              if (!(value.status == 403 && txt.includes('token expired'))) {
                reject(txt);
                return;
              }

              // We need to refresh our XSRF token!
              window.fetch('/api/v1/xsrf_token', {credentials: 'include'})
                  .then((respData) => {
                    return respData.json();
                  })
                  .then((jsonData) => {
                    // Clone options because sinon.spy args from different calls
                    // to window.fetch clobber each other in this scenario.
                    let opts = JSON.parse(JSON.stringify(options));
                    this.xsrfToken = jsonData['token'];
                    opts.body = JSON.stringify({
                      xsrf_token: this.xsrfToken,
                      data: data,
                    });
                    window.fetch(url, opts).then(resolve, reject);
                  });
            });
            return;
          }

          resolve(value);
        }, reject);
      })

    },

    _computeAnnotations: function(annotationsJson, localState) {
      let annotations = {};
      if (!annotationsJson) {
        annotationsJson = [];
      }

      Object.keys(localState).forEach((key) => {
        key = decodeURIComponent(key);
        annotations[key] = localState[key];
      });
      annotationsJson.forEach((annotation) => {
        // If we've already added something here through local state, copy that
        // over.
        let key = decodeURIComponent(annotation.key);
        if (annotations[key]) {
          Object.assign(annotation, annotations[key]);
        }
        annotations[key] = annotation;
      });
      return annotations;
    },

    _haveAnnotationError: function(annotationError) {
      return !!annotationError.base.message;
    },

    // Handle the result of the response of a post to the server.
    _postResponse: function(response) {
      let annotations = this.annotations;
      annotations[decodeURIComponent(response.key)] = response;
      let newArray = [];
      Object.keys(annotations).forEach((k) => {
        k = decodeURIComponent(k);
        newArray.push(annotations[k]);
      });
      this._annotationsResp = newArray;

      return response;
    },

    ////////////////////// Handlers ///////////////////////////

    handleOpenedChange: function(evt) {
      this.setLocalStateKey(evt.target.alert.key, {opened: evt.detail.value});
    },

    handleAnnotation: function(evt) {
      this.annotationError.action = 'Fetching all annotations';
      this.sendAnnotation(
              evt.target.alert.key, evt.detail.type, evt.detail.change)
          .then((response) => {
            this.setLocalStateKey(response.key, {opened: false});
          })
          .catch((error) => {
            let old = this.annotationError;
            this.annotationError.message = error;
            this.notifyPath('annotationError.message');
          });
    },

    handleComment: function(evt) {
      this._commentsModel = evt.target.alert;
      this._commentsErrorMessage = '';
      this.$.commentsDialog.open();
    },

    handleLinkBug: function(evt) {
      this._bugModel = evt.target.alert;
      this.$.fileBugLink.href =
          'https://bugs.chromium.org/p/chromium/issues/entry?status=Available&labels=' +
          this.bugQueueLabel + '&summary=' + this._bugModel.title +
          '&comment=' + encodeURIComponent(this._commentForBug(this._bugModel));
      this._filedBug = false;
      this._bugErrorMessage = '';
      this.$.bugDialog.open();
    },

    handleSnooze: function(evt) {
      this._snoozeModel = evt.target.alert;
      this.$.snoozeTime.value = DEFAULT_SNOOZE_TIME_MIN;
      this._snoozeErrorMessage = '';
      this.$.snoozeDialog.open();
    },

    ////////////////////// Bugs ///////////////////////////

    _commentForBug: function(bugModel) {
      let result = bugModel.title + '\n\n';
      if (bugModel.extension) {
        if (bugModel.extension.builders &&
            bugModel.extension.builders.length > 0) {
          result += 'Builders failed on: ';
          for (let i in bugModel.extension.builders) {
            result += '\n- ' + bugModel.extension.builders[i].name + ': \n  ' +
                bugModel.extension.builders[i].url;
          }
          result += '\n\n';
        }
        if (bugModel.extension.reasons &&
            bugModel.extension.reasons.length > 0) {
          result += 'Reasons: ';
          for (let i in bugModel.extension.reasons) {
            result += '\n' + bugModel.extension.reasons[i].url;
            if (bugModel.extension.reasons[i].test_names) {
              result += '\n' +
                  'Tests:';
              if (bugModel.extension.reasons[i].test_names) {
                result += '\n* ' +
                    bugModel.extension.reasons[i].test_names.join('\n* ');
              }
            }
          }
          result += '\n\n';
        }
      }
      return result;
    },

    _fileBugClicked: function() {
      this._filedBug = true;
    },

    _saveBug: function() {
      // TODO(add proper error handling)
      this.sendAnnotation(this._bugModel.key, 'add', {bugs: [this.$.bug.value]})
          .then(
              (response) => {
                this._bugErrorMessage = '';
                this.$.bug.value = '';
                this.$.bugDialog.close();

                this.setLocalStateKey(response.key, {opened: false});
              },
              (error) => {
                this._bugErrorMessage = error;
              });
    },

    ////////////////////// Snooze ///////////////////////////

    _snooze: function() {
      // TODO(add proper error handling)
      this.sendAnnotation(
              this._snoozeModel.key, 'add',
              {snoozeTime: Date.now() + ONE_MIN_MS * this.$.snoozeTime.value})
          .then(
              (response) => {
                this.$.snoozeTime.value = '';
                this.$.snoozeDialog.close();

                this.setLocalStateKey(response.key, {opened: false});
              },
              (error) => {
                this._snoozeErrorMessage = error;
              });
    },

    ////////////////////// Comments ///////////////////////////

    _addComment: function() {
      let text = this.$.commentText.value;
      if (!(text && /[^\s]/.test(text))) {
        this._commentsErrorMessage = 'Comment text cannot be blank!';
        return;
      }
      this.sendAnnotation(this._commentsModel.key, 'add', {
            comments: [text],
          })
          .then(
              (response) => {
                this.$.commentText.value = '';
                this._commentsErrorMessage = '';
                this.setLocalStateKey(response.key, {opened: false});
              },
              (error) => {
                this._commentsErrorMessage = error;
              });
    },

    _computeCommentsHidden: function(annotation) {
      return !(annotation && annotation.comments);
    },

    // This is mostly to make sure the comments in the modal get updated
    // properly if changed.
    _computeCommentsModelAnnotation: function(annotations, model) {
      if (!annotations || !model) {
        return null;
      }
      return this.computeAnnotation(annotations, model);
    },

    _computeHideDeleteComment(comment) {
      return comment.user != this.user;
    },

    _computeUsername(email) {
      if (!email) {
        return email;
      }
      let cutoff = email.indexOf('@');
      if (cutoff < 0) {
        return email;
      }
      return email.substring(0, cutoff);
    },

    _formatTimestamp: function(time) {
      if (!time) {
        return '';
      }
      return new Date(time).toLocaleString(false, {timeZoneName: 'short'});
    },

    _removeComment: function(evt) {
      let request = this.sendAnnotation(this._commentsModel.key, 'remove', {
        comments: [evt.model.index],
      });
      if (request) {
        request.then(
            (response) => {
              this.setLocalStateKey(response.key, {opened: false});
            },
            (error) => {
              this._commentsErrorMessage = error;
            });
      }
    },
  })
})();
