package kv

import "odydb/log"
import "odydb/entry"

type KV struct {
	log log.Log
	mem map[string][]byte
}

func (kv *KV) Open() error {
	if err := kv.log.Open(); err != nil {
		return err
	}

	kv.mem = map[string][]byte{}
	
	ent := entry.Entry{}
	for {
		eof, err := kv.log.Read(&ent)
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
	return kv.log.Close()
}

func (kv *KV) Get(key []byte) (val []byte, ok bool, err error) {
	val, exists := kv.mem[string(key)]
	return val, exists, nil
}

func (kv *KV) Set(key []byte, val []byte) (updated bool, err error) {
	kv.log.Write(&entry.Entry{Key: key, Value: val})
	kv.mem[string(key)] = val
	return true, nil
}

func (kv *KV) Del(key []byte) (deleted bool, err error) {
	_, exists, err := kv.Get(key)
	if err != nil {
		return false, err
	}

	if exists {
		kv.log.Write(&entry.Entry{Key: key, Deleted: true})
		delete(kv.mem, string(key))
	}

	return exists, nil
}

func (kv *KV) Size() int {
	return len(kv.mem)
}