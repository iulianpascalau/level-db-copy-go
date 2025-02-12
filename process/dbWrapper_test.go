package process

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewDBWrapper(t *testing.T) {
	t.Parallel()

	t.Run("path error should error", func(t *testing.T) {
		t.Parallel()

		wrapper, err := NewDBWrapper("/root/")
		assert.Nil(t, wrapper)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "permission denied for path /root/")
	})
	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		wrapper, err := NewDBWrapper(t.TempDir())
		assert.NotNil(t, wrapper)
		assert.Nil(t, err)

		_ = wrapper.Close()
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

	wrapper, _ := NewDBWrapper(t.TempDir())

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
}

func TestDbWrapper_PutRangeKeys(t *testing.T) {
	t.Parallel()

	wrapper, _ := NewDBWrapper(t.TempDir())

	wrapper.RangeKeys(func(key []byte, val []byte) bool {
		assert.Fail(t, "should have been an empty storage here")

		return false
	})

	time.Sleep(time.Second)

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
	time.Sleep(time.Second * 3)
	assert.Equal(t, expectedRangeKeys, rangeKeys)

	_ = wrapper.Close()
}
