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

package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/blox/blox/cluster-state-service/handler/api/v1/models"
	"github.com/blox/blox/cluster-state-service/handler/mocks"
	"github.com/blox/blox/cluster-state-service/handler/types"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	getTaskPrefix              = "/v1/tasks"
	listTasksPrefix            = "/v1/tasks"
	filterTasksByStatusPrefix  = "/v1/tasks?status="
	filterTasksByClusterPrefix = "/v1/tasks?cluster="

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
	suite.taskStore.EXPECT().FilterTasks(gomock.Any()).Times(0)
	extTasks := models.Tasks{
		Items: []*models.Task{&suite.extTask1, &suite.extTask2},
	}
	suite.listTasksTester(extTasks)
}

func (suite *TaskAPIsTestSuite) TestListTasksNoTasks() {
	emptyTaskList := make([]types.Task, 0)
	suite.taskStore.EXPECT().ListTasks().Return(emptyTaskList, nil)
	suite.taskStore.EXPECT().FilterTasks(gomock.Any()).Times(0)
	emptyExtTasks := models.Tasks{
		Items: []*models.Task{},
	}
	suite.listTasksTester(emptyExtTasks)
}

func (suite *TaskAPIsTestSuite) TestListTasksStoreReturnsError() {
	suite.taskStore.EXPECT().ListTasks().Return(nil, errors.New("Error when listing tasks"))
	suite.taskStore.EXPECT().FilterTasks(gomock.Any()).Times(0)

	request := suite.listTasksRequest()
	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateErrorResponseHeaderAndStatus(responseRecorder, http.StatusInternalServerError)
	suite.decodeErrorResponseAndValidate(responseRecorder, internalServerErrMsg)
}

func (suite *TaskAPIsTestSuite) TestListTasksInvalidFilter() {
	suite.taskStore.EXPECT().FilterTasks(gomock.Any()).Times(0)
	suite.taskStore.EXPECT().ListTasks().Times(0)

	url := listTasksPrefix + "?unsupportedFilter=val"
	request, err := http.NewRequest("GET", url, nil)
	assert.Nil(suite.T(), err, "Unexpected error creating list tasks request with invalid filter")

	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateErrorResponseHeaderAndStatus(responseRecorder, http.StatusBadRequest)
	suite.decodeErrorResponseAndValidate(responseRecorder, unsupportedFilterClientErrMsg)
}

func (suite *TaskAPIsTestSuite) TestListTasksBothStatusAndClusterFilter() {
	taskList := []types.Task{suite.task1}

	suite.taskStore.EXPECT().FilterTasks(map[string]string{taskStatusFilter: taskStatus1, taskClusterFilter: clusterARN1}).Return(taskList, nil)
	suite.taskStore.EXPECT().ListTasks().Times(0)

	url := listTasksPrefix + "?status=" + taskStatus1 + "&cluster=" + clusterARN1
	request, err := http.NewRequest("GET", url, nil)
	assert.Nil(suite.T(), err, "Unexpected error creating list tasks request with both status and cluster filter")

	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateSuccessfulResponseHeaderAndStatus(responseRecorder)
	extTasks := models.Tasks{
		Items: []*models.Task{&suite.extTask1},
	}
	suite.validateTasksInListTasksResponse(responseRecorder, extTasks)
}

func (suite *TaskAPIsTestSuite) TestListTasksWithStatusFilterReturnsTasks() {
	taskList := []types.Task{suite.task1}

	suite.taskStore.EXPECT().FilterTasks(map[string]string{taskStatusFilter: taskStatus1}).Return(taskList, nil)
	suite.taskStore.EXPECT().ListTasks().Times(0)

	extTasks := models.Tasks{
		Items: []*models.Task{&suite.extTask1},
	}
	suite.filterTasksByStatusTester(extTasks, taskStatus1)
}

func (suite *TaskAPIsTestSuite) TestListTasksWithCapitalizedStatusFilterReturnsTasks() {
	taskList := []types.Task{suite.task1}

	suite.taskStore.EXPECT().FilterTasks(map[string]string{taskStatusFilter: taskStatus1}).Return(taskList, nil)
	suite.taskStore.EXPECT().ListTasks().Times(0)

	extTasks := models.Tasks{
		Items: []*models.Task{&suite.extTask1},
	}
	suite.filterTasksByStatusTester(extTasks, strings.ToUpper(taskStatus1))
}

