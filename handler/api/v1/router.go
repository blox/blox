package v1

import (
	"github.com/gorilla/mux"
)

//TODO: add a map of path and query keys and use the map in task apis instead of hardcoding strings
func NewRouter(apis APIs) *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	s := r.Path("/v1").Subrouter()

	// tasks

	s.Path(`/task/{arn:(arn:aws:ecs:)([\-\w]+):[0-9]{12}:(task)\/[\-\w]+}`).
		Methods("GET").
		HandlerFunc(apis.TaskApis.GetTask)

	s.Path("/tasks").
		Methods("GET").
		HandlerFunc(apis.TaskApis.ListTasks)

	s.Path("/tasks/filter").
		Queries("status", "{status:pending|running|stopped}").
		Methods("GET").
		HandlerFunc(apis.TaskApis.FilterTasks)

	s.Path("/tasks/stream").
		Methods("GET").
		HandlerFunc(apis.TaskApis.StreamTasks)

	// instances

	s.Path(`/instance/{arn:(arn:aws:ecs:)([\-\w]+):[0-9]{12}:(container\-instance)\/[\-\w]+}`).
		Methods("GET").
		HandlerFunc(apis.ContainerInstanceApis.GetInstance)

	s.Path("/instances").
		Methods("GET").
		HandlerFunc(apis.ContainerInstanceApis.ListInstances)

	s.Path("/instances/filter").
		Queries("status", "{status:active|inactive}").
		Methods("GET").
		HandlerFunc(apis.ContainerInstanceApis.FilterInstances)

	s.Path("/instances/filter").
		Queries("cluster", "{cluster:[a-zA-Z0-9\\-_]{1,255}}").
		Methods("GET").
		HandlerFunc(apis.ContainerInstanceApis.FilterInstances)

	s.Path("/instances/filter").
		Queries("cluster", "{cluster:(arn:aws:ecs:)([\\-\\w]+):[0-9]{12}:(cluster)/[a-zA-Z0-9\\-_]{1,255}}").
		Methods("GET").
		HandlerFunc(apis.ContainerInstanceApis.FilterInstances)

	s.Path("/instances/stream").
		Methods("GET").
		HandlerFunc(apis.ContainerInstanceApis.StreamInstances)

	return s
}
