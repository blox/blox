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
	"github.com/blox/blox/cluster-state-service/swagger/v1/generated/client"
	"github.com/blox/blox/cluster-state-service/swagger/v1/generated/client/operations"
	"github.com/blox/blox/cluster-state-service/swagger/v1/generated/models"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/pkg/errors"
)

// CSSWrapper defines methods to access wrapper methods to call CSS APIs
type CSSWrapper interface {
	ListInstances(cluster *string) ([]*models.ContainerInstance, error)
	GetInstance(cluster *string, instanceARN *string) (*models.ContainerInstance, error)
	ListTasks(cluster *string) ([]*models.Task, error)
	GetTask(cluster *string, taskARN *string) (*models.Task, error)
}

type cssClientWrapper struct {
	client *client.BloxCSS
}

// NewCSSWrapper returns a new CSSWrapper for the canary
func NewCSSWrapper(clusterStateServiceEndpoint string) (CSSWrapper, error) {
	cssClient, err := newCSSClient(clusterStateServiceEndpoint)
	if err != nil {
		return nil, err
	}
	return cssClientWrapper{
		client: cssClient,
	}, nil
}

func newCSSClient(clusterStateServiceEndpoint string) (*client.BloxCSS, error) {
	if clusterStateServiceEndpoint == "" {
		return nil, errors.New("The address of the cluster-state-service endpoint had to be set to initialize the canary. ")
	}
	cssClient := client.NewHTTPClient(nil)
	cssTransport := httptransport.New(clusterStateServiceEndpoint, "/v1", []string{"http"})
	cssClient.SetTransport(cssTransport)
	return cssClient, nil
}

func (wrapper cssClientWrapper) ListInstances(cluster *string) ([]*models.ContainerInstance, error) {
	req := operations.NewListInstancesParams()
	req.SetCluster(cluster)

	resp, err := wrapper.client.Operations.ListInstances(req)
	if err != nil {
		return nil, errors.Wrapf(err, "Error calling ListInstances for cluster '%s'. ", *cluster)
	}

	return resp.Payload.Items, nil
}

func (wrapper cssClientWrapper) GetInstance(cluster *string, instanceARN *string) (*models.ContainerInstance, error) {
	req := operations.NewGetInstanceParams()
	req.SetCluster(*cluster)
	req.SetArn(*instanceARN)

	resp, err := wrapper.client.Operations.GetInstance(req)
	if err != nil {
		return nil, errors.Wrapf(err, "Error calling GetInstance for instance with ARN '%s' belonging to cluster '%s'. ",
			*instanceARN, *cluster)
	}

	return resp.Payload, nil
}

func (wrapper cssClientWrapper) ListTasks(cluster *string) ([]*models.Task, error) {
	req := operations.NewListTasksParams()
	req.SetCluster(cluster)

	resp, err := wrapper.client.Operations.ListTasks(req)
	if err != nil {
		return nil, errors.Wrapf(err, "Error calling ListTasks for cluster '%s'. ", *cluster)
	}

	return resp.Payload.Items, nil
}

func (wrapper cssClientWrapper) GetTask(cluster *string, taskARN *string) (*models.Task, error) {
	req := operations.NewGetTaskParams()
	req.SetCluster(*cluster)
	req.SetArn(*taskARN)

	resp, err := wrapper.client.Operations.GetTask(req)
	if err != nil {
		return nil, errors.Wrapf(err, "Error calling GetTask for task with ARN '%s' belonging to cluster '%s'. ",
			*taskARN, *cluster)
	}

	return resp.Payload, nil
}
