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

	"github.com/blox/blox/cluster-state-service/handler/mocks"
	"github.com/blox/blox/cluster-state-service/handler/types"
	"github.com/blox/blox/cluster-state-service/swagger/v1/generated/models"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	filterInstancesByStatusQueryValue  = "{status:active|inactive}"
	filterInstancesByClusterQueryValue = "{cluster:(arn:aws:ecs:)([\\-\\w]+):[0-9]{12}:(cluster)/[a-zA-Z0-9\\-_]{1,255}}"

	getInstancePrefix              = "/v1/instances"
	listInstancesPrefix            = "/v1/instances"
	filterInstancesByStatusPrefix  = "/v1/instances?status="
	filterInstancesByClusterPrefix = "/v1/instances?cluster="

	// Routing to GetInstance handler function without arn
	invalidGetInstancePath = "/instances/{cluster:[a-zA-Z0-9_]{1,255}}"

	// Routing to FilterInstances handler function
	unsupportedFilterInstancesKey        = "unsupportedFilter"
	unsupportedFilterInstancesQueryValue = "{unsupportedFilter:([\\-\\w]+)}"
	unsupportedFilterInstancesPrefix     = "/v1/instances/filter?unsupportedFilter="
)

type InstanceAPIsTestSuite struct {
	suite.Suite
	instanceStore  *mocks.MockContainerInstanceStore
	instanceAPIs   ContainerInstanceAPIs
	instance1      types.ContainerInstance
	extInstance1   models.ContainerInstance
	responseHeader http.Header

	// We need a router because some of the apis use mux.Vars() which uses the URL
	// parameters parsed and stored in a global map in the global context by the router.
	router *mux.Router
}

func (suite *InstanceAPIsTestSuite) SetupTest() {
	mockCtrl := gomock.NewController(suite.T())

	suite.instanceStore = mocks.NewMockContainerInstanceStore(mockCtrl)

	suite.instanceAPIs = NewContainerInstanceAPIs(suite.instanceStore)

	versionInfo := types.VersionInfo{}
	instanceDetail := types.InstanceDetail{
		AgentConnected:       &agentConnected1,
		ClusterARN:           &clusterARN1,
		ContainerInstanceARN: &instanceARN1,
		RegisteredResources:  []*types.Resource{},
		RemainingResources:   []*types.Resource{},
		Status:               &instanceStatus1,
		Version:              &version1,
		VersionInfo:          &versionInfo,
		UpdatedAt:            &updatedAt1,
	}
	suite.instance1 = types.ContainerInstance{
		ID:        &id1,
		Account:   &accountID,
		Time:      &time,
		Region:    &region,
		Resources: []string{instanceARN1},
		Detail:    &instanceDetail,
	}

	instanceModel, err := ToContainerInstance(suite.instance1)
	if err != nil {
		suite.T().Error("Cannot setup testSuite: Error when tranlating instance to external model")
	}
	suite.extInstance1 = instanceModel

	suite.responseHeader = http.Header{responseContentTypeKey: []string{responseContentTypeVal}}

	suite.router = suite.getRouter()
}

func TestInstanceAPIsTestSuite(t *testing.T) {
	suite.Run(t, new(InstanceAPIsTestSuite))
}

func (suite *InstanceAPIsTestSuite) TestGetInstanceReturnsInstance() {
	suite.instanceStore.EXPECT().GetContainerInstance(clusterName1, instanceARN1).Return(&suite.instance1, nil)

	request := suite.getInstanceRequest()
	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateSuccessfulResponseHeaderAndStatus(responseRecorder)

	reader := bytes.NewReader(responseRecorder.Body.Bytes())
	instanceInResponse := models.ContainerInstance{}
	err := json.NewDecoder(reader).Decode(&instanceInResponse)
	assert.Nil(suite.T(), err, "Unexpected error decoding response body")
	assert.Exactly(suite.T(), suite.extInstance1, instanceInResponse, "Instance in response is invalid")
}

func (suite *InstanceAPIsTestSuite) TestGetInstanceReturnsNoInstance() {
	suite.instanceStore.EXPECT().GetContainerInstance(clusterName1, instanceARN1).Return(nil, nil)

	request := suite.getInstanceRequest()
	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateErrorResponseHeaderAndStatus(responseRecorder, http.StatusNotFound)
	suite.decodeErrorResponseAndValidate(responseRecorder, instanceNotFoundClientErrMsg)
}

