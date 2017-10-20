Blox Deployments User Experience
================================

UI Conventions
--------------
- **Showing `Inactive` `EnvironmentRevision`s:** For all the output examples below, we're including `EnvironmentRevision`s in the `Inactive` state by specifying the `--all` flag to `describe-environment-status`. Without this flag, only `EnvironmentRevision`s that are not in `Inactive` will be shown.

Concepts
--------
Blox introduces these new concepts to the ECS APIs:

- `Environment`: A logical unit of deployment of a container-based application. An `Environment` declares how a given `Task Definition` should run on a given Instance Group, and how changes to these properties should be applied to a cluster.
- `EnvironmentRevision`: A versioned snapshot at a given point in time of an `Environment`. Every time a change is made to an `Environment` using `updateEnvironment`, a new `EnvironmentRevision` is created.
- `Deployment`: The process which actually makes changes to the state of a Cluster (i.e. run/terminate `Task`s) in order to make an `EnvironmentRevision` `Active`.
- Starting a `Deployment`: The event which changes which `EnvironmentRevision` is `Active` for a specific `Environment` in a controlled way.
- `Rollback`: A `Deployment` that changes the `Active` `EnvironmentRevision` for an `Environment` to an older `EnvironmentRevision`.


Use Cases
---------

The following are typical use cases for Deployments:

- Create a new `Environment` for an application, and `Deploy` its first `EnvironmentRevision` in order to launch `Tasks` in a cluster.
- Update an existing `Environment` to create a new `EnvironmentRevision`, and `Deploy` the newly created `EnvironmentRevision` in order to transition `Task`s in the cluster from the old `EnvironmentRevision` to the new one in a controlled way.
- Inspect the status of an `Environment`
- Rollback an `Environment` to an earlier `EnvironmentRevision` if its active `EnvironmentRevision` is causing problems.
- Pause an `Environment` in an emergency situation in order to prevent it from running/terminating any tasks in a cluster.
- Delete an `Environment` completely to stop managing its tasks in a cluster.
- (not applicable for Daemon) Scale up the `Environment` to facilitate more load.

### Creating and deploying an `Environment`

> `Environment`: A logical unit of deployment of a container-based application. An `Environment` declares how a given `Task Definition` should run on a given Instance Group, and how changes to these properties should be applied to a cluster.

Use `create-environment` to create a new `Environment` for an application in a given cluster:

```
$ aws ecs create-environment --yaml-file <<EOF
Name: "SomeEnvironment/Prod"
TaskDefinition: "log-daemon:1"
Cluster: "default"
InstanceGroup:
  Query: "attribute:stack == prod"
Deployment:
  Type: Daemon
  Method: "ReplaceThenTerminate" # (or "TerminateThenReplace")
  Configuration:
    MinHealthyPercent: 80
    MaxOverPercent: 25
EOF
{ "EnvironmentRevision": 1 }
```

At this point, nothing has actually happened to the cluster. We can verify this by inspecting the `Environment`'s status:

```
$ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
ENVIRONMENT               ACTIVE REVISION  STATUS     OUTDATED  DESIRED   HEALTHY  TOTAL
SomeEnvironment/Prod      N/A              Created    0         0         0        5

REVISION                  STATUS     DESIRED   HEALTHY  TOTAL
SomeEnvironment/Prod:1    Created    0         0        5
```

The meanings of each section are:
- `ENVIRONMENT`: The `Environment` we're inspecting
    - `ACTIVE REVISION`: The revision of the `EnvironmentRevision` that the `Environment` will try to maintain.
    - `STATUS`: The status of the currently `Active` `EnvironmentRevision`
    - `OUTDATED`: The number of tasks that do not belong to the `Active` `EnvironmentRevision`
    - `DESIRED`: The number of tasks that the current stage of the deployment will try to have running
    - `HEALTHY`: The number of healthy tasks across all `EnvironmentRevision`s.
    - `TOTAL`: The target number of tasks that should belong to this `EnvironmentRevision` when it's `Active`.
- `REVISION`: The `EnvironmentRevision`s of the `Environment` we're inspecting. If `--all` is specified it will list all revisions for the `Environment`, otherwise only revisions that are not `Inactive` will be shown.
    - `STATUS`: The current status of the `EnvironmentRevision`; one of `Active`, `Inactive`, `Deploying`, `Undeploying`, `Reverting` or `Reverted`.
    - `DESIRED`: The number of tasks that the scheduler will attempt to maintain for this revision, at this point in time.
    - `HEALTHY`: The number of tasks in this revision that are healthy.
    - `TOTAL`: The target number of tasks that should belong to this `EnvironmentRevision` in order for the `Environment` to reach a steady state.

In the above example, we can see that the `Environment` should be running 5 `Task`s in total (one for each `ContainerInstance` in the cluster), but it's not running any tasks at the moment (`HEALTHY` is 0), and it will not try to run any tasks right now (`DESIRED` is 0).

In order to actually run `Task`s to make this `Environment` `Active`, we first have to `Deploy` the newly created `EnvironmentRevision`:

```
$ aws ecs deploy-environment --environment SomeEnvironment/Prod --revision 1
{ "DeploymentId": "SomeEnvironment/Prod:1@2017090611472301" }
```

This will mark the `SomeEnvironment/Prod:1` `EnvironmentRevision` as the target revision, and start running `Task`s as quickly as possible:
```
$ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
ENVIRONMENT               TARGET REVISION  STATUS     OUTDATED  DESIRED   HEALTHY  TOTAL
SomeEnvironment/Prod      1                Deploying  0         5         0        5

REVISION                  STATUS     DESIRED   HEALTHY  TOTAL
SomeEnvironment/Prod:1    Deploying  5         0        5
```

After a while, the `Task`s run by the `Deployment` will become healthy, and the environment itself will become `Active` to indicate that it's reached a steady state.

```
$ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
ENVIRONMENT               TARGET REVISION  STATUS     OUTDATED  DESIRED   HEALTHY  TOTAL
SomeEnvironment/Prod      1                Active     0         5         5        5

REVISION                  STATUS     DESIRED   HEALTHY  TOTAL
SomeEnvironment/Prod:1    Active     5         5        5
```

