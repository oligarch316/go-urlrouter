package component

import "github.com/oligarch316/go-urlrouter/graph/priority"

var (
	DefaultKeyDecoder    KeyDecodeFunc        = decodeKeyDefault
	DefaultHostSegmenter PatternSegmenterFunc = segmentHostDefault
	DefaultPathSegmenter PatternSegmenterFunc = segmentPathDefault
)

func NewHostRouter[V any](opts ...func(*Router[V])) *Router[V] {
	res := &Router[V]{
		Decoder:   DefaultKeyDecoder,
		Segmenter: DefaultHostSegmenter,
		Tree:      new(priority.Tree[V]),
	}

	for _, opt := range opts {
		opt(res)
	}

	return res
}

func NewPathRouter[V any](opts ...func(*Router[V])) *Router[V] {
	res := &Router[V]{
		Decoder:   DefaultKeyDecoder,
		Segmenter: DefaultPathSegmenter,
		Tree:      new(priority.Tree[V]),
	}

	for _, opt := range opts {
		opt(res)
	}

	return res
}
