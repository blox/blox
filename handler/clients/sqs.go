package clients

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func NewSQSClient(session *session.Session) *sqs.SQS {
	return sqs.New(session)
}
