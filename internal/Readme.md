# Testing

## End to end tests
These tests call the Event Stream Handler (ESH) APIs using the swagger generated client. They use ECS APIs to start tasks etc. and exercise the ESH consumer to consume events from the ECS event-stream and update the local state in etcd. The ESH client is then used to exercise the APIs supported by the ESH.

### What are the assumptions made?
* ESH server is running locally.
* There is an ECS cluster named "eventStreamTestCluster" and it has at least one container instance registered to it. (This step can be automated in the future versions of the test suite by creating a test cluster and launching an EC2 instance with user data enabling to register itself to the created cluster.)

### How to run the test suite?
From a level above the ./internal directory, use the following commands depending on test suite you want to run.

**All e2e tests**
```
gucumber -tags=@e2e
```
**All e2e instance API tests**
```
gucumber -tags=@instance,@e2e
```
**All e2e task API tests**
```
gcucumber -tags=@task,@e2e
```
**e2e GetInstance API tests**
```
gcucumber -tags=@get-instance,@e2e
```
**e2e ListInstances API tests**
```
gcucumber -tags=@list-instances,@e2e
```
**e2e FilterInstances API tests**
```
gcucumber -tags=@filter-instances,@e2e
```
**e2e GetTask API tests**
```
gcucumber -tags=@get-task,@e2e
```
**e2e ListTasks API tests**
```
gcucumber -tags=@list-tasks,@e2e
```
**e2e FilterTasks API tests**
```
gcucumber -tags=@filter-tasks,@e2e
```

***Note:*** The the ESH client used by the tests are checked in. However, if there is any change in 'swagger.json' file in the handler, re-generate the models using the following command from inside the ./internal directory.
```
swagger generate client -f ../handler/api/v1/swagger/swagger.json -A amazon_ecs_esh
```
