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
	"io/ioutil"
	"net/http"

	"github.com/blox/blox/cluster-state-service/internal/features/wrappers"
	. "github.com/gucumber/gucumber"
)

const (
	invalidStatus  = "invalidStatus"
	invalidCluster = "cluster/cluster"

	badRequestHTTPResponse  = "400 Bad Request"
	listInstancesBadRequest = "ListInstancesBadRequest"
)

func init() {

	cssWrapper := wrappers.NewCSSWrapper()

	When(`^I list instances$`, func() {
		cssInstances, err := cssWrapper.ListInstances()
		if err != nil {
			T.Errorf(err.Error())
			return
		}
		for _, i := range cssInstances {
			cssContainerInstanceList = append(cssContainerInstanceList, *i)
		}
	})

	Then(`^the list instances response contains all the registered instances$`, func() {
		// cssContainerInstanceList can have instances from other clusters too
		if len(cssContainerInstanceList) < len(ecsContainerInstanceList) {
			T.Errorf("Unexpected number of instances in the list instances response. ")
			return
		}
		for _, ecsInstance := range ecsContainerInstanceList {
			err := ValidateListContainsInstance(ecsInstance, cssContainerInstanceList)
			if err != nil {
				T.Errorf(err.Error())
				return
			}
		}
	})

	When(`^I list instances with cluster filter set to the ECS cluster name$`, func() {
		clusterName, err := wrappers.GetClusterName()
		if err != nil {
			T.Errorf(err.Error())
			return
		}

		cssInstances, err := cssWrapper.FilterInstancesByClusterName(clusterName)
		if err != nil {
			T.Errorf(err.Error())
			return
		}
		for _, i := range cssInstances {
			cssContainerInstanceList = append(cssContainerInstanceList, *i)
		}
	})

	Then(`^the list instances response contains all the instances registered with the cluster$`, func() {
		if len(ecsContainerInstanceList) != len(cssContainerInstanceList) {
			T.Errorf("Unexpected number of instances in the filter instances response. ")
			return
		}
		for _, ecsInstance := range ecsContainerInstanceList {
			err := ValidateListContainsInstance(ecsInstance, cssContainerInstanceList)
			if err != nil {
				T.Errorf(err.Error())
				return
			}
		}
	})

	When(`^I try to list instances with an invalid status filter$`, func() {
		exceptionList = nil
		exceptionMsg, exceptionType, err := cssWrapper.TryListInstancesWithInvalidStatus(invalidStatus)
		if err != nil {
			T.Errorf(err.Error())
			return
		}
		exceptionList = append(exceptionList, Exception{exceptionType: exceptionType, exceptionMsg: exceptionMsg})
	})

	When(`^I try to list instances with an invalid cluster filter$`, func() {
		exceptionList = nil
		exceptionMsg, exceptionType, err := cssWrapper.TryListInstancesWithInvalidCluster(invalidCluster)
		if err != nil {
			T.Errorf(err.Error())
			return
		}
		exceptionList = append(exceptionList, Exception{exceptionType: exceptionType, exceptionMsg: exceptionMsg})
	})

	When(`^I try to list instances with redundant filters$`, func() {
		url := "http://localhost:3000/v1/instances?cluster=cluster1&cluster=cluster2"
		resp, err := http.Get(url)
		if err != nil {
			T.Errorf(err.Error())
			return
		}

		exceptionList = nil
		var exceptionType string
		if resp.Status == badRequestHTTPResponse {
			exceptionType = listInstancesBadRequest
		} else {
			T.Errorf("Unknown exception type '%s' when trying to list instances with redundant filters", resp.Status)
			return
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			T.Errorf("Error reading expection message when trying to list instances with redundant filters")
			return
		}
		exceptionMsg := string(body)
		exceptionList = append(exceptionList, Exception{exceptionType: exceptionType, exceptionMsg: exceptionMsg})
	})
}
