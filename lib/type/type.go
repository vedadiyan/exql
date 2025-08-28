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

// Type Checking Functions
// These functions provide comprehensive type validation and checking capabilities

// Basic Type Checking
func typeOf(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("type: expected 1 argument")
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

func isNull(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("is_null: expected 1 argument")
	}
	return lang.BoolValue(args[0] == nil), nil
}

func isDefined(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("is_defined: expected 1 argument")
	}
	return lang.BoolValue(args[0] != nil), nil
}

func isEmpty(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("is_empty: expected 1 argument")
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

func isNotEmpty(args []lang.Value) (lang.Value, error) {
	result, err := isEmpty(args)
	if err != nil {
		return nil, fmt.Errorf("is_not_empty: %w", err)
	}
	if boolVal, ok := result.(lang.BoolValue); ok {
		return lang.BoolValue(!bool(boolVal)), nil
	}
	return nil, errors.New("is_not_empty: unexpected result type")
}

// Primitive Type Checking
func isBool(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("is_bool: expected 1 argument")
	}
	_, ok := args[0].(lang.BoolValue)
	return lang.BoolValue(ok), nil
}

func isNumber(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("is_number: expected 1 argument")
	}
	_, ok := args[0].(lang.NumberValue)
	return lang.BoolValue(ok), nil
}

func isString(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("is_string: expected 1 argument")
	}
	_, ok := args[0].(lang.StringValue)
	return lang.BoolValue(ok), nil
}

func isList(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("is_list: expected 1 argument")
	}
	_, ok := args[0].(lang.ListValue)
	return lang.BoolValue(ok), nil
}

func isMap(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("is_map: expected 1 argument")
	}
	_, ok := args[0].(lang.MapValue)
	return lang.BoolValue(ok), nil
}

// Number Type Checking
func isInteger(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("is_integer: expected 1 argument")
	}
	if num, ok := args[0].(lang.NumberValue); ok {
		val := float64(num)
		return lang.BoolValue(val == math.Trunc(val) && !math.IsInf(val, 0) && !math.IsNaN(val)), nil
	}
	return lang.BoolValue(false), nil
}

func isFloat(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("is_float: expected 1 argument")
	}
	if num, ok := args[0].(lang.NumberValue); ok {
		val := float64(num)
		return lang.BoolValue(val != math.Trunc(val) && !math.IsInf(val, 0) && !math.IsNaN(val)), nil
	}
	return lang.BoolValue(false), nil
}

func isPositive(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("is_positive: expected 1 argument")
	}
	if num, ok := args[0].(lang.NumberValue); ok {
		return lang.BoolValue(float64(num) > 0), nil
	}
	return lang.BoolValue(false), nil
}

func isNegative(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("is_negative: expected 1 argument")
	}
	if num, ok := args[0].(lang.NumberValue); ok {
		return lang.BoolValue(float64(num) < 0), nil
	}
	return lang.BoolValue(false), nil
}

func isZero(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("is_zero: expected 1 argument")
	}
	if num, ok := args[0].(lang.NumberValue); ok {
		return lang.BoolValue(float64(num) == 0), nil
	}
	return lang.BoolValue(false), nil
}

func isEven(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("is_even: expected 1 argument")
	}
	if num, ok := args[0].(lang.NumberValue); ok {
		val := float64(num)
		if val == math.Trunc(val) {
			return lang.BoolValue(int64(val)%2 == 0), nil
		}
	}
	return lang.BoolValue(false), nil
}

func isOdd(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("is_odd: expected 1 argument")
	}
	if num, ok := args[0].(lang.NumberValue); ok {
		val := float64(num)
		if val == math.Trunc(val) {
			return lang.BoolValue(int64(val)%2 != 0), nil
		}
	}
	return lang.BoolValue(false), nil
}

func isNaN(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("is_nan: expected 1 argument")
	}
	if num, ok := args[0].(lang.NumberValue); ok {
		return lang.BoolValue(math.IsNaN(float64(num))), nil
	}
	return lang.BoolValue(false), nil
}

func isInfinite(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("is_infinite: expected 1 argument")
	}
	if num, ok := args[0].(lang.NumberValue); ok {
		return lang.BoolValue(math.IsInf(float64(num), 0)), nil
	}
	return lang.BoolValue(false), nil
}

func isFinite(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("is_finite: expected 1 argument")
	}
	if num, ok := args[0].(lang.NumberValue); ok {
		val := float64(num)
		return lang.BoolValue(!math.IsNaN(val) && !math.IsInf(val, 0)), nil
	}
	return lang.BoolValue(false), nil
}