func (suite *TaskAPIsTestSuite) TestListTasksWithStatusFilterNoTasks() {
	emptyTaskList := make([]types.Task, 0)

	suite.taskStore.EXPECT().FilterTasks(map[string]string{taskStatusFilter: taskStatus1}).Return(emptyTaskList, nil)
	suite.taskStore.EXPECT().ListTasks().Times(0)

	emptyExtTasks := models.Tasks{
		Items: []*models.Task{},
	}
	suite.filterTasksByStatusTester(emptyExtTasks, taskStatus1)
}

func (suite *TaskAPIsTestSuite) TestListTasksWithStatusFilterStoreReturnsError() {
	suite.taskStore.EXPECT().FilterTasks(map[string]string{taskStatusFilter: taskStatus1}).Return(nil, errors.New("Error when filtering tasks"))
	suite.taskStore.EXPECT().ListTasks().Times(0)

	request := suite.filterTasksByStatusRequest(taskStatus1)
	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateErrorResponseHeaderAndStatus(responseRecorder, http.StatusInternalServerError)
	suite.decodeErrorResponseAndValidate(responseRecorder, internalServerErrMsg)
}

func (suite *TaskAPIsTestSuite) TestListTasksWithInvalidStatusFilter() {
	suite.taskStore.EXPECT().FilterTasks(gomock.Any()).Times(0)
	suite.taskStore.EXPECT().ListTasks().Times(0)

	url := filterTasksByStatusPrefix + "invalidStatus"
	request, err := http.NewRequest("GET", url, nil)
	assert.Nil(suite.T(), err, "Unexpected error creating list tasks request with invalid status")

	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateErrorResponseHeaderAndStatus(responseRecorder, http.StatusBadRequest)
	suite.decodeErrorResponseAndValidate(responseRecorder, invalidStatusClientErrMsg)
}

func (suite *TaskAPIsTestSuite) TestListTasksWithClusterNameFilterReturnsTasks() {
	taskList := []types.Task{suite.task1}

	suite.taskStore.EXPECT().FilterTasks(map[string]string{taskClusterFilter: clusterName1}).Return(taskList, nil)
	suite.taskStore.EXPECT().ListTasks().Times(0)

	extTasks := models.Tasks{
		Items: []*models.Task{&suite.extTask1},
	}
	suite.filterTasksByClusterTester(clusterName1, extTasks)
}

func (suite *TaskAPIsTestSuite) TestListTasksWithClusterNameFilterNoTasks() {
	emptyTaskList := make([]types.Task, 0)

	suite.taskStore.EXPECT().FilterTasks(map[string]string{taskClusterFilter: clusterName1}).Return(emptyTaskList, nil)
	suite.taskStore.EXPECT().ListTasks().Times(0)

	emptyExtTasks := models.Tasks{
		Items: []*models.Task{},
	}
	suite.filterTasksByClusterTester(clusterName1, emptyExtTasks)
}

func (suite *TaskAPIsTestSuite) TestListTasksWithClusterNameFilterStoreReturnsError() {
	suite.taskStore.EXPECT().FilterTasks(map[string]string{taskClusterFilter: clusterName1}).Return(nil, errors.New("Error when filtering tasks"))
	suite.taskStore.EXPECT().ListTasks().Times(0)

	request := suite.filterTasksByClusterRequest(clusterName1)
	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateErrorResponseHeaderAndStatus(responseRecorder, http.StatusInternalServerError)
	suite.decodeErrorResponseAndValidate(responseRecorder, internalServerErrMsg)
}

func (suite *TaskAPIsTestSuite) TestListTasksWithClusterARNFilterReturnsTasks() {
	taskList := []types.Task{suite.task1}

	suite.taskStore.EXPECT().FilterTasks(map[string]string{taskClusterFilter: clusterARN1}).Return(taskList, nil)
	suite.taskStore.EXPECT().ListTasks().Times(0)

	extTasks := models.Tasks{
		Items: []*models.Task{&suite.extTask1},
	}
	suite.filterTasksByClusterTester(clusterARN1, extTasks)
}

