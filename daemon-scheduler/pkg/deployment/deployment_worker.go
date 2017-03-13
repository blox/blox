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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/blox/blox/daemon-scheduler/pkg/environment"
	deploymenttypes "github.com/blox/blox/daemon-scheduler/pkg/deployment/types"
	environmenttypes "github.com/blox/blox/daemon-scheduler/pkg/environment/types"
	"github.com/blox/blox/daemon-scheduler/pkg/facade"
	"github.com/blox/blox/daemon-scheduler/pkg/types"
	log "github.com/cihub/seelog"
	"github.com/pkg/errors"
)

const (
	TaskPending = "PENDING"
)

type DeploymentWorker interface {
	StartPendingDeployment(ctx context.Context, environmentName string) (*deploymenttypes.Deployment, error)
	// UpdateInProgressDeployment checks for in-progress deployments and moves them to complete when
	// the tasks started by the deployment have moved out of pending status
	UpdateInProgressDeployment(ctx context.Context, environmentName string) (*deploymenttypes.Deployment, error)
}

type deploymentWorker struct {
	environmentService environment.EnvironmentService
	environmentFacade  environment.EnvironmentFacade
	deploymentService  DeploymentService
	ecs                facade.ECS
	css                facade.ClusterState
}

func NewDeploymentWorker(
	environmentService environment.EnvironmentService,
	environmentFacade environment.EnvironmentFacade,
	deploymentService DeploymentService,
	ecs facade.ECS,
	css facade.ClusterState) DeploymentWorker {
	return deploymentWorker{
		environmentService: environmentService,
		environmentFacade:  environmentFacade,
		deploymentService:  deploymentService,
		ecs:                ecs,
		css:                css,
	}
}

func (d deploymentWorker) StartPendingDeployment(ctx context.Context,
	environmentName string) (*deploymenttypes.Deployment, error) {

	if environmentName == "" {
		return nil, errors.New("Environment name is missing")
	}

	environment, err := d.environmentService.GetEnvironment(ctx, environmentName)
	if err != nil {
		return nil, errors.Wrapf(err, "Error finding environment with name %s", environmentName)
	}

	if environment == nil {
		return nil, nil
	}

	instances, err := d.environmentFacade.InstanceARNs(environment)
	if err != nil {
		return nil, err
	}

	startedDeployment, err := d.deploymentService.StartDeployment(ctx, environmentName, instances)
	if err != nil {
		return nil, err
	}

	return startedDeployment, nil
}

func (d deploymentWorker) UpdateInProgressDeployment(ctx context.Context,
	environmentName string) (*deploymenttypes.Deployment, error) {

	if environmentName == "" {
		return nil, errors.New("Environment name is missing")
	}

	deployment, err := d.deploymentService.GetInProgressDeployment(ctx, environmentName)
	if err != nil {
		return nil, err
	}

	if deployment == nil {
		return nil, nil
	}

	environment, err := d.environmentService.GetEnvironment(ctx, environmentName)
	if err != nil {
		return nil, errors.Wrapf(err, "Error finding environment with name %s", environmentName)
	}

	if environment == nil {
		return nil, nil
	}

	taskProgress, err := d.checkDeploymentTaskProgress(environment, deployment)
	if err != nil {
		return nil, errors.Wrapf(err, "Error checking deployment %s progress in environment %s",
			deployment.ID, environment.Name)
	}

	updatedDeployment, err := d.updateDeployment(ctx, environment.Name, deployment, taskProgress)
	if err != nil {
		return nil, err
	}

	return updatedDeployment, nil
}

func (d deploymentWorker) checkDeploymentTaskProgress(environment *environmenttypes.Environment,
	deployment *deploymenttypes.Deployment) (*ecs.DescribeTasksOutput, error) {

	if environment.Cluster == "" {
		return nil, errors.New("Environment cluster should not be empty")
	}

	// TODO: replace with cluster state calls
	tasks, err := d.ecs.ListTasks(environment.Cluster, deployment.ID)
	if err != nil {
		return nil, err
	}

	resp, err := d.ecs.DescribeTasks(environment.Cluster, tasks)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (d deploymentWorker) updateDeployment(ctx context.Context,
	environmentName string, deployment *deploymenttypes.Deployment,
	resp *ecs.DescribeTasksOutput) (*deploymenttypes.Deployment, error) {

	updatedDeployment, err := d.updateDeploymentObject(deployment, resp)
	if err != nil {
		return nil, err
	}

	err = d.deploymentService.UpdateInProgressDeployment(ctx, environmentName, updatedDeployment)

	if err != nil {
		if _, ok := errors.Cause(err).(types.UnexpectedDeploymentStatusError); ok {
			// deployment updated is no longer the in-progress deployment of the environment
			log.Infof("Deployment %s is no longer the in-progress deployment", updatedDeployment.ID)
			return nil, nil
		}
		return nil, errors.Wrapf(err, "Error updating the deployment %v in the environment %v",
			*updatedDeployment, environmentName)
	}

	return updatedDeployment, nil
}

func (d deploymentWorker) updateDeploymentObject(deployment *deploymenttypes.Deployment,
	resp *ecs.DescribeTasksOutput) (*deploymenttypes.Deployment, error) {

	if d.deploymentCompleted(resp.Tasks, resp.Failures) {
		err := deployment.UpdateDeploymentToCompleted(resp.Failures)
		if err != nil {
			return nil, err
		}

		return deployment, nil
	}

	err := deployment.UpdateDeploymentToInProgress(deployment.DesiredTaskCount, resp.Failures)
	if err != nil {
		return nil, err
	}

	return deployment, nil
}

func (d deploymentWorker) deploymentCompleted(tasks []*ecs.Task, failures []*ecs.Failure) bool {
	if len(tasks) == 0 {
		return false
	}

	for _, t := range tasks {
		if aws.StringValue(t.LastStatus) == TaskPending {
			return false
		}
	}

	return true
}
