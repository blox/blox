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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/blox/blox/daemon-scheduler/generated/v1/models"
	"github.com/blox/blox/daemon-scheduler/pkg/deployment"
	"github.com/blox/blox/daemon-scheduler/pkg/facade"
	"github.com/blox/blox/daemon-scheduler/pkg/types"
	"github.com/aws/aws-sdk-go/service/ecs"
	log "github.com/cihub/seelog"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type API struct {
	environment deployment.Environment
	deployment  deployment.Deployment
	ecs         facade.ECS
}

func NewAPI(e deployment.Environment, d deployment.Deployment, ecs facade.ECS) API {
	return API{
		environment: e,
		deployment:  d,
		ecs:         ecs,
	}
}

func (api API) Ping(w http.ResponseWriter, r *http.Request) {
	setJSONContentType(w)
	w.WriteHeader(http.StatusOK)
}

func (api API) CreateEnvironment(w http.ResponseWriter, r *http.Request) {
	var createEnvReq models.CreateEnvironmentRequest
	b, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(b, &createEnvReq)

	err := createEnvReq.Validate(nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ecsCluster, err := api.validateCluster(&createEnvReq.InstanceGroup.Cluster)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	ecsTaskDefinition, err := api.validateTaskDefinition(createEnvReq.TaskDefinition)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	env, err := api.environment.CreateEnvironment(r.Context(), *createEnvReq.Name, *ecsTaskDefinition.TaskDefinitionArn, *ecsCluster.ClusterArn)
	if err != nil {
		handleBackendError(w, err)
		return
	}

	setJSONContentType(w)
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(toEnvironmentModel(*env))
	if err != nil {
		log.Errorf("Error sending response for CreateEnvironment: %+v", err)
	}
}

func (api API) GetEnvironment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	env, err := api.environment.GetEnvironment(r.Context(), name)
	if err != nil {
		writeInternalServerError(w, err)
		return
	}

	if env == nil {
		http.Error(w, fmt.Sprintf("Environment %s does not exist", name), http.StatusNotFound)
		return
	}

	setJSONContentType(w)
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(toEnvironmentModel(*env))
	if err != nil {
		log.Errorf("Error sending response for GetEnvironment: %+v", err)
	}
}

func (api API) ListEnvironments(w http.ResponseWriter, r *http.Request) {
	envs, err := api.environment.ListEnvironments(r.Context())
	if err != nil {
		writeInternalServerError(w, err)
		return
	}

	setJSONContentType(w)
	w.WriteHeader(http.StatusOK)
	envModels := []*models.Environment{}
	for _, envType := range envs {
		envModel := toEnvironmentModel(envType)
		envModels = append(envModels, &envModel)
	}
	environments := models.Environments{
		Items: envModels,
	}
	err = json.NewEncoder(w).Encode(environments)
	if err != nil {
		log.Errorf("Error sending response for ListEnvironments: %+v", err)
	}
}

func (api API) DeleteEnvironment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	err := api.environment.DeleteEnvironment(r.Context(), name)
	if err != nil {
		_, ok := errors.Cause(err).(types.BadRequestError)
		if ok {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		writeInternalServerError(w, err)
		return
	}

	//TODO: return something when successful?
}

func (api API) CreateDeployment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	token := vars[deploymentToken]

	d, err := api.deployment.CreateDeployment(r.Context(), name, token)
	if err != nil {
		handleBackendError(w, err)
		return
	}

	setJSONContentType(w)
	w.WriteHeader(http.StatusOK)

	depModel := toDeploymentModel(&name, *d)
	err = json.NewEncoder(w).Encode(depModel)
	if err != nil {
		log.Errorf("Error sending response for CreateDeployment: %+v", err)
	}
}

func (api API) GetDeployment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	id := vars["id"]

	d, err := api.deployment.GetDeployment(r.Context(), name, id)
	if err != nil {
		handleBackendError(w, err)
		return
	}

	if d == nil {
		http.Error(w, fmt.Sprintf("Deployment %s does not exist for environment %s", id, name), http.StatusNotFound)
		return
	}

	setJSONContentType(w)
	w.WriteHeader(http.StatusOK)

	depModel := toDeploymentModel(&name, *d)
	err = json.NewEncoder(w).Encode(depModel)
	if err != nil {
		log.Errorf("Error sending response for GetDeployment: %+v", err)
	}
}

func (api API) ListDeployments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	ds, err := api.deployment.ListDeployments(r.Context(), name)
	if err != nil {
		handleBackendError(w, err)
		return
	}

	setJSONContentType(w)
	w.WriteHeader(http.StatusOK)

	deploymentsModel := toDeploymentsModel(&name, ds)
	err = json.NewEncoder(w).Encode(deploymentsModel)
	if err != nil {
		log.Errorf("Error sending response for ListDeployments: %+v", err)
	}
}

func (api API) validateCluster(clusterName *string) (*ecs.Cluster, error) {
	cluster, err := api.ecs.DescribeCluster(clusterName)
	if err != nil {
		return nil, err
	}

	if *cluster.Status == "INACTIVE" {
		return nil, errors.New("Cluster is inactive")
	}

	return cluster, nil
}

func (api API) validateTaskDefinition(td *string) (*ecs.TaskDefinition, error) {
	taskDefinition, err := api.ecs.DescribeTaskDefinition(td)
	if err != nil {
		return nil, err
	}

	if *taskDefinition.Status == "INACTIVE" {
		return nil, errors.New("TaskDefinition is inactive")
	}

	return taskDefinition, nil
}
