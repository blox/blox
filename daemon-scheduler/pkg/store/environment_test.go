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

package store

import (
	"context"
	"testing"

	"github.com/blox/blox/daemon-scheduler/pkg/json"
	"github.com/blox/blox/daemon-scheduler/pkg/mocks"
	"github.com/blox/blox/daemon-scheduler/pkg/types"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	environmentName1 = "environmentName1"
	environmentName2 = "environmentName2"
	environmentKey1  = environmentKeyPrefix + environmentName1
	environmentKey2  = environmentKeyPrefix + environmentName2
	taskDefinition   = "arn:aws:ecs:us-east-1:12345678912:task-definition/test"
	cluster          = "arn:aws:ecs:us-east-1:123456789123:cluster/test"
)

type EnvironmentTestSuite struct {
	suite.Suite
	datastore        *mocks.MockDataStore
	environmentStore EnvironmentStore
	ctx              context.Context
	environment1     *types.Environment
	environment2     *types.Environment
	environment1JSON string
	environment2JSON string
}

func (suite *EnvironmentTestSuite) SetupTest() {
	mockCtrl := gomock.NewController(suite.T())
	var err error

	suite.datastore = mocks.NewMockDataStore(mockCtrl)
	suite.ctx = context.TODO()
	suite.environment1, err = types.NewEnvironment(environmentName1, taskDefinition, cluster)
	assert.Nil(suite.T(), err, "Cannot initialize EnvironmentTestSuite")
	suite.environment1JSON, err = json.MarshalJSON(suite.environment1)
	assert.Nil(suite.T(), err, "Cannot initialize EnvironmentTestSuite")
	suite.environment2, err = types.NewEnvironment(environmentName2, taskDefinition, cluster)
	assert.Nil(suite.T(), err, "Cannot initialize EnvironmentTestSuite")
	suite.environment2JSON, err = json.MarshalJSON(suite.environment2)
	assert.Nil(suite.T(), err, "Cannot initialize EnvironmentTestSuite")
	suite.environmentStore, err = NewEnvironmentStore(suite.datastore)
	assert.Nil(suite.T(), err, "Cannot initialize EnvironmentTestSuite")
}

func TestEnvironmentTestSuite(t *testing.T) {
	suite.Run(t, new(EnvironmentTestSuite))
}

func (suite *EnvironmentTestSuite) TestNewEnvironmentStoreEmptyDataStore() {
	_, err := NewEnvironmentStore(nil)
	assert.Error(suite.T(), err, "Expected an error when datastore is nil")
}

func (suite *EnvironmentTestSuite) TestNewEnvironmentStore() {
	es, err := NewEnvironmentStore(suite.datastore)
	assert.Nil(suite.T(), err, "Unexpected error when datastore is not nil")
	assert.NotNil(suite.T(), es, "Environment store should not be nil")
}

func (suite *EnvironmentTestSuite) TestPutWithMissingEnvironmentName() {
	err := suite.environmentStore.PutEnvironment(suite.ctx, types.Environment{})
	assert.Error(suite.T(), err, "Expected an error when environment name is missing")
}

func (suite *EnvironmentTestSuite) TestPutDataStorePutFails() {
	suite.datastore.EXPECT().Put(suite.ctx, environmentKey1, suite.environment1JSON).
		Return(errors.New("Put failed"))

	err := suite.environmentStore.PutEnvironment(suite.ctx, *suite.environment1)
	assert.Error(suite.T(), err, "Expected an error when datastore put fails")
}

func (suite *EnvironmentTestSuite) TestPut() {
	suite.datastore.EXPECT().Put(suite.ctx, environmentKey1, suite.environment1JSON).
		Return(nil)

	err := suite.environmentStore.PutEnvironment(suite.ctx, *suite.environment1)
	assert.Nil(suite.T(), err, "Unexpected error when datastore put succeeds")
}

func (suite *EnvironmentTestSuite) TestGetWithMissingEnvironmentName() {
	_, err := suite.environmentStore.GetEnvironment(suite.ctx, "")
	assert.Error(suite.T(), err, "Expected an error when environment name is missing")
}

func (suite *EnvironmentTestSuite) TestGetDataStoreGetFails() {
	suite.datastore.EXPECT().Get(suite.ctx, environmentKey1).
		Return(nil, errors.New("Get failed"))

	_, err := suite.environmentStore.GetEnvironment(suite.ctx, suite.environment1.Name)
	assert.Error(suite.T(), err, "Expected an error when datastore get fails")
}

