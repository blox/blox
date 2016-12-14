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

package store

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/blox/blox/cluster-state-service/handler/regex"
	storetypes "github.com/blox/blox/cluster-state-service/handler/store/types"
	"github.com/blox/blox/cluster-state-service/handler/types"
	log "github.com/cihub/seelog"
	"github.com/pkg/errors"
)

const (
	instanceKeyPrefix     = "ecs/instance/"
	instanceStatusFilter  = "status"
	instanceClusterFilter = "cluster"

	unversionedInstance = -1
)

// ContainerInstanceStore defines methods to access container instances from the datastore
type ContainerInstanceStore interface {
	AddContainerInstance(instance string) error
	AddUnversionedContainerInstance(instance string) error
	GetContainerInstance(cluster string, instanceARN string) (*types.ContainerInstance, error)
	ListContainerInstances() ([]types.ContainerInstance, error)
	FilterContainerInstances(filterKey string, filterValue string) ([]types.ContainerInstance, error)
	StreamContainerInstances(ctx context.Context) (chan storetypes.ContainerInstanceErrorWrapper, error)
	DeleteContainerInstance(cluster, instanceARN string) error
}

type eventInstanceStore struct {
	datastore   DataStore
	etcdTXStore EtcdTXStore
}

// NewContainerInstanceStore inistializes the eventInstanceStore struct
func NewContainerInstanceStore(ds DataStore, ts EtcdTXStore) (ContainerInstanceStore, error) {
	if ds == nil {
		return nil, errors.New("Datastore is not initialized")
	}
	if ts == nil {
		return nil, errors.New("Etcd transactional store is not initialized")
	}

	return eventInstanceStore{
		datastore:   ds,
		etcdTXStore: ts,
	}, nil
}

// AddContainerInstance adds a container instance represented in the instanceJSON to the datastore
func (instanceStore eventInstanceStore) AddContainerInstance(instanceJSON string) error {
	instance, key, err := instanceStore.unmarshalInstanceAndGenerateKey(instanceJSON)
	if err != nil {
		return err
	}

	log.Debugf("Instance store unmarshalled instance: %s, trying to add it to the store", instance.Detail.String())

	applier := &STMApplier{
		record:     types.ContainerInstance{},
		recordKey:  key,
		recordJSON: instanceJSON,
	}
	// TODO: NewSTMRepeatble panics if there's any error from the etcd
	// client. We should find a better way to handle that
	_, err = instanceStore.etcdTXStore.NewSTMRepeatable(context.TODO(),
		instanceStore.etcdTXStore.GetV3Client(),
		applier.applyVersionedRecord)

	return err
}

// AddUnversionedContainerInstance adds a container instance represented in the instanceJSON to the datastore only if the instance version is -1
func (instanceStore eventInstanceStore) AddUnversionedContainerInstance(instanceJSON string) error {
	instance, key, err := instanceStore.unmarshalInstanceAndGenerateKey(instanceJSON)
	if err != nil {
		return err
	}

	if instance.Detail.Version == nil || aws.Int64Value(instance.Detail.Version) != unversionedInstance {
		return errors.Errorf("Instance version while adding unversioned instance should be set to %d", unversionedTask)
	}

	log.Debugf("Instance store unmarshalled unversioned instance: %s, trying to add it to the store", instance.Detail.String())

	applier := &STMApplier{
		record:     types.ContainerInstance{},
		recordKey:  key,
		recordJSON: instanceJSON,
	}
	// TODO: NewSTMRepeatble panics if there's any error from the etcd
	// client. We should find a better way to handle that
	_, err = instanceStore.etcdTXStore.NewSTMRepeatable(context.TODO(),
		instanceStore.etcdTXStore.GetV3Client(),
		applier.applyUnversionedRecord)

	return err
}

// GetContainerInstance gets a container with ARN 'instanceARN' belonging to cluster 'cluster'
func (instanceStore eventInstanceStore) GetContainerInstance(cluster string, instanceARN string) (*types.ContainerInstance, error) {
	key, err := instanceStore.getInstanceKey(cluster, instanceARN)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not generate instance key for cluster '%s' and instance '%s'", cluster, instanceARN)
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
	case filterKey == instanceClusterFilter:
		return instanceStore.filterContainerInstancesByCluster(filterValue)
	default:
		return nil, errors.Errorf("Unsupported filter key: %s", filterKey)
	}
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

