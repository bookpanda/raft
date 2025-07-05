package raft

import "sync"

type Storage interface {
	Set(key string, value []byte)

	Get(key string) ([]byte, bool)

	HasData() bool
}

// simple in-memory storage implementation
type MapStorage struct {
	mu sync.Mutex
	m  map[string][]byte
}

func NewMapStorage() *MapStorage {
	s := &MapStorage{}
	s.m = make(map[string][]byte)
	return s
}

func (ms *MapStorage) Get(key string) ([]byte, bool) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	value, ok := ms.m[key]
	return value, ok
}

func (ms *MapStorage) Set(key string, value []byte) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.m[key] = value
}

func (ms *MapStorage) HasData() bool {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	return len(ms.m) > 0
}
