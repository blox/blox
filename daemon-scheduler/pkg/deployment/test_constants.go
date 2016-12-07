// Copyright 2016 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package deployment

const (
	environmentName  = "environmentName"
	environmentName1 = "environmentName1"
	environmentName2 = "environmentName2"
	environmentName3 = "environmentName3"
	taskDefinition   = "arn:aws:ecs:us-east-1:12345678912:task-definition/test"
	taskARN1         = "arn:aws:ecs:us-east-1:12345678912:task/c024d145-093b-499a-9b14-5baf273f5835"
	taskARN2         = "arn:aws:ecs:us-east-1:12345678912:task/a1d71628-01e3-4013-b18c-6e14032a9522"
	cluster1         = "arn:aws:ecs:us-east-1:123456789123:cluster/test1"
	cluster2         = "arn:aws:ecs:us-east-1:123456789123:cluster/test2"
	instanceARN      = "arn:aws:us-east-1:123456789123:container-instance/4b6d45ea-a4b4-4269-9d04-3af6ddfdc597"
	desiredTaskCount = 5
)
