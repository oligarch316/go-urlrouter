package memoized

import "github.com/oligarch316/go-urlrouter/graph/priority"

type Memo[V any] struct {
	Path  []priority.Key
	Value V
}

func (m Memo[_]) String() string { return priority.FormatPath(m.Value, m.Path...) }

type Tree[V any] struct {
	WithMemo priority.Tree[Memo[V]]
}

func (t *Tree[V]) Add(value V, path ...priority.Key) error {
	var (
		memo = Memo[V]{Path: path, Value: value}
		err  = t.WithMemo.Add(memo, path...)
	)

	if dupErr, ok := err.(priority.DuplicateValueError[Memo[V]]); ok {
		err = priority.DuplicateValueError[V]{ExistingValue: dupErr.ExistingValue.Value}
	}

	return err
}

func (t *Tree[V]) Search(searcher priority.Searcher[V], query ...string) bool {
	memoSearcher := func(memoResult *priority.Result[Memo[V]]) bool {
		return searcher.VisitSearch(&priority.Result[V]{
			Parameters: memoResult.Parameters,
			Tail:       memoResult.Tail,
			Value:      memoResult.Value.Value,
		})
	}

	return t.WithMemo.SearchFunc(memoSearcher, query...)
}

func (t *Tree[V]) SearchFunc(searcher func(result *priority.Result[V]) (done bool), query ...string) bool {
	return t.SearchFunc(priority.SearcherFunc[V](searcher), query...)
}

func (t *Tree[V]) Walk(walker priority.Walker[V]) bool {
	memoWalker := func(memo Memo[V]) bool {
		return walker.VisitWalk(memo.Value)
	}

	return t.WithMemo.WalkFunc(memoWalker)
}

func (t *Tree[V]) WalkFunc(walker func(value V) (done bool)) bool {
	return t.WalkFunc(priority.WalkerFunc[V](walker))
}
