package process

import (
	"sync"

	"github.com/multiversx/mx-chain-storage-go/leveldb"
	"github.com/multiversx/mx-chain-storage-go/types"
)

const (
	batchDelaySeconds = 1
	maxBatchSize      = 1000
	maxOpenFiles      = 10
)

type dbWrapper struct {
	mutDB sync.RWMutex
	db    types.Persister
}

// NewDBWrapper creates a new instance of type dbWrapper
func NewDBWrapper() *dbWrapper {
	return &dbWrapper{}
}

// Open will attempt to open the level DB from the provided path
// Errors if the inner DB is still opened
func (wrapper *dbWrapper) Open(path string) error {
	wrapper.mutDB.Lock()
	defer wrapper.mutDB.Unlock()

	if wrapper.db != nil {
		return errInnerDBIsNotClosed
	}

	lvdb, err := leveldb.NewDB(path, batchDelaySeconds, maxBatchSize, maxOpenFiles)
	if err != nil {
		return err
	}

	wrapper.db = lvdb

	return nil
}

// RangeKeys will call the provided handler for each key and value found in the storage
func (wrapper *dbWrapper) RangeKeys(handler func(key []byte, val []byte) bool) {
	wrapper.mutDB.RLock()
	defer wrapper.mutDB.RUnlock()

	if wrapper.db == nil {
		return
	}

	wrapper.db.RangeKeys(handler)
}

// Get gets the value associated to the key
func (wrapper *dbWrapper) Get(key []byte) ([]byte, error) {
	wrapper.mutDB.RLock()
	defer wrapper.mutDB.RUnlock()

	if wrapper.db == nil {
		return nil, errInnerDBIsNotOpened
	}

	return wrapper.db.Get(key)
}

// Put add the value to the (key, val) persistence medium
func (wrapper *dbWrapper) Put(key, val []byte) error {
	wrapper.mutDB.RLock()
	defer wrapper.mutDB.RUnlock()

	if wrapper.db == nil {
		return errInnerDBIsNotOpened
	}

	return wrapper.db.Put(key, val)
}

// Close closes the files/resources associated to the persistence medium
func (wrapper *dbWrapper) Close() error {
	wrapper.mutDB.Lock()
	defer wrapper.mutDB.Unlock()

	err := wrapper.db.Close()
	wrapper.db = nil

	return err
}

// IsInterfaceNil returns true if there is no value under the interface
func (wrapper *dbWrapper) IsInterfaceNil() bool {
	return wrapper == nil
}
