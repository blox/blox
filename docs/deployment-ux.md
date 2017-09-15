Blox Deployments User Experience
================================

Use Cases
---------

The following are typical use cases for Deployments:

-   Create an Environment to deploy a new environment revision. The
    Deployment creates Tasks in the background. Check the status of the
    deployment to see if it succeeds or not.
-   Declare the new state of the Environment by updating it, then deploy
    it by calling startDeployment. A new Environment revision is created
    and the Deployment manages moving the Tasks from the old revision to
    the new one at a controlled rate.
-   Rollback to an earlier Environment revision if the current state of
    the Environment is not stable. Each rollback updates the revision of
    the environment (it does not simply revert the environment
    revision).
-   (not applicable for Daemon) Scale up the Deployment to facilitate
    more load.
-   Pause the Deployment to apply multiple fixes to its Environment and
    then resume it to start a new deployment.
-   Use the status of the Environment as an indicator that a deployment
    has become stuck.

UI Conventions
--------------
- **Showing `Inactive` revisions:** For all the output examples below, we're including environment revisions in the `Inactive` state by specifying the `--all` flag to `describe-environment-status`. Without this flag, only revisions that are not in `Inactive` will be shown.
- For all the cluster state visualizations below, the legend is:

  | :large_blue_circle: pending | :white_check_mark: running | :red_circle: terminating | :no_entry: terminated |
  |:----------------------------|:---------------------------|:-------------------------|:----------------------|

Creating and deploying an Environment
-------------------------------------

In order to create a deployment, you first have to create an
environment that declares what should be running, where, and how the
deployment should be controlled.

```
$ aws ecs create-environment --yaml-file <<EOF
Name: "SomeEnvironment/Prod"
TaskDefinition: "log-daemon:1"
Cluster: "default"
InstanceGroup:
  Query: "attribute:stack == prod"
DeploymentType: Daemon
DeploymentConfiguration:
  Method: "Replace" # (or "Surge")
  MinHealthyPercent: "100"
EOF
{ "EnvironmentRevision": 1 }
```

At this point, nothing has actually happened to your cluster. You can
see this by showing your environment's status:

```
$ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
REVISION                  STATUS     DESIRED   CURRENT
SomeEnvironment/Prod:1    Created    0         0
```

In order to make this configuration Active, you first have to Deploy the
revision:

```
$ aws ecs deploy-environment --environment-name SomeEnvironment/Prod --environment-revision 1
{ "DeploymentId": "SomeEnvironment/Prod:1@2017090611472301" }
```

This will mark the revision as active, and start launching tasks to meet
the deployment configuration's constraints. In the case of a Daemon
environment, the desired target is to have one copy of the task for
every instance in the target Instance Group. For the initial deployment,
MinHealthyPercent is ignored, since we're starting from a situation that
already breaches MinHealthyPercent

```
$ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
REVISION                  STATUS          DESIRED   CURRENT
SomeEnvironment/Prod:1    Deploying       5         0
```

As the deployment progresses:

```
$ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
REVISION                  STATUS          DESIRED   CURRENT
SomeEnvironment/Prod:1    Deploying       5         0
```

```
$ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
REVISION                  STATUS          DESIRED   CURRENT
SomeEnvironment/Prod:1    Deploying       5         3
```

Failure:

```
$ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
REVISION                  STATUS          DESIRED   CURRENT
SomeEnvironment/Prod:1    Stuck           5         4
```

Success:

```
$ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
REVISION                  STATUS          DESIRED   CURRENT
SomeEnvironment/Prod:1    Active          5         5
```

Updating an environment
-----------------------

Typically, you'd update an environment in order to:

-   update the task definition that should be running
-   change the set of instances on which to run tasks
-   change the deployment configuration (i.e. to change health
    thresholds) (TODO: Maybe require new env for this?)

### Updating the task definition

If you don't specify all attributes, you must specify the revision you
wish to base your changes on. This prevents potential conflicts if
multiple callers attempt to update an environment at the same time.

