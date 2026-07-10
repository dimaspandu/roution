import { createRegexStrategy } from "./strategies/regex.js";
import { createTrieStrategy } from "./strategies/trie.js";

/**
 * Number of dynamic routes at or above which the trie strategy is selected
 * over the regex strategy. Considered an implementation detail.
 *
 * @type {number}
 */
export const DYNAMIC_ROUTE_THRESHOLD = 50;

/**
 * Select the most appropriate matching strategy for a compiled collection.
 *
 * The decision is automatic and based on the number of dynamic routes. A
 * small collection uses the regex strategy; a large collection uses the trie
 * strategy for better scaling. The public API and matching behavior remain
 * identical regardless of the chosen strategy.
 *
 * @param {import("./compile.js").CompiledRoutes} compiled - The compiled route collection.
 * @returns {import("./compile.js").MatchStrategy} The selected strategy.
 */
export function selectStrategy(compiled) {
  if (compiled.dynamicRoutes.length >= DYNAMIC_ROUTE_THRESHOLD) {
    return createTrieStrategy(compiled);
  }
  return createRegexStrategy(compiled);
}
