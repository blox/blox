package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/amazon-ecs-event-stream-handler/handler/api/v1/models"
	"github.com/aws/amazon-ecs-event-stream-handler/handler/mocks"
	"github.com/aws/amazon-ecs-event-stream-handler/handler/types"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	getTaskPath     = `/task/{arn:(arn:aws:ecs:)([\-\w]+):[0-9]{12}:(task)\/[\-\w]+}`
	listTasksPath   = "/tasks"
	filterTasksPath = "/tasks/filter"

	getTaskPrefix             = "/v1/task"
	listTasksPrefix           = "/v1/tasks"
	filterTasksByStatusPrefix = "/v1/tasks/filter?status="

	filterTasksByStatusQueryValue = "{status:pending|running|stopped}"

	taskStatusKey = "status"

	// Routing to GetInstance handler function without arn
	invalidGetTaskPath       = "/task"
	invalidFilterTasksPrefix = "/v1/tasks/filter"
)

type TaskAPIsTestSuite struct {
	suite.Suite
	taskStore      *mocks.MockTaskStore
	taskAPIs       TaskAPIs
	task1          types.Task
	task2          types.Task
	taskModel1     models.TaskModel
	taskModel2     models.TaskModel
	responseHeader http.Header

	// We need a router because some of the apis use mux.Vars() which uses the URL
	// parameters parsed and stored in a global map in the global context by the router.
	router *mux.Router
}

func (suite *TaskAPIsTestSuite) SetupTest() {
	mockCtrl := gomock.NewController(suite.T())

	suite.taskStore = mocks.NewMockTaskStore(mockCtrl)

	suite.taskAPIs = NewTaskAPIs(suite.taskStore)

	version := 1
	overrides := types.Overrides{
		ContainerOverrides: []*types.ContainerOverrides{},
	}
	taskDetail1 := types.TaskDetail{
		ClusterArn:           &clusterARN1,
		ContainerInstanceArn: &instanceARN1,
		Containers:           []*types.Container{},
		CreatedAt:            &createdAt,
		DesiredStatus:        &taskStatus1,
		LastStatus:           &taskStatus1,
		Overrides:            &overrides,
		TaskArn:              &taskARN1,
		TaskDefinitionArn:    &taskDefinitionARN,
		UpdatedAt:            &updatedAt1,
		Version:              &version,
	}
	suite.task1 = types.Task{
		Account:   &accountID,
		Detail:    &taskDetail1,
		ID:        &id1,
		Region:    &region,
		Resources: []string{taskARN1},
		Time:      &time,
	}

	taskModel, err := ToTaskModel(suite.task1)
	if err != nil {
		suite.T().Error("Cannot setup testSuite: Error when tranlating task to external model")
	}
	suite.taskModel1 = taskModel

	taskDetail2 := taskDetail1
	taskDetail2.TaskArn = &taskARN2
	taskDetail2.LastStatus = &taskStatus2

	suite.task2 = types.Task{
		Account:   &accountID,
		Detail:    &taskDetail2,
		ID:        &id1,
		Region:    &region,
		Resources: []string{taskARN2},
		Time:      &time,
	}

	suite.responseHeader = http.Header{responseContentTypeKey: []string{responseContentTypeVal}}

	suite.router = suite.getRouter()
}

// TODO - Add the following test cases once implementation is in place
// * arn validation fails on getTask
// * streaming api

func TestTaskAPIsTestSuite(t *testing.T) {
	suite.Run(t, new(TaskAPIsTestSuite))
}

func (suite *TaskAPIsTestSuite) TestGetTaskReturnsTask() {
	suite.taskStore.EXPECT().GetTask(taskARN1).Return(&suite.task1, nil)

	request := suite.getTaskRequest()
	responseRecorder := httptest.NewRecorder()

	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateSuccessfulResponseHeaderAndStatus(responseRecorder)

	reader := bytes.NewReader(responseRecorder.Body.Bytes())
	taskInResponse := types.Task{}
	err := json.NewDecoder(reader).Decode(&taskInResponse)
	assert.Nil(suite.T(), err, "Unexpected error decoding response body")
	assert.Exactly(suite.T(), suite.task1, taskInResponse, "Task in response is invalid")
}

func (suite *TaskAPIsTestSuite) TestGetTaskNoTask() {
	suite.taskStore.EXPECT().GetTask(taskARN1).Return(nil, nil)

	request := suite.getTaskRequest()
	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	suite.decodeErrorResponseAndValidate(responseRecorder, http.StatusNotFound, instanceNotFoundClientErrMsg)
}

func (suite *TaskAPIsTestSuite) TestGetTaskStoreReturnsError() {
	suite.taskStore.EXPECT().GetTask(taskARN1).Return(nil, errors.New("Error when getting task"))

	request := suite.getTaskRequest()
	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	suite.decodeErrorResponseAndValidate(responseRecorder, http.StatusInternalServerError, internalServerErrMsg)
}

