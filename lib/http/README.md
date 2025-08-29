# HTTP Package

The HTTP package provides functions for parsing and extracting information from HTTP request contexts, including headers, query parameters, cookies, and request metadata.

## Core Request Functions

### `header(context, headerName)`
Retrieves a specific HTTP header value from the request context.
- **Parameters:** 
  - `context` (map) - The HTTP request context
  - `headerName` (string) - Name of the header to retrieve (case-insensitive)
- **Returns:** Header value as string, or null if not found
- **Example:** `header(ctx, "Content-Type")` → `"application/json"`

### `headers(context)`
Retrieves all HTTP headers from the request context.
- **Parameters:** `context` (map) - The HTTP request context
- **Returns:** Map of all headers (key-value pairs)
- **Example:** `headers(ctx)` → `{"Content-Type": "application/json", "User-Agent": "curl/7.68.0"}`

### `method(context)`
Retrieves the HTTP method from the request context.
- **Parameters:** `context` (map) - The HTTP request context
- **Returns:** HTTP method as string (GET, POST, PUT, etc.)
- **Example:** `method(ctx)` → `"POST"`

### `path(context)`
Retrieves the request path from the request context.
- **Parameters:** `context` (map) - The HTTP request context
- **Returns:** URL path as string
- **Example:** `path(ctx)` → `"/api/users/123"`

### `body(context)`
Retrieves the request body from the request context.
- **Parameters:** `context` (map) - The HTTP request context
- **Returns:** Request body as string
- **Example:** `body(ctx)` → `'{"name": "John", "age": 30}'`

### `status(context)`
Retrieves the HTTP status code from the request context.
- **Parameters:** `context` (map) - The HTTP request context
- **Returns:** Status code as number, or 0 if not available
- **Example:** `status(ctx)` → `200`

## Query Parameter Functions

### `query(context)`
Retrieves all query parameters from the request context.
- **Parameters:** `context` (map) - The HTTP request context
- **Returns:** Map of all query parameters
- **Example:** `query(ctx)` → `{"page": "1", "limit": "10"}`

### `queryParam(context, paramName)`
Retrieves a specific query parameter value.
- **Parameters:** 
  - `context` (map) - The HTTP request context
  - `paramName` (string) - Name of the parameter to retrieve
- **Returns:** Parameter value as string, or null if not found
- **Example:** `queryParam(ctx, "page")` → `"1"`

## Client Information Functions

### `ip(context)`
Retrieves the client IP address, checking various headers and context fields.
- **Parameters:** `context` (map) - The HTTP request context
- **Returns:** IP address as string
- **Checks:** `X-Forwarded-For`, `X-Real-IP` headers and context IP fields
- **Example:** `ip(ctx)` → `"192.168.1.100"`

### `userAgent(context)`
Retrieves the User-Agent header from the request.
- **Parameters:** `context` (map) - The HTTP request context
- **Returns:** User-Agent string
- **Example:** `userAgent(ctx)` → `"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"`

### `host(context)`
Retrieves the Host header from the request.
- **Parameters:** `context` (map) - The HTTP request context
- **Returns:** Host value as string
- **Example:** `host(ctx)` → `"api.example.com"`

### `referer(context)`
Retrieves the Referer header from the request.
- **Parameters:** `context` (map) - The HTTP request context
- **Returns:** Referer URL as string
- **Example:** `referer(ctx)` → `"https://example.com/page"`

## Content Information Functions

### `contentType(context)`
Retrieves the Content-Type header, excluding parameters.
- **Parameters:** `context` (map) - The HTTP request context
- **Returns:** Content type as string (without charset or other parameters)
- **Example:** `contentType(ctx)` → `"application/json"` (from "application/json; charset=utf-8")

### `contentLength(context)`
Retrieves the Content-Length header as a number.
- **Parameters:** `context` (map) - The HTTP request context
- **Returns:** Content length as number, or 0 if not present/invalid
- **Example:** `contentLength(ctx)` → `1024`

### `accept(context)`
Retrieves the Accept header from the request.
- **Parameters:** `context` (map) - The HTTP request context
- **Returns:** Accept header value as string
- **Example:** `accept(ctx)` → `"application/json, text/plain, */*"`

## URL Information Functions

### `scheme(context)`
Retrieves the URL scheme (protocol) from the request.
- **Parameters:** `context` (map) - The HTTP request context
- **Returns:** Scheme as string, defaults to "https" if not found
- **Checks:** Context scheme field and `X-Forwarded-Proto` header
- **Example:** `scheme(ctx)` → `"https"`

### `port(context)`
Retrieves the port number from the request.
- **Parameters:** `context` (map) - The HTTP request context
- **Returns:** Port number, with defaults (443 for HTTPS, 80 for HTTP)
- **Example:** `port(ctx)` → `443`

## Cookie Functions

### `cookies(context)`
Retrieves all cookies from the request.
- **Parameters:** `context` (map) - The HTTP request context
- **Returns:** Map of cookie name-value pairs
- **Example:** `cookies(ctx)` → `{"session_id": "abc123", "theme": "dark"}`

### `cookie(context, cookieName)`
Retrieves a specific cookie value.
- **Parameters:** 
  - `context` (map) - The HTTP request context
  - `cookieName` (string) - Name of the cookie to retrieve
- **Returns:** Cookie value as string, or null if not found
- **Example:** `cookie(ctx, "session_id")` → `"abc123"`

## Authentication Functions

### `authorization(context)`
Retrieves the Authorization header from the request.
- **Parameters:** `context` (map) - The HTTP request context
- **Returns:** Authorization header value as string
- **Example:** `authorization(ctx)` → `"Bearer eyJhbGciOiJIUzI1NiIs..."`

## Usage Notes

### Context Structure
The HTTP context is expected to be a map containing request information with keys such as:
- `headers` - Map of HTTP headers
- `method` - HTTP method string  
- `path` - URL path string
- `query` - Map of query parameters
- `body` - Request body string
- `status` - HTTP status code
- `scheme` - URL scheme
- `port` - Port number
- Various IP-related fields (`remote_ip`, `client_ip`, `x_forwarded_for`, etc.)

### Header Handling
- Header names are case-insensitive when retrieving
- Multiple values in headers like `X-Forwarded-For` are handled by taking the first value
- Missing headers return empty strings or null values

### IP Address Resolution
The `ip()` function checks multiple sources in order:
1. Context fields: `remote_ip`, `client_ip`, `x_forwarded_for`, `x_real_ip`, `ip`
2. HTTP headers: `X-Forwarded-For`, `X-Real-IP`
3. For comma-separated values, returns the first IP address

### Cookie Parsing
Cookies are parsed from the `Cookie` header, handling multiple cookies separated by semicolons.

### Error Handling
Functions return appropriate default values when data is missing:
- Strings return empty string `""`
- Numbers return `0` 
- Maps return empty map `{}`
- Nullable values return `null`