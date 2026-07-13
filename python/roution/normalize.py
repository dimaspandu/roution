from urllib.parse import parse_qsl


def normalize_pathname(input):
    """Strip the URL fragment and query string from an incoming pathname.

    Everything after the first "#" or "?" is discarded. No other transformation
    is applied, so the result remains a raw pathname.
    """
    if input is None:
        return ""
    path = str(input)
    hash_index = path.find("#")
    if hash_index != -1:
        path = path[:hash_index]
    query_index = path.find("?")
    if query_index != -1:
        path = path[:query_index]
    return path


def split_segments(pathname):
    """Split a pathname into its non-empty segments."""
    return [segment for segment in pathname.split("/") if segment != ""]


def parse_query(input):
    """Parse the query string of an incoming pathname into a dict.

    Values follow standard HTTP semantics: single values are strings and
    repeated keys are collected into a list. URL encoding is decoded
    automatically and the URL fragment is ignored. Returns {} when there is
    no query string.
    """
    if input is None:
        return {}
    path = str(input)
    hash_index = path.find("#")
    if hash_index != -1:
        path = path[:hash_index]
    query_index = path.find("?")
    if query_index == -1:
        return {}
    query = path[query_index + 1:]
    if query == "":
        return {}

    pairs = parse_qsl(query, keep_blank_values=True)
    result = {}
    for key, value in pairs:
        if key in result:
            existing = result[key]
            if isinstance(existing, list):
                existing.append(value)
            else:
                result[key] = [existing, value]
        else:
            result[key] = value
    return result
