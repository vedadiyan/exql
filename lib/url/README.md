# URL Package

The URL package provides comprehensive functions for URL parsing, manipulation, validation, and construction operations including component extraction, query string handling, and path operations.

## URL Parsing and Analysis

### `parse(url)`
Parses a URL string into its component parts.
- **Parameters:** `url` (string) - URL string to parse
- **Returns:** Map containing URL components
- **Components:** scheme, host, port, path, query, fragment, user
- **Example:** 
  ```javascript
  parse("https://user@example.com:8080/path?param=value#section")
  // → {
  //     "scheme": "https",
  //     "host": "example.com", 
  //     "port": "8080",
  //     "path": "/path",
  //     "query": "param=value",
  //     "fragment": "section",
  //     "user": "user"
  //   }
  ```

### `scheme(url)`
Extracts the scheme/protocol from a URL.
- **Parameters:** `url` (string) - URL string
- **Returns:** Scheme as string
- **Example:** `scheme("https://example.com")` → `"https"`

### `host(url)`
Extracts the hostname from a URL (without port).
- **Parameters:** `url` (string) - URL string
- **Returns:** Hostname as string
- **Example:** `host("https://example.com:8080/path")` → `"example.com"`

### `port(url)`
Extracts the port number from a URL.
- **Parameters:** `url` (string) - URL string
- **Returns:** Port number, with defaults for standard schemes
- **Examples:**
  - `port("https://example.com:8080")` → `8080`
  - `port("https://example.com")` → `443` (HTTPS default)
  - `port("http://example.com")` → `80` (HTTP default)

### `path(url)`
Extracts the path component from a URL.
- **Parameters:** `url` (string) - URL string
- **Returns:** Path as string
- **Example:** `path("https://example.com/api/users")` → `"/api/users"`

### `fragment(url)`
Extracts the fragment (hash) component from a URL.
- **Parameters:** `url` (string) - URL string
- **Returns:** Fragment as string (without #)
- **Example:** `fragment("https://example.com/page#section")` → `"section"`

### `user(url)`
Extracts the username from a URL's user info.
- **Parameters:** `url` (string) - URL string
- **Returns:** Username as string
- **Example:** `user("https://username@example.com")` → `"username"`

## Query Parameter Operations

### `query(url)`
Extracts and parses query parameters from a URL.
- **Parameters:** `url` (string) - URL string
- **Returns:** Map of query parameters
- **Example:** `query("https://example.com?name=John&age=30")` → `{"name": "John", "age": "30"}`

### `query_param(url, paramName)`
Extracts a specific query parameter value.
- **Parameters:** 
  - `url` (string) - URL string
  - `paramName` (string) - Name of parameter to extract
- **Returns:** Parameter value as string, or null if not found
- **Example:** `query_param("https://example.com?name=John&age=30", "name")` → `"John"`

### `query_string(params)`
Converts a map of parameters to a query string.
- **Parameters:** `params` (object) - Map of parameter key-value pairs
- **Returns:** URL-encoded query string
- **Examples:**
  - `query_string({"name": "John", "age": 30})` → `"name=John&age=30"`
  - `query_string({"colors": ["red", "blue"]})` → `"colors=red&colors=blue"`

## URL Encoding

### `encode(string)`
URL-encodes a string for safe use in URLs.
- **Parameters:** `string` (string) - String to encode
- **Returns:** URL-encoded string
- **Example:** `encode("hello world!")` → `"hello%20world%21"`

### `decode(string)`
URL-decodes a string from URL encoding.
- **Parameters:** `string` (string) - URL-encoded string to decode
- **Returns:** Decoded string
- **Example:** `decode("hello%20world%21")` → `"hello world!"`

## URL Construction

### `build(components)`
Constructs a URL from component parts.
- **Parameters:** `components` (object) - Map of URL components
- **Returns:** Complete URL string
- **Components:** scheme, host, port, path, query, fragment, user, password
- **Example:** 
  ```javascript
  build({
    "scheme": "https",
    "host": "example.com",
    "port": "8080", 
    "path": "/api",
    "query": "v=1"
  })
  // → "https://example.com:8080/api?v=1"
  ```

### `join(baseUrl, ...pathSegments)`
Joins a base URL with additional path segments.
- **Parameters:** 
  - `baseUrl` (string) - Base URL
  - `...pathSegments` (string) - Path segments to append
- **Returns:** Complete joined URL
- **Example:** `join("https://api.example.com", "users", "123", "profile")` → `"https://api.example.com/users/123/profile"`

## URL Validation

### `is_absolute(url)`
Checks if a URL is absolute (has scheme).
- **Parameters:** `url` (string) - URL to check
- **Returns:** Boolean indicating if URL is absolute
- **Examples:**
  - `is_absolute("https://example.com")` → `true`
  - `is_absolute("/relative/path")` → `false`

## Path Operations

### `path_segments(urlOrPath)`
Splits a URL path into individual segments.
- **Parameters:** `urlOrPath` (string) - URL or path string
- **Returns:** Array of path segments (URL-decoded)
- **Examples:**
  - `path_segments("https://example.com/api/users/123")` → `["api", "users", "123"]`
  - `path_segments("/api/users/123")` → `["api", "users", "123"]`

### `clean(url)`
Cleans and normalizes a URL path by resolving . and .. segments.
- **Parameters:** `url` (string) - URL to clean
- **Returns:** URL with cleaned path
- **Examples:**
  - `clean("https://example.com/api/../users/./123")` → `"https://example.com/users/123"`
  - `clean("https://example.com/api/users/../admin")` → `"https://example.com/api/admin"`

## Usage Notes

### URL Component Defaults
- **Port defaults:**
  - HTTPS: 443
  - HTTP: 80
  - FTP: 21
  - SSH: 22
- Missing components return empty strings
- Invalid URLs throw parsing errors

### Query Parameter Handling
- Multiple values for the same parameter create arrays
- Empty parameters are preserved
- URL encoding/decoding is handled automatically
- Parameter order follows alphabetical sorting in output

### Path Operations
- Leading and trailing slashes are handled consistently
- Empty path segments are ignored
- URL decoding is applied to path segments
- Relative references are resolved properly

### URL Encoding Rules
- Space becomes `%20`
- Special characters are percent-encoded
- Safe characters (A-Z, a-z, 0-9, -, ., _, ~) are not encoded

### Error Handling
- Invalid URL strings throw descriptive parsing errors
- Missing components return appropriate defaults
- Malformed query strings are handled gracefully

### Common Patterns
```javascript
// Parse and modify URL
const parts = parse("https://api.example.com/users?page=1");
const newUrl = build({
    ...parts,
    "path": "/admin",
    "query": "section=dashboard"
});

// Build API endpoints
const baseUrl = "https://api.example.com";
const userEndpoint = join(baseUrl, "users", userId);
const queryUrl = userEndpoint + "?" + query_string({"include": "profile"});

// Validate and process
if (is_absolute(inputUrl)) {
    const segments = path_segments(inputUrl);
    // Process segments
}

// Clean user input
const cleanUrl = clean(userProvidedUrl);
```