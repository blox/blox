package store

import (
	"encoding/json"
	"github.com/aws/amazon-ecs-event-stream-handler/handler/mocks"
	"github.com/aws/amazon-ecs-event-stream-handler/handler/types"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

const (
	taskArn       = "arn:aws:ecs:us-east-1:159403520677:task/271022c0-f894-4aa2-b063-25bae55088d5"
	taskArn2      = "arn:aws:ecs:us-east-1:159403520677:task/345022c0-f894-4aa2-b063-25bae55088d5"
	pendingStatus = "pending"
)

type TaskStoreTestSuite struct {
	suite.Suite
	datastore *mocks.MockDataStore
	taskStore TaskStore
	taskKey   string
	task      types.Task
	task2     types.Task
	taskJSON  string
	task2JSON string
}

func (suite *TaskStoreTestSuite) SetupTest() {
	mockCtrl := gomock.NewController(suite.T())
	suite.datastore = mocks.NewMockDataStore(mockCtrl)
	suite.taskKey = taskKeyPrefix + taskArn

	var err error
	suite.taskStore, err = NewTaskStore(suite.datastore)
	assert.Nil(suite.T(), err, "Cannot setup testSuite: Unexpected error when calling NewTaskStore")

	suite.task = types.Task{}
	suite.task.Detail.TaskArn = taskArn
	suite.task.Detail.Version = 1
	suite.task.Detail.LastStatus = pendingStatus

	taskJSON, err := json.Marshal(suite.task)
	assert.Nil(suite.T(), err, "Cannot setup testSuite: Error when json marhsaling task %v", suite.task)

	suite.task2 = suite.task
	suite.task2.Detail.TaskArn = taskArn2
	suite.task2.Detail.Version++

	task2JSON, err := json.Marshal(suite.task2)
	assert.Nil(suite.T(), err, "Cannot setup testSuite: Error when json marhsaling task %v", suite.task2)

	suite.taskJSON = string(taskJSON)
	suite.task2JSON = string(task2JSON)
}

func TestTaskStoreTestSuite(t *testing.T) {
	suite.Run(t, new(TaskStoreTestSuite))
}

func (suite *TaskStoreTestSuite) TestNewTaskStoreNilDatastore() {
	_, err := NewTaskStore(nil)
	assert.Error(suite.T(), err, "Expected an error when datastore is nil")
}

func (suite *TaskStoreTestSuite) TestNewTaskStore() {
	taskStore, err := NewTaskStore(suite.datastore)
	assert.Nil(suite.T(), err, "Unexpected error when calling NewTaskStore")
	assert.NotNil(suite.T(), taskStore, "TaskStore should not be nil")
}

func (suite *TaskStoreTestSuite) TestAddTaskEmptyJSON() {
	err := suite.taskStore.AddTask("")
	assert.Error(suite.T(), err, "Expected an error when json empty in AddTask")
}

func (suite *TaskStoreTestSuite) TestAddTaskUnmarshalJSONError() {
	err := suite.taskStore.AddTask("invalidJson")
	assert.Error(suite.T(), err, "Expected an error when json invalid in AddTask")
}

func (suite *TaskStoreTestSuite) TestAddTaskEmptyTask() {
	task := types.Task{}
	taskJSON, err := json.Marshal(task)
	err = suite.taskStore.AddTask(string(taskJSON))
	assert.Error(suite.T(), err, "Expected an error when task arn is not set")
}

func (suite *TaskStoreTestSuite) TestAddTaskGetTaskFails() {
	suite.datastore.EXPECT().Get(suite.taskKey).Return(nil, errors.New("Error when getting key"))

	err := suite.taskStore.AddTask(suite.taskJSON)
	assert.Error(suite.T(), err, "Expected an error when get task fails")
}

func (suite *TaskStoreTestSuite) TestAddTaskGetTaskNoResults() {
	suite.datastore.EXPECT().Get(suite.taskKey).Return(make(map[string]string), nil)
	suite.datastore.EXPECT().Add(suite.taskKey, suite.taskJSON).Return(nil)

	err := suite.taskStore.AddTask(suite.taskJSON)
	assert.Nil(suite.T(), err, "Unexpected error when datastore returns empty results")
}

func (suite *TaskStoreTestSuite) TestAddTaskGetTaskMultipleResults() {
	resp := map[string]string{
		taskArn:  suite.taskJSON,
		taskArn2: suite.task2JSON,
	}

	suite.datastore.EXPECT().Get(suite.taskKey).Return(resp, nil)

	err := suite.taskStore.AddTask(suite.taskJSON)
	assert.Error(suite.T(), err, "Expected an error when datastore returns multiple results")
}

func (suite *TaskStoreTestSuite) TestAddTaskGetTaskInvalidJSONResult() {
	resp := map[string]string{
		taskArn: "invalidJson",
	}

	suite.datastore.EXPECT().Get(suite.taskKey).Return(resp, nil)

	err := suite.taskStore.AddTask(suite.taskJSON)
	assert.Error(suite.T(), err, "Expected an error when datastore returns invalid json results")
}

func (suite *TaskStoreTestSuite) TestAddTaskSameVersionTaskExists() {
	resp := map[string]string{
		taskArn: suite.taskJSON,
	}

	suite.datastore.EXPECT().Get(suite.taskKey).Return(resp, nil)
	suite.datastore.EXPECT().Add(suite.taskKey, suite.taskJSON).Times(0)

	err := suite.taskStore.AddTask(suite.taskJSON)
	assert.Nil(suite.T(), err, "Unexpected error when same version task exists")
}

func (suite *TaskStoreTestSuite) TestAddTaskHigherVersionTaskExists() {
	resp := map[string]string{
		taskArn: suite.task2JSON,
	}

	suite.datastore.EXPECT().Get(suite.taskKey).Return(resp, nil)
	suite.datastore.EXPECT().Add(suite.taskKey, gomock.Any()).Times(0)

	err := suite.taskStore.AddTask(suite.taskJSON)
	assert.Nil(suite.T(), err, "Unexpected error when higher version task exists")
}

func (suite *TaskStoreTestSuite) TestAddTaskLowerVersionTaskExists() {
	task := types.Task{}
	task.Detail.TaskArn = taskArn
	task.Detail.Version = suite.task.Detail.Version - 1

	taskJSON, err := json.Marshal(task)
	assert.Nil(suite.T(), err, "Error when json marhsaling task %v", task)

	resp := map[string]string{
		taskArn: string(taskJSON),
	}

	suite.datastore.EXPECT().Get(suite.taskKey).Return(resp, nil)
	suite.datastore.EXPECT().Add(suite.taskKey, suite.taskJSON).Return(nil)

	err = suite.taskStore.AddTask(suite.taskJSON)
	assert.Nil(suite.T(), err, "Unexpected error when lower version task exists")
}

func (suite *TaskStoreTestSuite) TestAddTaskFails() {
	suite.datastore.EXPECT().Get(suite.taskKey).Return(make(map[string]string), nil)
	suite.datastore.EXPECT().Add(suite.taskKey, suite.taskJSON).Return(errors.New("Add task failed"))

	err := suite.taskStore.AddTask(suite.taskJSON)
	assert.Error(suite.T(), err, "Expected an error when add task fails")
}

func (suite *TaskStoreTestSuite) TestGetTaskEmptyArn() {
	_, err := suite.taskStore.GetTask("")
	assert.Error(suite.T(), err, "Expected an error when arn empty in GetTask")
}

func (suite *TaskStoreTestSuite) TestGetTaskGetTaskFails() {
	suite.datastore.EXPECT().Get(suite.taskKey).Return(nil, errors.New("Error when getting key"))

	_, err := suite.taskStore.GetTask(taskArn)
	assert.Error(suite.T(), err, "Expected an error when get task fails")
}

func (suite *TaskStoreTestSuite) TestGetTaskGetTaskNoResults() {
	suite.datastore.EXPECT().Get(suite.taskKey).Return(make(map[string]string), nil)

	task, err := suite.taskStore.GetTask(taskArn)
	assert.Nil(suite.T(), err, "Unexpected error when datastore returns empty results")
	assert.Nil(suite.T(), task, "Unexpected object returned when datastore returns empty results")
}

func (suite *TaskStoreTestSuite) TestGetTaskGetMultipleResults() {
	resp := map[string]string{
		taskArn:  suite.taskJSON,
		taskArn2: suite.task2JSON,
	}

	suite.datastore.EXPECT().Get(suite.taskKey).Return(resp, nil)

	_, err := suite.taskStore.GetTask(taskArn)
	assert.Error(suite.T(), err, "Expected an error when datastore returns multiple results")
}

func (suite *TaskStoreTestSuite) TestGetTaskGetInvalidJSONResult() {
	resp := map[string]string{
		taskArn: "invalidJson",
	}

	suite.datastore.EXPECT().Get(suite.taskKey).Return(resp, nil)

	_, err := suite.taskStore.GetTask(taskArn)
	assert.Error(suite.T(), err, "Expected an error when datastore returns invalid json results")
}

func (suite *TaskStoreTestSuite) TestGetTask() {
	resp := map[string]string{
		taskArn: suite.taskJSON,
	}

	suite.datastore.EXPECT().Get(suite.taskKey).Return(resp, nil)

	task, err := suite.taskStore.GetTask(taskArn)
	assert.Nil(suite.T(), err, "Unexpected error when getting task")
	assert.NotNil(suite.T(), task, "Expected a non-nil task when calling GetTask")

	assert.Exactly(suite.T(), suite.task, *task, "Expected the returned task to match the one returned from the datastore")
}

func (suite *TaskStoreTestSuite) TestListTasksGetWithPrefixInvalidJSON() {
	resp := map[string]string{
		taskArn: "invalidJson",
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
		taskArn:  suite.taskJSON,
		taskArn2: suite.task2JSON,
	}

	suite.datastore.EXPECT().GetWithPrefix(taskKeyPrefix).Return(resp, nil)

	tasks, err := suite.taskStore.ListTasks()
	assert.Nil(suite.T(), err, "Unexpected error when GetWithPrefix returns multiple results")
	assert.NotNil(suite.T(), tasks, "Expected a non-nil result when GetWithPrefix returns multiple results")

	assert.Equal(suite.T(), len(resp), len(tasks), "Expected ListTasks result to be of the same length as GetWithPrefix result")

	for _, v := range tasks {
		//attempt to grab the same task from resp
		value, ok := resp[v.Detail.TaskArn]
		if !ok {
			suite.T().Errorf("Expected GetWithPrefix result to contain the same elements as ListTasks result. Missing %v", v)
		} else {
			var taskInResp types.Task
			json.Unmarshal([]byte(value), &taskInResp)
			assert.Exactly(suite.T(), taskInResp, v, "Expected GetWithPrefix result to contain the same elements as ListTasks result.")
		}
	}
}

func (suite *TaskStoreTestSuite) TestFilterTasksEmptyKey() {
	_, err := suite.taskStore.FilterTasks("", "value")
	assert.Error(suite.T(), err, "Expected an error when filterKey is empty")
}

func (suite *TaskStoreTestSuite) TestFilterTasksEmptyValue() {
	_, err := suite.taskStore.FilterTasks(statusFilter, "")
	assert.Error(suite.T(), err, "Expected an error when filterKey is empty")
}

func (suite *TaskStoreTestSuite) TestFilterTasksUnsupportedFilter() {
	_, err := suite.taskStore.FilterTasks("invalidFilter", "started")
	assert.Error(suite.T(), err, "Expected an error when unsupported filter key is provided")
}

func (suite *TaskStoreTestSuite) TestFilterTasksListTasksGetWithPrefixFails() {
	suite.datastore.EXPECT().GetWithPrefix(taskKeyPrefix).Return(nil, errors.New("GetWithPrefix failed"))

	_, err := suite.taskStore.FilterTasks(statusFilter, "randomFilter")
	assert.Error(suite.T(), err, "Expected an error when list tasks fails")
}

func (suite *TaskStoreTestSuite) TestFilterTasksListTasksReturnsNoResults() {
	suite.datastore.EXPECT().GetWithPrefix(taskKeyPrefix).Return(make(map[string]string), nil)

	tasks, err := suite.taskStore.FilterTasks(statusFilter, "randomFilter")

	assert.Nil(suite.T(), err, "Unexpected error when list tasks returns empty")
	assert.NotNil(suite.T(), tasks, "Result should be empty when lists tasks is empty")
	assert.Empty(suite.T(), tasks, "Result should be empty when lists tasks is empty")
}

func (suite *TaskStoreTestSuite) TestFilterTasksListTasksReturnsMultipleResultsNoneMatchFilter() {
	resp := map[string]string{
		taskArn:  suite.taskJSON,
		taskArn2: suite.task2JSON,
	}

	suite.datastore.EXPECT().GetWithPrefix(taskKeyPrefix).Return(resp, nil)

	tasks, err := suite.taskStore.FilterTasks(statusFilter, "randomFilter")

	assert.Nil(suite.T(), err, "Unexpected error when filter does not match")
	assert.NotNil(suite.T(), tasks, "Result should be empty when filter does not match")
	assert.Empty(suite.T(), tasks, "Result should be empty when filter does not match")
}

func (suite *TaskStoreTestSuite) TestFilterTasksListTasksReturnsMultipleResultsOneMatchesFilter() {
	filterStatus := "testStatus"

	taskMatchingStatus := suite.task2
	taskMatchingStatus.Detail.LastStatus = filterStatus
	taskMatchingStatusJSON, err := json.Marshal(taskMatchingStatus)

	resp := map[string]string{
		taskArn:  suite.taskJSON,
		taskArn2: string(taskMatchingStatusJSON),
	}

	suite.datastore.EXPECT().GetWithPrefix(taskKeyPrefix).Return(resp, nil)

	tasks, err := suite.taskStore.FilterTasks(statusFilter, filterStatus)

	assert.Nil(suite.T(), err, "Unexpected error when calling filter tasks")
	assert.Equal(suite.T(), 1, len(tasks), "Expected the length of the FilterTasks result to be 1")
	assert.Exactly(suite.T(), taskMatchingStatus, tasks[0], "Expected one result when one matches filter")
}

func (suite *TaskStoreTestSuite) TestFilterTasksListTasksReturnsMultipleResultsMultipleMatchFilter() {
	resp := map[string]string{
		taskArn:  suite.taskJSON,
		taskArn2: suite.task2JSON,
	}

	suite.datastore.EXPECT().GetWithPrefix(taskKeyPrefix).Return(resp, nil)

	tasks, err := suite.taskStore.FilterTasks(statusFilter, suite.task.Detail.LastStatus)

	assert.Nil(suite.T(), err, "Unexpected error when calling filter tasks")
	assert.Equal(suite.T(), 2, len(tasks), "Expected one result when multiple match filter")

	for _, v := range tasks {
		//attempt to grab the same task from resp
		value, ok := resp[v.Detail.TaskArn]
		if !ok {
			suite.T().Errorf("Expected GetWithPrefix result to contain the same elements as FilterTasks result. Missing %v", v)
		} else {
			var taskInResp types.Task
			json.Unmarshal([]byte(value), &taskInResp)
			assert.Exactly(suite.T(), taskInResp, v, "Expected GetWithPrefix result to contain the same elements as FilterTasks result.")
		}
	}
}
