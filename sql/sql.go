package sql

import (
	"odydb/cell"
	"odydb/schema"
)

type StmtSelect struct {
	table string
	cols  []string
	keys  []NamedCell
}

type StmtCreateTable struct {
	table string
	cols  []schema.Column
	pkey  []string
}

type StmtInsert struct {
	table string
	value []cell.Cell
}

type StmtUpdate struct {
	table string
	keys  []NamedCell
	value []NamedCell
}

type StmtDelete struct {
	table string
	keys []NamedCell
}

type NamedCell struct {
	column string
	value  cell.Cell
}
