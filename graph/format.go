package graph

import (
	"fmt"
	"strings"
)

func formatChain(items []string) string { return strings.Join(items, "→") }

func FormatKey(key Key) string {
	if key == nil {
		return "<nil>"
	}
	return key.String()
}

func FormatValue(val interface{}) string {
	if stringer, ok := val.(fmt.Stringer); ok {
		return stringer.String()
	}

	return fmt.Sprintf("value(%v)", val)
}

func FormatQuery(segs ...Segment) string {
	segChain := make([]string, len(segs))
	for i, seg := range segs {
		segChain[i] = string(seg)
	}

	return formatChain(segChain)
}

func FormatPath(val interface{}, keys ...Key) string {
	keyChain := make([]string, len(keys))
	for i, key := range keys {
		keyChain[i] = FormatKey(key)
	}

	return formatChain(append(keyChain, FormatValue(val)))
}
