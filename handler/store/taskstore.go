package store

import (
	"github.com/aws/amazon-ecs-event-stream-handler/handler/json"
	"github.com/aws/amazon-ecs-event-stream-handler/handler/types"
	log "github.com/cihub/seelog"
	"github.com/pkg/errors"
	"strings"
)

const (
	taskKeyPrefix = "ecs/task/"
	statusFilter  = "status"
)

// TaskStore defines methods to access tasks from the datastore
type TaskStore interface {
	AddTask(task string) error
	GetTask(arn string) (*types.Task, error)
	ListTasks() ([]types.Task, error)
	FilterTasks(filterKey string, filterValue string) ([]types.Task, error)
	StreamTasks() ([]types.Task, error)
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

func generateTaskKey(task types.Task) (string, error) {
	if len(task.Detail.TaskArn) == 0 {
		return "", errors.New("Task arn cannot be empty")
	}
	return taskKeyPrefix + task.Detail.TaskArn, nil
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

	key, err := generateTaskKey(task)
	if err != nil {
		return err
	}

	// check if record exists with higher version number
	existingTask, err := taskStore.getTaskByKey(key)
	if err != nil {
		return err
	}

	if existingTask != nil {
		if existingTask.Detail.Version >= task.Detail.Version {
			log.Infof("Higher or equal version %v of task %v with version %v already exists",
				existingTask.Detail.Version,
				task.Detail.TaskArn,
				task.Detail.Version)

			// do nothing. later version of the event has already been stored
			return nil
		}
	}

	err = taskStore.datastore.Add(key, taskJSON)
	if err != nil {
		return err
	}

	return nil
}

// GetTask gets a task with key 'arn' from the datastore
func (taskStore eventTaskStore) GetTask(arn string) (*types.Task, error) {
	if len(arn) == 0 {
		return nil, errors.New("Arn should not be empty")
	}

	var task types.Task
	task.Detail.TaskArn = arn

	key, err := generateTaskKey(task)
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
	if statusFilter != filterKey {
		return nil, errors.Errorf("Filter '%s' not supported", filterKey)
	}

	tasks, err := taskStore.ListTasks()
	if err != nil {
		return nil, err
	}

	tasksWithStatus := []types.Task{}
	for _, task := range tasks {
		if strings.ToLower(filterValue) == strings.ToLower(task.Detail.LastStatus) {
			tasksWithStatus = append(tasksWithStatus, task)
		}
	}

	return tasksWithStatus, nil
}

// StreamTasks returns a stream of all changes in the task keyspace
func (taskStore eventTaskStore) StreamTasks() ([]types.Task, error) {
	return nil, nil
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
		err = json.UnmarshalJSON(v, &task)
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
		var task types.Task
		err = json.UnmarshalJSON(v, &task)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}
