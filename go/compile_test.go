package roution

import (
	"reflect"
	"testing"
)

func TestCompileClassifiesRoot(t *testing.T) {
	compiled := Compile(map[string]any{"/": "home"})
	if got := compiled.StaticRoutes["/"].Value; got != "home" {
		t.Fatalf("static root value = %v, want home", got)
	}
	if len(compiled.DynamicRoutes) != 0 {
		t.Fatalf("dynamic routes = %d, want 0", len(compiled.DynamicRoutes))
	}
}

func TestCompileClassifiesStatic(t *testing.T) {
	compiled := Compile(map[string]any{"/articles": "a", "/about": "b"})
	if len(compiled.StaticRoutes) != 2 {
		t.Fatalf("static routes = %d, want 2", len(compiled.StaticRoutes))
	}
	if got := compiled.StaticRoutes["/about"].Value; got != "b" {
		t.Fatalf("about value = %v, want b", got)
	}
}

func TestCompileClassifiesDynamic(t *testing.T) {
	compiled := Compile(map[string]any{"/articles/:slug": "a"})
	if len(compiled.StaticRoutes) != 0 {
		t.Fatalf("static routes = %d, want 0", len(compiled.StaticRoutes))
	}
	if len(compiled.DynamicRoutes) != 1 {
		t.Fatalf("dynamic routes = %d, want 1", len(compiled.DynamicRoutes))
	}
	if got := compiled.DynamicRoutes[0].ParamNames; !reflect.DeepEqual(got, []string{"slug"}) {
		t.Fatalf("param names = %v, want [slug]", got)
	}
}

func TestCompileClassifiesWildcard(t *testing.T) {
	compiled := Compile(map[string]any{"*": "fallback"})
	if compiled.Wildcard == nil || compiled.Wildcard.Value != "fallback" {
		t.Fatalf("wildcard = %v, want fallback", compiled.Wildcard)
	}
}

func TestCompileMixedCollection(t *testing.T) {
	compiled := Compile(map[string]any{
		"/":            "home",
		"/articles/:slug": "dynamic",
		"*":            "fallback",
	})
	if len(compiled.StaticRoutes) != 1 {
		t.Fatalf("static routes = %d, want 1", len(compiled.StaticRoutes))
	}
	if len(compiled.DynamicRoutes) != 1 {
		t.Fatalf("dynamic routes = %d, want 1", len(compiled.DynamicRoutes))
	}
	if compiled.Wildcard == nil || compiled.Wildcard.Value != "fallback" {
		t.Fatalf("wildcard = %v, want fallback", compiled.Wildcard)
	}
}
