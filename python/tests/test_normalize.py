import unittest

from roution.normalize import normalize_pathname, split_segments, parse_query


class NormalizeTest(unittest.TestCase):
    def test_strips_query_string(self):
        self.assertEqual(normalize_pathname("/articles?page=1"), "/articles")

    def test_strips_url_fragment(self):
        self.assertEqual(normalize_pathname("/articles#comments"), "/articles")

    def test_strips_both(self):
        self.assertEqual(
            normalize_pathname("/articles?page=1&perPage=10#comments"), "/articles"
        )

    def test_returns_path_untouched(self):
        self.assertEqual(normalize_pathname("/articles/javascript"), "/articles/javascript")

    def test_none_returns_empty(self):
        self.assertEqual(normalize_pathname(None), "")

    def test_split_segments_ignores_slashes(self):
        self.assertEqual(split_segments("/a/b/c/"), ["a", "b", "c"])

    def test_split_segments_root(self):
        self.assertEqual(split_segments("/"), [])

    def test_parse_query_empty_without_query(self):
        self.assertEqual(parse_query("/articles"), {})

    def test_parse_query_empty_for_empty_query(self):
        self.assertEqual(parse_query("/articles?"), {})

    def test_parse_query_single(self):
        self.assertEqual(parse_query("/articles?page=1"), {"page": "1"})

    def test_parse_query_multiple(self):
        self.assertEqual(
            parse_query("/articles?page=1&perPage=10"), {"page": "1", "perPage": "10"}
        )

    def test_parse_query_repeated_keys(self):
        self.assertEqual(parse_query("/search?tag=a&tag=b"), {"tag": ["a", "b"]})

    def test_parse_query_ignores_fragment(self):
        self.assertEqual(parse_query("/articles?page=1#top"), {"page": "1"})

    def test_parse_query_decodes_encoding(self):
        self.assertEqual(parse_query("/search?q=hello%20world"), {"q": "hello world"})


if __name__ == "__main__":
    unittest.main()
