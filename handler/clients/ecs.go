package clients

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

func NewECSClient(session *session.Session) *ecs.ECS {
	return ecs.New(session)
}
