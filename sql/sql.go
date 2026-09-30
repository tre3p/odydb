package sql

import (
	"odydb/cell"
	"odydb/schema"
)

type SQLResult struct {
	Updated int
	Header  []string
	Values  []schema.Row
}

type StmtSelect struct {
	Table string
	Cols  []string
	Keys  []NamedCell
}

type StmtCreateTable struct {
	Table string
	Cols  []schema.Column
	Pkey  []string
}

type StmtInsert struct {
	Table string
	Value []cell.Cell
}

type StmtUpdate struct {
	Table string
	Keys  []NamedCell
	Value []NamedCell
}

type StmtDelete struct {
	Table string
	Keys  []NamedCell
}

type NamedCell struct {
	Column string
	Value  cell.Cell
}
