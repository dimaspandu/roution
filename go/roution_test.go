package roution

import (
	"reflect"
	"testing"
)

func TestCreateMatcherQuickStart(t *testing.T) {
	routes := map[string]any{
		"/":               "public/index.html",
		"/articles":       "public/articles/index.html",
		"/articles/:slug": "public/articles/[slug].html",
		"/composed":       []string{"public/header.html", "public/greetings.html", "public/footer.html"},
		"*":               "public/404.html",
	}

	matcher := CreateMatcher(routes)
	result := matcher.Match("/articles/javascript?page=1")

	want := MatchResult{
		Found:    true,
		Pathname: "/articles/javascript",
		Route:    "/articles/:slug",
		Params:   map[string]string{"slug": "javascript"},
		Query:    map[string]any{"page": "1"},
		Value:    "public/articles/[slug].html",
	}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("got %+v, want %+v", result, want)
	}
}

func TestCreateMatcherReusable(t *testing.T) {
	matcher := CreateMatcher(map[string]any{"/a": 1, "/b": 2, "*": 0})
	if matcher.Match("/a").Value != 1 {
		t.Fatal("expected /a -> 1")
	}
	if matcher.Match("/b").Value != 2 {
		t.Fatal("expected /b -> 2")
	}
	if matcher.Match("/unknown").Value != 0 {
		t.Fatal("expected /unknown -> 0")
	}
}

func TestCreateMatcherDefaultOptions(t *testing.T) {
	matcher := CreateMatcher(map[string]any{"/articles/:slug": "v"})
	if r := matcher.Match("/articles/javascript?page=1"); r.Route != "/articles/:slug" {
		t.Fatalf("route = %s", r.Route)
	}
}

func TestCreateMatcherQueryDisabled(t *testing.T) {
	opts := DefaultOptions()
	opts.Query = false
	matcher := CreateMatcher(map[string]any{"/a": "v"}, opts)
	if r := matcher.Match("/a?x=1"); len(r.Query) != 0 {
		t.Fatalf("query = %v, want empty", r.Query)
	}
}

func TestCreateMatcherExplicitTrie(t *testing.T) {
	opts := DefaultOptions()
	opts.Strategy = StrategyTrie
	matcher := CreateMatcher(map[string]any{"/articles/:slug": "v"}, opts)
	if r := matcher.Match("/articles/javascript?page=1"); r.Route != "/articles/:slug" {
		t.Fatalf("route = %s", r.Route)
	}
}

func TestCreateMatcherNilRoutes(t *testing.T) {
	matcher := CreateMatcher(nil)
	if r := matcher.Match("/anything"); r.Found {
		t.Fatalf("expected no match for nil routes, got %+v", r)
	}
}
