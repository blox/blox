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
	log "github.com/cihub/seelog"
)

const ecsEndpointEnvVarName = "ECS_ENDPOINT"

func NewECSClient(sess *session.Session) *ecs.ECS {
	// TODO: Use session passed in args and get rid of the env var
	endpoint := os.Getenv(ecsEndpointEnvVarName)
	if endpoint == "" {
		return ecs.New(sess)
	}
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Endpoint: aws.String(endpoint),
		},
	})
	if err != nil {
		log.Critical("Error initializing ecs client")
		return nil
	}

	return ecs.New(sess)
}
