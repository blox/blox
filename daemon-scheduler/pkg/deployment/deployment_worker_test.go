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
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/blox/blox/daemon-scheduler/pkg/facade"
	"github.com/blox/blox/daemon-scheduler/pkg/mocks"
	"github.com/blox/blox/daemon-scheduler/pkg/types"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	TaskRunning = "RUNNING"
)

type DeploymentWorkerTestSuite struct {
	suite.Suite
	environment                *mocks.MockEnvironment
	environmentFacade          *types.MockEnvironmentFacade
	deployment                 *mocks.MockDeployment
	clusterState               *facade.MockClusterState
	ecs                        *mocks.MockECS
	environmentObject          *types.Environment
	pendingDeploymentObject    *types.Deployment
	inProgressDeploymentObject *types.Deployment
	clusterTaskARNs            []*string
	emptyDescribeTasksOutput   *ecs.DescribeTasksOutput
	deploymentWorker           DeploymentWorker
	ctx                        context.Context
}

func (suite *DeploymentWorkerTestSuite) SetupTest() {
	mockCtrl := gomock.NewController(suite.T())

	suite.environment = mocks.NewMockEnvironment(mockCtrl)
	suite.environmentFacade = types.NewMockEnvironmentFacade(mockCtrl)
	suite.deployment = mocks.NewMockDeployment(mockCtrl)
	suite.clusterState = facade.NewMockClusterState(mockCtrl)
	suite.ecs = mocks.NewMockECS(mockCtrl)
	suite.deploymentWorker = NewDeploymentWorker(suite.environment, suite.environmentFacade,
		suite.deployment, suite.ecs, suite.clusterState)

	var err error
	suite.environmentObject, err = types.NewEnvironment(environmentName, taskDefinition, cluster1)
	assert.Nil(suite.T(), err, "Cannot initialize DeploymentWorkerTestSuite")

	suite.pendingDeploymentObject, err = types.NewDeployment(taskDefinition, suite.environmentObject.Token)
	assert.Nil(suite.T(), err, "Cannot initialize DeploymentWorkerTestSuite")

	inProgressDeployment := *suite.pendingDeploymentObject
	suite.inProgressDeploymentObject = &inProgressDeployment

	err = suite.inProgressDeploymentObject.UpdateDeploymentToInProgress(0, nil)
	assert.Nil(suite.T(), err, "Cannot initialize DeploymentWorkerTestSuite")

	suite.clusterTaskARNs = []*string{aws.String(taskARN1), aws.String(taskARN2)}
	suite.emptyDescribeTasksOutput = &ecs.DescribeTasksOutput{
		Tasks: []*ecs.Task{},
	}

	suite.ctx = context.TODO()
}

func TestDeploymentWorkerTestSuite(t *testing.T) {
	suite.Run(t, new(DeploymentWorkerTestSuite))
}

func (suite *DeploymentWorkerTestSuite) TestNewDeploymentWorker() {
	w := NewDeploymentWorker(suite.environment, suite.environmentFacade, suite.deployment, suite.ecs, suite.clusterState)
	assert.NotNil(suite.T(), w, "Worker should not be nil")
}

func (suite *DeploymentWorkerTestSuite) TestStartPendingDeploymentEmptyEnvironmentName() {
	_, err := suite.deploymentWorker.StartPendingDeployment(suite.ctx, "")
	assert.Error(suite.T(), err, "Expected an error when env name is missing")
}

func (suite *DeploymentWorkerTestSuite) TestStartPendingDeploymentGetInProgressDeploymentFails() {
	suite.deployment.EXPECT().GetInProgressDeployment(suite.ctx, environmentName).
		Return(nil, errors.New("Get in progress deployment failed"))

	_, err := suite.deploymentWorker.StartPendingDeployment(suite.ctx, environmentName)
	assert.Error(suite.T(), err, "Expected an error when get in progress deployment fails")
}

func (suite *DeploymentWorkerTestSuite) TestStartPendingDeploymentInProgressDeploymentExists() {
	suite.deployment.EXPECT().GetInProgressDeployment(suite.ctx, environmentName).
		Return(suite.inProgressDeploymentObject, nil)

	d, err := suite.deploymentWorker.StartPendingDeployment(suite.ctx, environmentName)
	assert.Nil(suite.T(), err, "Unexpected error when get in progress deployment exists")
	assert.Nil(suite.T(), d, "Deployment should be nil when in progress deployment exists")
}

