package component

import (
	"errors"
	"fmt"

	"github.com/oligarch316/go-urlrouter/graph"
)

const (
	decPrefixParam = ':'
	decPrefixWild  = '*'
)

var ErrInvalidSegment = errors.New("invalid segment")

// NOTE: KeyDecoder.Decode() is defined to process slice -> slice because:
//
// The ultimate plan is to route full urls by embedding PathRouter as the
// generic V type in HostRouter. Add() calls to the HostRouter will hopefully
// hide some "upsert" logic by responding to DuplicateValue errors appropriately.
// The problem will be managing multiple parameterized hosts pointing to the
// same location (and same PathRouter). Naively re-using an existing PathRouter
// will ignore different param names in the second parameterized host.
//
// The current thought is that if we can manage to make parameter types
// (key type in the Result.Parameter map) generic as originally intended,
// HostRouter[V=PathRouter[...]] becomes HostRouter[V=PathRouter[...], P=int].
// We then set parameter keys for HostRouter to general purpose incrementing ints while
// storing a map[int]<orig param key> as part of the PathRouter value V.
// Then we're free to reuse exiting PathRouter values found during HostRouter.Add(...)
// to implement "upsert", and a correct host parameter map can be constructed by
// essentially zipping HostRouter.Search(...) -> Result.Parameters (a map[int]<param value>)
// with PathRouter.Search(...) -> Result.Value.StoredHostParamMap (a map[int]<param key>)
//
// Doing such necessitates that KeyDecoder.Decode() can be stateful (scoped param counter),
// thus []string -> []graph.Key rather than string -> graph.Key

type KeyDecoder interface {
	Decode([]string) ([]graph.Key, error)
}

type KeyDecoderFunc func([]string) ([]graph.Key, error)

func (kdf KeyDecoderFunc) Decode(segs []string) ([]graph.Key, error) { return kdf(segs) }

type KeyDecodeFunc func(string) (graph.Key, error)

func (kdf KeyDecodeFunc) Decode(segs []string) ([]graph.Key, error) {
	var (
		err error
		res = make([]graph.Key, len(segs))
	)

	for i, seg := range segs {
		if res[i], err = kdf(seg); err != nil {
			return nil, err
		}
	}

	return res, nil
}

func decodeKeyDefault(raw string) (graph.Key, error) {
	rawLen := len(raw)
	if rawLen == 0 {
		return nil, fmt.Errorf("%w: empty segment", ErrInvalidSegment)
	}

	switch raw[0] {
	case decPrefixParam:
		if rawLen < 2 {
			return nil, fmt.Errorf("%w: empty parameter name", ErrInvalidSegment)
		}

		return graph.KeyParameter(raw[1:]), nil
	case decPrefixWild:
		return graph.KeyWildcard{}, nil
	}

	return graph.KeyConstant(raw), nil
}
