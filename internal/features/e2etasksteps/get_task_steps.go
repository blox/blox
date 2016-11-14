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

package e2etasksteps

import (
	"time"

	"github.com/aws/amazon-ecs-event-stream-handler/internal/features/wrappers"
	. "github.com/gucumber/gucumber"
)

const (
	nonExistentTaskARN = "arn:aws:ecs:us-east-1:123456789012:task/31900037-daf4-40c7-93f7-102ece023cef"
)

func init() {

	cssWrapper := wrappers.NewCSSWrapper()

	When(`^I get task with the cluster name and task ARN$`, func() {
		clusterName, err := wrappers.GetClusterName()
		if err != nil {
			T.Errorf(err.Error())
		}

		time.Sleep(15 * time.Second)
		if len(ecsTaskList) != 1 {
			T.Errorf("Error memorizing task started using ECS client. ")
		}
		taskARN := *ecsTaskList[0].TaskArn
		cssTask, err := cssWrapper.GetTask(clusterName, taskARN)
		if err != nil {
			T.Errorf(err.Error())
		}
		cssTaskList = append(cssTaskList, *cssTask)
	})

	Then(`^I get a task that matches the task started$`, func() {
		if len(ecsTaskList) != 1 || len(cssTaskList) != 1 {
			T.Errorf("Error memorizing results to validate them. ")
		}
		ecsTask := ecsTaskList[0]
		cssTask := cssTaskList[0]
		err := ValidateTasksMatch(ecsTask, cssTask)
		if err != nil {
			T.Errorf(err.Error())
		}
	})

	When(`^I try to get task with a non-existent ARN$`, func() {
		exceptionList = nil
		exception, err := cssWrapper.TryGetTask(nonExistentTaskARN)
		if err != nil {
			T.Errorf(err.Error())
		}
		exceptionList = append(exceptionList, exception)
	})

}
