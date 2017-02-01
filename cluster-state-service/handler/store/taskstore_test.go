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
	"encoding/json"
	"errors"
	"testing"

	"github.com/blox/blox/cluster-state-service/handler/mocks"
	storetypes "github.com/blox/blox/cluster-state-service/handler/store/types"
	"github.com/blox/blox/cluster-state-service/handler/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var (
	taskARN1      = "arn:aws:ecs:us-east-1:123456789012:task/271022c0-f894-4aa2-b063-25bae55088d5"
	taskARN2      = "arn:aws:ecs:us-east-1:123456789012:task/345022c0-f894-4aa2-b063-25bae55088d5"
	taskARN3      = "arn:aws:ecs:us-east-1:123456789012:task/345022c0-f894-4aa2-b063-25bae55088dd"
	pendingStatus = "pending"
	runningStatus = "running"
	someoneElse   = "someone-else"
)

type TaskStoreTestSuite struct {
	suite.Suite
	datastore                          *mocks.MockDataStore
	etcdTxStore                        *mocks.MockEtcdTXStore
	taskStore                          TaskStore
	taskKey1                           string
	firstPendingTask                   types.Task
	firstPendingTaskJSON               string
	secondPendingTask                  types.Task
	secondPendingTaskJSON              string
	firstTaskStartedBySomeoneElse      types.Task
	firstTaskStartedBySomeoneElseJSON  string
	secondTaskStartedBySomeoneElse     types.Task
	secondTaskStartedBySomeoneElseJSON string
	firstTaskOfFirstCluster            types.Task
	firstTaskOfFirstClusterJSON        string
	secondTaskOfFirstCluster           types.Task
	secondTaskOfFirstClusterJSON       string
}

func (suite *TaskStoreTestSuite) SetupTest() {
	mockCtrl := gomock.NewController(suite.T())
	suite.datastore = mocks.NewMockDataStore(mockCtrl)
	suite.etcdTxStore = mocks.NewMockEtcdTXStore(mockCtrl)

	suite.taskKey1 = taskKeyPrefix + clusterName1 + "/" + taskARN1

	var err error
	suite.taskStore, err = NewTaskStore(suite.datastore, suite.etcdTxStore)
	assert.Nil(suite.T(), err, "Cannot setup testSuite: Unexpected error when calling NewTaskStore")

	version1 := int64(1)
	suite.firstPendingTask = types.Task{
		Detail: &types.TaskDetail{
			TaskARN:    &taskARN1,
			ClusterARN: &clusterARN1,
			LastStatus: &pendingStatus,
			Version:    &version1,
		},
	}
	suite.firstPendingTaskJSON = suite.setupTask(suite.firstPendingTask)
	suite.firstTaskStartedBySomeoneElse = types.Task{
		Detail: &types.TaskDetail{
			TaskARN:    &taskARN1,
			ClusterARN: &clusterARN1,
			LastStatus: &pendingStatus,
			Version:    &version1,
			StartedBy:  someoneElse,
		},
	}
	suite.firstTaskStartedBySomeoneElseJSON = suite.setupTask(suite.firstTaskStartedBySomeoneElse)

	version2 := version1 + 1
	suite.secondPendingTask = types.Task{
		Detail: &types.TaskDetail{
			TaskARN:    &taskARN2,
			ClusterARN: &clusterARN2,
			LastStatus: &pendingStatus,
			Version:    &version2,
		},
	}
	suite.secondPendingTaskJSON = suite.setupTask(suite.secondPendingTask)
	suite.secondTaskStartedBySomeoneElse = types.Task{
		Detail: &types.TaskDetail{
			TaskARN:    &taskARN2,
			ClusterARN: &clusterARN2,
			LastStatus: &pendingStatus,
			Version:    &version2,
			StartedBy:  someoneElse,
		},
	}
	suite.secondTaskStartedBySomeoneElseJSON = suite.setupTask(suite.secondTaskStartedBySomeoneElse)
	suite.firstTaskOfFirstCluster = types.Task{
		Detail: &types.TaskDetail{
			TaskARN:    &taskARN1,
			ClusterARN: &clusterARN1,
			LastStatus: &pendingStatus,
			Version:    &version1,
		},
	}
	suite.firstTaskOfFirstClusterJSON = suite.setupTask(suite.firstTaskOfFirstCluster)

	suite.secondTaskOfFirstCluster = types.Task{
		Detail: &types.TaskDetail{
			TaskARN:    &taskARN2,
			ClusterARN: &clusterARN1,
			LastStatus: &pendingStatus,
			Version:    &version2,
		},
	}
	suite.secondTaskOfFirstClusterJSON = suite.setupTask(suite.secondTaskOfFirstCluster)
}

func (suite *TaskStoreTestSuite) setupTask(task types.Task) string {
	taskJSON, err := json.Marshal(task)
	assert.Nil(suite.T(), err, "Cannot setup testSuite: Error when json marhsaling task %v", task)
	taskJSONString := string(taskJSON)
	return taskJSONString
}

func TestTaskStoreTestSuite(t *testing.T) {
	suite.Run(t, new(TaskStoreTestSuite))
}

func (suite *TaskStoreTestSuite) TestNewTaskStoreNilDatastore() {
	_, err := NewTaskStore(nil, suite.etcdTxStore)
	assert.Error(suite.T(), err, "Expected an error when datastore is nil")
}

func (suite *TaskStoreTestSuite) TestNewTaskStoreNilEtcdTXStore() {
	_, err := NewTaskStore(suite.datastore, nil)
	assert.Error(suite.T(), err, "Expected an error when etcd transactional store is nil")
}

