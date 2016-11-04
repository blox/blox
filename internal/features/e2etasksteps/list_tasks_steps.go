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

package e2etasksteps

import (
	"time"

	"github.com/aws/amazon-ecs-event-stream-handler/internal/features/wrappers"
	. "github.com/gucumber/gucumber"
)

func init() {

	eshWrapper := wrappers.NewESHWrapper()

	When(`^I list tasks$`, func() {
		time.Sleep(15 * time.Second)
		eshTasks, err := eshWrapper.ListTasks()
		if err != nil {
			T.Errorf(err.Error())
		}
		for _, t := range eshTasks {
			eshTaskList = append(eshTaskList, *t)
		}
	})

	Then(`^the list tasks response contains at least (\d+) tasks$`, func(numTasks int) {
		if len(eshTaskList) < numTasks {
			T.Errorf("Number of tasks in list tasks response is less than expected")
		}
	})

	And(`^all (\d+) tasks are present in the list tasks response$`, func(numTasks int) {
		if len(ecsTaskList) != numTasks {
			T.Errorf("Error memorizing tasks started using ECS client")
		}
		for _, t := range ecsTaskList {
			err := ValidateListContainsTask(t, eshTaskList)
			if err != nil {
				T.Errorf(err.Error())
			}
		}
	})
}
