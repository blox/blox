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
	"context"

	"github.com/blox/blox/daemon-scheduler/pkg/json"
	storetypes "github.com/blox/blox/daemon-scheduler/pkg/store/types"
	"github.com/blox/blox/daemon-scheduler/pkg/types"
	"github.com/pkg/errors"
)

const (
	environmentKeyPrefix = "ecs/environment/"
)

// EnvironmentStore defines methods to handle interations with the datastore related to environments
type EnvironmentStore interface {
	// PutEnvironment performs a transactional put. It retrieves the environment using the 'name', validates the environment
	// based on the implementation of 'validateAndUpdateEnv' and updates the environment with the environment
	// returned by 'validateAndUpdateEnv' all within a transaction.
	PutEnvironment(ctx context.Context, name string, validateAndUpdateEnv storetypes.ValidateAndUpdateEnvironment) error

	// GetEnvironment retrieves an enrironment by name
	GetEnvironment(ctx context.Context, name string) (*types.Environment, error)

	// DeleteEnvironment 'deletes' an environment by name
	DeleteEnvironment(ctx context.Context, name string) error

	// ListEnvironments lists all environments
	ListEnvironments(ctx context.Context) ([]types.Environment, error)
}

type environmentStore struct {
	datastore   DataStore
	etcdTXStore EtcdTXStore
}

func NewEnvironmentStore(ds DataStore, ts EtcdTXStore) (EnvironmentStore, error) {
	if ds == nil {
		return nil, errors.Errorf("Datastore is not initialized")
	}
	if ts == nil {
		return nil, errors.Errorf("Etcd transactional store is not initialized")
	}

	return environmentStore{
		datastore:   ds,
		etcdTXStore: ts,
	}, nil
}

func generateEnvironmentKey(name string) (string, error) {
	if name == "" {
		return "", errors.New("Environment name is missing")
	}
	return environmentKeyPrefix + name, nil
}

func (e environmentStore) PutEnvironment(ctx context.Context, name string,
	validateAndUpdateEnv storetypes.ValidateAndUpdateEnvironment) error {
	key, err := generateEnvironmentKey(name)
	if err != nil {
		return err
	}

	applier := &EnvironmentSTMApplier{
		key:                  key,
		validateAndUpdateEnv: validateAndUpdateEnv,
	}

	_, err = e.etcdTXStore.NewSTMRepeatable(ctx, e.etcdTXStore.GetV3Client(), applier.updateEnvironment)

	return err
}

func (e environmentStore) GetEnvironment(ctx context.Context, name string) (*types.Environment, error) {
	key, err := generateEnvironmentKey(name)
	if err != nil {
		return nil, err
	}

	var environment types.Environment

	resp, err := e.datastore.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	if len(resp) == 0 {
		return nil, nil
	}

	if len(resp) > 1 {
		return nil, errors.Errorf("Multiple entries exist with the key %v", key)
	}

	for _, v := range resp {
		err = json.UnmarshalJSON(v, &environment)
		if err != nil {
			return nil, err
		}
		break
	}

	return &environment, nil
}

func (e environmentStore) DeleteEnvironment(ctx context.Context, name string) error {
	key, err := generateEnvironmentKey(name)
	if err != nil {
		return err
	}

	err = e.datastore.Delete(ctx, key)
	if err != nil {
		return err
	}

	return nil
}

func (e environmentStore) ListEnvironments(ctx context.Context) ([]types.Environment, error) {
	resp, err := e.datastore.GetWithPrefix(ctx, environmentKeyPrefix)
	if err != nil {
		return nil, err
	}

	environments := make([]types.Environment, 0, len(resp))

	for _, v := range resp {
		environment := types.Environment{}

		err = json.UnmarshalJSON(v, &environment)
		if err != nil {
			return nil, err
		}

		environments = append(environments, environment)
	}

	return environments, nil
}
