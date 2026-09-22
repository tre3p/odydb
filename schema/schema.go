package schema

import (
	"bytes"
	"errors"
	"odydb/cell"
	"slices"
)

type Schema struct {
	Table string
	Cols []Column
	PKey []int
}

type Column struct {
	Name string
	Type cell.CellType
}

type Row []cell.Cell

func (schema *Schema) NewRow() Row {
	return make(Row, len(schema.Cols))
}

func (row Row) EncodeKey(schema *Schema) (key []byte) {
	key = append([]byte(schema.Table), 0x00)

	for _, pKeyIdx := range schema.PKey {
		pKeyPart := row[pKeyIdx]
		key = pKeyPart.Encode(key)
	}

	return key
}

func (row Row) EncodeVal(schema *Schema) (val []byte) {
	for i, col := range row {
		// If it's primary key or its part - skip it
		if slices.Contains(schema.PKey, i) {
			continue
		}

		val = col.Encode(val)
	}

	return val
}

func (row Row) DecodeKey(schema *Schema, key []byte) (err error) {
	tableNameDelimIdx := bytes.IndexByte(key, 0x00)
	if tableNameDelimIdx == -1 {
		return errors.New("no table name delimiter found")
	}

	keysArr := key[tableNameDelimIdx+1:]

	for _, pKeyIdx := range schema.PKey {
		keyCell := cell.Cell{}
		keysArr, err = keyCell.Decode(keysArr)
		if err != nil {
			return err
		}

		row[pKeyIdx] = keyCell
	}

	return nil
}

func (row Row) DecodeVal(schema *Schema, val []byte) (err error) {
	for colIdx := range schema.Cols {
		// If it's primary key or its part - skip it
		if slices.Contains(schema.PKey, colIdx) {
			continue
		}

		valCell := cell.Cell{}
		val, err = valCell.Decode(val)
		if err != nil {
			return err
		}

		row[colIdx] = valCell
	}

	return nil
}