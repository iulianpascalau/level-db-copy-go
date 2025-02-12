package process

import (
	"os"
	"path"
)

type directoriesHandler struct {
	sourceDirs []string
	destDirs   []string
}

// NewDirectoriesHandler creates a new instance of type directoriesHandler
func NewDirectoriesHandler(sourceParentDir string, destParentDir string) (*directoriesHandler, error) {
	instance := &directoriesHandler{}

	var err error
	instance.sourceDirs, err = getInnerDirectories(sourceParentDir)
	if err != nil {
		return nil, err
	}

	instance.destDirs, err = getInnerDirectories(destParentDir)
	if err != nil {
		return nil, err
	}

	return instance, nil
}

func getInnerDirectories(parentDir string) ([]string, error) {
	dirInfo, err := os.ReadDir(parentDir)
	if err != nil {
		return nil, err
	}

	result := make([]string, 0, 1024)
	for _, directoryInfo := range dirInfo {
		if !directoryInfo.IsDir() {
			continue
		}

		result = append(result, path.Join(parentDir, directoryInfo.Name()))
	}

	return result, nil
}

// SourceDirectories returns the source directories
func (handler *directoriesHandler) SourceDirectories() []string {
	return handler.sourceDirs
}

// DestinationDirectories returns the destination directories
func (handler *directoriesHandler) DestinationDirectories() []string {
	return handler.destDirs
}

// IsInterfaceNil returns true if there is no value under the interface
func (handler *directoriesHandler) IsInterfaceNil() bool {
	return handler == nil
}
