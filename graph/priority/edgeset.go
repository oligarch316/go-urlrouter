package graph

import "sort"

type (
	DuplicateValueError[V any] struct{ ExistingValue V }
	InvalidContinuationError   struct{ Continuation []Key }
)

func (dve DuplicateValueError[_]) Error() string   { return "dupliate value" }
func (ice InvalidContinuationError) Error() string { return "invalid continuation" }

type edgeSetTerminal[V any] struct{ node *nodeValue[V] }

func (est *edgeSetTerminal[V]) add(state stateAdd[V]) error {
	if est.node != nil {
		return DuplicateValueError[V]{ExistingValue: est.node.value}
	}

	est.node = &nodeValue[V]{state}
	return nil
}

func (est edgeSetTerminal[V]) result(parameterValues []string) *Result[V] {
	if est.node == nil {
		return nil
	}

	return est.node.result(parameterValues)
}

func (est edgeSetTerminal[V]) walk(state stateWalk[V]) bool {
	if est.node == nil {
		return false
	}

	return state.visitor.VisitWalk(est.node.value)
}

type edgeSetValue[V any] struct{ term edgeSetTerminal[V] }

func (esv *edgeSetValue[V]) add(e edgeValue, state stateAdd[V]) error {
	return esv.term.add(state)
}

func (esv edgeSetValue[V]) search(state stateSearch[V]) bool {
	if result := esv.term.result(state.parameterValues); result != nil {
		return state.visitor.VisitSearch(result)
	}

	return false
}

func (esv edgeSetValue[V]) walk(state stateWalk[V]) bool {
	return esv.term.walk(state)
}

type edgeSetWildcard[V any] struct{ term edgeSetTerminal[V] }

func (esw *edgeSetWildcard[V]) add(e edgeWildcard, path []Key, state stateAdd[V]) error {
	if len(path) > 0 {
		return InvalidContinuationError{Continuation: path}
	}

	return esw.term.add(state)
}

func (esw edgeSetWildcard[V]) search(query []string, state stateSearch[V]) bool {
	if result := esw.term.result(state.parameterValues); result != nil {
		if len(query) > 0 {
			result.Tail = query
		}

		return state.visitor.VisitSearch(result)
	}

	return false
}

func (esw edgeSetWildcard[V]) walk(state stateWalk[V]) bool {
	return esw.term.walk(state)
}

type edgeSetConstant[V any] map[edgeConstant]*nodeConstant[V]

func (esc *edgeSetConstant[V]) add(e edgeConstant, path []Key, state stateAdd[V]) error {
	if *esc == nil {
		*esc = make(edgeSetConstant[V])
	}

	node, ok := (*esc)[e]
	if !ok {
		node = new(nodeConstant[V])
		(*esc)[e] = node
	}

	return node.add(path, state)
}

func (esc edgeSetConstant[V]) search(query []string, state stateSearch[V]) bool {
	head, tail := edgeConstant(query[0]), query[1:]

	if node, ok := esc[head]; ok {
		return node.search(tail, state)
	}

	return false
}

func (esc edgeSetConstant[V]) walk(state stateWalk[V]) bool {
	for _, node := range esc {
		if node.walk(state) {
			return true
		}
	}

	return false
}

type edgeSetParameter[V any] struct {
	nList sort.IntSlice
	nMap  map[int]*nodeParameter[V]
}

func (esp *edgeSetParameter[V]) createEntry(n int) *nodeParameter[V] {
	node := new(nodeParameter[V])

	esp.nMap[n] = node

	// TODO: Optimize via pes.nList.Search(n)
	esp.nList = append(esp.nList, n)
	esp.nList.Sort()

	return node
}

func (esp *edgeSetParameter[V]) add(e edgeParameter, path []Key, state stateAdd[V]) error {
	if esp.nMap == nil {
		esp.nMap = make(map[int]*nodeParameter[V])
	}

	state.parameterKeys = append(state.parameterKeys, e...)
	n := len(e)

	node, ok := esp.nMap[n]
	if !ok {
		node = esp.createEntry(n)
	}

	return node.add(path, state)
}

func (esp edgeSetParameter[V]) search(query []string, state stateSearch[V]) bool {
	var (
		nSegs        = len(query)
		wildSearches []func() bool
	)

	for _, nParams := range esp.nList {
		if nParams > nSegs {
			break
		}

		var (
			childNode  = esp.nMap[nParams]
			childQuery = query[nParams:]
			childState = stateSearch[V]{
				parameterValues: append(state.parameterValues, query[:nParams]...),
				visitor:         state.visitor,
			}
		)

		if childNode.searchStatic(childQuery, childState) {
			return true
		}

		wildSearches = append(wildSearches, func() bool {
			return childNode.searchWild(childQuery, childState)
		})
	}

	for i := len(wildSearches) - 1; i >= 0; i-- {
		if wildSearches[i]() {
			return true
		}
	}

	return false
}

func (esp edgeSetParameter[V]) walk(state stateWalk[V]) bool {
	for _, node := range esp.nMap {
		if node.walk(state) {
			return true
		}
	}

	return false
}
