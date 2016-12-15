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
	"errors"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/blox/blox/cluster-state-service/handler/api/v1/models"
	"github.com/blox/blox/cluster-state-service/handler/types"
)

func validateContainerInstance(instance types.ContainerInstance) error {
	// TODO: Validate inner structs in instance.Detail
	detail := instance.Detail
	if detail == nil {
		return errors.New("Instance detail cannot be empty")
	}
	if detail.AgentConnected == nil {
		return errors.New("Instance agent connected cannot be empty")
	}
	if detail.ClusterARN == nil {
		return errors.New("Instance cluster ARN cannot be empty")
	}
	if detail.ContainerInstanceARN == nil {
		return errors.New("Instance ARN cannot be empty")
	}
	if detail.RegisteredResources == nil {
		return errors.New("Instance registered resources cannot be empty")
	}
	if detail.RemainingResources == nil {
		return errors.New("Instance remaining resources cannot be empty")
	}
	if detail.Status == nil {
		return errors.New("Instance status cannot be empty")
	}
	if detail.VersionInfo == nil {
		return errors.New("Instance version info cannot be empty")
	}
	return nil
}

func toContainerInstanceResource(r *types.Resource) *models.ContainerInstanceResource {
	resource := &models.ContainerInstanceResource{
		Name: r.Name,
		Type: r.Type,
	}

	val := ""
	if r.DoubleValue != nil && aws.Float64Value(r.DoubleValue) != 0.0 {
		val = strconv.FormatFloat(aws.Float64Value(r.DoubleValue), 'f', 2, 64)
	} else if r.IntegerValue != nil && aws.Int64Value(r.IntegerValue) != 0 {
		val = strconv.FormatInt(aws.Int64Value(r.IntegerValue), 10)
	} else if r.LongValue != nil && aws.Int64Value(r.LongValue) != 0 {
		val = strconv.FormatInt(aws.Int64Value(r.LongValue), 10)
	} else if r.StringSetValue != nil {
		strVal := make([]string, len(r.StringSetValue))
		for i := range r.StringSetValue {
			strVal[i] = aws.StringValue(r.StringSetValue[i])
		}
		val = strings.Join(strVal, ",")
	}
	resource.Value = &val

	return resource
}

// ToContainerInstance tranlates a container instance represented by the internal structure (types.ContainerInstance) to it's external representation (models.ContainerInstance)
func ToContainerInstance(instance types.ContainerInstance) (models.ContainerInstance, error) {
	err := validateContainerInstance(instance)
	if err != nil {
		return models.ContainerInstance{}, err
	}
	regRes := make([]*models.ContainerInstanceResource, len(instance.Detail.RegisteredResources))
	for i := range instance.Detail.RegisteredResources {
		regRes[i] = toContainerInstanceResource(instance.Detail.RegisteredResources[i])
	}

	remRes := make([]*models.ContainerInstanceResource, len(instance.Detail.RemainingResources))
	for i := range instance.Detail.RegisteredResources {
		remRes[i] = toContainerInstanceResource(instance.Detail.RemainingResources[i])
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
		RegisteredResources:  regRes,
		RemainingResources:   remRes,
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
	if detail == nil {
		return errors.New("Task detail cannot be empty")
	}
	if detail.ClusterARN == nil {
		return errors.New("Task cluster ARN cannot be empty")
	}
	if detail.ContainerInstanceARN == nil {
		return errors.New("Task container instance ARN cannot be empty")
	}
	if detail.Containers == nil {
		return errors.New("Task containers cannot be empty")
	}
	if detail.CreatedAt == nil {
		return errors.New("Task created at cannot be empty")
	}
	if detail.DesiredStatus == nil {
		return errors.New("Task desired status cannot be empty")
	}
	if detail.LastStatus == nil {
		return errors.New("Task last status cannot be empty")
	}
	if detail.Overrides == nil {
		return errors.New("Task overrides cannot be empty")
	}
	if detail.TaskARN == nil {
		return errors.New("Task ARN cannot be empty")
	}
	if detail.TaskDefinitionARN == nil {
		return errors.New("Task definition ARN cannot be empty")
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
