package roution

// StaticRoute is a compiled static route entry.
type StaticRoute struct {
	Route string
	Value any
}

// DynamicRoute is a compiled dynamic route entry.
type DynamicRoute struct {
	Route      string
	Value      any
	Segments   []Segment
	ParamNames []string
}

// CompiledRoutes is the internal, reusable representation of route definitions.
type CompiledRoutes struct {
	StaticRoutes map[string]StaticRoute
	DynamicRoutes []DynamicRoute
	Wildcard     *StaticRoute
}

// Compile partitions route definitions into static routes, dynamic routes, and
// an optional wildcard fallback. Compilation happens once; the result is reused
// for every match.
func Compile(routes map[string]any) CompiledRoutes {
	staticRoutes := make(map[string]StaticRoute)
	var dynamicRoutes []DynamicRoute
	var wildcard *StaticRoute

	for pattern, value := range routes {
		parsed := ParsePattern(pattern)

		if parsed.IsWildCard {
			w := StaticRoute{Route: pattern, Value: value}
			wildcard = &w
			continue
		}

		onlyStatic := true
		for _, seg := range parsed.Segments {
			if seg.Type != SegmentStatic {
				onlyStatic = false
				break
			}
		}

		if onlyStatic {
			staticRoutes[parsed.Raw] = StaticRoute{Route: pattern, Value: value}
			continue
		}

		var paramNames []string
		for _, seg := range parsed.Segments {
			if seg.Type == SegmentParam {
				paramNames = append(paramNames, seg.Name)
			}
		}

		dynamicRoutes = append(dynamicRoutes, DynamicRoute{
			Route:      pattern,
			Value:      value,
			Segments:   parsed.Segments,
			ParamNames: paramNames,
		})
	}

	// Go maps do not preserve insertion order, so sort dynamic routes by
	// pattern for deterministic matching priority (first match wins).
	sortDynamicRoutes(dynamicRoutes)

	return CompiledRoutes{
		StaticRoutes: staticRoutes,
		DynamicRoutes: dynamicRoutes,
		Wildcard:     wildcard,
	}
}

func sortDynamicRoutes(routes []DynamicRoute) {
	// Simple insertion sort keeps the dependency small (no sort import needed
	// for typical route counts) while guaranteeing a stable, reproducible order.
	for i := 1; i < len(routes); i++ {
		for j := i; j > 0 && routes[j].Route < routes[j-1].Route; j-- {
			routes[j], routes[j-1] = routes[j-1], routes[j]
		}
	}
}
