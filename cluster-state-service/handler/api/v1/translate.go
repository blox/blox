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

package v1

import (
	"errors"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	storetypes "github.com/blox/blox/cluster-state-service/handler/store/types"
	"github.com/blox/blox/cluster-state-service/handler/types"
	"github.com/blox/blox/cluster-state-service/swagger/v1/generated/models"
)

const (
	resourceIntegerType   = "INTEGER"
	resourceDoubleType    = "DOUBLE"
	resourceLongType      = "LONG"
	resourceStringSetType = "STRINGSET"
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

	resourceType := aws.StringValue(r.Type)

	val := ""
	if resourceType == resourceDoubleType && r.DoubleValue != nil {
		val = strconv.FormatFloat(aws.Float64Value(r.DoubleValue), 'f', 2, 64)
	} else if resourceType == resourceIntegerType && r.IntegerValue != nil {
		val = strconv.FormatInt(aws.Int64Value(r.IntegerValue), 10)
	} else if resourceType == resourceLongType && r.LongValue != nil {
		val = strconv.FormatInt(aws.Int64Value(r.LongValue), 10)
	} else if resourceType == resourceStringSetType && r.StringSetValue != nil {
		strVal := make([]string, len(r.StringSetValue))
		for i := range r.StringSetValue {
			strVal[i] = aws.StringValue(r.StringSetValue[i])
		}
		val = strings.Join(strVal, ",")
	}
	resource.Value = &val

	return resource
}

// ToContainerInstance translates a container instance represented by the internal structure (storetypes.VersionedContainerInstance) to it's external representation (models.ContainerInstance)
func ToContainerInstance(versionedInstance storetypes.VersionedContainerInstance) (models.ContainerInstance, error) {
	c := versionedInstance.ContainerInstance
	err := validateContainerInstance(c)
	if err != nil {
		return models.ContainerInstance{}, err
	}
	regRes := make([]*models.ContainerInstanceResource, len(c.Detail.RegisteredResources))
	for i := range c.Detail.RegisteredResources {
		regRes[i] = toContainerInstanceResource(c.Detail.RegisteredResources[i])
	}

	remRes := make([]*models.ContainerInstanceResource, len(c.Detail.RemainingResources))
	for i := range c.Detail.RegisteredResources {
		remRes[i] = toContainerInstanceResource(c.Detail.RemainingResources[i])
	}

	versionInfo := models.ContainerInstanceVersionInfo{
		AgentHash:     c.Detail.VersionInfo.AgentHash,
		AgentVersion:  c.Detail.VersionInfo.AgentVersion,
		DockerVersion: c.Detail.VersionInfo.DockerVersion,
	}

	containerInstance := models.ContainerInstance{
		Metadata: &models.Metadata{
			EntityVersion: &versionedInstance.Version,
		},
		Entity: &models.ContainerInstanceDetail{
			AgentConnected:       c.Detail.AgentConnected,
			AgentUpdateStatus:    c.Detail.AgentUpdateStatus,
			ClusterARN:           c.Detail.ClusterARN,
			ContainerInstanceARN: c.Detail.ContainerInstanceARN,
			EC2InstanceID:        c.Detail.EC2InstanceID,
			RegisteredResources:  regRes,
			RemainingResources:   remRes,
			Status:               c.Detail.Status,
			VersionInfo:          &versionInfo,
		},
	}

	if c.Detail.Attributes != nil {
		attributes := make([]*models.ContainerInstanceAttribute, len(c.Detail.Attributes))
		for i := range c.Detail.Attributes {
			a := c.Detail.Attributes[i]
			attributes[i] = &models.ContainerInstanceAttribute{
				Name:  a.Name,
				Value: a.Value,
			}
		}
		containerInstance.Entity.Attributes = attributes
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

// ToTask translates a task represented by the internal structure (storetypes.VersionedTask) to it's external representation (models.Task)
func ToTask(versionedTask storetypes.VersionedTask) (models.Task, error) {
	t := versionedTask.Task
	err := validateTask(t)
	if err != nil {
		return models.Task{}, err
	}

	containers := make([]*models.TaskContainer, len(t.Detail.Containers))
	for i := range t.Detail.Containers {
		c := t.Detail.Containers[i]
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

	containerOverrides := make([]*models.TaskContainerOverride, len(t.Detail.Overrides.ContainerOverrides))
	for i := range t.Detail.Overrides.ContainerOverrides {
		c := t.Detail.Overrides.ContainerOverrides[i]
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
		TaskRoleArn:        t.Detail.Overrides.TaskRoleArn,
	}

	return models.Task{
		Metadata: &models.Metadata{
			EntityVersion: &versionedTask.Version,
		},
		Entity: &models.TaskDetail{
			ClusterARN:           t.Detail.ClusterARN,
			ContainerInstanceARN: t.Detail.ContainerInstanceARN,
			Containers:           containers,
			CreatedAt:            t.Detail.CreatedAt,
			DesiredStatus:        t.Detail.DesiredStatus,
			LastStatus:           t.Detail.LastStatus,
			Overrides:            &overrides,
			StartedAt:            t.Detail.StartedAt,
			StartedBy:            t.Detail.StartedBy,
			StoppedAt:            t.Detail.StoppedAt,
			StoppedReason:        t.Detail.StoppedReason,
			TaskARN:              t.Detail.TaskARN,
			TaskDefinitionARN:    t.Detail.TaskDefinitionARN,
		},
	}, nil
}
