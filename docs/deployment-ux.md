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


### Daemon deployments

There are some deployment peculiarities specific to Daemon Environments.

#### Deployment strategies
Daemon deployments have two update strategies, with different availability
impacts

1.  Terminate before replace

    If you're using this deployment strategy, then Blox will kill
    `(100 - MinHealthyPercent)%` of tasks, and only launch new tasks to
    replace as they are successfully terminated. This will ensure that
    there's never more than one version of the Task running at the same
    time on the same host, but may result in some downtime for the
    daemon. This is a good fit for daemons that would fail if there is
    more than one running per instance, and that can tolerate short
    periods of unavailability.

    After deployment starts:

    ```
    $ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
    REVISION                  STATUS       DESIRED   CURRENT
    SomeEnvironment/Prod:1    Undeploying  0         5
    SomeEnvironment/Prod:2    Deploying    5         0
    ```

    First task in old revision is terminated:

    ```
    $ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
    REVISION                  STATUS       DESIRED   CURRENT
    SomeEnvironment/Prod:1    Undeploying  0         4
    SomeEnvironment/Prod:2    Deploying    5         0
    ```

    First task in new revision is launched:

    ```
    $ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
    REVISION                  STATUS       DESIRED   CURRENT
    SomeEnvironment/Prod:1    Undeploying  0         4
    SomeEnvironment/Prod:2    Deploying    5         1
    ```

    Deployment progresses in a similar fashion until none of the old
    revision is left:

    ```
    $ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
    REVISION                  STATUS       DESIRED   CURRENT
    SomeEnvironment/Prod:1    Undeploying  0         0
    SomeEnvironment/Prod:2    Deploying    5         4
    ```

    The final task is launched and the new revision becomes active:

    ```
    $ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
    REVISION                  STATUS       DESIRED   CURRENT
    SomeEnvironment/Prod:1    Inactive     0         0
    SomeEnvironment/Prod:2    Active       5         5
    ```

2.  Terminate after replace

    If you're using this deployment strategy, then Blox will launch
    `MaxOver` tasks on as many instances, and only kill old tasks once
    the new task is confirmed to be running. This will ensure that
    there's never less than one Task running at the same time on the
    same Instance, but will result in two copies of the Task running on
    at most MaxSurge instances at any given time. This is a good fit for
    daemons that must always have at least one instance running, and
    that don't care if another instance of the daemon is running on the
    same host.

    After deployment starts:

    ```
    $ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
    REVISION                  STATUS       DESIRED   CURRENT
    SomeEnvironment/Prod:1    Undeploying  0         5
    SomeEnvironment/Prod:2    Deploying    5         0
    ```

    First task in new revision is launched, old revision is still
    running everywhere:

    ```
    $ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
    REVISION                  STATUS       DESIRED   CURRENT
    SomeEnvironment/Prod:1    Undeploying  0         5
    SomeEnvironment/Prod:2    Deploying    5         1
    ```

    First task in old revision is terminated:

    ```
    $ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
    REVISION                  STATUS       DESIRED   CURRENT
    SomeEnvironment/Prod:1    Undeploying  0         4
    SomeEnvironment/Prod:2    Deploying    5         1
    ```

    Deployment progresses in a similar fashion until the last task of
    the old revision is left:

    ```
    $ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
    REVISION                  STATUS       DESIRED   CURRENT
    SomeEnvironment/Prod:1    Undeploying  0         1
    SomeEnvironment/Prod:2    Deploying    5         5
    ```

    Deployment completes once the final task in the old revision is
    terminated:

    ```
    $ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
    REVISION                  STATUS       DESIRED   CURRENT
    SomeEnvironment/Prod:1    Inactive     0         0
    SomeEnvironment/Prod:2    Active       5         5
    ```

#### Changes to cluster membership while a deployment is ongoing

For larger clusters, or clusters that are backed by Autoscaling Groups, the available instances in the cluster might change (sometimes dramatically). Since the total number of available instances determines the desired end-state of the Environment, these state changes could affect the ability of an Environment to make progress while honoring its deployment constraints.

The behaviour we're aiming for is to minimize the number of instances that are not running the Daemon task at all at any given point in time. If a configurable number of instances are without the correct version of a Daemon task for a configurable amount of time, the deployment will enter the `Failing` state and emit an alarm. However, it will continue to try and make progress.

- If a deployment fails to start a task on a particular host (e.g. because the task definition is faulty, or the Task fails to start up in a nondeterministic way), priority will be given to restarting that task, rather than continuing to replace existing tasks.
- If a Daemon task terminates abnormally during a deployment on an instance that is not currently having its task replace, priority will be given to replacing that task over starting to terminate tasks on other instances.
- If new instances join the cluster while an environment is deploying, running the daemon task on those instances will be given priority over replacing existing instances of the daemon. The new instances will be visible only as a change in the `DESIRED` column for the `Deploying` environment revision:

  ```
  $ aws ecs describe-environment-status --all --environment-name SomeEnvironment/Prod --table
  REVISION                  STATUS       DESIRED   CURRENT
  SomeEnvironment/Prod:1    Undeploying  0         3
  SomeEnvironment/Prod:2    Deploying    5         2
  ```

  New instances matching the Instance Group join the cluster:

  ```
  REVISION                  STATUS       DESIRED   CURRENT
  SomeEnvironment/Prod:1    Undeploying  0         3
  SomeEnvironment/Prod:2    Deploying    10        8
  ```

- If instances in a cluster are terminated while a deployment is ongoing... (TODO: then what? Do we care?)

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

The reason that this alternative was discarded. However, this means that automated deployment systems could inadvertently cause an environment to be resumed while it is not actually safe to do so.
