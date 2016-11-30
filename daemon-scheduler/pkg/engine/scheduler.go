// Copyright 2016 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the License). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the license file accompanying this file. This file is distributed
// on an AS IS BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package engine

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/blox/blox/daemon-scheduler/pkg/clients/css/models"
	"github.com/blox/blox/daemon-scheduler/pkg/deployment"
	"github.com/blox/blox/daemon-scheduler/pkg/facade"
	"github.com/blox/blox/daemon-scheduler/pkg/types"
	log "github.com/cihub/seelog"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
)

const (
	schedulerTickerDuration = 10 * time.Second
	inactiveInstanceStatus  = "INACTIVE"
	runningTaskStatus       = "RUNNING"
	trackingInfoTTL         = 1 * time.Minute
)

// Scheduler loops through all the environments and makes sure they reach their eventual state.
type Scheduler struct {
	id             string
	ctx            context.Context
	environmentSvc deployment.Environment
	deploymentSvc  deployment.Deployment
	css            facade.ClusterState
	ecs            facade.ECS
	events         chan<- Event
	executionState map[string]environmentExecutionState
	inProgress     bool
	inProgressLock sync.RWMutex
}

type environmentExecutionState struct {
	environment    types.Environment
	trackingInfo   map[string]time.Time
	inProgress     bool
	inProgressLock sync.RWMutex
}

type instanceLookupResult struct {
	totalInstanceCount int
	newInstances       []*string
	deployedInstances  map[string][]*deployedTask
}

type deployedTask struct {
	taskARN                 string
	instanceARN             string
	deploymentID            string
	availableInClusterState bool
}

func NewScheduler(ctx context.Context, events chan<- Event, environmentSvc deployment.Environment,
	deploymentSvc deployment.Deployment, css facade.ClusterState, ecs facade.ECS) *Scheduler {
	return &Scheduler{
		id:             uuid.NewRandom().String(),
		ctx:            ctx,
		environmentSvc: environmentSvc,
		deploymentSvc:  deploymentSvc,
		css:            css,
		ecs:            ecs,
		events:         events,
		executionState: make(map[string]environmentExecutionState),
		inProgress:     false,
	}
}

func (scheduler *Scheduler) Start() {
	ticker := time.NewTicker(schedulerTickerDuration)
	go func(scheduler *Scheduler) {
		scheduler.runOnce()
		for {
			select {
			case <-ticker.C:
				scheduler.runOnce()
			case <-scheduler.ctx.Done():
				log.Infof("[s:%s] Shutting down scheduler", scheduler.id)
				ticker.Stop()
				return
			}
		}
	}(scheduler)
}

func (scheduler *Scheduler) runOnce() {
	if scheduler.isInProgress() {
		msg := fmt.Sprintf("[s:%s] Another instance of scheduler is already in progress, skipping", scheduler.id)
		log.Info(msg)
		scheduler.events <- SchedulerErrorEvent{
			Error: errors.Errorf(msg),
		}
		return
	}

	go func(scheduler *Scheduler) {
		err := scheduler.runOnceInternal()
		if err != nil {
			log.Errorf("[s:%s] Error running scheduler : %v", scheduler.id, err)
			scheduler.events <- SchedulerErrorEvent{
				Error: err,
			}
		}
	}(scheduler)
}