func (suite *TaskStoreTestSuite) TestNewTaskStore() {
	taskStore, err := NewTaskStore(suite.datastore, suite.etcdTxStore)
	assert.Nil(suite.T(), err, "Unexpected error when calling NewTaskStore")
	assert.NotNil(suite.T(), taskStore, "TaskStore should not be nil")
}

func (suite *TaskStoreTestSuite) TestAddTaskEmptyJSON() {
	err := suite.taskStore.AddTask("")
	assert.Error(suite.T(), err, "Expected an error when json empty in AddTask")
}

func (suite *TaskStoreTestSuite) TestAddTaskUnmarshalJSONError() {
	err := suite.taskStore.AddTask("invalidJSON")
	assert.Error(suite.T(), err, "Expected an error when json invalid in AddTask")
}

func (suite *TaskStoreTestSuite) TestAddTaskEmptyTask() {
	task := types.Task{}
	task1JSON, err := json.Marshal(task)
	err = suite.taskStore.AddTask(string(task1JSON))
	assert.Error(suite.T(), err, "Expected an error when task arn is not set")
}

func (suite *TaskStoreTestSuite) TestAddTaskTaskDetailNotSet() {
	task := types.Task{}
	taskJSON, err := json.Marshal(task)
	err = suite.taskStore.AddTask(string(taskJSON))
	assert.Error(suite.T(), err, "Expected an error when task detail is not set")
}

func (suite *TaskStoreTestSuite) TestAddTaskTaskARNNotSet() {
	taskDetail := types.TaskDetail{
		ClusterARN: &clusterARN1,
	}
	task := types.Task{
		Detail: &taskDetail,
	}
	taskJSON, err := json.Marshal(task)
	err = suite.taskStore.AddTask(string(taskJSON))
	assert.Error(suite.T(), err, "Expected an error when task ARN is not set")
}

func (suite *TaskStoreTestSuite) TestAddTaskClusterARNNotSet() {
	taskDetail := types.TaskDetail{
		TaskARN: &taskARN1,
	}
	task := types.Task{
		Detail: &taskDetail,
	}
	taskJSON, err := json.Marshal(task)
	err = suite.taskStore.AddTask(string(taskJSON))
	assert.Error(suite.T(), err, "Expected an error when cluster ARN is not set")
}

func (suite *TaskStoreTestSuite) TestAddTaskEmptyTaskARN() {
	taskARN := ""
	taskDetail := types.TaskDetail{
		ClusterARN: &clusterARN1,
		TaskARN:    &taskARN,
	}
	task := types.Task{
		Detail: &taskDetail,
	}
	taskJSON, err := json.Marshal(task)
	err = suite.taskStore.AddTask(string(taskJSON))
	assert.Error(suite.T(), err, "Expected an error when task ARN is empty")
}

func (suite *TaskStoreTestSuite) TestAddTaskEmptyClusterARN() {
	clusterARN := ""
	taskDetail := types.TaskDetail{
		ClusterARN: &clusterARN,
		TaskARN:    &taskARN1,
	}
	task := types.Task{
		Detail: &taskDetail,
	}
	taskJSON, err := json.Marshal(task)
	err = suite.taskStore.AddTask(string(taskJSON))
	assert.Error(suite.T(), err, "Expected an error when cluster ARN is empty")
}

func (suite *TaskStoreTestSuite) TestAddTaskInvalidClusterARNWithNoName() {
	clusterARN := "arn:aws:ecs:us-east-1:123456789123:cluster/"
	taskDetail := types.TaskDetail{
		ClusterARN: &clusterARN,
		TaskARN:    &taskARN1,
	}
	task := types.Task{
		Detail: &taskDetail,
	}
	taskJSON, err := json.Marshal(task)
	err = suite.taskStore.AddTask(string(taskJSON))
	assert.Error(suite.T(), err, "Expected an error when cluster ARN has no name")
}

func (suite *TaskStoreTestSuite) TestAddTaskSTMRepeatableFails() {
	suite.etcdTxStore.EXPECT().GetV3Client().Return(nil)
	suite.etcdTxStore.EXPECT().NewSTMRepeatable(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("Error when getting key"))
	err := suite.taskStore.AddTask(suite.firstPendingTaskJSON)
	assert.Error(suite.T(), err, "Expected error when STM repeatable fails to execute with an error")
}

func (suite *TaskStoreTestSuite) TestAddUnversionedTaskEmptyVersion() {
	suite.datastore.EXPECT().Get(gomock.Any()).Times(0)
	suite.datastore.EXPECT().Add(gomock.Any(), gomock.Any()).Times(0)

	task := suite.firstPendingTask
	task.Detail.Version = nil

	taskJSON, err := json.Marshal(task)
	assert.Nil(suite.T(), err, "Error when json marhsaling task %v", task)

	err = suite.taskStore.AddUnversionedTask(string(taskJSON))
	assert.NotNil(suite.T(), err, "Expected an error when adding unversioned task with empty version")
}

func (suite *TaskStoreTestSuite) TestAddUnversionedTaskInvalidVersion() {
	suite.datastore.EXPECT().Get(gomock.Any()).Times(0)
	suite.datastore.EXPECT().Add(gomock.Any(), gomock.Any()).Times(0)

	err := suite.taskStore.AddUnversionedTask(suite.firstPendingTaskJSON)
	assert.NotNil(suite.T(), err, "Expected ab error when adding unversioned task with invalid version")
}

