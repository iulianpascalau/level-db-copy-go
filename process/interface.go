package process

// DBWrapper defines the operations supported by a database wrapper
type DBWrapper interface {
	Open(path string) error
	RangeKeys(handler func(key []byte, val []byte) bool)
	Get(key []byte) ([]byte, error)
	Put(key, val []byte) error
	Close() error
	IsInterfaceNil() bool
}

// DirectoriesHandler defines the operations supported by a directories handler
type DirectoriesHandler interface {
	SourceDirectories() []string
	DestinationDirectories() []string
	IsInterfaceNil() bool
}