// runOnceInternal runs a single iteration of Scheduler.
func (scheduler *Scheduler) runOnceInternal() error {
	scheduler.setInProgress(true)
	defer scheduler.setInProgress(false)

	environments, err := scheduler.environmentSvc.ListEnvironments(scheduler.ctx)
	if err != nil {
		return errors.Wrapf(err, "Error getting environments", scheduler.id)
	}

	for _, environment := range environments {
		_, ok := scheduler.executionState[environment.Name]
		if !ok {
			scheduler.executionState[environment.Name] = environmentExecutionState{
				environment:  environment,
				trackingInfo: make(map[string]time.Time),
				inProgress:   false,
			}
		}

		// processing for environments is independent, so we can do them concurrently
		go func(scheduler *Scheduler, environment types.Environment) {
			state := scheduler.executionState[environment.Name]
			err := scheduler.runForEnvironment(&state)
			if err != nil {
				// TODO: we may want to report this for better ux
				log.Errorf("[s:%s, e:%s] Error running this iteration of Scheduler for environment : %v", scheduler.id, state.environment.Name, err)
				scheduler.events <- SchedulerErrorEvent{
					Error:       errors.Wrapf(err, "Error running scheduler for environment %s", state.environment.Name),
					Environment: state.environment,
				}
				return
			}
			msg := fmt.Sprintf("[s:%s, e:%s] Done running this iteration of scheduler for environment", scheduler.id, state.environment.Name)
			log.Info(msg)
			scheduler.events <- SchedulerEnvironmentEvent{
				Message:     msg,
				Environment: state.environment,
			}
		}(scheduler, environment)
	}

	return nil
}

func (scheduler *Scheduler) setInProgress(val bool) {
	scheduler.inProgressLock.Lock()
	defer scheduler.inProgressLock.Unlock()
	scheduler.inProgress = val
}

func (scheduler *Scheduler) isInProgress() bool {
	scheduler.inProgressLock.RLock()
	defer scheduler.inProgressLock.RUnlock()
	return scheduler.inProgress
}

func (state *environmentExecutionState) setInProgress(val bool) {
	lock := state.inProgressLock
	lock.Lock()
	defer lock.Unlock()
	state.inProgress = val
}

func (state *environmentExecutionState) isInProgress() bool {
	lock := state.inProgressLock
	lock.RLock()
	defer lock.RUnlock()
	return state.inProgress
}

func (scheduler *Scheduler) runForEnvironment(state *environmentExecutionState) error {
	environment := state.environment
	if state.isInProgress() {
		log.Infof("[s:%s, e:%s] Execution for environment is already in progress", scheduler.id, environment.Name)
		return nil
	}
	state.setInProgress(true)
	defer state.setInProgress(false)

	log.Debugf("[s:%s, e:%s] Instances tracked under environment = %d", scheduler.id, environment.Name, len(state.trackingInfo))

	currentDeployment, err := scheduler.getCurrentDeployment(&environment)
	if err != nil {
		return err
	}

	if currentDeployment == nil {
		return errors.Errorf("No deployment available for environment")
	}

	lookupResult, err := scheduler.lookupInstances(state)
	if err != nil {
		return errors.Wrapf(err, "Error finding instances to deploy for environment")
	}

	log.Debugf("[s:%s, e:%s] Instance lookup result: new=%d, deployed=%d, total=%d",
		scheduler.id, environment.Name, len(lookupResult.newInstances), len(lookupResult.deployedInstances), lookupResult.totalInstanceCount)

	scheduler.deployToNewInstances(state, lookupResult)

	err = scheduler.updateDeployedInstances(state, currentDeployment, lookupResult)
	if err != nil {
		return errors.Wrapf(err, "Error updating deployed instances for environment")
	}

	return nil
}

func (scheduler *Scheduler) getCurrentDeployment(environment *types.Environment) (*types.Deployment, error) {
	deployment, err := scheduler.environmentSvc.GetCurrentDeployment(scheduler.ctx, environment.Name)
	if err != nil {
		return nil, errors.Wrapf(err, "Error getting current deployment for cluster %s of environment", environment.Cluster)
	}

	return deployment, nil
}

