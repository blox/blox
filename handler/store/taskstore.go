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
	taskKeyPrefix    = "ecs/task/"
	taskStatusFilter = "status"
)

// TaskStore defines methods to access tasks from the datastore
type TaskStore interface {
	AddTask(task string) error
	GetTask(cluster string, taskARN string) (*types.Task, error)
	ListTasks() ([]types.Task, error)
	FilterTasks(filterKey string, filterValue string) ([]types.Task, error)
	StreamTasks(ctx context.Context) (chan storetypes.TaskErrorWrapper, error)
}

type eventTaskStore struct {
	datastore DataStore
}

func NewTaskStore(ds DataStore) (TaskStore, error) {
	if ds == nil {
		return nil, errors.Errorf("The datastore cannot be nil")
	}

	return eventTaskStore{
		datastore: ds,
	}, nil
}

func generateTaskKey(clusterName string, taskARN string) (string, error) {
	if !regex.IsClusterName(clusterName) {
		return "", errors.New("Cluster name does not match expected regex")
	}
	if !regex.IsTaskARN(taskARN) {
		return "", errors.New("Task ARN does not match expected regex")
	}
	return taskKeyPrefix + clusterName + "/" + taskARN, nil
}

// AddTask adds a task represented in the taskJSON to the datastore
func (taskStore eventTaskStore) AddTask(taskJSON string) error {
	if len(taskJSON) == 0 {
		return errors.New("Task json should not be empty")
	}

	var task types.Task
	err := json.UnmarshalJSON(taskJSON, &task)
	if err != nil {
		return err
	}

	if task.Detail == nil || task.Detail.ClusterArn == nil || task.Detail.TaskArn == nil {
		return errors.New("Cluster ARN and task ARN should not be empty in task JSON")
	}

	clusterName, err := regex.GetClusterNameFromARN(*task.Detail.ClusterArn)
	if err != nil {
		return err
	}

	key, err := generateTaskKey(clusterName, *task.Detail.TaskArn)
	if err != nil {
		return err
	}

	// check if record exists with higher version number
	existingTask, err := taskStore.getTaskByKey(key)
	if err != nil {
		return err
	}

	if existingTask != nil {
		existingTaskDetail := *existingTask.Detail
		currentTaskDetail := *task.Detail
		if *existingTaskDetail.Version >= *currentTaskDetail.Version {
			log.Infof("Higher or equal version %v of task %v with version %v already exists",
				existingTask.Detail.Version,
				task.Detail.TaskArn,
				task.Detail.Version)

			// do nothing. later version of the event has already been stored
			return nil
		}
	}

	compressedTaskJSON, err := compress.Compress(taskJSON)
	if err != nil {
		return err
	}

	err = taskStore.datastore.Add(key, string(compressedTaskJSON))
	if err != nil {
		return err
	}

	return nil
}

// GetTask gets a task with ARN 'taskARN' belonging to cluster 'cluster'
func (taskStore eventTaskStore) GetTask(cluster string, taskARN string) (*types.Task, error) {
	if len(cluster) == 0 {
		return nil, errors.New("Cluster should not be empty")
	}

	if len(taskARN) == 0 {
		return nil, errors.New("Task ARN should not be empty")
	}

	clusterName := cluster
	var err error
	if regex.IsClusterARN(cluster) {
		clusterName, err = regex.GetClusterNameFromARN(cluster)
		if err != nil {
			return nil, err
		}
	}

	key, err := generateTaskKey(clusterName, taskARN)
	if err != nil {
		return nil, err
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

	//TODO: make generic by finding the field name using reflection so we can filter
	//on arbitrary fields
	if taskStatusFilter != filterKey {
		return nil, errors.Errorf("Filter '%s' not supported", filterKey)
	}

	tasks, err := taskStore.ListTasks()
	if err != nil {
		return nil, err
	}

	tasksWithStatus := []types.Task{}
	for _, task := range tasks {
		if strings.ToLower(filterValue) == strings.ToLower(*task.Detail.LastStatus) {
			tasksWithStatus = append(tasksWithStatus, task)
		}
	}

	return tasksWithStatus, nil
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
				t, err := taskStore.uncompressAndUnmarshalString(v)
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
		task, err = taskStore.uncompressAndUnmarshalString(v)
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
		task, err := taskStore.uncompressAndUnmarshalString(v)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (taskStore eventTaskStore) uncompressAndUnmarshalString(val string) (types.Task, error) {
	var task types.Task

	uncompressedVal, err := compress.Uncompress([]byte(val))
	if err != nil {
		return task, err
	}
	err = json.UnmarshalJSON(uncompressedVal, &task)
	if err != nil {
		return task, err
	}

	return task, nil
}
