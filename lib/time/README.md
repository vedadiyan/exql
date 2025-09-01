# Time Package

The Time package provides comprehensive functions for time manipulation, parsing, formatting, calculations, and validation operations with support for timestamps, date arithmetic, and timezone handling.

## Current Time Functions

### `now()`
Returns the current Unix timestamp in seconds.
- **Parameters:** None
- **Returns:** Current timestamp as number
- **Example:** `now()` → `1693526400` (example timestamp)

### `nowMillis()`
Returns the current Unix timestamp in milliseconds.
- **Parameters:** None
- **Returns:** Current timestamp in milliseconds
- **Example:** `nowMillis()` → `1693526400000`

### `nowNanos()`
Returns the current Unix timestamp in nanoseconds.
- **Parameters:** None
- **Returns:** Current timestamp in nanoseconds
- **Example:** `nowNanos()` → `1693526400000000000`

## Time Parsing and Formatting

### `parse(timeString, layout?)`
Parses a time string into a Unix timestamp.
- **Parameters:** 
  - `timeString` (string) - Time string to parse
  - `layout` (string, optional) - Time format layout (default: RFC3339)
- **Returns:** Unix timestamp as number
- **Examples:**
  - `parse("2023-08-31T12:00:00Z")` → `1693483200`
  - `parse("2023-08-31 12:00:00", "YYYY-MM-DD HH:mm:ss")` → `1693483200`

### `format(timestamp, layout?)`
Formats a Unix timestamp into a string.
- **Parameters:** 
  - `timestamp` (number) - Unix timestamp
  - `layout` (string, optional) - Format layout (default: RFC3339)
- **Returns:** Formatted time string
- **Examples:**
  - `format(1693483200)` → `"2023-08-31T12:00:00Z"`
  - `format(1693483200, "YYYY-MM-DD HH:mm:ss")` → `"2023-08-31 12:00:00"`

## Time Arithmetic

### `add(timestamp, seconds)`
Adds seconds to a timestamp.
- **Parameters:** 
  - `timestamp` (number) - Original timestamp
  - `seconds` (number) - Seconds to add
- **Returns:** New timestamp
- **Example:** `add(1693483200, 3600)` → `1693486800` (1 hour later)

### `addDays(timestamp, days)`
Adds days to a timestamp.
- **Parameters:** 
  - `timestamp` (number) - Original timestamp
  - `days` (number) - Days to add
- **Returns:** New timestamp
- **Example:** `addDays(1693483200, 7)` → `1694088000` (7 days later)

### `addHours(timestamp, hours)`
Adds hours to a timestamp.
- **Parameters:** 
  - `timestamp` (number) - Original timestamp
  - `hours` (number) - Hours to add
- **Returns:** New timestamp
- **Example:** `addHours(1693483200, 24)` → `1693569600` (24 hours later)

### `addMinutes(timestamp, minutes)`
Adds minutes to a timestamp.
- **Parameters:** 
  - `timestamp` (number) - Original timestamp
  - `minutes` (number) - Minutes to add
- **Returns:** New timestamp
- **Example:** `addMinutes(1693483200, 30)` → `1693485000` (30 minutes later)

## Time Differences

### `diff(timestamp1, timestamp2)`
Calculates the difference between two timestamps in seconds.
- **Parameters:** 
  - `timestamp1` (number) - First timestamp
  - `timestamp2` (number) - Second timestamp
- **Returns:** Difference in seconds (timestamp1 - timestamp2)
- **Example:** `diff(1693486800, 1693483200)` → `3600` (1 hour difference)

### `diffDays(timestamp1, timestamp2)`
Calculates the difference between two timestamps in days.
- **Parameters:** Same as `diff()`
- **Returns:** Difference in days
- **Example:** `diffDays(1694088000, 1693483200)` → `7`

### `diffHours(timestamp1, timestamp2)`
Calculates the difference between two timestamps in hours.
- **Parameters:** Same as `diff()`
- **Returns:** Difference in hours
- **Example:** `diffHours(1693486800, 1693483200)` → `1`

### `diffMinutes(timestamp1, timestamp2)`
Calculates the difference between two timestamps in minutes.
- **Parameters:** Same as `diff()`
- **Returns:** Difference in minutes
- **Example:** `diffMinutes(1693485000, 1693483200)` → `30`

## Date Component Extraction

