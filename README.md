# ROUTION

A lightweight, runtime-independent route resolution engine for JavaScript.

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

The same matcher can be used in:

* Browser
* Node.js
* Bun
* Deno
* Cloudflare Workers
* Service Workers
* Edge runtimes
* Test environments
* Any JavaScript runtime supporting standard ECMAScript Modules (ESM)

ROUTION has zero runtime dependencies and makes no assumptions about where it runs.

## Philosophy

ROUTION follows a simple set of design principles.

* Compile once. Match many.
* Runtime independent.
* Framework agnostic.
* Single responsibility.
* Zero runtime dependencies.
* Unit-test friendly.

The public API is intentionally small while remaining suitable for projects of any size.

## Features

* Static route matching
* Dynamic route parameters
* Wildcard fallback routes
* Automatic pathname normalization
* Query strings are ignored
* URL fragments are ignored
* Compile once and reuse the matcher
* Runtime independent
* Framework agnostic
* Zero runtime dependencies
* Predictable matching results
* Supports arbitrary route values

## Installation

ROUTION has zero runtime dependencies.

Import the module directly into your application.

### ES Module

```javascript
import { createMatcher } from "./roution.js";
```

### Browser

```html
<script type="module">
import { createMatcher } from "./roution.js";
</script>
```

### Node.js

```javascript
import { createMatcher } from "./roution.js";
```

No package manager is required.

## Quick Start

```javascript
import { createMatcher } from "./roution.js";

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
  value: "public/articles/[slug].html"
}
```

## Demo

A runnable sample is available in the `demo/` folder. It defines a varied route
collection (static, dynamic, multi-parameter, wildcard, arrays, objects,
functions) and prints the matching result for several pathnames, including
paths that fall through to the wildcard route.

Run it with Node.js (no install step required):

```bash
node demo/index.js
```

Or via the npm script:

```bash
npm run demo
```

The demo imports directly from `src/roution.js`:

```javascript
import { createMatcher } from "../src/roution.js";

const matcher = createMatcher(routes);

const result = matcher.match("/articles/javascript?page=1");
```

## Route Definitions

Routes are defined as an object where each property key is a route pattern and each property value is application-defined data.

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

A route value may be any JavaScript value, including:

* String
* Number
* Boolean
* Object
* Array
* Function
* Async Function
* Class
* Promise
* Component
* Symbol
* Custom data structure

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

A successful match returns a structured object.

```javascript
{
  found: true,
  pathname: "/articles/javascript",
  route: "/articles/:slug",
  params: {
    slug: "javascript"
  },
  value: "public/articles/[slug].html"
}
```

| Property   | Description                                  |
| ---------- | -------------------------------------------- |
| `found`    | Whether a route was matched                  |
| `pathname` | The normalized pathname                      |
| `route`    | The matching route definition                |
| `params`   | Extracted route parameters                   |
| `value`    | The original value associated with the route |

If no route matches, `found` is `false`.

## Matching Strategy

ROUTION automatically selects the most appropriate matching strategy based on the compiled route collection.

This decision is automatic, requires no configuration, and is considered an implementation detail.

Regardless of the selected strategy, the public API and matching behavior remain the same.

## API Overview

### `createMatcher()`

Compiles route definitions into a reusable matcher.

```javascript
const matcher = createMatcher(routes);
```

Compilation happens only once.

The resulting matcher can be reused throughout the lifetime of the application.

### `match()`

Matches an incoming pathname.

```javascript
matcher.match("/articles/javascript?page=1");
```

The incoming pathname is automatically normalized before matching.

The method returns a structured matching result.

## Design Goals

ROUTION is intentionally small.

Its responsibility begins and ends with pathname resolution.

It intentionally avoids:

* Browser APIs
* Node.js APIs
* HTTP servers
* Framework-specific APIs
* Rendering systems
* Navigation
* History management
* Middleware execution
* Redirect handling

As long as a JavaScript runtime supports standard ECMAScript Modules, ROUTION behaves consistently.

The same matcher can be reused across browsers, servers, edge runtimes, and test environments without modification.

## Future Roadmap

The project aims to remain small and stable while improving the developer experience.

Potential future improvements include:

* TypeScript type definitions
* Improved type inference
* Additional route pattern capabilities
* Performance optimizations
* Expanded test coverage
* Enhanced developer tooling

The public API will continue to prioritize simplicity, portability, and long-term stability.

## License

MIT
