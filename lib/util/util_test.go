package util

import (
	"strings"
	"testing"
	"time"

	"github.com/vedadiyan/exql/lang"
)

func TestArgumentErrors(t *testing.T) {
	functions := []struct {
		fn                  func() (string, lang.Function)
		expectedArgs        int
		ignoreArgumentError bool
	}{
		{identity, 1, false},
		{noop, 0, true},
		{constant, 1, false},
		{defaultValue, 2, false},
		{debugPrint, -1, false}, // Variable arguments
		{inspect, 1, false},
		{dump, -1, false}, // Variable arguments
		{toString, 1, false},
		{toNumber, 1, false},
		{toBool, 1, false},
		{toList, -1, false}, // Variable arguments
		{tryOr, 2, false},
		{safe, 1, false},
		{uuid_, 0, false},
		{timestamp, 0, false},
		{randomString, -1, false}, // 0 or 1 arguments
	}

	for _, tf := range functions {
		name, fn := tf.fn()
		if tf.expectedArgs == -1 {
			// Skip variable argument functions
			continue
		}
		if tf.ignoreArgumentError {
			continue
		}
		t.Run(name+"_wrong_args", func(t *testing.T) {
			var args []lang.Value
			if tf.expectedArgs == 0 {
				args = []lang.Value{lang.NumberValue(1)} // Add extra arg for 0-arg functions
			} else {
				args = make([]lang.Value, tf.expectedArgs-1) // One less than expected
				for i := range args {
					args[i] = lang.StringValue("test")
				}
			}
			_, err := fn(args)
			if err == nil {
				t.Errorf("Expected error for wrong argument count in %s", name)
			}
		})
	}

	// Test range argument functions separately
	rangeFunctions := []struct {
		fn  func() (string, lang.Function)
		min int
		max int
	}{
		{conditionalIf, 2, 3},
		{conditionalUnless, 2, 3},
		{assert, 1, 2},
		{validate, 2, 3},
		{require, 1, 2},
	}

	for _, tf := range rangeFunctions {
		name, fn := tf.fn()

		// Test too few arguments
		t.Run(name+"_too_few_args", func(t *testing.T) {
			args := make([]lang.Value, tf.min-1)
			for i := range args {
				args[i] = lang.StringValue("test")
			}
			_, err := fn(args)
			if err == nil {
				t.Errorf("Expected error for too few arguments in %s", name)
			}
		})

		// Test too many arguments
		t.Run(name+"_too_many_args", func(t *testing.T) {
			args := make([]lang.Value, tf.max+1)
			for i := range args {
				args[i] = lang.StringValue("test")
			}
			_, err := fn(args)
			if err == nil {
				t.Errorf("Expected error for too many arguments in %s", name)
			}
		})
	}
}

