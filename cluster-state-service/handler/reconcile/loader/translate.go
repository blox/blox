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

package loader

import (
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/blox/blox/cluster-state-service/handler/types"
)

const (
	timeLayout = "2006-01-02T15:04:05Z"
	version    = int64(-1)
)

// ToTask tranlates an ECS task to the internal task type
func ToTask(ecsTask ecs.Task) types.Task {
	createdAt := ecsTask.CreatedAt.Format(timeLayout)
	updatedAt := currentTime()
	taskVersion := version
	taskDetail := types.TaskDetail{
		ClusterARN:           ecsTask.ClusterArn,
		ContainerInstanceARN: ecsTask.ContainerInstanceArn,
		Containers:           toContainers(ecsTask.Containers),
		CreatedAt:            &createdAt,
		DesiredStatus:        ecsTask.DesiredStatus,
		LastStatus:           ecsTask.LastStatus,
		Overrides:            toOverrides(ecsTask.Overrides),
		StartedBy:            aws.StringValue(ecsTask.StartedBy),
		StoppedReason:        aws.StringValue(ecsTask.StoppedReason),
		TaskARN:              ecsTask.TaskArn,
		TaskDefinitionARN:    ecsTask.TaskDefinitionArn,
		Version:              &taskVersion,
		UpdatedAt:            &updatedAt,
	}
	if ecsTask.StartedAt != nil {
		startedAt := ecsTask.StartedAt.Format(timeLayout)
		taskDetail.StartedAt = startedAt
	}
	if ecsTask.StoppedAt != nil {
		stoppedAt := ecsTask.StoppedAt.Format(timeLayout)
		taskDetail.StoppedAt = stoppedAt
	}
	return types.Task{
		Detail: &taskDetail,
	}
}

func toContainers(ecsContainers []*ecs.Container) []*types.Container {
	containers := make([]*types.Container, len(ecsContainers))
	for i := range ecsContainers {
		ecsContainer := ecsContainers[i]
		container := types.Container{
			ContainerARN:    ecsContainer.ContainerArn,
			ExitCode:        aws.Int64Value(ecsContainer.ExitCode),
			LastStatus:      ecsContainer.LastStatus,
			Name:            ecsContainer.Name,
			NetworkBindings: toNetworkBindings(ecsContainer.NetworkBindings),
			Reason:          aws.StringValue(ecsContainer.Reason),
		}
		containers[i] = &container
	}
	return containers
}

func toNetworkBindings(ecsNetworkBindings []*ecs.NetworkBinding) []*types.NetworkBinding {
	networkBindings := make([]*types.NetworkBinding, len(ecsNetworkBindings))
	for i := range ecsNetworkBindings {
		ecsNB := ecsNetworkBindings[i]
		nb := types.NetworkBinding{
			BindIP:        ecsNB.BindIP,
			ContainerPort: ecsNB.ContainerPort,
			HostPort:      ecsNB.HostPort,
			Protocol:      aws.StringValue(ecsNB.Protocol),
		}
		networkBindings[i] = &nb
	}
	return networkBindings
}

func toOverrides(ecsOverrides *ecs.TaskOverride) *types.Overrides {
	return &types.Overrides{
		ContainerOverrides: toContainerOverrides(ecsOverrides.ContainerOverrides),
		TaskRoleArn:        aws.StringValue(ecsOverrides.TaskRoleArn),
	}
}

func toContainerOverrides(ecsContainerOverrides []*ecs.ContainerOverride) []*types.ContainerOverrides {
	containerOverrides := make([]*types.ContainerOverrides, len(ecsContainerOverrides))
	for i := range ecsContainerOverrides {
		ecsCO := ecsContainerOverrides[i]
		co := types.ContainerOverrides{
			Command:     toCommand(ecsCO.Command),
			Environment: toEnvironment(ecsCO.Environment),
			Name:        ecsCO.Name,
		}
		containerOverrides[i] = &co
	}
	return containerOverrides
}

