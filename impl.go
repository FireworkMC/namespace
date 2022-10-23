package namespace

type ns struct {
	name string
	keys *syncMap[string, *nsk]
}

func newNamespace(n string) *ns {
	ns := &ns{name: n, keys: &syncMap[string, *nsk]{v: map[string]*nsk{}}}
	ns.keys.New = func(k string) *nsk {
		return &nsk{ns: ns, key: k, full: ns.name + ":" + k}
	}

	return ns
}

func (n *ns) String() string { return n.name }

func (n *ns) Key(k string) NSK {
	_, k, _ = parseNSK(k, false, true, false)
	return n.keys.Get(k)
}

func (n *ns) ParseKey(k string) (nsk NSK, err error) {
	_, k, err = parseNSK(k, true, true, false)
	if err != nil {
		return nil, err
	}
	return n.keys.Get(k), nil
}

type nsk struct {
	ns        *ns
	key, full string
}

func (n *nsk) Namespace() NS { return n.ns }

func (n *nsk) String() string { return n.full }

func (n *nsk) Key() string { return n.key }
