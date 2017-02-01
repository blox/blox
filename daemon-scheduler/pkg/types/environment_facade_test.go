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
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/blox/blox/cluster-state-service/swagger/v1/generated/models"
	"github.com/blox/blox/daemon-scheduler/pkg/facade"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type EnvironmentFacadeTestSuite struct {
	suite.Suite
	css               *facade.MockClusterState
	environment       *Environment
	environmentFacade EnvironmentFacade
}

func (suite *EnvironmentFacadeTestSuite) SetupTest() {
	mockCtrl := gomock.NewController(suite.T())
	var err error
	suite.css = facade.NewMockClusterState(mockCtrl)
	suite.environment, err = NewEnvironment(environmentName, taskDefinition, cluster)
	assert.Nil(suite.T(), err)
	suite.environmentFacade, err = NewEnvironmentFacade(suite.css)
	assert.Nil(suite.T(), err)
}

func TestEnvironmentFacadeTestSuite(t *testing.T) {
	suite.Run(t, new(EnvironmentFacadeTestSuite))
}

func (suite *EnvironmentFacadeTestSuite) TestNewEnvironmentFacadeNilCss() {
	_, err := NewEnvironmentFacade(nil)
	assert.Error(suite.T(), err)
}

func (suite *EnvironmentFacadeTestSuite) TestNewEnvironmentFacade() {
	f, err := NewEnvironmentFacade(suite.css)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), f)
}

func (suite *EnvironmentFacadeTestSuite) TestEnvironmentFacadeInstanceARNsEmptyClusterName() {
	suite.environment.Cluster = ""
	_, err := suite.environmentFacade.InstanceARNs(suite.environment)
	assert.Error(suite.T(), err)
}

func (suite *EnvironmentFacadeTestSuite) TestEnvironmentFacadeInstanceARNsListInstancesFails() {
	suite.css.EXPECT().ListInstances(suite.environment.Cluster).Return(nil, errors.New("List instances fails"))

	_, err := suite.environmentFacade.InstanceARNs(suite.environment)
	assert.Error(suite.T(), err)
}

func (suite *EnvironmentFacadeTestSuite) TestEnvironmentFacadeInstanceARNsListInstancesNil() {
	suite.css.EXPECT().ListInstances(suite.environment.Cluster).Return(nil, nil)

	instances, err := suite.environmentFacade.InstanceARNs(suite.environment)
	assert.Nil(suite.T(), err)
	assert.Empty(suite.T(), instances)
}

func (suite *EnvironmentFacadeTestSuite) TestEnvironmentFacadeInstanceARNsListInstancesEmpty() {
	suite.css.EXPECT().ListInstances(suite.environment.Cluster).Return([]*models.ContainerInstance{}, nil)

	instances, err := suite.environmentFacade.InstanceARNs(suite.environment)
	assert.Nil(suite.T(), err)
	assert.Empty(suite.T(), instances)
}

func (suite *EnvironmentFacadeTestSuite) TestEnvironmentFacadeInstanceARNs() {
	listInstancesResponse := make([]*models.ContainerInstance, 2)
	listInstancesResponse[0] = &models.ContainerInstance{
		ContainerInstanceARN: aws.String(instanceARN1),
	}
	listInstancesResponse[1] = &models.ContainerInstance{
		ContainerInstanceARN: aws.String(instanceARN2),
	}

	suite.css.EXPECT().ListInstances(suite.environment.Cluster).Return(listInstancesResponse, nil)
	expectedInstanceARNs := []*string{aws.String(instanceARN1), aws.String(instanceARN2)}

	instances, err := suite.environmentFacade.InstanceARNs(suite.environment)
	assert.Nil(suite.T(), err)
	assert.Exactly(suite.T(), expectedInstanceARNs, instances)
}
