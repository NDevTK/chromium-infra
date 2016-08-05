# Buildbucket-Swarming integration, aka Swarmbucket

[go/swarmbucket]

Buildbucket has native integration with swarming and recipes.
A bucket can be configured so a build for one of configured builders is
converted to a swarming task that runs a recipe. The results are reported back
to buildbucket when the task compeltes.

[TOC]

## Configuration

There are 3 swarmbucket configuration levels.
From global to concrete: application, bucket and build.
A more concrete config value can override a global one.

### Application

Luci-config file
[services/\<buildbucket-app-id\>:swarming_task_template.json][swarming_task_template.json]
specifies the body of POST request that buildbucket sends to swarming when
creating a task.

In particular it contains the command line parameters.
In general, you can specify anything you want, but to run a recipe,
see [kitchen].

#### Parameters

The template may contain parameters of form `$name` in string literals.
They will be expanded during task scheduling. Known parameters:

* bucket: value of `build.bucket`.
* builder: value of `builder_name` build parameter.
* repository: repository URL of the recipe.
* revision: revision of the recipe.
* recipe: name of the recipe.
* properties-json: a JSON string containing build properties.
* project: LUCI project id that the bucket is defined at, e.g. "chromium".

Example: [swarming_task_template.json].


#### Canary

If luci-config file
`services/\<buildbucket-app-id\>:swarming_task_template_canary.json` exists
then it is used as a task template in a fraction of builds specified by bucket
task_template_canary_percentage config value. The decision to use or not to use
the canary template for a build is deterministic.

See also "canary_template" build parameter.

### Bucket level

A bucket entry in buildbucket config may have `swarming` key, for example:

    buckets {
      name: "foobar"

      swarming {
        hostname: "chromium-swarm.appspot.com"
        builders {
          name: "Linux-Release"
          dimensions {
            key: "os"
            value: "Linux"
          }
          recipe {
            repository: "https://chromium.googlesource.com/chromium/tools/build"
            name: "chromium"
          }
        }
        builders {
          name: "Windows-Release"
          dimensions {
            key: "os"
            value: "Windows"
          }
          recipe {
            repository: "https://chromium.googlesource.com/chromium/tools/build"
            name: "chromium"
          }
        }
      }
    }

For the format and documentation see `Swarming` message in the
[project_config.proto](../proto/project_config.proto).
Real world configuration example:
["master.tryserver.infra" bucket](https://chromium.googlesource.com/infra/infra/+/infra/config/cr-buildbucket.cfg)

### Build level

A buildbucket build can have `"swarming"` parameter, which is a JSON object with
optional properties:

* `"recipe"`: specifies a recipe to run, an object with optional properties:
  * `"revision"`: recipe revision
* `"canary_template"`: specifies whether canary task template must be used:
  * `true`: use the canary template. If not found, respond with an error.
  * `false`: do not use canary template.
  * `null` (default): use canary template with some low probability if it
    exists.
* `"override_builder_cfg"`: can override builder configuration defined on the
  server. See also a section about it below.

#### Override configuration dynamically

`swarming.override_builder_cfg` parameter can override builder configuration
defined on the server. For example, value


```javascript
{
  "dimensions": ["cores:64"]
}
```

(re)defines "cores" dimension to be "64" for this particular build.

The format is defined by the Builder message in
[project_config.proto](../proto/project_config.proto); in practice, it is JSONPB
of the message.

## Tags

A swarming task created by buildbucket has extra tags:

* `buildbucket_hostname:<hostname>`
* `buildbucket_bucket:<bucket>`
* `recipe_repository:<repo url>`
* `recipe_revision:<revision>`
* `recipe_name:<name>`
* all tags in `swarming_tags` of builder config.
* all tags in build creation request.

A buildbucket build associated with a swarming task has extra tags:

* `swarming_hostname:<hostname>`
* `swarming_task_id:<task_id>`
* `swarming_tag:<tag>` for each swarming task tag, even if it was derived from
  a buildbucket build tag.
* `swarming_dimension:<dimension>` for each dimension.

## Life of a swarmbucket build

1. A user schedules a build on bucket `"foobar"` configured as above.
   The build has `"builder_name": "Linux-Release"` parameter.
1. `Linux-Release` config is matched. This build is for swarming.
1. Task template is rendered. Repository URL, recipe and other parameters
   are expanded to "https://chromium.googlesource.com/chromium/tools/build",
   "chromium", etc.
1. Swarming task is created.
1. Build entity is mutated to have extra tags, to point to the task,
   is marked as STARTED and saved.
1. Swarming task starts. Kitchen starts, runs the recipe.
1. Swarming task completes. Swarming sends a message to buildbucket via PubSub.
1. Buildbucket receives the message and marks the build as completed.

Note: swarming does not notify on task start, so buildbucket marks builds as
STARTED right after creation.

[go/swarmbucket]: https://goto.google.com/swarmbucket
[kitchen]: https://chromium.googlesource.com/infra/infra/+/master/go/src/infra/tools/kitchen/
[isolate_kitchen.py]: https://github.com/luci/recipes-py/blob/master/go/cmd/kitchen/isolate_kitchen.py
[swarming_task_template.json]: https://chrome-internal.googlesource.com/infradata/config/+/master/configs/cr-buildbucket/swarming_task_template.json
