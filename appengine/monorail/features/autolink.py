# Copyright 2016 The Chromium Authors. All rights reserved.
# Use of this source code is govered by a BSD-style
# license that can be found in the LICENSE file or at
# https://developers.google.com/open-source/licenses/bsd

"""Autolink helps auto-link references to artifacts in text.

This class maintains a registry of artifact autolink syntax specs and
callbacks. The structure of that registry is:
  { component_name: (lookup_callback,
                     { regex: substitution_callback, ...}),
    ...
  }

For example:
  { 'tracker':
     (GetReferencedIssues,
      ExtractProjectAndIssueIds,
      {_ISSUE_REF_RE: ReplaceIssueRef}),
    'versioncontrol':
     (GetReferencedRevisions,
      ExtractProjectAndRevNum,
      {_GIT_HASH_RE: ReplaceRevisionRef}),
  }

The dictionary of regexes is used here because, in the future, we
might add more regexes for each component rather than have one complex
regex per component.
"""

import logging
import re
import urllib
import urlparse

import settings
from framework import template_helpers
from framework import validate
from proto import project_pb2
from tracker import tracker_helpers


_CLOSING_TAG_RE = re.compile('</[a-z0-9]+>$', re.IGNORECASE)

_LINKIFY_SCHEMES = r'(https?://|ftp://|mailto:)'
# Also count a start-tag '<' as a url delimeter, since the autolinker
# is sometimes run against html fragments.
_IS_A_LINK_RE = re.compile(r'(%s)([^\s<]+)' % _LINKIFY_SCHEMES, re.UNICODE)

# These are allowed in links, but if any of closing delimiters appear
# at the end of the link, and the opening one is not part of the link,
# then trim off the closing delimiters.
_LINK_TRAILING_CHARS = [
    (None, ':'),
    (None, '.'),
    (None, ','),
    ('<', '>'),
    ('"', '"'),
    ('(', ')'),
    ('[', ']'),
    ('{', '}'),
    ]


def Linkify(_mr, autolink_regex_match,
            _component_ref_artifacts):
  """Examine a textual reference and replace it with a hyperlink or not.

  This is a callback for use with the autolink feature.

  Args:
    _mr: common info parsed from the user HTTP request.
    autolink_regex_match: regex match for the textual reference.
    _component_ref_artifacts: unused value

  Returns:
    A list of TextRuns with tag=a for all matched ftp, http, https and mailto
    links converted into HTML hyperlinks.
  """
  hyperlink = autolink_regex_match.group(0)

  trailing = ''
  for begin, end in _LINK_TRAILING_CHARS:
    if hyperlink.endswith(end):
      if not begin or hyperlink[:-len(end)].find(begin) == -1:
        trailing = end + trailing
        hyperlink = hyperlink[:-len(end)]

  tag_match = _CLOSING_TAG_RE.search(hyperlink)
  if tag_match:
    trailing = hyperlink[tag_match.start(0):] + trailing
    hyperlink = hyperlink[:tag_match.start(0)]

  if (not validate.IsValidURL(hyperlink) and
      not validate.IsValidEmail(hyperlink)):
    return [template_helpers.TextRun(autolink_regex_match.group(0))]

  result = [template_helpers.TextRun(hyperlink, tag='a', href=hyperlink)]
  if trailing:
    result.append(template_helpers.TextRun(trailing))

  return result


# Regular expression to detect git hashes.
# Used to auto-link to Git hashes on crrev.com when displaying issue details.
# Matches "rN", "r#N", and "revision N" when "rN" is not part of a larger word
# and N is a hexadecimal string of 40 chars.
_GIT_HASH_RE = re.compile(
    r'\b(?P<prefix>r(evision\s+#?)?)?(?P<revnum>([a-f0-9]{40}))\b',
    re.IGNORECASE | re.MULTILINE)

# This is for SVN revisions and Git commit posisitons.
_SVN_REF_RE = re.compile(
    r'\b(?P<prefix>r(evision\s+#?)?)(?P<revnum>([0-9]{1,7}))\b',
    re.IGNORECASE | re.MULTILINE)


