from __future__ import annotations

import json
from datetime import datetime, timezone

from roution import create_matcher


# Base routes shared by both scenarios. No wildcard is defined here so the
# difference between "without wildcard" and "with wildcard" is easy to observe.
BASE_ROUTES = {
    "/": "public/index.html",
    "/about": "public/about.html",
    "/articles": "public/articles/index.html",
    "/articles/:slug": "public/articles/[slug].html",
    "/users/:id/posts/:postId": "public/users/[id]/posts/[postId].html",
    "/composed": ["public/header.html", "public/greetings.html", "public/footer.html"],
    "/login": {"component": "LoginPage", "auth": False},
    "/dashboard": {"component": "DashboardPage", "auth": True, "layout": "main"},
    "/version": "1.0.0",
    "/numbers": [1, 2, 3],
    "/health": lambda: {"status": "ok"},
    "/time": lambda: datetime.now(timezone.utc).isoformat(),
}

SAMPLES = [
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
    "/time",
    # The following paths are intentionally unregistered to show the difference
    # between resolving without and with a wildcard fallback.
    "/articles/javascript/extra",
    "/missing/page",
    "/wild/card/path",
    "/not-a-real-page",
]


def main():
    # Scenario A: no wildcard. Any pathname that does not match an explicit route
    # resolves to found: false (route: null, value: null).
    without_wildcard = {**BASE_ROUTES}

    # Scenario B: with a wildcard. Unmatched pathnames still resolve to
    # found: true, but the value comes from the "*" fallback (public/404.html).
    with_wildcard = {**BASE_ROUTES, "*": "public/404.html"}

    # create_matcher accepts an optional Options (or dict) to tune behavior.
    # All fields are optional and fall back to safe defaults:
    #
    #   create_matcher(routes, {"query": False})
    #   create_matcher(routes, {"strategy": "regex"})
    #   create_matcher(routes, {"strategy": "trie"})
    #   create_matcher(routes, {"strategy": "auto", "dynamic_threshold": 20})

    run_scenario("WITHOUT wildcard (unmatched -> found: false)", without_wildcard)
    run_scenario(
        'WITH wildcard "*" (unmatched -> found: true, value: public/404.html)',
        with_wildcard,
    )


def run_scenario(title, routes):
    matcher = create_matcher(routes)

    print("=" * 70)
    print(title)
    print("=" * 70)
    print("Route collection:")
    print(json.dumps(list(routes.keys()), indent=2))
    print()
    print("Match samples:")
    print()

    for pathname in SAMPLES:
        result = matcher.match(pathname)
        print(f"pathname: {pathname}")
        print(f"  found:    {to_json(result['found'])}")
        print(f"  route:    {format_route(result['route'])}")
        print(f"  pathname: {result['pathname']}")
        print(f"  params:   {to_json(result['params'])}")
        print(f"  query:    {to_json(result['query'])}")
        print(f"  value:    {format_value(result['value'])}")
        print()


def to_json(value):
    """Render values as compact JSON so the output matches the JS/Go demos."""
    return json.dumps(value, separators=(",", ":"))


def format_route(route):
    """Render the matched route, using "null" when nothing matched."""
    return "null" if route is None else route


def format_value(value):
    """Render an opaque route value, mirroring the JavaScript/Go demos.

    Callables become "[Function]", composite values become JSON, and everything
    else uses the default string formatting.
    """
    if value is None:
        return "null"
    if callable(value):
        return "[Function]"
    if isinstance(value, (list, dict)):
        return json.dumps(value, separators=(",", ":"))
    return str(value)


if __name__ == "__main__":
    main()
