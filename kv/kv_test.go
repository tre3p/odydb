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
	updated, err = kv.Set([]byte("k2"), []byte("v2"))
	assert.True(t, updated && err == nil)

	kv.Close()
	err = kv.Open()
	assert.Nil(t, err)

	assert.Equal(t, 1, kv.Size())

	_, ok, err = kv.Get([]byte("k1"))
	assert.True(t, !ok && err == nil)

	val, ok, err = kv.Get([]byte("k2"))
	assert.True(t, string(val) == "v2" && ok && err == nil)
	// Reopen Case
}

func TestKVRecovery(t *testing.T) {
	kv := KV{}
	kv.log.FileName = ".test_db"
	defer os.Remove(kv.log.FileName)

	prepare := func() {
		os.Remove(kv.log.FileName)

		err := kv.Open()
		assert.Nil(t, err)
		defer kv.Close()

		updated, err := kv.Set([]byte("k1"), []byte("v1"))
		assert.True(t, updated && err == nil)
		updated, err = kv.Set([]byte("k2"), []byte("v2"))
		assert.True(t, updated && err == nil)
	}

	// Truncated log case
	prepare()
	fp, _ := os.OpenFile(kv.log.FileName, os.O_RDWR, 0o644)
	st, _ := fp.Stat()
	fp.Truncate(st.Size() - 1)
	fp.Close()

	err := kv.Open()
	assert.Nil(t, err)

	val, ok, err := kv.Get([]byte("k1"))
	assert.True(t, string(val) == "v1" && ok && err == nil)

	_, ok, err = kv.Get([]byte("k2"))
	assert.True(t, !ok && err == nil)
	kv.Close()
	// Truncated log case

	// Bad Checksum Case
	prepare()
	fp, _ = os.OpenFile(kv.log.FileName, os.O_RDWR, 0o644)
	st, _ = fp.Stat()
	fp.WriteAt([]byte{0}, st.Size() - 1)
	fp.Close()

	err = kv.Open()
	assert.Nil(t, err)

	val, ok, err = kv.Get([]byte("k1"))
	assert.True(t, string(val) == "v1" && ok && err == nil)

	_, ok, err = kv.Get([]byte("k2"))
	assert.True(t, !ok && err == nil)
	kv.Close()

	// Bad Checksum Case
}

func TestKvUpdateMode(t *testing.T) {
	kv := KV{}
	kv.log.FileName = ".test_db"
	defer os.Remove(kv.log.FileName)

	os.Remove(kv.log.FileName)
	err := kv.Open()
	assert.Nil(t, err)
	defer kv.Close()

	updated, err := kv.SetEx([]byte("k1"), []byte("v1"), ModeUpdate)
	assert.True(t, !updated && err == nil)

	updated, err = kv.SetEx([]byte("k1"), []byte("v1"), ModeUpdate)
	assert.True(t, !updated && err == nil)

	updated, err = kv.SetEx([]byte("k1"), []byte("v1"), ModeInsert)
	assert.True(t, updated && err == nil)

	updated, err = kv.SetEx([]byte("k1"), []byte("xx"), ModeInsert)
	assert.True(t, !updated && err == nil)

	updated, err = kv.SetEx([]byte("k1"), []byte("yy"), ModeUpdate)
	assert.True(t, updated && err == nil)

	updated, err = kv.SetEx([]byte("k1"), []byte("zz"), ModeUpsert)
	assert.True(t, updated && err == nil)

	updated, err = kv.SetEx([]byte("k2"), []byte("tt"), ModeUpsert)
	assert.True(t, updated && err == nil)
}