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
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/blox/blox/cluster-state-service/swagger/v1/generated/models"
	"github.com/blox/blox/daemon-scheduler/pkg/facade"
	"github.com/blox/blox/daemon-scheduler/pkg/mocks"
	"github.com/blox/blox/daemon-scheduler/pkg/types"
	"github.com/golang/mock/gomock"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	deploymentID = "deploymentID"
	instanceARN1 = "arn:aws:us-east-1:123456789123:container-instance/4b6d45ea-a4b4-4269-9d04-3af6ddfdc597"
	instanceARN2 = "arn:aws:us-east-1:123456789123:container-instance/15a1c3d8-e449-4377-9aed-affafb3da5eb"
)

type DeploymentTestSuite struct {
	suite.Suite
	environment       *mocks.MockEnvironment
	clusterState      *facade.MockClusterState
	ecs               *mocks.MockECS
	deployment        Deployment
	ctx               context.Context
	environmentObject *types.Environment
	deploymentObject  *types.Deployment
	token             string
	instanceARNs      []*string
	startTaskOutput   *ecs.StartTaskOutput
}

func (suite *DeploymentTestSuite) SetupTest() {
	mockCtrl := gomock.NewController(suite.T())

	suite.environment = mocks.NewMockEnvironment(mockCtrl)
	suite.clusterState = facade.NewMockClusterState(mockCtrl)
	suite.ecs = mocks.NewMockECS(mockCtrl)
	suite.deployment = NewDeployment(suite.environment, suite.clusterState, suite.ecs)

	var err error
	suite.environmentObject, err = types.NewEnvironment(environmentName, taskDefinition, cluster1)
	assert.Nil(suite.T(), err, "Cannot initialize DeploymentTestSuite")

	suite.token = suite.environmentObject.Token

	suite.deploymentObject, err = types.NewDeployment(taskDefinition, suite.token)
	assert.Nil(suite.T(), err, "Cannot initialize DeploymentTestSuite")

	suite.instanceARNs = []*string{aws.String(instanceARN1), aws.String(instanceARN2)}

	task := ecs.Task{
		TaskArn:           aws.String(taskARN1),
		TaskDefinitionArn: aws.String(taskDefinition),
	}

	failure := ecs.Failure{
		Arn: aws.String(instanceARN1),
	}

	suite.startTaskOutput = &ecs.StartTaskOutput{
		Tasks:    []*ecs.Task{&task},
		Failures: []*ecs.Failure{&failure},
	}

	suite.ctx = context.TODO()
}

func TestDeploymentTestSuite(t *testing.T) {
	suite.Run(t, new(DeploymentTestSuite))
}

func (suite *DeploymentTestSuite) TestNewDeployment() {
	d := NewDeployment(suite.environment, suite.clusterState, suite.ecs)
	assert.NotNil(suite.T(), d, "Expected an error when store is nil")
}

func (suite *DeploymentTestSuite) TestCreateDeploymentEmptyEnvironmentName() {
	_, err := suite.deployment.CreateDeployment(suite.ctx, "", suite.token)
	assert.Error(suite.T(), err, "Expected an error when environment name is empty")
}

func (suite *DeploymentTestSuite) TestCreateDeploymentEmptyToken() {
	_, err := suite.deployment.CreateDeployment(suite.ctx, environmentName, "")
	assert.Error(suite.T(), err, "Expected an error when token is empty")
}

func (suite *DeploymentTestSuite) TestCreateDeploymentGetEnvironmentFails() {
	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).
		Return(nil, errors.New("Get environment failed"))

	_, err := suite.deployment.CreateDeployment(suite.ctx, environmentName, suite.token)
	assert.Error(suite.T(), err, "Expected an error when get environment fails")
}

func (suite *DeploymentTestSuite) TestCreateDeploymentGetEnvironmentIsNil() {
	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).
		Return(nil, nil)

	_, err := suite.deployment.CreateDeployment(suite.ctx, environmentName, suite.token)
	assert.Error(suite.T(), err, "Expected an error when get environment is nil")
}

func (suite *DeploymentTestSuite) TestCreateDeploymentOutdatedToken() {
	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).
		Return(suite.environmentObject, nil)

	_, err := suite.deployment.CreateDeployment(suite.ctx, environmentName, "invalid")
	assert.Error(suite.T(), err, "Expected an error when token is outdated")
}

func (suite *DeploymentTestSuite) TestCreateDeploymentExistingToken() {
	suite.deploymentObject.Token = suite.token

	environment := *suite.environmentObject
	environment.Deployments[suite.deploymentObject.ID] = *suite.deploymentObject

	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).
		Return(suite.environmentObject, nil)

	_, err := suite.deployment.CreateDeployment(suite.ctx, environmentName, suite.token)
	assert.Error(suite.T(), err, "Expected an error when a deployment with the given token exists")
}

