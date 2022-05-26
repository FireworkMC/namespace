package namespace

type namespace struct {
	name string
	keys *syncMap[string, *key]
}

func newNamespace(n string) *namespace {
	ns := &namespace{name: n, keys: &syncMap[string, *key]{v: map[string]*key{}}}
	ns.keys.New = func(k string) *key {
		return &key{ns: ns, key: k, full: ns.name + ":" + k}
	}

	return ns
}

func (n *namespace) String() string { return n.name }

func (n *namespace) Key(k string) NSK {
	_, k, _ = parseNSK(k, false, true, false)
	return n.keys.Get(k)
}

func (n *namespace) ParseKey(k string) (nsk NSK, err error) {
	_, k, err = parseNSK(k, true, true, false)
	if err != nil {
		return nil, err
	}
	return n.keys.Get(k), nil
}

type key struct {
	ns        *namespace
	key, full string
}

func (n *key) Namespace() NS { return n.ns }

func (n *key) String() string { return n.full }

func (n *key) Key() string { return n.key }
