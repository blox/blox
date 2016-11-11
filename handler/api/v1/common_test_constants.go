// Copyright 2016 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the License). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the license file accompanying this file. This file is distributed
// on an AS IS BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

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
	pendingTaskCount1  = int64(0)
	runningTasksCount1 = int64(1)
	instanceStatus1    = "active"
	version1           = int64(1)
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
