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
	"testing"

	"github.com/blox/blox/cluster-state-service/handler/api/v1/models"
	"github.com/blox/blox/cluster-state-service/handler/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	agentUpdateStatus = "pending"
	ecsInstanceID     = "i-12345678"
)

var (
	attributeName = "Name"
	attributeVal  = "com.amazonaws.ecs.capability.privileged-container"
	resourceName  = "CPU"
	resourceType  = "INTEGER"
	resourceVal   = ""
)

type TranslateTestSuite struct {
	suite.Suite
	instance    types.ContainerInstance
	extInstance models.ContainerInstance
	task        types.Task
	extTask     models.Task
}

func (suite *TranslateTestSuite) SetupTest() {
	versionInfo := types.VersionInfo{}
	attribute := types.Attribute{
		Name:  &attributeName,
		Value: &attributeVal,
	}
	resource := types.Resource{
		Name:  &resourceName,
		Type:  &resourceType,
		Value: &resourceVal,
	}
	instanceDetail := types.InstanceDetail{
		AgentConnected:       &agentConnected1,
		AgentUpdateStatus:    agentUpdateStatus,
		Attributes:           []*types.Attribute{&attribute},
		ClusterARN:           &clusterARN1,
		ContainerInstanceARN: &instanceARN1,
		EC2InstanceID:        ecsInstanceID,
		RegisteredResources:  []*types.Resource{&resource},
		RemainingResources:   []*types.Resource{&resource},
		Status:               &instanceStatus1,
		VersionInfo:          &versionInfo,
	}
	suite.instance = types.ContainerInstance{
		ID:        &id1,
		Account:   &accountID,
		Time:      &time,
		Region:    &region,
		Resources: []string{instanceARN1},
		Detail:    &instanceDetail,
	}

	versionInfoModel := models.ContainerInstanceVersionInfo{}
	attributeModel := models.ContainerInstanceAttribute{
		Name:  &attributeName,
		Value: &attributeVal,
	}
	extResource := models.ContainerInstanceResource{
		Name:  &resourceName,
		Type:  &resourceType,
		Value: &resourceVal,
	}
	suite.extInstance = models.ContainerInstance{
		AgentConnected:       &agentConnected1,
		AgentUpdateStatus:    agentUpdateStatus,
		Attributes:           []*models.ContainerInstanceAttribute{&attributeModel},
		ClusterARN:           &clusterARN1,
		ContainerInstanceARN: &instanceARN1,
		EC2InstanceID:        ecsInstanceID,
		RegisteredResources:  []*models.ContainerInstanceResource{&extResource},
		RemainingResources:   []*models.ContainerInstanceResource{&extResource},
		Status:               &instanceStatus1,
		VersionInfo:          &versionInfoModel,
	}

	container := types.Container{
		ContainerARN: &containerARN1,
		LastStatus:   &taskStatus1,
		Name:         &taskName,
	}
	containerOverrides := types.ContainerOverrides{
		Name: &taskName,
	}
	overrides := types.Overrides{
		ContainerOverrides: []*types.ContainerOverrides{&containerOverrides},
	}
	taskDetail := types.TaskDetail{
		ClusterARN:           &clusterARN1,
		ContainerInstanceARN: &instanceARN1,
		Containers:           []*types.Container{&container},
		CreatedAt:            &createdAt,
		DesiredStatus:        &taskStatus1,
		LastStatus:           &taskStatus1,
		Overrides:            &overrides,
		TaskARN:              &taskARN1,
		TaskDefinitionARN:    &taskDefinitionARN,
	}
	suite.task = types.Task{
		ID:        &id1,
		Account:   &accountID,
		Time:      &time,
		Region:    &region,
		Resources: []string{taskARN1},
		Detail:    &taskDetail,
	}

	containerModel := models.TaskContainer{
		ContainerARN: &containerARN1,
		LastStatus:   &taskStatus1,
		Name:         &taskName,
	}
	containerOverridesModel := models.TaskContainerOverride{
		Name: &taskName,
	}
	overridesModel := models.TaskOverride{
		ContainerOverrides: []*models.TaskContainerOverride{&containerOverridesModel},
	}
	suite.extTask = models.Task{
		ClusterARN:           &clusterARN1,
		ContainerInstanceARN: &instanceARN1,
		Containers:           []*models.TaskContainer{&containerModel},
		CreatedAt:            &createdAt,
		DesiredStatus:        &taskStatus1,
		LastStatus:           &taskStatus1,
		Overrides:            &overridesModel,
		TaskARN:              &taskARN1,
		TaskDefinitionARN:    &taskDefinitionARN,
	}
}

func TestTranslateTestSuite(t *testing.T) {
	suite.Run(t, new(TranslateTestSuite))
}

func (suite *TranslateTestSuite) TestToContainerInstance() {
	translatedModel, err := ToContainerInstance(suite.instance)
	assert.Nil(suite.T(), err, "Unexpected error when translating container instance")
	assert.Equal(suite.T(), suite.extInstance, translatedModel, "Translated model does not match expected model")
}

func (suite *TranslateTestSuite) TestToContainerInstanceEmptyDetail() {
	instance := suite.instance
	instance.Detail = nil
	_, err := ToContainerInstance(instance)
	assert.NotNil(suite.T(), err, "Expected error when translating container instance with empty detail")
}

func (suite *TranslateTestSuite) TestToContainerInstanceEmptyAgentConnected() {
	instance := suite.instance
	instance.Detail.AgentConnected = nil
	_, err := ToContainerInstance(instance)
	assert.NotNil(suite.T(), err, "Expected error when translating container instance with empty agent connected")
}

