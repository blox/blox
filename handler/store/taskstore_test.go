package store

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/aws/amazon-ecs-event-stream-handler/handler/compress"
	"github.com/aws/amazon-ecs-event-stream-handler/handler/mocks"
	storetypes "github.com/aws/amazon-ecs-event-stream-handler/handler/store/types"
	"github.com/aws/amazon-ecs-event-stream-handler/handler/types"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	taskArn1      = "arn:aws:ecs:us-east-1:159403520677:task/271022c0-f894-4aa2-b063-25bae55088d5"
	taskArn2      = "arn:aws:ecs:us-east-1:159403520677:task/345022c0-f894-4aa2-b063-25bae55088d5"
	pendingStatus = "pending"
)

type TaskStoreTestSuite struct {
	suite.Suite
	datastore           *mocks.MockDataStore
	taskStore           TaskStore
	taskKey             string
	task1               types.Task
	task2               types.Task
	task1JSON           string
	task2JSON           string
	compressedTask1JSON string
	compressedTask2JSON string
}

func (suite *TaskStoreTestSuite) SetupTest() {
	mockCtrl := gomock.NewController(suite.T())
	suite.datastore = mocks.NewMockDataStore(mockCtrl)
	suite.taskKey = taskKeyPrefix + taskArn1

	var err error
	suite.taskStore, err = NewTaskStore(suite.datastore)
	assert.Nil(suite.T(), err, "Cannot setup testSuite: Unexpected error when calling NewTaskStore")

	suite.task1 = types.Task{}
	suite.task1.Detail.TaskArn = taskArn1
	suite.task1.Detail.Version = 1
	suite.task1.Detail.LastStatus = pendingStatus

	task1JSON, err := json.Marshal(suite.task1)
	assert.Nil(suite.T(), err, "Cannot setup testSuite: Error when json marhsaling task %v", suite.task1)

	suite.task2 = suite.task1
	suite.task2.Detail.TaskArn = taskArn2
	suite.task2.Detail.Version++

	task2JSON, err := json.Marshal(suite.task2)
	assert.Nil(suite.T(), err, "Cannot setup testSuite: Error when json marhsaling task %v", suite.task2)

	suite.task1JSON = string(task1JSON)
	suite.task2JSON = string(task2JSON)

	compressedTask1JSON, err := compress.Compress(suite.task1JSON)
	assert.Nil(suite.T(), err, "Cannot setup testSuite: Error when compressing task json %v", suite.task1JSON)
	suite.compressedTask1JSON = string(compressedTask1JSON)

	compressedTask2JSON, err := compress.Compress(suite.task2JSON)
	assert.Nil(suite.T(), err, "Cannot setup testSuite: Error when compressing task json %v", suite.task2JSON)
	suite.compressedTask2JSON = string(compressedTask2JSON)
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
	err := suite.taskStore.AddTask("invalidJSON")
	assert.Error(suite.T(), err, "Expected an error when json invalid in AddTask")
}

func (suite *TaskStoreTestSuite) TestAddTaskEmptyTask() {
	task := types.Task{}
	task1JSON, err := json.Marshal(task)
	err = suite.taskStore.AddTask(string(task1JSON))
	assert.Error(suite.T(), err, "Expected an error when task arn is not set")
}

func (suite *TaskStoreTestSuite) TestAddTaskGetTaskFails() {
	suite.datastore.EXPECT().Get(suite.taskKey).Return(nil, errors.New("Error when getting key"))

	err := suite.taskStore.AddTask(suite.task1JSON)
	assert.Error(suite.T(), err, "Expected an error when get task fails")
}

func (suite *TaskStoreTestSuite) TestAddTaskGetTaskNoResults() {
	suite.datastore.EXPECT().Get(suite.taskKey).Return(make(map[string]string), nil)
	suite.datastore.EXPECT().Add(suite.taskKey, suite.compressedTask1JSON).Return(nil)

	err := suite.taskStore.AddTask(suite.task1JSON)
	assert.Nil(suite.T(), err, "Unexpected error when datastore returns empty results")
}

func (suite *TaskStoreTestSuite) TestAddTaskGetTaskMultipleResults() {
	resp := map[string]string{
		taskArn1: suite.compressedTask1JSON,
		taskArn2: suite.compressedTask2JSON,
	}

	suite.datastore.EXPECT().Get(suite.taskKey).Return(resp, nil)

	err := suite.taskStore.AddTask(suite.task1JSON)
	assert.Error(suite.T(), err, "Expected an error when datastore returns multiple results")
}

func (suite *TaskStoreTestSuite) TestAddTaskGetTaskInvalidJSONResult() {
	compressedInvalidJSON := suite.compressString("invalidJSON")

	resp := map[string]string{
		taskArn1: compressedInvalidJSON,
	}

	suite.datastore.EXPECT().Get(suite.taskKey).Return(resp, nil)

	err := suite.taskStore.AddTask(suite.task1JSON)
	assert.Error(suite.T(), err, "Expected an error when datastore returns invalid json results")
}

