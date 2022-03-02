package graph

import (
	"fmt"
	"strings"
)

// TODO: These should become generic constraints
type (
	Parameter string
	Segment   string
	Value     interface{}
)

type Result struct {
	Parameters map[Parameter]Segment
	Tail       []Segment
	Value      Value
}

type (
	Key interface {
		fmt.Stringer
		sealedKey()
	}

	KeyConstant  Segment
	KeyParameter Parameter
	KeyWildcard  struct{}
)

func (KeyConstant) sealedKey()  {}
func (KeyParameter) sealedKey() {}
func (KeyWildcard) sealedKey()  {}

func (KeyWildcard) String() string     { return "wild" }
func (kc KeyConstant) String() string  { return fmt.Sprintf("const(%v)", Segment(kc)) }
func (kp KeyParameter) String() string { return fmt.Sprintf("param(%v)", Parameter(kp)) }

type (
	edge interface {
		fmt.Stringer
		sealedEdge()
	}

	edgeConstant  Segment
	edgeParameter []Parameter
	edgeValue     struct{}
	edgeWildcard  struct{}
)

func (edgeConstant) sealedEdge()  {}
func (edgeParameter) sealedEdge() {}
func (edgeValue) sealedEdge()     {}
func (edgeWildcard) sealedEdge()  {}

func (ep edgeValue) String() string    { return "value" }
func (ew edgeWildcard) String() string { return "wild" }
func (ec edgeConstant) String() string { return fmt.Sprintf("const(%v)", Segment(ec)) }

func (ep edgeParameter) String() string {
	strs := make([]string, len(ep))
	for i, param := range ep {
		strs[i] = fmt.Sprintf("%v", param)
	}
	return fmt.Sprintf("param(%s)", strings.Join(strs, ","))
}
