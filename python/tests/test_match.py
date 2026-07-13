import unittest

from roution.compile import compile
from roution.strategies.regex import create_regex_strategy
from roution.strategies.trie import create_trie_strategy


def build_strategies(routes):
    compiled = compile(routes)
    return [create_regex_strategy(compiled, True), create_trie_strategy(compiled, True)]


class MatchTest(unittest.TestCase):
    def test_static_root(self):
        for strategy in build_strategies({"/": "home"}):
            result = strategy.match("/")
            self.assertTrue(result["found"])
            self.assertEqual(result["value"], "home")
            self.assertEqual(result["params"], {})

    def test_static_route(self):
        for strategy in build_strategies({"/about": "about"}):
            self.assertEqual(strategy.match("/about")["route"], "/about")

    def test_dynamic_and_params(self):
        for strategy in build_strategies({"/articles/:slug": "page"}):
            result = strategy.match("/articles/javascript")
            self.assertEqual(result["route"], "/articles/:slug")
            self.assertEqual(result["params"], {"slug": "javascript"})

    def test_not_found(self):
        for strategy in build_strategies({"/about": "about"}):
            result = strategy.match("/missing")
            self.assertFalse(result["found"])
            self.assertIsNone(result["route"])
            self.assertIsNone(result["value"])

    def test_wildcard_fallback(self):
        for strategy in build_strategies({"*": "404"}):
            result = strategy.match("/articles/javascript")
            self.assertEqual(result["route"], "*")
            self.assertEqual(result["value"], "404")

    def test_static_priority_over_wildcard(self):
        for strategy in build_strategies({"/about": "about", "*": "404"}):
            self.assertEqual(strategy.match("/about")["value"], "about")

    def test_dynamic_priority_over_wildcard(self):
        for strategy in build_strategies({"/articles/:slug": "page", "*": "404"}):
            self.assertEqual(strategy.match("/articles/javascript")["value"], "page")

    def test_normalizes_query(self):
        for strategy in build_strategies({"/articles/:slug": "page"}):
            result = strategy.match("/articles/javascript?page=1")
            self.assertEqual(result["pathname"], "/articles/javascript")

    def test_normalizes_fragment(self):
        for strategy in build_strategies({"/articles": "page"}):
            result = strategy.match("/articles#section")
            self.assertEqual(result["pathname"], "/articles")

    def test_decodes_params(self):
        for strategy in build_strategies({"/files/:name": "file"}):
            result = strategy.match("/files/hello%20world")
            self.assertEqual(result["params"], {"name": "hello world"})

    def test_segment_count_must_match(self):
        for strategy in build_strategies({"/users/:id": "u"}):
            self.assertFalse(strategy.match("/users/1/posts")["found"])

    def test_static_priority_over_dynamic(self):
        for strategy in build_strategies(
            {"/users/admin": "admin", "/users/:id": "user"}
        ):
            result = strategy.match("/users/admin")
            self.assertEqual(result["route"], "/users/admin")
            self.assertEqual(result["params"], {})

    def test_identical_results(self):
        routes = {
            "/": "home",
            "/about": "about",
            "/articles": "list",
            "/articles/:slug": "article",
            "/users/:id/posts/:postId": "post",
            "/composed": ["a", "b"],
            "/login": {"auth": False},
            "/version": "1.0.0",
            "/numbers": [1, 2, 3],
            "*": "404",
        }
        compiled = compile(routes)
        regex = create_regex_strategy(compiled, True)
        trie = create_trie_strategy(compiled, True)
        samples = [
            "/",
            "/about",
            "/articles",
            "/articles/javascript?page=1&perPage=10",
            "/users/42/posts/7#comments",
            "/composed",
            "/login",
            "/version",
            "/numbers",
            "/articles/javascript/extra",
            "/missing/page",
        ]
        for pathname in samples:
            self.assertEqual(regex.match(pathname), trie.match(pathname))

    def test_query_disabled(self):
        compiled = compile({"/a": "v"})
        for strategy in [
            create_regex_strategy(compiled, False),
            create_trie_strategy(compiled, False),
        ]:
            self.assertEqual(strategy.match("/a?x=1")["query"], {})


if __name__ == "__main__":
    unittest.main()
