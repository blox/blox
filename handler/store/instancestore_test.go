package store

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/aws/amazon-ecs-event-stream-handler/handler/compress"
	"github.com/aws/amazon-ecs-event-stream-handler/handler/mocks"
	storetypes "github.com/aws/amazon-ecs-event-stream-handler/handler/store/types"
	"github.com/aws/amazon-ecs-event-stream-handler/handler/types"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
)

const (
	containerInstanceARN1 = "arn:aws:us-east-1:123456789123:container-instance/4b6d45ea-a4b4-4269-9d04-3af6ddfdc597"
	containerInstanceARN2 = "arn:aws:us-east-1:123456789123:container-instance/3af93452-d6b7-6759-0923-4f5123cfd025"
	clusterName1          = "cluster1"
	clusterName2          = "cluster2"
	clusterName3          = "cluster3"
	clusterARN1           = "arn:aws:ecs:us-east-1:123456789123:cluster/" + clusterName1
	clusterARN2           = "arn:aws:ecs:us-east-1:123456789123:cluster/" + clusterName2
	status1               = "active"
	status2               = "inactive"
)

type instanceStoreMockContext struct {
	mockCtrl                *gomock.Controller
	datastore               *mocks.MockDataStore
	instance1               types.ContainerInstance
	instance2               types.ContainerInstance
	instanceJSON1           string
	instanceJSON2           string
	compressedInstanceJSON1 string
	compressedInstanceJSON2 string
	instanceKey1            string
	instanceKey2            string
}

func NewContainerInstanceStoreMockContext(t *testing.T) *instanceStoreMockContext {
	context := instanceStoreMockContext{}
	context.mockCtrl = gomock.NewController(t)
	context.datastore = mocks.NewMockDataStore(context.mockCtrl)

	context.instance1 = types.ContainerInstance{}
	context.instance1.Detail.ContainerInstanceArn = containerInstanceARN1
	context.instance1.Detail.ClusterArn = clusterARN1
	context.instance1.Detail.Status = status1
	context.instance1.Detail.Version = 1
	context.instanceJSON1 = marshalInstance(t, context.instance1)
	context.compressedInstanceJSON1 = compressInstanceJSON(t, context.instanceJSON1)
	context.instanceKey1 = instanceKeyPrefix + containerInstanceARN1

	context.instance2 = types.ContainerInstance{}
	context.instance2.Detail.ContainerInstanceArn = containerInstanceARN2
	context.instance2.Detail.ClusterArn = clusterARN2
	context.instance2.Detail.Status = status2
	context.instance2.Detail.Version = 1
	context.instanceJSON2 = marshalInstance(t, context.instance2)
	context.compressedInstanceJSON2 = compressInstanceJSON(t, context.instanceJSON2)
	context.instanceKey2 = instanceKeyPrefix + containerInstanceARN2

	return &context
}

func TestInstanceStoreNilDatastore(t *testing.T) {
	_, err := NewContainerInstanceStore(nil)
	if err == nil {
		t.Error("Expected an error when datastore is nil")
	}
}

func TestInstanceStore(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore := instanceStore(t, context)

	if instanceStore == nil {
		t.Error("Instancestore should not be nil")
	}
}

func TestAddContainerInstanceEmptyInstanceJSON(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore := instanceStore(t, context)

	err := instanceStore.AddContainerInstance("")

	if err == nil {
		t.Error("Expected an error when instance JSON is empty in AddContainerInstance")
	}
}

func TestAddContainerInstanceJSONUnmarshalError(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore := instanceStore(t, context)

	err := instanceStore.AddContainerInstance("invalidJSON")

	if err == nil {
		t.Error("Expected an error when instance JSON is invalid in AddContainerInstance")
	}
}

