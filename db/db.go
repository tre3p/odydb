package db

import (
	"encoding/json"
	"errors"
	"odydb/kv"
	"odydb/schema"
	"odydb/sql"
	"slices"
)

type DB struct {
	KV kv.KV
	tables map[string]schema.Schema
}

func (db *DB) Open() error { 
	db.tables = make(map[string]schema.Schema)
	return db.KV.Open()
 }
 
func (db *DB) Close() error { return db.KV.Close() }

func (db *DB) ExecStmt(stmt interface{}) (r sql.SQLResult, err error) {
	switch ptr := stmt.(type) {
	case *sql.StmtCreateTable:
		err = db.execCreateTable(ptr)
	case *sql.StmtSelect:
		r.Header = ptr.Cols
		r.Values, err = db.execSelect(ptr)
	case *sql.StmtInsert:
		r.Updated, err = db.execInsert(ptr)
	case *sql.StmtUpdate:
		r.Updated, err = db.execUpdate(ptr)
	case *sql.StmtDelete:
		r.Updated, err = db.execDelete(ptr)
	default:
		panic("unreachable")
	}

	return
}

func (db *DB) GetSchema(table string) (schema.Schema, error) {
	scheme, ok := db.tables[table]
	if !ok {
		val, ok, err := db.KV.Get([]byte("@schema_"+table))
		if err == nil && ok {
			// todo this one shouldn't be json
			err = json.Unmarshal(val, &scheme)
		}
		if err != nil {
			return schema.Schema{}, err
		}
		if !ok {
			return schema.Schema{}, errors.New("table is not found")
		}
		db.tables[table] = scheme
	}

	return scheme, nil
}

func (db *DB) execCreateTable(stmt *sql.StmtCreateTable) error {
	if _, err := db.GetSchema(stmt.Table); err == nil {
		return errors.New("duplicate table")
	}

	columnIndices, err := lookupColumns(stmt.Cols, stmt.Pkey)
	if err != nil {
		return err
	}

	scheme := schema.Schema{
		Table: stmt.Table,
		Cols: stmt.Cols,
		PKey: columnIndices,
	}

	// todo this shouldn't be json
	serializedScheme, err := json.Marshal(scheme)
	if err != nil {
		return err
	}

	db.KV.Set([]byte("@schema_" + stmt.Table), serializedScheme)
	db.tables[stmt.Table] = scheme

	return nil
}

func (db *DB) execDelete(stmt *sql.StmtDelete) (int, error) {
	scheme, err := db.GetSchema(stmt.Table)
	if err != nil {
		return 0, errors.New("table is not found")
	}

	pKey, err := makePKey(&scheme, stmt.Keys)
	if err != nil {
		return 0, err
	}

	deleted, err := db.Delete(&scheme, pKey)
	if err != nil {
		return 0, err
	}

	if deleted {
		return 1, nil
	} else {
		return 0, nil
	}
}

func (db *DB) execUpdate(stmt *sql.StmtUpdate) (int, error) {
	scheme, err := db.GetSchema(stmt.Table)
	if err != nil {
		return 0, err
	}

	finalRow, err := makePKey(&scheme, stmt.Keys)
	if err != nil {
		return 0, err
	}

	var ok bool
	if ok, err = db.Select(&scheme, finalRow); err != nil {
		return 0, err
	}
	if !ok {
		return 0, errors.New("record doesn't exists")
	}

	if err = fillNonPKey(scheme, stmt.Value, finalRow); err != nil {
		return 0, err
	}

	if ok, err = db.Update(&scheme, finalRow); err != nil {
		return 0, err
	}

	if !ok {
		return 0, errors.New("unable to update")
	}

	return 1, nil
}

func fillNonPKey(scheme schema.Schema, values []sql.NamedCell, row schema.Row) error {
	if len(values) != (len(scheme.Cols) - len(scheme.PKey)) {
		return errors.New("invalid values size")
	}

	indices, err := lookupColumnsNamed(scheme.Cols, values)
	if err != nil {
		return err
	}

	for valIdx, val := range values {
		row[indices[valIdx]] = val.Value
	}

	return nil
}


