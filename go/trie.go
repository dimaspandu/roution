package roution

// TrieNode is a node in the path trie.
type TrieNode struct {
	StaticChildren map[string]*TrieNode
	ParamChild     *TrieParamChild
	Terminal       *StaticRoute
	Wildcard       *StaticRoute
}

// TrieParamChild is a single parameter child at one trie level.
type TrieParamChild struct {
	Name string
	Node *TrieNode
}

// TrieMatcher matches pathnames using a path trie.
//
// Suitable for collections with a large number of routes, where a single trie
// traversal scales better than scanning regular expressions. Static and
// dynamic routes live in the same trie; an optional wildcard is the fallback.
type TrieMatcher struct {
	root         *TrieNode
	includeQuery bool
}

// NewTrieMatcher builds a trie-backed matcher.
func NewTrieMatcher(compiled CompiledRoutes, includeQuery bool) *TrieMatcher {
	root := &TrieNode{
		StaticChildren: make(map[string]*TrieNode),
	}
	if compiled.Wildcard != nil {
		wc := *compiled.Wildcard
		root.Wildcard = &wc
	}
	for raw, sr := range compiled.StaticRoutes {
		insertTrie(root, ParsePattern(raw).Segments, sr)
	}
	for _, dr := range compiled.DynamicRoutes {
		insertTrie(root, dr.Segments, StaticRoute{Route: dr.Route, Value: dr.Value})
	}
	return &TrieMatcher{root: root, includeQuery: includeQuery}
}

func insertTrie(root *TrieNode, segments []Segment, entry StaticRoute) {
	node := root
	for _, seg := range segments {
		if seg.Type == SegmentParam {
			if node.ParamChild == nil {
				node.ParamChild = &TrieParamChild{
					Name: seg.Name,
					Node: &TrieNode{StaticChildren: make(map[string]*TrieNode)},
				}
			}
			node = node.ParamChild.Node
			continue
		}
		child, ok := node.StaticChildren[seg.Value]
		if !ok {
			child = &TrieNode{StaticChildren: make(map[string]*TrieNode)}
			node.StaticChildren[seg.Value] = child
		}
		node = child
	}
	terminal := entry
	node.Terminal = &terminal
}

func matchTrie(node *TrieNode, segments []string) (*StaticRoute, map[string]string) {
	params := make(map[string]string)
	for _, seg := range segments {
		if child, ok := node.StaticChildren[seg]; ok {
			node = child
			continue
		}
		if node.ParamChild != nil {
			params[node.ParamChild.Name] = decode(seg)
			node = node.ParamChild.Node
			continue
		}
		return nil, nil
	}
	if node.Terminal != nil {
		return node.Terminal, params
	}
	return nil, nil
}

// Match resolves an incoming pathname.
func (m *TrieMatcher) Match(pathname string) MatchResult {
	normalized := Normalize(pathname)
	query := m.query(pathname)

	segments := SplitSegments(normalized)
	terminal, params := matchTrie(m.root, segments)
	if terminal != nil {
		return MatchResult{
			Found:    true,
			Pathname: normalized,
			Route:    terminal.Route,
			Params:   params,
			Query:    query,
			Value:    terminal.Value,
		}
	}

	if m.root.Wildcard != nil {
		return MatchResult{
			Found:    true,
			Pathname: normalized,
			Route:    m.root.Wildcard.Route,
			Params:   map[string]string{},
			Query:    query,
			Value:    m.root.Wildcard.Value,
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

func (m *TrieMatcher) query(pathname string) map[string]any {
	if !m.includeQuery {
		return map[string]any{}
	}
	return ParseQuery(pathname)
}
