from dataclasses import dataclass

SEGMENT_STATIC = "static"
SEGMENT_PARAM = "param"


@dataclass
class Segment:
    """A single segment of a parsed route pattern."""

    type: str
    name: str = ""
    value: str = ""


@dataclass
class ParsedPattern:
    """The compiled representation of a route pattern."""

    raw: str
    is_wildcard: bool
    segments: list


def parse_pattern(pattern):
    """Parse a route pattern into structured segments.

    The wildcard pattern "*" is recognized as a standalone fallback route.
    Static segments become Segment(type="static", value=...) and ":name"
    segments become Segment(type="param", name=...).
    """
    if pattern == "*":
        return ParsedPattern(raw=pattern, is_wildcard=True, segments=[])

    segments = []
    for part in pattern.split("/"):
        if part == "":
            continue
        if part.startswith(":"):
            segments.append(Segment(type=SEGMENT_PARAM, name=part[1:]))
        else:
            segments.append(Segment(type=SEGMENT_STATIC, value=part))
    return ParsedPattern(raw=pattern, is_wildcard=False, segments=segments)
