# Types Package

The Types package provides comprehensive functions for type checking, validation, format detection, and type conversion operations to help identify and validate data types and formats.

## Basic Type Checking

### `type(value)`
Returns the type of a value as a string.
- **Parameters:** `value` (any) - Value to check
- **Returns:** Type name ("null", "boolean", "number", "string", "list", "map")
- **Example:** `type(42)` → `"number"`

### `isNull(value)`
Checks if a value is null.
- **Parameters:** `value` (any) - Value to check
- **Returns:** Boolean indicating null status
- **Example:** `isNull(null)` → `true`

### `isDefined(value)`
Checks if a value is defined (not null).
- **Parameters:** `value` (any) - Value to check
- **Returns:** Boolean indicating if value is defined
- **Example:** `isDefined("hello")` → `true`

### `isEmpty(value)`
Checks if a value is empty (null, empty string, empty collection, zero, false).
- **Parameters:** `value` (any) - Value to check
- **Returns:** Boolean indicating emptiness
- **Examples:**
  - `isEmpty("")` → `true`
  - `isEmpty([])` → `true`
  - `isEmpty({})` → `true`
  - `isEmpty(0)` → `true`

### `isNotEmpty(value)`
Checks if a value is not empty (opposite of `isEmpty`).
- **Parameters:** `value` (any) - Value to check
- **Returns:** Boolean indicating non-emptiness
- **Example:** `isNotEmpty("hello")` → `true`

## Primitive Type Checking

### `isBool(value)`
Checks if a value is a boolean.
- **Parameters:** `value` (any) - Value to check
- **Returns:** Boolean indicating boolean type
- **Example:** `isBool(true)` → `true`

### `isNumber(value)`
Checks if a value is a number.
- **Parameters:** `value` (any) - Value to check
- **Returns:** Boolean indicating number type
- **Example:** `isNumber(42)` → `true`

### `isString(value)`
Checks if a value is a string.
- **Parameters:** `value` (any) - Value to check
- **Returns:** Boolean indicating string type
- **Example:** `isString("hello")` → `true`

### `isList(value)`
Checks if a value is a list/array.
- **Parameters:** `value` (any) - Value to check
- **Returns:** Boolean indicating list type
- **Example:** `isList([1, 2, 3])` → `true`

### `isMap(value)`
Checks if a value is a map/object.
- **Parameters:** `value` (any) - Value to check
- **Returns:** Boolean indicating map type
- **Example:** `isMap({"key": "value"})` → `true`

### `isArray(value)`
Alias for `isList()` - checks if value is an array.
- **Parameters:** Same as `isList()`
- **Returns:** Boolean indicating array type
- **Example:** `isArray([1, 2, 3])` → `true`

### `isObject(value)`
Alias for `isMap()` - checks if value is an object.
- **Parameters:** Same as `isMap()`
- **Returns:** Boolean indicating object type
- **Example:** `isObject({"key": "value"})` → `true`

## Number Type Checking

### `isInteger(value)`
Checks if a value is an integer (whole number).
- **Parameters:** `value` (any) - Value to check
- **Returns:** Boolean indicating integer type
- **Example:** `isInteger(42)` → `true`, `isInteger(42.5)` → `false`

### `isFloat(value)`
Checks if a value is a floating-point number.
- **Parameters:** `value` (any) - Value to check
- **Returns:** Boolean indicating float type
- **Example:** `isFloat(42.5)` → `true`, `isFloat(42)` → `false`

### `isPositive(value)`
Checks if a number is positive (> 0).
- **Parameters:** `value` (any) - Value to check
- **Returns:** Boolean indicating positive number
- **Example:** `isPositive(42)` → `true`

### `isNegative(value)`
Checks if a number is negative (< 0).
- **Parameters:** `value` (any) - Value to check
- **Returns:** Boolean indicating negative number
- **Example:** `isNegative(-42)` → `true`

### `isZero(value)`
Checks if a number equals zero.
- **Parameters:** `value` (any) - Value to check
- **Returns:** Boolean indicating zero value
- **Example:** `isZero(0)` → `true`

### `isEven(value)`
Checks if an integer is even.
- **Parameters:** `value` (any) - Value to check
- **Returns:** Boolean indicating even integer
- **Example:** `isEven(4)` → `true`

