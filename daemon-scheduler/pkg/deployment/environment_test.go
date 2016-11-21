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

package deployment

import (
	"context"
	"testing"

	"github.com/blox/blox/daemon-scheduler/pkg/mocks"
	"github.com/blox/blox/daemon-scheduler/pkg/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type EnvironmentTestSuite struct {
	suite.Suite
	environmentStore    *mocks.MockEnvironmentStore
	environment         Environment
	ctx                 context.Context
	environment1        *types.Environment
	updatedEnvironment  *types.Environment
	environment2        *types.Environment
	deployment          *types.Deployment
	updatedDeployment   *types.Deployment
	unhealthyDeployment *types.Deployment
	currentTasks        []*ecs.Task
	taskMap             map[string]*ecs.Task
}

func (suite *EnvironmentTestSuite) SetupTest() {
	mockCtrl := gomock.NewController(suite.T())
	suite.environmentStore = mocks.NewMockEnvironmentStore(mockCtrl)
	suite.ctx = context.TODO()

	var err error
	suite.environment, err = NewEnvironment(suite.environmentStore)
	assert.Nil(suite.T(), err, "Cannot initialize EnvironmentTestSuite")

	task1 := ecs.Task{
		TaskArn:           aws.String(taskARN1),
		TaskDefinitionArn: aws.String(taskDefinition),
	}
	task2 := ecs.Task{
		TaskArn:           aws.String(taskARN2),
		TaskDefinitionArn: aws.String(taskDefinition),
	}
	suite.currentTasks = []*ecs.Task{&task1, &task2}

	suite.taskMap = make(map[string]*ecs.Task)
	for _, v := range suite.currentTasks {
		suite.taskMap[*v.TaskArn] = v
	}

	suite.deployment, err = types.NewDeployment(taskDefinition)
	assert.Nil(suite.T(), err, "Cannot initialize EnvironmentTestSuite")

	suite.updatedDeployment, err = suite.deployment.UpdateDeploymentInProgress(
		desiredTaskCount, nil, nil)
	assert.Nil(suite.T(), err, "Cannot initialize EnvironmentTestSuite")

	failedTask := ecs.Failure{
		Arn: aws.String(instanceARN),
	}

	suite.unhealthyDeployment, err = suite.deployment.UpdateDeploymentInProgress(
		desiredTaskCount, nil, []*ecs.Failure{&failedTask})
	assert.Nil(suite.T(), err, "Cannot initialize EnvironmentTestSuite")

	suite.environment1, err = types.NewEnvironment(environmentName1, taskDefinition, cluster)
	assert.Nil(suite.T(), err, "Cannot initialize EnvironmentTestSuite")

	suite.environment2, err = types.NewEnvironment(environmentName2, taskDefinition, cluster)
	assert.Nil(suite.T(), err, "Cannot initialize EnvironmentTestSuite")

	suite.updatedEnvironment, err = types.NewEnvironment(environmentName1, taskDefinition, cluster)
	assert.Nil(suite.T(), err, "Cannot initialize EnvironmentTestSuite")

	suite.updatedEnvironment.Deployments[suite.deployment.ID] = *suite.updatedDeployment
	suite.updatedEnvironment.DesiredTaskCount = suite.updatedDeployment.DesiredTaskCount
}

func TestEnvironmentTestSuite(t *testing.T) {
	suite.Run(t, new(EnvironmentTestSuite))
}

func (suite *EnvironmentTestSuite) TestNewEnvironmentEmptyStore() {
	_, err := NewEnvironment(nil)
	assert.Error(suite.T(), err, "Expected an error when store is nil")
}

func (suite *EnvironmentTestSuite) TestNewEnvironmentStore() {
	e, err := NewEnvironment(suite.environmentStore)
	assert.Nil(suite.T(), err, "Unexpected error when store is not nil")
	assert.NotNil(suite.T(), e, "Environment should not be nil")
}

func (suite *EnvironmentTestSuite) TestCreateEnvironmentEmptyName() {
	_, err := suite.environment.CreateEnvironment(suite.ctx, "", taskDefinition, cluster)
	assert.Error(suite.T(), err, "Expected an error when name is empty")
}

func (suite *EnvironmentTestSuite) TestCreateEnvironmentEmptyTaskDefinition() {
	_, err := suite.environment.CreateEnvironment(suite.ctx, environmentName1, "", cluster)
	assert.Error(suite.T(), err, "Expected an error when taskDefinition is empty")
}

func (suite *EnvironmentTestSuite) TestCreateEnvironmentEmptyCluster() {
	_, err := suite.environment.CreateEnvironment(suite.ctx, environmentName1, taskDefinition, "")
	assert.Error(suite.T(), err, "Expected an error when cluster is empty")
}

