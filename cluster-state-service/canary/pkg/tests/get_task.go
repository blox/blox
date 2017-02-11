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
	"github.com/blox/blox/cluster-state-service/swagger/v1/generated/client/operations"
	log "github.com/cihub/seelog"
	"github.com/pkg/errors"
)

const (
	getTaskMetric            = "getTask"
	getTaskClusterNamePrefix = canaryClusterNamePrefix + "GetTask_"
)

// GetTask defines the test for CSS GetTask API
func (c canary) GetTask() error {
	clusterName := getTaskClusterNamePrefix + util.TimeString()

	err := c.ecsWrapper.CreateCluster(&clusterName)
	if err != nil {
		return err
	}
	log.Infof("Get task - created cluster")

	numInstances := int64(1)
	ec2InstanceIDs, err := c.ec2Wrapper.LaunchInstances(&numInstances, &clusterName)
	if err != nil {
		return err
	}
	log.Infof("Get task - launched EC2 instance")

	instanceARN, err := util.ValidateECSContainerInstanceAndGetInstanceARN(c.ecsWrapper, clusterName)
	if err != nil {
		return err
	}
	log.Infof("Get task - instance registered with ECS")

	// TODO: Start more than 1 task for this test
	taskARN, err := c.ecsWrapper.StartTask(&clusterName, &instanceARN)
	if err != nil {
		return err
	}
	log.Infof("Get task - task started on ECS")

	err = validateGetTaskCall(c.cssWrapper, clusterName, *taskARN)
	if err != nil {
		return err
	}
	log.Infof("Get task - task in CSS GetTask response")

	err = util.DeleteCluster(c.ecsWrapper, clusterName)
	if err != nil {
		return err
	}

	err = util.TerminateInstances(c.ec2Wrapper, ec2InstanceIDs)
	if err != nil {
		return err
	}
	log.Infof("Get task - GetTask test resources cleaned up")

	return nil
}

// GetTaskMetric returns the metric name to which the GetTask test metric should be emitted into
func (c canary) GetTaskMetric() string {
	return getTaskMetric
}

func validateGetTaskCall(cssWrapper wrappers.CSSWrapper, clusterName string, taskARN string) error {
	// Takes some time for CSS to get the task.
	// Retry get call once every 10 seconds for 2 minutes.
	// TODO: Change sleep and retry related numbers to constants
	found := false
	for i := 0; i < 12; i++ {
		task, err := cssWrapper.GetTask(&clusterName, &taskARN)
		if err != nil {
			rootErr := errors.Cause(err)
			if _, ok := rootErr.(*operations.GetTaskNotFound); ok {
				time.Sleep(10 * time.Second)
				continue
			}
			return err
		}
		if task == nil {
			return errors.Errorf("Get task call with cluster name '%s' and task ARN '%s' returned a nil instance. ",
				clusterName, taskARN)
		}
		if task.Entity == nil {
			return errors.Errorf("Get task call with cluster name '%s' and task ARN '%s' returned a nil instance entity. ",
				clusterName, taskARN)
		}
		if aws.StringValue(task.Entity.TaskARN) == taskARN {
			found = true
			break
		}
		return errors.Errorf("Expected task belonging to cluster with name '%s' to have ARN '%s' but was '%s'. ",
			clusterName, taskARN, *task.Entity.TaskARN)
	}

	if !found {
		return errors.Errorf("Task with ARN '%s' belonging to cluster with name"+
			" '%s' not found in CSS. ", taskARN, clusterName)
	}

	return nil
}
