# ROUTION (Go)

The Go implementation of ROUTION, a lightweight, runtime-independent route
resolution engine. It follows the same specification and matching behavior as
the JavaScript reference implementation.

ROUTION compiles route definitions into a reusable matcher that resolves
incoming pathnames and returns structured matching results. It is framework
agnostic and has zero third-party dependencies (standard library only).

For the full specification and cross-language overview, see the
[root README](../README.md).

## Requirements

- Go 1.21 or newer (uses the standard library `maps` package in the demo).

## Installation

```bash
go get github.com/dimaspandu/roution
```

Then import it:

```go
import roution "github.com/dimaspandu/roution"
```

## Quick Start

```go
package main

import (
	"fmt"

	roution "github.com/dimaspandu/roution"
)

func main() {
	routes := map[string]any{
		"/":               "public/index.html",
		"/articles":       "public/articles/index.html",
		"/articles/:slug": "public/articles/[slug].html",
		"*":               "public/404.html",
	}

	matcher := roution.CreateMatcher(routes)

	result := matcher.Match("/articles/javascript?page=1")
	fmt.Println(result.Found, result.Route, result.Params, result.Query, result.Value)
	// true /articles/:slug map[slug:javascript] map[page:1] public/articles/[slug].html
}
```

## API

### `CreateMatcher(routes map[string]any, opts ...Options) Matcher`

Compiles route definitions into a reusable matcher. Compilation happens once;
the returned matcher can be reused for the lifetime of the application. A `nil`
routes map is tolerated (it yields no routes).

`Options` fields (all optional, with safe defaults via `DefaultOptions()`):

| Field              | Type     | Default  | Description                                                                    |
| ------------------ | -------- | -------- | ------------------------------------------------------------------------------ |
| `Query`            | `bool`   | `true`   | When `true`, the result `Query` is populated. When `false`, it is empty.       |
| `Strategy`         | `string` | `"auto"` | `StrategyAuto`, `StrategyRegex`, or `StrategyTrie`.                             |
| `DynamicThreshold` | `int`    | `50`     | Dynamic route count at or above which `"auto"` selects the trie strategy.       |

### `matcher.Match(pathname string) MatchResult`

Resolves an incoming pathname. The pathname is normalized (query string and URL
fragment are stripped) before matching.

`MatchResult` fields:

| Field      | Type                | Description                                   |
| ---------- | ------------------- | --------------------------------------------- |
| `Found`    | `bool`              | Whether a route was matched.                  |
| `Pathname` | `string`            | The normalized pathname.                      |
| `Route`    | `string`            | The matched route pattern, or `""`.           |
| `Params`   | `map[string]string` | Extracted route parameters.                   |
| `Query`    | `map[string]any`    | Parsed query string parameters.               |
| `Value`    | `any`               | The value associated with the matched route.  |

## Demo

The demo runs the same sample pathnames through two matchers, one without a
wildcard and one with a wildcard, so the difference is easy to observe.

```bash
go run ./demo
```

## Testing

Tests use the standard `testing` package.

```bash
go test ./...
```

## Project Structure

```text
go/
  normalize.go            # pathname normalization + query parsing
  pattern.go              # route pattern parsing
  compile.go              # one-time route compilation
  matcher.go              # public API (CreateMatcher, Options, MatchResult)
  strategy.go             # automatic strategy selection
  regex.go                # regex-backed matching strategy
  trie.go                 # trie-backed matching strategy
  demo/                   # runnable demo
  *_test.go               # unit tests
  go.mod
```

## Notes

Go map iteration order is not deterministic, so the compiler sorts dynamic
routes by pattern to guarantee stable, reproducible matching priority. The
public API and matching behavior are identical to the other implementations.

## License

MIT
