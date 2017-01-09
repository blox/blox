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
	// ID of the deployment created by the latest create-deployment call.
	// This need not necessarily correspond to the deployment picked up by the
	// background jobs launching tasks
	PendingDeploymentID string
	// The ID of the deployment that is being used by the background workers to
	// to launch tasks
	InProgressDeploymentID string

	// deploymentID -> deployment
	Deployments map[string]Deployment
}

type timeOrderedDeployments []Deployment

func (p timeOrderedDeployments) Len() int {
	return len(p)
}

// Less orders deployments reverse-chronologically: latest startTime first
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
	}, nil
}

// GetInProgressDeployment returns the in-progress deployment for the environment.
// There should be no more than one in progress deployments in an environment
func (e Environment) GetInProgressDeployment() (*Deployment, error) {
	if e.InProgressDeploymentID == "" {
		// TODO : We may want to change where the in progress deployment id is set
		// using the pending deployment id later
		if e.PendingDeploymentID == "" {
			return nil, nil
		}
		e.InProgressDeploymentID = e.PendingDeploymentID
	}
	d, ok := e.Deployments[e.InProgressDeploymentID]
	if !ok {
		return nil, errors.Errorf("Deployment with ID '%s' does not exist in the deployments for environment with name '%s'",
			e.InProgressDeploymentID, e.Name)
	}
	// If the status of the deployment is pending (this happens if we just set the
	// in progress deployment id from the pending deployment id), update the status
	// of the deployment to in progress
	if d.Status == DeploymentPending {
		d.Status = DeploymentInProgress
	} else if d.Status != DeploymentInProgress {
		return nil, nil
	}
	return &d, nil
}

// GetDeployments returns a list of deployments reverse-ordered by start time,
// i.e. lastest deployment first
func (e Environment) GetDeployments() ([]Deployment, error) {
	return e.SortDeploymentsReverseChronologically()
}

func (e Environment) GetCurrentDeployment() (*Deployment, error) {
	deployment, err := e.GetInProgressDeployment()
	if err != nil {
		return nil, err
	}

	if deployment != nil {
		return deployment, nil
	}

	// if there is no in-progress deployment then we take the latest completed deployment
	deployments, err := e.SortDeploymentsReverseChronologically()
	if err != nil {
		return nil, err
	}

	for _, d := range deployments {
		if d.Status == DeploymentCompleted {
			return &d, nil
		}
	}

	return nil, errors.Errorf("No deployment available for environment %s", e.Name)
}

// SortDeploymentsReverseChronologically returns deployments ordered reverse-chronologically: latest startTime first
func (e Environment) SortDeploymentsReverseChronologically() ([]Deployment, error) {
	deployments := make([]Deployment, 0, len(e.Deployments))
	for _, d := range e.Deployments {
		deployments = append(deployments, d)
	}

	sort.Sort(timeOrderedDeployments(deployments))
	return deployments, nil
}