def GetReferencedRevisions(_mr, _refs):
  """Load the referenced revision objects."""
  # For now we just autolink any revision hash without actually
  # checking that such a revision exists,
  # TODO(jrobbins): Hit crrev.com and check that the revision exists
  # and show a rollover with revision info.
  return None


def ExtractRevNums(_mr, autolink_regex_match):
  """Return internal representation of a rev reference."""
  ref = autolink_regex_match.group('revnum')
  logging.debug('revision ref = %s', ref)
  return [ref]


def ReplaceRevisionRef(
    mr, autolink_regex_match, _component_ref_artifacts):
  """Return HTML markup for an autolink reference."""
  prefix = autolink_regex_match.group('prefix')
  revnum = autolink_regex_match.group('revnum')
  url = _GetRevisionURLFormat(mr.project).format(revnum=revnum)
  content = revnum
  if prefix:
    content = '%s%s' % (prefix, revnum)
  return [template_helpers.TextRun(content, tag='a', href=url)]


def _GetRevisionURLFormat(project):
  # TODO(jrobbins): Expose a UI to customize it to point to whatever site
  # hosts the source code. Also, site-wide default.
  return (project.revision_url_format or settings.revision_url_format)


# Regular expression to detect issue references.
# Used to auto-link to other issues when displaying issue details.
# Matches "issue " when "issue" is not part of a larger word, or
# "issue #", or just a "#" when it is preceeded by a space.
_ISSUE_REF_RE = re.compile(r"""
    (?P<prefix>\b(issues?|bugs?)[ \t]*(:|=)?)
    ([ \t]*(?P<project_name>\b[-a-z0-9]+[:\#])?
     (?P<number_sign>\#?)
     (?P<local_id>\d+)\b
     (,?[ \t]*(and|or)?)?)+""", re.IGNORECASE | re.VERBOSE)

_SINGLE_ISSUE_REF_RE = re.compile(r"""
    (?P<prefix>\b(issue|bug)[ \t]*)?
    (?P<project_name>\b[-a-z0-9]+[:\#])?
    (?P<number_sign>\#?)
    (?P<local_id>\d+)\b""", re.IGNORECASE | re.VERBOSE)


def CurryGetReferencedIssues(services):
  """Return a function to get ref'd issues with these persist objects bound.

  Currying is a convienent way to give the callback access to the persist
  objects, but without requiring that all possible persist objects be passed
  through the autolink registry and functions.

  Args:
    services: connection to issue, config, and project persistence layers.

  Returns:
    A ready-to-use function that accepts the arguments that autolink
    expects to pass to it.
  """

  def GetReferencedIssues(mr, ref_tuples):
    """Return lists of open and closed issues referenced by these comments.

    Args:
      mr: commonly used info parsed from the request.
      ref_tuples: list of (project_name, local_id) tuples for each issue
          that is mentioned in the comment text. The project_name may be None,
          in which case the issue is assumed to be in the current project.

    Returns:
      A list of open and closed issue dicts.
    """
    ref_projects = services.project.GetProjectsByName(
        mr.cnxn,
        [(ref_pn or mr.project_name) for ref_pn, _ in ref_tuples])
    issue_ids = services.issue.ResolveIssueRefs(
        mr.cnxn, ref_projects, mr.project_name, ref_tuples)
    open_issues, closed_issues = (
        tracker_helpers.GetAllowedOpenedAndClosedIssues(
            mr, issue_ids, services))

    open_dict = {}
    for issue in open_issues:
      open_dict[_IssueProjectKey(issue.project_name, issue.local_id)] = issue

    closed_dict = {}
    for issue in closed_issues:
      closed_dict[_IssueProjectKey(issue.project_name, issue.local_id)] = issue

    logging.info('autolinking dicts %r and %r', open_dict, closed_dict)

    return open_dict, closed_dict

  return GetReferencedIssues


def _ParseProjectNameMatch(project_name):
  """Process the passed project name and determine the best representation.

  Args:
    project_name: a string with the project name matched in a regex

  Returns:
    A minimal representation of the project name, None if no valid content.
  """
  if not project_name:
    return None
  return project_name.lstrip().rstrip('#: \t\n')


