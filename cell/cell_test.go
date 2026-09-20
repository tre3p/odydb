package cell

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCellEncodeDecode(t *testing.T) {
	// I64 Tests
	int64Cell := Cell{Type: TypeI64, I64: -2}
	int64CellEncodedExpected := []byte{0x1, 0xfe, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	assert.Equal(t, int64CellEncodedExpected, int64Cell.Encode(nil))

	emptyIntCell := Cell{}
	rest, err := emptyIntCell.Decode(int64CellEncodedExpected)
	assert.True(t, len(rest) == 0 && err == nil)
	assert.Equal(t, emptyIntCell, int64Cell)


	// Str Tests
	strCell := Cell{Type: TypeStr, Str: []byte{'a', 's', 'd', 'f'}}
	strCellEncodedExpected := []byte{2, 4, 0, 0, 0, 'a', 's', 'd', 'f'}
	assert.Equal(t, strCellEncodedExpected, strCell.Encode(nil))
	emptyStrCell := Cell{}
	rest, err = emptyStrCell.Decode(strCellEncodedExpected)
	assert.True(t, len(rest) == 0 && err == nil)
	assert.Equal(t, strCell, emptyStrCell)
}
