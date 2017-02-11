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

package util

import (
	"time"

	"github.com/blox/blox/cluster-state-service/canary/pkg/wrappers"
	"github.com/pkg/errors"
)

// ValidateECSContainerInstanceAndGetInstanceARN validates that an ECS instance was launched
// within a cluster 'clusterName' and returns the ARN of the instance
func ValidateECSContainerInstanceAndGetInstanceARN(ecsWrapper wrappers.ECSWrapper, clusterName string) (string, error) {
	// Takes some time for EC2 instance to come up and for the ECS agent to
	// register the container instance with ECS.
	// Retry ECS list call once every minute for 10 minutes
	// TODO: Change sleep and retry related numbers to constants
	var instanceARN string
	found := false
	for i := 0; i < 10; i++ {
		instances, err := ecsWrapper.ListContainerInstances(&clusterName)
		if err != nil {
			return "", err
		}
		if len(instances) > 1 {
			return "", errors.Errorf("Expected a maximum of 1 instance registered with "+
				"ECS cluster with name '%s' but was '%d'. ", clusterName, len(instances))
		}
		if len(instances) == 1 {
			instanceARN = *instances[0]
			found = true
			break
		}
		time.Sleep(1 * time.Minute)
	}

	if !found {
		return "", errors.Errorf("No ECS instance found in cluster with name '%s'. ",
			clusterName)
	}

	return instanceARN, nil
}
