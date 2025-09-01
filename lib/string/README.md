# String Package

The String package provides comprehensive functions for string manipulation, formatting, validation, and transformation operations including length calculations, case conversion, trimming, searching, and pattern matching.

## Length and Size Functions

### `len(string)`
Returns the number of Unicode characters in a string.
- **Parameters:** `string` (string) - The string to measure
- **Returns:** Number of Unicode characters (not bytes)
- **Example:** `len("hello")` → `5`

### `size(string)`
Returns the number of bytes in a string.
- **Parameters:** `string` (string) - The string to measure
- **Returns:** Number of bytes
- **Example:** `size("hello")` → `5`

## String Building and Manipulation

### `concat(...strings)`
Concatenates multiple strings together.
- **Parameters:** `...strings` (any) - Values to concatenate (converted to strings)
- **Returns:** Combined string
- **Example:** `concat("Hello", " ", "World", "!")` → `"Hello World!"`

### `repeat(string, count)`
Repeats a string a specified number of times.
- **Parameters:** 
  - `string` (string) - String to repeat
  - `count` (number) - Number of repetitions
- **Returns:** Repeated string
- **Example:** `repeat("ha", 3)` → `"hahaha"`

### `reverse(string)`
Reverses the order of characters in a string.
- **Parameters:** `string` (string) - String to reverse
- **Returns:** Reversed string
- **Example:** `reverse("hello")` → `"olleh"`

## Case Conversion

### `upper(string)`
Converts a string to uppercase.
- **Parameters:** `string` (string) - String to convert
- **Returns:** Uppercase string
- **Example:** `upper("hello")` → `"HELLO"`

### `lower(string)`
Converts a string to lowercase.
- **Parameters:** `string` (string) - String to convert
- **Returns:** Lowercase string
- **Example:** `lower("HELLO")` → `"hello"`

### `title(string)`
Converts a string to title case (first letter of each word capitalized).
- **Parameters:** `string` (string) - String to convert
- **Returns:** Title case string
- **Example:** `title("hello world")` → `"Hello World"`

### `capitalize(string)`
Capitalizes the first letter and lowercases the rest.
- **Parameters:** `string` (string) - String to capitalize
- **Returns:** Capitalized string
- **Example:** `capitalize("hELLO")` → `"Hello"`

### `swapCase(string)`
Swaps the case of each character.
- **Parameters:** `string` (string) - String to swap case
- **Returns:** Case-swapped string
- **Example:** `swapCase("HeLLo")` → `"hEllO"`

## Trimming and Padding

### `trim(string, cutset?)`
Removes whitespace or specified characters from both ends.
- **Parameters:** 
  - `string` (string) - String to trim
  - `cutset` (string, optional) - Characters to remove (default: whitespace)
- **Returns:** Trimmed string
- **Examples:**
  - `trim("  hello  ")` → `"hello"`
  - `trim("...hello...", ".")` → `"hello"`

### `trimLeft(string, cutset?)`
Removes whitespace or specified characters from the left end.
- **Parameters:** Same as `trim()`
- **Returns:** Left-trimmed string
- **Example:** `trimLeft("  hello  ")` → `"hello  "`

### `trimRight(string, cutset?)`
Removes whitespace or specified characters from the right end.
- **Parameters:** Same as `trim()`
- **Returns:** Right-trimmed string
- **Example:** `trimRight("  hello  ")` → `"  hello"`

### `padLeft(string, totalLength, padChar?)`
Pads a string on the left to reach a target length.
- **Parameters:** 
  - `string` (string) - String to pad
  - `totalLength` (number) - Target length
  - `padChar` (string, optional) - Padding character (default: space)
- **Returns:** Left-padded string
- **Example:** `padLeft("hello", 10, "0")` → `"00000hello"`

### `padRight(string, totalLength, padChar?)`
Pads a string on the right to reach a target length.
- **Parameters:** Same as `padLeft()`
- **Returns:** Right-padded string
- **Example:** `padRight("hello", 10, ".")` → `"hello....."`

### `padCenter(string, totalLength, padChar?)`
Centers a string by padding both sides to reach a target length.
- **Parameters:** Same as `padLeft()`
- **Returns:** Center-padded string
- **Example:** `padCenter("hello", 11, "-")` → `"---hello---"`

## String Extraction