func (suite *TaskStoreTestSuite) TestGetTaskEmptyClusterName() {
	_, err := suite.taskStore.GetTask("", taskARN1)
	assert.Error(suite.T(), err, "Expected an error when cluster name is empty in GetTask")
}

func (suite *TaskStoreTestSuite) TestGetTaskEmptyTaskARN() {
	_, err := suite.taskStore.GetTask(clusterName1, "")
	assert.Error(suite.T(), err, "Expected an error when task ARN is empty in GetTask")
}

func (suite *TaskStoreTestSuite) TestGetTaskGetTaskFails() {
	suite.datastore.EXPECT().Get(suite.taskKey1).Return(nil, errors.New("Error when getting key"))

	_, err := suite.taskStore.GetTask(clusterName1, taskARN1)
	assert.Error(suite.T(), err, "Expected an error when get task fails")
}

func (suite *TaskStoreTestSuite) TestGetTaskGetTaskNoResults() {
	suite.datastore.EXPECT().Get(suite.taskKey1).Return(make(map[string]string), nil)

	task, err := suite.taskStore.GetTask(clusterName1, taskARN1)
	assert.Nil(suite.T(), err, "Unexpected error when datastore returns empty results")
	assert.Nil(suite.T(), task, "Unexpected object returned when datastore returns empty results")
}

func (suite *TaskStoreTestSuite) TestGetTaskGetMultipleResults() {
	resp := map[string]string{
		taskARN1: suite.firstPendingTaskJSON,
		taskARN2: suite.secondPendingTaskJSON,
	}

	suite.datastore.EXPECT().Get(suite.taskKey1).Return(resp, nil)

	_, err := suite.taskStore.GetTask(clusterName1, taskARN1)
	assert.Error(suite.T(), err, "Expected an error when datastore returns multiple results")
}

func (suite *TaskStoreTestSuite) TestGetTaskGetInvalidJSONResult() {
	resp := map[string]string{
		taskARN1: "invalidJSON",
	}
	suite.datastore.EXPECT().Get(suite.taskKey1).Return(resp, nil)
	_, err := suite.taskStore.GetTask(clusterName1, taskARN1)
	assert.Error(suite.T(), err, "Expected an error when datastore returns invalid json results")
}

func (suite *TaskStoreTestSuite) TestGetTaskWithClusterNameAndTaskARN() {
	resp := map[string]string{
		taskARN1: suite.firstPendingTaskJSON,
	}
	suite.datastore.EXPECT().Get(suite.taskKey1).Return(resp, nil)
	task, err := suite.taskStore.GetTask(clusterName1, taskARN1)
	assert.Nil(suite.T(), err, "Unexpected error when getting task")
	assert.NotNil(suite.T(), task, "Expected a non-nil task when calling GetTask")
	assert.Exactly(suite.T(), suite.firstPendingTask, *task, "Expected the returned task to match the one returned from the datastore")
}

func (suite *TaskStoreTestSuite) TestGetTaskWithClusterARNAndTaskARN() {
	resp := map[string]string{
		taskARN1: suite.firstPendingTaskJSON,
	}

	suite.datastore.EXPECT().Get(suite.taskKey1).Return(resp, nil)

	task, err := suite.taskStore.GetTask(clusterARN1, taskARN1)
	assert.Nil(suite.T(), err, "Unexpected error when getting task")
	assert.NotNil(suite.T(), task, "Expected a non-nil task when calling GetTask")

	assert.Exactly(suite.T(), suite.firstPendingTask, *task, "Expected the returned task to match the one returned from the datastore")
}

func (suite *TaskStoreTestSuite) TestListTasksGetWithPrefixInvalidJSON() {
	resp := map[string]string{
		taskARN1: "invalidJSON",
	}
	suite.datastore.EXPECT().GetWithPrefix(taskKeyPrefix).Return(resp, nil)

	_, err := suite.taskStore.ListTasks()
	assert.Error(suite.T(), err, "Expected an error when GetWithPrefix fails")
}

func (suite *TaskStoreTestSuite) TestListTasksGetWithPrefixFails() {
	suite.datastore.EXPECT().GetWithPrefix(taskKeyPrefix).Return(nil, errors.New("GetWithPrefix failed"))

	_, err := suite.taskStore.ListTasks()
	assert.Error(suite.T(), err, "Expected an error when GetWithPrefix fails")
}

func (suite *TaskStoreTestSuite) TestListTasksGetWithPrefixReturnsNoResults() {
	suite.datastore.EXPECT().GetWithPrefix(taskKeyPrefix).Return(make(map[string]string), nil)

	tasks, err := suite.taskStore.ListTasks()
	assert.Nil(suite.T(), err, "Unexpected error when GetWithPrefix returns empty results")
	assert.NotNil(suite.T(), tasks, "Expected a non-nil result when GetWithPrefix returns empty results")

	assert.Empty(suite.T(), tasks, "Tasks should be empty when GetWithPrefix returns empty results")
}

