package component

import "github.com/oligarch316/go-urlrouter/graph"

type Router[V any] struct {
	Decoder   KeyDecoder
	Segmenter PatternSegmenter
	Tree      graph.Tree[V]
}

func (r *Router[V]) Add(pattern string, value V) error {
	segs, err := r.Segmenter.Segment(pattern)
	if err != nil {
		return err
	}

	keys, err := r.Decoder.Decode(segs)
	if err != nil {
		return err
	}

	return r.Tree.Add(value, keys...)
}

func (r *Router[V]) Search(searcher graph.Searcher[V], query string) error {
	segs, err := r.Segmenter.Segment(query)
	if err != nil {
		return err
	}

	r.Tree.Search(searcher, segs...)
	return nil
}

func (r *Router[V]) SearchFunc(searcher func(result *graph.SearchResult[V]) (done bool), query string) error {
	return r.Search(graph.SearcherFunc[V](searcher), query)
}

func (r *Router[V]) Walk(walker graph.Walker[V])               { r.Tree.Walk(walker) }
func (r *Router[V]) WalkFunc(walker func(value V) (done bool)) { r.Walk(graph.WalkerFunc[V](walker)) }