### `substr(string, start, length?)`
Extracts a substring starting at a position.
- **Parameters:** 
  - `string` (string) - Source string
  - `start` (number) - Starting position (supports negative indexing)
  - `length` (number, optional) - Number of characters to extract
- **Returns:** Substring
- **Examples:**
  - `substr("hello world", 6)` → `"world"`
  - `substr("hello world", 0, 5)` → `"hello"`
  - `substr("hello world", -5)` → `"world"`

### `left(string, count)`
Gets the leftmost N characters.
- **Parameters:** 
  - `string` (string) - Source string
  - `count` (number) - Number of characters
- **Returns:** Leftmost characters
- **Example:** `left("hello world", 5)` → `"hello"`

### `right(string, count)`
Gets the rightmost N characters.
- **Parameters:** 
  - `string` (string) - Source string
  - `count` (number) - Number of characters
- **Returns:** Rightmost characters
- **Example:** `right("hello world", 5)` → `"world"`

## String Search

### `contains(string, substring)`
Checks if a string contains a substring.
- **Parameters:** 
  - `string` (string) - String to search in
  - `substring` (string) - Substring to find
- **Returns:** Boolean indicating presence
- **Example:** `contains("hello world", "world")` → `true`

### `startswith(string, prefix)`
Checks if a string starts with a prefix.
- **Parameters:** 
  - `string` (string) - String to check
  - `prefix` (string) - Prefix to look for
- **Returns:** Boolean indicating match
- **Example:** `startswith("hello world", "hello")` → `true`

### `endswith(string, suffix)`
Checks if a string ends with a suffix.
- **Parameters:** 
  - `string` (string) - String to check
  - `suffix` (string) - Suffix to look for
- **Returns:** Boolean indicating match
- **Example:** `endswith("hello world", "world")` → `true`

### `indexof(string, substring, start?)`
Finds the first index of a substring.
- **Parameters:** 
  - `string` (string) - String to search in
  - `substring` (string) - Substring to find
  - `start` (number, optional) - Starting search position
- **Returns:** Index of substring or -1 if not found
- **Example:** `indexof("hello world", "o")` → `4`

### `lastIndexOf(string, substring)`
Finds the last index of a substring.
- **Parameters:** 
  - `string` (string) - String to search in
  - `substring` (string) - Substring to find
- **Returns:** Last index of substring or -1 if not found
- **Example:** `lastIndexOf("hello world", "o")` → `7`

## String Replacement

### `replace(string, old, new, count?)`
Replaces occurrences of a substring.
- **Parameters:** 
  - `string` (string) - Source string
  - `old` (string) - Substring to replace
  - `new` (string) - Replacement string
  - `count` (number, optional) - Maximum replacements (-1 for all)
- **Returns:** String with replacements
- **Example:** `replace("hello world", "l", "x", 2)` → `"hexxo world"`

### `replaceAll(string, old, new)`
Replaces all occurrences of a substring.
- **Parameters:** 
  - `string` (string) - Source string
  - `old` (string) - Substring to replace
  - `new` (string) - Replacement string
- **Returns:** String with all replacements
- **Example:** `replaceAll("hello world", "l", "x")` → `"hexxo worxd"`

## String Splitting and Joining

### `split(string, separator, count?)`
Splits a string into an array by a separator.
- **Parameters:** 
  - `string` (string) - String to split
  - `separator` (string) - Separator to split on
  - `count` (number, optional) - Maximum splits (-1 for all)
- **Returns:** Array of string parts
- **Example:** `split("a,b,c,d", ",")` → `["a", "b", "c", "d"]`

### `join(separator, list)`
Joins an array of strings with a separator.
- **Parameters:** 
  - `separator` (string) - String to join with
  - `list` (array) - Array of strings to join
- **Returns:** Joined string
- **Example:** `join("-", ["a", "b", "c"])` → `"a-b-c"`

### `lines(string)`
Splits a string into lines.
- **Parameters:** `string` (string) - String to split into lines
- **Returns:** Array of lines
- **Example:** `lines("line1\nline2\nline3")` → `["line1", "line2", "line3"]`

### `fields(string)`
Splits a string by whitespace into words.
- **Parameters:** `string` (string) - String to split into words
- **Returns:** Array of words
- **Example:** `fields("hello   world  test")` → `["hello", "world", "test"]`

## Regular Expressions

