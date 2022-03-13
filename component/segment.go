package component

import (
	"errors"
	"fmt"
	"strings"
)

const (
	segHostSep = '.'
	segPathSep = '/'
)

var (
	ErrInvalidHost = errors.New("invalid host")
	ErrInvalidPath = errors.New("invalid path")
)

type PatternSegmenter interface {
	Segment(string) ([]string, error)
}

type PatternSegmenterFunc func(string) ([]string, error)

func (psf PatternSegmenterFunc) Segment(pattern string) ([]string, error) { return psf(pattern) }

func segmentHostDefault(pattern string) ([]string, error) {
	if pattern == "" {
		return nil, nil
	}

	var (
		fields  = strings.Split(pattern, string(segHostSep))
		nFields = len(fields)
		res     = make([]string, nFields)
	)

	for i, field := range fields {
		res[(nFields-1)-i] = field
	}

	return res, nil
}

func segmentPathDefault(pattern string) ([]string, error) {
	if pattern == "" || pattern[0] != segPathSep {
		return nil, fmt.Errorf("%w: missing leading slash", ErrInvalidPath)
	}

	pattern = pattern[1:]

	if pattern == "" {
		return nil, nil
	}

	if lastIdx := len(pattern) - 1; pattern[lastIdx] == segPathSep {
		if lastIdx == 0 {
			return nil, nil
		}

		pattern = pattern[:lastIdx]
	}

	return strings.Split(pattern, string(segPathSep)), nil
}
