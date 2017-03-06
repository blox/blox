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

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	"github.com/pkg/errors"
)

// EtcdTXStore defines methods to support etcd's STM
type EtcdTXStore interface {
	GetV3Client() *clientv3.Client
	NewSTMRepeatable(context.Context, *clientv3.Client, func(concurrency.STM) error) (*clientv3.TxnResponse, error)
}

type etcdTransactionalStore struct {
	v3Client *clientv3.Client
}

// NewEtcdTXStore initializs the etcdTransactionalStore struct
func NewEtcdTXStore(v3Client *clientv3.Client) (EtcdTXStore, error) {
	if v3Client == nil {
		return nil, errors.Errorf("Etcd client in not initialized")
	}
	return &etcdTransactionalStore{
		v3Client: v3Client,
	}, nil
}

func (txStore etcdTransactionalStore) NewSTMRepeatable(ctx context.Context, v3Client *clientv3.Client, apply func(concurrency.STM) error) (*clientv3.TxnResponse, error) {
	return concurrency.NewSTMRepeatable(ctx, v3Client, apply)
}

func (txStore etcdTransactionalStore) GetV3Client() *clientv3.Client {
	return txStore.v3Client
}
