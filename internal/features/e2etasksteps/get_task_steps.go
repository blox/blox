package e2etasksteps

import (
	"time"

	"github.com/aws/amazon-ecs-event-stream-handler/internal/features/wrappers"
	. "github.com/gucumber/gucumber"
)

func init() {

	eshWrapper := wrappers.NewESHWrapper()

	When(`^I get task with the same ARN$`, func() {
		time.Sleep(5 * time.Second)
		if len(ecsTaskList) != 1 {
			T.Errorf("Error memorizing task started using ECS client")
		}
		taskARN := *ecsTaskList[0].TaskArn
		eshTask, err := eshWrapper.GetTask(taskARN)
		if err != nil {
			T.Errorf(err.Error())
		}
		eshTaskList = append(eshTaskList, *eshTask)
	})

	Then(`^I get a task that matches the task started$`, func() {
		if len(ecsTaskList) != 1 || len(eshTaskList) != 1 {
			T.Errorf("Error memorizing results to validate them")
		}
		ecsTask := ecsTaskList[0]
		eshTask := eshTaskList[0]
		err := ValidateTasksMatch(ecsTask, eshTask)
		if err != nil {
			T.Errorf(err.Error())
		}
	})
}