func (suite *InstanceAPIsTestSuite) TestGetInstanceStoreReturnsError() {
	suite.instanceStore.EXPECT().GetContainerInstance(clusterName1, instanceARN1).Return(nil, errors.New("Error when getting instance"))

	request := suite.getInstanceRequest()
	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateErrorResponseHeaderAndStatus(responseRecorder, http.StatusInternalServerError)
	suite.decodeErrorResponseAndValidate(responseRecorder, internalServerErrMsg)
}

func (suite *InstanceAPIsTestSuite) TestGetInstanceWithoutARN() {
	suite.instanceStore.EXPECT().GetContainerInstance(gomock.Any(), gomock.Any()).Times(0)

	url := getInstancePrefix + "/" + clusterName1
	request, err := http.NewRequest("GET", url, nil)
	assert.Nil(suite.T(), err, "Unexpected error creating task get request")

	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateErrorResponseHeaderAndStatus(responseRecorder, http.StatusInternalServerError)
	suite.decodeErrorResponseAndValidate(responseRecorder, routingServerErrMsg)
}

func (suite *InstanceAPIsTestSuite) TestListInstancesReturnsInstances() {
	instanceList := []types.ContainerInstance{suite.instance1}
	suite.instanceStore.EXPECT().ListContainerInstances().Return(instanceList, nil)
	suite.instanceStore.EXPECT().FilterContainerInstances(gomock.Any()).Times(0)

	request := suite.listInstancesRequest()
	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateSuccessfulResponseHeaderAndStatus(responseRecorder)
	extInstances := models.ContainerInstances{
		Items: []*models.ContainerInstance{&suite.extInstance1},
	}
	suite.validateInstancesInListOrFilterInstancesResponse(responseRecorder, extInstances)
}

func (suite *InstanceAPIsTestSuite) TestListInstancesReturnsNoInstances() {
	emptyInstanceList := make([]types.ContainerInstance, 0)
	suite.instanceStore.EXPECT().ListContainerInstances().Return(emptyInstanceList, nil)
	suite.instanceStore.EXPECT().FilterContainerInstances(gomock.Any()).Times(0)

	request := suite.listInstancesRequest()
	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateSuccessfulResponseHeaderAndStatus(responseRecorder)
	emptyExtInstances := models.ContainerInstances{
		Items: []*models.ContainerInstance{},
	}
	suite.validateInstancesInListOrFilterInstancesResponse(responseRecorder, emptyExtInstances)
}

func (suite *InstanceAPIsTestSuite) TestListInstancesStoreReturnsError() {
	suite.instanceStore.EXPECT().ListContainerInstances().Return(nil, errors.New("Error when listing instances"))
	suite.instanceStore.EXPECT().FilterContainerInstances(gomock.Any()).Times(0)

	request := suite.listInstancesRequest()
	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateErrorResponseHeaderAndStatus(responseRecorder, http.StatusInternalServerError)
	suite.decodeErrorResponseAndValidate(responseRecorder, internalServerErrMsg)
}

func (suite *InstanceAPIsTestSuite) TestListInstancesInvalidFilter() {
	suite.instanceStore.EXPECT().ListContainerInstances().Times(0)
	suite.instanceStore.EXPECT().FilterContainerInstances(gomock.Any()).Times(0)

	url := listInstancesPrefix + "?unsupportedFilter=val"
	request, err := http.NewRequest("GET", url, nil)
	assert.Nil(suite.T(), err, "Unexpected error creating list instances request with invalid filter")

	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateErrorResponseHeaderAndStatus(responseRecorder, http.StatusBadRequest)
	suite.decodeErrorResponseAndValidate(responseRecorder, unsupportedFilterClientErrMsg)
}

func (suite *InstanceAPIsTestSuite) TestListInstancesStatusAndClusterARNFilter() {
	instanceList := []types.ContainerInstance{suite.instance1}
	filters := map[string]string{instanceStatusFilter: instanceStatus1, instanceClusterFilter: clusterARN1}
	suite.instanceStore.EXPECT().FilterContainerInstances(filters).
		Return(instanceList, nil)
	suite.instanceStore.EXPECT().ListContainerInstances().Times(0)

	request := suite.filterInstancesByStatusAndClusterRequest(instanceStatus1, clusterARN1)

	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateSuccessfulResponseHeaderAndStatus(responseRecorder)
	extInstances := models.ContainerInstances{
		Items: []*models.ContainerInstance{&suite.extInstance1},
	}
	suite.validateInstancesInListOrFilterInstancesResponse(responseRecorder, extInstances)
}

