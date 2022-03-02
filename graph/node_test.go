package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNodeAddInternalError(t *testing.T) {
	var (
		node nodeParameter[string]
		keys = []Key{KeyParameter("someParam")}
		val  = nodeValue[string]{value: "someVal"}
	)

	assert.ErrorIs(t, node.add(keys, &val), ErrInternal)
}

func TestNodePopEdgeError(t *testing.T) {
	// TODO: Remove this? It's checked in TestTreeAddError/nil_key

	keys := []Key{nil}
	_, _, err := popEdge(keys)

	assert.ErrorIs(t, err, ErrNilKey)
}

func TestNodePopEdgeSuccess(t *testing.T) {
	var (
		a = KeyConstant("constA")
		b = KeyConstant("constB")
		c = KeyConstant("constB")

		param1 = KeyParameter("paramA")
		param2 = KeyParameter("paramB")
		param3 = KeyParameter("paramC")

		wild KeyWildcard
	)

	subtests := []struct {
		name         string
		keys         []Key
		expectedHead edge
		expectedTail []Key
	}{
		{
			name:         "empty",
			keys:         nil,
			expectedHead: edgeValue{},
			expectedTail: nil,
		},
		{
			name:         "all constant",
			keys:         []Key{a, b, c},
			expectedHead: edgeConstant("constA"),
			expectedTail: []Key{b, c},
		},
		{
			name:         "all parameter",
			keys:         []Key{param1, param2, param3},
			expectedHead: edgeParameter{"paramA", "paramB", "paramC"},
			expectedTail: nil,
		},
		{
			name:         "constant first",
			keys:         []Key{a, param2, param3},
			expectedHead: edgeConstant("constA"),
			expectedTail: []Key{param2, param3},
		},
		{
			name:         "parameter first (single)",
			keys:         []Key{param1, b, c},
			expectedHead: edgeParameter{"paramA"},
			expectedTail: []Key{b, c},
		},
		{
			name:         "parameter first (multi)",
			keys:         []Key{param1, param2, c},
			expectedHead: edgeParameter{"paramA", "paramB"},
			expectedTail: []Key{c},
		},
		{
			name:         "wildcard first",
			keys:         []Key{wild, b, c},
			expectedHead: edgeWildcard{},
			expectedTail: []Key{b, c},
		},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			head, tail, err := popEdge(subtest.keys)

			require.NoError(t, err)

			assert.Equal(t, subtest.expectedHead, head)
			assert.Equal(t, subtest.expectedTail, tail)
		})
	}
}
