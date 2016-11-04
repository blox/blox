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

package wrappers

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/pkg/errors"
)

var (
	taskDefinitionSleep300 = "esh_test_sleep300"
)

type ECSWrapper struct {
	client *ecs.ECS
}

func NewECSWrapper() ECSWrapper {
	awsSession := newAWSSession()
	return ECSWrapper{
		client: ecs.New(awsSession),
	}
}

func newAWSSession() *session.Session {
	var sess *session.Session
	var err error
	if endpoint, err := getECSEndpoint(); err != nil {
		sess, err = session.NewSession()
	} else {
		sess, err = session.NewSessionWithOptions(session.Options{
			Config: aws.Config{
				Endpoint: aws.String(endpoint),
			},
		})
	}
	if err != nil {
		panic(err)
	}
	return sess
}

func (ecsWrapper ECSWrapper) RegisterSleep360TaskDefinition() (string, error) {
	taskDefnARN, err := ecsWrapper.DescribeTaskDefinition(taskDefinitionSleep300)
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

	resp, err := ecsWrapper.client.RegisterTaskDefinition(&in)
	if err != nil {
		return "", errors.New("Could not register sleep300 task definition")
	}

	return *resp.TaskDefinition.TaskDefinitionArn, nil
}

func (ecsWrapper ECSWrapper) DescribeTaskDefinition(taskDefn string) (string, error) {
	in := ecs.DescribeTaskDefinitionInput{
		TaskDefinition: &taskDefn,
	}
	resp, err := ecsWrapper.client.DescribeTaskDefinition(&in)
	if err != nil {
		return "", err
	}
	return *resp.TaskDefinition.TaskDefinitionArn, nil
}

func (ecsWrapper ECSWrapper) DeregisterTaskDefinition(taskDefnARN string) error {
	in := ecs.DeregisterTaskDefinitionInput{
		TaskDefinition: &taskDefnARN,
	}
	_, err := ecsWrapper.client.DeregisterTaskDefinition(&in)
	if err != nil {
		return errors.Errorf("Failed to deregister task definition with ARN '%s'", taskDefnARN)
	}
	return nil
}

func (ecsWrapper ECSWrapper) StartTask(clusterName string, taskDefn string) (ecs.Task, error) {
	containerInstances, err := ecsWrapper.ListContainerInstances(clusterName)
	if err != nil {
		return ecs.Task{}, err
	}
	if len(containerInstances) < 1 {
		return ecs.Task{}, errors.Errorf("No container instance registered to cluster '%s'", clusterName)
	}
	in := ecs.StartTaskInput{
		Cluster:            &clusterName,
		ContainerInstances: containerInstances[0:1],
		TaskDefinition:     &taskDefn,
	}
	resp, err := ecsWrapper.client.StartTask(&in)
	if err != nil {
		return ecs.Task{}, err
	}
	if len(resp.Failures) != 0 {
		reason := *resp.Failures[0].Reason
		return ecs.Task{}, errors.Errorf(
			"Failure starting task on cluster '%s' with '%d' container instances using task definition '%s'. Reason: %s",
			clusterName, len(in.ContainerInstances), taskDefn, reason)
	}
	if len(resp.Tasks) != 1 {
		return ecs.Task{}, errors.New("Invalid number of tasks started")
	}
	return *resp.Tasks[0], nil
}

func (ecsWrapper ECSWrapper) StopTask(clusterName string, taskARN string) error {
	in := ecs.StopTaskInput{
		Cluster: &clusterName,
		Task:    &taskARN,
	}
	_, err := ecsWrapper.client.StopTask(&in)
	if err != nil {
		return errors.Wrapf(err, "Failed to stop task with ARN '%s' on cluster: '%s'", taskARN, clusterName)
	}
	return nil
}

func (ecsWrapper ECSWrapper) ListTasks(clusterName string) ([]*string, error) {
	in := ecs.ListTasksInput{
		Cluster: &clusterName,
	}
	resp, err := ecsWrapper.client.ListTasks(&in)
	if err != nil {
		return nil, errors.New("Failed to list ECS tasks")
	}
	return resp.TaskArns, nil
}

func (ecsWrapper ECSWrapper) ListContainerInstances(clusterName string) ([]*string, error) {
	in := ecs.ListContainerInstancesInput{
		Cluster: &clusterName,
	}
	resp, err := ecsWrapper.client.ListContainerInstances(&in)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to list ECS container instances with cluster name: '%s'", clusterName)
	}
	return resp.ContainerInstanceArns, nil
}

func (ecsWrapper ECSWrapper) DescribeContainerInstance(clusterName string, instanceARN string) (ecs.ContainerInstance, error) {
	in := ecs.DescribeContainerInstancesInput{
		Cluster:            &clusterName,
		ContainerInstances: []*string{&instanceARN},
	}
	resp, err := ecsWrapper.client.DescribeContainerInstances(&in)
	if err != nil {
		return ecs.ContainerInstance{}, errors.Errorf("Failed to describe container instance with ARN '%s'", instanceARN)
	}
	if len(resp.Failures) != 0 {
		reason := *resp.Failures[0].Reason
		return ecs.ContainerInstance{}, errors.Errorf("Failed to describe container instance with ARN '%s'. Reason: %s", instanceARN, reason)
	}
	if len(resp.ContainerInstances) != 1 {
		return ecs.ContainerInstance{}, errors.New("Invalid number of instances in describe container instance response")
	}
	return *resp.ContainerInstances[0], nil
}
