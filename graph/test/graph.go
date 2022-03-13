package graphtest

import (
	"fmt"
	"strings"

	"github.com/oligarch316/go-urlrouter/graph"
	"github.com/oligarch316/go-urlrouter/graph/memoized"
)

type InfoMessage struct {
	Format string
	Args   []interface{}
}

func (im InfoMessage) String() string { return fmt.Sprintf(im.Format, im.Args...) }

type InfoList []fmt.Stringer

func (il InfoList) Note(message string) InfoList { return il.Notef(message) }

func (il InfoList) Notef(format string, a ...interface{}) InfoList {
	message := InfoMessage{Format: format, Args: a}
	return append(InfoList{message}, il...)
}

func (il InfoList) String() string {
	strs := make([]string, len(il))
	for i, item := range il {
		strs[i] = item.String()
	}
	return strings.Join(strs, "\n")
}

type PathItem struct {
	Keys  []graph.Key
	Value string
}

func (pi PathItem) String() string {
	return fmt.Sprintf("--------\nPath:\n> %s", graph.FormatPath(pi.Value, pi.Keys...))
}

type QueryItem []string

func (qi QueryItem) String() string {
	dataStr := "<empty>"
	if len(qi) > 0 {
		dataStr = graph.FormatQuery(qi...)
	}

	return fmt.Sprintf("--------\nQuery:\n> %s", dataStr)
}

type Tree struct{ memoized.Tree[string] }

func (t *Tree) String() string {
	var paths []string
	t.Memoized.WalkFunc(func(memo memoized.Memo[string]) bool {
		paths = append(paths, "> "+memo.String())
		return false
	})

	dataStr := "> <empty>"
	if len(paths) > 0 {
		dataStr = strings.Join(paths, "\n")
	}

	return fmt.Sprintf("--------\nTree:\n%s", dataStr)
}

func Info(items ...fmt.Stringer) InfoList           { return InfoList(items) }
func Query(items ...string) QueryItem               { return QueryItem(items) }
func Path(value string, keys ...graph.Key) PathItem { return PathItem{Keys: keys, Value: value} }
