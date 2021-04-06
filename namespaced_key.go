package namespace

type namespace struct {
	name string
}

func (n *namespace) String() string { return n.name }

func (n *namespace) Key(k string) (Key, error) { return GetKey(n.name + ":" + k) }

func (n *namespace) MustKey(k string) Key { return MustKey(n.name + ":" + k) }

func (n *namespace) namespace() *namespace { return n }

type namespacedKey struct {
	namespace *namespace
	key       string
	full      string
	valid     bool
}

func (n *namespacedKey) String() string { return n.full }

func (n *namespacedKey) Namespace() Namespace { return n.namespace }

func (n *namespacedKey) Key() string { return n.key }

func (n *namespacedKey) Valid() bool { return n.valid }

func (n *namespacedKey) namespacedKey() *namespacedKey { return n }
