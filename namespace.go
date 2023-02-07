package namespace

import (
	"encoding"

	"github.com/yehan2002/errors"
)

// See https://www.minecraft.net/en-us/article/minecraft-snapshot-17w43a

// Default this is the namespace used for minecraft items/effects ect.
// Plugins should use their own namespace if they aren't overriding vanilla features
var Default NS = Namespace("minecraft")

const (
	// ErrEmpty the namespaced key is empty.
	ErrEmpty = errors.Const("namespace: namespaced key is empty")
	// ErrNil the namespace or key is nil.
	// This is returned by methods on the zero value of [NS] or [NSK]
	ErrNil = errors.Const("namespace: nil namespace or key")
	// ErrTooLong the namespaced key exceeded the max allowed length.
	ErrTooLong = errors.Const("namespace: namespaced key is too long")
	// ErrInvalidChar the namespaced key contained an invalid character
	ErrInvalidChar = errors.Const("namespace: namespaced key contains illegal characters")
	// ErrTrailingSep the nsk contained a trailing `:`
	ErrTrailingSep = errors.Const("namespace: namespace contains trailing ':' character")
)

var (
	_ encoding.TextMarshaler   = (*NS)(nil)
	_ encoding.TextUnmarshaler = (*NS)(nil)
	_ encoding.TextMarshaler   = (*NSK)(nil)
	_ encoding.TextUnmarshaler = (*NSK)(nil)
)

var namespaces = syncMap[string, *ns]{
	v: map[string]*ns{},

	New: func(n string) *ns {
		ns := &ns{name: n, keys: &syncMap[string, *nsk]{v: map[string]*nsk{}}}
		ns.keys.New = func(k string) *nsk {
			return &nsk{ns: NS{ns}, key: k, full: ns.name + ":" + k}
		}

		return ns
	},
}

// NS is a namespace.
// A namespace can only contain digits, lowercase letters, underscores and hyphens.
//
// To compare namespaces use `==` operator directly (not on the pointer value) or use the [NS.Equal] method.
type NS struct{ ns *ns }

type ns struct {
	name string
	keys *syncMap[string, *nsk]
}

// Equal returns if this namespace and the given namespace are equal.
func (n NS) Equal(n2 NS) bool { return n.ns == n2.ns }

// IsNil returns if this nsk is nil.
// if this returns true, calling [NS.Key] will panic.
func (n NS) IsNil() bool { return n.ns == nil }

// Key creates a new key inside this namespace.
// All invalid characters in the will be replaced with an underscore.
//
// This panics if the length of the key is larger than [MaxLength] or is zero.
// If [NS.IsNil] returns true, this will panic.
func (n NS) Key(k string) NSK {
	if n.ns == nil {
		panic(ErrNil)
	}

	if nsk, ok := n.ns.keys.Get(k); ok {
		return NSK{nsk}
	}

	_, k, _ = parseNSK(k, false, true, false)
	return NSK{n.ns.keys.GetOrCreate(k)}
}

// ParseKey creates a new key inside this namespace.
// Unlike [NS.Key], this will return an error if any invalid characters are encountered.
//
// This returns an error if the length of the key is larger than [MaxLength] or is zero
// If [NS.IsNil] returns true, this will return [ErrNil].
func (n NS) ParseKey(k string) (nsk NSK, err error) {
	if n.ns == nil {
		return NSK{}, ErrNil
	}

	if nsk, ok := n.ns.keys.Get(k); ok {
		return NSK{nsk}, nil
	}

	_, k, err = parseNSK(k, true, true, false)
	if err != nil {
		return NSK{}, err
	}

	return NSK{n.ns.keys.GetOrCreate(k)}, nil
}

// MarshalText implements encoding.TextMarshaler
func (n *NS) MarshalText() (text []byte, err error) { return []byte(n.String()), nil }

// UnmarshalText implements encoding.TextUnmarshaler
func (n *NS) UnmarshalText(text []byte) (err error) {
	*n, err = ParseNamespace(string(text), true)
	return
}

func (n NS) String() string {
	if n.ns == nil {
		return ""
	}

	return n.ns.name
}

// NSK is a namespaced key.
// A namespace can only contain digits, lowercase letters, underscores, hyphens, forward slash and dots.
//
// To compare namespaces use `==` operator directly (not on the pointer value) or use the [NSK.Equal] method.
type NSK struct{ nsk *nsk }

type nsk struct {
	ns        NS
	key, full string
}

// Equal returns if this namespaced key is equal to the given namespaced key.
func (n NSK) Equal(n2 NSK) bool { return n.nsk == n2.nsk }

// IsNil returns if this nsk is nil.
// if this returns true, calling [NSK.Namespace] will panic.
func (n NSK) IsNil() bool { return n.nsk == nil }

// Namespace gets the namespace for this key.
// If [NSK.IsNil] returns true, this will panic.
func (n NSK) Namespace() NS {
	if n.nsk == nil {
		panic(ErrNil)
	}

	return n.nsk.ns
}

// Key gets the key part of the namespaced key (the part after the ':').
// If the namespaced key is nil, this returns an empty string.
func (n NSK) Key() string {
	if n.nsk == nil {
		return ""
	}

	return n.nsk.key
}

// MarshalText implements encoding.TextMarshaler
func (n *NSK) MarshalText() (text []byte, err error) { return []byte(n.String()), nil }

// UnmarshalText implements encoding.TextUnmarshaler
func (n *NSK) UnmarshalText(text []byte) (err error) {
	*n, err = ParseKey(string(text), true)
	return
}

func (n NSK) String() string {
	if n.nsk == nil {
		return ""
	}
	return n.nsk.full
}

// Namespace creates a new namespace from the given string.
//
// This panics if the length of the namespace is larger than [MaxLength] or is zero.
// All invalid characters in the namespace are replaced with underscores.
func Namespace(namespace string) NS {
	ns, err := ParseNamespace(namespace, false)
	if err != nil {
		panic(err)
	}
	return ns
}

// ParseNamespace creates a new namespace from the given string.
//
// This returns an error if the length of the namespace is larger than [MaxLength] or is zero.
// If strict mode is enabled, this returns an error when it encounters an invalid character.
// Otherwise, this replaces all invalid characters with underscores.
func ParseNamespace(v string, strict bool) (NS, error) {
	if ns, ok := namespaces.Get(v); ok {
		return NS{ns}, nil
	}

	ns, _, err := parseNSK(v, strict, true, true)
	if err != nil {
		return NS{}, err
	}

	return NS{namespaces.GetOrCreate(ns)}, nil
}

// Key creates a new namespaced key from the given string.
//
// This panics if the length of the key is larger than [MaxLength] or is zero.
// All invalid characters in the namespaced key are replaced with underscores.
func Key(namespacedKey string) NSK {
	nsk, err := ParseKey(namespacedKey, false)
	if err != nil {
		panic(err)
	}
	return nsk
}

// ParseKey creates a new namespaced key from the given string.
//
// This returns an error if the length of the key is larger than [MaxLength] or is zero.
// If strict mode is enabled, this returns an error when it encounters an invalid character.
// Otherwise, this replaces all invalid characters with underscores.
func ParseKey(namespacedKey string, strict bool) (NSK, error) {
	ns, k, err := parseNSK(namespacedKey, strict, false, false)
	if err != nil {
		return NSK{}, err
	}
	return NSK{namespaces.GetOrCreate(ns).keys.GetOrCreate(k)}, nil
}
