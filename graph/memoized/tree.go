package memoized

import (
	"github.com/oligarch316/go-urlrouter/graph"
	"github.com/oligarch316/go-urlrouter/graph/priority"
)

type Memo[V any] struct {
	Path  []graph.Key
	Value V
}

func (m Memo[_]) String() string { return graph.FormatPath(m.Value, m.Path...) }

func wrapSearcher[V any](searcher graph.Searcher[V]) graph.Searcher[Memo[V]] {
	wrapped := func(memoResult *graph.SearchResult[Memo[V]]) bool {
		return searcher.VisitSearch(&graph.SearchResult[V]{
			Parameters: memoResult.Parameters,
			Tail:       memoResult.Tail,
			Value:      memoResult.Value.Value,
		})
	}

	return graph.SearcherFunc[Memo[V]](wrapped)
}

func wrapWalker[V any](walker graph.Walker[V]) graph.Walker[Memo[V]] {
	wrapped := func(memo Memo[V]) bool {
		return walker.VisitWalk(memo.Value)
	}

	return graph.WalkerFunc[Memo[V]](wrapped)
}

type Tree[V any] struct{ Memoized priority.Tree[Memo[V]] }

func (t *Tree[V]) Add(value V, path ...graph.Key) error {
	var (
		memo = Memo[V]{Path: path, Value: value}
		err  = t.Memoized.Add(memo, path...)
	)

	if dupErr, ok := err.(graph.DuplicateValueError[Memo[V]]); ok {
		err = graph.DuplicateValueError[V]{ExistingValue: dupErr.ExistingValue.Value}
	}

	return err
}

func (t Tree[V]) Search(searcher graph.Searcher[V], query ...string) {
	t.Memoized.Search(wrapSearcher(searcher), query...)
}

func (t Tree[V]) SearchFunc(searcher func(result *graph.SearchResult[V]) (done bool), query ...string) {
	t.SearchFunc(graph.SearcherFunc[V](searcher), query...)
}

func (t Tree[V]) Walk(walker graph.Walker[V]) {
	t.Memoized.Walk(wrapWalker(walker))
}

func (t Tree[V]) WalkFunc(walker func(value V) (done bool)) {
	t.WalkFunc(graph.WalkerFunc[V](walker))
}