### `isOdd(value)`
Checks if an integer is odd.
- **Parameters:** `value` (any) - Value to check
- **Returns:** Boolean indicating odd integer
- **Example:** `isOdd(5)` → `true`

### `isNan(value)`
Checks if a number is NaN (Not a Number).
- **Parameters:** `value` (any) - Value to check
- **Returns:** Boolean indicating NaN status
- **Example:** `isNan(parseInt("abc"))` → `true`

### `isInfinite(value)`
Checks if a number is infinite.
- **Parameters:** `value` (any) - Value to check
- **Returns:** Boolean indicating infinite status
- **Example:** `isInfinite(1/0)` → `true`

### `isFinite(value)`
Checks if a number is finite (not NaN or infinite).
- **Parameters:** `value` (any) - Value to check
- **Returns:** Boolean indicating finite status
- **Example:** `isFinite(42)` → `true`

## String Content Validation

### `isNumericString(value)`
Checks if a string represents a valid number.
- **Parameters:** `value` (any) - Value to check
- **Returns:** Boolean indicating numeric string
- **Example:** `isNumericString("123.45")` → `true`

### `isAlpha(value)`
Checks if a string contains only alphabetic characters.
- **Parameters:** `value` (any) - Value to check
- **Returns:** Boolean indicating alphabetic content
- **Example:** `isAlpha("hello")` → `true`

### `isAlphanumeric(value)`
Checks if a string contains only alphanumeric characters.
- **Parameters:** `value` (any) - Value to check
- **Returns:** Boolean indicating alphanumeric content
- **Example:** `isAlphanumeric("hello123")` → `true`

### `isDigit(value)`
Checks if a string contains only digit characters.
- **Parameters:** `value` (any) - Value to check
- **Returns:** Boolean indicating digit-only content
- **Example:** `isDigit("12345")` → `true`

### `isLower(value)`
Checks if a string is in lowercase and contains letters.
- **Parameters:** `value` (any) - Value to check
- **Returns:** Boolean indicating lowercase string
- **Example:** `isLower("hello")` → `true`

### `isUpper(value)`
Checks if a string is in uppercase and contains letters.
- **Parameters:** `value` (any) - Value to check
- **Returns:** Boolean indicating uppercase string
- **Example:** `isUpper("HELLO")` → `true`

### `isWhitespace(value)`
Checks if a string contains only whitespace characters.
- **Parameters:** `value` (any) - Value to check
- **Returns:** Boolean indicating whitespace-only content
- **Example:** `isWhitespace("   \t\n")` → `true`

## Format Validation

### `isEmail(value)`
Validates if a string is a properly formatted email address.
- **Parameters:** `value` (any) - Value to validate
- **Returns:** Boolean indicating valid email format
- **Example:** `isEmail("user@example.com")` → `true`

### `isUrl(value)`
Validates if a string is a properly formatted URL.
- **Parameters:** `value` (any) - Value to validate
- **Returns:** Boolean indicating valid URL format
- **Example:** `isUrl("https://example.com")` → `true`

### `isIpAddress(value)`
Validates if a string is a valid IP address (IPv4 or IPv6).
- **Parameters:** `value` (any) - Value to validate
- **Returns:** Boolean indicating valid IP address
- **Examples:**
  - `isIpAddress("192.168.1.1")` → `true`
  - `isIpAddress("2001:db8::1")` → `true`

### `isUUID(value)`
Validates if a string is a properly formatted UUID.
- **Parameters:** `value` (any) - Value to validate
- **Returns:** Boolean indicating valid UUID format
- **Example:** `isUUID("550e8400-e29b-41d4-a716-446655440000")` → `true`

### `isJSON(value)`
Validates if a string is valid JSON.
- **Parameters:** `value` (any) - Value to validate
- **Returns:** Boolean indicating valid JSON format
- **Example:** `isJSON('{"name": "John"}')` → `true`

### `isBase64(value)`
Validates if a string is properly formatted Base64.
- **Parameters:** `value` (any) - Value to validate
- **Returns:** Boolean indicating valid Base64 format
- **Example:** `isBase64("aGVsbG8=")` → `true`

