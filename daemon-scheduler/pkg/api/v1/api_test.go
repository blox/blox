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
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/blox/blox/daemon-scheduler/generated/v1/models"
	"github.com/blox/blox/daemon-scheduler/pkg/mocks"
	"github.com/blox/blox/daemon-scheduler/pkg/types"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	clusterARN1       = "arn:aws:ecs:us-east-1:123456789123:cluster/test1"
	clusterARN2       = "arn:aws:ecs:us-east-1:123456789123:cluster/test2"
	taskDefinitionARN = "arn:aws:ecs:us-east-1:12345678912:task-definition/test"
)

type APITestSuite struct {
	suite.Suite
	environment *mocks.MockEnvironment
	deployment  *mocks.MockDeployment
	ecs         *mocks.MockECS
	api         API

	// We need a router because some of the apis use mux.Vars() which uses the URL
	// parameters parsed and stored in a global map in the global context by the router.
	router *mux.Router
}

func (suite *APITestSuite) SetupTest() {
	mockCtrl := gomock.NewController(suite.T())
	suite.environment = mocks.NewMockEnvironment(mockCtrl)
	suite.deployment = mocks.NewMockDeployment(mockCtrl)
	suite.ecs = mocks.NewMockECS(mockCtrl)
	suite.api = NewAPI(suite.environment, suite.deployment, suite.ecs)
	suite.router = suite.getRouter()
}

func TestAPITestSuite(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}

func (suite *APITestSuite) TestPing() {
	request, err := http.NewRequest("GET", "/v1/ping", nil)
	assert.Nil(suite.T(), err, "Unexpected error creating ping request")

	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	header := responseRecorder.Header()
	assert.NotNil(suite.T(), header, "Unexpected empty header in ping response")
	expectedHeader := http.Header{"Content-Type": []string{"application/json; charset=UTF-8"}}
	assert.Equal(suite.T(), expectedHeader, header, "Http header in ping response is  invalid")
	assert.Equal(suite.T(), http.StatusOK, responseRecorder.Code, "Http response status in ping response is invalid")
}

func (suite *APITestSuite) TestGetEnvironmentReturnsError() {
	name := "testEnv"
	err := errors.New("Error from GetEnvironment")
	suite.environment.EXPECT().GetEnvironment(gomock.Any(), name).Return(nil, err)
	request := suite.generateGetEnvironmentRequest(name)

	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	assert.Equal(suite.T(), http.StatusInternalServerError, responseRecorder.Code)
}

func (suite *APITestSuite) TestGetEnvironmentMissingReturnsError() {
	name := "testEnv"
	suite.environment.EXPECT().GetEnvironment(gomock.Any(), name).Return(nil, nil)
	request := suite.generateGetEnvironmentRequest(name)

	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	assert.Equal(suite.T(), http.StatusNotFound, responseRecorder.Code)
}

func (suite *APITestSuite) TestGetEnvironment() {
	name := "testEnv"
	environment := suite.createEnvironmentObject(name, taskDefinitionARN, clusterARN1)
	suite.environment.EXPECT().GetEnvironment(gomock.Any(), name).Return(environment, nil)
	request := suite.generateGetEnvironmentRequest(name)

	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	assert.Equal(suite.T(), http.StatusOK, responseRecorder.Code)

	var environmentModel models.Environment
	b, _ := ioutil.ReadAll(responseRecorder.Body)
	json.Unmarshal(b, &environmentModel)

	suite.assertSame(environment, &environmentModel)
}

func (suite *APITestSuite) TestListEnvironments() {
	e1 := suite.createEnvironmentObject("e1", taskDefinitionARN, clusterARN1)
	e2 := suite.createEnvironmentObject("e2", taskDefinitionARN, clusterARN2)
	environments := []types.Environment{*e1, *e2}
	suite.environment.EXPECT().ListEnvironments(gomock.Any()).Return(environments, nil)
	request, err := http.NewRequest("GET", "/v1/environments", nil)
	assert.Nil(suite.T(), err, "Unexpected error")

	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	assert.Equal(suite.T(), http.StatusOK, responseRecorder.Code)

	var environmentsModel models.Environments
	b, _ := ioutil.ReadAll(responseRecorder.Body)
	json.Unmarshal(b, &environmentsModel)

	for i := 0; i < len(environments); i++ {
		environment := environments[i]
		environmentModel := environmentsModel.Items[i]
		suite.assertSame(&environment, environmentModel)
	}
}

