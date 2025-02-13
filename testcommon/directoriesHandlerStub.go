package testcommon

// DirectoriesHandlerStub -
type DirectoriesHandlerStub struct {
	SourceDirectoriesCalled      func() []string
	DestinationDirectoriesCalled func() []string
}

// SourceDirectories -
func (stub *DirectoriesHandlerStub) SourceDirectories() []string {
	if stub.SourceDirectoriesCalled != nil {
		return stub.SourceDirectoriesCalled()
	}

	return make([]string, 0)
}

// DestinationDirectories -
func (stub *DirectoriesHandlerStub) DestinationDirectories() []string {
	if stub.DestinationDirectoriesCalled != nil {
		return stub.DestinationDirectoriesCalled()
	}

	return make([]string, 0)
}

// IsInterfaceNil -
func (stub *DirectoriesHandlerStub) IsInterfaceNil() bool {
	return stub == nil
}
