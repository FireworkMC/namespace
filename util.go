package namespace

import (
	"sync"
)

type syncMap[K comparable, V any] struct {
	mux sync.RWMutex
	v   map[K]V
	New func(K) V
}

func (s *syncMap[K, V]) Get(key K) (value V, exists bool) {
	s.mux.RLock()
	ns, ok := s.v[key]
	s.mux.RUnlock()
	return ns, ok
}

func (s *syncMap[K, V]) GetOrCreate(key K) V {
	s.mux.RLock()
	ns, ok := s.v[key]
	s.mux.RUnlock()
	if !ok {
		s.mux.Lock()
		ns, ok = s.v[key]
		if !ok {
			ns = s.New(key)
			s.v[key] = ns
		}
		s.mux.Unlock()
	}

	return ns
}
