package namespace

import (
	"fmt"
	"strings"
	"sync"
	"unicode"
)

//DefaultNamespace this is the namespace used for minecraft items/effects ect.
//Plugins should use their own namespace if they aren't overriding vanila features
var DefaultNamespace Namespace = func() Namespace { v, _ := GetNamespace("minecraft"); return v }()

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

type namespace struct {
	name string
}

func (n *namespace) String() string { return n.name }

func (n *namespace) Key(s string) (NamespacedKey, error) { return GetNamespacedKey(n, s) }

func (n *namespace) namespace() *namespace { return n }

type namespacedKey struct {
	namespace *namespace
	key       string
}

func (n *namespacedKey) String() string { return n.namespace.name + ":" + n.key }

func (n *namespacedKey) Namespace() Namespace { return n.namespace }

func (n *namespacedKey) Key() string { return n.key }

func (n *namespacedKey) namespacedKey() *namespacedKey { return n }

//GetNamespace gets the namespace or creates a new namespace if it does not exist.
func GetNamespace(s string) (Namespace, error) {
	s = strings.ToLower(s)
	if len(strings.TrimSpace(s)) == 0 || !isValidNamespace(s) {
		return nil, fmt.Errorf("Invalid namespace: '" + s + "'")
	}

	r, _ := namespaceMap.LoadOrStore(s, &namespace{name: s})
	return r.(*namespace), nil
}

//GetNamespacedKey get the namespaced key
func GetNamespacedKey(n Namespace, k string) (NamespacedKey, error) {
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

//NamespacedKeyFor get the given string
func NamespacedKeyFor(s string) (NamespacedKey, error) {
	v := strings.Split(s, ":")
	if len(v) != 2 {
		return nil, fmt.Errorf("namespace: NamespacedKeyFor:Invalid string provided")
	}
	namespace, err := GetNamespace(v[0])
	if err != nil {
		return nil, err
	}
	key, err := GetNamespacedKey(namespace, v[1])

	return key, err
}

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
