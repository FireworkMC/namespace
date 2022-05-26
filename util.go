package namespace

import (
	"sync"
)

type syncMap[K comparable, V any] struct {
	mux sync.RWMutex
	v   map[K]V
	New func(K) V
}

func (n *syncMap[K, V]) Get(v K) V {
	n.mux.RLock()
	ns, ok := n.v[v]
	n.mux.RUnlock()
	if !ok {
		n.mux.Lock()
		ns, ok = n.v[v]
		if !ok {
			ns = n.New(v)
			n.v[v] = ns
		}
		n.mux.Unlock()
	}

	return ns
}
