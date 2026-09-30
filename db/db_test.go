package db

import (
	"odydb/cell"
	"odydb/schema"
	"odydb/sql"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func parseStmt(t *testing.T, s string) interface{} {
	p := sql.NewParser(s)
	stmt, err := p.ParseStmt()
	assert.Nil(t, err)
	return stmt
}

func TestSQLByPKey(t *testing.T) {
	db := DB{}
	db.KV.Log.FileName = ".test_db"
	defer os.Remove(db.KV.Log.FileName)

	os.Remove(db.KV.Log.FileName)
	err := db.Open()
	assert.Nil(t, err)
	defer db.Close()

	s := "create table link (time int64, src string, dst string, primary key (src, dst));"
	_, err = db.ExecStmt(parseStmt(t, s))
	assert.Nil(t, err)

	s = "insert into link values(123, 'bob', 'alice');"
	r, err := db.ExecStmt(parseStmt(t, s))
	assert.Nil(t, err)
	assert.Equal(t, 1, r.Updated)

	s = "select time from link where dst = 'alice' and src = 'bob';"
	r, err = db.ExecStmt(parseStmt(t, s))
	assert.Nil(t, err)
	assert.Equal(t, []schema.Row{{cell.Cell{Type: cell.TypeI64, I64: 123}}}, r.Values)

	s = "update link set time = 456 where dst = 'alice' and src = 'bob';"
	r, err = db.ExecStmt(parseStmt(t, s))
	assert.Nil(t, err)
	assert.Equal(t, 1, r.Updated)

	s = "select time from link where dst = 'alice' and src = 'bob';"
	r, err = db.ExecStmt(parseStmt(t, s))
	assert.Nil(t, err)
	assert.Equal(t, []schema.Row{{cell.Cell{Type: cell.TypeI64, I64: 456}}}, r.Values)

	// reopen
	err = db.Close()
	require.Nil(t, err)
	db = DB{}
	db.KV.Log.FileName = ".test_db"
	err = db.Open()
	assert.Nil(t, err)

	s = "delete from link where src = 'bob' and dst = 'alice';"
	r, err = db.ExecStmt(parseStmt(t, s))
	assert.Nil(t, err)
	assert.Equal(t, 1, r.Updated)

	s = "select time from link where dst = 'alice' and src = 'bob';"
	r, err = db.ExecStmt(parseStmt(t, s))
	assert.Nil(t, err)
	assert.Equal(t, 0, len(r.Values))
}