### `year(timestamp)`
Extracts the year from a timestamp.
- **Parameters:** `timestamp` (number) - Unix timestamp
- **Returns:** Year as number
- **Example:** `year(1693483200)` → `2023`

### `month(timestamp)`
Extracts the month from a timestamp.
- **Parameters:** `timestamp` (number) - Unix timestamp
- **Returns:** Month as number (1-12)
- **Example:** `month(1693483200)` → `8` (August)

### `day(timestamp)`
Extracts the day of month from a timestamp.
- **Parameters:** `timestamp` (number) - Unix timestamp
- **Returns:** Day as number (1-31)
- **Example:** `day(1693483200)` → `31`

### `hour(timestamp)`
Extracts the hour from a timestamp.
- **Parameters:** `timestamp` (number) - Unix timestamp
- **Returns:** Hour as number (0-23)
- **Example:** `hour(1693483200)` → `12`

### `minute(timestamp)`
Extracts the minute from a timestamp.
- **Parameters:** `timestamp` (number) - Unix timestamp
- **Returns:** Minute as number (0-59)
- **Example:** `minute(1693483200)` → `0`

### `second(timestamp)`
Extracts the second from a timestamp.
- **Parameters:** `timestamp` (number) - Unix timestamp
- **Returns:** Second as number (0-59)
- **Example:** `second(1693483200)` → `0`

### `weekday(timestamp)`
Gets the day of the week from a timestamp.
- **Parameters:** `timestamp` (number) - Unix timestamp
- **Returns:** Weekday as number (0=Sunday, 1=Monday, ..., 6=Saturday)
- **Example:** `weekday(1693483200)` → `4` (Thursday)

### `yearday(timestamp)`
Gets the day of the year from a timestamp.
- **Parameters:** `timestamp` (number) - Unix timestamp
- **Returns:** Day of year as number (1-366)
- **Example:** `yearday(1693483200)` → `243`

### `week(timestamp)`
Gets the ISO week number from a timestamp.
- **Parameters:** `timestamp` (number) - Unix timestamp
- **Returns:** ISO week number (1-53)
- **Example:** `week(1693483200)` → `35`

## Time Boundaries

### `startOfDay(timestamp)`
Gets the timestamp for the start of the day (00:00:00).
- **Parameters:** `timestamp` (number) - Any timestamp in the day
- **Returns:** Timestamp for start of day
- **Example:** `startOfDay(1693526400)` → `1693440000`

### `endOfDay(timestamp)`
Gets the timestamp for the end of the day (23:59:59).
- **Parameters:** `timestamp` (number) - Any timestamp in the day
- **Returns:** Timestamp for end of day
- **Example:** `endOfDay(1693526400)` → `1693526399`

### `startOfWeek(timestamp)`
Gets the timestamp for the start of the week (Monday 00:00:00).
- **Parameters:** `timestamp` (number) - Any timestamp in the week
- **Returns:** Timestamp for start of week
- **Example:** `startOfWeek(1693526400)` → `1693180800`

### `startOfMonth(timestamp)`
Gets the timestamp for the start of the month (1st day 00:00:00).
- **Parameters:** `timestamp` (number) - Any timestamp in the month
- **Returns:** Timestamp for start of month
- **Example:** `startOfMonth(1693526400)` → `1690848000`

### `startOfYear(timestamp)`
Gets the timestamp for the start of the year (January 1st 00:00:00).
- **Parameters:** `timestamp` (number) - Any timestamp in the year
- **Returns:** Timestamp for start of year
- **Example:** `startOfYear(1693526400)` → `1672531200`

## Date Analysis

### `isWeekend(timestamp)`
Checks if a timestamp falls on a weekend.
- **Parameters:** `timestamp` (number) - Timestamp to check
- **Returns:** Boolean indicating weekend (Saturday or Sunday)
- **Example:** `isWeekend(1693526400)` → `false`

### `isLeapYear(timestamp)`
Checks if the year of a timestamp is a leap year.
- **Parameters:** `timestamp` (number) - Timestamp to check
- **Returns:** Boolean indicating leap year
- **Example:** `isLeapYear(1577836800)` → `true` (2020 was a leap year)

### `daysInMonth(timestamp)`
Returns the number of days in the month of a timestamp.
- **Parameters:** `timestamp` (number) - Timestamp in the month
- **Returns:** Number of days in that month
- **Example:** `daysInMonth(1693526400)` → `31` (August has 31 days)

