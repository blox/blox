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

package tests

import (
	"time"

  "github.com/aws/aws-sdk-go/aws"
	"github.com/blox/blox/cluster-state-service/canary/pkg/tests/util"
	"github.com/blox/blox/cluster-state-service/canary/pkg/wrappers"
	log "github.com/cihub/seelog"
	"github.com/pkg/errors"
)

const (
	listTasksMetric            = "listTasks"
	listTasksClusterNamePrefix = canaryClusterNamePrefix + "ListTasks_"
)

// ListTasks defines the test for CSS ListTasks API
func (c canary) ListTasks() error {
	clusterName := listTasksClusterNamePrefix + util.TimeString()

	err := c.ecsWrapper.CreateCluster(&clusterName)
	if err != nil {
		return err
	}
	log.Infof("List tasks - created cluster")

	numInstances := int64(1)
	ec2InstanceIDs, err := c.ec2Wrapper.LaunchInstances(&numInstances, &clusterName)
	if err != nil {
		return err
	}
	log.Infof("List tasks - launched EC2 instance")

	instanceARN, err := util.ValidateECSContainerInstanceAndGetInstanceARN(c.ecsWrapper, clusterName)
	if err != nil {
		return err
	}
	log.Infof("List tasks - instance registered with ECS")

	// TODO: Start more than 1 task for this test
	taskARN, err := c.ecsWrapper.StartTask(&clusterName, &instanceARN)
	if err != nil {
		return err
	}
	log.Infof("List tasks - task started on ECS")

	err = validateListTasksCall(c.cssWrapper, clusterName, *taskARN)
	if err != nil {
		return err
	}
	log.Infof("List tasks - instance in CSS ListTasks response")

	err = util.DeleteCluster(c.ecsWrapper, clusterName)
	if err != nil {
		return err
	}

	err = util.TerminateInstances(c.ec2Wrapper, ec2InstanceIDs)
	if err != nil {
		return err
	}
	log.Infof("List tasks - ListTasks test resources cleaned up")

	return nil
}

// ListTasksMetric returns the metric name to which the ListTasks test metric should be emitted into
func (c canary) ListTasksMetric() string {
	return listTasksMetric
}

func validateListTasksCall(cssWrapper wrappers.CSSWrapper, clusterName string, taskARN string) error {
	// Takes some time for CSS to get the task.
	// Retry list call once every 10 seconds for 2 minutes.
	// TODO: Change sleep and retry related numbers to constants
	found := false
	for i := 0; i < 12; i++ {
		tasks, err := cssWrapper.ListTasks(&clusterName)
		if err != nil {
			return err
		}
		for _, task := range tasks {
			if task.Entity != nil && aws.StringValue(task.Entity.TaskARN) == taskARN {
				found = true
				break
			}
		}
		if found {
			break
		}
		time.Sleep(10 * time.Second)
	}

	if !found {
		return errors.Errorf("Task with ARN '%s' belonging to cluster with name"+
			" '%s' not found in CSS. ", taskARN, clusterName)
	}

	return nil
}