### Updating an existing `Environment`

Typically, we'd update an `Environment` in order to:

-   update the task definition that should be running (i.e. deploy new code)
-   change the set of `ContainerInstance`s on which to run tasks
-   change the deployment configuration (i.e. to change health thresholds) (TODO: Maybe require new env for this?)

#### Updating the task definition

If we want to deploy new code to our `Environment`, we can update it with the new `TaskDefinition` we want to run:

```
$ aws ecs update-environment --environment SomeEnvironment/Prod --base-revision 1 --task-definition "log-daemon:2"
{ "EnvironmentRevision": 2 }
```

The `--base-revision` flag needs to be provided here because we are not fully specifying all the parameters of the `Environment`. This ensures that we're basing our changes on the `EnvironentRevision` that we expect, and that it hasn't been modified by someone else.

(TODO: We may want to hide this in the CLI/Console, and just assume the changes apply to the latest state.)

As before, just updating the `Environment` has not made any changes to the `Task`s running in the cluster:

```
$ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
ENVIRONMENT               TARGET REVISION  STATUS     OUTDATED  DESIRED   HEALTHY  TOTAL
SomeEnvironment/Prod      1                Active     0         5         5        5

REVISION                  STATUS     DESIRED   HEALTHY  TOTAL
SomeEnvironment/Prod:1    Active     5         5        5
SomeEnvironment/Prod:2    Created    0         0        0
```

In order to actually start replacing the tasks from `EnvironmentRevision` 1 with those from `EnvironmentRevision` 2, we have to `Deploy` the revision to the environment:

```
$ aws ecs deploy-environment --environment SomeEnvironment/Prod --revision 2
{ "DeploymentId": "SomeEnvironment/Prod:2@2017090611472301" }
```

Now the `Deployment` will gradually start killing `Task`s with the old `TaskDefinition` and launching tasks for the new `TaskDefinition`. Exactly how this gradual rollout happens depends on the `Environment`'s `Deployment` configuration; see the [section below]() for details.

Part-way through the `Deployment`, the `Environment`'s status could look something like this:

```
$ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
ENVIRONMENT               TARGET REVISION  STATUS     OUTDATED  DESIRED   HEALTHY  TOTAL
SomeEnvironment/Prod      2                Deploying  2         7         5        5

REVISION                  STATUS       DESIRED   HEALTHY  TOTAL
SomeEnvironment/Prod:1    Undeploying  2         2        0
SomeEnvironment/Prod:2    Deploying    5         3        5
```

Once the number of `HEALTHY` tasks matches the `TOTAL` tasks for the new `EnvironmentRevision`, and the number of `HEALTHY` `Task`s for the old `EnvironmentRevision` is 0, the `Deployment` is complete, and the `Environment` becomes `Active` again:

```
$ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
ENVIRONMENT               TARGET REVISION  STATUS     OUTDATED  DESIRED   HEALTHY  TOTAL
SomeEnvironment/Prod      2                Active     0         5         5        5

REVISION                  STATUS       DESIRED   HEALTHY  TOTAL
SomeEnvironment/Prod:1    Inactive     0         0        0
SomeEnvironment/Prod:2    Active       5         5        5
```

#### Updating the `InstanceGroup`

> TODO: Thinking about this more, I'm not sure if versioning the `InstanceGroup` along with the `EnvironmentRevision` is the right thing here. In particular, I'm not sure if it might be surprising that a rollback deployment will also change which instances form part of an environment. One way we could circumvent the weirdness is by going back to creating new `EnvironmentRevision`s for rollbacks, and only rolling back the `TaskDefinition` part (see Appendix of Alternatives Considered).

If we want to change the target group of `ContainerInstances` that the `Environment` should run on, we can modify the `InstanceGroup` of the `Environment`. For example, to modify the `Environment` from the previous example to also run on `ContainerInstance`s with the attribute `stack=gamma`, run:

```
$ aws ecs update-environment --environment SomeEnvironment/Prod --base-revision 2 --instance-group-query "attribute:stack in (prod, gamma)"
{ "EnvironmentRevision": 3 }
```

As before, just updating the `Environment` has not made any changes to the `Task`s running in the cluster:

```
$ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
ENVIRONMENT               TARGET REVISION  STATUS     OUTDATED  DESIRED   HEALTHY  TOTAL
SomeEnvironment/Prod      2                Active     0         5         5        5

REVISION                  STATUS     DESIRED   HEALTHY  TOTAL
SomeEnvironment/Prod:1    Inactive   0         0        0
SomeEnvironment/Prod:2    Active     5         5        5
SomeEnvironment/Prod:3    Created    0         0        0
```

In order to actually run `Tasks` on the `ContainerInstance`s that now match the `InstanceGroup`, we deploy the new `EnvironmentRevision`:

```
$ aws ecs deploy-environment --environment SomeEnvironment/Prod --revision 3
{ "DeploymentId": "SomeEnvironment/Prod:3@2017092211473204" }

$ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
ENVIRONMENT               TARGET REVISION  STATUS     OUTDATED  DESIRED   HEALTHY  TOTAL
SomeEnvironment/Prod      3                Deploying  2         7         5        8

REVISION                  STATUS      DESIRED   HEALTHY  TOTAL
SomeEnvironment/Prod:1    Inactive    0         0        0
SomeEnvironment/Prod:2    Undeploying 2         2        0
SomeEnvironment/Prod:3    Deploying   5         3        8
```

In the example above, 3 tasks from `SomeEnvironment/Prod:2` are also matched by the new `InstanceGroup`, while 2 of the tasks no longer match the new `InstanceGroup`.

Since the `TaskDefinition` for `SomeEnvironment/Prod` did not change from the previous revision, existing tasks from `SomeEnvironment/Prod:2` are immediately considered `HEALTHY` in `SomeEnvironment/Prod:3` (which removes them from `SomeEnvironment/Prod:2` as if they were terminated).

