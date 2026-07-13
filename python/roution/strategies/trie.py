from urllib.parse import unquote

from roution.compile import StaticRoute
from roution.normalize import normalize_pathname, split_segments, parse_query
from roution.pattern import parse_pattern
from roution.result import make_result


class _TrieNode:
    __slots__ = ("static_children", "param_child", "terminal", "wildcard")

    def __init__(self):
        self.static_children = {}
        self.param_child = None  # (name, _TrieNode)
        self.terminal = None  # StaticRoute
        self.wildcard = None  # StaticRoute


def _insert(root, segments, entry):
    node = root
    for segment in segments:
        if segment.type == "param":
            if node.param_child is None:
                node.param_child = (segment.name, _TrieNode())
            node = node.param_child[1]
            continue
        child = node.static_children.get(segment.value)
        if child is None:
            child = _TrieNode()
            node.static_children[segment.value] = child
        node = child
    node.terminal = entry


def _build_trie(compiled):
    root = _TrieNode()
    if compiled.wildcard is not None:
        root.wildcard = compiled.wildcard
    for raw, static_route in compiled.static_routes.items():
        _insert(root, parse_pattern(raw).segments, static_route)
    for dynamic_route in compiled.dynamic_routes:
        _insert(
            root,
            dynamic_route.segments,
            StaticRoute(route=dynamic_route.route, value=dynamic_route.value),
        )
    return root


def _match_trie(root, segments):
    node = root
    params = {}
    for segment in segments:
        child = node.static_children.get(segment)
        if child is not None:
            node = child
            continue
        if node.param_child is not None:
            name, child = node.param_child
            params[name] = unquote(segment)
            node = child
            continue
        return None, None
    if node.terminal is not None:
        return node.terminal, params
    return None, None


class TrieStrategy:
    """Matcher strategy backed by a path trie.

    Static and dynamic routes live in the same trie; an optional wildcard is
    the fallback. Suitable for collections with a large number of routes.
    """

    def __init__(self, compiled, include_query):
        self.compiled = compiled
        self.include_query = include_query
        self.root = _build_trie(compiled)

    def match(self, pathname):
        normalized = normalize_pathname(pathname)
        query = parse_query(pathname) if self.include_query else {}
        segments = split_segments(normalized)

        terminal, params = _match_trie(self.root, segments)
        if terminal is not None:
            return make_result(
                True, normalized, terminal.route, params, query, terminal.value
            )

        if self.compiled.wildcard is not None:
            wildcard = self.compiled.wildcard
            return make_result(
                True, normalized, wildcard.route, {}, query, wildcard.value
            )

        return make_result(False, normalized, None, {}, query, None)


def create_trie_strategy(compiled, include_query):
    return TrieStrategy(compiled, include_query)
