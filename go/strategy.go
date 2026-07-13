package roution

// DynamicRouteThreshold is the number of dynamic routes at or above which the
// "auto" strategy selects the trie implementation. Considered an
// implementation detail.
const DynamicRouteThreshold = 50

// SelectStrategy chooses the most appropriate matcher implementation for a
// compiled route collection.
//
// When Strategy is "auto" (the default), the decision is automatic and based on
// the number of dynamic routes: a small collection uses the regex strategy and
// a large collection uses the trie strategy. The explicit "regex" and "trie"
// values bypass the automatic decision. The public API and matching behavior
// remain identical regardless of the chosen strategy.
func SelectStrategy(compiled CompiledRoutes, opts Options) Matcher {
	switch opts.Strategy {
	case StrategyTrie:
		return NewTrieMatcher(compiled, opts.Query)
	case StrategyRegex:
		return NewRegexMatcher(compiled, opts.Query)
	default:
		if len(compiled.DynamicRoutes) >= opts.DynamicThreshold {
			return NewTrieMatcher(compiled, opts.Query)
		}
		return NewRegexMatcher(compiled, opts.Query)
	}
}
