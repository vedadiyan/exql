# Math Package

The Math package provides comprehensive mathematical functions including basic arithmetic, trigonometry, statistics, number theory, and random number generation.

## Basic Arithmetic Functions

### `abs(number)`
Returns the absolute value of a number.
- **Parameters:** `number` (number) - The number to get absolute value of
- **Returns:** Absolute value
- **Example:** `abs(-5)` → `5`

### `sign(number)`
Returns the sign of a number.
- **Parameters:** `number` (number) - The number to check
- **Returns:** 1 for positive, -1 for negative, 0 for zero
- **Example:** `sign(-42)` → `-1`

### `max(...numbers)`
Returns the largest of the given numbers.
- **Parameters:** `...numbers` (number) - Numbers to compare
- **Returns:** Maximum value
- **Example:** `max(1, 5, 3, 9, 2)` → `9`

### `min(...numbers)`
Returns the smallest of the given numbers.
- **Parameters:** `...numbers` (number) - Numbers to compare
- **Returns:** Minimum value
- **Example:** `min(1, 5, 3, 9, 2)` → `1`

### `clamp(value, min, max)`
Constrains a value within a range.
- **Parameters:** 
  - `value` (number) - Value to clamp
  - `min` (number) - Minimum allowed value
  - `max` (number) - Maximum allowed value
- **Returns:** Clamped value
- **Example:** `clamp(15, 1, 10)` → `10`

## Rounding Functions

### `ceil(number)`
Rounds a number up to the nearest integer.
- **Parameters:** `number` (number) - The number to round up
- **Returns:** Ceiling value
- **Example:** `ceil(4.2)` → `5`

### `floor(number)`
Rounds a number down to the nearest integer.
- **Parameters:** `number` (number) - The number to round down
- **Returns:** Floor value
- **Example:** `floor(4.8)` → `4`

### `round(number, precision?)`
Rounds a number to the nearest integer or specified decimal places.
- **Parameters:** 
  - `number` (number) - The number to round
  - `precision` (number, optional) - Decimal places (default: 0)
- **Returns:** Rounded value
- **Examples:**
  - `round(4.6)` → `5`
  - `round(3.14159, 2)` → `3.14`

### `trunc(number)`
Truncates the decimal part of a number.
- **Parameters:** `number` (number) - The number to truncate
- **Returns:** Integer part only
- **Example:** `trunc(4.8)` → `4`

## Power and Root Functions

### `pow(base, exponent)`
Raises a number to a power.
- **Parameters:** 
  - `base` (number) - Base number
  - `exponent` (number) - Power to raise to
- **Returns:** Result of base^exponent
- **Example:** `pow(2, 3)` → `8`

### `sqrt(number)`
Returns the square root of a number.
- **Parameters:** `number` (number) - The number (must be non-negative)
- **Returns:** Square root, or NaN for negative numbers
- **Example:** `sqrt(16)` → `4`

### `cbrt(number)`
Returns the cube root of a number.
- **Parameters:** `number` (number) - The number
- **Returns:** Cube root
- **Example:** `cbrt(27)` → `3`

## Exponential and Logarithmic Functions

### `exp(number)`
Returns e raised to the power of the given number.
- **Parameters:** `number` (number) - The exponent
- **Returns:** e^number
- **Example:** `exp(1)` → `2.718...`

### `exp2(number)`
Returns 2 raised to the power of the given number.
- **Parameters:** `number` (number) - The exponent
- **Returns:** 2^number
- **Example:** `exp2(3)` → `8`

### `log(number)`
Returns the natural logarithm (base e) of a number.
- **Parameters:** `number` (number) - The number (must be positive)
- **Returns:** Natural logarithm, or NaN for non-positive numbers
- **Example:** `log(Math.E)` → `1`

### `log10(number)`
Returns the base-10 logarithm of a number.
- **Parameters:** `number` (number) - The number (must be positive)
- **Returns:** Base-10 logarithm, or NaN for non-positive numbers
- **Example:** `log10(100)` → `2`