```
$ aws ecs update-environment --environment SomeEnvironment/Prod --revision 1 --task-definition "log-daemon:2"
{ "EnvironmentRevision": 1 }
```

At this point, nothing has actually happened to your cluster. You can
see this by showing your environment's status:

```
$ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
REVISION                  STATUS     DESIRED   CURRENT
SomeEnvironment/Prod:1    Active     5         5
SomeEnvironment/Prod:2    Created    0         0
```

In order to make this configuration Active, you first have to Deploy the
revision:

```
$ aws ecs deploy-environment --environment-name SomeEnvironment/Prod --environment-revision 2
{ "DeploymentId": "SomeEnvironment/Prod:2@2017090611472301" }
```

(Note that specifying a revision older than the latest `Active` revision is not supported; the only way to deploy an older revision than that is to use the `rollback-environment` command. This prevents inadvertently rolling back changes.)

This will mark the revision as Deploying, and start launching tasks to
meet the deployment configuration's constraints. In the case of a Daemon
environment, the desired target is to have one copy of the task for
every instance in the target Instance Group:

```
$ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
REVISION                  STATUS       DESIRED   CURRENT
SomeEnvironment/Prod:1    Undeploying  0         5
SomeEnvironment/Prod:2    Deploying    5         0
```

Once the deployment reaches a steady state, the old revision is marked
Inactive, and the new revision is marked Active:

```
$ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
REVISION                  STATUS       DESIRED   CURRENT
SomeEnvironment/Prod:1    Inactive     0         0
SomeEnvironment/Prod:2    Active       5         5
```

### Updating the Instance Group

If you don't specify all attributes, you must specify the revision you
wish to base your changes on. This prevents potential conflicts if
multiple callers attempt to update an environment at the same time.

```
$ aws ecs update-environment --environment SomeEnvironment/Prod --revision 1 --instance-group-query "attribute:stack in (prod, gamma)"
{ "EnvironmentRevision": 2 }
```

At this point, nothing has actually happened to your cluster. You can
see this by showing your environment's status:

```
$ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
REVISION                  STATUS     DESIRED   CURRENT
SomeEnvironment/Prod:1    Active     5         5
SomeEnvironment/Prod:2    Created    0         0
```

In order to make this configuration Active, you first have to Deploy the
revision:

```
$ aws ecs deploy-environment --environment-name SomeEnvironment/Prod --environment-revision 2
{ "DeploymentId": "SomeEnvironment/Prod:2@2017090611472301" }
```

This will mark the revision as Deploying, and start launching tasks to
meet the deployment configuration's constraints. In the case of a Daemon
environment, the desired target is to have one copy of the task for
every instance in the target Instance Group. Since the size of the
instance group is now different because it was changed, that's reflected
in the DESIRED column for the new revision:

```
$ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
REVISION                  STATUS       DESIRED   CURRENT
SomeEnvironment/Prod:1    Undeploying  0         5
SomeEnvironment/Prod:2    Deploying    6         0
```

Note that even though this doesn't change the task definition, it will
still stop and restart tasks. (TODO: Can we make it smarter, so that it
only launches/terminates tasks to end up in the right group?)

Once the deployment reaches a steady state, the old revision is marked
Inactive, and the new revision is marked Active:

```
$ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
REVISION                  STATUS       DESIRED   CURRENT
SomeEnvironment/Prod:1    Inactive     0         0
SomeEnvironment/Prod:2    Active       6         6
```

### Daemon Environments
Daemon Environments behave slightly differently from Service Environments (TBD). Daemon tasks typically provide instance-wide capabilities such as metrics/logging or network overlay, that the other tasks on the instance all depend upon. Any amount of time that the task is not present on the instance, could potentially result in downtime for the other tasks on the instance.

Because of this, at any time during the life of a Daemon environment, Blox will prioritize:
- minimizing the number of instances that are not running the Daemon task at all
- minimizing the number of instances that are not running the version of the Daemon task that they should be

