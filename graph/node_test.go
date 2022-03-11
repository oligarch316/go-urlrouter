package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGraphNodeAddInternalError(t *testing.T) {
	var (
		node  nodeParameter[string]
		keys  = []Key{KeyParameter("someParam")}
		state = stateAdd[string]{value: "someVal"}
	)

	assert.ErrorIs(t, node.add(keys, state), ErrInternal)
}
