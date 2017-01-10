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
	environment       *mocks.MockEnvironment
	deployment        *mocks.MockDeployment
	clusterState      *mocks.MockClusterState
	ecs               *mocks.MockECS
	environmentObject *types.Environment
	deploymentObject  *types.Deployment
	clusterTaskARNs   []*string
	deploymentWorker  DeploymentWorker
	ctx               context.Context
}

func (suite *DeploymentWorkerTestSuite) SetupTest() {
	mockCtrl := gomock.NewController(suite.T())

	suite.environment = mocks.NewMockEnvironment(mockCtrl)
	suite.deployment = mocks.NewMockDeployment(mockCtrl)
	suite.clusterState = mocks.NewMockClusterState(mockCtrl)
	suite.ecs = mocks.NewMockECS(mockCtrl)
	suite.deploymentWorker = NewDeploymentWorker(suite.environment, suite.deployment, suite.ecs, suite.clusterState)

	suite.clusterTaskARNs = []*string{aws.String(taskARN1), aws.String(taskARN2)}

	var err error
	suite.environmentObject, err = types.NewEnvironment(environmentName, taskDefinition, cluster1)
	assert.Nil(suite.T(), err, "Cannot initialize DeploymentWorkerTestSuite")

	suite.deploymentObject, err = types.NewDeployment(taskDefinition, suite.environmentObject.Token)
	assert.Nil(suite.T(), err, "Cannot initialize DeploymentWorkerTestSuite")
	assert.NotNil(suite.T(), suite.deploymentObject, "Cannot initialize DeploymentWorkerTestSuite")

	suite.deploymentObject, err = suite.deploymentObject.UpdateDeploymentInProgress(0, nil)
	assert.Nil(suite.T(), err, "Cannot initialize DeploymentWorkerTestSuite")

	suite.ctx = context.TODO()
}

func TestDeploymentWorkerTestSuite(t *testing.T) {
	suite.Run(t, new(DeploymentWorkerTestSuite))
}

func (suite *DeploymentWorkerTestSuite) TestNewDeploymentWorker() {
	w := NewDeploymentWorker(suite.environment, suite.deployment, suite.ecs, suite.clusterState)
	assert.NotNil(suite.T(), w, "Worker should not be nil")
}

func (suite *DeploymentWorkerTestSuite) TestUpdateInProgressDeploymentEmptyEnvironmentName() {
	_, err := suite.deploymentWorker.UpdateInProgressDeployment(suite.ctx, "")
	assert.Error(suite.T(), err, "Expected an error when env name is missing")
}

func (suite *DeploymentWorkerTestSuite) TestUpdateInProgressDeploymentGetInProgressDeploymentFails() {
	suite.deployment.EXPECT().GetInProgressDeployment(suite.ctx, environmentName).
		Return(nil, errors.New("Get in progress deployment failed"))

	_, err := suite.deploymentWorker.UpdateInProgressDeployment(suite.ctx, environmentName)
	assert.Error(suite.T(), err, "Expected an error when get in progress deployment fails")
}

func (suite *DeploymentWorkerTestSuite) TestUpdateInProgressDeploymentNoInProgressDeployment() {
	suite.deployment.EXPECT().GetInProgressDeployment(suite.ctx, environmentName).Return(nil, nil)

	d, err := suite.deploymentWorker.UpdateInProgressDeployment(suite.ctx, environmentName)
	assert.Nil(suite.T(), err, "Unexpected error when get in progress deployment returns empty")
	assert.Nil(suite.T(), d, "Deployment should be nil when get in progress Deployment returns empty")
}

func (suite *DeploymentWorkerTestSuite) TestUpdateInProgressDeploymentGetEnvironmentFails() {
	suite.deployment.EXPECT().GetInProgressDeployment(suite.ctx, environmentName).Return(suite.deploymentObject, nil)
	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).Return(nil, errors.New("Get environment failed"))

	_, err := suite.deploymentWorker.UpdateInProgressDeployment(suite.ctx, environmentName)
	assert.Error(suite.T(), err, "Expected an error when get environment fails")
}

func (suite *DeploymentWorkerTestSuite) TestUpdateInProgressDeploymentGetEnvironmentIsNil() {
	suite.deployment.EXPECT().GetInProgressDeployment(suite.ctx, environmentName).Return(suite.deploymentObject, nil)
	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).Return(nil, nil)

	d, err := suite.deploymentWorker.UpdateInProgressDeployment(suite.ctx, environmentName)
	assert.Nil(suite.T(), err, "Unxpected error when get environment is empty")
	assert.Nil(suite.T(), d, "Deployment should be nil when get environment returns empty")
}

func (suite *DeploymentWorkerTestSuite) TestUpdateInProgressDeploymentListTasksFails() {
	suite.deployment.EXPECT().GetInProgressDeployment(suite.ctx, environmentName).Return(suite.deploymentObject, nil)
	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).Return(suite.environmentObject, nil)
	suite.ecs.EXPECT().ListTasks(suite.environmentObject.Cluster, suite.deploymentObject.ID).
		Return(nil, errors.New("ListTasks failed"))

	_, err := suite.deploymentWorker.UpdateInProgressDeployment(suite.ctx, environmentName)
	assert.Error(suite.T(), err, "Expected an error when list tasks fails")
}

func (suite *DeploymentWorkerTestSuite) TestUpdateInProgressDeploymentDescribeTasksFails() {
	suite.deployment.EXPECT().GetInProgressDeployment(suite.ctx, environmentName).Return(suite.deploymentObject, nil)
	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).Return(suite.environmentObject, nil)
	suite.ecs.EXPECT().ListTasks(suite.environmentObject.Cluster, suite.deploymentObject.ID).
		Return(suite.clusterTaskARNs, nil)
	suite.ecs.EXPECT().DescribeTasks(suite.environmentObject.Cluster, suite.clusterTaskARNs).
		Return(nil, errors.New("DescribeTasks failed"))

	_, err := suite.deploymentWorker.UpdateInProgressDeployment(suite.ctx, environmentName)
	assert.Error(suite.T(), err, "Expected an error when describe tasks fails")
}
