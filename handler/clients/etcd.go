package clients

import (
	etcd "github.com/coreos/etcd/clientv3"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"time"
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
