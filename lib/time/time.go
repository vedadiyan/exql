package time

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/vedadiyan/exql/lang"
	"github.com/vedadiyan/exql/lib"
)

// Time/Date Functions
// These functions provide comprehensive time and date manipulation capabilities

func timeNow(args []lang.Value) (lang.Value, error) {
	if len(args) != 0 {
		return nil, errors.New("time_now: expected 0 arguments")
	}
	return lang.NumberValue(float64(time.Now().Unix())), nil
}

func timeNowMillis(args []lang.Value) (lang.Value, error) {
	if len(args) != 0 {
		return nil, errors.New("time_now_millis: expected 0 arguments")
	}
	return lang.NumberValue(float64(time.Now().UnixMilli())), nil
}

func timeNowNanos(args []lang.Value) (lang.Value, error) {
	if len(args) != 0 {
		return nil, errors.New("time_now_nanos: expected 0 arguments")
	}
	return lang.NumberValue(float64(time.Now().UnixNano())), nil
}

func timeParse(args []lang.Value) (lang.Value, error) {
	if len(args) < 1 || len(args) > 2 {
		return nil, errors.New("time_parse: expected 1 or 2 arguments")
	}

	timeStr, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("time_parse: time string %w", err)
	}

	layout := time.RFC3339 // Default layout
	if len(args) == 2 {
		layoutStr, err := lib.ToString(args[1])
		if err != nil {
			return nil, fmt.Errorf("time_parse: layout %w", err)
		}
		layout = convertTimeLayout(string(layoutStr))
	}

	t, err := time.Parse(layout, string(timeStr))
	if err != nil {
		return nil, fmt.Errorf("time_parse: failed to parse time '%s' with layout '%s': %w", string(timeStr), layout, err)
	}
	return lang.NumberValue(float64(t.Unix())), nil
}

func timeFormat(args []lang.Value) (lang.Value, error) {
	if len(args) < 1 || len(args) > 2 {
		return nil, errors.New("time_format: expected 1 or 2 arguments")
	}

	timestamp, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("time_format: timestamp %w", err)
	}

	layout := time.RFC3339 // Default layout
	if len(args) == 2 {
		layoutStr, err := lib.ToString(args[1])
		if err != nil {
			return nil, fmt.Errorf("time_format: layout %w", err)
		}
		layout = convertTimeLayout(string(layoutStr))
	}

	t := time.Unix(int64(timestamp), 0).UTC()
	return lang.StringValue(t.Format(layout)), nil
}

func timeAdd(args []lang.Value) (lang.Value, error) {
	if len(args) != 2 {
		return nil, errors.New("time_add: expected 2 arguments")
	}
	timestamp, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("time_add: timestamp %w", err)
	}
	seconds, err := lib.ToNumber(args[1])
	if err != nil {
		return nil, fmt.Errorf("time_add: seconds %w", err)
	}
	return lang.NumberValue(timestamp + seconds), nil
}

func timeAddDays(args []lang.Value) (lang.Value, error) {
	if len(args) != 2 {
		return nil, errors.New("time_add_days: expected 2 arguments")
	}
	timestamp, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("time_add_days: timestamp %w", err)
	}
	days, err := lib.ToNumber(args[1])
	if err != nil {
		return nil, fmt.Errorf("time_add_days: days %w", err)
	}
	return lang.NumberValue(timestamp + (days * 86400)), nil // 86400 seconds in a day
}

func timeAddHours(args []lang.Value) (lang.Value, error) {
	if len(args) != 2 {
		return nil, errors.New("time_add_hours: expected 2 arguments")
	}
	timestamp, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("time_add_hours: timestamp %w", err)
	}
	hours, err := lib.ToNumber(args[1])
	if err != nil {
		return nil, fmt.Errorf("time_add_hours: hours %w", err)
	}
	return lang.NumberValue(timestamp + (hours * 3600)), nil // 3600 seconds in an hour
}