func (suite *TaskStoreTestSuite) TestListTasksGetWithPrefixReturnsMultipleResults() {
	resp := map[string]string{
		taskARN1: suite.firstPendingTaskJSON,
		taskARN2: suite.secondPendingTaskJSON,
	}

	suite.datastore.EXPECT().GetWithPrefix(taskKeyPrefix).Return(resp, nil)

	tasks, err := suite.taskStore.ListTasks()
	assert.Nil(suite.T(), err, "Unexpected error when GetWithPrefix returns multiple results")
	assert.NotNil(suite.T(), tasks, "Expected a non-nil result when GetWithPrefix returns multiple results")

	assert.Equal(suite.T(), len(resp), len(tasks), "Expected ListTasks result to be of the same length as GetWithPrefix result")

	for _, v := range tasks {
		//attempt to grab the same task from resp
		value, ok := resp[*v.Detail.TaskARN]
		if !ok {
			suite.T().Errorf("Expected GetWithPrefix result to contain the same elements as ListTasks result. Missing %v", v)
		} else {
			var taskInResp types.Task
			json.Unmarshal([]byte(value), &taskInResp)
			assert.Exactly(suite.T(), taskInResp, v, "Expected GetWithPrefix result to contain the same elements as ListTasks result.")
		}
	}
}

func (suite *TaskStoreTestSuite) TestFilterTasksNoFilters() {
	var filters map[string]string
	_, err := suite.taskStore.FilterTasks(filters)
	assert.Error(suite.T(), err, "Expected an error when filter map is empty")
}

func (suite *TaskStoreTestSuite) TestFilterTasksNoFilterValues() {
	_, err := suite.taskStore.FilterTasks(map[string]string{taskStatusFilter: ""})
	assert.Error(suite.T(), err, "Expected an error when all filter values are empty")
}

func (suite *TaskStoreTestSuite) TestFilterTasksSomeFilterValues() {
	suite.datastore.EXPECT().GetWithPrefix(taskKeyPrefix).Return(make(map[string]string), nil)

	tasks, err := suite.taskStore.FilterTasks(map[string]string{taskStatusFilter: "randomVal", taskClusterFilter: ""})

	assert.Nil(suite.T(), err, "Unexpected error when GetWithPrefix returns empty")
	assert.NotNil(suite.T(), tasks, "Result should be empty when GetWithPrefix returns empty")
	assert.Empty(suite.T(), tasks, "Result should be empty whenGetWithPrefix returns empty")
}

func (suite *TaskStoreTestSuite) TestFilterTasksUnsupportedFilter() {
	_, err := suite.taskStore.FilterTasks(map[string]string{"invalidFilter": "value"})
	assert.Error(suite.T(), err, "Expected an error when unsupported filter key is provided")
}

func (suite *TaskStoreTestSuite) TestFilterTasksByStatusGetWithPrefixFails() {
	suite.datastore.EXPECT().GetWithPrefix(taskKeyPrefix).Return(nil, errors.New("GetWithPrefix failed"))

	_, err := suite.taskStore.FilterTasks(map[string]string{taskStatusFilter: "randomFilter"})
	assert.Error(suite.T(), err, "Expected an error when GetWithPrefix fails")
}

func (suite *TaskStoreTestSuite) TestFilterTasksByStatusGetWithPrefixReturnsNoResults() {
	suite.datastore.EXPECT().GetWithPrefix(taskKeyPrefix).Return(make(map[string]string), nil)

	tasks, err := suite.taskStore.FilterTasks(map[string]string{taskStatusFilter: "randomFilter"})

	assert.Nil(suite.T(), err, "Unexpected error when GetWithPrefix returns empty")
	assert.NotNil(suite.T(), tasks, "Result should be empty when GetWithPrefix returns empty")
	assert.Empty(suite.T(), tasks, "Result should be empty whenGetWithPrefix returns empty")
}

func (suite *TaskStoreTestSuite) TestFilterTasksByStatusGetWithPrefixReturnsMultipleResultsNoneMatchFilter() {
	resp := map[string]string{
		taskARN1: suite.firstPendingTaskJSON,
		taskARN2: suite.secondPendingTaskJSON,
	}
	suite.datastore.EXPECT().GetWithPrefix(taskKeyPrefix).Return(resp, nil)
	tasks, err := suite.taskStore.FilterTasks(map[string]string{taskStatusFilter: "randomFilter"})
	assert.Nil(suite.T(), err, "Unexpected error when filter does not match")
	assert.NotNil(suite.T(), tasks, "Result should be empty when filter does not match")
	assert.Empty(suite.T(), tasks, "Result should be empty when filter does not match")
}

func (suite *TaskStoreTestSuite) TestFilterTasksByStatusGetWithPrefixReturnsMultipleResultsOneMatchesFilter() {
	filterStatus := "testStatus"
	taskMatchingStatus := suite.secondPendingTask
	taskMatchingStatus.Detail.LastStatus = &filterStatus
	taskMatchingStatusJSON, err := json.Marshal(taskMatchingStatus)
	assert.Nil(suite.T(), err, "Error when json marhsaling task %v", taskMatchingStatusJSON)
	resp := map[string]string{
		taskARN1: suite.firstPendingTaskJSON,
		taskARN2: string(taskMatchingStatusJSON),
	}
	suite.datastore.EXPECT().GetWithPrefix(taskKeyPrefix).Return(resp, nil)
	tasks, err := suite.taskStore.FilterTasks(map[string]string{taskStatusFilter: filterStatus})
	assert.Nil(suite.T(), err, "Unexpected error when calling filter tasks")
	assert.Equal(suite.T(), 1, len(tasks), "Expected the length of the FilterTasks result to be 1")
	assert.Exactly(suite.T(), taskMatchingStatus, tasks[0], "Expected one result when one matches filter")
}

