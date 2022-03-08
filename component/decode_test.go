package component_test

import (
	"fmt"
	"testing"

	"github.com/oligarch316/go-urlrouter/component"
	"github.com/oligarch316/go-urlrouter/graph"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComponentDecodeKeyFailure(t *testing.T) {
	inputs := []string{
		"",
		":",
	}

	for _, input := range inputs {
		var (
			in   = input
			name = fmt.Sprintf("input '%s'", in)
		)

		t.Run(name, func(t *testing.T) {
			actual := component.DefaultKeyDecoder(in)
			assert.Nil(t, actual)
		})
	}
}

func TestComponentDecodeKeySuccess(t *testing.T) {
	subtests := []struct {
		input    string
		expected graph.Key
	}{
		{
			input:    "someConst",
			expected: graph.KeyConstant("someConst"),
		},
		{
			input:    ":someParam",
			expected: graph.KeyParameter("someParam"),
		},
		{
			input:    "*someWild",
			expected: graph.KeyWildcard{},
		},
	}

	for _, subtest := range subtests {
		var (
			st   = subtest
			name = fmt.Sprintf("input '%s'", st.input)
		)

		t.Run(name, func(t *testing.T) {
			actual := component.DefaultKeyDecoder(st.input)

			require.IsType(t, st.expected, actual)
			assert.EqualValues(t, st.expected, actual)
		})
	}
}
