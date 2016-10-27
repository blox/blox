package e2einstancesteps

import (
	"github.com/aws/amazon-ecs-event-stream-handler/internal/features/wrappers"
	"github.com/aws/amazon-ecs-event-stream-handler/internal/models"
	"github.com/aws/aws-sdk-go/service/ecs"
	. "github.com/gucumber/gucumber"
)

var (
	ecsContainerInstanceList = []ecs.ContainerInstance{}
	eshContainerInstanceList = []models.ContainerInstanceModel{}
)

func init() {

	ecsWrapper := wrappers.NewECSWrapper()

	Given(`^I have some instances registered with the ECS cluster$`, func() {
		ecsContainerInstanceList = nil
		eshContainerInstanceList = nil

		instanceARNs, err := ecsWrapper.ListContainerInstances(clusterName)
		if err != nil {
			T.Errorf(err.Error())
		}
		if len(instanceARNs) < 1 {
			T.Errorf("No container instances registered to the cluster '%s'", clusterName)
		}
		for _, instanceARN := range instanceARNs {
			ecsInstance, err := ecsWrapper.DescribeContainerInstance(clusterName, *instanceARN)
			if err != nil {
				T.Errorf(err.Error())
			}
			ecsContainerInstanceList = append(ecsContainerInstanceList, ecsInstance)
		}
	})
}