func (suite *TaskStoreTestSuite) TestFilterTasksByStatusGetWithPrefixReturnsMultipleResultsMultipleMatchFilter() {
	resp := map[string]string{
		taskARN1: suite.firstPendingTaskJSON,
		taskARN2: suite.secondPendingTaskJSON,
	}

	suite.datastore.EXPECT().GetWithPrefix(taskKeyPrefix).Return(resp, nil)
	tasks, err := suite.taskStore.FilterTasks(map[string]string{taskStatusFilter: pendingStatus})
	assert.Nil(suite.T(), err, "Unexpected error when calling filter tasks")
	assert.Equal(suite.T(), 2, len(tasks), "Expected one result when multiple match filter")
	for _, v := range tasks {
		//attempt to grab the same task from resp
		value, ok := resp[*v.Detail.TaskARN]
		if !ok {
			suite.T().Errorf("Expected GetWithPrefix result to contain the same elements as FilterTasks result. Missing %v", v)
		} else {
			var taskInResp types.Task
			json.Unmarshal([]byte(value), &taskInResp)
			assert.Exactly(suite.T(), taskInResp, v, "Expected GetWithPrefix result to contain the same elements as FilterTasks result.")
		}
	}
}

func (suite *TaskStoreTestSuite) TestFilterTasksByStartedByGetWithPrefixFails() {
	suite.datastore.EXPECT().GetWithPrefix(taskKeyPrefix).Return(nil, errors.New("GetWithPrefix failed"))

	_, err := suite.taskStore.FilterTasks(map[string]string{taskStartedByFilter: "randomFilter"})
	assert.Error(suite.T(), err, "Expected an error when GetWithPrefix fails")
}

func (suite *TaskStoreTestSuite) TestFilterTasksByStartedByListTasksReturnsNoResults() {
	suite.datastore.EXPECT().GetWithPrefix(taskKeyPrefix).Return(make(map[string]string), nil)

	tasks, err := suite.taskStore.FilterTasks(map[string]string{taskStartedByFilter: "randomFilter"})

	assert.Nil(suite.T(), err, "Unexpected error when list tasks returns empty")
	assert.NotNil(suite.T(), tasks, "Result should be empty when lists tasks is empty")
	assert.Empty(suite.T(), tasks, "Result should be empty when lists tasks is empty")
}

func (suite *TaskStoreTestSuite) TestFilterTasksByStartedByListTasksReturnsMultipleResultsNoneMatchFilter() {
	resp := map[string]string{
		taskARN1: suite.firstTaskStartedBySomeoneElseJSON,
		taskARN2: suite.secondTaskStartedBySomeoneElseJSON,
	}

	suite.datastore.EXPECT().GetWithPrefix(taskKeyPrefix).Return(resp, nil)

	tasks, err := suite.taskStore.FilterTasks(map[string]string{taskStartedByFilter: "randomFilter"})

	assert.Nil(suite.T(), err, "Unexpected error when filter does not match")
	assert.NotNil(suite.T(), tasks, "Result should be empty when filter does not match")
	assert.Empty(suite.T(), tasks, "Result should be empty when filter does not match")
}

func (suite *TaskStoreTestSuite) TestFilterTasksByStartedByListTasksReturnsMultipleResultsOneMatchesFilter() {
	filterStartedBy := "accident"

	taskMatchingStartedBy := suite.secondTaskStartedBySomeoneElse
	taskMatchingStartedBy.Detail.StartedBy = filterStartedBy
	taskMatchingStartedByJSON, err := json.Marshal(taskMatchingStartedBy)
	assert.Nil(suite.T(), err, "Error when json marhsaling task %v", taskMatchingStartedByJSON)

	resp := map[string]string{
		taskARN1: suite.firstTaskStartedBySomeoneElseJSON,
		taskARN2: string(taskMatchingStartedByJSON),
	}

	suite.datastore.EXPECT().GetWithPrefix(taskKeyPrefix).Return(resp, nil)

	tasks, err := suite.taskStore.FilterTasks(map[string]string{taskStartedByFilter: filterStartedBy})

	assert.Nil(suite.T(), err, "Unexpected error when calling filter tasks")
	assert.Equal(suite.T(), 1, len(tasks), "Expected the length of the FilterTasks result to be 1")
	assert.Exactly(suite.T(), taskMatchingStartedBy, tasks[0], "Expected one result when one matches filter")
}

func (suite *TaskStoreTestSuite) TestFilterTasksByStartedByListTasksReturnsMultipleResultsMultipleMatchFilter() {
	resp := map[string]string{
		taskARN1: suite.firstTaskStartedBySomeoneElseJSON,
		taskARN2: suite.secondTaskStartedBySomeoneElseJSON,
	}

	suite.datastore.EXPECT().GetWithPrefix(taskKeyPrefix).Return(resp, nil)

	tasks, err := suite.taskStore.FilterTasks(map[string]string{taskStartedByFilter: someoneElse})

	assert.Nil(suite.T(), err, "Unexpected error when calling filter tasks")
	assert.Equal(suite.T(), 2, len(tasks), "Expected multiple results when multiple match filter")

	for _, v := range tasks {
		//attempt to grab the same task from resp
		value, ok := resp[*v.Detail.TaskARN]
		if !ok {
			suite.T().Errorf("Expected GetWithPrefix result to contain the same elements as FilterTasks result. Missing %v", v)
		} else {
			var taskInResp types.Task
			json.Unmarshal([]byte(value), &taskInResp)
			assert.Exactly(suite.T(), taskInResp, v, "Expected GetWithPrefix result to contain the same elements as FilterTasks result.")
		}
	}
}

