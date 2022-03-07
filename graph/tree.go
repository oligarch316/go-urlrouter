package graph

type Tree[V any] struct{ root nodeConstant[V] }

func (t *Tree[V]) Add(val V, keys ...Key) error {
	return t.root.add(keys, &nodeValue[V]{value: val})
}

func (t *Tree[V]) Search(segs ...string) *Result[V] {
	return t.root.search(segs, nil)
}

func (t *Tree[V]) Values() []V {
	return t.root.values()
}
