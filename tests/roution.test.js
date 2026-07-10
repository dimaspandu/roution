import { test } from "node:test";
import assert from "node:assert/strict";
import { createMatcher } from "../src/roution.js";

test("createMatcher returns a matcher with match()", () => {
  const matcher = createMatcher({ "/": "home" });
  assert.equal(typeof matcher.match, "function");
});

test("throws on invalid routes argument", () => {
  assert.throws(() => createMatcher(null), TypeError);
  assert.throws(() => createMatcher([]), TypeError);
  assert.throws(() => createMatcher("routes"), TypeError);
});

test("reproduces the README quick start example", () => {
  const routes = {
    "/": "public/index.html",
    "/articles": "public/articles/index.html",
    "/articles/:slug": "public/articles/[slug].html",
    "/composed": ["public/header.html", "public/greetings.html", "public/footer.html"],
    "*": "public/404.html",
  };

  const matcher = createMatcher(routes);
  const result = matcher.match("/articles/javascript?page=1");

  assert.deepEqual(result, {
    found: true,
    pathname: "/articles/javascript",
    route: "/articles/:slug",
    params: { slug: "javascript" },
    value: "public/articles/[slug].html",
  });
});

test("match is reusable across many calls", () => {
  const matcher = createMatcher({ "/a": 1, "/b": 2, "*": 0 });
  assert.equal(matcher.match("/a").value, 1);
  assert.equal(matcher.match("/b").value, 2);
  assert.equal(matcher.match("/unknown").value, 0);
});

test("supports arbitrary opaque values", () => {
  const fn = () => "x";
  const matcher = createMatcher({ "/fn": fn });
  assert.equal(matcher.match("/fn").value, fn);
});