// String Type Checking
func isNumericString(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("is_numeric_string: expected 1 argument")
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

func isAlpha(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("is_alpha: expected 1 argument")
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

func isAlphaNumeric(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("is_alphanumeric: expected 1 argument")
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

func isDigit(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("is_digit: expected 1 argument")
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

func isLower(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("is_lower: expected 1 argument")
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

func isUpper(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("is_upper: expected 1 argument")
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

func isWhitespace(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("is_whitespace: expected 1 argument")
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

// Format Validation
func isEmail(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("is_email: expected 1 argument")
	}
	if str, ok := args[0].(lang.StringValue); ok {
		s := string(str)
		// Simple email regex
		emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
		matched, err := regexp.MatchString(emailRegex, s)
		if err != nil {
			return nil, fmt.Errorf("is_email: regex error: %w", err)
		}
		return lang.BoolValue(matched), nil
	}
	return lang.BoolValue(false), nil
}

func isURL(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("is_url: expected 1 argument")
	}
	if str, ok := args[0].(lang.StringValue); ok {
		s := string(str)
		// Simple URL regex
		urlRegex := `^https?://[a-zA-Z0-9.-]+(?:\.[a-zA-Z]{2,})?(?:/.*)?$`
		matched, err := regexp.MatchString(urlRegex, s)
		if err != nil {
			return nil, fmt.Errorf("is_url: regex error: %w", err)
		}
		return lang.BoolValue(matched), nil
	}
	return lang.BoolValue(false), nil
}

func isIPAddress(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("is_ip_address: expected 1 argument")
	}
	if str, ok := args[0].(lang.StringValue); ok {
		s := string(str)
		// IPv4 regex
		ipv4Regex := `^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`
		// IPv6 regex (simplified)
		ipv6Regex := `^(?:[0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}$|^::1$|^::$`

		ipv4Match, err1 := regexp.MatchString(ipv4Regex, s)
		if err1 != nil {
			return nil, fmt.Errorf("is_ip_address: IPv4 regex error: %w", err1)
		}
		ipv6Match, err2 := regexp.MatchString(ipv6Regex, s)
		if err2 != nil {
			return nil, fmt.Errorf("is_ip_address: IPv6 regex error: %w", err2)
		}

		return lang.BoolValue(ipv4Match || ipv6Match), nil
	}
	return lang.BoolValue(false), nil
}

func isUUID(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("is_uuid: expected 1 argument")
	}
	if str, ok := args[0].(lang.StringValue); ok {
		s := string(str)
		// UUID regex
		uuidRegex := `^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`
		matched, err := regexp.MatchString(uuidRegex, s)
		if err != nil {
			return nil, fmt.Errorf("is_uuid: regex error: %w", err)
		}
		return lang.BoolValue(matched), nil
	}
	return lang.BoolValue(false), nil
}

func isJSON(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("is_json: expected 1 argument")
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
				return nil, fmt.Errorf("is_json: error checking numeric string: %w", err)
			}
			if boolVal, ok := result.(lang.BoolValue); ok {
				isValidJSON = bool(boolVal)
			}
		}

		return lang.BoolValue(isValidJSON), nil
	}
	return lang.BoolValue(false), nil
}

func isBase64(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("is_base64: expected 1 argument")
	}
	if str, ok := args[0].(lang.StringValue); ok {
		s := string(str)
		// Base64 regex
		base64Regex := `^[A-Za-z0-9+/]*={0,2}$`
		matched, err := regexp.MatchString(base64Regex, s)
		if err != nil {
			return nil, fmt.Errorf("is_base64: regex error: %w", err)
		}
		if !matched {
			return lang.BoolValue(false), nil
		}
		// Check length (must be multiple of 4)
		return lang.BoolValue(len(s)%4 == 0), nil
	}
	return lang.BoolValue(false), nil
}

func isHex(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("is_hex: expected 1 argument")
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
			return nil, fmt.Errorf("is_hex: regex error: %w", err)
		}
		return lang.BoolValue(matched), nil
	}
	return lang.BoolValue(false), nil
}

// Collection Type Checking
func isArray(args []lang.Value) (lang.Value, error) {
	return isList(args) // Alias for consistency
}

func isObject(args []lang.Value) (lang.Value, error) {
	return isMap(args) // Alias for consistency
}

func hasLength(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("has_length: expected 1 argument")
	}

	switch args[0].(type) {
	case lang.StringValue, lang.ListValue, lang.MapValue:
		return lang.BoolValue(true), nil
	default:
		return lang.BoolValue(false), nil
	}
}

