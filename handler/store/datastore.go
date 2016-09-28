package store

import (
	"github.com/aws/amazon-ecs-event-stream-handler/handler/clients"
	etcd "github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"time"
)

const (
	requestTimeout = 5 * time.Second
)

// DataStore defines methods to access the database
type DataStore interface {
	GetWithPrefix(keyPrefix string) (map[string]string, error)
	Get(key string) (map[string]string, error)
	Add(key string, value string) error
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
		return errors.Errorf("Key cannot be empty")
	}

	if len(value) == 0 {
		return errors.Errorf("Value cannot be empty")
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
		return nil, errors.New("keyPrefix cannot be empty")
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
		return nil, errors.New("Key cannot be empty")
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
