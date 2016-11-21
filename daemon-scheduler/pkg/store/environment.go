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
	"context"

	"github.com/blox/blox/daemon-scheduler/pkg/json"
	"github.com/blox/blox/daemon-scheduler/pkg/types"
	"github.com/pkg/errors"
)

type EnvironmentStore interface {
	PutEnvironment(ctx context.Context, environment types.Environment) error
	GetEnvironment(ctx context.Context, name string) (*types.Environment, error)
	DeleteEnvironment(ctx context.Context, environment types.Environment) error
	ListEnvironments(ctx context.Context) ([]types.Environment, error)
}

type environmentStore struct {
	datastore DataStore
}

func NewEnvironmentStore(ds DataStore) (EnvironmentStore, error) {
	if ds == nil {
		return nil, errors.New("The datastore cannot be nil")
	}

	return environmentStore{
		datastore: ds,
	}, nil
}

func (e environmentStore) PutEnvironment(ctx context.Context, environment types.Environment) error {
	key, err := GenerateEnvironmentKey(environment.Name)
	if err != nil {
		return err
	}

	dataJSON, err := json.MarshalJSON(environment)
	if err != nil {
		return err
	}

	err = e.datastore.Put(ctx, key, dataJSON)
	if err != nil {
		return err
	}

	return nil
}

func (e environmentStore) GetEnvironment(ctx context.Context, name string) (*types.Environment, error) {
	if len(name) == 0 {
		return nil, errors.New("Environment name is missing")
	}

	key, err := GenerateEnvironmentKey(name)
	if err != nil {
		return nil, err
	}

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

	var environment types.Environment
	for _, v := range resp {
		err = json.UnmarshalJSON(v, &environment)
		if err != nil {
			return nil, err
		}
		break
	}

	return &environment, nil
}

func (e environmentStore) DeleteEnvironment(ctx context.Context, environment types.Environment) error {
	key, err := GenerateEnvironmentKey(environment.Name)
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
