## Subject Line
[Blox] Deploy release - v$$blox_version$$

## Activity Details
This activity is to: Release Blox version v$$blox_version$$ by merging the release-$$blox_version$$ branch into the dev and master branches on GitHub, and pushing new Blox images up to Docker Hub.

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

Refer to the Validation Activities section for more details about each activity and validation checklist step.

## Activity Checklist
- [ ] Prereqs to be performed prior to the release starting:
- [ ] Create a release-$$blox_version$$ branch off of the dev branch in your Github fork.
- [ ] Draft a v$$blox_version$$ release in GitHub against the master branch.
- [ ] Ensure that the cluster-state-service unit-tests pass.
- [ ] Ensure that the cluster-state-service end-to-end tests pass.
- [ ] Ensure that the daemon-scheduler unit-tests pass.
- [ ] Ensure that the daemon-scheduler end-to-end tests pass.
- [ ] Ensure that the daemon-scheduler integration tests pass.
- [ ] Ensure that `deploy/docker/conf/docker-compose.yml` has been updated with the correct image versions in the release-$$blox_version$$ branch.
- [ ] Ensure that `deploy/aws/conf/cloudformation_template.json` has been updated with the correct image versions in the release-$$blox_version$$ branch.
- [ ] Ensure that the `IamRoleTask` in `deploy/aws/conf/cloudformation_template.json` has the correct role permissions in the release-$$blox_version$$ branch.
- [ ] Ensure that a $$blox_version$$ change log entry exists in `CHANGELOG.md` in the release-$$blox_version$$ branch.
- [ ] Ensure that the version is set to $$blox_version$$ in `cluster-state-service/README.md` and `daemon-scheduler/README.md` in the release-$$blox_version$$ branch.
- [ ] Ensure that the version is set to $$blox_version$$ in `cluster-state-service/VERSION` and `daemon-scheduler/VERSION` in the release-$$blox_version$$ branch.
- [ ] Ensure that the version is set to $$blox_version$$ in `cluster-state-service/versioning/version.go` and `daemon-scheduler/versioning/version.go` in the release-$$blox_version$$ branch.
- [ ] Ensure that the release-$$blox_version$$ branch contains the latest changes in the dev and master branch.
- [ ] Ensure that the release-$$blox_version$$ branch commit hash is set to $$github_hash$$.
- [ ] Ensure that the technician has a GitHub account and permissions to merge into the Blox dev and master branches.
- [ ] Ensure that the technician has a Docker Hub account and permissions to push and remove images in the bloxoss repositories.
- [ ] Ensure that the technician has logged into Docker Hub via `docker logout` and `docker login`.
- [ ] Ensure that the technician has logged into ECR via `aws --region $$aws_region$$ ecr get-login`.
- [ ] Ensure that the $$ecr_css_repo_uri$$ ECR repository has been created.
- [ ] Ensure that the $$ecr_ds_repo_uri$$ ECR repository has been created.
- [ ] Delete all local Blox containers and images before building new images.
- [ ] Build the $$ecr_css_repo_uri$$:$$blox_version$$ image and publish to ECR.
- [ ] Build the $$ecr_ds_repo_uri$$:$$blox_version$$ image and publish to ECR.
- [ ] Ensure that the $$ecr_css_repo_uri$$:$$blox_version$$ image shows the correct version and commit hash.
- [ ] Ensure that the $$ecr_ds_repo_uri$$:$$blox_version$$ image shows the correct version and commit hash.
- [ ] Ensure that the cluster-state-service end-to-end tests pass against the $$ecr_css_repo_uri$$:$$blox_version$$ image.
- [ ] Ensure that the daemon-scheduler end-to-end tests pass against the $$ecr_ds_repo_uri$$:$$blox_version$$ image.
- [ ] Ensure that the daemon-scheduler integration tests pass against the $$ecr_ds_repo_uri$$:$$blox_version$$ image.
- [ ] Create a pull request from the release-$$blox_version$$ branch to the master branch in GitHub.
- [ ] Ensure that the pull request from the release-$$blox_version$$ branch to the master branch in GitHub is approved.
- [ ] Release activities:
- [ ] Push the release-$$blox_version$$ branch up to the Github dev and master branches.
- [ ] Publish the v$$blox_version$$ release in GitHub.
- [ ] Publish the bloxoss/cluster-state-service:latest,$$blox_version$$,$$github_hash$$ images to Docker Hub.
- [ ] Publish the bloxoss/daemon-scheduler:latest,$$blox_version$$,$$github_hash$$ images to Docker Hub.
- [ ] Log into the GitHub console and close the 'Release $$blox_version$$' milestone.

