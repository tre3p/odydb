package kv

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKVBasic(t *testing.T) {
	kv := KV{}
	kv.log.FileName = ".test_db"
	defer os.Remove(kv.log.FileName)

	os.Remove(kv.log.FileName) // cleanup before test
	err := kv.Open()
	assert.Nil(t, err)
	defer kv.Close()


	// Set And Get Case
	updated, err := kv.Set([]byte("k1"), []byte("v1"))
	assert.True(t, updated && err == nil)

	val, ok, err := kv.Get([]byte("k1"))
	assert.True(t, string(val) == "v1" && ok && err == nil)
	// Set And Get Case

	// Get Non-Existent Key Case
	_, ok, err = kv.Get([]byte("xxx"))
	assert.True(t, !ok && err == nil)
	// Get Non-Existent Key Case

	// Delete Non-Existent Key Case
	deleted, err := kv.Del([]byte("xxx"))
	assert.True(t, !deleted && err == nil)
	// Delete Non-Existent Key Case

	// Update Existing Key And Get Case
	updated, err = kv.Set([]byte("k1"), []byte("v2"))
	assert.True(t, updated && err == nil)

	val, ok, err = kv.Get([]byte("k1"))
	assert.True(t, string(val) == "v2" && ok && err == nil)
	// Update Existing Key And Get Case

	// Delete Existing Key And Get Case
	deleted, err = kv.Del([]byte("k1"))
	assert.True(t, deleted && err == nil)

	_, ok, err = kv.Get([]byte("k1"))
	assert.True(t, !ok && err == nil)
	// Delete Existing Key And Get Case

	// Reopen Case
	_, err = kv.Set([]byte("k1"), []byte("v1"))
	assert.Nil(t, err)
	val, ok, err = kv.Get([]byte("k1"))
	assert.True(t, string(val) == "v1" && ok && err == nil)

	kv.Close()
	err = kv.Open()
	assert.Nil(t, err)

	assert.Equal(t, 1, kv.Size())

	val, ok, err = kv.Get([]byte("k1"))
	assert.True(t, string(val) == "v1" && ok && err == nil)

	// Reopen Case
}