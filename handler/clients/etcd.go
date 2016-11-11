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

package clients

import (
	"time"

	etcd "github.com/coreos/etcd/clientv3"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

// EtcdInterface defines etcd methods that are used in the project to enable mocking
type EtcdInterface interface {
	// Close shuts down the client's etcd connections.
	Close() error

	// Put puts a key-value pair into etcd.
	// Note that key,value can be plain bytes array and string is
	// an immutable representation of that bytes array.
	// To get a string of bytes, do string([]byte(0x10, 0x20)).
	Put(ctx context.Context, key, val string, opts ...etcd.OpOption) (*etcd.PutResponse, error)

	// Get retrieves keys.
	// By default, Get will return the value for "key", if any.
	// When passed WithRange(end), Get will return the keys in the range [key, end).
	// When passed WithFromKey(), Get returns keys greater than or equal to key.
	// When passed WithRev(rev) with rev > 0, Get retrieves keys at the given revision;
	// if the required revision is compacted, the request will fail with ErrCompacted .
	// When passed WithLimit(limit), the number of returned keys is bounded by limit.
	// When passed WithSort(), the keys will be sorted.
	Get(ctx context.Context, key string, opts ...etcd.OpOption) (*etcd.GetResponse, error)

	// Watch watches on a key or prefix. The watched events will be returned
	// through the returned channel.
	// If the watch is slow or the required rev is compacted, the watch request
	// might be canceled from the server-side and the chan will be closed.
	// 'opts' can be: 'WithRev' and/or 'WithPrefix'.
	Watch(ctx context.Context, key string, opts ...etcd.OpOption) etcd.WatchChan

	// Delete deletes a key, or optionally using WithRange(end), [key, end).
	Delete(ctx context.Context, key string, opts ...etcd.OpOption) (*etcd.DeleteResponse, error)
}

var _ EtcdInterface = (*etcd.Client)(nil)

const (
	dialTimeout = 5 * time.Second
	endpoint    = "localhost:2379"
)

// NewEtcdClient initializes an etcd client
func NewEtcdClient() (*etcd.Client, error) {
	//TODO: attach a lease TTL
	etcd, err := etcd.New(etcd.Config{
		Endpoints:   []string{endpoint},
		DialTimeout: dialTimeout,
	})

	if err != nil {
		return nil, errors.Wrap(err, "Etcd connection error")
	}

	return etcd, nil
}
