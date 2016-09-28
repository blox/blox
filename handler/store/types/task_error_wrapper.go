package types

import (
	"github.com/aws/amazon-ecs-event-stream-handler/handler/types"
)

type TaskErrorWrapper struct {
	Task types.Task
	Err  error
}
