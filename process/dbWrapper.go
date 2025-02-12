package process

import (
	"github.com/multiversx/mx-chain-storage-go/leveldb"
	"github.com/multiversx/mx-chain-storage-go/types"
)

const (
	batchDelaySeconds = 1
	maxBatchSize      = 1000
	maxOpenFiles      = 10
)

type dbWrapper struct {
	db types.Persister
}

// NewDBWrapper creates a new instance of type dbWrapper
func NewDBWrapper(path string) (*dbWrapper, error) {
	lvdb, err := leveldb.NewDB(path, batchDelaySeconds, maxBatchSize, maxOpenFiles)
	if err != nil {
		return nil, err
	}

	return &dbWrapper{
		db: lvdb,
	}, nil
}

// RangeKeys will call the provided handler for each key and value found in the storage
func (wrapper *dbWrapper) RangeKeys(handler func(key []byte, val []byte) bool) {
	wrapper.db.RangeKeys(handler)
}

// Get gets the value associated to the key
func (wrapper *dbWrapper) Get(key []byte) ([]byte, error) {
	return wrapper.db.Get(key)
}

// Put add the value to the (key, val) persistence medium
func (wrapper *dbWrapper) Put(key, val []byte) error {
	return wrapper.db.Put(key, val)
}

// Close closes the files/resources associated to the persistence medium
func (wrapper *dbWrapper) Close() error {
	return wrapper.db.Close()
}

// IsInterfaceNil returns true if there is no value under the interface
func (wrapper *dbWrapper) IsInterfaceNil() bool {
	return wrapper == nil
}
