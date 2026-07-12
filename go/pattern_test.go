package roution

import (
	"reflect"
	"testing"
)

func TestParsePattern(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		want    ParsedPattern
	}{
		{
			name:    "parses wildcard",
			pattern: "*",
			want: ParsedPattern{
				Raw:        "*",
				IsWildCard: true,
				Segments:   []Segment{},
			},
		},
		{
			name:    "parses static root",
			pattern: "/",
			want: ParsedPattern{
				Raw:      "/",
				Segments: []Segment{},
			},
		},
		{
			name:    "parses static segments",
			pattern: "/articles",
			want: ParsedPattern{
				Raw: "/articles",
				Segments: []Segment{
					{
						Type:  SegmentStatic,
						Value: "articles",
					},
				},
			},
		},
		{
			name:    "parses dynamic segments",
			pattern: "/articles/:slug",
			want: ParsedPattern{
				Raw: "/articles/:slug",
				Segments: []Segment{
					{
						Type:  SegmentStatic,
						Value: "articles",
					},
					{
						Type: SegmentParam,
						Name: "slug",
					},
				},
			},
		},
		{
			name:    "parses multiple dynamic segments",
			pattern: "/users/:id/posts/:postId",
			want: ParsedPattern{
				Raw: "/users/:id/posts/:postId",
				Segments: []Segment{
					{
						Type:  SegmentStatic,
						Value: "users",
					},
					{
						Type: SegmentParam,
						Name: "id",
					},
					{
						Type:  SegmentStatic,
						Value: "posts",
					},
					{
						Type: SegmentParam,
						Name: "postId",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParsePattern(tt.pattern)

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf(
					"ParsePattern(%q)\n\n got: %#v\nwant: %#v",
					tt.pattern,
					got,
					tt.want,
				)
			}
		})
	}
}