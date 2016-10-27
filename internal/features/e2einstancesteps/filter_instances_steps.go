package e2einstancesteps

import (
	"github.com/aws/amazon-ecs-event-stream-handler/internal/features/wrappers"
	. "github.com/gucumber/gucumber"
)

func init() {

	eshWrapper := wrappers.NewESHWrapper()

	When(`^I filter instances by the same ECS cluster name$`, func() {
		eshInstances, err := eshWrapper.FilterInstancesByClusterName(clusterName)
		if err != nil {
			T.Errorf(err.Error())
		}
		for _, i := range eshInstances {
			eshContainerInstanceList = append(eshContainerInstanceList, *i)
		}
	})

	Then(`^the filter instances response contains all the instances registered with the cluster$`, func() {
		if len(ecsContainerInstanceList) != len(eshContainerInstanceList) {
			T.Errorf("Unexpected number of instances in the filter instances response")
		}
		for _, ecsInstance := range ecsContainerInstanceList {
			err := ValidateListContainsInstance(ecsInstance, eshContainerInstanceList)
			if err != nil {
				T.Errorf(err.Error())
			}
		}
	})

}
