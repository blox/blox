package v1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/amazon-ecs-event-stream-handler/handler/mocks"
	"github.com/aws/amazon-ecs-event-stream-handler/handler/types"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	getTaskPath    = `/task/{arn:(arn:aws:ecs:)([\-\w]+):[0-9]{12}:(task)\/[\-\w]+}`
	listTasksPath  = "/tasks"
	filterTaskPath = "/tasks/filter"

	getTaskPrefix             = "/v1/task"
	listTasksPrefix           = "/v1/tasks"
	filterTasksByStatusPrefix = "/v1/tasks/filter?status="

	filterTasksByStatusQueryValue = "{status:pending|running|stopped}"

	statusKey = "status"

	taskARN1 = "arn:aws:ecs:us-east-1:123456789012:task/271022c0-f894-4aa2-b063-25bae55088d5"
	taskARN2 = "arn:aws:ecs:us-east-1:123456789012:task/b6b9eace-958e-4f2a-a09c-8cf43b76cf97"

	taskStatus1 = "pending"
	taskStatus2 = "stopped"

	responseContentTypeKey = "Content-Type"
	responseContentTypeVal = "application/json; charset=UTF-8"
)

type TaskAPIsTestSuite struct {
	suite.Suite
	taskStore      *mocks.MockTaskStore
	taskAPIs       TaskAPIs
	task1          types.Task
	task2          types.Task
	responseHeader http.Header

	// We need a router because some of the apis use mux.Vars() which uses the URL
	// parameters parsed and stored in a global map in the global context by the router.
	router *mux.Router
}

func (suite *TaskAPIsTestSuite) SetupTest() {
	mockCtrl := gomock.NewController(suite.T())

	suite.taskStore = mocks.NewMockTaskStore(mockCtrl)

	suite.taskAPIs = NewTaskAPIs(suite.taskStore)

	suite.task1 = types.Task{}
	suite.task1.Detail.TaskArn = taskARN1
	suite.task1.Detail.Version = 1
	suite.task1.Detail.LastStatus = taskStatus1

	suite.task2 = types.Task{}
	suite.task2.Detail.TaskArn = taskARN2
	suite.task2.Detail.Version = 1
	suite.task2.Detail.LastStatus = taskStatus2

	suite.responseHeader = http.Header{responseContentTypeKey: []string{responseContentTypeVal}}

	suite.router = suite.getRouter()
}

// TODO - Add the following test cases once implementation is in place
// * taskStore returns an error on getTask
// * arn validation fails on getTask
// * taskStore returns an error on listTasks
// * taskStore returns an error on filterTasks
// * streaming api

func TestTaskAPIsTestSuite(t *testing.T) {
	suite.Run(t, new(TaskAPIsTestSuite))
}

func (suite *TaskAPIsTestSuite) TestGetTaskReturnsTask() {
	suite.taskStore.EXPECT().GetTask(taskARN1).Return(&suite.task1, nil)
	suite.getTaskTester(suite.task1)
}

func (suite *TaskAPIsTestSuite) TestGetTaskNoTask() {
	suite.taskStore.EXPECT().GetTask(taskARN1).Return(nil, nil)
	emptyTask := types.Task{}
	suite.getTaskTester(emptyTask)
}

func (suite *TaskAPIsTestSuite) TestListTasksReturnsTasks() {
	taskList := []types.Task{suite.task1, suite.task2}
	suite.taskStore.EXPECT().ListTasks().Return(taskList, nil)
	suite.listTaskTester(taskList)
}

func (suite *TaskAPIsTestSuite) TestListTasksNoTasks() {
	emptytaskList := make([]types.Task, 0)
	suite.taskStore.EXPECT().ListTasks().Return(emptytaskList, nil)
	suite.listTaskTester(emptytaskList)
}

func (suite *TaskAPIsTestSuite) TestFilterTasksReturnsTasks() {
	taskList := []types.Task{suite.task1}
	suite.taskStore.EXPECT().FilterTasks(statusKey, taskStatus1).Return(taskList, nil)
	suite.filterTasksTester(taskList)
}

