package roution

import (
	"reflect"
	"testing"
)

// buildStrategies returns both matching strategies for the same routes so
// every assertion can be verified against the regex and trie implementations.
func buildStrategies(t *testing.T, routes map[string]any) []Matcher {
	t.Helper()
	compiled := Compile(routes)
	return []Matcher{
		NewRegexMatcher(compiled, true),
		NewTrieMatcher(compiled, true),
	}
}

func TestStrategiesMatchStaticRoot(t *testing.T) {
	for _, m := range buildStrategies(t, map[string]any{"/": "home"}) {
		r := m.Match("/")
		if !r.Found || r.Value != "home" {
			t.Fatalf("%T: %+v", m, r)
		}
		if len(r.Params) != 0 {
			t.Fatalf("%T: params = %v, want empty", m, r.Params)
		}
	}
}

func TestStrategiesMatchStaticRoute(t *testing.T) {
	for _, m := range buildStrategies(t, map[string]any{"/about": "about"}) {
		r := m.Match("/about")
		if !r.Found || r.Route != "/about" {
			t.Fatalf("%T: %+v", m, r)
		}
	}
}

func TestStrategiesMatchDynamicAndExtractParams(t *testing.T) {
	for _, m := range buildStrategies(t, map[string]any{"/articles/:slug": "page"}) {
		r := m.Match("/articles/javascript")
		if !r.Found || r.Route != "/articles/:slug" {
			t.Fatalf("%T: %+v", m, r)
		}
		if !reflect.DeepEqual(r.Params, map[string]string{"slug": "javascript"}) {
			t.Fatalf("%T: params = %v", m, r.Params)
		}
	}
}

func TestStrategiesNotFound(t *testing.T) {
	for _, m := range buildStrategies(t, map[string]any{"/about": "about"}) {
		r := m.Match("/missing")
		if r.Found || r.Route != "" || r.Value != nil {
			t.Fatalf("%T: %+v", m, r)
		}
	}
}

func TestStrategiesWildcardFallback(t *testing.T) {
	for _, m := range buildStrategies(t, map[string]any{"*": "404"}) {
		r := m.Match("/articles/javascript")
		if !r.Found || r.Route != "*" || r.Value != "404" {
			t.Fatalf("%T: %+v", m, r)
		}
	}
}

func TestStrategiesStaticPriorityOverWildcard(t *testing.T) {
	for _, m := range buildStrategies(t, map[string]any{"/about": "about", "*": "404"}) {
		if r := m.Match("/about"); r.Value != "about" {
			t.Fatalf("%T: %+v", m, r)
		}
	}
}

func TestStrategiesDynamicPriorityOverWildcard(t *testing.T) {
	for _, m := range buildStrategies(t, map[string]any{"/articles/:slug": "page", "*": "404"}) {
		if r := m.Match("/articles/javascript"); r.Value != "page" {
			t.Fatalf("%T: %+v", m, r)
		}
	}
}

func TestStrategiesNormalizeQuery(t *testing.T) {
	for _, m := range buildStrategies(t, map[string]any{"/articles/:slug": "page"}) {
		r := m.Match("/articles/javascript?page=1")
		if !r.Found || r.Pathname != "/articles/javascript" {
			t.Fatalf("%T: %+v", m, r)
		}
	}
}

func TestStrategiesNormalizeFragment(t *testing.T) {
	for _, m := range buildStrategies(t, map[string]any{"/articles": "page"}) {
		r := m.Match("/articles#section")
		if !r.Found || r.Pathname != "/articles" {
			t.Fatalf("%T: %+v", m, r)
		}
	}
}

func TestStrategiesDecodeParams(t *testing.T) {
	for _, m := range buildStrategies(t, map[string]any{"/files/:name": "file"}) {
		r := m.Match("/files/hello%20world")
		if !reflect.DeepEqual(r.Params, map[string]string{"name": "hello world"}) {
			t.Fatalf("%T: params = %v", m, r.Params)
		}
	}
}

func TestStrategiesSegmentCountMustMatch(t *testing.T) {
	for _, m := range buildStrategies(t, map[string]any{"/users/:id": "u"}) {
		if r := m.Match("/users/1/posts"); r.Found {
			t.Fatalf("%T: expected no match, got %+v", m, r)
		}
	}
}

func TestStrategiesStaticPriorityOverDynamic(t *testing.T) {
	for _, m := range buildStrategies(t, map[string]any{
		"/users/admin": "admin",
		"/users/:id":   "user",
	}) {
		r := m.Match("/users/admin")
		if r.Route != "/users/admin" {
			t.Fatalf("%T: route = %s, want /users/admin", m, r.Route)
		}
		if len(r.Params) != 0 {
			t.Fatalf("%T: params = %v, want empty", m, r.Params)
		}
	}
}

func TestStrategiesIdenticalResults(t *testing.T) {
	routes := map[string]any{
		"/":                  "home",
		"/about":             "about",
		"/articles":          "list",
		"/articles/:slug":    "article",
		"/users/:id/posts/:postId": "post",
		"/composed":          []string{"a", "b"},
		"/login":             map[string]any{"auth": false},
		"/version":           "1.0.0",
		"/numbers":           []int{1, 2, 3},
		"*":                  "404",
	}
	compiled := Compile(routes)
	regex := NewRegexMatcher(compiled, true)
	trie := NewTrieMatcher(compiled, true)

	samples := []string{
		"/",
		"/about",
		"/articles",
		"/articles/javascript?page=1&perPage=10",
		"/users/42/posts/7#comments",
		"/composed",
		"/login",
		"/version",
		"/numbers",
		"/articles/javascript/extra",
		"/missing/page",
	}

	for _, p := range samples {
		a := regex.Match(p)
		b := trie.Match(p)
		if !reflect.DeepEqual(a, b) {
			t.Fatalf("mismatch for %q:\n regex=%+v\n trie =%+v", p, a, b)
		}
	}
}

func TestStrategiesQueryDisabled(t *testing.T) {
	compiled := Compile(map[string]any{"/a": "v"})
	regex := NewRegexMatcher(compiled, false)
	trie := NewTrieMatcher(compiled, false)
	for _, m := range []Matcher{regex, trie} {
		if r := m.Match("/a?x=1"); len(r.Query) != 0 {
			t.Fatalf("%T: query = %v, want empty", m, r.Query)
		}
	}
}
