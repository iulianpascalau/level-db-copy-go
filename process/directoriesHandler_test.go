package process

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDirectoriesHandler(t *testing.T) {
	t.Parallel()

	t.Run("can not read source parent directory should error", func(t *testing.T) {
		t.Parallel()

		handler, err := NewDirectoriesHandler("/no-root-dir", "./testdata/dir2")
		assert.Nil(t, handler)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "open /no-root-dir: no such file or directory")
	})
	t.Run("can not read destination parent directory should error", func(t *testing.T) {
		t.Parallel()

		handler, err := NewDirectoriesHandler("./testdata/dir1", "/no-root-dir")
		assert.Nil(t, handler)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "open /no-root-dir: no such file or directory")
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		handler, err := NewDirectoriesHandler("./testdata/dir1", "./testdata/dir2")
		assert.NotNil(t, handler)
		assert.Nil(t, err)

		sourceDirs := handler.SourceDirectories()
		expectedSourceDirs := []string{
			"testdata/dir1/aaaa",
			"testdata/dir1/bbbb",
		}

		destinationDirs := handler.DestinationDirectories()
		expectedDestinationDirs := []string{
			"testdata/dir2/aaaa",
			"testdata/dir2/cccc",
		}

		assert.Equal(t, expectedSourceDirs, sourceDirs)
		assert.Equal(t, expectedDestinationDirs, destinationDirs)
	})
}

func TestDirectoriesHandler_IsInterfaceNil(t *testing.T) {
	t.Parallel()

	var instance *directoriesHandler
	assert.True(t, instance.IsInterfaceNil())

	instance = &directoriesHandler{}
	assert.False(t, instance.IsInterfaceNil())
}
