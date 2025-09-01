# List Package

The List package provides comprehensive functions for list manipulation, transformation, and analysis operations including creation, modification, searching, sorting, and functional programming utilities.

## Basic List Information

### `length(list)`
Returns the number of elements in a list.
- **Parameters:** `list` (array) - The list to measure
- **Returns:** Number of elements
- **Example:** `length([1, 2, 3, 4])` → `4`

### `isEmpty(list)`
Checks if a list contains no elements.
- **Parameters:** `list` (array) - The list to check
- **Returns:** Boolean indicating if list is empty
- **Example:** `isEmpty([])` → `true`

## Element Access

### `get(list, index, default?)`
Retrieves an element at a specific index.
- **Parameters:** 
  - `list` (array) - The list to access
  - `index` (number) - Index position (supports negative indexing)
  - `default` (any, optional) - Value to return if index is out of bounds
- **Returns:** Element at index or default value
- **Examples:**
  - `get([1, 2, 3], 1)` → `2`
  - `get([1, 2, 3], -1)` → `3` (last element)
  - `get([1, 2, 3], 10, "default")` → `"default"`

### `set(list, index, value)`
Creates a new list with an element set at a specific index.
- **Parameters:** 
  - `list` (array) - The original list
  - `index` (number) - Index position (supports negative indexing)
  - `value` (any) - Value to set
- **Returns:** New list with updated value
- **Example:** `set([1, 2, 3], 1, 99)` → `[1, 99, 3]`

### `first(list, default?)`
Gets the first element of a list.
- **Parameters:** 
  - `list` (array) - The list
  - `default` (any, optional) - Value to return if list is empty
- **Returns:** First element or default
- **Example:** `first([1, 2, 3])` → `1`

### `last(list, default?)`
Gets the last element of a list.
- **Parameters:** 
  - `list` (array) - The list
  - `default` (any, optional) - Value to return if list is empty
- **Returns:** Last element or default
- **Example:** `last([1, 2, 3])` → `3`

### `head(list, default?)`
Alias for `first()` - gets the first element.
- **Parameters:** Same as `first()`
- **Returns:** First element or default
- **Example:** `head([1, 2, 3])` → `1`

## List Modification

### `append(list, ...values)`
Creates a new list with values added to the end.
- **Parameters:** 
  - `list` (array) - The original list
  - `...values` (any) - Values to append
- **Returns:** New list with appended values
- **Example:** `append([1, 2], 3, 4)` → `[1, 2, 3, 4]`

### `prepend(list, ...values)`
Creates a new list with values added to the beginning.
- **Parameters:** 
  - `list` (array) - The original list
  - `...values` (any) - Values to prepend
- **Returns:** New list with prepended values
- **Example:** `prepend([3, 4], 1, 2)` → `[1, 2, 3, 4]`

### `insert(list, index, value)`
Creates a new list with a value inserted at a specific index.
- **Parameters:** 
  - `list` (array) - The original list
  - `index` (number) - Position to insert (supports negative indexing)
  - `value` (any) - Value to insert
- **Returns:** New list with inserted value
- **Example:** `insert([1, 3, 4], 1, 2)` → `[1, 2, 3, 4]`

### `remove(list, index)`
Creates a new list with an element removed at a specific index.
- **Parameters:** 
  - `list` (array) - The original list
  - `index` (number) - Index of element to remove
- **Returns:** New list with element removed
- **Example:** `remove([1, 2, 3, 4], 1)` → `[1, 3, 4]`

## List Combination

### `concat(...lists)`
Combines multiple lists into one.
- **Parameters:** `...lists` (array) - Lists to concatenate
- **Returns:** New combined list
- **Example:** `concat([1, 2], [3, 4], [5])` → `[1, 2, 3, 4, 5]`

## List Sections

### `tail(list)`
Gets all elements except the first.
- **Parameters:** `list` (array) - The list
- **Returns:** New list without first element
- **Example:** `tail([1, 2, 3, 4])` → `[2, 3, 4]`

### `rest(list)`
Alias for `tail()` - gets all elements except the first.
- **Parameters:** Same as `tail()`
- **Returns:** New list without first element
- **Example:** `rest([1, 2, 3, 4])` → `[2, 3, 4]`

### `init(list)`
Gets all elements except the last.
- **Parameters:** `list` (array) - The list
- **Returns:** New list without last element
- **Example:** `init([1, 2, 3, 4])` → `[1, 2, 3]`

### `slice(list, start, end?)`
Extracts a section of a list.
- **Parameters:** 
  - `list` (array) - The list
  - `start` (number) - Starting index (supports negative)
  - `end` (number, optional) - Ending index (exclusive, supports negative)
- **Returns:** New list containing the slice
- **Example:** `slice([1, 2, 3, 4, 5], 1, 4)` → `[2, 3, 4]`

### `take(list, count)`
Gets the first N elements.
- **Parameters:** 
  - `list` (array) - The list
  - `count` (number) - Number of elements to take
- **Returns:** New list with first N elements
- **Example:** `take([1, 2, 3, 4, 5], 3)` → `[1, 2, 3]`

### `drop(list, count)`
Gets all elements except the first N.
- **Parameters:** 
  - `list` (array) - The list
  - `count` (number) - Number of elements to drop
- **Returns:** New list without first N elements
- **Example:** `drop([1, 2, 3, 4, 5], 2)` → `[3, 4, 5]`

## List Transformation

