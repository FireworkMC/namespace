package namespace

import (
	"fmt"
	"strings"
)

type namespace struct {
	name string
}

func (n *namespace) String() string { return n.name }

func (n *namespace) Key(k string) (Key, error) {
	k = strings.TrimSpace(strings.ToLower(k))
	if !IsValidKey(k) {
		return nil, fmt.Errorf("namespace: Invalid key")
	}
	return getKeyForNs(n, [2]string{n.name, k}), nil
}

func (n *namespace) MustKey(k string) (key Key) {
	var err error
	if key, err = n.Key(k); err == nil {
		return key
	}
	panic(err)
}

func (n *namespace) namespace() *namespace { return n }

type namespacedKey struct {
	namespace *namespace
	key       string
	full      string
}

func (n *namespacedKey) String() string { return n.full }

func (n *namespacedKey) Namespace() Namespace { return n.namespace }

func (n *namespacedKey) Key() string { return n.key }

func (n *namespacedKey) namespacedKey() *namespacedKey { return n }
