package structs

type KV struct {
	mem map[string][]byte
}

func (kv *KV) Open() error {
	kv.mem = map[string][]byte{} // empty
	return nil
}

func (kv *KV) Close() error { return nil }

func (kv *KV) Get(key []byte) (val []byte, ok bool, err error) {
	strKey := kv.getKey(key)
	val, ok = kv.mem[strKey]

	return
}

func (kv *KV) Set(key []byte, val []byte) (updated bool, err error) {
	strKey := kv.getKey(key)
	_, updated = kv.mem[strKey]

	kv.mem[strKey] = val

	return updated, nil
}

func (kv *KV) Del(key []byte) (deleted bool, err error) {
	strKey := kv.getKey(key)
	_, deleted = kv.mem[strKey]

	delete(kv.mem, strKey)

	return deleted, nil
}

func (kv *KV) getKey(key []byte) (strKey string) {
	return string(key)
}
