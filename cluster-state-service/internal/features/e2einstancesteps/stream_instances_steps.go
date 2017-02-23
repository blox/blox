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

package e2einstancesteps

import (
	"bufio"
	"encoding/json"

	"github.com/blox/blox/cluster-state-service/internal/features/e2etasksteps"
	"github.com/blox/blox/cluster-state-service/internal/features/wrappers"
	"github.com/blox/blox/cluster-state-service/swagger/v1/generated/models"
	. "github.com/gucumber/gucumber"
)

var (
	streamInstanceList = []models.ContainerInstance{}
)

func init() {
	cssWrapper := wrappers.NewCSSWrapper()
	stream := make(chan string)

	When(`^I start streaming all instance events$`, func() {
		streamInstanceList = nil

		r, err := cssWrapper.StreamInstances()
		if err != nil {
			T.Errorf(err.Error())
			return
		}

		go func() {
			scanner := bufio.NewScanner(r)
			for scanner.Scan() {
				instance := &models.ContainerInstance{}
				json.Unmarshal([]byte(scanner.Text()), instance)
				streamInstanceList = append(streamInstanceList, *instance)
			}
			stream <- "done"
		}()
	})

	When(`^I start streaming all instance events with past entity version$`, func() {
		streamInstanceList = nil

		if len(cssContainerInstanceList) != 1 {
			T.Errorf("Error memorizing instance retrieved using CSS client. ")
			return
		}

		r, err := cssWrapper.StreamInstancesWithEntityVersion(*cssContainerInstanceList[0].Metadata.EntityVersion)
		if err != nil {
			T.Errorf(err.Error())
			return
		}

		go func() {
			scanner := bufio.NewScanner(r)
			for scanner.Scan() {
				instance := &models.ContainerInstance{}
				json.Unmarshal([]byte(scanner.Text()), instance)
				streamInstanceList = append(streamInstanceList, *instance)
			}
			stream <- "done"
		}()
	})

	When(`^I get instance where the task was started$`, func() {
		cssContainerInstanceList = nil

		clusterName, err := wrappers.GetClusterName()
		if err != nil {
			T.Errorf(err.Error())
			return
		}

		if len(e2etasksteps.EcsTaskList) != 1 {
			T.Errorf("Error memorizing task retrieved using ECS client. ")
			return
		}

		cssInstance, err := cssWrapper.GetInstance(clusterName, *e2etasksteps.EcsTaskList[0].ContainerInstanceArn)
		if err != nil {
			T.Errorf(err.Error())
			return
		}
		cssContainerInstanceList = append(cssContainerInstanceList, *cssInstance)
	})

	Then(`^the stream instances response contains at least (\d+) instance$`, func(numInstances int) {
		<-stream
		if len(streamInstanceList) < numInstances {
			T.Errorf("Number of instances in stream instances response is less than expected. ")
		}
	})

	And(`^the stream instances response contains the instance where the task was started$`, func() {
		if len(cssContainerInstanceList) == 0 {
			T.Errorf("Error memorizing instances where the task was started. ")
			return
		}

		_, err := ValidateListContainsInstanceArn(*cssContainerInstanceList[0].Entity.ContainerInstanceARN, streamInstanceList)
		if err != nil {
			T.Errorf(err.Error())
		}
	})
}