func TestAddContainerInstanceEmptyContainerInstanceARN(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore := instanceStore(t, context)

	instance := types.ContainerInstance{}

	instanceJSON, err := json.Marshal(instance)
	err = instanceStore.AddContainerInstance(string(instanceJSON))

	if err == nil {
		t.Error("Expected an error when container instance arn is not set")
	}
}

func TestAddContainerInstanceGetContainerInstanceFails(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore := instanceStore(t, context)

	context.datastore.EXPECT().Get(context.instanceKey1).Return(nil, errors.New("Error when getting key"))

	err := instanceStore.AddContainerInstance(context.instanceJSON1)

	if err == nil {
		t.Error("Expected an error when datastore get fails")
	}
}

func TestAddContainerInstanceGetContainerInstanceNoResults(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore := instanceStore(t, context)

	context.datastore.EXPECT().Get(context.instanceKey1).Return(make(map[string]string), nil)
	context.datastore.EXPECT().Add(context.instanceKey1, context.compressedInstanceJSON1).Return(nil)

	err := instanceStore.AddContainerInstance(context.instanceJSON1)

	if err != nil {
		t.Error("Unexpected error when datastore get returns empty results")
	}
}

func TestAddContainerInstanceGetContainerInstanceMultipleResults(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore := instanceStore(t, context)

	resp := map[string]string{
		containerInstanceARN1: context.compressedInstanceJSON1,
		containerInstanceARN2: context.compressedInstanceJSON2,
	}
	context.datastore.EXPECT().Get(context.instanceKey1).Return(resp, nil)

	err := instanceStore.AddContainerInstance(context.instanceJSON1)

	if err == nil {
		t.Error("Expected an error when datastore get returns multiple results")
	}
}

func TestAddContainerInstanceGetContainerInstanceInvalidJsonResult(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore := instanceStore(t, context)

	resp := map[string]string{
		containerInstanceARN1: "invalidJSON",
	}
	context.datastore.EXPECT().Get(context.instanceKey1).Return(resp, nil)

	err := instanceStore.AddContainerInstance(context.instanceJSON1)

	if err == nil {
		t.Error("Expected an error when datastore get returns invalid JSON")
	}
}

func TestAddContainerInstanceSameVersionInstanceExists(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore := instanceStore(t, context)

	resp := map[string]string{
		containerInstanceARN1: context.compressedInstanceJSON1,
	}
	context.datastore.EXPECT().Get(context.instanceKey1).Return(resp, nil)
	context.datastore.EXPECT().Add(context.instanceKey1, context.compressedInstanceJSON1).Times(0)

	err := instanceStore.AddContainerInstance(context.instanceJSON1)

	if err != nil {
		t.Error("Unxpected error when adding instance and same version instance exists")
	}
}

func TestAddContainerInstanceHigherVersionInstanceExists(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore := instanceStore(t, context)

	instance := types.ContainerInstance{}
	instance.Detail.ContainerInstanceArn = containerInstanceARN1
	instance.Detail.Version = context.instance1.Detail.Version + 1

	instanceJSON := marshalInstance(t, instance)
	compressedJSON := compressInstanceJSON(t, instanceJSON)

	resp := map[string]string{
		containerInstanceARN1: compressedJSON,
	}
	context.datastore.EXPECT().Get(context.instanceKey1).Return(resp, nil)
	context.datastore.EXPECT().Add(context.instanceKey1, context.compressedInstanceJSON1).Times(0)

	err := instanceStore.AddContainerInstance(context.instanceJSON1)

	if err != nil {
		t.Error("Unxpected error when adding instance and higher version instance exists")
	}
}

