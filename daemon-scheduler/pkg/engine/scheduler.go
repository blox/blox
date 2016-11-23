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
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/blox/blox/daemon-scheduler/pkg/clients/css/models"
	"github.com/blox/blox/daemon-scheduler/pkg/deployment"
	"github.com/blox/blox/daemon-scheduler/pkg/facade"
	"github.com/blox/blox/daemon-scheduler/pkg/types"
	log "github.com/cihub/seelog"
	"github.com/pkg/errors"
)

const (
	schedulerTickerDuration = 10 * time.Second
	inactiveInstanceStatus  = "INACTIVE"
	runningTaskStatus       = "RUNNING"
	maxClusterStateLag      = 1 * time.Minute
)

// Scheduler loops through all the environments and makes sure they reach their eventual state.
type Scheduler struct {
	ctx            context.Context
	environmentSvc deployment.Environment
	deploymentSvc  deployment.Deployment
	css            facade.ClusterState
	events         chan<- Event
	inProgress     bool
	inProgressLock sync.RWMutex
}

type instanceLookupResult struct {
	newInstances      []*string
	deployedInstances map[string][]*deployedTask
}

type deployedTask struct {
	task       string
	deployment *types.Deployment
}

func StartScheduler(ctx context.Context, events chan<- Event, environmentSvc deployment.Environment,
	deploymentSvc deployment.Deployment, css facade.ClusterState) {
	scheduler := Scheduler{
		ctx:            ctx,
		environmentSvc: environmentSvc,
		deploymentSvc:  deploymentSvc,
		css:            css,
		events:         events,
		inProgress:     false,
	}
	scheduler.loop()
	log.Infof("Started scheduler")
}

func (scheduler Scheduler) loop() {
	ticker := time.NewTicker(schedulerTickerDuration)
	go func() {
		for {
			select {
			case <-ticker.C:
				if scheduler.isInProgress() {
					log.Info("Scheduler loop in progress, skipping")
					continue
				}
				go func() {
					err := scheduler.runOnce()
					if err != nil {
						log.Errorf("Error running scheduler", err)
					}
				}()
			case <-scheduler.ctx.Done():
				log.Infof("Shutting down scheduler")
				ticker.Stop()
				return
			}
		}
	}()
}

// runOnce runs a single iteration of Scheduler.
func (scheduler Scheduler) runOnce() error {
	scheduler.setInProgress(true)
	defer scheduler.setInProgress(false)

	environments, err := scheduler.environmentSvc.ListEnvironments(scheduler.ctx)
	if err != nil {
		return errors.Wrapf(err, "Error getting environment")
	}

	for _, environment := range environments {
		// processing for environments is independent, so we can do them concurrently
		go func(environment types.Environment) {
			err := scheduler.runForEnvironment(&environment)
			if err != nil {
				// TODO: we may want to report this for better ux
				log.Errorf("Error running this iteration of Scheduler for environment %s: %v", environment.Name, err)
				return
			}
			log.Infof("Done running this iteration of scheduler for environment %s", environment.Name)
		}(environment)
	}

	return nil
}

func (scheduler Scheduler) setInProgress(val bool) {
	scheduler.inProgressLock.Lock()
	defer scheduler.inProgressLock.Unlock()
	scheduler.inProgress = val
}

func (scheduler Scheduler) isInProgress() bool {
	scheduler.inProgressLock.RLock()
	defer scheduler.inProgressLock.RUnlock()
	return scheduler.inProgress
}

func (scheduler Scheduler) runForEnvironment(environment *types.Environment) error {
	currentDeployment, err := scheduler.getCurrentDeployment(environment)
	if err != nil {
		return err
	}

	if currentDeployment == nil {
		log.Infof("No deployment available for environment %s", environment.Name)
		return nil
	}

	lookupResult, err := scheduler.lookupInstances(environment)
	if err != nil {
		return errors.Wrapf(err, "Error finding instances to deploy for environment %s", environment.Name)
	}

	scheduler.deployToNewInstances(environment, lookupResult)

	err = scheduler.updateDeployedInstances(environment, currentDeployment, lookupResult)
	if err != nil {
		return errors.Wrapf(err, "Error updating deployed instances for environment %s", environment.Name)
	}

	return nil
}

func (scheduler Scheduler) getCurrentDeployment(environment *types.Environment) (*types.Deployment, error) {
	deployment, err := scheduler.environmentSvc.GetCurrentDeployment(scheduler.ctx, environment.Name)
	if err != nil {
		return nil, errors.Wrapf(err, "Error getting current deployment for cluster %s of environment %s",
			environment.Cluster, environment.Name)
	}

	return deployment, nil
}

