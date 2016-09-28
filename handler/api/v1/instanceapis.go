package v1

import (
	"encoding/json"
	"github.com/aws/amazon-ecs-event-stream-handler/handler/store"
	"github.com/gorilla/mux"
	"net/http"
)

type ContainerInstanceAPIs struct {
	instanceStore store.ContainerInstanceStore
}

func NewContainerInstanceAPIs(instanceStore store.ContainerInstanceStore) ContainerInstanceAPIs {
	return ContainerInstanceAPIs{
		instanceStore: instanceStore,
	}
}

//TODO: add arn validation
func (iApis ContainerInstanceAPIs) GetInstance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	instanceArn := vars["arn"]

	instance, err := iApis.instanceStore.GetContainerInstance(instanceArn)

	if err != nil {
		//TODO: return http error
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(instance); err != nil {
		//TODO
	}
}

func (iApis ContainerInstanceAPIs) ListInstances(w http.ResponseWriter, r *http.Request) {
	instances, err := iApis.instanceStore.ListContainerInstances()

	if err != nil {
		//TODO: return http error
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(instances); err != nil {
		//TODO
	}
}

func (iApis ContainerInstanceAPIs) FilterInstances(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	status := vars["status"]

	instances, err := iApis.instanceStore.FilterContainerInstances("status", status)

	if err != nil {
		//TODO: return http error
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(instances); err != nil {
		//TODO
	}
}

func (iApis ContainerInstanceAPIs) StreamInstances(w http.ResponseWriter, r *http.Request) {
	//TODO
}
