/**
 * Normalize an incoming pathname by removing the URL fragment and query string.
 *
 * All characters after the first "#" or "?" are discarded. No other
 * transformation is applied, so the result remains a raw pathname.
 *
 * @param {unknown} input - The incoming pathname, URL, or any value coercible to string.
 * @returns {string} The normalized pathname without query string or fragment.
 */
export function normalizePathname(input) {
  if (input == null) {
    return "";
  }

  let path = String(input);

  const hashIndex = path.indexOf("#");
  if (hashIndex !== -1) {
    path = path.slice(0, hashIndex);
  }

  const queryIndex = path.indexOf("?");
  if (queryIndex !== -1) {
    path = path.slice(0, queryIndex);
  }

  return path;
}

/**
 * Split a pathname into its non-empty segments.
 *
 * Leading and trailing slashes are ignored, so "/a/b/c/" yields ["a", "b", "c"]
 * and "/" yields [].
 *
 * @param {string} pathname - The normalized pathname to split.
 * @returns {string[]} The list of path segments in order.
 */
export function splitSegments(pathname) {
  return pathname.split("/").filter((segment) => segment !== "");
}
