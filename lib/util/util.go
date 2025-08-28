package util

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/vedadiyan/exql/lang"
	"github.com/vedadiyan/exql/lib"
)

// Conditional Functions
func conditionalIf() (string, func([]lang.Value) (lang.Value, error)) {
	name := "if"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) < 2 || len(args) > 3 {
			return nil, lib.RangeError(name, 2, 3)
		}

		condition, err := lib.ToBool(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: condition argument: %w", name, err)
		}

		if condition {
			return args[1], nil
		}

		if len(args) == 3 {
			return args[2], nil
		}

		return nil, nil
	}
	return name, fn
}

func conditionalUnless() (string, func([]lang.Value) (lang.Value, error)) {
	name := "unless"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) < 2 || len(args) > 3 {
			return nil, lib.RangeError(name, 2, 3)
		}

		condition, err := lib.ToBool(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: condition argument: %w", name, err)
		}

		if !condition {
			return args[1], nil
		}

		if len(args) == 3 {
			return args[2], nil
		}

		return nil, nil
	}
	return name, fn
}

func conditionalSwitch() (string, func([]lang.Value) (lang.Value, error)) {
	name := "switch"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) < 3 || len(args)%2 == 0 {
			return nil, fmt.Errorf("%s: expected an odd number of arguments (value, case1, result1, ... , default)", name)
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
	return name, fn
}

// Null Coalescing Functions
func coalesce() (string, func([]lang.Value) (lang.Value, error)) {
	name := "coalesce"
	fn := func(args []lang.Value) (lang.Value, error) {
		for _, arg := range args {
			if arg != nil && !isNull(arg) {
				return arg, nil
			}
		}
		return nil, nil
	}
	return name, fn
}

func defaultValue() (string, func([]lang.Value) (lang.Value, error)) {
	name := "default"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 2 {
			return nil, lib.ArgumentError(name, 2)
		}

		value := args[0]
		defaultVal := args[1]

		if value == nil || isNull(value) {
			return defaultVal, nil
		}

		return value, nil
	}
	return name, fn
}

func firstNonNull() (string, func([]lang.Value) (lang.Value, error)) {
	name := "first_non_null"
	_, coalesceFunc := coalesce()
	fn := func(args []lang.Value) (lang.Value, error) {
		return coalesceFunc(args)
	}
	return name, fn
}

func firstNonEmpty() (string, func([]lang.Value) (lang.Value, error)) {
	name := "first_non_empty"
	fn := func(args []lang.Value) (lang.Value, error) {
		for _, arg := range args {
			if arg != nil && !isEmpty(arg) {
				return arg, nil
			}
		}
		return nil, nil
	}
	return name, fn
}

// Comparison and Selection
func greatest() (string, func([]lang.Value) (lang.Value, error)) {
	name := "greatest"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) == 0 {
			return nil, nil
		}

		max := args[0]
		for i := 1; i < len(args); i++ {
			comparison, err := compare(args[i], max)
			if err != nil {
				return nil, fmt.Errorf("%s: comparison failed: %w", name, err)
			}
			if comparison > 0 {
				max = args[i]
			}
		}
		return max, nil
	}
	return name, fn
}

func least() (string, func([]lang.Value) (lang.Value, error)) {
	name := "least"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) == 0 {
			return nil, nil
		}

		min := args[0]
		for i := 1; i < len(args); i++ {
			comparison, err := compare(args[i], min)
			if err != nil {
				return nil, fmt.Errorf("%s: comparison failed: %w", name, err)
			}
			if comparison < 0 {
				min = args[i]
			}
		}
		return min, nil
	}
	return name, fn
}

func choose() (string, func([]lang.Value) (lang.Value, error)) {
	name := "choose"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("%s: expected at least 2 arguments", name)
		}

		index, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: index argument: %w", name, err)
		}

		intIndex := int(index)
		if intIndex < 1 || intIndex > len(args)-1 {
			return nil, fmt.Errorf("%s: index out of bounds", name)
		}

		return args[intIndex], nil // 1-based indexing
	}
	return name, fn
}

// Debugging and Inspection
func debugPrint() (string, func([]lang.Value) (lang.Value, error)) {
	name := "debug"
	fn := func(args []lang.Value) (lang.Value, error) {
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
	return name, fn
}

func inspect() (string, func([]lang.Value) (lang.Value, error)) {
	name := "inspect"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}

		return lang.StringValue(formatValueForInspect(args[0])), nil
	}
	return name, fn
}

func dump() (string, func([]lang.Value) (lang.Value, error)) {
	name := "dump"
	fn := func(args []lang.Value) (lang.Value, error) {
		var parts []string
		for _, arg := range args {
			parts = append(parts, formatValueForInspect(arg))
		}

		return lang.StringValue(strings.Join(parts, "\n")), nil
	}
	return name, fn
}

// Identity and Pass-through
func identity() (string, func([]lang.Value) (lang.Value, error)) {
	name := "identity"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		return args[0], nil
	}
	return name, fn
}

func noop() (string, func([]lang.Value) (lang.Value, error)) {
	name := "noop"
	fn := func(args []lang.Value) (lang.Value, error) {
		return nil, nil
	}
	return name, fn
}

func constant() (string, func([]lang.Value) (lang.Value, error)) {
	name := "constant"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		return args[0], nil
	}
	return name, fn
}

// Error Handling
func tryOr() (string, func([]lang.Value) (lang.Value, error)) {
	name := "try_or"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 2 {
			return nil, lib.ArgumentError(name, 2)
		}

		value := args[0]
		fallback := args[1]

		if value == nil {
			return fallback, nil
		}

		return value, nil
	}
	return name, fn
}

func safe() (string, func([]lang.Value) (lang.Value, error)) {
	name := "safe"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}

		value := args[0]
		if value == nil {
			return nil, nil
		}

		return value, nil
	}
	return name, fn
}

