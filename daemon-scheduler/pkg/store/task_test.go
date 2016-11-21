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
	"errors"
	"testing"

	"github.com/blox/blox/daemon-scheduler/pkg/json"
	"github.com/blox/blox/daemon-scheduler/pkg/mocks"
	"github.com/blox/blox/daemon-scheduler/pkg/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	env      = "testEnv"
	task     = "arn:aws:ecs:us-east-1:12345678912:task/c024d145-093b-499a-9b14-5baf273f5835"
	instance = "arn:aws:us-east-1:123456789123:container-instance/4b6d45ea-a4b4-4269-9d04-3af6ddfdc597"
)

type TaskTestSuite struct {
	suite.Suite
	datastore *mocks.MockDataStore
	taskStore TaskStore
	ctx       context.Context
	task      *types.Task
	taskJSON  string
	taskKey   string
}

func (suite *TaskTestSuite) SetupTest() {
	mockCtrl := gomock.NewController(suite.T())
	var err error

	suite.datastore = mocks.NewMockDataStore(mockCtrl)
	suite.ctx = context.TODO()

	suite.task, err = types.NewTask(task, instance)
	assert.Nil(suite.T(), err, "Cannot initialize TaskTestSuite")
	suite.taskJSON, err = json.MarshalJSON(suite.task)
	assert.Nil(suite.T(), err, "Cannot initialize TaskTestSuite")

	suite.taskStore, err = NewTaskStore(suite.datastore)
	assert.Nil(suite.T(), err, "Cannot initialize TaskTestSuite")

	suite.taskKey = environmentKeyPrefix + env + taskKeyConnector + task
}

func TestTaskTestSuite(t *testing.T) {
	suite.Run(t, new(TaskTestSuite))
}

func (suite *TaskTestSuite) TestNewTaskStoreEmptyDataStore() {
	_, err := NewTaskStore(nil)
	assert.Error(suite.T(), err, "Expected an error when datastore is nil")
}

func (suite *TaskTestSuite) TestNewNewTaskStore() {
	ts, err := NewTaskStore(suite.datastore)
	assert.Nil(suite.T(), err, "Unexpected error when datastore is not nil")
	assert.NotNil(suite.T(), ts, "Environment store should not be nil")
}

func (suite *TaskTestSuite) TestPutTaskEmptyEnvName() {
	err := suite.taskStore.PutTask(suite.ctx, "", *suite.task)
	assert.NotNil(suite.T(), err, "Expected an error putting a task when environment name is empty")
}

func (suite *TaskTestSuite) TestPutTaskEmptyTaskARN() {
	task := suite.task
	task.TaskARN = ""
	err := suite.taskStore.PutTask(suite.ctx, envName, *task)
	assert.NotNil(suite.T(), err, "Expected an error putting a task when task ARN is empty")
}

func (suite *TaskTestSuite) TestPutTaskDataStoreReturnsError() {
	suite.datastore.EXPECT().Put(suite.ctx, suite.taskKey, suite.taskJSON).Return(errors.New("Error adding task"))
	err := suite.taskStore.PutTask(suite.ctx, envName, *suite.task)
	assert.NotNil(suite.T(), err, "Expected an error putting a task when datastore returns an error")
}

func (suite *TaskTestSuite) TestPutTask() {
	suite.datastore.EXPECT().Put(suite.ctx, suite.taskKey, suite.taskJSON).Return(nil)
	err := suite.taskStore.PutTask(suite.ctx, envName, *suite.task)
	assert.Nil(suite.T(), err, "Unexpected error putting a task")
}
