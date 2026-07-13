# ROUTION

A lightweight, runtime-independent route resolution engine, available across multiple programming languages.

ROUTION compiles route definitions into a reusable matcher that resolves incoming pathnames and returns structured matching results.

It is framework agnostic, runtime independent, dependency free, and designed around a single responsibility:

> Resolve an incoming pathname against compiled route definitions and return the matching result.

ROUTION is **not** a router.

It does **not** perform:

* HTTP handling
* Rendering
* Browser navigation
* History management
* Middleware execution
* Redirects

It can serve as the route resolution layer of a router, but routing itself is outside the scope of the project.

Everything beyond pathname resolution is intentionally left to your application.

## Overview

Many applications need route matching without adopting a complete routing framework.

ROUTION provides a small, focused engine that compiles route definitions once and reuses them for efficient pathname matching.

Compile once. Match many.

The same matching behavior is shared across every language implementation, so the way routes are defined and resolved stays consistent whether you use JavaScript, Go, or another supported language.

## Language Implementations

ROUTION is implemented per language, each as a self-contained module inside its own folder. The JavaScript implementation is the reference; other languages follow the same specification and matching behavior.

| Language   | Path    | Status      | Notes                                                      |
| ---------- | ------- | ----------- | ---------------------------------------------------------- |
| JavaScript | `js/`   | Stable      | Reference implementation (ESM, zero dependencies).         |
| Go         | `go/`   | Stable      | Full implementation (compile, regex/trie strategies, options). |
| PHP        | `php/`  | Planned     | Not started yet.                                           |

Each language folder is independent: it ships its own source, tests, and (where applicable) examples, and has no dependency on the other folders.

## Project Layout

```text
roution/
  js/                 # JavaScript (ESM) reference implementation
    src/              # engine source
    tests/            # unit tests (node:test)
    examples/         # runnable example
    package.json
  go/                 # Go implementation (in progress)
    normalize.go      # pathname normalization + query parsing
    pattern.go        # route pattern parsing
    internal/         # matcher strategies (regex, trie) - in progress
    examples/         # runnable example (in progress)
    *_test.go         # unit tests
  README.md
  LICENSE
  CHANGELOG.md
```

## Philosophy

ROUTION follows a simple set of design principles.

* Compile once. Match many.
* Runtime independent.
* Framework agnostic.
* Language agnostic.
* Single responsibility.
* Zero runtime dependencies.
* Unit-test friendly.

The public API of each implementation is intentionally small while remaining suitable for projects of any size.

## Features

* Static route matching
* Dynamic route parameters
* Wildcard fallback routes
* Automatic pathname normalization
* Query string parsing (standard HTTP semantics)
* URL fragments are ignored
* Compile once and reuse the matcher
* Runtime independent
* Framework agnostic
* Language agnostic
* Zero runtime dependencies
* Predictable matching results
* Supports arbitrary route values

## Installation

Each implementation is dependency free. Import the module from the language folder directly into your application.

### JavaScript (ESM)

```javascript
import { createMatcher } from "./js/src/roution.js";
```

### Browser

```html
<script type="module">
import { createMatcher } from "./js/src/roution.js";
</script>
```

### Node.js

```javascript
import { createMatcher } from "./js/src/roution.js";
```

No package manager is required.

## Quick Start

```javascript
import { createMatcher } from "./js/src/roution.js";

const routes = {
  "/": "public/index.html",

  "/articles": "public/articles/index.html",

  "/articles/:slug": "public/articles/[slug].html",

  "/composed": [
    "public/header.html",
    "public/greetings.html",
    "public/footer.html",
  ],

  "*": "public/404.html",
};

const matcher = createMatcher(routes);

const result = matcher.match("/articles/javascript?page=1");
```

Returns:

```javascript
{
  found: true,
  pathname: "/articles/javascript",
  route: "/articles/:slug",
  params: {
    slug: "javascript"
  },
  query: {
    page: "1"
  },
  value: "public/articles/[slug].html"
}
```

### Go

```go
package main

import (
	roution "github.com/dimaspandu/roution"
)

func main() {
	routes := map[string]any{
		"/":               "public/index.html",
		"/articles/:slug": "public/articles/[slug].html",
		"*":               "public/404.html",
	}

	matcher := roution.CreateMatcher(routes)

	result := matcher.Match("/articles/javascript?page=1")
	_ = result
}
```

