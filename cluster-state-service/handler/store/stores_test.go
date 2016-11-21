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

package store

import (
	"testing"

	"github.com/blox/blox/cluster-state-service/handler/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type StoreTestSuite struct {
	suite.Suite
	datastore   *mocks.MockDataStore
	etcdTxStore *mocks.MockEtcdTXStore
}

func (testSuite *StoreTestSuite) SetupTest() {
	mockCtrl := gomock.NewController(testSuite.T())
	testSuite.datastore = mocks.NewMockDataStore(mockCtrl)
	testSuite.etcdTxStore = mocks.NewMockEtcdTXStore(mockCtrl)
}

func TestStoreTestSuite(t *testing.T) {
	suite.Run(t, new(StoreTestSuite))
}

func (testSuite *StoreTestSuite) TestNewStoresDatastoreNil() {
	_, err := NewStores(nil, testSuite.etcdTxStore)
	assert.Error(testSuite.T(), err, "Expected an error when NewStores is initialized with nil datastore")
}

func (testSuite *StoreTestSuite) TestNewStoresEtcdTxStoreNil() {
	_, err := NewStores(testSuite.datastore, nil)
	assert.Error(testSuite.T(), err, "Expected an error when NewStores is initialized with nil etcd transaction store")
}

func (testSuite *StoreTestSuite) TestNewStores() {
	stores, err := NewStores(testSuite.datastore, testSuite.etcdTxStore)
	assert.Nil(testSuite.T(), err, "Unexpected error when calling NewStores")
	assert.NotNil(testSuite.T(), stores, "Stores should not be nil")
	assert.NotNil(testSuite.T(), stores.TaskStore, "TaskStore should not be nil")
	assert.NotNil(testSuite.T(), stores.ContainerInstanceStore, "ContainerInstanceStores should not be nil")
}
