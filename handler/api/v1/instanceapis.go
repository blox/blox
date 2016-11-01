package v1

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/amazon-ecs-event-stream-handler/handler/api/v1/models"
	"github.com/aws/amazon-ecs-event-stream-handler/handler/regex"
	"github.com/aws/amazon-ecs-event-stream-handler/handler/store"
	"github.com/aws/amazon-ecs-event-stream-handler/handler/types"
	"github.com/gorilla/mux"
)

const (
	instanceARNKey     = "arn"
	instanceClusterKey = "cluster"

	instanceStatusFilter  = "status"
	instanceClusterFilter = "cluster"
)

type ContainerInstanceAPIs struct {
	instanceStore store.ContainerInstanceStore
}

func NewContainerInstanceAPIs(instanceStore store.ContainerInstanceStore) ContainerInstanceAPIs {
	return ContainerInstanceAPIs{
		instanceStore: instanceStore,
	}
}

func (instanceAPIs ContainerInstanceAPIs) GetInstance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	instanceARN := vars[instanceARNKey]
	cluster := vars[instanceClusterKey]

	if len(instanceARN) == 0 || len(cluster) == 0 || !regex.IsInstanceARN(instanceARN) || !regex.IsClusterName(cluster) {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(routingServerErrMsg)
		return
	}

	instance, err := instanceAPIs.instanceStore.GetContainerInstance(cluster, instanceARN)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(internalServerErrMsg)
		return
	}

	if instance == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(instanceNotFoundClientErrMsg)
		return
	}

	instanceModel, err := ToContainerInstanceModel(*instance)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(internalServerErrMsg)
		return
	}

	w.Header().Set(contentTypeKey, contentTypeVal)
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(instanceModel)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(encodingServerErrMsg)
		return
	}
}

func (instanceAPIs ContainerInstanceAPIs) ListInstances(w http.ResponseWriter, r *http.Request) {
	instances, err := instanceAPIs.instanceStore.ListContainerInstances()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(internalServerErrMsg)
		return
	}

	w.Header().Set(contentTypeKey, contentTypeVal)
	w.WriteHeader(http.StatusOK)

	instanceModels := make([]models.ContainerInstanceModel, len(instances))
	for i := range instances {
		instanceModels[i], err = ToContainerInstanceModel(instances[i])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(internalServerErrMsg)
			return
		}
	}
	err = json.NewEncoder(w).Encode(instanceModels)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(encodingServerErrMsg)
		return
	}
}

func (instanceAPIs ContainerInstanceAPIs) FilterInstances(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	status := vars[instanceStatusFilter]
	cluster := vars[instanceClusterFilter]

	if len(status) != 0 && len(cluster) != 0 {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(routingServerErrMsg)
		return
	}

	var instances []types.ContainerInstance
	var err error

	switch {
	case len(status) != 0:
		instances, err = instanceAPIs.instanceStore.FilterContainerInstances(instanceStatusFilter, status)
	case len(cluster) != 0:
		instances, err = instanceAPIs.instanceStore.FilterContainerInstances(instanceClusterFilter, cluster)
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

	instanceModels := make([]models.ContainerInstanceModel, len(instances))
	for i := range instances {
		instanceModels[i], err = ToContainerInstanceModel(instances[i])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(internalServerErrMsg)
			return
		}
	}

	err = json.NewEncoder(w).Encode(instanceModels)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(encodingServerErrMsg)
		return
	}
}

func (instanceAPIs ContainerInstanceAPIs) StreamInstances(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	instanceRespChan, err := instanceAPIs.instanceStore.StreamContainerInstances(ctx)
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

	w.Header().Set(connectionKey, contentTypeVal)
	w.Header().Set(transferEncodingKey, transferEncodingVal)

	for instanceResp := range instanceRespChan {
		if instanceResp.Err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(internalServerErrMsg)
			return
		}
		instanceModel, err := ToContainerInstanceModel(instanceResp.ContainerInstance)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(internalServerErrMsg)
			return
		}
		err = json.NewEncoder(w).Encode(instanceModel)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(encodingServerErrMsg)
			return
		}
		flusher.Flush()
	}

	// TODO: Handle client-side termination (Ctrl+C) using w.(http.CloseNotifier).closeNotify()
}