func TestAddContainerInstanceLowerVersionInstanceExists(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore := instanceStore(t, context)

	instance := types.ContainerInstance{}
	instance.Detail.ContainerInstanceArn = containerInstanceARN1
	instance.Detail.Version = context.instance1.Detail.Version - 1

	instanceJSON := marshalInstance(t, instance)
	compressedJSON := compressInstanceJSON(t, instanceJSON)

	resp := map[string]string{
		containerInstanceARN1: compressedJSON,
	}
	context.datastore.EXPECT().Get(context.instanceKey1).Return(resp, nil)
	context.datastore.EXPECT().Add(context.instanceKey1, context.compressedInstanceJSON1).Return(nil)

	err := instanceStore.AddContainerInstance(context.instanceJSON1)

	if err != nil {
		t.Error("Unxpected error when adding instance and higher version instance exists")
	}
}

func TestAddContainerInstanceFails(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore := instanceStore(t, context)

	instance := types.ContainerInstance{}
	instance.Detail.ContainerInstanceArn = containerInstanceARN1
	instance.Detail.Version = context.instance1.Detail.Version - 1

	instanceJSON := marshalInstance(t, instance)
	compressedJSON := compressInstanceJSON(t, instanceJSON)

	resp := map[string]string{
		containerInstanceARN1: compressedJSON,
	}
	context.datastore.EXPECT().Get(context.instanceKey1).Return(resp, nil)
	context.datastore.EXPECT().Add(context.instanceKey1, context.compressedInstanceJSON1).Return(errors.New("Add instance failed"))

	err := instanceStore.AddContainerInstance(context.instanceJSON1)

	if err == nil {
		t.Error("Expected an error when adding an instance fails")
	}
}

func TestGetContainerInstanceEmptyARN(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore := instanceStore(t, context)

	_, err := instanceStore.GetContainerInstance("")

	if err == nil {
		t.Error("Expected an error when arn is empty in GetContainerInstance")
	}
}

func TestGetContainerInstanceGetFails(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore := instanceStore(t, context)

	context.datastore.EXPECT().Get(context.instanceKey1).Return(nil, errors.New("Error when getting key"))

	_, err := instanceStore.GetContainerInstance(context.instance1.Detail.ContainerInstanceArn)

	if err == nil {
		t.Error("Expected an error when datastore get fails")
	}
}

func TestGetContainerInstanceGetNoResults(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore := instanceStore(t, context)

	context.datastore.EXPECT().Get(context.instanceKey1).Return(make(map[string]string), nil)

	instance, err := instanceStore.GetContainerInstance(context.instance1.Detail.ContainerInstanceArn)

	if err != nil {
		t.Error("Unexpected error when datastore get returns empty results")
	}

	if instance != nil {
		t.Error("Expected GetContainerInstance to return nil when get from datastore is empty")
	}
}

func TestGetContainerInstanceGetMultipleResults(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore := instanceStore(t, context)

	resp := map[string]string{
		containerInstanceARN1: context.compressedInstanceJSON1,
		containerInstanceARN2: context.compressedInstanceJSON2,
	}
	context.datastore.EXPECT().Get(context.instanceKey1).Return(resp, nil)

	_, err := instanceStore.GetContainerInstance(context.instance1.Detail.ContainerInstanceArn)

	if err == nil {
		t.Error("Expected an error when datastore get returns multiple results")
	}
}

func TestGetContainerInstanceGetInvalidJsonResult(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore := instanceStore(t, context)

	resp := map[string]string{
		containerInstanceARN1: "invalidJSON",
	}
	context.datastore.EXPECT().Get(context.instanceKey1).Return(resp, nil)

	_, err := instanceStore.GetContainerInstance(context.instance1.Detail.ContainerInstanceArn)

	if err == nil {
		t.Error("Expected an error when datastore get returns invalid JSON results")
	}
}

func TestGetContainerInstance(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore := instanceStore(t, context)

	resp := map[string]string{
		containerInstanceARN1: context.compressedInstanceJSON1,
	}
	context.datastore.EXPECT().Get(context.instanceKey1).Return(resp, nil)

	instance, err := instanceStore.GetContainerInstance(context.instance1.Detail.ContainerInstanceArn)

	if err != nil {
		t.Error("Unexpected error when getting an instance")
	}

	if !reflect.DeepEqual(*instance, context.instance1) {
		t.Error("Expected the returned instance to match the one returned from the datastore")
	}
}