### `reverse(list)`
Creates a new list with elements in reverse order.
- **Parameters:** `list` (array) - The list to reverse
- **Returns:** New reversed list
- **Example:** `reverse([1, 2, 3, 4])` → `[4, 3, 2, 1]`

### `sort(list)`
Creates a new sorted list in ascending order.
- **Parameters:** `list` (array) - The list to sort
- **Returns:** New sorted list
- **Example:** `sort([3, 1, 4, 2])` → `[1, 2, 3, 4]`

### `sortDesc(list)`
Creates a new sorted list in descending order.
- **Parameters:** `list` (array) - The list to sort
- **Returns:** New sorted list in descending order
- **Example:** `sortDesc([3, 1, 4, 2])` → `[4, 3, 2, 1]`

### `shuffle(list)`
Creates a new list with elements in random order.
- **Parameters:** `list` (array) - The list to shuffle
- **Returns:** New randomly shuffled list
- **Example:** `shuffle([1, 2, 3, 4])` → `[3, 1, 4, 2]` (random)

### `unique(list)`
Creates a new list with duplicate values removed.
- **Parameters:** `list` (array) - The list to deduplicate
- **Returns:** New list with unique values only
- **Example:** `unique([1, 2, 2, 3, 1, 4])` → `[1, 2, 3, 4]`

### `flatten(list, depth?)`
Flattens nested lists up to specified depth.
- **Parameters:** 
  - `list` (array) - The list to flatten
  - `depth` (number, optional) - Depth to flatten (default: 1)
- **Returns:** New flattened list
- **Examples:**
  - `flatten([[1, 2], [3, 4]])` → `[1, 2, 3, 4]`
  - `flatten([[[1, 2]], [[3, 4]]], 2)` → `[1, 2, 3, 4]`

## List Search and Analysis

### `contains(list, value)`
Checks if a list contains a specific value.
- **Parameters:** 
  - `list` (array) - The list to search
  - `value` (any) - Value to find
- **Returns:** Boolean indicating if value exists
- **Example:** `contains([1, 2, 3], 2)` → `true`

### `indexOf(list, value, start?)`
Finds the first index of a value in a list.
- **Parameters:** 
  - `list` (array) - The list to search
  - `value` (any) - Value to find
  - `start` (number, optional) - Starting search index
- **Returns:** Index of value or -1 if not found
- **Example:** `indexOf([1, 2, 3, 2], 2)` → `1`

### `lastIndexOf(list, value)`
Finds the last index of a value in a list.
- **Parameters:** 
  - `list` (array) - The list to search
  - `value` (any) - Value to find
- **Returns:** Last index of value or -1 if not found
- **Example:** `lastIndexOf([1, 2, 3, 2], 2)` → `3`

### `count(list, value)`
Counts occurrences of a value in a list.
- **Parameters:** 
  - `list` (array) - The list to search
  - `value` (any) - Value to count
- **Returns:** Number of occurrences
- **Example:** `count([1, 2, 2, 3, 2], 2)` → `3`

## List Generation

### `range(end)` / `range(start, end)` / `range(start, end, step)`
Creates a list of numbers within a range.
- **Parameters:** 
  - `end` (number) - End value (exclusive)
  - `start` (number, optional) - Start value (default: 0)
  - `step` (number, optional) - Step increment (default: 1)
- **Returns:** List of numbers in range
- **Examples:**
  - `range(5)` → `[0, 1, 2, 3, 4]`
  - `range(2, 7)` → `[2, 3, 4, 5, 6]`
  - `range(0, 10, 2)` → `[0, 2, 4, 6, 8]`

### `repeat(value, count)`
Creates a list with a value repeated N times.
- **Parameters:** 
  - `value` (any) - Value to repeat
  - `count` (number) - Number of repetitions
- **Returns:** List with repeated values
- **Example:** `repeat("hello", 3)` → `["hello", "hello", "hello"]`

## Advanced Operations

### `zip(...lists)`
Combines multiple lists element-wise into tuples.
- **Parameters:** `...lists` (array) - Lists to zip together
- **Returns:** List of tuples (arrays)
- **Example:** `zip([1, 2, 3], ["a", "b", "c"])` → `[[1, "a"], [2, "b"], [3, "c"]]`

### `filter(list)`
Removes null and falsy values from a list.
- **Parameters:** `list` (array) - The list to filter
- **Returns:** New list with non-null, truthy values
- **Example:** `filter([1, null, 2, "", 3, false, 4])` → `[1, 2, 3, 4]`

### `map(list)`
Identity function that returns a copy of the list.
- **Parameters:** `list` (array) - The list to copy
- **Returns:** Copy of the original list
- **Example:** `map([1, 2, 3])` → `[1, 2, 3]`

## Usage Notes

### Negative Indexing
Many functions support negative indexing where `-1` refers to the last element:
- `get([1, 2, 3], -1)` → `3`
- `set([1, 2, 3], -2, 99)` → `[1, 99, 3]`

### Immutability
Most functions return new lists rather than modifying the original:
- Original lists remain unchanged
- Chain operations safely: `reverse(sort(unique(list)))`

### Type Handling
- Mixed-type lists are supported
- Sorting attempts numeric comparison first, falls back to string comparison
- Comparison functions handle different data types gracefully

### Performance Considerations
- Large lists are handled efficiently
- Operations like `unique()` and `flatten()` may take longer on very large datasets
- `shuffle()` uses Fisher-Yates algorithm for uniform randomness