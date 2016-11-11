// Copyright 2016 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the License). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the license file accompanying this file. This file is distributed
// on an AS IS BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package loader

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/aws/amazon-ecs-event-stream-handler/handler/mocks"
	"github.com/aws/amazon-ecs-event-stream-handler/handler/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var (
	taskClusterARN1           = "arn:aws:ecs:us-east-1:123456789012:cluster/cluster1"
	taskClusterARN2           = "arn:aws:ecs:us-east-1:123456789012:cluster/cluster2"
	taskARN1                  = "arn:aws:ecs:us-east-1:123456789012:task/b6b9eace-958e-4f2a-a09c-8cf43b76cf97"
	taskARN2                  = "arn:aws:ecs:us-east-1:123456789012:task/ab345dfe-6578-2eab-c671-72847ffe8122"
	taskInstanceARN1          = "arn:aws:ecs:us-east-1:123456789012:container-task/b6b9eace-958e-4f2a-a09c-8cf43b76cf97"
	taskDefinitionARN1        = "arn:aws:ecs:us-east-1:123456789012:task-definition/testTask:1"
	redundantClusterARNOfTask = "arn:aws:ecs:us-east-1:123456789012:cluster/red-un-da-nt"
	redundantTaskARN          = "arn:aws:ecs:us-east-1:123456789012:task/red-un-da-nt"
)

type TaskLoaderTestSuite struct {
	suite.Suite
	taskStore      *mocks.MockTaskStore
	ecsWrapper     *mocks.MockECSWrapper
	taskLoader     TaskLoader
	clusterARNList []*string
	task           types.Task
	redundantTask  types.Task
	taskJSON       string
}

func (suite *TaskLoaderTestSuite) SetupTest() {
	mockCtrl := gomock.NewController(suite.T())

	suite.taskStore = mocks.NewMockTaskStore(mockCtrl)
	suite.ecsWrapper = mocks.NewMockECSWrapper(mockCtrl)
	suite.taskLoader = taskLoader{
		taskStore:  suite.taskStore,
		ecsWrapper: suite.ecsWrapper,
	}

	suite.clusterARNList = []*string{&taskClusterARN1, &taskClusterARN2}

	createdAt := "2016-11-07T15:30:00Z"
	startedAt := "2016-11-07T15:45:00Z"
	desiredStatus := "RUNNING"
	lastStatus := "PENDING"
	taskVersion := version
	taskInstanceARN1 := "arn:aws:ecs:us-east-1:123456789012:container-task/b6b9eace-958e-4f2a-a09c-8cf43b76cf97"
	taskDefinitionARN1 := "arn:aws:ecs:us-east-1:123456789012:task-definition/testTask:1"
	suite.task = types.Task{
		Detail: &types.TaskDetail{
			ClusterARN:           &taskClusterARN1,
			ContainerInstanceARN: &taskInstanceARN1,
			Containers:           []*types.Container{},
			CreatedAt:            &createdAt,
			DesiredStatus:        &desiredStatus,
			LastStatus:           &lastStatus,
			Overrides:            &types.Overrides{},
			StartedAt:            startedAt,
			TaskARN:              &taskARN1,
			TaskDefinitionARN:    &taskDefinitionARN1,
			Version:              &taskVersion,
		},
	}

	task, err := json.Marshal(suite.task)
	assert.Nil(suite.T(), err, "Cannot setup testSuite: Unexpected error when marshaling task")
	suite.taskJSON = string(task)

	suite.redundantTask = types.Task{
		Detail: &types.TaskDetail{
			ClusterARN: &redundantClusterARNOfTask,
			TaskARN:    &redundantTaskARN,
		},
	}
}

func TestTaskLoaderTestSuite(t *testing.T) {
	suite.Run(t, new(TaskLoaderTestSuite))
}

