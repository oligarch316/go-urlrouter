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

type segmenter func(string) ([]string, error)

func segmentHost(host string) ([]string, error) {
	if host == "" {
		return []string{}, nil
	}

	var (
		sections  = strings.Split(host, string(segHostSep))
		nSections = len(sections)
		res       = make([]string, nSections)
	)

	for i, section := range sections {
		if section == "" {
			return nil, fmt.Errorf("%w: empty segment", ErrInvalidHost)
		}

		res[(nSections-1)-i] = section
	}

	return res, nil
}

func segmentPath(path string) ([]string, error) {
	if len(path) == 0 || path[0] != segPathSep {
		return nil, fmt.Errorf("%w: missing leading slash", ErrInvalidPath)
	}

	path = path[1:]

	switch len(path) {
	case 0:
		return []string{}, nil
	case 1:
		if path[0] == segPathSep {
			return nil, fmt.Errorf("%w: empty segment", ErrInvalidPath)
		}
	default:
		if path[len(path)-1] == segPathSep {
			path = path[:len(path)-1]
		}
	}

	var (
		err      error
		sections = strings.Split(path, string(segPathSep))
	)

	for _, section := range sections {
		if section == "" {
			err = fmt.Errorf("%w: empty segment", ErrInvalidPath)
			break
		}
	}

	return sections, err
}
