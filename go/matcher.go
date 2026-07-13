package roution

// MatchResult is the structured result returned for an incoming pathname.
type MatchResult struct {
	Found    bool
	Pathname string
	Route    string
	Params   map[string]string
	Query    map[string]any
	Value    any
}

// Matcher resolves a pathname to a MatchResult.
type Matcher interface {
	Match(pathname string) MatchResult
}

// Strategy names.
const (
	StrategyAuto  = "auto"
	StrategyRegex = "regex"
	StrategyTrie  = "trie"
)

// Options tunes matcher behavior.
//
// All fields are optional. Use DefaultOptions to inherit the safe defaults
// (query enabled, automatic strategy selection, dynamic threshold of 50).
type Options struct {
	// Query enables parsing of the query string into the result.
	Query bool
	// Strategy selects the matching algorithm: "auto", "regex", or "trie".
	Strategy string
	// DynamicThreshold is the dynamic route count at/above which "auto"
	// selects the trie strategy.
	DynamicThreshold int
}

// DefaultOptions returns the safe default options.
func DefaultOptions() Options {
	return Options{
		Query:            true,
		Strategy:         StrategyAuto,
		DynamicThreshold: DynamicRouteThreshold,
	}
}

// CreateMatcher compiles route definitions into a reusable matcher.
//
// Route values are treated as opaque data. Compilation happens once. The
// matching strategy is selected automatically based on the compiled
// collection unless overridden through opts.
func CreateMatcher(routes map[string]any, opts ...Options) Matcher {
	o := DefaultOptions()
	if len(opts) > 0 {
		o = opts[0]
		if o.Strategy == "" {
			o.Strategy = StrategyAuto
		}
		if o.DynamicThreshold == 0 {
			o.DynamicThreshold = DynamicRouteThreshold
		}
	}

	compiled := Compile(routes)
	return SelectStrategy(compiled, o)
}

// CreateMatcherStrict is like CreateMatcher but requires a non-nil routes map.
//
// It is retained for callers that want an explicit precondition; the regular
// CreateMatcher already tolerates a nil map (yielding no routes).
func CreateMatcherStrict(routes map[string]any, opts ...Options) Matcher {
	if routes == nil {
		panic("CreateMatcherStrict expects a non-nil routes map")
	}
	return CreateMatcher(routes, opts...)
}