func timeAddMinutes(args []lang.Value) (lang.Value, error) {
	if len(args) != 2 {
		return nil, errors.New("time_add_minutes: expected 2 arguments")
	}
	timestamp, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("time_add_minutes: timestamp %w", err)
	}
	minutes, err := lib.ToNumber(args[1])
	if err != nil {
		return nil, fmt.Errorf("time_add_minutes: minutes %w", err)
	}
	return lang.NumberValue(timestamp + (minutes * 60)), nil // 60 seconds in a minute
}

func timeDiff(args []lang.Value) (lang.Value, error) {
	if len(args) != 2 {
		return nil, errors.New("time_diff: expected 2 arguments")
	}
	time1, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("time_diff: first time %w", err)
	}
	time2, err := lib.ToNumber(args[1])
	if err != nil {
		return nil, fmt.Errorf("time_diff: second time %w", err)
	}
	return lang.NumberValue(time1 - time2), nil
}

func timeDiffDays(args []lang.Value) (lang.Value, error) {
	diff, err := timeDiff(args)
	if err != nil {
		return nil, fmt.Errorf("time_diff_days: %w", err)
	}
	diffNum, err := lib.ToNumber(diff)
	if err != nil {
		return nil, fmt.Errorf("time_diff_days: %w", err)
	}
	return lang.NumberValue(diffNum / 86400), nil
}

func timeDiffHours(args []lang.Value) (lang.Value, error) {
	diff, err := timeDiff(args)
	if err != nil {
		return nil, fmt.Errorf("time_diff_hours: %w", err)
	}
	diffNum, err := lib.ToNumber(diff)
	if err != nil {
		return nil, fmt.Errorf("time_diff_hours: %w", err)
	}
	return lang.NumberValue(diffNum / 3600), nil
}

func timeDiffMinutes(args []lang.Value) (lang.Value, error) {
	diff, err := timeDiff(args)
	if err != nil {
		return nil, fmt.Errorf("time_diff_minutes: %w", err)
	}
	diffNum, err := lib.ToNumber(diff)
	if err != nil {
		return nil, fmt.Errorf("time_diff_minutes: %w", err)
	}
	return lang.NumberValue(diffNum / 60), nil
}

func timeYear(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("time_year: expected 1 argument")
	}
	timestamp, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("time_year: timestamp %w", err)
	}
	t := time.Unix(int64(timestamp), 0).UTC()
	return lang.NumberValue(float64(t.Year())), nil
}

func timeMonth(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("time_month: expected 1 argument")
	}
	timestamp, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("time_month: timestamp %w", err)
	}
	t := time.Unix(int64(timestamp), 0).UTC()
	return lang.NumberValue(float64(t.Month())), nil
}

func timeDay(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("time_day: expected 1 argument")
	}
	timestamp, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("time_day: timestamp %w", err)
	}
	t := time.Unix(int64(timestamp), 0).UTC()
	return lang.NumberValue(float64(t.Day())), nil
}

func timeHour(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("time_hour: expected 1 argument")
	}
	timestamp, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("time_hour: timestamp %w", err)
	}
	t := time.Unix(int64(timestamp), 0).UTC()
	return lang.NumberValue(float64(t.Hour())), nil
}

func timeMinute(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("time_minute: expected 1 argument")
	}
	timestamp, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("time_minute: timestamp %w", err)
	}
	t := time.Unix(int64(timestamp), 0).UTC()
	return lang.NumberValue(float64(t.Minute())), nil
}

func timeSecond(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("time_second: expected 1 argument")
	}
	timestamp, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("time_second: timestamp %w", err)
	}
	t := time.Unix(int64(timestamp), 0).UTC()
	return lang.NumberValue(float64(t.Second())), nil
}

