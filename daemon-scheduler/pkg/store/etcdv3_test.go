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
	"context"
	"testing"

	"github.com/blox/blox/daemon-scheduler/pkg/mocks"
	etcd "github.com/coreos/etcd/clientv3"
	mvccpb "github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	key          = "key"
	anotherKey   = key + "suffix"
	value        = "value"
	anotherValue = "anotherValue"
)

type DataStoreTestSuite struct {
	suite.Suite
	etcdInterface *mocks.MockEtcdInterface
	datastore     DataStore
	ctx           context.Context
}

func (testSuite *DataStoreTestSuite) SetupTest() {
	mockCtrl := gomock.NewController(testSuite.T())
	testSuite.etcdInterface = mocks.NewMockEtcdInterface(mockCtrl)

	var err error
	testSuite.datastore, err = NewDataStore(testSuite.etcdInterface)
	assert.Nil(testSuite.T(), err, "Cannot initialize DataStoreTestSuite")

	testSuite.ctx = context.TODO()
}

func TestDataStoreTestSuite(t *testing.T) {
	suite.Run(t, new(DataStoreTestSuite))
}

func (testSuite *DataStoreTestSuite) TestNewDataStoreEmptyEtcd() {
	_, err := NewDataStore(nil)
	assert.Error(testSuite.T(), err, "Expected an error when etcd client is nil")
}

func (testSuite *DataStoreTestSuite) TestNewDataStore() {
	ds, err := NewDataStore(testSuite.etcdInterface)
	assert.Nil(testSuite.T(), err, "Unexpected error when etcd client is not nil")
	assert.NotNil(testSuite.T(), ds, "Datastore should not be nil")
}

func (testSuite *DataStoreTestSuite) TestPutEmptyKey() {
	err := testSuite.datastore.Put(testSuite.ctx, "", value)
	assert.Error(testSuite.T(), err, "Expected an error when key is nil")
}

func (testSuite *DataStoreTestSuite) TestPutEmptyValue() {
	err := testSuite.datastore.Put(testSuite.ctx, key, "")
	assert.Error(testSuite.T(), err, "Expected an error when value is nil")
}

func (testSuite *DataStoreTestSuite) TestPutEtcdPutFails() {
	testSuite.etcdInterface.EXPECT().Put(testSuite.ctx, key, value).Return(nil, errors.New("Put failed"))

	err := testSuite.datastore.Put(testSuite.ctx, key, value)
	assert.Error(testSuite.T(), err, "Expected an error when etcd put fails")
}

func (testSuite *DataStoreTestSuite) TestPut() {
	testSuite.etcdInterface.EXPECT().Put(testSuite.ctx, key, value).Return(nil, nil)

	err := testSuite.datastore.Put(testSuite.ctx, key, value)
	assert.Nil(testSuite.T(), err, "Unexpected error when calling adding key %v, value %v", key, value)
}

func (testSuite *DataStoreTestSuite) TestGetWithPrefixEmptyKey() {
	_, err := testSuite.datastore.GetWithPrefix(testSuite.ctx, "")
	assert.Error(testSuite.T(), err, "Expected an error when key is nil")
}

func (testSuite *DataStoreTestSuite) TestGetWithPrefixEtcdGetFails() {
	testSuite.etcdInterface.EXPECT().Get(testSuite.ctx, key, gomock.Any()).Return(nil, errors.New("Get failed"))

	_, err := testSuite.datastore.GetWithPrefix(testSuite.ctx, key)
	assert.Error(testSuite.T(), err, "Expected an error when etcd get fails")
}

func (testSuite *DataStoreTestSuite) TestGetWithPrefixEtcdGetRespNil() {
	testSuite.etcdInterface.EXPECT().Get(testSuite.ctx, key, gomock.Any()).Return((*etcd.GetResponse)(nil), nil)

	resp, err := testSuite.datastore.GetWithPrefix(testSuite.ctx, key)
	assert.Nil(testSuite.T(), err, "Unexpected error when etcd get returns empty")
	assert.Empty(testSuite.T(), resp, "Expected an empty map")
}

func (testSuite *DataStoreTestSuite) TestGetWithPrefixEtcdGetRespKVNil() {
	var getResp etcd.GetResponse
	getResp.Kvs = nil
	testSuite.etcdInterface.EXPECT().Get(testSuite.ctx, key, gomock.Any()).Return(&getResp, nil)

	resp, err := testSuite.datastore.GetWithPrefix(testSuite.ctx, key)
	assert.Nil(testSuite.T(), err, "Unexpected error when etcd get returns empty")
	assert.Empty(testSuite.T(), resp, "Expected an empty map")
}

