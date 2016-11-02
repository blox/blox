package v1

import (
	"errors"

	"github.com/aws/amazon-ecs-event-stream-handler/handler/api/v1/models"
	"github.com/aws/amazon-ecs-event-stream-handler/handler/types"
)

func validateContainerInstance(instance types.ContainerInstance) error {
	if instance.Account == nil {
		return errors.New("Instance account cannot be empty")
	}
	// TODO: Validate inner structs in instance.Detail
	detail := instance.Detail
	if detail == nil || detail.AgentConnected == nil || detail.ClusterARN == nil ||
		detail.ContainerInstanceARN == nil || detail.PendingTasksCount == nil ||
		detail.RegisteredResources == nil || detail.RemainingResources == nil ||
		detail.RunningTasksCount == nil || detail.Status == nil || detail.Version == nil ||
		detail.VersionInfo == nil || detail.UpdatedAt == nil {
		return errors.New("Instance detail is invalid")
	}
	if instance.ID == nil {
		return errors.New("Instance id cannot be emoty")
	}
	if instance.Region == nil {
		return errors.New("Instance region cannot be empty")
	}
	if instance.Resources == nil {
		return errors.New("Instance resources cannot be empty")
	}
	if instance.Time == nil {
		return errors.New("Instance time cannot be empty")
	}
	return nil
}

func ToContainerInstanceModel(instance types.ContainerInstance) (models.ContainerInstanceModel, error) {
	err := validateContainerInstance(instance)
	if err != nil {
		return models.ContainerInstanceModel{}, err
	}
	regRes := make([]*models.ContainerInstanceDetailRegisteredResourceModel, len(instance.Detail.RegisteredResources))
	for i := range instance.Detail.RegisteredResources {
		r := instance.Detail.RegisteredResources[i]
		regRes[i] = &models.ContainerInstanceDetailRegisteredResourceModel{
			Name:  r.Name,
			Type:  r.Type,
			Value: r.Value,
		}
	}

	remRes := make([]*models.ContainerInstanceDetailRemainingResourceModel, len(instance.Detail.RemainingResources))
	for i := range instance.Detail.RegisteredResources {
		r := instance.Detail.RemainingResources[i]
		remRes[i] = &models.ContainerInstanceDetailRemainingResourceModel{
			Name:  r.Name,
			Type:  r.Type,
			Value: r.Value,
		}
	}

	versionInfo := models.ContainerInstanceDetailVersionInfoModel{
		AgentHash:     instance.Detail.VersionInfo.AgentHash,
		AgentVersion:  instance.Detail.VersionInfo.AgentVersion,
		DockerVersion: instance.Detail.VersionInfo.DockerVersion,
	}

	pendingTaskCount := int32(*instance.Detail.PendingTasksCount)
	runningTaskCount := int32(*instance.Detail.RunningTasksCount)
	version := int32(*instance.Detail.Version)
	detail := models.ContainerInstanceDetailModel{
		AgentConnected:       instance.Detail.AgentConnected,
		AgentUpdateStatus:    instance.Detail.AgentUpdateStatus,
		ClusterArn:           instance.Detail.ClusterARN,
		ContainerInstanceArn: instance.Detail.ContainerInstanceARN,
		Ec2InstanceID:        instance.Detail.EC2InstanceID,
		PendingTasksCount:    &pendingTaskCount,
		RegisteredResources:  regRes,
		RemainingResources:   remRes,
		RunningTasksCount:    &runningTaskCount,
		Status:               instance.Detail.Status,
		Version:              &version,
		VersionInfo:          &versionInfo,
		UpdatedAt:            instance.Detail.UpdatedAt,
	}

	if instance.Detail.Attributes != nil {
		attributes := make([]*models.ContainerInstanceDetailAttributeModel, len(instance.Detail.Attributes))
		for i := range instance.Detail.Attributes {
			a := instance.Detail.Attributes[i]
			attributes[i] = &models.ContainerInstanceDetailAttributeModel{
				Name:  a.Name,
				Value: a.Value,
			}
		}
		detail.Attributes = attributes
	}

	return models.ContainerInstanceModel{
		ID:        instance.ID,
		Account:   instance.Account,
		Time:      instance.Time,
		Region:    instance.Region,
		Resources: instance.Resources,
		Detail:    &detail,
	}, nil
}

