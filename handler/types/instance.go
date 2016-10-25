package types

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
	ClusterArn           *string      `json:"clusterArn"`
	ContainerInstanceArn *string      `json:"containerInstanceArn"`
	Ec2InstanceID        string       `json:"ec2InstanceId,omitempty"`
	PendingTasksCount    *int         `json:"pendingTasksCount"`
	RegisteredResources  []*Resource  `json:"registeredResources"`
	RemainingResources   []*Resource  `json:"remainingResources"`
	RunningTasksCount    *int         `json:"runningTasksCount"`
	Status               *string      `json:"status"`
	Version              *int         `json:"version"`
	VersionInfo          *VersionInfo `json:"versionInfo"`
	UpdatedAt            *string      `json:"updatedAt"`
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
