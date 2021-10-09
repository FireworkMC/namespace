package namespace

import (
	"strings"
	"unicode"
)

// See https://www.minecraft.net/en-us/article/minecraft-snapshot-17w43a

// ParseNamespacedKey parse the given namespaced key
func ParseNamespacedKey(s string) (nsk [2]string, valid bool) {
	v := strings.Split(strings.ToLower(s), ":")

	if len(v) != 2 {
		return
	}

	for i, s := range v {
		w := strings.TrimSpace(s)
		if len(w) == 0 {
			return
		}
		v[i] = w
	}

	if !IsValidNamespace(v[0]) || !IsValidKey(v[1]) {
		return
	}

	return [2]string{v[0], v[1]}, true

}

// IsValidNamespace returns if the given namspace is valid
func IsValidNamespace(s string) bool {
	for _, b := range []rune(s) {
		if !(unicode.IsLower(b) || unicode.IsDigit(b) || b == '-' || b == '_') {
			return false
		}
	}
	return true
}

// IsValidKey returns if the namespaced key is valid
func IsValidKey(s string) bool {
	for _, b := range []rune(s) {
		if !(unicode.IsLower(b) || unicode.IsDigit(b) || b == '-' || b == '_' || b == '/' || b == '.') {
			return false
		}
	}
	return true
}
