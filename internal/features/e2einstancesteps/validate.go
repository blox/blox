package e2einstancesteps

import (
	"github.com/aws/amazon-ecs-event-stream-handler/internal/models"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/pkg/errors"
)

func ValidateInstancesMatch(ecsInstance ecs.ContainerInstance, eshInstance models.ContainerInstanceModel) error {
	if *ecsInstance.ContainerInstanceArn != *eshInstance.ContainerInstanceARN ||
		*ecsInstance.Status != *eshInstance.Status {
		return errors.New("Container instances don't match")
	}
	return nil
}

func ValidateListContainsInstance(ecsInstance ecs.ContainerInstance, eshInstanceList []models.ContainerInstanceModel) error {
	instanceARN := *ecsInstance.ContainerInstanceArn
	var eshInstance models.ContainerInstanceModel
	for _, i := range eshInstanceList {
		if *i.ContainerInstanceARN == instanceARN {
			eshInstance = i
			break
		}
	}
	if eshInstance.ContainerInstanceARN == nil {
		return errors.Errorf("Instance with ARN '%s' not found in response", instanceARN)
	}
	return ValidateInstancesMatch(ecsInstance, eshInstance)
}
