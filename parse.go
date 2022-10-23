package namespace

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

const maxLength = 200
const defaultNamespace = "minecraft"

var translate [unicode.MaxASCII + 1]rune = func() (table [unicode.MaxASCII + 1]rune) {
	for i := rune(0); i <= unicode.MaxASCII; i++ {
		switch {
		case unicode.IsLower(i),
			unicode.IsDigit(i),
			i == '-',
			i == '_':
			table[i] = i

		case unicode.IsUpper(i):
			table[i] = unicode.ToLower(i)
		default:
			table[i] = utf8.RuneError
		}
	}
	return
}()

// parseNSK parses the given string and converts it into a namespaced key.
// If noSeparator if set, this parses the string as a key without a namespace.
// If nsOnly is also set, the string is parsed as a namespace without a key.
// If strict is set, this makes no attempt to correct the nsk and returns at the first
// error it encounters. If strict is false, this never returns an error unless the string is longer
// than [maxLength] or ends with a trailing ':' character.
func parseNSK(v string, strict, noSeparator, nsOnly bool) (ns, key string, err error) {
	var sep int
	var invalid bool

	if len(v) > maxLength {
	}

	if l := len(v); l == 0 {
		// return default namespace for empty string if we only need a namespace
		if nsOnly {
			return defaultNamespace, "", nil
		}

		return "", "", ErrEmpty

	} else if l > maxLength {
		return "", "", ErrTooLong
	} else if v[l-1] == ':' {
		return "", "", ErrTrailingSep
	}

	// there can be no separator if nsOnly is set
	if !noSeparator && nsOnly {
		noSeparator = true
	}

	// fast path. iterates over each character in the string without modifying it.
	var i int
	var char rune
	for i, char = range v {
		if r := translate[char&0x7f]; r == char {
			continue
		}

		// first ':' separates the namespace and the key
		if char == ':' && sep == 0 && !noSeparator {
			sep = i
			continue
		}

		// the `/` character is only allowed in the key.
		if char == '/' {
			if sep != 0 || (noSeparator && !nsOnly) {
				continue
			}

			if !noSeparator && sep == 0 {
				// check if the key contains a separator.
				// we do this to insure bare keys that have `/` are parsed correctly.
				noSeparator = strings.IndexByte(v[i:], ':') == -1
				// the string has no separator, we are parsing a key so `/` is allowed
				if noSeparator {
					continue
				}
			}

		}

		invalid = true
		break
	}

	if invalid {
		if strict {
			return "", "", ErrInvalidChar
		}

		// create a new string builder and copy the string we already validated
		var b strings.Builder
		b.Grow(len(v))
		b.WriteString(v[:i])

		if !noSeparator && sep == 0 {
			// check if the key contains a separator.
			// we do this to insure bare keys that have `/` are parsed correctly.
			noSeparator = strings.IndexByte(v[i:], ':') == -1
		}

		start := i
		for i, char = range v[start:] {
			replaced := translate[char&0x7f]

			// handle replacing the invalid characters
			if char > unicode.MaxASCII || replaced == utf8.RuneError {
				switch char {
				case ':':
					replaced = '_'

					// if this is the first ':' we encountered
					// it is valid and separates the namespace and key
					if sep == 0 && !noSeparator {
						replaced = ':'
						sep = i + start
					}

				case '/':
					replaced = '_'

					// `/` is allowed in keys
					if sep != 0 || noSeparator && !nsOnly {
						replaced = '/'
					}

				default:
					replaced = '_'
				}
			}

			b.WriteRune(replaced)
		}
		v = b.String()
	}

	if sep != 0 {
		if len(v)-sep > 1 {
			return v[:sep], v[sep+1:], nil
		}

		// this only happens when a string that ends with `:` is given.
		// we already checked this at the start, so this should be unreachable.
		return "", "", ErrTrailingSep
	}

	if nsOnly {
		return v, "", nil
	}

	return defaultNamespace, v, nil
}