The result mirrors the JavaScript shape: `Found`, `Pathname`, `Route`, `Params`
(`map[string]string`), `Query` (`map[string]any`), and `Value` (`any`).

## Demo

A runnable JavaScript sample is available in `js/examples/`. It defines a varied route collection (static, dynamic, multi-parameter, wildcard, arrays, objects, functions) and prints the matching result for several pathnames, including paths that fall through to the wildcard route.

Run it with Node.js (no install step required):

```bash
node js/examples/index.js
```

Or via the npm script from the `js/` folder:

```bash
cd js
npm run demo
```

The demo imports directly from `js/src/roution.js`:

```javascript
import { createMatcher } from "../src/roution.js";

const matcher = createMatcher(routes);

const result = matcher.match("/articles/javascript?page=1");
```

### Go

A runnable Go sample is available in `go/examples/basic/`. It builds the same
varied route collection and prints the matching result for several pathnames.

```bash
cd go
go run ./examples/basic
```

## Route Definitions

Routes are defined as a map where each key is a route pattern and each value is application-defined data. The exact syntax depends on the language (an object literal in JavaScript, a `map` in Go, an associative array in PHP), but the pattern grammar is identical across implementations.

```javascript
const routes = {
  "/": "public/index.html",

  "/articles": "public/articles/index.html",

  "/articles/:slug": "public/articles/[slug].html",

  "/composed": [
    "public/header.html",
    "public/greetings.html",
    "public/footer.html",
  ],

  "*": "public/404.html",
};
```

Unlike a traditional router, ROUTION does not assign any meaning to route values.

A value is treated as **opaque data**.

ROUTION does not assume that a value represents:

* an HTML file
* a component
* a request handler
* metadata
* configuration
* middleware
* or any other specific type

A route value may be any value supported by the host language, including strings, numbers, booleans, objects, arrays, functions, classes, or custom data structures.

ROUTION never executes, renders, imports, instantiates, awaits, or interprets the value.

Its only responsibility is returning the value associated with the matched route.

### Static File Server

```javascript
const routes = {
  "/": "public/index.html",
  "/about": "public/about.html",
};
```

The returned value can be used by a static file server to locate and serve a file.

### SPA Router

```javascript
const routes = {
  "/": HomePage,
  "/about": AboutPage,
  "/articles/:slug": ArticlePage,
};
```

The returned value can be passed to a rendering layer responsible for displaying the appropriate component.

### API Dispatcher

```javascript
const routes = {
  "/users": getUsers,
  "/users/:id": getUser,
};
```

The application may invoke the returned function after a successful match.

ROUTION never calls it automatically.

### Metadata

```javascript
const routes = {
  "/login": {
    component: LoginPage,
    auth: false,
  },

  "/dashboard": {
    component: DashboardPage,
    auth: true,
    layout: "main",
  },
};
```

Applications may associate arbitrary metadata with a route and decide how to use it after matching.

### Middleware / Callback

```javascript
const routes = {
  "/health": () => ({
    status: "ok",
  }),

  "/time": async () => {
    return new Date().toISOString();
  },
};
```

Functions and async functions are treated like any other value.

They are returned as-is and are never executed by ROUTION.

### Mixed Values

```javascript
const routes = {
  "/version": "1.0.0",

  "/config": {
    cache: true,
    ttl: 300,
  },

  "/numbers": [1, 2, 3],

  "/symbol": Symbol("example"),
};
```

Different routes may use completely different value types within the same matcher.

ROUTION makes no assumptions about the contents of those values.

All of the examples above are equally valid because ROUTION is only responsible for pathname resolution.

Once a matching route is found, the engine returns the associated value without modification.

What happens after that is entirely outside the scope of ROUTION.

Whether the value is rendered, executed, imported, awaited, inspected, or ignored is the responsibility of the application or framework built on top of ROUTION.

## Pathname Normalization

Incoming pathnames are normalized automatically before matching.

Query strings and URL fragments are ignored.

For example:

Input:

```text
/articles?page=1&perPage=10#comments
```

