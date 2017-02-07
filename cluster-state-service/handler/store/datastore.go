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
	"time"

	"github.com/blox/blox/cluster-state-service/handler/clients"
	"github.com/blox/blox/cluster-state-service/handler/regex"
	storetypes "github.com/blox/blox/cluster-state-service/handler/store/types"
	"github.com/blox/blox/cluster-state-service/handler/types"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
	"github.com/pkg/errors"
	"strconv"
)

const (
	// requestTimeout is timeout set when calling etcd APIs.
	// This timeout is set to 1 minute to support list APIs
	// with prefix match
	requestTimeout    = 1 * time.Minute
	streamIdleTimeout = 1 * time.Hour
)

// DataStore defines methods to access the database
type DataStore interface {
	GetWithPrefix(keyPrefix string) (map[string]storetypes.Entity, error)
	Get(key string) (map[string]storetypes.Entity, error)
	Add(key string, value string) error
	StreamWithPrefix(ctx context.Context, keyPrefix string, entityVersion string) (chan map[string]storetypes.Entity, error)
	Delete(key string) (int64, error)
}

type etcdDataStore struct {
	etcdInterface clients.EtcdInterface
}

// NewDataStore initializes the etcdDataStore struct
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
func (datastore etcdDataStore) GetWithPrefix(keyPrefix string) (map[string]storetypes.Entity, error) {
	if len(keyPrefix) == 0 {
		return nil, errors.New("Key prefix cannot be empty while getting data from datastore by prefix")
	}

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	resp, err := datastore.etcdInterface.Get(ctx, keyPrefix, clientv3.WithPrefix())
	defer cancel()

	if err != nil {
		return nil, handleEtcdError(err)
	}

	return handleGetResponse(resp), nil
}

// Get returns a map with one key-value pair where the key matches the provided key
func (datastore etcdDataStore) Get(key string) (map[string]storetypes.Entity, error) {
	if len(key) == 0 {
		return nil, errors.New("Key cannot be empty while getting data from datastore by key")
	}

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	resp, err := datastore.etcdInterface.Get(ctx, key)
	defer cancel()

	if err != nil {
		return nil, handleEtcdError(err)
	}

	return handleGetResponse(resp), nil
}

// StreamWithPrefix starts a go routine that streams key-value pairs whose keys start with keyPrefix into the channel returned
func (datastore etcdDataStore) StreamWithPrefix(ctx context.Context, keyPrefix string, entityVersion string) (chan map[string]storetypes.Entity, error) {
	if len(keyPrefix) == 0 {
		return nil, errors.New("Key prefix cannot be empty while streaming data from datastore by prefix")
	}

	// If entity version is specified, verify that it is not out of range.
	// This logic is here because the Watch channel does not throw these rpctypes errors, and blocks if you specify a revision in the future until it gets to that revision.
	// There is a small chance that Etcd could be compacted between this check and when the stream initializes, in which case the channel would close without any context as to why.
	// TODO: Look into a better way of handling this check in the datastore.stream method.
	if entityVersion != "" {
		revision, err := regex.GetEntityVersion(entityVersion)
		if err != nil {
			return nil, err
		}

		if _, err = datastore.etcdInterface.Get(ctx, keyPrefix, clientv3.WithRev(revision)); err != nil {
			if err == rpctypes.ErrCompacted || err == rpctypes.ErrFutureRev {
				return nil, types.NewOutOfRangeEntityVersion(err)
			}
			return nil, err
		}
	}

	kvChan := make(chan map[string]storetypes.Entity) // go routine datastore.stream() handles closing of this channel
	go datastore.stream(ctx, keyPrefix, entityVersion, kvChan)
	return kvChan, nil
}

// Delete returns a map with one key-value pair where the key matches the provided key
func (datastore etcdDataStore) Delete(key string) (int64, error) {
	if len(key) == 0 {
		return 0, errors.New("Key cannot be empty while deleting data from datastore by key")
	}

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	resp, err := datastore.etcdInterface.Delete(ctx, key)
	defer cancel()

	if err != nil {
		return 0, handleEtcdError(err)
	}

	if resp == nil {
		return 0, nil
	}
	return resp.Deleted, nil
}

func (datastore etcdDataStore) stream(ctx context.Context, keyPrefix string, entityVersion string, kvChan chan map[string]storetypes.Entity) {
	defer close(kvChan)

	etcdCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Revision will default to '0', meaning stream from now.
	var revision int64
	var err error
	if entityVersion != "" {
		if revision, err = regex.GetEntityVersion(entityVersion); err != nil {
			return
		}
	}

	watchChan := datastore.etcdInterface.Watch(etcdCtx, keyPrefix, clientv3.WithPrefix(), clientv3.WithRev(revision))
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
				// Skip empty events, such as Etcd deletes.
				// TODO: Look into whether we should return something here.
				if len(ev.Kv.Value) == 0 {
					continue
				}

				entity := storetypes.Entity{
					Key: string(ev.Kv.Key),
					Value: string(ev.Kv.Value),
					Version: strconv.FormatInt(ev.Kv.ModRevision, 10),
				}
				kv := map[string]storetypes.Entity{string(ev.Kv.Key): entity}
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

func handleGetResponse(resp *clientv3.GetResponse) map[string]storetypes.Entity {
	kv := make(map[string]storetypes.Entity)

	if resp == nil || resp.Kvs == nil {
		return kv
	}

	// response.Key = The object's key in Etcd.
	// response.Value = The object's value in Etcd.
	// response.ModRevision = The object's last modification identifier in Etcd (incrementing integer for every change in Etcd).
	for _, response := range resp.Kvs {
		entity := storetypes.Entity{
			Key: string(response.Key),
			Value: string(response.Value),
			Version: strconv.FormatInt(response.ModRevision, 10),
		}
		kv[string(response.Key)] = entity
	}

	return kv
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
