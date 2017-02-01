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
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/blox/blox/cluster-state-service/internal/features/wrappers"
	. "github.com/gucumber/gucumber"
)

const (
	invalidStatus  = "invalidStatus"
	invalidCluster = "cluster/cluster"

	badRequestHTTPResponse = "400 Bad Request"
	listTasksBadRequest    = "ListTasksBadRequest"
)

func init() {

	ecsWrapper := wrappers.NewECSWrapper()
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

	Then(`^the list tasks response contains at least (\d+) task(?:|s)$`, func(numTasks int) {
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

	When(`^I list tasks with cluster filter set to the ECS cluster name$`, func() {
		time.Sleep(15 * time.Second)
		clusterName, err := wrappers.GetClusterName()
		if err != nil {
			T.Errorf(err.Error())
		}
		cssTasks, err := cssWrapper.FilterTasksByCluster(clusterName)
		if err != nil {
			T.Errorf(err.Error())
		}
		for _, t := range cssTasks {
			cssTaskList = append(cssTaskList, *t)
		}
	})

	When(`^I list tasks with status filter set to (.+?)$`, func(status string) {
		time.Sleep(15 * time.Second)
		cssTasks, err := cssWrapper.FilterTasksByStatus(status)
		if err != nil {
			T.Errorf(err.Error())
		}
		for _, t := range cssTasks {
			cssTaskList = append(cssTaskList, *t)
		}
	})

	Given(`^I start (\d+) tasks in the ECS cluster with startedBy set to (.+?)$`, func(numTasks int, startedBy string) {
		startNTasks(numTasks, startedBy, ecsWrapper)
	})

	When(`^I list tasks with startedBy filter set to (.+?)$`, func(startedBy string) {
		time.Sleep(15 * time.Second)
		cssTasks, err := cssWrapper.FilterTasksByStartedBy(startedBy)
		if err != nil {
			T.Errorf(err.Error())
		}
		for _, t := range cssTasks {
			cssTaskList = append(cssTaskList, *t)
		}
	})

	And(`^I stop the (\d+) tasks in the ECS cluster$`, func(numTasks int) {
		clusterName, err := wrappers.GetClusterName()
		if err != nil {
			T.Errorf(err.Error())
		}
		if len(ecsTaskList) != numTasks {
			T.Errorf("Error memorizing tasks started using ECS client. ")
		}
		for _, t := range ecsTaskList {
			err := ecsWrapper.StopTask(clusterName, *t.TaskArn)
			if err != nil {
				T.Errorf(err.Error())
			}
		}
	})

	When(`^I list tasks with filters set to (.+?) status and cluster name$`, func(status string) {
		time.Sleep(15 * time.Second)
		clusterName, err := wrappers.GetClusterName()
		if err != nil {
			T.Errorf(err.Error())
		}
		cssTasks, err := cssWrapper.FilterTasksByStatusAndCluster(status, clusterName)
		if err != nil {
			T.Errorf(err.Error())
		}
		for _, t := range cssTasks {
			cssTaskList = append(cssTaskList, *t)
		}
	})

	And(`^all tasks in the list tasks response belong to the cluster and have status set to (.+?)$`, func(status string) {
		clusterName, err := wrappers.GetClusterName()
		if err != nil {
			T.Errorf(err.Error())
		}
		for _, t := range cssTaskList {
			if strings.ToLower(*t.LastStatus) != strings.ToLower(status) {
				T.Errorf("Task with ARN '%s' was expected to be '%s' but is '%s'", *t.TaskARN, status, *t.LastStatus)
			}
			if !strings.HasSuffix(*t.ClusterARN, "/"+clusterName) {
				T.Errorf("Task with ARN '%s' was expected to belong to cluster with name '%s' but belongs to cluster with ARN'%s'",
					*t.TaskARN, clusterName, *t.ClusterARN)
			}
		}
	})

	When(`^I list tasks with filters set to (.+?) status and a different cluster name$`, func(status string) {
		clusterName := "someCluster"
		cssTasks, err := cssWrapper.FilterTasksByStatusAndCluster(status, clusterName)
		if err != nil {
			T.Errorf(err.Error())
		}
		for _, t := range cssTasks {
			cssTaskList = append(cssTaskList, *t)
		}
	})

	Then(`^the list tasks response contains (\d+) tasks$`, func(numTasks int) {
		if len(cssTaskList) != numTasks {
			T.Errorf("Expected '%d' tasks in the list tasks response but got '%d'", numTasks, len(cssTaskList))
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

	When(`^I list tasks in the ECS cluster with status (.+?) and startedBy someone filters$`, func(status string) {
		time.Sleep(15 * time.Second)
		clusterName, err := wrappers.GetClusterName()
		if err != nil {
			T.Errorf(err.Error())
		}
		cssTasks, err := cssWrapper.ListTasksWithAllFilters(status,
			clusterName, "someone")
		if err != nil {
			T.Errorf(err.Error())
		}
		for _, t := range cssTasks {
			cssTaskList = append(cssTaskList, *t)
		}
	})

	And(`^all tasks in the list tasks response belong to the ECS cluster, are started by someone and have status set to (.+?)$`, func(status string) {
		clusterName, err := wrappers.GetClusterName()
		if err != nil {
			T.Errorf(err.Error())
		}
		for _, t := range cssTaskList {
			if t.StartedBy != "someone" {
				T.Errorf("Task with ARN '%s' was expected to be started by '%s' not '%s'", *t.TaskARN, "someone", t.StartedBy)
			}
			if strings.ToLower(*t.LastStatus) != strings.ToLower(status) {
				T.Errorf("Task with ARN '%s' was expected to be '%s' but is '%s'", *t.TaskARN, status, *t.LastStatus)
			}
			if !strings.HasSuffix(*t.ClusterARN, "/"+clusterName) {
				T.Errorf("Task with ARN '%s' was expected to belong to cluster with name '%s' but belongs to cluster with ARN'%s'",
					*t.TaskARN, clusterName, *t.ClusterARN)
			}
		}
	})

	When(`^I try to list tasks with redundant filters$`, func() {
		url := "http://localhost:3000/v1/tasks?cluster=cluster1&cluster=cluster2"
		resp, err := http.Get(url)
		if err != nil {
			T.Errorf(err.Error())
		}

		exceptionList = nil
		var exceptionType string
		if resp.Status == badRequestHTTPResponse {
			exceptionType = listTasksBadRequest
		} else {
			T.Errorf("Unknown exception type '%s' when trying to list tasks with redundant filters", resp.Status)
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			T.Errorf("Error reading expection message when trying to list tasks with redundant filters")
		}
		exceptionMsg := string(body)
		exceptionList = append(exceptionList, Exception{exceptionType: exceptionType, exceptionMsg: exceptionMsg})
	})
}
