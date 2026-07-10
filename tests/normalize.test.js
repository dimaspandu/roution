import { test } from "node:test";
import assert from "node:assert/strict";
import { normalizePathname, splitSegments } from "../src/normalize.js";

test("strips query string", () => {
  assert.equal(normalizePathname("/articles?page=1"), "/articles");
});

test("strips url fragment", () => {
  assert.equal(normalizePathname("/articles#comments"), "/articles");
});

test("strips both query and fragment", () => {
  assert.equal(normalizePathname("/articles?page=1&perPage=10#comments"), "/articles");
});

test("returns path untouched when no query or fragment", () => {
  assert.equal(normalizePathname("/articles/javascript"), "/articles/javascript");
});

test("converts null/undefined to empty string", () => {
  assert.equal(normalizePathname(null), "");
  assert.equal(normalizePathname(undefined), "");
});

test("coerces non-string input to string", () => {
  assert.equal(normalizePathname(123), "123");
});

test("splitSegments ignores leading and trailing slashes", () => {
  assert.deepEqual(splitSegments("/a/b/c/"), ["a", "b", "c"]);
});

test("splitSegments returns empty array for root", () => {
  assert.deepEqual(splitSegments("/"), []);
});