func (suite *EnvironmentTestSuite) TestGetDataStoreGetEmpty() {
	resp := make(map[string]string)
	suite.datastore.EXPECT().Get(suite.ctx, environmentKey1).Return(resp, nil)

	env, err := suite.environmentStore.GetEnvironment(suite.ctx, suite.environment1.Name)
	assert.Nil(suite.T(), err, "Unexpected error when datastore get is empty")
	assert.Nil(suite.T(), env, "Expected nil when datastore get is empty")
}

func (suite *EnvironmentTestSuite) TestGetDataStoreGetMultipleResults() {
	resp := map[string]string{
		environmentKey1: suite.environment1JSON,
		environmentKey2: suite.environment1JSON,
	}
	suite.datastore.EXPECT().Get(suite.ctx, environmentKey1).Return(resp, nil)

	env, err := suite.environmentStore.GetEnvironment(suite.ctx, suite.environment1.Name)
	assert.Error(suite.T(), err, "Expected an error when multiple results are returned from the datastore")
	assert.Nil(suite.T(), env, "Expected nil when multiple results are returned from the datastore")
}

func (suite *EnvironmentTestSuite) TestGetDataStoreInvalidJson() {
	resp := map[string]string{
		environmentKey1: "invalidJSON",
	}
	suite.datastore.EXPECT().Get(suite.ctx, environmentKey1).Return(resp, nil)

	env, err := suite.environmentStore.GetEnvironment(suite.ctx, suite.environment1.Name)
	assert.Error(suite.T(), err, "Expected an error when get returns invalid json")
	assert.Nil(suite.T(), env, "Expected nil when when get returns invalid json")
}

func (suite *EnvironmentTestSuite) TestGetDataStore() {
	resp := map[string]string{
		environmentKey1: suite.environment1JSON,
	}
	suite.datastore.EXPECT().Get(suite.ctx, environmentKey1).Return(resp, nil)

	env, err := suite.environmentStore.GetEnvironment(suite.ctx, suite.environment1.Name)
	assert.Nil(suite.T(), err, "Unexpected error when retrieving results")
	assert.Exactly(suite.T(), suite.environment1, env,
		"Expected the returned environment to be the same as the one returned by get")
}

func (suite *EnvironmentTestSuite) TestListGetWithPrefixFails() {
	suite.datastore.EXPECT().GetWithPrefix(suite.ctx, environmentKeyPrefix).
		Return(nil, errors.New("GetWithPrefix failed"))

	_, err := suite.environmentStore.ListEnvironments(suite.ctx)
	assert.Error(suite.T(), err, "Expected an error when datastore get with prefix fails")
}

func (suite *EnvironmentTestSuite) TestListGetWithPrefixInvalidJson() {
	resp := map[string]string{
		environmentKey1: "invalidJSON",
	}
	suite.datastore.EXPECT().GetWithPrefix(suite.ctx, environmentKeyPrefix).Return(resp, nil)

	envs, err := suite.environmentStore.ListEnvironments(suite.ctx)
	assert.Error(suite.T(), err, "Expected an error when getwithprefix returns invalid json")
	assert.Nil(suite.T(), envs, "Expected nil when when getwithprefix returns invalid json")
}

func (suite *EnvironmentTestSuite) TestList() {
	resp := map[string]string{
		environmentKey1: suite.environment1JSON,
		environmentKey2: suite.environment2JSON,
	}
	suite.datastore.EXPECT().GetWithPrefix(suite.ctx, environmentKeyPrefix).Return(resp, nil)

	envs, err := suite.environmentStore.ListEnvironments(suite.ctx)
	assert.Nil(suite.T(), err, "Unexpected error when listing environments")

	expectedEnvs := []types.Environment{*suite.environment1, *suite.environment2}
	assert.Equal(suite.T(), len(expectedEnvs), len(envs), "Expected listed environments(len=%d) to be of same count as what's returned from the store(len=%d)", len(envs), len(expectedEnvs))
	for _, expectedEnv := range expectedEnvs {
		assert.Contains(suite.T(), envs, expectedEnv, "Expected %s to be returned by ListEnvironments", expectedEnv)
	}
}