Since tasks sometimes fail, and cluster instance membership changes over time, it is expected that there will often be some window of time where the requirements above cannot be met. In particular, for larger clusters, or clusters that are backed by Autoscaling Groups, the available instances in the cluster might change frequently (and sometimes dramatically).

The goal then is to bound the size of the window where a instance is not running the Daemon environment, and to provide alarms on when this upper bound is breached. If a configurable number of instances are without the correct version of a Daemon task for a configurable amount of time, the environment revision will enter the `Failing` state and emit an alarm. However, it will continue to try and make progress.

- If a deployment fails to start a task on a particular host (e.g. because the task definition is faulty, or the Task fails to start up in a nondeterministic way), priority will be given to restarting that task, rather than continuing to replace existing tasks.
- If a Daemon task terminates abnormally during a deployment on an instance that is not currently having its task replaced, priority will be given to replacing that task over starting to terminate tasks on other instances.
- If new hosts join a cluster, new Daemon tasks will launch on them promptly (without waiting for the evaluation of deployment constraints or for other workflows to finish)

Daemon deployments have two update strategies, with different availability impacts:

#### Terminate before replace

If you're using this deployment strategy, then Blox will kill
`(100 - MinHealthyPercent)%` of tasks, and only launch new tasks to
replace as they are successfully terminated. This will ensure that
there's never more than one version of the Task running at the same
time on the same host, but may result in some downtime for the
daemon. This is a good fit for daemons that would fail if there is
more than one running per instance, and that can tolerate short
periods of unavailability.

##### Example 1: Steady State

