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

type Tree[V any] struct{ root nodeConstant[V] }

func (t *Tree[V]) Add(val V, keys ...Key) error {
	return t.root.add(keys, &nodeValue[V]{value: val})
}

func (t *Tree[V]) Search(segs ...string) *Result[V] {
	return t.root.search(segs, nil)
}

func (t *Tree[V]) Values() []V {
	return t.root.values()
}
