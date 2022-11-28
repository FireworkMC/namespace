package namespace

import (
	"errors"
	"strings"
	"testing"

	"github.com/yehan2002/is/v2"
)

type simpleKeyTest struct{}

func TestSimpleKey(t *testing.T) { is.Suite(t, &simpleKeyTest{}) }

func (s *simpleKeyTest) TestParse(is is.Is) {
	test := func(inp string, strict, noSeparator, nsOnly bool, ns, key string, err error) {
		is.T().Helper()
		ns2, k2, err2 := parseNSK(inp, strict, noSeparator, nsOnly)
		is(errors.Is(err, err2), "Expected error to be %v got %v", err, err2)
		is(ns2 == ns, "Expected ns to be %#v got %#v", ns, ns2)
		is(k2 == key, "Expected key to be %d got %d", key, k2)
	}

	test("minecraft:air", false, false, false, "minecraft", "air", nil)
	test("minecraft:blocks/air", false, false, false, "minecraft", "blocks/air", nil)
	test("minecraft:blocks/air.2", false, false, false, "minecraft", "blocks/air.2", nil)
	test("minecraft:AIR", false, false, false, "minecraft", "air", nil)
	test("abc:;;;123", false, false, false, "abc", "___123", nil)
	test("a;bc:a", false, false, false, "a_bc", "a", nil)
	test("a;bc:a/a", false, false, false, "a_bc", "a/a", nil)
	test("a/bc:a", false, false, false, "a_bc", "a", nil)
	test("a;:v/a", false, false, false, "a_", "v/a", nil)

	test("aa:aa", false, true, true, "aa_aa", "", nil)
	test("aa:aa", false, false, true, "aa_aa", "", nil)
	test("aa:aa", false, true, false, "minecraft", "aa_aa", nil)
	test("aa/aa", false, true, false, "minecraft", "aa/aa", nil)
	test("aa/aa", false, false, false, "minecraft", "aa/aa", nil)
	test("aa.aa", false, false, false, "minecraft", "aa.aa", nil)
	test("a/a:b", false, false, false, "a_a", "b", nil)
	test("a.a:b", false, false, false, "a_a", "b", nil)

	test("aa:", false, false, false, "", "", ErrTrailingSep)
	test("aa:", true, false, false, "", "", ErrTrailingSep)

	test("", true, false, false, "", "", ErrEmpty)
	test("", true, false, true, defaultNamespace, "", nil)

	test("a/a:a", true, false, false, "", "", ErrInvalidChar)
	test(strings.Repeat("a", maxLength+1), true, false, false, "", "", ErrTooLong)
}
