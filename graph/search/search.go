package search

import "github.com/oligarch316/go-urlrouter/graph"

type Searchable[V any] interface {
	Search(graph.Searcher[V], ...string) bool
}

func All[V any](tree Searchable[V], query ...string) []*graph.SearchResult[V] {
	visitor := new(VisitorAll[V])
	tree.Search(visitor, query...)
	return visitor.Results
}

func AllPredicate[V any](tree Searchable[V], predicate func(*graph.SearchResult[V]) bool, query ...string) []*graph.SearchResult[V] {
	visitor := &VisitorAllPredicate[V]{Predicate: predicate}
	tree.Search(visitor, query...)
	return visitor.Results
}

func First[V any](tree Searchable[V], query ...string) *graph.SearchResult[V] {
	visitor := new(VisitorFirst[V])
	tree.Search(visitor, query...)
	return visitor.Result
}

func FirstPredicate[V any](tree Searchable[V], predicate func(*graph.SearchResult[V]) bool, query ...string) *graph.SearchResult[V] {
	visitor := &VisitorFirstPredicate[V]{Predicate: predicate}
	tree.Search(visitor, query...)
	return visitor.Result
}