func (suite *DeploymentTestSuite) TestCreateDeploymentGetInProgressDeploymentFails() {
	// deployment does not exist in the deployments map -> GetCurrentDeployment fails
	suite.environmentObject.InProgressDeploymentID = uuid.NewUUID().String()

	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).Return(suite.environmentObject, nil).Times(2)

	_, err := suite.deployment.CreateDeployment(suite.ctx, environmentName, suite.token)
	assert.Error(suite.T(), err, "Expected an error when getting in progress deployment fails")
}

func (suite *DeploymentTestSuite) TestCreateDeploymentThereIsAnInProgressDeployment() {
	suite.environmentObject.InProgressDeploymentID = suite.deploymentObject.ID
	suite.environmentObject.Deployments[suite.deploymentObject.ID] = *suite.deploymentObject

	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).Return(suite.environmentObject, nil)
	suite.environment.EXPECT().AddPendingDeployment(suite.ctx, *suite.environmentObject, gomock.Any()).Times(0)

	_, err := suite.deployment.CreateDeployment(suite.ctx, environmentName, suite.environmentObject.Token)
	assert.Error(suite.T(), err, "Expected an error when there is an in-progress deployment")
}

func (suite *DeploymentTestSuite) TestCreateDeploymentAddDeploymentFails() {
	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).Return(suite.environmentObject, nil).Times(2)

	suite.environment.EXPECT().AddPendingDeployment(suite.ctx, *suite.environmentObject, gomock.Any()).Do(
		func(_ interface{}, _ interface{}, d types.Deployment) {
			verifyDeployment(suite.T(), suite.deploymentObject, &d)
		}).Return(nil, errors.New("Add deployment failed"))

	_, err := suite.deployment.CreateDeployment(suite.ctx, environmentName, suite.environmentObject.Token)
	assert.Error(suite.T(), err, "Expected an error when add deployment fails")
}

func (suite *DeploymentTestSuite) TestCreateDeployment() {
	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).Return(suite.environmentObject, nil).Times(2)

	suite.environment.EXPECT().AddPendingDeployment(suite.ctx, *suite.environmentObject, gomock.Any()).Do(
		func(_ interface{}, _ interface{}, d types.Deployment) {
			verifyDeployment(suite.T(), suite.deploymentObject, &d)
		}).Return(suite.environmentObject, nil)

	d, err := suite.deployment.CreateDeployment(suite.ctx, environmentName, suite.environmentObject.Token)
	assert.Nil(suite.T(), err, "Unexpected error when creating a deployment")
	verifyDeployment(suite.T(), suite.deploymentObject, d)
}

func verifyDeployment(t *testing.T, expected *types.Deployment, actual *types.Deployment) {
	assert.NotEmpty(t, actual.ID, "Deployment ID should not be empty")
	assert.Exactly(t, expected.Status, actual.Status, "Deployment status should match")
	assert.Exactly(t, expected.Health, actual.Health, "Deployment health should match")
	assert.Exactly(t, expected.TaskDefinition, actual.TaskDefinition, "Deployment task definition should match")
	assert.Exactly(t, expected.DesiredTaskCount, actual.DesiredTaskCount, "Deployment desired task count should match")
	assert.NotEmpty(t, actual.StartTime, "Deployment start time should not be empty")
	assert.Exactly(t, expected.EndTime, actual.EndTime, "Deployment end time should match")
}

func (suite *DeploymentTestSuite) TestGetDeploymentEmptyEnvironmentName() {
	_, err := suite.deployment.GetDeployment(suite.ctx, "", deploymentID)
	assert.Error(suite.T(), err, "Expected an error when environment name is empty")
}

func (suite *DeploymentTestSuite) TestGetDeploymentEmptyDeploymentID() {
	_, err := suite.deployment.GetDeployment(suite.ctx, environmentName, "")
	assert.Error(suite.T(), err, "Expected an error when deployment ID is empty")
}

func (suite *DeploymentTestSuite) TestGetDeploymentGetEnvironmentDeploymentsFails() {
	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).
		Return(nil, errors.New("Get environment failed"))

	_, err := suite.deployment.GetDeployment(suite.ctx, environmentName, deploymentID)
	assert.Error(suite.T(), err, "Expected an error when get environment fails")
}

func (suite *DeploymentTestSuite) TestGetDeploymentEnvironmentDoesNotHaveDeployments() {
	suite.environmentObject.Deployments = nil
	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).Return(suite.environmentObject, nil)

	d, err := suite.deployment.GetDeployment(suite.ctx, environmentName, deploymentID)
	assert.Nil(suite.T(), err, "Unexpected error when the environment does not have deployments")
	assert.Nil(suite.T(), d, "Expected a nil result when the environment does not have deployments")
}

func (suite *DeploymentTestSuite) TestGetDeploymentEnvironmentDoesNotHaveAMatchingDeployment() {
	suite.environmentObject.Deployments[suite.deploymentObject.ID] = *suite.deploymentObject
	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).Return(suite.environmentObject, nil)

	d, err := suite.deployment.GetDeployment(suite.ctx, environmentName, "non-existing-ID")
	assert.Nil(suite.T(), err, "Unexpected error when the environment does not have deployments")
	assert.Nil(suite.T(), d, "Expected a nil result when the environment does not have a matching deployment")
}

