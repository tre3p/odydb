package cell

import "encoding/binary"

type CellType uint8

const (
	TypeI64 CellType = 1
	TypeStr CellType = 2
)

type Cell struct {
	Type CellType
	I64 int64
	Str []byte
}

func (cell *Cell) Encode(toAppend []byte) []byte {
	var encoded []byte

	switch cell.Type {
	case TypeI64:
		encoded = make([]byte, 1 + 8)
		encoded[0] = byte(TypeI64)
		binary.LittleEndian.PutUint64(encoded[1:], uint64(cell.I64))
	case TypeStr:
		encoded = make([]byte, 1 + 4 + len(cell.Str))
		encoded[0] = byte(TypeStr)
		binary.LittleEndian.PutUint32(encoded[1:], uint32(len(cell.Str)))
		copy(encoded[5:], cell.Str)
	}

	return append(toAppend, encoded...)
}

func (cell *Cell) Decode(data []byte) (rest []byte, err error) {
	switch CellType(data[0]) {
	case TypeI64:
		cell.Type = TypeI64
		if _, err = binary.Decode(data[1:], binary.LittleEndian, &cell.I64); err != nil {
			return data, err
		}

		return data[1+8:], nil
	case TypeStr:
		cell.Type = TypeStr
		len := binary.LittleEndian.Uint32(data[1:])
		cell.Str = make([]byte, len)
		copy(cell.Str, data[5:5+len])

		return data[1+4+len:], nil
	}

	return nil, nil
}