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

	extTask, err := ToTask(*task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(internalServerErrMsg)
		return
	}

	err = json.NewEncoder(w).Encode(extTask)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(encodingServerErrMsg)
		return
	}
}

// ListTasks lists all tasks across all clusters
func (taskAPIs TaskAPIs) ListTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := taskAPIs.taskStore.ListTasks()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(internalServerErrMsg)
		return
	}

	w.Header().Set(contentTypeKey, contentTypeVal)
	w.WriteHeader(http.StatusOK)

	extTaskItems := make([]*models.Task, len(tasks))
	for i := range tasks {
		t, err := ToTask(tasks[i])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(internalServerErrMsg)
			return
		}
		extTaskItems[i] = &t
	}

	extTasks := models.Tasks{
		Items: extTaskItems,
	}

	err = json.NewEncoder(w).Encode(extTasks)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(encodingServerErrMsg)
		return
	}
}

// FilterTasks filters tasks across all clusters by status
func (taskAPIs TaskAPIs) FilterTasks(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	status := vars[taskStatusFilter]
	cluster := vars[taskClusterFilter]

	if len(status) != 0 && len(cluster) != 0 {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(routingServerErrMsg)
		return
	}

	var tasks []types.Task
	var err error

	switch {
	case len(status) != 0:
		tasks, err = taskAPIs.taskStore.FilterTasks(taskStatusFilter, status)
	case len(cluster) != 0:
		tasks, err = taskAPIs.taskStore.FilterTasks(taskClusterFilter, cluster)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(routingServerErrMsg)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(internalServerErrMsg)
		return
	}

	w.Header().Set(contentTypeKey, contentTypeVal)
	w.WriteHeader(http.StatusOK)

	extTaskItems := make([]*models.Task, len(tasks))
	for i := range tasks {
		t, err := ToTask(tasks[i])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(internalServerErrMsg)
			return
		}
		extTaskItems[i] = &t
	}

	extTasks := models.Tasks{
		Items: extTaskItems,
	}

	err = json.NewEncoder(w).Encode(extTasks)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(encodingServerErrMsg)
		return
	}
}

// StreamTasks streams tasks that change (status etc.) across all clusters
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
		extTask, err := ToTask(taskResp.Task)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(internalServerErrMsg)
			return
		}
		err = json.NewEncoder(w).Encode(extTask)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(encodingServerErrMsg)
			return
		}
		flusher.Flush()
	}

	// TODO: Handle client-side termination (Ctrl+C) using w.(http.CloseNotifier).closeNotify()
}
