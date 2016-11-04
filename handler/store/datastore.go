// Copyright 2016 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the License). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the license file accompanying this file. This file is distributed
// on an AS IS BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package store

import (
	"context"
	"time"

	"github.com/aws/amazon-ecs-event-stream-handler/handler/clients"
	etcd "github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
	"github.com/pkg/errors"
)

const (
	requestTimeout    = 5 * time.Second
	streamIdleTimeout = 1 * time.Hour
)

// DataStore defines methods to access the database
type DataStore interface {
	GetWithPrefix(keyPrefix string) (map[string]string, error)
	Get(key string) (map[string]string, error)
	Add(key string, value string) error
	StreamWithPrefix(ctx context.Context, keyPrefix string) (chan map[string]string, error)
}

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

// Add adds the provided key-value pair to the datastore
func (datastore etcdDataStore) Add(key string, value string) error {
	if len(key) == 0 {
		return errors.Errorf("Key cannot be empty while adding data into datastore")
	}

	if len(value) == 0 {
		return errors.Errorf("Value cannot be empty while adding data into datastore")
	}

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	_, err := datastore.etcdInterface.Put(ctx, key, value)
	defer cancel()

	if err != nil {
		return handleEtcdError(err)
	}

	return nil
}

// GetWithPrefix returns a map of key-value pairs where the key starts with keyPrefix
func (datastore etcdDataStore) GetWithPrefix(keyPrefix string) (map[string]string, error) {
	if len(keyPrefix) == 0 {
		return nil, errors.New("Key prefix cannot be empty while getting data from datastore by prefix")
	}

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	resp, err := datastore.etcdInterface.Get(ctx, keyPrefix, etcd.WithPrefix())
	defer cancel()

	if err != nil {
		return nil, handleEtcdError(err)
	}

	kv, err := handleGetResponse(resp)
	return kv, err
}

// Get returns a map with one key-value pair where the key matches the provided key
func (datastore etcdDataStore) Get(key string) (map[string]string, error) {
	if len(key) == 0 {
		return nil, errors.New("Key cannot be empty while getting data from datastore by key")
	}

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	resp, err := datastore.etcdInterface.Get(ctx, key)
	defer cancel()

	if err != nil {
		return nil, handleEtcdError(err)
	}

	kv, err := handleGetResponse(resp)
	return kv, err
}

// StreamWithPrefix starts a go routine that streams key-value pairs whose keys start with keyPrefix into the channel returned
func (datastore etcdDataStore) StreamWithPrefix(ctx context.Context, keyPrefix string) (chan map[string]string, error) {
	if len(keyPrefix) == 0 {
		return nil, errors.New("Key prefix cannot be empty while streaming data from datastore by prefix")
	}

	kvChan := make(chan map[string]string) // go routine datastore.stream() handles closing of this channel
	go datastore.stream(ctx, keyPrefix, kvChan)
	return kvChan, nil
}

func (datastore etcdDataStore) stream(ctx context.Context, keyPrefix string, kvChan chan map[string]string) {
	defer close(kvChan)

	etcdCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	watchChan := datastore.etcdInterface.Watch(etcdCtx, keyPrefix, etcd.WithPrefix())
	streamIdleTimer := time.NewTimer(streamIdleTimeout)
	defer streamIdleTimer.Stop()

	for {
		select {
		case event, ok := <-watchChan:
			if !ok {
				return
			}
			resetStreamIdleTimer(streamIdleTimer)
			for _, ev := range event.Events {
				kv := map[string]string{string(ev.Kv.Key): string(ev.Kv.Value)}
				kvChan <- kv
			}

		// TODO: Verify if this is needed or if we should we allow for infinite streaming even if the stream in idle
		case <-streamIdleTimer.C:
			return

		case <-etcdCtx.Done():
			return
		}
	}
}

func resetStreamIdleTimer(t *time.Timer) {
	if !t.Stop() {
		<-t.C
	}
	t.Reset(streamIdleTimeout)
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
