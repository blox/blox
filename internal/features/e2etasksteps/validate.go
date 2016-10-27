package e2etasksteps

import (
	"github.com/aws/amazon-ecs-event-stream-handler/internal/models"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/pkg/errors"
)

func ValidateTasksMatch(ecsTask ecs.Task, eshTask models.TaskModel) error {
	if *ecsTask.TaskArn != *eshTask.Detail.TaskArn ||
		*ecsTask.ClusterArn != *eshTask.Detail.ClusterArn ||
		*ecsTask.ContainerInstanceArn != *eshTask.Detail.ContainerInstanceArn {
		return errors.New("Tasks don't match")
	}
	return nil
}

func ValidateListContainsTask(ecsTask ecs.Task, eshTaskList []models.TaskModel) error {
	taskARN := *ecsTask.TaskArn
	var eshTask models.TaskModel
	for _, t := range eshTaskList {
		if *t.Detail.TaskArn == taskARN {
			eshTask = t
			break
		}
	}
	if eshTask.Detail == nil || eshTask.Detail.TaskArn == nil {
		return errors.Errorf("Task with ARN '%s' not found in response", taskARN)
	}
	return ValidateTasksMatch(ecsTask, eshTask)
}