func validateTaskModel(task types.Task) error {
	if task.Account == nil {
		return errors.New("Task account cannot be empty")
	}
	// TODO: Validate inner structs in task.Detail
	detail := task.Detail
	if detail == nil || detail.ClusterARN == nil || detail.ContainerInstanceARN == nil ||
		detail.Containers == nil || detail.CreatedAt == nil || detail.DesiredStatus == nil ||
		detail.LastStatus == nil || detail.Overrides == nil || detail.TaskArn == nil ||
		detail.TaskDefinitionARN == nil || detail.UpdatedAt == nil || detail.Version == nil {
		return errors.New("Task detail is invalid")
	}
	if task.ID == nil {
		return errors.New("Task id cannot be emoty")
	}
	if task.Region == nil {
		return errors.New("Task region cannot be empty")
	}
	if task.Resources == nil {
		return errors.New("Task resources cannot be empty")
	}
	if task.Time == nil {
		return errors.New("Task time cannot be empty")
	}
	return nil
}

func ToTaskModel(task types.Task) (models.TaskModel, error) {
	err := validateTaskModel(task)
	if err != nil {
		return models.TaskModel{}, err
	}

	containers := make([]*models.TaskDetailContainerModel, len(task.Detail.Containers))
	for i := range task.Detail.Containers {
		c := task.Detail.Containers[i]
		exitCode := int32(c.ExitCode)
		containers[i] = &models.TaskDetailContainerModel{
			ContainerArn: c.ContainerARN,
			ExitCode:     exitCode,
			LastStatus:   c.LastStatus,
			Name:         c.Name,
			Reason:       c.Reason,
		}
		if c.NetworkBindings != nil {
			networkBindings := make([]*models.TaskDetailNetworkBindingModel, len(c.NetworkBindings))
			for j := range c.NetworkBindings {
				n := c.NetworkBindings[j]
				containerPort := int32(*n.ContainerPort)
				hostPort := int32(*n.HostPort)
				networkBindings[j] = &models.TaskDetailNetworkBindingModel{
					BindIP:        n.BindIP,
					ContainerPort: &containerPort,
					HostPort:      &hostPort,
					Protocol:      n.Protocol,
				}
			}
			containers[i].NetworkBindings = networkBindings
		}
	}

	containerOverrides := make([]*models.TaskDetailContainerOverridesModel, len(task.Detail.Overrides.ContainerOverrides))
	for i := range task.Detail.Overrides.ContainerOverrides {
		c := task.Detail.Overrides.ContainerOverrides[i]
		containerOverrides[i] = &models.TaskDetailContainerOverridesModel{
			Command: c.Command,
			Name:    c.Name,
		}
		if c.Environment != nil {
			env := make([]*models.TaskDetailEnvironmentModel, len(c.Environment))
			for j := range c.Environment {
				e := c.Environment[j]
				env[j] = &models.TaskDetailEnvironmentModel{
					Name:  e.Name,
					Value: e.Value,
				}
			}
			containerOverrides[i].Environment = env
		}
	}

	overrides := models.TaskDetailOverridesModel{
		ContainerOverrides: containerOverrides,
		TaskRoleArn:        task.Detail.Overrides.TaskRoleArn,
	}

	version := int32(*task.Detail.Version)
	detail := models.TaskDetailModel{
		ClusterArn:           task.Detail.ClusterARN,
		ContainerInstanceArn: task.Detail.ContainerInstanceARN,
		Containers:           containers,
		CreatedAt:            task.Detail.CreatedAt,
		DesiredStatus:        task.Detail.DesiredStatus,
		LastStatus:           task.Detail.LastStatus,
		Overrides:            &overrides,
		StartedAt:            task.Detail.StartedAt,
		StartedBy:            task.Detail.StartedBy,
		StoppedAt:            task.Detail.StoppedAt,
		StoppedReason:        task.Detail.StoppedReason,
		TaskArn:              task.Detail.TaskArn,
		TaskDefinitionArn:    task.Detail.TaskDefinitionARN,
		UpdatedAt:            task.Detail.UpdatedAt,
		Version:              &version,
	}

	return models.TaskModel{
		ID:        task.ID,
		Account:   task.Account,
		Time:      task.Time,
		Region:    task.Region,
		Resources: task.Resources,
		Detail:    &detail,
	}, nil
}
