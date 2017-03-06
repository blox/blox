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

	"github.com/blox/blox/daemon-scheduler/pkg/facade"
	"github.com/blox/blox/daemon-scheduler/pkg/store"
	storetypes "github.com/blox/blox/daemon-scheduler/pkg/store/types"
	"github.com/blox/blox/daemon-scheduler/pkg/types"
	"github.com/pkg/errors"
)

// Deployment defines methods to handle deployments in an environment
// TODO: refactor to remove multiple environment retrievals from the database in some of the Get methods
// (check unit tests for GetEnvironment.Times(2 or 3))
type Deployment interface {
	// CreateDeployment creates a new deployment in the provided environment and updates the
	// environment's pending deployment ID to the ID of the deployment created. The environment
	// token must match the provided token, otherwise the deployment creation will fail.
	CreateDeployment(ctx context.Context, environmentName string, token string) (*types.Deployment, error)

	// CreateSubDeployment kicks off a deployment corresponding to the in progress deployment ID
	// in the environment to start tasks on given instances
	CreateSubDeployment(ctx context.Context, environmentName string, instanceARNs []*string) (*types.Deployment, error)

	// StartDeployment kicks off the provided deployment in the given environment (by starting tasks)
	StartDeployment(ctx context.Context, environmentName string, instanceARNs []*string) (*types.Deployment, error)

	// UpdateInProgressDeployment replaces an existing deployment in the environment with the
	// provided one if a deployment with the provided ID already exists
	UpdateInProgressDeployment(ctx context.Context, environmentName string, deployment *types.Deployment) error

	// GetDeployment returns the deployment with the provided id in the provided environment
	GetDeployment(ctx context.Context, environmentName string, id string) (*types.Deployment, error)

	// GetCurrentDeployment returns the deployment which needs to be used for starting tasks, i.e.
	// the in-progress deployment for the environment if one exists, otherwise the latest completed deployment.
	GetCurrentDeployment(ctx context.Context, environmentName string) (*types.Deployment, error)

	// GetPendingDeployment returns the pending deployment for the environment deployment.
	// There should be no more than one pending deployments in an environment.
	GetPendingDeployment(ctx context.Context, environmentName string) (*types.Deployment, error)

	// GetInProgressDeployment returns the in-progress deployment for the environmentName.
	// There should be no more than one in progress deployments in an environment.
	GetInProgressDeployment(ctx context.Context, environmentName string) (*types.Deployment, error)

	// ListDeploymentsSortedReverseChronologically returns a list of deployments reverse-ordered by start time,
	// i.e. lastest deployment first
	ListDeploymentsSortedReverseChronologically(ctx context.Context, environmentName string) ([]types.Deployment, error)

	// The followung functions are meant to be 'private' methods to be called by exported methods.
	// Adding it to the interface for the purpose of testing.
	// TODO: Change these to unexported methods. Currently unable to do so because mocking unexported methods with gomock fails
	// (https://github.com/golang/mock/issues/52).

	// ValidateAndCreateDeployment is a generator function for use by CreateDeployment().
	// It validates the environment corresponding to the deployment to be created and adds a
	// pending deployment to the environment if the validations succeed
	ValidateAndCreateDeployment(token string) (storetypes.ValidateAndUpdateEnvironment, *types.Deployment)

	// ValidateAndCreateSubDeployment is a generator function for use by CreateSubDeployment().
	// It validates the environment corresponding to the sub-deployment to be created. If the validations
	// succeed, tasks are started for the current deployment of the environment and the deployment information
	// is updated in the environment.
	ValidateAndCreateSubDeployment(instanceARNs []*string) (storetypes.ValidateAndUpdateEnvironment, *types.Deployment)

	// ValidateAndStartDeployment is a generator function for use by StartDeployment().
	// It validates the environment corresponding to the deployment to be started. If the validations
	// succeed, tasks are started for the pending deployment of the environment and the deployment information
	// is updated in the environment.
	ValidateAndStartDeployment(instanceARNs []*string) (storetypes.ValidateAndUpdateEnvironment, *types.Deployment)

	// ValidateAndUpdateInProgressDeployment is a generator function for use by UpdateInProgressDeployment().
	// It validates the environment corresponding to the deployment to be updated. If the validations
	// succeed, the deployment is updated using the deployment being passed in in the environment.
	ValidateAndUpdateInProgressDeployment(deployment *types.Deployment) storetypes.ValidateAndUpdateEnvironment
}

type deployment struct {
	environmentStore store.EnvironmentStore
	clusterState     facade.ClusterState
	ecs              facade.ECS
}

