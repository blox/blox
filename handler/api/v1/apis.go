package v1

import (
	"github.com/aws/amazon-ecs-event-stream-handler/handler/store"
)

type APIs struct {
	TaskApis              TaskAPIs
	ContainerInstanceApis ContainerInstanceAPIs
}

func NewAPIs(stores store.Stores) APIs {
	return APIs{
		TaskApis:              NewTaskAPIs(stores.TaskStore),
		ContainerInstanceApis: NewContainerInstanceAPIs(stores.ContainerInstanceStore),
	}
}
