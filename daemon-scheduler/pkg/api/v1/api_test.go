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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blox/blox/daemon-scheduler/pkg/mocks"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
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

func (suite *APITestSuite) getRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	s := r.Path("/v1").Subrouter()

	s.Path("/ping").
		Methods("GET").
		HandlerFunc(suite.api.Ping)

	return s
}
