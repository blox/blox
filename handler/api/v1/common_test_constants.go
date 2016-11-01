package v1

var (
	accountID          = "123456789012"
	region             = "us-east-1"
	time               = "2016-10-18T16:52:49Z"
	id1                = "4082c1f7-d572-4684-8b3b-a7dd637e8721"
	instanceARN1       = "arn:aws:ecs:us-east-1:123456789012:container-instance/b6b9eace-958e-4f2a-a09c-8cf43b76cf97"
	clusterName1       = "cluster1"
	clusterARN1        = "arn:aws:ecs:us-east-1:123456789012:cluster/" + clusterName1
	containerARN1      = "arn:aws:ecs:us-east-1:123456789012:container/57156e30-e410-4773-9a9e-ae8264c10bbd"
	agentConnected1    = true
	pendingTaskCount1  = 0
	runningTasksCount1 = 1
	instanceStatus1    = "active"
	version1           = 0
	updatedAt1         = "2016-10-20T18:53:29.005Z"
	taskStatus1        = "pending"
	taskStatus2        = "stopped"
	taskName           = "testTask"
	createdAt          = "2016-10-24T06:07:53.036Z"
	taskARN1           = "arn:aws:ecs:us-east-1:123456789012:task/271022c0-f894-4aa2-b063-25bae55088d5"
	taskARN2           = "arn:aws:ecs:us-east-1:123456789012:task/b6b9eace-958e-4f2a-a09c-8cf43b76cf97"
	taskDefinitionARN  = "arn:aws:ecs:us-east-1:123456789012:task-definition/testTask:1"
)

const (
	responseContentTypeKey = "Content-Type"
	responseContentTypeVal = "application/json; charset=UTF-8"
)
