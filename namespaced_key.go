package namespace

import (
	"fmt"
	"strings"
	"unicode"
)

type namespace struct {
	name string
}

func (n *namespace) String() string { return n.name }

func (n *namespace) Key(k string) (NamespacedKey, error) {
	k = strings.ToLower(k)
	if n == nil {
		return nil, fmt.Errorf("The provided namespace is nil")
	}
	if len(strings.TrimSpace(k)) == 0 || !isValidNamespacedKey(k) {
		return nil, fmt.Errorf("Invalid namespaced key: '" + k + "'")
	}

	r, _ := namespacedKeyMap.LoadOrStore(k, &namespacedKey{namespace: n.namespace(), key: k})
	return r.(*namespacedKey), nil
}

func (n *namespace) MustKey(k string) NamespacedKey {
	key, err := n.Key(k)
	if err != nil {
		panic(err)
	}
	return key
}

func (n *namespace) namespace() *namespace { return n }

type namespacedKey struct {
	namespace *namespace
	key       string
}

func (n *namespacedKey) String() string { return n.namespace.name + ":" + n.key }

func (n *namespacedKey) Namespace() Namespace { return n.namespace }

func (n *namespacedKey) Key() string { return n.key }

func (n *namespacedKey) namespacedKey() *namespacedKey { return n }

//See https://www.minecraft.net/en-us/article/minecraft-snapshot-17w43a

func isValidNamespace(s string) bool {
	for _, b := range []rune(s) {
		if !(unicode.IsLower(b) || unicode.IsDigit(b) || b == '-' || b == '_') {
			return false
		}
	}
	return true
}

func isValidNamespacedKey(s string) bool {
	for _, b := range []rune(s) {
		if !(unicode.IsLower(b) || unicode.IsDigit(b) || b == '-' || b == '_' || b == '/' || b == '.') {
			return false
		}
	}
	return true
}