## Rollback Procedure
- Log into the Docker Hub console and delete the bloxoss/cluster-state-service:latest,$$blox_version$$,$$github_hash$$ images.
- Log into the Docker Hub console and delete the bloxoss/daemon-scheduler:latest,$$blox_version$$,$$github_hash$$ images.
- Publish the bloxoss/cluster-state-service:latest tag pointed to the previous release.
- Publish the bloxoss/daemon-scheduler:latest tag pointed to the previous release.
- Revert Blox repository changes.
- Log into the GitHub console and reopen the 'Release $$blox_version$$' milestone.

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

#### Draft a v$$blox_version$$ release in GitHub against the master branch
```
Open up a new browser window to [here](https://github.com/blox/blox/releases/new)
Enter the following details:
 Tag version: v$$blox_version$$
 Target: master
 Release title: Release v$$blox_version$$
 Description: Release v$$blox_version$$
 This is a pre-release: Yes
Click: Save draft
```

#### Ensure that the cluster-state-service unit-tests pass
```
$ git checkout release-$$blox_version$$
$ cd <GitRepoBase>/cluster-state-service/
$ make unit-tests
...
$ echo $?
0  <- should return 0
```

#### Ensure that the cluster-state-service end-to-end tests pass
Required setup details are listed [here](https://github.com/blox/blox/blob/dev/cluster-state-service/internal/Readme.md). It is assumed that the cluster-state-service is running locally.

```
$ git checkout release-$$blox_version$$
$ cd <GitRepoBase>/cluster-state-service/
$ AWS_REGION=$$aws_region$$ AWS_PROFILE=<profile> ECS_CLUSTER=<cluster> gucumber -tags=@e2e
...
$ echo $?
0  <- should return 0
```

#### Ensure that the daemon-scheduler unit-tests pass
```
$ git checkout release-$$blox_version$$
$ cd <GitRepoBase>/daemon-scheduler/
$ make unit-tests
...
$ echo $?
0  <- should return 0
```

#### Ensure that the daemon-scheduler end-to-end tests pass
Required setup details are listed [here](https://github.com/blox/blox/blob/dev/daemon-scheduler/internal/features/README.md). It is assumed that the daemon-scheduler is running locally.

```
$ git checkout release-$$blox_version$$
$ cd <GitRepoBase>/daemon-scheduler/
$ AWS_REGION=$$aws_region$$ AWS_PROFILE=<profile> ECS_CLUSTER=<cluster> ECS_CLUSTER_ASG=<autoscale_group> gucumber -tags=@e2e
...
$ echo $?
0  <- should return 0
```

#### Ensure that the daemon-scheduler integration tests pass
Required setup details are listed [here](https://github.com/blox/blox/blob/dev/daemon-scheduler/internal/features/README.md). It is assumed that the daemon-scheduler is running locally.

```
$ git checkout release-$$blox_version$$
$ cd <GitRepoBase>/daemon-scheduler/
$ AWS_REGION=$$aws_region$$ AWS_PROFILE=<profile> ECS_CLUSTER=<cluster> ECS_CLUSTER_ASG=<autoscale_group> gucumber -tags=@integ
...
$ echo $?
0  <- should return 0
```

#### Ensure that the `IamRoleTask` in `deploy/aws/conf/cloudformation_template.json` has the correct role permissions in the release-$$blox_version$$ branch
- Open up `deploy/aws/conf/cloudformation_template.json` in a text editor and scroll down to the IamRoleTask role definition.
- Verify that all AWS API methods listed in IamRoleTask.Properties.Policies[0].PolicyDocument are still being called by the Blox source code.
- Verify that IamRoleTask.Properties.Policies[0].PolicyDocument is not missing any new AWS API methods that are being called by the Blox source code.

#### Ensure that the release-$$blox_version$$ branch contains the latest changes in the dev and master branch
```
$ git checkout release-$$blox_version$$
$ git branch --contains origin/dev | grep '^\* release-$$blox_version$$' || echo "ERROR: Current branch does not contain the latest changes in the 'dev' branch."
$ git branch --contains origin/master | grep '^\* release-$$blox_version$$' || echo "ERROR: Current branch does not contain the latest changes in the 'master' branch."
```

#### Ensure that the release-$$blox_version$$ branch commit hash is set to $$github_hash$$
```
$ git checkout release-$$blox_version$$
$ git rev-parse HEAD | grep '$$github_hash$$' || echo "Error: Commit hash should be '$$github_hash$$'."
```

#### Delete all local Blox containers and images before building new images
```
$ docker ps -a --no-trunc | grep 'cluster-state-service' | awk '{print $1}' | xargs docker rm -f
$ echo $?
0  <- should return 0

$ docker ps -a --no-trunc | grep 'daemon-scheduler' | awk '{print $1}' | xargs docker rm -f
$ echo $?
0  <- should return 0

$ docker images -a | grep 'bloxoss' | awk '{print $3}' | xargs docker rmi -f
$ echo $?
0  <- should return 0

$ docker images -a | grep '$$ecr_css_repo_uri$$' | awk '{print $3}' | xargs docker rmi -f
$ echo $?
0  <- should return 0

$ docker images -a | grep '$$ecr_ds_repo_uri$$' | awk '{print $3}' | xargs docker rmi -f
$ echo $?
0  <- should return 0
```

#### Build the $$ecr_css_repo_uri$$:$$blox_version$$ image and publish to ECR
```
$ git checkout release-$$blox_version$$
$ cd <GitRepoBase>/cluster-state-service/
$ make release
$ docker tag bloxoss/cluster-state-service:latest $$ecr_css_repo_uri$$:$$blox_version$$
$ docker push $$ecr_css_repo_uri$$:$$blox_version$$
```

#### Build the $$ecr_ds_repo_uri$$:$$blox_version$$ image and publish to ECR
```
$ git checkout release-$$blox_version$$
$ cd <GitRepoBase>/daemon-scheduler/
$ make release
$ docker tag bloxoss/daemon-scheduler:latest $$ecr_ds_repo_uri$$:$$blox_version$$
$ docker push $$ecr_ds_repo_uri$$:$$blox_version$$
```

#### Ensure that the $$ecr_css_repo_uri$$:$$blox_version$$ image shows the correct version and commit hash
```
$ docker run $$ecr_css_repo_uri$$:$$blox_version$$ --version
Blox Cluster State Service:
  Version: $$blox_version$$  <- Should show version '$$blox_version$$'
  Commit: $$github_hash$$  <- Should show commit '$$github_hash$$'
```

#### Ensure that the $$ecr_ds_repo_uri$$:$$blox_version$$ image shows the correct version and commit hash
```
$ docker run $$ecr_ds_repo_uri$$:$$blox_version$$ --version
Blox Daemon Scheduler:
  Version: $$blox_version$$  <- Should show version '$$blox_version$$'
  Commit: $$github_hash$$  <- Should show commit '$$github_hash$$'
```

#### Start the ECR images for e2e and integration testing
Required setup details are listed [here](https://github.com/blox/blox/blob/dev/deploy/README.md#local-installation). Stop all running Docker containers before proceeding.
```
$ docker ps | awk '{print $1}' | grep -v CONTAINER | xargs docker stop
$ git checkout release-$$blox_version$$
$ cd <GitRepoBase>/deploy/docker/conf/
$ sed -i '' 's#bloxoss/cluster-state-service:$$blox_version$$#$$ecr_css_repo_uri$$:$$blox_version$$#g' docker-compose.yml
$ sed -i '' 's#bloxoss/daemon-scheduler:$$blox_version$$#$$ecr_ds_repo_uri$$:$$blox_version$$#g' docker-compose.yml
$ sed -i '' 's/<region>/$$aws_region$$/g' docker-compose.yml
$ docker-compose up -d
$ docker ps
CONTAINER ID   IMAGE                                   STATUS
70d2ca6c5de7   $$ecr_ds_repo_uri$$:$$blox_version$$    Up      <- Should show 'Up' and version '$$blox_version$$' 
e2214884f981   $$ecr_css_repo_uri$$:$$blox_version$$   Up      <- Should show 'Up' and version '$$blox_version$$'
088f0d7c20e8   quay.io/coreos/etcd:v3.x.y              Up      <- Should show 'Up'
$ git checkout docker-compose.yml  # Revert local changes to file
```

#### Ensure that the cluster-state-service end-to-end tests pass against the $$ecr_css_repo_uri$$:$$blox_version$$ image
- Perform the `Start the ECR images for e2e and integration testing` steps above
- Perform the `Ensure that the cluster-state-service end-to-end tests pass` steps above

#### Ensure that the daemon-scheduler end-to-end tests pass against the $$ecr_ds_repo_uri$$:$$blox_version$$ image
- Perform the `Start the ECR images for e2e and integration testing` steps above
- Perform the `Ensure that the daemon-scheduler end-to-end tests pass` steps above

#### Ensure that the daemon-scheduler integration tests pass against the $$ecr_ds_repo_uri$$:$$blox_version$$ image
- Perform the `Start the ECR images for e2e and integration testing` steps above
- Perform the `Ensure that the daemon-scheduler integration tests pass` steps above

#### Push the release-$$blox_version$$ branch up to the Github dev and master branches
```
$ git checkout release-$$blox_version$$
$ git remote add upstream git@github.com:blox/blox.git
$ git push -v upstream 'release-$$blox_version$$:dev'
$ git push -v upstream 'release-$$blox_version$$:master'
```

#### Publish the v$$blox_version$$ release in GitHub
```
Open up a new browser window to [here](https://github.com/blox/blox/releases)
Click on the Edit button next to 'Release v$$blox_version$$'
Ensure the following details:
 Tag version: v$$blox_version$$
 Target: master
 Release title: Release v$$blox_version$$
 Description: Release v$$blox_version$$
 This is a pre-release: Yes
Click: Publish release
```

#### Publish the bloxoss/cluster-state-service:latest,$$blox_version$$,$$github_hash$$ images to Docker Hub
```
$ docker pull $$ecr_css_repo_uri$$:$$blox_version$$
$ docker tag $$ecr_css_repo_uri$$:$$blox_version$$ bloxoss/cluster-state-service:latest
$ docker tag $$ecr_css_repo_uri$$:$$blox_version$$ bloxoss/cluster-state-service:$$blox_version$$
$ docker tag $$ecr_css_repo_uri$$:$$blox_version$$ bloxoss/cluster-state-service:$$github_hash$$
$ docker push bloxoss/cluster-state-service:latest
$ docker push bloxoss/cluster-state-service:$$blox_version$$
$ docker push bloxoss/cluster-state-service:$$github_hash$$
```

#### Publish the bloxoss/daemon-scheduler:latest,$$blox_version$$,$$github_hash$$ images to Docker Hub
```
$ docker pull $$ecr_ds_repo_uri$$:$$blox_version$$
$ docker tag $$ecr_ds_repo_uri$$:$$blox_version$$ bloxoss/daemon-scheduler:latest
$ docker tag $$ecr_ds_repo_uri$$:$$blox_version$$ bloxoss/daemon-scheduler:$$blox_version$$
$ docker tag $$ecr_ds_repo_uri$$:$$blox_version$$ bloxoss/daemon-scheduler:$$github_hash$$
$ docker push bloxoss/daemon-scheduler:latest
$ docker push bloxoss/daemon-scheduler:$$blox_version$$
$ docker push bloxoss/daemon-scheduler:$$github_hash$$
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

#### Verify that doing a docker run of all six tags shows the correct version and commit hash
```
$ docker run bloxoss/cluster-state-service:latest --version
Blox Cluster State Service:
  Version: $$blox_version$$  <- Should show version '$$blox_version$$'
  Commit: $$github_hash$$  <- Should show commit '$$github_hash$$'

$ docker run bloxoss/cluster-state-service:$$blox_version$$ --version
Blox Cluster State Service:
  Version: $$blox_version$$  <- Should show version '$$blox_version$$'
  Commit: $$github_hash$$  <- Should show commit '$$github_hash$$'

$ docker run bloxoss/cluster-state-service:$$github_hash$$ --version
Blox Cluster State Service:
  Version: $$blox_version$$  <- Should show version '$$blox_version$$'
  Commit: $$github_hash$$  <- Should show commit '$$github_hash$$'

$ docker run bloxoss/daemon-scheduler:latest --version
Blox Daemon Scheduler:
  Version: $$blox_version$$  <- Should show version '$$blox_version$$'
  Commit: $$github_hash$$  <- Should show commit '$$github_hash$$'

$ docker run bloxoss/daemon-scheduler:$$blox_version$$ --version
Blox Daemon Scheduler:
  Version: $$blox_version$$  <- Should show version '$$blox_version$$'
  Commit: $$github_hash$$  <- Should show commit '$$github_hash$$'

$ docker run bloxoss/daemon-scheduler:$$github_hash$$ --version
Blox Daemon Scheduler:
  Version: $$blox_version$$  <- Should show version '$$blox_version$$'
  Commit: $$github_hash$$  <- Should show commit '$$github_hash$$'
```

#### Verify that doing a Local Deployment of the v$$blox_version$$ tag shows the correct versions
Required setup details are listed [here](https://github.com/blox/blox/blob/dev/deploy/README.md#local-installation). Stop all running Docker containers before proceeding.
```
$ docker ps | awk '{print $1}' | grep -v CONTAINER | xargs docker stop
$ git clone https://github.com/blox/blox.git
$ cd ./blox/deploy/docker/conf/
$ git checkout v$$blox_version$$
$ sed -i '' 's/<region>/$$aws_region$$/g' docker-compose.yml
$ docker-compose up -d
$ docker ps
CONTAINER ID   IMAGE                                            STATUS
70d2ca6c5de7   bloxoss/daemon-scheduler:$$blox_version$$        Up      <- Should show 'Up' and version '$$blox_version$$' 
e2214884f981   bloxoss/cluster-state-service:$$blox_version$$   Up      <- Should show 'Up' and version '$$blox_version$$'
088f0d7c20e8   quay.io/coreos/etcd:v3.x.y                       Up      <- Should show 'Up'
```

#### Verify that doing a CloudFormation Deployment of the v$$blox_version$$ tag shows the correct versions
Required setup details are listed [here](https://github.com/blox/blox/blob/dev/deploy/README.md#aws-installation).
```
# Create /tmp/blox_parameters.json following the instructions on the README.md URL above.
$ git clone https://github.com/blox/blox.git
$ cd ./blox/deploy/aws/conf/
$ git checkout v$$blox_version$$
$ aws --region $$aws_region$$ cloudformation create-stack --stack-name BloxAws --template-body file://./cloudformation_template.json --capabilities CAPABILITY_NAMED_IAM --parameters file:///tmp/blox_parameters.json
# After CloudFormation Deployment completes, SSH into the EC2 instance created.
$ docker ps
CONTAINER ID   IMAGE                                            STATUS
70d2ca6c5de7   bloxoss/daemon-scheduler:$$blox_version$$        Up      <- Should show 'Up' and version '$$blox_version$$' 
e2214884f981   bloxoss/cluster-state-service:$$blox_version$$   Up      <- Should show 'Up' and version '$$blox_version$$'
088f0d7c20e8   quay.io/coreos/etcd:v3.x.y                       Up      <- Should show 'Up'
```

## Validation Checklist
- [ ] Verify that the Blox GitHub release < https://github.com/blox/blox/releases/tag/v$$blox_version$$ > looks correct and points to the correct revision. You should see the git hash '$$github_hash$$' on this page.
- [ ] Verify that the pull request from the release-$$blox_version$$ branch to the master branch in GitHub is closed.
- [ ] Verify that doing a docker pull of all six tags works.
- [ ] Verify that doing a docker run of all six tags shows the correct version and commit hash.
- [ ] Verify that doing a Local Deployment of the v$$blox_version$$ tag shows the correct versions.
- [ ] Verify that doing a CloudFormation Deployment of the v$$blox_version$$ tag shows the correct versions.
