package graph

import "sort"

type (
	DuplicateValueError[V any] struct{ ExistingValue V }
	InvalidContinuationError   struct{ Continuation []Key }
)

func (dve DuplicateValueError[_]) Error() string   { return "dupliate value" }
func (ice InvalidContinuationError) Error() string { return "invalid continuation" }

type terminalEdge[V any] struct{ val *nodeValue[V] }

func (te *terminalEdge[V]) add(val *nodeValue[V]) error {
	if te.val != nil {
		return DuplicateValueError[V]{ExistingValue: te.val.value}
	}

	te.val = val
	return nil
}

func (te *terminalEdge[V]) result(paramVals []Segment) *Result[V] {
	if te.val == nil {
		return nil
	}

	return te.val.result(paramVals)
}

func (te *terminalEdge[V]) values() []V {
	if te.val == nil {
		return nil
	}

	return []V{te.val.value}
}

type valueEdgeSet[V any] struct{ term terminalEdge[V] }

func (ves *valueEdgeSet[V]) add(e edgeValue, val *nodeValue[V]) error {
	return ves.term.add(val)
}

func (ves *valueEdgeSet[V]) result(paramVals []Segment) *Result[V] {
	return ves.term.result(paramVals)
}

func (ves *valueEdgeSet[V]) values() []V {
	return ves.term.values()
}

type wildcardEdgeSet[V any] struct{ term terminalEdge[V] }

func (wes *wildcardEdgeSet[V]) add(e edgeWildcard, keys []Key, val *nodeValue[V]) error {
	if len(keys) > 0 {
		return InvalidContinuationError{Continuation: keys}
	}

	return wes.term.add(val)
}

func (wes *wildcardEdgeSet[V]) result(segs, paramVals []Segment) *Result[V] {
	res := wes.term.result(paramVals)

	if res != nil && len(segs) > 0 {
		res.Tail = segs
	}

	return res
}

func (wes *wildcardEdgeSet[V]) values() []V {
	return wes.term.values()
}

type constantEdgeSet[V any] map[edgeConstant]*nodeConstant[V]

func (ces *constantEdgeSet[V]) add(e edgeConstant, keys []Key, val *nodeValue[V]) error {
	if *ces == nil {
		*ces = make(constantEdgeSet[V])
	}

	node, ok := (*ces)[e]
	if !ok {
		node = new(nodeConstant[V])
		(*ces)[e] = node
	}

	return node.add(keys, val)
}

func (ces constantEdgeSet[V]) search(segs, paramVals []Segment) *Result[V] {
	head, tail := edgeConstant(segs[0]), segs[1:]

	node, ok := ces[head]
	if !ok {
		return nil
	}

	return node.search(tail, paramVals)
}

func (ces constantEdgeSet[V]) values() []V {
	var res []V

	for _, node := range ces {
		res = append(res, node.values()...)
	}

	return res
}

type parameterEdgeSet[V any] struct {
	nList sort.IntSlice
	nMap  map[int]*nodeParameter[V]
}

func (pes *parameterEdgeSet[V]) createEntry(n int) *nodeParameter[V] {
	node := new(nodeParameter[V])

	pes.nMap[n] = node

	// TODO: Optimize via pes.nList.Search(n)
	pes.nList = append(pes.nList, n)
	pes.nList.Sort()

	return node
}

func (pes *parameterEdgeSet[V]) add(e edgeParameter, keys []Key, val *nodeValue[V]) error {
	if pes.nMap == nil {
		pes.nMap = make(map[int]*nodeParameter[V])
	}

	val.parameterKeys = append(val.parameterKeys, e...)
	n := len(e)

	node, ok := pes.nMap[n]
	if !ok {
		node = pes.createEntry(n)
	}

	return node.add(keys, val)
}

func (pes parameterEdgeSet[V]) search(segs, paramVals []Segment) *Result[V] {
	var (
		nSegs        = len(segs)
		wildSearches []func() *Result[V]
	)

	for _, nParams := range pes.nList {
		if nParams > nSegs {
			break
		}

		var (
			childNode      = pes.nMap[nParams]
			childSegs      = segs[nParams:]
			childParamVals = append(paramVals, segs[:nParams]...)
		)

		if res := childNode.searchStatic(childSegs, childParamVals); res != nil {
			return res
		}

		wildSearches = append(wildSearches, func() *Result[V] {
			return childNode.searchWild(childSegs, childParamVals)
		})
	}

	for i := len(wildSearches) - 1; i >= 0; i-- {
		if res := wildSearches[i](); res != nil {
			return res
		}
	}

	return nil
}

func (pes parameterEdgeSet[V]) values() []V {
	var res []V

	for _, node := range pes.nMap {
		res = append(res, node.values()...)
	}

	return res
}
