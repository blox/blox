package v1

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/amazon-ecs-event-stream-handler/handler/store"
	"github.com/gorilla/mux"
)

const (
	contentTypeKey      = "Content-Type"
	contentTypeVal      = "application/json; charset=UTF-8"
	connectionKey       = "Connection"
	connectionVal       = "Keep-Alive"
	transferEncodingKey = "Transfer-Encoding"
	transferEncodingVal = "chunked"
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
	taskArn := vars["arn"]

	task, err := taskAPIs.taskStore.GetTask(taskArn)

	if err != nil {
		//TODO: return http error
	}

	w.Header().Set(contentTypeKey, contentTypeVal)
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(task); err != nil {
		//TODO
	}
}

func (taskAPIs TaskAPIs) ListTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := taskAPIs.taskStore.ListTasks()

	if err != nil {
		//TODO: return http error
	}

	w.Header().Set(contentTypeKey, contentTypeVal)
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		//TODO
	}
}

func (taskAPIs TaskAPIs) FilterTasks(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	status := vars["status"]

	tasks, err := taskAPIs.taskStore.FilterTasks("status", status)

	if err != nil {
		//TODO: return http error
	}

	w.Header().Set(contentTypeKey, contentTypeVal)
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		//TODO
	}
}

func (taskAPIs TaskAPIs) StreamTasks(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	taskRespChan, err := taskAPIs.taskStore.StreamTasks(ctx)
	if err != nil {
		//TODO
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		//TODO
	}

	w.Header().Set(connectionKey, connectionVal)
	w.Header().Set(transferEncodingKey, transferEncodingVal)

	for taskResp := range taskRespChan {
		if taskResp.Err != nil {
			//TODO
		}
		if err := json.NewEncoder(w).Encode(taskResp.Task); err != nil {
			//TODO
		}
		flusher.Flush()
	}

	// TODO: Handle client-side termination (Ctrl+C) using w.(http.CloseNotifier).closeNotify()
}
