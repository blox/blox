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
		detail.RunningTasksCount == nil || detail.Status == nil || detail.VersionInfo == nil {
		return errors.New("Instance detail is invalid")
	}
	return nil
}

// ToContainerInstance tranlates a container instance represented by the internal structure (types.ContainerInstance) to it's external representation (models.ContainerInstance)
func ToContainerInstance(instance types.ContainerInstance) (models.ContainerInstance, error) {
	err := validateContainerInstance(instance)
	if err != nil {
		return models.ContainerInstance{}, err
	}
	regRes := make([]*models.ContainerInstanceResource, len(instance.Detail.RegisteredResources))
	for i := range instance.Detail.RegisteredResources {
		r := instance.Detail.RegisteredResources[i]
		regRes[i] = &models.ContainerInstanceResource{
			Name:  r.Name,
			Type:  r.Type,
			Value: r.Value,
		}
	}

	remRes := make([]*models.ContainerInstanceResource, len(instance.Detail.RemainingResources))
	for i := range instance.Detail.RegisteredResources {
		r := instance.Detail.RemainingResources[i]
		remRes[i] = &models.ContainerInstanceResource{
			Name:  r.Name,
			Type:  r.Type,
			Value: r.Value,
		}
	}

	versionInfo := models.ContainerInstanceVersionInfo{
		AgentHash:     instance.Detail.VersionInfo.AgentHash,
		AgentVersion:  instance.Detail.VersionInfo.AgentVersion,
		DockerVersion: instance.Detail.VersionInfo.DockerVersion,
	}

	containerInstance := models.ContainerInstance{
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
		VersionInfo:          &versionInfo,
	}

	if instance.Detail.Attributes != nil {
		attributes := make([]*models.ContainerInstanceAttribute, len(instance.Detail.Attributes))
		for i := range instance.Detail.Attributes {
			a := instance.Detail.Attributes[i]
			attributes[i] = &models.ContainerInstanceAttribute{
				Name:  a.Name,
				Value: a.Value,
			}
		}
		containerInstance.Attributes = attributes
	}

	return containerInstance, nil
}

func validateTask(task types.Task) error {
	// TODO: Validate inner structs in task.Detail
	detail := task.Detail
	if detail == nil || detail.ClusterARN == nil || detail.ContainerInstanceARN == nil ||
		detail.Containers == nil || detail.CreatedAt == nil || detail.DesiredStatus == nil ||
		detail.LastStatus == nil || detail.Overrides == nil || detail.TaskARN == nil ||
		detail.TaskDefinitionARN == nil {
		return errors.New("Task detail is invalid")
	}
	return nil
}

// ToTask tranlates a task represented by the internal structure (types.Task) to it's external representation (models.Task)
func ToTask(task types.Task) (models.Task, error) {
	err := validateTask(task)
	if err != nil {
		return models.Task{}, err
	}

	containers := make([]*models.TaskContainer, len(task.Detail.Containers))
	for i := range task.Detail.Containers {
		c := task.Detail.Containers[i]
		containers[i] = &models.TaskContainer{
			ContainerARN: c.ContainerARN,
			ExitCode:     c.ExitCode,
			LastStatus:   c.LastStatus,
			Name:         c.Name,
			Reason:       c.Reason,
		}
		if c.NetworkBindings != nil {
			networkBindings := make([]*models.TaskNetworkBinding, len(c.NetworkBindings))
			for j := range c.NetworkBindings {
				n := c.NetworkBindings[j]
				networkBindings[j] = &models.TaskNetworkBinding{
					BindIP:        n.BindIP,
					ContainerPort: n.ContainerPort,
					HostPort:      n.HostPort,
					Protocol:      n.Protocol,
				}
			}
			containers[i].NetworkBindings = networkBindings
		}
	}

	containerOverrides := make([]*models.TaskContainerOverride, len(task.Detail.Overrides.ContainerOverrides))
	for i := range task.Detail.Overrides.ContainerOverrides {
		c := task.Detail.Overrides.ContainerOverrides[i]
		containerOverrides[i] = &models.TaskContainerOverride{
			Command: c.Command,
			Name:    c.Name,
		}
		if c.Environment != nil {
			env := make([]*models.TaskEnvironment, len(c.Environment))
			for j := range c.Environment {
				e := c.Environment[j]
				env[j] = &models.TaskEnvironment{
					Name:  e.Name,
					Value: e.Value,
				}
			}
			containerOverrides[i].Environment = env
		}
	}

	overrides := models.TaskOverride{
		ContainerOverrides: containerOverrides,
		TaskRoleArn:        task.Detail.Overrides.TaskRoleArn,
	}

	return models.Task{
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
	}, nil
}
