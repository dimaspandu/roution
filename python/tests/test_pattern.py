import unittest

from roution.pattern import (
    SEGMENT_PARAM,
    SEGMENT_STATIC,
    ParsedPattern,
    Segment,
    parse_pattern,
)


class PatternTest(unittest.TestCase):
    def test_wildcard(self):
        self.assertEqual(
            parse_pattern("*"),
            ParsedPattern(raw="*", is_wildcard=True, segments=[]),
        )

    def test_static_root(self):
        self.assertEqual(parse_pattern("/").segments, [])

    def test_static_segments(self):
        self.assertEqual(
            parse_pattern("/articles").segments,
            [Segment(type=SEGMENT_STATIC, value="articles")],
        )

    def test_dynamic_segments(self):
        self.assertEqual(
            parse_pattern("/articles/:slug").segments,
            [
                Segment(type=SEGMENT_STATIC, value="articles"),
                Segment(type=SEGMENT_PARAM, name="slug"),
            ],
        )

    def test_multiple_dynamic(self):
        self.assertEqual(
            parse_pattern("/users/:id/posts/:postId").segments,
            [
                Segment(type=SEGMENT_STATIC, value="users"),
                Segment(type=SEGMENT_PARAM, name="id"),
                Segment(type=SEGMENT_STATIC, value="posts"),
                Segment(type=SEGMENT_PARAM, name="postId"),
            ],
        )


if __name__ == "__main__":
    unittest.main()
