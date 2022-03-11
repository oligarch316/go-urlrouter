package graph

import (
	"errors"
	"fmt"
)

var ErrInternal = errors.New("internal")

type internalError string

func internalErrorf(format string, a ...interface{}) internalError {
	message := fmt.Sprintf(format, a...)
	return internalError(message)
}

func (ie internalError) Error() string        { return string(ie) }
func (ie internalError) Is(target error) bool { return target == ErrInternal }

type (
	Key interface {
		fmt.Stringer
		sealedKey()
	}

	KeyConstant  string
	KeyParameter string
	KeyWildcard  struct{}
)

func (KeyConstant) sealedKey()  {}
func (KeyParameter) sealedKey() {}
func (KeyWildcard) sealedKey()  {}

func (KeyWildcard) String() string     { return "wild" }
func (kc KeyConstant) String() string  { return fmt.Sprintf("const(%s)", string(kc)) }
func (kp KeyParameter) String() string { return fmt.Sprintf("param(%s)", string(kp)) }

type Result[V any] struct {
	Parameters map[string]string
	Tail       []string
	Value      V
}

type (
	Searcher[V any] interface {
		VisitSearch(result *Result[V]) (done bool)
	}

	SearcherFunc[V any] func(result *Result[V]) (done bool)

	Walker[V any] interface {
		VisitWalk(value V) (done bool)
	}

	WalkerFunc[V any] func(value V) (done bool)
)

func (sf SearcherFunc[V]) VisitSearch(result *Result[V]) bool { return sf(result) }
func (wf WalkerFunc[V]) VisitWalk(value V) bool               { return wf(value) }

type Tree[V any] struct{ root nodeConstant[V] }

func (t *Tree[V]) Add(value V, path ...Key) error {
	return t.root.add(path, stateAdd[V]{value: value})
}

func (t *Tree[V]) Search(searcher Searcher[V], query ...string) bool {
	return t.root.search(query, stateSearch[V]{visitor: searcher})
}

func (t *Tree[V]) SearchFunc(visitor func(result *Result[V]) (done bool), query ...string) bool {
	return t.Search(SearcherFunc[V](visitor), query...)
}

func (t *Tree[V]) Walk(walker Walker[V]) bool {
	return t.root.walk(stateWalk[V]{visitor: walker})
}

func (t *Tree[V]) WalkFunc(walker func(value V) (done bool)) bool {
	return t.Walk(WalkerFunc[V](walker))
}
