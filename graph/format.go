package graph

import (
	"fmt"
	"strings"
)

func formatChain(items []string) string { return strings.Join(items, "→") }

func FormatValue(val interface{}) string {
	if stringer, ok := val.(fmt.Stringer); ok {
		return stringer.String()
	}

	return fmt.Sprintf("value(%v)", val)
}

func ForamtQuery(segs ...string) string { return formatChain(segs) }

func FormatPath(val interface{}, keys ...Key) string {
	keyChain := make([]string, len(keys))
	for i, key := range keys {
		keyChain[i] = key.String()
	}

	return formatChain(append(keyChain, FormatValue(val)))
}