def ExtractProjectAndIssueIds(_mr, autolink_regex_match):
  """Convert a regex match for a textual reference into our internal form."""
  whole_str = autolink_regex_match.group(0)
  refs = []
  for submatch in _SINGLE_ISSUE_REF_RE.finditer(whole_str):
    ref = (_ParseProjectNameMatch(submatch.group('project_name')),
           int(submatch.group('local_id')))
    refs.append(ref)
    logging.info('issue ref = %s', ref)

  return refs


# This uses project name to avoid a lookup on project ID in a function
# that has no services object.
def _IssueProjectKey(project_name, local_id):
  """Make a dictionary key to identify a referenced issue."""
  return '%s:%d' % (project_name, local_id)


class IssueRefRun(object):
  """A text run that links to a referenced issue."""

  def __init__(self, issue, is_closed, project_name, prefix):
    self.tag = 'a'
    self.css_class = 'closed_ref' if is_closed else None
    self.title = issue.summary
    self.href = '/p/%s/issues/detail?id=%d' % (project_name, issue.local_id)

    self.content = '%s%d' % (prefix, issue.local_id)
    if is_closed and not prefix:
      self.content = ' %s ' % self.content


def ReplaceIssueRef(mr, autolink_regex_match, component_ref_artifacts):
  """Examine a textual reference and replace it with an autolink or not.

  Args:
    mr: commonly used info parsed from the request
    autolink_regex_match: regex match for the textual reference.
    component_ref_artifacts: result of earlier call to GetReferencedIssues.

  Returns:
    A list of IssueRefRuns and TextRuns to replace the textual
    reference.  If there is an issue to autolink to, we return an HTML
    hyperlink.  Otherwise, we the run will have the original plain
    text.
  """
  open_dict, closed_dict = component_ref_artifacts
  original = autolink_regex_match.group(0)
  logging.info('called ReplaceIssueRef on %r', original)
  result_runs = []
  pos = 0
  for submatch in _SINGLE_ISSUE_REF_RE.finditer(original):
    if submatch.start() >= pos:
      if original[pos: submatch.start()]:
        result_runs.append(template_helpers.TextRun(
            original[pos: submatch.start()]))
      replacement_run = _ReplaceSingleIssueRef(
          mr, submatch, open_dict, closed_dict)
      result_runs.append(replacement_run)
      pos = submatch.end()

  if original[pos:]:
    result_runs.append(template_helpers.TextRun(original[pos:]))

  return result_runs


def _ReplaceSingleIssueRef(mr, submatch, open_dict, closed_dict):
  """Replace one issue reference with a link, or the original text."""
  prefix = submatch.group('prefix') or ''
  project_name = submatch.group('project_name')
  if project_name:
    prefix += project_name
    project_name = project_name.lstrip().rstrip(':#')
  else:
    # We need project_name for the URL, even if it is not in the text.
    project_name = mr.project_name

  number_sign = submatch.group('number_sign')
  if number_sign:
    prefix += number_sign
  local_id = int(submatch.group('local_id'))
  issue_key = _IssueProjectKey(project_name or mr.project_name, local_id)

  if issue_key in open_dict:
    return IssueRefRun(open_dict[issue_key], False, project_name, prefix)
  elif issue_key in closed_dict:
    return IssueRefRun(closed_dict[issue_key], True, project_name, prefix)
  else:  # Don't link to non-existent issues.
    return template_helpers.TextRun('%s%d' % (prefix, local_id))


