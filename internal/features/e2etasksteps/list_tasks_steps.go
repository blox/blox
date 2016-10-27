package e2etasksteps

import (
	"time"

	"github.com/aws/amazon-ecs-event-stream-handler/internal/features/wrappers"
	. "github.com/gucumber/gucumber"
)

func init() {

	eshWrapper := wrappers.NewESHWrapper()

	When(`^I list tasks$`, func() {
		time.Sleep(5 * time.Second)
		eshTasks, err := eshWrapper.ListTasks()
		if err != nil {
			T.Errorf(err.Error())
		}
		for _, t := range eshTasks {
			eshTaskList = append(eshTaskList, *t)
		}
	})

	Then(`^the list tasks response contains at least (\d+) tasks$`, func(numTasks int) {
		if len(eshTaskList) < numTasks {
			T.Errorf("Number of tasks in list tasks response is less than expected")
		}
	})

	And(`^all (\d+) tasks are present in the list tasks response$`, func(numTasks int) {
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
}
