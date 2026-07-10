import { compile } from "./compile.js";
import { selectStrategy } from "./strategy.js";

/**
 * Create a reusable matcher from route definitions.
 *
 * Route definitions are compiled once. The matching strategy (regex or trie)
 * is selected automatically based on the compiled collection; the public API
 * and matching behavior stay the same either way.
 *
 * @param {Object<string, any>} routes - Map of pattern strings to opaque values.
 * @returns {import("./compile.js").MatchStrategy} The matcher instance exposing a match method.
 * @throws {TypeError} If routes is not a non-null plain object.
 */
export function createMatcher(routes) {
  if (routes == null || typeof routes !== "object" || Array.isArray(routes)) {
    throw new TypeError("createMatcher expects a non-null object of route definitions");
  }

  const compiled = compile(routes);
  const strategy = selectStrategy(compiled);

  return {
    match(pathname) {
      return strategy.match(pathname);
    },
  };
}

export default createMatcher;