func NewDeployment(
	environmentStore store.EnvironmentStore,
	clusterState facade.ClusterState,
	ecs facade.ECS) Deployment {

	return deployment{
		environmentStore: environmentStore,
		clusterState:     clusterState,
		ecs:              ecs,
	}
}

func (d deployment) CreateDeployment(ctx context.Context,
	environmentName string, token string) (*types.Deployment, error) {

	if environmentName == "" {
		return nil, types.NewBadRequestError(errors.New("Environment name is missing when creating a deployment"))
	}

	if token == "" {
		return nil, types.NewBadRequestError(errors.New("Token is missing when creating a deployment"))
	}

	validateAndCreateDep, deployment := d.ValidateAndCreateDeployment(token)
	err := d.environmentStore.PutEnvironment(ctx, environmentName, validateAndCreateDep)

	if err != nil {
		return nil, err
	}

	return deployment, nil
}

func (d deployment) ValidateAndCreateDeployment(token string) (storetypes.ValidateAndUpdateEnvironment, *types.Deployment) {
	deployment := &types.Deployment{}
	return func(existingEnv *types.Environment) (*types.Environment, error) {
		if existingEnv == nil {
			return nil, types.NewNotFoundError(errors.Errorf("Environment with name '%s' does not exist", environmentName))
		}

		err := d.verifyToken(*existingEnv, token)
		if err != nil {
			return nil, err
		}

		err = d.verifyNoInProgressDeploymentExists(*existingEnv)
		if err != nil {
			return nil, err
		}

		// Create and add a pending deployment to the environment
		dep, err := types.NewDeployment(existingEnv.DesiredTaskDefinition, existingEnv.Token)
		if err != nil {
			return nil, err
		}
		err = existingEnv.AddPendingDeployment(*dep)
		if err != nil {
			return nil, err
		}

		*deployment = *dep

		return existingEnv, nil
	}, deployment
}

func (d deployment) verifyToken(env types.Environment, token string) error {
	if len(token) > 0 && env.Token != token {
		return types.NewBadRequestError(errors.Errorf("Token '%s' is outdated and does not match the environment token '%s'", token, env.Token))
	}

	for _, deployment := range env.Deployments {
		if deployment.Token == token {
			return types.NewBadRequestError(errors.Errorf("Deployment with token '%s' already exists", token))
		}
	}

	return nil
}

func (d deployment) verifyNoInProgressDeploymentExists(env types.Environment) error {
	if env.InProgressDeploymentID == "" {
		return nil
	}
	return types.NewBadRequestError(errors.Errorf("There is already a deployment in progress: %s", env.InProgressDeploymentID))
}

func (d deployment) CreateSubDeployment(ctx context.Context, environmentName string, instanceARNs []*string) (*types.Deployment, error) {
	if environmentName == "" {
		return nil, errors.New("Environment name is missing when creating a deployment")
	}

	validateAndCreateSubDep, deployment := d.ValidateAndCreateSubDeployment(instanceARNs)
	err := d.environmentStore.PutEnvironment(ctx, environmentName, validateAndCreateSubDep)

	if err != nil {
		return nil, err
	}
	return deployment, nil
}

func (d deployment) ValidateAndCreateSubDeployment(instanceARNs []*string) (storetypes.ValidateAndUpdateEnvironment, *types.Deployment) {
	deployment := &types.Deployment{}
	return func(existingEnv *types.Environment) (*types.Environment, error) {
		if existingEnv == nil {
			return nil, types.NewNotFoundError(errors.Errorf("Environment with name '%s' does not exist", environmentName))
		}

		curDeployment, intErr := d.getCurrentDeploymentFromEnvironment(existingEnv)
		if intErr != nil {
			return nil, errors.Wrapf(intErr,
				"Unable to retrieve deployment for environment with name '%s' to create a sub-deployment", environmentName)
		}
		if curDeployment == nil {
			return nil, errors.Errorf(
				"There is no deployment for environment with name '%s' to create a sub-deployment", environmentName)
		}

		dep, intErr := d.startTasksAndUpdateDeploymentInfo(existingEnv, curDeployment, instanceARNs)
		if intErr != nil {
			return nil, intErr
		}

		*deployment = *dep

		return existingEnv, nil
	}, deployment
}

func (d deployment) getCurrentDeploymentFromEnvironment(env *types.Environment) (*types.Deployment, error) {
	// If there is an in-progress deployment, use that as the current deployment
	if env.InProgressDeploymentID != "" {
		inProgress, ok := env.Deployments[env.InProgressDeploymentID]
		if !ok {
			return nil, errors.Errorf("Deployment with ID '%s' does not exist in the deployments for environment with name '%s'",
				env.InProgressDeploymentID, env.Name)
		}

		return &inProgress, nil
	}

	// If there is no in-progress deployment, use the latest completed deployment
	deployments, err := env.SortDeploymentsReverseChronologically()
	if err != nil {
		return nil, err
	}

	for _, d := range deployments {
		if d.Status == types.DeploymentCompleted {
			return &d, nil
		}
	}

	return nil, nil
}

