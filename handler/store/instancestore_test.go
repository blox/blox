package store

import (
	"github.com/aws/amazon-ecs-event-stream-handler/handler/mocks"
	"github.com/aws/amazon-ecs-event-stream-handler/handler/types"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"reflect"
	"testing"
)

const (
	containerInstanceArn = "arn:aws:us-east-1:123456789123:container-instance/4b6d45ea-a4b4-4269-9d04-3af6ddfdc597"
)

type instanceStoreMockContext struct {
	mockCtrl     *gomock.Controller
	datastore    *mocks.MockDataStore
	instance     types.ContainerInstance
	instanceJson string
	instanceKey  string
}

func NewContainerInstanceStoreMockContext(t *testing.T) *instanceStoreMockContext {
	context := instanceStoreMockContext{}
	context.mockCtrl = gomock.NewController(t)
	context.datastore = mocks.NewMockDataStore(context.mockCtrl)

	context.instance = types.ContainerInstance{}
	context.instance.Detail.ContainerInstanceArn = containerInstanceArn
	context.instance.Detail.Version = 1

	json, err := json.Marshal(context.instance)
	if err != nil {
		t.Error("Cannot initialize instanceStoreMockContext")
	}
	context.instanceJson = string(json)

	context.instanceKey = instanceKeyPrefix + containerInstanceArn

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

	instanceStore, err := NewContainerInstanceStore(context.datastore)

	if err != nil {
		t.Error("Unexpected error when calling NewContainerInstanceStore")
	}

	if instanceStore == nil {
		t.Error("Instancestore should not be nil")
	}
}

func TestAddContainerInstanceEmptyInstanceJson(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore, err := NewContainerInstanceStore(context.datastore)

	if err != nil {
		t.Error("Unexpected error when calling NewContainerInstanceStore")
	}

	err = instanceStore.AddContainerInstance("")

	if err == nil {
		t.Error("Expected an error when instance json is empty in AddContainerInstance")
	}
}

func TestAddContainerInstanceJsonUnmarshalError(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore, err := NewContainerInstanceStore(context.datastore)

	if err != nil {
		t.Error("Unexpected error when calling NewContainerInstanceStore")
	}

	err = instanceStore.AddContainerInstance("invalidJson")

	if err == nil {
		t.Error("Expected an error when instance json is invalid in AddContainerInstance")
	}
}

func TestAddContainerInstanceEmptyContainerInstanceArn(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore, err := NewContainerInstanceStore(context.datastore)

	if err != nil {
		t.Error("Unexpected error when calling NewContainerInstanceStore")
	}

	instance := types.ContainerInstance{}

	instanceJson, err := json.Marshal(instance)
	err = instanceStore.AddContainerInstance(string(instanceJson))

	if err == nil {
		t.Error("Expected an error when container instance arn is not set")
	}
}

func TestAddContainerInstanceGetContainerInstanceFails(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore, err := NewContainerInstanceStore(context.datastore)

	if err != nil {
		t.Error("Unexpected error when calling NewContainerInstanceStore")
	}

	context.datastore.EXPECT().Get(context.instanceKey).Return(nil, errors.New("Error when getting key"))

	err = instanceStore.AddContainerInstance(context.instanceJson)

	if err == nil {
		t.Error("Expected an error when datastore get fails")
	}
}

func TestAddContainerInstanceGetContainerInstanceNoResults(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore, err := NewContainerInstanceStore(context.datastore)

	if err != nil {
		t.Error("Unexpected error when calling NewContainerInstanceStore")
	}

	context.datastore.EXPECT().Get(context.instanceKey).Return(make(map[string]string), nil)
	context.datastore.EXPECT().Add(context.instanceKey, context.instanceJson).Return(nil)

	err = instanceStore.AddContainerInstance(context.instanceJson)

	if err != nil {
		t.Error("Unexpected error when datastore get returns empty results")
	}
}

func TestAddContainerInstanceGetContainerInstanceMultipleResults(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore, err := NewContainerInstanceStore(context.datastore)

	if err != nil {
		t.Error("Unexpected error when calling NewContainerInstanceStore")
	}

	resp := map[string]string{
		containerInstanceArn:       "result0",
		containerInstanceArn + "1": "result1",
	}
	context.datastore.EXPECT().Get(context.instanceKey).Return(resp, nil)

	err = instanceStore.AddContainerInstance(context.instanceJson)

	if err == nil {
		t.Error("Expected an error when datastore get returns multiple results")
	}
}

