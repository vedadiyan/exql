# Util Package

The Util package provides general-purpose utility functions including conditional logic, null handling, comparison operations, debugging tools, type conversion, validation, and functional programming utilities.

## Conditional Functions

### `if(condition, trueValue, falseValue?)`
Returns a value based on a condition (ternary operator).
- **Parameters:** 
  - `condition` (boolean) - Condition to evaluate
  - `trueValue` (any) - Value to return if condition is true
  - `falseValue` (any, optional) - Value to return if condition is false
- **Returns:** trueValue if condition is true, falseValue (or null) otherwise
- **Example:** `if(age >= 18, "adult", "minor")` → `"adult"` or `"minor"`

### `unless(condition, falseValue, trueValue?)`
Opposite of `if` - returns a value when condition is false.
- **Parameters:** 
  - `condition` (boolean) - Condition to evaluate
  - `falseValue` (any) - Value to return if condition is false
  - `trueValue` (any, optional) - Value to return if condition is true
- **Returns:** falseValue if condition is false, trueValue (or null) otherwise
- **Example:** `unless(isEmpty(name), name, "Anonymous")` → name if not empty, "Anonymous" if empty

### `switch(value, case1, result1, case2, result2, ..., default)`
Multi-way conditional based on value matching.
- **Parameters:** 
  - `value` (any) - Value to test
  - `case1, result1, ...` (any) - Case-result pairs
  - `default` (any) - Default value if no cases match
- **Returns:** Result for matching case, or default value
- **Example:** `switch(grade, "A", "Excellent", "B", "Good", "C", "Average", "Unknown")` → grade-based result

## Null Coalescing Functions

### `coalesce(...values)`
Returns the first non-null value from a list.
- **Parameters:** `...values` (any) - Values to check
- **Returns:** First non-null value, or null if all are null
- **Example:** `coalesce(null, "", "default")` → `""` (first non-null)

### `default(value, defaultValue)`
Returns a default value if the first value is null.
- **Parameters:** 
  - `value` (any) - Value to check
  - `defaultValue` (any) - Default to use if value is null
- **Returns:** value if not null, defaultValue otherwise
- **Example:** `default(userInput, "No input provided")` → userInput or default message

### `firstNonNull(...values)`
Alias for `coalesce()` - returns first non-null value.
- **Parameters:** Same as `coalesce()`
- **Returns:** First non-null value
- **Example:** `firstNonNull(null, undefined, "found")` → `"found"`

### `firstNonEmpty(...values)`
Returns the first non-empty value (not null and not empty string/collection).
- **Parameters:** `...values` (any) - Values to check
- **Returns:** First non-empty value, or null if all are empty
- **Example:** `firstNonEmpty("", null, [], "hello")` → `"hello"`

## Comparison and Selection

### `greatest(...values)`
Returns the largest value from the arguments.
- **Parameters:** `...values` (number) - Values to compare
- **Returns:** Maximum value, or null if no arguments
- **Example:** `greatest(1, 5, 3, 9, 2)` → `9`

### `least(...values)`
Returns the smallest value from the arguments.
- **Parameters:** `...values` (number) - Values to compare
- **Returns:** Minimum value, or null if no arguments
- **Example:** `least(1, 5, 3, 9, 2)` → `1`

### `choose(index, ...options)`
Selects a value by 1-based index from options.
- **Parameters:** 
  - `index` (number) - 1-based index to select
  - `...options` (any) - Options to choose from
- **Returns:** Value at specified index
- **Example:** `choose(2, "first", "second", "third")` → `"second"`

## Debugging and Inspection

### `debug(...values)`
Prints debug information to console and returns first value.
- **Parameters:** `...values` (any) - Values to debug
- **Returns:** First argument (for chaining)
- **Side Effect:** Prints "[DEBUG] value1 value2 ..." to console
- **Example:** `debug("Processing", user, "step 1")` → prints debug info, returns "Processing"

### `inspect(value)`
Returns a detailed string description of a value's type and content.
- **Parameters:** `value` (any) - Value to inspect
- **Returns:** Detailed description string
- **Example:** `inspect([1, 2, 3])` → `"list: length 3, elements: [1, 2, 3]"`

### `dump(...values)`
Returns detailed string descriptions of multiple values, separated by newlines.
- **Parameters:** `...values` (any) - Values to dump
- **Returns:** Multi-line description string
- **Example:** `dump("hello", 42, [1, 2])` → detailed descriptions separated by newlines

## Identity and Pass-through

### `identity(value)`
Returns the input value unchanged (identity function).
- **Parameters:** `value` (any) - Value to return
- **Returns:** The same value
- **Example:** `identity(42)` → `42`

### `noop(...args)`
No-operation function that ignores all arguments.
- **Parameters:** `...args` (any) - Arguments (ignored)
- **Returns:** Always returns null
- **Example:** `noop("anything")` → `null`

### `constant(value)`
Returns the input value (alias for `identity`).
- **Parameters:** `value` (any) - Value to return
- **Returns:** The same value
- **Example:** `constant("hello")` → `"hello"`

## Error Handling

### `tryOr(value, fallback)`
Returns value if not null, otherwise returns fallback.
- **Parameters:** 
  - `value` (any) - Value to try
  - `fallback` (any) - Fallback value
- **Returns:** value if not null, fallback otherwise
- **Example:** `tryOr(riskyOperation(), "safe default")` → value or fallback

### `safe(value)`
Returns value if not null, otherwise returns null (null-safe wrapper).
- **Parameters:** `value` (any) - Value to make safe
- **Returns:** value if not null, null otherwise
- **Example:** `safe(potentiallyNull)` → value or null

