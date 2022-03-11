package graph

import "sort"

type (
	DuplicateValueError[V any] struct{ ExistingValue V }
	InvalidContinuationError   struct{ Continuation []Key }
)

func (dve DuplicateValueError[_]) Error() string   { return "dupliate value" }
func (ice InvalidContinuationError) Error() string { return "invalid continuation" }

type edgeSetTerminal[V any] struct{ val *nodeValue[V] }

func (est *edgeSetTerminal[V]) add(val *nodeValue[V]) error {
	if est.val != nil {
		return DuplicateValueError[V]{ExistingValue: est.val.value}
	}

	est.val = val
	return nil
}

func (est *edgeSetTerminal[V]) result(paramVals []string) *Result[V] {
	if est.val == nil {
		return nil
	}

	return est.val.result(paramVals)
}

func (est *edgeSetTerminal[V]) values() []V {
	if est.val == nil {
		return nil
	}

	return []V{est.val.value}
}

type edgeSetValue[V any] struct{ term edgeSetTerminal[V] }

func (esv *edgeSetValue[V]) add(e edgeValue, val *nodeValue[V]) error {
	return esv.term.add(val)
}

func (esv *edgeSetValue[V]) result(paramVals []string) *Result[V] {
	return esv.term.result(paramVals)
}

func (esv *edgeSetValue[V]) values() []V {
	return esv.term.values()
}

type edgeSetWildcard[V any] struct{ term edgeSetTerminal[V] }

func (esw *edgeSetWildcard[V]) add(e edgeWildcard, keys []Key, val *nodeValue[V]) error {
	if len(keys) > 0 {
		return InvalidContinuationError{Continuation: keys}
	}

	return esw.term.add(val)
}

func (esw *edgeSetWildcard[V]) result(segs, paramVals []string) *Result[V] {
	res := esw.term.result(paramVals)

	if res != nil && len(segs) > 0 {
		res.Tail = segs
	}

	return res
}

func (esw *edgeSetWildcard[V]) values() []V {
	return esw.term.values()
}

type edgeSetConstant[V any] map[edgeConstant]*nodeConstant[V]

func (esc *edgeSetConstant[V]) add(e edgeConstant, keys []Key, val *nodeValue[V]) error {
	if *esc == nil {
		*esc = make(edgeSetConstant[V])
	}

	node, ok := (*esc)[e]
	if !ok {
		node = new(nodeConstant[V])
		(*esc)[e] = node
	}

	return node.add(keys, val)
}

func (esc edgeSetConstant[V]) search(segs, paramVals []string) *Result[V] {
	head, tail := edgeConstant(segs[0]), segs[1:]

	node, ok := esc[head]
	if !ok {
		return nil
	}

	return node.search(tail, paramVals)
}

func (esc edgeSetConstant[V]) values() []V {
	var res []V

	for _, node := range esc {
		res = append(res, node.values()...)
	}

	return res
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

func (esp *edgeSetParameter[V]) add(e edgeParameter, keys []Key, val *nodeValue[V]) error {
	if esp.nMap == nil {
		esp.nMap = make(map[int]*nodeParameter[V])
	}

	val.parameterKeys = append(val.parameterKeys, e...)
	n := len(e)

	node, ok := esp.nMap[n]
	if !ok {
		node = esp.createEntry(n)
	}

	return node.add(keys, val)
}

func (esp edgeSetParameter[V]) search(segs, paramVals []string) *Result[V] {
	var (
		nSegs        = len(segs)
		wildSearches []func() *Result[V]
	)

	for _, nParams := range esp.nList {
		if nParams > nSegs {
			break
		}

		var (
			childNode      = esp.nMap[nParams]
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

func (esp edgeSetParameter[V]) values() []V {
	var res []V

	for _, node := range esp.nMap {
		res = append(res, node.values()...)
	}

	return res
}