func (suite *TaskStoreTestSuite) TestAddTaskSameVersionTaskExists() {
	resp := map[string]string{
		taskArn1: suite.compressedTask1JSON,
	}

	suite.datastore.EXPECT().Get(suite.taskKey).Return(resp, nil)
	suite.datastore.EXPECT().Add(suite.taskKey, suite.task1JSON).Times(0)

	err := suite.taskStore.AddTask(suite.task1JSON)
	assert.Nil(suite.T(), err, "Unexpected error when same version task exists")
}

func (suite *TaskStoreTestSuite) TestAddTaskHigherVersionTaskExists() {
	resp := map[string]string{
		taskArn1: suite.compressedTask2JSON,
	}

	suite.datastore.EXPECT().Get(suite.taskKey).Return(resp, nil)
	suite.datastore.EXPECT().Add(suite.taskKey, gomock.Any()).Times(0)

	err := suite.taskStore.AddTask(suite.task1JSON)
	assert.Nil(suite.T(), err, "Unexpected error when higher version task exists")
}

func (suite *TaskStoreTestSuite) TestAddTaskLowerVersionTaskExists() {
	task := types.Task{}
	task.Detail.TaskArn = taskArn1
	task.Detail.Version = suite.task1.Detail.Version - 1

	taskJSON, err := json.Marshal(task)
	assert.Nil(suite.T(), err, "Error when json marhsaling task %v", task)

	compressedTaskJSON := suite.compressString(string(taskJSON))

	resp := map[string]string{
		taskArn1: compressedTaskJSON,
	}

	suite.datastore.EXPECT().Get(suite.taskKey).Return(resp, nil)
	suite.datastore.EXPECT().Add(suite.taskKey, suite.compressedTask1JSON).Return(nil)

	err = suite.taskStore.AddTask(suite.task1JSON)
	assert.Nil(suite.T(), err, "Unexpected error when lower version task exists")
}

func (suite *TaskStoreTestSuite) TestAddTaskFails() {
	suite.datastore.EXPECT().Get(suite.taskKey).Return(make(map[string]string), nil)
	suite.datastore.EXPECT().Add(suite.taskKey, suite.compressedTask1JSON).Return(errors.New("Add task failed"))

	err := suite.taskStore.AddTask(suite.task1JSON)
	assert.Error(suite.T(), err, "Expected an error when add task fails")
}

func (suite *TaskStoreTestSuite) TestGetTaskEmptyArn() {
	_, err := suite.taskStore.GetTask("")
	assert.Error(suite.T(), err, "Expected an error when arn empty in GetTask")
}

func (suite *TaskStoreTestSuite) TestGetTaskGetTaskFails() {
	suite.datastore.EXPECT().Get(suite.taskKey).Return(nil, errors.New("Error when getting key"))

	_, err := suite.taskStore.GetTask(taskArn1)
	assert.Error(suite.T(), err, "Expected an error when get task fails")
}

func (suite *TaskStoreTestSuite) TestGetTaskGetTaskNoResults() {
	suite.datastore.EXPECT().Get(suite.taskKey).Return(make(map[string]string), nil)

	task, err := suite.taskStore.GetTask(taskArn1)
	assert.Nil(suite.T(), err, "Unexpected error when datastore returns empty results")
	assert.Nil(suite.T(), task, "Unexpected object returned when datastore returns empty results")
}

func (suite *TaskStoreTestSuite) TestGetTaskGetMultipleResults() {
	resp := map[string]string{
		taskArn1: suite.compressedTask1JSON,
		taskArn2: suite.compressedTask2JSON,
	}

	suite.datastore.EXPECT().Get(suite.taskKey).Return(resp, nil)

	_, err := suite.taskStore.GetTask(taskArn1)
	assert.Error(suite.T(), err, "Expected an error when datastore returns multiple results")
}

func (suite *TaskStoreTestSuite) TestGetTaskGetInvalidJSONResult() {
	compressedInvalidJSON := suite.compressString("invalidJSON")

	resp := map[string]string{
		taskArn1: compressedInvalidJSON,
	}

	suite.datastore.EXPECT().Get(suite.taskKey).Return(resp, nil)

	_, err := suite.taskStore.GetTask(taskArn1)
	assert.Error(suite.T(), err, "Expected an error when datastore returns invalid json results")
}

func (suite *TaskStoreTestSuite) TestGetTask() {
	resp := map[string]string{
		taskArn1: suite.compressedTask1JSON,
	}

	suite.datastore.EXPECT().Get(suite.taskKey).Return(resp, nil)

	task, err := suite.taskStore.GetTask(taskArn1)
	assert.Nil(suite.T(), err, "Unexpected error when getting task")
	assert.NotNil(suite.T(), task, "Expected a non-nil task when calling GetTask")

	assert.Exactly(suite.T(), suite.task1, *task, "Expected the returned task to match the one returned from the datastore")
}

