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

package facade

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/blox/blox/daemon-scheduler/pkg/clients/css/client"
	"github.com/blox/blox/daemon-scheduler/pkg/clients/css/client/operations"
	"github.com/blox/blox/daemon-scheduler/pkg/clients/css/models"
	"github.com/pkg/errors"
)

// ClusterState defines methods to get cluster and task state
type ClusterState interface {
	ListInstances(cluster string) ([]*models.ContainerInstance, error)
	ListTasks(cluster string) ([]*models.Task, error)
}

type clusterState struct {
	client *client.BloxCSS
}

func NewClusterState(css *client.BloxCSS) (ClusterState, error) {
	if css == nil {
		return nil, errors.New("CSS client should not be nil")
	}
	return clusterState{
		client: css,
	}, nil
}

func (c clusterState) ListInstances(cluster string) ([]*models.ContainerInstance, error) {
	req := operations.NewFilterInstancesParams()
	req.SetCluster(cluster)

	resp, err := c.client.Operations.FilterInstances(req)
	if err != nil {
		return nil, errors.Wrap(err, "List instances failed")
	}

	return resp.Payload.Items, nil
}

func (c clusterState) ListTasks(cluster string) ([]*models.Task, error) {
	req := operations.NewListTasksParams()
	resp, err := c.client.Operations.ListTasks(req)
	if err != nil {
		return nil, errors.Wrapf(err, "Error calling ListTasks with req %v", req)
	}
	tasks := []*models.Task{}
	for _, task := range resp.Payload.Items {
		if aws.StringValue(task.ClusterARN) == cluster {
			tasks = append(tasks, task)
		}
	}
	return tasks, nil
}
