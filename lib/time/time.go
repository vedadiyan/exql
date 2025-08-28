package time

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/vedadiyan/exql/lang"
	"github.com/vedadiyan/exql/lib"
)

func now() (string, func([]lang.Value) (lang.Value, error)) {
	name := "now"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 0 {
			return nil, lib.ArgumentError(name, 0)
		}
		return lang.NumberValue(float64(time.Now().Unix())), nil
	}
	return name, fn
}

func nowMillis() (string, func([]lang.Value) (lang.Value, error)) {
	name := "now_millis"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 0 {
			return nil, lib.ArgumentError(name, 0)
		}
		return lang.NumberValue(float64(time.Now().UnixMilli())), nil
	}
	return name, fn
}

func nowNanos() (string, func([]lang.Value) (lang.Value, error)) {
	name := "now_nanos"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 0 {
			return nil, lib.ArgumentError(name, 0)
		}
		return lang.NumberValue(float64(time.Now().UnixNano())), nil
	}
	return name, fn
}

func parse() (string, func([]lang.Value) (lang.Value, error)) {
	name := "parse"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) < 1 || len(args) > 2 {
			return nil, fmt.Errorf("%s: expected 1 or 2 arguments", name)
		}
		timeStr, err := lib.ToString(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: time string %w", name, err)
		}
		layout := time.RFC3339
		if len(args) == 2 {
			layoutStr, err := lib.ToString(args[1])
			if err != nil {
				return nil, fmt.Errorf("%s: layout %w", name, err)
			}
			layout = convertTimeLayout(string(layoutStr))
		}
		t, err := time.Parse(layout, string(timeStr))
		if err != nil {
			return nil, fmt.Errorf("%s: failed to parse time '%s' with layout '%s': %w", name, string(timeStr), layout, err)
		}
		return lang.NumberValue(float64(t.Unix())), nil
	}
	return name, fn
}

func toFormat() (string, func([]lang.Value) (lang.Value, error)) {
	name := "format"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) < 1 || len(args) > 2 {
			return nil, fmt.Errorf("%s: expected 1 or 2 arguments", name)
		}
		timestamp, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: timestamp %w", name, err)
		}
		layout := time.RFC3339
		if len(args) == 2 {
			layoutStr, err := lib.ToString(args[1])
			if err != nil {
				return nil, fmt.Errorf("%s: layout %w", name, err)
			}
			layout = convertTimeLayout(string(layoutStr))
		}
		t := time.Unix(int64(timestamp), 0).UTC()
		return lang.StringValue(t.Format(layout)), nil
	}
	return name, fn
}

func add() (string, func([]lang.Value) (lang.Value, error)) {
	name := "add"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 2 {
			return nil, lib.ArgumentError(name, 2)
		}
		timestamp, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: timestamp %w", name, err)
		}
		seconds, err := lib.ToNumber(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: seconds %w", name, err)
		}
		return lang.NumberValue(timestamp + seconds), nil
	}
	return name, fn
}

func addDays() (string, func([]lang.Value) (lang.Value, error)) {
	name := "add_days"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 2 {
			return nil, lib.ArgumentError(name, 2)
		}
		timestamp, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: timestamp %w", name, err)
		}
		days, err := lib.ToNumber(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: days %w", name, err)
		}
		return lang.NumberValue(timestamp + (days * 86400)), nil
	}
	return name, fn
}

func addHours() (string, func([]lang.Value) (lang.Value, error)) {
	name := "add_hours"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 2 {
			return nil, lib.ArgumentError(name, 2)
		}
		timestamp, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: timestamp %w", name, err)
		}
		hours, err := lib.ToNumber(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: hours %w", name, err)
		}
		return lang.NumberValue(timestamp + (hours * 3600)), nil
	}
	return name, fn
}

func addMinutes() (string, func([]lang.Value) (lang.Value, error)) {
	name := "add_minutes"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 2 {
			return nil, lib.ArgumentError(name, 2)
		}
		timestamp, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: timestamp %w", name, err)
		}
		minutes, err := lib.ToNumber(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: minutes %w", name, err)
		}
		return lang.NumberValue(timestamp + (minutes * 60)), nil
	}
	return name, fn
}