func (suite *EnvironmentTestSuite) TestCreateEnvironmentGetEnvironmentFails() {
	suite.environmentStore.EXPECT().GetEnvironment(suite.ctx, environmentName1).
		Return(nil, errors.New("Get environment failed"))

	_, err := suite.environment.CreateEnvironment(suite.ctx, environmentName1, taskDefinition, cluster)
	assert.Error(suite.T(), err, "Expected an error when get environment fails")
}

func (suite *EnvironmentTestSuite) TestCreateEnvironmentEnvironmentExists() {
	suite.environmentStore.EXPECT().GetEnvironment(suite.ctx, environmentName1).
		Return(suite.environment1, nil)

	_, err := suite.environment.CreateEnvironment(suite.ctx, environmentName1, taskDefinition, cluster)
	assert.Error(suite.T(), err, "Expected an error when environment exists")
}

func (suite *EnvironmentTestSuite) TestCreateEnvironmentPutEnvironmentFails() {
	suite.environmentStore.EXPECT().GetEnvironment(suite.ctx, environmentName1).
		Return(nil, nil)

	suite.environmentStore.EXPECT().PutEnvironment(suite.ctx, gomock.Any()).Do(func(_ interface{}, e types.Environment) {
		verifyEnvironment(suite.T(), suite.environment1, &e)
	}).Return(errors.New("Put environment failed"))

	_, err := suite.environment.CreateEnvironment(suite.ctx, environmentName1, taskDefinition, cluster)
	assert.Error(suite.T(), err, "Expected an error when put environment fails")
}

func (suite *EnvironmentTestSuite) TestCreateEnvironment() {
	suite.environmentStore.EXPECT().GetEnvironment(suite.ctx, environmentName1).Return(nil, nil)

	suite.environmentStore.EXPECT().PutEnvironment(suite.ctx, gomock.Any()).Do(func(_ interface{}, e types.Environment) {
		verifyEnvironment(suite.T(), suite.environment1, &e)
	}).Return(nil)

	env, err := suite.environment.CreateEnvironment(suite.ctx, environmentName1, taskDefinition, cluster)
	assert.Nil(suite.T(), err, "Unexpected error when creating an environment")
	verifyEnvironment(suite.T(), suite.environment1, env)
}

func (suite *EnvironmentTestSuite) TestGetEnvironmentEmptyName() {
	_, err := suite.environment.GetEnvironment(suite.ctx, "")
	assert.Error(suite.T(), err, "Expected an error when name is empty")
}

func (suite *EnvironmentTestSuite) TestGetEnvironmentGetFromStoreFails() {
	suite.environmentStore.EXPECT().GetEnvironment(suite.ctx, environmentName1).
		Return(nil, errors.New("Get environment failed"))

	_, err := suite.environment.GetEnvironment(suite.ctx, environmentName1)
	assert.Error(suite.T(), err, "Expected an error when get from store fails")
}

func (suite *EnvironmentTestSuite) TestGetEnvironment() {
	suite.environmentStore.EXPECT().GetEnvironment(suite.ctx, environmentName1).
		Return(suite.environment1, nil)

	env, err := suite.environment.GetEnvironment(suite.ctx, environmentName1)
	assert.Nil(suite.T(), err, "Unexpected error when getting an environment")
	assert.Exactly(suite.T(), suite.environment1, env, "Expected the environment to match the expected one")
}

func (suite *EnvironmentTestSuite) TestListEnvironmentsListFromStoreFails() {
	suite.environmentStore.EXPECT().ListEnvironments(suite.ctx).
		Return(nil, errors.New("List failed"))

	_, err := suite.environment.ListEnvironments(suite.ctx)
	assert.Error(suite.T(), err, "Expected an error when list from store fails")
}

func (suite *EnvironmentTestSuite) TestListEnvironments() {
	expectedEnvs := []types.Environment{*suite.environment1, *suite.environment2}
	suite.environmentStore.EXPECT().ListEnvironments(suite.ctx).
		Return(expectedEnvs, nil)

	envs, err := suite.environment.ListEnvironments(suite.ctx)
	assert.Nil(suite.T(), err, "Unexpected error when listing environments")
	assert.Exactly(suite.T(), expectedEnvs, envs, "Expected listed environments to match what's returned from the store")
}

func (suite *EnvironmentTestSuite) TestAddDeploymentEmptyDeploymentID() {
	deployment := types.Deployment{}

	_, err := suite.environment.AddDeployment(suite.ctx, *suite.environment1, deployment)
	assert.Error(suite.T(), err, "Expected an error when deployment ID is missing")
}

