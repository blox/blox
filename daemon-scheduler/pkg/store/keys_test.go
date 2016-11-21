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

package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	envName = "testEnv"
	envKey  = environmentKeyPrefix + envName
	taskARN = "arn:aws:ecs:us-east-1:12345678912:task/c024d145-093b-499a-9b14-5baf273f5835"
	taskKey = envKey + taskKeyConnector + taskARN
)

func TestGenerateEnvironmentKeyEmptyEnvName(t *testing.T) {
	_, err := GenerateEnvironmentKey("")
	assert.NotNil(t, err, "Expected an error generating environment key when environment name is empty")
}

func TestGenerateEnvironmentKey(t *testing.T) {
	key, err := GenerateEnvironmentKey(envName)
	assert.Nil(t, err, "Unexpected error when generating environment key")
	assert.Equal(t, envKey, key)
}

func TestGenerateTaskKeyEmptyEnvName(t *testing.T) {
	_, err := GenerateTaskKey("", taskARN)
	assert.NotNil(t, err, "Expected an error generating task key when environment name is empty")
}

func TestGenerateTaskKeyEmptyTaskARN(t *testing.T) {
	_, err := GenerateTaskKey(envName, "")
	assert.NotNil(t, err, "Expected an error generating task key when task ARN is empty")
}

func TestGenerateTaskKey(t *testing.T) {
	key, err := GenerateTaskKey(envName, taskARN)
	assert.Nil(t, err, "Unexpected error when generating task key")
	assert.Equal(t, taskKey, key)
}
