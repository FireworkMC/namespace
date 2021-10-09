package namespace

import (
	"fmt"
	"strings"
	"sync"
)

// DefaultNamespace this is the namespace used for minecraft items/effects ect.
// Plugins should use their own namespace if they aren't overriding vanila features
var DefaultNamespace Namespace = MustNamespace("minecraft")

var namespaceMap, namespacedKeyMap sync.Map

// Namespace a namespace.
// The name space may contain lowercase letters, numbers, hyphens and underscores.
// Namespaces can be directly compared with eachother.
type Namespace interface {
	// String gets the namespace name
	String() string
	// Key creates a new namspaced key within this namespace.
	// This returns an error if the given string is invalid.
	// A valid key may only contain lowercase letters, numbers, hyphens, underscores, forward-slashes and dots.
	Key(s string) (Key, error)
	// MustKey like `Key` but panics instead of returning an error if the provided key is invalid.
	MustKey(s string) Key
	namespace() *namespace
}

// Key a namespaced key
// A valid key may only contain lowercase letters, numbers, hyphens, underscores, forward-slashes and dots.
// Keys can also be directly compared with eachother.
type Key interface {
	// Namespace gets the namspace that this key belongs to.
	Namespace() Namespace
	// Key gets the key
	Key() string
	// String gets the namespacedkey as a string in the form of `namespace:key`
	String() string
	namespacedKey() *namespacedKey
}

// GetNamespace gets the namespace or creates a new namespace if it does not exist.
func GetNamespace(s string) (Namespace, error) {
	r, ok := namespaceMap.Load(s)
	if !ok {
		s = strings.TrimSpace(strings.ToLower(s))
		if !IsValidNamespace(s) {
			return nil, fmt.Errorf("Invalid namespace: '" + s + "'")
		}

		r, _ = namespaceMap.LoadOrStore(s, &namespace{name: s})
	}

	return r.(*namespace), nil
}

// GetKey parses the given string and returns it as a namespaced key
func GetKey(s string) (Key, error) {
	ns, ok := ParseNamespacedKey(s)
	if !ok {
		return nil, fmt.Errorf("namespace: GetKey: Invalid string provided")
	}

	var r interface{}

	if r, ok = namespacedKeyMap.Load(ns); !ok {

		if r, ok = namespaceMap.Load(ns[0]); !ok {
			r, _ = namespaceMap.LoadOrStore(ns[0], &namespace{name: ns[0]})
		}

		r = getKeyForNs(r.(Namespace).namespace(), ns)
	}

	return r.(*namespacedKey), nil
}

// getKeyForNs gets the namespaced key for the given namespace and key.
// This function assumes that the given namespace and key are valid.
func getKeyForNs(n *namespace, ns [2]string) Key {
	var r interface{}
	var ok bool
	if r, ok = namespaceMap.Load(ns); !ok {
		v := &namespacedKey{namespace: n, key: ns[1], full: n.name + ":" + ns[1]}
		r, _ = namespacedKeyMap.LoadOrStore(ns, v)
	}
	return r.(Key)
}

// MustNamespace gets or creates the given namespace.
// This panics if the given namespace is not valid
func MustNamespace(s string) Namespace {
	ns, err := GetNamespace(s)
	if err != nil {
		panic(err)
	}
	return ns
}

// MustKey gets the namespaced key for the give string.
// panics if the given key is invalid.
func MustKey(s string) Key {
	key, err := GetKey(s)
	if err != nil {
		panic(err)
	}
	return key
}