## Type Conversion Utilities

### `tostring(value)`
Converts a value to its string representation.
- **Parameters:** `value` (any) - Value to convert
- **Returns:** String representation
- **Example:** `tostring(123)` → `"123"`

### `tonumber(value)`
Converts a value to a number.
- **Parameters:** `value` (any) - Value to convert
- **Returns:** Numeric representation
- **Example:** `tonumber("123.45")` → `123.45`

### `tobool(value)`
Converts a value to a boolean.
- **Parameters:** `value` (any) - Value to convert
- **Returns:** Boolean representation
- **Example:** `tobool("true")` → `true`

### `tolist(...values)`
Converts arguments to a list.
- **Parameters:** `...values` (any) - Values to convert to list
- **Returns:** List containing the values
- **Examples:**
  - `tolist()` → `[]` (empty list)
  - `tolist([1, 2, 3])` → `[1, 2, 3]` (unchanged if already list)
  - `tolist(1, 2, 3)` → `[1, 2, 3]` (multiple args to list)

## Validation Utilities

### `assert(condition, message?)`
Throws an error if condition is false.
- **Parameters:** 
  - `condition` (boolean) - Condition to assert
  - `message` (string, optional) - Custom error message
- **Returns:** true if assertion passes
- **Throws:** Error if condition is false
- **Example:** `assert(age >= 0, "Age cannot be negative")` → true or throws error

### `validate(value, condition, fallback?)`
Returns value if condition is true, otherwise returns fallback.
- **Parameters:** 
  - `value` (any) - Value to validate
  - `condition` (boolean) - Validation condition
  - `fallback` (any, optional) - Fallback value
- **Returns:** value if condition is true, fallback otherwise
- **Example:** `validate(email, isValidEmail(email), null)` → email if valid, null if invalid

### `require(value, message?)`
Throws an error if value is null or undefined.
- **Parameters:** 
  - `value` (any) - Required value
  - `message` (string, optional) - Custom error message
- **Returns:** The value if not null
- **Throws:** Error if value is null
- **Example:** `require(userId, "User ID is required")` → userId or throws error

## Functional Utilities

### `apply(...args)`
Identity function for function application patterns.
- **Parameters:** `...args` (any) - Arguments to apply
- **Returns:** First argument
- **Example:** `apply(value)` → `value`

### `pipe(...args)`
Identity function for piping patterns.
- **Parameters:** `...args` (any) - Arguments to pipe
- **Returns:** First argument
- **Example:** `pipe(value)` → `value`

### `compose(...args)`
Identity function for composition patterns.
- **Parameters:** `...args` (any) - Arguments to compose
- **Returns:** First argument
- **Example:** `compose(value)` → `value`

## Miscellaneous Utilities

### `uuid()`
Generates a new UUID (Universally Unique Identifier).
- **Parameters:** None
- **Returns:** UUID string
- **Example:** `uuid()` → `"550e8400-e29b-41d4-a716-446655440000"`

### `timestamp()`
Returns the current Unix timestamp in seconds.
- **Parameters:** None
- **Returns:** Current timestamp as number
- **Example:** `timestamp()` → `1693526400`

### `randomString(length?)`
Generates a random alphanumeric string.
- **Parameters:** `length` (number, optional) - String length (default: 10, max: 1000)
- **Returns:** Random string
- **Example:** `randomString(8)` → `"A7xK9mPq"`

### `memoize(value)`
Placeholder for memoization (currently returns input unchanged).
- **Parameters:** `value` (any) - Value to memoize
- **Returns:** The input value
- **Example:** `memoize(expensiveFunction)` → returns input

### `benchmark(value)`
Placeholder for benchmarking (returns mock timing).
- **Parameters:** `value` (any) - Value to benchmark
- **Returns:** Mock timing (0.001)
- **Example:** `benchmark(operation)` → `0.001`

## Usage Notes

### Null vs Empty
- **Null:** `null` or `undefined` values
- **Empty:** Empty strings (`""`), empty arrays (`[]`), empty objects (`{}`), zero (`0`), or `false`
- `coalesce()` checks for null, `firstNonEmpty()` checks for both null and empty

### Conditional Logic Patterns
```javascript
// Traditional ternary
const result = if(isValid, processData(), getDefault());

// Multi-case selection
const grade = switch(score, 
    90, "A", 
    80, "B", 
    70, "C", 
    "F"
);

// Null coalescing chain
const name = coalesce(user.nickname, user.firstName, user.email, "Anonymous");
```

### Debugging Workflow
```javascript
// Chain debugging without breaking flow
const result = debug("Starting processing", inputData)
    |> process()
    |> debug("After processing")
    |> validate();

// Detailed inspection
const info = inspect(complexObject);
console.log(info); // "map: 5 keys, content: {...}"
```

### Validation Patterns
```javascript
// Assert preconditions
assert(users.length > 0, "No users found");
assert(isValidEmail(email), "Invalid email format");

// Validate with fallbacks
const safeAge = validate(inputAge, inputAge >= 0 && inputAge <= 150, 0);
const requiredName = require(userName, "Name is required for this operation");
```

### Type Conversion Safety
```javascript
// Safe conversions with error handling
try {
    const num = tonumber(userInput);
    const str = tostring(dataObject);
    const bool = tobool(flagValue);
} catch (error) {
    // Handle conversion errors
}
```

### Error Handling
- Functions throw descriptive errors for invalid arguments
- Validation functions provide clear error messages
- Type conversion functions validate input compatibility