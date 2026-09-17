package entry

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntryBasicEncodeDecode(t *testing.T) {
	ent := Entry{Key: []byte("k1"), Value: []byte("v1"), Deleted: false}
	entExpectedEncode := []byte{2, 0, 0, 0, 2, 0, 0, 0, 0, 'k', '1', 'v', '1'}
	assert.Equal(t, entExpectedEncode, ent.Encode())

	decoded := Entry{}
	err := decoded.Decode(bytes.NewBuffer(entExpectedEncode))
	assert.Nil(t, err)
	assert.Equal(t, ent, decoded)

	ent = Entry{Key: []byte("k1"), Deleted: true}
	data := []byte{2, 0, 0, 0, 0, 0, 0, 0, 1, 'k', '1'}
	assert.Equal(t, data, ent.Encode())

	decoded = Entry{}
	err = decoded.Decode(bytes.NewBuffer(data))
	assert.Nil(t, err)
	assert.Equal(t, ent, decoded)
}