// updateDeployedInstances performs deployment on instances which already have some version of environment deployed
func (scheduler *Scheduler) updateDeployedInstances(state *environmentExecutionState, currentDeployment *types.Deployment, result *instanceLookupResult) error {
	environment := state.environment
	// go through already deployed instances and select the tasks which need to be replaced
	for instanceARN, deployedTasks := range result.deployedInstances {
		shouldDeploy := true
		tasksToStop := make([]string, 0)
		for _, dt := range deployedTasks {
			if dt.availableInClusterState {
				if dt.deploymentID == currentDeployment.ID {
					//now that we know this deployment made it to instance we can safely delete from our tracking data
					delete(state.trackingInfo, instanceARN)
					//NOTE: in a pathological scenario there can be > 1 tasks corresponding
					//to current deployment running, we will stop all but one of those
					if shouldDeploy {
						shouldDeploy = false
						continue
					}
				}
				log.Debugf("[s:%s, e:%s] Adding task %s to stop tasks list", scheduler.id, environment.Name, dt.taskARN)
				tasksToStop = append(tasksToStop, dt.taskARN)
			} else {
				deployedAt, ok := state.trackingInfo[instanceARN]
				if ok {
					//if deployment happened a while ago and
					//we haven't heard from cluster-state we
					//ask ECS if the deployment succeeded
					if time.Now().UTC().Sub(deployedAt) > trackingInfoTTL {
						deployed, err := scheduler.isDeployedToInstance(state, currentDeployment, instanceARN)
						if err != nil {
							return err
						}
						shouldDeploy = !deployed
					} else {
						shouldDeploy = false
					}
				}
			}
		}

		// order is to stop existing task(s) and start new one
		if len(tasksToStop) > 0 {
			log.Debugf("[s:%s, e:%s] Sending StopTasksEvent with %d tasks", scheduler.id, environment.Name, len(tasksToStop))
			scheduler.events <- StopTasksEvent{
				Cluster:     environment.Cluster,
				Tasks:       tasksToStop,
				Environment: environment,
			}
		}

		if shouldDeploy {
			log.Debugf("[s:%s, e:%s] Sending StartDeploymentEvent for deployment %s to instance %s",
				scheduler.id, environment.Name, currentDeployment.ID, instanceARN)
			state.trackingInfo[instanceARN] = time.Now().UTC()
			scheduler.events <- StartDeploymentEvent{
				Environment: environment,
				Instances:   []*string{aws.String(instanceARN)},
			}
		}
	}

	return nil
}

func (scheduler *Scheduler) isDeployedToInstance(state *environmentExecutionState, currentDeployment *types.Deployment, instanceARN string) (bool, error) {
	taskARNs, err := scheduler.ecs.ListTasksByInstance(state.environment.Cluster, instanceARN)
	if err != nil {
		return false, errors.Wrapf(err, "Error listing tasks for instance %s in cluster %s for environment", instanceARN, state.environment.Cluster)
	}

	if len(taskARNs) > 0 {
		output, err := scheduler.ecs.DescribeTasks(state.environment.Cluster, taskARNs)
		if err != nil {
			return false, errors.Wrapf(err, "Error describing tasks in cluster %s for environment", state.environment.Cluster)
		}

		for _, task := range output.Tasks {
			if aws.StringValue(task.StartedBy) == currentDeployment.ID {
				return true, nil
			}
		}
	}

	return false, nil
}

// deployToNewInstances performs deployment on instances which never got any deployment for the given environment
func (scheduler *Scheduler) deployToNewInstances(state *environmentExecutionState, result *instanceLookupResult) {
	if len(result.newInstances) > 0 {
		for _, instanceARN := range result.newInstances {
			state.trackingInfo[aws.StringValue(instanceARN)] = time.Now().UTC()
		}
		event := StartDeploymentEvent{
			Environment: state.environment,
			Instances:   result.newInstances,
		}
		scheduler.events <- event
		log.Debugf("Sent event to start tasks on %d instances", len(result.newInstances))
	}
}

