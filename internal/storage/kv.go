package storage

import "sync"

type KV struct {
	log Log
	mem map[string][]byte
	mu  sync.RWMutex
}

func (kv *KV) Open() error {
	kv.mem = map[string][]byte{} // empty
	return nil
}

func (kv *KV) Close() error { return nil }

func (kv *KV) Get(key []byte) (val []byte, ok bool, err error) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	strKey := kv.getKey(key)
	val, ok = kv.mem[strKey]

	return
}

func (kv *KV) Set(key []byte, val []byte) (updated bool, err error) {
	strKey := kv.getKey(key)

	kv.mu.Lock()
	defer kv.mu.Unlock()

	_, updated = kv.mem[strKey]

	kv.mem[strKey] = val

	return updated, nil
}

func (kv *KV) Del(key []byte) (deleted bool, err error) {
	strKey := kv.getKey(key)

	kv.mu.Lock()
	defer kv.mu.Unlock()

	_, deleted = kv.mem[strKey]

	delete(kv.mem, strKey)

	return deleted, nil
}

func (kv *KV) getKey(key []byte) (strKey string) {
	return string(key)
}