func diff() (string, func([]lang.Value) (lang.Value, error)) {
	name := "diff"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 2 {
			return nil, lib.ArgumentError(name, 2)
		}
		time1, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: first time %w", name, err)
		}
		time2, err := lib.ToNumber(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: second time %w", name, err)
		}
		return lang.NumberValue(time1 - time2), nil
	}
	return name, fn
}

func diffDays() (string, func([]lang.Value) (lang.Value, error)) {
	name := "diff_days"
	_, diff := diff()
	fn := func(args []lang.Value) (lang.Value, error) {
		diff, err := diff(args)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		diffNum, err := lib.ToNumber(diff)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		return lang.NumberValue(diffNum / 86400), nil
	}
	return name, fn
}

func diffHours() (string, func([]lang.Value) (lang.Value, error)) {
	name := "diff_hours"
	_, diff := diff()
	fn := func(args []lang.Value) (lang.Value, error) {
		diff, err := diff(args)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		diffNum, err := lib.ToNumber(diff)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		return lang.NumberValue(diffNum / 3600), nil
	}
	return name, fn
}

func diffMinutes() (string, func([]lang.Value) (lang.Value, error)) {
	name := "diff_minutes"
	_, diff := diff()
	fn := func(args []lang.Value) (lang.Value, error) {
		diff, err := diff(args)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		diffNum, err := lib.ToNumber(diff)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		return lang.NumberValue(diffNum / 60), nil
	}
	return name, fn
}

func year() (string, func([]lang.Value) (lang.Value, error)) {
	name := "year"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		timestamp, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: timestamp %w", name, err)
		}
		t := time.Unix(int64(timestamp), 0).UTC()
		return lang.NumberValue(float64(t.Year())), nil
	}
	return name, fn
}

func month() (string, func([]lang.Value) (lang.Value, error)) {
	name := "month"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		timestamp, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: timestamp %w", name, err)
		}
		t := time.Unix(int64(timestamp), 0).UTC()
		return lang.NumberValue(float64(t.Month())), nil
	}
	return name, fn
}

func day() (string, func([]lang.Value) (lang.Value, error)) {
	name := "day"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		timestamp, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: timestamp %w", name, err)
		}
		t := time.Unix(int64(timestamp), 0).UTC()
		return lang.NumberValue(float64(t.Day())), nil
	}
	return name, fn
}

func hour() (string, func([]lang.Value) (lang.Value, error)) {
	name := "hour"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		timestamp, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: timestamp %w", name, err)
		}
		t := time.Unix(int64(timestamp), 0).UTC()
		return lang.NumberValue(float64(t.Hour())), nil
	}
	return name, fn
}

func minute() (string, func([]lang.Value) (lang.Value, error)) {
	name := "minute"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		timestamp, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: timestamp %w", name, err)
		}
		t := time.Unix(int64(timestamp), 0).UTC()
		return lang.NumberValue(float64(t.Minute())), nil
	}
	return name, fn
}

func second() (string, func([]lang.Value) (lang.Value, error)) {
	name := "second"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		timestamp, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: timestamp %w", name, err)
		}
		t := time.Unix(int64(timestamp), 0).UTC()
		return lang.NumberValue(float64(t.Second())), nil
	}
	return name, fn
}

func weekday() (string, func([]lang.Value) (lang.Value, error)) {
	name := "weekday"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		timestamp, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: timestamp %w", name, err)
		}
		t := time.Unix(int64(timestamp), 0).UTC()
		return lang.NumberValue(float64(t.Weekday())), nil
	}
	return name, fn
}

func yearday() (string, func([]lang.Value) (lang.Value, error)) {
	name := "yearday"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		timestamp, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: timestamp %w", name, err)
		}
		t := time.Unix(int64(timestamp), 0).UTC()
		return lang.NumberValue(float64(t.YearDay())), nil
	}
	return name, fn
}

func week() (string, func([]lang.Value) (lang.Value, error)) {
	name := "week"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		timestamp, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: timestamp %w", name, err)
		}
		t := time.Unix(int64(timestamp), 0).UTC()
		_, week := t.ISOWeek()
		return lang.NumberValue(float64(week)), nil
	}
	return name, fn
}

func startOfDay() (string, func([]lang.Value) (lang.Value, error)) {
	name := "start_of_day"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		timestamp, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: timestamp %w", name, err)
		}
		t := time.Unix(int64(timestamp), 0).UTC()
		startOfDay := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
		return lang.NumberValue(float64(startOfDay.Unix())), nil
	}
	return name, fn
}

