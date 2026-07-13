import re
from urllib.parse import unquote

from roution.normalize import normalize_pathname, split_segments, parse_query
from roution.result import make_result


def _build_regex(dynamic_route):
    parts = []
    for index, segment in enumerate(dynamic_route.segments):
        if index > 0:
            parts.append("/")
        if segment.type == "static":
            parts.append(re.escape(segment.value))
        else:
            parts.append("([^/]+)")
    return re.compile("^/" + "".join(parts) + "$")


class RegexStrategy:
    """Matcher strategy backed by precompiled regular expressions.

    Static routes use an exact dict lookup, dynamic routes use precompiled
    regular expressions, and an optional wildcard is used as a fallback.
    """

    def __init__(self, compiled, include_query):
        self.compiled = compiled
        self.include_query = include_query
        self.regex_routes = [
            (_build_regex(route), route.param_names, route.route, route.value)
            for route in compiled.dynamic_routes
        ]

    def match(self, pathname):
        normalized = normalize_pathname(pathname)
        query = parse_query(pathname) if self.include_query else {}

        static_hit = self.compiled.static_routes.get(normalized)
        if static_hit is not None:
            return make_result(
                True, normalized, static_hit.route, {}, query, static_hit.value
            )

        segments = split_segments(normalized)
        reconstructed = "/" + "/".join(segments)
        for regex, param_names, route, value in self.regex_routes:
            match = regex.match(reconstructed)
            if match is None:
                continue
            params = {
                name: unquote(match.group(index + 1))
                for index, name in enumerate(param_names)
            }
            return make_result(True, normalized, route, params, query, value)

        if self.compiled.wildcard is not None:
            wildcard = self.compiled.wildcard
            return make_result(
                True, normalized, wildcard.route, {}, query, wildcard.value
            )

        return make_result(False, normalized, None, {}, query, None)


def create_regex_strategy(compiled, include_query):
    return RegexStrategy(compiled, include_query)
