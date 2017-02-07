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
	"testing"

	storetypes "github.com/blox/blox/cluster-state-service/handler/store/types"
	"github.com/blox/blox/cluster-state-service/handler/types"
	"github.com/blox/blox/cluster-state-service/swagger/v1/generated/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	agentUpdateStatus = "pending"
	ecsInstanceID     = "i-12345678"
)

var (
	attributeName         = "Name"
	attributeVal          = "com.amazonaws.ecs.capability.privileged-container"
	resourceName          = "CPU"
	intResourceType       = "INTEGER"
	longResourceType      = "LONG"
	doubleResourceType    = "DOUBLE"
	stringSetResourceType = "STRINGSET"
	intResourceVal        = int64(1024)
	intResourceValStr     = "1024"
)

type TranslateTestSuite struct {
	suite.Suite
	instance          types.ContainerInstance
	versionedInstance storetypes.VersionedContainerInstance
	extInstance       models.ContainerInstance
	task              types.Task
	versionedTask     storetypes.VersionedTask
	extTask           models.Task
}

func (suite *TranslateTestSuite) SetupTest() {
	versionInfo := types.VersionInfo{}
	attribute := types.Attribute{
		Name:  &attributeName,
		Value: &attributeVal,
	}
	resource := types.Resource{
		Name:         &resourceName,
		Type:         &intResourceType,
		IntegerValue: &intResourceVal,
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
	suite.versionedInstance = storetypes.VersionedContainerInstance{
		ContainerInstance: suite.instance,
		Version: entityVersion,
	}

	versionInfoModel := models.ContainerInstanceVersionInfo{}
	attributeModel := models.ContainerInstanceAttribute{
		Name:  &attributeName,
		Value: &attributeVal,
	}
	extResource := models.ContainerInstanceResource{
		Name:  &resourceName,
		Type:  &intResourceType,
		Value: &intResourceValStr,
	}
	suite.extInstance = models.ContainerInstance{
		Metadata: &models.Metadata{
			EntityVersion: &entityVersion,
		},
		Entity: &models.ContainerInstanceDetail{
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
		},
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
	suite.versionedTask = storetypes.VersionedTask{
		Task: suite.task,
		Version: entityVersion,
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
		Metadata: &models.Metadata{
			EntityVersion: &entityVersion,
		},
		Entity: &models.TaskDetail{
			ClusterARN:           &clusterARN1,
			ContainerInstanceARN: &instanceARN1,
			Containers:           []*models.TaskContainer{&containerModel},
			CreatedAt:            &createdAt,
			DesiredStatus:        &taskStatus1,
			LastStatus:           &taskStatus1,
			Overrides:            &overridesModel,
			TaskARN:              &taskARN1,
			TaskDefinitionARN:    &taskDefinitionARN,
		},
	}
}

func TestTranslateTestSuite(t *testing.T) {
	suite.Run(t, new(TranslateTestSuite))
}

func (suite *TranslateTestSuite) TestToContainerInstanceIntResourceType() {
	translatedModel, err := ToContainerInstance(suite.versionedInstance)
	assert.Nil(suite.T(), err, "Unexpected error when translating container instance")
	assert.Equal(suite.T(), suite.extInstance, translatedModel, "Translated model does not match expected model")
}

func (suite *TranslateTestSuite) TestToContainerInstanceLongResourceType() {
	longResourceVal := int64(1024)
	longResourceValStr := "1024"
	versionedInstance := suite.versionedInstance
	resource := types.Resource{
		Name:      &resourceName,
		Type:      &longResourceType,
		LongValue: &longResourceVal,
	}
	versionedInstance.ContainerInstance.Detail.RegisteredResources = []*types.Resource{&resource}
	extInstance := suite.extInstance
	extResource := models.ContainerInstanceResource{
		Name:  &resourceName,
		Type:  &longResourceType,
		Value: &longResourceValStr,
	}
	extInstance.Entity.RegisteredResources = []*models.ContainerInstanceResource{&extResource}
	translatedModel, err := ToContainerInstance(versionedInstance)
	assert.Nil(suite.T(), err, "Unexpected error when translating container instance")
	assert.Equal(suite.T(), extInstance, translatedModel, "Translated model does not match expected model")
}

func (suite *TranslateTestSuite) TestToContainerInstanceDoubleResourceType() {
	doubleResourceVal := 10.30
	doubleResourceValStr := "10.30"
	versionedInstance := suite.versionedInstance
	resource := types.Resource{
		Name:        &resourceName,
		Type:        &doubleResourceType,
		DoubleValue: &doubleResourceVal,
	}
	versionedInstance.ContainerInstance.Detail.RegisteredResources = []*types.Resource{&resource}
	extInstance := suite.extInstance
	extResource := models.ContainerInstanceResource{
		Name:  &resourceName,
		Type:  &doubleResourceType,
		Value: &doubleResourceValStr,
	}
	extInstance.Entity.RegisteredResources = []*models.ContainerInstanceResource{&extResource}
	translatedModel, err := ToContainerInstance(versionedInstance)
	assert.Nil(suite.T(), err, "Unexpected error when translating container instance")
	assert.Equal(suite.T(), extInstance, translatedModel, "Translated model does not match expected model")
}

func (suite *TranslateTestSuite) TestToContainerInstanceStringSetResourceType() {
	str1 := "2376"
	str2 := "22"
	stringSetResourceVal := []*string{&str1, &str2}
	stringSetResourceValStr := str1 + "," + str2
	versionedInstance := suite.versionedInstance
	resource := types.Resource{
		Name:           &resourceName,
		Type:           &stringSetResourceType,
		StringSetValue: stringSetResourceVal,
	}
	versionedInstance.ContainerInstance.Detail.RegisteredResources = []*types.Resource{&resource}
	extInstance := suite.extInstance
	extResource := models.ContainerInstanceResource{
		Name:  &resourceName,
		Type:  &stringSetResourceType,
		Value: &stringSetResourceValStr,
	}
	extInstance.Entity.RegisteredResources = []*models.ContainerInstanceResource{&extResource}
	translatedModel, err := ToContainerInstance(versionedInstance)
	assert.Nil(suite.T(), err, "Unexpected error when translating container instance")
	assert.Equal(suite.T(), extInstance, translatedModel, "Translated model does not match expected model")
}

func (suite *TranslateTestSuite) TestToContainerInstanceEmptyDetail() {
	versionedInstance := suite.versionedInstance
	versionedInstance.ContainerInstance.Detail = nil
	_, err := ToContainerInstance(versionedInstance)
	assert.NotNil(suite.T(), err, "Expected error when translating container instance with empty detail")
}

func (suite *TranslateTestSuite) TestToContainerInstanceEmptyAgentConnected() {
	versionedInstance := suite.versionedInstance
	versionedInstance.ContainerInstance.Detail.AgentConnected = nil
	_, err := ToContainerInstance(versionedInstance)
	assert.NotNil(suite.T(), err, "Expected error when translating container instance with empty agent connected")
}

func (suite *TranslateTestSuite) TestToContainerInstanceEmptyAttributes() {
	versionedInstance := suite.versionedInstance
	versionedInstance.ContainerInstance.Detail.Attributes = nil
	translatedModel, err := ToContainerInstance(versionedInstance)

	assert.Nil(suite.T(), err, "Unexpected error when translating container instance with empty attributes")
	expectedModel := suite.extInstance
	expectedModel.Entity.Attributes = nil
	assert.Equal(suite.T(), expectedModel, translatedModel, "Translated model does not match expected model with empty attributes")
}

func (suite *TranslateTestSuite) TestToContainerInstanceEmptyClusterARN() {
	versionedInstance := suite.versionedInstance
	versionedInstance.ContainerInstance.Detail.ClusterARN = nil
	_, err := ToContainerInstance(versionedInstance)
	assert.NotNil(suite.T(), err, "Expected error when translating container instance with empty cluster ARN")
}

func (suite *TranslateTestSuite) TestToContainerInstanceEmptyContainerInstanceARN() {
	versionedInstance := suite.versionedInstance
	versionedInstance.ContainerInstance.Detail.ContainerInstanceARN = nil
	_, err := ToContainerInstance(versionedInstance)
	assert.NotNil(suite.T(), err, "Expected error when translating container instance with empty instance ARN")
}

func (suite *TranslateTestSuite) TestToContainerInstanceEmptyRegisteredResources() {
	versionedInstance := suite.versionedInstance
	versionedInstance.ContainerInstance.Detail.RegisteredResources = nil
	_, err := ToContainerInstance(versionedInstance)
	assert.NotNil(suite.T(), err, "Expected error when translating container instance with empty registered resources")
}

func (suite *TranslateTestSuite) TestToContainerInstanceEmptyRemainingResources() {
	versionedInstance := suite.versionedInstance
	versionedInstance.ContainerInstance.Detail.RemainingResources = nil
	_, err := ToContainerInstance(versionedInstance)
	assert.NotNil(suite.T(), err, "Expected error when translating container instance with empty remaining resources")
}

func (suite *TranslateTestSuite) TestToContainerInstanceEmptyStatus() {
	versionedInstance := suite.versionedInstance
	versionedInstance.ContainerInstance.Detail.Status = nil
	_, err := ToContainerInstance(versionedInstance)
	assert.NotNil(suite.T(), err, "Expected error when translating container instance with empty status")
}

func (suite *TranslateTestSuite) TestToContainerInstanceEmptyVersionInfo() {
	versionedInstance := suite.versionedInstance
	versionedInstance.ContainerInstance.Detail.VersionInfo = nil
	_, err := ToContainerInstance(versionedInstance)
	assert.NotNil(suite.T(), err, "Expected error when translating container instance with empty version info")
}

func (suite *TranslateTestSuite) TestToTask() {
	translatedModel, err := ToTask(suite.versionedTask)
	assert.Nil(suite.T(), err, "Unexpected error when translating container instance")
	assert.Equal(suite.T(), suite.extTask, translatedModel, "Translated model does not match expected model")
}

func (suite *TranslateTestSuite) TestToTaskEmptyDetail() {
	versionedTask := suite.versionedTask
	versionedTask.Task.Detail = nil
	_, err := ToTask(versionedTask)
	assert.NotNil(suite.T(), err, "Expected error when translating task with empty detail")
}

func (suite *TranslateTestSuite) TestToTaskEmptyClusterARN() {
	versionedTask := suite.versionedTask
	versionedTask.Task.Detail.ClusterARN = nil
	_, err := ToTask(versionedTask)
	assert.NotNil(suite.T(), err, "Expected error when translating task with empty cluster ARN")
}

func (suite *TranslateTestSuite) TestToTaskEmptyContainerInstanceARN() {
	versionedTask := suite.versionedTask
	versionedTask.Task.Detail.ContainerInstanceARN = nil
	_, err := ToTask(versionedTask)
	assert.NotNil(suite.T(), err, "Expected error when translating task with empty container instance ARN")
}

func (suite *TranslateTestSuite) TestToTaskEmptyContainers() {
	versionedTask := suite.versionedTask
	versionedTask.Task.Detail.Containers = nil
	_, err := ToTask(versionedTask)
	assert.NotNil(suite.T(), err, "Expected error when translating task with empty containers")
}

func (suite *TranslateTestSuite) TestToTaskEmptyCreatedAt() {
	versionedTask := suite.versionedTask
	versionedTask.Task.Detail.CreatedAt = nil
	_, err := ToTask(versionedTask)
	assert.NotNil(suite.T(), err, "Expected error when translating task with empty created at")
}

func (suite *TranslateTestSuite) TestToTaskEmptyDesiredStatus() {
	versionedTask := suite.versionedTask
	versionedTask.Task.Detail.DesiredStatus = nil
	_, err := ToTask(versionedTask)
	assert.NotNil(suite.T(), err, "Expected error when translating task with empty desired status")
}

func (suite *TranslateTestSuite) TestToTaskEmptyLastStatus() {
	versionedTask := suite.versionedTask
	versionedTask.Task.Detail.LastStatus = nil
	_, err := ToTask(versionedTask)
	assert.NotNil(suite.T(), err, "Expected error when translating task with empty last status")
}

func (suite *TranslateTestSuite) TestToTaskEmptyOverrides() {
	versionedTask := suite.versionedTask
	versionedTask.Task.Detail.Overrides = nil
	_, err := ToTask(versionedTask)
	assert.NotNil(suite.T(), err, "Expected error when translating task with empty overrides")
}

func (suite *TranslateTestSuite) TestToTaskEmptyTaskARN() {
	versionedTask := suite.versionedTask
	versionedTask.Task.Detail.TaskARN = nil
	_, err := ToTask(versionedTask)
	assert.NotNil(suite.T(), err, "Expected error when translating task with empty task ARN")
}

func (suite *TranslateTestSuite) TestToTaskEmptyTaskDefinitionARN() {
	versionedTask := suite.versionedTask
	versionedTask.Task.Detail.TaskDefinitionARN = nil
	_, err := ToTask(versionedTask)
	assert.NotNil(suite.T(), err, "Expected error when translating task with empty task definition ARN")
}
