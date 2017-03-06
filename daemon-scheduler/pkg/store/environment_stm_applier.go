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
	"github.com/blox/blox/daemon-scheduler/pkg/json"
	"github.com/blox/blox/daemon-scheduler/pkg/types"
	"github.com/coreos/etcd/clientv3/concurrency"
	"github.com/pkg/errors"
)

// EnvironmentSTMApplier supports transactional updates to an environment
type EnvironmentSTMApplier struct {
	// environment key
	key string

	// function that validates the environment currently in the datastore and
	// returns the updated environment to be put into the datastore
	validateAndUpdateEnv func(*types.Environment) (*types.Environment, error)
}

// updateEnvironment performs the following 3 steps as part of a transaction
// 1. Get existing environment from datastore
// 2. Validate environment retrieved from the datastore & update environment
// 3. Put updated environment into datastore
func (envApplier EnvironmentSTMApplier) updateEnvironment(stm concurrency.STM) error {
	err := envApplier.validateApplier()
	if err != nil {
		return err
	}

	existingEnvJSON := stm.Get(envApplier.key)

	existingEnv, err := envApplier.unmarshalEnv(existingEnvJSON)
	if err != nil {
		return err
	}

	updatedEnv, err := envApplier.validateAndUpdateEnv(existingEnv)
	if err != nil {
		return err
	}
	if updatedEnv == nil {
		return errors.New("Environment to be updated cannot be nil")
	}

	updatedEnvJSON, err := envApplier.marshalEnv(*updatedEnv)
	if err != nil {
		return err
	}

	stm.Put(envApplier.key, updatedEnvJSON)
	return nil
}

func (envApplier EnvironmentSTMApplier) unmarshalEnv(envJSON string) (*types.Environment, error) {
	if len(envJSON) == 0 {
		return nil, nil
	}

	var env *types.Environment
	err := json.UnmarshalJSON(envJSON, &env)
	if err != nil {
		return nil, errors.Wrapf(err, "Error unmarshaling environment JSON '%s'", envJSON)
	}

	return env, nil
}

func (envApplier EnvironmentSTMApplier) marshalEnv(env types.Environment) (string, error) {
	envJSON, err := json.MarshalJSON(env)
	if err != nil {
		return "", errors.Wrapf(err, "Error marshaling environment '%v'", env)
	}

	return envJSON, nil
}

func (envApplier EnvironmentSTMApplier) validateApplier() error {
	if envApplier.key == "" {
		return errors.New("Environment key has to be initialized")
	}
	if envApplier.validateAndUpdateEnv == nil {
		return errors.New("Environment validate and update funtion has to be initialized")
	}
	return nil
}
