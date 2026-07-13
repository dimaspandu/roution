package roution

import (
	"net/url"
	"regexp"
	"strings"
)

// RegexMatcher matches pathnames using precompiled regular expressions.
//
// Suitable for collections with a moderate number of dynamic routes. Static
// routes use an exact map lookup, dynamic routes use precompiled regular
// expressions, and an optional wildcard is used as a fallback.
type RegexMatcher struct {
	compiled     CompiledRoutes
	regexRoutes  []regexRoute
	includeQuery bool
}

type regexRoute struct {
	re         *regexp.Regexp
	paramNames []string
	route      string
	value      any
}

// NewRegexMatcher builds a regex-backed matcher.
func NewRegexMatcher(compiled CompiledRoutes, includeQuery bool) *RegexMatcher {
	regexRoutes := make([]regexRoute, 0, len(compiled.DynamicRoutes))
	for _, dr := range compiled.DynamicRoutes {
		regexRoutes = append(regexRoutes, regexRoute{
			re:         buildRegex(dr),
			paramNames: dr.ParamNames,
			route:      dr.Route,
			value:      dr.Value,
		})
	}
	return &RegexMatcher{
		compiled:     compiled,
		regexRoutes:  regexRoutes,
		includeQuery: includeQuery,
	}
}

func buildRegex(dr DynamicRoute) *regexp.Regexp {
	var b strings.Builder
	b.WriteString("^/")
	for i, seg := range dr.Segments {
		if i > 0 {
			b.WriteString("/")
		}
		if seg.Type == SegmentStatic {
			b.WriteString(regexp.QuoteMeta(seg.Value))
			continue
		}
		b.WriteString("([^/]+)")
	}
	b.WriteString("$")
	return regexp.MustCompile(b.String())
}

// Match resolves an incoming pathname.
func (m *RegexMatcher) Match(pathname string) MatchResult {
	normalized := Normalize(pathname)
	query := m.query(pathname)

	if sr, ok := m.compiled.StaticRoutes[normalized]; ok {
		return MatchResult{
			Found:    true,
			Pathname: normalized,
			Route:    sr.Route,
			Params:   map[string]string{},
			Query:    query,
			Value:    sr.Value,
		}
	}

	segments := SplitSegments(normalized)
	reconstructed := "/" + strings.Join(segments, "/")

	for _, rr := range m.regexRoutes {
		matches := rr.re.FindStringSubmatch(reconstructed)
		if matches == nil {
			continue
		}
		params := make(map[string]string, len(rr.paramNames))
		for i, name := range rr.paramNames {
			params[name] = decode(matches[i+1])
		}
		return MatchResult{
			Found:    true,
			Pathname: normalized,
			Route:    rr.route,
			Params:   params,
			Query:    query,
			Value:    rr.value,
		}
	}

	if m.compiled.Wildcard != nil {
		return MatchResult{
			Found:    true,
			Pathname: normalized,
			Route:    m.compiled.Wildcard.Route,
			Params:   map[string]string{},
			Query:    query,
			Value:    m.compiled.Wildcard.Value,
		}
	}

	return MatchResult{
		Found:    false,
		Pathname: normalized,
		Route:    "",
		Params:   map[string]string{},
		Query:    query,
		Value:    nil,
	}
}

func (m *RegexMatcher) query(pathname string) map[string]any {
	if !m.includeQuery {
		return map[string]any{}
	}
	return ParseQuery(pathname)
}

// decode applies URL decoding to a captured segment, falling back to the raw
// value when decoding fails.
func decode(value string) string {
	decoded, err := url.QueryUnescape(value)
	if err != nil {
		return value
	}
	return decoded
}
