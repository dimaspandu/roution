# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.2] - 2026-07-10

### Added

- Optional `options` argument for `createMatcher(routes, options)`:
  - `query` (boolean, default `true`): when `false`, the `query` field in the
    matching result is always an empty object.
  - `strategy` (`"auto"` | `"regex"` | `"trie"`, default `"auto"`): overrides
    the automatic strategy selection with an explicit matching algorithm.
  - `dynamicThreshold` (number, default `50`): dynamic route count at or above
    which the `"auto"` strategy selects the trie strategy.
- All options are optional and keep the previous default behavior, so existing
  usage is unchanged.

## [1.0.1] - 2026-07-10

### Added

- `query` field in the matching result. The query string is now parsed and
  exposed as `query` using standard HTTP semantics: values are strings, repeated
  keys become arrays, URL encoding is decoded, and the URL fragment is ignored.
  When there is no query string, `query` is an empty object.

## [1.0.0] - 2026-07-10

### Added

- Initial release of ROUTION, a lightweight, runtime-independent route
  resolution engine for JavaScript.
- Public API: `createMatcher(routes)` returning a matcher with a `match(pathname)`
  method.
- Static route matching (exact pathname lookup).
- Dynamic route parameters using the `:name` syntax (e.g. `/articles/:slug`).
- Wildcard fallback route (`*`) for unmatched pathnames.
- Automatic pathname normalization that strips query strings and URL fragments
  before matching.
- Structured matching result containing `found`, `pathname`, `route`, `params`,
  and `value`.
- Opaque route values: any JavaScript value (string, object, array, function,
  etc.) is returned as-is and never executed or interpreted by the engine.
- Layered, unit-testable architecture:
  - `src/normalize.js` - pathname normalization and segment splitting.
  - `src/pattern.js` - route pattern parsing.
  - `src/compile.js` - one-time compilation into a reusable route collection.
  - `src/strategy.js` - automatic strategy selection.
  - `src/strategies/regex.js` - regex-backed matching strategy.
  - `src/strategies/trie.js` - trie-backed matching strategy for large route sets.
  - `src/roution.js` - public API entry point.
- Automatic matching strategy selection based on the number of dynamic routes,
  switching from the regex strategy to the trie strategy at a configurable
  threshold (`DYNAMIC_ROUTE_THRESHOLD`).
- Full unit test suite under `tests/` using the built-in Node.js test runner
  (zero runtime dependencies, no `npm install` required).
- Runnable demo under `demo/` showcasing various route value types and match
  results, including wildcard fallback.
- `package.json` scripts: `npm test` and `npm run demo`.
- MIT license.
