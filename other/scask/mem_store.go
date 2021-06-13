package main

import "sync"

type memStore struct {
	sync.RWMutex
	storage map[string]string
}

func newMemStore() *memStore {
	return &memStore{storage: make(map[string]string)}
}

func (mem *memStore) get(key string) (string, error) {
	mem.RLock()
	defer mem.RUnlock()
	value, ok := mem.storage[key]
	if !ok {
		return "", errKeyNotFount
	}

	return value, nil
}

func (mem *memStore) put(key string, value string) error {
	mem.Lock()
	defer mem.Unlock()
	mem.storage[key] = value

	return nil
}

func (mem *memStore) del(key string) error {
	mem.Lock()
	defer mem.Unlock()
	delete(mem.storage, key)

	return nil
}