func TestAddContainerInstanceGetContainerInstanceInvalidJsonResult(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore, err := NewContainerInstanceStore(context.datastore)

	if err != nil {
		t.Error("Unexpected error when calling NewContainerInstanceStore")
	}

	resp := map[string]string{
		containerInstanceArn: "invalidJson",
	}
	context.datastore.EXPECT().Get(context.instanceKey).Return(resp, nil)

	err = instanceStore.AddContainerInstance(context.instanceJson)

	if err == nil {
		t.Error("Expected an error when datastore get returns invalid json")
	}
}

func TestAddContainerInstanceSameVersionInstanceExists(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore, err := NewContainerInstanceStore(context.datastore)

	if err != nil {
		t.Error("Unexpected error when calling NewContainerInstanceStore")
	}

	resp := map[string]string{
		containerInstanceArn: context.instanceJson,
	}
	context.datastore.EXPECT().Get(context.instanceKey).Return(resp, nil)
	context.datastore.EXPECT().Add(context.instanceKey, context.instanceJson).Times(0)

	err = instanceStore.AddContainerInstance(context.instanceJson)

	if err != nil {
		t.Error("Unxpected error when adding instance and same version instance exists")
	}
}

func TestAddContainerInstanceHigherVersionInstanceExists(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore, err := NewContainerInstanceStore(context.datastore)

	if err != nil {
		t.Error("Unexpected error when calling NewContainerInstanceStore")
	}

	instance := types.ContainerInstance{}
	instance.Detail.ContainerInstanceArn = containerInstanceArn
	instance.Detail.Version = context.instance.Detail.Version + 1

	json, _ := json.Marshal(instance)

	resp := map[string]string{
		containerInstanceArn: string(json),
	}
	context.datastore.EXPECT().Get(context.instanceKey).Return(resp, nil)
	context.datastore.EXPECT().Add(context.instanceKey, context.instanceJson).Times(0)

	err = instanceStore.AddContainerInstance(context.instanceJson)

	if err != nil {
		t.Error("Unxpected error when adding instance and higher version instance exists")
	}
}

func TestAddContainerInstanceLowerVersionInstanceExists(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore, err := NewContainerInstanceStore(context.datastore)

	if err != nil {
		t.Error("Unexpected error when calling NewContainerInstanceStore")
	}

	instance := types.ContainerInstance{}
	instance.Detail.ContainerInstanceArn = containerInstanceArn
	instance.Detail.Version = context.instance.Detail.Version - 1

	json, _ := json.Marshal(instance)

	resp := map[string]string{
		containerInstanceArn: string(json),
	}
	context.datastore.EXPECT().Get(context.instanceKey).Return(resp, nil)
	context.datastore.EXPECT().Add(context.instanceKey, context.instanceJson).Return(nil)

	err = instanceStore.AddContainerInstance(context.instanceJson)

	if err != nil {
		t.Error("Unxpected error when adding instance and higher version instance exists")
	}
}

func TestAddContainerInstanceFails(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore, err := NewContainerInstanceStore(context.datastore)

	if err != nil {
		t.Error("Unexpected error when calling NewContainerInstanceStore")
	}

	instance := types.ContainerInstance{}
	instance.Detail.ContainerInstanceArn = containerInstanceArn
	instance.Detail.Version = context.instance.Detail.Version - 1

	json, _ := json.Marshal(instance)

	resp := map[string]string{
		containerInstanceArn: string(json),
	}
	context.datastore.EXPECT().Get(context.instanceKey).Return(resp, nil)
	context.datastore.EXPECT().Add(context.instanceKey, context.instanceJson).Return(errors.New("Add instance failed"))

	err = instanceStore.AddContainerInstance(context.instanceJson)

	if err == nil {
		t.Error("Expected an error when adding an instance fails")
	}
}

func TestGetContainerInstanceEmptyArn(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore, err := NewContainerInstanceStore(context.datastore)

	if err != nil {
		t.Error("Unexpected error when calling NewContainerInstanceStore")
	}

	_, err = instanceStore.GetContainerInstance("")

	if err == nil {
		t.Error("Expected an error when arn is empty in GetContainerInstance")
	}
}

func TestGetContainerInstanceGetFails(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore, err := NewContainerInstanceStore(context.datastore)

	if err != nil {
		t.Error("Unexpected error when calling NewContainerInstanceStore")
	}

	context.datastore.EXPECT().Get(context.instanceKey).Return(nil, errors.New("Error when getting key"))

	_, err = instanceStore.GetContainerInstance(context.instance.Detail.ContainerInstanceArn)

	if err == nil {
		t.Error("Expected an error when datastore get fails")
	}
}

