package memoized

import "github.com/oligarch316/go-urlrouter/graph"

type Memo[V any] struct {
	Path  []graph.Key
	Value V
}

func (m Memo[_]) String() string { return graph.FormatPath(m.Value, m.Path...) }

type Tree[V any] struct {
	WithMemo graph.Tree[Memo[V]]
}

func (t *Tree[V]) Add(value V, path ...graph.Key) error {
	var (
		memo = Memo[V]{Path: path, Value: value}
		err  = t.WithMemo.Add(memo, path...)
	)

	if dupErr, ok := err.(graph.DuplicateValueError[Memo[V]]); ok {
		err = graph.DuplicateValueError[V]{ExistingValue: dupErr.ExistingValue.Value}
	}

	return err
}

func (t *Tree[V]) Search(searcher graph.Searcher[V], query ...string) bool {
	memoSearcher := func(memoResult *graph.Result[Memo[V]]) bool {
		return searcher.VisitSearch(&graph.Result[V]{
			Parameters: memoResult.Parameters,
			Tail:       memoResult.Tail,
			Value:      memoResult.Value.Value,
		})
	}

	return t.WithMemo.SearchFunc(memoSearcher, query...)
}

func (t *Tree[V]) SearchFunc(searcher func(result *graph.Result[V]) (done bool), query ...string) bool {
	return t.SearchFunc(graph.SearcherFunc[V](searcher), query...)
}

func (t *Tree[V]) Walk(walker graph.Walker[V]) bool {
	memoWalker := func(memo Memo[V]) bool {
		return walker.VisitWalk(memo.Value)
	}

	return t.WithMemo.WalkFunc(memoWalker)
}

func (t *Tree[V]) WalkFunc(walker func(value V) (done bool)) bool {
	return t.WalkFunc(graph.WalkerFunc[V](walker))
}
