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
	taskKeyPrefix       = "ecs/task/"
	taskStatusFilter    = "status"
	taskStartedByFilter = "startedBy"
	taskClusterFilter   = "cluster"

	unversionedTask = -1
)

// TaskStore defines methods to access tasks from the datastore
type TaskStore interface {
	AddTask(task string) error
	AddUnversionedTask(task string) error
	GetTask(cluster string, taskARN string) (*types.Task, error)
	ListTasks() ([]types.Task, error)
	FilterTasks(filterKey string, filterValue string) ([]types.Task, error)
	StreamTasks(ctx context.Context) (chan storetypes.TaskErrorWrapper, error)
	DeleteTask(cluster, taskARN string) error
}

type eventTaskStore struct {
	datastore   DataStore
	etcdTXStore EtcdTXStore
}

// NewTaskStore initializes the eventTaskStore struct
func NewTaskStore(ds DataStore, ts EtcdTXStore) (TaskStore, error) {
	if ds == nil {
		return nil, errors.Errorf("Datastore is not initialized")
	}
	if ts == nil {
		return nil, errors.Errorf("Etcd transactional store is not initialized")
	}

	return eventTaskStore{
		datastore:   ds,
		etcdTXStore: ts,
	}, nil
}

// AddTask adds a task represented in the taskJSON to the datastore
func (taskStore eventTaskStore) AddTask(taskJSON string) error {
	task, key, err := taskStore.unmarshalTaskAndGenerateKey(taskJSON)
	if err != nil {
		return err
	}

	log.Debugf("Task store unmarshalled task: %s, trying to add it to the store", task.Detail.String())

	applier := &STMApplier{
		record:     types.Task{},
		recordKey:  key,
		recordJSON: taskJSON,
	}
	// TODO: NewSTMRepeatble panics if there's any error from the etcd
	// client. We should find a better way to handle that
	_, err = taskStore.etcdTXStore.NewSTMRepeatable(context.TODO(),
		taskStore.etcdTXStore.GetV3Client(),
		applier.applyVersionedRecord)
	return err
}

// AddUnversionedTask adds a task represented in the taskJSON to the datastore only if the task version is set to -1
func (taskStore eventTaskStore) AddUnversionedTask(taskJSON string) error {
	task, key, err := taskStore.unmarshalTaskAndGenerateKey(taskJSON)
	if err != nil {
		return err
	}

	if task.Detail.Version == nil || aws.Int64Value(task.Detail.Version) != unversionedTask {
		return errors.Errorf("Task version while adding unversioned task should be set to %d", unversionedTask)
	}

	log.Debugf("Task store unmarshalled unversioned task: %s, trying to add it to the store", task.Detail.String())

	applier := &STMApplier{
		record:     types.Task{},
		recordKey:  key,
		recordJSON: taskJSON,
	}
	// TODO: NewSTMRepeatble panics if there's any error from the etcd
	// client. We should find a better way to handle that
	_, err = taskStore.etcdTXStore.NewSTMRepeatable(context.TODO(),
		taskStore.etcdTXStore.GetV3Client(),
		applier.applyUnversionedRecord)
	return err
}

// GetTask gets a task with ARN 'taskARN' belonging to cluster 'cluster'
func (taskStore eventTaskStore) GetTask(cluster string, taskARN string) (*types.Task, error) {
	key, err := taskStore.getTaskKey(cluster, taskARN)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not generate task key for cluster '%s' and task '%s'",
			cluster, taskARN)
	}

	return taskStore.getTaskByKey(key)
}

// ListTasks lists all the tasks existing in the datastore
func (taskStore eventTaskStore) ListTasks() ([]types.Task, error) {
	return taskStore.getTasksByKeyPrefix(taskKeyPrefix)
}

// FilterTasks returns all the tasks from the datastore that match the provided filters
func (taskStore eventTaskStore) FilterTasks(filterKey string, filterValue string) ([]types.Task, error) {
	if len(filterKey) == 0 || len(filterValue) == 0 {
		return nil, errors.New("Filter key and value cannot be empty")
	}

	switch {
	case filterKey == taskStatusFilter:
		return taskStore.filterTasks(isTaskStatus, filterValue)
	case filterKey == taskStartedByFilter:
		return taskStore.filterTasks(isTaskStartedBy, filterValue)
	case filterKey == taskClusterFilter:
		return taskStore.filterTasksByCluster(filterValue)
	default:
		return nil, errors.Errorf("Unsupported filter key: %s", filterKey)
	}

}

// StreamTasks streams all changes in the task keyspace into a channel
func (taskStore eventTaskStore) StreamTasks(ctx context.Context) (chan storetypes.TaskErrorWrapper, error) {
	taskStoreCtx, cancel := context.WithCancel(ctx) // go routine taskStore.pipeBetweenChannels() handles canceling this context

	dsChan, err := taskStore.datastore.StreamWithPrefix(taskStoreCtx, taskKeyPrefix)
	if err != nil {
		cancel()
		return nil, err
	}

	taskRespChan := make(chan storetypes.TaskErrorWrapper) // go routine taskStore.pipeBetweenChannels() handles closing of this channel
	go taskStore.pipeBetweenChannels(taskStoreCtx, cancel, dsChan, taskRespChan)
	return taskRespChan, nil
}

// DeleteTask deletes the task key from the data store
func (taskStore eventTaskStore) DeleteTask(cluster string, taskARN string) error {
	key, err := taskStore.getTaskKey(cluster, taskARN)
	if err != nil {
		return errors.Wrapf(err, "Could not generate task key for cluster '%s' and task '%s'",
			cluster, taskARN)
	}

	numKeysDeleted, err := taskStore.datastore.Delete(key)
	log.Debugf("Deleted '%d' key(s) from the store for task '%s', belonging to cluster '%s'",
		numKeysDeleted, taskARN, cluster)
	// TODO: Should numKeysDeleted != 1 cause an error as well?
	return err
}

