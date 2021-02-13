package namespace

type namespace struct {
	name string
}

func (n *namespace) String() string { return n.name }

func (n *namespace) Key(k string) (Key, error) { return GetKey(k) }

func (n *namespace) MustKey(k string) Key { return MustKey(k) }

func (n *namespace) namespace() *namespace { return n }

type namespacedKey struct {
	namespace *namespace
	key       string
}

func (n *namespacedKey) String() string { return n.namespace.name + ":" + n.key }

func (n *namespacedKey) Namespace() Namespace { return n.namespace }

func (n *namespacedKey) Key() string { return n.key }

func (n *namespacedKey) namespacedKey() *namespacedKey { return n }
