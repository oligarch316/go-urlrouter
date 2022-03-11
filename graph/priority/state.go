package priority

import "github.com/oligarch316/go-urlrouter/graph"

type stateAdd[V any] struct {
	parameterKeys []string
	value         V
}

type stateSearch[V any] struct {
	parameterValues []string
	visitor         graph.Searcher[V]
}

type stateWalk[V any] struct {
	visitor graph.Walker[V]
}
