import unittest

from roution.compile import compile


class CompileTest(unittest.TestCase):
    def test_classifies_root(self):
        compiled = compile({"/": "home"})
        self.assertEqual(compiled.static_routes["/"].value, "home")
        self.assertEqual(len(compiled.dynamic_routes), 0)

    def test_classifies_static(self):
        compiled = compile({"/articles": "a", "/about": "b"})
        self.assertEqual(len(compiled.static_routes), 2)
        self.assertEqual(compiled.static_routes["/about"].value, "b")

    def test_classifies_dynamic(self):
        compiled = compile({"/articles/:slug": "a"})
        self.assertEqual(len(compiled.static_routes), 0)
        self.assertEqual(len(compiled.dynamic_routes), 1)
        self.assertEqual(compiled.dynamic_routes[0].param_names, ["slug"])

    def test_classifies_wildcard(self):
        compiled = compile({"*": "fallback"})
        self.assertIsNotNone(compiled.wildcard)
        self.assertEqual(compiled.wildcard.value, "fallback")

    def test_mixed_collection(self):
        compiled = compile(
            {
                "/": "home",
                "/articles/:slug": "dynamic",
                "*": "fallback",
            }
        )
        self.assertEqual(len(compiled.static_routes), 1)
        self.assertEqual(len(compiled.dynamic_routes), 1)
        self.assertEqual(compiled.wildcard.value, "fallback")


if __name__ == "__main__":
    unittest.main()
