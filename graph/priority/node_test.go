package priority

import (
	"testing"

	"github.com/oligarch316/go-urlrouter/graph"
	"github.com/stretchr/testify/assert"
)

func TestGraphPriorityInternalError(t *testing.T) {
	var (
		node  nodeParameter[string]
		keys  = []graph.Key{graph.KeyParameter("someParam")}
		state = stateAdd[string]{value: "someVal"}
	)

	assert.ErrorIs(t, node.add(keys, state), graph.ErrInternal)
}
