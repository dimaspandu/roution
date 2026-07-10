import { test } from "node:test";
import assert from "node:assert/strict";
import { compile } from "../src/compile.js";
import { selectStrategy, DYNAMIC_ROUTE_THRESHOLD } from "../src/strategy.js";

function routesWithDynamic(count) {
  const routes = { "*": "404" };
  for (let i = 0; i < count; i++) {
    routes[`/item/${i}/:slug`] = i;
  }
  return routes;
}

test("selects regex strategy below the threshold", () => {
  const compiled = compile(routesWithDynamic(DYNAMIC_ROUTE_THRESHOLD - 1));
  const strategy = selectStrategy(compiled);
  assert.equal(strategy.match("/item/0/a").route, "/item/0/:slug");
});

test("selects trie strategy at or above the threshold", () => {
  const compiled = compile(routesWithDynamic(DYNAMIC_ROUTE_THRESHOLD));
  const strategy = selectStrategy(compiled);
  assert.equal(strategy.match("/item/0/a").route, "/item/0/:slug");
});

test("both strategies resolve identically near the threshold", () => {
  const compiledLow = compile(routesWithDynamic(DYNAMIC_ROUTE_THRESHOLD - 1));
  const compiledHigh = compile(routesWithDynamic(DYNAMIC_ROUTE_THRESHOLD));
  const low = selectStrategy(compiledLow);
  const high = selectStrategy(compiledHigh);

  for (let i = 0; i < 5; i++) {
    const path = `/item/${i}/slug-${i}`;
    assert.deepEqual(low.match(path), high.match(path));
  }
});
