package sql

import (
	"odydb/cell"
	"odydb/schema"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseName(t *testing.T) {
	p := NewParser(" a b0 _0_ 123 ")
	name, ok := p.tryName()
	assert.True(t, name == "a" && ok)

	name, ok = p.tryName()
	assert.True(t, name == "b0" && ok)

	name, ok = p.tryName()
	assert.True(t, name == "_0_" && ok)

	_, ok = p.tryName()
	assert.False(t, ok)
}

func TestParseKeyword(t *testing.T) {
	p := NewParser(" select  HELLO ")
	assert.False(t, p.tryKeyword("sel"))
	assert.True(t, p.tryKeyword("SELECT"))
	assert.True(t, p.tryKeyword("hello") && p.isEnd())
}

func TestParseValue(t *testing.T) {
	testParseValue(t, " -123 ", cell.Cell{Type: cell.TypeI64, I64: -123})
	testParseValue(t, ` 'abc\'\"d' `, cell.Cell{Type: cell.TypeStr, Str: []byte("abc'\"d")})
	testParseValue(t, ` "abc\'\"d" `, cell.Cell{Type: cell.TypeStr, Str: []byte("abc'\"d")})
}

func TestParseStmt(t *testing.T) {
	var stmt interface{}

	s := "    select a  from  t  where  c=1;  "
	stmt = &StmtSelect{
		table: "t",
		cols:  []string{"a"},
		keys:  []NamedCell{{column: "c", value: cell.Cell{Type: cell.TypeI64, I64: 1}}},
	}
	testParseStatement(t, s, stmt)

	s = "select a,b_02 from T where c = 1 and d='e' ;"
	stmt = &StmtSelect{
		table: "T",
		cols:  []string{"a", "b_02"},
		keys: []NamedCell{
			{column: "c", value: cell.Cell{Type: cell.TypeI64, I64: 1}},
			{column: "d", value: cell.Cell{Type: cell.TypeStr, Str: []byte("e")}},
		},
	}
	testParseStatement(t, s, stmt)

	s = "create table t (a string, b int64, primary key (b));"
	stmt = &StmtCreateTable{
		table: "t",
		cols: []schema.Column{
			{Name: "a", Type: cell.TypeStr},
			{Name: "b", Type: cell.TypeI64},
		},
		pkey: []string{"b"},
	}
	testParseStatement(t, s, stmt)

	s = "insert into t values(1, 'hi');"
	stmt = &StmtInsert{
		table: "t",
		value: []cell.Cell{
			{Type: cell.TypeI64, I64: 1},
			{Type: cell.TypeStr, Str: []byte("hi")},
		},
	}
	//testParseStatement(t, s, stmt)

	s = "update t set a = 1, b = 2 where c = 3 and d = 4;"
	stmt = &StmtUpdate{
		table: "t",
		value: []NamedCell{
			{column: "a", value: cell.Cell{Type: cell.TypeI64, I64: 1}},
			{column: "b", value: cell.Cell{Type: cell.TypeI64, I64: 2}},
		},
		keys: []NamedCell{
			{column: "c", value: cell.Cell{Type: cell.TypeI64, I64: 3}},
			{column: "d", value: cell.Cell{Type: cell.TypeI64, I64: 4}},
		},
	}
	//testParseStatement(t, s, stmt)

	s = "delete from t where c = 3 and d = 4;"
	stmt = &StmtDelete{
		table: "t",
		keys: []NamedCell{
			{column: "c", value: cell.Cell{Type: cell.TypeI64, I64: 3}},
			{column: "d", value: cell.Cell{Type: cell.TypeI64, I64: 4}},
		},
	}
	//testParseStatement(t, s, stmt)
}

func testParseValue(t *testing.T, s string, ref cell.Cell) {
	p := NewParser(s)
	out := cell.Cell{}
	err := p.parseValue(&out)
	assert.Nil(t, err)
	assert.True(t, p.isEnd())
	assert.Equal(t, ref, out)
}

func testParseStatement(t *testing.T, s string, ref interface{}) {
	p := NewParser(s)
	out, err := p.parseStmt()
	assert.Nil(t, err)
	assert.True(t, p.isEnd())
	assert.Equal(t, ref, out)
}
