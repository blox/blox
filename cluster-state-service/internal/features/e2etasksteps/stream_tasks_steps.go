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
	"bufio"
	"encoding/json"

	"github.com/blox/blox/cluster-state-service/internal/features/wrappers"
	"github.com/blox/blox/cluster-state-service/swagger/v1/generated/models"
	. "github.com/gucumber/gucumber"
)

var (
	streamTaskList = []models.Task{}
)

func init() {

	cssWrapper := wrappers.NewCSSWrapper()
	stream := make(chan string)

	When(`^I start streaming all task events$`, func() {
		r, err := cssWrapper.StreamTasks()
		if err != nil {
			T.Errorf(err.Error())
		}

		go func() {
			scanner := bufio.NewScanner(r)
			for scanner.Scan() {
				task := &models.Task{}
				json.Unmarshal([]byte(scanner.Text()), task)
				streamTaskList = append(streamTaskList, *task)
			}
			stream <- "done"
		}()
	})

	Then(`^the stream tasks response contains at least (\d+) task$`, func(numTasks int) {
		if len(ecsTaskList) != 1 {
			T.Errorf("Error memorizing task started using ECS client. ")
		}

		<-stream
		if len(streamTaskList) < numTasks {
			T.Errorf("Number of tasks in stream tasks response is less than expected. ")
		}
	})

	And(`^the stream tasks response contains the task started$`, func() {
		err := ValidateListContainsTaskWithStatus(ecsTaskList[0], streamTaskList, "PENDING", "RUNNING")
		if err != nil {
			T.Errorf(err.Error())
		}
	})

}
