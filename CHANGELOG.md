# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
