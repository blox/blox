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
	getTaskPrefix              = "/v1/tasks"
	listTasksPrefix            = "/v1/tasks"
	filterTasksByStatusPrefix  = "/v1/tasks/filter?status="
	filterTasksByClusterPrefix = "/v1/tasks/filter?cluster="

	filterTasksByStatusQueryValue = "{status:pending|running|stopped}"

	// Routing to GetInstance handler function without task ARN
	invalidGetTaskPath       = "/tasks/{cluster:[a-zA-Z0-9_]{1,255}}"
	invalidFilterTasksPrefix = "/v1/tasks/filter"
)

type TaskAPIsTestSuite struct {
	suite.Suite
	taskStore      *mocks.MockTaskStore
	taskAPIs       TaskAPIs
	task1          types.Task
	task2          types.Task
	extTask1       models.Task
	extTask2       models.Task
	responseHeader http.Header

	// We need a router because some of the apis use mux.Vars() which uses the URL
	// parameters parsed and stored in a global map in the global context by the router.
	router *mux.Router
}

func (suite *TaskAPIsTestSuite) SetupTest() {
	mockCtrl := gomock.NewController(suite.T())

	suite.taskStore = mocks.NewMockTaskStore(mockCtrl)

	suite.taskAPIs = NewTaskAPIs(suite.taskStore)

	overrides := types.Overrides{
		ContainerOverrides: []*types.ContainerOverrides{},
	}
	taskDetail1 := types.TaskDetail{
		ClusterARN:           &clusterARN1,
		ContainerInstanceARN: &instanceARN1,
		Containers:           []*types.Container{},
		CreatedAt:            &createdAt,
		DesiredStatus:        &taskStatus1,
		LastStatus:           &taskStatus1,
		Overrides:            &overrides,
		TaskARN:              &taskARN1,
		TaskDefinitionARN:    &taskDefinitionARN,
		UpdatedAt:            &updatedAt1,
		Version:              &version1,
	}
	suite.task1 = types.Task{
		Account:   &accountID,
		Detail:    &taskDetail1,
		ID:        &id1,
		Region:    &region,
		Resources: []string{taskARN1},
		Time:      &time,
	}

	extTask, err := ToTask(suite.task1)
	if err != nil {
		suite.T().Error("Cannot setup testSuite: Error when tranlating task to external task")
	}
	suite.extTask1 = extTask

	taskDetail2 := taskDetail1
	taskDetail2.TaskARN = &taskARN2
	taskDetail2.LastStatus = &taskStatus2

	suite.task2 = types.Task{
		Account:   &accountID,
		Detail:    &taskDetail2,
		ID:        &id1,
		Region:    &region,
		Resources: []string{taskARN2},
		Time:      &time,
	}

	extTask, err = ToTask(suite.task2)
	if err != nil {
		suite.T().Error("Cannot setup testSuite: Error when tranlating task to external model")
	}
	suite.extTask2 = extTask

	suite.responseHeader = http.Header{responseContentTypeKey: []string{responseContentTypeVal}}

	suite.router = suite.getRouter()
}

// TODO - Add the following test cases once implementation is in place
// * arn validation fails on getTask
// * streaming api

func TestTaskAPIsTestSuite(t *testing.T) {
	suite.Run(t, new(TaskAPIsTestSuite))
}

func (suite *TaskAPIsTestSuite) TestGetReturnsTask() {
	suite.taskStore.EXPECT().GetTask(clusterName1, taskARN1).Return(&suite.task1, nil)

	request := suite.getTaskRequest()
	responseRecorder := httptest.NewRecorder()

	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateSuccessfulResponseHeaderAndStatus(responseRecorder)

	reader := bytes.NewReader(responseRecorder.Body.Bytes())
	taskInResponse := models.Task{}
	err := json.NewDecoder(reader).Decode(&taskInResponse)
	assert.Nil(suite.T(), err, "Unexpected error decoding response body")
	assert.Exactly(suite.T(), suite.extTask1, taskInResponse, "Task in response is invalid")
}