func TestListContainerInstancesGetWithPrefixInvalidJson(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore := instanceStore(t, context)

	resp := map[string]string{
		containerInstanceARN1: "invalidJSON",
	}
	context.datastore.EXPECT().GetWithPrefix(instanceKeyPrefix).Return(resp, nil)

	_, err := instanceStore.ListContainerInstances()

	if err == nil {
		t.Error("Expected an error when datastore GetWithPrefix fails")
	}
}

func TestListContainerInstancesGetWithPrefixFails(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore := instanceStore(t, context)

	context.datastore.EXPECT().GetWithPrefix(instanceKeyPrefix).Return(nil, errors.New("GetWithPrefix failed"))
	_, err := instanceStore.ListContainerInstances()

	if err == nil {
		t.Error("Expected an error when datastore GetWithPrefix fails")
	}
}

func TestListContainerInstancesGetWithPrefixNoResults(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore := instanceStore(t, context)

	context.datastore.EXPECT().GetWithPrefix(instanceKeyPrefix).Return(make(map[string]string), nil)
	instances, err := instanceStore.ListContainerInstances()

	if err != nil {
		t.Error("Unexpected error when datastore GetWithPrefix returns empty")
	}

	if len(instances) > 0 {
		t.Error("Expected ListContainerInstances result to be empty when GetWithPrefix result is empty")
	}
}

func TestListContainerInstancesGetWithPrefixMultipleResults(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore := instanceStore(t, context)

	resp := map[string]string{
		containerInstanceARN1: context.compressedInstanceJSON1,
		containerInstanceARN2: context.compressedInstanceJSON2,
	}
	context.datastore.EXPECT().GetWithPrefix(instanceKeyPrefix).Return(resp, nil)
	instances, err := instanceStore.ListContainerInstances()

	if err != nil {
		t.Error("Unexpected error when datastore GetWithPrefix returns empty")
	}

	if len(instances) != len(resp) {
		t.Error("Expected ListContainerInstances result to be the same length as the GetWithPrefix result")
	}

	for _, v := range instances {
		value, ok := resp[v.Detail.ContainerInstanceArn]
		if !ok {
			t.Errorf("Expected GetWithPrefix result to contain the same elements as ListContainerInstances result. Missing %v", v)
		} else {
			instance := uncompressAndUnmarshalString(t, value)
			if !reflect.DeepEqual(v, instance) {
				t.Errorf("Expected GetWithPrefix result to contain the same elements as ListContainerInstances result. %v does not match %v", v, instance)
			}
		}
	}
}

func TestFilterContainerInstancesEmptyKey(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore := instanceStore(t, context)
	_, err := instanceStore.FilterContainerInstances("", "value")
	if err == nil {
		t.Error("Expected an error when filterKey is empty in FilterContainerInstances")
	}
}

func TestFilterContainerInstancesEmptyValue(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore := instanceStore(t, context)
	_, err := instanceStore.FilterContainerInstances(key, "")
	if err == nil {
		t.Error("Expected an error when filterValue is empty in FilterContainerInstances")
	}
}

func TestFilterContainerInstancesUnsupportedFilter(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore := instanceStore(t, context)
	_, err := instanceStore.FilterContainerInstances("invalidFilter", "value")
	if err == nil {
		t.Error("Expected an error when unsupported filter key is provided in FilterContainerInstances")
	}
}

func TestFilterContainerInstancesDatastoreGetWithPrefixFails(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	context.datastore.EXPECT().GetWithPrefix(instanceKeyPrefix).Return(nil, errors.New("GetWithPrefix failed"))

	instanceStore := instanceStore(t, context)
	_, err := instanceStore.FilterContainerInstances(instanceStatusFilter, status1)
	if err == nil {
		t.Error("Expected an error when datastore GetWithPrefix fails in FilterContainerInstances")
	}
}

