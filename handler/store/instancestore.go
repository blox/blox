package store

import (
	"context"
	"strings"

	"github.com/aws/amazon-ecs-event-stream-handler/handler/compress"
	"github.com/aws/amazon-ecs-event-stream-handler/handler/json"
	"github.com/aws/amazon-ecs-event-stream-handler/handler/regex"
	storetypes "github.com/aws/amazon-ecs-event-stream-handler/handler/store/types"
	"github.com/aws/amazon-ecs-event-stream-handler/handler/types"
	log "github.com/cihub/seelog"
	"github.com/pkg/errors"
)

const (
	instanceKeyPrefix    = "ecs/instance/"
	instanceStatusFilter = "status"
	clusterFilter        = "cluster"
)

// ContainerInstanceStore defines methods to access container instances from the datastore
type ContainerInstanceStore interface {
	AddContainerInstance(instance string) error
	GetContainerInstance(cluster string, instanceARN string) (*types.ContainerInstance, error)
	ListContainerInstances() ([]types.ContainerInstance, error)
	FilterContainerInstances(filterKey string, filterValue string) ([]types.ContainerInstance, error)
	StreamContainerInstances(ctx context.Context) (chan storetypes.ContainerInstanceErrorWrapper, error)
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

func generateInstanceKey(clusterName string, instanceARN string) (string, error) {
	if len(clusterName) == 0 {
		return "", errors.New("Cluster name cannot be empty")
	}
	if len(instanceARN) == 0 {
		return "", errors.New("Instance ARN cannot be empty")
	}
	return instanceKeyPrefix + clusterName + "/" + instanceARN, nil
}

// AddContainerInstance adds a container instance represented in the instanceJSON to the datastore
func (instanceStore eventInstanceStore) AddContainerInstance(instanceJSON string) error {
	if len(instanceJSON) == 0 {
		return errors.New("Instance JSON should not be empty")
	}

	var instance types.ContainerInstance
	err := json.UnmarshalJSON(instanceJSON, &instance)
	if err != nil {
		return err
	}

	if instance.Detail == nil || instance.Detail.ClusterArn == nil || instance.Detail.ContainerInstanceArn == nil {
		return errors.New("Cluster ARN and container instance ARN should not be empty in instance JSON")
	}

	clusterName, err := regex.GetClusterNameFromARN(*instance.Detail.ClusterArn)
	if err != nil {
		return err
	}

	key, err := generateInstanceKey(clusterName, *instance.Detail.ContainerInstanceArn)
	if err != nil {
		return err
	}

	// check if record exists with higher version number
	existingInstance, err := instanceStore.getInstanceByKey(key)
	if err != nil {
		return err
	}

	if existingInstance != nil {
		existingInstanceDetail := *existingInstance.Detail
		currentInstanceDetail := *instance.Detail
		if *existingInstanceDetail.Version >= *currentInstanceDetail.Version {
			log.Infof("Higher or equal version %v of instance %v with version %v already exists",
				existingInstance.Detail.Version,
				instance.Detail.ContainerInstanceArn,
				instance.Detail.Version)

			// do nothing. later version of the event has already been stored
			return nil
		}
	}

	compressedInstanceJSON, err := compress.Compress(instanceJSON)
	if err != nil {
		return err
	}

	err = instanceStore.datastore.Add(key, string(compressedInstanceJSON))
	if err != nil {
		return err
	}

	return nil
}

// GetContainerInstance gets a container with ARN 'instanceARN' belonging to cluster 'cluster'
func (instanceStore eventInstanceStore) GetContainerInstance(cluster string, instanceARN string) (*types.ContainerInstance, error) {
	if len(cluster) == 0 {
		return nil, errors.New("Cluster should not be empty")
	}
	if len(instanceARN) == 0 {
		return nil, errors.New("Instance ARN should not be empty")
	}

	clusterName := cluster
	var err error
	if regex.IsClusterARN(cluster) {
		clusterName, err = regex.GetClusterNameFromARN(cluster)
		if err != nil {
			return nil, err
		}
	}

	key, err := generateInstanceKey(clusterName, instanceARN)
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
	if len(filterKey) == 0 || len(filterValue) == 0 {
		return nil, errors.New("Filter key and value cannot be empty")
	}

	switch {
	case filterKey == instanceStatusFilter:
		return instanceStore.filterContainerInstancesByStatus(filterValue)
	case filterKey == clusterFilter:
		return instanceStore.filterContainerInstancesByCluster(filterValue)
	default:
		return nil, errors.New("Unsupported filter key")
	}
}

func (instanceStore eventInstanceStore) filterContainerInstancesByStatus(status string) ([]types.ContainerInstance, error) {
	instances, err := instanceStore.ListContainerInstances()
	if err != nil {
		return nil, err
	}
	filteredInstances := make([]types.ContainerInstance, 0, len(instances))
	for _, instance := range instances {
		if strings.ToLower(status) == strings.ToLower(*instance.Detail.Status) {
			filteredInstances = append(filteredInstances, instance)
		}
	}
	return filteredInstances, nil
}

func (instanceStore eventInstanceStore) filterContainerInstancesByCluster(cluster string) ([]types.ContainerInstance, error) {
	clusterName := cluster
	var err error
	if regex.IsClusterARN(cluster) {
		clusterName, err = regex.GetClusterNameFromARN(cluster)
		if err != nil {
			return nil, err
		}
	}

	instancesForClusterPrefix := instanceKeyPrefix + clusterName + "/"
	return instanceStore.getInstancesByKeyPrefix(instancesForClusterPrefix)
}

// StreamContainerInstances returns a stream of all changes in the container instance keyspace
func (instanceStore eventInstanceStore) StreamContainerInstances(ctx context.Context) (chan storetypes.ContainerInstanceErrorWrapper, error) {
	instanceStoreCtx, cancel := context.WithCancel(ctx) // go routine instanceStore.pipeBetweenChannels() handles canceling this context

	dsChan, err := instanceStore.datastore.StreamWithPrefix(instanceStoreCtx, instanceKeyPrefix)
	if err != nil {
		cancel()
		return nil, err
	}

	instanceRespChan := make(chan storetypes.ContainerInstanceErrorWrapper) // go routine instanceStore.pipeBetweenChannels() handles closing of this channel
	go instanceStore.pipeBetweenChannels(instanceStoreCtx, cancel, dsChan, instanceRespChan)
	return instanceRespChan, nil
}

func (instanceStore eventInstanceStore) pipeBetweenChannels(ctx context.Context, cancel context.CancelFunc, dsChan chan map[string]string, instanceRespChan chan storetypes.ContainerInstanceErrorWrapper) {
	defer close(instanceRespChan)
	defer cancel()

	for {
		select {
		case resp, ok := <-dsChan:
			if !ok {
				return
			}
			for _, v := range resp {
				ins, err := instanceStore.uncompressAndUnmarshalInstance(v)
				if err != nil {
					instanceRespChan <- storetypes.ContainerInstanceErrorWrapper{ContainerInstance: types.ContainerInstance{}, Err: err}
					return
				}
				instanceRespChan <- storetypes.ContainerInstanceErrorWrapper{ContainerInstance: ins, Err: nil}
			}

		case <-ctx.Done():
			return
		}
	}
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
		uncompressedVal, err := compress.Uncompress([]byte(v))
		if err != nil {
			return nil, err
		}
		err = json.UnmarshalJSON(uncompressedVal, &instance)
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
		uncompressedVal, err := compress.Uncompress([]byte(v))
		if err != nil {
			return nil, err
		}

		var instance types.ContainerInstance
		err = json.UnmarshalJSON(uncompressedVal, &instance)
		if err != nil {
			return nil, err
		}
		instances = append(instances, instance)
	}
	return instances, nil
}

func (instanceStore eventInstanceStore) uncompressAndUnmarshalInstance(val string) (types.ContainerInstance, error) {
	var instance types.ContainerInstance

	uncompressedVal, err := compress.Uncompress([]byte(val))
	if err != nil {
		return instance, err
	}
	err = json.UnmarshalJSON(uncompressedVal, &instance)
	if err != nil {
		return instance, err
	}

	return instance, nil
}
