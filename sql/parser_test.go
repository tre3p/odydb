package sql

import (
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