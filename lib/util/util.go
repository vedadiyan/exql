package util

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/vedadiyan/exql/lang"
	"github.com/vedadiyan/exql/lib"
)

// Utility/Control Functions
// These functions provide general utility operations and control flow capabilities

// Conditional Functions
func conditionalIf(args []lang.Value) (lang.Value, error) {
	if len(args) < 2 || len(args) > 3 {
		return nil, errors.New("if: expected 2 or 3 arguments")
	}

	condition, err := lib.ToBool(args[0])
	if err != nil {
		return nil, fmt.Errorf("if: condition argument: %w", err)
	}

	trueValue := args[1]

	if condition {
		return trueValue, nil
	}

	if len(args) == 3 {
		return args[2], nil // false value
	}

	return nil, nil
}

func conditionalUnless(args []lang.Value) (lang.Value, error) {
	if len(args) < 2 || len(args) > 3 {
		return nil, errors.New("unless: expected 2 or 3 arguments")
	}

	condition, err := lib.ToBool(args[0])
	if err != nil {
		return nil, fmt.Errorf("unless: condition argument: %w", err)
	}

	trueValue := args[1]

	if !condition {
		return trueValue, nil
	}

	if len(args) == 3 {
		return args[2], nil // false value
	}

	return nil, nil
}

func conditionalSwitch(args []lang.Value) (lang.Value, error) {
	if len(args) < 3 || len(args)%2 == 0 {
		return nil, errors.New("switch: expected an odd number of arguments (value, case1, result1, ... , default)")
	}

	testValue := args[0]

	// Check each case pair
	for i := 1; i < len(args)-1; i += 2 {
		caseValue := args[i]
		result := args[i+1]

		if deepEqual(testValue, caseValue) {
			return result, nil
		}
	}

	// Return default value (last argument)
	return args[len(args)-1], nil
}

// Null Coalescing Functions
func coalesce(args []lang.Value) (lang.Value, error) {
	for _, arg := range args {
		if arg != nil && !isNull(arg) {
			return arg, nil
		}
	}
	return nil, nil
}

func defaultValue(args []lang.Value) (lang.Value, error) {
	if len(args) != 2 {
		return nil, errors.New("default: expected 2 arguments")
	}

	value := args[0]
	defaultVal := args[1]

	if value == nil || isNull(value) {
		return defaultVal, nil
	}

	return value, nil
}

func firstNonNull(args []lang.Value) (lang.Value, error) {
	return coalesce(args)
}

func firstNonEmpty(args []lang.Value) (lang.Value, error) {
	for _, arg := range args {
		if arg != nil && !isEmpty(arg) {
			return arg, nil
		}
	}
	return nil, nil
}

// Comparison and Selection
func greatest(args []lang.Value) (lang.Value, error) {
	if len(args) == 0 {
		return nil, nil
	}

	max := args[0]
	for i := 1; i < len(args); i++ {
		comparison, err := compare(args[i], max)
		if err != nil {
			return nil, fmt.Errorf("greatest: comparison failed: %w", err)
		}
		if comparison > 0 {
			max = args[i]
		}
	}
	return max, nil
}

func least(args []lang.Value) (lang.Value, error) {
	if len(args) == 0 {
		return nil, nil
	}

	min := args[0]
	for i := 1; i < len(args); i++ {
		comparison, err := compare(args[i], min)
		if err != nil {
			return nil, fmt.Errorf("least: comparison failed: %w", err)
		}
		if comparison < 0 {
			min = args[i]
		}
	}
	return min, nil
}

func choose(args []lang.Value) (lang.Value, error) {
	if len(args) < 2 {
		return nil, errors.New("choose: expected at least 2 arguments")
	}

	index, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("choose: index argument: %w", err)
	}

	intIndex := int(index)
	if intIndex < 1 || intIndex > len(args)-1 {
		return nil, errors.New("choose: index out of bounds")
	}

	return args[intIndex], nil // 1-based indexing
}

// Debugging and Inspection
func debugPrint(args []lang.Value) (lang.Value, error) {
	var parts []string
	for _, arg := range args {
		parts = append(parts, formatValueForDebug(arg))
	}

	message := strings.Join(parts, " ")
	fmt.Println("[DEBUG]", message)

	// Return the first argument (or nil if no args)
	if len(args) > 0 {
		return args[0], nil
	}
	return nil, nil
}

func inspect(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("inspect: expected 1 argument")
	}

	return lang.StringValue(formatValueForInspect(args[0])), nil
}

func dump(args []lang.Value) (lang.Value, error) {
	var parts []string
	for _, arg := range args {
		parts = append(parts, formatValueForInspect(arg))
	}

	return lang.StringValue(strings.Join(parts, "\n")), nil
}

// Identity and Pass-through
func identity(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("identity: expected 1 argument")
	}
	return args[0], nil
}

