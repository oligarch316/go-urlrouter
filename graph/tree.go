package graph

type Tree struct{ root nodeConstant }

func (t *Tree) Add(val Value, keys ...Key) error {
	return t.root.add(keys, &nodeValue{value: val})
}

func (t *Tree) Search(segs ...Segment) *Result {
	return t.root.search(segs, nil)
}