func TestExport(t *testing.T) {
	functions := Export()

	expectedFunctions := []string{
		"if", "unless", "switch", "coalesce", "default", "firstNonNull", "firstNonEmpty",
		"greatest", "least", "choose", "debug", "inspect", "dump", "identity", "noop",
		"constant", "tryOr", "safe", "tostring", "tonumber", "tobool", "tolist",
		"assert", "validate", "require", "apply", "pipe", "compose", "uuid", "timestamp",
		"randomString", "memoize", "benchmark",
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

func TestConditionalIf(t *testing.T) {
	_, fn := conditionalIf()

	tests := []struct {
		name      string
		args      []lang.Value
		expected  lang.Value
		expectErr bool
	}{
		{"true condition with value", []lang.Value{lang.BoolValue(true), lang.StringValue("yes")}, lang.StringValue("yes"), false},
		{"false condition with value", []lang.Value{lang.BoolValue(false), lang.StringValue("yes")}, nil, false},
		{"true condition with else", []lang.Value{lang.BoolValue(true), lang.StringValue("yes"), lang.StringValue("no")}, lang.StringValue("yes"), false},
		{"false condition with else", []lang.Value{lang.BoolValue(false), lang.StringValue("yes"), lang.StringValue("no")}, lang.StringValue("no"), false},
		{"number condition true", []lang.Value{lang.NumberValue(1), lang.StringValue("yes")}, lang.StringValue("yes"), false},
		{"number condition false", []lang.Value{lang.NumberValue(0), lang.StringValue("yes")}, nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn(tt.args)
			if tt.expectErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !deepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestConditionalUnless(t *testing.T) {
	_, fn := conditionalUnless()

	tests := []struct {
		name     string
		args     []lang.Value
		expected lang.Value
	}{
		{"false condition with value", []lang.Value{lang.BoolValue(false), lang.StringValue("yes")}, lang.StringValue("yes")},
		{"true condition with value", []lang.Value{lang.BoolValue(true), lang.StringValue("yes")}, nil},
		{"false condition with else", []lang.Value{lang.BoolValue(false), lang.StringValue("yes"), lang.StringValue("no")}, lang.StringValue("yes")},
		{"true condition with else", []lang.Value{lang.BoolValue(true), lang.StringValue("yes"), lang.StringValue("no")}, lang.StringValue("no")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn(tt.args)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !deepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestConditionalSwitch(t *testing.T) {
	_, fn := conditionalSwitch()

	tests := []struct {
		name      string
		args      []lang.Value
		expected  lang.Value
		expectErr bool
	}{
		{
			"match first case",
			[]lang.Value{lang.NumberValue(1), lang.NumberValue(1), lang.StringValue("one"), lang.StringValue("default")},
			lang.StringValue("one"),
			false,
		},
		{
			"no match return default",
			[]lang.Value{lang.NumberValue(3), lang.NumberValue(1), lang.StringValue("one"), lang.StringValue("default")},
			lang.StringValue("default"),
			false,
		},
		{
			"string match",
			[]lang.Value{lang.StringValue("a"), lang.StringValue("a"), lang.StringValue("first"), lang.StringValue("b"), lang.StringValue("second"), lang.StringValue("default")},
			lang.StringValue("first"),
			false,
		},
		{
			"even number of args error",
			[]lang.Value{lang.NumberValue(1), lang.StringValue("one")},
			nil,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn(tt.args)
			if tt.expectErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !deepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestCoalesce(t *testing.T) {
	_, fn := coalesce()

	tests := []struct {
		name     string
		args     []lang.Value
		expected lang.Value
	}{
		{"first non-null", []lang.Value{nil, lang.StringValue("hello"), lang.StringValue("world")}, lang.StringValue("hello")},
		{"all null", []lang.Value{nil, nil, nil}, nil},
		{"no nulls", []lang.Value{lang.StringValue("first"), lang.StringValue("second")}, lang.StringValue("first")},
		{"empty args", []lang.Value{}, nil},
		{"single non-null", []lang.Value{lang.NumberValue(42)}, lang.NumberValue(42)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn(tt.args)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !deepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestDefaultValue(t *testing.T) {
	_, fn := defaultValue()

	tests := []struct {
		name     string
		args     []lang.Value
		expected lang.Value
	}{
		{"null value uses default", []lang.Value{nil, lang.StringValue("default")}, lang.StringValue("default")},
		{"non-null value returns value", []lang.Value{lang.StringValue("value"), lang.StringValue("default")}, lang.StringValue("value")},
		{"zero value returns value", []lang.Value{lang.NumberValue(0), lang.StringValue("default")}, lang.NumberValue(0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn(tt.args)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !deepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestFirstNonEmpty(t *testing.T) {
	_, fn := firstNonEmpty()

	tests := []struct {
		name     string
		args     []lang.Value
		expected lang.Value
	}{
		{"first non-empty string", []lang.Value{lang.StringValue(""), lang.StringValue("hello")}, lang.StringValue("hello")},
		{"first non-empty list", []lang.Value{lang.ListValue{}, lang.ListValue{lang.NumberValue(1)}}, lang.ListValue{lang.NumberValue(1)}},
		{"all empty", []lang.Value{lang.StringValue(""), lang.ListValue{}, lang.MapValue{}}, nil},
		{"whitespace is empty", []lang.Value{lang.StringValue("   "), lang.StringValue("hello")}, lang.StringValue("hello")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn(tt.args)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !deepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestGreatest(t *testing.T) {
	_, fn := greatest()

	tests := []struct {
		name     string
		args     []lang.Value
		expected lang.Value
	}{
		{"numbers", []lang.Value{lang.NumberValue(1), lang.NumberValue(5), lang.NumberValue(3)}, lang.NumberValue(5)},
		{"single value", []lang.Value{lang.NumberValue(42)}, lang.NumberValue(42)},
		{"empty args", []lang.Value{}, nil},
		{"negative numbers", []lang.Value{lang.NumberValue(-1), lang.NumberValue(-5), lang.NumberValue(-3)}, lang.NumberValue(-1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn(tt.args)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !deepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestLeast(t *testing.T) {
	_, fn := least()

	tests := []struct {
		name     string
		args     []lang.Value
		expected lang.Value
	}{
		{"numbers", []lang.Value{lang.NumberValue(1), lang.NumberValue(5), lang.NumberValue(3)}, lang.NumberValue(1)},
		{"single value", []lang.Value{lang.NumberValue(42)}, lang.NumberValue(42)},
		{"empty args", []lang.Value{}, nil},
		{"negative numbers", []lang.Value{lang.NumberValue(-1), lang.NumberValue(-5), lang.NumberValue(-3)}, lang.NumberValue(-5)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn(tt.args)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !deepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestChoose(t *testing.T) {
	_, fn := choose()

	tests := []struct {
		name      string
		args      []lang.Value
		expected  lang.Value
		expectErr bool
	}{
		{"choose first", []lang.Value{lang.NumberValue(1), lang.StringValue("a"), lang.StringValue("b")}, lang.StringValue("a"), false},
		{"choose second", []lang.Value{lang.NumberValue(2), lang.StringValue("a"), lang.StringValue("b")}, lang.StringValue("b"), false},
		{"index out of bounds", []lang.Value{lang.NumberValue(3), lang.StringValue("a"), lang.StringValue("b")}, nil, true},
		{"index zero", []lang.Value{lang.NumberValue(0), lang.StringValue("a"), lang.StringValue("b")}, nil, true},
		{"negative index", []lang.Value{lang.NumberValue(-1), lang.StringValue("a"), lang.StringValue("b")}, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn(tt.args)
			if tt.expectErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !deepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestInspect(t *testing.T) {
	_, fn := inspect()

	tests := []struct {
		name     string
		input    lang.Value
		contains string
	}{
		{"string", lang.StringValue("hello"), "string:"},
		{"number", lang.NumberValue(42), "number:"},
		{"boolean", lang.BoolValue(true), "boolean:"},
		{"list", lang.ListValue{lang.NumberValue(1)}, "list:"},
		{"map", lang.MapValue{"key": lang.StringValue("value")}, "map:"},
		{"null", nil, "null"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{tt.input})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			resultStr := string(result.(lang.StringValue))
			if !strings.Contains(resultStr, tt.contains) {
				t.Errorf("expected result to contain %s, got %s", tt.contains, resultStr)
			}
		})
	}
}

func TestIdentity(t *testing.T) {
	_, fn := identity()

	tests := []lang.Value{
		lang.StringValue("hello"),
		lang.NumberValue(42),
		lang.BoolValue(true),
		nil,
		lang.ListValue{lang.NumberValue(1)},
		lang.MapValue{"key": lang.StringValue("value")},
	}

	for _, tt := range tests {
		t.Run("identity", func(t *testing.T) {
			result, err := fn([]lang.Value{tt})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !deepEqual(result, tt) {
				t.Errorf("expected %v, got %v", tt, result)
			}
		})
	}
}

func TestNoop(t *testing.T) {
	_, fn := noop()

	result, err := fn([]lang.Value{lang.StringValue("anything")})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestToString(t *testing.T) {
	_, fn := toString()

	tests := []struct {
		name     string
		input    lang.Value
		expected string
	}{
		{"string", lang.StringValue("hello"), "hello"},
		{"number", lang.NumberValue(42), "42"},
		{"boolean true", lang.BoolValue(true), "true"},
		{"boolean false", lang.BoolValue(false), "false"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{tt.input})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if string(result.(lang.StringValue)) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, string(result.(lang.StringValue)))
			}
		})
	}
}

func TestToNumber(t *testing.T) {
	_, fn := toNumber()

	tests := []struct {
		name     string
		input    lang.Value
		expected float64
	}{
		{"number", lang.NumberValue(42), 42},
		{"string number", lang.StringValue("42"), 42},
		{"boolean true", lang.BoolValue(true), 1},
		{"boolean false", lang.BoolValue(false), 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{tt.input})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if float64(result.(lang.NumberValue)) != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, float64(result.(lang.NumberValue)))
			}
		})
	}
}

func TestToBool(t *testing.T) {
	_, fn := toBool()

	tests := []struct {
		name     string
		input    lang.Value
		expected bool
	}{
		{"boolean true", lang.BoolValue(true), true},
		{"boolean false", lang.BoolValue(false), false},
		{"number zero", lang.NumberValue(0), false},
		{"number non-zero", lang.NumberValue(42), true},
		{"empty string", lang.StringValue(""), false},
		{"non-empty string", lang.StringValue("hello"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{tt.input})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if bool(result.(lang.BoolValue)) != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, bool(result.(lang.BoolValue)))
			}
		})
	}
}

func TestToList(t *testing.T) {
	_, fn := toList()

	tests := []struct {
		name     string
		args     []lang.Value
		expected lang.ListValue
	}{
		{"empty args", []lang.Value{}, lang.ListValue{}},
		{"single list arg", []lang.Value{lang.ListValue{lang.NumberValue(1)}}, lang.ListValue{lang.NumberValue(1)}},
		{"single non-list arg", []lang.Value{lang.StringValue("hello")}, lang.ListValue{lang.StringValue("hello")}},
		{"multiple args", []lang.Value{lang.StringValue("a"), lang.NumberValue(1)}, lang.ListValue{lang.StringValue("a"), lang.NumberValue(1)}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn(tt.args)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			resultList := result.(lang.ListValue)
			if len(resultList) != len(tt.expected) {
				t.Errorf("expected length %d, got %d", len(tt.expected), len(resultList))
			}
			for i, expected := range tt.expected {
				if i < len(resultList) && !deepEqual(resultList[i], expected) {
					t.Errorf("index %d: expected %v, got %v", i, expected, resultList[i])
				}
			}
		})
	}
}

func TestAssert(t *testing.T) {
	_, fn := assert()

	tests := []struct {
		name      string
		args      []lang.Value
		expected  lang.Value
		expectErr bool
	}{
		{"true assertion", []lang.Value{lang.BoolValue(true)}, lang.BoolValue(true), false},
		{"false assertion", []lang.Value{lang.BoolValue(false)}, nil, true},
		{"true with message", []lang.Value{lang.BoolValue(true), lang.StringValue("msg")}, lang.BoolValue(true), false},
		{"false with message", []lang.Value{lang.BoolValue(false), lang.StringValue("custom")}, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn(tt.args)
			if tt.expectErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !deepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestRequire(t *testing.T) {
	_, fn := require()

	tests := []struct {
		name      string
		args      []lang.Value
		expected  lang.Value
		expectErr bool
	}{
		{"non-null value", []lang.Value{lang.StringValue("value")}, lang.StringValue("value"), false},
		{"null value", []lang.Value{nil}, nil, true},
		{"null with message", []lang.Value{nil, lang.StringValue("custom message")}, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn(tt.args)
			if tt.expectErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !deepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestUUID(t *testing.T) {
	_, fn := uuid_()

	result, err := fn([]lang.Value{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	uuidStr := string(result.(lang.StringValue))
	if len(uuidStr) != 36 {
		t.Errorf("expected UUID length 36, got %d", len(uuidStr))
	}

	// Test that consecutive calls return different UUIDs
	result2, err := fn([]lang.Value{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if uuidStr == string(result2.(lang.StringValue)) {
		t.Error("expected different UUIDs on consecutive calls")
	}
}

func TestTimestamp(t *testing.T) {
	_, fn := timestamp()

	before := time.Now().Unix()
	result, err := fn([]lang.Value{})
	after := time.Now().Unix()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	timestamp := int64(result.(lang.NumberValue))
	if timestamp < before || timestamp > after {
		t.Errorf("timestamp %d not in range [%d, %d]", timestamp, before, after)
	}
}

func TestRandomString(t *testing.T) {
	_, fn := randomString()

	// Test default length
	result, err := fn([]lang.Value{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	randomStr := string(result.(lang.StringValue))
	if len(randomStr) != 10 {
		t.Errorf("expected default length 10, got %d", len(randomStr))
	}

	// Test custom length
	result2, err := fn([]lang.Value{lang.NumberValue(5)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	randomStr2 := string(result2.(lang.StringValue))
	if len(randomStr2) != 5 {
		t.Errorf("expected custom length 5, got %d", len(randomStr2))
	}

	// Test that results are different
	if randomStr == randomStr2 {
		t.Error("expected different random strings")
	}
}

// Test Helper Functions
func TestIsNull(t *testing.T) {
	tests := []struct {
		name     string
		input    lang.Value
		expected bool
	}{
		{"null", nil, true},
		{"string", lang.StringValue("hello"), false},
		{"number", lang.NumberValue(42), false},
		{"boolean", lang.BoolValue(false), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isNull(tt.input)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestIsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		input    lang.Value
		expected bool
	}{
		{"null", nil, true},
		{"empty string", lang.StringValue(""), true},
		{"whitespace string", lang.StringValue("   "), true},
		{"non-empty string", lang.StringValue("hello"), false},
		{"empty list", lang.ListValue{}, true},
		{"non-empty list", lang.ListValue{lang.NumberValue(1)}, false},
		{"empty map", lang.MapValue{}, true},
		{"non-empty map", lang.MapValue{"key": lang.StringValue("value")}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isEmpty(tt.input)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestCompare(t *testing.T) {
	tests := []struct {
		name     string
		a, b     lang.Value
		expected int
	}{
		{"equal numbers", lang.NumberValue(5), lang.NumberValue(5), 0},
		{"a less than b", lang.NumberValue(3), lang.NumberValue(5), -1},
		{"a greater than b", lang.NumberValue(7), lang.NumberValue(5), 1},
		{"convert strings", lang.StringValue("5"), lang.StringValue("3"), 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := compare(tt.a, tt.b)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestFormatValueForDebug(t *testing.T) {
	tests := []struct {
		name     string
		input    lang.Value
		expected string
	}{
		{"null", nil, "null"},
		{"boolean true", lang.BoolValue(true), "true"},
		{"boolean false", lang.BoolValue(false), "false"},
		{"number", lang.NumberValue(42), "42"},
		{"string", lang.StringValue("hello"), "\"hello\""},
		{"list", lang.ListValue{lang.NumberValue(1), lang.StringValue("a")}, "[1, \"a\"]"},
		{"map", lang.MapValue{"key": lang.StringValue("value")}, "{\"key\": \"value\"}"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatValueForDebug(tt.input)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestFormatValueForInspect(t *testing.T) {
	tests := []struct {
		name     string
		input    lang.Value
		contains []string
	}{
		{"null", nil, []string{"null"}},
		{"boolean", lang.BoolValue(true), []string{"boolean:", "true"}},
		{"number", lang.NumberValue(42), []string{"number:", "42"}},
		{"string", lang.StringValue("hello"), []string{"string:", "\"hello\"", "length: 5"}},
		{"list", lang.ListValue{lang.NumberValue(1)}, []string{"list:", "length 1"}},
		{"map", lang.MapValue{"key": lang.StringValue("value")}, []string{"map:", "1 keys"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatValueForInspect(tt.input)
			for _, expected := range tt.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("expected result to contain %s, got %s", expected, result)
				}
			}
		})
	}
}

func TestValidate(t *testing.T) {
	_, fn := validate()

	tests := []struct {
		name      string
		args      []lang.Value
		expected  lang.Value
		expectErr bool
	}{
		{"valid condition", []lang.Value{lang.StringValue("test"), lang.BoolValue(true)}, lang.StringValue("test"), false},
		{"invalid condition", []lang.Value{lang.StringValue("test"), lang.BoolValue(false)}, nil, false},
		{"invalid with fallback", []lang.Value{lang.StringValue("test"), lang.BoolValue(false), lang.StringValue("fallback")}, lang.StringValue("fallback"), false},
		{"valid with fallback", []lang.Value{lang.StringValue("test"), lang.BoolValue(true), lang.StringValue("fallback")}, lang.StringValue("test"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn(tt.args)
			if tt.expectErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !deepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestTryOr(t *testing.T) {
	_, fn := tryOr()

	tests := []struct {
		name     string
		args     []lang.Value
		expected lang.Value
	}{
		{"non-null value", []lang.Value{lang.StringValue("value"), lang.StringValue("fallback")}, lang.StringValue("value")},
		{"null value", []lang.Value{nil, lang.StringValue("fallback")}, lang.StringValue("fallback")},
		{"zero value", []lang.Value{lang.NumberValue(0), lang.StringValue("fallback")}, lang.NumberValue(0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn(tt.args)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !deepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestSafe(t *testing.T) {
	_, fn := safe()

	tests := []struct {
		name     string
		input    lang.Value
		expected lang.Value
	}{
		{"non-null value", lang.StringValue("value"), lang.StringValue("value")},
		{"null value", nil, nil},
		{"zero value", lang.NumberValue(0), lang.NumberValue(0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{tt.input})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !deepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestFunctionalUtilities(t *testing.T) {
	// Test apply, pipe, compose - they all just return the first argument
	functionalFuncs := []func() (string, lang.Function){apply, pipe, compose}

	for _, fnCreator := range functionalFuncs {
		name, fn := fnCreator()
		t.Run(name, func(t *testing.T) {
			testValue := lang.StringValue("test")
			result, err := fn([]lang.Value{testValue, lang.StringValue("extra")})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !deepEqual(result, testValue) {
				t.Errorf("expected %v, got %v", testValue, result)
			}
		})
	}
}

func TestMemoizeAndBenchmark(t *testing.T) {
	// Test memoize and benchmark - they return simplified results
	t.Run("memoize", func(t *testing.T) {
		_, fn := memoize()
		testValue := lang.StringValue("test")
		result, err := fn([]lang.Value{testValue})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !deepEqual(result, testValue) {
			t.Errorf("expected %v, got %v", testValue, result)
		}
	})

	t.Run("benchmark", func(t *testing.T) {
		_, fn := benchmark()
		result, err := fn([]lang.Value{lang.StringValue("test")})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if float64(result.(lang.NumberValue)) != 0.001 {
			t.Errorf("expected 0.001, got %v", float64(result.(lang.NumberValue)))
		}
	})
}

// Test Error Cases
func TestInvalidInputs(t *testing.T) {
	testCases := []struct {
		name string
		fn   func() (string, lang.Function)
		args []lang.Value
	}{
		{"toString with unconvertible", toString, []lang.Value{lang.ListValue{}}},
		{"toNumber with non-numeric string", toNumber, []lang.Value{lang.StringValue("abc")}},
		{"choose with invalid index", choose, []lang.Value{lang.StringValue("not a number"), lang.StringValue("a")}},
		{"greatest with uncomparable", greatest, []lang.Value{lang.StringValue("abc"), lang.NumberValue(1)}},
		{"least with uncomparable", least, []lang.Value{lang.StringValue("abc"), lang.NumberValue(1)}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, fn := tc.fn()
			_, err := fn(tc.args)
			if err == nil {
				t.Error("expected error for invalid input")
			}
		})
	}
}

func TestRandomStringBounds(t *testing.T) {
	_, fn := randomString()

	// Test zero length
	result, err := fn([]lang.Value{lang.NumberValue(0)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	randomStr := string(result.(lang.StringValue))
	if len(randomStr) != 10 { // Should default to 10
		t.Errorf("expected default length 10 for zero input, got %d", len(randomStr))
	}

	// Test negative length
	result2, err := fn([]lang.Value{lang.NumberValue(-5)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	randomStr2 := string(result2.(lang.StringValue))
	if len(randomStr2) != 10 { // Should default to 10
		t.Errorf("expected default length 10 for negative input, got %d", len(randomStr2))
	}

	// Test max length cap
	result3, err := fn([]lang.Value{lang.NumberValue(2000)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	randomStr3 := string(result3.(lang.StringValue))
	if len(randomStr3) != 1000 { // Should cap at 1000
		t.Errorf("expected capped length 1000 for large input, got %d", len(randomStr3))
	}
}

// Benchmark Tests
func BenchmarkConditionalIf(b *testing.B) {
	_, fn := conditionalIf()
	input := []lang.Value{lang.BoolValue(true), lang.StringValue("yes"), lang.StringValue("no")}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fn(input)
	}
}

func BenchmarkCoalesce(b *testing.B) {
	_, fn := coalesce()
	input := []lang.Value{nil, nil, lang.StringValue("value"), lang.StringValue("other")}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fn(input)
	}
}

func BenchmarkUUID(b *testing.B) {
	_, fn := uuid_()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fn([]lang.Value{})
	}
}

func BenchmarkRandomString(b *testing.B) {
	_, fn := randomString()
	input := []lang.Value{lang.NumberValue(50)}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fn(input)
	}
}
