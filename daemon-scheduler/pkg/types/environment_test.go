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

package types

import (
	"testing"
	"time"

	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	environmentName = "environmentName"
)

type EnvironmentTestSuite struct {
	suite.Suite
	environment *Environment
}

func (suite *EnvironmentTestSuite) SetupTest() {
	var err error
	suite.environment, err = NewEnvironment(environmentName, taskDefinition, cluster)
	assert.Nil(suite.T(), err, "Cannot initialize EnvironmentTestSuite")
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

func generateToken() string {
	return uuid.NewRandom().String()
}
