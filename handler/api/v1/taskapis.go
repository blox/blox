package v1

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/amazon-ecs-event-stream-handler/handler/api/v1/models"
	"github.com/aws/amazon-ecs-event-stream-handler/handler/store"
	"github.com/gorilla/mux"
)

const (
	taskStatusFilter = "status"
)

type TaskAPIs struct {
	taskStore store.TaskStore
}

func NewTaskAPIs(taskStore store.TaskStore) TaskAPIs {
	return TaskAPIs{
		taskStore: taskStore,
	}
}

//TODO: add arn validation
func (taskAPIs TaskAPIs) GetTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskARN := vars["arn"]
	cluster := vars["cluster"]

	if len(taskARN) == 0 || len(cluster) == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(routingServerErrMsg)
		return
	}

	task, err := taskAPIs.taskStore.GetTask(cluster, taskARN)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(internalServerErrMsg)
		return
	}

	if task == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(taskNotFoundClientErrMsg)
		return
	}

	w.Header().Set(contentTypeKey, contentTypeVal)
	w.WriteHeader(http.StatusOK)

	taskModel, err := ToTaskModel(*task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(internalServerErrMsg)
		return
	}

	err = json.NewEncoder(w).Encode(taskModel)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(encodingServerErrMsg)
		return
	}
}

func (taskAPIs TaskAPIs) ListTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := taskAPIs.taskStore.ListTasks()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(internalServerErrMsg)
		return
	}

	w.Header().Set(contentTypeKey, contentTypeVal)
	w.WriteHeader(http.StatusOK)

	taskModels := make([]models.TaskModel, len(tasks))
	for i := range tasks {
		taskModels[i], err = ToTaskModel(tasks[i])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(internalServerErrMsg)
			return
		}
	}

	err = json.NewEncoder(w).Encode(taskModels)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(encodingServerErrMsg)
		return
	}
}

func (taskAPIs TaskAPIs) FilterTasks(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	status := vars[taskStatusFilter]

	if len(status) == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(routingServerErrMsg)
		return
	}

	tasks, err := taskAPIs.taskStore.FilterTasks(taskStatusFilter, status)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(internalServerErrMsg)
		return
	}

	w.Header().Set(contentTypeKey, contentTypeVal)
	w.WriteHeader(http.StatusOK)

	taskModels := make([]models.TaskModel, len(tasks))
	for i := range tasks {
		taskModels[i], err = ToTaskModel(tasks[i])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(internalServerErrMsg)
			return
		}
	}

	err = json.NewEncoder(w).Encode(taskModels)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(encodingServerErrMsg)
		return
	}
}

func (taskAPIs TaskAPIs) StreamTasks(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	taskRespChan, err := taskAPIs.taskStore.StreamTasks(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(internalServerErrMsg)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(internalServerErrMsg)
		return
	}

	w.Header().Set(connectionKey, connectionVal)
	w.Header().Set(transferEncodingKey, transferEncodingVal)

	for taskResp := range taskRespChan {
		if taskResp.Err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(internalServerErrMsg)
			return
		}
		taskModel, err := ToTaskModel(taskResp.Task)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(internalServerErrMsg)
			return
		}
		err = json.NewEncoder(w).Encode(taskModel)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(encodingServerErrMsg)
			return
		}
		flusher.Flush()
	}

	// TODO: Handle client-side termination (Ctrl+C) using w.(http.CloseNotifier).closeNotify()
}
