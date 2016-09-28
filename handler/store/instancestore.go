package store

import (
	"github.com/aws/amazon-ecs-event-stream-handler/handler/json"
	"github.com/aws/amazon-ecs-event-stream-handler/handler/types"
	log "github.com/cihub/seelog"
	"github.com/pkg/errors"
)

const (
	instanceKeyPrefix = "ecs/instance/"
)

// ContainerInstanceStore defines methods to access container instances from the datastore
type ContainerInstanceStore interface {
	AddContainerInstance(instance string) error
	GetContainerInstance(arn string) (*types.ContainerInstance, error)
	ListContainerInstances() ([]types.ContainerInstance, error)
	FilterContainerInstances(filterKey string, filterValue string) ([]types.ContainerInstance, error)
	StreamContainerInstances() ([]types.ContainerInstance, error)
}

type eventInstanceStore struct {
	datastore DataStore
}

func NewContainerInstanceStore(ds DataStore) (ContainerInstanceStore, error) {
	if ds == nil {
		return nil, errors.New("The datastore cannot be nil")
	}

	return eventInstanceStore{
		datastore: ds,
	}, nil
}

func generateInstanceKey(instance types.ContainerInstance) (string, error) {
	if len(instance.Detail.ContainerInstanceArn) == 0 {
		return "", errors.New("Container instance arn cannot be empty")
	}
	return instanceKeyPrefix + instance.Detail.ContainerInstanceArn, nil
}

// AddContainerInstance adds a container instance represented in the instanceJSON to the datastore
func (instanceStore eventInstanceStore) AddContainerInstance(instanceJSON string) error {
	if len(instanceJSON) == 0 {
		return errors.New("Instance json should not be empty")
	}

	var instance types.ContainerInstance
	err := json.UnmarshalJSON(instanceJSON, &instance)
	if err != nil {
		return err
	}

	key, err := generateInstanceKey(instance)
	if err != nil {
		return err
	}

	// check if record exists with higher version number
	existingInstance, err := instanceStore.getInstanceByKey(key)
	if err != nil {
		return err
	}

	if existingInstance != nil {
		if existingInstance.Detail.Version >= instance.Detail.Version {
			log.Infof("Higher or equal version %v of instance %v with version %v already exists",
				existingInstance.Detail.Version,
				instance.Detail.ContainerInstanceArn,
				instance.Detail.Version)

			// do nothing. later version of the event has already been stored
			return nil
		}
	}

	err = instanceStore.datastore.Add(key, instanceJSON)
	if err != nil {
		return err
	}

	return nil
}

// GetContainerInstance gets a container instance with key 'arn' from the datastore
func (instanceStore eventInstanceStore) GetContainerInstance(arn string) (*types.ContainerInstance, error) {
	if len(arn) == 0 {
		return nil, errors.New("Arn should not be empty")
	}

	var instance types.ContainerInstance
	instance.Detail.ContainerInstanceArn = arn

	key, err := generateInstanceKey(instance)
	if err != nil {
		return nil, err
	}

	return instanceStore.getInstanceByKey(key)
}

// ListContainerInstances lists all container instances existing in the datastore
func (instanceStore eventInstanceStore) ListContainerInstances() ([]types.ContainerInstance, error) {
	return instanceStore.getInstancesByKeyPrefix(instanceKeyPrefix)
}

// FilterContainerInstances returns all container instances from the datastore that match the provided filters
func (instanceStore eventInstanceStore) FilterContainerInstances(filterKey string, filterValue string) ([]types.ContainerInstance, error) {
	return nil, nil
}

// StreamContainerInstances returns a stream of all changes in the container instance keyspace
func (instanceStore eventInstanceStore) StreamContainerInstances() ([]types.ContainerInstance, error) {
	return nil, nil
}

func (instanceStore eventInstanceStore) getInstanceByKey(key string) (*types.ContainerInstance, error) {
	if len(key) == 0 {
		return nil, errors.New("Key cannot be empty")
	}

	resp, err := instanceStore.datastore.Get(key)
	if err != nil {
		return nil, err
	}

	if len(resp) == 0 {
		return nil, nil
	}

	if len(resp) > 1 {
		return nil, errors.Errorf("Multiple entries exist in the datastore with key %v", key)
	}

	var instance types.ContainerInstance
	for _, v := range resp {
		err = json.UnmarshalJSON(v, &instance)
		if err != nil {
			return nil, err
		}
		break
	}
	return &instance, nil
}

func (instanceStore eventInstanceStore) getInstancesByKeyPrefix(key string) ([]types.ContainerInstance, error) {
	if len(key) == 0 {
		return nil, errors.New("Key cannot be empty")
	}

	resp, err := instanceStore.datastore.GetWithPrefix(key)
	if err != nil {
		return nil, err
	}

	if len(resp) == 0 {
		return make([]types.ContainerInstance, 0), nil
	}

	instances := []types.ContainerInstance{}
	for _, v := range resp {
		var instance types.ContainerInstance
		err = json.UnmarshalJSON(v, &instance)
		if err != nil {
			return nil, err
		}
		instances = append(instances, instance)
	}
	return instances, nil
}
