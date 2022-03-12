package graphtest

import (
	"testing"

	"github.com/oligarch316/go-urlrouter/graph"
	"github.com/stretchr/testify/assert"
)

type searchResult graph.SearchResult[string]

type searchResultList []searchResult

func (srl searchResultList) values() []string {
	res := make([]string, len(srl))
	for i, result := range srl {
		res[i] = result.Value
	}
	return res
}

type searchVisitor struct{ actual searchResultList }

func (sv *searchVisitor) VisitSearch(result *graph.SearchResult[string]) bool {
	sv.actual = append(sv.actual, searchResult(*result))
	return false
}

func TestGraphSearch(t *testing.T) {
	var (
		a = graph.KeyConstant("a")
		b = graph.KeyConstant("b")
		c = graph.KeyConstant("c")

		param1 = graph.KeyParameter("param1")
		param2 = graph.KeyParameter("param2")
		param3 = graph.KeyParameter("param3")

		wild = graph.KeyWildcard{}
	)

	type Result graph.SearchResult[string]

	type searchTest struct {
		query    QueryItem
		expected searchResultList
	}

	subtests := []struct {
		paths    []PathItem
		searches []searchTest
	}{
		{
			// Basic constant matching

			paths: []PathItem{
				Path("val1", a, b),
				Path("val2", a, c),
			},
			searches: []searchTest{
				{
					query: Query("a", "b"),
					expected: searchResultList{
						{
							Value: "val1",
						},
					},
				},
				{
					query: Query("a", "c"),
					expected: searchResultList{
						{
							Value: "val2",
						},
					},
				},
				{
					query:    Query("a", "d"),
					expected: nil,
				},
			},
		},
		{
			// One of each key flavor

			paths: []PathItem{
				Path("valRoot"),
				Path("valA", a),
				Path("valParam", param1),
				Path("valWild", wild),
			},
			searches: []searchTest{
				{
					query: Query(),
					expected: searchResultList{
						{
							Value: "valRoot",
						},
						{
							Value: "valWild",
						},
					},
				},
				{
					query: Query("a"),
					expected: searchResultList{
						{
							Value: "valA",
						},
						{
							Value: "valParam",
							Parameters: map[string]string{
								"param1": "a",
							},
						},
						{
							Value: "valWild",
							Tail:  []string{"a"},
						},
					},
				},
			},
		},
		{
			// Only match requirement is minimum query length

			paths: []PathItem{
				Path("val0", wild),
				Path("val1", param1, wild),
				Path("val2", param1, param2, wild),
				Path("val3", param1, param2, param3, wild),
			},
			searches: []searchTest{
				{
					query: Query(),
					expected: searchResultList{
						{
							Value: "val0",
						},
					},
				},
				{
					query: Query("a"),
					expected: searchResultList{
						{
							Value: "val1",
							Parameters: map[string]string{
								"param1": "a",
							},
						},
						{
							Value: "val0",
							Tail:  []string{"a"},
						},
					},
				},
				{
					query: Query("a", "b"),
					expected: searchResultList{
						{
							Value: "val2",
							Parameters: map[string]string{
								"param1": "a",
								"param2": "b",
							},
						},
						{
							Value: "val1",
							Parameters: map[string]string{
								"param1": "a",
							},
							Tail: []string{"b"},
						},
						{
							Value: "val0",
							Tail:  []string{"a", "b"},
						},
					},
				},
				{
					query: Query("a", "b", "c"),
					expected: searchResultList{
						{
							Value: "val3",
							Parameters: map[string]string{
								"param1": "a",
								"param2": "b",
								"param3": "c",
							},
						},
						{
							Value: "val2",
							Parameters: map[string]string{
								"param1": "a",
								"param2": "b",
							},
							Tail: []string{"c"},
						},
						{
							Value: "val1",
							Parameters: map[string]string{
								"param1": "a",
							},
							Tail: []string{"b", "c"},
						},
						{
							Value: "val0",
							Tail:  []string{"a", "b", "c"},
						},
					},
				},
				{
					query: Query("a", "b", "c", "d"),
					expected: searchResultList{
						{
							Value: "val3",
							Parameters: map[string]string{
								"param1": "a",
								"param2": "b",
								"param3": "c",
							},
							Tail: []string{"d"},
						},
						{
							Value: "val2",
							Parameters: map[string]string{
								"param1": "a",
								"param2": "b",
							},
							Tail: []string{"c", "d"},
						},
						{
							Value: "val1",
							Parameters: map[string]string{
								"param1": "a",
							},
							Tail: []string{"b", "c", "d"},
						},
						{
							Value: "val0",
							Tail:  []string{"a", "b", "c", "d"},
						},
					},
				},
			},
		},
		{
			// Final constant key dictates parameter names

			paths: []PathItem{
				Path("valA", param1, param2, a),
				Path("valB", param2, param3, b),
				Path("valC", param3, param1, c),
			},
			searches: []searchTest{
				{
					query: Query("seg1", "seg2", "a"),
					expected: searchResultList{
						{
							Value: "valA",
							Parameters: map[string]string{
								"param1": "seg1",
								"param2": "seg2",
							},
						},
					},
				},
				{
					query: Query("seg1", "seg2", "b"),
					expected: searchResultList{
						{
							Value: "valB",
							Parameters: map[string]string{
								"param2": "seg1",
								"param3": "seg2",
							},
						},
					},
				},
				{
					query: Query("seg1", "seg2", "c"),
					expected: searchResultList{
						{
							Value: "valC",
							Parameters: map[string]string{
								"param3": "seg1",
								"param1": "seg2",
							},
						},
					},
				},
				{
					query:    Query("seg1", "seg2", "d"),
					expected: nil,
				},
			},
		},
	}

L:
	for _, subtest := range subtests {
		var tree Tree

		for _, path := range subtest.paths {
			err := tree.Add(path.Value, path.Keys...)

			if !assert.NoError(t, err, Info(path, &tree)) {
				continue L
			}
		}

		for _, searchTest := range subtest.searches {
			var (
				visitor      = new(searchVisitor)
				info         = Info(searchTest.query, &tree)
				expectedVals = searchTest.expected.values()
			)

			tree.Search(visitor, searchTest.query...)

			if !assert.Equal(t, expectedVals, visitor.actual.values(), info.Note("check values")) {
				continue
			}

			for i, expected := range searchTest.expected {
				var (
					val    = expected.Value
					actual = visitor.actual[i]
				)

				assert.Equal(t, expected.Parameters, actual.Parameters, info.Notef("check params - index %d (%s)", i, val))
				assert.Equal(t, expected.Tail, actual.Tail, info.Notef("check tail - index %d (%s)", i, val))
			}
		}
	}
}
