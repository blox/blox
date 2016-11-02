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
	filterInstancesByStatusQueryValue  = "{status:active|inactive}"
	filterInstancesByClusterQueryValue = "{cluster:(arn:aws:ecs:)([\\-\\w]+):[0-9]{12}:(cluster)/[a-zA-Z0-9\\-_]{1,255}}"

	getInstancePrefix              = "/v1/instances"
	listInstancesPrefix            = "/v1/instances"
	filterInstancesByStatusPrefix  = "/v1/instances/filter?status="
	filterInstancesByClusterPrefix = "/v1/instances/filter?cluster="

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
	instanceModel1 models.ContainerInstanceModel
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
		PendingTasksCount:    &pendingTaskCount1,
		RegisteredResources:  []*types.Resource{},
		RemainingResources:   []*types.Resource{},
		RunningTasksCount:    &runningTasksCount1,
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

	instanceModel, err := ToContainerInstanceModel(suite.instance1)
	if err != nil {
		suite.T().Error("Cannot setup testSuite: Error when tranlating instance to external model")
	}
	suite.instanceModel1 = instanceModel

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
	instanceInResponse := models.ContainerInstanceModel{}
	err := json.NewDecoder(reader).Decode(&instanceInResponse)
	assert.Nil(suite.T(), err, "Unexpected error decoding response body")
	assert.Exactly(suite.T(), suite.instanceModel1, instanceInResponse, "Instance in response is invalid")
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

	request := suite.listInstancesRequest()
	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateSuccessfulResponseHeaderAndStatus(responseRecorder)
	instanceModelList := []models.ContainerInstanceModel{suite.instanceModel1}
	suite.validateInstancesInListOrFilterInstancesResponse(responseRecorder, instanceModelList)
}

func (suite *InstanceAPIsTestSuite) TestListInstancesReturnsNoInstances() {
	emptyInstanceList := make([]types.ContainerInstance, 0)
	suite.instanceStore.EXPECT().ListContainerInstances().Return(emptyInstanceList, nil)

	request := suite.listInstancesRequest()
	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateSuccessfulResponseHeaderAndStatus(responseRecorder)
	emptyInstanceModelList := make([]models.ContainerInstanceModel, 0)
	suite.validateInstancesInListOrFilterInstancesResponse(responseRecorder, emptyInstanceModelList)
}

func (suite *InstanceAPIsTestSuite) TestListInstancesStoreReturnsError() {
	suite.instanceStore.EXPECT().ListContainerInstances().Return(nil, errors.New("Error when listing instances"))

	request := suite.listInstancesRequest()
	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateErrorResponseHeaderAndStatus(responseRecorder, http.StatusInternalServerError)
	suite.decodeErrorResponseAndValidate(responseRecorder, internalServerErrMsg)
}

func (suite *InstanceAPIsTestSuite) TestFilterInstancesByStatusReturnsInstances() {
	instanceList := []types.ContainerInstance{suite.instance1}
	suite.instanceStore.EXPECT().FilterContainerInstances(instanceStatusFilter, instanceStatus1).Return(instanceList, nil)
	suite.instanceStore.EXPECT().FilterContainerInstances(instanceClusterFilter, gomock.Any()).Times(0)

	request := suite.filterInstancesByStatusRequest()
	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateSuccessfulResponseHeaderAndStatus(responseRecorder)
	instanceModelList := []models.ContainerInstanceModel{suite.instanceModel1}
	suite.validateInstancesInListOrFilterInstancesResponse(responseRecorder, instanceModelList)
}

func (suite *InstanceAPIsTestSuite) TestFilterInstancesByStatusReturnsNoInstances() {
	emptyInstanceList := []types.ContainerInstance{}
	suite.instanceStore.EXPECT().FilterContainerInstances(instanceStatusFilter, instanceStatus1).Return(emptyInstanceList, nil)
	suite.instanceStore.EXPECT().FilterContainerInstances(instanceClusterFilter, gomock.Any()).Times(0)

	request := suite.filterInstancesByStatusRequest()
	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateSuccessfulResponseHeaderAndStatus(responseRecorder)
	emptyInstanceModelList := []models.ContainerInstanceModel{}
	suite.validateInstancesInListOrFilterInstancesResponse(responseRecorder, emptyInstanceModelList)
}

func (suite *InstanceAPIsTestSuite) TestFilterInstancesByStatusStoreReturnsError() {
	suite.instanceStore.EXPECT().FilterContainerInstances(instanceStatusFilter, instanceStatus1).Return(nil, errors.New("Error when filtering instances"))
	suite.instanceStore.EXPECT().FilterContainerInstances(instanceClusterFilter, gomock.Any()).Times(0)

	request := suite.filterInstancesByStatusRequest()
	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateErrorResponseHeaderAndStatus(responseRecorder, http.StatusInternalServerError)
	suite.decodeErrorResponseAndValidate(responseRecorder, internalServerErrMsg)
}

