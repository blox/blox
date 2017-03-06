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

package deployment

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/blox/blox/daemon-scheduler/pkg/mocks"
	"github.com/blox/blox/daemon-scheduler/pkg/types"
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

	failedTask := ecs.Failure{
		Arn: aws.String(instanceARN),
	}

	suite.environment1, err = types.NewEnvironment(environmentName1, taskDefinition, cluster1)
	assert.Nil(suite.T(), err, "Cannot initialize EnvironmentTestSuite")

	suite.environment2, err = types.NewEnvironment(environmentName2, taskDefinition, cluster2)
	assert.Nil(suite.T(), err, "Cannot initialize EnvironmentTestSuite")

	suite.deployment, err = types.NewDeployment(taskDefinition, suite.environment1.Token)
	assert.Nil(suite.T(), err, "Cannot initialize EnvironmentTestSuite")

	suite.updatedDeployment, err = types.NewDeployment(taskDefinition, suite.environment1.Token)
	assert.Nil(suite.T(), err, "Cannot initialize EnvironmentTestSuite")
	*suite.updatedDeployment = *suite.deployment
	err = suite.updatedDeployment.UpdateDeploymentToInProgress(desiredTaskCount, nil)
	assert.Nil(suite.T(), err, "Cannot initialize EnvironmentTestSuite")

	suite.unhealthyDeployment, err = types.NewDeployment(taskDefinition, suite.environment1.Token)
	assert.Nil(suite.T(), err, "Cannot initialize EnvironmentTestSuite")
	*suite.unhealthyDeployment = *suite.deployment
	err = suite.unhealthyDeployment.UpdateDeploymentToInProgress(desiredTaskCount, []*ecs.Failure{&failedTask})
	assert.Nil(suite.T(), err, "Cannot initialize EnvironmentTestSuite")

	suite.updatedEnvironment, err = types.NewEnvironment(environmentName1, taskDefinition, cluster1)
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
	_, err := suite.environment.CreateEnvironment(suite.ctx, "", taskDefinition, cluster1)
	assert.Error(suite.T(), err, "Expected an error when name is empty")
}

func (suite *EnvironmentTestSuite) TestCreateEnvironmentEmptyTaskDefinition() {
	_, err := suite.environment.CreateEnvironment(suite.ctx, environmentName1, "", cluster1)
	assert.Error(suite.T(), err, "Expected an error when taskDefinition is empty")
}

func (suite *EnvironmentTestSuite) TestCreateEnvironmentEmptyCluster() {
	_, err := suite.environment.CreateEnvironment(suite.ctx, environmentName1, taskDefinition, "")
	assert.Error(suite.T(), err, "Expected an error when cluster is empty")
}

func (suite *EnvironmentTestSuite) TestCreateEnvironmentPutEnvironmentTxFails() {
	suite.environmentStore.EXPECT().PutEnvironment(suite.ctx, environmentName1, gomock.Any()).
		Return(errors.New("Put environment failed"))

	_, err := suite.environment.CreateEnvironment(suite.ctx, environmentName1, taskDefinition, cluster1)
	assert.Error(suite.T(), err, "Expected an error when put environment fails")
}

func (suite *EnvironmentTestSuite) TestCreateEnvironment() {
	suite.environmentStore.EXPECT().PutEnvironment(suite.ctx, environmentName1, gomock.Any()).
		Return(nil)

	env, err := suite.environment.CreateEnvironment(suite.ctx, environmentName1, taskDefinition, cluster1)
	assert.Nil(suite.T(), err, "Unexpected error when creating an environment")
	verifyEnvironment(suite.T(), suite.environment1, env)
}

// Testing validateAndCreateEnvironment used by CreateEnvironment - environment already exists
func (suite *EnvironmentTestSuite) TestValidateAndCreateEnvironmentEnvironmentExists() {
	validateAndCreate := suite.environment.ValidateAndCreateEnvironment(suite.environment1)
	_, err := validateAndCreate(suite.environment2)
	assert.Error(suite.T(), err, "Expected an error when environment exists")
}

