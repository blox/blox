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
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/blox/blox/cluster-state-service/swagger/v1/generated/models"
	"github.com/pkg/errors"
)

func ValidateInstancesMatch(ecsInstance ecs.ContainerInstance, cssInstance models.ContainerInstance) error {
	if *ecsInstance.ContainerInstanceArn != *cssInstance.Entity.ContainerInstanceARN ||
		*ecsInstance.Status != *cssInstance.Entity.Status {
		return errors.New("Container instances don't match. ")
	}
	return nil
}

func ValidateListContainsInstance(ecsInstance ecs.ContainerInstance, cssInstanceList []models.ContainerInstance) error {
	cssInstance, err := ValidateListContainsInstanceArn(*ecsInstance.ContainerInstanceArn, cssInstanceList)
	if err != nil {
		return err
	}

	return ValidateInstancesMatch(ecsInstance, *cssInstance)
}

func ValidateListContainsInstanceArn(instanceARN string, cssInstanceList []models.ContainerInstance) (*models.ContainerInstance, error) {
	for _, i := range cssInstanceList {
		if *i.Entity.ContainerInstanceARN == instanceARN {
			return &i, nil
		}
	}
	return nil, errors.Errorf("Instance with ARN '%s' not found in response. ", instanceARN)
}
