package types

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/vedadiyan/exql/lang"
	"github.com/vedadiyan/exql/lib"
)

// Basic Type Checking
func typeOf() (string, lang.Function) {
	name := "type"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		switch args[0].(type) {
		case nil:
			return lang.StringValue("null"), nil
		case lang.BoolValue:
			return lang.StringValue("boolean"), nil
		case lang.NumberValue:
			return lang.StringValue("number"), nil
		case lang.StringValue:
			return lang.StringValue("string"), nil
		case lang.ListValue:
			return lang.StringValue("list"), nil
		case lang.MapValue:
			return lang.StringValue("map"), nil
		default:
			return lang.StringValue("unknown"), nil
		}
	}
	return name, fn
}

func isNull() (string, lang.Function) {
	name := "is_null"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		return lang.BoolValue(args[0] == nil), nil
	}
	return name, fn
}

func isDefined() (string, lang.Function) {
	name := "is_defined"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		return lang.BoolValue(args[0] != nil), nil
	}
	return name, fn
}

func isEmpty() (string, lang.Function) {
	name := "is_empty"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		switch v := args[0].(type) {
		case nil:
			return lang.BoolValue(true), nil
		case lang.StringValue:
			return lang.BoolValue(string(v) == ""), nil
		case lang.ListValue:
			return lang.BoolValue(len(v) == 0), nil
		case lang.MapValue:
			return lang.BoolValue(len(v) == 0), nil
		case lang.NumberValue:
			return lang.BoolValue(float64(v) == 0), nil
		case lang.BoolValue:
			return lang.BoolValue(!bool(v)), nil
		default:
			return lang.BoolValue(true), nil
		}
	}
	return name, fn
}

func isNotEmpty() (string, lang.Function) {
	name := "is_not_empty"
	_, isEmpty := isEmpty()
	fn := func(args []lang.Value) (lang.Value, error) {
		result, err := isEmpty(args)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		if boolVal, ok := result.(lang.BoolValue); ok {
			return lang.BoolValue(!bool(boolVal)), nil
		}
		return nil, errors.New("is_not_empty: unexpected result type")
	}
	return name, fn
}

// Primitive Type Checking
func isBool() (string, lang.Function) {
	name := "is_bool"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		_, ok := args[0].(lang.BoolValue)
		return lang.BoolValue(ok), nil
	}
	return name, fn
}

func isNumber() (string, lang.Function) {
	name := "is_number"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		_, ok := args[0].(lang.NumberValue)
		return lang.BoolValue(ok), nil
	}
	return name, fn
}

func isString() (string, lang.Function) {
	name := "is_string"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		_, ok := args[0].(lang.StringValue)
		return lang.BoolValue(ok), nil
	}
	return name, fn
}

func isList() (string, lang.Function) {
	name := "is_list"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		_, ok := args[0].(lang.ListValue)
		return lang.BoolValue(ok), nil
	}
	return name, fn
}

func isMap() (string, lang.Function) {
	name := "is_map"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		_, ok := args[0].(lang.MapValue)
		return lang.BoolValue(ok), nil
	}
	return name, fn
}

func isArray() (string, lang.Function) {
	name := "is_array"
	_, isList := isList()
	fn := func(args []lang.Value) (lang.Value, error) {
		return isList(args)
	}
	return name, fn
}

func isObject() (string, lang.Function) {
	name := "is_object"
	_, isMap := isMap()
	fn := func(args []lang.Value) (lang.Value, error) {
		return isMap(args)
	}
	return name, fn
}

// Number Type Checking
func isInteger() (string, lang.Function) {
	name := "is_integer"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		if num, ok := args[0].(lang.NumberValue); ok {
			val := float64(num)
			return lang.BoolValue(val == math.Trunc(val) && !math.IsInf(val, 0) && !math.IsNaN(val)), nil
		}
		return lang.BoolValue(false), nil
	}
	return name, fn
}

func isFloat() (string, lang.Function) {
	name := "is_float"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		if num, ok := args[0].(lang.NumberValue); ok {
			val := float64(num)
			return lang.BoolValue(val != math.Trunc(val) && !math.IsInf(val, 0) && !math.IsNaN(val)), nil
		}
		return lang.BoolValue(false), nil
	}
	return name, fn
}

