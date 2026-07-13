package main

import (
	"fmt"

	roution "github.com/dimaspandu/roution"
)

func main() {
	routes := map[string]any{
		"/":                  "public/index.html",
		"/about":             "public/about.html",
		"/articles":          "public/articles/index.html",
		"/articles/:slug":    "public/articles/[slug].html",
		"/users/:id/posts/:postId": "public/users/[id]/posts/[postId].html",
		"/composed":          []string{"public/header.html", "public/greetings.html", "public/footer.html"},
		"/login":             map[string]any{"component": "LoginPage", "auth": false},
		"/dashboard":         map[string]any{"component": "DashboardPage", "auth": true, "layout": "main"},
		"/version":           "1.0.0",
		"/numbers":           []int{1, 2, 3},
		"/health":            func() map[string]string { return map[string]string{"status": "ok"} },
		"*":                  "public/404.html",
	}

	// Default options: query enabled, automatic strategy selection.
	matcher := roution.CreateMatcher(routes)

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
		"/articles/javascript/extra",
		"/missing/page",
		"/wild/card/path",
	}

	for _, pathname := range samples {
		result := matcher.Match(pathname)
		fmt.Printf("pathname: %s\n", pathname)
		fmt.Printf("  found:    %v\n", result.Found)
		fmt.Printf("  route:    %s\n", result.Route)
		fmt.Printf("  pathname: %s\n", result.Pathname)
		fmt.Printf("  params:   %v\n", result.Params)
		fmt.Printf("  query:    %v\n", result.Query)
		fmt.Printf("  value:    %v\n", result.Value)
		fmt.Println()
	}
}