func (suite *InstanceAPIsTestSuite) TestListInstancesStatusAndClusterNameFilter() {
	instanceList := []types.ContainerInstance{suite.instance1}
	filters := map[string]string{instanceStatusFilter: instanceStatus1, instanceClusterFilter: clusterName1}
	suite.instanceStore.EXPECT().FilterContainerInstances(filters).
		Return(instanceList, nil)
	suite.instanceStore.EXPECT().ListContainerInstances().Times(0)

	request := suite.filterInstancesByStatusAndClusterRequest(instanceStatus1, clusterName1)

	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateSuccessfulResponseHeaderAndStatus(responseRecorder)
	extInstances := models.ContainerInstances{
		Items: []*models.ContainerInstance{&suite.extInstance1},
	}
	suite.validateInstancesInListOrFilterInstancesResponse(responseRecorder, extInstances)
}

func (suite *InstanceAPIsTestSuite) TestListInstancesWithStatusFilterReturnsInstances() {
	instanceList := []types.ContainerInstance{suite.instance1}
	suite.instanceStore.EXPECT().FilterContainerInstances(map[string]string{instanceStatusFilter: instanceStatus1}).
		Return(instanceList, nil)
	suite.instanceStore.EXPECT().ListContainerInstances().Times(0)

	request := suite.filterInstancesByStatusRequest(instanceStatus1)
	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateSuccessfulResponseHeaderAndStatus(responseRecorder)
	extInstances := models.ContainerInstances{
		Items: []*models.ContainerInstance{&suite.extInstance1},
	}
	suite.validateInstancesInListOrFilterInstancesResponse(responseRecorder, extInstances)
}

func (suite *InstanceAPIsTestSuite) TestListInstancesWithCapitalizedStatusFilterReturnsInstances() {
	instanceList := []types.ContainerInstance{suite.instance1}
	suite.instanceStore.EXPECT().FilterContainerInstances(map[string]string{instanceStatusFilter: instanceStatus1}).
		Return(instanceList, nil)
	suite.instanceStore.EXPECT().ListContainerInstances().Times(0)

	request := suite.filterInstancesByStatusRequest(strings.ToUpper(instanceStatus1))
	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateSuccessfulResponseHeaderAndStatus(responseRecorder)
	extInstances := models.ContainerInstances{
		Items: []*models.ContainerInstance{&suite.extInstance1},
	}
	suite.validateInstancesInListOrFilterInstancesResponse(responseRecorder, extInstances)
}

func (suite *InstanceAPIsTestSuite) TestListInstancesWithStatusFilterReturnsNoInstances() {
	emptyInstanceList := []types.ContainerInstance{}
	suite.instanceStore.EXPECT().FilterContainerInstances(map[string]string{instanceStatusFilter: instanceStatus1}).
		Return(emptyInstanceList, nil)
	suite.instanceStore.EXPECT().ListContainerInstances().Times(0)

	request := suite.filterInstancesByStatusRequest(instanceStatus1)
	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateSuccessfulResponseHeaderAndStatus(responseRecorder)
	extInstances := models.ContainerInstances{
		Items: []*models.ContainerInstance{},
	}
	suite.validateInstancesInListOrFilterInstancesResponse(responseRecorder, extInstances)
}

func (suite *InstanceAPIsTestSuite) TestFilterInstancesWithStatusFilterStoreReturnsError() {
	suite.instanceStore.EXPECT().FilterContainerInstances(map[string]string{instanceStatusFilter: instanceStatus1}).
		Return(nil, errors.New("Error when filtering instances"))
	suite.instanceStore.EXPECT().ListContainerInstances().Times(0)

	request := suite.filterInstancesByStatusRequest(instanceStatus1)
	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateErrorResponseHeaderAndStatus(responseRecorder, http.StatusInternalServerError)
	suite.decodeErrorResponseAndValidate(responseRecorder, internalServerErrMsg)
}

func (suite *InstanceAPIsTestSuite) TestListInstancesWithInvalidStatusFilter() {
	suite.instanceStore.EXPECT().FilterContainerInstances(gomock.Any()).Times(0)
	suite.instanceStore.EXPECT().ListContainerInstances().Times(0)

	url := filterInstancesByStatusPrefix + "invalidStatus"
	request, err := http.NewRequest("GET", url, nil)
	assert.Nil(suite.T(), err, "Unexpected error creating list instances request with invalid status")

	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateErrorResponseHeaderAndStatus(responseRecorder, http.StatusBadRequest)
	suite.decodeErrorResponseAndValidate(responseRecorder, invalidStatusClientErrMsg)
}

