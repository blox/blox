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

package facade

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
	"github.com/pkg/errors"
)

const (
	startTaskPartitionSize    = 10
	describeTaskPartitionSize = 100
)

type ECS interface {
	StartTask(
		clusterArn string,
		containerInstances []*string,
		startedBy string,
		taskDefinition string) (*ecs.StartTaskOutput, error)

	ListClusters() ([]*string, error)
	DescribeCluster(cluster *string) (*ecs.Cluster, error)
	DescribeTaskDefinition(taskDefinition *string) (*ecs.TaskDefinition, error)
	ListTasks(cluster string, startedBy string) ([]*string, error)
	ListTasksByInstance(cluster string, instanceARN string) ([]*string, error)
	DescribeTasks(cluster string, tasks []*string) (*ecs.DescribeTasksOutput, error)
	StopTask(clusterArn string, taskArn string) error
}

type ecsClient struct {
	ecs ecsiface.ECSAPI
}

func NewECS(ecs ecsiface.ECSAPI) ECS {
	return ecsClient{
		ecs: ecs,
	}
}

func (c ecsClient) StartTask(
	clusterArn string,
	containerInstances []*string,
	startedBy string,
	taskDefinition string) (*ecs.StartTaskOutput, error) {
	output := &ecs.StartTaskOutput{
		Failures: []*ecs.Failure{},
		Tasks:    []*ecs.Task{},
	}

	//NOTE: StartTask takes 10 instances at a time
	for i := 0; i < len(containerInstances); i += startTaskPartitionSize {
		high := i + startTaskPartitionSize
		if high > len(containerInstances) {
			high = len(containerInstances)
		}
		partition := containerInstances[i:high]
		input := &ecs.StartTaskInput{
			Cluster:            aws.String(clusterArn),
			ContainerInstances: partition,
			StartedBy:          aws.String(startedBy),
			TaskDefinition:     aws.String(taskDefinition),
		}

		resp, err := c.ecs.StartTask(input)
		if err != nil {
			return nil, errors.Wrapf(err, "Could not start taskDefinition %v on cluster %v",
				taskDefinition, clusterArn)
		}

		output.Failures = append(output.Failures, resp.Failures...)
		output.Tasks = append(output.Tasks, resp.Tasks...)

	}
	return output, nil
}

func (c ecsClient) DescribeCluster(cluster *string) (*ecs.Cluster, error) {
	input := &ecs.DescribeClustersInput{
		Clusters: []*string{cluster},
	}
	resp, err := c.ecs.DescribeClusters(input)
	if err != nil {
		return nil, errors.Wrapf(err, "Error calling DescribeClusters for cluster %s", *cluster)
	}

	if len(resp.Clusters) == 0 {
		return nil, errors.Wrapf(err, "Cluster with name %s is missing", *cluster)
	}
	return resp.Clusters[0], nil
}

func (c ecsClient) DescribeTaskDefinition(td *string) (*ecs.TaskDefinition, error) {
	input := &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: td,
	}
	resp, err := c.ecs.DescribeTaskDefinition(input)
	if err != nil {
		return nil, errors.Wrapf(err, "Error calling DescribeTaskDefinition for taskDefinition %s", *td)
	}

	return resp.TaskDefinition, nil
}

func (c ecsClient) DescribeTasks(cluster string, tasks []*string) (*ecs.DescribeTasksOutput, error) {
	output := &ecs.DescribeTasksOutput{
		Failures: []*ecs.Failure{},
		Tasks:    []*ecs.Task{},
	}

	//NOTE: DescribeTasks takes 100 tasks at a time
	for i := 0; i < len(tasks); i += describeTaskPartitionSize {
		high := i + describeTaskPartitionSize
		if high > len(tasks) {
			high = len(tasks)
		}
		partition := tasks[i:high]

		input := &ecs.DescribeTasksInput{
			Cluster: aws.String(cluster),
			Tasks:   partition,
		}

		resp, err := c.ecs.DescribeTasks(input)
		if err != nil {
			return nil, errors.Wrapf(err, "Error calling DescribeTasks for cluster %s and tasks %v", cluster, tasks)
		}

		output.Failures = append(output.Failures, resp.Failures...)
		output.Tasks = append(output.Tasks, resp.Tasks...)
	}

	return output, nil
}

func (c ecsClient) ListClusters() ([]*string, error) {
	clusters := []*string{}
	var nextToken *string

	for {
		input := &ecs.ListClustersInput{}
		resp, err := c.ecs.ListClusters(input)
		if err != nil {
			return nil, errors.Wrap(err, "Error list-clusters")
		}

		if resp.ClusterArns != nil {
			clusters = append(clusters, resp.ClusterArns...)
		}

		if aws.StringValue(nextToken) == "" {
			break
		}

		input.NextToken = nextToken
	}

	return clusters, nil
}

func (c ecsClient) ListTasks(cluster string, startedBy string) ([]*string, error) {
	input := &ecs.ListTasksInput{
		Cluster:   aws.String(cluster),
		StartedBy: aws.String(startedBy),
	}

	resp, err := c.ecs.ListTasks(input)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to list ECS tasks in cluster %v startedBy %v", cluster, startedBy)
	}

	return resp.TaskArns, nil
}

func (c ecsClient) ListTasksByInstance(cluster string, instanceARN string) ([]*string, error) {
	input := &ecs.ListTasksInput{
		Cluster:           aws.String(cluster),
		ContainerInstance: aws.String(instanceARN),
	}

	resp, err := c.ecs.ListTasks(input)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to list ECS tasks in cluster %s instance %s", cluster, instanceARN)
	}

	return resp.TaskArns, nil
}

func (c ecsClient) StopTask(clusterArn string, taskArn string) error {
	input := &ecs.StopTaskInput{
		Cluster: aws.String(clusterArn),
		Task:    aws.String(taskArn),
	}
	_, err := c.ecs.StopTask(input)
	if err != nil {
		return errors.Wrapf(err, "Error stopping task %s in cluster %s", taskArn, clusterArn)
	}
	return nil
}
