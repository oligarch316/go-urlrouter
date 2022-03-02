package graph

import "errors"

var ErrNilKey = errors.New("nil key")

func popEdge(keys []Key) (edge, []Key, error) {
	if len(keys) < 1 {
		return edgeValue{}, nil, nil
	}

	var paramEdge edgeParameter

	for i, key := range keys {
		if key == nil {
			return nil, nil, ErrNilKey
		}

		switch t := key.(type) {
		case KeyParameter:
			paramEdge = append(paramEdge, Parameter(t))
		case KeyConstant:
			if i == 0 {
				return edgeConstant(t), keys[1:], nil
			}

			return paramEdge, keys[i:], nil
		case KeyWildcard:
			if i == 0 {
				return edgeWildcard{}, keys[1:], nil
			}

			return paramEdge, keys[i:], nil
		}
	}

	return paramEdge, nil, nil
}

type nodeValue struct {
	parameterKeys []Parameter
	value         Value
}

func (nv nodeValue) result(paramVals []Segment) *Result {
	res := &Result{
		Parameters: make(map[Parameter]Segment),
		Value:      nv.value,
	}

	for i, key := range nv.parameterKeys {
		res.Parameters[key] = paramVals[i]
	}

	return res
}

type nodeConstant struct {
	constantEdges  constantEdgeSet
	parameterEdges parameterEdgeSet
	valueEdges     valueEdgeSet
	wildcardEdges  wildcardEdgeSet
}

func (nc *nodeConstant) add(keys []Key, val *nodeValue) error {
	head, tail, err := popEdge(keys)
	if err != nil {
		return err
	}

	switch e := head.(type) {
	case edgeValue:
		return nc.valueEdges.add(e, val)
	case edgeConstant:
		return nc.constantEdges.add(e, tail, val)
	case edgeParameter:
		return nc.parameterEdges.add(e, tail, val)
	case edgeWildcard:
		return nc.wildcardEdges.add(e, tail, val)
	}

	return internalErrorf("constant node: invalid edge type %T: %s", head, head)
}

func (nc nodeConstant) search(segs, paramVals []Segment) *Result {
	if len(segs) < 1 {
		if res := nc.valueEdges.result(paramVals); res != nil {
			return res
		}

		return nc.wildcardEdges.result(nil, paramVals)
	}

	if res := nc.constantEdges.search(segs, paramVals); res != nil {
		return res
	}

	if res := nc.parameterEdges.search(segs, paramVals); res != nil {
		return res
	}

	return nc.wildcardEdges.result(segs, paramVals)
}

type nodeParameter struct {
	constantEdges constantEdgeSet
	valueEdges    valueEdgeSet
	wildcardEdges wildcardEdgeSet
}

func (np *nodeParameter) add(keys []Key, val *nodeValue) error {
	head, tail, err := popEdge(keys)
	if err != nil {
		return err
	}

	switch e := head.(type) {
	case edgeValue:
		return np.valueEdges.add(e, val)
	case edgeConstant:
		return np.constantEdges.add(e, tail, val)
	case edgeWildcard:
		return np.wildcardEdges.add(e, tail, val)
	}

	return internalErrorf("parameter node: invalid edge type %T: %s", head, head)
}

func (np nodeParameter) searchStatic(segs, paramVals []Segment) *Result {
	if len(segs) < 1 {
		return np.valueEdges.result(paramVals)
	}

	return np.constantEdges.search(segs, paramVals)
}

func (np nodeParameter) searchWild(segs, paramVals []Segment) *Result {
	return np.wildcardEdges.result(segs, paramVals)
}
