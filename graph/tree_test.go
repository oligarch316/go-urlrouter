package graph_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/oligarch316/go-urlrouter/graph"
	"github.com/stretchr/testify/assert"
)

type Action struct {
	name string
	data []string
}

func (a Action) String() string {
	dataStr := "<empty>"
	if len(a.data) > 0 {
		dataStr = strings.Join(a.data, "→")
	}

	return fmt.Sprintf("%s %s", a.name, dataStr)
}

type Log []Action

func (l Log) Title(info ...string) string {
	strs := append(info, "--------", "actions:", l.String())
	return strings.Join(strs, "\n")
}

func (l Log) String() string {
	if len(l) == 0 {
		return "> <no actions>"
	}

	res := "> " + l[0].String()
	for _, action := range l[1:] {
		res += "\n> "
		res += action.String()
	}
	return res
}

type LoggedTree struct {
	tree graph.Tree
	Log
}

func (lt *LoggedTree) Add(val graph.Value, keys ...graph.Key) error {
	action := Action{name: "add"}
	for _, key := range keys {
		str := "<nil>"
		if key != nil {
			str = key.String()
		}

		action.data = append(action.data, str)
	}

	lt.Log = append(lt.Log, action)
	return lt.tree.Add(val, keys...)
}

func (lt *LoggedTree) Search(segs ...graph.Segment) (*graph.Result, Log) {
	action := Action{name: "search"}
	for _, seg := range segs {
		action.data = append(action.data, string(seg))
	}

	return lt.tree.Search(segs...), append(lt.Log, action)
}

func TestTreeAddError(t *testing.T) {
	var (
		a = graph.KeyConstant("A")
		b = graph.KeyConstant("B")
		c = graph.KeyConstant("C")

		param1 = graph.KeyParameter("param1")
		param2 = graph.KeyParameter("param2")
		param3 = graph.KeyParameter("param3")

		wild graph.KeyWildcard
	)

	t.Run("nil key", func(t *testing.T) {
		subtests := [][]graph.Key{
			{nil},
			{a, nil},
			{param1, nil},
			{a, b, nil, c},
		}

		for _, subtest := range subtests {
			var tree LoggedTree

			assert.ErrorIs(t, tree.Add("someVal", subtest...), graph.ErrNilKey, tree.Title())
		}
	})

	t.Run("non-terminal wildcard", func(t *testing.T) {
		subtests := []struct{ keys, expectedContinuation []graph.Key }{
			{
				keys:                 []graph.Key{wild, a},
				expectedContinuation: []graph.Key{a},
			},
			{
				keys:                 []graph.Key{wild, param1},
				expectedContinuation: []graph.Key{param1},
			},
			{
				keys:                 []graph.Key{wild, wild},
				expectedContinuation: []graph.Key{wild},
			},
			{
				keys:                 []graph.Key{a, wild, b, c},
				expectedContinuation: []graph.Key{b, c},
			},
			{
				keys:                 []graph.Key{param1, wild, b, c},
				expectedContinuation: []graph.Key{b, c},
			},
		}

		for _, subtest := range subtests {
			var (
				tree      LoggedTree
				targetErr graph.InvalidContinuationError
			)

			if !assert.ErrorAs(t, tree.Add("someVal", subtest.keys...), &targetErr, tree.Title()) {
				continue
			}

			assert.Equal(t, subtest.expectedContinuation, targetErr.Continuation, tree.Title())
		}
	})

	t.Run("duplicate", func(t *testing.T) {
		subtests := []struct{ first, second []graph.Key }{
			{
				first:  []graph.Key{wild},
				second: []graph.Key{wild},
			},
			{
				first:  []graph.Key{a, b, c},
				second: []graph.Key{a, b, c},
			},
			{
				first:  []graph.Key{a, b, c, wild},
				second: []graph.Key{a, b, c, wild},
			},
			{
				first:  []graph.Key{param1, param2},
				second: []graph.Key{param1, param3},
			},
			{
				first:  []graph.Key{param1, param2, wild},
				second: []graph.Key{param1, param3, wild},
			},
			{
				first:  []graph.Key{param1, param2, c},
				second: []graph.Key{param1, param3, c},
			},
			{
				first:  []graph.Key{param1, param2, c, wild},
				second: []graph.Key{param1, param3, c, wild},
			},
		}

		for _, subtest := range subtests {
			var (
				tree      LoggedTree
				targetErr graph.DuplicateValueError
			)

			if !assert.NoError(t, tree.Add("firstVal", subtest.first...), tree.Title()) {
				continue
			}

			if !assert.ErrorAs(t, tree.Add("secondVal", subtest.second...), &targetErr, tree.Title()) {
				continue
			}

			assert.Equal(t, "firstVal", targetErr.ExistingValue, tree.Title())
		}
	})
}

func TestTreeSearchFailure(t *testing.T) {
	// TODO
	t.Skip("TODO")
}

func TestTreeSearchSuccess(t *testing.T) {
	var (
		a = graph.KeyConstant("a")
		b = graph.KeyConstant("b")
		c = graph.KeyConstant("c")

		param1 = graph.KeyParameter("param1")
		param2 = graph.KeyParameter("param2")
		param3 = graph.KeyParameter("param3")

		wild graph.KeyWildcard
	)

	type (
		addItem struct {
			value graph.Value
			keys  []graph.Key
		}

		expectItem struct {
			value  graph.Value
			params map[graph.Parameter]graph.Segment
			tail   []graph.Segment
		}

		searchItem struct {
			query  []graph.Segment
			expect expectItem
		}
	)

	type expectF func(*expectItem)

	var (
		expectValue = func(val graph.Value) expectF { return func(ei *expectItem) { ei.value = val } }
		expectParam = func(k graph.Parameter, v graph.Segment) expectF { return func(ei *expectItem) { ei.params[k] = v } }
		expectTail  = func(tail ...graph.Segment) expectF { return func(ei *expectItem) { ei.tail = tail } }

		add = func(val graph.Value, keys ...graph.Key) addItem { return addItem{value: val, keys: keys} }

		search = func(query ...graph.Segment) func(...expectF) searchItem {
			item := searchItem{
				query: query,
				expect: expectItem{
					params: make(map[graph.Parameter]graph.Segment),
				},
			}

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
					// TODO: Fails due to nil <-> []Segment{} mismatch
					expectValue("val1"),
					expectParam("param1", "a"),
				),
				search("a", "b")(
					// TODO: Fails due to nil <-> []Segment{} mismatch
					expectValue("val2"),
					expectParam("param1", "a"),
					expectParam("param2", "b"),
				),
				search("a", "b", "c")(
					// TODO: Fails due to nil <-> []Segment{} mismatch
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
		var tree LoggedTree

		for _, item := range subtest.addItems {
			if !assert.NoError(t, tree.Add(item.value, item.keys...), tree.Title()) {
				continue L
			}
		}

		for _, item := range subtest.searchItems {
			result, log := tree.Search(item.query...)

			if !assert.NotNil(t, result, log.Title()) {
				continue
			}

			if !assert.Equal(t, item.expect.value, result.Value, log.Title("check value")) {
				continue
			}

			assert.Equal(t, item.expect.params, result.Parameters, log.Title("check params"))
			assert.Equal(t, item.expect.tail, result.Tail, log.Title("check tail"))
		}
	}
}
