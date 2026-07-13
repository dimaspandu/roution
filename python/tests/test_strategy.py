import unittest

from roution.compile import compile
from roution.strategy import (
    DYNAMIC_ROUTE_THRESHOLD,
    STRATEGY_AUTO,
    STRATEGY_REGEX,
    STRATEGY_TRIE,
    Options,
    select_strategy,
)


def routes_with_dynamic(count):
    routes = {"*": "404"}
    for i in range(count):
        routes[f"/item/{i}/:slug"] = i
    return routes


class StrategyTest(unittest.TestCase):
    def test_regex_below_threshold(self):
        compiled = compile(routes_with_dynamic(DYNAMIC_ROUTE_THRESHOLD - 1))
        strategy = select_strategy(compiled, Options(strategy=STRATEGY_AUTO))
        self.assertEqual(strategy.match("/item/0/a")["route"], "/item/0/:slug")

    def test_trie_at_threshold(self):
        compiled = compile(routes_with_dynamic(DYNAMIC_ROUTE_THRESHOLD))
        strategy = select_strategy(compiled, Options(strategy=STRATEGY_AUTO))
        self.assertEqual(strategy.match("/item/0/a")["route"], "/item/0/:slug")

    def test_explicit_regex(self):
        compiled = compile(routes_with_dynamic(DYNAMIC_ROUTE_THRESHOLD))
        strategy = select_strategy(compiled, Options(strategy=STRATEGY_REGEX))
        self.assertEqual(strategy.match("/item/0/a")["route"], "/item/0/:slug")

    def test_explicit_trie(self):
        compiled = compile(routes_with_dynamic(DYNAMIC_ROUTE_THRESHOLD - 1))
        strategy = select_strategy(compiled, Options(strategy=STRATEGY_TRIE))
        self.assertEqual(strategy.match("/item/0/a")["route"], "/item/0/:slug")

    def test_custom_threshold(self):
        compiled = compile(routes_with_dynamic(3))
        strategy = select_strategy(
            compiled, Options(strategy=STRATEGY_AUTO, dynamic_threshold=3)
        )
        self.assertEqual(strategy.match("/item/0/a")["route"], "/item/0/:slug")


if __name__ == "__main__":
    unittest.main()
