package graph

import "errors"

var ErrEmptyKeys = errors.New("empty keys")

type Tree struct{ root nodeConstant }

func (t *Tree) Add(val Value, keys ...Key) error {
	if len(keys) < 1 {
		return ErrEmptyKeys
	}

	return t.root.add(keys, &nodeValue{value: val})
}

func (t *Tree) Search(segs ...Segment) *Result {
	return t.root.search(segs, nil)
}
