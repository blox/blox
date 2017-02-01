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
	"errors"

	"io"
	"time"

	"github.com/blox/blox/cluster-state-service/swagger/v1/generated/client"
	"github.com/blox/blox/cluster-state-service/swagger/v1/generated/client/operations"
	"github.com/blox/blox/cluster-state-service/swagger/v1/generated/models"
)

const (
	getTaskNotFoundException         = "GetTaskNotFound"
	getInstanceNotFoundException     = "GetInstanceNotFound"
	listTasksBadRequestException     = "ListTasksBadRequest"
	listInstancesBadRequestException = "ListInstancesBadRequest"
)

type CSSWrapper struct {
	client *client.BloxCSS
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

func (wrapper CSSWrapper) TryGetTask(taskARN string) (string, string, error) {
	in := operations.NewGetTaskParams()
	in.SetArn(taskARN)
	_, err := wrapper.client.Operations.GetTask(in)
	if err != nil {
		if _, ok := err.(*operations.GetTaskNotFound); ok {
			return err.(*operations.GetTaskNotFound).Payload, getTaskNotFoundException, nil
		}
		return "", "", errors.New("Unknown exception when calling GetTask")
	}
	return "", "", errors.New("Expected an exception when calling GetTask, but none received")
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

func (wrapper CSSWrapper) TryListTasksWithInvalidStatus(status string) (string, string, error) {
	in := operations.NewListTasksParams()
	in.SetStatus(&status)
	_, err := wrapper.client.Operations.ListTasks(in)
	if err != nil {
		if _, ok := err.(*operations.ListTasksBadRequest); ok {
			return err.(*operations.ListTasksBadRequest).Payload, listTasksBadRequestException, nil
		}
		return "", "", errors.New("Unknown exception when calling ListTasks with invalid status")
	}
	return "", "", errors.New("Expected an exception when calling ListTasks with invalid status, but none received")
}

func (wrapper CSSWrapper) TryListTasksWithInvalidCluster(cluster string) (string, string, error) {
	in := operations.NewListTasksParams()
	in.SetCluster(&cluster)
	_, err := wrapper.client.Operations.ListTasks(in)
	if err != nil {
		if _, ok := err.(*operations.ListTasksBadRequest); ok {
			return err.(*operations.ListTasksBadRequest).Payload, listTasksBadRequestException, nil
		}
		return "", "", errors.New("Unknown exception when calling ListTasks with invalid cluster")
	}
	return "", "", errors.New("Expected an exception when calling ListTasks with invalid cluster, but none received")
}

func (wrapper CSSWrapper) ListTasksWithAllFilters(status string, cluster string, startedBy string) ([]*models.Task, error) {
	in := operations.NewListTasksParams()
	in.SetStatus(&status)
	in.SetCluster(&cluster)
	in.SetStartedBy(&startedBy)
	resp, err := wrapper.client.Operations.ListTasks(in)
	if err != nil {
		return nil, err
	}
	tasks := resp.Payload
	return tasks.Items, nil
}

func (wrapper CSSWrapper) FilterTasksByStatus(status string) ([]*models.Task, error) {
	in := operations.NewListTasksParams()
	in.SetStatus(&status)
	resp, err := wrapper.client.Operations.ListTasks(in)
	if err != nil {
		return nil, err
	}
	tasks := resp.Payload
	return tasks.Items, nil
}

func (wrapper CSSWrapper) FilterTasksByCluster(cluster string) ([]*models.Task, error) {
	in := operations.NewListTasksParams()
	in.SetCluster(&cluster)
	resp, err := wrapper.client.Operations.ListTasks(in)
	if err != nil {
		return nil, err
	}
	tasks := resp.Payload
	return tasks.Items, nil
}

func (wrapper CSSWrapper) FilterTasksByStartedBy(startedBy string) ([]*models.Task, error) {
	in := operations.NewListTasksParams()
	in.SetStartedBy(&startedBy)
	resp, err := wrapper.client.Operations.ListTasks(in)
	if err != nil {
		return nil, err
	}
	tasks := resp.Payload
	return tasks.Items, nil
}

func (wrapper CSSWrapper) FilterTasksByStatusAndCluster(status string, cluster string) ([]*models.Task, error) {
	in := operations.NewListTasksParams()
	in.SetStatus(&status)
	in.SetCluster(&cluster)
	resp, err := wrapper.client.Operations.ListTasks(in)
	if err != nil {
		return nil, err
	}
	tasks := resp.Payload
	return tasks.Items, nil
}

func (wrapper CSSWrapper) StreamTasks() (*io.PipeReader, error) {
	r, w := io.Pipe()
	in := operations.NewStreamTasksParams()
	in.SetTimeout(10 * time.Second)
	go func() {
		defer w.Close()
		wrapper.client.Operations.StreamTasks(in, w)
	}()
	return r, nil
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

func (wrapper CSSWrapper) TryGetInstance(instanceARN string) (string, string, error) {
	in := operations.NewGetInstanceParams()
	in.SetArn(instanceARN)
	_, err := wrapper.client.Operations.GetInstance(in)
	if err != nil {
		if _, ok := err.(*operations.GetInstanceNotFound); ok {
			return err.(*operations.GetInstanceNotFound).Payload, getInstanceNotFoundException, nil
		}
		return "", "", errors.New("Unknown exception when calling Get Instance")
	}
	return "", "", errors.New("Expected an exception when calling Get Instance, but none received")
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

func (wrapper CSSWrapper) TryListInstancesWithInvalidStatus(status string) (string, string, error) {
	in := operations.NewListInstancesParams()
	in.SetStatus(&status)
	_, err := wrapper.client.Operations.ListInstances(in)
	if err != nil {
		if _, ok := err.(*operations.ListInstancesBadRequest); ok {
			return err.(*operations.ListInstancesBadRequest).Payload, listInstancesBadRequestException, nil
		}
		return "", "", errors.New("Unknown exception when calling ListInstances with invalid status")
	}
	return "", "", errors.New("Expected an exception when calling ListInstances with invalid status, but none received")
}

func (wrapper CSSWrapper) TryListInstancesWithInvalidCluster(cluster string) (string, string, error) {
	in := operations.NewListInstancesParams()
	in.SetCluster(&cluster)
	_, err := wrapper.client.Operations.ListInstances(in)
	if err != nil {
		if _, ok := err.(*operations.ListInstancesBadRequest); ok {
			return err.(*operations.ListInstancesBadRequest).Payload, listInstancesBadRequestException, nil
		}
		return "", "", errors.New("Unknown exception when calling ListInstances with invalid cluster")
	}
	return "", "", errors.New("Expected an exception when calling ListInstances with invalid cluster, but none received")
}

func (wrapper CSSWrapper) FilterInstancesByClusterName(clusterName string) ([]*models.ContainerInstance, error) {
	in := operations.NewListInstancesParams()
	in.SetCluster(&clusterName)
	resp, err := wrapper.client.Operations.ListInstances(in)
	if err != nil {
		return nil, err
	}
	instances := resp.Payload
	return instances.Items, nil
}

func (wrapper CSSWrapper) StreamInstances() (*io.PipeReader, error) {
	r, w := io.Pipe()
	in := operations.NewStreamInstancesParams()
	in.SetTimeout(10 * time.Second)
	go func() {
		defer w.Close()
		wrapper.client.Operations.StreamInstances(in, w)
	}()
	return r, nil
}