func (suite *TaskStoreTestSuite) TestFilterTasksByClusterNameGetWithPrefixFails() {
	clusterKey := taskKeyPrefix + clusterName1 + "/"
	suite.datastore.EXPECT().GetWithPrefix(clusterKey).Return(nil, errors.New("GetTasksByKeyPrefix failed"))

	_, err := suite.taskStore.FilterTasks(map[string]string{taskClusterFilter: clusterName1})
	assert.Error(suite.T(), err, "Expected an error when GetWithPrefix fails")
}

func (suite *TaskStoreTestSuite) TestFilterTasksByClusterARNGetWithPrefixFails() {
	clusterKey := taskKeyPrefix + clusterName1 + "/"
	suite.datastore.EXPECT().GetWithPrefix(clusterKey).Return(nil, errors.New("GetTasksByKeyPrefix failed"))

	_, err := suite.taskStore.FilterTasks(map[string]string{taskClusterFilter: clusterARN1})
	assert.Error(suite.T(), err, "Expected an error when GetWithPrefix fails")
}

func (suite *TaskStoreTestSuite) TestFilterTasksByClusterNameGetWithPrefixReturnsTasks() {
	resp := map[string]string{
		taskARN1: suite.firstTaskOfFirstClusterJSON,
		taskARN2: suite.secondTaskOfFirstClusterJSON,
	}

	clusterKey := taskKeyPrefix + clusterName1 + "/"
	suite.datastore.EXPECT().GetWithPrefix(clusterKey).Return(resp, nil)

	tasks, err := suite.taskStore.FilterTasks(map[string]string{taskClusterFilter: clusterName1})

	assert.Nil(suite.T(), err, "Unexpected error when calling filter tasks")
	assert.Equal(suite.T(), len(resp), len(tasks), "Expected multiple results when multiple match filter")

	for _, v := range tasks {
		//attempt to grab the same task from resp
		value, ok := resp[*v.Detail.TaskARN]
		if !ok {
			suite.T().Errorf("Expected GetWithPrefix result to contain the same elements as FilterTasks result. Missing %v", v)
		} else {
			var taskInResp types.Task
			json.Unmarshal([]byte(value), &taskInResp)
			assert.Exactly(suite.T(), taskInResp, v, "Expected GetWithPrefix result to contain the same elements as FilterTasks result.")
		}
	}
}

func (suite *TaskStoreTestSuite) TestFilterTasksByClusterARNGetWithPrefixReturnsTasks() {
	resp := map[string]string{
		taskARN1: suite.firstTaskOfFirstClusterJSON,
		taskARN2: suite.secondTaskOfFirstClusterJSON,
	}

	clusterKey := taskKeyPrefix + clusterName1 + "/"
	suite.datastore.EXPECT().GetWithPrefix(clusterKey).Return(resp, nil)

	tasks, err := suite.taskStore.FilterTasks(map[string]string{taskClusterFilter: clusterARN1})

	assert.Nil(suite.T(), err, "Unexpected error when calling filter tasks")
	assert.Equal(suite.T(), len(resp), len(tasks), "Expected multiple results when multiple match filter")

	for _, v := range tasks {
		//attempt to grab the same task from resp
		value, ok := resp[*v.Detail.TaskARN]
		if !ok {
			suite.T().Errorf("Expected GetWithPrefix result to contain the same elements as FilterTasks result. Missing %v", v)
		} else {
			var taskInResp types.Task
			json.Unmarshal([]byte(value), &taskInResp)
			assert.Exactly(suite.T(), taskInResp, v, "Expected GetWithPrefix result to contain the same elements as FilterTasks result.")
		}
	}
}

func (suite *TaskStoreTestSuite) TestFilterTasksByClusterAndStartedBy() {
	cluster1SomeoneTask := types.Task{
		Detail: &types.TaskDetail{
			TaskARN:    &taskARN1,
			ClusterARN: &clusterARN1,
			LastStatus: &pendingStatus,
			StartedBy:  someoneElse,
		},
	}
	cluster1SomeoneTaskJSON := suite.setupTask(cluster1SomeoneTask)

	cluster1RandomTask := types.Task{
		Detail: &types.TaskDetail{
			TaskARN:    &taskARN2,
			ClusterARN: &clusterARN1,
			LastStatus: &pendingStatus,
			StartedBy:  "random",
		},
	}
	cluster1RandomTaskJSON := suite.setupTask(cluster1RandomTask)
	resp := map[string]string{
		taskARN1: cluster1SomeoneTaskJSON,
		taskARN2: cluster1RandomTaskJSON,
	}

	clusterKey := taskKeyPrefix + clusterName1 + "/"
	suite.datastore.EXPECT().GetWithPrefix(clusterKey).Return(resp, nil)

	tasks, err := suite.taskStore.FilterTasks(
		map[string]string{taskClusterFilter: clusterARN1, taskStartedByFilter: someoneElse})

	assert.Nil(suite.T(), err, "Unexpected error when calling filter tasks")
	assert.Equal(suite.T(), 1, len(tasks))
	assert.Exactly(suite.T(), cluster1SomeoneTask, tasks[0])
}

