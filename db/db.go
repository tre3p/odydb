package db

import "odydb/kv"
import "odydb/schema"

type DB struct {
	KV kv.KV
}

func (db *DB) Open() error { return db.KV.Open() }
func (db *DB) Close() error { return db.KV.Close() }

func (db *DB) Select(schema *schema.Schema, row schema.Row) (ok bool, err error) {
	// todo
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