package e2einstancesteps

import (
	"github.com/aws/amazon-ecs-event-stream-handler/internal/features/wrappers"
	. "github.com/gucumber/gucumber"
)

const (
	nonExistentInstanceARN = "arn:aws:ecs:us-east-1:123456789012:container-instance/950bc492-81a5-4fed-9419-8edeab4769e5"
)

func init() {

	ecsWrapper := wrappers.NewECSWrapper()
	eshWrapper := wrappers.NewESHWrapper()

	Given(`^I have an instance registered with the ECS cluster$`, func() {
		ecsContainerInstanceList = nil
		eshContainerInstanceList = nil

		instanceARNs, err := ecsWrapper.ListContainerInstances(clusterName)
		if err != nil {
			T.Errorf(err.Error())
		}
		if len(instanceARNs) < 1 {
			T.Errorf("No container instances registered to the cluster '%s'", clusterName)
		}
		ecsInstance, err := ecsWrapper.DescribeContainerInstance(clusterName, *instanceARNs[0])
		if err != nil {
			T.Errorf(err.Error())
		}
		ecsContainerInstanceList = append(ecsContainerInstanceList, ecsInstance)
	})

	When(`^I get instance with the cluster name and instance ARN$`, func() {
		if len(ecsContainerInstanceList) != 1 {
			T.Errorf("Error memorizing container instance registered to ECS")
		}
		instanceARN := *ecsContainerInstanceList[0].ContainerInstanceArn
		eshInstance, err := eshWrapper.GetInstance(clusterName, instanceARN)
		if err != nil {
			T.Errorf(err.Error())
		}
		eshContainerInstanceList = append(eshContainerInstanceList, *eshInstance)
	})

	Then(`^I get an instance that matches the registered instance$`, func() {
		if len(ecsContainerInstanceList) != 1 || len(eshContainerInstanceList) != 1 {
			T.Errorf("Error memorizing results to validate them")
		}
		ecsInstance := ecsContainerInstanceList[0]
		eshInstance := eshContainerInstanceList[0]
		err := ValidateInstancesMatch(ecsInstance, eshInstance)
		if err != nil {
			T.Errorf(err.Error())
		}
	})

	When(`^I try to get instance with a non-existent ARN$`, func() {
		exceptionList = nil
		exception, err := eshWrapper.TryGetInstance(nonExistentInstanceARN)
		if err != nil {
			T.Errorf(err.Error())
		}
		exceptionList = append(exceptionList, exception)
	})
}
