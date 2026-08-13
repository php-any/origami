package runtime

import "sync"

func syncMapLoad[V any](m *sync.Map, key string) (V, bool) {
	v, ok := m.Load(key)
	if !ok {
		var zero V
		return zero, false
	}
	return v.(V), true
}

func syncMapStore[V any](m *sync.Map, key string, val V) {
	m.Store(key, val)
}

func syncMapClear(m *sync.Map) {
	m.Range(func(key, _ any) bool {
		m.Delete(key)
		return true
	})
}
