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

package types

import (
	"testing"
	"time"

	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type EnvironmentTestSuite struct {
	suite.Suite
	environment *Environment
	deployment  *Deployment
}

func (suite *EnvironmentTestSuite) SetupTest() {
	var err error
	suite.environment, err = NewEnvironment(environmentName, taskDefinition, cluster)
	assert.Nil(suite.T(), err, "Cannot initialize EnvironmentTestSuite")
	suite.deployment, err = NewDeployment(taskDefinition, generateToken())
	assert.Nil(suite.T(), err, "Unexpected error when creating deployment")
}

func TestEnvironmentTestSuite(t *testing.T) {
	suite.Run(t, new(EnvironmentTestSuite))
}

func (suite *EnvironmentTestSuite) TestNewEnvironmentEmptyName() {
	_, err := NewEnvironment("", taskDefinition, cluster)
	assert.Error(suite.T(), err, "Expected an error when name is empty")
}

func (suite *EnvironmentTestSuite) TestNewEnvironmentEmptyTaskDefinition() {
	_, err := NewEnvironment(environmentName, "", cluster)
	assert.Error(suite.T(), err, "Expected an error when taskDefinition is empty")
}

func (suite *EnvironmentTestSuite) TestNewEnvironmentEmptyCluster() {
	_, err := NewEnvironment(environmentName, taskDefinition, "")
	assert.Error(suite.T(), err, "Expected an error when cluster is empty")
}

func (suite *EnvironmentTestSuite) TestNewEnvironment() {
	e, err := NewEnvironment(environmentName, taskDefinition, cluster)
	assert.Nil(suite.T(), err, "Unexpected error when creating a new environment")
	assert.NotNil(suite.T(), e, "Environment should not be nil")
	assert.NotEmpty(suite.T(), e.Token, "Token should not be empty")
	assert.Exactly(suite.T(), environmentName, e.Name, "Name should match")
	assert.Exactly(suite.T(), taskDefinition, e.DesiredTaskDefinition, "TaskDefinition should match")
	assert.Exactly(suite.T(), 0, e.DesiredTaskCount, "Task count should match")
	assert.Exactly(suite.T(), cluster, e.Cluster, "Cluster should match")
	assert.Exactly(suite.T(), EnvironmentHealthy, e.Health, "Should be healthy")
	assert.Empty(suite.T(), e.Deployments, "Deployments should be empty")
}

func (suite *EnvironmentTestSuite) TestSortDeploymentsReverseChronologically() {
	deployment1, err := NewDeployment(taskDefinition, generateToken())
	assert.Nil(suite.T(), err, "Unexpected error when creating deployment")

	deployment2, err := NewDeployment(taskDefinition, generateToken())
	assert.Nil(suite.T(), err, "Unexpected error when creating deployment")
	deployment2.StartTime = deployment1.StartTime.Add(time.Minute)

	suite.environment.Deployments[deployment1.ID] = *deployment1
	suite.environment.Deployments[deployment2.ID] = *deployment2

	deployments, err := suite.environment.SortDeploymentsReverseChronologically()
	assert.Nil(suite.T(), err, "Unexpected error when sorting deployments")
	assert.Exactly(suite.T(), *deployment2, deployments[0], "Expected the deployments to match")
	assert.Exactly(suite.T(), *deployment1, deployments[1], "Expected the deployments to match")
}

func (suite *EnvironmentTestSuite) TestAddPendingDeploymentStatusNotPending() {
	suite.deployment.Status = DeploymentInProgress

	err := suite.environment.AddPendingDeployment(*suite.deployment)
	assert.Error(suite.T(), err, "Expected an error when the deployment status is not pending")
}

func (suite *EnvironmentTestSuite) TestAddPendingDeployment() {
	err := suite.environment.AddPendingDeployment(*suite.deployment)
	assert.Nil(suite.T(), err, "Unexpected error when adding a pending deployment")
	assert.Exactly(suite.T(), *suite.deployment, suite.environment.Deployments[suite.deployment.ID], "")
	assert.Exactly(suite.T(), suite.deployment.ID, suite.environment.PendingDeploymentID, "")
}

func (suite *EnvironmentTestSuite) TestUpdatePendingDeploymentToInProgressNoPendingDeployment() {
	suite.environment.PendingDeploymentID = ""
	err := suite.environment.UpdatePendingDeploymentToInProgress()
	assert.Error(suite.T(), err, "Expected an error when there is no pending deployment")
}

func (suite *EnvironmentTestSuite) TestUpdatePendingDeploymentToInProgressPendingDeploymentNotInMap() {
	suite.environment.PendingDeploymentID = generateToken()

	err := suite.environment.UpdatePendingDeploymentToInProgress()
	assert.Error(suite.T(), err, "")
}

func (suite *EnvironmentTestSuite) TestUpdatePendingDeploymentToInProgressPendingDeploymentStatusIncorrect() {
	err := suite.environment.AddPendingDeployment(*suite.deployment)
	assert.Nil(suite.T(), err, "Unexpected error when adding a pending deployment")
	d := suite.environment.Deployments[suite.deployment.ID]
	d.Status = DeploymentCompleted
	suite.environment.Deployments[d.ID] = d

	err = suite.environment.UpdatePendingDeploymentToInProgress()
	assert.Error(suite.T(), err, "")
}

func (suite *EnvironmentTestSuite) TestUpdatePendingDeploymentToInProgress() {
	err := suite.environment.AddPendingDeployment(*suite.deployment)
	assert.Nil(suite.T(), err, "Unexpected error when adding a pending deployment")

	err = suite.environment.UpdatePendingDeploymentToInProgress()
	assert.Nil(suite.T(), err, "")
	assert.Empty(suite.T(), suite.environment.PendingDeploymentID, "")
	assert.Exactly(suite.T(), suite.deployment.ID, suite.environment.InProgressDeploymentID, "")
	assert.Exactly(suite.T(), *suite.deployment, suite.environment.Deployments[suite.environment.InProgressDeploymentID], "")
}

func generateToken() string {
	return uuid.NewRandom().String()
}
