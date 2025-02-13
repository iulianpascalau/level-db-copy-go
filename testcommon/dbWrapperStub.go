package testcommon

// DBWrapperStub -
type DBWrapperStub struct {
	OpenCalled      func(path string) error
	RangeKeysCalled func(handler func(key []byte, val []byte) bool)
	GetCalled       func(key []byte) ([]byte, error)
	PutCalled       func(key, val []byte) error
	CloseCalled     func() error
}

// Open -
func (stub *DBWrapperStub) Open(path string) error {
	if stub.OpenCalled != nil {
		return stub.OpenCalled(path)
	}

	return nil
}

// RangeKeys -
func (stub *DBWrapperStub) RangeKeys(handler func(key []byte, val []byte) bool) {
	if stub.RangeKeysCalled != nil {
		stub.RangeKeysCalled(handler)
	}
}

// Get -
func (stub *DBWrapperStub) Get(key []byte) ([]byte, error) {
	if stub.GetCalled != nil {
		return stub.GetCalled(key)
	}

	return make([]byte, 0), nil
}

// Put -
func (stub *DBWrapperStub) Put(key, val []byte) error {
	if stub.PutCalled != nil {
		return stub.PutCalled(key, val)
	}

	return nil
}

// Close -
func (stub *DBWrapperStub) Close() error {
	if stub.CloseCalled != nil {
		return stub.CloseCalled()
	}

	return nil
}

// IsInterfaceNil -
func (stub *DBWrapperStub) IsInterfaceNil() bool {
	return stub == nil
}
