package log

import "github.com/oligarch316/go-urlrouter/graph"

type Path[V any] struct {
	Keys  []graph.Key
	Value V
}

func (p Path[_]) String() string { return graph.FormatPath(p.Value, p.Keys...) }

type Tree[V any] struct {
	tree graph.Tree[Path[V]]
}

func (t *Tree[V]) Paths() []Path[V] { return t.tree.Values() }

func (t *Tree[V]) Add(val V, keys ...graph.Key) error {
	path := Path[V]{Keys: keys, Value: val}
	return t.tree.Add(path, keys...)
}

func (t *Tree[V]) Search(segs ...graph.Segment) *graph.Result[V] {
	if pathResult := t.tree.Search(segs...); pathResult != nil {
		return &graph.Result[V]{
			Parameters: pathResult.Parameters,
			Tail:       pathResult.Tail,
			Value:      pathResult.Value.Value,
		}
	}

	return nil
}

func (t *Tree[V]) Values() []V {
	var res []V

	for _, pathVal := range t.tree.Values() {
		res = append(res, pathVal.Value)
	}

	return res
}
