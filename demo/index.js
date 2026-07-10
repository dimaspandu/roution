import { createMatcher } from "../src/roution.js";

const routes = {
  "/": "public/index.html",

  "/about": "public/about.html",

  "/articles": "public/articles/index.html",

  "/articles/:slug": "public/articles/[slug].html",

  "/users/:id/posts/:postId": "public/users/[id]/posts/[postId].html",

  "/composed": ["public/header.html", "public/greetings.html", "public/footer.html"],

  "/login": {
    component: "LoginPage",
    auth: false,
  },

  "/dashboard": {
    component: "DashboardPage",
    auth: true,
    layout: "main",
  },

  "/version": "1.0.0",

  "/numbers": [1, 2, 3],

  "/health": () => ({ status: "ok" }),

  "/time": async () => new Date().toISOString(),

  "*": "public/404.html",
};

const matcher = createMatcher(routes);

const samples = [
  "/",
  "/about",
  "/articles",
  "/articles/javascript?page=1",
  "/users/42/posts/7#comments",
  "/composed",
  "/login",
  "/dashboard",
  "/version",
  "/numbers",
  "/health",
  "/time",
  "/articles/javascript/extra",
  "/missing/page",
  "/wild/card/path",
];

console.log("Route collection:");
console.log(JSON.stringify(Object.keys(routes), null, 2));
console.log("\nMatch samples:\n");

for (const pathname of samples) {
  const result = matcher.match(pathname);
  console.log(`pathname: ${pathname}`);
  console.log(`  found:    ${result.found}`);
  console.log(`  route:    ${result.route}`);
  console.log(`  pathname: ${result.pathname}`);
  console.log(`  params:   ${JSON.stringify(result.params)}`);
  console.log(`  value:    ${formatValue(result.value)}`);
  console.log("");
}

function formatValue(value) {
  if (typeof value === "function") {
    return `[Function ${value.name || "anonymous"}]`;
  }
  if (typeof value === "object" && value !== null) {
    return JSON.stringify(value);
  }
  return String(value);
}
