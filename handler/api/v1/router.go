package v1

import (
	"github.com/aws/amazon-ecs-event-stream-handler/handler/regex"
	"github.com/gorilla/mux"
)

// TODO: add a map of path and query keys and use the map in task apis instead of hardcoding strings
const (
	getTaskPath     = "/tasks/{cluster:" + regex.ClusterNameRegex + "}/{arn:" + regex.TaskARNRegex + "}"
	listTasksPath   = "/tasks"
	filterTasksPath = "/tasks/filter"
	streamTasksPath = "/tasks/stream"

	getInstancePath     = "/instances/{cluster:" + regex.ClusterNameRegex + "}/{arn:" + regex.InstanceARNRegex + "}"
	listInstancesPath   = "/instances"
	filterInstancesPath = "/instances/filter"
	streamInstancesPath = "/instances/stream"

	clusterKey     = "cluster"
	clusterNameVal = "{" + clusterKey + ":" + regex.ClusterNameRegex + "}"
	clusterARNVal  = "{" + clusterKey + ":" + regex.ClusterARNRegex + "}"

	taskKey    = "task"
	taskARNVal = "{" + taskKey + ":" + regex.TaskARNRegex + "}"

	instanceKey    = "instance"
	instanceARNVal = "{" + instanceKey + ":" + regex.InstanceARNRegex + "}"

	statusKey         = "status"
	taskStatusVal     = "{" + statusKey + ":pending|running|stopped}"
	instanceStatusVal = "{" + statusKey + ":active|inactive}"
)

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

	// Filter tasks by status
	s.Path(filterTasksPath).
		Queries(statusKey, taskStatusVal).
		Methods("GET").
		HandlerFunc(apis.TaskApis.FilterTasks)

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

	// Filter instances by status
	s.Path(filterInstancesPath).
		Queries(statusKey, instanceStatusVal).
		Methods("GET").
		HandlerFunc(apis.ContainerInstanceApis.FilterInstances)

	// Filter instances by cluser name
	s.Path(filterInstancesPath).
		Queries(clusterKey, clusterNameVal).
		Methods("GET").
		HandlerFunc(apis.ContainerInstanceApis.FilterInstances)

	// Filter instances by cluster ARN
	s.Path(filterInstancesPath).
		Queries(clusterKey, clusterARNVal).
		Methods("GET").
		HandlerFunc(apis.ContainerInstanceApis.FilterInstances)

	// Stream instances
	s.Path(streamInstancesPath).
		Methods("GET").
		HandlerFunc(apis.ContainerInstanceApis.StreamInstances)

	return s
}