func (suite *EnvironmentTestSuite) TestAddDeploymentDeploymentExists() {
	suite.environment1.Deployments[suite.deployment.ID] = *suite.deployment

	_, err := suite.environment.AddDeployment(suite.ctx, *suite.environment1, *suite.deployment)
	assert.Error(suite.T(), err, "Expected an error when deployment exists")
}

func (suite *EnvironmentTestSuite) TestAddDeploymentPutEnvironmentFails() {
	suite.environmentStore.EXPECT().PutEnvironment(suite.ctx, gomock.Any()).Do(func(_ interface{}, e types.Environment) {
		verifyEnvironment(suite.T(), suite.environment1, &e)
	}).Return(errors.New("Put environment failed"))

	_, err := suite.environment.AddDeployment(suite.ctx, *suite.environment1, *suite.deployment)
	assert.Error(suite.T(), err, "Expected an error when put environment fails")
}

func (suite *EnvironmentTestSuite) TestAddDeployment() {
	suite.environmentStore.EXPECT().PutEnvironment(suite.ctx, gomock.Eq(*suite.environment1)).Return(nil)

	env, err := suite.environment.AddDeployment(suite.ctx, *suite.environment1, *suite.deployment)
	assert.Nil(suite.T(), err, "Unexpected error when adding a deployment")

	suite.environment1.Deployments[suite.deployment.ID] = *suite.deployment
	assert.Exactly(suite.T(), suite.environment1, env, "Environment does not match the expected environment")
}

func (suite *EnvironmentTestSuite) TestUpdateDeploymentEmptyDeploymentID() {
	deployment := types.Deployment{}

	_, err := suite.environment.UpdateDeployment(suite.ctx, *suite.environment1, deployment)
	assert.Error(suite.T(), err, "Expected an error when deployment ID is missing")
}

func (suite *EnvironmentTestSuite) TestUpdateDeploymentDeploymentDoesNotExist() {
	_, err := suite.environment.UpdateDeployment(suite.ctx, *suite.environment1, *suite.deployment)
	assert.Error(suite.T(), err, "Expected an error when deployment does not exist")
}

func (suite *EnvironmentTestSuite) TestUpdateDeploymentPutEnvironmentFails() {
	suite.environment1.Deployments[suite.deployment.ID] = *suite.deployment

	suite.environmentStore.EXPECT().PutEnvironment(suite.ctx, gomock.Any()).Do(func(_ interface{}, e types.Environment) {
		verifyEnvironment(suite.T(), suite.updatedEnvironment, &e)
	}).Return(errors.New("Put environment failed"))

	_, err := suite.environment.UpdateDeployment(suite.ctx, *suite.environment1, *suite.updatedDeployment)
	assert.Error(suite.T(), err, "Expected an error when put fails")
}

func (suite *EnvironmentTestSuite) TestUpdateDeploymentUnhealthy() {
	suite.environment1.Deployments[suite.deployment.ID] = *suite.deployment

	suite.updatedEnvironment.Deployments[suite.unhealthyDeployment.ID] = *suite.unhealthyDeployment
	suite.updatedEnvironment.Health = types.EnvironmentUnhealthy

	suite.environmentStore.EXPECT().PutEnvironment(suite.ctx, gomock.Any()).Do(func(_ interface{}, e types.Environment) {
		verifyEnvironment(suite.T(), suite.updatedEnvironment, &e)
	}).Return(nil)

	env, err := suite.environment.UpdateDeployment(suite.ctx, *suite.environment1, *suite.unhealthyDeployment)
	assert.Nil(suite.T(), err, "Unexpected error when unhealthy deployment is being updated")
	verifyEnvironment(suite.T(), suite.updatedEnvironment, env)
}

func (suite *EnvironmentTestSuite) TestUpdateDeploymentInvalidCurrentTaskState() {
	suite.environment1.Deployments[suite.deployment.ID] = *suite.deployment

	//update TaskDefinition on the deployment so it doesn't match the currentTasks task def
	suite.updatedDeployment.TaskDefinition = "invalid"
	suite.updatedDeployment.CurrentTasks = suite.currentTasks
	suite.updatedEnvironment.Deployments[suite.deployment.ID] = *suite.updatedDeployment

	suite.environmentStore.EXPECT().PutEnvironment(suite.ctx, gomock.Any()).Do(func(_ interface{}, e types.Environment) {
		verifyEnvironment(suite.T(), suite.updatedEnvironment, &e)
	}).Return(nil)

	_, err := suite.environment.UpdateDeployment(suite.ctx, *suite.environment1, *suite.updatedDeployment)
	assert.Error(suite.T(), err, "Expected an error when deployment task definition does not match current task task definition")
}

