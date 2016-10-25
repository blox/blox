package types

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
	ClusterArn           *string      `json:"clusterArn"`
	ContainerInstanceArn *string      `json:"containerInstanceArn"`
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
	TaskDefinitionArn    *string      `json:"taskDefinitionArn"`
	UpdatedAt            *string      `json:"updatedAt"`
	Version              *int         `json:"version"`
}

type Container struct {
	ContainerArn    *string           `json:"containerArn"`
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
