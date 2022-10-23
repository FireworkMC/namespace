package namespace

import (
	"github.com/yehan2002/errors"
)

// See https://www.minecraft.net/en-us/article/minecraft-snapshot-17w43a

// Default this is the namespace used for minecraft items/effects ect.
// Plugins should use their own namespace if they aren't overriding vanilla features
var Default NS = Namespace("minecraft")

const (
	// ErrEmpty the namespaced key is empty.
	ErrEmpty = errors.Const("namespace: namespaced key is empty")
	// ErrTooLong the namespaced key exceeded the max allowed length.
	ErrTooLong = errors.Const("namespace: namespaced key is too long")
	// ErrInvalidChar the namespaced key contained an invalid character
	ErrInvalidChar = errors.Const("namespace: namespaced key contains illegal characters")
	// ErrTrailingSep the nsk contained a trailing `:`
	ErrTrailingSep = errors.Const("namespace: namespace contains trailing ':' character")
)

var namespaces = syncMap[string, *ns]{
	v:   map[string]*ns{},
	New: newNamespace,
}

// NS a namespace.
type NS interface {
	// String gets the namespace as a string
	String() string

	// Key creates a new key inside this namespace.
	// This panics if the length of the key is larger than `maxLength`
	Key(k string) NSK

	// ParseKey parses the given string and returns a key if it is a valid key.
	ParseKey(k string) (key NSK, err error)
}

// NSK a namespaced key
type NSK interface {
	Namespace() NS
	String() string
	Key() string
}

// Namespace creates a new namespace from the given string.
func Namespace(v string) NS {
	ns, _, _ := parseNSK(v, false, true, true)
	return namespaces.Get(ns)
}

// ParseNamespace creates a new namespace
func ParseNamespace(v string) (NS, error) {
	ns, _, err := parseNSK(v, true, true, true)
	if err != nil {
		return nil, err
	}

	return namespaces.Get(ns), nil
}

// Key creates a new key
func Key(v string) NSK {
	ns, k, _ := parseNSK(v, false, false, false)
	return namespaces.Get(ns).keys.Get(k)
}

// ParseKey parses a key
func ParseKey(v string) (NSK, error) {
	ns, k, err := parseNSK(v, true, false, false)
	if err != nil {
		return nil, err
	}
	return namespaces.Get(ns).keys.Get(k), nil
}