// lookupInstances returns instanceLookupResult struct containing the state of all instances in the cluster corresponding to environment
func (scheduler *Scheduler) lookupInstances(state *environmentExecutionState) (*instanceLookupResult, error) {
	environment := state.environment

	// get all instances in Cluster
	instances, err := scheduler.css.ListInstances(environment.Cluster)
	if err != nil {
		return nil, errors.Wrapf(err, "Error getting instances for cluster %s of environment", environment.Cluster)
	}

	instanceARNToInstance := make(map[string]*models.ContainerInstance)
	for _, instance := range instances {
		instanceARNToInstance[aws.StringValue(instance.ContainerInstanceARN)] = instance
	}

	result := &instanceLookupResult{
		totalInstanceCount: 0,
		newInstances:       make([]*string, 0),
		deployedInstances:  make(map[string][]*deployedTask),
	}

	result, err = scheduler.loadInstancesAlreadyDeployed(state, instanceARNToInstance, result)
	if err != nil {
		return nil, errors.Wrapf(err, "Error finding instances where environment is already deployed")
	}

	// collect all the instances which do not have this environment installed
	for _, i := range instances {
		instanceARN := aws.StringValue(i.ContainerInstanceARN)
		if aws.StringValue(i.Status) == inactiveInstanceStatus {
			delete(result.deployedInstances, instanceARN)
			continue
		}
		result.totalInstanceCount++
		_, ok := result.deployedInstances[instanceARN]
		if !ok {
			result.newInstances = append(result.newInstances, i.ContainerInstanceARN)
		}
	}
	return result, nil
}

// loadInstancesAlreadyDeployed populates instanceLookupResult struct with the state of instances derived from
// state of environments, deployments and cluster
func (scheduler *Scheduler) loadInstancesAlreadyDeployed(state *environmentExecutionState,
	instanceARNToInstance map[string]*models.ContainerInstance,
	result *instanceLookupResult) (*instanceLookupResult, error) {

	environment := state.environment

	tasks, err := scheduler.getRunningTasks(environment.Cluster)
	if err != nil {
		return result, err
	}

	deployments, err := scheduler.deploymentSvc.ListDeployments(scheduler.ctx, environment.Name)
	if err != nil {
		return result, errors.Wrapf(err, "Error calling ListDeployments with environment")

	}

	// preparing a map for easy lookup
	deploymentsMap := make(map[string]*types.Deployment)
	for _, d := range deployments {
		deploymentsMap[d.ID] = &d

	}

	// for each task find the deployment it corresponds to and tag the instance of the task as deployed
	for _, task := range tasks {
		// ignore if task does not belong to this environment
		_, ok := deploymentsMap[task.StartedBy]
		if !ok {
			continue
		}
		instanceARN := aws.StringValue(task.ContainerInstanceARN)
		instance, ok := instanceARNToInstance[instanceARN]
		if !ok || aws.StringValue(instance.Status) == inactiveInstanceStatus {
			continue
		}

		deployedTasks, ok := result.deployedInstances[instanceARN]
		if !ok {
			deployedTasks = make([]*deployedTask, 0)
		}

		deployedTasks = append(deployedTasks, &deployedTask{
			instanceARN:             instanceARN,
			taskARN:                 aws.StringValue(task.TaskARN),
			deploymentID:            task.StartedBy,
			availableInClusterState: true,
		})

		result.deployedInstances[instanceARN] = deployedTasks
	}

	// Also add tasks which are not yet available in cluster-state, this happens when events are delayed.
	for instanceARN, _ := range state.trackingInfo {
		deployedTasks, ok := result.deployedInstances[instanceARN]
		if !ok {
			deployedTasks = make([]*deployedTask, 0)
		}
		deployedTasks = append(deployedTasks, &deployedTask{
			instanceARN:             instanceARN,
			availableInClusterState: false,
		})
		result.deployedInstances[instanceARN] = deployedTasks
	}

	return result, nil
}

// getRunningTasks returns a map of taskARN -> task where task is -probably- running
func (scheduler *Scheduler) getRunningTasks(cluster string) (map[string]*models.Task, error) {
	resp, err := scheduler.css.ListTasks(cluster)
	if err != nil {
		return nil, errors.Wrapf(err, "Error getting tasks for cluster %s", cluster)
	}

	tasks := make(map[string]*models.Task)
	for _, task := range resp {
		if aws.StringValue(task.DesiredStatus) == runningTaskStatus {
			tasks[aws.StringValue(task.TaskARN)] = task
		}
	}

	return tasks, nil
}

// setExecutionState provides a way for tests to set initial state of scheduler. Not to be used by regular scheduler flow
func (scheduler *Scheduler) setExecutionState(state map[string]environmentExecutionState) {
	scheduler.executionState = state
}
