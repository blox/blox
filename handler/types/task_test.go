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
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetVersionValidTask(t *testing.T) {
	tsk := Task{}

	version := int64(1)
	task := Task{
		Detail: &TaskDetail{
			Version: &version,
		},
	}
	taskJSON := marshalTask(t, task)

	extractedVersion, err := tsk.GetVersion(taskJSON)
	assert.Nil(t, err, "Unexpected error getting task version")
	assert.Equal(t, version, extractedVersion, "Invalid version extracted from task JSON")
}

func TestGetVersionEmptyTaskDetail(t *testing.T) {
	tsk := Task{}

	task := Task{}
	taskJSON := marshalTask(t, task)

	_, err := tsk.GetVersion(taskJSON)
	assert.NotNil(t, err, "Expected an error getting task version for an task with no detail")
}

func TestGetVersionEmptyTaskVersion(t *testing.T) {
	tsk := Task{}

	task := Task{
		Detail: &TaskDetail{},
	}
	taskJSON := marshalTask(t, task)

	_, err := tsk.GetVersion(taskJSON)
	assert.NotNil(t, err, "Expected an error getting task version for an task with no version")
}

func TestGerVersionInvalidTask(t *testing.T) {
	tsk := Task{}
	_, err := tsk.GetVersion("invalidTaskJSON")
	assert.NotNil(t, err, "Expected an error getting task version for an invalid task")
}

func marshalTask(t *testing.T, tsk Task) string {
	tskJSON, err := json.Marshal(tsk)
	assert.Nil(t, err, "Unexpected error marshaling task")
	return string(tskJSON)
}