func (suite *TaskStoreTestSuite) TestFilterTasksByClusterAndStatus() {
	cluster1PendingTask := types.Task{
		Detail: &types.TaskDetail{
			TaskARN:    &taskARN1,
			ClusterARN: &clusterARN1,
			LastStatus: &pendingStatus,
		},
	}
	cluster1PendingTaskJSON := suite.setupTask(cluster1PendingTask)

	cluster1RunningTask := types.Task{
		Detail: &types.TaskDetail{
			TaskARN:    &taskARN2,
			ClusterARN: &clusterARN1,
			LastStatus: &runningStatus,
		},
	}
	cluster1RunningTaskJSON := suite.setupTask(cluster1RunningTask)
	resp := map[string]string{
		taskARN1: cluster1PendingTaskJSON,
		taskARN2: cluster1RunningTaskJSON,
	}

	clusterKey := taskKeyPrefix + clusterName1 + "/"
	suite.datastore.EXPECT().GetWithPrefix(clusterKey).Return(resp, nil)

	tasks, err := suite.taskStore.FilterTasks(
		map[string]string{taskClusterFilter: clusterARN1, taskStatusFilter: pendingStatus})

	assert.Nil(suite.T(), err, "Unexpected error when calling filter tasks")
	assert.Equal(suite.T(), 1, len(tasks))
	assert.Exactly(suite.T(), cluster1PendingTask, tasks[0])
}

func (suite *TaskStoreTestSuite) TestFilterTasksByStartedByAndStatus() {
	cluster1PendingSomeoneTask := types.Task{
		Detail: &types.TaskDetail{
			TaskARN:    &taskARN1,
			ClusterARN: &clusterARN1,
			LastStatus: &pendingStatus,
			StartedBy:  someoneElse,
		},
	}
	cluster1PendingSomeoneTaskJSON := suite.setupTask(cluster1PendingSomeoneTask)

	cluster1RunningRandomTask := types.Task{
		Detail: &types.TaskDetail{
			TaskARN:    &taskARN2,
			ClusterARN: &clusterARN1,
			LastStatus: &runningStatus,
			StartedBy:  "random",
		},
	}
	cluster1RunningRandomTaskJSON := suite.setupTask(cluster1RunningRandomTask)

	cluster1PendingRandomTask := types.Task{
		Detail: &types.TaskDetail{
			TaskARN:    &taskARN2,
			ClusterARN: &clusterARN1,
			LastStatus: &pendingStatus,
			StartedBy:  "random",
		},
	}
	cluster1PendingRandomTaskJSON := suite.setupTask(cluster1PendingRandomTask)

	resp := map[string]string{
		taskARN1: cluster1PendingSomeoneTaskJSON,
		taskARN2: cluster1RunningRandomTaskJSON,
		taskARN3: cluster1PendingRandomTaskJSON,
	}

	suite.datastore.EXPECT().GetWithPrefix(taskKeyPrefix).Return(resp, nil)

	tasks, err := suite.taskStore.FilterTasks(
		map[string]string{taskStartedByFilter: "random", taskStatusFilter: pendingStatus})

	assert.Nil(suite.T(), err, "Unexpected error when calling filter tasks")
	assert.Equal(suite.T(), 1, len(tasks))
	assert.Exactly(suite.T(), cluster1PendingRandomTask, tasks[0])
}

func (suite *TaskStoreTestSuite) TestFilterTasksByStartedByAndClusterAndStatus() {
	cluster1PendingSomeoneTask := types.Task{
		Detail: &types.TaskDetail{
			TaskARN:    &taskARN1,
			ClusterARN: &clusterARN1,
			LastStatus: &pendingStatus,
			StartedBy:  someoneElse,
		},
	}
	cluster1PendingSomeoneTaskJSON := suite.setupTask(cluster1PendingSomeoneTask)

	cluster1RunningRandomTask := types.Task{
		Detail: &types.TaskDetail{
			TaskARN:    &taskARN2,
			ClusterARN: &clusterARN1,
			LastStatus: &runningStatus,
			StartedBy:  "random",
		},
	}
	cluster1RunningRandomTaskJSON := suite.setupTask(cluster1RunningRandomTask)

	cluster1PendingRandomTask := types.Task{
		Detail: &types.TaskDetail{
			TaskARN:    &taskARN2,
			ClusterARN: &clusterARN1,
			LastStatus: &pendingStatus,
			StartedBy:  "random",
		},
	}
	cluster1PendingRandomTaskJSON := suite.setupTask(cluster1PendingRandomTask)

	resp := map[string]string{
		taskARN1: cluster1PendingSomeoneTaskJSON,
		taskARN2: cluster1RunningRandomTaskJSON,
		taskARN3: cluster1PendingRandomTaskJSON,
	}

	clusterKey := taskKeyPrefix + clusterName1 + "/"
	suite.datastore.EXPECT().GetWithPrefix(clusterKey).Return(resp, nil)

	tasks, err := suite.taskStore.FilterTasks(
		map[string]string{taskClusterFilter: clusterARN1,
			taskStartedByFilter: "random",
			taskStatusFilter:    pendingStatus})

	assert.Nil(suite.T(), err, "Unexpected error when calling filter tasks")
	assert.Equal(suite.T(), 1, len(tasks))
	assert.Exactly(suite.T(), cluster1PendingRandomTask, tasks[0])
}

func (suite *TaskStoreTestSuite) TestStreamTasksDataStoreStreamReturnsError() {
	ctx := context.Background()
	suite.datastore.EXPECT().StreamWithPrefix(gomock.Any(), taskKeyPrefix).Return(nil, errors.New("StreamWithPrefix failed"))

	taskRespChan, err := suite.taskStore.StreamTasks(ctx)
	assert.Error(suite.T(), err, "Expected an error when datastore StreamWithPrefix returns an error")
	assert.Nil(suite.T(), taskRespChan, "Unexpected task response channel when there is a datastore channel setup error")
}

