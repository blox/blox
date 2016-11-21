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
	"fmt"

	"github.com/blox/blox/daemon-scheduler/pkg/facade"
	"github.com/blox/blox/daemon-scheduler/pkg/types"
	"github.com/pkg/errors"
)

type Deployment interface {
	// CreateDeployment kicks off a deployment in the provided environment. If token is provided
	// the environment token must match the provided token, otherwise the deployment will fail.
	CreateDeployment(ctx context.Context, environmentName string, token string) (*types.Deployment, error)
	// CreateSubDeployment kicks off latest deployment to start tasks on given instances
	CreateSubDeployment(ctx context.Context, environmentName string, instanceARNs []*string) (*types.Deployment, error)
	// GetDeployment returns the deployment with the provided id in the provided environment
	GetDeployment(ctx context.Context, environmentName string, id string) (*types.Deployment, error)
	// ListDeployments returns a list of all the deployments in the provided environment
	ListDeployments(ctx context.Context, environmentName string) ([]types.Deployment, error)
}

type deployment struct {
	environment  Environment
	clusterState facade.ClusterState
	ecs          facade.ECS
}

func NewDeployment(
	environment Environment,
	clusterState facade.ClusterState,
	ecs facade.ECS) Deployment {

	return deployment{
		environment:  environment,
		clusterState: clusterState,
		ecs:          ecs,
	}
}

func (d deployment) CreateDeployment(ctx context.Context,
	environmentName string, token string) (*types.Deployment, error) {

	if len(environmentName) == 0 {
		return nil, errors.New("Environment name is missing when creating a deployment")
	}

	env, err := d.getEnvironment(ctx, environmentName)
	if err != nil {
		return nil, errors.Wrapf(err, "Error retrieving environment with name %s", environmentName)
	}

	if len(token) > 0 && env.Token != token {
		return nil, errors.Errorf("Token %v is outdated. Token on the environment is %v",
			token, env.Token)
	}

	// create and add a pending deployment to the environment
	deployment, err := types.NewDeployment(env.DesiredTaskDefinition)
	if err != nil {
		return nil, err
	}

	env, err = d.environment.AddDeployment(ctx, *env, *deployment)
	if err != nil {
		return nil, errors.Wrapf(err, "Error adding deployment %v to environment %s", deployment, environmentName)
	}

	inprogress, err := env.GetInProgressDeployment()
	if err != nil {
		return nil, err
	}

	if inprogress != nil {
		//TODO: add deployment to the pending deployment queue
		return deployment, nil
	}

	//TODO: switch back to using css when bootstrapping is done
	//instanceARNs, err := d.getInstanceARNs(*env)
	instanceARNs, err := d.getECSClusterInstanceArns(*env)
	if err != nil {
		return nil, errors.Wrapf(err, "Error getting instances in cluster for environment %s", environmentName)
	}

	if len(instanceARNs) == 0 {
		return nil, errors.Errorf("There are no instances found that match the cluster: %s", env.Cluster)
	}

	return d.startDeployment(ctx, env, deployment, instanceARNs)
}

func (d deployment) CreateSubDeployment(ctx context.Context, environmentName string, instanceARNs []*string) (*types.Deployment, error) {
	if len(environmentName) == 0 {
		return nil, errors.New("Environment name is missing when creating a deployment")
	}

	env, err := d.getEnvironment(ctx, environmentName)
	if err != nil {
		return nil, errors.Wrapf(err, "Error retrieving environment with name %s", environmentName)
	}

	deployment, err := d.environment.GetCurrentDeployment(ctx, environmentName)

	if err != nil {
		return nil, errors.Wrap(err, "Unable to create sub-deployment")
	}

	if deployment == nil {
		return nil, fmt.Errorf("No deployment found for environment %s", environmentName)
	}

	return d.startDeployment(ctx, env, deployment, instanceARNs)
}

func (d deployment) startDeployment(ctx context.Context, env *types.Environment, deployment *types.Deployment, instanceARNs []*string) (*types.Deployment, error) {
	//TODO: Deprecate this in favor of sub-deployments/activities. It is safer to leave activities immutable after they are marked completed
	//Also move this to async

	resp, err := d.ecs.StartTask(env.Cluster, instanceARNs, deployment.ID, deployment.TaskDefinition)
	if err != nil {
		return nil, errors.Wrap(err, "Error starting tasks")
	}

	tasks := resp.Tasks
	if deployment.CurrentTasks != nil {
		tasks = append(tasks, deployment.CurrentTasks...)
	}
	failures := resp.Failures
	if deployment.FailedInstances != nil {
		failures = append(failures, deployment.FailedInstances...)
	}
	updatedDeployment, err := deployment.UpdateDeploymentInProgress(
		len(instanceARNs), tasks, failures)

	if err != nil {
		return nil, errors.Wrap(err, "Error updating deployment")
	}

	env, err = d.environment.UpdateDeployment(ctx, *env, *updatedDeployment)
	if err != nil {
		return nil, errors.Wrap(err, "Error updating deployment in environment")
	}

	return updatedDeployment, nil
}

func (d deployment) GetDeployment(ctx context.Context,
	environmentName string, id string) (*types.Deployment, error) {

	if len(environmentName) == 0 {
		return nil, errors.New("Environment name is missing when getting a deployment")
	}

	if len(id) == 0 {
		return nil, errors.New("ID is missing when getting a deployment")
	}

	deployments, err := d.getEnvironmentDeployments(ctx, environmentName)
	if err != nil {
		return nil, err
	}

	for _, deployment := range deployments {
		if deployment.ID == id {
			return &deployment, nil
		}
	}

	return nil, nil
}

func (d deployment) ListDeployments(ctx context.Context,
	environmentName string) ([]types.Deployment, error) {

	if len(environmentName) == 0 {
		return nil, errors.New("Environment name is missing when listing deployments")
	}

	deployments, err := d.getEnvironmentDeployments(ctx, environmentName)
	if err != nil {
		return nil, err
	}

	return deployments, nil
}

func (d deployment) getEnvironmentDeployments(ctx context.Context,
	environmentName string) ([]types.Deployment, error) {

	env, err := d.getEnvironment(ctx, environmentName)
	if err != nil {
		return nil, err
	}

	deployments, err := env.GetDeployments()
	if err != nil {
		return nil, err
	}

	return deployments, nil
}

func (d deployment) getECSClusterInstanceArns(env types.Environment) ([]*string, error) {
	return d.ecs.ListInstances(env.Cluster)
}

func (d deployment) getInstanceARNs(env types.Environment) ([]*string, error) {
	instances, err := d.clusterState.ListInstances(env.Cluster)
	if err != nil {
		return nil, err
	}

	instanceARNs := make([]*string, 0, len(instances))
	for _, v := range instances {
		instanceARNs = append(instanceARNs, v.ContainerInstanceARN)
	}

	return instanceARNs, nil
}

func (d deployment) getEnvironment(ctx context.Context,
	environmentName string) (*types.Environment, error) {

	env, err := d.environment.GetEnvironment(ctx, environmentName)
	if err != nil {
		return nil, errors.Wrapf(err, "Error finding environment with name %s", environmentName)
	}

	if env == nil {
		return nil, types.NewNotFoundError(fmt.Errorf("Environment with name %s is missing", environmentName))
	}

	return env, err
}