func timeWeekday(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("time_weekday: expected 1 argument")
	}
	timestamp, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("time_weekday: timestamp %w", err)
	}
	t := time.Unix(int64(timestamp), 0).UTC()
	return lang.NumberValue(float64(t.Weekday())), nil // 0=Sunday, 1=Monday, etc.
}

func timeYearday(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("time_yearday: expected 1 argument")
	}
	timestamp, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("time_yearday: timestamp %w", err)
	}
	t := time.Unix(int64(timestamp), 0).UTC()
	return lang.NumberValue(float64(t.YearDay())), nil
}

func timeWeek(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("time_week: expected 1 argument")
	}
	timestamp, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("time_week: timestamp %w", err)
	}
	t := time.Unix(int64(timestamp), 0).UTC()
	_, week := t.ISOWeek()
	return lang.NumberValue(float64(week)), nil
}

func timeStartOfDay(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("time_start_of_day: expected 1 argument")
	}
	timestamp, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("time_start_of_day: timestamp %w", err)
	}
	t := time.Unix(int64(timestamp), 0).UTC()
	startOfDay := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	return lang.NumberValue(float64(startOfDay.Unix())), nil
}

func timeEndOfDay(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("time_end_of_day: expected 1 argument")
	}
	timestamp, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("time_end_of_day: timestamp %w", err)
	}
	t := time.Unix(int64(timestamp), 0).UTC()
	endOfDay := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, time.UTC)
	return lang.NumberValue(float64(endOfDay.Unix())), nil
}

func timeStartOfWeek(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("time_start_of_week: expected 1 argument")
	}
	timestamp, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("time_start_of_week: timestamp %w", err)
	}
	t := time.Unix(int64(timestamp), 0).UTC()

	// Find Monday of this week
	days := int(t.Weekday())
	if days == 0 {
		days = 7 // Sunday
	}
	days-- // Make Monday = 0

	startOfWeek := t.AddDate(0, 0, -days)
	startOfWeek = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, time.UTC)
	return lang.NumberValue(float64(startOfWeek.Unix())), nil
}

func timeStartOfMonth(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("time_start_of_month: expected 1 argument")
	}
	timestamp, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("time_start_of_month: timestamp %w", err)
	}
	t := time.Unix(int64(timestamp), 0).UTC()
	startOfMonth := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
	return lang.NumberValue(float64(startOfMonth.Unix())), nil
}

func timeStartOfYear(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("time_start_of_year: expected 1 argument")
	}
	timestamp, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("time_start_of_year: timestamp %w", err)
	}
	t := time.Unix(int64(timestamp), 0).UTC()
	startOfYear := time.Date(t.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
	return lang.NumberValue(float64(startOfYear.Unix())), nil
}

func timeIsWeekend(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("time_is_weekend: expected 1 argument")
	}
	timestamp, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("time_is_weekend: timestamp %w", err)
	}
	t := time.Unix(int64(timestamp), 0).UTC()
	weekday := t.Weekday()
	return lang.BoolValue(weekday == time.Saturday || weekday == time.Sunday), nil
}

func timeIsLeapYear(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("time_is_leap_year: expected 1 argument")
	}
	timestamp, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("time_is_leap_year: timestamp %w", err)
	}
	t := time.Unix(int64(timestamp), 0).UTC()
	year := t.Year()
	return lang.BoolValue((year%4 == 0 && year%100 != 0) || (year%400 == 0)), nil
}

func timeDaysInMonth(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("time_days_in_month: expected 1 argument")
	}
	timestamp, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("time_days_in_month: timestamp %w", err)
	}
	t := time.Unix(int64(timestamp), 0).UTC()

	// Get first day of next month, then subtract one day
	nextMonth := t.AddDate(0, 1, 0)
	firstOfNextMonth := time.Date(nextMonth.Year(), nextMonth.Month(), 1, 0, 0, 0, 0, time.UTC)
	lastOfThisMonth := firstOfNextMonth.AddDate(0, 0, -1)

	return lang.NumberValue(float64(lastOfThisMonth.Day())), nil
}

