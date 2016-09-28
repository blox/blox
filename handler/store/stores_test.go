package store

import (
	"github.com/aws/amazon-ecs-event-stream-handler/handler/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type StoreTestSuite struct {
	suite.Suite
	datastore *mocks.MockDataStore
}

func (testSuite *StoreTestSuite) SetupTest() {
	mockCtrl := gomock.NewController(testSuite.T())
	testSuite.datastore = mocks.NewMockDataStore(mockCtrl)
}

func TestStoreTestSuite(t *testing.T) {
	suite.Run(t, new(StoreTestSuite))
}

func (testSuite *StoreTestSuite) TestNewStoresDatastoreNil() {
	_, err := NewStores(nil)
	assert.Error(testSuite.T(), err, "Expected an error when NewStores is initialized with nil")
}

func (testSuite *StoreTestSuite) TestNewStores() {
	stores, err := NewStores(testSuite.datastore)
	assert.Nil(testSuite.T(), err, "Unexpected error when calling NewStores")
	assert.NotNil(testSuite.T(), stores, "Stores should not be nil")
	assert.NotNil(testSuite.T(), stores.TaskStore, "TaskStore should not be nil")
	assert.NotNil(testSuite.T(), stores.ContainerInstanceStore, "ContainerInstanceStores should not be nil")
}
