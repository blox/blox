package types

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
)

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
	TaskArn              *string      `json:"taskArn"`
	TaskDefinitionARN    *string      `json:"taskDefinitionArn"`
	UpdatedAt            *string      `json:"updatedAt"`
	Version              *int         `json:"version"`
}

func (taskDetail *TaskDetail) String() string {
	return fmt.Sprintf("Task %s; Version: %d; Task Definition: %s; %s -> %s; Cluster: %s; Container Instance: %s; Started By: %s; Updated At: %s",
		aws.StringValue(taskDetail.TaskArn),
		aws.IntValue(taskDetail.Version),
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
	ExitCode        int               `json:"exitCode,omitempty"`
	LastStatus      *string           `json:"lastStatus"`
	Name            *string           `json:"name"`
	NetworkBindings []*NetworkBinding `json:"networkBindings,omitempty"`
	Reason          string            `json: "reason,omitempty"`
}

type NetworkBinding struct {
	BindIP        *string `json:"bindIP"`
	ContainerPort *int    `json: "containerPort"`
	HostPort      *int    `json: "hostPort"`
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
