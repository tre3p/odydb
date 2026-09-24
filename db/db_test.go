package db

import (
	"odydb/cell"
	"odydb/schema"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTableByPKey(t *testing.T) {
	db := DB{}
	db.KV.Log.FileName = ".test_db"
	defer os.Remove(db.KV.Log.FileName)

	os.Remove(db.KV.Log.FileName)
	err := db.Open()
	assert.Nil(t, err)
	defer db.Close()

	scheme := &schema.Schema{
		Table: "link",
		Cols: []schema.Column{
			{Name: "time", Type: cell.TypeI64},
			{Name: "src", Type: cell.TypeStr},
			{Name: "dst", Type: cell.TypeStr},
		},
		PKey: []int{1, 2},
	}

	row := schema.Row{
		cell.Cell{Type: cell.TypeI64, I64: 123},
		cell.Cell{Type: cell.TypeStr, Str: []byte("a")},
		cell.Cell{Type: cell.TypeStr, Str: []byte("b")},
	}

	ok, err := db.Select(scheme, row)
	assert.True(t, !ok && err == nil)

	updated, err := db.Insert(scheme, row)
	assert.True(t, updated && err == nil)

	out := schema.Row{
		cell.Cell{},
		cell.Cell{Type: cell.TypeStr, Str: []byte("a")},
		cell.Cell{Type: cell.TypeStr, Str: []byte("b")},
	}
	ok, err = db.Select(scheme, out)
	assert.True(t, ok && err == nil)
	assert.Equal(t, row, out)

	row[0].I64 = 456
	updated, err = db.Update(scheme, row)
	assert.True(t, updated && err == nil)

	ok, err = db.Select(scheme, out)
	assert.True(t, ok && err == nil)
	assert.Equal(t, row, out)

	deleted, err := db.Delete(scheme, row)
	assert.True(t, deleted && err == nil)

	ok, err = db.Select(scheme, row)
	assert.True(t, !ok && err == nil)
}