package graph

import "errors"

var (
	ErrInternal = errors.New("internal")
	ErrNilKey   = errors.New("nil key")
)

type DuplicateValueError[V any] struct{ ExistingValue V }

func (dve DuplicateValueError[_]) Error() string { return "dupliate value" }

type InvalidContinuationError struct{ Continuation []Key }

func (ice InvalidContinuationError) Error() string { return "invalid continuation" }