Normalized pathname:

```text
/articles
```

This normalization happens automatically before route matching.

## Matching Result

A successful match returns a structured result (the concrete shape is expressed in the host language, but the fields are the same everywhere).

```javascript
{
  found: true,
  pathname: "/articles/javascript",
  route: "/articles/:slug",
  params: {
    slug: "javascript"
  },
  query: {
    page: "1"
  },
  value: "public/articles/[slug].html"
}
```

| Property   | Description                                         |
| ---------- | --------------------------------------------------- |
| `found`    | Whether a route was matched                         |
| `pathname` | The normalized pathname                             |
| `route`    | The matching route definition                       |
| `params`   | Extracted route parameters                          |
| `query`    | Parsed query string parameters                      |
| `value`    | The original value associated with the route        |

If no route matches, `found` is `false`.

## Query Parameters

In addition to `pathname` normalization, ROUTION parses the query string and exposes it as the `query` property of the matching result.

Query values follow standard HTTP semantics:

* Values are kept as strings.
* A key that appears more than once is collected into an array.
* URL encoding is decoded automatically.
* The URL fragment is ignored.

Example:

```javascript
matcher.match("/search?q=hello%20world&tag=a&tag=b");
```

Returns:

```javascript
{
  found: true,
  pathname: "/search",
  route: "/search",
  params: {},
  query: {
    q: "hello world",
    tag: ["a", "b"]
  },
  value: "public/search.html"
}
```

When there is no query string, `query` is an empty object/map.

## Matching Strategy

ROUTION automatically selects the most appropriate matching strategy based on the compiled route collection.

This decision is automatic, requires no configuration, and is considered an implementation detail.

Regardless of the selected strategy, the public API and matching behavior remain the same.

## API Overview (JavaScript reference)

### `createMatcher()`

Compiles route definitions into a reusable matcher.

```javascript
const matcher = createMatcher(routes);
```

Compilation happens only once.

The resulting matcher can be reused throughout the lifetime of the application.

### `createMatcher(routes, options)`

`createMatcher()` accepts an optional second argument to tune behavior. All options are optional and have safe defaults, so existing usage is unchanged.

```javascript
const matcher = createMatcher(routes, {
  query: true,          // include the parsed query in the result (default: true)
  strategy: "auto",     // "auto" (default) | "regex" | "trie"
  dynamicThreshold: 50  // dynamic route count at/above which "auto" selects trie
});
```

| Option            | Type                            | Default | Description                                                                 |
| ----------------- | ------------------------------- | ------- | --------------------------------------------------------------------------- |
| `query`           | `boolean`                       | `true`  | When `true`, the result includes the parsed `query` object. When `false`, `query` is always `{}`. |
| `strategy`        | `"auto" \| "regex" \| "trie"`   | `"auto"`| `"auto"` picks based on route count; `"regex"`/`"trie"` force a specific algorithm. |
| `dynamicThreshold`| `number`                        | `50`    | Dynamic route count at or above which `"auto"` selects the trie strategy.  |

### `match()`

Matches an incoming pathname.

```javascript
matcher.match("/articles/javascript?page=1");
```

The incoming pathname is automatically normalized before matching.

The method returns a structured matching result.

## Testing

Each language folder ships its own test suite.

### JavaScript

```bash
cd js
npm test
```

### Go

```bash
cd go
go test ./...
```

## Design Goals

ROUTION is intentionally small.

Its responsibility begins and ends with pathname resolution.

It intentionally avoids:

* Language-specific standard libraries where avoidable
* HTTP servers
* Framework-specific APIs
* Rendering systems
* Navigation
* History management
* Middleware execution
* Redirect handling

As long as a runtime supports the implementation's standard module system, ROUTION behaves consistently.

The same matcher can be reused across browsers, servers, edge runtimes, and test environments without modification.

## Future Roadmap

The project aims to remain small and stable while improving the developer experience and expanding language coverage.

Potential future improvements include:

* Additional language implementations (PHP and others)
* TypeScript type definitions
* Improved type inference
* Additional route pattern capabilities
* Performance optimizations
* Expanded test coverage
* Enhanced developer tooling

The public API will continue to prioritize simplicity, portability, and long-term stability.

## License

MIT