func toEndOfDay() (string, func([]lang.Value) (lang.Value, error)) {
	name := "end_of_day"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		timestamp, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: timestamp %w", name, err)
		}
		t := time.Unix(int64(timestamp), 0).UTC()
		endOfDay := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, time.UTC)
		return lang.NumberValue(float64(endOfDay.Unix())), nil
	}
	return name, fn
}

func startOfWeek() (string, func([]lang.Value) (lang.Value, error)) {
	name := "start_of_week"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		timestamp, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: timestamp %w", name, err)
		}
		t := time.Unix(int64(timestamp), 0).UTC()
		days := int(t.Weekday())
		if days == 0 {
			days = 7
		}
		days--
		startOfWeek := t.AddDate(0, 0, -days)
		startOfWeek = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, time.UTC)
		return lang.NumberValue(float64(startOfWeek.Unix())), nil
	}
	return name, fn
}

func startOfMonth() (string, func([]lang.Value) (lang.Value, error)) {
	name := "start_of_month"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		timestamp, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: timestamp %w", name, err)
		}
		t := time.Unix(int64(timestamp), 0).UTC()
		startOfMonth := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
		return lang.NumberValue(float64(startOfMonth.Unix())), nil
	}
	return name, fn
}

func startOfYear() (string, func([]lang.Value) (lang.Value, error)) {
	name := "start_of_year"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		timestamp, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: timestamp %w", name, err)
		}
		t := time.Unix(int64(timestamp), 0).UTC()
		startOfYear := time.Date(t.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
		return lang.NumberValue(float64(startOfYear.Unix())), nil
	}
	return name, fn
}

func isWeekend() (string, func([]lang.Value) (lang.Value, error)) {
	name := "is_weekend"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		timestamp, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: timestamp %w", name, err)
		}
		t := time.Unix(int64(timestamp), 0).UTC()
		weekday := t.Weekday()
		return lang.BoolValue(weekday == time.Saturday || weekday == time.Sunday), nil
	}
	return name, fn
}

func isLeapYear() (string, func([]lang.Value) (lang.Value, error)) {
	name := "is_leap_year"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		timestamp, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: timestamp %w", name, err)
		}
		t := time.Unix(int64(timestamp), 0).UTC()
		year := t.Year()
		return lang.BoolValue((year%4 == 0 && year%100 != 0) || (year%400 == 0)), nil
	}
	return name, fn
}

func daysInMonth() (string, func([]lang.Value) (lang.Value, error)) {
	name := "days_in_month"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		timestamp, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: timestamp %w", name, err)
		}
		t := time.Unix(int64(timestamp), 0).UTC()
		nextMonth := t.AddDate(0, 1, 0)
		firstOfNextMonth := time.Date(nextMonth.Year(), nextMonth.Month(), 1, 0, 0, 0, 0, time.UTC)
		lastOfThisMonth := firstOfNextMonth.AddDate(0, 0, -1)
		return lang.NumberValue(float64(lastOfThisMonth.Day())), nil
	}
	return name, fn
}

func age() (string, func([]lang.Value) (lang.Value, error)) {
	name := "age"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) < 1 || len(args) > 2 {
			return nil, fmt.Errorf("%s: expected 1 or 2 arguments", name)
		}
		birthTimestamp, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: birth timestamp %w", name, err)
		}
		currentTimestamp := float64(time.Now().Unix())
		if len(args) == 2 {
			currentTimestamp, err = lib.ToNumber(args[1])
			if err != nil {
				return nil, fmt.Errorf("%s: current timestamp %w", name, err)
			}
		}
		birth := time.Unix(int64(birthTimestamp), 0).UTC()
		current := time.Unix(int64(currentTimestamp), 0).UTC()
		years := current.Year() - birth.Year()
		if current.Month() < birth.Month() || (current.Month() == birth.Month() && current.Day() < birth.Day()) {
			years--
		}
		return lang.NumberValue(float64(years)), nil
	}
	return name, fn
}