func (suite *TaskStoreTestSuite) TestStreamTasksValidJSONInDSChannel() {
	ctx := context.Background()
	dsChan := make(chan map[string]string)
	defer close(dsChan)
	suite.datastore.EXPECT().StreamWithPrefix(gomock.Any(), taskKeyPrefix).Return(dsChan, nil)

	taskRespChan, err := suite.taskStore.StreamTasks(ctx)
	assert.Nil(suite.T(), err, "Unexpected error when calling stream tasks")
	assert.NotNil(suite.T(), taskRespChan)

	taskResp := addTaskToDSChanAndReadFromTaskRespChan(suite.firstPendingTaskJSON, dsChan, taskRespChan)

	assert.Nil(suite.T(), taskResp.Err, "Unexpected error when reading task from channel")
	assert.Equal(suite.T(), suite.firstPendingTask, taskResp.Task, "Expected task in task response to match that in the stream")
}

func (suite *TaskStoreTestSuite) TestStreamTasksInvalidJSONInDSChannel() {
	ctx := context.Background()
	dsChan := make(chan map[string]string)
	defer close(dsChan)
	suite.datastore.EXPECT().StreamWithPrefix(gomock.Any(), taskKeyPrefix).Return(dsChan, nil)

	taskRespChan, err := suite.taskStore.StreamTasks(ctx)
	assert.Nil(suite.T(), err, "Unexpected error when calling stream tasks")
	assert.NotNil(suite.T(), taskRespChan)

	taskResp := addTaskToDSChanAndReadFromTaskRespChan("invalidJSON", dsChan, taskRespChan)

	assert.Error(suite.T(), taskResp.Err, "Expected an error when dsChannel returns an invalid task json")
	assert.Equal(suite.T(), types.Task{}, taskResp.Task, "Expected empty task in response when there is a decode error")

	_, ok := <-taskRespChan
	assert.False(suite.T(), ok, "Expected task response channel to be closed")
}

func (suite *TaskStoreTestSuite) TestStreamTasksCancelUpstreamContext() {
	ctx, cancel := context.WithCancel(context.Background())
	dsChan := make(chan map[string]string)
	defer close(dsChan)
	suite.datastore.EXPECT().StreamWithPrefix(gomock.Any(), taskKeyPrefix).Return(dsChan, nil)

	taskRespChan, err := suite.taskStore.StreamTasks(ctx)
	assert.Nil(suite.T(), err, "Unexpected error when calling stream tasks")
	assert.NotNil(suite.T(), taskRespChan)

	cancel()

	_, ok := <-taskRespChan
	assert.False(suite.T(), ok, "Expected task response channel to be closed")
}

func (suite *TaskStoreTestSuite) TestStreamTasksCloseDownstreamChannel() {
	ctx := context.Background()
	dsChan := make(chan map[string]string)
	suite.datastore.EXPECT().StreamWithPrefix(gomock.Any(), taskKeyPrefix).Return(dsChan, nil)

	taskRespChan, err := suite.taskStore.StreamTasks(ctx)
	assert.Nil(suite.T(), err, "Unexpected error when calling stream tasks")
	assert.NotNil(suite.T(), taskRespChan)

	close(dsChan)

	_, ok := <-taskRespChan
	assert.False(suite.T(), ok, "Expected task response channel to be closed")
}

func (suite *TaskStoreTestSuite) TestDeleteTaskEmptyClusterName() {
	err := suite.taskStore.DeleteTask("", taskARN1)
	assert.Error(suite.T(), err, "Expected an error when cluster name is empty in DeleteTask")
}

func (suite *TaskStoreTestSuite) TestDeleteTaskEmptyTaskARN() {
	err := suite.taskStore.DeleteTask(clusterName1, "")
	assert.Error(suite.T(), err, "Expected an error when task ARN is empty in DeleteTask")
}

func (suite *TaskStoreTestSuite) TestDeleteTaskDeleteTaskFails() {
	suite.datastore.EXPECT().Delete(suite.taskKey1).Return(int64(0), errors.New("Error when deleting key"))

	err := suite.taskStore.DeleteTask(clusterName1, taskARN1)
	assert.Error(suite.T(), err, "Expected an error when delete task fails")
}

func (suite *TaskStoreTestSuite) TestDeleteTaskDeleteNoError() {
	suite.datastore.EXPECT().Delete(suite.taskKey1).Return(int64(1), nil)

	err := suite.taskStore.DeleteTask(clusterName1, taskARN1)
	assert.NoError(suite.T(), err, "Error when deleting task")
}

func (suite *TaskStoreTestSuite) TestDeleteTaskDeleteWithClusterNameAndTaskARN() {
	suite.datastore.EXPECT().Delete(suite.taskKey1).Return(int64(1), nil)

	err := suite.taskStore.DeleteTask(clusterARN1, taskARN1)
	assert.NoError(suite.T(), err, "Error when deleting task")
}

func addTaskToDSChanAndReadFromTaskRespChan(taskToAdd string, dsChan chan map[string]string, taskRespChan chan storetypes.TaskErrorWrapper) storetypes.TaskErrorWrapper {
	var taskResp storetypes.TaskErrorWrapper

	doneChan := make(chan bool)
	defer close(doneChan)
	go func() {
		taskResp = <-taskRespChan
		doneChan <- true
	}()

	dsVal := map[string]string{taskARN1: taskToAdd}
	dsChan <- dsVal
	<-doneChan

	return taskResp
}