func (suite *TranslateTestSuite) TestToContainerInstanceEmptyAttributes() {
	instance := suite.instance
	instance.Detail.Attributes = nil
	translatedModel, err := ToContainerInstance(instance)

	assert.Nil(suite.T(), err, "Unexpected error when translating container instance with empty attributes")
	expectedModel := suite.extInstance
	expectedModel.Attributes = nil
	assert.Equal(suite.T(), expectedModel, translatedModel, "Translated model does not match expected model with empty attributes")
}

func (suite *TranslateTestSuite) TestToContainerInstanceEmptyClusterARN() {
	instance := suite.instance
	instance.Detail.ClusterARN = nil
	_, err := ToContainerInstance(instance)
	assert.NotNil(suite.T(), err, "Expected error when translating container instance with empty cluster ARN")
}

func (suite *TranslateTestSuite) TestToContainerInstanceEmptyContainerInstanceARN() {
	instance := suite.instance
	instance.Detail.ContainerInstanceARN = nil
	_, err := ToContainerInstance(instance)
	assert.NotNil(suite.T(), err, "Expected error when translating container instance with empty instance ARN")
}

func (suite *TranslateTestSuite) TestToContainerInstanceEmptyRegisteredResources() {
	instance := suite.instance
	instance.Detail.RegisteredResources = nil
	_, err := ToContainerInstance(instance)
	assert.NotNil(suite.T(), err, "Expected error when translating container instance with empty registered resources")
}

func (suite *TranslateTestSuite) TestToContainerInstanceEmptyRemainingResources() {
	instance := suite.instance
	instance.Detail.RemainingResources = nil
	_, err := ToContainerInstance(instance)
	assert.NotNil(suite.T(), err, "Expected error when translating container instance with empty remaining resources")
}

func (suite *TranslateTestSuite) TestToContainerInstanceEmptyStatus() {
	instance := suite.instance
	instance.Detail.Status = nil
	_, err := ToContainerInstance(instance)
	assert.NotNil(suite.T(), err, "Expected error when translating container instance with empty status")
}

func (suite *TranslateTestSuite) TestToContainerInstanceEmptyVersionInfo() {
	instance := suite.instance
	instance.Detail.VersionInfo = nil
	_, err := ToContainerInstance(instance)
	assert.NotNil(suite.T(), err, "Expected error when translating container instance with empty version info")
}

func (suite *TranslateTestSuite) TestToTask() {
	translatedModel, err := ToTask(suite.task)
	assert.Nil(suite.T(), err, "Unexpected error when translating container instance")
	assert.Equal(suite.T(), suite.extTask, translatedModel, "Translated model does not match expected model")
}

func (suite *TranslateTestSuite) TestToTaskEmptyDetail() {
	task := suite.task
	task.Detail = nil
	_, err := ToTask(task)
	assert.NotNil(suite.T(), err, "Expected error when translating task with empty detail")
}

func (suite *TranslateTestSuite) TestToTaskEmptyClusterARN() {
	task := suite.task
	task.Detail.ClusterARN = nil
	_, err := ToTask(task)
	assert.NotNil(suite.T(), err, "Expected error when translating task with empty cluster ARN")
}

func (suite *TranslateTestSuite) TestToTaskEmptyContainerInstanceARN() {
	task := suite.task
	task.Detail.ContainerInstanceARN = nil
	_, err := ToTask(task)
	assert.NotNil(suite.T(), err, "Expected error when translating task with empty container instance ARN")
}

func (suite *TranslateTestSuite) TestToTaskEmptyContainers() {
	task := suite.task
	task.Detail.Containers = nil
	_, err := ToTask(task)
	assert.NotNil(suite.T(), err, "Expected error when translating task with empty containers")
}

func (suite *TranslateTestSuite) TestToTaskEmptyCreatedAt() {
	task := suite.task
	task.Detail.CreatedAt = nil
	_, err := ToTask(task)
	assert.NotNil(suite.T(), err, "Expected error when translating task with empty created at")
}

func (suite *TranslateTestSuite) TestToTaskEmptyDesiredStatus() {
	task := suite.task
	task.Detail.DesiredStatus = nil
	_, err := ToTask(task)
	assert.NotNil(suite.T(), err, "Expected error when translating task with empty desired status")
}

func (suite *TranslateTestSuite) TestToTaskEmptyLastStatus() {
	task := suite.task
	task.Detail.LastStatus = nil
	_, err := ToTask(task)
	assert.NotNil(suite.T(), err, "Expected error when translating task with empty last status")
}

func (suite *TranslateTestSuite) TestToTaskEmptyOverrides() {
	task := suite.task
	task.Detail.Overrides = nil
	_, err := ToTask(task)
	assert.NotNil(suite.T(), err, "Expected error when translating task with empty overrides")
}

func (suite *TranslateTestSuite) TestToTaskEmptyTaskARN() {
	task := suite.task
	task.Detail.TaskARN = nil
	_, err := ToTask(task)
	assert.NotNil(suite.T(), err, "Expected error when translating task with empty task ARN")
}

func (suite *TranslateTestSuite) TestToTaskEmptyTaskDefinitionARN() {
	task := suite.task
	task.Detail.TaskDefinitionARN = nil
	_, err := ToTask(task)
	assert.NotNil(suite.T(), err, "Expected error when translating task with empty task definition ARN")
}