Let's consider a Daemon update with `MinHealthyPercent = 60`, which completes without:
- any changes in cluster membership (i.e. no instances join or leave the cluster)
- any tasks failing to start up
- any existing tasks terminating abnormally

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
REVISION                  STATUS       DESIRED   CURRENT
SomeEnvironment/Prod:1    Undeploying  3         5
SomeEnvironment/Prod:2    Deploying    0         0
```

At *T=1*, two tasks in the old revision get terminated, and replaced with two pending tasks:

```
$ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
REVISION                  STATUS       DESIRED   CURRENT
SomeEnvironment/Prod:1    Undeploying  3         3
SomeEnvironment/Prod:2    Deploying    2         0
```

At *T=4*, one of the tasks finish terminating, and it's replaced by a new task, so we terminate another old task:

```
$ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
REVISION                  STATUS       DESIRED   CURRENT
SomeEnvironment/Prod:1    Undeploying  2         3
SomeEnvironment/Prod:2    Deploying    3         1
```

Deployment progresses in a similar fashion until none of the old revision is left at *T=10*:

```
$ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
REVISION                  STATUS       DESIRED   CURRENT
SomeEnvironment/Prod:1    Undeploying  0         0
SomeEnvironment/Prod:2    Deploying    5         4
```

The final task is launched at *T=11* and the new revision becomes active:

```
$ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
REVISION                  STATUS       DESIRED   CURRENT
SomeEnvironment/Prod:1    Inactive     0         0
SomeEnvironment/Prod:2    Active       5         5
```

##### Example 2: Scale up

If new instances join the cluster while an environment is deploying, running the daemon task on those instances will be given priority over replacing existing instances of the daemon. The new instances will be visible only as a change in the `DESIRED` column for the `Deploying` environment revision.

For example, suppose we add 4 new instances at T=3 from the previous example:

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

At T=3, the new tasks are launched immediately because the new hosts are not running a daemon instance:
```
$ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
REVISION                  STATUS       DESIRED   CURRENT
SomeEnvironment/Prod:1    Undeploying  3         3
SomeEnvironment/Prod:2    Deploying    5         0
```

At T=4, the flow diverges from the happy path. Since the 3 new tasks have not come up yet, the environment is now at `HealthyPercent < 60`. This prevents new tasks from being terminated:

```
REVISION                  STATUS       DESIRED   CURRENT
SomeEnvironment/Prod:1    Undeploying  3        3
SomeEnvironment/Prod:2    Deploying    5        1
```

As the new tasks come up, by T=7, enough new tasks are running to start terminating old tasks again:

```
REVISION                  STATUS       DESIRED   CURRENT
SomeEnvironment/Prod:1    Undeploying  2        3
SomeEnvironment/Prod:2    Deploying    5        4
```

From T=9 onwards, the flow looks the same as the happy path again.

#### Terminate after replace

If you're using this deployment strategy, then Blox will launch `MaxOver` tasks on as many instances, and only kill old tasks once the new task is confirmed to be running. This will ensure that there's never less than one Task running at the same time on the same Instance, but will result in two copies of the Task running on at most `MaxOver` instances at any given time. This is a good fit for daemons that must always have at least one instance running, and that don't care if another instance of the daemon is running on the same host.

##### Example 1: Steady state

Let's consider a Daemon update with `MaxOver = 2`, which completes without:
- any changes in cluster membership (i.e. no instances join or leave the cluster)
- any tasks failing to start up
- any existing tasks terminating abnormally

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

After the deployment starts at T=1, the only visible difference is that there are now 2 desired tasks for revision 2:

```
$ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
REVISION                  STATUS       DESIRED   CURRENT
SomeEnvironment/Prod:1    Undeploying  5         5
SomeEnvironment/Prod:2    Deploying    2         0
```

At  T=2, the first of the new tasks are running, and so the first task in the old revision is terminated and another task is started in the new revision:

```
$ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
REVISION                  STATUS       DESIRED   CURRENT
SomeEnvironment/Prod:1    Undeploying  4         5
SomeEnvironment/Prod:2    Deploying    3         1
```

Deployment progresses in a similar fashion until the last task of the old revision is left at T=5:

```
$ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
REVISION                  STATUS       DESIRED   CURRENT
SomeEnvironment/Prod:1    Undeploying  0         1
SomeEnvironment/Prod:2    Deploying    5         4
```

Deployment completes once the final task in the old revision is terminated at T=6 when the final task in the new revision comes up:

```
$ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
REVISION                  STATUS       DESIRED   CURRENT
SomeEnvironment/Prod:1    Inactive     0         0
SomeEnvironment/Prod:2    Active       5         5
```

##### Example 2: Scale up

For this example, suppose we add a new instance at T=3 from the previous example:

| T | Event                           | i<sub>1</sub>                                            | i<sub>2</sub>                                              | i<sub>3</sub>                                                     | i<sub>4</sub>                                                     | i<sub>5</sub>                                                     | i<sub>6</sub>                     |
|:--|:--------------------------------|:---------------------------------------------------------|:-----------------------------------------------------------|:------------------------------------------------------------------|:------------------------------------------------------------------|:------------------------------------------------------------------|:----------------------------------|
| 3 | Running 2; Terminate 2, Start 4 | :no_entry:<sub>v1</sub>  :white_check_mark:<sub>v2</sub> | :red_circle:<sub>v1</sub>  :white_check_mark:<sub>v2</sub> | :white_check_mark:<sub>v1</sub> :large_blue_circle: <sub>v2</sub> | :white_check_mark:<sub>v1</sub> :large_blue_circle: <sub>v2</sub> | :white_check_mark:<sub>v1</sub>                                   | :large_blue_circle: <sub>v2</sub> |
| 4 | Running 3; Terminate 3, Start 5 | :no_entry:<sub>v1</sub>  :white_check_mark:<sub>v2</sub> | :no_entry:<sub>v1</sub>  :white_check_mark:<sub>v2</sub>   | :red_circle:<sub>v1</sub> :white_check_mark: <sub>v2</sub>        | :white_check_mark:<sub>v1</sub> :large_blue_circle: <sub>v2</sub> | :white_check_mark:<sub>v1</sub> :large_blue_circle: <sub>v2</sub> | :white_check_mark: <sub>v2</sub>  |
| 5 | Running 4; Terminate 4          | :no_entry:<sub>v1</sub>  :white_check_mark:<sub>v2</sub> | :no_entry:<sub>v1</sub>  :white_check_mark:<sub>v2</sub>   | :no_entry:<sub>v1</sub> :white_check_mark: <sub>v2</sub>          | :red_circle:<sub>v1</sub> :white_check_mark: <sub>v2</sub>        | :white_check_mark:<sub>v1</sub> :large_blue_circle: <sub>v2</sub> | :white_check_mark: <sub>v2</sub>  |
| 6 | Running 5; Terminate 5          | :no_entry:<sub>v1</sub>  :white_check_mark:<sub>v2</sub> | :no_entry:<sub>v1</sub>  :white_check_mark:<sub>v2</sub>   | :no_entry:<sub>v1</sub> :white_check_mark: <sub>v2</sub>          | :no_entry:<sub>v1</sub> :white_check_mark: <sub>v2</sub>          | :red_circle:<sub>v1</sub> :white_check_mark: <sub>v2</sub>        | :white_check_mark: <sub>v2</sub>  |
| 7 | Deployment Completes            | :no_entry:<sub>v1</sub>  :white_check_mark:<sub>v2</sub> | :no_entry:<sub>v1</sub>  :white_check_mark:<sub>v2</sub>   | :no_entry:<sub>v1</sub> :white_check_mark: <sub>v2</sub>          | :no_entry:<sub>v1</sub> :white_check_mark: <sub>v2</sub>          | :no_entry:<sub>v1</sub> :white_check_mark: <sub>v2</sub>          | :white_check_mark: <sub>v2</sub>  |

Since this deployment method doesn't have a `MinHealthyPercent` setting, and the `MaxOver` threshold is not breached, the deployment can continue as normal, even though a new, empty instance has joined the cluster. It will still promptly schedule an instance to be run on the new instance. The new instances will be visible only as a change in the `DESIRED` column for the `Deploying` environment revision at T=3:

```
$ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
REVISION                  STATUS       DESIRED   CURRENT
SomeEnvironment/Prod:1    Undeploying  3         4
SomeEnvironment/Prod:2    Deploying    5         2
```

Rollback to an earlier revision
-------------------------------

Let's say that revision 2 introduced an Environment change that is actually broken, and we need to roll back to revision 1:

```
$ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
REVISION                  STATUS       DESIRED   CURRENT
SomeEnvironment/Prod:1    Inactive     0         0
SomeEnvironment/Prod:2    Active       5         5
```


In order to do a prompt rollback to an earlier revision, we can use the `rollback-environment` command. The specified revision becomes active again, and we deactivate the bad revision. The bad revision will enter the `Reverted` state once it is done `Reverting`.

```
$ aws ecs deploy-environment --environment-name SomeEnvironment/Prod --environment-revision 1
{ "DeploymentId": "SomeEnvironment/Prod:1@2017090611472301" }

$ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
REVISION                  STATUS       DESIRED   CURRENT
SomeEnvironment/Prod:1    Deploying    5         0
SomeEnvironment/Prod:2    Reverting    0         5

$ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
REVISION                  STATUS       DESIRED   CURRENT
SomeEnvironment/Prod:1    Active       5         5
SomeEnvironment/Prod:2    Reverted     0         0
```

### In-progress deployments

Note that even though the above example shows a revision that's already `Active`, even a `Deploying` revision could still be rolled back. You could even specify a rollback revision that's older than the one that is currently `Undeploying`; in this case, both the `Deploying` and `Undeploying` revisions will be reverted.

### Rollback deployment configuration

Rollback deployments may have different deployment configurations to allow for quicker recovery, or to avoid getting stuck due to breached deployment constraints such as `MinHealhtyPercent`.

Pause ongoing deployments
-------------------------

In some scenarios, it may be dangerous to continue running or terminating tasks in a cluster. For example, there might be some problem with the underlying cluster that causes new tasks to fail, or some external dependency being unavailable might make it impossible for tasks to start/stop cleanly.

For these cases, all actions taken in an Environment can be suspended using the `freeze-environment` command:

```
$ aws ecs freeze-environment --environment-name SomeEnvironment/Prod
$ aws ecs describe-environment --environment-name SomeEnvironment/Prod --table
ENVIRONMENT               STATUS       DESIRED   CURRENT
SomeEnvironment/Prod      Suspended    6         3
```

Once it is safe for the Environment to take action again, it can be resumed using the `thaw-environment` command:

```
$ aws ecs thaw-environment --environment-name SomeEnvironment/Prod
$ aws ecs describe-environment --environment-name SomeEnvironment/Prod --table
ENVIRONMENT               STATUS       DESIRED   CURRENT
SomeEnvironment/Prod      Deploying    6         3
```

Note that suspending an environment doesn't prevent you from starting a new deployment; however, in order for any deployments to progress, the Environment must be resumed first.

Detect stuck deployments
------------------------

In some cases, it might be impossible for an environment to make progress towards making a new Environment revision `Active`. For example, if a Daemon Environment is trying to launch a new copy of a Daemon task in replace-before-terminate mode, it could fail to do so because the instance does not have enough free resources.


In this case, the environment will enter into the `Stuck` (come up with a better name) state:

```
$ aws ecs describe-environment-status --environment-name SomeEnvironment/Prod --table
REVISION                  STATUS       DESIRED   CURRENT
SomeEnvironment/Prod:2    Undeploying  0         3
SomeEnvironment/Prod:3    Stuck        6         3
```

Once an Environment has been in this state for a user-configurable amount of time, it will enter into the `Failing` (come up with a better name) state, and optionally emit a CloudWatch Event that you can alarm on:

```
$ aws ecs describe-environment-status --environment-name SomeEnvironment/Prod --table
REVISION                  STATUS       DESIRED   CURRENT
SomeEnvironment/Prod:2    Stuck        0         3
SomeEnvironment/Prod:3    Failing      6         3
```

Alternatives considered
-----------------------
This section documents alternative paths we considered, but discarded, and why.

### Rollbacks create new revisions

An alternative approach considered for handling rollbacks is to have rollback copy the earlier revision to a new revision, and deploy that. The old revision remains inactive, and is not reused.

```
$ aws ecs deploy-environment --environment-name SomeEnvironment/Prod --environment-revision 1
{ "DeploymentId": "SomeEnvironment/Prod:3@2017090611472301" }
$ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
REVISION                  STATUS       DESIRED   CURRENT
SomeEnvironment/Prod:1    Inactive     0         0
SomeEnvironment/Prod:2    Undeploying  6         6
SomeEnvironment/Prod:3    Deploying    0         0
$ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
REVISION                  STATUS       DESIRED   CURRENT
SomeEnvironment/Prod:1    Inactive     0         0
SomeEnvironment/Prod:2    Inactive     0         0
SomeEnvironment/Prod:3    Active       6         6
```

The reason this alternative was discarded, was that we valued the ease with which users can see what's deployed in terms of what they actually requested. Since the rollback revisions were never explicitly created by a human, it's not super clear that they're identical to the older revisions they replace.

### Implicitly resume Environment on deployment

We considered automatically resuming an environment when a new `deploy-environment` or `rollback-environment` command is issued. This would prevent customers from getting stuck because they forgot to resume their Environment.

The reason that this alternative was discarded is that automated deployment systems could inadvertently cause an environment to be resumed while it is not actually safe to do so. However, we could consider achieving the same thing by including a `--force` flag or similar on the `deploy-environment`/`rollback-environment` commands.

### Don't prioritize placing tasks on new instances

We considered not prioritizing placing tasks on new instances. However, this would result in actually decreasing the cluster's `HealthyPercent` from T=4-13, until we finally place new tasks on the new instances after all existing workflows are done.

