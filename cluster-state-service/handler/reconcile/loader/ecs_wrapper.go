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

package loader

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
	"github.com/blox/blox/cluster-state-service/handler/types"
	"github.com/pkg/errors"
)

const (
	describeInstancesPageSize = 100
	describeTasksPageSize     = 100
)

// ECSWrapper defines methods to access wrapper methods to call ECS APIs
type ECSWrapper interface {
	ListAllClusters() ([]*string, error)
	ListAllTasks(clusterARN *string) ([]*string, error)
	DescribeTasks(clusterARN *string, taskARNs []*string) ([]types.Task, []string, error)
	ListAllContainerInstances(clusterARN *string) ([]*string, error)
	DescribeContainerInstances(clusterARN *string, instanceARNs []*string) ([]types.ContainerInstance, []string, error)
}

type clientWrapper struct {
	client ecsiface.ECSAPI
}

func NewECSWrapper(ecsClient ecsiface.ECSAPI) ECSWrapper {
	return clientWrapper{
		client: ecsClient,
	}
}

// ListAllClusters retrieves a list of all cluster ARNS by making one or more calls to ECS
func (wrapper clientWrapper) ListAllClusters() ([]*string, error) {
	var clusterARNs []*string
	var nextToken *string
	nextToken = nil
	for {
		c, n, err := wrapper.listClusters(nextToken)
		if err != nil {
			return nil, err
		}
		clusterARNs = append(clusterARNs, c...)
		if aws.StringValue(n) == "" {
			break
		}
		nextToken = n
	}

	return clusterARNs, nil
}

func (wrapper clientWrapper) listClusters(nextToken *string) ([]*string, *string, error) {
	in := ecs.ListClustersInput{
		NextToken: nextToken,
	}

	resp, err := wrapper.client.ListClusters(&in)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "Failed to list ECS clusters.")
	}

	return resp.ClusterArns, resp.NextToken, nil
}

// ListAllTasks retrieves a list of all task ARNS in the cluster identified by 'clusterARN' by making one or more calls to ECS
func (wrapper clientWrapper) ListAllTasks(clusterARN *string) ([]*string, error) {
	var taskARNs []*string
	var nextToken *string
	nextToken = nil
	for {
		t, n, err := wrapper.listTasks(clusterARN, nextToken)
		if err != nil {
			return nil, err
		}
		taskARNs = append(taskARNs, t...)
		if aws.StringValue(n) == "" {
			break
		}
		nextToken = n
	}
	return taskARNs, nil
}

func (wrapper clientWrapper) listTasks(clusterARN *string, nextToken *string) ([]*string, *string, error) {
	if aws.StringValue(clusterARN) == "" {
		return nil, nil, errors.New("Failed to list ECS tasks. Error: Cluster cannot be empty")
	}

	in := ecs.ListTasksInput{
		Cluster:   clusterARN,
		NextToken: nextToken,
	}

	resp, err := wrapper.client.ListTasks(&in)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "Failed to list ECS tasks. Error: %s")
	}

	return resp.TaskArns, resp.NextToken, nil
}

// DescribeTasks desribes all tasks identified by 'taskARNs' belonging to cluster identified by 'clusterARN'
func (wrapper clientWrapper) DescribeTasks(clusterARN *string, taskARNs []*string) ([]types.Task, []string, error) {
	if aws.StringValue(clusterARN) == "" {
		return nil, nil, errors.New("Failed to describe ECS tasks. Error: Cluster cannot be empty")
	}
	tasks := make([]types.Task, 0)
	failedTaskARNS := make([]string, 0)

	for i := 0; i < len(taskARNs); i += describeTasksPageSize {
		high := i + describeTasksPageSize
		if high > len(taskARNs) {
			high = len(taskARNs)
		}

		in := ecs.DescribeTasksInput{
			Cluster: clusterARN,
			Tasks:   taskARNs[i:high],
		}

		resp, err := wrapper.client.DescribeTasks(&in)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "Failed to describe ECS tasks.")
		}
		for i := range resp.Tasks {
			task := ToTask(*resp.Tasks[i])
			tasks = append(tasks, task)
		}
		for i := range resp.Failures {
			failedTaskARNS = append(failedTaskARNS, aws.StringValue(resp.Failures[i].Arn))
		}
	}
	return tasks, failedTaskARNS, nil
}

// ListAllContainerInstances retrieves a list of all container instance ARNS in the cluster identified by 'clusterARN' by making one or more calls to ECS
func (wrapper clientWrapper) ListAllContainerInstances(clusterARN *string) ([]*string, error) {
	var instanceARNs []*string
	var nextToken *string
	nextToken = nil
	for {
		c, n, err := wrapper.listContainerInstances(clusterARN, nextToken)
		if err != nil {
			return nil, err
		}
		instanceARNs = append(instanceARNs, c...)
		if aws.StringValue(n) == "" {
			break
		}
		nextToken = n
	}
	return instanceARNs, nil
}

func (wrapper clientWrapper) listContainerInstances(clusterARN *string, nextToken *string) ([]*string, *string, error) {
	if aws.StringValue(clusterARN) == "" {
		return nil, nil, errors.New("Failed to list ECS container instances. Error: Cluster cannot be empty")
	}

	in := ecs.ListContainerInstancesInput{
		Cluster:   clusterARN,
		NextToken: nextToken,
	}

	resp, err := wrapper.client.ListContainerInstances(&in)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "Failed to list ECS container instances.")
	}

	return resp.ContainerInstanceArns, resp.NextToken, nil
}

// DescribeContainerInstances desribes all container instances identified by 'instanceARNs' belonging to cluster identified by 'clusterARN'
func (wrapper clientWrapper) DescribeContainerInstances(clusterARN *string, instanceARNs []*string) ([]types.ContainerInstance, []string, error) {
	if aws.StringValue(clusterARN) == "" {
		return nil, nil, errors.New("Failed to describe ECS container instances. Error: Cluster cannot be empty")
	}
	instances := make([]types.ContainerInstance, 0)
	failedInstanceARNS := make([]string, 0)

	for i := 0; i < len(instanceARNs); i += describeInstancesPageSize {
		high := i + describeInstancesPageSize
		if high > len(instanceARNs) {
			high = len(instanceARNs)
		}
		in := ecs.DescribeContainerInstancesInput{
			Cluster:            clusterARN,
			ContainerInstances: instanceARNs[i:high],
		}

		resp, err := wrapper.client.DescribeContainerInstances(&in)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "Failed to describe ECS container instances.")
		}
		for i := range resp.ContainerInstances {
			ins := ToContainerInstance(*resp.ContainerInstances[i], *clusterARN)
			instances = append(instances, ins)
		}
		for i := range resp.Failures {
			failedInstanceARNS = append(failedInstanceARNS, aws.StringValue(resp.Failures[i].Arn))
		}
	}
	return instances, failedInstanceARNS, nil
}
