package v1

import (
	"encoding/json"
	"github.com/aws/amazon-ecs-event-stream-handler/handler/store"
	"github.com/gorilla/mux"
	"net/http"
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
func (tApis TaskAPIs) GetTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskArn := vars["arn"]

	task, err := tApis.taskStore.GetTask(taskArn)

	if err != nil {
		//TODO: return http error
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(task); err != nil {
		//TODO
	}
}

func (tApis TaskAPIs) ListTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := tApis.taskStore.ListTasks()

	if err != nil {
		//TODO: return http error
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		//TODO
	}
}

func (tApis TaskAPIs) FilterTasks(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	status := vars["status"]

	tasks, err := tApis.taskStore.FilterTasks("status", status)

	if err != nil {
		//TODO: return http error
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		//TODO
	}
}

func (tApis TaskAPIs) StreamTasks(w http.ResponseWriter, r *http.Request) {
	//TODO
}
