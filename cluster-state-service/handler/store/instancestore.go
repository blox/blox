// Copyright 2016-2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.
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
)

var (
	supportedInstanceFilters = map[string]string{instanceStatusFilter: "", instanceClusterFilter: ""}
)

// ContainerInstanceStore defines methods to access container instances from the datastore
type ContainerInstanceStore interface {
	AddContainerInstance(instance string) error
	GetContainerInstance(cluster string, instanceARN string) (*storetypes.VersionedContainerInstance, error)
	ListContainerInstances() ([]storetypes.VersionedContainerInstance, error)
	FilterContainerInstances(filterMap map[string]string) ([]storetypes.VersionedContainerInstance, error)
	StreamContainerInstances(ctx context.Context, entityVersion string) (chan storetypes.VersionedContainerInstance, error)
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
		applier.applyRecord)

	return err
}

// GetContainerInstance gets a container with ARN 'instanceARN' belonging to cluster 'cluster'
func (instanceStore eventInstanceStore) GetContainerInstance(cluster string, instanceARN string) (*storetypes.VersionedContainerInstance, error) {
	key, err := instanceStore.getInstanceKey(cluster, instanceARN)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not generate instance key for cluster '%s' and instance '%s'", cluster, instanceARN)
	}
	return instanceStore.getInstanceByKey(key)
}

// ListContainerInstances lists all container instances existing in the datastore
func (instanceStore eventInstanceStore) ListContainerInstances() ([]storetypes.VersionedContainerInstance, error) {
	return instanceStore.getInstancesByKeyPrefix(instanceKeyPrefix)
}

// FilterContainerInstances returns all container instances from the datastore that match the provided filters
func (instanceStore eventInstanceStore) FilterContainerInstances(filterMap map[string]string) ([]storetypes.VersionedContainerInstance, error) {
	if len(filterMap) == 0 {
		return nil, errors.New("There has to be at least one filter")
	}

	filters := make([]string, 0, len(filterMap))
	for k := range filterMap {
		filters = append(filters, k)
	}

	if !instanceStore.areFiltersValid(filters) {
		return nil, errors.Errorf("At least one of the provided filters '%v' is not supported.", filters)
	}

	for key, val := range filterMap {
		if val == "" {
			return nil, errors.Errorf("Filter value for filter '%s' is empty", key)
		}
	}

	status, statusFilterExists := filterMap[instanceStatusFilter]
	cluster, clusterFilterExists := filterMap[instanceClusterFilter]
	switch {
	case statusFilterExists && clusterFilterExists:
		return instanceStore.filterContainerInstancesByStatusAndCluster(status, cluster)
	case statusFilterExists:
		return instanceStore.filterContainerInstancesByStatus(status)
	case clusterFilterExists:
		return instanceStore.filterContainerInstancesByCluster(cluster)
	default:
		return nil, errors.Errorf("Unsupported filter combination '%v'", filters)
	}
}

// StreamContainerInstances returns a stream of all changes in the container instance keyspace
func (instanceStore eventInstanceStore) StreamContainerInstances(ctx context.Context, entityVersion string) (chan storetypes.VersionedContainerInstance, error) {
	instanceStoreCtx, cancel := context.WithCancel(ctx) // go routine instanceStore.pipeBetweenChannels() handles canceling this context

	dsChan, err := instanceStore.datastore.StreamWithPrefix(instanceStoreCtx, instanceKeyPrefix, entityVersion)
	if err != nil {
		cancel()
		return nil, err
	}

	instanceRespChan := make(chan storetypes.VersionedContainerInstance) // go routine instanceStore.pipeBetweenChannels() handles closing of this channel
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

func (instanceStore eventInstanceStore) areFiltersValid(filters []string) bool {
	if len(filters) > len(supportedInstanceFilters) {
		return false
	}
	for _, f := range filters {
		_, ok := supportedInstanceFilters[f]
		if !ok {
			return false
		}
	}
	return true
}

func (instanceStore eventInstanceStore) filterContainerInstancesByStatus(status string) ([]storetypes.VersionedContainerInstance, error) {
	instances, err := instanceStore.ListContainerInstances()
	if err != nil {
		return nil, err
	}
	return instanceStore.filterContainerInstancesByStatusFromList(status, instances), nil
}

func (instanceStore eventInstanceStore) filterContainerInstancesByStatusFromList(status string, instances []storetypes.VersionedContainerInstance) []storetypes.VersionedContainerInstance {
	filteredInstances := make([]storetypes.VersionedContainerInstance, 0, len(instances))
	for _, instance := range instances {
		if strings.ToLower(status) == strings.ToLower(aws.StringValue(instance.ContainerInstance.Detail.Status)) {
			filteredInstances = append(filteredInstances, instance)
		}
	}
	return filteredInstances
}

func (instanceStore eventInstanceStore) filterContainerInstancesByCluster(cluster string) ([]storetypes.VersionedContainerInstance, error) {
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

func (instanceStore eventInstanceStore) filterContainerInstancesByStatusAndCluster(status string, cluster string) ([]storetypes.VersionedContainerInstance, error) {
	instancesFilteredByCluster, err := instanceStore.filterContainerInstancesByCluster(cluster)
	if err != nil {
		return nil, err
	}
	return instanceStore.filterContainerInstancesByStatusFromList(status, instancesFilteredByCluster), nil
}

func (instanceStore eventInstanceStore) getInstanceByKey(key string) (*storetypes.VersionedContainerInstance, error) {
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

	var versionedInstance storetypes.VersionedContainerInstance
	for _, entity := range resp {
		versionedInstance.ContainerInstance, err = instanceStore.unmarshalInstance(entity.Value)
		versionedInstance.Version = entity.Version
		if err != nil {
			return nil, err
		}
		break
	}
	return &versionedInstance, nil
}

func (instanceStore eventInstanceStore) getInstancesByKeyPrefix(key string) ([]storetypes.VersionedContainerInstance, error) {
	if len(key) == 0 {
		return nil, errors.New("Key cannot be empty")
	}

	resp, err := instanceStore.datastore.GetWithPrefix(key)
	if err != nil {
		return nil, err
	}

	if len(resp) == 0 {
		return make([]storetypes.VersionedContainerInstance, 0), nil
	}

	versionedInstances := []storetypes.VersionedContainerInstance{}
	for _, entity := range resp {
		var versionedInstance storetypes.VersionedContainerInstance
		versionedInstance.ContainerInstance, err = instanceStore.unmarshalInstance(entity.Value)
		versionedInstance.Version = entity.Version
		if err != nil {
			return nil, err
		}
		versionedInstances = append(versionedInstances, versionedInstance)
	}
	return versionedInstances, nil
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

func (instanceStore eventInstanceStore) pipeBetweenChannels(ctx context.Context, cancel context.CancelFunc, dsChan chan map[string]storetypes.Entity, instanceRespChan chan storetypes.VersionedContainerInstance) {
	defer close(instanceRespChan)
	defer cancel()

	for {
		select {
		case resp, ok := <-dsChan:
			if !ok {
				return
			}
			for _, entity := range resp {
				var versionedInstance storetypes.VersionedContainerInstance
				ins, err := instanceStore.unmarshalInstance(entity.Value)
				if err != nil {
					versionedInstance.Err = err
					instanceRespChan <- versionedInstance
					return
				}
				versionedInstance.ContainerInstance = ins
				versionedInstance.Version = entity.Version
				instanceRespChan <- versionedInstance
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
