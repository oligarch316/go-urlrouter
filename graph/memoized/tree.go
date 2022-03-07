package memoized

import "github.com/oligarch316/go-urlrouter/graph"

type Memo[V any] struct {
	Keys  []graph.Key
	Value V
}

func (m Memo[_]) String() string { return graph.FormatPath(m.Value, m.Keys...) }

type Tree[V any] struct {
	tree graph.Tree[Memo[V]]
}

func (t *Tree[V]) Memos() []Memo[V] { return t.tree.Values() }

func (t *Tree[V]) Add(val V, keys ...graph.Key) error {
	memo := Memo[V]{Keys: keys, Value: val}
	return t.tree.Add(memo, keys...)
}

func (t *Tree[V]) Search(segs ...graph.Segment) *graph.Result[V] {
	if memoResult := t.tree.Search(segs...); memoResult != nil {
		return &graph.Result[V]{
			Parameters: memoResult.Parameters,
			Tail:       memoResult.Tail,
			Value:      memoResult.Value.Value,
		}
	}

	return nil
}

func (t *Tree[V]) Values() []V {
	var res []V

	for _, memoVal := range t.tree.Values() {
		res = append(res, memoVal.Value)
	}

	return res
}