func (suite *TaskStoreTestSuite) TestListTasksGetWithPrefixInvalidJSON() {
	compressedInvalidJSON := suite.compressString("invalidJSON")

	resp := map[string]string{
		taskArn1: compressedInvalidJSON,
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
		taskArn1: suite.compressedTask1JSON,
		taskArn2: suite.compressedTask2JSON,
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
			uncompressedTaskJSON, err := compress.Uncompress([]byte(value))
			assert.Nil(suite.T(), err, "Unexpected error when uncompressing task json %v", uncompressedTaskJSON)
			var taskInResp types.Task
			json.Unmarshal([]byte(uncompressedTaskJSON), &taskInResp)
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
		taskArn1: suite.compressedTask1JSON,
		taskArn2: suite.compressedTask2JSON,
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
	assert.Nil(suite.T(), err, "Error when json marhsaling task %v", taskMatchingStatusJSON)

	compressedTaskMatchingStatusJSON := suite.compressString(string(taskMatchingStatusJSON))

	resp := map[string]string{
		taskArn1: suite.compressedTask1JSON,
		taskArn2: compressedTaskMatchingStatusJSON,
	}

	suite.datastore.EXPECT().GetWithPrefix(taskKeyPrefix).Return(resp, nil)

	tasks, err := suite.taskStore.FilterTasks(statusFilter, filterStatus)

	assert.Nil(suite.T(), err, "Unexpected error when calling filter tasks")
	assert.Equal(suite.T(), 1, len(tasks), "Expected the length of the FilterTasks result to be 1")
	assert.Exactly(suite.T(), taskMatchingStatus, tasks[0], "Expected one result when one matches filter")
}

func (suite *TaskStoreTestSuite) TestFilterTasksListTasksReturnsMultipleResultsMultipleMatchFilter() {
	resp := map[string]string{
		taskArn1: suite.compressedTask1JSON,
		taskArn2: suite.compressedTask2JSON,
	}

	suite.datastore.EXPECT().GetWithPrefix(taskKeyPrefix).Return(resp, nil)

	tasks, err := suite.taskStore.FilterTasks(statusFilter, suite.task1.Detail.LastStatus)

	assert.Nil(suite.T(), err, "Unexpected error when calling filter tasks")
	assert.Equal(suite.T(), 2, len(tasks), "Expected one result when multiple match filter")

	for _, v := range tasks {
		//attempt to grab the same task from resp
		value, ok := resp[v.Detail.TaskArn]
		if !ok {
			suite.T().Errorf("Expected GetWithPrefix result to contain the same elements as FilterTasks result. Missing %v", v)
		} else {
			uncompressedTaskJSON, err := compress.Uncompress([]byte(value))
			assert.Nil(suite.T(), err, "Unexpected error when uncompresseing task json %v", uncompressedTaskJSON)
			var taskInResp types.Task
			json.Unmarshal([]byte(uncompressedTaskJSON), &taskInResp)
			assert.Exactly(suite.T(), taskInResp, v, "Expected GetWithPrefix result to contain the same elements as FilterTasks result.")
		}
	}
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

	taskResp := addTaskToDSChanAndReadFromTaskRespChan(suite.compressedTask1JSON, dsChan, taskRespChan)

	assert.Nil(suite.T(), taskResp.Err, "Unexpected error when reading task from channel")
	assert.Equal(suite.T(), suite.task1, taskResp.Task, "Expected task in task response to match that in the stream")
}

func (suite *TaskStoreTestSuite) TestStreamTasksInvalidJSONInDSChannel() {
	ctx := context.Background()
	dsChan := make(chan map[string]string)
	defer close(dsChan)
	suite.datastore.EXPECT().StreamWithPrefix(gomock.Any(), taskKeyPrefix).Return(dsChan, nil)

	taskRespChan, err := suite.taskStore.StreamTasks(ctx)
	assert.Nil(suite.T(), err, "Unexpected error when calling stream tasks")
	assert.NotNil(suite.T(), taskRespChan)

	compressedInvalidJSON := suite.compressString("invalidJSON")
	taskResp := addTaskToDSChanAndReadFromTaskRespChan(compressedInvalidJSON, dsChan, taskRespChan)

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

func addTaskToDSChanAndReadFromTaskRespChan(taskToAdd string, dsChan chan map[string]string, taskRespChan chan storetypes.TaskErrorWrapper) storetypes.TaskErrorWrapper {
	var taskResp storetypes.TaskErrorWrapper

	doneChan := make(chan bool)
	defer close(doneChan)
	go func() {
		taskResp = <-taskRespChan
		doneChan <- true
	}()

	dsVal := map[string]string{taskArn1: taskToAdd}
	dsChan <- dsVal
	<-doneChan

	return taskResp
}

func (suite *TaskStoreTestSuite) compressString(str string) string {
	compressedVal, err := compress.Compress(str)
	assert.Nil(suite.T(), err, "Error when compressing string %v", str)
	return string(compressedVal)
}