func TestFilterContainerInstancesDatastoreGetWithPrefixReturnsNoResults(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	context.datastore.EXPECT().GetWithPrefix(instanceKeyPrefix).Return(make(map[string]string), nil)

	instanceStore := instanceStore(t, context)
	instances, err := instanceStore.FilterContainerInstances(instanceStatusFilter, status1)

	if err != nil {
		t.Error("Unexpected error when datastore GetWithPrefix returns empty map in FilterContainerInstances")
	}

	if instances == nil || len(instances) != 0 {
		t.Error("Result should be empty when datastore GetWithPrefix returns empty map in FilterContainerInstances")
	}
}

func TestFilterContainerInstancesNoResultsMatchStatusFilter(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	resp := map[string]string{
		containerInstanceARN1: context.compressedInstanceJSON1,
	}
	context.datastore.EXPECT().GetWithPrefix(instanceKeyPrefix).Return(resp, nil)

	instanceStore := instanceStore(t, context)
	instances, err := instanceStore.FilterContainerInstances(instanceStatusFilter, status2)

	if err != nil {
		t.Error("Unexpected error when datastore GetWithPrefix returns results in FilterContainerInstances")
	}

	if instances == nil || len(instances) != 0 {
		t.Error("Result should be empty when status filter does not match any results in FilterContainerInstances")
	}
}

func TestFilterContainerInstancesNoResultsMatchClusterFilter(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	resp := map[string]string{
		containerInstanceARN1: context.compressedInstanceJSON1,
		containerInstanceARN2: context.compressedInstanceJSON2,
	}
	context.datastore.EXPECT().GetWithPrefix(instanceKeyPrefix).Return(resp, nil)

	instanceStore := instanceStore(t, context)
	instances, err := instanceStore.FilterContainerInstances(clusterFilter, clusterName3)

	if err != nil {
		t.Error("Unexpected error when datastore GetWithPrefix returns results in FilterContainerInstances")
	}

	if instances == nil || len(instances) != 0 {
		t.Error("Result should be empty when cluster filter does not match any results in FilterContainerInstances")
	}
}

func TestFilterContainerInstancesMultipleResultsOneMatchesStatusFilter(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	resp := map[string]string{
		containerInstanceARN1: context.compressedInstanceJSON1,
		containerInstanceARN2: context.compressedInstanceJSON2,
	}
	context.datastore.EXPECT().GetWithPrefix(instanceKeyPrefix).Return(resp, nil)

	instanceStore := instanceStore(t, context)
	instances, err := instanceStore.FilterContainerInstances(instanceStatusFilter, status1)

	if err != nil {
		t.Error("Unexpected error when datastore GetWithPrefix returns results in FilterContainerInstances")
	}
	validateFilterContainerInstancesMultipleResultsOneMatchesFilterResult(t, instances, context.instance1)
}

func TestFilterContainerInstancesMultipleResultsOneMatchesClusterNameFilter(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	resp := map[string]string{
		containerInstanceARN1: context.compressedInstanceJSON1,
		containerInstanceARN2: context.compressedInstanceJSON2,
	}
	context.datastore.EXPECT().GetWithPrefix(instanceKeyPrefix).Return(resp, nil)

	instanceStore := instanceStore(t, context)
	instances, err := instanceStore.FilterContainerInstances(clusterFilter, clusterName1)

	if err != nil {
		t.Error("Unexpected error when datastore GetWithPrefix returns results in FilterContainerInstances")
	}
	validateFilterContainerInstancesMultipleResultsOneMatchesFilterResult(t, instances, context.instance1)
}

