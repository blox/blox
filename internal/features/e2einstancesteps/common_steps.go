package e2einstancesteps

import (
	"github.com/aws/amazon-ecs-event-stream-handler/internal/features/wrappers"
	"github.com/aws/amazon-ecs-event-stream-handler/internal/models"
	"github.com/aws/aws-sdk-go/service/ecs"
	. "github.com/gucumber/gucumber"
)

var (
	// Lists to memorize results required for the subsequent steps
	ecsContainerInstanceList = []ecs.ContainerInstance{}
	eshContainerInstanceList = []models.ContainerInstanceModel{}
	exceptionList            = []string{}
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

	Then(`^I get a (.+?) instance exception$`, func(exception string) {
		if len(exceptionList) != 1 {
			T.Errorf("Error memorizing exception")
		}
		if exception != exceptionList[0] {
			T.Errorf("Expected exception '%s' but got '%s'", exception, exceptionList[0])
		}
	})
}
