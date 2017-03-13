// Copyright 2016-2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.
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
	"errors"
	"testing"

	"github.com/blox/blox/daemon-scheduler/pkg/environment/types"
	"github.com/blox/blox/daemon-scheduler/pkg/json"
	"github.com/blox/blox/daemon-scheduler/pkg/mocks"
	"github.com/coreos/etcd/clientv3"
	"github.com/stretchr/testify/assert"
)

func TestValidateApplierNoEnvironmentKey(t *testing.T) {
	applier := &EnvironmentSTMApplier{}

	err := applier.validateApplier()
	assert.Error(t, err, "Expected error when record key is not set in applier")
}

func TestValidateApplierNoValidateAndUpdateFunction(t *testing.T) {
	applier := &EnvironmentSTMApplier{
		key: "someKey",
	}

	err := applier.validateApplier()
	assert.Error(t, err, "Expected error when record validateAndUpdateFunction is not set in applier")
}

func TestUpdateEnvironment(t *testing.T) {
	key := "key"

	getEnv := environment(t)
	getEnvJSON := environmentJSON(t, getEnv)

	putEnv := environment(t)
	putEnvJSON := environmentJSON(t, putEnv)

	mockSTM := &mocks.MockSTM{
		GetFunc: func(k string) string {
			assert.Equal(t, key, k, "Unexpected key for Get")
			return getEnvJSON
		},
		PutFunc: func(k string, val string, opts ...clientv3.OpOption) {
			assert.Equal(t, key, k, "Unexpected key in Put")
			assert.Equal(t, val, putEnvJSON, "Unexpected record in Put")
		},
	}

	applier := &EnvironmentSTMApplier{
		key: key,
		validateAndUpdateEnv: func(existingEnv *types.Environment) (*types.Environment, error) {
			return putEnv, nil
		},
	}

	err := applier.updateEnvironment(mockSTM)
	assert.NoError(t, err, "Unexpected error updating environment")
}

func TestUpdateEnvironmentGetReturnsInvalidEnvironment(t *testing.T) {
	key := "key"

	getEnvJSON := "invalidEnvironment"

	mockSTM := &mocks.MockSTM{
		GetFunc: func(k string) string {
			assert.Equal(t, key, k, "Unexpected key for Get")
			return getEnvJSON
		},
	}

	applier := &EnvironmentSTMApplier{
		key: key,
		validateAndUpdateEnv: func(existingEnv *types.Environment) (*types.Environment, error) {
			return nil, nil
		},
	}

	err := applier.updateEnvironment(mockSTM)
	assert.NotNil(t, err, "Expected error when get returns invalid environment")
}

func TestUpdateEnvironmentValidateFunctionReturnsError(t *testing.T) {
	key := "key"

	getEnv := environment(t)
	getEnvJSON := environmentJSON(t, getEnv)

	mockSTM := &mocks.MockSTM{
		GetFunc: func(k string) string {
			assert.Equal(t, key, k, "Unexpected key for Get")
			return getEnvJSON
		},
	}

	applier := &EnvironmentSTMApplier{
		key: key,
		validateAndUpdateEnv: func(existingEnv *types.Environment) (*types.Environment, error) {
			return nil, errors.New("Error validating")
		},
	}

	err := applier.updateEnvironment(mockSTM)
	assert.NotNil(t, err, "Expected error updating environment when validate function returns error")
}

func TestUpdateEnvironmentValidateFunctionReturnsNilEnvironment(t *testing.T) {
	key := "key"

	getEnv := environment(t)
	getEnvJSON := environmentJSON(t, getEnv)

	mockSTM := &mocks.MockSTM{
		GetFunc: func(k string) string {
			assert.Equal(t, key, k, "Unexpected key for Get")
			return getEnvJSON
		},
	}

	applier := &EnvironmentSTMApplier{
		key: key,
		validateAndUpdateEnv: func(existingEnv *types.Environment) (*types.Environment, error) {
			return nil, nil
		},
	}

	err := applier.updateEnvironment(mockSTM)
	assert.NotNil(t, err, "Expected error updating environment when validate function returns invalid environment")
}

func environment(t *testing.T) *types.Environment {
	env, err := types.NewEnvironment("envName", "td", "clusterName")
	assert.Nil(t, err, "Unexpected error creating an environment")
	return env
}

func environmentJSON(t *testing.T, env *types.Environment) string {
	envJSON, err := json.MarshalJSON(env)
	assert.Nil(t, err, "Unexpected error marshaling environment")
	return envJSON
}
