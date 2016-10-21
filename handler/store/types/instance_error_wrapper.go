package types

import (
	"github.com/aws/amazon-ecs-event-stream-handler/handler/types"
)

type ContainerInstanceErrorWrapper struct {
	ContainerInstance types.ContainerInstance
	Err               error
}
