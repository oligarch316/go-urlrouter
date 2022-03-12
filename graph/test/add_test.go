package graphtest

import (
	"testing"

	"github.com/oligarch316/go-urlrouter/graph"
	"github.com/stretchr/testify/assert"
)

func TestGraphAddError(t *testing.T) {
	var (
		a = graph.KeyConstant("a")
		b = graph.KeyConstant("b")
		c = graph.KeyConstant("c")

		param1 = graph.KeyParameter("param1")
		param2 = graph.KeyParameter("param2")
		param3 = graph.KeyParameter("param3")

		wild = graph.KeyWildcard{}
	)

	t.Run("nil key", func(t *testing.T) {
		subtests := []PathItem{
			Path("someVal", nil),
			Path("someVal", a, nil),
			Path("someVal", param1, nil),
			Path("someVal", a, b, nil, c),
		}

		for _, path := range subtests {
			var (
				tree Tree
				info = Info(path, &tree)
			)

			err := tree.Add(path.Value, path.Keys...)
			assert.ErrorIs(t, err, graph.ErrNilKey, info)
		}
	})

	t.Run("invalid continuation", func(t *testing.T) {
		subtests := []struct {
			path                 PathItem
			expectedContinuation []graph.Key
		}{
			{
				path:                 Path("someVal", wild, a),
				expectedContinuation: []graph.Key{a},
			},
			{
				path:                 Path("someVal", wild, param1),
				expectedContinuation: []graph.Key{param1},
			},
			{
				path:                 Path("someVal", wild, wild),
				expectedContinuation: []graph.Key{wild},
			},
			{
				path:                 Path("someVal", a, wild, b, c),
				expectedContinuation: []graph.Key{b, c},
			},
			{
				path:                 Path("someVal", param1, wild, b, c),
				expectedContinuation: []graph.Key{b, c},
			},
		}

		for _, subtest := range subtests {
			var (
				tree      Tree
				targetErr graph.InvalidContinuationError
				info      = Info(subtest.path, &tree)
			)

			err := tree.Add(subtest.path.Value, subtest.path.Keys...)
			if !assert.ErrorAs(t, err, &targetErr, info.Note("check error type")) {
				continue
			}

			assert.Equal(t, subtest.expectedContinuation, targetErr.Continuation, info.Note("check continuation"))
		}
	})

	t.Run("duplicate value", func(t *testing.T) {
		subtests := []struct{ first, second PathItem }{
			{
				first:  Path("firstVal", wild),
				second: Path("secondVal", wild),
			},
			{
				first:  Path("firstVal", a, b, c),
				second: Path("secondVal", a, b, c),
			},
			{
				first:  Path("firstVal", a, b, c, wild),
				second: Path("secondVal", a, b, c, wild),
			},
			{
				first:  Path("firstVal", param1, param2),
				second: Path("secondVal", param1, param3),
			},
			{
				first:  Path("firstVal", param1, param2, wild),
				second: Path("secondVal", param1, param3, wild),
			},
			{
				first:  Path("firstVal", param1, param2, c),
				second: Path("secondVal", param1, param3, c),
			},
			{
				first:  Path("firstVal", param1, param2, c, wild),
				second: Path("secondVal", param1, param3, c, wild),
			},
		}

		for _, subtest := range subtests {
			var (
				tree      Tree
				targetErr graph.DuplicateValueError[string]
			)

			var (
				firstErr  = tree.Add(subtest.first.Value, subtest.first.Keys...)
				firstInfo = Info(subtest.first, &tree)
			)

			if !assert.NoError(t, firstErr, firstInfo.Note("check first add error")) {
				continue
			}

			var (
				secondErr  = tree.Add(subtest.second.Value, subtest.second.Keys...)
				secondInfo = Info(subtest.second, &tree)
			)

			if !assert.ErrorAs(t, secondErr, &targetErr, secondInfo.Note("check second add err")) {
				continue
			}

			assert.Equal(t, subtest.first.Value, targetErr.ExistingValue, secondInfo.Note("check existing error"))
		}
	})
}
