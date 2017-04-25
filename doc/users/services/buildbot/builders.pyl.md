# The builders.pyl File Format

[go/builders-pyl]

builders.pyl is a declarative definition of a buildbot master. It is
intended to hide all of the buildbot-specific implementation details
from the user and just expose the features and settings a
non-buildbot-guru cares about.

[TOC]

## What is the .pyl format?

`.pyl` is short for PYthon Literal. It is a subset of Python syntax
intended to capture pure declarations of expressions. It is roughly
analogous to JSON: you can specify any Python object, but should limit
yourself to things like dicts, arrays, strings, numbers, and booleans.
It is basically JSON except that Python-style comments and trailing
commas are allowed.

## Overview

Each builders.pyl describes a single *waterfall*, which is a collection
of buildbot *builders* that talk to a single buildbot *master*. Each
builder may be implemented by multiple *bots*; you can think of a
bot as a single VM or machine.

Each master has one or more builders. A builder is basically a single
configuration running a single series of steps, collected together into
a [recipe](recipes.md). Each builder may have per-builder properties
set for it (to control the logic the recipe executes), and each builder
may also pass along properties from the bot, so there are four types
of configuration:

1.  overall per-master
2.  per-builder
3.  per-scheduler
4.  per-bot

The keys in the dict should follow that order; within each section, all
required keys should appear first (sorted alphabetically), then all
optional keys (sorted alphabetically).

Bots are usually collected into *pools*, so that they can be load
balanced. Every bot in the pool has the same configuration.

*** aside
Side note: buildbot used to call things "slaves" instead of "bots", and
lots of Chromium docs still use "slave".

For compatibility with older versions of builders.pyl, if the file contains
a field called "slave_port", then any field named "bot_X" mustbe called
"slave_X" instead. Once all of the builders.pyl files have been updated,
this support will be dropped.
***

## Example

Here's a simple file containing all of the required fields:

```python
{
  "master_base_class": "Master1",
  "master_port": 20100,
  "master_port_alt": 40100,
  "bot_port": 30100,
  "templates": ["../master.chromium/templates"],

  "builders": {
     "Chromium Mojo Linux": {
       "recipe": "chromium_mojo",
       "scheduler": "chromium_src_commits",
       "bot_pools": ["linux_precise"],
       "category": "0builders",
     },
  },

  "schedulers": {
    "chromium_src_commits": {
      "type": "git_poller",
      "git_repo_url": "https://chromium.googlesource.com/chromium/src.git",
    },
  },

  "bot_pools": {
    "linux_precise": {
      "bot_data": {
        "bits": 64,
        "os": "linux",
        "version": "precise",
      },
      "bots": ["vm{1..50}-m1"],
    },
  },
}
```

## Top-level keys

At the top-level, builders.pyl files contain a single Python dictionary
containing things that are configured per-master.

### master_base_class