Once the `Task`s on the new `ContainerInstance`s have become healthy, and the tasks on the old `ContainerInstance`s have terminated, `SomeEnvironment/Prod:3` becomes `Active`:

```
$ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
ENVIRONMENT               TARGET REVISION  STATUS     OUTDATED  DESIRED   HEALTHY  TOTAL
SomeEnvironment/Prod      3                Active     0         8         8        8

REVISION                  STATUS       DESIRED   HEALTHY  TOTAL
SomeEnvironment/Prod:1    Inactive     0         0        0
SomeEnvironment/Prod:2    Inactive     0         0        0
SomeEnvironment/Prod:3    Active       8         8        8
```

### Rollback to an earlier `EnvironmentRevision`

Usually, trying to deploy an `EnvironmentRevision` older than the latest `Active` `EnvironmentRevision` will result in an error. This is to prevent inadvertently rolling back changes that have been deployed by another user.

In the case where we *do* want to roll back changes, we have to pass the `--rollback` flag to `deploy-environment`.

Let's say that after several revisions, `SomeEnvironment/Prod:5` introduced a change that is actually broken:

```
$ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
ENVIRONMENT               TARGET REVISION  STATUS     OUTDATED  DESIRED   HEALTHY  TOTAL
SomeEnvironment/Prod      3                Active     0         8         8        8

REVISION                  STATUS       DESIRED   HEALTHY  TOTAL
SomeEnvironment/Prod:1    Inactive     0         0        0
SomeEnvironment/Prod:2    Inactive     0         0        0
SomeEnvironment/Prod:3    Inactive     0         0        0
SomeEnvironment/Prod:4    Inactive     0         0        0
SomeEnvironment/Prod:5    Active       8         8        8
```

We can promptly roll back to any earlier `EnvironmentRevision` using `deploy-environment` with `--rollback` to deploy it:

```
$ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
ENVIRONMENT               TARGET REVISION  STATUS     OUTDATED  DESIRED   HEALTHY  TOTAL
SomeEnvironment/Prod      3                Active     8         8         8        8

REVISION                  STATUS       DESIRED   HEALTHY  TOTAL
SomeEnvironment/Prod:1    Inactive     0         0        0
SomeEnvironment/Prod:2    Inactive     0         0        0
SomeEnvironment/Prod:3    Deploying    3         0        8
SomeEnvironment/Prod:4    Inactive     0         0        0
SomeEnvironment/Prod:5    Reverting    5         8        0
```

Note that even though the above example shows rollback of an `EnvironmentRevision` that's already `Active`, even a `Deploying` `EnvironmentRevision` could still be rolled back. You could even specify a rollback `EnvironmentRevision` that's older than the one that is currently `Undeploying`; in this case, both the `Deploying` and `Undeploying` `EnvironmentRevision`s will be reverted.

#### Rollback deployment configuration

Rollback deployments may have different deployment configurations to allow for quicker recovery, or to avoid getting stuck due to breached deployment constraints such as `MinHealhtyPercent`. These can be specified in a `RollbackConfiguration` section of the environment when it is created.

### Pause ongoing deployments

In some scenarios, it may be dangerous to continue running or terminating `Task`s in a cluster. For example, there might be some problem with the underlying cluster that causes new `Task`s to fail, or some external dependency being unavailable might make it impossible for `Task`s to start/stop cleanly.

For these cases, all actions taken by an `Environment` can be suspended using the `freeze-environment`* command:

```
$ aws ecs freeze-environment --environment-name SomeEnvironment/Prod

$ aws ecs describe-environment-status --short --environment-name SomeEnvironment/Prod --table
ENVIRONMENT               TARGET REVISION  STATUS     OUTDATED  DESIRED   HEALTHY  TOTAL
SomeEnvironment/Prod      3                Frozen     2         8         3        8
```

(`--short` omits details about each revision, and just shows aggregate environment status).

Once it is safe for the `Environment` to take action again, it can be resumed using the `thaw-environment`* command:

```
$ aws ecs thaw-environment --environment-name SomeEnvironment/Prod

$ aws ecs describe-environment-status --short --environment-name SomeEnvironment/Prod --table
ENVIRONMENT               TARGET REVISION  STATUS     OUTDATED  DESIRED   HEALTHY  TOTAL
SomeEnvironment/Prod      3                Frozen     0         8         8        8
```

Note that suspending an `Environment` doesn't prevent you from starting a new deployment; however, in order for any deployments to progress, the `Environment` must be thawed first.

<small>* `freeze` and `thaw` here are placeholders; `suspend` and `resume` already have meanings in the container ecosystem, and using those terms here would be confusing.</small>

### Failure modes
#### `Stuck`
While the scheduler is running or terminating `Task`s to try and bring an `EnvironmentRevision`'s `HEALTHY` count up to its `DESIRED` count, external factors might prevent these actions from completing successfully. For example, if the deployment needs a `Task` to be running on a particular `ContainerInstance` and that instance has no more resources available, these counts might never converge.

Once an `EnvironmentRevision` has had `HEALTHY != DESIRED` for a configurable amount of time, it will enter the `Stuck`* state to indicate that it is unable to progress. While in this state, the scheduler will still try to make progress by running/terminating other `Task`s to meet the `EnvironmentRevision`s `DESIRED` state:

```
$ aws ecs describe-environment-status --short --environment-name SomeEnvironment/Prod --table
ENVIRONMENT               TARGET REVISION  STATUS     OUTDATED  DESIRED   HEALTHY  TOTAL
SomeEnvironment/Prod      3                Stuck      2         8         3        8
```

(Potentially out of scope): When an EnvironmentRevision enters the `Stuck` state, a Cloudwatch Event will be emitted.

If any of the `EnvironmentRevision`s of an environment's `STATUS` is `Stuck`, the overall `Environment` will also show as `Stuck`.

