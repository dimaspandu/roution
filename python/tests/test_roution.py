import unittest

from roution import create_matcher
from roution.strategy import Options


class RoutionTest(unittest.TestCase):
    def test_quick_start(self):
        routes = {
            "/": "public/index.html",
            "/articles": "public/articles/index.html",
            "/articles/:slug": "public/articles/[slug].html",
            "/composed": ["public/header.html", "public/greetings.html", "public/footer.html"],
            "*": "public/404.html",
        }
        result = create_matcher(routes).match("/articles/javascript?page=1")
        self.assertEqual(
            result,
            {
                "found": True,
                "pathname": "/articles/javascript",
                "route": "/articles/:slug",
                "params": {"slug": "javascript"},
                "query": {"page": "1"},
                "value": "public/articles/[slug].html",
            },
        )

    def test_reusable(self):
        matcher = create_matcher({"/a": 1, "/b": 2, "*": 0})
        self.assertEqual(matcher.match("/a")["value"], 1)
        self.assertEqual(matcher.match("/b")["value"], 2)
        self.assertEqual(matcher.match("/unknown")["value"], 0)

    def test_default_options(self):
        matcher = create_matcher({"/articles/:slug": "v"})
        self.assertEqual(
            matcher.match("/articles/javascript?page=1")["route"], "/articles/:slug"
        )

    def test_options_as_dict(self):
        matcher = create_matcher({"/articles/:slug": "v"}, {"strategy": "trie"})
        self.assertEqual(
            matcher.match("/articles/javascript?page=1")["route"], "/articles/:slug"
        )

    def test_query_disabled(self):
        matcher = create_matcher({"/a": "v"}, Options(query=False))
        self.assertEqual(matcher.match("/a?x=1")["query"], {})

    def test_invalid_routes(self):
        with self.assertRaises(TypeError):
            create_matcher(None)
        with self.assertRaises(TypeError):
            create_matcher([])

    def test_nil_routes(self):
        matcher = create_matcher({})
        self.assertFalse(matcher.match("/anything")["found"])


if __name__ == "__main__":
    unittest.main()
