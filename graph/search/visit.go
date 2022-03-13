package search

import "github.com/oligarch316/go-urlrouter/graph"

type VisitorAll[V any] struct{ Results []*graph.SearchResult[V] }

func (va *VisitorAll[V]) VisitSearch(result *graph.SearchResult[V]) bool {
	va.Results = append(va.Results, result)
	return false
}

type VisitorAllPredicate[V any] struct {
	Predicate func(*graph.SearchResult[V]) bool
	Results   []*graph.SearchResult[V]
}

func (vap *VisitorAllPredicate[V]) VisitSearch(result *graph.SearchResult[V]) bool {
	if vap.Predicate(result) {
		vap.Results = append(vap.Results, result)
	}

	return false
}

type VisitorFirst[V any] struct{ Result *graph.SearchResult[V] }

func (vf *VisitorFirst[V]) VisitSearch(result *graph.SearchResult[V]) bool {
	vf.Result = result
	return true
}

type VisitorFirstPredicate[V any] struct {
	Predicate func(*graph.SearchResult[V]) bool
	Result    *graph.SearchResult[V]
}

func (vfp *VisitorFirstPredicate[V]) VisitSearch(result *graph.SearchResult[V]) bool {
	if vfp.Predicate(result) {
		vfp.Result = result
		return true
	}

	return false
}
