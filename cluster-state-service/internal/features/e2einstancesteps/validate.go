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
	"github.com/blox/blox/cluster-state-service/internal/models"
	"github.com/pkg/errors"
)

func ValidateInstancesMatch(ecsInstance ecs.ContainerInstance, cssInstance models.ContainerInstance) error {
	if *ecsInstance.ContainerInstanceArn != *cssInstance.ContainerInstanceARN ||
		*ecsInstance.Status != *cssInstance.Status {
		return errors.New("Container instances don't match. ")
	}
	return nil
}

func ValidateListContainsInstance(ecsInstance ecs.ContainerInstance, cssInstanceList []models.ContainerInstance) error {
	instanceARN := *ecsInstance.ContainerInstanceArn
	var cssInstance models.ContainerInstance
	for _, i := range cssInstanceList {
		if *i.ContainerInstanceARN == instanceARN {
			cssInstance = i
			break
		}
	}
	if cssInstance.ContainerInstanceARN == nil {
		return errors.Errorf("Instance with ARN '%s' not found in response. ", instanceARN)
	}
	return ValidateInstancesMatch(ecsInstance, cssInstance)
}
