package graph

import (
	"fmt"
	"strings"
)

type Result[V any] struct {
	Parameters map[string]string
	Tail       []string
	Value      V
}

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

type (
	edge interface {
		fmt.Stringer
		sealedEdge()
	}

	edgeConstant  string
	edgeParameter []string
	edgeValue     struct{}
	edgeWildcard  struct{}
)

func (edgeConstant) sealedEdge()  {}
func (edgeParameter) sealedEdge() {}
func (edgeValue) sealedEdge()     {}
func (edgeWildcard) sealedEdge()  {}

func (ep edgeValue) String() string    { return "value" }
func (ew edgeWildcard) String() string { return "wild" }
func (ec edgeConstant) String() string { return fmt.Sprintf("const(%s)", string(ec)) }

func (ep edgeParameter) String() string {
	strs := make([]string, len(ep))
	for i, param := range ep {
		strs[i] = fmt.Sprintf("%s", param)
	}
	return fmt.Sprintf("param(%s)", strings.Join(strs, ","))
}