This is a *required* field. It must be set to the name of the Python
class of the buildbot master that this master is based on. This is
usually one of the classes defined in
[site_config/config_default.py](https://chromium.googlesource.com/chromium/tools/build/+/master/site_config/config_default.py).

For example, if you were setting up a new master in the -m1 VLAN, you
would be subclassing Master.Master1, so this value would be `'Master1'`.

### master_port

This is a *required* field. It must be set to the main IP port that the
buildbot master instance runs on. You should set this to the port
obtained from the admins.

### master_port_alt

This is a *required* field. It must be set to the alternate IP port that
the buildbot master instance runs on. You should set this to the port
obtained from the admins.

### bot_port

This is a *required* field. It must be set to the port that the buildbot
bots will attempt to connect to on the master.

### templates

This is a *required* field. It must be set to a list of directory paths
(relative to the master directory) that contains the HTML templates that
will be used to display the builds. Each directory is searched in order
for templates as needed (so earlier directories override later
directories).

### buildbot_url

This is an *optional* field. It can be set to customize the URL
the HTML templates use to refer to the top-level web page. If it
is not provided, we will synthesize one for build.chromium.org based
on the master name.

### buildbucket_bucket

This is an *optional* field but must be present if the builders on the
master are intended to be scheduled through buildbucket (i.e., they are
tryservers or triggered from other bots). Such builders normally have
their scheduler set to `None`, so, equivalently, if any of the builders
have their scheduler set to `None`, this field must be present.

If set, it should contain the string value of the
[buildbucket bucket](/appengine/cr-buildbucket/README.md) created for this
buildbot master. If it is not set, it defaults to `None`. By convention,
buckets are named to match the master name, e.g. "master.tryserver.nacl".

### master_classname

This is an *optional* field. If it is not specified, it is synthesized
from the name of the directory containing the builders.pyl file.

For example, if the builders.pyl file was in
[masters/master.client.crashpad](https://chromium.googlesource.com/chromium/tools/build/+/master/masters/master.client.crashpad/builders.pyl),
the master-classname would default to ClientCrashpad.

### service_account_file

This is an *optional* field but must be present if the builders on the
master are intended to be scheduled through buildbucket (i.e., they are
tryservers or triggered from other builders, possibly on other masters).

Such builders normally have their scheduler set to `None`, so,
equivalently, if any of the builders have their scheduler set to `None`,
this field must be present.

If set, it should point to the filename in the credentials directory on
the bot machine (i.e., just the basename + extension, no directory
part), that contains the [OAuth service account
info](../master_auth.md) the bot will use to connect to buildbucket.
By convention, the value is "service-account-\<project\>.json". If not
set, it defaults to None.

### pubsub_service_account_file

Similar to service_account_file, this is also an *optional* field but
must be present if the builders on the master are intended to send build data
to pubsub.

If set, it should point to the filename in the credentials directory on
the bot machine (i.e., just the basename + extension, no directory
part), that contains the [OAuth service account
info](../master_auth.md) the bot will use to connect to pubsub.
By convention, the value is "service-account-\<project\>.json".
The <project> field is usually "luci-milo" for most masters.  If not
set, it defaults to None.

### builders

This is a *required* field and must be a dict of builder names and their
respective configurations; valid values for those configurations are
described in the per-builder configurations section, below.

### schedulers

This is a *required* field and must be a dict of scheduler names and
their respective configurations; valid values for those configurations
are described in the per-scheduler configurations section, below. The
dict may be empty, if there are no scheduled builders, only tryservers,
but it must be present even in that case.

### bot_pools

This is a *required* field and must be a dict of pool names and
properties, as described below.

## Per-builder configurations

Each builder is described by a dict that contains three or four fields:

### recipe

This is a *required* field that specifies the [recipe
name](/doc/users/recipes.md).

### scheduler

This is a *required* field that indicates which scheduler will be used
to schedule builds on the builder.

The field have must be set to either `None` or to one of the keys in the
top-level `schedulers` dict. If it is set to None, then the builder will
only be schedulable via buildbucket; in this situation, the master must
have top-level `buildbucket_bucket` and `service_account_file` values
set (as noted above).

A builder that has a scheduler specified may also potentially be
scheduled via buildbucket, but that doing so would be unusual (builders
should normally only have one purpose).

### bot_pools

This is a *required* field that specifies one or more pools of bots
that can be builders.

### mergeRequests

This is an *optional* field that specifies whether buildbot will merge duplicate
requests together. If unspecified, this field defaults to True if a named
scheduler is specified, and False otherwise.

You might want to merge builds if you have a waterfall builder that is polling
a repository, because you want to always test the most current revision.
You would not want to merge builds for tryservers because you want to test each
revision in isolation.

### auto_reboot

This is an *optional* field that specifies whether the builder should
reboot after each build. If not specified, it defaults to `True`.

### properties

This is an *optional* field that is a dict of settings that will be
passed to the [recipe](recipes.md) as key/value properties.

### botbuilddir

This is an *optional* field; if it is not set, it defaults to the
builder name. This field can be used to share a single build directory
between multiple builders (so, for example, you don't have to check out
the source tree twice for a debug builder and a release builder).

### category

This is an *optional* field that specifies a category for the builder, so you
can group builders visually on the master.  The categories are sorted
left-to-right in ascending order, and for display any initial number is
stripped.  So categories will often be specified like `"0builders"`,
`"1testers"`, etc.

### builder_timeout_s

If set, forcibly kill builds that run longer than this many seconds. If unset
(or None), builds may run indefinitely.

## Per-scheduler configurations

### type

This is a *required* field used to the type of scheduler this is; it
must have one of the following three values: `"cron"`, `"git_poller"`, or
`"repo_poller"`.

`cron` indicates that builds will be scheduled periodically (one or
more times every day). The scheduler dict must also have the "hour" and
"minute" fields.

`git_poller` indicates that builds will be scheduled when there are new
commits to the given repo. The scheduler dict must also have the "git-repo-url"
field.

`repo_poller` behaves the same as `git_poller`, but uses repo rather than git
(repo being the meta repository used in projects such as Android or ChromiumOS).
The scheduler dict must also have the `"repo_url"` field.

### git_repo_url

This is a *required* field if the scheduler type is "git_poller". It must
not be present otherwise.

It must contain a string value that is the URL for a repo to be cloned
and polled for changes.

### branch

This is an *optional* field that is used if the scheduler type is
"git_poller" or "repo_poller". It must not be present otherwise.

It must contain a string value that is the branch name in the repo to watch.
If it is not specified, it defaults to "master".

### repo_url

This is a *required* field if the scheduler type is "repo_poller". It must
not be present otherwise.

The URL that is the base of the repo tree. It is assumed that the manifest is
located in the `manifest` subdirectory of this path. For example, Android would
use `"https://android.googlesource.com/platform"` rather than
`"https://android.googlesource.com/platform/manifest"`.

### rev_link_template

This is an *optional* field that may be used with the "repo_poller" scheduler
type. It is a format string that will be used to generate a link to the change
being built in the build page.

The format string expects two string arguments: the project path and the
revision SHA. For example, Android uses
"https://android.googlesource.com/platform/%s/+/%s".

### hour

This is a *required* if the scheduler type is "cron". It must not be
present otherwise.

This field and the `minute` field control when cron jobs are scheduled
on the builder.

The field may have a value of either `"*"`, an integer, or a list of
integers, where integers must be in the range \[0, 23). The value `"*"`
is equivalent to specifying a list containing every value in the range.
This matches the syntax used for the `Nightly` scheduler in buildbot.

### minute

This is a *required* field if the scheduler type is "cron". It must
not be present otherwise.

This field and the `hour` field control when cron jobs are scheduled on
the builder.

The field may have a value of either `"*"`, an integer, or a list of
integers, where integers must be in the range \[0, 60). The value `"*"`
is equivalent to specifying a list containing every value in the range.
This matches the syntax used for the `Nightly` scheduler in buildbot.

## Per-pool configurations

Each pool (or group) of bots consists of a set of machines that all
have the same characteristics. The pool is described by a dict that
contains two fields

### bot_data

This is a *required* field that contains a dict describing the
configuration of every bot in the pool, as described below.

### bots

This is a *required* field that contains list of either individual hostnames,
one for each VM (do not specify the domain, just the basename), or a
string that can specify a range of hostnames, expanded as the bash shell
would expand them. So, for example, `vm{1..3}-m1` would expand to `vm1-m1`,
`vm2-m1`, `vm3-m1`.

## Per-bot configurations

The bot-data dict provides a bare description of the physical
characteristics of each machine: operating system name, version, and
architecture, with the following keys:

### bits

This is a *required* field and must have either the value 32 or 64 (as
numbers, not strings).

### os

This is a *required* field that must have one of the following values:
`"mac"`, `"linux"`, or `"win"`.

### version

This is a *required* field and must have one of the following values:

os       | Valid values
---------|-------------
`"mac"`  | `"10.6"`, `"10.7"`, `"10.8"`, `"10.9"`, `"10.10"`, `"10.11"`
`"linux"`| `"precise"`, `"trusty"`, `"xenial"`
`"win"`  | `"xp"`, `"vista"`, `"win7"`, `"win8"`, `"win10"`, `"2008"`

## Feedback

[crbug](https://crbug.com) label:
[Infra-MasterGen](https://crbug.com?q=label:Infra-MasterGen)

[go/builders-pyl]: http://go/builders-pyl
