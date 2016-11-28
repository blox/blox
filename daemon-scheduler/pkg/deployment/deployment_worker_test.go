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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/blox/blox/daemon-scheduler/pkg/mocks"
	"github.com/blox/blox/daemon-scheduler/pkg/types"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type DeploymentWorkerTestSuite struct {
	suite.Suite
	environment           *mocks.MockEnvironment
	clusterState          *mocks.MockClusterState
	ecs                   *mocks.MockECS
	deploymentEnvironment *types.Environment
	deployment            *types.Deployment
	clusterTaskARNs       []*string
	deploymentWorker      DeploymentWorker
	ctx                   context.Context
}

func (suite *DeploymentWorkerTestSuite) SetupTest() {
	mockCtrl := gomock.NewController(suite.T())

	suite.environment = mocks.NewMockEnvironment(mockCtrl)
	suite.clusterState = mocks.NewMockClusterState(mockCtrl)
	suite.ecs = mocks.NewMockECS(mockCtrl)
	suite.deploymentWorker = NewDeploymentWorker(suite.environment, suite.ecs, suite.clusterState)

	suite.clusterTaskARNs = []*string{aws.String(taskARN1), aws.String(taskARN2)}

	var err error
	suite.deploymentEnvironment, err = types.NewEnvironment(environmentName, taskDefinition, cluster)
	assert.Nil(suite.T(), err, "Cannot initialize DeploymentWorkerTestSuite")

	suite.deployment, err = types.NewDeployment(taskDefinition, suite.deploymentEnvironment.Token)
	assert.Nil(suite.T(), err, "Cannot initialize DeploymentWorkerTestSuite")
	assert.NotNil(suite.T(), suite.deployment, "Cannot initialize DeploymentWorkerTestSuite")

	suite.deployment, err = suite.deployment.UpdateDeploymentInProgress(0, nil)
	assert.Nil(suite.T(), err, "Cannot initialize DeploymentWorkerTestSuite")

	suite.deploymentEnvironment.Deployments[suite.deployment.ID] = *suite.deployment
	suite.deploymentEnvironment.InProgressDeploymentID = suite.deployment.ID

	suite.ctx = context.TODO()
}

func TestDeploymentWorkerTestSuite(t *testing.T) {
	suite.Run(t, new(DeploymentWorkerTestSuite))
}

func (suite *DeploymentWorkerTestSuite) TestNewDeploymentWorker() {
	w := NewDeploymentWorker(suite.environment, suite.ecs, suite.clusterState)
	assert.NotNil(suite.T(), w, "Worker should not be nil")
}

func (suite *DeploymentWorkerTestSuite) TestUpdateInProgressDeploymentEmptyEnvironmentName() {
	_, err := suite.deploymentWorker.UpdateInProgressDeployment(suite.ctx, "")
	assert.Error(suite.T(), err, "Expected an error when env name is missing")
}

func (suite *DeploymentWorkerTestSuite) TestUpdateInProgressDeploymentGetEnvironmentFails() {
	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).
		Return(nil, errors.New("Get environment failed"))

	_, err := suite.deploymentWorker.UpdateInProgressDeployment(suite.ctx, environmentName)
	assert.Error(suite.T(), err, "Expected an error when get environment fails")
}

func (suite *DeploymentWorkerTestSuite) TestUpdateInProgressDeploymentGetEnvironmentEmpty() {
	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).
		Return(nil, nil)

	d, err := suite.deploymentWorker.UpdateInProgressDeployment(suite.ctx, environmentName)
	assert.Nil(suite.T(), err, "Unexpected error when environment is missing")
	assert.Nil(suite.T(), d, "Deployment should be nil when environment is missing")
}

func (suite *DeploymentWorkerTestSuite) TestUpdateInProgressDeploymentNoInProgressDeployment() {
	suite.deploymentEnvironment.Deployments = nil
	suite.deploymentEnvironment.InProgressDeploymentID = ""

	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).
		Return(suite.deploymentEnvironment, nil)

	d, err := suite.deploymentWorker.UpdateInProgressDeployment(suite.ctx, environmentName)
	assert.Nil(suite.T(), err, "Unexpected error when there is no in progress deployment")
	assert.Nil(suite.T(), d, "Deployment should be nil when there is no in progress deployment")
}

func (suite *DeploymentWorkerTestSuite) TestUpdateInProgressDeploymentListTasksFails() {
	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).
		Return(suite.deploymentEnvironment, nil)
	suite.ecs.EXPECT().ListTasks(suite.deploymentEnvironment.Cluster, suite.deployment.ID, "").
		Return(nil, errors.New("ListTasks failed"))

	_, err := suite.deploymentWorker.UpdateInProgressDeployment(suite.ctx, environmentName)
	assert.Error(suite.T(), err, "Expected an error when list tasks fails")
}

func (suite *DeploymentWorkerTestSuite) TestUpdateInProgressDeploymentDescribeTasksFails() {
	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).
		Return(suite.deploymentEnvironment, nil)
	suite.ecs.EXPECT().ListTasks(suite.deploymentEnvironment.Cluster, suite.deployment.ID, "").
		Return(suite.clusterTaskARNs, nil)
	suite.ecs.EXPECT().DescribeTasks(suite.deploymentEnvironment.Cluster, suite.clusterTaskARNs).
		Return(nil, errors.New("DescribeTasks failed"))

	_, err := suite.deploymentWorker.UpdateInProgressDeployment(suite.ctx, environmentName)
	assert.Error(suite.T(), err, "Expected an error when describe tasks fails")
}