### `log2(number)`
Returns the base-2 logarithm of a number.
- **Parameters:** `number` (number) - The number (must be positive)
- **Returns:** Base-2 logarithm, or NaN for non-positive numbers
- **Example:** `log2(8)` → `3`

## Trigonometric Functions

### `sin(radians)`
Returns the sine of an angle in radians.
- **Parameters:** `radians` (number) - Angle in radians
- **Returns:** Sine value
- **Example:** `sin(Math.PI / 2)` → `1`

### `cos(radians)`
Returns the cosine of an angle in radians.
- **Parameters:** `radians` (number) - Angle in radians
- **Returns:** Cosine value
- **Example:** `cos(0)` → `1`

### `tan(radians)`
Returns the tangent of an angle in radians.
- **Parameters:** `radians` (number) - Angle in radians
- **Returns:** Tangent value
- **Example:** `tan(Math.PI / 4)` → `1`

### `asin(value)`
Returns the arcsine of a value in radians.
- **Parameters:** `value` (number) - Value between -1 and 1
- **Returns:** Arcsine in radians, or NaN if out of range
- **Example:** `asin(1)` → `1.570...` (π/2)

### `acos(value)`
Returns the arccosine of a value in radians.
- **Parameters:** `value` (number) - Value between -1 and 1
- **Returns:** Arccosine in radians, or NaN if out of range
- **Example:** `acos(0)` → `1.570...` (π/2)

### `atan(value)`
Returns the arctangent of a value in radians.
- **Parameters:** `value` (number) - The value
- **Returns:** Arctangent in radians
- **Example:** `atan(1)` → `0.785...` (π/4)

### `atan2(y, x)`
Returns the angle from the X-axis to a point (x,y) in radians.
- **Parameters:** 
  - `y` (number) - Y coordinate
  - `x` (number) - X coordinate
- **Returns:** Angle in radians
- **Example:** `atan2(1, 1)` → `0.785...` (π/4)

## Hyperbolic Functions

### `sinh(number)`
Returns the hyperbolic sine of a number.
- **Parameters:** `number` (number) - The number
- **Returns:** Hyperbolic sine
- **Example:** `sinh(0)` → `0`

### `cosh(number)`
Returns the hyperbolic cosine of a number.
- **Parameters:** `number` (number) - The number
- **Returns:** Hyperbolic cosine
- **Example:** `cosh(0)` → `1`

### `tanh(number)`
Returns the hyperbolic tangent of a number.
- **Parameters:** `number` (number) - The number
- **Returns:** Hyperbolic tangent
- **Example:** `tanh(0)` → `0`

## Angle Conversion

### `radians(degrees)`
Converts degrees to radians.
- **Parameters:** `degrees` (number) - Angle in degrees
- **Returns:** Angle in radians
- **Example:** `radians(180)` → `3.141...` (π)

### `degrees(radians)`
Converts radians to degrees.
- **Parameters:** `radians` (number) - Angle in radians
- **Returns:** Angle in degrees
- **Example:** `degrees(Math.PI)` → `180`

## Statistical Functions

### `sum(...values)`
Calculates the sum of numbers or arrays.
- **Parameters:** `...values` (number|array) - Numbers or arrays to sum
- **Returns:** Sum of all values
- **Example:** `sum(1, 2, 3, [4, 5])` → `15`

### `mean(...values)`
Calculates the arithmetic mean (average).
- **Parameters:** `...values` (number|array) - Numbers or arrays
- **Returns:** Mean value
- **Example:** `mean(1, 2, 3, 4, 5)` → `3`

### `median(...values)`
Calculates the median (middle value).
- **Parameters:** `...values` (number|array) - Numbers or arrays
- **Returns:** Median value
- **Example:** `median(1, 2, 3, 4, 5)` → `3`

### `mode(...values)`
Finds the most frequently occurring value.
- **Parameters:** `...values` (number|array) - Numbers or arrays
- **Returns:** Most frequent value
- **Example:** `mode(1, 2, 2, 3, 2)` → `2`

