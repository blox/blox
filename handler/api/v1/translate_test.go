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
	"testing"

	"github.com/aws/amazon-ecs-event-stream-handler/handler/api/v1/models"
	"github.com/aws/amazon-ecs-event-stream-handler/handler/types"
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
	instance      types.ContainerInstance
	instanceModel models.ContainerInstanceModel
	task          types.Task
	taskModel     models.TaskModel
}

func (suite *TranslateTestSuite) SetupTest() {
	versionInfo := types.VersionInfo{}
	attribute := types.Attribute{
		Name:  &attributeName,
		Value: &attributeVal,
	}
	resoure := types.Resource{
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
		PendingTasksCount:    &pendingTaskCount1,
		RegisteredResources:  []*types.Resource{&resoure},
		RemainingResources:   []*types.Resource{&resoure},
		RunningTasksCount:    &runningTasksCount1,
		Status:               &instanceStatus1,
		Version:              &version1,
		VersionInfo:          &versionInfo,
		UpdatedAt:            &updatedAt1,
	}
	suite.instance = types.ContainerInstance{
		ID:        &id1,
		Account:   &accountID,
		Time:      &time,
		Region:    &region,
		Resources: []string{instanceARN1},
		Detail:    &instanceDetail,
	}

	versionInfoModel := models.ContainerInstanceVersionInfoModel{}
	attributeModel := models.ContainerInstanceAttributeModel{
		Name:  &attributeName,
		Value: &attributeVal,
	}
	pendingTaskCount := int32(pendingTaskCount1)
	regResoureModel := models.ContainerInstanceRegisteredResourceModel{
		Name:  &resourceName,
		Type:  &resourceType,
		Value: &resourceVal,
	}
	remResoureModel := models.ContainerInstanceRemainingResourceModel{
		Name:  &resourceName,
		Type:  &resourceType,
		Value: &resourceVal,
	}
	runningTasksCount := int32(runningTasksCount1)
	version := int32(version1)
	suite.instanceModel = models.ContainerInstanceModel{
		AgentConnected:       &agentConnected1,
		AgentUpdateStatus:    agentUpdateStatus,
		Attributes:           []*models.ContainerInstanceAttributeModel{&attributeModel},
		ClusterARN:           &clusterARN1,
		ContainerInstanceARN: &instanceARN1,
		EC2InstanceID:        ecsInstanceID,
		PendingTasksCount:    &pendingTaskCount,
		RegisteredResources:  []*models.ContainerInstanceRegisteredResourceModel{&regResoureModel},
		RemainingResources:   []*models.ContainerInstanceRemainingResourceModel{&remResoureModel},
		RunningTasksCount:    &runningTasksCount,
		Status:               &instanceStatus1,
		Version:              &version,
		VersionInfo:          &versionInfoModel,
		UpdatedAt:            &updatedAt1,
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
		UpdatedAt:            &updatedAt1,
		Version:              &version1,
	}
	suite.task = types.Task{
		ID:        &id1,
		Account:   &accountID,
		Time:      &time,
		Region:    &region,
		Resources: []string{taskARN1},
		Detail:    &taskDetail,
	}

	containerModel := models.TaskContainerModel{
		ContainerARN: &containerARN1,
		LastStatus:   &taskStatus1,
		Name:         &taskName,
	}
	containerOverridesModel := models.TaskContainerOverrideModel{
		Name: &taskName,
	}
	overridesModel := models.TaskOverrideModel{
		ContainerOverrides: []*models.TaskContainerOverrideModel{&containerOverridesModel},
	}
	suite.taskModel = models.TaskModel{
		ClusterARN:           &clusterARN1,
		ContainerInstanceARN: &instanceARN1,
		Containers:           []*models.TaskContainerModel{&containerModel},
		CreatedAt:            &createdAt,
		DesiredStatus:        &taskStatus1,
		LastStatus:           &taskStatus1,
		Overrides:            &overridesModel,
		TaskARN:              &taskARN1,
		TaskDefinitionARN:    &taskDefinitionARN,
		UpdatedAt:            &updatedAt1,
		Version:              &version,
	}
}

func TestTranslateTestSuite(t *testing.T) {
	suite.Run(t, new(TranslateTestSuite))
}