func (suite *DeploymentTestSuite) TestGetDeployment() {
	deployment1, err := types.NewDeployment(suite.environmentObject.DesiredTaskDefinition, suite.environmentObject.Token)
	assert.Nil(suite.T(), err, "Could not create a new deployment")

	deployment2, err := types.NewDeployment(suite.environmentObject.DesiredTaskDefinition, suite.environmentObject.Token)
	assert.Nil(suite.T(), err, "Could not create a new deployment")

	suite.environmentObject.Deployments[deployment1.ID] = *deployment1
	suite.environmentObject.Deployments[deployment2.ID] = *deployment2

	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).Return(suite.environmentObject, nil)

	d, err := suite.deployment.GetDeployment(suite.ctx, environmentName, deployment2.ID)
	assert.Nil(suite.T(), err, "Unexpected error when the environment has multiple deployments")
	assert.Exactly(suite.T(), deployment2, d, "Deployment does not match the one in the environment")
}

func (suite *DeploymentTestSuite) TestCreateSubDeploymentEmptyEnvironmentName() {
	_, err := suite.deployment.CreateSubDeployment(suite.ctx, "", suite.instanceARNs)
	assert.NotNil(suite.T(), err, "Expected an error creating a sub-deployment without an environment name")
}

func (suite *DeploymentTestSuite) TestCreateSubDeploymentGetEnvironmentFails() {
	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).Return(nil, errors.New("Get environment failed"))

	_, err := suite.deployment.CreateSubDeployment(suite.ctx, environmentName, suite.instanceARNs)
	assert.NotNil(suite.T(), err, "Expected an error creating a sub-deployment when get environment fails")
}

func (suite *DeploymentTestSuite) TestCreateSubDeploymentGetEnvironmentReturnsNil() {
	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).Return(nil, nil)

	_, err := suite.deployment.CreateSubDeployment(suite.ctx, environmentName, suite.instanceARNs)
	assert.NotNil(suite.T(), err, "Expected an error creating a sub-deployment when get environment returns nil")
}

func (suite *DeploymentTestSuite) TestCreateSubDeploymentGetCurrentDeploymentReturnsError() {
	// deployment does not exist in the deployments map -> GetCurrentDeployment fails
	suite.environmentObject.InProgressDeploymentID = uuid.NewUUID().String()

	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).Return(suite.environmentObject, nil).Times(2)

	_, err := suite.deployment.CreateSubDeployment(suite.ctx, environmentName, suite.instanceARNs)
	assert.NotNil(suite.T(), err, "Expected an error creating a sub-deployment when get in progress deployment returns an error")
}

func (suite *DeploymentTestSuite) TestCreateSubDeploymentGetCurrentDeploymentReturnsNoDeployment() {
	suite.environmentObject.InProgressDeploymentID = ""
	suite.environmentObject.PendingDeploymentID = ""
	suite.environmentObject.Deployments = nil

	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).Return(suite.environmentObject, nil).Times(3)

	_, err := suite.deployment.CreateSubDeployment(suite.ctx, environmentName, suite.instanceARNs)
	assert.NotNil(suite.T(), err, "Expected an error creating a sub-deployment when get in progress deployment returns an error")
}

func (suite *DeploymentTestSuite) TestCreateSubDeploymentStartTasksFails() {
	err := suite.deploymentObject.UpdateDeploymentToInProgress(0, nil)
	assert.Nil(suite.T(), err, "Unexpected error when moving deployment to in-progress")

	env := suite.environmentObject
	env.InProgressDeploymentID = suite.deploymentObject.ID
	env.Deployments[suite.deploymentObject.ID] = *suite.deploymentObject

	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).Return(env, nil).Times(2)
	suite.ecs.EXPECT().StartTask(env.Cluster, suite.instanceARNs, suite.deploymentObject.ID, suite.deploymentObject.TaskDefinition).
		Return(nil, errors.New("Error starting tasks"))

	_, err = suite.deployment.CreateSubDeployment(suite.ctx, environmentName, suite.instanceARNs)
	assert.NotNil(suite.T(), err, "Expected an error creating a sub-deployment when start tasks fails")
}

func (suite *DeploymentTestSuite) TestCreateSubDeploymentUpdateDeploymentFails() {
	inprogressDeployment, err := types.NewDeployment(suite.environmentObject.DesiredTaskDefinition, suite.environmentObject.Token)
	assert.Nil(suite.T(), err, "Deployment creation failed")
	inprogressDeployment.Status = types.DeploymentInProgress

	env := suite.environmentObject
	env.Name = environmentName
	env.InProgressDeploymentID = inprogressDeployment.ID
	env.Deployments[inprogressDeployment.ID] = *inprogressDeployment

	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).Return(env, nil).Times(2)
	suite.ecs.EXPECT().StartTask(env.Cluster, suite.instanceARNs, inprogressDeployment.ID, inprogressDeployment.TaskDefinition).Return(suite.startTaskOutput, nil)

	updatedDeployment := *inprogressDeployment
	updatedDeployment.DesiredTaskCount = len(suite.instanceARNs)
	updatedDeployment.Health = types.DeploymentUnhealthy
	updatedDeployment.Status = types.DeploymentInProgress
	updatedDeployment.FailedInstances = suite.startTaskOutput.Failures

	suite.environment.EXPECT().UpdateDeployment(suite.ctx, *env, gomock.Any()).Do(
		func(_ interface{}, _ interface{}, d types.Deployment) {
			verifyDeployment(suite.T(), &updatedDeployment, &d)
		}).Return(nil, errors.New("Error updating deployment"))

	_, err = suite.deployment.CreateSubDeployment(suite.ctx, environmentName, suite.instanceARNs)
	assert.NotNil(suite.T(), err, "Expected an error creating a sub-deployment when update deployment fails")
}

