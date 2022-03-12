package walk

import "github.com/oligarch316/go-urlrouter/graph"

type Walkable[V any] interface{ Walk(graph.Walker[V]) }

func All[V any](tree Walkable[V]) []V {
	visitor := new(VisitorAll[V])
	tree.Walk(visitor)
	return visitor.Values
}

func AllPredicate[V any](tree Walkable[V], predicate func(V) bool) []V {
	visitor := &VisitorAllPredicate[V]{Predicate: predicate}
	tree.Walk(visitor)
	return visitor.Values
}