func (suite *TaskLoaderTestSuite) TestLoadTasksListAllClustersReturnsError() {
	gomock.InOrder(
		suite.taskStore.EXPECT().ListTasks().Return(make([]types.Task, 0), nil),
		suite.ecsWrapper.EXPECT().ListAllClusters().Return(nil, errors.New("Error while listing all clusters")),
		suite.ecsWrapper.EXPECT().ListAllTasks(gomock.Any()).Times(0),
		suite.ecsWrapper.EXPECT().DescribeTasks(gomock.Any(), gomock.Any()).Times(0),
	)

	err := suite.taskLoader.LoadTasks()
	assert.Error(suite.T(), err, "Expected an error when ecs returns an error when listing all clusters")
}

func (suite *TaskLoaderTestSuite) TestLoadTasksListAllTasksReturnsError() {
	gomock.InOrder(
		suite.taskStore.EXPECT().ListTasks().Return(make([]types.Task, 0), nil),
		suite.ecsWrapper.EXPECT().ListAllClusters().Return(suite.clusterARNList, nil),
		suite.ecsWrapper.EXPECT().ListAllTasks(suite.clusterARNList[0]).Return(nil, errors.New("Error while listing all tasks")),
		suite.ecsWrapper.EXPECT().ListAllTasks(suite.clusterARNList[1]).Times(0),
		suite.ecsWrapper.EXPECT().DescribeTasks(gomock.Any(), gomock.Any()).Times(0),
	)

	err := suite.taskLoader.LoadTasks()
	assert.Error(suite.T(), err, "Expected an error when ecs returns an error when listing all tasks in a cluster")
}

func (suite *TaskLoaderTestSuite) TestLoadTasksDescribeTasksReturnsError() {
	taskARNList := []*string{&taskARN1}
	gomock.InOrder(
		suite.taskStore.EXPECT().ListTasks().Return(make([]types.Task, 0), nil),
		suite.ecsWrapper.EXPECT().ListAllClusters().Return(suite.clusterARNList, nil),
		suite.ecsWrapper.EXPECT().ListAllTasks(suite.clusterARNList[0]).Return(taskARNList, nil),
		suite.ecsWrapper.EXPECT().ListAllTasks(suite.clusterARNList[1]).Times(0),
		suite.ecsWrapper.EXPECT().DescribeTasks(suite.clusterARNList[0], taskARNList).Return(nil, nil, errors.New("Error while desribing task")),
	)

	err := suite.taskLoader.LoadTasks()
	assert.Error(suite.T(), err, "Expected an error when ecs returns an error when describing tasks")
}

func (suite *TaskLoaderTestSuite) TestLoadTasksStoreReturnsError() {
	taskARNList := []*string{&taskARN1}
	taskList := []types.Task{suite.task}
	emptyTaskARNList := []*string{}
	gomock.InOrder(
		suite.taskStore.EXPECT().ListTasks().Return(make([]types.Task, 0), nil),
		suite.ecsWrapper.EXPECT().ListAllClusters().Return(suite.clusterARNList, nil),
		suite.ecsWrapper.EXPECT().ListAllTasks(suite.clusterARNList[0]).Return(taskARNList, nil),
		suite.ecsWrapper.EXPECT().DescribeTasks(suite.clusterARNList[0], taskARNList).Return(taskList, nil, nil),
		suite.ecsWrapper.EXPECT().ListAllTasks(suite.clusterARNList[1]).Return(emptyTaskARNList, nil),
		suite.ecsWrapper.EXPECT().DescribeTasks(suite.clusterARNList[1], gomock.Any()).Times(0),
		suite.taskStore.EXPECT().AddUnversionedTask(suite.taskJSON).Return(errors.New("Error while adding task to store")),
	)

	err := suite.taskLoader.LoadTasks()
	assert.Error(suite.T(), err, "Expected an error when store returns an error when adding task")
}