func (suite *TranslateTestSuite) TestToContainerInstanceModel() {
	translatedModel, err := ToContainerInstanceModel(suite.instance)
	assert.Nil(suite.T(), err, "Unexpected error when translating container instance")
	assert.Equal(suite.T(), suite.instanceModel, translatedModel, "Translated model does not match expected model")
}

func (suite *TranslateTestSuite) TestToContainerInstanceModelEmptyDetail() {
	instance := suite.instance
	instance.Detail = nil
	_, err := ToContainerInstanceModel(instance)
	assert.NotNil(suite.T(), err, "Expected error when translating container instance with empty detail")
}

func (suite *TranslateTestSuite) TestToContainerInstanceModelEmptyAgentConnected() {
	instance := suite.instance
	instance.Detail.AgentConnected = nil
	_, err := ToContainerInstanceModel(instance)
	assert.NotNil(suite.T(), err, "Expected error when translating container instance with empty agent connected")
}

func (suite *TranslateTestSuite) TestToContainerInstanceModelEmptyAttributes() {
	instance := suite.instance
	instance.Detail.Attributes = nil
	translatedModel, err := ToContainerInstanceModel(instance)

	assert.Nil(suite.T(), err, "Unexpected error when translating container instance with empty attributes")
	expectedModel := suite.instanceModel
	expectedModel.Attributes = nil
	assert.Equal(suite.T(), expectedModel, translatedModel, "Translated model does not match expected model with empty attributes")
}

func (suite *TranslateTestSuite) TestToContainerInstanceModelEmptyClusterARN() {
	instance := suite.instance
	instance.Detail.ClusterARN = nil
	_, err := ToContainerInstanceModel(instance)
	assert.NotNil(suite.T(), err, "Expected error when translating container instance with empty cluster ARN")
}

func (suite *TranslateTestSuite) TestToContainerInstanceModelEmptyContainerInstanceARN() {
	instance := suite.instance
	instance.Detail.ContainerInstanceARN = nil
	_, err := ToContainerInstanceModel(instance)
	assert.NotNil(suite.T(), err, "Expected error when translating container instance with empty instance ARN")
}

func (suite *TranslateTestSuite) TestToContainerInstanceModelEmptyPendingTasksCount() {
	instance := suite.instance
	instance.Detail.PendingTasksCount = nil
	_, err := ToContainerInstanceModel(instance)
	assert.NotNil(suite.T(), err, "Expected error when translating container instance with empty pending tasks count")
}

func (suite *TranslateTestSuite) TestToContainerInstanceModelEmptyRegisteredResources() {
	instance := suite.instance
	instance.Detail.RegisteredResources = nil
	_, err := ToContainerInstanceModel(instance)
	assert.NotNil(suite.T(), err, "Expected error when translating container instance with empty registered resources")
}

func (suite *TranslateTestSuite) TestToContainerInstanceModelEmptyRemainingResources() {
	instance := suite.instance
	instance.Detail.RemainingResources = nil
	_, err := ToContainerInstanceModel(instance)
	assert.NotNil(suite.T(), err, "Expected error when translating container instance with empty remaining resources")
}

func (suite *TranslateTestSuite) TestToContainerInstanceModelEmptyRunningTasksCount() {
	instance := suite.instance
	instance.Detail.RunningTasksCount = nil
	_, err := ToContainerInstanceModel(instance)
	assert.NotNil(suite.T(), err, "Expected error when translating container instance with empty running tasks count")
}

func (suite *TranslateTestSuite) TestToContainerInstanceModelEmptyStatus() {
	instance := suite.instance
	instance.Detail.Status = nil
	_, err := ToContainerInstanceModel(instance)
	assert.NotNil(suite.T(), err, "Expected error when translating container instance with empty status")
}

func (suite *TranslateTestSuite) TestToContainerInstanceModelEmptyVersion() {
	instance := suite.instance
	instance.Detail.Version = nil
	_, err := ToContainerInstanceModel(instance)
	assert.NotNil(suite.T(), err, "Expected error when translating container instance with empty version")
}