class Autolink(object):
  """Maintains a registry of autolink syntax and can apply it to comments."""

  def __init__(self):
    self.registry = {}

  def RegisterComponent(self, component_name, artifact_lookup_function,
                        match_to_reference_function, autolink_re_subst_dict):
    """Register all the autolink info for a software component.

    Args:
      component_name: string name of software component, must be unique.
      artifact_lookup_function: function to batch lookup all artifacts that
          might have been referenced in a set of comments:
          function(all_matches) -> referenced_artifacts
          the referenced_artifacts will be pased to each subst function.
      match_to_reference_function: convert a regex match object to
          some internal representation of the artifact reference.
      autolink_re_subst_dict: dictionary of regular expressions and
          the substitution function that should be called for each match:
          function(match, referenced_artifacts) -> replacement_markup
    """
    self.registry[component_name] = (artifact_lookup_function,
                                     match_to_reference_function,
                                     autolink_re_subst_dict)

  def GetAllReferencedArtifacts(self, mr, comment_text_list):
    """Call callbacks to lookup all artifacts possibly referenced.

    Args:
      mr: information parsed out of the user HTTP request.
      comment_text_list: list of comment content strings.

    Returns:
      Opaque object that can be pased to MarkupAutolinks.  It's
      structure happens to be {component_name: artifact_list, ...}.
    """
    all_referenced_artifacts = {}
    for comp, (lookup, match_to_refs, re_dict) in self.registry.iteritems():
      refs = set()
      for comment_text in comment_text_list:
        for regex in re_dict:
          for match in regex.finditer(comment_text):
            additional_refs = match_to_refs(mr, match)
            if additional_refs:
              refs.update(additional_refs)

      all_referenced_artifacts[comp] = lookup(mr, refs)

    return all_referenced_artifacts

  def MarkupAutolinks(self, mr, text_runs, all_referenced_artifacts):
    """Loop over components and regexes, applying all substitutions.

    Args:
      mr: info parsed from the user's HTTP request.
      text_runs: List of text runs for the user's comment.
      all_referenced_artifacts: result of previous call to
        GetAllReferencedArtifacts.

    Returns:
      List of text runs for the entire user comment, some of which may have
      attribures that cause them to render as links in render-rich-text.ezt.
    """
    items = self.registry.items()
    items.sort()  # Process components in determinate alphabetical order.
    for component, (_lookup, _match_ref, re_subst_dict) in items:
      component_ref_artifacts = all_referenced_artifacts[component]
      for regex, subst_fun in re_subst_dict.iteritems():
        text_runs = self._ApplySubstFunctionToRuns(
            text_runs, regex, subst_fun, mr, component_ref_artifacts)

    return text_runs

  def _ApplySubstFunctionToRuns(
      self, text_runs, regex, subst_fun, mr, component_ref_artifacts):
    """Apply autolink regex and substitution function to each text run.

    Args:
      text_runs: list of TextRun objects with parts of the original comment.
      regex: Regular expression for detecting textual references to artifacts.
      subst_fun: function to return autolink markup, or original text.
      mr: common info parsed from the user HTTP request.
      component_ref_artifacts: already-looked-up destination artifacts to use
        when computing substitution text.

    Returns:
      A new list with more and smaller runs, some of which may have tag
      and link attributes set.
    """
    result_runs = []
    for run in text_runs:
      content = run.content
      if run.tag:
        # This chunk has already been substituted, don't allow nested
        # autolinking to mess up our output.
        result_runs.append(run)
      else:
        pos = 0
        for match in regex.finditer(content):
          if match.start() > pos:
            result_runs.append(template_helpers.TextRun(
                content[pos: match.start()]))
          replacement_runs = subst_fun(mr, match, component_ref_artifacts)
          result_runs.extend(replacement_runs)
          pos = match.end()

        if run.content[pos:]:  # Keep any text that came after the last match
          result_runs.append(template_helpers.TextRun(run.content[pos:]))

    # TODO(jrobbins): ideally we would merge consecutive plain text runs
    # so that regexes can match across those run boundaries.

    return result_runs


def RegisterAutolink(services):
  """Register all the autolink hooks."""
  services.autolink.RegisterComponent(
      '01-linkify',
      lambda request, mr: None,
      lambda mr, match: None,
      {_IS_A_LINK_RE: Linkify})

  services.autolink.RegisterComponent(
      '02-tracker',
      CurryGetReferencedIssues(services),
      ExtractProjectAndIssueIds,
      {_ISSUE_REF_RE: ReplaceIssueRef})

  services.autolink.RegisterComponent(
      '03-versioncontrol',
      GetReferencedRevisions,
      ExtractRevNums,
      {_GIT_HASH_RE: ReplaceRevisionRef,
       _SVN_REF_RE: ReplaceRevisionRef})