func (suite *DeploymentWorkerTestSuite) TestStartPendingDeploymentGetPendingDeploymentFails() {
	suite.deployment.EXPECT().GetInProgressDeployment(suite.ctx, environmentName).Return(nil, nil)

	suite.deployment.EXPECT().GetPendingDeployment(suite.ctx, environmentName).
		Return(nil, errors.New("Get pending deployment failed"))

	_, err := suite.deploymentWorker.StartPendingDeployment(suite.ctx, environmentName)
	assert.Error(suite.T(), err, "Expected an error when get pending deployment fails")
}

func (suite *DeploymentWorkerTestSuite) TestStartPendingDeploymentNoPendingDeploymentExists() {
	suite.deployment.EXPECT().GetInProgressDeployment(suite.ctx, environmentName).Return(nil, nil)
	suite.deployment.EXPECT().GetPendingDeployment(suite.ctx, environmentName).Return(nil, nil)

	d, err := suite.deploymentWorker.StartPendingDeployment(suite.ctx, environmentName)
	assert.Nil(suite.T(), err, "Unexpected error when no pending deployment exists")
	assert.Nil(suite.T(), d, "Expected no deployments when no pending deployment exists")
}

func (suite *DeploymentWorkerTestSuite) TestStartPendingDeploymentGetEnvironmentFails() {
	suite.deployment.EXPECT().GetInProgressDeployment(suite.ctx, environmentName).Return(nil, nil)
	suite.deployment.EXPECT().GetPendingDeployment(suite.ctx, environmentName).Return(suite.pendingDeploymentObject, nil)
	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).Return(nil, errors.New("Get environment failed"))

	_, err := suite.deploymentWorker.StartPendingDeployment(suite.ctx, environmentName)
	assert.Error(suite.T(), err, "Expected an error when get environment fails")
}

func (suite *DeploymentWorkerTestSuite) TestStartPendingDeploymentGetEnvironmentIsEmpty() {
	suite.deployment.EXPECT().GetInProgressDeployment(suite.ctx, environmentName).Return(nil, nil)
	suite.deployment.EXPECT().GetPendingDeployment(suite.ctx, environmentName).Return(suite.pendingDeploymentObject, nil)
	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).Return(nil, nil)

	d, err := suite.deploymentWorker.StartPendingDeployment(suite.ctx, environmentName)
	assert.Nil(suite.T(), err, "Unexpected error when get environment is empty")
	assert.Nil(suite.T(), d)
}

func (suite *DeploymentWorkerTestSuite) TestStartPendingDeploymentInstanceARNsFails() {
	suite.deployment.EXPECT().GetInProgressDeployment(suite.ctx, environmentName).Return(nil, nil)
	suite.deployment.EXPECT().GetPendingDeployment(suite.ctx, environmentName).Return(suite.pendingDeploymentObject, nil)
	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).Return(suite.environmentObject, nil)
	suite.environmentFacade.EXPECT().InstanceARNs(suite.environmentObject).Return(nil, errors.New("Instance ARNs fails"))

	_, err := suite.deploymentWorker.StartPendingDeployment(suite.ctx, environmentName)
	assert.Error(suite.T(), err, "Expected an error when get instance arns fails")
}

func (suite *DeploymentWorkerTestSuite) TestStartPendingDeploymentStartDeploymentFails() {
	suite.deployment.EXPECT().GetInProgressDeployment(suite.ctx, environmentName).Return(nil, nil)
	suite.deployment.EXPECT().GetPendingDeployment(suite.ctx, environmentName).Return(suite.pendingDeploymentObject, nil)
	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).Return(suite.environmentObject, nil)
	suite.environmentFacade.EXPECT().InstanceARNs(suite.environmentObject).Return(suite.clusterTaskARNs, nil)
	suite.deployment.EXPECT().StartDeployment(suite.ctx, suite.environmentObject, suite.pendingDeploymentObject, suite.clusterTaskARNs).
		Return(nil, errors.New("Start deployment fails"))

	_, err := suite.deploymentWorker.StartPendingDeployment(suite.ctx, environmentName)
	assert.Error(suite.T(), err, "Expected an error when start deployment fails")
}

