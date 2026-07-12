package roution

import (
	"net/url"
	"strings"
)

// Normalize removes the query string and URL fragment from an incoming
// pathname.
//
// Everything after the first '#' or '?' is discarded. No other transformation
// is applied, so the returned value remains the original pathname.
//
// Examples:
//
//	/articles?page=1#comments -> /articles
//	/articles#comments        -> /articles
//	/articles?page=1          -> /articles
func Normalize(input string) string {
	path := input

	if index := strings.IndexByte(path, '#'); index >= 0 {
		path = path[:index]
	}

	if index := strings.IndexByte(path, '?'); index >= 0 {
		path = path[:index]
	}

	return path
}

// ParseQuery parses the query string of a pathname.
//
// The returned map stores single values as string and repeated values as
// []string.
//
// Examples:
//
//	?a=1
//	    -> "a": "1"
//
//	?a=1&a=2
//	    -> "a": []string{"1", "2"}
func ParseQuery(input string) map[string]any {
	path := input

	if index := strings.IndexByte(path, '#'); index >= 0 {
		path = path[:index]
	}

	index := strings.IndexByte(path, '?')
	if index < 0 {
		return map[string]any{}
	}

	query := path[index+1:]
	if query == "" {
		return map[string]any{}
	}

	values, err := url.ParseQuery(query)
	if err != nil {
		return map[string]any{}
	}

	result := make(map[string]any, len(values))

	for key, value := range values {
		if len(value) == 1 {
			result[key] = value[0]
			continue
		}

		result[key] = value
	}

	return result
}

// SplitSegments splits a pathname into non-empty path segments.
//
// Leading and trailing slashes are ignored.
//
// Examples:
//
//	"/users/42"
//	    -> []string{"users", "42"}
//
//	"/"
//	    -> []string{}
func SplitSegments(pathname string) []string {
	parts := strings.Split(pathname, "/")
	segments := make([]string, 0, len(parts))

	for _, part := range parts {
		if part == "" {
			continue
		}

		segments = append(segments, part)
	}

	return segments
}