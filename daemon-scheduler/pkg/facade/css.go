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
	"github.com/blox/blox/cluster-state-service/swagger/v1/generated/client"
	"github.com/blox/blox/cluster-state-service/swagger/v1/generated/client/operations"
	"github.com/blox/blox/cluster-state-service/swagger/v1/generated/models"
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
	req := operations.NewListInstancesParams()
	req.SetCluster(&cluster)

	resp, err := c.client.Operations.ListInstances(req)
	if err != nil {
		return nil, errors.Wrapf(err, "Error calling ListInstances with cluster %v", cluster)
	}

	return resp.Payload.Items, nil
}

func (c clusterState) ListTasks(cluster string) ([]*models.Task, error) {
	req := operations.NewListTasksParams()
	req.SetCluster(&cluster)

	resp, err := c.client.Operations.ListTasks(req)
	if err != nil {
		return nil, errors.Wrapf(err, "Error calling ListTasks with cluster %v", cluster)
	}
	return resp.Payload.Items, nil
}