// DeleteContainerInstance deletes the container instance record from the data store
func (instanceStore eventInstanceStore) DeleteContainerInstance(cluster string, instanceARN string) error {
	key, err := instanceStore.getInstanceKey(cluster, instanceARN)
	if err != nil {
		return errors.Wrapf(err, "Could not generate instance key for cluster '%s' and instance '%s'",
			cluster, instanceARN)
	}
	numKeysDeleted, err := instanceStore.datastore.Delete(key)
	log.Debugf("Deleted '%d' key(s) from the store for container instance '%s', belonging to cluster '%s'",
		numKeysDeleted, instanceARN, cluster)
	// TODO: Should numKeysDeleted != 1 cause an error as well?
	return err
}

func (instanceStore eventInstanceStore) unmarshalInstanceAndGenerateKey(instanceJSON string) (*types.ContainerInstance, string, error) {
	if len(instanceJSON) == 0 {
		return nil, "", errors.New("Instance JSON should not be empty")
	}

	instance, err := instanceStore.unmarshalInstance(instanceJSON)
	if err != nil {
		return nil, "", err
	}

	if instance.Detail == nil {
		return nil, "", errors.New("Instance detail not initialized in JSON")
	}
	if aws.StringValue(instance.Detail.ClusterARN) == "" {
		return nil, "", errors.New("Cluster ARN should not be empty in instance JSON")
	}
	if aws.StringValue(instance.Detail.ContainerInstanceARN) == "" {
		return nil, "", errors.New("Container instance ARN should not be empty in instance JSON")
	}

	clusterARN := aws.StringValue(instance.Detail.ClusterARN)
	clusterName, err := regex.GetClusterNameFromARN(clusterARN)
	if err != nil {
		return nil, "", errors.Wrapf(err, "Error retrieving cluster name from ARN '%s' for instance", clusterARN)
	}

	key, err := generateInstanceKey(clusterName, aws.StringValue(instance.Detail.ContainerInstanceARN))
	if err != nil {
		return nil, "", err
	}

	return &instance, key, nil
}

func (instanceStore eventInstanceStore) filterContainerInstancesByStatus(status string) ([]types.ContainerInstance, error) {
	instances, err := instanceStore.ListContainerInstances()
	if err != nil {
		return nil, err
	}
	filteredInstances := make([]types.ContainerInstance, 0, len(instances))
	for _, instance := range instances {
		if strings.ToLower(status) == strings.ToLower(aws.StringValue(instance.Detail.Status)) {
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
		instance, err = instanceStore.unmarshalInstance(v)
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
		instance, err := instanceStore.unmarshalInstance(string(v))
		if err != nil {
			return nil, err
		}
		instances = append(instances, instance)
	}
	return instances, nil
}

func (instanceStore eventInstanceStore) getInstanceKey(cluster string, instanceARN string) (string, error) {
	if len(cluster) == 0 {
		return "", errors.New("Cluster should not be empty")
	}
	if len(instanceARN) == 0 {
		return "", errors.New("Instance ARN should not be empty")
	}

	clusterName := cluster
	var err error
	if regex.IsClusterARN(cluster) {
		clusterName, err = regex.GetClusterNameFromARN(cluster)
		if err != nil {
			return "", err
		}
	}

	return generateInstanceKey(clusterName, instanceARN)
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
				ins, err := instanceStore.unmarshalInstance(v)
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

func (instanceStore eventInstanceStore) unmarshalInstance(val string) (types.ContainerInstance, error) {
	var instance types.ContainerInstance
	err := json.Unmarshal([]byte(val), &instance)
	if err != nil {
		return instance, errors.Wrapf(err, "Error unmarshaling instance '%s'", val)
	}

	return instance, nil
}

func generateInstanceKey(clusterName string, instanceARN string) (string, error) {
	if !regex.IsClusterName(clusterName) {
		return "", errors.Errorf("Error generating instance key. Cluster name '%s' does not match expected regex", clusterName)
	}
	if !regex.IsInstanceARN(instanceARN) {
		return "", errors.Errorf("Error generating instance key. Instance ARN '%s' does not match expected regex", instanceARN)
	}
	return instanceKeyPrefix + clusterName + "/" + instanceARN, nil
}
