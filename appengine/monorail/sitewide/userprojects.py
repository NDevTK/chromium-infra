# Copyright 2016 The Chromium Authors. All rights reserved.
# Use of this source code is govered by a BSD-style
# license that can be found in the LICENSE file or at
# https://developers.google.com/open-source/licenses/bsd

"""Classes for the user projects feed."""

from framework import jsonfeed
from sitewide import sitewide_helpers


class ProjectsJsonFeed(jsonfeed.JsonFeed):
  """Servlet to get all of a user's projects in JSON format."""

  def HandleRequest(self, mr):
    """Retrieve list of a user's projects for the "My projects" menu.

    Args:
      mr: common information parsed from the HTTP request.

    Returns:
      Results dictionary in JSON format
    """
    if not mr.auth.user_id:
      return {'error': 'User is not logged in.'}

    json_data = {}

    with self.profiler.Phase('page processing'):
      json_data.update(self._GatherProjects(mr))

    return json_data

  def _GatherProjects(self, mr):
    """Return a dict of project names the current user is involved in."""
    with self.profiler.Phase('GetUserProjects'):
      project_lists = sitewide_helpers.GetUserProjects(
          mr.cnxn, self.services, mr.auth.user_pb, mr.auth.effective_ids,
          mr.auth.effective_ids)
      (visible_ownership, _visible_deleted, visible_membership,
       visible_contrib) = project_lists

    with self.profiler.Phase('GetStarredProjects'):
      starred_projects = sitewide_helpers.GetViewableStarredProjects(
          mr.cnxn, self.services, mr.auth.user_id,
          mr.auth.effective_ids, mr.auth.user_pb)

    projects_dict = {
        'memberof': [p.project_name for p in visible_membership],
        'ownerof': [p.project_name for p in visible_ownership],
        'contributorto': [p.project_name for p in visible_contrib],
        'starred_projects': [p.project_name for p in starred_projects],
    }

    return projects_dict
