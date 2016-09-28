package store

import (
	"context"
	"testing"
	"time"

	"github.com/aws/amazon-ecs-event-stream-handler/handler/mocks"
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
}

func (testSuite *DataStoreTestSuite) SetupTest() {
	mockCtrl := gomock.NewController(testSuite.T())
	testSuite.etcdInterface = mocks.NewMockEtcdInterface(mockCtrl)

	var err error
	testSuite.datastore, err = NewDataStore(testSuite.etcdInterface)
	assert.Nil(testSuite.T(), err, "Cannot initialize DataStoreTestSuite")
}

func TestDataStoreTestSuite(t *testing.T) {
	suite.Run(t, new(DataStoreTestSuite))
}

func (testSuite *DataStoreTestSuite) TestNewDataStoreEmptyEtcd() {
	_, err := NewDataStore(nil)
	assert.Error(testSuite.T(), err, "Expected an error when etcd client is nil")
}

func (testSuite *DataStoreTestSuite) TestNewDataStore() {
	_, err := NewDataStore(nil)
	assert.Error(testSuite.T(), err, "Expected an error when etcd client is nil")
}

func (testSuite *DataStoreTestSuite) TestAddEmptyKey() {
	err := testSuite.datastore.Add("", "test")
	assert.Error(testSuite.T(), err, "Expected an error when key is nil")
}

func (testSuite *DataStoreTestSuite) TestAddEmptyValue() {
	err := testSuite.datastore.Add("test", "")
	assert.Error(testSuite.T(), err, "Expected an error when value is nil")
}

func (testSuite *DataStoreTestSuite) TestAddEtcdPutFails() {
	testSuite.etcdInterface.EXPECT().Put(gomock.Any(), key, value).Return(nil, errors.New("Put failed"))

	err := testSuite.datastore.Add(key, value)
	assert.Error(testSuite.T(), err, "Expected an error when etcd put fails")
}

func (testSuite *DataStoreTestSuite) TestAdd() {
	testSuite.etcdInterface.EXPECT().Put(gomock.Any(), key, value).Return(nil, nil)

	err := testSuite.datastore.Add(key, value)
	assert.Nil(testSuite.T(), err, "Unexpected error when calling adding key %v, value %v", key, value)
}

func (testSuite *DataStoreTestSuite) TestGetWithPrefixEmptyKey() {
	_, err := testSuite.datastore.GetWithPrefix("")
	assert.Error(testSuite.T(), err, "Expected an error when key is nil")
}

func (testSuite *DataStoreTestSuite) TestGetWithPrefixEtcdGetFails() {
	testSuite.etcdInterface.EXPECT().Get(gomock.Any(), key, gomock.Any()).Return(nil, errors.New("Get failed"))

	_, err := testSuite.datastore.GetWithPrefix(key)
	assert.Error(testSuite.T(), err, "Expected an error when etcd get fails")
}

func (testSuite *DataStoreTestSuite) TestGetWithPrefixEtcdGetRespNil() {
	testSuite.etcdInterface.EXPECT().Get(gomock.Any(), key, gomock.Any()).Return((*etcd.GetResponse)(nil), nil)

	resp, err := testSuite.datastore.GetWithPrefix(key)
	assert.Nil(testSuite.T(), err, "Unexpected error when etcd get returns empty")
	assert.Empty(testSuite.T(), resp, "Expected an empty map")
}

func (testSuite *DataStoreTestSuite) TestGetWithPrefixEtcdGetRespKVNil() {
	var getResp etcd.GetResponse
	getResp.Kvs = nil
	testSuite.etcdInterface.EXPECT().Get(gomock.Any(), key, gomock.Any()).Return(&getResp, nil)

	resp, err := testSuite.datastore.GetWithPrefix(key)
	assert.Nil(testSuite.T(), err, "Unexpected error when etcd get returns empty")
	assert.Empty(testSuite.T(), resp, "Expected an empty map")
}

func (testSuite *DataStoreTestSuite) TestGetWithPrefixEtcdGetRespKVEmpty() {
	var getResp etcd.GetResponse
	getResp.Kvs = []*mvccpb.KeyValue{}
	testSuite.etcdInterface.EXPECT().Get(gomock.Any(), key, gomock.Any()).Return(&getResp, nil)

	resp, err := testSuite.datastore.GetWithPrefix(key)
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
	testSuite.etcdInterface.EXPECT().Get(gomock.Any(), key, gomock.Any()).Return(&getResp, nil)

	resp, err := testSuite.datastore.GetWithPrefix(key)
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
	_, err := testSuite.datastore.Get("")
	assert.Error(testSuite.T(), err, "Expected an error when key is nil")
}

func (testSuite *DataStoreTestSuite) TestGetEtcdGetFails() {
	testSuite.etcdInterface.EXPECT().Get(gomock.Any(), key).Return(nil, errors.New("Get failed"))

	_, err := testSuite.datastore.Get(key)
	assert.Error(testSuite.T(), err, "Expected an error when etcd get fails")
}

func (testSuite *DataStoreTestSuite) TestGetEtcdGetRespNil() {
	testSuite.etcdInterface.EXPECT().Get(gomock.Any(), key).Return((*etcd.GetResponse)(nil), nil)

	resp, err := testSuite.datastore.Get(key)
	assert.Nil(testSuite.T(), err, "Unexpected error when etcd get returns empty")
	assert.Empty(testSuite.T(), resp, "Expected an empty map")
}

