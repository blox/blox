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

package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	taskARN     = "arn:aws:ecs:us-east-1:12345678912:task/c024d145-093b-499a-9b14-5baf273f5835"
	instanceARN = "arn:aws:us-east-1:123456789123:container-instance/4b6d45ea-a4b4-4269-9d04-3af6ddfdc597"
)

func TestNewTaskEmptyTaskARN(t *testing.T) {
	_, err := NewTask("", instanceARN)
	assert.NotNil(t, err, "Expected an error creating a new task with empty task ARN")
}

func TestNewTaskEmptyInstanceARN(t *testing.T) {
	_, err := NewTask(taskARN, "")
	assert.NotNil(t, err, "Expected an error creating a new task with empty instance ARN")
}

func TestNewTask(t *testing.T) {
	task, err := NewTask(taskARN, instanceARN)
	assert.Nil(t, err, "Unexpected error creating a new task")
	expectedTask := &Task{
		TaskARN:     taskARN,
		InstanceARN: instanceARN,
	}
	assert.Equal(t, expectedTask, task)
}
