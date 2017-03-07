# Testing

## End to end tests
These tests call the cluster-state-service (CSS) APIs using the swagger generated client. They use ECS APIs to start tasks etc. and exercise the CSS consumer to consume events from the ECS event-stream and update the local state in etcd. The CSS client is then used to exercise the APIs supported by the CSS.

### What are the assumptions made?
* CSS server is running locally.
* The test will automatically set up the environment for you. By default, one EC2 instance is launched with the latest ECS-optimized AMI and without a key pair. The instance is set to auto-terminate itself after 1 hour in case the test fails to clean up. The instance is then registered to an ECS cluster named `E2ETestCluster`.

### How to run the test suite?
From a level above the ./internal directory, use the following commands depending on test suite you want to run.

## Customization
* If you're using a custom ECS Endpoint, you can use `ECS_ENDPOINT=<endpoint>` to specify the name.
* If you're using a custom ECS cluster, you can use `ECS_CLUSTER=<cluster>` to specify the name.
* If you want the EC2 instances to be launched with your key pair, you can use `EC2_KEY_PAIR=<key_pair>` to specify the key pair.

Note: The examples commands here make the following assumptions about your setup:

1. The AWS Credentials are saved under a profile named `test-profile`
2. The ECS Cluster is setup in the `us-east-1` region


**All e2e tests**
```
AWS_REGION=us-east-1 AWS_PROFILE=test-profile gucumber -tags=@e2e
```

**All e2e instance API tests**
```
AWS_REGION=us-east-1 AWS_PROFILE=test-profile gucumber -tags=@instance,@e2e
```
**All e2e task API tests**
```
AWS_REGION=us-east-1 AWS_PROFILE=test-profile gucumber -tags=@task,@e2e
```
**e2e GetInstance API tests**
```
AWS_REGION=us-east-1 AWS_PROFILE=test-profile gucumber -tags=@get-instance,@e2e
```
**e2e ListInstances API tests**
```
AWS_REGION=us-east-1 AWS_PROFILE=test-profile gucumber -tags=@list-instances,@e2e
```
**e2e FilterInstances API tests**
```
AWS_REGION=us-east-1 AWS_PROFILE=test-profile gucumber -tags=@filter-instances,@e2e
```
**e2e GetTask API tests**
```
AWS_REGION=us-east-1 AWS_PROFILE=test-profile gucumber -tags=@get-task,@e2e
```
**e2e ListTasks API tests**
```
AWS_REGION=us-east-1 AWS_PROFILE=test-profile gucumber -tags=@list-tasks,@e2e
```
**e2e FilterTasks API tests**
```
AWS_REGION=us-east-1 AWS_PROFILE=test-profile gucumber -tags=@filter-tasks,@e2e
```
