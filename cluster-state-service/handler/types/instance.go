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

package types

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/pkg/errors"
)

const invalidInstanceVersion = int64(-100)

// ContainerInstance defines the structure of the container instance json received from the event stream
type ContainerInstance struct {
	ID        *string         `json:"id"`
	Account   *string         `json:"account"`
	Time      *string         `json:"time"`
	Region    *string         `json:"region"`
	Resources []string        `json:"resources"`
	Detail    *InstanceDetail `json: "detail"`
}

type InstanceDetail struct {
	AgentConnected       *bool        `json:"agentConnected"`
	AgentUpdateStatus    string       `json:"agentUpdateStatus,omitempty"`
	Attributes           []*Attribute `json:"attributes,omitempty"`
	ClusterARN           *string      `json:"clusterArn"`
	ContainerInstanceARN *string      `json:"containerInstanceArn"`
	EC2InstanceID        string       `json:"ec2InstanceId,omitempty"`
	RegisteredResources  []*Resource  `json:"registeredResources"`
	RemainingResources   []*Resource  `json:"remainingResources"`
	Status               *string      `json:"status"`
	Version              *int64       `json:"version"`
	VersionInfo          *VersionInfo `json:"versionInfo"`
	UpdatedAt            *string      `json:"updatedAt"`
}

func (instanceDetail *InstanceDetail) String() string {
	return fmt.Sprintf("Instance %s; Version: %d; Cluster: %s; EC2 Instance ID: %s; AgentConnected: %t; Status: %s; Updated at: %s",
		aws.StringValue(instanceDetail.ContainerInstanceARN),
		aws.Int64Value(instanceDetail.Version),
		aws.StringValue(instanceDetail.ClusterARN),
		instanceDetail.EC2InstanceID,
		aws.BoolValue(instanceDetail.AgentConnected),
		aws.StringValue(instanceDetail.Status),
		aws.StringValue(instanceDetail.UpdatedAt))
}

type Attribute struct {
	Name  *string `json:"name`
	Value *string `json: "value"`
}

type Resource struct {
	Name  *string `json:"name"`
	Type  *string `json:"type"`
	Value *string `json:"value"`
}

type VersionInfo struct {
	AgentHash     string `json:"agentHash,omitempty"`
	AgentVersion  string `json:"agentVersion,omitempty"`
	DockerVersion string `json:"dockerVersion,omitempty"`
}

// GetVersion retrieces the version of the the container instance represented by the instanceJSON string
func (instance ContainerInstance) GetVersion(instanceJSON string) (int64, error) {
	ins := &ContainerInstance{}
	err := json.Unmarshal([]byte(instanceJSON), ins)
	if err != nil {
		return invalidInstanceVersion, errors.Wrapf(err, "Error unmarshaling instance")
	}
	if ins.Detail == nil {
		return invalidInstanceVersion, errors.New("Instance detail is not set")
	}
	if ins.Detail.Version == nil {
		return invalidInstanceVersion, errors.New("Instance version is not set")
	}
	return aws.Int64Value(ins.Detail.Version), nil
}