func isPositive() (string, lang.Function) {
	name := "is_positive"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		if num, ok := args[0].(lang.NumberValue); ok {
			return lang.BoolValue(float64(num) > 0), nil
		}
		return lang.BoolValue(false), nil
	}
	return name, fn
}

func isNegative() (string, lang.Function) {
	name := "is_negative"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		if num, ok := args[0].(lang.NumberValue); ok {
			return lang.BoolValue(float64(num) < 0), nil
		}
		return lang.BoolValue(false), nil
	}
	return name, fn
}

func isZero() (string, lang.Function) {
	name := "is_zero"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		if num, ok := args[0].(lang.NumberValue); ok {
			return lang.BoolValue(float64(num) == 0), nil
		}
		return lang.BoolValue(false), nil
	}
	return name, fn
}

func isEven() (string, lang.Function) {
	name := "is_even"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		if num, ok := args[0].(lang.NumberValue); ok {
			val := float64(num)
			if val == math.Trunc(val) {
				return lang.BoolValue(int64(val)%2 == 0), nil
			}
		}
		return lang.BoolValue(false), nil
	}
	return name, fn
}

func isOdd() (string, lang.Function) {
	name := "is_odd"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		if num, ok := args[0].(lang.NumberValue); ok {
			val := float64(num)
			if val == math.Trunc(val) {
				return lang.BoolValue(int64(val)%2 != 0), nil
			}
		}
		return lang.BoolValue(false), nil
	}
	return name, fn
}

func isNaN() (string, lang.Function) {
	name := "is_nan"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		if num, ok := args[0].(lang.NumberValue); ok {
			return lang.BoolValue(math.IsNaN(float64(num))), nil
		}
		return lang.BoolValue(false), nil
	}
	return name, fn
}

func isInfinite() (string, lang.Function) {
	name := "is_infinite"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		if num, ok := args[0].(lang.NumberValue); ok {
			return lang.BoolValue(math.IsInf(float64(num), 0)), nil
		}
		return lang.BoolValue(false), nil
	}
	return name, fn
}

func isFinite() (string, lang.Function) {
	name := "is_finite"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		if num, ok := args[0].(lang.NumberValue); ok {
			val := float64(num)
			return lang.BoolValue(!math.IsNaN(val) && !math.IsInf(val, 0)), nil
		}
		return lang.BoolValue(false), nil
	}
	return name, fn
}

// String Type Checking
func isNumericString() (string, lang.Function) {
	name := "is_numeric_string"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		if str, ok := args[0].(lang.StringValue); ok {
			s := strings.TrimSpace(string(str))
			if s == "" {
				return lang.BoolValue(false), nil
			}
			_, err := strconv.ParseFloat(s, 64)
			return lang.BoolValue(err == nil), nil
		}
		return lang.BoolValue(false), nil
	}
	return name, fn
}

func isAlpha() (string, lang.Function) {
	name := "is_alpha"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		if str, ok := args[0].(lang.StringValue); ok {
			s := string(str)
			if s == "" {
				return lang.BoolValue(false), nil
			}
			for _, r := range s {
				if !unicode.IsLetter(r) {
					return lang.BoolValue(false), nil
				}
			}
			return lang.BoolValue(true), nil
		}
		return lang.BoolValue(false), nil
	}
	return name, fn
}

func isAlphaNumeric() (string, lang.Function) {
	name := "is_alphanumeric"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		if str, ok := args[0].(lang.StringValue); ok {
			s := string(str)
			if s == "" {
				return lang.BoolValue(false), nil
			}
			for _, r := range s {
				if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
					return lang.BoolValue(false), nil
				}
			}
			return lang.BoolValue(true), nil
		}
		return lang.BoolValue(false), nil
	}
	return name, fn
}

func isDigit() (string, lang.Function) {
	name := "is_digit"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		if str, ok := args[0].(lang.StringValue); ok {
			s := string(str)
			if s == "" {
				return lang.BoolValue(false), nil
			}
			for _, r := range s {
				if !unicode.IsDigit(r) {
					return lang.BoolValue(false), nil
				}
			}
			return lang.BoolValue(true), nil
		}
		return lang.BoolValue(false), nil
	}
	return name, fn
}