func (suite *APITestSuite) TestListEnvironmentsServerError() {
	err := errors.New("Error when calling ListEnvironments")
	suite.environment.EXPECT().ListEnvironments(gomock.Any()).Return(nil, err)
	request, err := http.NewRequest("GET", "/v1/environments", nil)
	assert.Nil(suite.T(), err, "Unexpected error")

	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	assert.Equal(suite.T(), http.StatusInternalServerError, responseRecorder.Code)
}

func (suite *APITestSuite) TestFilterEnvironments() {
	e1 := suite.createEnvironmentObject("e1", taskDefinitionARN, clusterARN1)
	e2 := suite.createEnvironmentObject("e2", taskDefinitionARN, clusterARN1)
	environments := []types.Environment{*e1, *e2}
	suite.environment.EXPECT().FilterEnvironments(gomock.Any(), gomock.Any(), gomock.Any()).Return(environments, nil)

	request := suite.generateFilterEnvironmentRequest(clusterARN1)

	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	assert.Equal(suite.T(), http.StatusOK, responseRecorder.Code)

	var environmentsModel models.Environments
	b, _ := ioutil.ReadAll(responseRecorder.Body)
	json.Unmarshal(b, &environmentsModel)

	for i := 0; i < len(environments); i++ {
		environment := environments[i]
		environmentModel := environmentsModel.Items[i]
		suite.assertSame(&environment, environmentModel)
	}
}

func (suite *APITestSuite) TestFilterEnvironmentsServerError() {
	err := errors.New("Error when calling FilterEnvironments")
	suite.environment.EXPECT().FilterEnvironments(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, err)
	request := suite.generateFilterEnvironmentRequest(clusterARN1)

	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	assert.Equal(suite.T(), http.StatusInternalServerError, responseRecorder.Code)
}

func (suite *APITestSuite) TestDeleteEnvironment() {
	name := "testEnv"
	suite.environment.EXPECT().DeleteEnvironment(gomock.Any(), name).Return(nil)

	request := suite.generateDeleteEnvironmentRequest(name)

	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	assert.Equal(suite.T(), http.StatusOK, responseRecorder.Code)
}

func (suite *APITestSuite) TestDeleteEnvironmentMissingEnvironment() {
	name := "testEnv"
	notfounderr := types.NewNotFoundError(errors.New("Environment is missing"))
	suite.environment.EXPECT().DeleteEnvironment(gomock.Any(), name).Return(notfounderr)

	request := suite.generateDeleteEnvironmentRequest(name)

	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	assert.Equal(suite.T(), http.StatusNotFound, responseRecorder.Code)
}

func (suite *APITestSuite) assertSame(environment *types.Environment, environmentModel *models.Environment) {
	assert.Equal(suite.T(), environment.Name, aws.StringValue(environmentModel.Name))
	assert.Equal(suite.T(), environment.Cluster, environmentModel.InstanceGroup.Cluster)
	assert.Equal(suite.T(), environment.DesiredTaskDefinition, environmentModel.TaskDefinition)
}

func (suite *APITestSuite) getRouter() *mux.Router {
	return NewRouter(suite.api)
}

func (suite *APITestSuite) generateGetEnvironmentRequest(name string) *http.Request {
	request, err := http.NewRequest("GET", "/v1/environments/"+name, nil)
	assert.Nil(suite.T(), err, "Unexpected error")
	return request
}

func (suite *APITestSuite) generateDeleteEnvironmentRequest(name string) *http.Request {
	request, err := http.NewRequest("DELETE", "/v1/environments/"+name, nil)
	assert.Nil(suite.T(), err, "Unexpected error")
	return request
}

func (suite *APITestSuite) generateFilterEnvironmentRequest(cluster string) *http.Request {
	url := "/v1/environments?cluster=" + cluster
	request, err := http.NewRequest("GET", url, nil)
	assert.Nil(suite.T(), err, "Unexpected error generating a filter environments request")
	return request
}

func (suite *APITestSuite) createEnvironmentObject(name string, td string, cluster string) *types.Environment {
	environment, err := types.NewEnvironment(name, td, cluster)
	assert.Nil(suite.T(), err, "Unexpected error generating an environment object")
	return environment
}
