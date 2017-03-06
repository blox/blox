// Copyright 2016-2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.
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
	"strings"

	"github.com/blox/blox/daemon-scheduler/pkg/store"
	storetypes "github.com/blox/blox/daemon-scheduler/pkg/store/types"
	"github.com/blox/blox/daemon-scheduler/pkg/types"
	"github.com/blox/blox/daemon-scheduler/pkg/validate"
	log "github.com/cihub/seelog"
	"github.com/pkg/errors"
)

const (
	clusterFilter = "cluster"
)

var (
	supportedFilterKeys = []string{clusterFilter}
)

// Environment defines methods to handle environments
type Environment interface {
	// CreateEnvironment stores a new environment in the database
	CreateEnvironment(ctx context.Context, name string, taskDefinition string, cluster string) (*types.Environment, error)
	// GetEnvironment gets the environment with the provided name from the database
	GetEnvironment(ctx context.Context, name string) (*types.Environment, error)
	// DeleteEnvironment deletes the environment with the provided name from the database
	DeleteEnvironment(ctx context.Context, name string) error
	// ListEnvironments returns a list with all the existing environments
	ListEnvironments(ctx context.Context) ([]types.Environment, error)
	// FilterEnvironments returns a list of all environments that match the filters
	FilterEnvironments(ctx context.Context, filterKey string, filterVal string) ([]types.Environment, error)

	// This is meant to be a 'private' method to be called by exported methods. Adding it to the interface for the purpose of testing.
	// TODO: Change these to unexported methods. Currently unable to do so because mocking unexported methods with gomock fails
	// (https://github.com/golang/mock/issues/52).
	// ValidateAndCreateEnvironment is a generator function for use by CreateEnvironment().
	// It validates that the environment does not already exist before creating the new environment.
	ValidateAndCreateEnvironment(newEnv *types.Environment) storetypes.ValidateAndUpdateEnvironment
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

	environment, err := types.NewEnvironment(name, taskDefinition, cluster)
	if err != nil {
		return nil, err
	}

	err = e.environmentStore.PutEnvironment(ctx, name, e.ValidateAndCreateEnvironment(environment))

	if err != nil {
		return nil, errors.Wrapf(err, "Error saving environment %s to store", name)
	}

	return environment, nil
}

func (e environment) ValidateAndCreateEnvironment(newEnv *types.Environment) storetypes.ValidateAndUpdateEnvironment {
	return func(existingEnv *types.Environment) (*types.Environment, error) {
		if existingEnv == nil {
			return newEnv, nil
		}
		log.Errorf("An environment with name %s already exists", existingEnv.Name)
		return nil, types.NewBadRequestError(errors.Errorf("An environment with name %s already exists", existingEnv.Name))
	}
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

	err := e.environmentStore.DeleteEnvironment(ctx, name)
	if err != nil {
		return errors.Wrapf(err, "Error deleting environment %s from store", name)
	}

	return nil
}

func (e environment) ListEnvironments(ctx context.Context) ([]types.Environment, error) {
	//TODO: should we sort the deployments by time before returning?
	return e.environmentStore.ListEnvironments(ctx)
}

func (e environment) FilterEnvironments(ctx context.Context, filterKey string, filterVal string) ([]types.Environment, error) {
	if filterKey == "" {
		return nil, errors.New("Filter key is missing")
	}
	if filterVal == "" {
		return nil, errors.New("Filter value is missing")
	}
	switch filterKey {
	case clusterFilter:
		return e.filterEnvironmentsByCluster(ctx, filterVal)
	default:
		return nil, errors.Errorf("Unsupported filter key '%s'. Supported filters are '%v'", filterKey, supportedFilterKeys)
	}
}

func (e environment) filterEnvironmentsByCluster(ctx context.Context, cluster string) ([]types.Environment, error) {
	if validate.IsClusterARN(cluster) {
		return e.filterEnvironmentsByClusterARN(ctx, cluster)
	}
	if validate.IsClusterName(cluster) {
		return e.filterEnvironmentsByClusterName(ctx, cluster)
	}
	return nil, errors.Errorf("'%s' is neither a cluster name nor a cluster ARN", cluster)
}

func (e environment) filterEnvironmentsByClusterARN(ctx context.Context, clusterARN string) ([]types.Environment, error) {
	envs, err := e.environmentStore.ListEnvironments(ctx)
	if err != nil {
		return nil, err
	}

	filteredEnvs := make([]types.Environment, 0, len(envs))
	for _, env := range envs {
		if clusterARN == env.Cluster {
			filteredEnvs = append(filteredEnvs, env)
		}
	}
	return filteredEnvs, nil
}

func (e environment) filterEnvironmentsByClusterName(ctx context.Context, clusterName string) ([]types.Environment, error) {
	envs, err := e.environmentStore.ListEnvironments(ctx)
	if err != nil {
		return nil, err
	}

	filteredEnvs := make([]types.Environment, 0, len(envs))
	for _, env := range envs {
		clusterARNSuffix := "/" + clusterName
		if strings.HasSuffix(env.Cluster, clusterARNSuffix) {
			filteredEnvs = append(filteredEnvs, env)
		}
	}
	return filteredEnvs, nil
}