// Testing validateAndCreateEnvironment used by CreateEnvironment - successful creation
func (suite *EnvironmentTestSuite) TestValidateAndCreateEnvironment() {
	validateAndCreate := suite.environment.ValidateAndCreateEnvironment(suite.environment1)
	env, err := validateAndCreate(nil)
	assert.Nil(suite.T(), err, "Unexpected error while creating an environment when there is no existing environment")
	assert.Exactly(suite.T(), suite.environment1, env, "Invalid new environment created")
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

func (suite *EnvironmentTestSuite) TestDeleteEnvironment() {
	environment, err := types.NewEnvironment("TestDeleteEnvironment", taskDefinition, cluster1)
	assert.Nil(suite.T(), err, "Unexpected error when creating environment")
	suite.environmentStore.EXPECT().DeleteEnvironment(suite.ctx, environment.Name).
		Return(nil).Times(1)

	err = suite.environment.DeleteEnvironment(suite.ctx, environment.Name)
	assert.Nil(suite.T(), err, "Unexpected error when deleting environment")
}

func (suite *EnvironmentTestSuite) TestDeleteEnvironmentReturnsError() {
	environment, err := types.NewEnvironment("TestDeleteEnvironmentReturnsError", taskDefinition, cluster1)
	assert.Nil(suite.T(), err, "Unexpected error when creating environment")
	err = errors.New("Error calling DeleteEnvironment")
	suite.environmentStore.EXPECT().DeleteEnvironment(suite.ctx, environment.Name).
		Return(err).Times(1)

	observedErr := suite.environment.DeleteEnvironment(suite.ctx, environment.Name)
	assert.Exactly(suite.T(), err, errors.Cause(observedErr))
}

func (suite *EnvironmentTestSuite) TestDeleteEnvironmentEmptyName() {
	suite.environmentStore.EXPECT().DeleteEnvironment(suite.ctx, gomock.Any()).
		Times(0)

	_, ok := suite.environment.DeleteEnvironment(suite.ctx, "").(types.BadRequestError)
	assert.True(suite.T(), ok, "Expecting BadRequestError when deleting environment with empty name")
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

func (suite *EnvironmentTestSuite) TestFilterEnvironmentsEmptyFilterKey() {
	suite.environmentStore.EXPECT().ListEnvironments(gomock.Any()).Times(0)

	_, err := suite.environment.FilterEnvironments(suite.ctx, "", "filterVal")
	assert.Error(suite.T(), err, "Expected an error when filter key is empty")
}

func (suite *EnvironmentTestSuite) TestFilterEnvironmentsEmptyFilterVal() {
	suite.environmentStore.EXPECT().ListEnvironments(gomock.Any()).Times(0)

	_, err := suite.environment.FilterEnvironments(suite.ctx, clusterFilter, "")
	assert.Error(suite.T(), err, "Expected an error when filter val is empty")
}

func (suite *EnvironmentTestSuite) TestFilterEnvironmentsUnsupportedFilterKey() {
	suite.environmentStore.EXPECT().ListEnvironments(gomock.Any()).Times(0)

	_, err := suite.environment.FilterEnvironments(suite.ctx, "unsupportedFilter", "filterVal")
	assert.Error(suite.T(), err, "Expected an error when filter key is unsupported")
}

func (suite *EnvironmentTestSuite) TestFilterEnvironmentsListEnvironmentsReturnsError() {
	suite.environmentStore.EXPECT().ListEnvironments(gomock.Any()).
		Return(nil, errors.New("Error listing environments"))

	_, err := suite.environment.FilterEnvironments(suite.ctx, clusterFilter, "filterVal")
	assert.Error(suite.T(), err, "Expected an error filtering environments when store list environments returns an error")
}

func (suite *EnvironmentTestSuite) TestFilterEnvironmentsByClusterARN() {
	suite.filterEnvironmentsByCluster(cluster1)
}

func (suite *EnvironmentTestSuite) TestFilterEnvironmentsByClusterName() {
	suite.filterEnvironmentsByCluster(clusterName1)
}

func (suite *EnvironmentTestSuite) filterEnvironmentsByCluster(cluster string) {
	cluster1Env1 := suite.environment1
	cluster2Env1 := suite.environment2
	cluster1Env2, err := types.NewEnvironment(environmentName3, taskDefinition, cluster1)
	assert.Nil(suite.T(), err, "Unexpected error creating new environment")

	allEnvs := []types.Environment{*cluster1Env1, *cluster2Env1, *cluster1Env2}
	suite.environmentStore.EXPECT().ListEnvironments(suite.ctx).
		Return(allEnvs, nil)

	envs, err := suite.environment.FilterEnvironments(suite.ctx, clusterFilter, cluster)
	assert.Nil(suite.T(), err, "Unexpected error when filtering environments")
	expectedEnvs := []types.Environment{*cluster1Env1, *cluster1Env2}
	assert.Exactly(suite.T(), expectedEnvs, envs, "Returned filtered environments does not match expected environments")
}

func (suite *EnvironmentTestSuite) TestFilterEnvironmentsByClusterInvalidCluster() {
	suite.environmentStore.EXPECT().ListEnvironments(gomock.Any()).Times(0)

	invalidCluster := "cluster/cluster"
	_, err := suite.environment.FilterEnvironments(suite.ctx, clusterFilter, invalidCluster)
	assert.Error(suite.T(), err, "Expected an error when filtering using invalid cluster")
}

func verifyEnvironment(t *testing.T, expected *types.Environment, actual *types.Environment) {
	assert.NotNil(t, actual, "Environment should not be nil")
	assert.Exactly(t, expected.Name, actual.Name, "Name should match")
	assert.Exactly(t, expected.DesiredTaskDefinition, actual.DesiredTaskDefinition, "Task definition should match")
	assert.Exactly(t, expected.DesiredTaskCount, actual.DesiredTaskCount, "Desired task count should match")
	assert.Exactly(t, expected.Cluster, actual.Cluster, "Cluster should match")
	assert.Exactly(t, expected.Health, actual.Health, "Health should match")
	assert.Exactly(t, expected.Deployments, actual.Deployments, "Deployments should match")
}
