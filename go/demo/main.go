package main

import (
	"encoding/json"
	"fmt"
	"maps"
	"reflect"
	"sort"
	"strings"
	"time"

	roution "github.com/dimaspandu/roution"
)

func main() {
	// Base routes shared by both scenarios. No wildcard is defined here so the
	// difference between "without wildcard" and "with wildcard" is easy to observe.
	baseRoutes := map[string]any{
		"/":                        "public/index.html",
		"/about":                   "public/about.html",
		"/articles":                "public/articles/index.html",
		"/articles/:slug":          "public/articles/[slug].html",
		"/users/:id/posts/:postId": "public/users/[id]/posts/[postId].html",
		"/composed":                []string{"public/header.html", "public/greetings.html", "public/footer.html"},
		"/login":                   map[string]any{"component": "LoginPage", "auth": false},
		"/dashboard":               map[string]any{"component": "DashboardPage", "auth": true, "layout": "main"},
		"/version":                 "1.0.0",
		"/numbers":                 []int{1, 2, 3},
		"/health":                  func() map[string]string { return map[string]string{"status": "ok"} },
		"/time":                    func() string { return time.Now().Format(time.RFC3339) },
	}

	// Scenario A: no wildcard. Any pathname that does not match an explicit route
	// resolves to found: false (route: null, value: null).
	withoutWildcard := maps.Clone(baseRoutes)

	// Scenario B: with a wildcard. Unmatched pathnames still resolve to
	// found: true, but the value comes from the "*" fallback (public/404.html).
	withWildcard := maps.Clone(baseRoutes)
	withWildcard["*"] = "public/404.html"

	// CreateMatcher accepts optional Options to tune behavior. All fields are
	// optional and fall back to safe defaults:
	//
	//   roution.CreateMatcher(routes, roution.Options{Query: false})
	//   roution.CreateMatcher(routes, roution.Options{Strategy: roution.StrategyRegex})
	//   roution.CreateMatcher(routes, roution.Options{Strategy: roution.StrategyTrie})
	//   roution.CreateMatcher(routes, roution.Options{Strategy: roution.StrategyAuto, DynamicThreshold: 20})

	samples := []string{
		"/",
		"/about",
		"/articles",
		"/articles/javascript?page=1",
		"/users/42/posts/7#comments",
		"/composed",
		"/login",
		"/dashboard",
		"/version",
		"/numbers",
		"/health",
		"/time",
		// The following paths are intentionally unregistered to show the difference
		// between resolving without and with a wildcard fallback.
		"/articles/javascript/extra",
		"/missing/page",
		"/wild/card/path",
		"/not-a-real-page",
	}

	runScenario("WITHOUT wildcard (unmatched -> found: false)", withoutWildcard, samples)
	runScenario("WITH wildcard \"*\" (unmatched -> found: true, value: public/404.html)", withWildcard, samples)
}

func runScenario(title string, routes map[string]any, samples []string) {
	matcher := roution.CreateMatcher(routes)

	fmt.Println(strings.Repeat("=", 70))
	fmt.Println(title)
	fmt.Println(strings.Repeat("=", 70))
	fmt.Println("Route collection:")
	fmt.Println(formatKeys(routes))
	fmt.Println()
	fmt.Println("Match samples:")
	fmt.Println()

	for _, pathname := range samples {
		result := matcher.Match(pathname)
		fmt.Printf("pathname: %s\n", pathname)
		fmt.Printf("  found:    %v\n", result.Found)
		fmt.Printf("  route:    %s\n", formatRoute(result.Route))
		fmt.Printf("  pathname: %s\n", result.Pathname)
		fmt.Printf("  params:   %s\n", toJSON(result.Params))
		fmt.Printf("  query:    %s\n", toJSON(result.Query))
		fmt.Printf("  value:    %s\n", formatValue(result.Value))
		fmt.Println()
	}
}

// formatKeys returns the route patterns as a pretty-printed JSON array with a
// stable (sorted) order, since Go map iteration order is not deterministic.
func formatKeys(routes map[string]any) string {
	keys := make([]string, 0, len(routes))
	for k := range routes {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	b, _ := json.MarshalIndent(keys, "", "  ")
	return string(b)
}

// formatRoute renders the matched route, using "null" when nothing matched.
func formatRoute(route string) string {
	if route == "" {
		return "null"
	}
	return route
}

// toJSON renders params and query as compact JSON so the output matches the
// JavaScript demo ("{}" and {"slug":"javascript"}) instead of Go's map syntax.
func toJSON(value any) string {
	b, err := json.Marshal(value)
	if err != nil {
		return fmt.Sprintf("%v", value)
	}
	return string(b)
}

// formatValue renders an opaque route value for display, mirroring the
// JavaScript demo: functions become "[Function]", composite values become JSON,
// and everything else uses the default formatting.
func formatValue(value any) string {
	if value == nil {
		return "null"
	}
	switch reflect.ValueOf(value).Kind() {
	case reflect.Func:
		return "[Function]"
	case reflect.Slice, reflect.Array, reflect.Map, reflect.Struct:
		if b, err := json.Marshal(value); err == nil {
			return string(b)
		}
	}
	return fmt.Sprintf("%v", value)
}