func TestGetContainerInstanceGetNoResults(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore, err := NewContainerInstanceStore(context.datastore)

	if err != nil {
		t.Error("Unexpected error when calling NewContainerInstanceStore")
	}

	context.datastore.EXPECT().Get(context.instanceKey).Return(make(map[string]string), nil)

	instance, err := instanceStore.GetContainerInstance(context.instance.Detail.ContainerInstanceArn)

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

	instanceStore, err := NewContainerInstanceStore(context.datastore)

	if err != nil {
		t.Error("Unexpected error when calling NewContainerInstanceStore")
	}

	resp := map[string]string{
		containerInstanceArn:       "result0",
		containerInstanceArn + "1": "result1",
	}
	context.datastore.EXPECT().Get(context.instanceKey).Return(resp, nil)

	_, err = instanceStore.GetContainerInstance(context.instance.Detail.ContainerInstanceArn)

	if err == nil {
		t.Error("Expected an error when datastore get returns multiple results")
	}
}

func TestGetContainerInstanceGetInvalidJsonResult(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore, err := NewContainerInstanceStore(context.datastore)

	if err != nil {
		t.Error("Unexpected error when calling NewContainerInstanceStore")
	}

	resp := map[string]string{
		containerInstanceArn: "invalidJson",
	}
	context.datastore.EXPECT().Get(context.instanceKey).Return(resp, nil)

	_, err = instanceStore.GetContainerInstance(context.instance.Detail.ContainerInstanceArn)

	if err == nil {
		t.Error("Expected an error when datastore get returns invalid json results")
	}
}

func TestGetContainerInstance(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore, err := NewContainerInstanceStore(context.datastore)

	if err != nil {
		t.Error("Unexpected error when calling NewContainerInstanceStore")
	}

	resp := map[string]string{
		containerInstanceArn: context.instanceJson,
	}
	context.datastore.EXPECT().Get(context.instanceKey).Return(resp, nil)

	instance, err := instanceStore.GetContainerInstance(context.instance.Detail.ContainerInstanceArn)

	if err != nil {
		t.Error("Unexpected error when getting an instance")
	}

	if !reflect.DeepEqual(*instance, context.instance) {
		t.Error("Expected the returned instance to match the one returned from the datastore")
	}
}

func TestListContainerInstancesGetWithPrefixInvalidJson(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore, err := NewContainerInstanceStore(context.datastore)

	if err != nil {
		t.Error("Unexpected error when calling NewContainerInstanceStore")
	}

	resp := map[string]string{
		containerInstanceArn: "invalidJson",
	}
	context.datastore.EXPECT().GetWithPrefix(instanceKeyPrefix).Return(resp, nil)

	_, err = instanceStore.ListContainerInstances()

	if err == nil {
		t.Error("Expected an error when datastore GetWithPrefix fails")
	}
}

func TestListContainerInstancesGetWithPrefixFails(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore, err := NewContainerInstanceStore(context.datastore)

	if err != nil {
		t.Error("Unexpected error when calling NewContainerInstanceStore")
	}

	context.datastore.EXPECT().GetWithPrefix(instanceKeyPrefix).Return(nil, errors.New("GetWithPrefix failed"))
	_, err = instanceStore.ListContainerInstances()

	if err == nil {
		t.Error("Expected an error when datastore GetWithPrefix fails")
	}
}

func TestListContainerInstancesGetWithPrefixNoResults(t *testing.T) {
	context := NewContainerInstanceStoreMockContext(t)
	defer context.mockCtrl.Finish()

	instanceStore, err := NewContainerInstanceStore(context.datastore)

	if err != nil {
		t.Error("Unexpected error when calling NewContainerInstanceStore")
	}

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

	instanceStore, err := NewContainerInstanceStore(context.datastore)

	if err != nil {
		t.Error("Unexpected error when calling NewContainerInstanceStore")
	}

	instance := types.ContainerInstance{}
	instance.Detail.ContainerInstanceArn = containerInstanceArn + "1"

	instance2Json, _ := json.Marshal(instance)

	resp := map[string]string{
		containerInstanceArn:       context.instanceJson,
		containerInstanceArn + "1": string(instance2Json),
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
			var valueInstance types.ContainerInstance
			json.Unmarshal([]byte(value), &valueInstance)
			if !reflect.DeepEqual(v, valueInstance) {
				t.Errorf("Expected GetWithPrefix result to contain the same elements as ListContainerInstances result. %v does not match %v", v, valueInstance)
			}
		}
	}
}
