package walk

type VisitorAll[V any] struct{ Values []V }

func (va *VisitorAll[V]) VisitWalk(value V) bool {
	va.Values = append(va.Values, value)
	return false
}

type VisitorAllPredicate[V any] struct {
	Predicate func(V) bool
	Values    []V
}

func (vap *VisitorAllPredicate[V]) VisitWalk(value V) bool {
	if vap.Predicate(value) {
		vap.Values = append(vap.Values, value)
	}

	return false
}
