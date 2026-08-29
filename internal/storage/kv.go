package storage

import (
	"sync"
)

type KV struct {
	log Log
	mem map[string][]byte
	mu  sync.RWMutex
}

func (kv *KV) Open() error {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	if err := kv.log.Open(); err != nil {
		return err
	}

	kv.mem = map[string][]byte{}

	for {
		var ent Entry

		eof, err := kv.log.Read(&ent)

		if err != nil {
			return err
		}
		if eof {
			break
		}

		if ent.deleted {
			delete(kv.mem, string(ent.key))
		} else {
			kv.mem[string(ent.key)] = ent.val
		}
	}

	return nil
}

func (kv *KV) Close() error {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	return kv.log.Close()
}

func (kv *KV) Get(key []byte) (val []byte, ok bool, err error) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	strKey := string(key)
	val, ok = kv.mem[strKey]

	return
}

func (kv *KV) Set(key []byte, val []byte) (updated bool, err error) {
	strKey := string(key)

	kv.mu.Lock()
	defer kv.mu.Unlock()

	_, updated = kv.mem[strKey]

	ent := Entry{
		key:     key,
		val:     val,
		deleted: false,
	}

	if err := kv.log.Write(&ent); err != nil {
		return false, err
	}

	kv.mem[strKey] = val

	return updated, nil
}

func (kv *KV) Del(key []byte) (deleted bool, err error) {
	strKey := string(key)

	kv.mu.Lock()
	defer kv.mu.Unlock()

	_, deleted = kv.mem[strKey]

	ent := Entry{
		key:     key,
		val:     nil,
		deleted: true,
	}

	if err := kv.log.Write(&ent); err != nil {
		return false, err
	}

	delete(kv.mem, strKey)

	return deleted, nil
}