### `age(birthTimestamp, currentTimestamp?)`
Calculates age in years between two timestamps.
- **Parameters:** 
  - `birthTimestamp` (number) - Birth date timestamp
  - `currentTimestamp` (number, optional) - Current time (default: now)
- **Returns:** Age in years
- **Example:** `age(946684800, 1693526400)` → `23` (born 2000, checked in 2023)

## Timezone Operations

### `toTimezone(timestamp, timezone)`
Converts a UTC timestamp to a specific timezone.
- **Parameters:** 
  - `timestamp` (number) - UTC timestamp
  - `timezone` (string) - Target timezone (e.g., "America/New_York")
- **Returns:** Timestamp adjusted for timezone
- **Example:** `toTimezone(1693483200, "America/New_York")` → `1693468800`

### `fromTimezone(timestamp, timezone)`
Converts a timestamp from a specific timezone to UTC.
- **Parameters:** 
  - `timestamp` (number) - Timestamp in specified timezone
  - `timezone` (string) - Source timezone
- **Returns:** UTC timestamp
- **Example:** `fromTimezone(1693468800, "America/New_York")` → `1693483200`

## Utility Functions

### `sleep(seconds)`
Pauses execution for a specified number of seconds.
- **Parameters:** `seconds` (number) - Duration to sleep
- **Returns:** Always returns true after sleeping
- **Example:** `sleep(2)` (pauses for 2 seconds)

### `validate(timeString, layout?)`
Validates if a string can be parsed as a valid time.
- **Parameters:** 
  - `timeString` (string) - Time string to validate
  - `layout` (string, optional) - Expected format (default: RFC3339)
- **Returns:** Boolean indicating validity
- **Example:** `validate("2023-08-31T12:00:00Z")` → `true`

### `range(start, end, step)`
Creates a range of timestamps with specified step.
- **Parameters:** 
  - `start` (number) - Start timestamp
  - `end` (number) - End timestamp
  - `step` (number) - Step in seconds
- **Returns:** Array of timestamps
- **Example:** `range(1693483200, 1693490400, 3600)` → `[1693483200, 1693486800, 1693490400]` (hourly)

## Time Formats

### Standard Layout Tokens
- `YYYY` - 4-digit year
- `YY` - 2-digit year
- `MM` - 2-digit month
- `DD` - 2-digit day
- `HH` - 24-hour format hour
- `mm` - Minutes
- `ss` - Seconds
- `SSS` - Milliseconds

### Predefined Layouts
- `"ISO8601"` or `"RFC3339"` - ISO 8601 format
- `"RFC822"` - RFC 822 format
- `"RFC850"` - RFC 850 format
- `"RFC1123"` - RFC 1123 format
- `"Kitchen"` - Kitchen timer format
- `"Stamp"` - Simple timestamp format

## Usage Notes

### Timestamp Format
- All timestamps are Unix timestamps (seconds since January 1, 1970 UTC)
- Functions work with UTC time by default
- Use timezone functions for local time conversions

### Date Arithmetic
- All arithmetic functions work with seconds as the base unit
- Helper functions provide convenience for common intervals (days, hours, minutes)
- Results maintain Unix timestamp format

### Leap Year Rules
- Divisible by 4: leap year
- Divisible by 100: not leap year
- Divisible by 400: leap year
- Example: 2000 (leap), 1900 (not leap), 2004 (leap)

### Weekend Definition
- Saturday and Sunday are considered weekend days
- Monday is considered the start of the week for `startOfWeek()`

### Timezone Handling
- Timezone names follow IANA Time Zone Database format
- Examples: "America/New_York", "Europe/London", "Asia/Tokyo"
- Use "UTC" or empty string for UTC timezone

### Error Handling
- Invalid time strings return parsing errors
- Invalid timestamps may produce unexpected results
- Timezone functions validate timezone names

### Common Patterns
```javascript
// Get current day boundaries
const dayStart = startOfDay(now());
const dayEnd = endOfDay(now());

// Calculate time until deadline
const deadline = parse("2023-12-31T23:59:59Z");
const timeLeft = diff(deadline, now());
const daysLeft = diffDays(deadline, now());

// Format for display
const formatted = format(now(), "YYYY-MM-DD HH:mm:ss");

// Work with different timezones
const utcTime = now();
const localTime = toTimezone(utcTime, "America/New_York");
const backToUtc = fromTimezone(localTime, "America/New_York");
```