func (suite *DeploymentTestSuite) TestCreateSubDeployment() {
	inprogressDeployment, err := types.NewDeployment(suite.environmentObject.DesiredTaskDefinition, suite.environmentObject.Token)
	assert.Nil(suite.T(), err, "Deployment creation failed")
	inprogressDeployment.Status = types.DeploymentInProgress

	env := suite.environmentObject
	env.InProgressDeploymentID = inprogressDeployment.ID
	env.Deployments[inprogressDeployment.ID] = *inprogressDeployment

	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).Return(env, nil).Times(2)
	suite.ecs.EXPECT().StartTask(env.Cluster, suite.instanceARNs, inprogressDeployment.ID, inprogressDeployment.TaskDefinition).Return(suite.startTaskOutput, nil)

	updatedDeployment := *inprogressDeployment
	updatedDeployment.DesiredTaskCount = len(suite.instanceARNs)
	updatedDeployment.Health = types.DeploymentUnhealthy
	updatedDeployment.Status = types.DeploymentInProgress
	updatedDeployment.FailedInstances = suite.startTaskOutput.Failures

	suite.environment.EXPECT().UpdateDeployment(suite.ctx, *env, gomock.Any()).Do(
		func(_ interface{}, _ interface{}, d types.Deployment) {
			verifyDeployment(suite.T(), &updatedDeployment, &d)
		}).Return(env, nil)

	d, err := suite.deployment.CreateSubDeployment(suite.ctx, environmentName, suite.instanceARNs)
	assert.Nil(suite.T(), err, "Unexpected error creating a sub-deployment")
	verifyDeployment(suite.T(), &updatedDeployment, d)
}

func (suite *DeploymentTestSuite) TestCreateSubDeploymentWithCompletedDeployment() {
	currentDeployment, err := types.NewDeployment(suite.environmentObject.DesiredTaskDefinition, suite.environmentObject.Token)
	assert.Nil(suite.T(), err, "Deployment creation failed")
	currentDeployment.Status = types.DeploymentCompleted

	env := suite.environmentObject
	env.Deployments[currentDeployment.ID] = *currentDeployment

	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).Return(env, nil).Times(3)
	suite.ecs.EXPECT().StartTask(env.Cluster, suite.instanceARNs, currentDeployment.ID, currentDeployment.TaskDefinition).Return(suite.startTaskOutput, nil)

	d, err := suite.deployment.CreateSubDeployment(suite.ctx, environmentName, suite.instanceARNs)
	assert.Nil(suite.T(), err, "Unexpected error creating a sub-deployment")
	verifyDeployment(suite.T(), currentDeployment, d)
}

func createContainerInstances(instanceARNs []*string) []*models.ContainerInstance {
	containerInstances := make([]*models.ContainerInstance, 0, len(instanceARNs))
	for _, i := range instanceARNs {
		containerInstance := &models.ContainerInstance{
			Metadata: &models.Metadata{EntityVersion: aws.String("123")},
			Entity: &models.ContainerInstanceDetail{
				ContainerInstanceARN: i,
			},
		}
		containerInstances = append(containerInstances, containerInstance)
	}

	return containerInstances
}

func (suite *DeploymentTestSuite) TestStartDeploymentPendingDeployment() {
	suite.environmentObject.PendingDeploymentID = suite.deploymentObject.ID
	suite.environmentObject.Deployments[suite.deploymentObject.ID] = *suite.deploymentObject
	suite.ecs.EXPECT().StartTask(suite.environmentObject.Cluster, suite.instanceARNs,
		suite.deploymentObject.ID, suite.deploymentObject.TaskDefinition).Return(suite.startTaskOutput, nil)

	updatedDeployment := *suite.deploymentObject
	updatedDeployment.DesiredTaskCount = len(suite.instanceARNs)
	updatedDeployment.Status = types.DeploymentInProgress
	updatedDeployment.Health = types.DeploymentUnhealthy
	updatedDeployment.FailedInstances = suite.startTaskOutput.Failures

	suite.environment.EXPECT().UpdateDeployment(suite.ctx, gomock.Any(), gomock.Any()).Do(
		func(_ interface{}, e types.Environment, d types.Deployment) {
			assert.Empty(suite.T(), suite.environmentObject.PendingDeploymentID)
			assert.Exactly(suite.T(), suite.deploymentObject.ID, suite.environmentObject.InProgressDeploymentID)
			verifyDeployment(suite.T(), &updatedDeployment, &d)
		}).Return(nil, nil)

	d, err := suite.deployment.StartDeployment(suite.ctx, suite.environmentObject, suite.deploymentObject, suite.instanceARNs)
	assert.Nil(suite.T(), err)
	assert.Exactly(suite.T(), types.DeploymentInProgress, d.Status)
}