func (testSuite *DataStoreTestSuite) TestGetWithPrefixEtcdGetRespKVEmpty() {
	var getResp etcd.GetResponse
	getResp.Kvs = []*mvccpb.KeyValue{}
	testSuite.etcdInterface.EXPECT().Get(testSuite.ctx, key, gomock.Any()).Return(&getResp, nil)

	resp, err := testSuite.datastore.GetWithPrefix(testSuite.ctx, key)
	assert.Nil(testSuite.T(), err, "Unexpected error when etcd get returns empty")
	assert.Empty(testSuite.T(), resp, "Expected an empty map")
}

func (testSuite *DataStoreTestSuite) TestGetWithPrefixEtcd() {
	var getResp etcd.GetResponse
	getResp.Kvs = make([]*mvccpb.KeyValue, 2)
	getResp.Kvs[0] = &mvccpb.KeyValue{
		Key:   []byte(key),
		Value: []byte(value),
	}
	getResp.Kvs[1] = &mvccpb.KeyValue{
		Key:   []byte(anotherKey),
		Value: []byte(anotherValue),
	}
	testSuite.etcdInterface.EXPECT().Get(testSuite.ctx, key, gomock.Any()).Return(&getResp, nil)

	resp, err := testSuite.datastore.GetWithPrefix(testSuite.ctx, key)
	assert.Nil(testSuite.T(), err, "Unexpected error when etcd get returns results")
	assert.Equal(testSuite.T(), len(getResp.Kvs), len(resp), "Expected lengths of resp and getResp to be the same")

	for i := 0; i < len(getResp.Kvs); i++ {
		expectedKey := string(getResp.Kvs[i].Key)
		value, ok := resp[expectedKey]
		if !ok {
			testSuite.T().Errorf("Expected key %v does not exist in resp", expectedKey)
		} else {
			assert.Exactly(testSuite.T(), string(getResp.Kvs[i].Value), value, "Expected value does not match the received response")
		}
	}
}

func (testSuite *DataStoreTestSuite) TestGetEmptyKey() {
	_, err := testSuite.datastore.Get(testSuite.ctx, "")
	assert.Error(testSuite.T(), err, "Expected an error when key is nil")
}

func (testSuite *DataStoreTestSuite) TestGetEtcdGetFails() {
	testSuite.etcdInterface.EXPECT().Get(testSuite.ctx, key).Return(nil, errors.New("Get failed"))

	_, err := testSuite.datastore.Get(testSuite.ctx, key)
	assert.Error(testSuite.T(), err, "Expected an error when etcd get fails")
}

func (testSuite *DataStoreTestSuite) TestGetEtcdGetRespNil() {
	testSuite.etcdInterface.EXPECT().Get(testSuite.ctx, key).Return((*etcd.GetResponse)(nil), nil)

	resp, err := testSuite.datastore.Get(testSuite.ctx, key)
	assert.Nil(testSuite.T(), err, "Unexpected error when etcd get returns empty")
	assert.Empty(testSuite.T(), resp, "Expected an empty map")
}

func (testSuite *DataStoreTestSuite) TestGetEtcdGetRespKVNil() {
	var getResp etcd.GetResponse
	getResp.Kvs = nil
	testSuite.etcdInterface.EXPECT().Get(testSuite.ctx, key).Return(&getResp, nil)

	resp, err := testSuite.datastore.Get(testSuite.ctx, key)
	assert.Nil(testSuite.T(), err, "Unexpected error when etcd get returns empty")
	assert.Empty(testSuite.T(), resp, "Expected an empty map")
}

func (testSuite *DataStoreTestSuite) TestGetEtcdGetRespKVEmpty() {
	var getResp etcd.GetResponse
	getResp.Kvs = []*mvccpb.KeyValue{}
	testSuite.etcdInterface.EXPECT().Get(testSuite.ctx, key).Return(&getResp, nil)

	resp, err := testSuite.datastore.Get(testSuite.ctx, key)
	assert.Nil(testSuite.T(), err, "Unexpected error when etcd get returns empty")
	assert.Empty(testSuite.T(), resp, "Expected an empty map")
}

func (testSuite *DataStoreTestSuite) TestGetEtcd() {
	var getResp etcd.GetResponse
	getResp.Kvs = make([]*mvccpb.KeyValue, 1)
	getResp.Kvs[0] = &mvccpb.KeyValue{
		Key:   []byte(key),
		Value: []byte(value),
	}
	testSuite.etcdInterface.EXPECT().Get(testSuite.ctx, key).Return(&getResp, nil)

	resp, err := testSuite.datastore.Get(testSuite.ctx, key)
	assert.Nil(testSuite.T(), err, "Unexpected error when etcd get returns results")
	assert.Equal(testSuite.T(), len(getResp.Kvs), len(resp), "Expected lengths of resp and getResp to be the same")

	expectedKey := string(getResp.Kvs[0].Key)
	value, ok := resp[expectedKey]
	if !ok {
		testSuite.T().Errorf("Expected key %v does not exist in resp", expectedKey)
	} else {
		assert.Exactly(testSuite.T(), string(getResp.Kvs[0].Value), value, "Expected value does not match the received response")
	}
}
