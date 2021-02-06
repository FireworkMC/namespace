package namespace

import (
	"fmt"
	"strings"
	"sync"
)

//DefaultNamespace this is the namespace used for minecraft items/effects ect.
//Plugins should use their own namespace if they aren't overriding vanila features
var DefaultNamespace Namespace = MustGetNamespace("minecraft")

var namespaceMap, namespacedKeyMap sync.Map

//Namespace a namespace
type Namespace interface {
	//String gets the namespace name
	String() string

	Key(s string) (NamespacedKey, error)
	MustKey(s string) NamespacedKey

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
	s = strings.ToLower(s)
	if len(strings.TrimSpace(s)) == 0 || !isValidNamespace(s) {
		return nil, fmt.Errorf("Invalid namespace: '" + s + "'")
	}

	r, _ := namespaceMap.LoadOrStore(s, &namespace{name: s})
	return r.(*namespace), nil
}

//GetNamespacedKey get the given string
func GetNamespacedKey(s string) (NamespacedKey, error) {
	s = strings.ToLower(strings.TrimSpace(s))
	v := strings.Split(s, ":")
	if len(v) != 2 {
		return nil, fmt.Errorf("namespace: NamespacedKeyFor:Invalid string provided")
	}
	namespace, err := GetNamespace(v[0])
	if err != nil {
		return nil, err
	}

	key, err := namespace.Key(v[1])

	return key, err
}

//MustGetNamespace gets or creates the given namespace. This panics if the given namespace is not valid
func MustGetNamespace(s string) Namespace {
	ns, err := GetNamespace(s)
	if err != nil {
		panic(err)
	}
	return ns
}

//MustGetNamespacedKey gets the namespaced key for the give string.
//panics if the given key is invalid.
func MustGetNamespacedKey(s string) NamespacedKey {
	key, err := GetNamespacedKey(s)
	if err != nil {
		panic(err)
	}
	return key
}