func (db *DB) execInsert(stmt *sql.StmtInsert) (int, error) {
	scheme, err := db.GetSchema(stmt.Table)
	if err != nil {
		return 0, err
	}

	if len(scheme.Cols) != len(stmt.Value) {
		return 0, errors.New("insert columns count doesn't match schema")
	}


	ok, err := db.Insert(&scheme, stmt.Value)

	if err != nil {
		return 0, err
	}

	if ok {
		return 1, nil
	} else {
		return 0, nil
	}
}


func (db *DB) execSelect(stmt *sql.StmtSelect)  ([]schema.Row, error) {
	scheme, err := db.GetSchema(stmt.Table)
	if err != nil {
		return nil, err
	}

	indices, err := lookupColumns(scheme.Cols, stmt.Cols)
	if err != nil {
		return nil, err
	}

	row, err := makePKey(&scheme, stmt.Keys)
	if err != nil {
		return nil, err
	}

	var ok bool
	if ok, err = db.Select(&scheme, row); err != nil {
		return nil, err
	}

	if !ok {
		return nil, nil
	}

	row = subsetRow(row, indices)
	return []schema.Row{row}, nil
}

func makePKey(scheme *schema.Schema, pkey []sql.NamedCell) (schema.Row, error) {
	if len(scheme.PKey) != len(pkey) {
		return nil, errors.New("not primary key")
	}

	row := scheme.NewRow()
	for _, idx1 := range scheme.PKey {
		col := scheme.Cols[idx1]
		idx2 := slices.IndexFunc(pkey, func(expr sql.NamedCell) bool {
			return expr.Column == col.Name && expr.Value.Type == col.Type
		})
		if idx2 < 0 {
			return nil, errors.New("not primary key")
		}

		row[idx1] = pkey[idx2].Value
	}

	return row, nil
}

func subsetRow(row schema.Row, indices []int) (out schema.Row) {
	for _, idx := range indices {
		out = append(out, row[idx])
	}

	return
}

func lookupColumnsNamed(schemaCols []schema.Column, colNamedCells []sql.NamedCell) ([]int, error) {
	var colNames []string

	for _, col := range colNamedCells {
		colNames = append(colNames, col.Column)
	}

	return lookupColumns(schemaCols, colNames)
}

func lookupColumns(schemaCols []schema.Column, colNames []string) ([]int, error) {
	var indices []int

	for _, name := range colNames {
		idx := slices.IndexFunc(schemaCols, func(col schema.Column) bool {
			return col.Name == name
		})
		if idx < 0 {
			return nil, errors.New("column is not found")
		}
		
		indices = append(indices, idx)
	}

	return indices, nil
}

func (db *DB) Select(schema *schema.Schema, row schema.Row) (ok bool, err error) {
	key := row.EncodeKey(schema)
	val, ok, err := db.KV.Get(key)
	if err != nil || !ok {
		return false, err
	}

	if err = row.DecodeVal(schema, val); err != nil {
		return false, err
	}

	return true, nil
}

func (db *DB) Insert(schema *schema.Schema, row schema.Row) (updated bool, err error) {
	key := row.EncodeKey(schema)
	value := row.EncodeVal(schema)
	return db.KV.SetEx(key, value, kv.ModeInsert)
}

func (db *DB) Upsert(schema *schema.Schema, row schema.Row) (updated bool, err error) {
	key := row.EncodeKey(schema)
	val := row.EncodeVal(schema)

	return db.KV.SetEx(key, val, kv.ModeUpsert)
}

func (db *DB) Update(schema *schema.Schema, row schema.Row) (updated bool, err error) {
	key := row.EncodeKey(schema)
	val := row.EncodeVal(schema)

	return db.KV.SetEx(key, val, kv.ModeUpdate)
}

func (db *DB) Delete(schema *schema.Schema, row schema.Row) (deleted bool, err error) {
	key := row.EncodeKey(schema)

	return db.KV.Del(key)
}