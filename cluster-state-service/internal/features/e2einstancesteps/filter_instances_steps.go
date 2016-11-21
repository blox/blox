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
	"github.com/blox/blox/cluster-state-service/internal/features/wrappers"
	. "github.com/gucumber/gucumber"
)

func init() {
	cssWrapper := wrappers.NewCSSWrapper()

	When(`^I filter instances by the same ECS cluster name$`, func() {
		clusterName, err := wrappers.GetClusterName()
		if err != nil {
			T.Errorf(err.Error())
		}

		cssInstances, err := cssWrapper.FilterInstancesByClusterName(clusterName)
		if err != nil {
			T.Errorf(err.Error())
		}
		for _, i := range cssInstances {
			cssContainerInstanceList = append(cssContainerInstanceList, *i)
		}
	})

	Then(`^the filter instances response contains all the instances registered with the cluster$`, func() {
		if len(ecsContainerInstanceList) != len(cssContainerInstanceList) {
			T.Errorf("Unexpected number of instances in the filter instances response. ")
		}
		for _, ecsInstance := range ecsContainerInstanceList {
			err := ValidateListContainsInstance(ecsInstance, cssContainerInstanceList)
			if err != nil {
				T.Errorf(err.Error())
			}
		}
	})

}
