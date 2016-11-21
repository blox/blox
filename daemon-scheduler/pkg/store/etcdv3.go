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

	"github.com/blox/blox/daemon-scheduler/pkg/clients"
	etcd "github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
	"github.com/pkg/errors"
)

type etcdDataStore struct {
	etcdInterface clients.EtcdInterface
}

func NewDataStore(etcdInterface clients.EtcdInterface) (DataStore, error) {
	if etcdInterface == nil {
		return nil, errors.Errorf("Invalid etcd input")
	}
	return &etcdDataStore{
		etcdInterface: etcdInterface,
	}, nil
}

// Put puts the provided key-value pair into the datastore
func (datastore etcdDataStore) Put(ctx context.Context, key string, value string) error {
	if len(key) == 0 {
		return errors.Errorf("Key cannot be empty")
	}

	if len(value) == 0 {
		return errors.Errorf("Value cannot be empty")
	}

	_, err := datastore.etcdInterface.Put(ctx, key, value)

	if err != nil {
		return handleEtcdError(err)
	}

	return nil
}

// Get returns a map with one key-value pair where the key matches the provided key
func (datastore etcdDataStore) Get(ctx context.Context, key string) (map[string]string, error) {
	if len(key) == 0 {
		return nil, errors.New("Key cannot be empty")
	}

	resp, err := datastore.etcdInterface.Get(ctx, key)

	if err != nil {
		return nil, handleEtcdError(err)
	}

	kv, err := handleGetResponse(resp)
	return kv, err
}

// GetWithPrefix returns a map of key-value pairs where the key starts with keyPrefix
func (datastore etcdDataStore) GetWithPrefix(ctx context.Context, keyPrefix string) (map[string]string, error) {
	if len(keyPrefix) == 0 {
		return nil, errors.New("Key prefix cannot be empty while getting data from datastore by prefix")
	}

	resp, err := datastore.etcdInterface.Get(ctx, keyPrefix, etcd.WithPrefix())

	if err != nil {
		return nil, handleEtcdError(err)
	}

	kv, err := handleGetResponse(resp)
	return kv, err
}

func handleGetResponse(resp *etcd.GetResponse) (map[string]string, error) {
	kv := make(map[string]string)

	if resp == nil || resp.Kvs == nil {
		return kv, nil
	}

	for _, response := range resp.Kvs {
		kv[string(response.Key)] = string(response.Value)
	}

	return kv, nil
}

// Delete deletes the record with the matching key from the database
func (datastore etcdDataStore) Delete(ctx context.Context, key string) error {
	if len(key) == 0 {
		return errors.New("Key cannot be empty")
	}

	_, err := datastore.etcdInterface.Delete(ctx, key)

	if err != nil {
		return handleEtcdError(err)
	}

	return nil
}

func handleEtcdError(err error) error {
	switch err {
	case context.Canceled:
		return errors.Wrapf(err, "Context is canceled by another routine")
	case context.DeadlineExceeded:
		return errors.Wrapf(err, "Context deadline is exceeded")
	case rpctypes.ErrEmptyKey:
		return errors.Wrapf(err, "Client-side error")
	default:
		return errors.Wrapf(err, "Bad cluster endpoints, which are not etcd servers")
	}
}