// Range Checking
func isInRange(args []lang.Value) (lang.Value, error) {
	if len(args) != 3 {
		return nil, errors.New("is_in_range: expected 3 arguments")
	}

	value, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("is_in_range: value %w", err)
	}

	min, err := lib.ToNumber(args[1])
	if err != nil {
		return nil, fmt.Errorf("is_in_range: min %w", err)
	}

	max, err := lib.ToNumber(args[2])
	if err != nil {
		return nil, fmt.Errorf("is_in_range: max %w", err)
	}

	if min > max {
		return nil, errors.New("is_in_range: min cannot be greater than max")
	}

	return lang.BoolValue(value >= min && value <= max), nil
}

func isLengthInRange(args []lang.Value) (lang.Value, error) {
	if len(args) != 3 {
		return nil, errors.New("is_length_in_range: expected 3 arguments")
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
		return nil, fmt.Errorf("is_length_in_range: min %w", err)
	}

	max, err := lib.ToNumber(args[2])
	if err != nil {
		return nil, fmt.Errorf("is_length_in_range: max %w", err)
	}

	if min < 0 {
		return nil, errors.New("is_length_in_range: min cannot be negative")
	}

	if min > max {
		return nil, errors.New("is_length_in_range: min cannot be greater than max")
	}

	return lang.BoolValue(length >= min && length <= max), nil
}

// Type Conversion Checking
func canConvertToNumber(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("can_convert_to_number: expected 1 argument")
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

func canConvertToString(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("can_convert_to_string: expected 1 argument")
	}

	// Everything except nil can be converted to string in some way
	switch args[0].(type) {
	case nil:
		return lang.BoolValue(false), nil
	default:
		return lang.BoolValue(true), nil
	}
}

func canConvertToBool(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("can_convert_to_bool: expected 1 argument")
	}

	// Everything has a boolean representation
	switch args[0].(type) {
	case nil:
		return lang.BoolValue(false), nil
	default:
		return lang.BoolValue(true), nil
	}
}

// Comparison Helpers
func areEqual(args []lang.Value) (lang.Value, error) {
	if len(args) != 2 {
		return nil, errors.New("are_equal: expected 2 arguments")
	}

	return lang.BoolValue(deepEqual(args[0], args[1])), nil
}

func areStrictEqual(args []lang.Value) (lang.Value, error) {
	if len(args) != 2 {
		return nil, errors.New("are_strict_equal: expected 2 arguments")
	}

	// Check type first
	type1, err := typeOf([]lang.Value{args[0]})
	if err != nil {
		return nil, fmt.Errorf("are_strict_equal: error getting type of first argument: %w", err)
	}

	type2, err := typeOf([]lang.Value{args[1]})
	if err != nil {
		return nil, fmt.Errorf("are_strict_equal: error getting type of second argument: %w", err)
	}

	if type1 != type2 {
		return lang.BoolValue(false), nil
	}

	return lang.BoolValue(deepEqual(args[0], args[1])), nil
}

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

// Functions that would be in the BuiltinFunctions map:
var TypeFunctions = map[string]func([]lang.Value) (lang.Value, error){
	// Basic type checking
	"type":         typeOf,
	"is_null":      isNull,
	"is_defined":   isDefined,
	"is_empty":     isEmpty,
	"is_not_empty": isNotEmpty,

	// Primitive types
	"is_bool":   isBool,
	"is_number": isNumber,
	"is_string": isString,
	"is_list":   isList,
	"is_map":    isMap,
	"is_array":  isArray,  // Alias
	"is_object": isObject, // Alias

	// Number types
	"is_integer":  isInteger,
	"is_float":    isFloat,
	"is_positive": isPositive,
	"is_negative": isNegative,
	"is_zero":     isZero,
	"is_even":     isEven,
	"is_odd":      isOdd,
	"is_nan":      isNaN,
	"is_infinite": isInfinite,
	"is_finite":   isFinite,

	// String types
	"is_numeric_string": isNumericString,
	"is_alpha":          isAlpha,
	"is_alphanumeric":   isAlphaNumeric,
	"is_digit":          isDigit,
	"is_lower":          isLower,
	"is_upper":          isUpper,
	"is_whitespace":     isWhitespace,

	// Format validation
	"is_email":      isEmail,
	"is_url":        isURL,
	"is_ip_address": isIPAddress,
	"is_uuid":       isUUID,
	"is_json":       isJSON,
	"is_base64":     isBase64,
	"is_hex":        isHex,

	// Collection checking
	"has_length": hasLength,

	// Range checking
	"is_in_range":        isInRange,
	"is_length_in_range": isLengthInRange,

	// Conversion checking
	"can_convert_to_number": canConvertToNumber,
	"can_convert_to_string": canConvertToString,
	"can_convert_to_bool":   canConvertToBool,

	// Comparison
	"are_equal":        areEqual,
	"are_strict_equal": areStrictEqual,
}