#### `Failing`
If an environment defines the `MinHealthyPercent` deployment configuration option, then the environment as a whole can enter the `Failing` state. In particular, if `HEALTHY / TOTAL < MinHealthyPercent / 100` for a configurable amount of time, then the `Environment` will enter the `Failing` state, and emit a Cloudwatch Event.

`Environment` Deployment Configuration
--------------------------------------

### `MinHealthyPercent`
For all supported deployment methods, we can define the `MinHealthyPercent` configuration option. If set,
 - the scheduler will never take actions that would reduce the actual percentage of healthy tasks below this threshold
 - if the threshold is breached for some other reason, the scheduler will refuse to take any actions that would further reduce the number of healthy tasks (but it will continue to take actions that will increase the number of healthy tasks

### `DeploymentType`s
Blox supports two types of deployment:
- `Daemon`: Ensures that a single `Task` of a given `TaskDefinition` is running on every `ContainerInstance` in an `InstanceGroup`.
- `Service` (not in V1): Ensures that `N` `Tasks` of a given `TaskDefinition` is running within the `InstanceGroup`, subject to configurable placement constraints.

#### `DeploymentType: Daemon`

Daemon `Environment`s behave slightly differently from Service `Environment`s. Daemon `Task`s typically provide `ContainerInstance`-wide capabilities such as metrics/logging or network overlay, that the other `Task`s on the `ContainerInstance` all depend upon. Any amount of time that the `Task` is not present on the `ContainerInstance`, could potentially result in downtime for the other `Task`s on the `ContainerInstance`.

Since `Task`s sometimes fail, and cluster `ContainerInstance` membership changes over time, it is expected that there will often be some window of time where the requirements above cannot be met. In particular, for larger clusters, or clusters that are backed by Autoscaling Groups, the available `ContainerInstance`s in the cluster might change frequently (and sometimes dramatically).

The goal then is to bound the size of the window where a `ContainerInstance` is not running the Daemon `Environment`, and to provide alarms on when this upper bound is breached. If a configurable number of `ContainerInstance`s are without the correct version of a Daemon `Task` for a configurable amount of time, the `EnvironmentRevision` will enter the `Failing` state and emit an alarm. However, it will continue to try and make progress.

Daemon deployments have two update strategies, with different availability impacts:

##### `DeploymentMethod: TerminateThenReplace`

If you're using this deployment strategy, then Blox will kill `(100 - MinHealthyPercent)%` of `Task`s, and only launch new `Task`s to replace them as they are successfully terminated. This will ensure that there's never more than one version of the `Task` running at the same time on the same `ContainerInstance`, but may result in some downtime for the daemon. This is a good fit for daemons that would fail if there is more than one running per `ContainerInstance`, and that can tolerate short periods of unavailability.

##### `DeploymentMethod: ReplaceThenTerminate`

If you're using this deployment strategy, then Blox will launch `(100 + MaxOverPercent)% * TOTAL` `Task`s on as many `ContainerInstance`s, and only kill old `Task`s once the new `Task` is healthy. This will ensure that there's never less than one `Task` running at the same time on the same `ContainerInstance`, but will result in two copies of the `Task` running on at most `MaxOver` `ContainerInstance`s at any given time. This is a good fit for daemons that must always have at least one copy of the `Task` running, and that don't care if another copy of the `Task` is running on the same `ContainerInstance`.

### Deleting an Environment
Once an `Environment` is no longer needed, it can be deleted using the `delete-environment` command:

```
$ aws ecs delete-environment --name SomeEnvironment/Prod
```

Deleting an `Environment` will not make any changes to cluster state by default, it will only stop taking deployment actions. This means that any `Task`s that were running at the time of deletion will still be running.

If we want to also terminate any running tasks in the environment, we can specify the `--terminate-tasks` option:

```
$ aws ecs delete-environment --name SomeEnvironment/Prod --terminate-tasks
```

This will place the environment into the `Deleting` state, as all revisions that still have tasks running are undeployed:

```
$ aws ecs describe-environment-status --short --environment-name SomeEnvironment/Prod --table
ENVIRONMENT               TARGET REVISION  STATUS     OUTDATED  DESIRED   HEALTHY  TOTAL
SomeEnvironment/Prod      3                Deleting   5         0         5        0
```

Either way, once it is successfully deleted, the environment will still be visible in the `Deleted` state for 24-48 hours after deletion:

```
$ aws ecs describe-environment-status --short --environment-name SomeEnvironment/Prod --table
ENVIRONMENT               TARGET REVISION  STATUS     OUTDATED  DESIRED   HEALTHY  TOTAL
SomeEnvironment/Prod      3                Deleted    0         0         0        0
```

## Appendix A: Detailed deployment examples

> For all the cluster state visualizations below, the legend is:
>
> | :large_blue_circle: pending task | :white_check_mark: running task | :red_circle: terminating task | :no_entry: terminated task |
> |:---------------------------------|:--------------------------------|:------------------------------|:---------------------------|

Here are some examples that illustrate what exactly will occur in a cluster when a new `EnvironmentRevision` is deployed under various combinations of deployment method and cluster state changes.

> TODO: This entire section needs rework to take into account some changes to `MinHealthyPercent` handling.

### Example 1: TerminateThenReplace; Steady State

Let's consider a Daemon `Environment` update with `Method = TerminateThenReplace, MinHealthyPercent = 60`, which completes without:
- any changes in cluster membership (i.e. no `ContainerInstance`s join or leave the cluster)
- any `Task`s failing to start up
- any existing `Task`s terminating abnormally

This is what the overall deployment progress would look like while inspecting the cluster:

| T  | Event                  | i<sub>1</sub>                                            | i<sub>2</sub>                                            | i<sub>3</sub>                                            | i<sub>4</sub>                                            | i<sub>5</sub>                                             |
|:---|:-----------------------|:---------------------------------------------------------|:---------------------------------------------------------|:---------------------------------------------------------|:---------------------------------------------------------|:----------------------------------------------------------|
| 0  | Deployment Starts      | :white_check_mark:<sub>v1</sub>                          | :white_check_mark:<sub>v1</sub>                          | :white_check_mark:<sub>v1</sub>                          | :white_check_mark:<sub>v1</sub>                          | :white_check_mark:<sub>v1</sub>                           |
| 1  | Terminate 1, 2         | :red_circle:<sub>v1</sub>                                | :red_circle:<sub>v1</sub>                                | :white_check_mark:<sub>v1</sub>                          | :white_check_mark:<sub>v1</sub>                          | :white_check_mark:<sub>v1</sub>                           |
| 2  | Terminated 1; start 1  | :no_entry:<sub>v1</sub> :large_blue_circle:<sub>v2</sub> | :red_circle:<sub>v1</sub>                                | :white_check_mark:<sub>v1</sub>                          | :white_check_mark:<sub>v1</sub>                          | :white_check_mark:<sub>v1</sub>                           |
| 3  | Terminated 2; start 2  | :no_entry:<sub>v1</sub> :large_blue_circle:<sub>v2</sub> | :no_entry:<sub>v1</sub> :large_blue_circle:<sub>v2</sub> | :white_check_mark:<sub>v1</sub>                          | :white_check_mark:<sub>v1</sub>                          | :white_check_mark:<sub>v1</sub>                           |
| 4  | Running 1; Terminate 3 | :no_entry:<sub>v1</sub> :white_check_mark:<sub>v2</sub>  | :no_entry:<sub>v1</sub> :large_blue_circle:<sub>v2</sub> | :red_circle:<sub>v1</sub>                                | :white_check_mark:<sub>v1</sub>                          | :white_check_mark:<sub>v1</sub>                           |
| 5  | Running 2; Terminate 4 | :no_entry:<sub>v1</sub> :white_check_mark:<sub>v2</sub>  | :no_entry:<sub>v1</sub> :white_check_mark:<sub>v2</sub>  | :red_circle:<sub>v1</sub>                                | :red_circle:<sub>v1</sub>                                | :white_check_mark:<sub>v1</sub>                           |
| 6  | Terminated 3; start 3  | :no_entry:<sub>v1</sub> :white_check_mark:<sub>v2</sub>  | :no_entry:<sub>v1</sub> :white_check_mark:<sub>v2</sub>  | :no_entry:<sub>v1</sub> :large_blue_circle:<sub>v2</sub> | :red_circle:<sub>v1</sub>                                | :white_check_mark:<sub>v1</sub>                           |
| 7  | Terminated 4; start 4  | :no_entry:<sub>v1</sub> :white_check_mark:<sub>v2</sub>  | :no_entry:<sub>v1</sub> :white_check_mark:<sub>v2</sub>  | :no_entry:<sub>v1</sub> :large_blue_circle:<sub>v2</sub> | :no_entry:<sub>v1</sub> :large_blue_circle:<sub>v2</sub> | :white_check_mark:<sub>v1</sub>                           |
| 8  | Running 3; Terminate 5 | :no_entry:<sub>v1</sub> :white_check_mark:<sub>v2</sub>  | :no_entry:<sub>v1</sub> :white_check_mark:<sub>v2</sub>  | :no_entry:<sub>v1</sub> :white_check_mark:<sub>v2</sub>  | :no_entry:<sub>v1</sub> :large_blue_circle:<sub>v2</sub> | :red_circle:<sub>v1</sub>                                 |
| 9  | Running 4              | :no_entry:<sub>v1</sub> :white_check_mark:<sub>v2</sub>  | :no_entry:<sub>v1</sub> :white_check_mark:<sub>v2</sub>  | :no_entry:<sub>v1</sub> :white_check_mark:<sub>v2</sub>  | :no_entry:<sub>v1</sub> :white_check_mark:<sub>v2</sub>  | :red_circle:<sub>v1</sub>                                 |
| 10 | Terminated 5; start 5  | :no_entry:<sub>v1</sub> :white_check_mark:<sub>v2</sub>  | :no_entry:<sub>v1</sub> :white_check_mark:<sub>v2</sub>  | :no_entry:<sub>v1</sub> :white_check_mark:<sub>v2</sub>  | :no_entry:<sub>v1</sub> :white_check_mark:<sub>v2</sub>  | :no_entry:<sub>v1</sub>  :large_blue_circle:<sub>v2</sub> |
| 11 | Running 5              | :no_entry:<sub>v1</sub> :white_check_mark:<sub>v2</sub>  | :no_entry:<sub>v1</sub> :white_check_mark:<sub>v2</sub>  | :no_entry:<sub>v1</sub> :white_check_mark:<sub>v2</sub>  | :no_entry:<sub>v1</sub> :white_check_mark:<sub>v2</sub>  | :no_entry:<sub>v1</sub>  :white_check_mark:<sub>v2</sub>  |

At *T=0*, the only visible effect of the deployment is updated Desired counts in the status display:

```
$ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
REVISION                  STATUS       DESIRED   HEALTHY
SomeEnvironment/Prod:1    Undeploying  3         5
SomeEnvironment/Prod:2    Deploying    0         0
```

At *T=1*, two `Task`s in the old `EnvironmentRevision` get terminated, and replaced with two pending `Task`s:

```
$ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
REVISION                  STATUS       DESIRED   HEALTHY
SomeEnvironment/Prod:1    Undeploying  3         3
SomeEnvironment/Prod:2    Deploying    2         0
```

At *T=4*, one of the `Task`s finish terminating, and it's replaced by a new `Task`, so we terminate another old `Task`:

```
$ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
REVISION                  STATUS       DESIRED   HEALTHY
SomeEnvironment/Prod:1    Undeploying  2         3
SomeEnvironment/Prod:2    Deploying    3         1
```

`Deployment` progresses in a similar fashion until none of the old `EnvironmentRevision` is left at *T=10*:

```
$ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
REVISION                  STATUS       DESIRED   HEALTHY
SomeEnvironment/Prod:1    Undeploying  0         0
SomeEnvironment/Prod:2    Deploying    5         4
```

The final `Task` is launched at *T=11* and the new `EnvironmentRevision` becomes active:

```
$ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
REVISION                  STATUS       DESIRED   HEALTHY
SomeEnvironment/Prod:1    Inactive     0         0
SomeEnvironment/Prod:2    Active       5         5
```

##### Example 2: Scale up

If new `ContainerInstance`s join the cluster while an `Environment` is deploying, running the daemon `Task` on those `ContainerInstance`s will be given priority over replacing existing `ContainerInstance`s of the daemon. The new `ContainerInstance`s will be visible only as a change in the `DESIRED` column for the `Deploying` `EnvironmentRevision`.

For example, suppose we add 4 new `ContainerInstance`s at T=3 from the previous example:

| T | Event; actions                 | i<sub>1</sub>                                            | i<sub>2</sub>                                            | i<sub>3</sub>                                            | i<sub>4</sub>                   | i<sub>5</sub>                   | i<sub>6</sub>                    | i<sub>7</sub>                    | i<sub>8</sub>                    |
|:--|:-------------------------------|:---------------------------------------------------------|:---------------------------------------------------------|:---------------------------------------------------------|:--------------------------------|:--------------------------------|:---------------------------------|:---------------------------------|:---------------------------------|
| 3 | Terminated 2; start 2, 6, 7, 8 | :no_entry:<sub>v1</sub> :large_blue_circle:<sub>v2</sub> | :no_entry:<sub>v1</sub> :large_blue_circle:<sub>v2</sub> | :white_check_mark:<sub>v1</sub>                          | :white_check_mark:<sub>v1</sub> | :white_check_mark:<sub>v1</sub> | :large_blue_circle:<sub>v2</sub> | :large_blue_circle:<sub>v2</sub> | :large_blue_circle:<sub>v2</sub> |
| 4 | Running 1                      | :no_entry:<sub>v1</sub> :white_check_mark:<sub>v2</sub>  | :no_entry:<sub>v1</sub> :large_blue_circle:<sub>v2</sub> | :white_check_mark:<sub>v1</sub>                          | :white_check_mark:<sub>v1</sub> | :white_check_mark:<sub>v1</sub> | :large_blue_circle:<sub>v2</sub> | :large_blue_circle:<sub>v2</sub> | :large_blue_circle:<sub>v2</sub> |
| 5 | Running 2                      | :no_entry:<sub>v1</sub> :white_check_mark:<sub>v2</sub>  | :no_entry:<sub>v1</sub> :white_check_mark:<sub>v2</sub>  | :white_check_mark:<sub>v1</sub>                          | :white_check_mark:<sub>v1</sub> | :white_check_mark:<sub>v1</sub> | :large_blue_circle:<sub>v2</sub> | :large_blue_circle:<sub>v2</sub> | :large_blue_circle:<sub>v2</sub> |
| 6 | Running 6                      | :no_entry:<sub>v1</sub> :white_check_mark:<sub>v2</sub>  | :no_entry:<sub>v1</sub> :white_check_mark:<sub>v2</sub>  | :white_check_mark:<sub>v1</sub>                          | :white_check_mark:<sub>v1</sub> | :white_check_mark:<sub>v1</sub> | :white_check_mark:<sub>v2</sub>  | :large_blue_circle:<sub>v2</sub> | :large_blue_circle:<sub>v2</sub> |
| 7 | Running 7; terminate 3         | :no_entry:<sub>v1</sub> :white_check_mark:<sub>v2</sub>  | :no_entry:<sub>v1</sub> :white_check_mark:<sub>v2</sub>  | :red_circle:<sub>v1</sub>                                | :white_check_mark:<sub>v1</sub> | :white_check_mark:<sub>v1</sub> | :white_check_mark:<sub>v2</sub>  | :white_check_mark:<sub>v2</sub>  | :large_blue_circle:<sub>v2</sub> |
| 8 | Terminated 3; start 3          | :no_entry:<sub>v1</sub> :white_check_mark:<sub>v2</sub>  | :no_entry:<sub>v1</sub> :white_check_mark:<sub>v2</sub>  | :no_entry:<sub>v1</sub> :large_blue_circle:<sub>v2</sub> | :white_check_mark:<sub>v1</sub> | :white_check_mark:<sub>v1</sub> | :white_check_mark:<sub>v2</sub>  | :white_check_mark:<sub>v2</sub>  | :large_blue_circle:<sub>v2</sub> |
| 9 | Running 8; terminate 4         | :no_entry:<sub>v1</sub> :white_check_mark:<sub>v2</sub>  | :no_entry:<sub>v1</sub> :white_check_mark:<sub>v2</sub>  | :no_entry:<sub>v1</sub> :large_blue_circle:<sub>v2</sub> | :red_circle:<sub>v1</sub>       | :white_check_mark:<sub>v1</sub> | :white_check_mark:<sub>v2</sub>  | :white_check_mark:<sub>v2</sub>  | :white_check_mark:<sub>v2</sub>  |
| 9 | Running 3; terminate 5         | :no_entry:<sub>v1</sub> :white_check_mark:<sub>v2</sub>  | :no_entry:<sub>v1</sub> :white_check_mark:<sub>v2</sub>  | :no_entry:<sub>v1</sub> :white_check_mark:<sub>v2</sub>  | :red_circle:<sub>v1</sub>       | :red_circle:<sub>v1</sub>       | :white_check_mark:<sub>v2</sub>  | :white_check_mark:<sub>v2</sub>  | :white_check_mark:<sub>v2</sub>  |

At T=3, the new `Task`s are launched immediately because the new `ContainerInstance`s are not running a daemon instance:
```
$ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
REVISION                  STATUS       DESIRED   HEALTHY
SomeEnvironment/Prod:1    Undeploying  3         3
SomeEnvironment/Prod:2    Deploying    5         0
```

At T=4, the flow diverges from the happy path. Since the 3 new `Task`s have not come up yet, the `Environment` is now at `HealthyPercent < 60`. This prevents new `Task`s from being terminated:

```
REVISION                  STATUS       DESIRED   HEALTHY
SomeEnvironment/Prod:1    Undeploying  3        3
SomeEnvironment/Prod:2    Deploying    5        1
```

As the new `Task`s come up, by T=7, enough new `Task`s are running to start terminating old `Task`s again:

```
REVISION                  STATUS       DESIRED   HEALTHY
SomeEnvironment/Prod:1    Undeploying  2        3
SomeEnvironment/Prod:2    Deploying    5        4
```

From T=9 onwards, the flow looks the same as the happy path again.

##### Example 1: Steady state

Let's consider a Daemon update with `MaxOver = 2`, which completes without:
- any changes in cluster membership (i.e. no `ContainerInstance`s join or leave the cluster)
- any `Task`s failing to start up
- any existing `Task`s terminating abnormally

| T | Event                           | i<sub>1</sub>                                                    | i<sub>2</sub>                                                     | i<sub>3</sub>                                                     | i<sub>4</sub>                                                     | i<sub>5</sub>                                                     |
|:--|:--------------------------------|:-----------------------------------------------------------------|:------------------------------------------------------------------|:------------------------------------------------------------------|:------------------------------------------------------------------|:------------------------------------------------------------------|
| 0 | Deployment Starts               | :white_check_mark:<sub>v1</sub>                                  | :white_check_mark:<sub>v1</sub>                                   | :white_check_mark:<sub>v1</sub>                                   | :white_check_mark:<sub>v1</sub>                                   | :white_check_mark:<sub>v1</sub>                                   |
| 1 | Start 1, 2                      | :white_check_mark:<sub>v1</sub> :large_blue_circle:<sub>v2</sub> | :white_check_mark:<sub>v1</sub>  :large_blue_circle:<sub>v2</sub> | :white_check_mark:<sub>v1</sub>                                   | :white_check_mark:<sub>v1</sub>                                   | :white_check_mark:<sub>v1</sub>                                   |
| 2 | Running 1; Terminate 1, Start 3 | :red_circle:<sub>v1</sub>  :white_check_mark:<sub>v2</sub>       | :white_check_mark:<sub>v1</sub>  :large_blue_circle:<sub>v2</sub> | :white_check_mark:<sub>v1</sub> :large_blue_circle: <sub>v2</sub> | :white_check_mark:<sub>v1</sub>                                   | :white_check_mark:<sub>v1</sub>                                   |
| 3 | Running 2; Terminate 2, Start 4 | :no_entry:<sub>v1</sub>  :white_check_mark:<sub>v2</sub>         | :red_circle:<sub>v1</sub>  :white_check_mark:<sub>v2</sub>        | :white_check_mark:<sub>v1</sub> :large_blue_circle: <sub>v2</sub> | :white_check_mark:<sub>v1</sub> :large_blue_circle: <sub>v2</sub> | :white_check_mark:<sub>v1</sub>                                   |
| 4 | Running 3; Terminate 3, Start 5 | :no_entry:<sub>v1</sub>  :white_check_mark:<sub>v2</sub>         | :no_entry:<sub>v1</sub>  :white_check_mark:<sub>v2</sub>          | :red_circle:<sub>v1</sub> :white_check_mark: <sub>v2</sub>        | :white_check_mark:<sub>v1</sub> :large_blue_circle: <sub>v2</sub> | :white_check_mark:<sub>v1</sub> :large_blue_circle: <sub>v2</sub> |
| 5 | Running 4; Terminate 4          | :no_entry:<sub>v1</sub>  :white_check_mark:<sub>v2</sub>         | :no_entry:<sub>v1</sub>  :white_check_mark:<sub>v2</sub>          | :no_entry:<sub>v1</sub> :white_check_mark: <sub>v2</sub>          | :red_circle:<sub>v1</sub> :white_check_mark: <sub>v2</sub>        | :white_check_mark:<sub>v1</sub> :large_blue_circle: <sub>v2</sub> |
| 6 | Running 5; Terminate 5          | :no_entry:<sub>v1</sub>  :white_check_mark:<sub>v2</sub>         | :no_entry:<sub>v1</sub>  :white_check_mark:<sub>v2</sub>          | :no_entry:<sub>v1</sub> :white_check_mark: <sub>v2</sub>          | :no_entry:<sub>v1</sub> :white_check_mark: <sub>v2</sub>          | :red_circle:<sub>v1</sub> :white_check_mark: <sub>v2</sub>        |
| 7 | Deployment Completes            | :no_entry:<sub>v1</sub>  :white_check_mark:<sub>v2</sub>         | :no_entry:<sub>v1</sub>  :white_check_mark:<sub>v2</sub>          | :no_entry:<sub>v1</sub> :white_check_mark: <sub>v2</sub>          | :no_entry:<sub>v1</sub> :white_check_mark: <sub>v2</sub>          | :no_entry:<sub>v1</sub> :white_check_mark: <sub>v2</sub>          |

After the deployment starts at T=1, the only visible difference is that there are now 2 desired `Task`s for `EnvironmentRevision` 2:

```
$ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
REVISION                  STATUS       DESIRED   HEALTHY
SomeEnvironment/Prod:1    Undeploying  5         5
SomeEnvironment/Prod:2    Deploying    2         0
```

At T=2, the first of the new `Task`s are running, and so the first `Task` in the old `EnvironmentRevision` is terminated and another `Task` is started in the new `EnvironmentRevision`:

```
$ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
REVISION                  STATUS       DESIRED   HEALTHY
SomeEnvironment/Prod:1    Undeploying  4         5
SomeEnvironment/Prod:2    Deploying    3         1
```

`Deployment` progresses in a similar fashion until the last `Task` of the old `EnvironmentRevision` is left at T=5:

```
$ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
REVISION                  STATUS       DESIRED   HEALTHY
SomeEnvironment/Prod:1    Undeploying  0         1
SomeEnvironment/Prod:2    Deploying    5         4
```

`Deployment` completes once the final `Task` in the old `EnvironmentRevision` is terminated at T=6 when the final `Task` in the new `EnvironmentRevision` comes up:

```
$ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
REVISION                  STATUS       DESIRED   HEALTHY
SomeEnvironment/Prod:1    Inactive     0         0
SomeEnvironment/Prod:2    Active       5         5
```

##### Example 2: Scale up

For this example, suppose we add a new `ContainerInstance` at T=3 from the previous example:

| T | Event                           | i<sub>1</sub>                                            | i<sub>2</sub>                                              | i<sub>3</sub>                                                     | i<sub>4</sub>                                                     | i<sub>5</sub>                                                     | i<sub>6</sub>                     |
|:--|:--------------------------------|:---------------------------------------------------------|:-----------------------------------------------------------|:------------------------------------------------------------------|:------------------------------------------------------------------|:------------------------------------------------------------------|:----------------------------------|
| 3 | Running 2; Terminate 2, Start 4 | :no_entry:<sub>v1</sub>  :white_check_mark:<sub>v2</sub> | :red_circle:<sub>v1</sub>  :white_check_mark:<sub>v2</sub> | :white_check_mark:<sub>v1</sub> :large_blue_circle: <sub>v2</sub> | :white_check_mark:<sub>v1</sub> :large_blue_circle: <sub>v2</sub> | :white_check_mark:<sub>v1</sub>                                   | :large_blue_circle: <sub>v2</sub> |
| 4 | Running 3; Terminate 3, Start 5 | :no_entry:<sub>v1</sub>  :white_check_mark:<sub>v2</sub> | :no_entry:<sub>v1</sub>  :white_check_mark:<sub>v2</sub>   | :red_circle:<sub>v1</sub> :white_check_mark: <sub>v2</sub>        | :white_check_mark:<sub>v1</sub> :large_blue_circle: <sub>v2</sub> | :white_check_mark:<sub>v1</sub> :large_blue_circle: <sub>v2</sub> | :white_check_mark: <sub>v2</sub>  |
| 5 | Running 4; Terminate 4          | :no_entry:<sub>v1</sub>  :white_check_mark:<sub>v2</sub> | :no_entry:<sub>v1</sub>  :white_check_mark:<sub>v2</sub>   | :no_entry:<sub>v1</sub> :white_check_mark: <sub>v2</sub>          | :red_circle:<sub>v1</sub> :white_check_mark: <sub>v2</sub>        | :white_check_mark:<sub>v1</sub> :large_blue_circle: <sub>v2</sub> | :white_check_mark: <sub>v2</sub>  |
| 6 | Running 5; Terminate 5          | :no_entry:<sub>v1</sub>  :white_check_mark:<sub>v2</sub> | :no_entry:<sub>v1</sub>  :white_check_mark:<sub>v2</sub>   | :no_entry:<sub>v1</sub> :white_check_mark: <sub>v2</sub>          | :no_entry:<sub>v1</sub> :white_check_mark: <sub>v2</sub>          | :red_circle:<sub>v1</sub> :white_check_mark: <sub>v2</sub>        | :white_check_mark: <sub>v2</sub>  |
| 7 | Deployment Completes            | :no_entry:<sub>v1</sub>  :white_check_mark:<sub>v2</sub> | :no_entry:<sub>v1</sub>  :white_check_mark:<sub>v2</sub>   | :no_entry:<sub>v1</sub> :white_check_mark: <sub>v2</sub>          | :no_entry:<sub>v1</sub> :white_check_mark: <sub>v2</sub>          | :no_entry:<sub>v1</sub> :white_check_mark: <sub>v2</sub>          | :white_check_mark: <sub>v2</sub>  |

Since this deployment method doesn't have a `MinHealthyPercent` setting, and the `MaxOver` threshold is not breached, the deployment can continue as normal, even though a new, empty `ContainerInstance` has joined the cluster. It will still promptly schedule a `Task` to be run on the new `ContainerInstance`. The new `ContainerInstance`s will be visible only as a change in the `DESIRED` column for the `Deploying` `EnvironmentRevision` at T=3:

```
$ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
REVISION                  STATUS       DESIRED   HEALTHY
SomeEnvironment/Prod:1    Undeploying  3         4
SomeEnvironment/Prod:2    Deploying    5         2
```

Appendix B: Alternatives considered
-----------------------------------
This section documents alternative paths we considered, but discarded, and why.

### Rollbacks create new `EnvironmentRevision`s

An alternative approach considered for handling rollbacks is to have rollback copy the earlier `EnvironmentRevision` to a new `EnvironmentRevision`, and deploy that. The old `EnvironmentRevision` remains inactive, and is not reused.

```
$ aws ecs deploy-environment --environment-name SomeEnvironment/Prod --environment-revision 1
{ "DeploymentId": "SomeEnvironment/Prod:3@2017090611472301" }
$ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
REVISION                  STATUS       DESIRED   HEALTHY
SomeEnvironment/Prod:1    Inactive     0         0
SomeEnvironment/Prod:2    Undeploying  6         6
SomeEnvironment/Prod:3    Deploying    0         0
$ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
REVISION                  STATUS       DESIRED   HEALTHY
SomeEnvironment/Prod:1    Inactive     0         0
SomeEnvironment/Prod:2    Inactive     0         0
SomeEnvironment/Prod:3    Active       6         6
```

The reason this alternative was discarded, was that we valued the ease with which users can see what's deployed in terms of what they actually requested. Since the rollback `EnvironmentRevision`s were never explicitly created by a human, it's not super clear that they're identical to the older `EnvironmentRevision`s they replace.

Potential reasons to reconsider this, would be if we want to only roll back some part of the `Environment`, e.g. only the `TaskDefinition`, and not the `InstanceGroup`.

### Implicitly resume `Environment` on deployment

We considered automatically resuming an `Environment` when a new `deploy-environment` or `rollback-environment` command is issued. This would prevent customers from getting stuck because they forgot to resume their `Environment`.

The reason that this alternative was discarded is that automated deployment systems could inadvertently cause an `Environment` to be resumed while it is not actually safe to do so. However, we could consider achieving the same thing by including a `--force` flag or similar on the `deploy-environment`/`rollback-environment` commands.