func (d deployment) startTasksAndUpdateDeploymentInfo(env *types.Environment, deployment *types.Deployment, instanceARNs []*string) (*types.Deployment, error) {
	resp, err := d.ecs.StartTask(env.Cluster, instanceARNs, deployment.ID, deployment.TaskDefinition)
	if err != nil {
		return nil, errors.Wrapf(
			err, "Error starting tasks for deployment with ID '%s' in environment with name '%s'", deployment.ID)
	}

	failures := resp.Failures
	if deployment.FailedInstances != nil {
		failures = append(failures, deployment.FailedInstances...)
	}

	// if deployment is already completed (in the case where a new instance joins
	// and there are no in-progress deployments) then we do not update the deployment object
	// TODO: Figure out how we want to track failures in sub-deployments
	if deployment.Status == types.DeploymentCompleted {
		return deployment, nil
	}

	if deployment.Status == types.DeploymentPending {
		err = deployment.UpdateDeploymentToInProgress(len(instanceARNs), failures)
		if err != nil {
			return nil, err
		}
		err = env.UpdatePendingDeploymentToInProgress()
		if err != nil {
			return nil, err
		}
	}

	if deployment.Status == types.DeploymentInProgress {
		// just update succeeded and failed instances
		err = deployment.UpdateDeploymentToInProgress(len(instanceARNs), failures)
		if err != nil {
			return nil, err
		}
	}

	// replace deployment with updated version
	env.Deployments[deployment.ID] = *deployment
	env.DesiredTaskCount = deployment.DesiredTaskCount

	if deployment.Health == types.DeploymentHealthy {
		env.Health = types.EnvironmentHealthy
	} else {
		env.Health = types.EnvironmentUnhealthy
	}

	return deployment, nil
}

func (d deployment) StartDeployment(ctx context.Context, environmentName string, instanceARNs []*string) (*types.Deployment, error) {
	if environmentName == "" {
		return nil, types.NewBadRequestError(errors.New("Environment name is missing when starting a deployment"))
	}

	validateAndStartDep, deployment := d.ValidateAndStartDeployment(instanceARNs)
	err := d.environmentStore.PutEnvironment(ctx, environmentName, validateAndStartDep)

	if err != nil {
		return nil, err
	}

	return deployment, nil
}

func (d deployment) ValidateAndStartDeployment(instanceARNs []*string) (storetypes.ValidateAndUpdateEnvironment, *types.Deployment) {
	deployment := &types.Deployment{}
	return func(existingEnv *types.Environment) (*types.Environment, error) {
		if existingEnv == nil {
			return nil, types.NewNotFoundError(errors.Errorf("Environment with name '%s' does not exist", environmentName))
		}

		if existingEnv.InProgressDeploymentID != "" {
			return nil, errors.Errorf("There is already a deployment in-progress '%s'", existingEnv.InProgressDeploymentID)
		}

		if existingEnv.PendingDeploymentID == "" {
			return nil, errors.New("There is no pending deployment")
		}

		pendingDeployment, ok := existingEnv.Deployments[existingEnv.PendingDeploymentID]
		if !ok {
			return nil, errors.Errorf("Deployment with ID '%s' does not exist in the deployments for environment with name '%s'",
				existingEnv.PendingDeploymentID, existingEnv.Name)
		}

		dep, intErr := d.startTasksAndUpdateDeploymentInfo(existingEnv, &pendingDeployment, instanceARNs)
		if intErr != nil {
			return nil, intErr
		}

		*deployment = *dep

		return existingEnv, nil
	}, deployment
}

func (d deployment) UpdateInProgressDeployment(ctx context.Context, environmentName string, deployment *types.Deployment) error {
	if environmentName == "" {
		return types.NewBadRequestError(errors.New("Environment name is missing when updating a deployment"))
	}
	if deployment.ID == "" {
		return errors.New("Deployment id cannot be empty")
	}

	validateAndUpdateDep := d.ValidateAndUpdateInProgressDeployment(deployment)
	err := d.environmentStore.PutEnvironment(ctx, environmentName, validateAndUpdateDep)

	return err
}