func toCommand(ecsCommand []*string) []string {
	command := make([]string, len(ecsCommand))
	for i := range ecsCommand {
		command[i] = aws.StringValue(ecsCommand[i])
	}
	return command
}

func toEnvironment(ecsEnvironment []*ecs.KeyValuePair) []*types.Environment {
	environment := make([]*types.Environment, len(ecsEnvironment))
	for i := range ecsEnvironment {
		ecsEnv := ecsEnvironment[i]
		env := types.Environment{
			Name:  ecsEnv.Name,
			Value: ecsEnv.Value,
		}
		environment[i] = &env
	}
	return environment
}

// ToContainerInstance tranlates an ECS container instance to the internal container instance type
func ToContainerInstance(ecsInstance ecs.ContainerInstance, clusterARN string) types.ContainerInstance {
	instanceVersion := version
	updatedAt := currentTime()
	insDetail := types.InstanceDetail{
		AgentConnected:       ecsInstance.AgentConnected,
		AgentUpdateStatus:    aws.StringValue(ecsInstance.AgentUpdateStatus),
		Attributes:           toAttributes(ecsInstance.Attributes),
		ClusterARN:           &clusterARN,
		ContainerInstanceARN: ecsInstance.ContainerInstanceArn,
		EC2InstanceID:        aws.StringValue(ecsInstance.Ec2InstanceId),
		RegisteredResources:  toResources(ecsInstance.RegisteredResources),
		RemainingResources:   toResources(ecsInstance.RemainingResources),
		Status:               ecsInstance.Status,
		Version:              &instanceVersion,
		VersionInfo:          toVersionInfo(ecsInstance.VersionInfo),
		UpdatedAt:            &updatedAt,
	}
	return types.ContainerInstance{
		Detail: &insDetail,
	}
}

func toAttributes(ecsAttributes []*ecs.Attribute) []*types.Attribute {
	attributes := make([]*types.Attribute, len(ecsAttributes))
	for i := range ecsAttributes {
		ecsAttribute := ecsAttributes[i]
		a := types.Attribute{
			Name:  ecsAttribute.Name,
			Value: ecsAttribute.Value,
		}
		attributes[i] = &a
	}
	return attributes
}

func toResources(ecsResources []*ecs.Resource) []*types.Resource {
	resources := make([]*types.Resource, len(ecsResources))
	for i := range ecsResources {
		ecsRes := ecsResources[i]
		r := types.Resource{
			Name: ecsRes.Name,
			Type: ecsRes.Type,
		}
		val := ""
		if ecsRes.DoubleValue != nil {
			val = strconv.FormatFloat(aws.Float64Value(ecsRes.DoubleValue), 'f', 2, 64)
		} else if ecsRes.IntegerValue != nil {
			val = strconv.FormatInt(aws.Int64Value(ecsRes.IntegerValue), 10)
		} else if ecsRes.LongValue != nil {
			val = strconv.FormatInt(aws.Int64Value(ecsRes.LongValue), 10)
		} else if ecsRes.StringSetValue != nil {
			strVal := make([]string, len(ecsRes.StringSetValue))
			for i := range ecsRes.StringSetValue {
				strVal[i] = aws.StringValue(ecsRes.StringSetValue[i])
			}
			val = strings.Join(strVal, ",")
		}
		r.Value = &val
		resources[i] = &r
	}
	return resources
}

func toVersionInfo(ecsVersionInfo *ecs.VersionInfo) *types.VersionInfo {
	return &types.VersionInfo{
		AgentHash:     aws.StringValue(ecsVersionInfo.AgentHash),
		AgentVersion:  aws.StringValue(ecsVersionInfo.AgentVersion),
		DockerVersion: aws.StringValue(ecsVersionInfo.DockerVersion),
	}
}

func currentTime() string {
	return time.Now().Format(timeLayout)
}
