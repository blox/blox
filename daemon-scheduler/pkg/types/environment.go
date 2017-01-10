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

// SortDeploymentsReverseChronologically returns deployments ordered reverse-chronologically: latest startTime first
func (e Environment) SortDeploymentsReverseChronologically() ([]Deployment, error) {
	deployments := make([]Deployment, 0, len(e.Deployments))
	for _, d := range e.Deployments {
		deployments = append(deployments, d)
	}

	sort.Sort(timeOrderedDeployments(deployments))
	return deployments, nil
}
