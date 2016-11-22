// Copyright 2016 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package clients

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/pkg/errors"
)

const (
	ecsEndpointEnvVarName = "ECS_ENDPOINT"
)

func getECSEndpoint() (string, error) {
	endpoint := os.Getenv(ecsEndpointEnvVarName)
	if endpoint == "" {
		return "", errors.Errorf("Empty endpoint. Please specify the ECS endpoint name using the '%s' environment variable", ecsEndpointEnvVarName)
	}

	return endpoint, nil
}

func newAWSSession() (*session.Session, error) {
	var sess *session.Session
	var err error
	if endpoint, err := getECSEndpoint(); err != nil {
		sess, err = session.NewSession()
	} else {
		sess, err = session.NewSessionWithOptions(session.Options{
			Config: aws.Config{
				Endpoint: aws.String(endpoint),
			},
		})
	}
	return sess, err
}

func NewECSClient() (*ecs.ECS, error) {
	sess, err := newAWSSession()
	if err != nil {
		return nil, err
	}
	return ecs.New(sess), nil
}
