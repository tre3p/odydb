package kv

type KV struct {
	mem map[string][]byte
}

func (kv *KV) Open() error {
	kv.mem = map[string][]byte{}
	return nil
}

func (kv *KV) Close() error {
	return nil
}

func (kv *KV) Get(key []byte) (val []byte, ok bool, err error) {
	val, exists := kv.mem[string(key)]
	return val, exists, nil
}

func (kv *KV) Set(key []byte, val []byte) (updated bool, err error) {
	kv.mem[string(key)] = val
	return true, nil
}

func (kv *KV) Del(key []byte) (deleted bool, err error) {
	_, exists, err := kv.Get(key)
	if err != nil {
		return false, err
	}

	delete(kv.mem, string(key))
	return exists, nil
}