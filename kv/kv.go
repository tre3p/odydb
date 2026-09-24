package kv

import "odydb/log"
import "odydb/entry"

type KV struct {
	Log log.Log
	mem map[string][]byte
}

type UpdateMode int

const (
	ModeUpsert UpdateMode = 0
	ModeInsert UpdateMode = 1
	ModeUpdate UpdateMode = 2
)

func (kv *KV) Open() error {
	if err := kv.Log.Open(); err != nil {
		return err
	}

	kv.mem = map[string][]byte{}
	
	ent := entry.Entry{}
	for {
		eof, err := kv.Log.Read(&ent)
		if err != nil {
			return err
		}
		if eof {
			break
		}

		if ent.Deleted {
			delete(kv.mem, string(ent.Key))
		} else {
			kv.mem[string(ent.Key)] = ent.Value
		}
	}

	return nil
}

func (kv *KV) Close() error {
	return kv.Log.Close()
}

func (kv *KV) SetEx(key[] byte, val []byte, mode UpdateMode) (bool, error) {
	_, exists, err := kv.Get(key)

	if err != nil {
		return false, err
	}

	switch mode {
	case ModeUpsert:
		return kv.Set(key, val)
	case ModeInsert:
		if !exists {
			return kv.Set(key, val)
		}
	case ModeUpdate:
		if exists {
			return kv.Set(key, val)
		}
	}

	return false, nil
}

func (kv *KV) Get(key []byte) (val []byte, ok bool, err error) {
	val, exists := kv.mem[string(key)]
	return val, exists, nil
}

func (kv *KV) Set(key []byte, val []byte) (updated bool, err error) {
	kv.Log.Write(&entry.Entry{Key: key, Value: val})
	kv.mem[string(key)] = val
	return true, nil
}

func (kv *KV) Del(key []byte) (deleted bool, err error) {
	_, exists, err := kv.Get(key)
	if err != nil {
		return false, err
	}

	if exists {
		kv.Log.Write(&entry.Entry{Key: key, Deleted: true})
		delete(kv.mem, string(key))
	}

	return exists, nil
}

func (kv *KV) Size() int {
	return len(kv.mem)
}