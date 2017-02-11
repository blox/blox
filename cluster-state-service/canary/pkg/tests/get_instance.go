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
	getInstanceMetric            = "getInstance"
	getInstanceClusterNamePrefix = canaryClusterNamePrefix + "GetInstance_"
)

// GetInstance defines the test for CSS GetInstance API
func (c canary) GetInstance() error {
	clusterName := getInstanceClusterNamePrefix + util.TimeString()

	err := c.ecsWrapper.CreateCluster(&clusterName)
	if err != nil {
		return err
	}
	log.Infof("Get instance - created cluster")

	// TODO: Launch more than 1 instance for this test
	numInstances := int64(1)
	ec2InstanceIDs, err := c.ec2Wrapper.LaunchInstances(&numInstances, &clusterName)
	if err != nil {
		return err
	}
	log.Infof("Get instance - launched EC2 instance")

	instanceARN, err := util.ValidateECSContainerInstanceAndGetInstanceARN(c.ecsWrapper, clusterName)
	if err != nil {
		return err
	}
	log.Infof("Get instance - instance registered with ECS")

	err = validateGetInstanceCall(c.cssWrapper, clusterName, instanceARN)
	if err != nil {
		return err
	}
	log.Infof("Get instance - instance in CSS GetInstance response")

	err = util.DeleteCluster(c.ecsWrapper, clusterName)
	if err != nil {
		return err
	}

	err = util.TerminateInstances(c.ec2Wrapper, ec2InstanceIDs)
	if err != nil {
		return err
	}
	log.Infof("Get instance - GetInstance test resources cleaned up")

	return nil
}

// GetInstanceMetric returns the metric name to which the GetInstance test metric should be emitted into
func (c canary) GetInstanceMetric() string {
	return getInstanceMetric
}

func validateGetInstanceCall(cssWrapper wrappers.CSSWrapper, clusterName string, instanceARN string) error {
	// Takes some time for CSS to get the instance.
	// Retry get call once every 10 seconds for 2 minutes.
	// TODO: Change sleep and retry related numbers to constants
	found := false
	for i := 0; i < 12; i++ {
		instance, err := cssWrapper.GetInstance(&clusterName, &instanceARN)
		if err != nil {
			rootErr := errors.Cause(err)
			if _, ok := rootErr.(*operations.GetInstanceNotFound); ok {
				time.Sleep(10 * time.Second)
				continue
			}
			return err
		}
		if instance == nil {
			return errors.Errorf("Get instance call with cluster name '%s' and instance ARN '%s' returned a nil instance. ",
				clusterName, instanceARN)
		}
		if instance.Entity == nil {
			return errors.Errorf("Get instance call with cluster name '%s' and instance ARN '%s' returned a nil instance entity. ",
				clusterName, instanceARN)
		}
		if aws.StringValue(instance.Entity.ContainerInstanceARN) == instanceARN {
			found = true
			break
		}
		return errors.Errorf("Expected instance belonging to cluster with name '%s' to have ARN '%s' but was '%s'. ",
			clusterName, instanceARN, *instance.Entity.ContainerInstanceARN)
	}

	if !found {
		return errors.Errorf("Instance with ARN '%s' belonging to cluster with name"+
			" '%s' not found in CSS. ", instanceARN, clusterName)
	}

	return nil
}
