from roution.compile import compile
from roution.strategy import select_strategy, default_options


def create_matcher(routes, options=None):
    """Create a reusable matcher from route definitions.

    Route definitions are compiled once. The matching strategy is selected
    automatically based on the compiled collection unless overridden through
    options (an Options instance or a plain dict).

    Raises TypeError if routes is not a non-null dict.
    """
    if routes is None or not isinstance(routes, dict) or isinstance(routes, list):
        raise TypeError("create_matcher expects a non-null dict of route definitions")

    if options is None:
        resolved = default_options()
    elif isinstance(options, dict):
        resolved = default_options()
        resolved.query = options.get("query", resolved.query)
        resolved.strategy = options.get("strategy", resolved.strategy)
        resolved.dynamic_threshold = options.get(
            "dynamic_threshold", resolved.dynamic_threshold
        )
    else:
        resolved = options

    compiled = compile(routes)
    return select_strategy(compiled, resolved)
