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
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/blox/blox/cluster-state-service/handler/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TranslateTestSuite struct {
	suite.Suite
	instance    types.ContainerInstance
	ecsInstance ecs.ContainerInstance
	task        types.Task
	ecsTask     ecs.Task
	clusterARN  string
}

func (suite *TranslateTestSuite) SetupTest() {
	agentConnected := true
	agentUpdateStatus := "PENDING"
	clusterARN := "arn:aws:ecs:us-east-1:123456789012:cluster/clusterName1"
	containerInstanceARN := "arn:aws:ecs:us-east-1:123456789012:container-instance/b6b9eace-958e-4f2a-a09c-8cf43b76cf97"
	containerARN := "arn:aws:ecs:us-east-1:123456789012:container/57156e30-e410-4773-9a9e-ae8264c10bbd"
	taskARN := "arn:aws:ecs:us-east-1:123456789012:task/271022c0-f894-4aa2-b063-25bae55088d5"
	taskDefinitionARN := "arn:aws:ecs:us-east-1:123456789012:task-definition/testTask:1"
	ec2InstanceID := "i-12345678"
	attributeName := "Name"
	attributeVal := "com.amazonaws.ecs.capability.privileged-container"
	resourceName := "CPU"
	resourceType := "INTEGER"
	resourceIntVal := int64(1024)
	containerStatus := "ACTIVE"
	agentHash := "2ad18e61a4ba696287902b0bc177a031a7112fb6"
	agentVersion := "1.13.0"
	dockerVersion := "1.12.3"
	containerName := "testContainer"
	createdAt := "2016-11-07T15:30:00Z"
	ecsCreatedAt, err := time.Parse(timeLayout, createdAt)
	assert.Nil(suite.T(), err, "Cannot setup testSuite: Unexpected error when parsing time")
	desiredStatus := "RUNNING"
	lastStatus := "PENDING"
	command := []string{"sleep", "300"}
	startedAt := "2016-11-07T15:45:00Z"
	ecsStartedAt, err := time.Parse(timeLayout, startedAt)
	assert.Nil(suite.T(), err, "Cannot setup testSuite: Unexpected error when parsing time")

	suite.ecsInstance = ecs.ContainerInstance{
		AgentConnected:    &agentConnected,
		AgentUpdateStatus: &agentUpdateStatus,
		Attributes: []*ecs.Attribute{
			&ecs.Attribute{
				Name:  &attributeName,
				Value: &attributeVal,
			},
		},
		ContainerInstanceArn: &containerInstanceARN,
		Ec2InstanceId:        &ec2InstanceID,
		RegisteredResources: []*ecs.Resource{
			&ecs.Resource{
				Name:         &resourceName,
				Type:         &resourceType,
				IntegerValue: &resourceIntVal,
			},
		},
		RemainingResources: []*ecs.Resource{
			&ecs.Resource{
				Name:         &resourceName,
				Type:         &resourceType,
				IntegerValue: &resourceIntVal,
			},
		},
		Status: &containerStatus,
		VersionInfo: &ecs.VersionInfo{
			AgentHash:     &agentHash,
			AgentVersion:  &agentVersion,
			DockerVersion: &dockerVersion,
		},
	}

	resourceVal := strconv.FormatInt(resourceIntVal, 10)
	instanceVersion := version
	suite.instance = types.ContainerInstance{
		Detail: &types.InstanceDetail{
			AgentConnected:    &agentConnected,
			AgentUpdateStatus: agentUpdateStatus,
			Attributes: []*types.Attribute{
				&types.Attribute{
					Name:  &attributeName,
					Value: &attributeVal,
				},
			},
			ClusterARN:           &clusterARN,
			ContainerInstanceARN: &containerInstanceARN,
			EC2InstanceID:        ec2InstanceID,
			RegisteredResources: []*types.Resource{
				&types.Resource{
					Name:  &resourceName,
					Type:  &resourceType,
					Value: &resourceVal,
				},
			},
			RemainingResources: []*types.Resource{
				&types.Resource{
					Name:  &resourceName,
					Type:  &resourceType,
					Value: &resourceVal,
				},
			},
			Status:  &containerStatus,
			Version: &instanceVersion,
			VersionInfo: &types.VersionInfo{
				AgentHash:     agentHash,
				AgentVersion:  agentVersion,
				DockerVersion: dockerVersion,
			},
		},
	}

	suite.ecsTask = ecs.Task{
		ClusterArn:           &clusterARN,
		ContainerInstanceArn: &containerInstanceARN,
		Containers: []*ecs.Container{
			&ecs.Container{
				ContainerArn: &containerARN,
				LastStatus:   &containerStatus,
				Name:         &containerName,
			},
		},
		CreatedAt:     &ecsCreatedAt,
		DesiredStatus: &desiredStatus,
		LastStatus:    &lastStatus,
		Overrides: &ecs.TaskOverride{
			ContainerOverrides: []*ecs.ContainerOverride{
				&ecs.ContainerOverride{
					Command: []*string{&command[0], &command[1]},
					Name:    &containerName,
				},
			},
		},
		StartedAt:         &ecsStartedAt,
		TaskArn:           &taskARN,
		TaskDefinitionArn: &taskDefinitionARN,
	}

	taskVersion := version
	suite.task = types.Task{
		Detail: &types.TaskDetail{
			ClusterARN:           &clusterARN,
			ContainerInstanceARN: &containerInstanceARN,
			Containers: []*types.Container{
				&types.Container{
					ContainerARN:    &containerARN,
					LastStatus:      &containerStatus,
					Name:            &containerName,
					NetworkBindings: []*types.NetworkBinding{},
				},
			},
			CreatedAt:     &createdAt,
			DesiredStatus: &desiredStatus,
			LastStatus:    &lastStatus,
			Overrides: &types.Overrides{
				ContainerOverrides: []*types.ContainerOverrides{
					&types.ContainerOverrides{
						Command:     command,
						Environment: []*types.Environment{},
						Name:        &containerName,
					},
				},
			},
			StartedAt:         startedAt,
			TaskARN:           &taskARN,
			TaskDefinitionARN: &taskDefinitionARN,
			Version:           &taskVersion,
		},
	}

	suite.clusterARN = clusterARN
}

func TestTranslateTestSuite(t *testing.T) {
	suite.Run(t, new(TranslateTestSuite))
}

func (suite *TranslateTestSuite) TestToContainerInstance() {
	instance := ToContainerInstance(suite.ecsInstance, suite.clusterARN)
	// TODO: Mock out Time.Now() and avoid this hack
	expectedInstance := suite.instance
	updatedAt := *instance.Detail.UpdatedAt
	expectedInstance.Detail.UpdatedAt = &updatedAt
	assert.Equal(suite.T(), expectedInstance, instance, "Translated instance does not match expected instance")
}

func (suite *TranslateTestSuite) TestToTask() {
	task := ToTask(suite.ecsTask)
	// TODO: Mock out Time.Now() and avoid this hack
	expectedTask := suite.task
	updatedAt := *task.Detail.UpdatedAt
	expectedTask.Detail.UpdatedAt = &updatedAt
	assert.Equal(suite.T(), expectedTask, task, "Translated task does not match expected task")
}
