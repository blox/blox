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
	"sort"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
)

type EnvironmentHealth uint8

const (
	EnvironmentHealthy EnvironmentHealth = iota
	EnvironmentUnhealthy
)

type Environment struct {
	Token                 string
	Name                  string
	DesiredTaskDefinition string
	DesiredTaskCount      int
	Cluster               string
	Health                EnvironmentHealth

	//TODO: use an internal type instead of ECS Task
	// taskdef -> [taskArn -> Task]
	CurrentTasks map[string]map[string]*ecs.Task

	// deploymentID -> deployment
	Deployments map[string]Deployment
}

type timeOrderedDeployments []Deployment

func (p timeOrderedDeployments) Len() int {
	return len(p)
}

// Less orders by latest startTime first
func (p timeOrderedDeployments) Less(i, j int) bool {
	return p[i].StartTime.After(p[j].StartTime)
}

func (p timeOrderedDeployments) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func NewEnvironment(name string, taskDefinition string, cluster string) (*Environment, error) {
	if len(name) == 0 {
		return nil, errors.New("Name should not be empty")
	}

	if len(taskDefinition) == 0 {
		return nil, errors.New("TaskDefinition should not be empty")
	}

	if len(cluster) == 0 {
		return nil, errors.New("Cluster should not be empty")
	}

	return &Environment{
		Token: uuid.NewRandom().String(),
		Name:  name,
		DesiredTaskDefinition: taskDefinition,
		Cluster:               cluster,
		Health:                EnvironmentHealthy,
		Deployments:           make(map[string]Deployment),
		CurrentTasks:          make(map[string]map[string]*ecs.Task),
	}, nil
}

// GetInProgressDeployment returns the in-progress deployment for the environment.
// There should only be one in-progress deployment.
func (e Environment) GetInProgressDeployment() (*Deployment, error) {
	deployments, err := e.sortDeploymentsByStartTime()
	if err != nil {
		return nil, err
	}

	// there should only be one in-progress deployment and it
	// should be close to the top. There might be some pending
	// deployments started after it so it's not always the latest one.
	for _, d := range deployments {
		if d.Status == DeploymentInProgress {
			return &d, nil
		}

		// once we reach completed deployments, we know there are no in-progress
		// deployments in the list after that since they're sorted by time and there
		// can only be one in-progress deployment
		if d.Status == DeploymentCompleted {
			return nil, nil
		}
	}

	return nil, nil
}

// GetDeployments returns a list of deployments reverse-ordered by start time,
// i.e. lastest deployment first
func (e Environment) GetDeployments() ([]Deployment, error) {
	return e.sortDeploymentsByStartTime()
}

func (e Environment) sortDeploymentsByStartTime() ([]Deployment, error) {
	deployments := make([]Deployment, 0, len(e.Deployments))
	for _, d := range e.Deployments {
		deployments = append(deployments, d)
	}

	sort.Sort(timeOrderedDeployments(deployments))
	return deployments, nil
}
