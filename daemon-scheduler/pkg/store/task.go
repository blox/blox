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
	"encoding/json"

	"github.com/blox/blox/daemon-scheduler/pkg/types"
	"github.com/pkg/errors"
)

type TaskStore interface {
	PutTask(ctx context.Context, envName string, task types.Task) error
}

type taskStore struct {
	datastore DataStore
}

func NewTaskStore(ds DataStore) (TaskStore, error) {
	if ds == nil {
		return nil, errors.New("Datastore cannot be empty while initializing taskstore")
	}

	return taskStore{
		datastore: ds,
	}, nil
}

func (ts taskStore) PutTask(ctx context.Context, envName string, task types.Task) error {
	key, err := GenerateTaskKey(envName, task.TaskARN)
	if err != nil {
		return err
	}

	taskJSON, err := json.Marshal(task)
	if err != nil {
		return errors.Wrap(err, "Error while mashaling task")
	}

	err = ts.datastore.Put(ctx, key, string(taskJSON))
	if err != nil {
		return errors.Wrapf(err, "Error while adding task '%s' to the store", taskJSON)
	}

	return nil
}
