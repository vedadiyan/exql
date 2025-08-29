# JSON Package

The JSON package provides comprehensive functions for JSON parsing, manipulation, validation, and transformation operations.

## Core JSON Operations

### `parse(jsonString)`
Parses a JSON string into a structured value.
- **Parameters:** `jsonString` (string) - Valid JSON string to parse
- **Returns:** Parsed JSON data (object, array, string, number, boolean, or null)
- **Example:** `parse('{"name": "John", "age": 30}')` → `{"name": "John", "age": 30}`

### `string(data, pretty?)`
Converts data to a JSON string representation.
- **Parameters:** 
  - `data` (any) - Data to convert to JSON
  - `pretty` (boolean, optional) - If true, formats with indentation
- **Returns:** JSON string representation
- **Examples:**
  - `string({"name": "John"})` → `'{"name":"John"}'`
  - `string({"name": "John"}, true)` → `'{\n  "name": "John"\n}'`

### `valid(jsonString)`
Checks if a string is valid JSON.
- **Parameters:** `jsonString` (string) - String to validate
- **Returns:** Boolean indicating JSON validity
- **Example:** `valid('{"name": "John"}')` → `true`

## JSON Data Access

### `get(data, path)`
Retrieves a value from JSON data using dot notation path.
- **Parameters:** 
  - `data` (object|string) - JSON data or JSON string
  - `path` (string) - Dot notation path (e.g., "user.profile.name")
- **Returns:** Value at the specified path, or null if not found
- **Examples:**
  - `get({"user": {"name": "John"}}, "user.name")` → `"John"`
  - `get({"items": [1, 2, 3]}, "items.0")` → `1`

### `set(data, path, value)`
Sets a value in JSON data using dot notation path.
- **Parameters:** 
  - `data` (object|string) - JSON data or JSON string
  - `path` (string) - Dot notation path where to set the value
  - `value` (any) - Value to set
- **Returns:** Modified JSON data
- **Example:** `set({"user": {}}, "user.name", "John")` → `{"user": {"name": "John"}}`

### `delete(data, path)`
Removes a value from JSON data using dot notation path.
- **Parameters:** 
  - `data` (object|string) - JSON data or JSON string
  - `path` (string) - Dot notation path of value to remove
- **Returns:** Modified JSON data with value removed
- **Example:** `delete({"user": {"name": "John", "age": 30}}, "user.age")` → `{"user": {"name": "John"}}`

### `has(data, path)`
Checks if a path exists in JSON data.
- **Parameters:** 
  - `data` (object|string) - JSON data or JSON string
  - `path` (string) - Dot notation path to check
- **Returns:** Boolean indicating if path exists
- **Example:** `has({"user": {"name": "John"}}, "user.name")` → `true`

## JSON Structure Analysis

### `keys(data)`
Retrieves all keys from a JSON object.
- **Parameters:** `data` (object|string) - JSON object or JSON string
- **Returns:** Array of key strings
- **Example:** `keys({"name": "John", "age": 30})` → `["name", "age"]`

### `values(data)`
Retrieves all values from a JSON object or array.
- **Parameters:** `data` (object|array|string) - JSON data or JSON string
- **Returns:** Array of values
- **Examples:**
  - `values({"name": "John", "age": 30})` → `["John", 30]`
  - `values([1, 2, 3])` → `[1, 2, 3]`

### `length(data)`
Gets the length/size of JSON data.
- **Parameters:** `data` (object|array|string) - JSON data or JSON string
- **Returns:** Number representing length/size
- **Examples:**
  - `length({"a": 1, "b": 2})` → `2`
  - `length([1, 2, 3])` → `3`
  - `length("hello")` → `5`

### `type(data)`
Determines the JSON type of data.
- **Parameters:** `data` (any|string) - JSON data or JSON string
- **Returns:** String indicating type ("null", "boolean", "number", "string", "array", "object")
- **Examples:**
  - `type(42)` → `"number"`
  - `type({"a": 1})` → `"object"`
  - `type([1, 2, 3])` → `"array"`

## JSON Manipulation

### `merge(object1, object2, ...)`
Merges multiple JSON objects into one (shallow merge).
- **Parameters:** Multiple JSON objects or JSON strings
- **Returns:** Merged JSON object
- **Note:** Later objects override earlier ones for duplicate keys
- **Example:** `merge({"a": 1}, {"b": 2}, {"a": 3})` → `{"a": 3, "b": 2}`

## Path Operations

JSON paths use dot notation with support for:
- **Object properties:** `"user.profile.name"`
- **Array indices:** `"items.0"` (first element), `"items.1"` (second element)
- **Nested paths:** `"users.0.profile.settings.theme"`

### Path Examples
```json
{
  "users": [
    {
      "name": "John",
      "profile": {
        "email": "john@example.com",
        "settings": {
          "theme": "dark"
        }
      }
    }
  ]
}
```

- `get(data, "users.0.name")` → `"John"`
- `get(data, "users.0.profile.email")` → `"john@example.com"`  
- `get(data, "users.0.profile.settings.theme")` → `"dark"`

## Usage Notes

### Data Input Flexibility
Most functions accept either:
- Parsed JSON data (objects, arrays, primitives)
- JSON strings that will be automatically parsed

### Path Creation
- The `set()` function automatically creates intermediate objects/arrays as needed
- Numeric path segments create arrays, string segments create objects
- Example: `set({}, "items.0.name", "test")` creates `{"items": [{"name": "test"}]}`

### Type Handling
- All JSON types are supported: null, boolean, number, string, array, object
- Invalid JSON strings return appropriate errors
- Path operations return `null` for non-existent paths

### Error Handling
- Invalid JSON strings throw parsing errors
- Invalid paths in get operations return `null`
- Type mismatches (e.g., accessing array index on object) return `null`
- Functions validate input types and provide meaningful error messages

### Performance Notes
- String inputs are parsed each time - consider parsing once for multiple operations
- Deep paths are traversed efficiently
- Large JSON structures are handled appropriately

### Common Patterns
```javascript
// Parse once, use many times
const data = parse(jsonString);
const name = get(data, "user.name");
const email = get(data, "user.email");

// Chain operations
let user = {};
user = set(user, "profile.name", "John");
user = set(user, "profile.email", "john@example.com");
user = set(user, "settings.theme", "dark");

// Validation before access
if (has(data, "user.profile.settings")) {
    const settings = get(data, "user.profile.settings");
    // Use settings safely
}
```