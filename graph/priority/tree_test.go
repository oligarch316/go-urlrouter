package priority_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/oligarch316/go-urlrouter/graph/memoized"
	"github.com/oligarch316/go-urlrouter/graph/priority"
	"github.com/stretchr/testify/assert"
)

// TODO: Deprecated, Remove me

type stringerFunc func() string

func (sf stringerFunc) String() string { return sf() }

type stringerList []fmt.Stringer

func (sl stringerList) String() string {
	strs := make([]string, len(sl))
	for i, item := range sl {
		strs[i] = item.String()
	}
	return strings.Join(strs, "\n")
}

func info(name string, data interface{}) stringerFunc {
	return func() string {
		return fmt.Sprintf("--------\n%s:\n%s", name, data)
	}
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

func (tt *TestTree) Add(val string, keys ...priority.Key) (error, fmt.Stringer) {
	var (
		err     = tt.Tree.Add(val, keys...)
		addInfo = stringerList{
			info("Add", "> "+priority.FormatPath(val, keys...)),
			info("Tree", tt),
		}
	)

	return err, addInfo
}

func (tt *TestTree) Search(segs ...string) (*priority.Result[string], fmt.Stringer) {
	var (
		res        = tt.Tree.Search(segs...)
		searchInfo = stringerList{
			info("Search", "> "+priority.FormatQuery(segs...)),
			info("Tree", tt),
		}
	)

	return res, searchInfo
}

func TestGraphTreeAddError(t *testing.T) {
	var (
		a = priority.KeyConstant("A")
		b = priority.KeyConstant("B")
		c = priority.KeyConstant("C")

		param1 = priority.KeyParameter("param1")
		param2 = priority.KeyParameter("param2")
		param3 = priority.KeyParameter("param3")

		wild priority.KeyWildcard
	)

	t.Run("nil key", func(t *testing.T) {
		subtests := [][]priority.Key{
			{nil},
			{a, nil},
			{param1, nil},
			{a, b, nil, c},
		}

		for _, subtest := range subtests {
			var tree TestTree

			err, info := tree.Add("someVal", subtest...)
			assert.ErrorIs(t, err, priority.ErrNilKey, info)
		}
	})

	t.Run("non-terminal wildcard", func(t *testing.T) {
		subtests := []struct{ keys, expectedContinuation []priority.Key }{
			{
				keys:                 []priority.Key{wild, a},
				expectedContinuation: []priority.Key{a},
			},
			{
				keys:                 []priority.Key{wild, param1},
				expectedContinuation: []priority.Key{param1},
			},
			{
				keys:                 []priority.Key{wild, wild},
				expectedContinuation: []priority.Key{wild},
			},
			{
				keys:                 []priority.Key{a, wild, b, c},
				expectedContinuation: []priority.Key{b, c},
			},
			{
				keys:                 []priority.Key{param1, wild, b, c},
				expectedContinuation: []priority.Key{b, c},
			},
		}

		for _, subtest := range subtests {
			var (
				tree      TestTree
				targetErr priority.InvalidContinuationError
			)

			err, info := tree.Add("someVal", subtest.keys...)
			if !assert.ErrorAs(t, err, &targetErr, "check error\n%s", info) {
				continue
			}

			assert.Equal(t, subtest.expectedContinuation, targetErr.Continuation, "check continuation\n%s", info)
		}
	})

	t.Run("duplicate", func(t *testing.T) {
		subtests := []struct{ first, second []priority.Key }{
			{
				first:  []priority.Key{wild},
				second: []priority.Key{wild},
			},
			{
				first:  []priority.Key{a, b, c},
				second: []priority.Key{a, b, c},
			},
			{
				first:  []priority.Key{a, b, c, wild},
				second: []priority.Key{a, b, c, wild},
			},
			{
				first:  []priority.Key{param1, param2},
				second: []priority.Key{param1, param3},
			},
			{
				first:  []priority.Key{param1, param2, wild},
				second: []priority.Key{param1, param3, wild},
			},
			{
				first:  []priority.Key{param1, param2, c},
				second: []priority.Key{param1, param3, c},
			},
			{
				first:  []priority.Key{param1, param2, c, wild},
				second: []priority.Key{param1, param3, c, wild},
			},
		}

		for _, subtest := range subtests {
			var (
				tree      TestTree
				targetErr priority.DuplicateValueError[string]
			)

			err, info := tree.Add("firstVal", subtest.first...)
			if !assert.NoError(t, err, "check first add error\n%s", info) {
				continue
			}

			err, info = tree.Add("secondVal", subtest.second...)
			if !assert.ErrorAs(t, err, &targetErr, "check second add error\n%s", info) {
				continue
			}

			assert.Equal(t, "firstVal", targetErr.ExistingValue, "check existing value", info)
		}
	})
}

func TestGraphTreeSearchFailure(t *testing.T) {
	// TODO
	t.Skip("TODO")
}

func TestGraphTreeSearchSuccess(t *testing.T) {
	var (
		a = priority.KeyConstant("a")
		b = priority.KeyConstant("b")
		c = priority.KeyConstant("c")

		param1 = priority.KeyParameter("param1")
		param2 = priority.KeyParameter("param2")
		param3 = priority.KeyParameter("param3")

		wild priority.KeyWildcard
	)

	type (
		addItem struct {
			value string
			keys  []priority.Key
		}

		expectItem struct {
			value  string
			params map[string]string
			tail   []string
		}

		searchItem struct {
			query  []string
			expect expectItem
		}
	)

	type expectF func(*expectItem)

	var (
		expectValue = func(val string) expectF {
			return func(ei *expectItem) { ei.value = val }
		}

		expectParam = func(k string, v string) expectF {
			return func(ei *expectItem) {
				if ei.params == nil {
					ei.params = make(map[string]string)
				}
				ei.params[k] = v
			}
		}

		expectTail = func(tail ...string) expectF {
			return func(ei *expectItem) { ei.tail = tail }
		}

		add = func(val string, keys ...priority.Key) addItem {
			return addItem{value: val, keys: keys}
		}

		search = func(query ...string) func(...expectF) searchItem {
			item := searchItem{query: query}

			return func(expectations ...expectF) searchItem {
				for _, e := range expectations {
					e(&item.expect)
				}
				return item
			}
		}
	)

	subtests := []struct {
		addItems    []addItem
		searchItems []searchItem
	}{
		// TODO: Systematize these tests, need to exhaust scenarios involving
		// - Correct prioritization: const > param > wildcard
		// - Correct parameter mapping
		// - Correct wildcard tails
		{
			addItems: []addItem{
				add("valB", a, b),
				add("valC", a, c),
			},
			searchItems: []searchItem{
				search("a", "b")(
					expectValue("valB"),
				),
				search("a", "c")(
					expectValue("valC"),
				),
			},
		},
		{
			addItems: []addItem{
				add("valRoot"),
				add("valConst", a),
				add("valParam", param1),
				add("valWild", wild),
			},
			searchItems: []searchItem{
				search()(
					expectValue("valRoot"),
				),
				search("a")(
					expectValue("valConst"),
				),
				search("b")(
					expectValue("valParam"),
					expectParam("param1", "b"),
				),
				search("a", "b")(
					expectValue("valWild"),
					expectTail("a", "b"),
				),
			},
		},
		{
			addItems: []addItem{
				add("val0", wild),
				add("val1", param1, wild),
				add("val2", param1, param2, wild),
				add("val3", param1, param2, param3, wild),
			},
			searchItems: []searchItem{
				search()(
					expectValue("val0"),
				),
				search("a")(
					expectValue("val1"),
					expectParam("param1", "a"),
				),
				search("a", "b")(
					expectValue("val2"),
					expectParam("param1", "a"),
					expectParam("param2", "b"),
				),
				search("a", "b", "c")(
					expectValue("val3"),
					expectParam("param1", "a"),
					expectParam("param2", "b"),
					expectParam("param3", "c"),
				),
				search("a", "b", "c", "d")(
					expectValue("val3"),
					expectParam("param1", "a"),
					expectParam("param2", "b"),
					expectParam("param3", "c"),
					expectTail("d"),
				),
			},
		},
		{
			addItems: []addItem{
				add("valA", param1, param2, a),
				add("valB", param2, param3, b),
				add("valC", param3, param1, c),
			},
			searchItems: []searchItem{
				search("seg1", "seg2", "a")(
					expectValue("valA"),
					expectParam("param1", "seg1"),
					expectParam("param2", "seg2"),
				),
				search("seg1", "seg2", "b")(
					expectValue("valB"),
					expectParam("param2", "seg1"),
					expectParam("param3", "seg2"),
				),
				search("seg1", "seg2", "c")(
					expectValue("valC"),
					expectParam("param3", "seg1"),
					expectParam("param1", "seg2"),
				),
			},
		},
	}

L:
	for _, subtest := range subtests {
		var tree TestTree

		for _, item := range subtest.addItems {
			err, info := tree.Add(item.value, item.keys...)

			if !assert.NoError(t, err, info) {
				continue L
			}
		}

		for _, item := range subtest.searchItems {
			result, info := tree.Search(item.query...)
			if !assert.NotNil(t, result, "check nil\n%s", info) {
				continue
			}

			if !assert.Equal(t, item.expect.value, result.Value, "check value\n%s", info) {
				continue
			}

			assert.Equal(t, item.expect.params, result.Parameters, "check params\n%s", info)
			assert.Equal(t, item.expect.tail, result.Tail, "check tail\n%s", info)
		}
	}
}