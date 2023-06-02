package cache

import (
	"sync"
	"time"
)

type Map[K comparable, V any] struct {
	entries map[K]V
	lock    sync.Mutex
	timers  map[K]*time.Timer
	ttl     time.Duration
}

func NewMap[K comparable, V any](ttl time.Duration) *Map[K, V] {
	return &Map[K, V]{
		entries: map[K]V{},
		timers:  map[K]*time.Timer{},
		ttl:     ttl,
	}
}

func (m *Map[K, V]) Delete(key K) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.reset(key)
	delete(m.entries, key)
}

func (m *Map[K, V]) Exists(key K) bool {
	m.lock.Lock()
	defer m.lock.Unlock()

	_, ok := m.entries[key]
	return ok
}

func (m *Map[K, V]) Get(key K) V {
	m.lock.Lock()
	defer m.lock.Unlock()

	return m.entries[key]
}

func (m *Map[K, V]) Set(key K, value V) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.entries[key] = value
	m.timers[key] = time.AfterFunc(m.ttl, func() { m.Delete(key) })
	m.reset(key)
}

func (m *Map[K, V]) reset(key K) {
	if t, ok := m.timers[key]; ok {
		t.Stop()
		delete(m.timers, key)
	}
}
