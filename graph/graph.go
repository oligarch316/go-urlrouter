package graph

type SearchResult[V any] struct {
	Parameters map[string]string
	Tail       []string
	Value      V
}

type Searcher[V any] interface {
	VisitSearch(result *SearchResult[V]) (done bool)
}

type SearcherFunc[V any] func(result *SearchResult[V]) (done bool)

func (sf SearcherFunc[V]) VisitSearch(result *SearchResult[V]) bool { return sf(result) }

type Walker[V any] interface {
	VisitWalk(value V) (done bool)
}

type WalkerFunc[V any] func(value V) (done bool)

func (wf WalkerFunc[V]) VisitWalk(value V) bool { return wf(value) }

type Tree[V any] interface {
	Add(V, ...Key) error
	Search(Searcher[V], ...string) bool
	Walk(Walker[V]) bool
}