func (suite *DeploymentTestSuite) TestGetCurrentDeploymentOnlyInProgressExists() {
	environment, err := types.NewEnvironment("TestGetCurrent", taskDefinition, cluster1)
	assert.Nil(suite.T(), err, "Unexpected error when creating environment")
	deployment, err := types.NewDeployment(taskDefinition, uuid.NewRandom().String())
	assert.Nil(suite.T(), err, "Unexpected error when creating deployment")
	deployment.UpdateDeploymentToInProgress(0, nil)
	assert.Nil(suite.T(), err, "Unexpected error when moving deployment to in-progress")
	environment.InProgressDeploymentID = deployment.ID
	environment.Deployments[deployment.ID] = *deployment

	suite.environment.EXPECT().GetEnvironment(suite.ctx, environment.Name).Return(environment, nil)

	d, err := suite.deployment.GetCurrentDeployment(suite.ctx, environment.Name)
	assert.Nil(suite.T(), err, "Unexpected error when calling GetCurrentDeployment")
	assert.Exactly(suite.T(), deployment.ID, d.ID, "Expected the deployment to match the in-progress deployment")
	assert.Exactly(suite.T(), types.DeploymentInProgress, d.Status, "Expected the deployment status to be in-progress")
}

func (suite *DeploymentTestSuite) TestGetCurrentDeploymentOnlyPendingExists() {
	environment, err := types.NewEnvironment("TestGetCurrent", taskDefinition, cluster1)
	assert.Nil(suite.T(), err, "Unexpected error when creating environment")
	deployment, err := types.NewDeployment(taskDefinition, uuid.NewRandom().String())
	assert.Nil(suite.T(), err, "Unexpected error when creating deployment")
	environment.PendingDeploymentID = deployment.ID
	environment.Deployments[deployment.ID] = *deployment

	suite.environment.EXPECT().GetEnvironment(suite.ctx, environment.Name).Return(environment, nil).Times(2)

	d, err := suite.deployment.GetCurrentDeployment(suite.ctx, environment.Name)
	assert.Nil(suite.T(), err, "Unexpected error when calling GetCurrentDeployment")
	assert.Nil(suite.T(), d, "Did not expect a deployment when there is only a pending deployment available")
}

func (suite *DeploymentTestSuite) TestGetCurrentDeploymentPendingAndCompletedExist() {
	environment, err := types.NewEnvironment("TestGetCurrent", taskDefinition, cluster1)
	assert.Nil(suite.T(), err, "Unexpected error when creating environment")
	pending, err := types.NewDeployment(taskDefinition, uuid.NewRandom().String())
	assert.Nil(suite.T(), err, "Unexpected error when creating deployment")
	environment.PendingDeploymentID = pending.ID
	environment.Deployments[pending.ID] = *pending
	completed, err := types.NewDeployment(taskDefinition, uuid.NewRandom().String())
	assert.Nil(suite.T(), err, "Unexpected error when creating deployment")
	err = completed.UpdateDeploymentToCompleted(nil)
	assert.Nil(suite.T(), err)
	environment.Deployments[completed.ID] = *completed

	suite.environment.EXPECT().GetEnvironment(suite.ctx, environment.Name).Return(environment, nil).Times(2)

	d, err := suite.deployment.GetCurrentDeployment(suite.ctx, environment.Name)
	assert.Nil(suite.T(), err, "Unexpected error when calling GetCurrentDeployment")
	assert.NotNil(suite.T(), d)
	verifyDeployment(suite.T(), completed, d)
}

func (suite *DeploymentTestSuite) TestGetCurrentDeploymentPendingAndInProgressExist() {
	environment, err := types.NewEnvironment("TestGetCurrent", taskDefinition, cluster1)
	assert.Nil(suite.T(), err, "Unexpected error when creating environment")
	pending, err := types.NewDeployment(taskDefinition, uuid.NewRandom().String())
	assert.Nil(suite.T(), err, "Unexpected error when creating deployment")
	environment.PendingDeploymentID = pending.ID
	environment.Deployments[pending.ID] = *pending
	inProgress, err := types.NewDeployment(taskDefinition, uuid.NewRandom().String())
	assert.Nil(suite.T(), err, "Unexpected error when creating deployment")
	err = inProgress.UpdateDeploymentToInProgress(1, nil)
	assert.Nil(suite.T(), err)
	environment.InProgressDeploymentID = inProgress.ID
	environment.Deployments[inProgress.ID] = *inProgress

	suite.environment.EXPECT().GetEnvironment(suite.ctx, environment.Name).Return(environment, nil)

	d, err := suite.deployment.GetCurrentDeployment(suite.ctx, environment.Name)
	assert.Nil(suite.T(), err, "Unexpected error when calling GetCurrentDeployment")
	assert.NotNil(suite.T(), d)
	verifyDeployment(suite.T(), inProgress, d)
}

