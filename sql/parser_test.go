package sql

import (
	"odydb/cell"
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

func testParseValue(t *testing.T, s string, ref cell.Cell) {
	p := NewParser(s)
	out := cell.Cell{}
	err := p.parseValue(&out)
	assert.Nil(t, err)
	assert.True(t, p.isEnd())
	assert.Equal(t, ref, out)
}