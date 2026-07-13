from roution import create_matcher

ROUTES = {
    "/": "public/index.html",
    "/about": "public/about.html",
    "/articles": "public/articles/index.html",
    "/articles/:slug": "public/articles/[slug].html",
    "/users/:id/posts/:postId": "public/users/[id]/posts/[postId].html",
    "/composed": ["public/header.html", "public/greetings.html", "public/footer.html"],
    "/login": {"component": "LoginPage", "auth": False},
    "/dashboard": {"component": "DashboardPage", "auth": True, "layout": "main"},
    "/version": "1.0.0",
    "/numbers": [1, 2, 3],
    "/health": lambda: {"status": "ok"},
    "*": "public/404.html",
}

SAMPLES = [
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
    "/articles/javascript/extra",
    "/missing/page",
    "/wild/card/path",
]


def main():
    matcher = create_matcher(ROUTES)
    for pathname in SAMPLES:
        result = matcher.match(pathname)
        print(f"pathname: {pathname}")
        print(f"  found:    {result['found']}")
        print(f"  route:    {result['route']}")
        print(f"  pathname: {result['pathname']}")
        print(f"  params:   {result['params']}")
        print(f"  query:    {result['query']}")
        print(f"  value:    {result['value']}")
        print()


if __name__ == "__main__":
    main()