func noop(args []lang.Value) (lang.Value, error) {
	return nil, nil
}

func constant(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("constant: expected 1 argument")
	}
	return args[0], nil
}

// Error Handling
func tryOr(args []lang.Value) (lang.Value, error) {
	if len(args) != 2 {
		return nil, errors.New("try_or: expected 2 arguments")
	}

	value := args[0]
	fallback := args[1]

	if value == nil {
		return fallback, nil
	}

	return value, nil
}

func safe(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("safe: expected 1 argument")
	}

	value := args[0]
	if value == nil {
		return nil, nil
	}

	return value, nil
}

// Type Conversion Utilities
func toString(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("tostring: expected 1 argument")
	}
	val, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("tostring: %w", err)
	}
	return val, nil
}

func toNumber(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("tonumber: expected 1 argument")
	}
	val, err := lib.ToNumber(args[0])
	if err != nil {
		return nil, fmt.Errorf("tonumber: %w", err)
	}
	return lang.NumberValue(val), nil
}

func toBool(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("tobool: expected 1 argument")
	}
	val, err := lib.ToBool(args[0])
	if err != nil {
		return nil, fmt.Errorf("tobool: %w", err)
	}
	return lang.BoolValue(val), nil
}

func toList(args []lang.Value) (lang.Value, error) {
	if len(args) == 0 {
		return lang.ListValue{}, nil
	}

	if len(args) == 1 {
		if list, ok := args[0].(lang.ListValue); ok {
			return list, nil
		}
		return lang.ListValue{args[0]}, nil
	}

	result := make(lang.ListValue, len(args))
	copy(result, args)
	return result, nil
}

// Validation Utilities
func assert(args []lang.Value) (lang.Value, error) {
	if len(args) < 1 || len(args) > 2 {
		return nil, errors.New("assert: expected 1 or 2 arguments")
	}

	condition, err := lib.ToBool(args[0])
	if err != nil {
		return nil, fmt.Errorf("assert: condition argument: %w", err)
	}

	if !condition {
		message := "Assertion failed"
		if len(args) == 2 {
			msg, err := lib.ToString(args[1])
			if err != nil {
				return nil, fmt.Errorf("assert: message argument: %w", err)
			}
			message = string(msg)
		}
		return nil, fmt.Errorf("assertion failed: %s", message)
	}

	return lang.BoolValue(true), nil
}

func validate(args []lang.Value) (lang.Value, error) {
	if len(args) < 2 || len(args) > 3 {
		return nil, errors.New("validate: expected 2 or 3 arguments")
	}

	value := args[0]
	condition, err := lib.ToBool(args[1])
	if err != nil {
		return nil, fmt.Errorf("validate: condition argument: %w", err)
	}

	if condition {
		return value, nil
	}

	if len(args) == 3 {
		return args[2], nil // Return fallback
	}

	return nil, nil
}

func require(args []lang.Value) (lang.Value, error) {
	if len(args) < 1 || len(args) > 2 {
		return nil, errors.New("require: expected 1 or 2 arguments")
	}

	value := args[0]
	if value == nil || isNull(value) {
		message := "Required value is null"
		if len(args) == 2 {
			msg, err := lib.ToString(args[1])
			if err != nil {
				return nil, fmt.Errorf("require: message argument: %w", err)
			}
			message = string(msg)
		}
		return nil, fmt.Errorf("required value missing: %s", message)
	}

	return value, nil
}

// Functional Utilities
func apply(args []lang.Value) (lang.Value, error) {
	if len(args) < 1 {
		return nil, errors.New("apply: expected at least 1 argument")
	}
	return args[0], nil
}

func pipe(args []lang.Value) (lang.Value, error) {
	if len(args) < 1 {
		return nil, errors.New("pipe: expected at least 1 argument")
	}
	return args[0], nil
}

func compose(args []lang.Value) (lang.Value, error) {
	if len(args) < 1 {
		return nil, errors.New("compose: expected at least 1 argument")
	}
	return args[0], nil
}

// Miscellaneous Utilities
func uuid_(args []lang.Value) (lang.Value, error) {
	if len(args) != 0 {
		return nil, errors.New("uuid: expected 0 arguments")
	}
	return lang.StringValue(uuid.NewString()), nil
}

func timestamp(args []lang.Value) (lang.Value, error) {
	if len(args) != 0 {
		return nil, errors.New("timestamp: expected 0 arguments")
	}
	return lang.NumberValue(float64(time.Now().Unix())), nil
}