func (taskStore eventTaskStore) unmarshalTaskAndGenerateKey(taskJSON string) (*types.Task, string, error) {
	if len(taskJSON) == 0 {
		return nil, "", errors.New("Task json should not be empty")
	}

	task, err := taskStore.unmarshalString(taskJSON)
	if err != nil {
		return nil, "", err
	}

	if task.Detail == nil || task.Detail.ClusterARN == nil || task.Detail.TaskARN == nil {
		return nil, "", errors.New("Cluster ARN and task ARN should not be empty in task JSON")
	}

	clusterARN := aws.StringValue(task.Detail.ClusterARN)
	clusterName, err := regex.GetClusterNameFromARN(clusterARN)
	if err != nil {
		return nil, "", errors.Wrapf(err, "Error retrieving cluster name from ARN '%s' for task", clusterARN)
	}

	key, err := generateTaskKey(clusterName, aws.StringValue(task.Detail.TaskARN))
	if err != nil {
		return nil, "", err
	}

	return &task, key, nil
}

type taskFilter func(string, types.Task) bool

func isTaskStatus(status string, task types.Task) bool {
	return strings.ToLower(status) == strings.ToLower(aws.StringValue(task.Detail.LastStatus))
}

func isTaskStartedBy(startedBy string, task types.Task) bool {
	return startedBy == task.Detail.StartedBy
}

func (taskStore eventTaskStore) filterTasks(filter taskFilter, filterValue string) ([]types.Task, error) {
	tasks, err := taskStore.ListTasks()
	if err != nil {
		return nil, err
	}

	result := []types.Task{}
	for _, task := range tasks {
		if filter(filterValue, task) {
			result = append(result, task)
		}
	}

	return result, nil
}

func (taskStore eventTaskStore) filterTasksByCluster(cluster string) ([]types.Task, error) {
	clusterName := cluster
	var err error
	if regex.IsClusterARN(cluster) {
		clusterName, err = regex.GetClusterNameFromARN(cluster)
		if err != nil {
			return nil, err
		}
	}

	tasksForClusterPrefix := taskKeyPrefix + clusterName + "/"
	return taskStore.getTasksByKeyPrefix(tasksForClusterPrefix)
}

func (taskStore eventTaskStore) getTaskKey(cluster string, taskARN string) (string, error) {
	if len(cluster) == 0 {
		return "", errors.New("Cluster should not be empty")
	}

	if len(taskARN) == 0 {
		return "", errors.New("Task ARN should not be empty")
	}

	clusterName := cluster
	var err error
	if regex.IsClusterARN(cluster) {
		clusterName, err = regex.GetClusterNameFromARN(cluster)
		if err != nil {
			return "", err
		}
	}

	return generateTaskKey(clusterName, taskARN)
}

func (taskStore eventTaskStore) pipeBetweenChannels(ctx context.Context, cancel context.CancelFunc, dsChan chan map[string]string, taskRespChan chan storetypes.TaskErrorWrapper) {
	defer close(taskRespChan)
	defer cancel()

	for {
		select {
		case resp, ok := <-dsChan:
			if !ok {
				return
			}
			for _, v := range resp {
				t, err := taskStore.unmarshalString(v)
				if err != nil {
					taskRespChan <- storetypes.TaskErrorWrapper{Task: types.Task{}, Err: err}
					return
				}
				taskRespChan <- storetypes.TaskErrorWrapper{Task: t, Err: nil}
			}

		case <-ctx.Done():
			return
		}
	}
}

func (taskStore eventTaskStore) getTaskByKey(key string) (*types.Task, error) {
	if len(key) == 0 {
		return nil, errors.New("Key cannot be empty")
	}

	resp, err := taskStore.datastore.Get(key)
	if err != nil {
		return nil, err
	}

	if len(resp) == 0 {
		return nil, nil
	}

	if len(resp) > 1 {
		return nil, errors.Errorf("Multiple entries exist in the datastore with key %v", key)
	}

	var task types.Task
	for _, v := range resp {
		task, err = taskStore.unmarshalString(v)
		if err != nil {
			return nil, err
		}
		break
	}
	return &task, nil
}

func (taskStore eventTaskStore) getTasksByKeyPrefix(key string) ([]types.Task, error) {
	if len(key) == 0 {
		return nil, errors.New("Key cannot be empty")
	}

	resp, err := taskStore.datastore.GetWithPrefix(key)
	if err != nil {
		return nil, err
	}

	if len(resp) == 0 {
		return make([]types.Task, 0), nil
	}

	tasks := []types.Task{}
	for _, v := range resp {
		task, err := taskStore.unmarshalString(v)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (taskStore eventTaskStore) unmarshalString(val string) (types.Task, error) {
	var task types.Task
	err := json.Unmarshal([]byte(val), &task)
	if err != nil {
		return task, errors.Wrapf(err, "Error unmarshaling task '%s'", val)
	}

	return task, nil
}

func generateTaskKey(clusterName string, taskARN string) (string, error) {
	if !regex.IsClusterName(clusterName) {
		return "", errors.Errorf("Error generating task key. Cluster name '%s' does not match expected regex", clusterName)
	}
	if !regex.IsTaskARN(taskARN) {
		return "", errors.Errorf("Error generating task key. Task ARN '%s' does not match expected regex", taskARN)
	}
	return taskKeyPrefix + clusterName + "/" + taskARN, nil
}
