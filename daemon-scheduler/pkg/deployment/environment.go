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

package deployment

import (
	"context"

	"github.com/blox/blox/daemon-scheduler/pkg/store"
	"github.com/blox/blox/daemon-scheduler/pkg/types"
	log "github.com/cihub/seelog"
	"github.com/pkg/errors"
)

type Environment interface {
	// CreateEnvironment stores a new environment in the database
	CreateEnvironment(ctx context.Context, name string, taskDefinition string, cluster string) (*types.Environment, error)
	// GetEnvironment gets the environment with the provided name from the database
	GetEnvironment(ctx context.Context, name string) (*types.Environment, error)
	// DeleteEnvironment deletes the environment with the provided name from the database
	DeleteEnvironment(ctx context.Context, name string) error
	// ListEnvironments returns a list with all the existing environments
	ListEnvironments(ctx context.Context) ([]types.Environment, error)

	// AddDeployment adds a deployment to the environment if a deployment with
	// the provided ID does not exist
	AddDeployment(ctx context.Context, environment types.Environment, deployment types.Deployment) (*types.Environment, error)
	// UpdateDeployment replaces an existing deployment in the environment with the
	// provided one if a deployment with the provided ID already exists
	UpdateDeployment(ctx context.Context, environment types.Environment, deployment types.Deployment) (*types.Environment, error)
	//GetCurrentDeployment returns the deployment which needs to be used for starting tasks
	GetCurrentDeployment(ctx context.Context, name string) (*types.Deployment, error)
}

type environment struct {
	environmentStore store.EnvironmentStore
}

func NewEnvironment(environmentStore store.EnvironmentStore) (Environment, error) {
	if environmentStore == nil {
		return nil, errors.New("Environment is not initialized")
	}
	return environment{
		environmentStore: environmentStore,
	}, nil
}

func (e environment) CreateEnvironment(ctx context.Context,
	name string, taskDefinition string, cluster string) (*types.Environment, error) {

	if len(name) == 0 {
		return nil, errors.New("Environment name is missing")
	}

	if len(taskDefinition) == 0 {
		return nil, errors.New("Environment task definition is missing")
	}

	if len(cluster) == 0 {
		return nil, errors.New("Environment cluster is missing")
	}

	env, err := e.GetEnvironment(ctx, name)
	if err != nil {
		return nil, errors.Wrapf(err, "Error getting environment with name %s", name)
	}

	if env != nil {
		log.Errorf("An environment with name %s already exists", name)
		return nil, types.NewBadRequestError(errors.Errorf("An environment with name %s already exists", name))
	}

	environment, err := types.NewEnvironment(name, taskDefinition, cluster)
	if err != nil {
		return nil, err
	}

	err = e.environmentStore.PutEnvironment(ctx, *environment)
	if err != nil {
		return nil, errors.Wrapf(err, "Error saving environment %s to store", name)
	}

	return environment, nil
}

func (e environment) GetEnvironment(ctx context.Context, name string) (*types.Environment, error) {
	if len(name) == 0 {
		return nil, types.NewBadRequestError(errors.New("Environment name is missing"))
	}

	//TODO: should we sort the deployments by time before returning?
	env, err := e.environmentStore.GetEnvironment(ctx, name)
	if err != nil {
		return nil, errors.Wrapf(err, "Error loading environment %s from store", name)
	}

	return env, nil
}

func (e environment) DeleteEnvironment(ctx context.Context, name string) error {
	if len(name) == 0 {
		return types.NewBadRequestError(errors.New("Environment name is missing"))
	}

	env, err := e.environmentStore.GetEnvironment(ctx, name)
	if err != nil {
		return err
	}

	if env == nil {
		log.Infof("Environment %s does not exist", name)
		return nil
	}

	err = e.environmentStore.DeleteEnvironment(ctx, *env)
	if err != nil {
		return errors.Wrapf(err, "Error deleting environment %s from store", name)
	}

	return nil
}

func (e environment) ListEnvironments(ctx context.Context) ([]types.Environment, error) {
	//TODO: should we sort the deployments by time before returning?
	return e.environmentStore.ListEnvironments(ctx)
}

func (e environment) AddDeployment(ctx context.Context, environment types.Environment,
	deployment types.Deployment) (*types.Environment, error) {

	if len(deployment.ID) == 0 {
		return nil, errors.New("Deployment id cannot be empty")
	}

	_, ok := environment.Deployments[deployment.ID]
	if ok {
		return nil, errors.Errorf("Deployment %s already exists: %+v", deployment.ID, deployment)
	}

	environment.Deployments[deployment.ID] = deployment
	environment.PendingDeploymentID = deployment.ID

	err := e.environmentStore.PutEnvironment(ctx, environment)
	if err != nil {
		return nil, errors.Wrapf(err, "Error saving environment %s to store", environment.Name)
	}

	return &environment, nil
}

func (e environment) UpdateDeployment(ctx context.Context, environment types.Environment,
	deployment types.Deployment) (*types.Environment, error) {

	if len(deployment.ID) == 0 {
		return nil, errors.New("Deployment id cannot be empty")
	}

	_, ok := environment.Deployments[deployment.ID]
	if !ok {
		return nil, errors.Errorf("Deployment %s does not exist", deployment.ID)
	}

	// replace deployment with updated version
	environment.Deployments[deployment.ID] = deployment
	environment.DesiredTaskCount = deployment.DesiredTaskCount

	if deployment.Health == types.DeploymentHealthy {
		environment.Health = types.EnvironmentHealthy
	} else {
		environment.Health = types.EnvironmentUnhealthy
	}

	err := e.environmentStore.PutEnvironment(ctx, environment)
	if err != nil {
		return nil, err
	}

	return &environment, nil
}

func (e environment) GetCurrentDeployment(ctx context.Context, name string) (*types.Deployment, error) {
	env, err := e.GetEnvironment(ctx, name)
	if err != nil {
		return nil, err
	}

	if env == nil {
		return nil, types.NewNotFoundError(errors.Errorf("Environment %s is not available", name))
	}

	return env.GetCurrentDeployment()
}