func (suite *TaskAPIsTestSuite) TestListTasksWithClusterARNFilterNoTasks() {
	emptyTaskList := make([]types.Task, 0)

	suite.taskStore.EXPECT().FilterTasks(map[string]string{taskClusterFilter: clusterARN1}).Return(emptyTaskList, nil)
	suite.taskStore.EXPECT().ListTasks().Times(0)

	emptyExtTasks := models.Tasks{
		Items: []*models.Task{},
	}
	suite.filterTasksByClusterTester(clusterARN1, emptyExtTasks)
}

func (suite *TaskAPIsTestSuite) TestListTasksWithClusterARNFilterStoreReturnsError() {
	suite.taskStore.EXPECT().FilterTasks(map[string]string{taskClusterFilter: clusterARN1}).Return(nil, errors.New("Error when filtering tasks"))
	suite.taskStore.EXPECT().ListTasks().Times(0)

	request := suite.filterTasksByClusterRequest(clusterARN1)
	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateErrorResponseHeaderAndStatus(responseRecorder, http.StatusInternalServerError)
	suite.decodeErrorResponseAndValidate(responseRecorder, internalServerErrMsg)
}

func (suite *TaskAPIsTestSuite) TestListTasksWithInvalidClusterFilter() {
	url := filterTasksByClusterPrefix + "cluster/cluster"
	request, err := http.NewRequest("GET", url, nil)
	assert.Nil(suite.T(), err, "Unexpected error creating list tasks request with invalid cluster")

	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateErrorResponseHeaderAndStatus(responseRecorder, http.StatusBadRequest)
	suite.decodeErrorResponseAndValidate(responseRecorder, invalidClusterClientErrMsg)
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

func (suite *TaskAPIsTestSuite) filterTasksByStatusRequest(status string) *http.Request {
	url := filterTasksByStatusPrefix + status
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

func (suite *TaskAPIsTestSuite) listTasksTester(tasks models.Tasks) {
	request := suite.listTasksRequest()
	responseRecorder := httptest.NewRecorder()

	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateSuccessfulResponseHeaderAndStatus(responseRecorder)
	suite.validateTasksInListTasksResponse(responseRecorder, tasks)
}

func (suite *TaskAPIsTestSuite) filterTasksByStatusTester(tasks models.Tasks, status string) {
	request := suite.filterTasksByStatusRequest(status)
	responseRecorder := httptest.NewRecorder()

	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateSuccessfulResponseHeaderAndStatus(responseRecorder)
	suite.validateTasksInListTasksResponse(responseRecorder, tasks)
}

func (suite *TaskAPIsTestSuite) filterTasksByClusterTester(cluster string, tasks models.Tasks) {
	request := suite.filterTasksByClusterRequest(cluster)
	responseRecorder := httptest.NewRecorder()

	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateSuccessfulResponseHeaderAndStatus(responseRecorder)
	suite.validateTasksInListTasksResponse(responseRecorder, tasks)
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

	// Invalid router paths to make sure handler functions handle them
	s.Path(invalidGetTaskPath).
		Methods("GET").
		HandlerFunc(suite.taskAPIs.GetTask)

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

func (suite *TaskAPIsTestSuite) validateTasksInListTasksResponse(responseRecorder *httptest.ResponseRecorder, expectedTasks models.Tasks) {
	reader := bytes.NewReader(responseRecorder.Body.Bytes())
	tasksInResponse := new(models.Tasks)
	err := json.NewDecoder(reader).Decode(tasksInResponse)
	assert.Nil(suite.T(), err, "Unexpected error decoding response body")
	assert.Exactly(suite.T(), expectedTasks, *tasksInResponse, "Tasks in response is invalid")
}

func (suite *TaskAPIsTestSuite) decodeErrorResponseAndValidate(responseRecorder *httptest.ResponseRecorder, expectedErrMsg string) {
	actualMsg := responseRecorder.Body.String()
	assert.Equal(suite.T(), expectedErrMsg+"\n", actualMsg, "Error message is invalid")
}
