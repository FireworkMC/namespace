package namespace

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/yehan2002/errors"
)

const maxLength = 200
const defaultNamespace = "minecraft"
const validChars = "abcdefghijklmnopqrstuvwxyz1234567890-_"
const upperChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

var translate [unicode.MaxASCII + 1]rune = func() (table [unicode.MaxASCII + 1]rune) {
	for _, char := range validChars {
		table[char] = char
	}

	for _, char := range upperChars {
		table[char] = unicode.ToLower(char)
	}

	for i := range table {
		if table[i] == 0 {
			table[i] = utf8.RuneError
		}
	}

	return
}()

// validBytes is a array of bytes in which validBytes[r] == r for all valid
// bytes in a namespace/namespaced key (except . and /).
var validBytes [256]byte = func() (table [256]byte) {
	for _, char := range []byte(validChars) {
		table[char] = char
	}

	table[0] = 255

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
	var charByte byte
	for i, charByte = range []byte(v) {
		switch charByte {
		case validBytes[charByte]:
			// char is a valid character.
			continue
		case ':':
			// first ':' separates the namespace and the key
			if sep == 0 && !noSeparator {
				sep = i
				continue
			}
		case '/', '.':
			// the `/` and `.` characters are only allowed in the key.
			if sep != 0 || (noSeparator && !nsOnly) {
				continue
			}

			if !noSeparator && sep == 0 {
				// check if the key contains a separator.
				// we do this to insure bare keys that have `/` or `.` are parsed correctly.
				noSeparator = strings.IndexByte(v[i:], ':') == -1
				// the string has no separator, we are parsing a key so `/` or `.` is allowed
				if noSeparator {
					continue
				}
			}
		}

		invalid = true
		break
	}

	var char rune
	if invalid {
		if strict {
			return "", "", getCharError(v, i, noSeparator)
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
				replaced = '_'

				switch char {
				case ':':
					// if this is the first ':' we encountered
					// it is valid and separates the namespace and key
					if sep == 0 && !noSeparator {
						replaced = ':'
						sep = i + start
					}

				case '/', '.':
					// `/` is allowed in keys
					if sep != 0 || noSeparator && !nsOnly {
						replaced = char
					}
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

// getCharError returns [ErrInvalidChar] wrapped with more information about the error that
// occurred.
func getCharError(v string, i int, inNamespace bool) error {
	invalidChar, _ := utf8.DecodeRuneInString(v[i:])
	var msg string

	if invalidChar == '/' || invalidChar == '.' {
		msg = fmt.Sprintf(`%q is not allowed in a key`, invalidChar)
	} else if invalidChar == ':' {
		if inNamespace {
			msg = `":" is not allowed in a namespace`
		} else {
			msg = `found multiple ":" characters`
		}
	} else {
		msg = fmt.Sprintf("invalid character %q (%c)", invalidChar, invalidChar)
	}

	return errors.CauseStr(ErrInvalidChar, msg)
}