func (suite *EnvironmentTestSuite) TestUpdateDeploymentEnvironmentDoesNotHaveCurrentTasks() {
	suite.environment1.Deployments[suite.deployment.ID] = *suite.deployment

	suite.updatedDeployment.CurrentTasks = suite.currentTasks
	suite.updatedEnvironment.CurrentTasks = map[string]map[string]*ecs.Task{
		suite.updatedDeployment.TaskDefinition: suite.taskMap,
	}
	suite.updatedEnvironment.Deployments[suite.deployment.ID] = *suite.updatedDeployment

	suite.environmentStore.EXPECT().PutEnvironment(suite.ctx, gomock.Any()).Do(func(_ interface{}, e types.Environment) {
		verifyEnvironment(suite.T(), suite.updatedEnvironment, &e)
	}).Return(nil)

	env, err := suite.environment.UpdateDeployment(suite.ctx, *suite.environment1, *suite.updatedDeployment)
	assert.Nil(suite.T(), err, "Unexpected error when environment does not have current tasks")
	verifyEnvironment(suite.T(), suite.updatedEnvironment, env)
}

func (suite *EnvironmentTestSuite) TestUpdateDeploymentEnvironmentCurrentTasksContainsADifferentTaskDef() {
	suite.environment1.Deployments[suite.deployment.ID] = *suite.deployment

	differentTaskDef := suite.updatedDeployment.TaskDefinition + "1"
	suite.environment1.CurrentTasks = map[string]map[string]*ecs.Task{
		differentTaskDef: suite.taskMap,
	}

	suite.updatedDeployment.CurrentTasks = suite.currentTasks

	suite.updatedEnvironment.CurrentTasks[taskDefinition] = suite.taskMap
	suite.updatedEnvironment.CurrentTasks[differentTaskDef] = suite.taskMap
	suite.updatedEnvironment.Deployments[suite.deployment.ID] = *suite.updatedDeployment

	suite.environmentStore.EXPECT().PutEnvironment(suite.ctx, gomock.Any()).Do(func(_ interface{}, e types.Environment) {
		verifyEnvironment(suite.T(), suite.updatedEnvironment, &e)
	}).Return(nil)

	env, err := suite.environment.UpdateDeployment(suite.ctx, *suite.environment1, *suite.updatedDeployment)
	assert.Nil(suite.T(), err, "Unexpected error when updating a deployment and the environment current tasks contain a different task definition")
	verifyEnvironment(suite.T(), suite.updatedEnvironment, env)
}

func (suite *EnvironmentTestSuite) TestUpdateDeploymentEnvironmentCurrentTasksContainsTaskDefWithOverlapingTasks() {
	suite.environment1.Deployments[suite.deployment.ID] = *suite.deployment

	taskMapWithOneTask := make(map[string]*ecs.Task)
	taskMapWithOneTask[*suite.currentTasks[0].TaskArn] = suite.currentTasks[0]

	suite.environment1.CurrentTasks = map[string]map[string]*ecs.Task{
		taskDefinition: taskMapWithOneTask,
	}

	suite.updatedDeployment.CurrentTasks = suite.currentTasks

	//update map to include the new taskDef -> taskMap
	suite.updatedEnvironment.CurrentTasks[taskDefinition] = suite.taskMap
	suite.updatedEnvironment.Deployments[suite.deployment.ID] = *suite.updatedDeployment

	suite.environmentStore.EXPECT().PutEnvironment(suite.ctx, gomock.Any()).Do(func(_ interface{}, e types.Environment) {
		verifyEnvironment(suite.T(), suite.updatedEnvironment, &e)
	}).Return(nil)

	env, err := suite.environment.UpdateDeployment(suite.ctx, *suite.environment1, *suite.updatedDeployment)
	assert.Nil(suite.T(), err, "Unexpected error when updating a deployment and the environment current tasks contain some of the deployment task arns")
	verifyEnvironment(suite.T(), suite.updatedEnvironment, env)
}

func verifyEnvironment(t *testing.T, expected *types.Environment, actual *types.Environment) {
	assert.NotNil(t, actual, "Environment should not be nil")
	assert.Exactly(t, expected.Name, actual.Name, "Name should match")
	assert.Exactly(t, expected.DesiredTaskDefinition, actual.DesiredTaskDefinition, "Task definition should match")
	assert.Exactly(t, expected.DesiredTaskCount, actual.DesiredTaskCount, "Desired task count should match")
	assert.Exactly(t, expected.Cluster, actual.Cluster, "Cluster should match")
	assert.Exactly(t, expected.Health, actual.Health, "Health should match")
	assert.Exactly(t, expected.CurrentTasks, actual.CurrentTasks, "Current tasks should match")
	assert.Exactly(t, expected.Deployments, actual.Deployments, "Deployments should match")
}