func isLower() (string, lang.Function) {
	name := "is_lower"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		if str, ok := args[0].(lang.StringValue); ok {
			s := string(str)
			if s == "" {
				return lang.BoolValue(false), nil
			}
			return lang.BoolValue(s == strings.ToLower(s) && strings.ContainsAny(s, "abcdefghijklmnopqrstuvwxyz")), nil
		}
		return lang.BoolValue(false), nil
	}
	return name, fn
}

func isUpper() (string, lang.Function) {
	name := "is_upper"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		if str, ok := args[0].(lang.StringValue); ok {
			s := string(str)
			if s == "" {
				return lang.BoolValue(false), nil
			}
			return lang.BoolValue(s == strings.ToUpper(s) && strings.ContainsAny(s, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")), nil
		}
		return lang.BoolValue(false), nil
	}
	return name, fn
}

func isWhitespace() (string, lang.Function) {
	name := "is_whitespace"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		if str, ok := args[0].(lang.StringValue); ok {
			s := string(str)
			if s == "" {
				return lang.BoolValue(false), nil
			}
			for _, r := range s {
				if !unicode.IsSpace(r) {
					return lang.BoolValue(false), nil
				}
			}
			return lang.BoolValue(true), nil
		}
		return lang.BoolValue(false), nil
	}
	return name, fn
}

// Format Validation
func isEmail() (string, lang.Function) {
	name := "is_email"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		if str, ok := args[0].(lang.StringValue); ok {
			s := string(str)
			// Simple email regex
			emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
			matched, err := regexp.MatchString(emailRegex, s)
			if err != nil {
				return nil, fmt.Errorf("%s: regex error: %w", name, err)
			}
			return lang.BoolValue(matched), nil
		}
		return lang.BoolValue(false), nil
	}
	return name, fn
}

func isURL() (string, lang.Function) {
	name := "is_url"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		if str, ok := args[0].(lang.StringValue); ok {
			s := string(str)
			// Simple URL regex
			urlRegex := `^https?://[a-zA-Z0-9.-]+(?:\.[a-zA-Z]{2,})?(?:/.*)?$`
			matched, err := regexp.MatchString(urlRegex, s)
			if err != nil {
				return nil, fmt.Errorf("%s: regex error: %w", name, err)
			}
			return lang.BoolValue(matched), nil
		}
		return lang.BoolValue(false), nil
	}
	return name, fn
}

func isIPAddress() (string, lang.Function) {
	name := "is_ip_address"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		if str, ok := args[0].(lang.StringValue); ok {
			s := string(str)
			// IPv4 regex
			ipv4Regex := `^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`
			// IPv6 regex (simplified)
			ipv6Regex := `^(?:[0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}$|^::1$|^::$`

			ipv4Match, err1 := regexp.MatchString(ipv4Regex, s)
			if err1 != nil {
				return nil, fmt.Errorf("%s: IPv4 regex error: %w", name, err1)
			}
			ipv6Match, err2 := regexp.MatchString(ipv6Regex, s)
			if err2 != nil {
				return nil, fmt.Errorf("%s: IPv6 regex error: %w", name, err2)
			}

			return lang.BoolValue(ipv4Match || ipv6Match), nil
		}
		return lang.BoolValue(false), nil
	}
	return name, fn
}

func isUUID() (string, lang.Function) {
	name := "is_uuid"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		if str, ok := args[0].(lang.StringValue); ok {
			s := string(str)
			// UUID regex
			uuidRegex := `^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`
			matched, err := regexp.MatchString(uuidRegex, s)
			if err != nil {
				return nil, fmt.Errorf("%s: regex error: %w", name, err)
			}
			return lang.BoolValue(matched), nil
		}
		return lang.BoolValue(false), nil
	}
	return name, fn
}

