package clients

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/pkg/errors"
)

func NewAWSSession() (*session.Session, error) {
	sess, err := session.NewSession()

	if err != nil {
		return nil, errors.Wrap(err, "Could not load aws session")
	}

	return sess, nil
}
