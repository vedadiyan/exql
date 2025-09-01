package lib

import (
	"fmt"
	"strconv"

	"github.com/vedadiyan/exql/lang"
)

// Helper functions - these need to be implemented based on your lang package
func ToString(v lang.Value) (lang.StringValue, error) {
	switch val := v.(type) {
	case string:
		{
			return lang.StringValue(val), nil
		}
	case lang.StringValue:
		return val, nil
	case lang.NumberValue:
		if val == lang.NumberValue(int64(val)) {
			return lang.StringValue(strconv.FormatInt(int64(val), 10)), nil
		}
		return lang.StringValue(strconv.FormatFloat(float64(val), 'g', -1, 64)), nil
	case lang.BoolValue:
		if bool(val) {
			return lang.StringValue("true"), nil
		}
		return lang.StringValue("false"), nil
	case nil:
		return lang.StringValue(""), nil
	default:
		return lang.StringValue(""), fmt.Errorf("cannot convert %T to string", v)
	}
}

func ToNumber(v lang.Value) (float64, error) {
	switch val := v.(type) {
	case int:
		{
			return float64(val), nil
		}
	case int16:
		{
			return float64(val), nil
		}
	case int32:
		{
			return float64(val), nil
		}
	case int64:
		{
			return float64(val), nil
		}
	case int8:
		{
			return float64(val), nil
		}
	case uint:
		{
			return float64(val), nil
		}
	case uint16:
		{
			return float64(val), nil
		}
	case uint32:
		{
			return float64(val), nil
		}
	case uint64:
		{
			return float64(val), nil
		}
	case uint8:
		{
			return float64(val), nil
		}
	case float32:
		{
			return float64(val), nil
		}
	case float64:
		{
			return val, nil
		}
	case lang.NumberValue:
		return float64(val), nil
	case lang.StringValue:
		if f, err := strconv.ParseFloat(string(val), 64); err == nil {
			return f, nil
		} else {
			return 0, fmt.Errorf("cannot convert string '%s' to number: %w", string(val), err)
		}
	case lang.BoolValue:
		if bool(val) {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("cannot convert %T to number", v)
	}
}

func ToBool(v lang.Value) (bool, error) {
	switch val := v.(type) {
	case bool:
		{
			return val, nil
		}
	case lang.BoolValue:
		return bool(val), nil
	case lang.NumberValue:
		return float64(val) != 0, nil
	case lang.StringValue:
		return string(val) != "", nil
	case lang.ListValue:
		return len(val) > 0, nil
	case lang.MapValue:
		return len(val) > 0, nil
	case nil:
		return false, nil
	default:
		return false, fmt.Errorf("cannot convert %T to bool", v)
	}
}
