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

// TaskLoader defines the interface to load container tasks from
// the data store and ECS and to merge the same.
type TaskLoader interface {
	LoadTasks() error
}

// taskLoader implements the TaskLoader interface.
type taskLoader struct {
	taskStore  store.TaskStore
	ecsWrapper ECSWrapper
}

// taskARNLookup maps task ARNs to a struct. This is to facilitate easy lookup
// of task ARNs.
type taskARNLookup map[string]struct{}

// clusterARNsToInstances maps cluster ARNs to the taskARNLookup map. This is to
// faciliate easy lookup of cluster ARNs to task ARNs.
type clusterARNsToTasks map[string]taskARNLookup

// taskKeyToDelete is a wrapper for task and cluster ARNs to delete.
type taskKeyToDelete struct {
	taskARN    string
	clusterARN string
}

func NewTaskLoader(taskStore store.TaskStore, ecsClient ecsiface.ECSAPI) TaskLoader {
	return taskLoader{
		taskStore:  taskStore,
		ecsWrapper: NewECSWrapper(ecsClient),
	}
}

// LoadTasks retrieves all tasks belonging to all clusters in ECS and loads them into data store
func (loader taskLoader) LoadTasks() error {
	// Construct a map of clusters to tasks for tasks in local data store.
	localState, err := loader.loadLocalClusterStateFromStore()
	if err != nil {
		return errors.Wrapf(err, "Error loading tasks from data store")
	}
	// TODO: We do this in both instance and task stores. Optimize to do it in only one place.
	clusterARNs, err := loader.ecsWrapper.ListAllClusters()
	if err != nil {
		return errors.Wrapf(err, "Error listing clusters from ECS")
	}
	ecsState := make(clusterARNsToTasks)
	for _, cluster := range clusterARNs {
		// TODO Parallelize this so that tasks across clusters can be
		// gathered in parallel.
		tasks, err := loader.getTasksFromECS(cluster)
		if err != nil {
			return errors.Wrapf(err,
				"Error getting tasks from ECS for cluster '%s'", aws.StringValue(cluster))
		}
		clusterARN := aws.StringValue(cluster)
		// Add the cluster ARN to the lookup map.
		ecsState[clusterARN] = make(taskARNLookup)
		for _, task := range tasks {
			err := loader.putTask(task)
			if err != nil {
				return err
			}
			// Populate the entries for the cluster ARN in the lookup map.
			ecsState[clusterARN][aws.StringValue(task.Detail.TaskARN)] = struct{}{}
		}
	}
	// Get a list of keys to delete from the local store.
	keys := getTaskKeysNotInECS(localState, ecsState)
	log.Debugf("Tasks to delete: %v", keys)
	for _, key := range keys {
		// Not handling returned error because we want as many cleanup operations to succeed as possible.
		if err := loader.taskStore.DeleteTask(key.clusterARN, key.taskARN); err != nil {
			log.Infof("Error deleting task '%s' belonging to cluster '%s' from data store",
				key.taskARN, key.clusterARN)
		}
	}
	return nil
}

// loadLocalClusterStateFromStore loads task records from local store into a map for
// easy lookup and comparison
func (loader taskLoader) loadLocalClusterStateFromStore() (clusterARNsToTasks, error) {
	tasks, err := loader.taskStore.ListTasks()
	if err != nil {
		return nil, errors.Wrapf(err, "Error loading tasks from store")
	}

	state := make(clusterARNsToTasks)
	for _, task := range tasks {
		clusterARN := aws.StringValue(task.Detail.ClusterARN)
		if _, ok := state[clusterARN]; !ok {
			state[clusterARN] = make(taskARNLookup)
		}
		state[clusterARN][aws.StringValue(task.Detail.TaskARN)] = struct{}{}
	}

	return state, nil
}

// getTasksFromECS gets a list of tasks from ECS for the specified cluster.
func (loader taskLoader) getTasksFromECS(cluster *string) ([]types.Task, error) {
	var tasks []types.Task
	taskARNs, err := loader.ecsWrapper.ListAllTasks(cluster)
	if err != nil {
		return tasks, errors.Wrapf(err,
			"Error listing all tasks for cluster '%s'", aws.StringValue(cluster))
	}
	if len(taskARNs) == 0 {
		return tasks, nil
	}
	tasks, failedTaskARNs, err := loader.ecsWrapper.DescribeTasks(cluster, taskARNs)
	if err != nil {
		return tasks, errors.Wrapf(err,
			"Error describing tasks for cluster '%s'", aws.StringValue(cluster))
	}
	if len(failedTaskARNs) != 0 {
		// If we're unable to describe listed tasks, just print the list out. Since
		// we treat ECS as the source of truth, it should be fine to make this assumption.
		log.Infof("Failed to describe listed tasks: %s", strings.Join(failedTaskARNs[:], " "))
	}
	return tasks, nil
}

// putTask puts the task record into the data store
func (loader taskLoader) putTask(task types.Task) error {
	mTask, err := json.Marshal(task)
	if err != nil {
		return errors.Wrapf(err, "Failed to marshal task JSON")
	}
	taskJSON := string(mTask)
	err = loader.taskStore.AddUnversionedTask(taskJSON)
	if err != nil {
		return errors.Wrapf(err, "Failed to add unversioned task '%s'", taskJSON)
	}
	return nil
}

// getTaskKeysNotInECS gets a list of task keys to delete from the local store. This is
// the set of keys that are in the local store, but not in ECS
func getTaskKeysNotInECS(localState, ecsState clusterARNsToTasks) []taskKeyToDelete {
	var taskKeysNotInECS []taskKeyToDelete
	// For each cluster in local state, get all task records
	for clusterARN, taskRecords := range localState {
		// Check if cluster in local state exists in ecs state
		ecsTaskRecords, ok := ecsState[clusterARN]
		if !ok {
			// Cluster in local state not found in ECS state
			// Add all task records to the to-be-deleted list
			for taskARN, _ := range taskRecords {
				taskKeysNotInECS = append(taskKeysNotInECS, taskKeyToDelete{
					taskARN:    taskARN,
					clusterARN: clusterARN,
				})
			}
			continue
		}
		// Cluster in local state found in ECS state. Compare all
		// tasks that belong to the cluster to those in ECS
		for taskARN, _ := range taskRecords {
			if _, ok := ecsTaskRecords[taskARN]; !ok {
				taskKeysNotInECS = append(taskKeysNotInECS, taskKeyToDelete{
					taskARN:    taskARN,
					clusterARN: clusterARN,
				})
			}
		}
	}
	return taskKeysNotInECS
}