func isJSON() (string, lang.Function) {
	name := "is_json"
	_, isNumericString := isNumericString()
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		if str, ok := args[0].(lang.StringValue); ok {
			s := strings.TrimSpace(string(str))
			if s == "" {
				return lang.BoolValue(false), nil
			}

			// Simple JSON validation - check if it starts/ends with proper brackets/braces
			isValidJSON := (strings.HasPrefix(s, "{") && strings.HasSuffix(s, "}")) ||
				(strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]")) ||
				(strings.HasPrefix(s, "\"") && strings.HasSuffix(s, "\"")) ||
				s == "true" || s == "false" || s == "null"

			if !isValidJSON {
				// Check if it's a valid numeric string
				result, err := isNumericString([]lang.Value{lang.StringValue(s)})
				if err != nil {
					return nil, fmt.Errorf("%s: error checking numeric string: %w", name, err)
				}
				if boolVal, ok := result.(lang.BoolValue); ok {
					isValidJSON = bool(boolVal)
				}
			}

			return lang.BoolValue(isValidJSON), nil
		}
		return lang.BoolValue(false), nil
	}
	return name, fn
}

func isBase64() (string, lang.Function) {
	name := "is_base64"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		if str, ok := args[0].(lang.StringValue); ok {
			s := string(str)
			// Base64 regex
			base64Regex := `^[A-Za-z0-9+/]*={0,2}$`
			matched, err := regexp.MatchString(base64Regex, s)
			if err != nil {
				return nil, fmt.Errorf("%s: regex error: %w", name, err)
			}
			if !matched {
				return lang.BoolValue(false), nil
			}
			// Check length (must be multiple of 4)
			return lang.BoolValue(len(s)%4 == 0), nil
		}
		return lang.BoolValue(false), nil
	}
	return name, fn
}

func isHex() (string, lang.Function) {
	name := "is_hex"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		if str, ok := args[0].(lang.StringValue); ok {
			s := string(str)
			if s == "" {
				return lang.BoolValue(false), nil
			}
			// Hex regex
			hexRegex := `^[0-9a-fA-F]+$`
			matched, err := regexp.MatchString(hexRegex, s)
			if err != nil {
				return nil, fmt.Errorf("%s: regex error: %w", name, err)
			}
			return lang.BoolValue(matched), nil
		}
		return lang.BoolValue(false), nil
	}
	return name, fn
}

// Collection Type Checking
func hasLength() (string, lang.Function) {
	name := "has_length"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		switch args[0].(type) {
		case lang.StringValue, lang.ListValue, lang.MapValue:
			return lang.BoolValue(true), nil
		default:
			return lang.BoolValue(false), nil
		}
	}
	return name, fn
}

// Range Checking
func isInRange() (string, lang.Function) {
	name := "is_in_range"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 3 {
			return nil, lib.ArgumentError(name, 3)
		}
		value, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: value %w", name, err)
		}
		min, err := lib.ToNumber(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: min %w", name, err)
		}
		max, err := lib.ToNumber(args[2])
		if err != nil {
			return nil, fmt.Errorf("%s: max %w", name, err)
		}
		if min > max {
			return nil, errors.New("is_in_range: min cannot be greater than max")
		}
		return lang.BoolValue(value >= min && value <= max), nil
	}
	return name, fn
}

func isLengthInRange() (string, lang.Function) {
	name := "is_length_in_range"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 3 {
			return nil, lib.ArgumentError(name, 3)
		}
		var length float64
		switch v := args[0].(type) {
		case lang.StringValue:
			length = float64(len(string(v)))
		case lang.ListValue:
			length = float64(len(v))
		case lang.MapValue:
			length = float64(len(v))
		default:
			return nil, errors.New("is_length_in_range: argument must be string, list, or map")
		}
		min, err := lib.ToNumber(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: min %w", name, err)
		}
		max, err := lib.ToNumber(args[2])
		if err != nil {
			return nil, fmt.Errorf("%s: max %w", name, err)
		}
		if min < 0 {
			return nil, errors.New("is_length_in_range: min cannot be negative")
		}
		if min > max {
			return nil, errors.New("is_length_in_range: min cannot be greater than max")
		}
		return lang.BoolValue(length >= min && length <= max), nil
	}
	return name, fn
}

