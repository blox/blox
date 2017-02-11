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

package wrappers

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
	"github.com/pkg/errors"
)

var (
	taskDefinitionSleep300 = "bloxCSSCanarySleep300"
)

// ECSWrapper defines methods to access wrapper methods to call ECS APIs
type ECSWrapper interface {
	CreateCluster(clusterName *string) error
	DeleteCluster(clusterName *string) error
	ListContainerInstances(clusterName *string) ([]*string, error)
	DeregisterContainerInstances(clusterName *string, instanceARNs []*string) error
	StartTask(clusterName *string, instanceARN *string) (*string, error)
}

type ecsClientWrapper struct {
	client ecsiface.ECSAPI
}

// NewECSWrapper returns a new ECSWrapper for the canary
func NewECSWrapper(sess *session.Session) (ECSWrapper, error) {
	if sess == nil {
		return nil, errors.New("AWS session for has to be initialized to initialize the ECS client. ")
	}
	ecsClient := ecs.New(sess)
	return ecsClientWrapper{
		client: ecsClient,
	}, nil
}

func (wrapper ecsClientWrapper) CreateCluster(clusterName *string) error {
	in := ecs.CreateClusterInput{
		ClusterName: clusterName,
	}

	_, err := wrapper.client.CreateCluster(&in)
	if err != nil {
		return errors.Wrapf(err, "Error creating ECS cluster with name '%s'. ", *clusterName)
	}

	return nil
}

func (wrapper ecsClientWrapper) DeleteCluster(clusterName *string) error {
	in := ecs.DeleteClusterInput{
		Cluster: clusterName,
	}

	_, err := wrapper.client.DeleteCluster(&in)
	if err != nil {
		return errors.Wrapf(err, "Error deleting ECS cluster with name '%s'. ", *clusterName)
	}

	return nil
}

func (wrapper ecsClientWrapper) ListContainerInstances(clusterName *string) ([]*string, error) {
	in := ecs.ListContainerInstancesInput{
		Cluster: clusterName,
	}

	resp, err := wrapper.client.ListContainerInstances(&in)
	if err != nil {
		return nil, errors.Wrapf(err, "Error listing container instances for ECS cluster with name '%s'. ",
			*clusterName)
	}

	return resp.ContainerInstanceArns, nil
}

func (wrapper ecsClientWrapper) DeregisterContainerInstances(clusterName *string, instanceARNs []*string) error {
	for _, instanceARN := range instanceARNs {
		err := wrapper.deregisterContainerInstance(clusterName, instanceARN)
		if err != nil {
			return err
		}
	}
	return nil
}

func (wrapper ecsClientWrapper) deregisterContainerInstance(clusterName *string, instanceARN *string) error {
	forceDeregister := true
	in := ecs.DeregisterContainerInstanceInput{
		Cluster:           clusterName,
		ContainerInstance: instanceARN,
		Force:             &forceDeregister,
	}

	_, err := wrapper.client.DeregisterContainerInstance(&in)
	if err != nil {
		return errors.Wrapf(err, "Error deregistering container instance with ARN '%s' in cluster with name '%s'. ",
			*instanceARN, *clusterName)
	}

	return nil
}

// TODO: Avoid hardcoding sleep time in the function name and make it a configurable parameter instead
func (wrapper ecsClientWrapper) registerSleep300TaskDefinition() (string, error) {
	taskDefnARN, err := wrapper.describeTaskDefinition(taskDefinitionSleep300)
	if err == nil {
		return taskDefnARN, nil
	}

	name := "sleep300"
	image := "busybox"
	cpu := int64(100)
	memory := int64(10)
	sleepCmd := "sleep"
	sleepTime := "300"
	command := []*string{&sleepCmd, &sleepTime}

	containerDefn := ecs.ContainerDefinition{
		Name:    &name,
		Image:   &image,
		Cpu:     &cpu,
		Memory:  &memory,
		Command: command,
	}

	in := ecs.RegisterTaskDefinitionInput{
		Family:               &taskDefinitionSleep300,
		ContainerDefinitions: []*ecs.ContainerDefinition{&containerDefn},
	}

	resp, err := wrapper.client.RegisterTaskDefinition(&in)
	if err != nil {
		return "", errors.Wrapf(err, "Could not register task definition '%s'. ", taskDefinitionSleep300)
	}

	return *resp.TaskDefinition.TaskDefinitionArn, nil
}

func (wrapper ecsClientWrapper) StartTask(clusterName *string, instanceARN *string) (*string, error) {
	_, err := wrapper.registerSleep300TaskDefinition()
	if err != nil {
		return nil, err
	}
	in := ecs.StartTaskInput{
		Cluster:            clusterName,
		ContainerInstances: []*string{instanceARN},
		TaskDefinition:     &taskDefinitionSleep300,
	}
	resp, err := wrapper.client.StartTask(&in)
	if err != nil {
		return nil, errors.Wrapf(err, "Error starting task on cluster '%s' on container instance '%s' using task definition '%s'. ",
			*clusterName, *instanceARN, taskDefinitionSleep300)
	}
	if len(resp.Failures) != 0 {
		reason := *resp.Failures[0].Reason
		return nil, errors.Errorf(
			"Failure starting task on cluster '%s' on container instance '%s' using task definition '%s'. Reason: %s. ",
			*clusterName, *instanceARN, taskDefinitionSleep300, reason)
	}
	if len(resp.Tasks) != 1 {
		return nil, errors.Errorf("'%d' tasks were started on cluster with name '%s' but expected only 1. ",
			len(resp.Tasks), *clusterName)
	}
	return resp.Tasks[0].TaskArn, nil
}

func (wrapper ecsClientWrapper) describeTaskDefinition(taskDefn string) (string, error) {
	in := ecs.DescribeTaskDefinitionInput{
		TaskDefinition: &taskDefn,
	}
	resp, err := wrapper.client.DescribeTaskDefinition(&in)
	if err != nil {
		return "", errors.Wrapf(err, "Could not describe task definition '%s'. ", taskDefn)
	}
	return *resp.TaskDefinition.TaskDefinitionArn, nil
}
