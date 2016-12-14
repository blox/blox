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
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/blox/blox/cluster-state-service/internal/features/wrappers"
	"github.com/blox/blox/cluster-state-service/internal/models"
	. "github.com/gucumber/gucumber"
)

var (
	// Lists to memorize results required for the subsequent steps
	ecsContainerInstanceList = []ecs.ContainerInstance{}
	cssContainerInstanceList = []models.ContainerInstance{}
	exceptionList            = []string{}
)

func init() {

	ecsWrapper := wrappers.NewECSWrapper()

	Given(`^I have some instances registered with the ECS cluster$`, func() {
		ecsContainerInstanceList = nil
		cssContainerInstanceList = nil

		clusterName, err := wrappers.GetClusterName()
		if err != nil {
			T.Errorf(err.Error())
		}

		instanceARNs, err := ecsWrapper.ListContainerInstances(clusterName)
		if err != nil {
			T.Errorf(err.Error())
		}
		if len(instanceARNs) < 1 {
			T.Errorf("No container instances registered to the cluster '%s'. ", clusterName)
		}
		for _, instanceARN := range instanceARNs {
			ecsInstance, err := ecsWrapper.DescribeContainerInstance(clusterName, *instanceARN)
			if err != nil {
				T.Errorf(err.Error())
			}
			ecsContainerInstanceList = append(ecsContainerInstanceList, ecsInstance)
		}
	})

	Then(`^I get a (.+?) instance exception$`, func(exception string) {
		if len(exceptionList) != 1 {
			T.Errorf("Error memorizing exception")
		}
		if exception != exceptionList[0] {
			T.Errorf("Expected exception '%s' but got '%s'. ", exception, exceptionList[0])
		}
	})
}
