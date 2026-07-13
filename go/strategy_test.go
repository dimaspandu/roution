package roution

import "testing"

func routesWithDynamic(count int) map[string]any {
	routes := map[string]any{"*": "404"}
	for i := 0; i < count; i++ {
		routes["/item/"+itoa(i)+"/:slug"] = i
	}
	return routes
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	neg := i < 0
	if neg {
		i = -i
	}
	var buf [20]byte
	pos := len(buf)
	for i > 0 {
		pos--
		buf[pos] = byte('0' + i%10)
		i /= 10
	}
	if neg {
		pos--
		buf[pos] = '-'
	}
	return string(buf[pos:])
}

func TestSelectRegexBelowThreshold(t *testing.T) {
	compiled := Compile(routesWithDynamic(DynamicRouteThreshold - 1))
	m := SelectStrategy(compiled, Options{Strategy: StrategyAuto})
	if r := m.Match("/item/0/a"); r.Route != "/item/0/:slug" {
		t.Fatalf("route = %s, want /item/0/:slug", r.Route)
	}
}

func TestSelectTrieAtThreshold(t *testing.T) {
	compiled := Compile(routesWithDynamic(DynamicRouteThreshold))
	m := SelectStrategy(compiled, Options{Strategy: StrategyAuto})
	if r := m.Match("/item/0/a"); r.Route != "/item/0/:slug" {
		t.Fatalf("route = %s, want /item/0/:slug", r.Route)
	}
}

func TestSelectExplicitRegex(t *testing.T) {
	compiled := Compile(routesWithDynamic(DynamicRouteThreshold))
	m := SelectStrategy(compiled, Options{Strategy: StrategyRegex})
	if r := m.Match("/item/0/a"); r.Route != "/item/0/:slug" {
		t.Fatalf("route = %s, want /item/0/:slug", r.Route)
	}
}

func TestSelectExplicitTrie(t *testing.T) {
	compiled := Compile(routesWithDynamic(DynamicRouteThreshold - 1))
	m := SelectStrategy(compiled, Options{Strategy: StrategyTrie})
	if r := m.Match("/item/0/a"); r.Route != "/item/0/:slug" {
		t.Fatalf("route = %s, want /item/0/:slug", r.Route)
	}
}

func TestSelectCustomThreshold(t *testing.T) {
	compiled := Compile(routesWithDynamic(3))
	m := SelectStrategy(compiled, Options{Strategy: StrategyAuto, DynamicThreshold: 3})
	if r := m.Match("/item/0/a"); r.Route != "/item/0/:slug" {
		t.Fatalf("route = %s, want /item/0/:slug", r.Route)
	}
}
