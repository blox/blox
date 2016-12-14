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

package loader

import (
	"encoding/json"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
	"github.com/blox/blox/cluster-state-service/handler/store"
	"github.com/blox/blox/cluster-state-service/handler/types"
	log "github.com/cihub/seelog"
	"github.com/pkg/errors"
)

// ContainerInstanceLoader defines the interface to load container instances from
// the data store and ECS and to merge the same.
type ContainerInstanceLoader interface {
	LoadContainerInstances() error
}

// instanceLoader implements the ContainerInstanceLoader interface.
type instanceLoader struct {
	instanceStore store.ContainerInstanceStore
	ecsWrapper    ECSWrapper
}

// instanceARNLookup maps instance ARNs to a struct. This is to facilitate easy lookup
// of instance ARNs.
type instanceARNLookup map[string]struct{}

// clusterARNsToInstances maps cluster ARNs to the instanceARNLookup map. This is to
// faciliate easy lookup of cluster ARNs to instance ARNs.
type clusterARNsToInstances map[string]instanceARNLookup

// instanceKeyToDelete is a wrapper for instance and cluster ARNs to delete.
type instanceKeyToDelete struct {
	instanceARN string
	clusterARN  string
}

func NewContainerInstanceLoader(instanceStore store.ContainerInstanceStore, ecsClient ecsiface.ECSAPI) ContainerInstanceLoader {
	return instanceLoader{
		instanceStore: instanceStore,
		ecsWrapper:    NewECSWrapper(ecsClient),
	}
}

// LoadContainerInstances retrieves all instances belonging to all clusters in ECS and loads them into data store
func (loader instanceLoader) LoadContainerInstances() error {
	// Construct a map of clusters to instances for instances in local data store.
	localState, err := loader.loadLocalClusterStateFromStore()
	if err != nil {
		return errors.Wrapf(err, "Error loading instances from data store")
	}
	clusterARNs, err := loader.ecsWrapper.ListAllClusters()
	if err != nil {
		return errors.Wrapf(err, "Error listing clusters from ECS")
	}
	ecsState := make(clusterARNsToInstances)
	for _, cluster := range clusterARNs {
		// TODO Parallelize this so that instances across clusters can be
		// gathered in parallel.
		instances, err := loader.getContainerInstancesFromECS(cluster)
		if err != nil {
			return errors.Wrapf(err,
				"Error getting container instances from ECS for cluster '%s'", aws.StringValue(cluster))
		}
		clusterARN := aws.StringValue(cluster)
		// Add the cluster ARN to the lookup map.
		ecsState[clusterARN] = make(instanceARNLookup)
		for _, instance := range instances {
			err := loader.putContainerInstance(instance)
			if err != nil {
				return err
			}
			// Populate the entries for the cluster ARN in the lookup map.
			ecsState[clusterARN][aws.StringValue(instance.Detail.ContainerInstanceARN)] = struct{}{}
		}
	}
	// Get a list of keys to delete from the local store.
	keys := getInstanceKeysNotInECS(localState, ecsState)
	log.Debugf("Instances to delete: %v", keys)
	for _, key := range keys {
		// Not handling returned error because we want as many cleanup operations to succeed as possible.
		if err := loader.instanceStore.DeleteContainerInstance(key.clusterARN, key.instanceARN); err != nil {
			log.Infof("Error deleting container instance '%s' belonging to cluster '%s' from data store",
				key.instanceARN, key.clusterARN)
		}
	}
	return nil
}

// loadLocalClusterStateFromStore loads container instance records from local store into a
// map for easy lookup and comparison
func (loader instanceLoader) loadLocalClusterStateFromStore() (clusterARNsToInstances, error) {
	instances, err := loader.instanceStore.ListContainerInstances()
	if err != nil {
		return nil, errors.Wrapf(err, "Error loading instances from store")
	}

	state := make(clusterARNsToInstances)
	for _, instance := range instances {
		clusterARN := aws.StringValue(instance.Detail.ClusterARN)
		if _, ok := state[clusterARN]; !ok {
			state[clusterARN] = make(instanceARNLookup)
		}
		state[clusterARN][aws.StringValue(instance.Detail.ContainerInstanceARN)] = struct{}{}
	}

	return state, nil
}

// getContainerInstancesFromECS gets a list of container instances from ECS for the specified cluster.
func (loader instanceLoader) getContainerInstancesFromECS(cluster *string) ([]types.ContainerInstance, error) {
	var instances []types.ContainerInstance
	instanceARNs, err := loader.ecsWrapper.ListAllContainerInstances(cluster)
	if err != nil {
		return instances, errors.Wrapf(err,
			"Error listing all container instances for cluster '%s'", aws.StringValue(cluster))
	}
	if len(instanceARNs) == 0 {
		return instances, nil
	}
	instances, failedInstanceARNs, err := loader.ecsWrapper.DescribeContainerInstances(cluster, instanceARNs)
	if err != nil {
		return instances, errors.Wrapf(err,
			"Error describing container instances for cluster '%s'", aws.StringValue(cluster))
	}
	if len(failedInstanceARNs) != 0 {
		// If we're unable to describe listed container instances, just print the list out.
		// Since we treat ECS as the source of truth, it should be fine to make this assumption.
		log.Infof("Failed to describe listed instances: %s", strings.Join(failedInstanceARNs[:], " "))
	}
	return instances, nil
}

// putContainerInstance puts the container instance record to the data store
func (loader instanceLoader) putContainerInstance(instance types.ContainerInstance) error {
	ins, err := json.Marshal(instance)
	if err != nil {
		return errors.Wrapf(err, "Failed to marshal instance JSON")
	}
	instanceJSON := string(ins)
	err = loader.instanceStore.AddUnversionedContainerInstance(instanceJSON)
	if err != nil {
		return errors.Wrapf(err, "Failed to add unversioned container instance '%s'", instanceJSON)
	}
	return nil
}

// getInstanceKeysNotInECS gets a list of instance keys to delete from the local store. This is
// the set of keys that are in the local store, but not in ECS
func getInstanceKeysNotInECS(localState, ecsState clusterARNsToInstances) []instanceKeyToDelete {
	var instanceKeysNotInECS []instanceKeyToDelete
	// For each cluster in local state, get all instance records
	for clusterARN, instanceRecords := range localState {
		// Check if cluster in local state exists in ecs state
		ecsInstanceRecords, ok := ecsState[clusterARN]
		if !ok {
			// Cluster in local state not found in ECS state
			// Add all instance records to the to-be-deleted list
			for instanceARN, _ := range instanceRecords {
				instanceKeysNotInECS = append(instanceKeysNotInECS, instanceKeyToDelete{
					instanceARN: instanceARN,
					clusterARN:  clusterARN,
				})
			}
			continue
		}
		// Cluster in local state found in ECS state. Compare all
		// instances that belong to the cluster to those in ECS
		for instanceARN, _ := range instanceRecords {
			if _, ok := ecsInstanceRecords[instanceARN]; !ok {
				instanceKeysNotInECS = append(instanceKeysNotInECS, instanceKeyToDelete{
					instanceARN: instanceARN,
					clusterARN:  clusterARN,
				})
			}
		}
	}
	return instanceKeysNotInECS
}