### `isHex(value)`
Validates if a string contains only hexadecimal characters.
- **Parameters:** `value` (any) - Value to validate
- **Returns:** Boolean indicating valid hex format
- **Example:** `isHex("1a2b3c")` → `true`

## Collection Checking

### `hasLength(value)`
Checks if a value has a length property (string, list, or map).
- **Parameters:** `value` (any) - Value to check
- **Returns:** Boolean indicating if value has length
- **Example:** `hasLength("hello")` → `true`

## Range Checking

### `isInRange(value, min, max)`
Checks if a number falls within a specified range (inclusive).
- **Parameters:** 
  - `value` (number) - Number to check
  - `min` (number) - Minimum value
  - `max` (number) - Maximum value
- **Returns:** Boolean indicating if value is in range
- **Example:** `isInRange(5, 1, 10)` → `true`

### `isLengthInRange(value, min, max)`
Checks if the length of a string/list/map falls within a range.
- **Parameters:** 
  - `value` (string|array|object) - Value to check length of
  - `min` (number) - Minimum length
  - `max` (number) - Maximum length
- **Returns:** Boolean indicating if length is in range
- **Example:** `isLengthInRange("hello", 3, 10)` → `true`

## Type Conversion Checking

### `canConvertToNumber(value)`
Checks if a value can be converted to a number.
- **Parameters:** `value` (any) - Value to check
- **Returns:** Boolean indicating convertibility to number
- **Examples:**
  - `canConvertToNumber("123")` → `true`
  - `canConvertToNumber(true)` → `true`
  - `canConvertToNumber("abc")` → `false`

### `canConvertToString(value)`
Checks if a value can be converted to a string.
- **Parameters:** `value` (any) - Value to check
- **Returns:** Boolean indicating convertibility to string
- **Example:** `canConvertToString(123)` → `true`

### `canConvertToBool(value)`
Checks if a value can be converted to a boolean.
- **Parameters:** `value` (any) - Value to check
- **Returns:** Boolean indicating convertibility to boolean
- **Example:** `canConvertToBool("true")` → `true`

## Comparison Functions

### `areEqual(value1, value2)`
Performs deep equality comparison between two values.
- **Parameters:** 
  - `value1` (any) - First value
  - `value2` (any) - Second value
- **Returns:** Boolean indicating deep equality
- **Example:** `areEqual([1, 2, 3], [1, 2, 3])` → `true`

### `areStrictEqual(value1, value2)`
Performs strict equality comparison (same type and value).
- **Parameters:** 
  - `value1` (any) - First value
  - `value2` (any) - Second value
- **Returns:** Boolean indicating strict equality
- **Example:** `areStrictEqual(42, "42")` → `false`

## Usage Notes

### Type Safety
- Functions return `false` for type mismatches rather than errors
- Type checking is strict - numbers don't match strings
- Use conversion checking functions before attempting conversions

### String Validation
- Format validation uses regular expressions for pattern matching
- Email validation uses a simplified but practical regex
- URL validation requires http/https protocols
- IP validation supports both IPv4 and IPv6 formats

### Number Validation
- Integer checking verifies whole numbers (no decimal part)
- Float checking requires decimal part
- NaN and Infinity are handled separately from finite numbers
- Even/odd checking only applies to integers

### Range Validation
- `isInRange()` uses inclusive bounds (min ≤ value ≤ max)
- `isLengthInRange()` works with strings, arrays, and objects
- Negative ranges are validated (min must not exceed max)

### Comparison Semantics
- `areEqual()` performs deep comparison for nested structures
- `areStrictEqual()` requires identical types
- Both handle null values appropriately

### Performance Considerations
- Type checking functions are lightweight and fast
- Format validation with regex may be slower for complex patterns
- Deep equality comparison scales with object complexity

### Common Patterns
```javascript
// Input validation
if (isString(input) && isEmail(input)) {
    // Process valid email
}

// Type-safe operations
if (isNumber(value) && isInRange(value, 0, 100)) {
    // Process percentage
}

// Format validation chain
if (isString(data) && isJSON(data)) {
    const parsed = JSON.parse(data);
    // Process JSON data
}

// Collection validation
if (isList(items) && isLengthInRange(items, 1, 10)) {
    // Process valid list
}
```