func (testSuite *DataStoreTestSuite) TestGetEtcdGetRespKVNil() {
	var getResp etcd.GetResponse
	getResp.Kvs = nil
	testSuite.etcdInterface.EXPECT().Get(gomock.Any(), key).Return(&getResp, nil)

	resp, err := testSuite.datastore.Get(key)
	assert.Nil(testSuite.T(), err, "Unexpected error when etcd get returns empty")
	assert.Empty(testSuite.T(), resp, "Expected an empty map")
}

func (testSuite *DataStoreTestSuite) TestGetEtcdGetRespKVEmpty() {
	var getResp etcd.GetResponse
	getResp.Kvs = []*mvccpb.KeyValue{}
	testSuite.etcdInterface.EXPECT().Get(gomock.Any(), key).Return(&getResp, nil)

	resp, err := testSuite.datastore.Get(key)
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
	testSuite.etcdInterface.EXPECT().Get(gomock.Any(), key).Return(&getResp, nil)

	resp, err := testSuite.datastore.Get(key)
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

func (testSuite *DataStoreTestSuite) TestStreamWithPrefixEmptyKeyPrefix() {
	ctx := context.Background()
	_, err := testSuite.datastore.StreamWithPrefix(ctx, "")
	assert.Error(testSuite.T(), err, "Expected an error when key is nil")
}

func (testSuite *DataStoreTestSuite) TestStreamWithPrefix() {
	ctx := context.Background()
	watchChan := make(chan etcd.WatchResponse)
	defer close(watchChan)
	testSuite.etcdInterface.EXPECT().Watch(gomock.Any(), key, gomock.Any()).Return(watchChan)

	dsChan, err := testSuite.datastore.StreamWithPrefix(ctx, key)
	assert.Nil(testSuite.T(), err, "Unexpected error when setting up streaming")
	assert.NotNil(testSuite.T(), dsChan, "Expected valid channel for streaming")

	dsVal := addToWatchChanAndReadFromDataChan(watchChan, dsChan)
	expectedDsVal := map[string]string{
		key: value,
	}
	assert.Equal(testSuite.T(), expectedDsVal, dsVal, "Expected key-val read from dsChan to match what was put into watchChan")
}

func (testSuite *DataStoreTestSuite) TestStreamWithPrefixCancelUpstreamContext() {
	ctx, cancel := context.WithCancel(context.Background())
	var watchChan etcd.WatchChan
	testSuite.etcdInterface.EXPECT().Watch(gomock.Any(), key, gomock.Any()).Return(watchChan)

	dsChan, err := testSuite.datastore.StreamWithPrefix(ctx, key)
	assert.Nil(testSuite.T(), err, "Unexpected error when setting up streaming")
	assert.NotNil(testSuite.T(), dsChan, "Expected valid channel for streaming")

	cancel()

	_, ok := <-dsChan
	assert.False(testSuite.T(), ok, "Expected dschan to be closed")
}

func (testSuite *DataStoreTestSuite) TestStreamWithPrefixCloseDownstreamChannel() {
	ctx := context.Background()
	watchChan := make(chan etcd.WatchResponse)
	testSuite.etcdInterface.EXPECT().Watch(gomock.Any(), key, gomock.Any()).Return(watchChan)

	dsChan, err := testSuite.datastore.StreamWithPrefix(ctx, key)
	assert.Nil(testSuite.T(), err, "Unexpected error when setting up streaming")
	assert.NotNil(testSuite.T(), dsChan, "Expected valid channel for streaming")

	close(watchChan)

	_, ok := <-dsChan
	assert.False(testSuite.T(), ok, "Expected dschan to be closed")
}

func (testSuite *DataStoreTestSuite) TestStreamWithPrefixStreamTimeout() {
	if testing.Short() {
		testSuite.T().Skip("Skipping TestStreamWithPrefixStreamTimeout in short mode")
	}

	ctx := context.Background()
	var watchChan etcd.WatchChan
	testSuite.etcdInterface.EXPECT().Watch(gomock.Any(), key, gomock.Any()).Return(watchChan)

	dsChan, err := testSuite.datastore.StreamWithPrefix(ctx, key)
	assert.Nil(testSuite.T(), err, "Unexpected error when setting up streaming")
	assert.NotNil(testSuite.T(), dsChan, "Expected valid channel for streaming")

	time.Sleep(streamIdleTimeout)

	_, ok := <-dsChan
	assert.False(testSuite.T(), ok, "Expected dschan to be closed")
}

func addToWatchChanAndReadFromDataChan(watchChan chan etcd.WatchResponse, dsChan chan map[string]string) map[string]string {
	var dsVal map[string]string

	doneChan := make(chan bool)
	defer close(doneChan)
	go func() {
		dsVal = <-dsChan
		doneChan <- true
	}()

	var event etcd.Event
	event.Kv = &mvccpb.KeyValue{
		Key:   []byte(key),
		Value: []byte(value),
	}
	var watchResp etcd.WatchResponse
	watchResp.Events = make([]*etcd.Event, 1)
	watchResp.Events[0] = &event

	watchChan <- watchResp
	<-doneChan

	return dsVal
}