func (suite *DeploymentTestSuite) TestGetCurrentDeploymentOnlyCompletedExists() {
	environment, err := types.NewEnvironment("TestGetCurrent", taskDefinition, cluster1)
	assert.Nil(suite.T(), err, "Unexpected error when creating environment")
	deployment, err := types.NewDeployment(taskDefinition, uuid.NewRandom().String())
	assert.Nil(suite.T(), err, "Unexpected error when creating deployment")
	deployment.Status = types.DeploymentCompleted
	environment.Deployments[deployment.ID] = *deployment

	suite.environment.EXPECT().GetEnvironment(suite.ctx, environment.Name).Return(environment, nil).Times(2)

	d, err := suite.deployment.GetCurrentDeployment(suite.ctx, environment.Name)
	assert.Nil(suite.T(), err, "Unexpected error when calling GetCurrentDeployment")
	assert.Exactly(suite.T(), deployment.ID, d.ID, "Expected the deployment to match the completed deployment")
}

func (suite *DeploymentTestSuite) TestGetCurrentDeploymentEmpty() {
	suite.environmentObject.PendingDeploymentID = ""
	suite.environmentObject.InProgressDeploymentID = ""
	suite.environmentObject.Deployments = nil

	suite.environment.EXPECT().GetEnvironment(suite.ctx, suite.environmentObject.Name).Return(suite.environmentObject, nil).Times(2)

	d, err := suite.deployment.GetCurrentDeployment(suite.ctx, suite.environmentObject.Name)
	assert.Nil(suite.T(), err, "Unexpected error when calling GetCurrentDeployment")
	assert.Nil(suite.T(), d, "Unexpected deployment when there are no in progress or completed deployments")
}

func (suite *DeploymentTestSuite) TestGetCurrentDeploymentGetEnvironmentReturnsErrors() {
	err := errors.New("Error getting environment")
	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).Return(nil, err)

	_, observedErr := suite.deployment.GetCurrentDeployment(suite.ctx, environmentName)
	assert.Exactly(suite.T(), err, errors.Cause(observedErr))
}

func (suite *DeploymentTestSuite) TestGetCurrentDeploymentGetEnvironmentReturnsEmpty() {
	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).Return(nil, nil)

	_, observedErr := suite.deployment.GetCurrentDeployment(suite.ctx, environmentName)
	_, ok := observedErr.(types.NotFoundError)
	assert.True(suite.T(), ok, "Expected NotFoundError when GetCurrentDeployment is called with name of environment which does not exist")
}

func (suite *DeploymentTestSuite) TestGetCurrentDeploymentMissingName() {
	_, observedErr := suite.deployment.GetCurrentDeployment(suite.ctx, "")
	assert.Error(suite.T(), observedErr, "Expected an error when GetCurrentDeployment is called with a missing environment name")
	_, ok := observedErr.(types.BadRequestError)
	assert.True(suite.T(), ok, "Expected BadRequestError when GetEnvironment is called with a missing environment name")
}

func (suite *DeploymentTestSuite) TestGetInProgressDeploymentMissingName() {
	_, observedErr := suite.deployment.GetInProgressDeployment(suite.ctx, "")
	assert.Error(suite.T(), observedErr, "Expected an error when GetInProgressDeployment is called with a missing environment name")
	_, ok := observedErr.(types.BadRequestError)
	assert.True(suite.T(), ok, "Expected BadRequestError when GetInProgressDeployment is called with a missing environment name")
}

func (suite *DeploymentTestSuite) TestGetInProgressDeploymentGetEnvironmentReturnsErrors() {
	err := errors.New("Error getting environment")
	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).Return(nil, err)

	_, observedErr := suite.deployment.GetInProgressDeployment(suite.ctx, environmentName)
	assert.Exactly(suite.T(), err, errors.Cause(observedErr))
}

func (suite *DeploymentTestSuite) TestGetInProgressDeploymentGetEnvironmentReturnsEmpty() {
	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).Return(nil, nil)

	_, observedErr := suite.deployment.GetInProgressDeployment(suite.ctx, environmentName)
	_, ok := observedErr.(types.NotFoundError)
	assert.True(suite.T(), ok, "Expected NotFoundError when GetInProgressDeployment is called with name of environment which does not exist")
}

func (suite *DeploymentTestSuite) TestGetInProgressDeploymentNoInProgressDeployments() {
	environment, err := types.NewEnvironment("TestGetInProgressDeployment", taskDefinition, cluster1)
	assert.Nil(suite.T(), err, "Unexpected error when creating environment")
	environment.PendingDeploymentID = ""
	environment.InProgressDeploymentID = ""

	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).Return(environment, nil)

	d, err := suite.deployment.GetInProgressDeployment(suite.ctx, environmentName)
	assert.Nil(suite.T(), err, "Did not expect errors when there are no in-progress deployments")
	assert.Nil(suite.T(), d, "There should be no in-progress deployments")
}

