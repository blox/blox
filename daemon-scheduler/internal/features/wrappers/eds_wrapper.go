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
	"context"
	"os"

	"github.com/blox/blox/daemon-scheduler/generated/v1/client"
	"github.com/blox/blox/daemon-scheduler/generated/v1/client/operations"
	"github.com/blox/blox/daemon-scheduler/generated/v1/models"
	httptransport "github.com/go-openapi/runtime/client"
)

const (
	defaultSchedulerEndpoint = "localhost:2000"
)

type EDSWrapper struct {
	client *client.BloxDaemonScheduler
}

func NewEDSWrapper() EDSWrapper {
	endpoint := os.Getenv("SCHEDULER_ENDPOINT")
	if len(endpoint) == 0 {
		endpoint = defaultSchedulerEndpoint
	}
	transport := httptransport.New(endpoint, "/v1", []string{"http"})
	httpclient := client.New(transport, nil)
	return EDSWrapper{
		client: httpclient,
	}
}

func (eds EDSWrapper) Ping() error {
	params := operations.NewPingParams()
	_, err := eds.client.Operations.Ping(params)
	if err != nil {
		return err
	}
	return nil
}

func (eds EDSWrapper) CreateEnvironment(in *models.CreateEnvironmentRequest) (*models.Environment, error) {
	params := operations.NewCreateEnvironmentParams()
	params.Body = in
	resp, err := eds.client.Operations.CreateEnvironment(params)
	if err != nil {
		return nil, err
	}
	return resp.Payload, nil
}

func (eds EDSWrapper) GetEnvironment(name *string) (*models.Environment, error) {
	params := operations.NewGetEnvironmentParams()
	params.Name = *name
	resp, err := eds.client.Operations.GetEnvironment(params)
	if err != nil {
		return nil, err
	}
	return resp.Payload, nil
}

func (eds EDSWrapper) DeleteEnvironment(name *string) error {
	params := operations.NewDeleteEnvironmentParams()
	params.Name = *name
	_, err := eds.client.Operations.DeleteEnvironment(params)
	if err != nil {
		return err
	}
	return nil
}

func (eds EDSWrapper) GetDeployment(name *string, id *string) (*models.Deployment, error) {
	params := operations.NewGetDeploymentParams()
	params.Name = *name
	params.ID = *id
	resp, err := eds.client.Operations.GetDeployment(params)
	if err != nil {
		return nil, err
	}
	return resp.Payload, nil
}

func (eds EDSWrapper) ListEnvironments() ([]*models.Environment, error) {
	//TODO: Handle pagination when available
	params := operations.NewListEnvironmentsParams()
	resp, err := eds.client.Operations.ListEnvironments(params)
	if err != nil {
		return nil, err
	}
	return resp.Payload.Items, nil
}

func (eds EDSWrapper) FilterEnvironments(cluster string) ([]*models.Environment, error) {
	//TODO: Handle pagination when available
	params := operations.NewListEnvironmentsParams()
	params.SetCluster(&cluster)
	resp, err := eds.client.Operations.ListEnvironments(params)
	if err != nil {
		return nil, err
	}
	return resp.Payload.Items, nil
}

func (eds EDSWrapper) ListDeployments(name *string) ([]*models.Deployment, error) {
	//TODO: Handle pagination when available
	params := operations.NewListDeploymentsParams()
	params.Name = *name
	resp, err := eds.client.Operations.ListDeployments(params)
	if err != nil {
		return nil, err
	}
	return resp.Payload.Items, nil
}

func (eds EDSWrapper) CreateDeployment(ctx context.Context, envName *string,
	deploymentToken *string) (*models.Deployment, error) {
	params := &operations.CreateDeploymentParams{
		Name:            *envName,
		DeploymentToken: *deploymentToken,
		Context:         ctx,
	}
	resp, err := eds.client.Operations.CreateDeployment(params)
	if err != nil {
		return nil, err
	}
	return resp.Payload, nil
}