func toTimezone() (string, func([]lang.Value) (lang.Value, error)) {
	name := "to_timezone"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 2 {
			return nil, lib.ArgumentError(name, 2)
		}
		timestamp, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: timestamp %w", name, err)
		}
		timezone, err := lib.ToString(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: timezone %w", name, err)
		}
		t := time.Unix(int64(timestamp), 0).UTC()
		if string(timezone) == "UTC" || string(timezone) == "" {
			return lang.NumberValue(timestamp), nil
		}
		loc, err := time.LoadLocation(string(timezone))
		if err != nil {
			return nil, fmt.Errorf("%s: invalid timezone '%s': %w", name, string(timezone), err)
		}
		localTime := t.In(loc)
		return lang.NumberValue(float64(localTime.Unix())), nil
	}
	return name, fn
}

func toFromTimezone() (string, func([]lang.Value) (lang.Value, error)) {
	name := "from_timezone"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 2 {
			return nil, lib.ArgumentError(name, 2)
		}
		timestamp, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: timestamp %w", name, err)
		}
		timezone, err := lib.ToString(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: timezone %w", name, err)
		}
		if string(timezone) == "UTC" || string(timezone) == "" {
			return lang.NumberValue(timestamp), nil
		}
		loc, err := time.LoadLocation(string(timezone))
		if err != nil {
			return nil, fmt.Errorf("%s: invalid timezone '%s': %w", name, string(timezone), err)
		}
		localTime := time.Unix(int64(timestamp), 0).In(loc)
		utcTime := localTime.UTC()
		return lang.NumberValue(float64(utcTime.Unix())), nil
	}
	return name, fn
}

func sleep() (string, func([]lang.Value) (lang.Value, error)) {
	name := "sleep"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		seconds, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: seconds %w", name, err)
		}
		if seconds < 0 {
			return nil, errors.New("sleep: seconds cannot be negative")
		}
		duration := time.Duration(seconds * float64(time.Second))
		time.Sleep(duration)
		return lang.BoolValue(true), nil
	}
	return name, fn
}

func validate() (string, func([]lang.Value) (lang.Value, error)) {
	name := "validate"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) < 1 || len(args) > 2 {
			return nil, fmt.Errorf("%s: expected 1 or 2 arguments", name)
		}
		timeStr, err := lib.ToString(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: time string %w", name, err)
		}
		layout := time.RFC3339
		if len(args) == 2 {
			layoutStr, err := lib.ToString(args[1])
			if err != nil {
				return nil, fmt.Errorf("%s: layout %w", name, err)
			}
			layout = convertTimeLayout(string(layoutStr))
		}
		_, err = time.Parse(layout, string(timeStr))
		return lang.BoolValue(err == nil), nil
	}
	return name, fn
}

func rrange() (string, func([]lang.Value) (lang.Value, error)) {
	name := "range"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 3 {
			return nil, lib.ArgumentError(name, 3)
		}
		start, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: start %w", name, err)
		}
		end, err := lib.ToNumber(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: end %w", name, err)
		}
		step, err := lib.ToNumber(args[2])
		if err != nil {
			return nil, fmt.Errorf("%s: step %w", name, err)
		}
		if step <= 0 {
			return nil, errors.New("range: step must be positive")
		}
		if start > end {
			return nil, errors.New("range: start must be less than or equal to end")
		}
		var result lang.ListValue
		for current := start; current <= end; current += step {
			result = append(result, lang.NumberValue(current))
		}
		return result, nil
	}
	return name, fn
}

func convertTimeLayout(layout string) string {
	layout = strings.ReplaceAll(layout, "YYYY", "2006")
	layout = strings.ReplaceAll(layout, "YY", "06")
	layout = strings.ReplaceAll(layout, "MM", "01")
	layout = strings.ReplaceAll(layout, "DD", "02")
	layout = strings.ReplaceAll(layout, "HH", "15")
	layout = strings.ReplaceAll(layout, "mm", "04")
	layout = strings.ReplaceAll(layout, "ss", "05")
	layout = strings.ReplaceAll(layout, "SSS", "000")
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

var TimeFunctions = []func() (string, func([]lang.Value) (lang.Value, error)){
	now,
	nowMillis,
	nowNanos,
	parse,
	toFormat,
	add,
	addDays,
	addHours,
	addMinutes,
	diff,
	diffDays,
	diffHours,
	diffMinutes,
	year,
	month,
	day,
	hour,
	minute,
	second,
	weekday,
	yearday,
	week,
	startOfDay,
	toEndOfDay,
	startOfWeek,
	startOfMonth,
	startOfYear,
	isWeekend,
	isLeapYear,
	daysInMonth,
	age,
	toTimezone,
	toFromTimezone,
	sleep,
	validate,
	rrange,
}
