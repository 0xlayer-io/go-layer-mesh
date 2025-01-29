package syncmap

import "github.com/puzpuzpuz/xsync/v3"

type IntegerConstraint interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

func NewIntegerMapOf[K IntegerConstraint, V any]() *SyncMap[K, V] {
	return &SyncMap[K, V]{
		m: xsync.NewMapOf[K, V](),
	}
}

func NewTypedMapOf[K comparable, V any](hasher Hasher[K]) *SyncMap[K, V] {
	return &SyncMap[K, V]{
		m: xsync.NewMapOfWithHasher[K, V](hasher),
	}
}

type SyncMap[K comparable, V any] struct {
	m *xsync.MapOf[K, V]
}

func (m *SyncMap[K, V]) Size() int {
	return m.m.Size()
}

func (m *SyncMap[K, V]) Has(key K) (exists bool) {
	_, exists = m.m.Load(key)
	return
}

func (m *SyncMap[K, V]) Load(key K) (val V, exists bool) {
	return m.m.Load(key)
}

func (m *SyncMap[K, V]) LoadOrStore(key K, val V) (actual V, loaded bool) {
	return m.m.LoadOrStore(key, val)
}

func (m *SyncMap[K, V]) Store(key K, val V) {
	m.m.Store(key, val)
}

func (m *SyncMap[K, V]) Range(f func(key K, value V) bool) {
	m.m.Range(f)
}

func (m *SyncMap[K, V]) Keys() (keys []K) {
	m.Range(func(key K, value V) bool {
		keys = append(keys, key)
		return true
	})
	return
}

func (m *SyncMap[K, V]) Pop() (key K, val V, exists bool) {
	m.Range(func(k K, v V) bool {
		key, val, exists = k, v, true
		return false
	})

	if exists {
		m.Delete(key)
	}
	return
}

func (m *SyncMap[K, V]) Delete(key K) {
	m.m.Delete(key)
}

func (m *SyncMap[K, V]) Clear() {
	m.m.Clear()
}
