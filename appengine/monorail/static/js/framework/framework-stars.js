/* Copyright 2016 The Chromium Authors. All Rights Reserved.
 *
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file or at
 * https://developers.google.com/open-source/licenses/bsd
 */

/**
 * This file contains JS functions that support setting and showing
 * stars throughout Monorail.
 */


/**
 * The character to display when the user has starred an issue.
 */
var TKR_STAR_ON = '\u2605';


/**
 * The character to display when the user has not starred an issue.
 */
var TKR_STAR_OFF = '\u2606';


/**
 * Function to toggle the star on an issue.  Does both an update of the
 * DOM and hit the server to record the star.
 *
 * @param {Element} el The star <a> element.
 * @param {String} projectName name of the project to be starred, or name of
 *                 the project containing the issue to be starred.
 * @param {Integer} localId number of the issue to be starred.
 * @param {String} projectName number of the user to be starred.
 * @param {string} token The security token.
 */
function TKR_toggleStar(el, projectName, localId, userId, hotlistId, token) {
  var starred = (el.textContent.trim() == TKR_STAR_OFF) ? 1 : 0;
  TKR_toggleStarLocal(el);

  var type;
  if (userId) type = 'users';
  if (projectName) type = 'projects';
  if (projectName && localId) type = 'issues';
  if (hotlistId) type = 'hotlists';

  args = {'starred': starred};
  if (type == 'issues') {
    url = '/p/' + projectName + '/issues/setstar.do';
    args['id'] = localId;
  } else {
    url = '/hosting/stars.do';
    args['scope'] = type;
    if (type == 'hotlists') {
      args['item'] = hotlistId;
    } else if (type == 'projects'){
      args['item'] = projectName;
    } else {
      args['item'] =  userId;
    };
  };

  TKR_setStar(el, url, args, token, url);
}


/**
 * Just update the display state of a star, without contacting the server.
 * Optionally update the value of a form element as well. Useful for when
 * a user is entering a new issue and wants to set its initial starred state.
 * @param {Element} el Star <img> element.
 * @param {string} opt_formElementId HTML ID of the hidden form element for
 *      stars.
 */
function TKR_toggleStarLocal(el, opt_formElementId) {
  var starred = (el.textContent.trim() == TKR_STAR_OFF) ? 1 : 0;

  el.textContent = starred ? TKR_STAR_ON : TKR_STAR_OFF;
  el.style.color = starred ? 'cornflowerblue' : 'grey';
  el.title = starred ? 'You have starred this item' : 'Click to star this item';

  if (opt_formElementId) {
    $(opt_formElementId).value = '' + starred; // convert to string
  }
}


/**
 * Send the new star state to the server and create a callback for its response.
 * @param {Element} el The star <a> element.
 * @param {String} url The server URL to post to.
 * @param {Dict} args The arguments to send in the POST request.
 * @param {String} opt_token The security token to send in the request.
 */
function TKR_setStar(el, url, args, opt_token) {
  if (opt_token) {
    CS_doPost(url, function(event) { TKR_gotSetStar(el, event); }, args,
              opt_token, url);
  } else {
    CS_doPost(url, function(event) { TKR_gotSetStar(el, event); }, args);
  }
}


/**
 * Evaluates the server response after a starring operation completed.
 * @param {Element} el <a> element containing the star which was clicked.
 * @param {event} event with xhr JSON response from the server.
 */
function TKR_gotSetStar(el, event) {
  var xhr = event.target;
  if (xhr.readyState != 4 || xhr.status != 200)
    return;
  var args = CS_parseJSON(xhr);
  var localStarred = (el.textContent.trim() == TKR_STAR_ON) ? 1 : 0;
  var serverStarred = args['starred'];
  if (localStarred != serverStarred) {
    TKR_toggleStarLocal(el);
  }
}


/**
 * When we show two star icons on the same details page, keep them
 * in sync with each other. And, update a message about starring
 * that is displayed near the issue update form.
 * @param {Element} clickedStar The star that the user clicked on.
 * @param {string} otherStarId ID of the other star icon.
 */
function TKR_syncStarIcons(clickedStar, otherStarId) {
  var otherStar = document.getElementById(otherStarId);
  if (!otherStar) {
    return;
  }
  TKR_toggleStarLocal(otherStar);

  var vote_feedback = document.getElementById('vote_feedback');
  if (!vote_feedback) {
    return;
  }

  if (clickedStar.textContent == TKR_STAR_OFF) {
    vote_feedback.textContent =
        'Vote for this issue and get email change notifications.';
  } else {
    vote_feedback.textContent = 'Your vote has been recorded.';
  }
}


// Exports
_TKR_toggleStar = TKR_toggleStar;
_TKR_toggleStarLocal = TKR_toggleStarLocal;
_TKR_syncStarIcons = TKR_syncStarIcons;
