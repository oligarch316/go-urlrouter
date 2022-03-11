package graph

type nodeValue[V any] struct {
	parameterKeys []string
	value         V
}

func (nv nodeValue[V]) result(paramVals []string) *Result[V] {
	res := &Result[V]{Value: nv.value}

	if len(nv.parameterKeys) > 0 {
		res.Parameters = make(map[string]string)

		for i, key := range nv.parameterKeys {
			res.Parameters[key] = paramVals[i]
		}
	}

	return res
}

type nodeConstant[V any] struct {
	constantEdges  edgeSetConstant[V]
	parameterEdges edgeSetParameter[V]
	valueEdges     edgeSetValue[V]
	wildcardEdges  edgeSetWildcard[V]
}

func (nc *nodeConstant[V]) add(keys []Key, val *nodeValue[V]) error {
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

func (nc nodeConstant[V]) search(segs, paramVals []string) *Result[V] {
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

func (nc nodeConstant[V]) values() []V {
	res := nc.valueEdges.values()
	res = append(res, nc.constantEdges.values()...)
	res = append(res, nc.parameterEdges.values()...)
	res = append(res, nc.wildcardEdges.values()...)

	return res
}

type nodeParameter[V any] struct {
	constantEdges edgeSetConstant[V]
	valueEdges    edgeSetValue[V]
	wildcardEdges edgeSetWildcard[V]
}

func (np *nodeParameter[V]) add(keys []Key, val *nodeValue[V]) error {
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

func (np nodeParameter[V]) searchStatic(segs, paramVals []string) *Result[V] {
	if len(segs) < 1 {
		return np.valueEdges.result(paramVals)
	}

	return np.constantEdges.search(segs, paramVals)
}

func (np nodeParameter[V]) searchWild(segs, paramVals []string) *Result[V] {
	return np.wildcardEdges.result(segs, paramVals)
}

func (np nodeParameter[V]) values() []V {
	res := np.valueEdges.values()
	res = append(res, np.constantEdges.values()...)
	res = append(res, np.wildcardEdges.values()...)

	return res
}
