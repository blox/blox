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

package v1

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/blox/blox/daemon-scheduler/pkg/types"
	"github.com/blox/blox/daemon-scheduler/swagger/v1/generated/models"
)

func toEnvironmentModel(envType types.Environment) models.Environment {
	health := models.HealthStatusHealthy
	if envType.Health == types.EnvironmentUnhealthy {
		health = models.HealthStatusUnhealthy
	}
	return models.Environment{
		Name: &envType.Name,
		InstanceGroup: &models.InstanceGroup{
			Cluster: envType.Cluster,
		},
		Health:          health,
		DeploymentToken: envType.Token,
		TaskDefinition:  envType.DesiredTaskDefinition,
	}
}

func toDeploymentModel(envName *string, depType types.Deployment) *models.Deployment {
	instanceArns := []string{}
	for _, failure := range depType.FailedInstances {
		instanceArns = append(instanceArns, aws.StringValue(failure.Arn))
	}

	return &models.Deployment{
		EnvironmentName: envName,
		ID:              &depType.ID,
		Status:          aws.String(toDeploymentStatus(depType.Status)),
		TaskDefinition:  aws.String(depType.TaskDefinition),
		FailedInstances: instanceArns,
	}
}

func toDeploymentsModel(envName *string, depTypes []types.Deployment) *models.Deployments {
	depModels := []*models.Deployment{}
	for _, depType := range depTypes {
		depModel := toDeploymentModel(envName, depType)
		depModels = append(depModels, depModel)
	}
	return &models.Deployments{
		Items: depModels,
	}
}

func toDeploymentStatus(statusType types.DeploymentStatus) string {
	switch {
	case types.DeploymentPending == statusType:
		return "pending"
	case types.DeploymentInProgress == statusType:
		return "running"
	case types.DeploymentCompleted == statusType:
		return "completed"
	default:
		return "unknown"
	}
}
