package component_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/oligarch316/go-urlrouter/component"
	"github.com/oligarch316/go-urlrouter/graph"
	"github.com/oligarch316/go-urlrouter/graph/memoized"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO: More comprehensive tests

const infoFormat = `--------
%s:
> %s
--------
Tree:
%s`

type actionInfo struct {
	action, data string
	tree         TestTree
}

func (ai actionInfo) String() string {
	return fmt.Sprintf(infoFormat, ai.action, ai.data, ai.tree)
}

type TestTree struct{ memoized.Tree[string] }

func (tt TestTree) String() string {
	var (
		strs  []string
		memos = tt.Memos()
	)

	if len(memos) == 0 {
		return "<empty>"
	}

	for _, memo := range memos {
		strs = append(strs, "> "+memo.String())
	}

	return strings.Join(strs, "\n")
}

func (tt *TestTree) asOption(p *component.Params[string]) { p.Tree = tt }

func (tt TestTree) addInfo(data string) fmt.Stringer {
	return actionInfo{action: "Add", data: data, tree: tt}
}

func (tt TestTree) searchInfo(data string) fmt.Stringer {
	return actionInfo{action: "Search", data: data, tree: tt}
}

func TestComponentRouteHost(t *testing.T) {
	var (
		tree   TestTree
		router = component.NewHostRouter(tree.asOption)
	)

	addItems := []struct{ host, value string }{
		{
			host:  "a.b.com",
			value: "valCom1",
		},
		{
			host:  ":some.b.com",
			value: "valCom2",
		},
		{
			host:  "*any.b.com",
			value: "valCom3",
		},
		{
			host:  "a.b.:other",
			value: "valOther1",
		},
		{
			host:  ":some.b.:other",
			value: "valOther2",
		},
		{
			host:  "*any.b.:other",
			value: "valOther3",
		},
		{
			host:  "*any",
			value: "valAny",
		},
	}

	for _, item := range addItems {
		err := router.Add(item.host, item.value)

		require.NoError(t, err, tree.addInfo(item.host))
	}

	searchItems := []struct {
		host     string
		expected graph.Result[string]
	}{
		{
			host: "a.b.com",
			expected: graph.Result[string]{
				Value: "valCom1",
			},
		},
		{
			host: "x.b.com",
			expected: graph.Result[string]{
				Value: "valCom2",
				Parameters: map[string]string{
					"some": "x",
				},
			},
		},
		{
			host: "x.y.b.com",
			expected: graph.Result[string]{
				Value: "valCom3",
				Tail:  []string{"y", "x"},
			},
		},
		{
			host: "a.b.org",
			expected: graph.Result[string]{
				Value: "valOther1",
				Parameters: map[string]string{
					"other": "org",
				},
			},
		},
		{
			host: "x.b.org",
			expected: graph.Result[string]{
				Value: "valOther2",
				Parameters: map[string]string{
					"some":  "x",
					"other": "org",
				},
			},
		},
		{
			host: "x.y.b.org",
			expected: graph.Result[string]{
				Value: "valOther3",
				Parameters: map[string]string{
					"other": "org",
				},
				Tail: []string{"y", "x"},
			},
		},
		{
			host: "x.y.com",
			expected: graph.Result[string]{
				Value: "valAny",
				Tail:  []string{"com", "y", "x"},
			},
		},
		{
			host: "x.y.org",
			expected: graph.Result[string]{
				Value: "valAny",
				Tail:  []string{"org", "y", "x"},
			},
		},
	}

	for _, item := range searchItems {
		var (
			result, err = router.Search(item.host)
			info        = tree.searchInfo(item.host)
		)

		if !assert.NoError(t, err, info) {
			continue
		}

		if !assert.Equal(t, item.expected.Value, result.Value, "check value\n%s", info) {
			continue
		}

		assert.Equal(t, item.expected.Parameters, result.Parameters, "check params\n%s", info)
		assert.Equal(t, item.expected.Tail, result.Tail, "check tail\n%s", info)
	}
}

func TestComponentRoutePath(t *testing.T) {
	var (
		tree   TestTree
		router = component.NewPathRouter(tree.asOption)
	)

	addItems := []struct{ path, value string }{
		{
			path:  "/foo/a/b",
			value: "valFoo1",
		},
		{
			path:  "/foo/a/:some",
			value: "valFoo2",
		},
		{
			path:  "/foo/a/*any",
			value: "valFoo3",
		},
		{
			path:  "/:other/a/b",
			value: "valOther1",
		},
		{
			path:  "/:other/a/:some",
			value: "valOther2",
		},
		{
			path:  "/:other/a/*any",
			value: "valOther3",
		},
		{
			path:  "/*any",
			value: "valAny",
		},
	}

	for _, item := range addItems {
		err := router.Add(item.path, item.value)

		require.NoError(t, err, tree.addInfo(item.path))
	}

	searchItems := []struct {
		path     string
		expected graph.Result[string]
	}{
		{
			path: "/foo/a/b",
			expected: graph.Result[string]{
				Value: "valFoo1",
			},
		},
		{
			path: "/foo/a/x",
			expected: graph.Result[string]{
				Value: "valFoo2",
				Parameters: map[string]string{
					"some": "x",
				},
			},
		},
		{
			path: "/foo/a/x/y",
			expected: graph.Result[string]{
				Value: "valFoo3",
				Tail:  []string{"x", "y"},
			},
		},
		{
			path: "/bar/a/b",
			expected: graph.Result[string]{
				Value: "valOther1",
				Parameters: map[string]string{
					"other": "bar",
				},
			},
		},
		{
			path: "/bar/a/x",
			expected: graph.Result[string]{
				Value: "valOther2",
				Parameters: map[string]string{
					"other": "bar",
					"some":  "x",
				},
			},
		},
		{
			path: "/bar/a/x/y",
			expected: graph.Result[string]{
				Value: "valOther3",
				Parameters: map[string]string{
					"other": "bar",
				},
				Tail: []string{"x", "y"},
			},
		},
		{
			path: "/foo/x/y",
			expected: graph.Result[string]{
				Value: "valAny",
				Tail:  []string{"foo", "x", "y"},
			},
		},
		{
			path: "/bar/x/y",
			expected: graph.Result[string]{
				Value: "valAny",
				Tail:  []string{"bar", "x", "y"},
			},
		},
	}

	for _, item := range searchItems {
		var (
			result, err = router.Search(item.path)
			info        = tree.searchInfo(item.path)
		)

		if !assert.NoError(t, err, info) {
			continue
		}

		if !assert.Equal(t, item.expected.Value, result.Value, "check value\n%s", info) {
			continue
		}

		assert.Equal(t, item.expected.Parameters, result.Parameters, "check params\n%s", info)
		assert.Equal(t, item.expected.Tail, result.Tail, "check tail\n%s", info)
	}
}