func (suite *InstanceAPIsTestSuite) TestListInstancesWithClusterNameReturnsInstances() {
	instanceList := []types.ContainerInstance{suite.instance1}
	suite.instanceStore.EXPECT().FilterContainerInstances(map[string]string{instanceClusterFilter: clusterName1}).
		Return(instanceList, nil)
	suite.instanceStore.EXPECT().ListContainerInstances().Times(0)

	request := suite.filterInstancesByClusterRequest(clusterName1)
	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateSuccessfulResponseHeaderAndStatus(responseRecorder)
	extInstances := models.ContainerInstances{
		Items: []*models.ContainerInstance{&suite.extInstance1},
	}
	suite.validateInstancesInListOrFilterInstancesResponse(responseRecorder, extInstances)
}

func (suite *InstanceAPIsTestSuite) TestListInstancesWithClusterNameReturnsNoInstances() {
	emptyInstanceList := []types.ContainerInstance{}
	suite.instanceStore.EXPECT().FilterContainerInstances(map[string]string{instanceClusterFilter: clusterName1}).
		Return(emptyInstanceList, nil)
	suite.instanceStore.EXPECT().ListContainerInstances().Times(0)

	request := suite.filterInstancesByClusterRequest(clusterName1)
	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateSuccessfulResponseHeaderAndStatus(responseRecorder)
	emptyExtInstances := models.ContainerInstances{
		Items: []*models.ContainerInstance{},
	}
	suite.validateInstancesInListOrFilterInstancesResponse(responseRecorder, emptyExtInstances)
}

func (suite *InstanceAPIsTestSuite) TestListInstancesWithClusterNameFilterStoreReturnsError() {
	suite.instanceStore.EXPECT().FilterContainerInstances(map[string]string{instanceClusterFilter: clusterName1}).
		Return(nil, errors.New("Error when filtering instances"))
	suite.instanceStore.EXPECT().ListContainerInstances().Times(0)

	request := suite.filterInstancesByClusterRequest(clusterName1)
	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateErrorResponseHeaderAndStatus(responseRecorder, http.StatusInternalServerError)
	suite.decodeErrorResponseAndValidate(responseRecorder, internalServerErrMsg)
}

func (suite *InstanceAPIsTestSuite) TestListInstancesWithClusterARNReturnsInstances() {
	instanceList := []types.ContainerInstance{suite.instance1}
	suite.instanceStore.EXPECT().FilterContainerInstances(map[string]string{instanceClusterFilter: clusterARN1}).
		Return(instanceList, nil)
	suite.instanceStore.EXPECT().ListContainerInstances().Times(0)

	request := suite.filterInstancesByClusterRequest(clusterARN1)
	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateSuccessfulResponseHeaderAndStatus(responseRecorder)
	extInstances := models.ContainerInstances{
		Items: []*models.ContainerInstance{&suite.extInstance1},
	}
	suite.validateInstancesInListOrFilterInstancesResponse(responseRecorder, extInstances)
}

func (suite *InstanceAPIsTestSuite) TestListInstancesWithClusterARNReturnsNoInstances() {
	emptyInstanceList := []types.ContainerInstance{}
	suite.instanceStore.EXPECT().FilterContainerInstances(map[string]string{instanceClusterFilter: clusterARN1}).
		Return(emptyInstanceList, nil)
	suite.instanceStore.EXPECT().ListContainerInstances().Times(0)

	request := suite.filterInstancesByClusterRequest(clusterARN1)
	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateSuccessfulResponseHeaderAndStatus(responseRecorder)
	emptyExtInstances := models.ContainerInstances{
		Items: []*models.ContainerInstance{},
	}
	suite.validateInstancesInListOrFilterInstancesResponse(responseRecorder, emptyExtInstances)
}

func (suite *InstanceAPIsTestSuite) TestListInstancesWithClusterARNFilterStoreReturnsError() {
	suite.instanceStore.EXPECT().FilterContainerInstances(map[string]string{instanceClusterFilter: clusterARN1}).
		Return(nil, errors.New("Error when filtering instances"))
	suite.instanceStore.EXPECT().ListContainerInstances().Times(0)

	request := suite.filterInstancesByClusterRequest(clusterARN1)
	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateErrorResponseHeaderAndStatus(responseRecorder, http.StatusInternalServerError)
	suite.decodeErrorResponseAndValidate(responseRecorder, internalServerErrMsg)
}

