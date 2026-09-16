package entry

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestEntryBasicEncodeDecode(t *testing.T) {
	ent := Entry{key: []byte("k1"), value: []byte("v1")}
	entExpectedEncode := []byte{2, 0, 0, 0, 2, 0, 0, 0, 'k', '1', 'v', '1'}
	assert.Equal(t, entExpectedEncode, ent.Encode())

	decoded := Entry{}
	err := decoded.Decode(bytes.NewBuffer(entExpectedEncode))
	assert.Nil(t, err)
	assert.Equal(t, ent, decoded)
}