func (suite *DeploymentTestSuite) TestGetInProgressDeploymentOnlyPendingExists() {
	environment, err := types.NewEnvironment("TestGetInProgressDeployment", taskDefinition, cluster1)
	assert.Nil(suite.T(), err, "Unexpected error when creating environment")
	environment.InProgressDeploymentID = ""

	deployment, err := types.NewDeployment(taskDefinition, uuid.NewRandom().String())
	assert.Nil(suite.T(), err, "Unexpected error when creating deployment")
	environment.PendingDeploymentID = deployment.ID
	environment.Deployments[deployment.ID] = *deployment

	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).Return(environment, nil)

	d, err := suite.deployment.GetInProgressDeployment(suite.ctx, environmentName)
	assert.Nil(suite.T(), err, "Did not expect errors when there are no in-progress deployments")
	assert.Nil(suite.T(), d, "Did not expect an in-progress deployment")
}

func (suite *DeploymentTestSuite) TestGetInProgressDeploymentMissingDeploymentInEnvironment() {
	environment, err := types.NewEnvironment("TestGetInProgressDeployment", taskDefinition, cluster1)
	assert.Nil(suite.T(), err, "Unexpected error when creating environment")

	deployment, err := types.NewDeployment(taskDefinition, uuid.NewRandom().String())
	assert.Nil(suite.T(), err, "Unexpected error when creating deployment")
	err = deployment.UpdateDeploymentToInProgress(0, nil)
	assert.Nil(suite.T(), err, "Unexpected error when moving deployment to in progress")

	environment.InProgressDeploymentID = deployment.ID

	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).Return(environment, nil)

	_, err = suite.deployment.GetInProgressDeployment(suite.ctx, environmentName)
	assert.Error(suite.T(), err, "Expected an error when the in-progress deployment is not in the deployment map")
}

func (suite *DeploymentTestSuite) TestGetInProgressDeployment() {
	environment, err := types.NewEnvironment("TestGetInProgressDeployment", taskDefinition, cluster1)
	assert.Nil(suite.T(), err, "Unexpected error when creating environment")

	deployment, err := types.NewDeployment(taskDefinition, uuid.NewRandom().String())
	assert.Nil(suite.T(), err, "Unexpected error when creating deployment")
	err = deployment.UpdateDeploymentToInProgress(0, nil)
	assert.Nil(suite.T(), err, "Unexpected error when moving deployment to in progress")

	environment.InProgressDeploymentID = deployment.ID
	environment.Deployments[deployment.ID] = *deployment

	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).Return(environment, nil)

	d, err := suite.deployment.GetInProgressDeployment(suite.ctx, environmentName)
	assert.Nil(suite.T(), err, "Did not expect an error when getting in-progress deployment")
	assert.Exactly(suite.T(), deployment, d, "Expected deployments to match")
}

func (suite *DeploymentTestSuite) TestListDeploymentsSortedReverseChronologicallyMissingName() {
	_, observedErr := suite.deployment.ListDeploymentsSortedReverseChronologically(suite.ctx, "")
	assert.Error(suite.T(), observedErr, "Expected an error when ListDeploymentsSortedReverseChronologically is called with a missing environment name")
	_, ok := observedErr.(types.BadRequestError)
	assert.True(suite.T(), ok, "Expected BadRequestError when ListDeploymentsSortedReverseChronologically is called with a missing environment name")
}

func (suite *DeploymentTestSuite) TestListDeploymentsSortedReverseChronologicallyGetEnvironmentReturnsErrors() {
	err := errors.New("Error getting environment")
	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).Return(nil, err)

	_, observedErr := suite.deployment.ListDeploymentsSortedReverseChronologically(suite.ctx, environmentName)
	assert.Exactly(suite.T(), err, errors.Cause(observedErr))
}

func (suite *DeploymentTestSuite) TestListDeploymentsSortedReverseChronologicallyGetEnvironmentReturnsEmpty() {
	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).Return(nil, nil)

	_, observedErr := suite.deployment.ListDeploymentsSortedReverseChronologically(suite.ctx, environmentName)
	_, ok := observedErr.(types.NotFoundError)
	assert.True(suite.T(), ok, "Expected NotFoundError when ListDeploymentsSortedReverseChronologically is called with name of environment which does not exist")
}

