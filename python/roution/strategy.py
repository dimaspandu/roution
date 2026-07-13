from dataclasses import dataclass

from roution.compile import CompiledRoutes
from roution.strategies.regex import create_regex_strategy
from roution.strategies.trie import create_trie_strategy

DYNAMIC_ROUTE_THRESHOLD = 50

STRATEGY_AUTO = "auto"
STRATEGY_REGEX = "regex"
STRATEGY_TRIE = "trie"


@dataclass
class Options:
    """Tunes matcher behavior.

    All fields are optional. Use default_options() to inherit the safe
    defaults (query enabled, automatic strategy selection, dynamic threshold
    of 50).
    """

    query: bool = True
    strategy: str = STRATEGY_AUTO
    dynamic_threshold: int = DYNAMIC_ROUTE_THRESHOLD


def default_options():
    return Options()


def select_strategy(compiled, options=None):
    """Select the most appropriate matching strategy for a compiled collection.

    When strategy is "auto" (the default), the decision is automatic and based
    on the number of dynamic routes. The explicit "regex" and "trie" values
    bypass the automatic decision. The matching behavior remains identical
    regardless of the chosen strategy.
    """
    if options is None:
        options = default_options()

    strategy = options.strategy
    if strategy == STRATEGY_TRIE:
        return create_trie_strategy(compiled, options.query)
    if strategy == STRATEGY_REGEX:
        return create_regex_strategy(compiled, options.query)
    if len(compiled.dynamic_routes) >= options.dynamic_threshold:
        return create_trie_strategy(compiled, options.query)
    return create_regex_strategy(compiled, options.query)
