package priority

import "github.com/oligarch316/go-urlrouter/graph"

type Tree[V any] struct{ root nodeConstant[V] }

func (t *Tree[V]) Add(value V, path ...graph.Key) error {
	return t.root.add(path, stateAdd[V]{value: value})
}

func (t Tree[V]) Search(searcher graph.Searcher[V], query ...string) {
	t.root.search(query, stateSearch[V]{visitor: searcher})
}

func (t Tree[V]) SearchFunc(visitor func(result *graph.SearchResult[V]) (done bool), query ...string) {
	t.Search(graph.SearcherFunc[V](visitor), query...)
}

func (t Tree[V]) Walk(walker graph.Walker[V]) {
	t.root.walk(stateWalk[V]{visitor: walker})
}

func (t Tree[V]) WalkFunc(walker func(value V) (done bool)) {
	t.Walk(graph.WalkerFunc[V](walker))
}
