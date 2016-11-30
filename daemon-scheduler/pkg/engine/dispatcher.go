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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/blox/blox/daemon-scheduler/pkg/clients/css/models"
	"github.com/blox/blox/daemon-scheduler/pkg/deployment"
	"github.com/blox/blox/daemon-scheduler/pkg/facade"
	log "github.com/cihub/seelog"
	"github.com/pkg/errors"
)

const (
	capacity = 100
)

// StartDispatcher starts dispatcher. dispatcher listens to events on channel and forwards them to workers
func StartDispatcher(ctx context.Context, environmentSvc deployment.Environment,
	deploymentSvc deployment.Deployment, ecs facade.ECS, css facade.ClusterState,
	deploymentWorker deployment.DeploymentWorker) chan<- Event {

	events := make(chan Event, capacity)

	go func() {
		for {
			select {
			case event := <-events:
				go func(event Event) {
					worker := worker{
						environmentSvc:   environmentSvc,
						deploymentSvc:    deploymentSvc,
						deploymentWorker: deploymentWorker,
						ecs:              ecs,
						css:              css,
					}
					err := worker.handleEvent(ctx, event)
					if err != nil {
						log.Errorf("Error handling event: %v", err)
					}
				}(event)
			case <-ctx.Done():
				log.Infof("Shutting down dispatcher")
				return
			}
		}
	}()

	log.Infof("Started dispatcher")

	return events
}

// Worker is actor which handles an event appropriately
type worker struct {
	environmentSvc   deployment.Environment
	deploymentSvc    deployment.Deployment
	deploymentWorker deployment.DeploymentWorker
	ecs              facade.ECS
	css              facade.ClusterState
}

func (w *worker) handleEvent(ctx context.Context, event Event) error {
	switch event.GetType() {
	case StartDeploymentEventType:
		return w.handleStartDeploymentEvent(ctx, event)
	case StopTasksEventType:
		return w.handleStopTasksEvent(ctx, event)
	case UpdateInProgressDeploymentEventType:
		return w.handleUpdateInProgressDeploymentEvent(ctx, event)
	default:
		return w.handleDefaultEvent(ctx, event)
	}
}

func (w *worker) handleDefaultEvent(ctx context.Context, event Event) error {
	log.Debugf("Received event : %s", event.GetType())
	return nil
}

func (w *worker) handleUpdateInProgressDeploymentEvent(ctx context.Context, event Event) error {
	deploymentEvent, ok := event.(UpdateInProgressDeploymentEvent)
	if !ok {
		return errors.Errorf("Expected event with event-type %v to be of struct-type UpdateInProgressDeploymentEvent",
			event.GetType())
	}

	_, err := w.deploymentWorker.UpdateInProgressDeployment(ctx, deploymentEvent.Environment.Name)
	if err != nil {
		return err
	}

	return nil
}

func (w *worker) handleStartDeploymentEvent(ctx context.Context, event Event) error {
	deploymentEvent, ok := event.(StartDeploymentEvent)
	if !ok {
		return errors.Errorf("Expected event with event-type %s to be of struct-type StartDeploymentEvent", event.GetType())
	}

	deployment, err := w.deploymentSvc.CreateSubDeployment(ctx, deploymentEvent.Environment.Name, deploymentEvent.Instances)
	if err != nil {
		return errors.Wrapf(err, "Error starting deployment using environment %s on %d instances",
			deploymentEvent.Environment.Name, len(deploymentEvent.Instances))
	}

	log.Infof("Succesfully created a deployment with %s on %d instances",
		deployment.ID, len(deploymentEvent.Instances))
	return nil
}

func (w *worker) handleStopTasksEvent(ctx context.Context, event Event) error {
	stopTasksEvent, ok := event.(StopTasksEvent)
	if !ok {
		return errors.Errorf("Expected event with event-type %s to be of struct-type StopTasksEvent", event.GetType())
	}

	tasksInCluster, err := w.css.ListTasks(stopTasksEvent.Cluster)
	if err != nil {
		return errors.Wrapf(err, "Error getting tasks in cluster %s", stopTasksEvent.Cluster)
	}
	taskMap := make(map[string]*models.Task)
	for _, task := range tasksInCluster {
		taskMap[aws.StringValue(task.TaskARN)] = task
	}

	stoppedTasks := []string{}
	for _, task := range stopTasksEvent.Tasks {
		knownTask, ok := taskMap[task]
		if !ok {
			continue
		}
		if aws.StringValue(knownTask.DesiredStatus) == "STOPPED" {
			stoppedTasks = append(stoppedTasks, task)
			continue
		}
		err := w.ecs.StopTask(stopTasksEvent.Cluster, task)
		if err != nil {
			log.Errorf("Error stopping task %s : %v", task, err)
			continue
		}
		stoppedTasks = append(stoppedTasks, task)
	}

	// TODO: Clear the tasks from environment

	log.Infof("Successfully stopped %d tasks out of %d tasks under environment %s",
		len(stoppedTasks), len(stopTasksEvent.Tasks), stopTasksEvent.Environment.Name)

	return nil
}