func (suite *InstanceAPIsTestSuite) TestListInstancesWithInvalidClusterFilter() {
	url := filterInstancesByClusterPrefix + "cluster/cluster"
	request, err := http.NewRequest("GET", url, nil)
	assert.Nil(suite.T(), err, "Unexpected error creating list instances request with invalid cluster")

	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateErrorResponseHeaderAndStatus(responseRecorder, http.StatusBadRequest)
	suite.decodeErrorResponseAndValidate(responseRecorder, invalidClusterClientErrMsg)
}

// Helper functions

func (suite *InstanceAPIsTestSuite) getRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	s := r.Path("/v1").Subrouter()

	s.Path(getInstancePath).
		Methods("GET").
		HandlerFunc(suite.instanceAPIs.GetInstance)

	s.Path(listInstancesPath).Methods("GET").
		HandlerFunc(suite.instanceAPIs.ListInstances)

	// Invalid router paths to make sure handler functions handle them
	s.Path(invalidGetInstancePath).
		Methods("GET").
		HandlerFunc(suite.instanceAPIs.GetInstance)

	return s
}

func (suite *InstanceAPIsTestSuite) getInstanceRequest() *http.Request {
	url := getInstancePrefix + "/" + clusterName1 + "/" + instanceARN1
	request, err := http.NewRequest("GET", url, nil)
	assert.Nil(suite.T(), err, "Unexpected error creating get instance request")
	return request
}

func (suite *InstanceAPIsTestSuite) listInstancesRequest() *http.Request {
	request, err := http.NewRequest("GET", listInstancesPrefix, nil)
	assert.Nil(suite.T(), err, "Unexpected error creating list instances request")
	return request
}

func (suite *InstanceAPIsTestSuite) filterInstancesByStatusRequest(status string) *http.Request {
	url := filterInstancesByStatusPrefix + status
	request, err := http.NewRequest("GET", url, nil)
	assert.Nil(suite.T(), err, "Unexpected error creating filter instances by status request")
	return request
}

func (suite *InstanceAPIsTestSuite) filterInstancesByClusterRequest(cluster string) *http.Request {
	url := filterInstancesByClusterPrefix + cluster
	request, err := http.NewRequest("GET", url, nil)
	assert.Nil(suite.T(), err, "Unexpected error creating filter instances by cluster request")
	return request
}

func (suite *InstanceAPIsTestSuite) filterInstancesByStatusAndClusterRequest(status string, cluster string) *http.Request {
	url := "/v1/instances?status=" + status + "&cluster=" + cluster
	request, err := http.NewRequest("GET", url, nil)
	assert.Nil(suite.T(), err, "Unexpected error creating filter instances by status and cluster request")
	return request
}

func (suite *InstanceAPIsTestSuite) validateSuccessfulResponseHeaderAndStatus(responseRecorder *httptest.ResponseRecorder) {
	h := responseRecorder.Header()
	assert.NotNil(suite.T(), h, "Unexpected empty header")
	assert.Equal(suite.T(), suite.responseHeader, h, "Http header is invalid")
	assert.Equal(suite.T(), http.StatusOK, responseRecorder.Code, "Http response status is invalid")
}

func (suite *InstanceAPIsTestSuite) validateErrorResponseHeaderAndStatus(responseRecorder *httptest.ResponseRecorder, errorCode int) {
	h := responseRecorder.Header()
	assert.NotNil(suite.T(), h, "Unexpected empty header")
	assert.Equal(suite.T(), errorCode, responseRecorder.Code, "Http response status is invalid")
}

func (suite *InstanceAPIsTestSuite) decodeErrorResponseAndValidate(responseRecorder *httptest.ResponseRecorder, expectedErrMsg string) {
	actualMsg := responseRecorder.Body.String()
	assert.Equal(suite.T(), expectedErrMsg+"\n", actualMsg, "Error message is invalid")
}

func (suite *InstanceAPIsTestSuite) validateInstancesInListOrFilterInstancesResponse(responseRecorder *httptest.ResponseRecorder, expectedInstances models.ContainerInstances) {
	reader := bytes.NewReader(responseRecorder.Body.Bytes())
	instancesInResponse := new(models.ContainerInstances)
	err := json.NewDecoder(reader).Decode(instancesInResponse)
	assert.Nil(suite.T(), err, "Unexpected error decoding response body")
	assert.Exactly(suite.T(), expectedInstances, *instancesInResponse, "Instances in response are invalid")
}
