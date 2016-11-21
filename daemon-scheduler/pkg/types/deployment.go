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

package types

import (
	"time"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

type DeploymentStatus uint8

const (
	DeploymentPending DeploymentStatus = iota
	DeploymentInProgress
	DeploymentCompleted
)

type DeploymentHealth uint8

const (
	DeploymentHealthy DeploymentHealth = iota
	DeploymentUnhealthy
)

type Deployment struct {
	ID               string
	Status           DeploymentStatus
	Health           DeploymentHealth
	TaskDefinition   string
	DesiredTaskCount int

	//TODO: use an internal type instead of ECS Task
	CurrentTasks    []*ecs.Task
	FailedInstances []*ecs.Failure
	StartTime       time.Time
	EndTime         time.Time
}

func NewDeployment(taskDefinition string) (*Deployment, error) {
	if len(taskDefinition) == 0 {
		return nil, errors.New("Task definition cannot be empty")
	}

	return &Deployment{
		ID:             uuid.NewV4().String(),
		Status:         DeploymentPending,
		Health:         DeploymentHealthy,
		StartTime:      time.Now(),
		TaskDefinition: taskDefinition,
	}, nil
}

func (d Deployment) UpdateDeploymentInProgress(
	desiredTaskCount int,
	currentTasks []*ecs.Task,
	failedInstances []*ecs.Failure) (*Deployment, error) {

	if d.Status == DeploymentCompleted {
		return nil, errors.New("Deployment cannot move from completed to in-progress")
	}

	if len(failedInstances) == 0 {
		d.Health = DeploymentHealthy
	} else {
		d.Health = DeploymentUnhealthy
	}

	d.Status = DeploymentInProgress
	d.DesiredTaskCount = desiredTaskCount
	d.CurrentTasks = currentTasks
	d.FailedInstances = failedInstances

	return &d, nil
}

func (d Deployment) UpdateDeploymentCompleted(
	CurrentTasks []*ecs.Task) *Deployment {

	d.Status = DeploymentCompleted
	d.Health = DeploymentHealthy
	d.CurrentTasks = CurrentTasks
	d.FailedInstances = nil
	d.EndTime = time.Now()

	return &d
}
