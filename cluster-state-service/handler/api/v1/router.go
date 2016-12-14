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
	"github.com/blox/blox/cluster-state-service/handler/regex"
	"github.com/gorilla/mux"
)

// TODO: add a map of path and query keys and use the map in task apis instead of hardcoding strings
var (
	// Stripping off '^' and '$' from the beginning and end of regexes respectively for the router
	clusterNameRegex = string(regex.ClusterNameRegex[1 : len(regex.ClusterNameRegex)-1])
	clusterARNRegex  = string(regex.ClusterARNRegex[1 : len(regex.ClusterARNRegex)-1])
	taskARNRegex     = string(regex.TaskARNRegex[1 : len(regex.TaskARNRegex)-1])
	instanceARNRegex = string(regex.InstanceARNRegex[1 : len(regex.InstanceARNRegex)-1])

	getTaskPath     = "/tasks/{cluster:" + clusterNameRegex + "}/{arn:" + taskARNRegex + "}"
	listTasksPath   = "/tasks"
	streamTasksPath = "/tasks/stream"

	getInstancePath     = "/instances/{cluster:" + clusterNameRegex + "}/{arn:" + instanceARNRegex + "}"
	listInstancesPath   = "/instances"
	streamInstancesPath = "/instances/stream"
)

// NewRouter initializes a new router with registered routes redirected to appropriate handler functions
func NewRouter(apis APIs) *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	s := r.Path("/v1").Subrouter()

	// Tasks

	// Get task using cluster name and task ARN
	s.Path(getTaskPath).
		Methods("GET").
		HandlerFunc(apis.TaskApis.GetTask)

	// List tasks
	s.Path(listTasksPath).
		Methods("GET").
		HandlerFunc(apis.TaskApis.ListTasks)

	// Stream tasks
	s.Path(streamTasksPath).
		Methods("GET").
		HandlerFunc(apis.TaskApis.StreamTasks)

	// Instances

	// Get instance using cluster name and instance ARN
	s.Path(getInstancePath).
		Methods("GET").
		HandlerFunc(apis.ContainerInstanceApis.GetInstance)

	// List instances
	s.Path(listInstancesPath).
		Methods("GET").
		HandlerFunc(apis.ContainerInstanceApis.ListInstances)

	// Stream instances
	s.Path(streamInstancesPath).
		Methods("GET").
		HandlerFunc(apis.ContainerInstanceApis.StreamInstances)

	return s
}
