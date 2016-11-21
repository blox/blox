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

const invalidTaskVersion = int64(-100)

// Task defines the structure of the task json received from the event stream
type Task struct {
	ID        *string     `json:"id"`
	Account   *string     `json:"account"`
	Time      *string     `json:"time"`
	Region    *string     `json:"region"`
	Resources []string    `json:"resources"`
	Detail    *TaskDetail `json: "detail"`
}

type TaskDetail struct {
	ClusterARN           *string      `json:"clusterArn"`
	ContainerInstanceARN *string      `json:"containerInstanceArn"`
	Containers           []*Container `json:"containers"`
	CreatedAt            *string      `json:"createdAt"`
	DesiredStatus        *string      `json:"desiredStatus"`
	LastStatus           *string      `json:"lastStatus"`
	Overrides            *Overrides   `json:"overrides"`
	StartedAt            string       `json:"startedAt,omitempty"`
	StartedBy            string       `json:"startedBy,omitempty"`
	StoppedAt            string       `json:"stoppedAt,omitempty"`
	StoppedReason        string       `json:"stoppedReason,omitempty"`
	TaskARN              *string      `json:"taskArn"`
	TaskDefinitionARN    *string      `json:"taskDefinitionArn"`
	UpdatedAt            *string      `json:"updatedAt"`
	Version              *int64       `json:"version"`
}

func (taskDetail *TaskDetail) String() string {
	return fmt.Sprintf("Task %s; Version: %d; Task Definition: %s; %s -> %s; Cluster: %s; Container Instance: %s; Started By: %s; Updated At: %s",
		aws.StringValue(taskDetail.TaskARN),
		aws.Int64Value(taskDetail.Version),
		aws.StringValue(taskDetail.TaskDefinitionARN),
		aws.StringValue(taskDetail.LastStatus),
		aws.StringValue(taskDetail.DesiredStatus),
		aws.StringValue(taskDetail.ClusterARN),
		aws.StringValue(taskDetail.ContainerInstanceARN),
		taskDetail.StartedBy,
		aws.StringValue(taskDetail.UpdatedAt))
}

type Container struct {
	ContainerARN    *string           `json:"containerArn"`
	ExitCode        int64             `json:"exitCode,omitempty"`
	LastStatus      *string           `json:"lastStatus"`
	Name            *string           `json:"name"`
	NetworkBindings []*NetworkBinding `json:"networkBindings,omitempty"`
	Reason          string            `json: "reason,omitempty"`
}

type NetworkBinding struct {
	BindIP        *string `json:"bindIP"`
	ContainerPort *int64  `json: "containerPort"`
	HostPort      *int64  `json: "hostPort"`
	Protocol      string  `json: "protocol,omitempty"`
}

type Overrides struct {
	ContainerOverrides []*ContainerOverrides `json:"containerOverrides"`
	TaskRoleArn        string                `json:"taskRoleArn,omitempty"`
}

type ContainerOverrides struct {
	Command     []string       `json:"command,omitempty"`
	Environment []*Environment `json: "environment,omitempty"`
	Name        *string        `json: "name"`
}

type Environment struct {
	Name  *string `json: "name"`
	Value *string `json: "value"`
}

// GetVersion retrieces the version of the the container instance represented by the instanceJSON string
func (task Task) GetVersion(taskJSON string) (int64, error) {
	t := &Task{}
	err := json.Unmarshal([]byte(taskJSON), t)
	if err != nil {
		return invalidTaskVersion, errors.Wrapf(err, "Error unmarshaling task")
	}
	if t.Detail == nil {
		return invalidTaskVersion, errors.New("Task detail is not set")
	}
	if t.Detail.Version == nil {
		return invalidTaskVersion, errors.New("Task version is not set")
	}
	return aws.Int64Value(t.Detail.Version), nil
}
