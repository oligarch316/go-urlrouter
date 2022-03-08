package component

import (
	"errors"

	"github.com/oligarch316/go-urlrouter/graph"
)

var ErrNotFound = errors.New("not found")

type Option[V any] func(*Params[V])

type Params[V any] struct {
	Decoder KeyDecoder
	Tree    interface {
		Add(V, ...graph.Key) error
		Search(...string) *graph.Result[V]
	}
}

type router[V any] struct {
	params  Params[V]
	segment segmenter
}

func newRouter[V any](segment segmenter, opts []Option[V]) router[V] {
	params := Params[V]{
		Decoder: DefaultKeyDecoder,
		Tree:    new(graph.Tree[V]),
	}

	for _, opt := range opts {
		opt(&params)
	}

	return router[V]{params: params, segment: segment}
}

func (r *router[V]) Add(item string, val V) error {
	segs, err := r.segment(item)
	if err != nil {
		return err
	}

	keys := r.params.Decoder.Decode(segs)

	return r.params.Tree.Add(val, keys...)
}

func (r *router[V]) Search(item string) (*graph.Result[V], error) {
	segs, err := r.segment(item)
	if err != nil {
		return nil, err
	}

	res := r.params.Tree.Search(segs...)
	if res == nil {
		err = ErrNotFound
	}

	return res, err
}
