package e2etasksteps

import (
	"github.com/aws/amazon-ecs-event-stream-handler/internal/models"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/pkg/errors"
)

func ValidateTasksMatch(ecsTask ecs.Task, eshTask models.TaskModel) error {
	if *ecsTask.TaskArn != *eshTask.TaskARN ||
		*ecsTask.ClusterArn != *eshTask.ClusterARN ||
		*ecsTask.ContainerInstanceArn != *eshTask.ContainerInstanceARN {
		return errors.New("Tasks don't match")
	}
	return nil
}

func ValidateListContainsTask(ecsTask ecs.Task, eshTaskList []models.TaskModel) error {
	taskARN := *ecsTask.TaskArn
	var eshTask models.TaskModel
	for _, t := range eshTaskList {
		if *t.TaskARN == taskARN {
			eshTask = t
			break
		}
	}
	if eshTask.TaskARN == nil {
		return errors.Errorf("Task with ARN '%s' not found in response", taskARN)
	}
	return ValidateTasksMatch(ecsTask, eshTask)
}
