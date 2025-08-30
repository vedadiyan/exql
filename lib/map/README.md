# Maps Package

The Maps package provides comprehensive functions for map (object) manipulation, transformation, and analysis operations including creation, modification, filtering, merging, and conversion utilities.

## Basic Map Operations

### `keys(map)`
Returns all keys from a map as a sorted list.
- **Parameters:** `map` (object) - The map to extract keys from
- **Returns:** List of key strings in sorted order
- **Example:** `keys({"name": "John", "age": 30})` → `["age", "name"]`

### `values(map)`
Returns all values from a map in key-sorted order.
- **Parameters:** `map` (object) - The map to extract values from
- **Returns:** List of values corresponding to sorted keys
- **Example:** `values({"name": "John", "age": 30})` → `[30, "John"]`

### `size(map)`
Returns the number of key-value pairs in a map.
- **Parameters:** `map` (object) - The map to measure
- **Returns:** Number of entries
- **Example:** `size({"a": 1, "b": 2, "c": 3})` → `3`

### `isEmpty(map)`
Checks if a map contains no key-value pairs.
- **Parameters:** `map` (object) - The map to check
- **Returns:** Boolean indicating if map is empty
- **Example:** `isEmpty({})` → `true`

### `has(map, key)`
Checks if a map contains a specific key.
- **Parameters:** 
  - `map` (object) - The map to search
  - `key` (string) - Key to look for
- **Returns:** Boolean indicating key existence
- **Example:** `has({"name": "John", "age": 30}, "name")` → `true`

### `get(map, key, default?)`
Retrieves a value for a specific key.
- **Parameters:** 
  - `map` (object) - The map to access
  - `key` (string) - Key to retrieve
  - `default` (any, optional) - Value to return if key doesn't exist
- **Returns:** Value at key or default value
- **Example:** `get({"name": "John"}, "age", 25)` → `25`

### `set(map, key, value)`
Creates a new map with a key-value pair added or updated.
- **Parameters:** 
  - `map` (object) - The original map
  - `key` (string) - Key to set
  - `value` (any) - Value to assign
- **Returns:** New map with updated key-value pair
- **Example:** `set({"name": "John"}, "age", 30)` → `{"name": "John", "age": 30}`

### `delete(map, key)`
Creates a new map with a specific key removed.
- **Parameters:** 
  - `map` (object) - The original map
  - `key` (string) - Key to remove
- **Returns:** New map without the specified key
- **Example:** `delete({"name": "John", "age": 30}, "age")` → `{"name": "John"}`

## Map Merging

### `merge(...maps)`
Combines multiple maps into one (shallow merge).
- **Parameters:** `...maps` (object) - Maps to merge
- **Returns:** New merged map (later maps override earlier ones)
- **Example:** `merge({"a": 1, "b": 2}, {"b": 3, "c": 4})` → `{"a": 1, "b": 3, "c": 4}`

### `mergeDeep(...maps)`
Performs deep merge of multiple maps, recursively merging nested objects.
- **Parameters:** `...maps` (object) - Maps to merge deeply
- **Returns:** New deeply merged map
- **Example:** 
  ```javascript
  mergeDeep(
    {"user": {"name": "John", "settings": {"theme": "dark"}}}, 
    {"user": {"age": 30, "settings": {"lang": "en"}}}
  )
  // → {"user": {"name": "John", "age": 30, "settings": {"theme": "dark", "lang": "en"}}}
  ```

## Map Transformation

### `invert(map)`
Creates a new map with keys and values swapped.
- **Parameters:** `map` (object) - The map to invert
- **Returns:** New map with values as keys and keys as values
- **Example:** `invert({"a": "x", "b": "y"})` → `{"x": "a", "y": "b"}`

### `filter(map)`
Creates a new map with null/falsy values removed.
- **Parameters:** `map` (object) - The map to filter
- **Returns:** New map without null/falsy values
- **Example:** `filter({"a": 1, "b": null, "c": "", "d": 2})` → `{"a": 1, "d": 2}`

### `filterKeys(map, keys)`
Creates a new map containing only specified keys.
- **Parameters:** 
  - `map` (object) - The original map
  - `keys` (array) - List of keys to keep
- **Returns:** New map with only the specified keys
- **Example:** `filterKeys({"a": 1, "b": 2, "c": 3}, ["a", "c"])` → `{"a": 1, "c": 3}`

### `omitKeys(map, keys)`
Creates a new map excluding specified keys.
- **Parameters:** 
  - `map` (object) - The original map
  - `keys` (array) - List of keys to exclude
- **Returns:** New map without the specified keys
- **Example:** `omitKeys({"a": 1, "b": 2, "c": 3}, ["b"])` → `{"a": 1, "c": 3}`

### `rename(map, renameMap)`
Creates a new map with keys renamed according to a mapping.
- **Parameters:** 
  - `map` (object) - The original map
  - `renameMap` (object) - Map of old key to new key mappings
