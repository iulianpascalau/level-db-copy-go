package process

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewDBWrapper(t *testing.T) {
	t.Parallel()

	wrapper := NewDBWrapper()
	assert.NotNil(t, wrapper)
}

func TestDbWrapper_Open(t *testing.T) {
	t.Run("path error should error", func(t *testing.T) {
		t.Parallel()

		wrapper := NewDBWrapper()
		err := wrapper.Open("/root/")
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "permission denied for path /root/")
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		wrapper := NewDBWrapper()
		err := wrapper.Open(t.TempDir())
		assert.Nil(t, err)

		_ = wrapper.Close()
	})
	t.Run("double open should not be allowed", func(t *testing.T) {
		t.Parallel()

		wrapper := NewDBWrapper()
		err := wrapper.Open(t.TempDir())
		assert.Nil(t, err)

		err = wrapper.Open(t.TempDir())
		assert.Equal(t, errInnerDBIsNotClosed, err)
	})
}

func TestDbWrapper_IsInterfaceNil(t *testing.T) {
	t.Parallel()

	var instance *dbWrapper
	assert.True(t, instance.IsInterfaceNil())

	instance = &dbWrapper{}
	assert.False(t, instance.IsInterfaceNil())
}

func TestDbWrapper_GetPutClose(t *testing.T) {
	t.Parallel()

	wrapper := NewDBWrapper()

	t.Run("Put in an unopened DB should error", func(t *testing.T) {
		err := wrapper.Put([]byte("key1"), []byte("val1"))
		assert.Equal(t, errInnerDBIsNotOpened, err)
	})
	t.Run("Get from an unopened DB should error", func(t *testing.T) {
		value, err := wrapper.Get([]byte("key1"))
		assert.Equal(t, errInnerDBIsNotOpened, err)
		assert.Nil(t, value)
	})
	t.Run("should work", func(t *testing.T) {
		_ = wrapper.Open(t.TempDir())

		wrapper.RangeKeys(func(key []byte, val []byte) bool {
			assert.Fail(t, "should have been an empty storage here")

			return false
		})

		err := wrapper.Put([]byte("key1"), []byte("value1"))
		assert.Nil(t, err)

		err = wrapper.Put([]byte("key2"), []byte("value2"))
		assert.Nil(t, err)

		// we need this sleep to put the data inside the files
		time.Sleep(time.Second * 3)

		recoveredValue, err := wrapper.Get([]byte("key1"))
		assert.Nil(t, err)
		assert.Equal(t, "value1", string(recoveredValue))

		recoveredValue, err = wrapper.Get([]byte("key2"))
		assert.Nil(t, err)
		assert.Equal(t, "value2", string(recoveredValue))

		err = wrapper.Close()
		assert.Nil(t, err)
	})
}

func TestDbWrapper_PutRangeKeys(t *testing.T) {
	t.Parallel()

	wrapper := NewDBWrapper()
	t.Run("RangeKeys should not call handler if the DB is not opened", func(t *testing.T) {
		wrapper.RangeKeys(func(key []byte, val []byte) bool {
			assert.Fail(t, "should have not called the handler")

			return false
		})

		time.Sleep(time.Second)
	})
	t.Run("should work", func(t *testing.T) {
		_ = wrapper.Open(t.TempDir())

		wrapper.RangeKeys(func(key []byte, val []byte) bool {
			assert.Fail(t, "should have been an empty storage here")

			return false
		})

		err := wrapper.Put([]byte("key1"), []byte("value1"))
		assert.Nil(t, err)

		err = wrapper.Put([]byte("key2"), []byte("value2"))
		assert.Nil(t, err)

		time.Sleep(time.Second * 3)

		rangeKeys := make(map[string]string, 2)
		wrapper.RangeKeys(func(key []byte, val []byte) bool {
			rangeKeys[string(key)] = string(val)

			return true
		})
		expectedRangeKeys := map[string]string{
			"key1": "value1",
			"key2": "value2",
		}
		assert.Equal(t, expectedRangeKeys, rangeKeys)

		_ = wrapper.Close()
	})
}
