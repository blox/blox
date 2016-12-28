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

package e2einstancesteps

import (
	"bufio"
	"encoding/json"

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
		r, err := cssWrapper.StreamInstances()
		if err != nil {
			T.Errorf(err.Error())
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

	Then(`^the stream instances response contains at least (\d+) instance$`, func(numInstances int) {
		<-stream
		if len(streamInstanceList) < numInstances {
			T.Errorf("Number of instances in stream instances response is less than expected. ")
		}
	})

	And(`^the stream instances response contains the cluster where the task was started$`, func() {
		clusterName, err := wrappers.GetClusterName()
		if err != nil {
			T.Errorf(err.Error())
		}

		err = ValidateListContainsCluster(clusterName, streamInstanceList)
		if err != nil {
			T.Errorf(err.Error())
		}
	})
}
