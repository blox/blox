package types

// Task defines the structure of the task json received from the event stream
type Task struct {
	ID        string   `json:"id"`
	Account   string   `json:"account"`
	Time      string   `json:"time"`
	Region    string   `json:"region"`
	Resources []string `json:"resources"`
	Detail    struct {
		ClusterArn           string `json:"clusterArn"`
		ContainerInstanceArn string `json:"containerInstanceArn"`
		Containers           []struct {
			ContainerArn    string `json:"containerArn"`
			ExitCode        *int   `json:"exitCode,omitempty"`
			LastStatus      string `json:"lastStatus"`
			Name            string `json:"name"`
			NetworkBindings []struct {
				BindIP        string  `json:"bindIP"`
				ContainerPort int     `json: "containerPort"`
				HostPort      int     `json: "hostPort"`
				Protocol      *string `json: "protocol,omitempty"`
			} `json:"networkBindings,omitempty"`
			Reason *string `json: "reason,omitempty"`
		} `json:"containers"`
		CreatedAt     string `json:"createdAt"`
		DesiredStatus string `json:"desiredStatus"`
		LastStatus    string `json:"lastStatus"`
		Overrides     struct {
			ContainerOverrides []struct {
				Command     []string `json:"command,omitempty"`
				Environment []struct {
					Name  string `json: "name"`
					Value string `json: "value"`
				} `json: "environment,omitempty"`
				Name string `json: "name"`
			} `json:"containerOverrides"`
			TaskRoleArn *string `json:"taskRoleArn,omitempty"`
		} `json:"overrides"`
		StartedAt         *string `json:"startedAt,omitempty"`
		StartedBy         *string `json:"startedBy,omitempty"`
		StoppedAt         *string `json:"stoppedAt,omitempty"`
		StoppedReason     *string `json:"stoppedReason,omitempty"`
		UpdatedAt         string  `json:"updatedAt"`
		TaskArn           string  `json:"taskArn"`
		TaskDefinitionArn string  `json:"taskDefinitionArn"`
		Version           int     `json:"version"`
	} `json: "detail"`
}
