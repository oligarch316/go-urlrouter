package component

type (
	HostRouter[V any] struct{ router[V] }
	PathRouter[V any] struct{ router[V] }
)

func NewHostRouter[V any](opts ...Option[V]) *HostRouter[V] {
	return &HostRouter[V]{newRouter(segmentHost, opts)}
}

func NewPathRouter[V any](opts ...Option[V]) *PathRouter[V] {
	return &PathRouter[V]{newRouter(segmentPath, opts)}
}
