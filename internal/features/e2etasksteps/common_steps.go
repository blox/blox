package e2etasksteps

import (
	"github.com/aws/amazon-ecs-event-stream-handler/internal/features/wrappers"
	"github.com/aws/amazon-ecs-event-stream-handler/internal/models"
	"github.com/aws/aws-sdk-go/service/ecs"
	. "github.com/gucumber/gucumber"
)

var (
	ecsTaskList = []ecs.Task{}
	eshTaskList = []models.TaskModel{}
	taskDefnARN = ""
)

func init() {

	ecsWrapper := wrappers.NewECSWrapper()

	BeforeAll(func() {
		var err error
		taskDefnARN, err = ecsWrapper.RegisterSleep360TaskDefinition()
		if err != nil {
			T.Errorf(err.Error())
		}
	})

	Before("@task", func() {
		err := stopAllTasks(ecsWrapper)
		if err != nil {
			T.Errorf("Failed to stop all ECS tasks. Error: %s", err)
		}
	})

	AfterAll(func() {
		err := stopAllTasks(ecsWrapper)
		if err != nil {
			T.Errorf("Failed to stop all ECS tasks. Error: %s", err)
		}
		err = ecsWrapper.DeregisterTaskDefinition(taskDefnARN)
		if err != nil {
			T.Errorf("Failed to deregister task definition. Error: %s", err)
		}
	})

	Given(`^I start (\d+) task(?:|s) in the ECS cluster$`, func(numTasks int) {
		ecsTaskList = nil
		eshTaskList = nil
		for i := 0; i < numTasks; i++ {
			ecsTask, err := ecsWrapper.StartTask(clusterName, taskDefinitionSleep300)
			if err != nil {
				T.Errorf(err.Error())
			}
			ecsTaskList = append(ecsTaskList, ecsTask)
		}
	})
}

func stopAllTasks(ecsWrapper wrappers.ECSWrapper) error {
	taskARNList, err := ecsWrapper.ListTasks(clusterName)
	if err != nil {
		return err
	}
	for _, t := range taskARNList {
		err = ecsWrapper.StopTask(clusterName, *t)
		if err != nil {
			return err
		}
	}
	return nil
}
