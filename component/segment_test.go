package component_test

import (
	"fmt"
	"testing"

	"github.com/oligarch316/go-urlrouter/component"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComponentSegmentPathError(t *testing.T) {
	inputs := []string{
		"",
		"noslash",
	}

	for _, input := range inputs {
		var (
			in   = input
			name = fmt.Sprintf("input '%s'", in)
		)

		t.Run(name, func(t *testing.T) {
			_, err := component.DefaultPathSegmenter(in)
			assert.ErrorIs(t, err, component.ErrInvalidPath)
		})
	}
}

func TestComponentSegmentHostSuccess(t *testing.T) {
	subtests := []struct {
		input    string
		expected []string
	}{
		{
			input:    "",
			expected: nil,
		},
		{
			input:    ".",
			expected: []string{"", ""},
		},
		{
			input:    "tld",
			expected: []string{"tld"},
		},
		{
			input:    "sub1.sub2.tld",
			expected: []string{"tld", "sub2", "sub1"},
		},
	}

	for _, subtest := range subtests {
		var (
			st   = subtest
			name = fmt.Sprintf("input '%s'", st.input)
		)

		t.Run(name, func(t *testing.T) {
			actual, err := component.DefaultHostSegmenter(st.input)

			require.NoError(t, err)
			assert.Equal(t, st.expected, actual)
		})
	}
}

func TestComponentSegmentPathSuccess(t *testing.T) {
	subtests := []struct {
		input    string
		expected []string
	}{
		{
			input:    "/",
			expected: nil,
		},
		{
			input:    "//",
			expected: nil,
		},
		{
			input:    "///",
			expected: []string{"", ""},
		},
		{
			input:    "/a",
			expected: []string{"a"},
		},
		{
			input:    "/a/",
			expected: []string{"a"},
		},
		{
			input:    "/a/b/c",
			expected: []string{"a", "b", "c"},
		},
		{
			input:    "/a/b/c/",
			expected: []string{"a", "b", "c"},
		},
	}

	for _, subtest := range subtests {
		var (
			st   = subtest
			name = fmt.Sprintf("input '%s'", st.input)
		)

		t.Run(name, func(t *testing.T) {
			actual, err := component.DefaultPathSegmenter(st.input)

			require.NoError(t, err)
			assert.Equal(t, st.expected, actual)
		})
	}
}
