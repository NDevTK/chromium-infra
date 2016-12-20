# Copyright 2016 The Chromium Authors. All rights reserved.
# Use of this source code is govered by a BSD-style
# license that can be found in the LICENSE file or at
# https://developers.google.com/open-source/licenses/bsd

"""Service manager to initialize all services."""

from features import autolink
from services import cachemanager_svc
from services import config_svc
from services import features_svc
from services import issue_svc
from services import project_svc
from services import spam_svc
from services import star_svc
from services import user_svc
from services import usergroup_svc


svcs = None


class Services(object):
  """A simple container for widely-used service objects."""

  def __init__(
      self, project=None, user=None, issue=None, config=None,
      usergroup=None, cache_manager=None, autolink_obj=None,
      user_star=None, project_star=None, issue_star=None, features=None,
      spam=None, hotlist_star=None):
    # Persistence services
    self.project = project
    self.user = user
    self.usergroup = usergroup
    self.issue = issue
    self.config = config
    self.user_star = user_star
    self.project_star = project_star
    self.hotlist_star = hotlist_star
    self.issue_star = issue_star
    self.features = features

    # Misc. services
    self.cache_manager = cache_manager
    self.autolink = autolink_obj
    self.spam = spam


def set_up_services():
  """Set up all services."""

  global svcs
  if svcs is None:
    # Sorted as: cache_manager first, everything which depends on it,
    # issue (which depends on project and config), things with no deps.
    cache_manager = cachemanager_svc.CacheManager()
    config = config_svc.ConfigService(cache_manager)
    features = features_svc.FeaturesService(cache_manager)
    hotlist_star = star_svc.HotlistStarService(cache_manager)
    issue_star = star_svc.IssueStarService(cache_manager)
    project = project_svc.ProjectService(cache_manager)
    project_star = star_svc.ProjectStarService(cache_manager)
    user = user_svc.UserService(cache_manager)
    user_star = star_svc.UserStarService(cache_manager)
    usergroup = usergroup_svc.UserGroupService(cache_manager)
    issue = issue_svc.IssueService(project, config, cache_manager)
    autolink_obj = autolink.Autolink()
    spam = spam_svc.SpamService()
    svcs = Services(
      cache_manager=cache_manager, config=config, features=features,
      issue_star=issue_star, project=project, project_star=project_star,
      user=user, user_star=user_star, usergroup=usergroup, issue=issue,
      autolink_obj=autolink_obj, spam=spam, hotlist_star=hotlist_star)
  return svcs