func (suite *TaskAPIsTestSuite) TestFilterTasksNoTasks() {
	emptytaskList := make([]types.Task, 0)
	suite.taskStore.EXPECT().FilterTasks(statusKey, taskStatus1).Return(emptytaskList, nil)
	suite.filterTasksTester(emptytaskList)
}

// Helper functions

func (suite *TaskAPIsTestSuite) getTaskTester(task types.Task) {
	url := getTaskPrefix + "/" + taskARN1
	request, err := http.NewRequest("GET", url, nil)
	assert.Nil(suite.T(), err, "Unexpected error creating task get request")

	responseRecorder := httptest.NewRecorder()

	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateSuccessfulResponseHeaderAndStatus(responseRecorder)
	suite.validateTaskInGetTaskResponse(responseRecorder, task)
}

func (suite *TaskAPIsTestSuite) listTaskTester(taskList []types.Task) {
	request, err := http.NewRequest("GET", listTasksPrefix, nil)
	assert.Nil(suite.T(), err, "Unexpected error creating task get request")

	responseRecorder := httptest.NewRecorder()

	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateSuccessfulResponseHeaderAndStatus(responseRecorder)
	suite.validateTasksInListOrFilterTasksResponse(responseRecorder, taskList)
}

func (suite *TaskAPIsTestSuite) filterTasksTester(taskList []types.Task) {

	url := filterTasksByStatusPrefix + taskStatus1
	request, err := http.NewRequest("GET", url, nil)
	assert.Nil(suite.T(), err, "Unexpected error creating task get request")

	responseRecorder := httptest.NewRecorder()

	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateSuccessfulResponseHeaderAndStatus(responseRecorder)
	suite.validateTasksInListOrFilterTasksResponse(responseRecorder, taskList)
}

func (suite *TaskAPIsTestSuite) getRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	s := r.Path("/v1").Subrouter()

	s.Path(getTaskPath).
		Methods("GET").
		HandlerFunc(suite.taskAPIs.GetTask)

	s.Path(listTasksPath).
		Methods("GET").
		HandlerFunc(suite.taskAPIs.ListTasks)

	s.Path(filterTaskPath).
		Queries(statusKey, filterTasksByStatusQueryValue).
		Methods("GET").
		HandlerFunc(suite.taskAPIs.FilterTasks)

	return s
}

func (suite *TaskAPIsTestSuite) validateSuccessfulResponseHeaderAndStatus(responseRecorder *httptest.ResponseRecorder) {
	h := responseRecorder.Header()
	assert.NotNil(suite.T(), h, "Unexpected empty header")
	assert.Equal(suite.T(), suite.responseHeader, h, "Http header is invalid")
	assert.Equal(suite.T(), http.StatusOK, responseRecorder.Code, "Http response status is invalid")
}

func (suite *TaskAPIsTestSuite) validateTaskInGetTaskResponse(responseRecorder *httptest.ResponseRecorder, expectedTask types.Task) {
	reader := bytes.NewReader(responseRecorder.Body.Bytes())
	taskInResponse := types.Task{}
	err := json.NewDecoder(reader).Decode(&taskInResponse)
	assert.Nil(suite.T(), err, "Unexpected error decoding response body")
	assert.Exactly(suite.T(), expectedTask, taskInResponse, "Task in response is invalid")
}

func (suite *TaskAPIsTestSuite) validateTasksInListOrFilterTasksResponse(responseRecorder *httptest.ResponseRecorder, expectedTasks []types.Task) {
	reader := bytes.NewReader(responseRecorder.Body.Bytes())
	tasksInResponse := new([]types.Task)
	err := json.NewDecoder(reader).Decode(tasksInResponse)
	assert.Nil(suite.T(), err, "Unexpected error decoding response body")
	assert.Exactly(suite.T(), expectedTasks, *tasksInResponse, "Tasks in response is invalid")
}