### `variance(...values)`
Calculates the sample variance.
- **Parameters:** `...values` (number|array) - Numbers or arrays
- **Returns:** Sample variance
- **Example:** `variance(1, 2, 3, 4, 5)` → `2.5`

### `stddev(...values)`
Calculates the sample standard deviation.
- **Parameters:** `...values` (number|array) - Numbers or arrays
- **Returns:** Sample standard deviation
- **Example:** `stddev(1, 2, 3, 4, 5)` → `1.58...`

## Random Number Generation

### `random()` / `random(max)` / `random(min, max)`
Generates random numbers.
- **Parameters:** 
  - No args: Random float 0-1
  - `max` (number): Random integer 0 to max-1
  - `min, max` (number): Random integer min to max-1
- **Returns:** Random number
- **Examples:**
  - `random()` → `0.123...`
  - `random(10)` → `7`
  - `random(5, 15)` → `12`

### `randomSeed(seed?)`
Sets the random number generator seed.
- **Parameters:** `seed` (number, optional) - Seed value (uses current time if omitted)
- **Returns:** Always returns true
- **Example:** `randomSeed(12345)`

### `randomFloat()` / `randomFloat(max)` / `randomFloat(min, max)`
Generates random floating-point numbers.
- **Parameters:** Similar to `random()` but returns floats
- **Returns:** Random floating-point number
- **Example:** `randomFloat(1, 10)` → `7.234...`

## Number Validation

### `isNan(number)`
Checks if a value is NaN (Not a Number).
- **Parameters:** `number` (number) - Value to check
- **Returns:** Boolean indicating if value is NaN
- **Example:** `isNan(parseInt("abc"))` → `true`

### `isInf(number)`
Checks if a value is infinite.
- **Parameters:** `number` (number) - Value to check
- **Returns:** Boolean indicating if value is infinite
- **Example:** `isInf(1/0)` → `true`

### `isFinite(number)`
Checks if a value is finite (not NaN or infinite).
- **Parameters:** `number` (number) - Value to check
- **Returns:** Boolean indicating if value is finite
- **Example:** `isFinite(42)` → `true`

## Number Theory

### `gcd(...numbers)`
Calculates the Greatest Common Divisor.
- **Parameters:** `...numbers` (number) - Numbers to find GCD of
- **Returns:** Greatest common divisor
- **Example:** `gcd(12, 18, 24)` → `6`

### `lcm(...numbers)`
Calculates the Least Common Multiple.
- **Parameters:** `...numbers` (number) - Numbers to find LCM of
- **Returns:** Least common multiple
- **Example:** `lcm(4, 6, 8)` → `24`

### `factorial(number)`
Calculates the factorial of a number.
- **Parameters:** `number` (number) - Non-negative integer
- **Returns:** Factorial, or NaN for negative numbers
- **Example:** `factorial(5)` → `120`

## Mathematical Constants

### `pi()`
Returns the value of π (pi).
- **Parameters:** None
- **Returns:** Value of π
- **Example:** `pi()` → `3.141592653589793`

### `e()`
Returns the value of e (Euler's number).
- **Parameters:** None
- **Returns:** Value of e
- **Example:** `e()` → `2.718281828459045`

### `phi()`
Returns the value of φ (golden ratio).
- **Parameters:** None
- **Returns:** Value of φ
- **Example:** `phi()` → `1.618033988749895`

## Usage Notes

### Input Flexibility
Statistical functions accept:
- Multiple number arguments: `mean(1, 2, 3, 4, 5)`
- Arrays: `mean([1, 2, 3, 4, 5])`
- Mixed: `mean(1, 2, [3, 4], 5)`

### Error Handling
- Invalid operations return `NaN` (e.g., `sqrt(-1)`, `log(0)`)
- Division by zero returns `Infinity` or `-Infinity`
- Statistical functions require at least one numeric value

### Precision
- Floating-point arithmetic limitations apply
- Use appropriate rounding for display purposes
- Be aware of precision loss in very large calculations

### Random Numbers
- `random()` functions use JavaScript's Math.random() internally
- `randomSeed()` provides reproducible sequences
- Integer random functions are inclusive of min, exclusive of max