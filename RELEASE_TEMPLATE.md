## Subject Line
[Blox] Deploy release - v$$blox_version$$

## Activity Details
This activity is to: Release Blox version v$$blox_version$$ by merging the dev branch into the master branch on GitHub, and pushing new Blox images up to Docker Hub.

The purpose of this change is to: Release the latest features and bug fixes for Blox to allow consumers to start using the new functionality. For more details about the specific changes being released, refer to the $$blox_version$$ release notes [here](https://github.com/blox/blox/blob/dev/CHANGELOG.md).

The Blox framework version is: v$$blox_version$$ with a git hash of $$github_hash$$

## Impact Details
What will happen if this release doesn't happen?  
Blox consumers will not be able to take advantage of the new features and bug fixes in Blox.

Why is this the correct time/day to complete the release?  
This corresponds to the completion date of all GitHub issues planned for this release.

Are there any related, prerequisite changes upon which this release hinges?  
No.

## Worst Case Scenario
What could happen if this change causes impact?  
Blox consumers could start running this new version of the Blox framework and encounter random errors leading to unexpected behavior.

Where are the most likely places this change will fail?  
Errors encountered building the new Blox images or pushing them up to Docker Hub.

## Hostname or Service
GitHub: https://github.com/blox/blox  
Docker Hub: https://hub.docker.com/u/bloxoss/

## Timeline / Activity Plan
Times are relative to the start of the release.

- 00:00 Change the release status to "In Progress".
- 00:05 Perform the Activity Checklist steps below.
- 00:25 Perform the Validation Checklist steps below.
- 01:00 Activity complete. Change the release status to "Complete".

## Activity Checklist
Refer to the Validation Activities section for more details about each release step.

- [ ] Prereqs to be performed prior to the release starting:
- [ ] Ensure that the cluster-state-service unit-tests pass.
- [ ] Ensure that the cluster-state-service end-to-end tests pass.
- [ ] Ensure that the daemon-scheduler unit-tests pass.
- [ ] Ensure that the daemon-scheduler end-to-end tests pass.
- [ ] Ensure that the daemon-scheduler integration tests pass.
- [ ] Ensure that `deploy/docker/conf/docker-compose.yml` has been updated with the correct image versions.
- [ ] Ensure that `deploy/aws/conf/cloudformation_template.json` has been updated with the correct image versions.
- [ ] Ensure that `IamRoleInstance,IamRoleService,IamRoleTask,IamRoleLambda` in `deploy/aws/conf/cloudformation_template.json` have the correct role permissions.
- [ ] Ensure that a $$blox_version$$ change log entry exists in `CHANGELOG.md`.
- [ ] Ensure that the version is set to $$blox_version$$ in `cluster-state-service/README.md` and `daemon-scheduler/README.md`.
- [ ] Ensure that the version is set to $$blox_version$$ in `cluster-state-service/VERSION` and `daemon-scheduler/VERSION`.
- [ ] Ensure that the version is set to $$blox_version$$ in `cluster-state-service/versioning/version.go` and `daemon-scheduler/versioning/version.go`.
- [ ] Ensure that the technician has a GitHub account and permissions to merge into the Blox master branch.
- [ ] Ensure that the technician has a Docker Hub account and permissions to push into the bloxoss repositories.
- [ ] Ensure that the technician has logged into Docker Hub via `docker logout` and `docker login`.
- [ ] Release activities:
- [ ] Merge the GitHub pull request into the master branch.
- [ ] Create a v$$blox_version$$ release tag in GitHub with commit $$github_hash$$.
- [ ] Delete all bloxoss local images before building new ones.
- [ ] Build the bloxoss/cluster-state-service:latest,$$blox_version$$,$$github_hash$$ images and publish to Docker Hub.
- [ ] Build the bloxoss/daemon-scheduler:latest,$$blox_version$$,$$github_hash$$ images and publish to Docker Hub.

## Rollback Procedure
- Log into the Docker Hub console and delete the bloxoss/cluster-state-service:latest,$$blox_version$$,$$github_hash$$ images.
- Log into the Docker Hub console and delete the bloxoss/daemon-scheduler:latest,$$blox_version$$,$$github_hash$$ images.
- Publish the bloxoss/cluster-state-service:latest tag pointed to the previous release.
- Publish the bloxoss/daemon-scheduler:latest tag pointed to the previous release.
- Revert Blox repository changes.

#### Publish the bloxoss/cluster-state-service:latest tag pointed to the previous release
```
$ docker pull bloxoss/cluster-state-service:$$previous_version$$
$ docker tag bloxoss/cluster-state-service:$$previous_version$$ bloxoss/cluster-state-service:latest
$ docker push bloxoss/cluster-state-service:latest
```

#### Publish the bloxoss/daemon-scheduler:latest tag pointed to the previous release
```
$ docker pull bloxoss/daemon-scheduler:$$previous_version$$
$ docker tag bloxoss/daemon-scheduler:$$previous_version$$ bloxoss/daemon-scheduler:latest
$ docker push bloxoss/daemon-scheduler:latest
```

#### Revert Blox repository changes
```
$ cd <GitRepoBase>
$ git revert v$$blox_version$$
$ git push
$ git fetch --tags origin
$ git tag -d v$$blox_version$$
$ git push origin :refs/tags/v$$blox_version$$
```

## Validation Activities

#### Delete all bloxoss local images before building new ones
```
$ docker images -a | grep bloxoss | awk '{print $3}' | xargs docker rmi -f
$ echo $?
0  <- should return 0
```

#### Build the bloxoss/cluster-state-service:latest,$$blox_version$$,$$github_hash$$ images and publish to Docker Hub
```
$ cd <GitRepoBase>/cluster-state-service/
$ make release
$ docker tag bloxoss/cluster-state-service:latest bloxoss/cluster-state-service:$$blox_version$$
$ docker tag bloxoss/cluster-state-service:latest bloxoss/cluster-state-service:$$github_hash$$
$ docker push bloxoss/cluster-state-service:latest
$ docker push bloxoss/cluster-state-service:$$blox_version$$
$ docker push bloxoss/cluster-state-service:$$github_hash$$
```

#### Build the bloxoss/daemon-scheduler:latest,$$blox_version$$,$$github_hash$$ images and publish to Docker Hub
```
$ cd <GitRepoBase>/daemon-scheduler/
$ make release
$ docker tag bloxoss/daemon-scheduler:latest bloxoss/daemon-scheduler:$$blox_version$$
$ docker tag bloxoss/daemon-scheduler:latest bloxoss/daemon-scheduler:$$github_hash$$
$ docker push bloxoss/daemon-scheduler:latest
$ docker push bloxoss/daemon-scheduler:$$blox_version$$
$ docker push bloxoss/daemon-scheduler:$$github_hash$$
```

#### Ensure that the cluster-state-service unit-tests pass
```
$ cd <GitRepoBase>/cluster-state-service/
$ make unit-tests
...
$ echo $?
0  <- should return 0
```

#### Ensure that the cluster-state-service end-to-end tests pass
Required setup details are listed [here](https://github.com/blox/blox/blob/dev/cluster-state-service/internal/Readme.md). It is assumed that the cluster-state-service is running locally.

```
$ cd <GitRepoBase>/cluster-state-service/
$ AWS_REGION=<region> AWS_PROFILE=<profile> ECS_CLUSTER=<cluster> gucumber -tags=@e2e
...
$ echo $?
0  <- should return 0
```

#### Ensure that the daemon-scheduler unit-tests pass
```
$ cd <GitRepoBase>/daemon-scheduler/
$ make unit-tests
...
$ echo $?
0  <- should return 0
```

#### Ensure that the daemon-scheduler end-to-end tests pass
Required setup details are listed [here](https://github.com/blox/blox/blob/dev/daemon-scheduler/internal/features/README.md). It is assumed that the daemon-scheduler is running locally.

```
$ cd <GitRepoBase>/daemon-scheduler/
$ AWS_REGION=<region> AWS_PROFILE=<profile> ECS_CLUSTER=<cluster> ECS_CLUSTER_ASG=<autoscale_group> gucumber -tags=@e2e
...
$ echo $?
0  <- should return 0
```

#### Ensure that the daemon-scheduler integration tests pass
Required setup details are listed [here](https://github.com/blox/blox/blob/dev/daemon-scheduler/internal/features/README.md). It is assumed that the daemon-scheduler is running locally.

```
$ cd <GitRepoBase>/daemon-scheduler/
$ AWS_REGION=<region> AWS_PROFILE=<profile> ECS_CLUSTER=<cluster> ECS_CLUSTER_ASG=<autoscale_group> gucumber -tags=@integ
...
$ echo $?
0  <- should return 0
```

#### Verify that doing a docker pull of all six tags works.
```
$ docker pull bloxoss/cluster-state-service:latest
$ echo $?
0  <- should return 0

$ docker pull bloxoss/cluster-state-service:$$blox_version$$
$ echo $?
0  <- should return 0

$ docker pull bloxoss/cluster-state-service:$$github_hash$$
$ echo $?
0  <- should return 0

$ docker pull bloxoss/daemon-scheduler:latest
$ echo $?
0  <- should return 0

$ docker pull bloxoss/daemon-scheduler:$$blox_version$$
$ echo $?
0  <- should return 0

$ docker pull bloxoss/daemon-scheduler:$$github_hash$$
$ echo $?
0  <- should return 0
```

#### Verify that doing a Local Deployment of the v$$blox_version$$ tag shows the correct versions
Required setup details are listed [here](https://github.com/blox/blox/blob/dev/deploy/README.md#local-installation). Stop all running Docker containers before proceeding.
```
$ cd <GitRepoBase>/deploy/docker/conf/
$ git checkout v$$blox_version$$
$ sed -i.bak 's/<region>/us-east-1/g' docker-compose.yml  # Replace us-east-1 with desired region
$ docker-compose up -d
$ docker ps
CONTAINER ID   IMAGE                                            STATUS
70d2ca6c5de7   bloxoss/daemon-scheduler:$$blox_version$$        Up      <- Should show 'Up' and version '$$blox_version$$' 
e2214884f981   bloxoss/cluster-state-service:$$blox_version$$   Up      <- Should show 'Up' and version '$$blox_version$$'
088f0d7c20e8   quay.io/coreos/etcd:v3.x.y                       Up      <- Should show 'Up'
```

#### Verify that doing a CloudFormation Deployment of the v$$blox_version$$ tag shows the correct versions
Required setup details are listed [here](https://github.com/blox/blox/blob/dev/deploy/README.md#aws-installation). Stop all running Docker containers before proceeding.
```
# Create /tmp/blox_parameters.json following the instructions on the README.md URL above.
$ cd <GitRepoBase>/deploy/aws/conf/
$ git checkout v$$blox_version$$
$ aws --region <region> cloudformation create-stack --stack-name BloxAws --template-body file://./cloudformation_template.json --capabilities CAPABILITY_NAMED_IAM --parameters file:///tmp/blox_parameters.json
# After CloudFormation Deployment completes, SSH into the EC2 instance created.
$ docker ps
CONTAINER ID   IMAGE                                            STATUS
70d2ca6c5de7   bloxoss/daemon-scheduler:$$blox_version$$        Up      <- Should show 'Up' and version '$$blox_version$$' 
e2214884f981   bloxoss/cluster-state-service:$$blox_version$$   Up      <- Should show 'Up' and version '$$blox_version$$'
088f0d7c20e8   quay.io/coreos/etcd:v3.x.y                       Up      <- Should show 'Up'
```

#### Verify that running the cluster-state-service end-to-end tests against the Local Deployment passes
Perform the `Ensure that the cluster-state-service end-to-end tests pass` steps above against the cluster-state-service started by the Local Deployment process.

#### Verify that running the daemon-scheduler end-to-end tests against the Local Deployment passes
Perform the `Ensure that the daemon-scheduler end-to-end tests pass` steps above against the daemon-scheduler started by the Local Deployment process.

#### Verify that running the daemon-scheduler integration tests against the Local Deployment passes
Perform the `Ensure that the daemon-scheduler integration tests pass` steps above against the daemon-scheduler started by the Local Deployment process.

#### Verify that running the cluster-state-service end-to-end tests against the CloudFormation Deployment passes
Perform the `Ensure that the cluster-state-service end-to-end tests pass` steps above against the cluster-state-service started by the CloudFormation Deployment process.

#### Verify that running the daemon-scheduler end-to-end tests against the CloudFormation Deployment passes
Perform the `Ensure that the daemon-scheduler end-to-end tests pass` steps above against the daemon-scheduler started by the CloudFormation Deployment process.

#### Verify that running the daemon-scheduler integration tests against the CloudFormation Deployment passes
Perform the `Ensure that the daemon-scheduler integration tests pass` steps above against the daemon-scheduler started by the CloudFormation Deployment process.

## Validation Checklist
Refer to the Validation Activities section for more details about each validation step.

- [ ] Verify that the Blox GitHub release < https://github.com/blox/blox/releases/tag/v$$blox_version$$ > looks correct and points to the correct revision. You should see the git hash '$$github_hash$$' on this page.
- [ ] Verify that doing a docker pull of all six tags works.
- [ ] Verify that doing a Local Deployment of the v$$blox_version$$ tag shows the correct versions.
- [ ] Verify that running the cluster-state-service end-to-end tests against the Local Deployment passes.
- [ ] Verify that running the daemon-scheduler end-to-end tests against the Local Deployment passes.
- [ ] Verify that running the daemon-scheduler integration tests against the Local Deployment passes.
- [ ] Verify that doing a CloudFormation Deployment of the v$$blox_version$$ tag shows the correct versions.
- [ ] Verify that running the cluster-state-service end-to-end tests against the CloudFormation Deployment passes.
- [ ] Verify that running the daemon-scheduler end-to-end tests against the CloudFormation Deployment passes.
- [ ] Verify that running the daemon-scheduler integration tests against the CloudFormation Deployment passes.