func TestFilterContainerInstancesMultipleResultsOneMatchesClusterArnFilter(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	resp := map[string]string{
		containerInstanceARN1: context.compressedInstanceJSON1,
		containerInstanceARN2: context.compressedInstanceJSON2,
	}
	context.datastore.EXPECT().GetWithPrefix(instanceKeyPrefix).Return(resp, nil)

	instanceStore := instanceStore(t, context)
	instances, err := instanceStore.FilterContainerInstances(clusterFilter, clusterARN1)

	if err != nil {
		t.Error("Unexpected error when datastore GetWithPrefix returns results in FilterContainerInstances")
	}
	validateFilterContainerInstancesMultipleResultsOneMatchesFilterResult(t, instances, context.instance1)
}

func validateFilterContainerInstancesMultipleResultsOneMatchesFilterResult(t *testing.T, instances []types.ContainerInstance, expectedInstance types.ContainerInstance) {
	if instances == nil || len(instances) != 1 {
		t.Error("Result should have 1 instance when 1 instance matches results in FilterContainerInstances")
	}

	if !reflect.DeepEqual(instances[0], expectedInstance) {
		t.Error("Expected the returned instance to match the instance with cluster name " + clusterName1)
	}
}

func TestFilterContainerInstancesMultipleResultsMatchStatusFilter(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instance := context.instance2
	instance.Detail.Status = status1
	instanceJSON := marshalInstance(t, instance)
	compressedJSON := compressInstanceJSON(t, instanceJSON)

	resp := map[string]string{
		containerInstanceARN1: context.compressedInstanceJSON1,
		containerInstanceARN2: compressedJSON,
	}
	context.datastore.EXPECT().GetWithPrefix(instanceKeyPrefix).Return(resp, nil)

	instanceStore := instanceStore(t, context)
	instances, err := instanceStore.FilterContainerInstances(instanceStatusFilter, status1)

	if err != nil {
		t.Error("Unexpected error when datastore GetWithPrefix returns results in FilterContainerInstances")
	}
	validateFilterContainerInstancecMultipleResultsMatchFilterResult(t, instances, resp)
}

func TestFilterContainerInstancesMultipleResultsMatchClusterNameFilter(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instance := context.instance2
	instance.Detail.ClusterArn = clusterARN1
	instanceJSON := marshalInstance(t, instance)
	compressedJSON := compressInstanceJSON(t, instanceJSON)

	resp := map[string]string{
		containerInstanceARN1: context.compressedInstanceJSON1,
		containerInstanceARN2: compressedJSON,
	}
	context.datastore.EXPECT().GetWithPrefix(instanceKeyPrefix).Return(resp, nil)

	instanceStore := instanceStore(t, context)
	instances, err := instanceStore.FilterContainerInstances(clusterFilter, clusterName1)

	if err != nil {
		t.Error("Unexpected error when datastore GetWithPrefix returns results in FilterContainerInstances")
	}
	validateFilterContainerInstancecMultipleResultsMatchFilterResult(t, instances, resp)
}

func TestFilterContainerInstancesMultipleResultsMatchClusterArnFilter(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instance := context.instance2
	instance.Detail.ClusterArn = clusterARN1
	instanceJSON := marshalInstance(t, instance)
	compressedJSON := compressInstanceJSON(t, instanceJSON)

	resp := map[string]string{
		containerInstanceARN1: context.compressedInstanceJSON1,
		containerInstanceARN2: compressedJSON,
	}
	context.datastore.EXPECT().GetWithPrefix(instanceKeyPrefix).Return(resp, nil)

	instanceStore := instanceStore(t, context)
	instances, err := instanceStore.FilterContainerInstances(clusterFilter, clusterARN1)

	if err != nil {
		t.Error("Unexpected error when datastore GetWithPrefix returns results in FilterContainerInstances")
	}
	validateFilterContainerInstancecMultipleResultsMatchFilterResult(t, instances, resp)
}

