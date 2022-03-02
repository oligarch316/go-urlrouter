package graph

import (
	"errors"
	"fmt"
)

var ErrInternal = errors.New("internal")

type internalError string

func internalErrorf(format string, a ...interface{}) internalError {
	message := fmt.Sprintf(format, a...)
	return internalError(message)
}

func (ie internalError) Error() string        { return string(ie) }
func (ie internalError) Is(target error) bool { return target == ErrInternal }