func (suite *TaskAPIsTestSuite) TestGetTaskNoTask() {
	suite.taskStore.EXPECT().GetTask(clusterName1, taskARN1).Return(nil, nil)

	request := suite.getTaskRequest()
	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateErrorResponseHeaderAndStatus(responseRecorder, http.StatusNotFound)
	suite.decodeErrorResponseAndValidate(responseRecorder, taskNotFoundClientErrMsg)
}

func (suite *TaskAPIsTestSuite) TestGetTaskStoreReturnsError() {
	suite.taskStore.EXPECT().GetTask(clusterName1, taskARN1).Return(nil, errors.New("Error when getting task"))

	request := suite.getTaskRequest()
	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateErrorResponseHeaderAndStatus(responseRecorder, http.StatusInternalServerError)
	suite.decodeErrorResponseAndValidate(responseRecorder, internalServerErrMsg)
}

func (suite *TaskAPIsTestSuite) TestGetTaskWithoutTaskARN() {
	suite.taskStore.EXPECT().GetTask(gomock.Any(), gomock.Any()).Times(0)

	url := getTaskPrefix + "/" + clusterName1
	request, err := http.NewRequest("GET", url, nil)
	assert.Nil(suite.T(), err, "Unexpected error creating task get request")

	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateErrorResponseHeaderAndStatus(responseRecorder, http.StatusInternalServerError)
	suite.decodeErrorResponseAndValidate(responseRecorder, routingServerErrMsg)
}

func (suite *TaskAPIsTestSuite) TestListTasksReturnsTasks() {
	taskList := []types.Task{suite.task1, suite.task2}
	suite.taskStore.EXPECT().ListTasks().Return(taskList, nil)
	extTasks := models.Tasks{
		Items: []*models.Task{&suite.extTask1, &suite.extTask2},
	}
	suite.listTaskTester(extTasks)
}

func (suite *TaskAPIsTestSuite) TestListTasksNoTasks() {
	emptyTaskList := make([]types.Task, 0)
	suite.taskStore.EXPECT().ListTasks().Return(emptyTaskList, nil)
	emptyExtTasks := models.Tasks{
		Items: []*models.Task{},
	}
	suite.listTaskTester(emptyExtTasks)
}

func (suite *TaskAPIsTestSuite) TestListTasksStoreReturnsError() {
	suite.taskStore.EXPECT().ListTasks().Return(nil, errors.New("Error when listing tasks"))

	request := suite.listTasksRequest()
	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateErrorResponseHeaderAndStatus(responseRecorder, http.StatusInternalServerError)
	suite.decodeErrorResponseAndValidate(responseRecorder, internalServerErrMsg)
}

func (suite *TaskAPIsTestSuite) TestFilterTasksByStatusReturnsTasks() {
	taskList := []types.Task{suite.task1}
	suite.taskStore.EXPECT().FilterTasks(taskStatusFilter, taskStatus1).Return(taskList, nil)
	extTasks := models.Tasks{
		Items: []*models.Task{&suite.extTask1},
	}
	suite.filterTasksByStatusTester(extTasks)
}

func (suite *TaskAPIsTestSuite) TestFilterTasksByStatusNoTasks() {
	emptyTaskList := make([]types.Task, 0)
	suite.taskStore.EXPECT().FilterTasks(taskStatusFilter, taskStatus1).Return(emptyTaskList, nil)
	emptyExtTasks := models.Tasks{
		Items: []*models.Task{},
	}
	suite.filterTasksByStatusTester(emptyExtTasks)
}

func (suite *TaskAPIsTestSuite) TestFilterTasksByStatusStoreReturnsError() {
	suite.taskStore.EXPECT().FilterTasks(taskStatusFilter, taskStatus1).Return(nil, errors.New("Error when filtering tasks"))

	request := suite.filterTasksByStatusRequest()
	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateErrorResponseHeaderAndStatus(responseRecorder, http.StatusInternalServerError)
	suite.decodeErrorResponseAndValidate(responseRecorder, internalServerErrMsg)
}

