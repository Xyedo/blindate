package rwmap

import "sync"

func New[K comparable, V any]() *RwMap[K, V] {
	return &RwMap[K, V]{
		mu: sync.RWMutex{},
		m:  make(map[K]V),
	}
}

type RwMap[K comparable, V any] struct {
	mu sync.RWMutex
	m  map[K]V
}

func (rw *RwMap[K, V]) Set(key K, value V) {
	rw.mu.Lock()
	defer rw.mu.Unlock()
	rw.m[key] = value

}

func (rw *RwMap[K, V]) Get(key K) (V, bool) {
	rw.mu.RLock()
	defer rw.mu.RUnlock()
	value, ok := rw.m[key]
	return value, ok
}

func (rw *RwMap[K, V]) Delete(key K) {
	rw.mu.Lock()
	defer rw.mu.Unlock()
	delete(rw.m, key)
}
