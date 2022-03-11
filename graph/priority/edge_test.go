package priority

import (
	"testing"

	"github.com/oligarch316/go-urlrouter/graph"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGraphPopEdge(t *testing.T) {
	var (
		a = graph.KeyConstant("constA")
		b = graph.KeyConstant("constB")
		c = graph.KeyConstant("constB")

		param1 = graph.KeyParameter("paramA")
		param2 = graph.KeyParameter("paramB")
		param3 = graph.KeyParameter("paramC")

		wild graph.KeyWildcard
	)

	subtests := []struct {
		name         string
		keys         []graph.Key
		expectedHead edge
		expectedTail []graph.Key
	}{
		{
			name:         "empty",
			keys:         nil,
			expectedHead: edgeValue{},
			expectedTail: nil,
		},
		{
			name:         "all constant",
			keys:         []graph.Key{a, b, c},
			expectedHead: edgeConstant("constA"),
			expectedTail: []graph.Key{b, c},
		},
		{
			name:         "all parameter",
			keys:         []graph.Key{param1, param2, param3},
			expectedHead: edgeParameter{"paramA", "paramB", "paramC"},
			expectedTail: nil,
		},
		{
			name:         "constant first",
			keys:         []graph.Key{a, param2, param3},
			expectedHead: edgeConstant("constA"),
			expectedTail: []graph.Key{param2, param3},
		},
		{
			name:         "parameter first (single)",
			keys:         []graph.Key{param1, b, c},
			expectedHead: edgeParameter{"paramA"},
			expectedTail: []graph.Key{b, c},
		},
		{
			name:         "parameter first (multi)",
			keys:         []graph.Key{param1, param2, c},
			expectedHead: edgeParameter{"paramA", "paramB"},
			expectedTail: []graph.Key{c},
		},
		{
			name:         "wildcard first",
			keys:         []graph.Key{wild, b, c},
			expectedHead: edgeWildcard{},
			expectedTail: []graph.Key{b, c},
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
