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
	listInstancesMetric            = "listInstances"
	listInstancesClusterNamePrefix = canaryClusterNamePrefix + "ListInstances_"
)

// ListInstances defines the test for CSS ListInstances API
func (c canary) ListInstances() error {
	clusterName := listInstancesClusterNamePrefix + util.TimeString()

	err := c.ecsWrapper.CreateCluster(&clusterName)
	if err != nil {
		return err
	}
	log.Infof("List instances - created cluster")

	// TODO: Launch more than 1 instance for this test
	numInstances := int64(1)
	ec2InstanceIDs, err := c.ec2Wrapper.LaunchInstances(&numInstances, &clusterName)
	if err != nil {
		return err
	}
	log.Infof("List instances - launched EC2 instance")

	instanceARN, err := util.ValidateECSContainerInstanceAndGetInstanceARN(c.ecsWrapper, clusterName)
	if err != nil {
		return err
	}
	log.Infof("List instances - instance registered with ECS")

	err = validateListInstancesCall(c.cssWrapper, clusterName, instanceARN)
	if err != nil {
		return err
	}
	log.Infof("List instances - instance in CSS ListInstances response")

	err = util.DeleteCluster(c.ecsWrapper, clusterName)
	if err != nil {
		return err
	}

	err = util.TerminateInstances(c.ec2Wrapper, ec2InstanceIDs)
	if err != nil {
		return err
	}
	log.Infof("List instances - ListInstances test resources cleaned up")

	return nil
}

// ListInstancesMetric returns the metric name to which the ListInstances test metric should be emitted into
func (c canary) ListInstancesMetric() string {
	return listInstancesMetric
}

func validateListInstancesCall(cssWrapper wrappers.CSSWrapper, clusterName string, instanceARN string) error {
	// Takes some time for CSS to get the instance.
	// Retry list call once every 10 seconds for 2 minutes.
	// TODO: Change sleep and retry related numbers to constants
	found := false
	for i := 0; i < 12; i++ {
		instances, err := cssWrapper.ListInstances(&clusterName)
		if err != nil {
			return err
		}
		for _, ins := range instances {
			if ins.Entity != nil && aws.StringValue(ins.Entity.ContainerInstanceARN) == instanceARN {
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
		return errors.Errorf("Instance with ARN '%s' belonging to cluster with name"+
			" '%s' not found in CSS. ", instanceARN, clusterName)
	}

	return nil
}