func (suite *TaskAPIsTestSuite) TestFilterTasksByClusterNameReturnsTasks() {
	taskList := []types.Task{suite.task1}
	suite.taskStore.EXPECT().FilterTasks(taskClusterFilter, clusterName1).Return(taskList, nil)
	extTasks := models.Tasks{
		Items: []*models.Task{&suite.extTask1},
	}
	suite.filterTasksByClusterTester(clusterName1, extTasks)
}

func (suite *TaskAPIsTestSuite) TestFilterTasksByClusterNameNoTasks() {
	emptyTaskList := make([]types.Task, 0)
	suite.taskStore.EXPECT().FilterTasks(taskClusterFilter, clusterName1).Return(emptyTaskList, nil)
	emptyExtTasks := models.Tasks{
		Items: []*models.Task{},
	}
	suite.filterTasksByClusterTester(clusterName1, emptyExtTasks)
}

func (suite *TaskAPIsTestSuite) TestFilterTasksByClusterNameStoreReturnsError() {
	suite.taskStore.EXPECT().FilterTasks(taskClusterFilter, clusterName1).Return(nil, errors.New("Error when filtering tasks"))

	request := suite.filterTasksByClusterRequest(clusterName1)
	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateErrorResponseHeaderAndStatus(responseRecorder, http.StatusInternalServerError)
	suite.decodeErrorResponseAndValidate(responseRecorder, internalServerErrMsg)
}

func (suite *TaskAPIsTestSuite) TestFilterTasksByClusterARNReturnsTasks() {
	taskList := []types.Task{suite.task1}
	suite.taskStore.EXPECT().FilterTasks(taskClusterFilter, clusterARN1).Return(taskList, nil)
	extTasks := models.Tasks{
		Items: []*models.Task{&suite.extTask1},
	}
	suite.filterTasksByClusterTester(clusterARN1, extTasks)
}

func (suite *TaskAPIsTestSuite) TestFilterTasksByClusterARNNoTasks() {
	emptyTaskList := make([]types.Task, 0)
	suite.taskStore.EXPECT().FilterTasks(taskClusterFilter, clusterARN1).Return(emptyTaskList, nil)
	emptyExtTasks := models.Tasks{
		Items: []*models.Task{},
	}
	suite.filterTasksByClusterTester(clusterARN1, emptyExtTasks)
}

func (suite *TaskAPIsTestSuite) TestFilterTasksByClusterARNStoreReturnsError() {
	suite.taskStore.EXPECT().FilterTasks(taskClusterFilter, clusterARN1).Return(nil, errors.New("Error when filtering tasks"))

	request := suite.filterTasksByClusterRequest(clusterARN1)
	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateErrorResponseHeaderAndStatus(responseRecorder, http.StatusInternalServerError)
	suite.decodeErrorResponseAndValidate(responseRecorder, internalServerErrMsg)
}

func (suite *TaskAPIsTestSuite) TestFilterTasksNoKey() {
	suite.taskStore.EXPECT().FilterTasks(taskStatusFilter, gomock.Any()).Times(0)

	request, err := http.NewRequest("GET", invalidFilterTasksPrefix, nil)
	assert.Nil(suite.T(), err, "Unexpected error creating filter tasks request")

	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateErrorResponseHeaderAndStatus(responseRecorder, http.StatusInternalServerError)
	suite.decodeErrorResponseAndValidate(responseRecorder, routingServerErrMsg)
}

// Helper functions

func (suite *TaskAPIsTestSuite) getTaskRequest() *http.Request {
	url := getTaskPrefix + "/" + clusterName1 + "/" + taskARN1
	request, err := http.NewRequest("GET", url, nil)
	assert.Nil(suite.T(), err, "Unexpected error creating get task request")
	return request
}

