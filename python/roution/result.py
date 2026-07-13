def make_result(found, pathname, route, params, query, value):
    """Build a structured matching result (mirrors the JavaScript shape)."""
    return {
        "found": found,
        "pathname": pathname,
        "route": route,
        "params": params,
        "query": query,
        "value": value,
    }