### `match(string, pattern)`
Tests if a string matches a regular expression pattern.
- **Parameters:** 
  - `string` (string) - String to test
  - `pattern` (string) - Regular expression pattern
- **Returns:** Boolean indicating match
- **Example:** `match("hello123", "\\d+")` → `true`

### `findAll(string, pattern, count?)`
Finds all matches of a regular expression pattern.
- **Parameters:** 
  - `string` (string) - String to search
  - `pattern` (string) - Regular expression pattern
  - `count` (number, optional) - Maximum matches (-1 for all)
- **Returns:** Array of matched strings
- **Example:** `findAll("hello123world456", "\\d+")` → `["123", "456"]`

### `replaceRegex(string, pattern, replacement)`
Replaces text matching a regular expression pattern.
- **Parameters:** 
  - `string` (string) - Source string
  - `pattern` (string) - Regular expression pattern
  - `replacement` (string) - Replacement string
- **Returns:** String with regex replacements
- **Example:** `replaceRegex("hello123world", "\\d+", "XXX")` → `"helloXXXworld"`

## Character Operations

### `charAt(string, index)`
Gets the character at a specific index.
- **Parameters:** 
  - `string` (string) - Source string
  - `index` (number) - Character index
- **Returns:** Character at index or empty string if out of bounds
- **Example:** `charAt("hello", 1)` → `"e"`

### `charCode(string, index)`
Gets the Unicode code point of a character at an index.
- **Parameters:** 
  - `string` (string) - Source string
  - `index` (number) - Character index
- **Returns:** Unicode code point or 0 if out of bounds
- **Example:** `charCode("hello", 0)` → `104`

### `fromCharCode(...codes)`
Creates a string from Unicode code points.
- **Parameters:** `...codes` (number) - Unicode code points
- **Returns:** String created from code points
- **Example:** `fromCharCode(72, 101, 108, 108, 111)` → `"Hello"`

## String Validation

### `isEmpty(string)`
Checks if a string is empty or contains only whitespace.
- **Parameters:** `string` (string) - String to check
- **Returns:** Boolean indicating emptiness
- **Example:** `isEmpty("   ")` → `true`

### `isNumeric(string)`
Checks if a string represents a valid number.
- **Parameters:** `string` (string) - String to check
- **Returns:** Boolean indicating numeric validity
- **Example:** `isNumeric("123.45")` → `true`

### `isAlpha(string)`
Checks if a string contains only alphabetic characters.
- **Parameters:** `string` (string) - String to check
- **Returns:** Boolean indicating alphabetic content
- **Example:** `isAlpha("hello")` → `true`

### `isAlphanumeric(string)`
Checks if a string contains only alphanumeric characters.
- **Parameters:** `string` (string) - String to check
- **Returns:** Boolean indicating alphanumeric content
- **Example:** `isAlphanumeric("hello123")` → `true`

### `isSpace(string)`
Checks if a string contains only whitespace characters.
- **Parameters:** `string` (string) - String to check
- **Returns:** Boolean indicating whitespace-only content
- **Example:** `isSpace("   \t\n")` → `true`

## Type Conversion

### `toString(value)`
Converts a value to its string representation.
- **Parameters:** `value` (any) - Value to convert
- **Returns:** String representation
- **Example:** `toString(123)` → `"123"`

### `toNumber(string)`
Converts a string to a number.
- **Parameters:** `string` (string) - String to convert
- **Returns:** Number value or 0 if invalid
- **Example:** `toNumber("123.45")` → `123.45`

## Usage Notes

### Unicode Handling
- `len()` counts Unicode characters, not bytes
- `size()` returns byte count (useful for UTF-8 strings)
- Character operations work with Unicode properly

### Regular Expressions
- Use double backslashes in patterns: `"\\d+"` for digits
- JavaScript regex syntax is supported
- Case-sensitive by default

### Negative Indexing
- `substr()` supports negative start positions (-1 = last character)
- Negative indices count from the end of the string

### Performance Considerations
- String operations create new strings (immutable)
- Regular expressions are compiled each time (consider caching for repeated use)
- Large string manipulations may impact performance

### Common Patterns
```javascript
// Clean and format input
const clean = trim(lower(input));

// Extract filename from path
const filename = right(path, len(path) - lastIndexOf(path, "/") - 1);

// Validate and parse
if (isNumeric(input)) {
    const number = toNumber(input);
    // Use number
}
```