// updateDeployedInstances performs deployment on instances which already have some version of environment deployed
func (scheduler Scheduler) updateDeployedInstances(environment *types.Environment, currentDeployment *types.Deployment, result *instanceLookupResult) error {
	// go through already deployed instances and select the tasks which need to be replaced
	for instanceARN, deployedTasks := range result.deployedInstances {
		instanceRunningCurrent := false
		tasks := make([]string, 0)
		for _, dt := range deployedTasks {
			if dt.deployment.ID == currentDeployment.ID {
				//NOTE: in a pathological scenario there can be > 1 tasks corresponding
				//to current deployment running, we will stop all but one of those
				if !instanceRunningCurrent {
					instanceRunningCurrent = true
					continue
				}
			}
			tasks = append(tasks, dt.task)
		}

		// order is to stop existing task(s) and start new one
		if len(tasks) > 0 {
			log.Debugf("Sending StopTasksEvent with %d tasks", len(tasks))
			scheduler.events <- StopTasksEvent{
				Cluster:     environment.Cluster,
				Tasks:       tasks,
				Environment: *environment,
			}
		}

		if !instanceRunningCurrent {
			log.Debugf("Sending StartDeploymentEvent for deployment %s to instance %s", currentDeployment.ID, instanceARN)
			scheduler.events <- StartDeploymentEvent{
				Environment: *environment,
				Instances:   []*string{aws.String(instanceARN)},
			}
		}
	}

	return nil
}

// deployToNewInstances performs deployment on instances which never got any deployment for the given environment
func (scheduler Scheduler) deployToNewInstances(environment *types.Environment, result *instanceLookupResult) {
	if len(result.newInstances) > 0 {
		event := StartDeploymentEvent{
			Environment: *environment,
			Instances:   result.newInstances,
		}
		scheduler.events <- event
	}
}

// lookupInstances returns instanceLookupResult struct containing the state of all instances in the cluster corresponding to environment
func (scheduler Scheduler) lookupInstances(environment *types.Environment) (*instanceLookupResult, error) {

	// get all instances in Cluster
	instances, err := scheduler.css.ListInstances(environment.Cluster)
	if err != nil {
		return nil, errors.Wrapf(err, "Error getting instances for cluster %s of environment %s",
			environment.Cluster, environment.Name)
	}

	instanceARNToInstance := make(map[string]*models.ContainerInstance)
	for _, instance := range instances {
		instanceARNToInstance[aws.StringValue(instance.ContainerInstanceARN)] = instance
	}

	result := &instanceLookupResult{
		newInstances:      make([]*string, 0),
		deployedInstances: make(map[string][]*deployedTask),
	}

	result, err = scheduler.loadInstancesAlreadyDeployed(environment, instanceARNToInstance, result)
	if err != nil {
		return nil, errors.Wrapf(err, "Error finding instances where environment %s is already deployed", environment.Name)
	}

	// collect all the instances which do not have this environment installed
	for _, i := range instances {
		instanceARN := aws.StringValue(i.ContainerInstanceARN)
		if aws.StringValue(i.Status) == inactiveInstanceStatus {
			delete(result.deployedInstances, instanceARN)
			continue
		}
		_, ok := result.deployedInstances[instanceARN]
		if !ok {
			result.newInstances = append(result.newInstances, i.ContainerInstanceARN)
		}
	}
	return result, nil
}

// loadInstancesAlreadyDeployed populates instanceLookupResult struct with the state of instances derived from
// state of environments, deployments and cluster
func (scheduler Scheduler) loadInstancesAlreadyDeployed(environment *types.Environment,
	instanceARNToInstance map[string]*models.ContainerInstance, result *instanceLookupResult) (*instanceLookupResult, error) {
	tasks, err := scheduler.getRunningTasks(environment.Cluster)
	if err != nil {
		return result, err
	}

	deployments, err := scheduler.deploymentSvc.ListDeployments(scheduler.ctx, environment.Name)
	if err != nil {
		return result, errors.Wrapf(err, "Error calling ListDeployments with environment %s", environment.Name)
	}

	// preparing a map for easy lookup
	deploymentsMap := make(map[string]*types.Deployment)
	for _, d := range deployments {
		deploymentsMap[d.ID] = &d
	}

	// for each task find the deployment it corresponds to and tag the instance of the task as deployed
	for _, task := range tasks {
		deployment, ok := deploymentsMap[task.StartedBy]
		if ok {
			instanceARN := aws.StringValue(task.ContainerInstanceARN)
			instance, ok := instanceARNToInstance[instanceARN]
			if !ok || aws.StringValue(instance.Status) == inactiveInstanceStatus {
				continue
			}
			deployedTasks, ok := result.deployedInstances[instanceARN]
			if !ok {
				deployedTasks = make([]*deployedTask, 0)
			}

			taskARN := aws.StringValue(task.TaskARN)
			_, ok = tasks[taskARN]
			if ok {
				deployedTasks = append(deployedTasks, &deployedTask{task: taskARN, deployment: deployment})
			}

			result.deployedInstances[instanceARN] = deployedTasks
		}
	}

	return result, nil
}

// getRunningTasks returns a map of taskARN -> task where task is -probably- running
func (scheduler Scheduler) getRunningTasks(cluster string) (map[string]*models.Task, error) {
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