func (suite *TaskLoaderTestSuite) TestLoadTasksEmptyLocalStore() {
	taskARNList := []*string{&taskARN1}
	taskList := []types.Task{suite.task}
	emptyTaskARNList := []*string{}
	gomock.InOrder(
		suite.taskStore.EXPECT().ListTasks().Return(make([]types.Task, 0), nil),
		suite.ecsWrapper.EXPECT().ListAllClusters().Return(suite.clusterARNList, nil),
		suite.ecsWrapper.EXPECT().ListAllTasks(suite.clusterARNList[0]).Return(taskARNList, nil),
		suite.ecsWrapper.EXPECT().DescribeTasks(suite.clusterARNList[0], taskARNList).Return(taskList, nil, nil),
		suite.ecsWrapper.EXPECT().ListAllTasks(suite.clusterARNList[1]).Return(emptyTaskARNList, nil),
		suite.ecsWrapper.EXPECT().DescribeTasks(suite.clusterARNList[1], gomock.Any()).Times(0),
		suite.taskStore.EXPECT().AddUnversionedTask(suite.taskJSON).Return(nil),
	)
	err := suite.taskLoader.LoadTasks()
	assert.Nil(suite.T(), err, "Unexpected error when loading tasks")
}

func (suite *TaskLoaderTestSuite) TestLoadTasksLocalStoreSameAsECS() {
	taskARNList := []*string{&taskARN1}
	emptyTaskARNList := []*string{}
	taskListInStore := []types.Task{suite.task}
	taskList := []types.Task{suite.task}
	// taskListInStore == taskList, which should mean that there shouldn't
	// be a call to DeleteTask()
	suite.taskStore.EXPECT().DeleteTask(gomock.Any(), gomock.Any()).Return(nil).Times(0)
	gomock.InOrder(
		suite.taskStore.EXPECT().ListTasks().Return(taskListInStore, nil),
		suite.ecsWrapper.EXPECT().ListAllClusters().Return(suite.clusterARNList, nil),
		suite.ecsWrapper.EXPECT().ListAllTasks(suite.clusterARNList[0]).Return(taskARNList, nil),
		suite.ecsWrapper.EXPECT().DescribeTasks(suite.clusterARNList[0], taskARNList).Return(taskList, nil, nil),
		suite.ecsWrapper.EXPECT().ListAllTasks(suite.clusterARNList[1]).Return(emptyTaskARNList, nil),
		suite.ecsWrapper.EXPECT().DescribeTasks(suite.clusterARNList[1], gomock.Any()).Times(0),
		suite.taskStore.EXPECT().AddUnversionedTask(suite.taskJSON).Return(nil),
	)
	err := suite.taskLoader.LoadTasks()
	assert.Nil(suite.T(), err, "Unexpected error when loading tasks")
}

func (suite *TaskLoaderTestSuite) TestLoadTasksRedundantEntriesInLocalStore() {
	taskARNList := []*string{&taskARN1}
	emptyTaskARNList := []*string{}
	taskListInStore := []types.Task{suite.task, suite.redundantTask}
	taskList := []types.Task{suite.task}
	gomock.InOrder(
		suite.taskStore.EXPECT().ListTasks().Return(taskListInStore, nil),
		suite.ecsWrapper.EXPECT().ListAllClusters().Return(suite.clusterARNList, nil),
		suite.ecsWrapper.EXPECT().ListAllTasks(suite.clusterARNList[0]).Return(taskARNList, nil),
		suite.ecsWrapper.EXPECT().DescribeTasks(suite.clusterARNList[0], taskARNList).Return(taskList, nil, nil),
		suite.ecsWrapper.EXPECT().ListAllTasks(suite.clusterARNList[1]).Return(emptyTaskARNList, nil),
		suite.ecsWrapper.EXPECT().DescribeTasks(suite.clusterARNList[1], gomock.Any()).Times(0),
		suite.taskStore.EXPECT().AddUnversionedTask(suite.taskJSON).Return(nil),
		// Expect delete task for the redundant task
		suite.taskStore.EXPECT().DeleteTask(redundantClusterARNOfTask, redundantTaskARN).Return(nil),
	)
	err := suite.taskLoader.LoadTasks()
	assert.Nil(suite.T(), err, "Unexpected error when loading tasks")
}
