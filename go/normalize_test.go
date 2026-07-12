package roution

import (
	"reflect"
	"testing"
)

func TestNormalize(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "strips query string",
			in:   "/articles?page=1",
			want: "/articles",
		},
		{
			name: "strips url fragment",
			in:   "/articles#comments",
			want: "/articles",
		},
		{
			name: "strips both query and fragment",
			in:   "/articles?page=1&perPage=10#comments",
			want: "/articles",
		},
		{
			name: "returns path untouched",
			in:   "/articles/javascript",
			want: "/articles/javascript",
		},
		{
			name: "empty string",
			in:   "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.in)

			if got != tt.want {
				t.Fatalf("Normalize(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestSplitSegments(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want []string
	}{
		{
			name: "ignores leading and trailing slashes",
			in:   "/a/b/c/",
			want: []string{"a", "b", "c"},
		},
		{
			name: "returns empty slice for root",
			in:   "/",
			want: []string{},
		},
		{
			name: "ignores duplicate slashes",
			in:   "//a///b//",
			want: []string{"a", "b"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SplitSegments(tt.in)

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("SplitSegments(%q) = %#v, want %#v", tt.in, got, tt.want)
			}
		})
	}
}

func TestParseQuery(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want map[string]any
	}{
		{
			name: "returns empty map without query",
			in:   "/articles",
			want: map[string]any{},
		},
		{
			name: "returns empty map for empty query",
			in:   "/articles?",
			want: map[string]any{},
		},
		{
			name: "single key",
			in:   "/articles?page=1",
			want: map[string]any{
				"page": "1",
			},
		},
		{
			name: "multiple keys",
			in:   "/articles?page=1&perPage=10",
			want: map[string]any{
				"page":    "1",
				"perPage": "10",
			},
		},
		{
			name: "repeated keys",
			in:   "/search?tag=a&tag=b",
			want: map[string]any{
				"tag": []string{"a", "b"},
			},
		},
		{
			name: "ignores fragment",
			in:   "/articles?page=1#top",
			want: map[string]any{
				"page": "1",
			},
		},
		{
			name: "decodes url encoding",
			in:   "/search?q=hello%20world",
			want: map[string]any{
				"q": "hello world",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseQuery(tt.in)

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("ParseQuery(%q) = %#v, want %#v", tt.in, got, tt.want)
			}
		})
	}
}