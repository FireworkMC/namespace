package namespace

import (
	"fmt"
	"strings"
	"sync"
)

//DefaultNamespace this is the namespace used for minecraft items/effects ect.
//Plugins should use their own namespace if they aren't overriding vanila features
var DefaultNamespace Namespace = MustNamespace("minecraft")

var namespaceMap, namespacedKeyMap sync.Map

//Namespace a namespace
type Namespace interface {
	//String gets the namespace name
	String() string

	Key(s string) (NamespacedKey, error)

	namespace() *namespace
}

//NamespacedKey a namespaced key
type NamespacedKey interface {
	//Namespace gets the current namespace
	Namespace() Namespace
	//Key get the key
	Key() string

	String() string

	namespacedKey() *namespacedKey
}

//GetNamespace gets the namespace or creates a new namespace if it does not exist.
func GetNamespace(s string) (Namespace, error) {
	r, ok := namespaceMap.Load(s)
	if !ok {
		s = strings.TrimSpace(strings.ToLower(s))
		if len(s) == 0 || !IsValidNamespace(s) {
			return nil, fmt.Errorf("Invalid namespace: '" + s + "'")
		}
		r, _ = namespaceMap.LoadOrStore(s, &namespace{name: s})
	}

	return r.(*namespace), nil
}

//GetNamespacedKey get the given string
func GetNamespacedKey(s string) (NamespacedKey, error) {
	ns, ok := ParseNamespacedKey(s)
	if !ok {
		return nil, fmt.Errorf("namespace: NamespacedKeyFor: Invalid string provided")
	}

	var r interface{}

	if r, ok = namespacedKeyMap.Load(ns); !ok {

		if r, ok = namespaceMap.Load(ns[0]); !ok {
			r, _ = namespaceMap.LoadOrStore(ns[0], &namespace{name: ns[0]})
		}

		namespace := r.(Namespace)

		if r, ok = namespacedKeyMap.Load(ns); !ok {
			v := &namespacedKey{namespace: namespace.namespace(), key: ns[1]}
			r, _ = namespacedKeyMap.LoadOrStore(ns, v)
		}
	}

	return r.(*namespacedKey), nil
}

//MustNamespace gets or creates the given namespace. This panics if the given namespace is not valid
func MustNamespace(s string) Namespace {
	ns, err := GetNamespace(s)
	if err != nil {
		panic(err)
	}
	return ns
}

//MustNamespacedKey gets the namespaced key for the give string.
//panics if the given key is invalid.
func MustNamespacedKey(s string) NamespacedKey {
	key, err := GetNamespacedKey(s)
	if err != nil {
		panic(err)
	}
	return key
}
