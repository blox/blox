# Testing

## End to end tests
These tests call the cluster-state-service (CSS) APIs using the swagger generated client. They use ECS APIs to start tasks etc. and exercise the CSS consumer to consume events from the ECS event-stream and update the local state in etcd. The CSS client is then used to exercise the APIs supported by the CSS.

### What are the assumptions made?
* CSS server is running locally.
* An ECS cluster has already been created with at least one Container Instance registered to it. You can specify the cluster name by using the `ECS_CLUSTER` envrionment variable (This step can be automated in the future versions of the test suite by creating a test cluster and launching an EC2 instance with user data enabling to register itself to the created cluster).

### How to run the test suite?
From a level above the ./internal directory, use the following commands depending on test suite you want to run.

Note: The examples commands here make the following assumptions about your setup:

1. The ECS Cluster is named `test`
2. The AWS Credentials are saved under a profile named `test-profile`
3. The ECS Cluster is setup in the `us-east-1` region

If you're using a custom ECS Endpoint, then you can use `ECS_ENDPOINT=<endpoint>` to specify the same.

**All e2e tests**
```
AWS_REGION=us-east-1 AWS_PROFILE=test-profile ECS_CLUSTER=test gucumber -tags=@e2e
```

**All e2e instance API tests**
```
AWS_REGION=us-east-1 AWS_PROFILE=test-profile ECS_CLUSTER=test gucumber -tags=@instance,@e2e
```
**All e2e task API tests**
```
AWS_REGION=us-east-1 AWS_PROFILE=test-profile ECS_CLUSTER=test gcucumber -tags=@task,@e2e
```
**e2e GetInstance API tests**
```
AWS_REGION=us-east-1 AWS_PROFILE=test-profile ECS_CLUSTER=test gcucumber -tags=@get-instance,@e2e
```
**e2e ListInstances API tests**
```
AWS_REGION=us-east-1 AWS_PROFILE=test-profile ECS_CLUSTER=test gcucumber -tags=@list-instances,@e2e
```
**e2e FilterInstances API tests**
```
AWS_REGION=us-east-1 AWS_PROFILE=test-profile ECS_CLUSTER=test gcucumber -tags=@filter-instances,@e2e
```
**e2e GetTask API tests**
```
AWS_REGION=us-east-1 AWS_PROFILE=test-profile ECS_CLUSTER=test gcucumber -tags=@get-task,@e2e
```
**e2e ListTasks API tests**
```
AWS_REGION=us-east-1 AWS_PROFILE=test-profile ECS_CLUSTER=test gcucumber -tags=@list-tasks,@e2e
```
**e2e FilterTasks API tests**
```
AWS_REGION=us-east-1 AWS_PROFILE=test-profile ECS_CLUSTER=test gcucumber -tags=@filter-tasks,@e2e
```

***Note:*** The the CSS client used by the tests are checked in. However, if there is any change in 'swagger.json' file in the handler, re-generate the models using the following command from inside the ./internal directory.
```
swagger generate client -f ../handler/api/v1/swagger/swagger.json -A amazon_css
```
