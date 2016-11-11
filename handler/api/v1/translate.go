// Copyright 2016 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the License). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the license file accompanying this file. This file is distributed
// on an AS IS BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package v1

import (
	"errors"

	"github.com/aws/amazon-ecs-event-stream-handler/handler/api/v1/models"
	"github.com/aws/amazon-ecs-event-stream-handler/handler/types"
)

func validateContainerInstance(instance types.ContainerInstance) error {
	// TODO: Validate inner structs in instance.Detail
	detail := instance.Detail
	if detail == nil || detail.AgentConnected == nil || detail.ClusterARN == nil ||
		detail.ContainerInstanceARN == nil || detail.PendingTasksCount == nil ||
		detail.RegisteredResources == nil || detail.RemainingResources == nil ||
		detail.RunningTasksCount == nil || detail.Status == nil || detail.Version == nil ||
		detail.VersionInfo == nil || detail.UpdatedAt == nil {
		return errors.New("Instance detail is invalid")
	}
	return nil
}

func ToContainerInstanceModel(instance types.ContainerInstance) (models.ContainerInstanceModel, error) {
	err := validateContainerInstance(instance)
	if err != nil {
		return models.ContainerInstanceModel{}, err
	}
	regRes := make([]*models.ContainerInstanceRegisteredResourceModel, len(instance.Detail.RegisteredResources))
	for i := range instance.Detail.RegisteredResources {
		r := instance.Detail.RegisteredResources[i]
		regRes[i] = &models.ContainerInstanceRegisteredResourceModel{
			Name:  r.Name,
			Type:  r.Type,
			Value: r.Value,
		}
	}

	remRes := make([]*models.ContainerInstanceRemainingResourceModel, len(instance.Detail.RemainingResources))
	for i := range instance.Detail.RegisteredResources {
		r := instance.Detail.RemainingResources[i]
		remRes[i] = &models.ContainerInstanceRemainingResourceModel{
			Name:  r.Name,
			Type:  r.Type,
			Value: r.Value,
		}
	}

	versionInfo := models.ContainerInstanceVersionInfoModel{
		AgentHash:     instance.Detail.VersionInfo.AgentHash,
		AgentVersion:  instance.Detail.VersionInfo.AgentVersion,
		DockerVersion: instance.Detail.VersionInfo.DockerVersion,
	}

	containerInstance := models.ContainerInstanceModel{
		AgentConnected:       instance.Detail.AgentConnected,
		AgentUpdateStatus:    instance.Detail.AgentUpdateStatus,
		ClusterARN:           instance.Detail.ClusterARN,
		ContainerInstanceARN: instance.Detail.ContainerInstanceARN,
		EC2InstanceID:        instance.Detail.EC2InstanceID,
		PendingTasksCount:    instance.Detail.PendingTasksCount,
		RegisteredResources:  regRes,
		RemainingResources:   remRes,
		RunningTasksCount:    instance.Detail.RunningTasksCount,
		Status:               instance.Detail.Status,
		Version:              instance.Detail.Version,
		VersionInfo:          &versionInfo,
		UpdatedAt:            instance.Detail.UpdatedAt,
	}

	if instance.Detail.Attributes != nil {
		attributes := make([]*models.ContainerInstanceAttributeModel, len(instance.Detail.Attributes))
		for i := range instance.Detail.Attributes {
			a := instance.Detail.Attributes[i]
			attributes[i] = &models.ContainerInstanceAttributeModel{
				Name:  a.Name,
				Value: a.Value,
			}
		}
		containerInstance.Attributes = attributes
	}

	return containerInstance, nil
}

func validateTaskModel(task types.Task) error {
	// TODO: Validate inner structs in task.Detail
	detail := task.Detail
	if detail == nil || detail.ClusterARN == nil || detail.ContainerInstanceARN == nil ||
		detail.Containers == nil || detail.CreatedAt == nil || detail.DesiredStatus == nil ||
		detail.LastStatus == nil || detail.Overrides == nil || detail.TaskARN == nil ||
		detail.TaskDefinitionARN == nil || detail.UpdatedAt == nil || detail.Version == nil {
		return errors.New("Task detail is invalid")
	}
	return nil
}

func ToTaskModel(task types.Task) (models.TaskModel, error) {
	err := validateTaskModel(task)
	if err != nil {
		return models.TaskModel{}, err
	}

	containers := make([]*models.TaskContainerModel, len(task.Detail.Containers))
	for i := range task.Detail.Containers {
		c := task.Detail.Containers[i]
		containers[i] = &models.TaskContainerModel{
			ContainerARN: c.ContainerARN,
			ExitCode:     c.ExitCode,
			LastStatus:   c.LastStatus,
			Name:         c.Name,
			Reason:       c.Reason,
		}
		if c.NetworkBindings != nil {
			networkBindings := make([]*models.TaskNetworkBindingModel, len(c.NetworkBindings))
			for j := range c.NetworkBindings {
				n := c.NetworkBindings[j]
				networkBindings[j] = &models.TaskNetworkBindingModel{
					BindIP:        n.BindIP,
					ContainerPort: n.ContainerPort,
					HostPort:      n.HostPort,
					Protocol:      n.Protocol,
				}
			}
			containers[i].NetworkBindings = networkBindings
		}
	}

	containerOverrides := make([]*models.TaskContainerOverrideModel, len(task.Detail.Overrides.ContainerOverrides))
	for i := range task.Detail.Overrides.ContainerOverrides {
		c := task.Detail.Overrides.ContainerOverrides[i]
		containerOverrides[i] = &models.TaskContainerOverrideModel{
			Command: c.Command,
			Name:    c.Name,
		}
		if c.Environment != nil {
			env := make([]*models.TaskEnvironmentModel, len(c.Environment))
			for j := range c.Environment {
				e := c.Environment[j]
				env[j] = &models.TaskEnvironmentModel{
					Name:  e.Name,
					Value: e.Value,
				}
			}
			containerOverrides[i].Environment = env
		}
	}

	overrides := models.TaskOverrideModel{
		ContainerOverrides: containerOverrides,
		TaskRoleArn:        task.Detail.Overrides.TaskRoleArn,
	}

	return models.TaskModel{
		ClusterARN:           task.Detail.ClusterARN,
		ContainerInstanceARN: task.Detail.ContainerInstanceARN,
		Containers:           containers,
		CreatedAt:            task.Detail.CreatedAt,
		DesiredStatus:        task.Detail.DesiredStatus,
		LastStatus:           task.Detail.LastStatus,
		Overrides:            &overrides,
		StartedAt:            task.Detail.StartedAt,
		StartedBy:            task.Detail.StartedBy,
		StoppedAt:            task.Detail.StoppedAt,
		StoppedReason:        task.Detail.StoppedReason,
		TaskARN:              task.Detail.TaskARN,
		TaskDefinitionARN:    task.Detail.TaskDefinitionARN,
		UpdatedAt:            task.Detail.UpdatedAt,
		Version:              task.Detail.Version,
	}, nil
}