func (suite *InstanceAPIsTestSuite) TestFilterInstancesByClusterReturnsInstances() {
	instanceList := []types.ContainerInstance{suite.instance1}
	suite.instanceStore.EXPECT().FilterContainerInstances(instanceClusterFilter, clusterARN1).Return(instanceList, nil)
	suite.instanceStore.EXPECT().FilterContainerInstances(instanceStatusFilter, gomock.Any()).Times(0)

	request := suite.filterInstancesByClusterRequest()
	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateSuccessfulResponseHeaderAndStatus(responseRecorder)
	instanceModelList := []models.ContainerInstanceModel{suite.instanceModel1}
	suite.validateInstancesInListOrFilterInstancesResponse(responseRecorder, instanceModelList)
}

func (suite *InstanceAPIsTestSuite) TestFilterInstancesByClusterReturnsNoInstances() {
	emptyInstanceList := []types.ContainerInstance{}
	suite.instanceStore.EXPECT().FilterContainerInstances(instanceClusterFilter, clusterARN1).Return(emptyInstanceList, nil)
	suite.instanceStore.EXPECT().FilterContainerInstances(instanceStatusFilter, gomock.Any()).Times(0)

	request := suite.filterInstancesByClusterRequest()
	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateSuccessfulResponseHeaderAndStatus(responseRecorder)
	emptyInstanceModelList := []models.ContainerInstanceModel{}
	suite.validateInstancesInListOrFilterInstancesResponse(responseRecorder, emptyInstanceModelList)
}

func (suite *InstanceAPIsTestSuite) TestFilterInstancesByClusterStoreReturnsError() {
	suite.instanceStore.EXPECT().FilterContainerInstances(instanceClusterFilter, clusterARN1).Return(nil, errors.New("Error when filtering instances"))
	suite.instanceStore.EXPECT().FilterContainerInstances(instanceStatusFilter, gomock.Any()).Times(0)

	request := suite.filterInstancesByClusterRequest()
	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateErrorResponseHeaderAndStatus(responseRecorder, http.StatusInternalServerError)
	suite.decodeErrorResponseAndValidate(responseRecorder, internalServerErrMsg)
}

func (suite *InstanceAPIsTestSuite) TestFilterInstancesByUnsupportedKey() {
	suite.instanceStore.EXPECT().FilterContainerInstances(instanceStatusFilter, gomock.Any()).Times(0)
	suite.instanceStore.EXPECT().FilterContainerInstances(instanceClusterFilter, gomock.Any()).Times(0)

	url := unsupportedFilterInstancesPrefix + "filterVal"
	request, err := http.NewRequest("GET", url, nil)
	assert.Nil(suite.T(), err, "Unexpected error creating filter instances request")

	responseRecorder := httptest.NewRecorder()
	suite.router.ServeHTTP(responseRecorder, request)

	suite.validateErrorResponseHeaderAndStatus(responseRecorder, http.StatusInternalServerError)
	suite.decodeErrorResponseAndValidate(responseRecorder, routingServerErrMsg)
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

	s.Path(filterInstancesPath).
		Queries(instanceStatusFilter, filterInstancesByStatusQueryValue).
		Methods("GET").
		HandlerFunc(suite.instanceAPIs.FilterInstances)

	s.Path(filterInstancesPath).
		Queries(instanceClusterFilter, filterInstancesByClusterQueryValue).
		Methods("GET").
		HandlerFunc(suite.instanceAPIs.FilterInstances)

	// Invalid router paths to make sure handler functions handle them
	s.Path(invalidGetInstancePath).
		Methods("GET").
		HandlerFunc(suite.instanceAPIs.GetInstance)

	s.Path(filterInstancesPath).
		Queries(unsupportedFilterInstancesKey, unsupportedFilterInstancesQueryValue).
		Methods("GET").
		HandlerFunc(suite.instanceAPIs.FilterInstances)

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

func (suite *InstanceAPIsTestSuite) filterInstancesByStatusRequest() *http.Request {
	url := filterInstancesByStatusPrefix + instanceStatus1
	request, err := http.NewRequest("GET", url, nil)
	assert.Nil(suite.T(), err, "Unexpected error creating filter instances by status request")
	return request
}

func (suite *InstanceAPIsTestSuite) filterInstancesByClusterRequest() *http.Request {
	url := filterInstancesByClusterPrefix + clusterARN1
	request, err := http.NewRequest("GET", url, nil)
	assert.Nil(suite.T(), err, "Unexpected error creating filter instances by cluster request")
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
	reader := bytes.NewReader(responseRecorder.Body.Bytes())
	var str string
	err := json.NewDecoder(reader).Decode(&str)
	assert.Nil(suite.T(), err, "Unexpected error decoding response body")
	assert.Equal(suite.T(), expectedErrMsg, str)
}

func (suite *InstanceAPIsTestSuite) validateInstancesInListOrFilterInstancesResponse(responseRecorder *httptest.ResponseRecorder, expectedInstanceModels []models.ContainerInstanceModel) {
	reader := bytes.NewReader(responseRecorder.Body.Bytes())
	instancesInResponse := new([]models.ContainerInstanceModel)
	err := json.NewDecoder(reader).Decode(instancesInResponse)
	assert.Nil(suite.T(), err, "Unexpected error decoding response body")
	assert.Exactly(suite.T(), expectedInstanceModels, *instancesInResponse, "Instances in response are invalid")
}
