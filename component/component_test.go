package component_test

import (
	"fmt"
	"testing"

	"github.com/oligarch316/go-urlrouter/component"
	"github.com/oligarch316/go-urlrouter/graph"
	"github.com/oligarch316/go-urlrouter/graph/search"
	graphtest "github.com/oligarch316/go-urlrouter/graph/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Query string

func (q Query) String() string {
	return fmt.Sprintf("--------\nQuery:\n> %s", string(q))
}

type RouteItem struct{ pattern, value string }

func (ri RouteItem) String() string {
	return fmt.Sprintf("--------\nRoute:\n> %s→%s", ri.pattern, ri.value)
}

type Tree struct{ graphtest.Tree }

func (t *Tree) AsOption(r *component.Router[string]) { r.Tree = t }

func Route(pattern, value string) RouteItem {
	return RouteItem{pattern: pattern, value: value}
}

func TestComponentRouterHostSearch(t *testing.T) {
	var (
		tree   Tree
		router = component.NewHostRouter(tree.AsOption)
	)

	routes := []RouteItem{
		Route("a.b.com", "valCom1"),
		Route(":some.b.com", "valCom2"),
		Route("*any.b.com", "valCom3"),

		Route("a.b.:other", "valOther1"),
		Route(":some.b.:other", "valOther2"),
		Route("*any.b.:other", "valOther3"),

		Route("*any", "valAny"),
	}

	for _, route := range routes {
		err := router.Add(route.pattern, route.value)

		require.NoError(t, err, graphtest.Info(route, &tree).Note("check add error"))
	}

	searchTests := []struct {
		query    Query
		expected graph.SearchResult[string]
	}{
		{
			query: "a.b.com",
			expected: graph.SearchResult[string]{
				Value: "valCom1",
			},
		},
		{
			query: "x.b.com",
			expected: graph.SearchResult[string]{
				Value: "valCom2",
				Parameters: map[string]string{
					"some": "x",
				},
			},
		},
		{
			query: "x.y.b.com",
			expected: graph.SearchResult[string]{
				Value: "valCom3",
				Tail:  []string{"y", "x"},
			},
		},
		{
			query: "a.b.org",
			expected: graph.SearchResult[string]{
				Value: "valOther1",
				Parameters: map[string]string{
					"other": "org",
				},
			},
		},
		{
			query: "x.b.org",
			expected: graph.SearchResult[string]{
				Value: "valOther2",
				Parameters: map[string]string{
					"some":  "x",
					"other": "org",
				},
			},
		},
		{
			query: "x.y.b.org",
			expected: graph.SearchResult[string]{
				Value: "valOther3",
				Parameters: map[string]string{
					"other": "org",
				},
				Tail: []string{"y", "x"},
			},
		},

		{
			query: "x.y.com",
			expected: graph.SearchResult[string]{
				Value: "valAny",
				Tail:  []string{"com", "y", "x"},
			},
		},
		{
			query: "x.y.org",
			expected: graph.SearchResult[string]{
				Value: "valAny",
				Tail:  []string{"org", "y", "x"},
			},
		},
	}

	for _, searchTest := range searchTests {
		var (
			query    = searchTest.query
			expected = searchTest.expected

			visitor = new(search.VisitorFirst[string])
			info    = graphtest.Info(query, &tree)
		)

		err := router.Search(visitor, string(query))
		if !assert.NoError(t, err, info.Note("check search error")) {
			continue
		}

		actual := visitor.Result

		if !assert.NotNil(t, actual, info.Note("check nil result")) {
			continue
		}

		if !assert.Equal(t, expected.Value, actual.Value, info.Note("check value")) {
			continue
		}

		assert.Equal(t, expected.Parameters, actual.Parameters, info.Note("check params"))
		assert.Equal(t, expected.Tail, actual.Tail, info.Note("check tail"))
	}
}

func TestComponentRouterPathSearch(t *testing.T) {
	var (
		tree   Tree
		router = component.NewPathRouter(tree.AsOption)
	)

	routes := []RouteItem{
		Route("/foo/a/b", "valFoo1"),
		Route("/foo/a/:some", "valFoo2"),
		Route("/foo/a/*any", "valFoo3"),

		Route("/:other/a/b", "valOther1"),
		Route("/:other/a/:some", "valOther2"),
		Route("/:other/a/*any", "valOther3"),

		Route("/*any", "valAny"),
	}

	for _, route := range routes {
		err := router.Add(route.pattern, route.value)

		require.NoError(t, err, graphtest.Info(route, &tree).Note("check add error"))
	}

	searchTests := []struct {
		query    Query
		expected graph.SearchResult[string]
	}{
		{
			query: "/foo/a/b",
			expected: graph.SearchResult[string]{
				Value: "valFoo1",
			},
		},
		{
			query: "/foo/a/x",
			expected: graph.SearchResult[string]{
				Value: "valFoo2",
				Parameters: map[string]string{
					"some": "x",
				},
			},
		},
		{
			query: "/foo/a/x/y",
			expected: graph.SearchResult[string]{
				Value: "valFoo3",
				Tail:  []string{"x", "y"},
			},
		},
		{
			query: "/bar/a/b",
			expected: graph.SearchResult[string]{
				Value: "valOther1",
				Parameters: map[string]string{
					"other": "bar",
				},
			},
		},
		{
			query: "/bar/a/x",
			expected: graph.SearchResult[string]{
				Value: "valOther2",
				Parameters: map[string]string{
					"some":  "x",
					"other": "bar",
				},
			},
		},
		{
			query: "/bar/a/x/y",
			expected: graph.SearchResult[string]{
				Value: "valOther3",
				Parameters: map[string]string{
					"other": "bar",
				},
				Tail: []string{"x", "y"},
			},
		},
		{
			query: "/foo/x/y",
			expected: graph.SearchResult[string]{
				Value: "valAny",
				Tail:  []string{"foo", "x", "y"},
			},
		},
		{
			query: "/bar/x/y",
			expected: graph.SearchResult[string]{
				Value: "valAny",
				Tail:  []string{"bar", "x", "y"},
			},
		},
	}

	for _, searchTest := range searchTests {
		var (
			query    = searchTest.query
			expected = searchTest.expected

			visitor = new(search.VisitorFirst[string])
			info    = graphtest.Info(query, &tree)
		)

		err := router.Search(visitor, string(query))
		if !assert.NoError(t, err, info.Note("check search error")) {
			continue
		}

		actual := visitor.Result

		if !assert.NotNil(t, actual, info.Note("check nil result")) {
			continue
		}

		if !assert.Equal(t, expected.Value, actual.Value, info.Note("check value")) {
			continue
		}

		assert.Equal(t, expected.Parameters, actual.Parameters, info.Note("check params"))
		assert.Equal(t, expected.Tail, actual.Tail, info.Note("check tail"))
	}
}