// Type Conversion Utilities
func toString() (string, func([]lang.Value) (lang.Value, error)) {
	name := "tostring"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		val, err := lib.ToString(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		return val, nil
	}
	return name, fn
}

func toNumber() (string, func([]lang.Value) (lang.Value, error)) {
	name := "tonumber"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		val, err := lib.ToNumber(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		return lang.NumberValue(val), nil
	}
	return name, fn
}

func toBool() (string, func([]lang.Value) (lang.Value, error)) {
	name := "tobool"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		val, err := lib.ToBool(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		return lang.BoolValue(val), nil
	}
	return name, fn
}

func toList() (string, func([]lang.Value) (lang.Value, error)) {
	name := "tolist"
	fn := func(args []lang.Value) (lang.Value, error) {
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
	return name, fn
}

// Validation Utilities
func assert() (string, func([]lang.Value) (lang.Value, error)) {
	name := "assert"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) < 1 || len(args) > 2 {
			return nil, lib.RangeError(name, 1, 2)
		}

		condition, err := lib.ToBool(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: condition argument: %w", name, err)
		}

		if !condition {
			message := "Assertion failed"
			if len(args) == 2 {
				msg, err := lib.ToString(args[1])
				if err != nil {
					return nil, fmt.Errorf("%s: message argument: %w", name, err)
				}
				message = string(msg)
			}
			return nil, fmt.Errorf("assertion failed: %s", message)
		}

		return lang.BoolValue(true), nil
	}
	return name, fn
}

func validate() (string, func([]lang.Value) (lang.Value, error)) {
	name := "validate"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) < 2 || len(args) > 3 {
			return nil, lib.RangeError(name, 2, 3)
		}

		value := args[0]
		condition, err := lib.ToBool(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: condition argument: %w", name, err)
		}

		if condition {
			return value, nil
		}

		if len(args) == 3 {
			return args[2], nil // Return fallback
		}

		return nil, nil
	}
	return name, fn
}

func require() (string, func([]lang.Value) (lang.Value, error)) {
	name := "require"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) < 1 || len(args) > 2 {
			return nil, lib.RangeError(name, 1, 2)
		}

		value := args[0]
		if value == nil || isNull(value) {
			message := "Required value is null"
			if len(args) == 2 {
				msg, err := lib.ToString(args[1])
				if err != nil {
					return nil, fmt.Errorf("%s: message argument: %w", name, err)
				}
				message = string(msg)
			}
			return nil, fmt.Errorf("required value missing: %s", message)
		}

		return value, nil
	}
	return name, fn
}

// Functional Utilities
func apply() (string, func([]lang.Value) (lang.Value, error)) {
	name := "apply"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("%s: expected at least 1 argument", name)
		}
		return args[0], nil
	}
	return name, fn
}

func pipe() (string, func([]lang.Value) (lang.Value, error)) {
	name := "pipe"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("%s: expected at least 1 argument", name)
		}
		return args[0], nil
	}
	return name, fn
}

func compose() (string, func([]lang.Value) (lang.Value, error)) {
	name := "compose"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("%s: expected at least 1 argument", name)
		}
		return args[0], nil
	}
	return name, fn
}

// Miscellaneous Utilities
func uuid_() (string, func([]lang.Value) (lang.Value, error)) {
	name := "uuid"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 0 {
			return nil, lib.ArgumentError(name, 0)
		}
		return lang.StringValue(uuid.NewString()), nil
	}
	return name, fn
}

func timestamp() (string, func([]lang.Value) (lang.Value, error)) {
	name := "timestamp"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 0 {
			return nil, lib.ArgumentError(name, 0)
		}
		return lang.NumberValue(float64(time.Now().Unix())), nil
	}
	return name, fn
}

func randomString() (string, func([]lang.Value) (lang.Value, error)) {
	name := "random_string"
	fn := func(args []lang.Value) (lang.Value, error) {
		length := 10 // default length
		if len(args) == 1 {
			l, err := lib.ToNumber(args[0])
			if err != nil {
				return nil, fmt.Errorf("%s: length argument: %w", name, err)
			}
			length = int(l)
			if length <= 0 {
				length = 10
			}
			if length > 1000 {
				length = 1000 // cap at 1000 chars
			}
		} else if len(args) != 0 {
			return nil, lib.RangeError(name, 0, 1)
		}

		chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		result := make([]byte, length)
		rand.Seed(time.Now().UnixNano())
		for i := range result {
			result[i] = chars[rand.Intn(len(chars))]
		}

		return lang.StringValue(string(result)), nil
	}
	return name, fn
}

func memoize() (string, func([]lang.Value) (lang.Value, error)) {
	name := "memoize"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("%s: expected at least 1 argument", name)
		}
		return args[0], nil
	}
	return name, fn
}

func benchmark() (string, func([]lang.Value) (lang.Value, error)) {
	name := "benchmark"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) < 1 {
			return nil, fmt.Errorf("%s: expected at least 1 argument", name)
		}
		return lang.NumberValue(0.001), nil
	}
	return name, fn
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

var UtilityFunctions = []func() (string, func([]lang.Value) (lang.Value, error)){
	conditionalIf,
	conditionalUnless,
	conditionalSwitch,
	coalesce,
	defaultValue,
	firstNonNull,
	firstNonEmpty,
	greatest,
	least,
	choose,
	debugPrint,
	inspect,
	dump,
	identity,
	noop,
	constant,
	tryOr,
	safe,
	toString,
	toNumber,
	toBool,
	toList,
	assert,
	validate,
	require,
	apply,
	pipe,
	compose,
	uuid_,
	timestamp,
	randomString,
	memoize,
	benchmark,
}
