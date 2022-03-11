package priority

type nodeValue[V any] struct{ stateAdd[V] }

func (nv nodeValue[V]) result(parameterValues []string) *Result[V] {
	res := &Result[V]{Value: nv.value}

	if len(nv.parameterKeys) > 0 {
		res.Parameters = make(map[string]string)

		for i, key := range nv.parameterKeys {
			res.Parameters[key] = parameterValues[i]
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

func (nc *nodeConstant[V]) add(path []Key, state stateAdd[V]) error {
	head, tail, err := popEdge(path)
	if err != nil {
		return err
	}

	switch e := head.(type) {
	case edgeValue:
		return nc.valueEdges.add(e, state)
	case edgeConstant:
		return nc.constantEdges.add(e, path, state)
	case edgeParameter:
		return nc.parameterEdges.add(e, path, state)
	case edgeWildcard:
		return nc.wildcardEdges.add(e, tail, state)
	}

	return internalErrorf("constant node: invalid edge type %T: %s", head, head)
}

func (nc nodeConstant[V]) search(query []string, state stateSearch[V]) bool {
	if len(query) < 1 {
		if nc.valueEdges.search(state) {
			return true
		}

		return nc.wildcardEdges.search(nil, state)
	}

	if nc.constantEdges.search(query, state) {
		return true
	}

	if nc.parameterEdges.search(query, state) {
		return true
	}

	return nc.wildcardEdges.search(query, state)
}

func (nc nodeConstant[V]) walk(state stateWalk[V]) bool {
	if nc.valueEdges.walk(state) {
		return true
	}

	if nc.constantEdges.walk(state) {
		return true
	}

	if nc.parameterEdges.walk(state) {
		return true
	}

	return nc.wildcardEdges.walk(state)
}

type nodeParameter[V any] struct {
	constantEdges edgeSetConstant[V]
	valueEdges    edgeSetValue[V]
	wildcardEdges edgeSetWildcard[V]
}

func (np *nodeParameter[V]) add(path []Key, state stateAdd[V]) error {
	head, tail, err := popEdge(path)
	if err != nil {
		return err
	}

	switch e := head.(type) {
	case edgeValue:
		return np.valueEdges.add(e, state)
	case edgeConstant:
		return np.constantEdges.add(e, tail, state)
	case edgeWildcard:
		return np.wildcardEdges.add(e, tail, state)
	}

	return internalErrorf("parameter node: invalid edge type %T: %s", head, head)
}

func (np nodeParameter[V]) searchStatic(query []string, state stateSearch[V]) bool {
	if len(query) < 1 {
		return np.valueEdges.search(state)
	}

	return np.constantEdges.search(query, state)
}

func (np nodeParameter[V]) searchWild(query []string, state stateSearch[V]) bool {
	return np.wildcardEdges.search(query, state)
}

func (np nodeParameter[V]) walk(state stateWalk[V]) bool {
	if np.valueEdges.walk(state) {
		return true
	}

	if np.constantEdges.walk(state) {
		return true
	}

	return np.wildcardEdges.walk(state)
}
