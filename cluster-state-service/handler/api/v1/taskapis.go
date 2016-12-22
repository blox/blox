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

package v1

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/blox/blox/cluster-state-service/handler/api/v1/models"
	"github.com/blox/blox/cluster-state-service/handler/regex"
	"github.com/blox/blox/cluster-state-service/handler/store"
	"github.com/blox/blox/cluster-state-service/handler/types"
	"github.com/gorilla/mux"
)

const (
	taskARNKey     = "arn"
	taskClusterKey = "cluster"

	taskStatusFilter  = "status"
	taskClusterFilter = "cluster"
)

var (
	// Using maps because arrays don't support easy lookup
	supportedTaskFilters  = map[string]string{taskStatusFilter: "", taskClusterFilter: ""}
	supportedTaskStatuses = map[string]string{"pending": "", "running": "", "stopped": ""}
)

// TaskAPIs encapsulates the backend datastore with which the task APIs interact
type TaskAPIs struct {
	taskStore store.TaskStore
}

// NewTaskAPIs initializes the TaskAPIs struct
func NewTaskAPIs(taskStore store.TaskStore) TaskAPIs {
	return TaskAPIs{
		taskStore: taskStore,
	}
}

// GetTask gets a task using the cluster name to which the task belongs to and the task ARN
func (taskAPIs TaskAPIs) GetTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskARN := vars[taskARNKey]
	cluster := vars[taskClusterKey]

	if len(taskARN) == 0 || len(cluster) == 0 || !regex.IsTaskARN(taskARN) || !regex.IsClusterName(cluster) {
		http.Error(w, routingServerErrMsg, http.StatusInternalServerError)
		return
	}

	task, err := taskAPIs.taskStore.GetTask(cluster, taskARN)

	if err != nil {
		http.Error(w, internalServerErrMsg, http.StatusInternalServerError)
		return
	}

	if task == nil {
		http.Error(w, taskNotFoundClientErrMsg, http.StatusNotFound)
		return
	}

	w.Header().Set(contentTypeKey, contentTypeJSON)
	w.WriteHeader(http.StatusOK)

	extTask, err := ToTask(*task)
	if err != nil {
		http.Error(w, internalServerErrMsg, http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(extTask)
	if err != nil {
		http.Error(w, encodingServerErrMsg, http.StatusInternalServerError)
		return
	}
}

// ListTasks lists all tasks across all clusters after applying filters, if any
func (taskAPIs TaskAPIs) ListTasks(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	if taskAPIs.hasUnsupportedFilters(query) {
		http.Error(w, unsupportedFilterClientErrMsg, http.StatusBadRequest)
		return
	}

	status := query.Get(taskStatusFilter)
	cluster := query.Get(taskClusterFilter)

	// TODO: Support filtering by both status and cluster
	if status != "" && cluster != "" {
		http.Error(w, unsupportedFilterCombinationClientErrMsg, http.StatusBadRequest)
		return
	}

	if status != "" {
		if !taskAPIs.isValidStatus(status) {
			http.Error(w, invalidStatusClientErrMsg, http.StatusBadRequest)
			return
		}
	}

	if cluster != "" {
		if !regex.IsClusterARN(cluster) && !regex.IsClusterName(cluster) {
			http.Error(w, invalidClusterClientErrMsg, http.StatusBadRequest)
			return
		}
	}

	var tasks []types.Task
	var err error
	switch {
	case status != "":
		tasks, err = taskAPIs.taskStore.FilterTasks(taskStatusFilter, status)
	case cluster != "":
		tasks, err = taskAPIs.taskStore.FilterTasks(taskClusterFilter, cluster)
	default:
		tasks, err = taskAPIs.taskStore.ListTasks()
	}

	if err != nil {
		http.Error(w, internalServerErrMsg, http.StatusInternalServerError)
		return
	}

	w.Header().Set(contentTypeKey, contentTypeJSON)
	w.WriteHeader(http.StatusOK)

	extTaskItems := make([]*models.Task, len(tasks))
	for i := range tasks {
		t, err := ToTask(tasks[i])
		if err != nil {
			http.Error(w, internalServerErrMsg, http.StatusInternalServerError)
			return
		}
		extTaskItems[i] = &t
	}

	extTasks := models.Tasks{
		Items: extTaskItems,
	}

	err = json.NewEncoder(w).Encode(extTasks)
	if err != nil {
		http.Error(w, encodingServerErrMsg, http.StatusInternalServerError)
		return
	}
}

// StreamTasks streams tasks that change (status etc.) across all clusters
func (taskAPIs TaskAPIs) StreamTasks(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	taskRespChan, err := taskAPIs.taskStore.StreamTasks(ctx)
	if err != nil {
		http.Error(w, internalServerErrMsg, http.StatusInternalServerError)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, internalServerErrMsg, http.StatusInternalServerError)
		return
	}

	w.Header().Set(contentTypeKey, contentTypeStream)
	w.Header().Set(connectionKey, connectionVal)
	w.Header().Set(transferEncodingKey, transferEncodingVal)

	for taskResp := range taskRespChan {
		if taskResp.Err != nil {
			http.Error(w, internalServerErrMsg, http.StatusInternalServerError)
			return
		}
		extTask, err := ToTask(taskResp.Task)
		if err != nil {
			http.Error(w, internalServerErrMsg, http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(extTask)
		if err != nil {
			http.Error(w, encodingServerErrMsg, http.StatusInternalServerError)
			return
		}
		flusher.Flush()
	}

	// TODO: Handle client-side termination (Ctrl+C) using w.(http.CloseNotifier).closeNotify()
}

func (taskAPIs TaskAPIs) isValidStatus(status string) bool {
	_, ok := supportedTaskStatuses[status]
	return ok
}

func (taskAPIs TaskAPIs) hasUnsupportedFilters(filters map[string][]string) bool {
	if len(filters) > len(supportedTaskFilters) {
		return true
	}

	for f := range filters {
		_, ok := supportedTaskFilters[f]
		if !ok {
			return true
		}
	}
	return false
}
