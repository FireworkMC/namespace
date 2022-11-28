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

// NS a namespace.
// A namespace can only contain digits, lowercase letters, underscores and hyphens.
// To compare namespaces use `==` operator directly (not on the pointer value) or use the Equal method.
type NS struct{ ns *ns }

type ns struct {
	name string
	keys *syncMap[string, *nsk]
}

// Equal returns if `n` is equal to `n2`
func (n NS) Equal(n2 NS) bool { return n.ns == n2.ns }

// IsNil returns if this nsk is nil.
// if this returns true, calling [NS.Key] will panic.
func (n NS) IsNil() bool { return n.ns == nil }

// Key creates a new key inside this namespace.
// This panics if the length of the key is larger than `maxLength`.
// If [NS.IsNil] returns true, this will panic.
func (n NS) Key(k string) NSK {
	if n.ns == nil {
		panic(ErrNil)
	}

	_, k, _ = parseNSK(k, false, true, false)
	return NSK{n.ns.keys.Get(k)}
}

// ParseKey parses the given string and returns a key if it is a valid key.
// If the namespace is nil, the default namespace will be used.
func (n NS) ParseKey(k string) (nsk NSK, err error) {
	if n.ns == nil {
		return NSK{}, ErrNil
	}

	_, k, err = parseNSK(k, true, true, false)
	if err != nil {
		return NSK{}, err
	}

	return NSK{n.ns.keys.Get(k)}, nil
}

// MarshalText implements encoding.TextMarshaler
func (n *NS) MarshalText() (text []byte, err error) { return []byte(n.String()), nil }

// UnmarshalText implements encoding.TextUnmarshaler
func (n *NS) UnmarshalText(text []byte) (err error) {
	*n, err = ParseNamespace(string(text))
	return
}

func (n NS) String() string {
	if n.ns == nil {
		return ""
	}

	return n.ns.name
}

// NSK a namespaced key.
// A namespace can only contain digits, lowercase letters, underscores, hyphens, forward slash and dots.
// To compare namespaces use `==` operator directly (not on the pointer value) or use the Equal method.
type NSK struct{ nsk *nsk }

type nsk struct {
	ns        NS
	key, full string
}

// Equal returns if `n` is equal to `n2`
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

// Key gets the key part of the namespaced key (the part after the ':')
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
	*n, err = ParseKey(string(text))
	return
}

func (n NSK) String() string {
	if n.nsk == nil {
		return ""
	}
	return n.nsk.full
}

// Namespace creates a new namespace from the given string.
// This panics if `v` is longer than 200 characters.
func Namespace(v string) NS {
	ns, _, err := parseNSK(v, false, true, true)
	if err != nil {
		panic(err)
	}
	return NS{namespaces.Get(ns)}
}

// ParseNamespace creates a new namespace
func ParseNamespace(v string) (NS, error) {
	ns, _, err := parseNSK(v, true, true, true)
	if err != nil {
		return NS{}, err
	}

	return NS{namespaces.Get(ns)}, nil
}

// Key creates a new key.
// This panics if len(v) > 200 or len(v) == 0.
func Key(v string) NSK {
	ns, k, err := parseNSK(v, false, false, false)
	if err != nil {
		panic(err)
	}
	return NSK{namespaces.Get(ns).keys.Get(k)}
}

// ParseKey parses a key
func ParseKey(v string) (NSK, error) {
	ns, k, err := parseNSK(v, true, false, false)
	if err != nil {
		return NSK{}, err
	}
	return NSK{namespaces.Get(ns).keys.Get(k)}, nil
}