func validateFilterContainerInstancecMultipleResultsMatchFilterResult(t *testing.T, instances []types.ContainerInstance, datastoreResp map[string]string) {
	if instances == nil || len(instances) != 2 {
		t.Error("Result should have 2 instances when 2 instances match results in FilterContainerInstances")
	}

	for _, v := range instances {
		value, ok := datastoreResp[v.Detail.ContainerInstanceArn]
		if !ok {
			t.Errorf("Expected FilterContainerInstances result to contain the same elements as datastore GetWithPrefix result. Missing %v", v)
		} else {
			instance := uncompressAndUnmarshalString(t, value)
			if !reflect.DeepEqual(v, instance) {
				t.Errorf("Expected FilterContainerInstances result to contain the same elements as GetWithPrefix result. %v does not match %v", v, instance)
			}
		}
	}
}

func TestStreamContainerInstancesDataStoreStreamReturnsError(t *testing.T) {
	ctx := NewContainerInstanceStoreMockContext(t)
	defer ctx.mockCtrl.Finish()

	tstCtx := context.Background()
	ctx.datastore.EXPECT().StreamWithPrefix(gomock.Any(), instanceKeyPrefix).Return(nil, errors.New("StreamWithPrefix failed"))

	instanceStore := instanceStore(t, ctx)
	instaceRespChan, err := instanceStore.StreamContainerInstances(tstCtx)
	if err == nil {
		t.Error("Expected an error when datastore StreamWithPrefix returns an error")
	}
	if instaceRespChan != nil {
		t.Error("Unexpected instance response channel when there is a datastore channel setup error")
	}
}

func TestStreamContainerInstancesValidJSONInDSChannel(t *testing.T) {
	ctx := NewContainerInstanceStoreMockContext(t)
	defer ctx.mockCtrl.Finish()

	tstCtx := context.Background()
	dsChan := make(chan map[string]string)
	defer close(dsChan)
	ctx.datastore.EXPECT().StreamWithPrefix(gomock.Any(), instanceKeyPrefix).Return(dsChan, nil)

	instanceStore := instanceStore(t, ctx)
	instanceRespChan, err := instanceStore.StreamContainerInstances(tstCtx)
	if err != nil {
		t.Error("Unexpected error when calling stream instances")
	}
	if instanceRespChan == nil {
		t.Error("Expected valid non-nil instanceRespChannel")
	}

	instanceResp := addContainerInstanceToDSChanAndReadFromInstanceRespChan(ctx.compressedInstanceJSON1, dsChan, instanceRespChan)
	if instanceResp.Err != nil {
		t.Error("Unexpected error when reading instance from channel")
	}
	if !reflect.DeepEqual(ctx.instance1, instanceResp.ContainerInstance) {
		t.Error("Expected instance in instance response to match that in the stream")
	}
}

func TestStreamContainerInstancesInvalidJSONInDSChannel(t *testing.T) {
	ctx := NewContainerInstanceStoreMockContext(t)
	defer ctx.mockCtrl.Finish()

	tstCtx := context.Background()
	dsChan := make(chan map[string]string)
	defer close(dsChan)
	ctx.datastore.EXPECT().StreamWithPrefix(gomock.Any(), instanceKeyPrefix).Return(dsChan, nil)

	instanceStore := instanceStore(t, ctx)
	instanceRespChan, err := instanceStore.StreamContainerInstances(tstCtx)
	if err != nil {
		t.Error("Unexpected error when calling stream instances")
	}
	if instanceRespChan == nil {
		t.Error("Expected valid non-nil instanceRespChannel")
	}

	compressedInvalidJSON := compressInstanceJSON(t, "invalidJSON")
	instanceResp := addContainerInstanceToDSChanAndReadFromInstanceRespChan(compressedInvalidJSON, dsChan, instanceRespChan)

	if instanceResp.Err == nil {
		t.Error("Expected an error when dsChannel returns an invalid instance json")
	}
	if !reflect.DeepEqual(types.ContainerInstance{}, instanceResp.ContainerInstance) {
		t.Error("Expected empty instance in response when there is a decode error")
	}

	_, ok := <-instanceRespChan
	if ok {
		t.Error("Expected instance response channel to be closed")
	}
}

