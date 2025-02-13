package process

import (
	"fmt"
	"path"
	"strings"
	"sync"

	"github.com/multiversx/mx-chain-core-go/core/check"
	logger "github.com/multiversx/mx-chain-logger-go"
)

var log = logger.GetOrCreate("process")

type paths struct {
	src  string
	dest string
}

type dataCopyHandler struct {
	mutCriticalArea    sync.Mutex
	directoriesHandler DirectoriesHandler
	srcDBWrapper       DBWrapper
	destDBWrapper      DBWrapper
}

// NewDataCopyHandler creates a new instance of type data copy handler
func NewDataCopyHandler(
	directoriesHandler DirectoriesHandler,
	srcDBWrapper DBWrapper,
	destDBWrapper DBWrapper,
) (*dataCopyHandler, error) {
	if check.IfNil(directoriesHandler) {
		return nil, errNilDirectoriesHandler
	}
	if check.IfNil(srcDBWrapper) {
		return nil, fmt.Errorf("%w for the source DB wrapper", errNilDBWrapper)
	}
	if check.IfNil(destDBWrapper) {
		return nil, fmt.Errorf("%w for the destination DB wrapper", errNilDBWrapper)
	}

	return &dataCopyHandler{
		directoriesHandler: directoriesHandler,
		srcDBWrapper:       srcDBWrapper,
		destDBWrapper:      destDBWrapper,
	}, nil
}

// Process will attempt to complete the DB copy process
func (handler *dataCopyHandler) Process() error {
	handler.mutCriticalArea.Lock()
	defer handler.mutCriticalArea.Unlock()

	commonDirs, names := handler.computeCommonDirs()
	log.Info("Common directories between the source and destination parent paths", "sub-directories", names)

	counter := 1
	for name, pathInfo := range commonDirs {
		log.Info("now processing sub-directory", "name", name, "overall progress", fmt.Sprintf("%d/%d", counter, len(commonDirs)))

		numInserts, err := handler.processDB(pathInfo)
		if err != nil {
			return err
		}

		log.Info("successfully processed DB", "name", name, "missing info added", numInserts)
		counter++
	}

	return nil
}

func (handler *dataCopyHandler) computeCommonDirs() (map[string]paths, string) {
	srcDirs := convertDirStrings(handler.directoriesHandler.SourceDirectories())
	destDirs := convertDirStrings(handler.directoriesHandler.DestinationDirectories())

	commonDirs := make(map[string]paths, len(srcDirs)+len(destDirs))
	names := make([]string, 0, len(srcDirs)+len(destDirs))
	for name, srcFullPath := range srcDirs {
		destFullPath, found := destDirs[name]
		if found {
			commonDirs[name] = paths{
				src:  srcFullPath,
				dest: destFullPath,
			}
			names = append(names, name)
		}
	}

	return commonDirs, strings.Join(names, ", ")
}

func convertDirStrings(dirStrings []string) map[string]string {
	mapDirs := make(map[string]string, len(dirStrings))
	for _, dir := range dirStrings {
		_, lastDirElement := path.Split(dir)
		mapDirs[lastDirElement] = dir
	}

	return mapDirs
}

func (handler *dataCopyHandler) processDB(pathInfo paths) (int, error) {
	err := handler.srcDBWrapper.Open(pathInfo.src)
	if err != nil {
		return 0, err
	}

	err = handler.destDBWrapper.Open(pathInfo.dest)
	if err != nil {
		return 0, err
	}

	numInserts := 0
	handlerFunc := func(key []byte, val []byte) bool {
		existingValue, _ := handler.destDBWrapper.Get(key)
		if existingValue == nil {
			err = handler.destDBWrapper.Put(key, val)
			if err != nil {
				log.Error("error encountered while processing a DB put operation",
					"dest path", pathInfo.dest, "key", key)
			} else {
				numInserts++
			}
		}

		return true
	}

	handler.srcDBWrapper.RangeKeys(handlerFunc)

	errClose1 := handler.srcDBWrapper.Close()
	errClose2 := handler.destDBWrapper.Close()

	if errClose1 != nil {
		return numInserts, errClose1
	}

	return numInserts, errClose2
}
