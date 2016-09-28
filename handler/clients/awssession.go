package clients

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/pkg/errors"
)

func NewAWSSession() (*session.Session, error) {
	//TODO: env vars? cli inputs?
	sess, err := session.NewSessionWithOptions(session.Options{
		Config:  aws.Config{Region: aws.String("us-east-1")},
		Profile: "event-handler",
	})

	if err != nil {
		return nil, errors.Wrap(err, "Could not load aws session")
	}

	return sess, nil
}