func (suite *TaskAPIsTestSuite) TestGetTaskWithoutArn() {
	suite.taskStore.EXPECT().GetTask(gomock.Any()).Times(0)

	request, err := http.NewRequest("GET", getTaskPrefix, nil)
	assert.Nil(suite.T(), err, "Unexpected error creating task get request")

	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	suite.decodeErrorResponseAndValidate(responseRecorder, http.StatusInternalServerError, routingServerErrMsg)
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

func (suite *TaskAPIsTestSuite) TestListTasksStoreReturnsError() {
	suite.taskStore.EXPECT().ListTasks().Return(nil, errors.New("Error when listing tasks"))

	request := suite.listTasksRequest()
	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	suite.decodeErrorResponseAndValidate(responseRecorder, http.StatusInternalServerError, internalServerErrMsg)
}

func (suite *TaskAPIsTestSuite) TestFilterTasksReturnsTasks() {
	taskList := []types.Task{suite.task1}
	suite.taskStore.EXPECT().FilterTasks(taskStatusKey, taskStatus1).Return(taskList, nil)
	suite.filterTasksTester(taskList)
}

func (suite *TaskAPIsTestSuite) TestFilterTasksNoTasks() {
	emptytaskList := make([]types.Task, 0)
	suite.taskStore.EXPECT().FilterTasks(taskStatusKey, taskStatus1).Return(emptytaskList, nil)
	suite.filterTasksTester(emptytaskList)
}

func (suite *TaskAPIsTestSuite) TestFilterTasksStoreReturnsError() {
	suite.taskStore.EXPECT().FilterTasks(taskStatusKey, taskStatus1).Return(nil, errors.New("Error when filtering tasks"))

	request := suite.filterTasksRequest()
	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	suite.decodeErrorResponseAndValidate(responseRecorder, http.StatusInternalServerError, internalServerErrMsg)
}

func (suite *TaskAPIsTestSuite) TestFilterTasksNoKey() {
	suite.taskStore.EXPECT().FilterTasks(taskStatusKey, gomock.Any()).Times(0)

	request, err := http.NewRequest("GET", invalidFilterTasksPrefix, nil)
	assert.Nil(suite.T(), err, "Unexpected error creating filter tasks request")

	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	suite.decodeErrorResponseAndValidate(responseRecorder, http.StatusInternalServerError, routingServerErrMsg)
}

// Helper functions

func (suite *TaskAPIsTestSuite) getTaskRequest() *http.Request {
	url := getTaskPrefix + "/" + taskARN1
	request, err := http.NewRequest("GET", url, nil)
	assert.Nil(suite.T(), err, "Unexpected error creating get task request")
	return request
}

func (suite *TaskAPIsTestSuite) listTasksRequest() *http.Request {
	request, err := http.NewRequest("GET", listTasksPrefix, nil)
	assert.Nil(suite.T(), err, "Unexpected error creating list tasks request")
	return request
}

func (suite *TaskAPIsTestSuite) filterTasksRequest() *http.Request {
	url := filterTasksByStatusPrefix + taskStatus1
	request, err := http.NewRequest("GET", url, nil)
	assert.Nil(suite.T(), err, "Unexpected error creating task get request")
	return request
}

func (suite *TaskAPIsTestSuite) listTaskTester(taskList []types.Task) {
	request := suite.listTasksRequest()
	responseRecorder := httptest.NewRecorder()

	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateSuccessfulResponseHeaderAndStatus(responseRecorder)
	suite.validateTasksInListOrFilterTasksResponse(responseRecorder, taskList)
}

func (suite *TaskAPIsTestSuite) filterTasksTester(taskList []types.Task) {
	request := suite.filterTasksRequest()
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

	s.Path(filterTasksPath).
		Queries(taskStatusKey, filterTasksByStatusQueryValue).
		Methods("GET").
		HandlerFunc(suite.taskAPIs.FilterTasks)

	// Invalid router paths to make sure handler functions handle them
	s.Path(invalidGetTaskPath).
		Methods("GET").
		HandlerFunc(suite.taskAPIs.GetTask)

	s.Path(filterTasksPath).
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

func (suite *TaskAPIsTestSuite) validateTasksInListOrFilterTasksResponse(responseRecorder *httptest.ResponseRecorder, expectedTasks []types.Task) {
	reader := bytes.NewReader(responseRecorder.Body.Bytes())
	tasksInResponse := new([]types.Task)
	err := json.NewDecoder(reader).Decode(tasksInResponse)
	assert.Nil(suite.T(), err, "Unexpected error decoding response body")
	assert.Exactly(suite.T(), expectedTasks, *tasksInResponse, "Tasks in response is invalid")
}

func (suite *TaskAPIsTestSuite) decodeErrorResponseAndValidate(responseRecorder *httptest.ResponseRecorder, expectedErrCode int, expectedErrMsg string) {
	reader := bytes.NewReader(responseRecorder.Body.Bytes())
	errorModel := models.ErrorModel{}
	err := json.NewDecoder(reader).Decode(&errorModel)
	assert.Nil(suite.T(), err, "Unexpected error decoding response body")
	assert.Equal(suite.T(), int32(expectedErrCode), *errorModel.Code)
	assert.Equal(suite.T(), expectedErrMsg, *errorModel.Message)
}