func TestStreamContainerInstancesCancelUpstreamContext(t *testing.T) {
	ctx := NewContainerInstanceStoreMockContext(t)
	defer ctx.mockCtrl.Finish()

	tstCtx, cancel := context.WithCancel(context.Background())
	dsChan := make(chan map[string]string)
	defer close(dsChan)
	ctx.datastore.EXPECT().StreamWithPrefix(gomock.Any(), instanceKeyPrefix).Return(dsChan, nil)

	instanceStore := instanceStore(t, ctx)
	instanceRespChan, err := instanceStore.StreamContainerInstances(tstCtx)
	if err != nil {
		t.Error("Unexpected error when calling stream instances")
	}
	if instanceRespChan == nil {
		t.Error("Expected valid non-nil instanceRespChannel")
	}

	cancel()

	_, ok := <-instanceRespChan
	if ok {
		t.Error("Expected instance response channel to be closed")
	}
}

func TestStreamContainerInstancesCloseDownstreamChannel(t *testing.T) {
	ctx := NewContainerInstanceStoreMockContext(t)
	defer ctx.mockCtrl.Finish()

	tstCtx := context.Background()
	dsChan := make(chan map[string]string)
	ctx.datastore.EXPECT().StreamWithPrefix(gomock.Any(), instanceKeyPrefix).Return(dsChan, nil)

	instanceStore := instanceStore(t, ctx)
	instanceRespChan, err := instanceStore.StreamContainerInstances(tstCtx)
	if err != nil {
		t.Error("Unexpected error when calling stream instances")
	}
	if instanceRespChan == nil {
		t.Error("Expected valid non-nil instanceRespChannel")
	}

	close(dsChan)

	_, ok := <-instanceRespChan
	if ok {
		t.Error("Expected instance response channel to be closed")
	}
}

func instanceStore(t *testing.T, context *instanceStoreMockContext) ContainerInstanceStore {
	instanceStore, err := NewContainerInstanceStore(context.datastore)
	if err != nil {
		t.Error("Unexpected error when calling NewContainerInstanceStore")
	}
	return instanceStore
}

func marshalInstance(t *testing.T, instance types.ContainerInstance) string {
	instanceJSON, err := json.Marshal(instance)
	if err != nil {
		t.Error("Failed to marshal instance: ", err)
	}
	return string(instanceJSON)
}

func compressInstanceJSON(t *testing.T, instanceJSON string) string {
	compressedJSON, err := compress.Compress(instanceJSON)
	if err != nil {
		t.Error("Failed to compress instanceJSON: ", err)
	}
	return string(compressedJSON)
}

func uncompressAndUnmarshalString(t *testing.T, str string) types.ContainerInstance {
	var instance types.ContainerInstance
	uncompressedStr, err := compress.Uncompress([]byte(str))
	if err != nil {
		t.Error("Failed to uncompress string: ", err)
	}
	err = json.Unmarshal([]byte(uncompressedStr), &instance)
	if err != nil {
		t.Error("Failed to unmarshal compressed string: ", err)
	}
	return instance
}

func addContainerInstanceToDSChanAndReadFromInstanceRespChan(instanceToAdd string, dsChan chan map[string]string, instanceRespChan chan storetypes.ContainerInstanceErrorWrapper) storetypes.ContainerInstanceErrorWrapper {
	var instanceResp storetypes.ContainerInstanceErrorWrapper

	doneChan := make(chan bool)
	defer close(doneChan)
	go func() {
		instanceResp = <-instanceRespChan
		doneChan <- true
	}()

	dsVal := map[string]string{containerInstanceARN1: instanceToAdd}
	dsChan <- dsVal
	<-doneChan

	return instanceResp
}