// Type Conversion Checking
func canConvertToNumber() (string, lang.Function) {
	name := "can_convert_to_number"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		switch v := args[0].(type) {
		case lang.NumberValue:
			return lang.BoolValue(true), nil
		case lang.StringValue:
			s := strings.TrimSpace(string(v))
			if s == "" {
				return lang.BoolValue(false), nil
			}
			_, err := strconv.ParseFloat(s, 64)
			return lang.BoolValue(err == nil), nil
		case lang.BoolValue:
			return lang.BoolValue(true), nil
		default:
			return lang.BoolValue(false), nil
		}
	}
	return name, fn
}

func canConvertToString() (string, lang.Function) {
	name := "can_convert_to_string"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		switch args[0].(type) {
		case nil:
			return lang.BoolValue(false), nil
		default:
			return lang.BoolValue(true), nil
		}
	}
	return name, fn
}

func canConvertToBool() (string, lang.Function) {
	name := "can_convert_to_bool"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		switch args[0].(type) {
		case nil:
			return lang.BoolValue(false), nil
		default:
			return lang.BoolValue(true), nil
		}
	}
	return name, fn
}

// Comparison Helpers
func deepEqual(a, b lang.Value) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	switch va := a.(type) {
	case lang.BoolValue:
		if vb, ok := b.(lang.BoolValue); ok {
			return bool(va) == bool(vb)
		}
	case lang.NumberValue:
		if vb, ok := b.(lang.NumberValue); ok {
			return float64(va) == float64(vb)
		}
	case lang.StringValue:
		if vb, ok := b.(lang.StringValue); ok {
			return string(va) == string(vb)
		}
	case lang.ListValue:
		if vb, ok := b.(lang.ListValue); ok {
			if len(va) != len(vb) {
				return false
			}
			for i := range va {
				if !deepEqual(va[i], vb[i]) {
					return false
				}
			}
			return true
		}
	case lang.MapValue:
		if vb, ok := b.(lang.MapValue); ok {
			if len(va) != len(vb) {
				return false
			}
			for key, valueA := range va {
				if valueB, exists := vb[key]; exists {
					if !deepEqual(valueA, valueB) {
						return false
					}
				} else {
					return false
				}
			}
			return true
		}
	}
	return false
}

func areEqual() (string, lang.Function) {
	name := "are_equal"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 2 {
			return nil, lib.ArgumentError(name, 2)
		}
		return lang.BoolValue(deepEqual(args[0], args[1])), nil
	}
	return name, fn
}

func areStrictEqual() (string, lang.Function) {
	name := "are_strict_equal"
	_, toTypeOf := typeOf()
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 2 {
			return nil, lib.ArgumentError(name, 2)
		}
		// Check type first
		type1, err := toTypeOf([]lang.Value{args[0]})
		if err != nil {
			return nil, fmt.Errorf("%s: error getting type of first argument: %w", name, err)
		}
		type2, err := toTypeOf([]lang.Value{args[1]})
		if err != nil {
			return nil, fmt.Errorf("%s: error getting type of second argument: %w", name, err)
		}
		if type1 != type2 {
			return lang.BoolValue(false), nil
		}
		return lang.BoolValue(deepEqual(args[0], args[1])), nil
	}
	return name, fn
}

var typeFunctions = []func() (string, lang.Function){
	// Basic type checking
	typeOf,
	isNull,
	isDefined,
	isEmpty,
	isNotEmpty,

	// Primitive types
	isBool,
	isNumber,
	isString,
	isList,
	isMap,
	isArray,  // Alias
	isObject, // Alias

	// Number types
	isInteger,
	isFloat,
	isPositive,
	isNegative,
	isZero,
	isEven,
	isOdd,
	isNaN,
	isInfinite,
	isFinite,

	// String types
	isNumericString,
	isAlpha,
	isAlphaNumeric,
	isDigit,
	isLower,
	isUpper,
	isWhitespace,

	// Format validation
	isEmail,
	isURL,
	isIPAddress,
	isUUID,
	isJSON,
	isBase64,
	isHex,

	// Collection checking
	hasLength,

	// Range checking
	isInRange,
	isLengthInRange,

	// Conversion checking
	canConvertToNumber,
	canConvertToString,
	canConvertToBool,

	// Comparison
	areEqual,
	areStrictEqual,
}

func Export() map[string]lang.Function {
	out := make(map[string]lang.Function)
	for _, value := range typeFunctions {
		name, fn := value()
		out[name] = fn
	}
	return out
}