func (suite *TaskAPIsTestSuite) listTasksRequest() *http.Request {
	request, err := http.NewRequest("GET", listTasksPrefix, nil)
	assert.Nil(suite.T(), err, "Unexpected error creating list tasks request")
	return request
}

func (suite *TaskAPIsTestSuite) filterTasksByStatusRequest() *http.Request {
	url := filterTasksByStatusPrefix + taskStatus1
	request, err := http.NewRequest("GET", url, nil)
	assert.Nil(suite.T(), err, "Unexpected error creating filter tasks by status request")
	return request
}

func (suite *TaskAPIsTestSuite) filterTasksByClusterRequest(cluster string) *http.Request {
	url := filterTasksByClusterPrefix + cluster
	request, err := http.NewRequest("GET", url, nil)
	assert.Nil(suite.T(), err, "Unexpected error creating filter tasks by cluster request")
	return request
}

func (suite *TaskAPIsTestSuite) listTaskTester(tasks models.Tasks) {
	request := suite.listTasksRequest()
	responseRecorder := httptest.NewRecorder()

	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateSuccessfulResponseHeaderAndStatus(responseRecorder)
	suite.validateTasksInListOrFilterTasksResponse(responseRecorder, tasks)
}

func (suite *TaskAPIsTestSuite) filterTasksByStatusTester(tasks models.Tasks) {
	request := suite.filterTasksByStatusRequest()
	responseRecorder := httptest.NewRecorder()

	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateSuccessfulResponseHeaderAndStatus(responseRecorder)
	suite.validateTasksInListOrFilterTasksResponse(responseRecorder, tasks)
}

func (suite *TaskAPIsTestSuite) filterTasksByClusterTester(cluster string, tasks models.Tasks) {
	request := suite.filterTasksByClusterRequest(cluster)
	responseRecorder := httptest.NewRecorder()

	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateSuccessfulResponseHeaderAndStatus(responseRecorder)
	suite.validateTasksInListOrFilterTasksResponse(responseRecorder, tasks)
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
		Queries(statusKey, taskStatusVal).
		Methods("GET").
		HandlerFunc(suite.taskAPIs.FilterTasks)

	s.Path(filterTasksPath).
		Queries(clusterKey, clusterNameVal).
		Methods("GET").
		HandlerFunc(suite.taskAPIs.FilterTasks)

	s.Path(filterTasksPath).
		Queries(clusterKey, clusterARNVal).
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

func (suite *TaskAPIsTestSuite) validateErrorResponseHeaderAndStatus(responseRecorder *httptest.ResponseRecorder, errorCode int) {
	h := responseRecorder.Header()
	assert.NotNil(suite.T(), h, "Unexpected empty header")
	assert.Equal(suite.T(), errorCode, responseRecorder.Code, "Http response status is invalid")
}

func (suite *TaskAPIsTestSuite) validateTasksInListOrFilterTasksResponse(responseRecorder *httptest.ResponseRecorder, expectedTasks models.Tasks) {
	reader := bytes.NewReader(responseRecorder.Body.Bytes())
	tasksInResponse := new(models.Tasks)
	err := json.NewDecoder(reader).Decode(tasksInResponse)
	assert.Nil(suite.T(), err, "Unexpected error decoding response body")
	assert.Exactly(suite.T(), expectedTasks, *tasksInResponse, "Tasks in response is invalid")
}

func (suite *TaskAPIsTestSuite) decodeErrorResponseAndValidate(responseRecorder *httptest.ResponseRecorder, expectedErrMsg string) {
	reader := bytes.NewReader(responseRecorder.Body.Bytes())
	var str string
	err := json.NewDecoder(reader).Decode(&str)
	assert.Nil(suite.T(), err, "Unexpected error decoding response body")
	assert.Equal(suite.T(), expectedErrMsg, str)
}