func (suite *DeploymentTestSuite) TestListDeploymentsSortedReverseChronologically() {
	environment, err := types.NewEnvironment("TestGetDeployments", taskDefinition, cluster1)
	assert.Nil(suite.T(), err, "Unexpected error when creating environment")

	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).Return(environment, nil)

	deployment1, err := types.NewDeployment(taskDefinition, uuid.NewRandom().String())
	assert.Nil(suite.T(), err, "Unexpected error when creating deployment")

	deployment2, err := types.NewDeployment(taskDefinition, uuid.NewRandom().String())
	assert.Nil(suite.T(), err, "Unexpected error when creating deployment")
	deployment2.StartTime = deployment1.StartTime.Add(time.Minute)

	environment.Deployments[deployment1.ID] = *deployment1
	environment.Deployments[deployment2.ID] = *deployment2

	deployments, err := suite.deployment.ListDeploymentsSortedReverseChronologically(suite.ctx, environmentName)
	assert.Nil(suite.T(), err, "Did not expect errors when getting deployments")
	assert.Exactly(suite.T(), *deployment2, deployments[0], "Expected deployments to match")
	assert.Exactly(suite.T(), *deployment1, deployments[1], "Expected deployments to match")
}

func (suite *DeploymentTestSuite) TestGetPendingDeploymentMissingName() {
	_, observedErr := suite.deployment.GetPendingDeployment(suite.ctx, "")
	assert.Error(suite.T(), observedErr, "Expected an error when GetPendingDeployment is called with a missing environment name")
	_, ok := observedErr.(types.BadRequestError)
	assert.True(suite.T(), ok, "Expected BadRequestError when GetPendingDeployment is called with a missing environment name")
}

func (suite *DeploymentTestSuite) TestGetPendingDeploymentGetEnvironmentReturnsErrors() {
	err := errors.New("Error getting environment")
	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).Return(nil, err)

	_, observedErr := suite.deployment.GetPendingDeployment(suite.ctx, environmentName)
	assert.Exactly(suite.T(), err, errors.Cause(observedErr))
}

func (suite *DeploymentTestSuite) TestGetPendingDeploymentGetEnvironmentReturnsEmpty() {
	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).Return(nil, nil)

	_, observedErr := suite.deployment.GetPendingDeployment(suite.ctx, environmentName)
	_, ok := observedErr.(types.NotFoundError)
	assert.True(suite.T(), ok, "Expected NotFoundError when GetPendingDeployment is called with name of environment which does not exist")
}

func (suite *DeploymentTestSuite) TestGetPendingDeploymentNoPendingDeployments() {
	environment, err := types.NewEnvironment("TestGetPendingDeployment", taskDefinition, cluster1)
	assert.Nil(suite.T(), err, "Unexpected error when creating environment")
	environment.PendingDeploymentID = ""
	environment.InProgressDeploymentID = ""

	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).Return(environment, nil)

	d, err := suite.deployment.GetPendingDeployment(suite.ctx, environmentName)
	assert.Nil(suite.T(), err, "Did not expect errors when there are no pending deployments")
	assert.Nil(suite.T(), d, "There should be no pending deployments")
}

func (suite *DeploymentTestSuite) TestGetPendingDeploymentOnlyInProgressExists() {
	environment, err := types.NewEnvironment("TestGetPendingDeployment", taskDefinition, cluster1)
	assert.Nil(suite.T(), err, "Unexpected error when creating environment")
	environment.PendingDeploymentID = ""

	deployment, err := types.NewDeployment(taskDefinition, uuid.NewRandom().String())
	assert.Nil(suite.T(), err, "Unexpected error when creating deployment")
	err = deployment.UpdateDeploymentToInProgress(1, nil)
	assert.Nil(suite.T(), err)
	environment.InProgressDeploymentID = deployment.ID
	environment.Deployments[deployment.ID] = *deployment

	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).Return(environment, nil)

	d, err := suite.deployment.GetPendingDeployment(suite.ctx, environmentName)
	assert.Nil(suite.T(), err, "Did not expect errors when there are no pending deployments")
	assert.Nil(suite.T(), d, "Did not expect an pending deployment")
}

func (suite *DeploymentTestSuite) TestGetPendingDeploymentMissingDeploymentInEnvironment() {
	environment, err := types.NewEnvironment("TestGetPendingDeployment", taskDefinition, cluster1)
	assert.Nil(suite.T(), err, "Unexpected error when creating environment")

	environment.PendingDeploymentID = uuid.NewUUID().String()

	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).Return(environment, nil)

	_, err = suite.deployment.GetPendingDeployment(suite.ctx, environmentName)
	assert.Error(suite.T(), err, "Expected an error when the pending deployment is not in the deployment map")
}

func (suite *DeploymentTestSuite) TestGetPendingDeployment() {
	environment, err := types.NewEnvironment("TestGetPendingDeployment", taskDefinition, cluster1)
	assert.Nil(suite.T(), err, "Unexpected error when creating environment")

	deployment, err := types.NewDeployment(taskDefinition, uuid.NewRandom().String())
	assert.Nil(suite.T(), err, "Unexpected error when creating deployment")

	environment.PendingDeploymentID = deployment.ID
	environment.Deployments[deployment.ID] = *deployment

	suite.environment.EXPECT().GetEnvironment(suite.ctx, environmentName).Return(environment, nil)

	d, err := suite.deployment.GetPendingDeployment(suite.ctx, environmentName)
	assert.Nil(suite.T(), err, "Did not expect an error when getting pending deployment")
	assert.Exactly(suite.T(), deployment, d, "Expected deployments to match")
}
