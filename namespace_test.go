package namespace

import (
	"testing"

	"github.com/yehan2002/is/v2"
)

type namespaceTest struct{}

func TestNamespace(t *testing.T) { is.Suite(t, &namespaceTest{}) }

func (n *namespaceTest) TestNamespace(is is.Is) {
	_, err := ParseNamespace("a/b", true)
	is.Err(err, ErrInvalidChar, "Expected ParseNamespace to return an error")
	_, err = ParseKey("a/b:c", true)
	is.Err(err, ErrInvalidChar, "Expected ParseKey to return an error")

	ns, err := ParseNamespace("minecraft", true)
	is.Err(err, nil, "Expected error to be nil")

	k1, k2 := Default.Key("air"), Key("minecraft:air")
	is(k1 == k2, "keys should be equal")

	k3, err := ParseKey("minecraft:air", true)
	is.Err(err, nil, "Expected error to be nil")
	is(k1 == k3, "keys should be equal")

	k4, err := ns.ParseKey("air")
	is.Err(err, nil, "Expected error to be nil")
	is(k1 == k4, "keys should be equal")

	is(ns == k1.Namespace(), "Expected namespace to be equal")
	is(ns == k2.Namespace(), "Expected namespace to be equal")
	is(ns == k3.Namespace(), "Expected namespace to be equal")

}

func BenchmarkParseNSK(b *testing.B) {
	k := "namespace_name:a_key_with_a_really_long_name"
	b.ReportAllocs()
	b.SetBytes(int64(len(k)))
	for i := 0; i < b.N; i++ {
		parseNSK(k, false, false, false)
	}
}
