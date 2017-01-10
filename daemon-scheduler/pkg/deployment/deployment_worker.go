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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/blox/blox/daemon-scheduler/pkg/facade"
	"github.com/blox/blox/daemon-scheduler/pkg/types"
	"github.com/pkg/errors"
)

const (
	TaskPending = "PENDING"
)

type DeploymentWorker interface {
	// UpdateInProgressDeployment checks for in-progress deployments and moves them to complete when
	// the tasks started by the deployment have moved out of pending status
	UpdateInProgressDeployment(ctx context.Context, environmentName string) (*types.Deployment, error)
}

type deploymentWorker struct {
	environment Environment
	deployment  Deployment
	ecs         facade.ECS
	css         facade.ClusterState
}

func NewDeploymentWorker(
	environment Environment,
	deployment Deployment,
	ecs facade.ECS,
	css facade.ClusterState) DeploymentWorker {
	return deploymentWorker{
		environment: environment,
		deployment:  deployment,
		ecs:         ecs,
		css:         css,
	}
}

func (d deploymentWorker) UpdateInProgressDeployment(ctx context.Context,
	environmentName string) (*types.Deployment, error) {

	if environmentName == "" {
		return nil, errors.New("Environment name is missing")
	}

	deployment, err := d.deployment.GetInProgressDeployment(ctx, environmentName)
	if err != nil {
		return nil, err
	}

	if deployment == nil {
		return nil, nil
	}

	environment, err := d.environment.GetEnvironment(ctx, environmentName)
	if err != nil {
		return nil, errors.Wrapf(err, "Error finding environment with name %s", environmentName)
	}

	if environment == nil {
		return nil, nil
	}

	updatedDeployment, err := d.updateDeployment(environment, deployment)
	if err != nil {
		return nil, errors.Wrap(err, "Error updating the deployment")
	}

	// retrieve in-progress again to make sure it has not been updated by another process
	// TODO: wrap the in-progress check and updateDeployment in a transaction
	deployment, err = d.deployment.GetInProgressDeployment(ctx, environmentName)
	if err != nil {
		return nil, err
	}

	if deployment == nil || deployment.ID != updatedDeployment.ID {
		return nil, nil
	}

	_, err = d.environment.UpdateDeployment(ctx, *environment, *updatedDeployment)
	if err != nil {
		return nil, errors.Wrapf(err, "Error updating the deployment %v in the environment %v",
			*updatedDeployment, environment.Name)
	}

	return updatedDeployment, nil
}

func (d deploymentWorker) updateDeployment(environment *types.Environment,
	deployment *types.Deployment) (*types.Deployment, error) {

	if environment.Cluster == "" {
		return nil, errors.New("Environment cluster should not be empty")
	}

	tasks, err := d.ecs.ListTasks(environment.Cluster, deployment.ID)
	if err != nil {
		return nil, err
	}

	resp, err := d.ecs.DescribeTasks(environment.Cluster, tasks)
	if err != nil {
		return nil, err
	}

	if d.deploymentCompleted(resp.Tasks, resp.Failures) {
		return deployment.UpdateDeploymentCompleted(resp.Failures)
	}

	updatedDeployment, err := deployment.UpdateDeploymentInProgress(
		deployment.DesiredTaskCount, resp.Failures)
	if err != nil {
		return nil, err
	}

	return updatedDeployment, nil
}

func (d deploymentWorker) deploymentCompleted(tasks []*ecs.Task, failures []*ecs.Failure) bool {
	for _, t := range tasks {
		if aws.StringValue(t.LastStatus) == TaskPending {
			return false
		}
	}

	return true
}
