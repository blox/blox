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

package wrappers

import (
	"errors"

	"github.com/aws/amazon-ecs-event-stream-handler/internal/client"
	"github.com/aws/amazon-ecs-event-stream-handler/internal/client/operations"
	"github.com/aws/amazon-ecs-event-stream-handler/internal/models"
)

const (
	getTaskNotFoundException     = "GetTaskNotFound"
	getInstanceNotFoundException = "GetInstanceNotFound"
)

type CSSWrapper struct {
	client *client.AmazonCSS
}

func NewCSSWrapper() CSSWrapper {
	return CSSWrapper{
		client: client.NewHTTPClient(nil),
	}
}

func (wrapper CSSWrapper) GetTask(clusterName string, taskARN string) (*models.Task, error) {
	in := operations.NewGetTaskParams()
	in.SetCluster(clusterName)
	in.SetArn(taskARN)
	resp, err := wrapper.client.Operations.GetTask(in)
	if err != nil {
		return nil, err
	}
	task := resp.Payload
	return task, nil
}

func (wrapper CSSWrapper) TryGetTask(taskARN string) (string, error) {
	in := operations.NewGetTaskParams()
	in.SetArn(taskARN)
	_, err := wrapper.client.Operations.GetTask(in)
	if err != nil {
		if _, ok := err.(*operations.GetTaskNotFound); ok {
			return getTaskNotFoundException, nil
		}
		return "", errors.New("Unknown exception when calling Get Task")
	}
	return "", errors.New("Expected an exception when calling Get Task, but none received")
}

func (wrapper CSSWrapper) ListTasks() ([]*models.Task, error) {
	in := operations.NewListTasksParams()
	resp, err := wrapper.client.Operations.ListTasks(in)
	if err != nil {
		return nil, err
	}
	tasks := resp.Payload
	return tasks.Items, nil
}

func (wrapper CSSWrapper) FilterTasksByStatus(status string) ([]*models.Task, error) {
	in := operations.NewFilterTasksParams()
	in.SetStatus(status)
	resp, err := wrapper.client.Operations.FilterTasks(in)
	if err != nil {
		return nil, err
	}
	tasks := resp.Payload
	return tasks.Items, nil
}

func (wrapper CSSWrapper) FilterTasksByCluster(cluster string) ([]*models.Task, error) {
	in := operations.NewFilterTasksParams()
	in.SetCluster(cluster)
	resp, err := wrapper.client.Operations.FilterTasks(in)
	if err != nil {
		return nil, err
	}
	tasks := resp.Payload
	return tasks.Items, nil
}

func (wrapper CSSWrapper) GetInstance(clusterName string, instanceARN string) (*models.ContainerInstance, error) {
	in := operations.NewGetInstanceParams()
	in.SetCluster(clusterName)
	in.SetArn(instanceARN)
	resp, err := wrapper.client.Operations.GetInstance(in)
	if err != nil {
		return nil, err
	}
	instance := resp.Payload
	return instance, nil
}

func (wrapper CSSWrapper) TryGetInstance(instanceARN string) (string, error) {
	in := operations.NewGetInstanceParams()
	in.SetArn(instanceARN)
	_, err := wrapper.client.Operations.GetInstance(in)
	if err != nil {
		if _, ok := err.(*operations.GetInstanceNotFound); ok {
			return getInstanceNotFoundException, nil
		}
		return "", errors.New("Unknown exception when calling Get Instance")
	}
	return "", errors.New("Expected an exception when calling Get Instance, but none received")
}

func (wrapper CSSWrapper) ListInstances() ([]*models.ContainerInstance, error) {
	in := operations.NewListInstancesParams()
	resp, err := wrapper.client.Operations.ListInstances(in)
	if err != nil {
		return nil, err
	}
	instances := resp.Payload
	return instances.Items, nil
}

func (wrapper CSSWrapper) FilterInstancesByClusterName(clusterName string) ([]*models.ContainerInstance, error) {
	in := operations.NewFilterInstancesParams()
	in.SetCluster(clusterName)
	resp, err := wrapper.client.Operations.FilterInstances(in)
	if err != nil {
		return nil, err
	}
	instances := resp.Payload
	return instances.Items, nil
}