func timeAge(args []lang.Value) (lang.Value, error) {
	if len(args) < 1 || len(args) > 2 {
		return nil, errors.New("time_age: expected 1 or 2 arguments")
	}

	birthTimestamp, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("time_age: birth timestamp %w", err)
	}

	currentTimestamp := float64(time.Now().Unix())
	if len(args) == 2 {
		currentTimestamp, err = lib.ToNumber(args[1])
		if err != nil {
			return nil, fmt.Errorf("time_age: current timestamp %w", err)
		}
	}

	birth := time.Unix(int64(birthTimestamp), 0).UTC()
	current := time.Unix(int64(currentTimestamp), 0).UTC()

	years := current.Year() - birth.Year()

	// Adjust if birthday hasn't occurred this year
	if current.Month() < birth.Month() ||
		(current.Month() == birth.Month() && current.Day() < birth.Day()) {
		years--
	}

	return lang.NumberValue(float64(years)), nil
}

func timeToTimezone(args []lang.Value) (lang.Value, error) {
	if len(args) != 2 {
		return nil, errors.New("time_to_timezone: expected 2 arguments")
	}

	timestamp, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("time_to_timezone: timestamp %w", err)
	}

	timezone, err := lib.ToString(args[1])
	if err != nil {
		return nil, fmt.Errorf("time_to_timezone: timezone %w", err)
	}

	t := time.Unix(int64(timestamp), 0).UTC()

	if string(timezone) == "UTC" || string(timezone) == "" {
		return lang.NumberValue(timestamp), nil
	}

	loc, err := time.LoadLocation(string(timezone))
	if err != nil {
		return nil, fmt.Errorf("time_to_timezone: invalid timezone '%s': %w", string(timezone), err)
	}

	localTime := t.In(loc)
	return lang.NumberValue(float64(localTime.Unix())), nil
}

func timeFromTimezone(args []lang.Value) (lang.Value, error) {
	if len(args) != 2 {
		return nil, errors.New("time_from_timezone: expected 2 arguments")
	}

	timestamp, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("time_from_timezone: timestamp %w", err)
	}

	timezone, err := lib.ToString(args[1])
	if err != nil {
		return nil, fmt.Errorf("time_from_timezone: timezone %w", err)
	}

	if string(timezone) == "UTC" || string(timezone) == "" {
		return lang.NumberValue(timestamp), nil
	}

	loc, err := time.LoadLocation(string(timezone))
	if err != nil {
		return nil, fmt.Errorf("time_from_timezone: invalid timezone '%s': %w", string(timezone), err)
	}

	localTime := time.Unix(int64(timestamp), 0).In(loc)
	utcTime := localTime.UTC()
	return lang.NumberValue(float64(utcTime.Unix())), nil
}

func timeSleep(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("time_sleep: expected 1 argument")
	}

	seconds, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("time_sleep: seconds %w", err)
	}

	if seconds < 0 {
		return nil, errors.New("time_sleep: seconds cannot be negative")
	}

	duration := time.Duration(seconds * float64(time.Second))
	time.Sleep(duration)
	return lang.BoolValue(true), nil
}

func timeValidate(args []lang.Value) (lang.Value, error) {
	if len(args) < 1 || len(args) > 2 {
		return nil, errors.New("time_validate: expected 1 or 2 arguments")
	}

	timeStr, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("time_validate: time string %w", err)
	}

	layout := time.RFC3339
	if len(args) == 2 {
		layoutStr, err := lib.ToString(args[1])
		if err != nil {
			return nil, fmt.Errorf("time_validate: layout %w", err)
		}
		layout = convertTimeLayout(string(layoutStr))
	}

	_, err = time.Parse(layout, string(timeStr))
	return lang.BoolValue(err == nil), nil
}

