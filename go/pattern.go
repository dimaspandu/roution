package roution

import "strings"

const (
	SegmentStatic = "static"
	SegmentParam  = "param"
)

// Segment represents a single segment in a parsed route pattern.
//
// A segment is either:
//
//   - a literal path segment, represented by Type == SegmentStatic and Value
//   - a parameter segment, represented by Type == SegmentParam and Name
type Segment struct {
	Type  string
	Name  string
	Value string
}

// ParsedPattern is the compiled representation of a route pattern.
//
// The wildcard pattern "*" is represented by IsWildCard == true and an empty
// segment list.
type ParsedPattern struct {
	Raw        string
	IsWildCard bool
	Segments   []Segment
}

// ParsePattern parses a route pattern into structured segments.
//
// The wildcard pattern "*" is recognized as a standalone fallback route.
//
// Static segments become:
//
//	Segment{
//	    Type: SegmentStatic,
//	    Value: "...",
//	}
//
// Parameter segments beginning with ':' become:
//
//	Segment{
//	    Type: SegmentParam,
//	    Name: "...",
//	}
//
// Examples:
//
//	"/articles"
//	    -> [{Type: SegmentStatic, Value: "articles"}]
//
//	"/articles/:slug"
//	    -> [
//	         {Type: SegmentStatic, Value: "articles"},
//	         {Type: SegmentParam, Name: "slug"},
//	       ]
func ParsePattern(pattern string) ParsedPattern {
	if pattern == "*" {
		return ParsedPattern{
			Raw:        pattern,
			IsWildCard: true,
			Segments:   []Segment{},
		}
	}

	parts := strings.Split(pattern, "/")
	segments := make([]Segment, 0, len(parts))

	for _, part := range parts {
		if part == "" {
			continue
		}

		if part[0] == ':' {
			segments = append(segments, Segment{
				Type: SegmentParam,
				Name: part[1:],
			})
			continue
		}

		segments = append(segments, Segment{
			Type:  SegmentStatic,
			Value: part,
		})
	}

	return ParsedPattern{
		Raw:      pattern,
		Segments: segments,
	}
}