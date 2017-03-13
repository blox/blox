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

	"github.com/blox/blox/daemon-scheduler/pkg/environment/types"
	"github.com/blox/blox/daemon-scheduler/pkg/json"
	"github.com/blox/blox/daemon-scheduler/pkg/mocks"
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
	txstore          *mocks.MockEtcdTXStore
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
	suite.txstore = mocks.NewMockEtcdTXStore(mockCtrl)

	suite.ctx = context.TODO()
	suite.environment1, err = types.NewEnvironment(environmentName1, taskDefinition, cluster)
	assert.Nil(suite.T(), err, "Cannot initialize EnvironmentTestSuite")
	suite.environment1JSON, err = json.MarshalJSON(suite.environment1)
	assert.Nil(suite.T(), err, "Cannot initialize EnvironmentTestSuite")

	suite.environment2, err = types.NewEnvironment(environmentName2, taskDefinition, cluster)
	assert.Nil(suite.T(), err, "Cannot initialize EnvironmentTestSuite")
	suite.environment2JSON, err = json.MarshalJSON(suite.environment2)
	assert.Nil(suite.T(), err, "Cannot initialize EnvironmentTestSuite")

	suite.environmentStore, err = NewEnvironmentStore(suite.datastore, suite.txstore)
	assert.Nil(suite.T(), err, "Cannot initialize EnvironmentTestSuite")
}

func TestEnvironmentTestSuite(t *testing.T) {
	suite.Run(t, new(EnvironmentTestSuite))
}

func (suite *EnvironmentTestSuite) TestNewEnvironmentStoreEmptyDataStore() {
	_, err := NewEnvironmentStore(nil, suite.txstore)
	assert.Error(suite.T(), err, "Expected an error when datastore is nil")
}

func (suite *EnvironmentTestSuite) TestNewEnvironmentStoreEmptyEtcdTransactionalStore() {
	_, err := NewEnvironmentStore(suite.datastore, nil)
	assert.Error(suite.T(), err, "Expected an error when etcd transactional store is nil")
}

func (suite *EnvironmentTestSuite) TestNewEnvironmentStore() {
	es, err := NewEnvironmentStore(suite.datastore, suite.txstore)
	assert.Nil(suite.T(), err, "Unexpected error when datastore and etcd transactional store are not nil")
	assert.NotNil(suite.T(), es, "Environment store should not be nil")
}

func (suite *EnvironmentTestSuite) TestPutEnvironmentWithNoName() {
	f := func(env *types.Environment) (*types.Environment, error) {
		return nil, nil
	}
	err := suite.environmentStore.PutEnvironment(suite.ctx, "", f)
	assert.Error(suite.T(), err, "Expected an error when environment name is missing")
}

func (suite *EnvironmentTestSuite) TestPutEnvironmentSTMRepeatableFails() {
	f := func(env *types.Environment) (*types.Environment, error) {
		return nil, nil
	}

	suite.txstore.EXPECT().GetV3Client().Return(nil)
	suite.txstore.EXPECT().NewSTMRepeatable(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("STMRepeatable error"))

	err := suite.environmentStore.PutEnvironment(suite.ctx, suite.environment1.Name, f)
	assert.Error(suite.T(), err, "Expected an error when STMRepeatable fails")
}

func (suite *EnvironmentTestSuite) TestPutEnvironment() {
	f := func(env *types.Environment) (*types.Environment, error) {
		return nil, nil
	}

	suite.txstore.EXPECT().GetV3Client().Return(nil)
	suite.txstore.EXPECT().NewSTMRepeatable(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)

	err := suite.environmentStore.PutEnvironment(suite.ctx, suite.environment1.Name, f)
	assert.Nil(suite.T(), err, "Unexpeced error putting environment")
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

func (suite *EnvironmentTestSuite) TestDeleteEnvironmentNoName() {
	suite.datastore.EXPECT().Delete(gomock.Any(), gomock.Any()).Times(0)

	err := suite.environmentStore.DeleteEnvironment(suite.ctx, "")
	assert.Error(suite.T(), err, "Expected an error when deleting an environment with no name")
}

func (suite *EnvironmentTestSuite) TestDeleteEnvironmentStoreReturnsError() {
	suite.datastore.EXPECT().Delete(suite.ctx, gomock.Any()).Return(errors.New("Delete failed"))

	err := suite.environmentStore.DeleteEnvironment(suite.ctx, suite.environment1.Name)
	assert.Error(suite.T(), err, "Expected an error when datastore delete returned an error")
}

func (suite *EnvironmentTestSuite) TestDeleteEnvironment() {
	suite.datastore.EXPECT().Delete(suite.ctx, gomock.Any()).Return(nil)

	err := suite.environmentStore.DeleteEnvironment(suite.ctx, suite.environment1.Name)
	assert.Nil(suite.T(), err, "Unexpected error deleting an environment")
}
