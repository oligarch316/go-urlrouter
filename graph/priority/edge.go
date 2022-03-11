package graph

import (
	"errors"
	"fmt"
	"strings"
)

var ErrNilKey = errors.New("nil key")

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
			paramEdge = append(paramEdge, string(t))
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
