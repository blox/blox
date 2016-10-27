package e2etasksteps

import (
	"time"

	"github.com/aws/amazon-ecs-event-stream-handler/internal/features/wrappers"
	. "github.com/gucumber/gucumber"
)

func init() {

	ecsWrapper := wrappers.NewECSWrapper()
	eshWrapper := wrappers.NewESHWrapper()

	When(`^I filter tasks by (.+?) status$`, func(status string) {
		time.Sleep(5 * time.Second)
		eshTasks, err := eshWrapper.FilterTasksByStatus(status)
		if err != nil {
			T.Errorf(err.Error())
		}
		for _, t := range eshTasks {
			eshTaskList = append(eshTaskList, *t)
		}
	})

	Then(`^the filter tasks response contains at least (\d+) tasks$`, func(numTasks int) {
		if len(eshTaskList) < numTasks {
			T.Errorf("Number of tasks in filter tasks response is less than expected")
		}
	})

	And(`^all (\d+) tasks are present in the filter tasks response$`, func(numTasks int) {
		if len(ecsTaskList) != numTasks {
			T.Errorf("Error memorizing tasks started using ECS client")
		}
		for _, t := range ecsTaskList {
			err := ValidateListContainsTask(t, eshTaskList)
			if err != nil {
				T.Errorf(err.Error())
			}
		}
	})

	And(`^I stop the (\d+) tasks in the ECS cluster$`, func(numTasks int) {
		if len(ecsTaskList) != numTasks {
			T.Errorf("Error memorizing tasks started using ECS client")
		}
		for _, t := range ecsTaskList {
			err := ecsWrapper.StopTask(clusterName, *t.TaskArn)
			if err != nil {
				T.Errorf(err.Error())
			}
		}
	})
}