func timeRange(args []lang.Value) (lang.Value, error) {
	if len(args) != 3 {
		return nil, errors.New("time_range: expected 3 arguments")
	}

	start, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("time_range: start %w", err)
	}

	end, err := lib.ToNumber(args[1])
	if err != nil {
		return nil, fmt.Errorf("time_range: end %w", err)
	}

	step, err := lib.ToNumber(args[2])
	if err != nil {
		return nil, fmt.Errorf("time_range: step %w", err)
	}

	if step <= 0 {
		return nil, errors.New("time_range: step must be positive")
	}

	if start > end {
		return nil, errors.New("time_range: start must be less than or equal to end")
	}

	var result lang.ListValue
	for current := start; current <= end; current += step {
		result = append(result, lang.NumberValue(current))
	}

	return result, nil
}

// Helper function to convert common time format patterns to Go's format
func convertTimeLayout(layout string) string {
	// Convert common patterns to Go's reference time format
	layout = strings.ReplaceAll(layout, "YYYY", "2006")
	layout = strings.ReplaceAll(layout, "YY", "06")
	layout = strings.ReplaceAll(layout, "MM", "01")
	layout = strings.ReplaceAll(layout, "DD", "02")
	layout = strings.ReplaceAll(layout, "HH", "15")
	layout = strings.ReplaceAll(layout, "mm", "04")
	layout = strings.ReplaceAll(layout, "ss", "05")
	layout = strings.ReplaceAll(layout, "SSS", "000")

	// Common formats
	switch layout {
	case "ISO8601":
		return time.RFC3339
	case "RFC3339":
		return time.RFC3339
	case "RFC822":
		return time.RFC822
	case "RFC850":
		return time.RFC850
	case "RFC1123":
		return time.RFC1123
	case "RFC1123Z":
		return time.RFC1123Z
	case "RFC3339Nano":
		return time.RFC3339Nano
	case "Kitchen":
		return time.Kitchen
	case "Stamp":
		return time.Stamp
	case "StampMilli":
		return time.StampMilli
	case "StampMicro":
		return time.StampMicro
	case "StampNano":
		return time.StampNano
	default:
		return layout
	}
}

// Functions that would be in the BuiltinFunctions map:
var TimeFunctions = map[string]func([]lang.Value) (lang.Value, error){
	// Current time
	"time_now":        timeNow,
	"time_now_millis": timeNowMillis,
	"time_now_nanos":  timeNowNanos,

	// Parsing and formatting
	"time_parse":    timeParse,
	"time_format":   timeFormat,
	"time_validate": timeValidate,

	// Arithmetic
	"time_add":          timeAdd,
	"time_add_days":     timeAddDays,
	"time_add_hours":    timeAddHours,
	"time_add_minutes":  timeAddMinutes,
	"time_diff":         timeDiff,
	"time_diff_days":    timeDiffDays,
	"time_diff_hours":   timeDiffHours,
	"time_diff_minutes": timeDiffMinutes,

	// Component extraction
	"time_year":    timeYear,
	"time_month":   timeMonth,
	"time_day":     timeDay,
	"time_hour":    timeHour,
	"time_minute":  timeMinute,
	"time_second":  timeSecond,
	"time_weekday": timeWeekday,
	"time_yearday": timeYearday,
	"time_week":    timeWeek,

	// Time boundaries
	"time_start_of_day":   timeStartOfDay,
	"time_end_of_day":     timeEndOfDay,
	"time_start_of_week":  timeStartOfWeek,
	"time_start_of_month": timeStartOfMonth,
	"time_start_of_year":  timeStartOfYear,

	// Time properties
	"time_is_weekend":    timeIsWeekend,
	"time_is_leap_year":  timeIsLeapYear,
	"time_days_in_month": timeDaysInMonth,
	"time_age":           timeAge,

	// Timezone operations
	"time_to_timezone":   timeToTimezone,
	"time_from_timezone": timeFromTimezone,

	// Utilities
	"time_sleep": timeSleep,
	"time_range": timeRange,
}
