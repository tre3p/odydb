package schema

import (
	"odydb/cell"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSchemaEncode(t *testing.T) {
	schema := &Schema{
		Table: "link",
		Cols: []Column{
			{Name: "time", Type: cell.TypeI64},
			{Name: "src", Type: cell.TypeStr},
			{Name: "dst", Type: cell.TypeStr},
		},
		PKey: []int{1, 2},
	}

	row := Row{
		cell.Cell{Type: cell.TypeI64, I64: 123},
		cell.Cell{Type: cell.TypeStr, Str: []byte("a")},
		cell.Cell{Type: cell.TypeStr, Str: []byte("b")},
	}
	expectedKey := []byte{'l', 'i', 'n', 'k', 0, 2, 1, 0, 0, 0, 'a', 2, 1, 0, 0, 0, 'b'}
	expectedVal := []byte{1, 123, 0, 0, 0, 0, 0, 0, 0}
	assert.Equal(t, expectedKey, row.EncodeKey(schema))
	assert.Equal(t, expectedVal, row.EncodeVal(schema))

	decoded := schema.NewRow()
	err := decoded.DecodeKey(schema, expectedKey)
	assert.Nil(t, err)
	err = decoded.DecodeVal(schema, expectedVal)
	assert.Nil(t, err)
	
	assert.Equal(t, row, decoded)
}