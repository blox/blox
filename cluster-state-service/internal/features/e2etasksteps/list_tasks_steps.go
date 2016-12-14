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

	"github.com/blox/blox/cluster-state-service/internal/features/wrappers"
	. "github.com/gucumber/gucumber"
)

const (
	invalidStatus  = "invalidStatus"
	invalidCluster = "cluster/cluster"
)

func init() {

	cssWrapper := wrappers.NewCSSWrapper()

	When(`^I list tasks$`, func() {
		time.Sleep(15 * time.Second)
		cssTasks, err := cssWrapper.ListTasks()
		if err != nil {
			T.Errorf(err.Error())
		}
		for _, t := range cssTasks {
			cssTaskList = append(cssTaskList, *t)
		}
	})

	Then(`^the list tasks response contains at least (\d+) tasks$`, func(numTasks int) {
		if len(cssTaskList) < numTasks {
			T.Errorf("Number of tasks in list tasks response is less than expected. ")
		}
	})

	And(`^all (\d+) tasks are present in the list tasks response$`, func(numTasks int) {
		if len(ecsTaskList) != numTasks {
			T.Errorf("Error memorizing tasks started using ECS client. ")
		}
		for _, t := range ecsTaskList {
			err := ValidateListContainsTask(t, cssTaskList)
			if err != nil {
				T.Errorf(err.Error())
			}
		}
	})

	When(`^I try to list tasks with an invalid status filter$`, func() {
		exceptionList = nil
		exceptionMsg, exceptionType, err := cssWrapper.TryListTasksWithInvalidStatus(invalidStatus)
		if err != nil {
			T.Errorf(err.Error())
		}
		exceptionList = append(exceptionList, Exception{exceptionType: exceptionType, exceptionMsg: exceptionMsg})
	})

	When(`^I try to list tasks with an invalid cluster filter$`, func() {
		exceptionList = nil
		exceptionMsg, exceptionType, err := cssWrapper.TryListTasksWithInvalidCluster(invalidCluster)
		if err != nil {
			T.Errorf(err.Error())
		}
		exceptionList = append(exceptionList, Exception{exceptionType: exceptionType, exceptionMsg: exceptionMsg})
	})
}