func (suite *DeploymentWorkerTestSuite) TestStartPendingDeployment() {
	suite.deployment.EXPECT().GetInProgressDeployment(suite.ctx, environmentName).Return(nil, nil)
	suite.deployment.EXPECT().GetPendingDeployment(suite.ctx, environmentName).Return(suite.pendingDeploymentObject, nil)
	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).Return(suite.environmentObject, nil)
	suite.environmentFacade.EXPECT().InstanceARNs(suite.environmentObject).Return(suite.clusterTaskARNs, nil)
	suite.deployment.EXPECT().StartDeployment(suite.ctx, suite.environmentObject, suite.pendingDeploymentObject, suite.clusterTaskARNs).
		Return(suite.inProgressDeploymentObject, nil)

	d, err := suite.deploymentWorker.StartPendingDeployment(suite.ctx, environmentName)
	assert.Nil(suite.T(), err)
	verifyDeployment(suite.T(), suite.inProgressDeploymentObject, d)
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
	suite.deployment.EXPECT().GetInProgressDeployment(suite.ctx, environmentName).Return(suite.inProgressDeploymentObject, nil)
	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).Return(nil, errors.New("Get environment failed"))

	_, err := suite.deploymentWorker.UpdateInProgressDeployment(suite.ctx, environmentName)
	assert.Error(suite.T(), err, "Expected an error when get environment fails")
}

func (suite *DeploymentWorkerTestSuite) TestUpdateInProgressDeploymentGetEnvironmentIsNil() {
	suite.deployment.EXPECT().GetInProgressDeployment(suite.ctx, environmentName).Return(suite.inProgressDeploymentObject, nil)
	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).Return(nil, nil)

	d, err := suite.deploymentWorker.UpdateInProgressDeployment(suite.ctx, environmentName)
	assert.Nil(suite.T(), err, "Unxpected error when get environment is empty")
	assert.Nil(suite.T(), d, "Deployment should be nil when get environment returns empty")
}

func (suite *DeploymentWorkerTestSuite) TestUpdateInProgressDeploymentListTasksFails() {
	suite.deployment.EXPECT().GetInProgressDeployment(suite.ctx, environmentName).Return(suite.inProgressDeploymentObject, nil)
	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).Return(suite.environmentObject, nil)
	suite.ecs.EXPECT().ListTasks(suite.environmentObject.Cluster, suite.inProgressDeploymentObject.ID).
		Return(nil, errors.New("ListTasks failed"))

	_, err := suite.deploymentWorker.UpdateInProgressDeployment(suite.ctx, environmentName)
	assert.Error(suite.T(), err, "Expected an error when list tasks fails")
}

func (suite *DeploymentWorkerTestSuite) TestUpdateInProgressDeploymentDescribeTasksFails() {
	suite.deployment.EXPECT().GetInProgressDeployment(suite.ctx, environmentName).Return(suite.inProgressDeploymentObject, nil)
	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).Return(suite.environmentObject, nil)
	suite.ecs.EXPECT().ListTasks(suite.environmentObject.Cluster, suite.inProgressDeploymentObject.ID).
		Return(suite.clusterTaskARNs, nil)
	suite.ecs.EXPECT().DescribeTasks(suite.environmentObject.Cluster, suite.clusterTaskARNs).
		Return(nil, errors.New("DescribeTasks failed"))

	_, err := suite.deploymentWorker.UpdateInProgressDeployment(suite.ctx, environmentName)
	assert.Error(suite.T(), err, "Expected an error when describe tasks fails")
}

func (suite *DeploymentWorkerTestSuite) TestUpdateInProgressDeploymentNoTasksStartedByTheDeployment() {
	suite.deployment.EXPECT().GetInProgressDeployment(suite.ctx, environmentName).Return(suite.inProgressDeploymentObject, nil).Times(2)
	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).Return(suite.environmentObject, nil)
	suite.ecs.EXPECT().ListTasks(suite.environmentObject.Cluster, suite.inProgressDeploymentObject.ID).
		Return(suite.clusterTaskARNs, nil)

	noTasks := &ecs.DescribeTasksOutput{
		Tasks: []*ecs.Task{},
	}

	suite.ecs.EXPECT().DescribeTasks(suite.environmentObject.Cluster, suite.clusterTaskARNs).Return(noTasks, nil)
	suite.environment.EXPECT().UpdateDeployment(suite.ctx, *suite.environmentObject, *suite.inProgressDeploymentObject).
		Return(suite.environmentObject, nil)

	d, err := suite.deploymentWorker.UpdateInProgressDeployment(suite.ctx, environmentName)
	assert.Nil(suite.T(), err, "Unexpected error when there are no tasks started by the deployment")
	assert.Equal(suite.T(), suite.inProgressDeploymentObject, d, "Expected deployments to match")
}