func (d deployment) ValidateAndUpdateInProgressDeployment(deployment *types.Deployment) storetypes.ValidateAndUpdateEnvironment {
	return func(existingEnv *types.Environment) (*types.Environment, error) {
		if existingEnv == nil {
			return nil, types.NewNotFoundError(errors.Errorf("Environment with name '%s' does not exist", environmentName))
		}

		if existingEnv.InProgressDeploymentID != deployment.ID {
			return nil, types.NewUnexpectedDeploymentStatusError(errors.Errorf("The in-progress deployment of environment with name '%s' is '%s' and not '%s'",
				environmentName, existingEnv.InProgressDeploymentID, deployment.ID))
		}

		_, ok := existingEnv.Deployments[deployment.ID]
		if !ok {
			return nil, errors.Errorf("Deployment with ID '%s' does not exist in environment with name '%s'", deployment.ID, environmentName)
		}

		// replace deployment with updated version
		existingEnv.Deployments[deployment.ID] = *deployment
		existingEnv.DesiredTaskCount = deployment.DesiredTaskCount

		if deployment.Health == types.DeploymentHealthy {
			existingEnv.Health = types.EnvironmentHealthy
		} else {
			existingEnv.Health = types.EnvironmentUnhealthy
		}

		return existingEnv, nil
	}
}

func (d deployment) GetDeployment(ctx context.Context,
	environmentName string, id string) (*types.Deployment, error) {

	if environmentName == "" {
		return nil, errors.New("Environment name is missing when getting a deployment")
	}

	if id == "" {
		return nil, errors.New("ID is missing when getting a deployment")
	}

	deployments, err := d.ListDeploymentsSortedReverseChronologically(ctx, environmentName)
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

func (d deployment) getEnvironmentOrFailIfDoesNotExist(ctx context.Context, name string) (*types.Environment, error) {
	env, err := d.environmentStore.GetEnvironment(ctx, name)
	if err != nil {
		return nil, err
	}

	if env == nil {
		return nil, types.NewNotFoundError(errors.Errorf("Environment %s does not exist", name))
	}

	return env, nil
}

func (d deployment) GetCurrentDeployment(ctx context.Context, environmentName string) (*types.Deployment, error) {
	if len(environmentName) == 0 {
		return nil, types.NewBadRequestError(errors.New("Environment name is missing"))
	}

	deployment, err := d.GetInProgressDeployment(ctx, environmentName)
	if err != nil {
		return nil, err
	}

	if deployment != nil {
		return deployment, nil
	}

	env, err := d.getEnvironmentOrFailIfDoesNotExist(ctx, environmentName)
	if err != nil {
		return nil, err
	}

	// if there is no in-progress deployment then we take the latest completed deployment
	deployments, err := env.SortDeploymentsReverseChronologically()
	if err != nil {
		return nil, err
	}

	for _, d := range deployments {
		if d.Status == types.DeploymentCompleted {
			return &d, nil
		}
	}

	return nil, nil
}

func (d deployment) GetPendingDeployment(ctx context.Context, environmentName string) (*types.Deployment, error) {
	if len(environmentName) == 0 {
		return nil, types.NewBadRequestError(errors.New("Environment name is missing"))
	}

	env, err := d.getEnvironmentOrFailIfDoesNotExist(ctx, environmentName)
	if err != nil {
		return nil, err
	}

	if env.PendingDeploymentID == "" {
		return nil, nil
	}

	pending, ok := env.Deployments[env.PendingDeploymentID]
	if !ok {
		return nil, errors.Errorf("Deployment with ID '%s' does not exist in the deployments for environment with name '%s'",
			env.PendingDeploymentID, env.Name)
	}

	return &pending, nil
}

func (d deployment) GetInProgressDeployment(ctx context.Context, environmentName string) (*types.Deployment, error) {
	if len(environmentName) == 0 {
		return nil, types.NewBadRequestError(errors.New("Environment name is missing"))
	}

	env, err := d.getEnvironmentOrFailIfDoesNotExist(ctx, environmentName)
	if err != nil {
		return nil, err
	}

	if env.InProgressDeploymentID == "" {
		return nil, nil
	}

	inProgress, ok := env.Deployments[env.InProgressDeploymentID]
	if !ok {
		return nil, errors.Errorf("Deployment with ID '%s' does not exist in the deployments for environment with name '%s'",
			env.InProgressDeploymentID, env.Name)
	}

	return &inProgress, nil
}

func (d deployment) ListDeploymentsSortedReverseChronologically(ctx context.Context, environmentName string) ([]types.Deployment, error) {
	if len(environmentName) == 0 {
		return nil, types.NewBadRequestError(errors.New("Environment name is missing"))
	}

	env, err := d.getEnvironmentOrFailIfDoesNotExist(ctx, environmentName)
	if err != nil {
		return nil, err
	}

	return env.SortDeploymentsReverseChronologically()
}