func (suite *TranslateTestSuite) TestToContainerInstanceModelEmptyVersionInfo() {
	instance := suite.instance
	instance.Detail.VersionInfo = nil
	_, err := ToContainerInstanceModel(instance)
	assert.NotNil(suite.T(), err, "Expected error when translating container instance with empty version info")
}

func (suite *TranslateTestSuite) TestToContainerInstanceModelEmptyUpdatedAt() {
	instance := suite.instance
	instance.Detail.UpdatedAt = nil
	_, err := ToContainerInstanceModel(instance)
	assert.NotNil(suite.T(), err, "Expected error when translating container instance with empty updated at")
}

func (suite *TranslateTestSuite) TestToTaskModel() {
	translatedModel, err := ToTaskModel(suite.task)
	assert.Nil(suite.T(), err, "Unexpected error when translating container instance")
	assert.Equal(suite.T(), suite.taskModel, translatedModel, "Translated model does not match expected model")
}

func (suite *TranslateTestSuite) TestToTaskModelEmptyDetail() {
	task := suite.task
	task.Detail = nil
	_, err := ToTaskModel(task)
	assert.NotNil(suite.T(), err, "Expected error when translating task with empty detail")
}

func (suite *TranslateTestSuite) TestToTaskModelEmptyClusterARN() {
	task := suite.task
	task.Detail.ClusterARN = nil
	_, err := ToTaskModel(task)
	assert.NotNil(suite.T(), err, "Expected error when translating task with empty cluster ARN")
}

func (suite *TranslateTestSuite) TestToTaskModelEmptyContainerInstanceARN() {
	task := suite.task
	task.Detail.ContainerInstanceARN = nil
	_, err := ToTaskModel(task)
	assert.NotNil(suite.T(), err, "Expected error when translating task with empty container instance ARN")
}

func (suite *TranslateTestSuite) TestToTaskModelEmptyContainers() {
	task := suite.task
	task.Detail.Containers = nil
	_, err := ToTaskModel(task)
	assert.NotNil(suite.T(), err, "Expected error when translating task with empty containers")
}

func (suite *TranslateTestSuite) TestToTaskModelEmptyCreatedAt() {
	task := suite.task
	task.Detail.CreatedAt = nil
	_, err := ToTaskModel(task)
	assert.NotNil(suite.T(), err, "Expected error when translating task with empty created at")
}

func (suite *TranslateTestSuite) TestToTaskModelEmptyDesiredStatus() {
	task := suite.task
	task.Detail.DesiredStatus = nil
	_, err := ToTaskModel(task)
	assert.NotNil(suite.T(), err, "Expected error when translating task with empty desired status")
}

func (suite *TranslateTestSuite) TestToTaskModelEmptyLastStatus() {
	task := suite.task
	task.Detail.LastStatus = nil
	_, err := ToTaskModel(task)
	assert.NotNil(suite.T(), err, "Expected error when translating task with empty last status")
}

func (suite *TranslateTestSuite) TestToTaskModelEmptyOverrides() {
	task := suite.task
	task.Detail.Overrides = nil
	_, err := ToTaskModel(task)
	assert.NotNil(suite.T(), err, "Expected error when translating task with empty overrides")
}

func (suite *TranslateTestSuite) TestToTaskModelEmptyTaskARN() {
	task := suite.task
	task.Detail.TaskARN = nil
	_, err := ToTaskModel(task)
	assert.NotNil(suite.T(), err, "Expected error when translating task with empty task ARN")
}

func (suite *TranslateTestSuite) TestToTaskModelEmptyTaskDefinitionARN() {
	task := suite.task
	task.Detail.TaskDefinitionARN = nil
	_, err := ToTaskModel(task)
	assert.NotNil(suite.T(), err, "Expected error when translating task with empty task definition ARN")
}

func (suite *TranslateTestSuite) TestToTaskModelEmptyUpdatedAt() {
	task := suite.task
	task.Detail.UpdatedAt = nil
	_, err := ToTaskModel(task)
	assert.NotNil(suite.T(), err, "Expected error when translating task with empty updated at")
}

func (suite *TranslateTestSuite) TestToTaskModelEmptyVersion() {
	task := suite.task
	task.Detail.Version = nil
	_, err := ToTaskModel(task)
	assert.NotNil(suite.T(), err, "Expected error when translating task with empty version")
}