func (suite *DeploymentWorkerTestSuite) TestUpdateInProgressDeploymentTasksArePending() {
	suite.deployment.EXPECT().GetInProgressDeployment(suite.ctx, environmentName).Return(suite.inProgressDeploymentObject, nil).Times(2)
	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).Return(suite.environmentObject, nil)
	suite.ecs.EXPECT().ListTasks(suite.environmentObject.Cluster, suite.inProgressDeploymentObject.ID).
		Return(suite.clusterTaskARNs, nil)

	pendingTask := &ecs.Task{
		TaskArn:    aws.String(taskARN1),
		LastStatus: aws.String(TaskPending),
	}

	runningTask := &ecs.Task{
		TaskArn:    aws.String(taskARN2),
		LastStatus: aws.String(TaskRunning),
	}

	tasks := &ecs.DescribeTasksOutput{
		Tasks: []*ecs.Task{runningTask, pendingTask},
	}

	suite.ecs.EXPECT().DescribeTasks(suite.environmentObject.Cluster, suite.clusterTaskARNs).Return(tasks, nil)
	suite.environment.EXPECT().UpdateDeployment(suite.ctx, *suite.environmentObject, *suite.inProgressDeploymentObject).
		Return(suite.environmentObject, nil)

	d, err := suite.deploymentWorker.UpdateInProgressDeployment(suite.ctx, environmentName)
	assert.Nil(suite.T(), err, "Unexpected error when there is a pending task started by the deployment")
	assert.Equal(suite.T(), suite.inProgressDeploymentObject, d, "Expected deployments to match")
}

func (suite *DeploymentWorkerTestSuite) TestUpdateInProgressDeploymentDeploymentCompleted() {
	suite.deployment.EXPECT().GetInProgressDeployment(suite.ctx, environmentName).Return(suite.inProgressDeploymentObject, nil).Times(2)
	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).Return(suite.environmentObject, nil)
	suite.ecs.EXPECT().ListTasks(suite.environmentObject.Cluster, suite.inProgressDeploymentObject.ID).
		Return(suite.clusterTaskARNs, nil)

	runningTask1 := &ecs.Task{
		TaskArn:    aws.String(taskARN1),
		LastStatus: aws.String(TaskRunning),
	}

	runningTask2 := &ecs.Task{
		TaskArn:    aws.String(taskARN2),
		LastStatus: aws.String(TaskRunning),
	}

	tasks := &ecs.DescribeTasksOutput{
		Tasks: []*ecs.Task{runningTask1, runningTask2},
	}

	suite.ecs.EXPECT().DescribeTasks(suite.environmentObject.Cluster, suite.clusterTaskARNs).Return(tasks, nil)

	err := suite.inProgressDeploymentObject.UpdateDeploymentToCompleted(nil)
	completedDeployment := suite.inProgressDeploymentObject
	assert.Nil(suite.T(), err, "Unexpected error when moving deployment to completed")

	suite.environment.EXPECT().UpdateDeployment(suite.ctx, *suite.environmentObject, gomock.Any()).Do(
		func(_ interface{}, _ interface{}, d types.Deployment) {
			verifyDeploymentCompleted(suite.T(), completedDeployment, &d)
		}).Return(suite.environmentObject, nil)

	d, err := suite.deploymentWorker.UpdateInProgressDeployment(suite.ctx, environmentName)
	assert.Nil(suite.T(), err, "Unexpected error when the deployment is completed")
	verifyDeploymentCompleted(suite.T(), completedDeployment, d)
}

func (suite *DeploymentWorkerTestSuite) TestUpdateInProgressDeploymentInProgressDeploymentCheckFailsAfterUpdatingDeploymentObject() {
	gomock.InOrder(
		suite.deployment.EXPECT().GetInProgressDeployment(suite.ctx, environmentName).Return(suite.inProgressDeploymentObject, nil),
		suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).Return(suite.environmentObject, nil),
		suite.ecs.EXPECT().ListTasks(suite.environmentObject.Cluster, suite.inProgressDeploymentObject.ID).
			Return(suite.clusterTaskARNs, nil),
		suite.ecs.EXPECT().DescribeTasks(suite.environmentObject.Cluster, suite.clusterTaskARNs).
			Return(suite.emptyDescribeTasksOutput, nil),
		suite.deployment.EXPECT().GetInProgressDeployment(suite.ctx, environmentName).Return(nil, errors.New("Second in-progress deployment check fails")),
	)
	_, err := suite.deploymentWorker.UpdateInProgressDeployment(suite.ctx, environmentName)
	assert.Error(suite.T(), err, "Expected an error when in progress deployment check fails")
}

