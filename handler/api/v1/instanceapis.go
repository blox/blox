package v1

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/amazon-ecs-event-stream-handler/handler/store"
	"github.com/aws/amazon-ecs-event-stream-handler/handler/types"
	"github.com/gorilla/mux"
)

const (
	statusFilter  = "status"
	clusterFilter = "cluster"
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
func (instanceAPIs ContainerInstanceAPIs) GetInstance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	instanceArn := vars["arn"]

	instance, err := instanceAPIs.instanceStore.GetContainerInstance(instanceArn)

	if err != nil {
		//TODO: return http error
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(instance); err != nil {
		//TODO
	}
}

func (instanceAPIs ContainerInstanceAPIs) ListInstances(w http.ResponseWriter, r *http.Request) {
	instances, err := instanceAPIs.instanceStore.ListContainerInstances()

	if err != nil {
		//TODO: return http error
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(instances); err != nil {
		//TODO
	}
}

func (instanceAPIs ContainerInstanceAPIs) FilterInstances(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	status := vars[statusFilter]
	cluster := vars[clusterFilter]

	if len(status) != 0 && len(cluster) != 0 {
		// TODO: return http error
	}

	var instances []types.ContainerInstance
	var err error

	switch {
	case len(status) != 0:
		instances, err = instanceAPIs.instanceStore.FilterContainerInstances(statusFilter, status)
	case len(cluster) != 0:
		instances, err = instanceAPIs.instanceStore.FilterContainerInstances(clusterFilter, cluster)
	default:
		// TODO: return http error
	}

	if err != nil {
		//TODO: return http error
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(instances); err != nil {
		//TODO
	}
}

func (instanceAPIs ContainerInstanceAPIs) StreamInstances(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	instanceRespChan, err := instanceAPIs.instanceStore.StreamContainerInstances(ctx)
	if err != nil {
		//TODO
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		//TODO
	}

	w.Header().Set(connectionKey, connectionVal)
	w.Header().Set(transferEncodingKey, transferEncodingVal)

	for instanceResp := range instanceRespChan {
		if instanceResp.Err != nil {
			//TODO
		}
		if err := json.NewEncoder(w).Encode(instanceResp.ContainerInstance); err != nil {
			//TODO
		}
		flusher.Flush()
	}

	// TODO: Handle client-side termination (Ctrl+C) using w.(http.CloseNotifier).closeNotify()
}
