// Copyright 2016-2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.
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
	"github.com/blox/blox/cluster-state-service/internal/features/wrappers"
	. "github.com/gucumber/gucumber"
)

const (
	nonExistentTaskARN = "arn:aws:ecs:us-east-1:123456789012:task/31900037-daf4-40c7-93f7-102ece023cef"
)

func init() {

	cssWrapper := wrappers.NewCSSWrapper()

	Then(`^I get a task that matches the task started$`, func() {
		if len(EcsTaskList) != 1 || len(cssTaskList) != 1 {
			T.Errorf("Error memorizing results to validate them. ")
			return
		}
		ecsTask := EcsTaskList[0]
		cssTask := cssTaskList[0]
		err := ValidateTasksMatch(ecsTask, cssTask)
		if err != nil {
			T.Errorf(err.Error())
		}
	})

	When(`^I try to get task with a non-existent ARN$`, func() {
		exceptionList = nil
		exceptionMsg, exceptionType, err := cssWrapper.TryGetTask(nonExistentTaskARN)
		if err != nil {
			T.Errorf(err.Error())
			return
		}
		exceptionList = append(exceptionList, Exception{exceptionType: exceptionType, exceptionMsg: exceptionMsg})
	})

}
