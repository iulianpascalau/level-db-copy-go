package integrationTests

import (
	"os"
	"path"
	"testing"

	"iulianpascalau/level-db-copy-go/process"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDBCopy(t *testing.T) {
	srcParentDir, destParentDir := setupDirs(t)

	dirHandler, err := process.NewDirectoriesHandler(srcParentDir, destParentDir)
	assert.Nil(t, err)

	copyHandler, err := process.NewDataCopyHandler(
		dirHandler,
		process.NewDBWrapper(),
		process.NewDBWrapper(),
	)
	assert.Nil(t, err)

	err = copyHandler.Process()
	assert.Nil(t, err)

	expectedAdata := map[string]string{
		"A-key1": "A-value-d-1",
		"A-key2": "A-value-d-2",
		"A-key3": "A-value-d-3",
	}
	expectedBdata := map[string]string{
		"B-key1": "B-value-d-1",
		"B-key2": "B-value-d-2",
		"B-key3": "B-value-s-3", // copied from src
		"B-key4": "B-value-d-4",
	}
	expectedCdata := map[string]string{
		"C-key1": "C-value-d-1",
		"C-key2": "C-value-d-2",
		"C-key3": "C-value-d-3",
	}
	expectedDdata := make(map[string]string)
	var expectedEdata map[string]string = nil
	expectedFdata := map[string]string{
		"F-key1": "F-value-d-1",
	}

	assert.Equal(t, expectedAdata, getAllData(t, path.Join(destParentDir, "A")))
	assert.Equal(t, expectedBdata, getAllData(t, path.Join(destParentDir, "B")))
	assert.Equal(t, expectedCdata, getAllData(t, path.Join(destParentDir, "C")))
	assert.Equal(t, expectedDdata, getAllData(t, path.Join(destParentDir, "D")))
	assert.Equal(t, expectedEdata, getAllData(t, path.Join(destParentDir, "E")))
	assert.Equal(t, expectedFdata, getAllData(t, path.Join(destParentDir, "F")))
}

func setupDirs(t *testing.T) (string, string) {
	srcParentDir := t.TempDir()
	destParentDir := t.TempDir()

	putData(t,
		path.Join(srcParentDir, "A"),
		[]string{"A-key1", "A-key2", "A-key3"},
		[]string{"A-value-s-1", "A-value-s-2", "A-value-s-3"},
	)
	putData(t,
		path.Join(srcParentDir, "B"),
		[]string{"B-key1", "B-key2", "B-key3", "B-key4"},
		[]string{"B-value-s-1", "B-value-s-2", "B-value-s-3", "B-value-s-4"},
	)
	putData(t,
		path.Join(srcParentDir, "C"),
		[]string{"C-key1", "C-key2"},
		[]string{"C-value-s-1", "C-value-s-2"},
	)
	putData(t,
		path.Join(srcParentDir, "D"),
		make([]string, 0),
		make([]string, 0),
	)
	putData(t,
		path.Join(srcParentDir, "E"),
		[]string{"E-key1"},
		[]string{"E-value-s-1"},
	)

	// same keys
	putData(t,
		path.Join(destParentDir, "A"),
		[]string{"A-key1", "A-key2", "A-key3"},
		[]string{"A-value-d-1", "A-value-d-2", "A-value-d-3"},
	)
	// missing key3
	putData(t,
		path.Join(destParentDir, "B"),
		[]string{"B-key1", "B-key2", "B-key4"},
		[]string{"B-value-d-1", "B-value-d-2", "B-value-d-4"},
	)
	// nothing missing, but dest has more keys
	putData(t,
		path.Join(destParentDir, "C"),
		[]string{"C-key1", "C-key2", "C-key3"},
		[]string{"C-value-d-1", "C-value-d-2", "C-value-d-3"},
	)
	putData(t,
		path.Join(destParentDir, "D"),
		make([]string, 0),
		make([]string, 0),
	)
	putData(t,
		path.Join(destParentDir, "F"),
		[]string{"F-key1"},
		[]string{"F-value-d-1"},
	)

	return srcParentDir, destParentDir
}

func putData(t *testing.T, path string, keys []string, values []string) {
	require.Equal(t, len(keys), len(values))

	wrapper := process.NewDBWrapper()
	err := wrapper.Open(path)
	require.Nil(t, err)

	for i := 0; i < len(keys); i++ {
		err = wrapper.Put([]byte(keys[i]), []byte(values[i]))
		require.Nil(t, err)
	}

	err = wrapper.Close()
	require.Nil(t, err)
}

func getAllData(t *testing.T, path string) map[string]string {
	_, err := os.ReadDir(path)
	if err != nil {
		return nil
	}

	wrapper := process.NewDBWrapper()
	err = wrapper.Open(path)
	require.Nil(t, err)

	result := make(map[string]string)
	wrapper.RangeKeys(func(key []byte, val []byte) bool {
		result[string(key)] = string(val)

		return true
	})

	return result
}
