package priority

type stateAdd[V any] struct {
	parameterKeys []string
	value         V
}

type stateSearch[V any] struct {
	parameterValues []string
	visitor         Searcher[V]
}

type stateWalk[V any] struct {
	visitor Walker[V]
}
