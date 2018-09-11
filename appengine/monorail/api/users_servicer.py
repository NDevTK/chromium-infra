# Copyright 2018 The Chromium Authors. All rights reserved.
# Use of this source code is govered by a BSD-style
# license that can be found in the LICENSE file or at
# https://developers.google.com/open-source/licenses/bsd

from api import converters
from api import monorail_servicer
from api import converters
from api.api_proto import users_pb2
from api.api_proto import users_prpc_pb2
from api.api_proto import user_objects_pb2
from businesslogic import work_env
from framework import framework_views


class UsersServicer(monorail_servicer.MonorailServicer):
  """Handle API requests related to User objects.

  Each API request is implemented with a method as defined in the
  .proto file that does any request-specific validation, uses work_env
  to safely operate on business objects, and returns a response proto.
  """

  DESCRIPTION = users_prpc_pb2.UsersServiceDescription

  @monorail_servicer.PRPCMethod
  def GetUser(self, _mc, request):
    # Do any request-specific validation:
    # None.

    # Use work_env to safely operate on business objects.
    email = request.display_name
    user_id = request.user_id

    # Return a response proto.
    return users_pb2.User(
        email=email,
        user_id=user_id)

  @monorail_servicer.PRPCMethod
  def ListReferencedUsers(self, mc, request):
    """Return the list of existing users in a response proto."""
    emails = request.emails
    with work_env.WorkEnv(mc, self.services) as we:
      users = we.ListReferencedUsers(emails)

    with mc.profiler.Phase('converting to response objects'):
      response_users = converters.ConvertUsers(users)
      response = users_pb2.ListReferencedUsersResponse(users=response_users)

    return response

  @monorail_servicer.PRPCMethod
  def GetMemberships(self, mc, request):
    """Return the user groups for the given user visible to the requester."""
    user_id = converters.IngestUserRefs(
        mc.cnxn, [request.user_ref], self.services.user)[0]

    with work_env.WorkEnv(mc, self.services) as we:
      group_ids = we.GetMemberships(user_id)

    with mc.profiler.Phase('converting to response objects'):
      groups_by_id = framework_views.MakeAllUserViews(
          mc.cnxn, self.services.user, group_ids)
      group_refs = converters.ConvertUserRefs(
          group_ids, [], groups_by_id, True)

      return users_pb2.GetMembershipsResponse(group_refs=group_refs)

  @monorail_servicer.PRPCMethod
  def GetUserCommits(self, mc, request):
    """Return a user's commits in a response proto."""
    with work_env.WorkEnv(mc, self.services) as we:
      user_commits = we.GetUserCommits(request.email, request.from_timestamp,
          request.until_timestamp)

    with mc.profiler.Phase('converting to response objects'):
      converted_commits = converters.ConvertCommitList(
          user_commits)
      response = users_pb2.GetUserCommitsResponse(
          user_commits=converted_commits)
      return response

  @monorail_servicer.PRPCMethod
  def GetUserStarCount(self, mc, request):
    """Return the star count for a given user."""
    user_id = converters.IngestUserRefs(
        mc.cnxn, [request.user_ref], self.services.user)[0]

    with work_env.WorkEnv(mc, self.services) as we:
      star_count = we.GetUserStarCount(user_id)

    result = users_pb2.GetUserStarCountResponse(star_count=star_count)
    return result

  @monorail_servicer.PRPCMethod
  def StarUser(self, mc, request):
    """Star a given user."""
    user_id = converters.IngestUserRefs(
        mc.cnxn, [request.user_ref], self.services.user)[0]

    with work_env.WorkEnv(mc, self.services) as we:
      we.StarUser(user_id, request.starred)
      star_count = we.GetUserStarCount(user_id)

    result = users_pb2.StarUserResponse(star_count=star_count)
    return result
