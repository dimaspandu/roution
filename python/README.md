# ROUTION (Python)

The Python implementation of ROUTION, a lightweight, runtime-independent route
resolution engine. It follows the same specification and matching behavior as
the JavaScript and Go implementations.

ROUTION compiles route definitions into a reusable matcher that resolves
incoming pathnames and returns structured matching results. It is framework
agnostic and has zero third-party dependencies (standard library only).

For the full specification and cross-language overview, see the
[root README](../README.md).

## Requirements

- Python 3.10 or newer.

## Installation

Install from source (editable, for local development):

```bash
pip install -e .
```

Or copy the `roution/` package into your project. The package depends only on
the Python standard library.

```python
from roution import create_matcher
```

## Quick Start

```python
from roution import create_matcher

routes = {
    "/":               "public/index.html",
    "/articles":       "public/articles/index.html",
    "/articles/:slug": "public/articles/[slug].html",
    "*":               "public/404.html",
}

matcher = create_matcher(routes)

result = matcher.match("/articles/javascript?page=1")
print(result.found, result.route, result.params, result.query, result.value)
# True /articles/:slug {'slug': 'javascript'} {'page': '1'} public/articles/[slug].html
```

## API

### `create_matcher(routes: dict[str, Any], options: Options | None = None) Matcher`

Compiles route definitions into a reusable matcher. Compilation happens once;
the returned matcher can be reused for the lifetime of the application. A `None`
routes value is tolerated (it yields no routes).

`Options` fields (all optional, with safe defaults via `Options()`):

| Field              | Type   | Default  | Description                                                               |
| ------------------ | ------ | -------- | ------------------------------------------------------------------------- |
| `query`            | `bool` | `True`   | When `True`, the result `query` is populated. When `False`, it is empty.  |
| `strategy`         | `str`  | `"auto"` | `Strategy.AUTO`, `Strategy.REGEX`, or `Strategy.TRIE`.                    |
| `dynamic_threshold`| `int`  | `50`     | Dynamic route count at or above which `"auto"` selects the trie strategy. |

### `matcher.match(pathname: str) MatchResult`

Resolves an incoming pathname. The pathname is normalized (query string and URL
fragment are stripped) before matching.

`MatchResult` fields:

| Field      | Type                | Description                                   |
| ---------- | ------------------- | --------------------------------------------- |
| `found`    | `bool`              | Whether a route was matched.                  |
| `pathname` | `str`               | The normalized pathname.                      |
| `route`    | `str`               | The matched route pattern, or `""`.           |
| `params`   | `dict[str, str]`    | Extracted route parameters.                   |
| `query`    | `dict[str, Any]`    | Parsed query string parameters.               |
| `value`    | `Any`               | The value associated with the matched route.  |

## Demo

The demo runs the same sample pathnames through two matchers, one without a
wildcard and one with a wildcard, so the difference is easy to observe.

```bash
python examples/basic/main.py
```

## Testing

Tests use the standard `unittest` package.

```bash
python -m unittest discover -s tests
```

## Project Structure

```text
python/
  roution/
    __init__.py           # public API (create_matcher, Options, MatchResult)
    normalize.py          # pathname normalization + query parsing
    pattern.py            # route pattern parsing
    compile.py            # one-time route compilation
    matcher.py            # matcher implementation
    strategy.py           # automatic strategy selection
    result.py             # MatchResult dataclass
    strategies/
      regex.py            # regex-backed matching strategy
      trie.py             # trie-backed matching strategy
  examples/
    basic/main.py         # runnable demo
  tests/                  # unit tests
  pyproject.toml
```

## Notes

The compiler sorts dynamic routes by pattern to guarantee stable, reproducible
matching priority. The public API and matching behavior are identical to the
other implementations.

## License

MIT