func randomString(args []lang.Value) (lang.Value, error) {
	length := 10 // default length
	if len(args) == 1 {
		l, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("random_string: length argument: %w", err)
		}
		length = int(l)
		if length <= 0 {
			length = 10
		}
		if length > 1000 {
			length = 1000 // cap at 1000 chars
		}
	} else if len(args) != 0 {
		return nil, errors.New("random_string: expected 0 or 1 argument")
	}

	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	rand.Seed(time.Now().UnixNano())
	for i := range result {
		result[i] = chars[rand.Intn(len(chars))]
	}

	return lang.StringValue(string(result)), nil
}

func memoize(args []lang.Value) (lang.Value, error) {
	if len(args) < 1 {
		return nil, errors.New("memoize: expected at least 1 argument")
	}
	return args[0], nil
}

func benchmark(args []lang.Value) (lang.Value, error) {
	if len(args) < 1 {
		return nil, errors.New("benchmark: expected at least 1 argument")
	}
	return lang.NumberValue(0.001), nil
}

// Helper functions for utility operations
func isNull(v lang.Value) bool {
	return v == nil
}

func isEmpty(v lang.Value) bool {
	if v == nil {
		return true
	}

	switch val := v.(type) {
	case lang.StringValue:
		return strings.TrimSpace(string(val)) == ""
	case lang.ListValue:
		return len(val) == 0
	case lang.MapValue:
		return len(val) == 0
	default:
		return isNull(v)
	}
}

func compare(a, b lang.Value) (int, error) {
	aNum, err := lib.ToNumber(a)
	if err != nil {
		return 0, fmt.Errorf("failed to convert first argument to number: %w", err)
	}
	bNum, err := lib.ToNumber(b)
	if err != nil {
		return 0, fmt.Errorf("failed to convert second argument to number: %w", err)
	}

	if aNum < bNum {
		return -1, nil
	}
	if aNum > bNum {
		return 1, nil
	}
	return 0, nil
}

func deepEqual(a, b lang.Value) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	return formatValueForDebug(a) == formatValueForDebug(b)
}

func formatValueForDebug(v lang.Value) string {
	if v == nil {
		return "null"
	}

	switch val := v.(type) {
	case lang.BoolValue:
		if bool(val) {
			return "true"
		}
		return "false"
	case lang.NumberValue:
		return strconv.FormatFloat(float64(val), 'g', -1, 64)
	case lang.StringValue:
		return fmt.Sprintf("\"%s\"", string(val))
	case lang.ListValue:
		var items []string
		for _, item := range val {
			items = append(items, formatValueForDebug(item))
		}
		return "[" + strings.Join(items, ", ") + "]"
	case lang.MapValue:
		var pairs []string
		for k, v := range val {
			pairs = append(pairs, fmt.Sprintf("\"%s\": %s", k, formatValueForDebug(v)))
		}
		return "{" + strings.Join(pairs, ", ") + "}"
	default:
		return "unknown"
	}
}

func formatValueForInspect(v lang.Value) string {
	if v == nil {
		return "null"
	}

	switch val := v.(type) {
	case lang.BoolValue:
		return fmt.Sprintf("boolean: %v", bool(val))
	case lang.NumberValue:
		return fmt.Sprintf("number: %g", float64(val))
	case lang.StringValue:
		return fmt.Sprintf("string: \"%s\" (length: %d)", string(val), len(string(val)))
	case lang.ListValue:
		return fmt.Sprintf("list: length %d, elements: %s", len(val), formatValueForDebug(v))
	case lang.MapValue:
		return fmt.Sprintf("map: %d keys, content: %s", len(val), formatValueForDebug(v))
	default:
		return "unknown type"
	}
}

// Functions that would be in the BuiltinFunctions map:
type Function func(args []lang.Value) (lang.Value, error)

var UtilityFunctions = map[string]Function{
	// Conditional functions
	"if":     conditionalIf,
	"unless": conditionalUnless,
	"switch": conditionalSwitch,

	// Null coalescing
	"coalesce":        coalesce,
	"default":         defaultValue,
	"first_non_null":  firstNonNull,
	"first_non_empty": firstNonEmpty,

	// Comparison and selection
	"greatest": greatest,
	"least":    least,
	"choose":   choose,

	// Debugging and inspection
	"debug":   debugPrint,
	"inspect": inspect,
	"dump":    dump,

	// Identity and pass-through
	"identity": identity,
	"noop":     noop,
	"constant": constant,

	// Error handling
	"try_or": tryOr,
	"safe":   safe,

	// Type conversion utilities
	"tostring": toString,
	"tonumber": toNumber,
	"tobool":   toBool,
	"tolist":   toList,

	// Validation utilities
	"assert":   assert,
	"validate": validate,
	"require":  require,

	// Functional utilities
	"apply":   apply,
	"pipe":    pipe,
	"compose": compose,

	// Miscellaneous utilities
	"uuid":          uuid_,
	"timestamp":     timestamp,
	"random_string": randomString,
	"memoize":       memoize,
	"benchmark":     benchmark,
}