- **Returns:** New map with renamed keys
- **Example:** `rename({"firstName": "John", "age": 30}, {"firstName": "name"})` → `{"name": "John", "age": 30}`

## Map Conversion

### `toList(map)`
Converts a map to a list of key-value pairs.
- **Parameters:** `map` (object) - The map to convert
- **Returns:** List of [key, value] pairs in sorted key order
- **Example:** `toList({"b": 2, "a": 1})` → `[["a", 1], ["b", 2]]`

### `fromList(list)`
Creates a map from a list of key-value pairs.
- **Parameters:** `list` (array) - List of [key, value] pairs
- **Returns:** New map created from pairs
- **Example:** `fromList([["a", 1], ["b", 2]])` → `{"a": 1, "b": 2}`

### `toQueryString(map)`
Converts a map to a URL query string.
- **Parameters:** `map` (object) - The map to convert
- **Returns:** URL-encoded query string
- **Example:** `toQueryString({"name": "John Doe", "age": 30})` → `"age=30&name=John Doe"`

### `fromQueryString(queryString)`
Parses a URL query string into a map.
- **Parameters:** `queryString` (string) - Query string to parse
- **Returns:** Map of parameter key-value pairs
- **Note:** Duplicate parameters create arrays
- **Examples:**
  - `fromQueryString("name=John&age=30")` → `{"name": "John", "age": "30"}`
  - `fromQueryString("color=red&color=blue")` → `{"color": ["red", "blue"]}`

## Path Operations (Dot Notation)

### `getPath(map, path, default?)`
Retrieves a value using dot notation path.
- **Parameters:** 
  - `map` (object) - The map to access
  - `path` (string) - Dot notation path (e.g., "user.profile.name")
  - `default` (any, optional) - Value to return if path doesn't exist
- **Returns:** Value at path or default value
- **Example:** `getPath({"user": {"profile": {"name": "John"}}}, "user.profile.name")` → `"John"`

### `setPath(map, path, value)`
Sets a value using dot notation path, creating intermediate objects as needed.
- **Parameters:** 
  - `map` (object) - The original map
  - `path` (string) - Dot notation path where to set value
  - `value` (any) - Value to set
- **Returns:** New map with value set at path
- **Example:** `setPath({}, "user.profile.name", "John")` → `{"user": {"profile": {"name": "John"}}}`

### `hasPath(map, path)`
Checks if a dot notation path exists in the map.
- **Parameters:** 
  - `map` (object) - The map to check
  - `path` (string) - Dot notation path to verify
- **Returns:** Boolean indicating path existence
- **Example:** `hasPath({"user": {"name": "John"}}, "user.name")` → `true`

### `deletePath(map, path)`
Removes a value at a dot notation path.
- **Parameters:** 
  - `map` (object) - The original map
  - `path` (string) - Dot notation path to remove
- **Returns:** New map with path removed
- **Example:** `deletePath({"user": {"name": "John", "age": 30}}, "user.age")` → `{"user": {"name": "John"}}`

## Usage Notes

### Immutability
All functions return new maps rather than modifying originals:
- Original maps remain unchanged
- Safe to chain operations: `filterKeys(set(map, "new", "value"), ["key1", "key2"])`

### Key Ordering
- `keys()` and `values()` return results in sorted order for consistency
- `toList()` preserves sorted key order
- Original insertion order is not maintained

### Path Operations
Dot notation paths support nested object access:
- `"user.profile.settings.theme"` accesses deeply nested values
- `setPath()` automatically creates intermediate objects
- Empty paths (`""`) refer to the root object

### Deep vs Shallow Operations
- `merge()` performs shallow merge (only top-level keys)
- `mergeDeep()` recursively merges nested objects
- Most other operations are shallow by default

### Type Handling
- All value types are supported (strings, numbers, booleans, objects, arrays, null)
- Key names are always converted to strings
- `null` and `undefined` values are handled appropriately

### Query String Handling
- Parameters with multiple values become arrays
- URL encoding/decoding is handled automatically
- Empty parameters are preserved
- Parameter order in output follows sorted key order

### Performance Considerations
- Operations create copies, suitable for functional programming patterns
- Deep operations (`mergeDeep`, complex paths) may be slower on large nested structures
- `filterKeys()` and `omitKeys()` are efficient for key-based filtering

### Common Patterns
```javascript
// Building complex objects
let user = {};
user = setPath(user, "profile.personal.name", "John");
user = setPath(user, "profile.personal.email", "john@example.com");
user = setPath(user, "settings.theme", "dark");

// Safe property access
const theme = getPath(user, "settings.theme", "light");

// Filtering and transformation
const publicData = omitKeys(user, ["password", "ssn"]);
const renamedData = rename(publicData, {"email": "emailAddress"});
```