func (suite *DeploymentWorkerTestSuite) TestUpdateInProgressDeploymentNoInProgressDeploymentAfterUpdatingDeploymentObject() {
	gomock.InOrder(
		suite.deployment.EXPECT().GetInProgressDeployment(suite.ctx, environmentName).Return(suite.inProgressDeploymentObject, nil),
		suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).Return(suite.environmentObject, nil),
		suite.ecs.EXPECT().ListTasks(suite.environmentObject.Cluster, suite.inProgressDeploymentObject.ID).
			Return(suite.clusterTaskARNs, nil),
		suite.ecs.EXPECT().DescribeTasks(suite.environmentObject.Cluster, suite.clusterTaskARNs).
			Return(suite.emptyDescribeTasksOutput, nil),
		suite.deployment.EXPECT().GetInProgressDeployment(suite.ctx, environmentName).Return(nil, nil),
	)
	d, err := suite.deploymentWorker.UpdateInProgressDeployment(suite.ctx, environmentName)
	assert.Nil(suite.T(), err, "Unexpected error when there is no in-progress deployment")
	assert.Nil(suite.T(), d, "No deployment should be updated if there is no in-progress deployment")
}

func (suite *DeploymentWorkerTestSuite) TestUpdateInProgressDeploymentInProgressDeploymentIsDifferentAfterUpdatingDeploymentObject() {
	newInProgressDeployment, err := types.NewDeployment(taskDefinition, suite.environmentObject.Token)
	assert.Nil(suite.T(), err, "Could not create a new deployment")

	gomock.InOrder(
		suite.deployment.EXPECT().GetInProgressDeployment(suite.ctx, environmentName).Return(suite.inProgressDeploymentObject, nil),
		suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).Return(suite.environmentObject, nil),
		suite.ecs.EXPECT().ListTasks(suite.environmentObject.Cluster, suite.inProgressDeploymentObject.ID).
			Return(suite.clusterTaskARNs, nil),
		suite.ecs.EXPECT().DescribeTasks(suite.environmentObject.Cluster, suite.clusterTaskARNs).
			Return(suite.emptyDescribeTasksOutput, nil),
		suite.deployment.EXPECT().GetInProgressDeployment(suite.ctx, environmentName).Return(newInProgressDeployment, nil),
	)

	d, err := suite.deploymentWorker.UpdateInProgressDeployment(suite.ctx, environmentName)
	assert.Nil(suite.T(), err, "Unexpected error when the in-progress deployment has changed")
	assert.Nil(suite.T(), d, "No deployment should be updated if the in-progress deployment has changed")
}

func (suite *DeploymentWorkerTestSuite) TestUpdateInProgressDeploymentUpdateDeploymentFails() {
	gomock.InOrder(
		suite.deployment.EXPECT().GetInProgressDeployment(suite.ctx, environmentName).Return(suite.inProgressDeploymentObject, nil),
		suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).Return(suite.environmentObject, nil),
		suite.ecs.EXPECT().ListTasks(suite.environmentObject.Cluster, suite.inProgressDeploymentObject.ID).
			Return(suite.clusterTaskARNs, nil),
		suite.ecs.EXPECT().DescribeTasks(suite.environmentObject.Cluster, suite.clusterTaskARNs).
			Return(suite.emptyDescribeTasksOutput, nil),
		suite.deployment.EXPECT().GetInProgressDeployment(suite.ctx, environmentName).Return(suite.inProgressDeploymentObject, nil),
		suite.environment.EXPECT().UpdateDeployment(suite.ctx, *suite.environmentObject, *suite.inProgressDeploymentObject).
			Return(nil, errors.New("Update deployment failed")),
	)

	_, err := suite.deploymentWorker.UpdateInProgressDeployment(suite.ctx, environmentName)
	assert.Error(suite.T(), err, "Expected an error when update deployment fails")
}

func verifyDeploymentCompleted(t *testing.T, expected *types.Deployment, actual *types.Deployment) {
	assert.Exactly(t, expected.ID, actual.ID, "Deployment ids should match")
	assert.Exactly(t, types.DeploymentCompleted, actual.Status, "Deployment status should be completed")
	assert.Exactly(t, expected.Health, actual.Health, "Deployment health should match")
	assert.Exactly(t, expected.TaskDefinition, actual.TaskDefinition, "Deployment task definition should match")
	assert.Exactly(t, expected.DesiredTaskCount, actual.DesiredTaskCount, "Deployment desired task count should match")
	assert.Exactly(t, expected.FailedInstances, actual.FailedInstances, "Deployment failed instances should match")
}
