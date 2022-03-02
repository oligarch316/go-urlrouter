package graph

import "sort"

type (
	DuplicateValueError      struct{ ExistingValue Value }
	InvalidContinuationError struct{ Continuation []Key }
)

func (dve DuplicateValueError) Error() string      { return "dupliate value" }
func (ice InvalidContinuationError) Error() string { return "invalid continuation" }

type terminalEdge struct{ val *nodeValue }

func (te *terminalEdge) add(val *nodeValue) error {
	if te.val != nil {
		return DuplicateValueError{ExistingValue: te.val.value}
	}

	te.val = val
	return nil
}

func (te *terminalEdge) result(paramVals []Segment) *Result {
	if te.val == nil {
		return nil
	}

	return te.val.result(paramVals)
}

type valueEdgeSet struct{ term terminalEdge }

func (ves *valueEdgeSet) add(e edgeValue, val *nodeValue) error {
	return ves.term.add(val)
}

func (ves *valueEdgeSet) result(paramVals []Segment) *Result {
	return ves.term.result(paramVals)
}

type wildcardEdgeSet struct{ term terminalEdge }

func (wes *wildcardEdgeSet) add(e edgeWildcard, keys []Key, val *nodeValue) error {
	if len(keys) > 0 {
		return InvalidContinuationError{Continuation: keys}
	}

	return wes.term.add(val)
}

func (wes *wildcardEdgeSet) result(segs, paramVals []Segment) *Result {
	res := wes.term.result(paramVals)
	if res != nil {
		res.Tail = segs
	}

	return res
}

type constantEdgeSet map[edgeConstant]*nodeConstant

func (ces *constantEdgeSet) add(e edgeConstant, keys []Key, val *nodeValue) error {
	if *ces == nil {
		*ces = make(constantEdgeSet)
	}

	node, ok := (*ces)[e]
	if !ok {
		node = new(nodeConstant)
		(*ces)[e] = node
	}

	return node.add(keys, val)
}

func (ces constantEdgeSet) search(segs, paramVals []Segment) *Result {
	head, tail := edgeConstant(segs[0]), segs[1:]

	node, ok := ces[head]
	if !ok {
		return nil
	}

	return node.search(tail, paramVals)
}

/* TODO
 * The search logic here currently iterates over candidate nodes in order of shortest -> longest
 * number of param keys need to reach said node. This results in correct prioritization when following
 * that node through to a constant or value edge, but is incorrect for following to a wildcard edge.
 * In the wildcard case, longer param key chains should take priority over shorter.
 *
 * Current plan for a fix is to:
 * 1.) Remove wildcardEdges from nodeParameter
 * 2.) Have nMap value be struct{ nodeParameter, wildcardEdgeSet } (separate map would be optimized but don't bother now)
 * 3.) Modify add to check for wildcard edge from popKeys(...) and add appropriately
 * 4.) Modify search to perform wildcard search after node search and in reverse order
 */

type parameterEdgeSet struct {
	nList sort.IntSlice
	nMap  map[int]*nodeParameter
}

func (pes *parameterEdgeSet) createEntry(n int) *nodeParameter {
	node := new(nodeParameter)

	pes.nMap[n] = node

	// TODO: Optimize via pes.nList.Search(n)
	pes.nList = append(pes.nList, n)
	pes.nList.Sort()

	return node
}

func (pes *parameterEdgeSet) add(e edgeParameter, keys []Key, val *nodeValue) error {
	if pes.nMap == nil {
		pes.nMap = make(map[int]*nodeParameter)
	}

	val.parameterKeys = append(val.parameterKeys, e...)
	n := len(e)

	node, ok := pes.nMap[n]
	if !ok {
		node = pes.createEntry(n)
	}

	return node.add(keys, val)
}

func (pes parameterEdgeSet) search(segs, paramVals []Segment) *Result {
	nSegs := len(segs)

	for _, nParams := range pes.nList {
		if nParams > nSegs {
			break
		}

		var (
			childNode      = pes.nMap[nParams]
			childSegs      = segs[nParams:]
			childParamVals = append(paramVals, segs[:nParams]...)
		)

		if res := childNode.search(childSegs, childParamVals); res != nil {
			return res
		}
	}

	return nil
}
