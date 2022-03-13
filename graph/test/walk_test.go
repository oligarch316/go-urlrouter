package graphtest

import (
	"testing"

	"github.com/oligarch316/go-urlrouter/graph"
	"github.com/stretchr/testify/assert"
)

type walkVisitor struct{ actual []string }

func (wv *walkVisitor) VisitWalk(value string) bool {
	wv.actual = append(wv.actual, value)
	return false
}

func TestGraphWalk(t *testing.T) {
	var (
		a     = graph.KeyConstant("a")
		param = graph.KeyParameter("param")
		wild  = graph.KeyWildcard{}
	)

	subtests := []struct {
		paths    []PathItem
		expected []string
	}{
		{
			paths:    []PathItem{},
			expected: []string{},
		},
		{
			paths: []PathItem{
				Path("valRoot"),
				Path("valRootWild", wild),

				Path("valA", a),
				Path("valAWild", a, wild),

				Path("valParam", param),
				Path("valParamWild", param, wild),

				Path("valAParam", a, param),
				Path("valAParamWild", a, param, wild),

				Path("valParamA", param, a),
				Path("valParamAWild", param, a, wild),
			},
			expected: []string{
				"valRoot",
				"valRootWild",

				"valA",
				"valAWild",

				"valParam",
				"valParamWild",

				"valAParam",
				"valAParamWild",

				"valParamA",
				"valParamAWild",
			},
		},
	}

L:
	for _, subtest := range subtests {
		var (
			tree    Tree
			visitor = new(walkVisitor)
		)

		for _, path := range subtest.paths {
			err := tree.Add(path.Value, path.Keys...)

			if !assert.NoError(t, err, Info(path, &tree)) {
				continue L
			}
		}

		tree.Walk(visitor)
		assert.ElementsMatch(t, subtest.expected, visitor.actual, Info(&tree))
	}
}
