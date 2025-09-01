/*
 * Copyright 2025 Pouya Vedadiyan
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package time

import (
	"math"
	"testing"
	"time"

	"github.com/vedadiyan/exql/lang"
)

func TestNow(t *testing.T) {
	_, fn := now()

	result, err := fn([]lang.Value{})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	timestamp := float64(result.(lang.NumberValue))
	currentTime := float64(time.Now().Unix())

	// Should be within 1 second of current time
	if math.Abs(timestamp-currentTime) > 1 {
		t.Errorf("Timestamp too far from current time: got %f, expected around %f", timestamp, currentTime)
	}
}

func TestNowMillis(t *testing.T) {
	_, fn := nowMillis()

	result, err := fn([]lang.Value{})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	timestamp := float64(result.(lang.NumberValue))
	currentTime := float64(time.Now().UnixMilli())

	// Should be within 1000ms of current time
	if math.Abs(timestamp-currentTime) > 1000 {
		t.Errorf("Timestamp too far from current time: got %f, expected around %f", timestamp, currentTime)
	}
}

func TestNowNanos(t *testing.T) {
	_, fn := nowNanos()

	result, err := fn([]lang.Value{})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	timestamp := float64(result.(lang.NumberValue))
	if timestamp <= 0 {
		t.Errorf("Expected positive nanosecond timestamp, got %f", timestamp)
	}
}

func TestParse(t *testing.T) {
	_, fn := parse()

	tests := []struct {
		name     string
		timeStr  string
		layout   *string
		hasError bool
	}{
		{
			"RFC3339 default",
			"2023-01-15T10:30:00Z",
			nil,
			false,
		},
		{
			"custom layout",
			"2023-01-15 10:30:00",
			stringPtr("YYYY-MM-DD HH:mm:ss"),
			false,
		},
		{
			"RFC822 layout",
			"15 Jan 23 10:30 UTC",
			stringPtr("RFC822"),
			false,
		},
		{
			"invalid time string",
			"not-a-time",
			nil,
			true,
		},
		{
			"wrong format",
			"2023-01-15",
			stringPtr("YYYY-MM-DD HH:mm:ss"),
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []lang.Value{lang.StringValue(tt.timeStr)}
			if tt.layout != nil {
				args = append(args, lang.StringValue(*tt.layout))
			}

			result, err := fn(args)
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			timestamp := float64(result.(lang.NumberValue))
			if timestamp <= 0 {
				t.Errorf("Expected positive timestamp, got %f", timestamp)
			}
		})
	}
}

func TestConvertTimeLayout(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"YYYY conversion", "YYYY-MM-DD", "2006-01-02"},
		{"HH:mm:ss conversion", "HH:mm:ss", "15:04:05"},
		{"RFC3339 constant", "RFC3339", time.RFC3339},
		{"RFC822 constant", "RFC822", time.RFC822},
		{"Kitchen constant", "Kitchen", time.Kitchen},
		{"custom pattern", "DD/MM/YYYY HH:mm", "02/01/2006 15:04"},
		{"milliseconds", "YYYY-MM-DD HH:mm:ss.SSS", "2006-01-02 15:04:05.000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertTimeLayout(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestParseFormatRoundtrip(t *testing.T) {
	_, parseFn := parse()
	_, formatFn := toFormat()

	tests := []struct {
		name    string
		timeStr string
		layout  string
	}{
		{"RFC3339", "2023-06-15T14:30:45Z", ""},
		{"custom format", "2023-06-15 14:30:45", "YYYY-MM-DD HH:mm:ss"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse time string
			parseArgs := []lang.Value{lang.StringValue(tt.timeStr)}
			if tt.layout != "" {
				parseArgs = append(parseArgs, lang.StringValue(tt.layout))
			}

			parsed, err := parseFn(parseArgs)
			if err != nil {
				t.Errorf("Parse error: %v", err)
				return
			}

			// Format back to string
			formatArgs := []lang.Value{parsed}
			if tt.layout != "" {
				formatArgs = append(formatArgs, lang.StringValue(tt.layout))
			}

			formatted, err := formatFn(formatArgs)
			if err != nil {
				t.Errorf("Format error: %v", err)
				return
			}

			if string(formatted.(lang.StringValue)) != tt.timeStr {
				t.Errorf("Round trip failed: expected %s, got %s", tt.timeStr, string(formatted.(lang.StringValue)))
			}
		})
	}
}

func TestArgumentErrors(t *testing.T) {
	functions := []struct {
		fn           func() (string, lang.Function)
		expectedArgs int
	}{
		{now, 0},
		{nowMillis, 0},
		{nowNanos, 0},
		{add, 2},
		{addDays, 2},
		{addHours, 2},
		{addMinutes, 2},
		{diff, 2},
		{diffDays, 2},
		{diffHours, 2},
		{diffMinutes, 2},
		{year, 1},
		{month, 1},
		{day, 1},
		{hour, 1},
		{minute, 1},
		{second, 1},
		{weekday, 1},
		{yearday, 1},
		{week, 1},
		{startOfDay, 1},
		{toEndOfDay, 1},
		{startOfWeek, 1},
		{startOfMonth, 1},
		{startOfYear, 1},
		{isWeekend, 1},
		{isLeapYear, 1},
		{daysInMonth, 1},
		{toTimezone, 2},
		{toFromTimezone, 2},
		{sleep, 1},
		{rrange, 3},
	}

	for _, tf := range functions {
		name, fn := tf.fn()
		t.Run(name+"_wrong_args", func(t *testing.T) {
			var args []lang.Value
			if tf.expectedArgs == 0 {
				args = []lang.Value{lang.NumberValue(1)} // Add extra arg for 0-arg functions
			} else {
				args = make([]lang.Value, tf.expectedArgs-1) // One less than expected
				for i := range args {
					args[i] = lang.NumberValue(1000)
				}
			}
			_, err := fn(args)
			if err == nil {
				t.Errorf("Expected error for wrong argument count in %s", name)
			}
		})
	}
}

func TestNonNumericTimestamps(t *testing.T) {
	functions := []func() (string, lang.Function){
		add, addDays, addHours, addMinutes, diff, diffDays, diffHours, diffMinutes,
		year, month, day, hour, minute, second, weekday, yearday, week,
		startOfDay, toEndOfDay, startOfWeek, startOfMonth, startOfYear,
		isWeekend, isLeapYear, daysInMonth, age, sleep,
	}

	nonNumericInput := lang.StringValue("not a number")

	for _, getFn := range functions {
		name, fn := getFn()
		t.Run(name+"_non_numeric", func(t *testing.T) {
			args := []lang.Value{nonNumericInput}
			if name == "add" || name == "addDays" || name == "addHours" || name == "addMinutes" ||
				name == "diff" || name == "diffDays" || name == "diffHours" || name == "diffMinutes" {
				args = append(args, lang.NumberValue(100))
			}

			_, err := fn(args)
			if err == nil {
				t.Errorf("Expected error for non-numeric input in %s", name)
			}
		})
	}
}

func TestComplexScenarios(t *testing.T) {
	t.Run("time_arithmetic_consistency", func(t *testing.T) {
		baseTime := float64(time.Date(2023, 6, 15, 12, 0, 0, 0, time.UTC).Unix())

		_, addDaysFn := addDays()
		_, addHoursFn := addHours()
		_, diffFn := diff()

		// Add 1 day using days function
		dayResult, err := addDaysFn([]lang.Value{lang.NumberValue(baseTime), lang.NumberValue(1)})
		if err != nil {
			t.Errorf("Error adding days: %v", err)
			return
		}

		// Add 24 hours using hours function
		hourResult, err := addHoursFn([]lang.Value{lang.NumberValue(baseTime), lang.NumberValue(24)})
		if err != nil {
			t.Errorf("Error adding hours: %v", err)
			return
		}

		// Both should be equal
		if dayResult != hourResult {
			t.Errorf("Adding 1 day and 24 hours should be equal: %v vs %v", dayResult, hourResult)
		}

		// Verify with diff
		diffResult, err := diffFn([]lang.Value{dayResult, lang.NumberValue(baseTime)})
		if err != nil {
			t.Errorf("Error calculating diff: %v", err)
			return
		}

		expectedDiff := 86400.0 // 1 day in seconds
		if float64(diffResult.(lang.NumberValue)) != expectedDiff {
			t.Errorf("Expected diff of %f, got %f", expectedDiff, float64(diffResult.(lang.NumberValue)))
		}
	})

	t.Run("start_of_periods", func(t *testing.T) {
		// Test with mid-month, mid-week time
		testTime := time.Date(2023, 6, 15, 14, 30, 45, 0, time.UTC) // Thursday
		timestamp := float64(testTime.Unix())

		_, startDayFn := startOfDay()
		_, startWeekFn := startOfWeek()
		_, startMonthFn := startOfMonth()
		_, startYearFn := startOfYear()

		// Get all start times
		startDay, _ := startDayFn([]lang.Value{lang.NumberValue(timestamp)})
		startWeek, _ := startWeekFn([]lang.Value{lang.NumberValue(timestamp)})
		startMonth, _ := startMonthFn([]lang.Value{lang.NumberValue(timestamp)})
		startYear, _ := startYearFn([]lang.Value{lang.NumberValue(timestamp)})

		// Verify hierarchy: year <= month <= week <= day <= original
		dayTime := float64(startDay.(lang.NumberValue))
		weekTime := float64(startWeek.(lang.NumberValue))
		monthTime := float64(startMonth.(lang.NumberValue))
		yearTime := float64(startYear.(lang.NumberValue))

		if !(yearTime <= monthTime && monthTime <= weekTime && weekTime <= dayTime && dayTime <= timestamp) {
			t.Errorf("Start of periods not in correct order: year=%f, month=%f, week=%f, day=%f, original=%f",
				yearTime, monthTime, weekTime, dayTime, timestamp)
		}
	})
}

func TestExport(t *testing.T) {
	functions := Export()

	expectedFunctions := []string{
		"now", "nowMillis", "nowNanos", "parse", "format",
		"add", "addDays", "addHours", "addMinutes",
		"diff", "diffDays", "diffHours", "diffMinutes",
		"year", "month", "day", "hour", "minute", "second",
		"weekday", "yearday", "week",
		"startOfDay", "endOfDay", "startOfWeek", "startOfMonth", "startOfYear",
		"isWeekend", "isLeapYear", "daysInMonth", "age",
		"toTimezone", "fromTimezone", "sleep", "validate", "range",
	}

	if len(functions) != len(expectedFunctions) {
		t.Errorf("Expected %d functions, got %d", len(expectedFunctions), len(functions))
	}

	for _, name := range expectedFunctions {
		if _, exists := functions[name]; !exists {
			t.Errorf("Expected function %s not found", name)
		}
	}
}

func BenchmarkNow(b *testing.B) {
	_, fn := now()

	for i := 0; i < b.N; i++ {
		fn([]lang.Value{})
	}
}

func BenchmarkFormat(b *testing.B) {
	_, fn := toFormat()
	timestamp := float64(time.Now().Unix())
	args := []lang.Value{lang.NumberValue(timestamp)}

	for i := 0; i < b.N; i++ {
		fn(args)
	}
}

func BenchmarkParse(b *testing.B) {
	_, fn := parse()
	timeStr := "2023-06-15T14:30:45Z"
	args := []lang.Value{lang.StringValue(timeStr)}

	for i := 0; i < b.N; i++ {
		fn(args)
	}
}

func BenchmarkYear(b *testing.B) {
	_, fn := year()
	timestamp := float64(time.Now().Unix())
	args := []lang.Value{lang.NumberValue(timestamp)}

	for i := 0; i < b.N; i++ {
		fn(args)
	}
}

func stringPtr(s string) *string {
	return &s
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func TestFormat(t *testing.T) {
	_, fn := toFormat()
	testTime := time.Date(2023, 1, 15, 10, 30, 45, 0, time.UTC)
	timestamp := float64(testTime.Unix())

	tests := []struct {
		name     string
		layout   *string
		expected string
	}{
		{
			"RFC3339 default",
			nil,
			"2023-01-15T10:30:45Z",
		},
		{
			"custom layout",
			stringPtr("YYYY-MM-DD HH:mm:ss"),
			"2023-01-15 10:30:45",
		},
		{
			"RFC822 layout",
			stringPtr("RFC822"),
			"15 Jan 23 10:30 UTC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []lang.Value{lang.NumberValue(timestamp)}
			if tt.layout != nil {
				args = append(args, lang.StringValue(*tt.layout))
			}

			result, err := fn(args)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			formatted := string(result.(lang.StringValue))
			if formatted != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, formatted)
			}
		})
	}
}

func TestAdd(t *testing.T) {
	_, fn := add()
	baseTime := float64(1000)

	tests := []struct {
		name     string
		seconds  float64
		expected float64
	}{
		{"add positive seconds", 100, 1100},
		{"add negative seconds", -100, 900},
		{"add zero", 0, 1000},
		{"add fractional", 0.5, 1000.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.NumberValue(baseTime), lang.NumberValue(tt.seconds)})
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if float64(result.(lang.NumberValue)) != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, float64(result.(lang.NumberValue)))
			}
		})
	}
}

func TestAddDays(t *testing.T) {
	_, fn := addDays()
	baseTime := float64(1000)

	tests := []struct {
		name     string
		days     float64
		expected float64
	}{
		{"add one day", 1, 1000 + 86400},
		{"add multiple days", 7, 1000 + (7 * 86400)},
		{"subtract days", -1, 1000 - 86400},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.NumberValue(baseTime), lang.NumberValue(tt.days)})
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if float64(result.(lang.NumberValue)) != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, float64(result.(lang.NumberValue)))
			}
		})
	}
}

func TestAddHours(t *testing.T) {
	_, fn := addHours()
	baseTime := float64(1000)

	tests := []struct {
		name     string
		hours    float64
		expected float64
	}{
		{"add one hour", 1, 1000 + 3600},
		{"add multiple hours", 24, 1000 + (24 * 3600)},
		{"subtract hours", -1, 1000 - 3600},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.NumberValue(baseTime), lang.NumberValue(tt.hours)})
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if float64(result.(lang.NumberValue)) != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, float64(result.(lang.NumberValue)))
			}
		})
	}
}

func TestAddMinutes(t *testing.T) {
	_, fn := addMinutes()
	baseTime := float64(1000)

	tests := []struct {
		name     string
		minutes  float64
		expected float64
	}{
		{"add one minute", 1, 1000 + 60},
		{"add multiple minutes", 60, 1000 + (60 * 60)},
		{"subtract minutes", -30, 1000 - (30 * 60)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.NumberValue(baseTime), lang.NumberValue(tt.minutes)})
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if float64(result.(lang.NumberValue)) != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, float64(result.(lang.NumberValue)))
			}
		})
	}
}

func TestDiff(t *testing.T) {
	_, fn := diff()

	tests := []struct {
		name     string
		time1    float64
		time2    float64
		expected float64
	}{
		{"positive diff", 2000, 1000, 1000},
		{"negative diff", 1000, 2000, -1000},
		{"zero diff", 1000, 1000, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.NumberValue(tt.time1), lang.NumberValue(tt.time2)})
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if float64(result.(lang.NumberValue)) != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, float64(result.(lang.NumberValue)))
			}
		})
	}
}

func TestDiffDays(t *testing.T) {
	_, fn := diffDays()
	time1 := float64(1000 + 86400) // 1 day later
	time2 := float64(1000)

	result, err := fn([]lang.Value{lang.NumberValue(time1), lang.NumberValue(time2)})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	expected := 1.0
	if float64(result.(lang.NumberValue)) != expected {
		t.Errorf("Expected %f, got %f", expected, float64(result.(lang.NumberValue)))
	}
}

func TestDiffHours(t *testing.T) {
	_, fn := diffHours()
	time1 := float64(1000 + 3600) // 1 hour later
	time2 := float64(1000)

	result, err := fn([]lang.Value{lang.NumberValue(time1), lang.NumberValue(time2)})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	expected := 1.0
	if float64(result.(lang.NumberValue)) != expected {
		t.Errorf("Expected %f, got %f", expected, float64(result.(lang.NumberValue)))
	}
}

func TestDiffMinutes(t *testing.T) {
	_, fn := diffMinutes()
	time1 := float64(1000 + 60) // 1 minute later
	time2 := float64(1000)

	result, err := fn([]lang.Value{lang.NumberValue(time1), lang.NumberValue(time2)})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	expected := 1.0
	if float64(result.(lang.NumberValue)) != expected {
		t.Errorf("Expected %f, got %f", expected, float64(result.(lang.NumberValue)))
	}
}

func TestTimeComponents(t *testing.T) {
	// Test with a known time: 2023-06-15 14:30:45 UTC
	testTime := time.Date(2023, 6, 15, 14, 30, 45, 0, time.UTC)
	timestamp := float64(testTime.Unix())

	tests := []struct {
		name     string
		fn       func() (string, lang.Function)
		expected float64
	}{
		{"year", year, 2023},
		{"month", month, 6},
		{"day", day, 15},
		{"hour", hour, 14},
		{"minute", minute, 30},
		{"second", second, 45},
		{"weekday", weekday, float64(testTime.Weekday())},
		{"yearday", yearday, float64(testTime.YearDay())},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, fn := tt.fn()
			result, err := fn([]lang.Value{lang.NumberValue(timestamp)})
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if float64(result.(lang.NumberValue)) != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, float64(result.(lang.NumberValue)))
			}
		})
	}
}

func TestWeek(t *testing.T) {
	_, fn := week()

	// Test with a known date
	testTime := time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC)
	timestamp := float64(testTime.Unix())

	result, err := fn([]lang.Value{lang.NumberValue(timestamp)})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	_, expectedWeek := testTime.ISOWeek()
	if float64(result.(lang.NumberValue)) != float64(expectedWeek) {
		t.Errorf("Expected week %d, got %f", expectedWeek, float64(result.(lang.NumberValue)))
	}
}

func TestStartOfDay(t *testing.T) {
	_, fn := startOfDay()

	// Test with 2023-06-15 14:30:45
	testTime := time.Date(2023, 6, 15, 14, 30, 45, 0, time.UTC)
	timestamp := float64(testTime.Unix())

	result, err := fn([]lang.Value{lang.NumberValue(timestamp)})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	expectedStart := time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC)
	expected := float64(expectedStart.Unix())

	if float64(result.(lang.NumberValue)) != expected {
		t.Errorf("Expected %f, got %f", expected, float64(result.(lang.NumberValue)))
	}
}

func TestEndOfDay(t *testing.T) {
	_, fn := toEndOfDay()

	// Test with 2023-06-15 14:30:45
	testTime := time.Date(2023, 6, 15, 14, 30, 45, 0, time.UTC)
	timestamp := float64(testTime.Unix())

	result, err := fn([]lang.Value{lang.NumberValue(timestamp)})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	expectedEnd := time.Date(2023, 6, 15, 23, 59, 59, 999999999, time.UTC)
	expected := float64(expectedEnd.Unix())

	if float64(result.(lang.NumberValue)) != expected {
		t.Errorf("Expected %f, got %f", expected, float64(result.(lang.NumberValue)))
	}
}

func TestStartOfWeek(t *testing.T) {
	_, fn := startOfWeek()

	// Test with a Thursday (2023-06-15)
	testTime := time.Date(2023, 6, 15, 14, 30, 45, 0, time.UTC) // Thursday
	timestamp := float64(testTime.Unix())

	result, err := fn([]lang.Value{lang.NumberValue(timestamp)})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	// Should return Monday of that week (2023-06-12)
	expectedStart := time.Date(2023, 6, 12, 0, 0, 0, 0, time.UTC)
	expected := float64(expectedStart.Unix())

	if float64(result.(lang.NumberValue)) != expected {
		t.Errorf("Expected %f, got %f", expected, float64(result.(lang.NumberValue)))
	}
}

func TestStartOfMonth(t *testing.T) {
	_, fn := startOfMonth()

	testTime := time.Date(2023, 6, 15, 14, 30, 45, 0, time.UTC)
	timestamp := float64(testTime.Unix())

	result, err := fn([]lang.Value{lang.NumberValue(timestamp)})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	expectedStart := time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC)
	expected := float64(expectedStart.Unix())

	if float64(result.(lang.NumberValue)) != expected {
		t.Errorf("Expected %f, got %f", expected, float64(result.(lang.NumberValue)))
	}
}

func TestStartOfYear(t *testing.T) {
	_, fn := startOfYear()

	testTime := time.Date(2023, 6, 15, 14, 30, 45, 0, time.UTC)
	timestamp := float64(testTime.Unix())

	result, err := fn([]lang.Value{lang.NumberValue(timestamp)})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	expectedStart := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	expected := float64(expectedStart.Unix())

	if float64(result.(lang.NumberValue)) != expected {
		t.Errorf("Expected %f, got %f", expected, float64(result.(lang.NumberValue)))
	}
}

func TestIsWeekend(t *testing.T) {
	_, fn := isWeekend()

	tests := []struct {
		name     string
		date     time.Time
		expected bool
	}{
		{"Monday", time.Date(2023, 6, 12, 0, 0, 0, 0, time.UTC), false},
		{"Friday", time.Date(2023, 6, 16, 0, 0, 0, 0, time.UTC), false},
		{"Saturday", time.Date(2023, 6, 17, 0, 0, 0, 0, time.UTC), true},
		{"Sunday", time.Date(2023, 6, 18, 0, 0, 0, 0, time.UTC), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timestamp := float64(tt.date.Unix())
			result, err := fn([]lang.Value{lang.NumberValue(timestamp)})
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if bool(result.(lang.BoolValue)) != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, bool(result.(lang.BoolValue)))
			}
		})
	}
}

func TestIsLeapYear(t *testing.T) {
	_, fn := isLeapYear()

	tests := []struct {
		name     string
		year     int
		expected bool
	}{
		{"regular leap year", 2020, true},
		{"non-leap year", 2021, false},
		{"century non-leap", 1900, false},
		{"century leap", 2000, true},
		{"regular non-leap", 2023, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testTime := time.Date(tt.year, 6, 15, 0, 0, 0, 0, time.UTC)
			timestamp := float64(testTime.Unix())

			result, err := fn([]lang.Value{lang.NumberValue(timestamp)})
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if bool(result.(lang.BoolValue)) != tt.expected {
				t.Errorf("Expected %v for year %d, got %v", tt.expected, tt.year, bool(result.(lang.BoolValue)))
			}
		})
	}
}

func TestDaysInMonth(t *testing.T) {
	_, fn := daysInMonth()

	tests := []struct {
		name     string
		date     time.Time
		expected float64
	}{
		{"January", time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC), 31},
		{"February non-leap", time.Date(2023, 2, 15, 0, 0, 0, 0, time.UTC), 28},
		{"February leap", time.Date(2020, 2, 15, 0, 0, 0, 0, time.UTC), 29},
		{"April", time.Date(2023, 4, 15, 0, 0, 0, 0, time.UTC), 30},
		{"December", time.Date(2023, 12, 15, 0, 0, 0, 0, time.UTC), 31},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timestamp := float64(tt.date.Unix())
			result, err := fn([]lang.Value{lang.NumberValue(timestamp)})
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if float64(result.(lang.NumberValue)) != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, float64(result.(lang.NumberValue)))
			}
		})
	}
}

func TestAge(t *testing.T) {
	_, fn := age()

	// Birth date: 1990-06-15
	birthTime := time.Date(1990, 6, 15, 0, 0, 0, 0, time.UTC)
	birthTimestamp := float64(birthTime.Unix())

	tests := []struct {
		name        string
		currentTime *time.Time
		expected    float64
	}{
		{
			"exact birthday",
			&time.Time{},
			0, // Will be calculated
		},
		{
			"before birthday this year",
			timePtr(time.Date(2023, 6, 14, 0, 0, 0, 0, time.UTC)),
			32,
		},
		{
			"after birthday this year",
			timePtr(time.Date(2023, 6, 16, 0, 0, 0, 0, time.UTC)),
			33,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var args []lang.Value
			args = []lang.Value{lang.NumberValue(birthTimestamp)}

			if tt.name == "exact birthday" {
				// 33 years later exactly
				exactBirthday := time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC)
				args = append(args, lang.NumberValue(float64(exactBirthday.Unix())))
				tt.expected = 33
			} else {
				args = append(args, lang.NumberValue(float64(tt.currentTime.Unix())))
			}

			result, err := fn(args)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if float64(result.(lang.NumberValue)) != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, float64(result.(lang.NumberValue)))
			}
		})
	}
}

func TestTimezone(t *testing.T) {
	_, toTzFn := toTimezone()
	_, fromTzFn := toFromTimezone()

	testTime := time.Date(2023, 6, 15, 12, 0, 0, 0, time.UTC)
	timestamp := float64(testTime.Unix())

	tests := []struct {
		name     string
		timezone string
		hasError bool
	}{
		{"UTC timezone", "UTC", false},
		{"empty timezone", "", false},
		{"America/New_York", "America/New_York", false},
		{"invalid timezone", "Invalid/Timezone", true},
	}

	for _, tt := range tests {
		t.Run("toTimezone_"+tt.name, func(t *testing.T) {
			result, err := toTzFn([]lang.Value{lang.NumberValue(timestamp), lang.StringValue(tt.timezone)})
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// For UTC and empty, should return same timestamp
			if tt.timezone == "UTC" || tt.timezone == "" {
				if float64(result.(lang.NumberValue)) != timestamp {
					t.Errorf("Expected same timestamp for UTC, got %f", float64(result.(lang.NumberValue)))
				}
			}
		})

		t.Run("fromTimezone_"+tt.name, func(t *testing.T) {
			result, err := fromTzFn([]lang.Value{lang.NumberValue(timestamp), lang.StringValue(tt.timezone)})
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// For UTC and empty, should return same timestamp
			if tt.timezone == "UTC" || tt.timezone == "" {
				if float64(result.(lang.NumberValue)) != timestamp {
					t.Errorf("Expected same timestamp for UTC, got %f", float64(result.(lang.NumberValue)))
				}
			}
		})
	}
}

func TestSleep(t *testing.T) {
	_, fn := sleep()

	tests := []struct {
		name     string
		seconds  float64
		hasError bool
	}{
		{"short sleep", 0.001, false}, // 1ms
		{"zero sleep", 0, false},
		{"negative sleep", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := time.Now()
			result, err := fn([]lang.Value{lang.NumberValue(tt.seconds)})
			elapsed := time.Since(start)

			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if !bool(result.(lang.BoolValue)) {
				t.Errorf("Expected true result")
			}

			expectedDuration := time.Duration(tt.seconds * float64(time.Second))
			if elapsed < expectedDuration {
				t.Errorf("Sleep was too short: expected at least %v, got %v", expectedDuration, elapsed)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	_, fn := validate()

	tests := []struct {
		name     string
		timeStr  string
		layout   *string
		expected bool
	}{
		{
			"valid RFC3339",
			"2023-01-15T10:30:00Z",
			nil,
			true,
		},
		{
			"invalid RFC3339",
			"not-a-time",
			nil,
			false,
		},
		{
			"valid custom layout",
			"2023-01-15 10:30:00",
			stringPtr("YYYY-MM-DD HH:mm:ss"),
			true,
		},
		{
			"invalid custom layout",
			"2023-01-15",
			stringPtr("YYYY-MM-DD HH:mm:ss"),
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []lang.Value{lang.StringValue(tt.timeStr)}
			if tt.layout != nil {
				args = append(args, lang.StringValue(*tt.layout))
			}

			result, err := fn(args)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if bool(result.(lang.BoolValue)) != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, bool(result.(lang.BoolValue)))
			}
		})
	}
}

func TestRange(t *testing.T) {
	_, fn := rrange()

	tests := []struct {
		name        string
		start       float64
		end         float64
		step        float64
		expectedLen int
		hasError    bool
	}{
		{"simple range", 1, 5, 1, 5, false},
		{"step of 2", 0, 10, 2, 6, false},
		{"fractional step", 0, 2, 0.5, 5, false},
		{"zero step", 1, 5, 0, 0, true},
		{"negative step", 1, 5, -1, 0, true},
		{"start > end", 5, 1, 1, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []lang.Value{
				lang.NumberValue(tt.start),
				lang.NumberValue(tt.end),
				lang.NumberValue(tt.step),
			}

			result, err := fn(args)
			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			resultList := result.(lang.ListValue)
			if len(resultList) != tt.expectedLen {
				t.Errorf("Expected length %d, got %d", tt.expectedLen, len(resultList))
			}
		})
	}
}
