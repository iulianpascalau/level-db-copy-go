package process

import (
	"fmt"
	"testing"

	"iulianpascalau/level-db-copy-go/testcommon"

	"github.com/stretchr/testify/assert"
)

type recorder struct {
	srcOpenedDBs  []string
	srcClosedDBs  []string
	destOpenedDBs []string
	destClosedDBs []string
	putOps        map[string]string
}

type testHandler struct {
	getOps map[string]string
}

func TestNewDataCopyHandler(t *testing.T) {
	t.Parallel()

	t.Run("nil directories handler should error", func(t *testing.T) {
		t.Parallel()

		handler, err := NewDataCopyHandler(
			nil,
			&testcommon.DBWrapperStub{},
			&testcommon.DBWrapperStub{},
		)

		assert.Nil(t, handler)
		assert.Equal(t, errNilDirectoriesHandler, err)
	})
	t.Run("nil source DB wrapper should error", func(t *testing.T) {
		t.Parallel()

		handler, err := NewDataCopyHandler(
			&testcommon.DirectoriesHandlerStub{},
			nil,
			&testcommon.DBWrapperStub{},
		)

		assert.Nil(t, handler)
		assert.ErrorIs(t, err, errNilDBWrapper)
		assert.Contains(t, err.Error(), "for the source DB wrapper")
	})
	t.Run("nil destination DB wrapper should error", func(t *testing.T) {
		t.Parallel()

		handler, err := NewDataCopyHandler(
			&testcommon.DirectoriesHandlerStub{},
			&testcommon.DBWrapperStub{},
			nil,
		)

		assert.Nil(t, handler)
		assert.ErrorIs(t, err, errNilDBWrapper)
		assert.Contains(t, err.Error(), "for the destination DB wrapper")
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		handler, err := NewDataCopyHandler(
			&testcommon.DirectoriesHandlerStub{},
			&testcommon.DBWrapperStub{},
			&testcommon.DBWrapperStub{},
		)

		assert.NotNil(t, handler)
		assert.Nil(t, err)
	})
}

func TestDataCopyHandler_Process(t *testing.T) {
	t.Parallel()

	t.Run("should work", func(t *testing.T) {
		test := &testHandler{
			getOps: map[string]string{
				"A-key-0": "dest",
				"A-key-1": "dest",
				"A-key-4": "dest",

				"B-key-1": "dest",
				"B-key-2": "dest",
				"B-key-3": "dest",
			},
		}

		rec := &recorder{
			putOps: make(map[string]string),
		}

		expectedDBOperationOrder := []string{"A", "B"}
		expectedPutOperations := map[string]string{
			"A-key-2": "A-val-s-2",
			"A-key-3": "A-val-s-3",
			"B-key-0": "B-val-s-0",
			"B-key-4": "B-val-s-4",
		}

		handler, _ := NewDataCopyHandler(setupForProcess(t, test, rec))
		err := handler.Process()
		assert.Nil(t, err)

		assert.ElementsMatch(t, expectedDBOperationOrder, rec.srcOpenedDBs)
		assert.ElementsMatch(t, expectedDBOperationOrder, rec.srcClosedDBs)
		assert.ElementsMatch(t, expectedDBOperationOrder, rec.destOpenedDBs)
		assert.ElementsMatch(t, expectedDBOperationOrder, rec.destClosedDBs)
		assert.Equal(t, expectedPutOperations, rec.putOps)
	})
}

func setupForProcess(t *testing.T, test *testHandler, recorder *recorder) (DirectoriesHandler, DBWrapper, DBWrapper) {
	directoriesHandlerInstance := &testcommon.DirectoriesHandlerStub{
		SourceDirectoriesCalled: func() []string {
			return []string{"A", "B", "C"}
		},
		DestinationDirectoriesCalled: func() []string {
			return []string{"A", "B", "D"}
		},
	}

	var currentSrcDB []byte
	srcDbWrapper := &testcommon.DBWrapperStub{
		OpenCalled: func(path string) error {
			currentSrcDB = []byte(path)
			recorder.srcOpenedDBs = append(recorder.srcOpenedDBs, path)

			return nil
		},
		RangeKeysCalled: func(handler func(key []byte, val []byte) bool) {
			for i := 0; i < 5; i++ {
				handler(
					[]byte(fmt.Sprintf("%s-key-%d", string(currentSrcDB), i)),
					[]byte(fmt.Sprintf("%s-val-s-%d", string(currentSrcDB), i)),
				)
			}
		},
		GetCalled: func(key []byte) ([]byte, error) {
			assert.Fail(t, "should have not called Get on the src DB wrapper")
			return nil, nil
		},
		PutCalled: func(key, val []byte) error {
			assert.Fail(t, "should have not called Put on the src DB wrapper")
			return nil
		},
		CloseCalled: func() error {
			recorder.srcClosedDBs = append(recorder.srcClosedDBs, string(currentSrcDB))
			currentSrcDB = nil

			return nil
		},
	}

	var currentDestDB []byte
	destDbWrapper := &testcommon.DBWrapperStub{
		OpenCalled: func(path string) error {
			currentDestDB = []byte(path)
			recorder.destOpenedDBs = append(recorder.destOpenedDBs, path)

			return nil
		},
		RangeKeysCalled: func(handler func(key []byte, val []byte) bool) {
			assert.Fail(t, "should have not called RangeKeysCalled on the dest DB wrapper")
		},
		GetCalled: func(key []byte) ([]byte, error) {
			val, found := test.getOps[string(key)]
			if !found {
				return nil, fmt.Errorf("not found")
			}

			return []byte(val), nil
		},
		PutCalled: func(key, val []byte) error {
			recorder.putOps[string(key)] = string(val)
			return nil
		},
		CloseCalled: func() error {
			recorder.destClosedDBs = append(recorder.destClosedDBs, string(currentDestDB))
			currentDestDB = nil

			return nil
		},
	}

	return directoriesHandlerInstance, srcDbWrapper, destDbWrapper
}
