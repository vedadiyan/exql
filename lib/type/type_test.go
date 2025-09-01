package types

import (
	"testing"

	"github.com/vedadiyan/exql/lang"
)

// Test Basic Type Checking
func TestTypeOf(t *testing.T) {
	_, fn := typeOf()

	tests := []struct {
		name     string
		input    lang.Value
		expected string
	}{
		{"null", nil, "null"},
		{"boolean", lang.BoolValue(true), "boolean"},
		{"number", lang.NumberValue(42), "number"},
		{"string", lang.StringValue("hello"), "string"},
		{"list", lang.ListValue{lang.NumberValue(1), lang.NumberValue(2)}, "list"},
		{"map", lang.MapValue{"key": lang.StringValue("value")}, "map"},
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

	// Test wrong number of arguments
	_, err := fn([]lang.Value{})
	if err == nil {
		t.Error("expected error for no arguments")
	}
}

func TestIsNull(t *testing.T) {
	_, fn := isNull()

	tests := []struct {
		name     string
		input    lang.Value
		expected bool
	}{
		{"null value", nil, true},
		{"non-null value", lang.StringValue("test"), false},
		{"zero number", lang.NumberValue(0), false},
		{"false boolean", lang.BoolValue(false), false},
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

func TestIsDefined(t *testing.T) {
	_, fn := isDefined()

	tests := []struct {
		name     string
		input    lang.Value
		expected bool
	}{
		{"null value", nil, false},
		{"defined value", lang.StringValue("test"), true},
		{"zero number", lang.NumberValue(0), true},
		{"empty string", lang.StringValue(""), true},
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

func TestIsEmpty(t *testing.T) {
	_, fn := isEmpty()

	tests := []struct {
		name     string
		input    lang.Value
		expected bool
	}{
		{"null", nil, true},
		{"empty string", lang.StringValue(""), true},
		{"non-empty string", lang.StringValue("hello"), false},
		{"empty list", lang.ListValue{}, true},
		{"non-empty list", lang.ListValue{lang.NumberValue(1)}, false},
		{"empty map", lang.MapValue{}, true},
		{"non-empty map", lang.MapValue{"key": lang.StringValue("value")}, false},
		{"zero number", lang.NumberValue(0), true},
		{"non-zero number", lang.NumberValue(42), false},
		{"false boolean", lang.BoolValue(false), true},
		{"true boolean", lang.BoolValue(true), false},
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

func TestIsNotEmpty(t *testing.T) {
	_, fn := isNotEmpty()

	tests := []struct {
		name     string
		input    lang.Value
		expected bool
	}{
		{"null", nil, false},
		{"empty string", lang.StringValue(""), false},
		{"non-empty string", lang.StringValue("hello"), true},
		{"zero number", lang.NumberValue(0), false},
		{"non-zero number", lang.NumberValue(42), true},
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

// Test Primitive Type Checking
func TestIsBool(t *testing.T) {
	_, fn := isBool()

	tests := []struct {
		name     string
		input    lang.Value
		expected bool
	}{
		{"true boolean", lang.BoolValue(true), true},
		{"false boolean", lang.BoolValue(false), true},
		{"number", lang.NumberValue(1), false},
		{"string", lang.StringValue("true"), false},
		{"null", nil, false},
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

func TestIsNumber(t *testing.T) {
	_, fn := isNumber()

	tests := []struct {
		name     string
		input    lang.Value
		expected bool
	}{
		{"integer", lang.NumberValue(42), true},
		{"float", lang.NumberValue(3.14), true},
		{"zero", lang.NumberValue(0), true},
		{"negative", lang.NumberValue(-5), true},
		{"string", lang.StringValue("42"), false},
		{"boolean", lang.BoolValue(true), false},
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

func TestIsString(t *testing.T) {
	_, fn := isString()

	tests := []struct {
		name     string
		input    lang.Value
		expected bool
	}{
		{"string", lang.StringValue("hello"), true},
		{"empty string", lang.StringValue(""), true},
		{"number", lang.NumberValue(42), false},
		{"boolean", lang.BoolValue(true), false},
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

func TestIsList(t *testing.T) {
	_, fn := isList()

	tests := []struct {
		name     string
		input    lang.Value
		expected bool
	}{
		{"list", lang.ListValue{lang.NumberValue(1), lang.NumberValue(2)}, true},
		{"empty list", lang.ListValue{}, true},
		{"string", lang.StringValue("hello"), false},
		{"map", lang.MapValue{"key": lang.StringValue("value")}, false},
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

func TestIsMap(t *testing.T) {
	_, fn := isMap()

	tests := []struct {
		name     string
		input    lang.Value
		expected bool
	}{
		{"map", lang.MapValue{"key": lang.StringValue("value")}, true},
		{"empty map", lang.MapValue{}, true},
		{"list", lang.ListValue{lang.NumberValue(1)}, false},
		{"string", lang.StringValue("hello"), false},
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

// Test Number Type Checking
func TestIsInteger(t *testing.T) {
	_, fn := isInteger()

	tests := []struct {
		name     string
		input    lang.Value
		expected bool
	}{
		{"positive integer", lang.NumberValue(42), true},
		{"negative integer", lang.NumberValue(-5), true},
		{"zero", lang.NumberValue(0), true},
		{"positive float", lang.NumberValue(3.14), false},
		{"negative float", lang.NumberValue(-2.5), false},
		{"string", lang.StringValue("42"), false},
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

func TestIsFloat(t *testing.T) {
	_, fn := isFloat()

	tests := []struct {
		name     string
		input    lang.Value
		expected bool
	}{
		{"positive float", lang.NumberValue(3.14), true},
		{"negative float", lang.NumberValue(-2.5), true},
		{"integer", lang.NumberValue(42), false},
		{"zero", lang.NumberValue(0), false},
		{"string", lang.StringValue("3.14"), false},
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

func TestIsPositive(t *testing.T) {
	_, fn := isPositive()

	tests := []struct {
		name     string
		input    lang.Value
		expected bool
	}{
		{"positive number", lang.NumberValue(42), true},
		{"positive float", lang.NumberValue(3.14), true},
		{"zero", lang.NumberValue(0), false},
		{"negative number", lang.NumberValue(-5), false},
		{"string", lang.StringValue("42"), false},
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

func TestIsNegative(t *testing.T) {
	_, fn := isNegative()

	tests := []struct {
		name     string
		input    lang.Value
		expected bool
	}{
		{"negative number", lang.NumberValue(-5), true},
		{"negative float", lang.NumberValue(-3.14), true},
		{"zero", lang.NumberValue(0), false},
		{"positive number", lang.NumberValue(42), false},
		{"string", lang.StringValue("-5"), false},
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

func TestIsZero(t *testing.T) {
	_, fn := isZero()

	tests := []struct {
		name     string
		input    lang.Value
		expected bool
	}{
		{"zero", lang.NumberValue(0), true},
		{"positive number", lang.NumberValue(42), false},
		{"negative number", lang.NumberValue(-5), false},
		{"string", lang.StringValue("0"), false},
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

func TestIsEven(t *testing.T) {
	_, fn := isEven()

	tests := []struct {
		name     string
		input    lang.Value
		expected bool
	}{
		{"even positive", lang.NumberValue(42), true},
		{"even negative", lang.NumberValue(-4), true},
		{"zero", lang.NumberValue(0), true},
		{"odd positive", lang.NumberValue(43), false},
		{"odd negative", lang.NumberValue(-5), false},
		{"float", lang.NumberValue(3.14), false},
		{"string", lang.StringValue("42"), false},
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

func TestIsOdd(t *testing.T) {
	_, fn := isOdd()

	tests := []struct {
		name     string
		input    lang.Value
		expected bool
	}{
		{"odd positive", lang.NumberValue(43), true},
		{"odd negative", lang.NumberValue(-5), true},
		{"even positive", lang.NumberValue(42), false},
		{"even negative", lang.NumberValue(-4), false},
		{"zero", lang.NumberValue(0), false},
		{"float", lang.NumberValue(3.14), false},
		{"string", lang.StringValue("43"), false},
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

// Test String Type Checking
func TestIsNumericString(t *testing.T) {
	_, fn := isNumericString()

	tests := []struct {
		name     string
		input    lang.Value
		expected bool
	}{
		{"integer string", lang.StringValue("42"), true},
		{"float string", lang.StringValue("3.14"), true},
		{"negative string", lang.StringValue("-5"), true},
		{"scientific notation", lang.StringValue("1e10"), true},
		{"with whitespace", lang.StringValue("  42  "), true},
		{"non-numeric", lang.StringValue("hello"), false},
		{"empty string", lang.StringValue(""), false},
		{"mixed", lang.StringValue("42abc"), false},
		{"number type", lang.NumberValue(42), false},
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

func TestIsAlpha(t *testing.T) {
	_, fn := isAlpha()

	tests := []struct {
		name     string
		input    lang.Value
		expected bool
	}{
		{"lowercase", lang.StringValue("hello"), true},
		{"uppercase", lang.StringValue("HELLO"), true},
		{"mixed case", lang.StringValue("HeLLo"), true},
		{"with numbers", lang.StringValue("hello123"), false},
		{"with spaces", lang.StringValue("hello world"), false},
		{"empty string", lang.StringValue(""), false},
		{"special chars", lang.StringValue("hello!"), false},
		{"number type", lang.NumberValue(42), false},
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

func TestIsAlphaNumeric(t *testing.T) {
	_, fn := isAlphaNumeric()

	tests := []struct {
		name     string
		input    lang.Value
		expected bool
	}{
		{"letters only", lang.StringValue("hello"), true},
		{"numbers only", lang.StringValue("123"), true},
		{"mixed", lang.StringValue("hello123"), true},
		{"with spaces", lang.StringValue("hello 123"), false},
		{"empty string", lang.StringValue(""), false},
		{"special chars", lang.StringValue("hello!"), false},
		{"number type", lang.NumberValue(42), false},
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

func TestIsDigit(t *testing.T) {
	_, fn := isDigit()

	tests := []struct {
		name     string
		input    lang.Value
		expected bool
	}{
		{"digits only", lang.StringValue("123"), true},
		{"single digit", lang.StringValue("5"), true},
		{"with letters", lang.StringValue("123abc"), false},
		{"with decimal", lang.StringValue("12.3"), false},
		{"empty string", lang.StringValue(""), false},
		{"number type", lang.NumberValue(123), false},
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

func TestIsLower(t *testing.T) {
	_, fn := isLower()

	tests := []struct {
		name     string
		input    lang.Value
		expected bool
	}{
		{"lowercase", lang.StringValue("hello"), true},
		{"lowercase with numbers", lang.StringValue("hello123"), true},
		{"uppercase", lang.StringValue("HELLO"), false},
		{"mixed case", lang.StringValue("Hello"), false},
		{"no letters", lang.StringValue("123"), false},
		{"empty string", lang.StringValue(""), false},
		{"number type", lang.NumberValue(42), false},
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

func TestIsUpper(t *testing.T) {
	_, fn := isUpper()

	tests := []struct {
		name     string
		input    lang.Value
		expected bool
	}{
		{"uppercase", lang.StringValue("HELLO"), true},
		{"uppercase with numbers", lang.StringValue("HELLO123"), true},
		{"lowercase", lang.StringValue("hello"), false},
		{"mixed case", lang.StringValue("Hello"), false},
		{"no letters", lang.StringValue("123"), false},
		{"empty string", lang.StringValue(""), false},
		{"number type", lang.NumberValue(42), false},
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

func TestIsWhitespace(t *testing.T) {
	_, fn := isWhitespace()

	tests := []struct {
		name     string
		input    lang.Value
		expected bool
	}{
		{"spaces", lang.StringValue("   "), true},
		{"tabs", lang.StringValue("\t\t"), true},
		{"newlines", lang.StringValue("\n\n"), true},
		{"mixed whitespace", lang.StringValue(" \t\n "), true},
		{"with text", lang.StringValue(" hello "), false},
		{"empty string", lang.StringValue(""), false},
		{"number type", lang.NumberValue(42), false},
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

// Test Format Validation
func TestIsEmail(t *testing.T) {
	_, fn := isEmail()

	tests := []struct {
		name     string
		input    lang.Value
		expected bool
	}{
		{"valid email", lang.StringValue("test@example.com"), true},
		{"email with subdomain", lang.StringValue("user@mail.example.com"), true},
		{"email with numbers", lang.StringValue("user123@example.org"), true},
		{"invalid - no @", lang.StringValue("testexample.com"), false},
		{"invalid - no domain", lang.StringValue("test@"), false},
		{"invalid - no tld", lang.StringValue("test@example"), false},
		{"empty string", lang.StringValue(""), false},
		{"number type", lang.NumberValue(42), false},
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

func TestIsURL(t *testing.T) {
	_, fn := isURL()

	tests := []struct {
		name     string
		input    lang.Value
		expected bool
	}{
		{"http url", lang.StringValue("http://example.com"), true},
		{"https url", lang.StringValue("https://example.com"), true},
		{"url with path", lang.StringValue("https://example.com/path"), true},
		{"url with query", lang.StringValue("https://example.com/path?query=value"), true},
		{"invalid - no protocol", lang.StringValue("example.com"), false},
		{"invalid - ftp", lang.StringValue("ftp://example.com"), false},
		{"empty string", lang.StringValue(""), false},
		{"number type", lang.NumberValue(42), false},
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

func TestIsIPAddress(t *testing.T) {
	_, fn := isIPAddress()

	tests := []struct {
		name     string
		input    lang.Value
		expected bool
	}{
		{"valid IPv4", lang.StringValue("192.168.1.1"), true},
		{"valid IPv4 - localhost", lang.StringValue("127.0.0.1"), true},
		{"valid IPv6 - full", lang.StringValue("2001:0db8:85a3:0000:0000:8a2e:0370:7334"), true},
		{"valid IPv6 - localhost", lang.StringValue("::1"), true},
		{"valid IPv6 - any", lang.StringValue("::"), true},
		{"invalid IPv4 - out of range", lang.StringValue("256.1.1.1"), false},
		{"invalid IPv4 - incomplete", lang.StringValue("192.168.1"), false},
		{"invalid format", lang.StringValue("not.an.ip"), false},
		{"empty string", lang.StringValue(""), false},
		{"number type", lang.NumberValue(42), false},
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

func TestIsUUID(t *testing.T) {
	_, fn := isUUID()

	tests := []struct {
		name     string
		input    lang.Value
		expected bool
	}{
		{"valid UUID", lang.StringValue("550e8400-e29b-41d4-a716-446655440000"), true},
		{"valid UUID - lowercase", lang.StringValue("550e8400-e29b-41d4-a716-446655440000"), true},
		{"valid UUID - uppercase", lang.StringValue("550E8400-E29B-41D4-A716-446655440000"), true},
		{"invalid - no dashes", lang.StringValue("550e8400e29b41d4a716446655440000"), false},
		{"invalid - wrong length", lang.StringValue("550e8400-e29b-41d4-a716-44665544000"), false},
		{"invalid - wrong format", lang.StringValue("550e8400-e29b-41d4-a716"), false},
		{"empty string", lang.StringValue(""), false},
		{"number type", lang.NumberValue(42), false},
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

func TestIsJSON(t *testing.T) {
	_, fn := isJSON()

	tests := []struct {
		name     string
		input    lang.Value
		expected bool
	}{
		{"object", lang.StringValue(`{"key": "value"}`), true},
		{"array", lang.StringValue(`[1, 2, 3]`), true},
		{"string", lang.StringValue(`"hello"`), true},
		{"number", lang.StringValue("42"), true},
		{"boolean true", lang.StringValue("true"), true},
		{"boolean false", lang.StringValue("false"), true},
		{"null", lang.StringValue("null"), true},
		{"float", lang.StringValue("3.14"), true},
		{"invalid object", lang.StringValue(`{key: "value"}`), false},
		{"invalid array", lang.StringValue(`[1, 2, 3`), false},
		{"empty string", lang.StringValue(""), false},
		{"plain text", lang.StringValue("hello"), false},
		{"number type", lang.NumberValue(42), false},
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

func TestIsBase64(t *testing.T) {
	_, fn := isBase64()

	tests := []struct {
		name     string
		input    lang.Value
		expected bool
	}{
		{"valid base64", lang.StringValue("SGVsbG8gV29ybGQ="), true},
		{"valid base64 no padding", lang.StringValue("SGVsbG8"), true},
		{"valid base64 single padding", lang.StringValue("SGVsbG8g"), true},
		{"invalid - wrong chars", lang.StringValue("SGVsbG8@"), false},
		{"empty string", lang.StringValue(""), true}, // Empty string is valid base64
		{"number type", lang.NumberValue(42), false},
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

func TestIsHex(t *testing.T) {
	_, fn := isHex()

	tests := []struct {
		name     string
		input    lang.Value
		expected bool
	}{
		{"hex lowercase", lang.StringValue("deadbeef"), true},
		{"hex uppercase", lang.StringValue("DEADBEEF"), true},
		{"hex mixed", lang.StringValue("DeAdBeEf"), true},
		{"hex with numbers", lang.StringValue("123abc"), true},
		{"invalid - with g", lang.StringValue("deadbeeG"), false},
		{"invalid - with symbols", lang.StringValue("dead-beef"), false},
		{"empty string", lang.StringValue(""), false},
		{"number type", lang.NumberValue(42), false},
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

// Test Collection Type Checking
func TestHasLength(t *testing.T) {
	_, fn := hasLength()

	tests := []struct {
		name     string
		input    lang.Value
		expected bool
	}{
		{"string", lang.StringValue("hello"), true},
		{"empty string", lang.StringValue(""), true},
		{"list", lang.ListValue{lang.NumberValue(1)}, true},
		{"empty list", lang.ListValue{}, true},
		{"map", lang.MapValue{"key": lang.StringValue("value")}, true},
		{"empty map", lang.MapValue{}, true},
		{"number", lang.NumberValue(42), false},
		{"boolean", lang.BoolValue(true), false},
		{"null", nil, false},
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

// Test Range Checking
func TestIsInRange(t *testing.T) {
	_, fn := isInRange()

	tests := []struct {
		name     string
		value    lang.Value
		min      lang.Value
		max      lang.Value
		expected bool
		hasError bool
	}{
		{"in range", lang.NumberValue(5), lang.NumberValue(1), lang.NumberValue(10), true, false},
		{"at min", lang.NumberValue(1), lang.NumberValue(1), lang.NumberValue(10), true, false},
		{"at max", lang.NumberValue(10), lang.NumberValue(1), lang.NumberValue(10), true, false},
		{"below range", lang.NumberValue(0), lang.NumberValue(1), lang.NumberValue(10), false, false},
		{"above range", lang.NumberValue(11), lang.NumberValue(1), lang.NumberValue(10), false, false},
		{"string convertible", lang.StringValue("5"), lang.NumberValue(1), lang.NumberValue(10), true, false},
		{"invalid min > max", lang.NumberValue(5), lang.NumberValue(10), lang.NumberValue(1), false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{tt.value, tt.min, tt.max})
			if tt.hasError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if bool(result.(lang.BoolValue)) != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, bool(result.(lang.BoolValue)))
			}
		})
	}

	// Test wrong number of arguments
	_, err := fn([]lang.Value{lang.NumberValue(5)})
	if err == nil {
		t.Error("expected error for wrong number of arguments")
	}
}

func TestIsLengthInRange(t *testing.T) {
	_, fn := isLengthInRange()

	tests := []struct {
		name     string
		value    lang.Value
		min      lang.Value
		max      lang.Value
		expected bool
		hasError bool
	}{
		{"string in range", lang.StringValue("hello"), lang.NumberValue(3), lang.NumberValue(10), true, false},
		{"string at min", lang.StringValue("hi"), lang.NumberValue(2), lang.NumberValue(10), true, false},
		{"string at max", lang.StringValue("hello"), lang.NumberValue(1), lang.NumberValue(5), true, false},
		{"string too short", lang.StringValue("hi"), lang.NumberValue(5), lang.NumberValue(10), false, false},
		{"string too long", lang.StringValue("hello world"), lang.NumberValue(1), lang.NumberValue(5), false, false},
		{"list in range", lang.ListValue{lang.NumberValue(1), lang.NumberValue(2)}, lang.NumberValue(1), lang.NumberValue(5), true, false},
		{"map in range", lang.MapValue{"a": lang.NumberValue(1), "b": lang.NumberValue(2)}, lang.NumberValue(1), lang.NumberValue(5), true, false},
		{"invalid min > max", lang.StringValue("hello"), lang.NumberValue(10), lang.NumberValue(1), false, true},
		{"invalid negative min", lang.StringValue("hello"), lang.NumberValue(-1), lang.NumberValue(10), false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{tt.value, tt.min, tt.max})
			if tt.hasError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if bool(result.(lang.BoolValue)) != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, bool(result.(lang.BoolValue)))
			}
		})
	}
}

// Test Type Conversion Checking
func TestCanConvertToNumber(t *testing.T) {
	_, fn := canConvertToNumber()

	tests := []struct {
		name     string
		input    lang.Value
		expected bool
	}{
		{"number", lang.NumberValue(42), true},
		{"numeric string", lang.StringValue("42"), true},
		{"float string", lang.StringValue("3.14"), true},
		{"negative string", lang.StringValue("-5"), true},
		{"boolean true", lang.BoolValue(true), true},
		{"boolean false", lang.BoolValue(false), true},
		{"non-numeric string", lang.StringValue("hello"), false},
		{"empty string", lang.StringValue(""), false},
		{"null", nil, false},
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

func TestCanConvertToString(t *testing.T) {
	_, fn := canConvertToString()

	tests := []struct {
		name     string
		input    lang.Value
		expected bool
	}{
		{"string", lang.StringValue("hello"), true},
		{"number", lang.NumberValue(42), true},
		{"boolean", lang.BoolValue(true), true},
		{"list", lang.ListValue{lang.NumberValue(1)}, true},
		{"map", lang.MapValue{"key": lang.StringValue("value")}, true},
		{"null", nil, false},
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

func TestCanConvertToBool(t *testing.T) {
	_, fn := canConvertToBool()

	tests := []struct {
		name     string
		input    lang.Value
		expected bool
	}{
		{"boolean", lang.BoolValue(true), true},
		{"string", lang.StringValue("hello"), true},
		{"number", lang.NumberValue(42), true},
		{"list", lang.ListValue{lang.NumberValue(1)}, true},
		{"map", lang.MapValue{"key": lang.StringValue("value")}, true},
		{"null", nil, false},
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

// Test Comparison Functions
func TestAreEqual(t *testing.T) {
	_, fn := areEqual()

	tests := []struct {
		name     string
		a        lang.Value
		b        lang.Value
		expected bool
	}{
		{"equal numbers", lang.NumberValue(42), lang.NumberValue(42), true},
		{"equal strings", lang.StringValue("hello"), lang.StringValue("hello"), true},
		{"equal booleans", lang.BoolValue(true), lang.BoolValue(true), true},
		{"both null", nil, nil, true},
		{"equal lists", lang.ListValue{lang.NumberValue(1), lang.NumberValue(2)}, lang.ListValue{lang.NumberValue(1), lang.NumberValue(2)}, true},
		{"equal maps", lang.MapValue{"key": lang.StringValue("value")}, lang.MapValue{"key": lang.StringValue("value")}, true},
		{"different numbers", lang.NumberValue(42), lang.NumberValue(43), false},
		{"different types", lang.NumberValue(42), lang.StringValue("42"), false},
		{"null vs non-null", nil, lang.StringValue("hello"), false},
		{"different lists", lang.ListValue{lang.NumberValue(1)}, lang.ListValue{lang.NumberValue(2)}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{tt.a, tt.b})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if bool(result.(lang.BoolValue)) != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, bool(result.(lang.BoolValue)))
			}
		})
	}

	// Test wrong number of arguments
	_, err := fn([]lang.Value{lang.NumberValue(42)})
	if err == nil {
		t.Error("expected error for wrong number of arguments")
	}
}

func TestAreStrictEqual(t *testing.T) {
	_, fn := areStrictEqual()

	tests := []struct {
		name     string
		a        lang.Value
		b        lang.Value
		expected bool
	}{
		{"equal numbers same type", lang.NumberValue(42), lang.NumberValue(42), true},
		{"equal strings same type", lang.StringValue("hello"), lang.StringValue("hello"), true},
		{"different types", lang.NumberValue(42), lang.StringValue("42"), false},
		{"equal values different types", lang.BoolValue(true), lang.NumberValue(1), false},
		{"both null", nil, nil, true},
		{"null vs string", nil, lang.StringValue("null"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{tt.a, tt.b})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if bool(result.(lang.BoolValue)) != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, bool(result.(lang.BoolValue)))
			}
		})
	}
}

// Test Alias Functions
func TestIsArray(t *testing.T) {
	_, fn := isArray()

	// This should behave exactly like isList
	tests := []struct {
		name     string
		input    lang.Value
		expected bool
	}{
		{"list", lang.ListValue{lang.NumberValue(1), lang.NumberValue(2)}, true},
		{"empty list", lang.ListValue{}, true},
		{"string", lang.StringValue("hello"), false},
		{"map", lang.MapValue{"key": lang.StringValue("value")}, false},
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

func TestIsObject(t *testing.T) {
	_, fn := isObject()

	// This should behave exactly like isMap
	tests := []struct {
		name     string
		input    lang.Value
		expected bool
	}{
		{"map", lang.MapValue{"key": lang.StringValue("value")}, true},
		{"empty map", lang.MapValue{}, true},
		{"list", lang.ListValue{lang.NumberValue(1)}, false},
		{"string", lang.StringValue("hello"), false},
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

// Test Special Number Cases
func TestIsNaN(t *testing.T) {
	_, fn := isNaN()

	tests := []struct {
		name     string
		input    lang.Value
		expected bool
	}{
		{"normal number", lang.NumberValue(42), false},
		{"zero", lang.NumberValue(0), false},
		{"string", lang.StringValue("hello"), false},
		{"boolean", lang.BoolValue(true), false},
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

func TestIsInfinite(t *testing.T) {
	_, fn := isInfinite()

	tests := []struct {
		name     string
		input    lang.Value
		expected bool
	}{
		{"normal number", lang.NumberValue(42), false},
		{"zero", lang.NumberValue(0), false},
		{"large number", lang.NumberValue(1e100), false},
		{"string", lang.StringValue("hello"), false},
		{"boolean", lang.BoolValue(true), false},
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

func TestIsFinite(t *testing.T) {
	_, fn := isFinite()

	tests := []struct {
		name     string
		input    lang.Value
		expected bool
	}{
		{"normal number", lang.NumberValue(42), true},
		{"zero", lang.NumberValue(0), true},
		{"negative number", lang.NumberValue(-5), true},
		{"float", lang.NumberValue(3.14), true},
		{"string", lang.StringValue("hello"), false},
		{"boolean", lang.BoolValue(true), false},
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

func TestArgumentErrors(t *testing.T) {
	functions := []struct {
		fn           func() (string, lang.Function)
		expectedArgs int
	}{
		{typeOf, 1},
		{isNull, 1},
		{isDefined, 1},
		{isEmpty, 1},
		{isNotEmpty, 1},
		{isBool, 1},
		{isNumber, 1},
		{isString, 1},
		{isList, 1},
		{isMap, 1},
		{isArray, 1},
		{isObject, 1},
		{isInteger, 1},
		{isFloat, 1},
		{isPositive, 1},
		{isNegative, 1},
		{isZero, 1},
		{isEven, 1},
		{isOdd, 1},
		{isNaN, 1},
		{isInfinite, 1},
		{isFinite, 1},
		{isNumericString, 1},
		{isAlphaNumeric, 1},
		{isDigit, 1},
		{isLower, 1},
		{isUpper, 1},
		{isWhitespace, 1},
		{isAlpha, 1},
		{isEmail, 1},
		{isURL, 1},
		{isIPAddress, 1},
		{isUUID, 1},
		{isJSON, 1},
		{isBase64, 1},
		{isHex, 1},
		{hasLength, 1},
		{canConvertToNumber, 1},
		{canConvertToString, 1},
		{canConvertToBool, 1},
		{areEqual, 2},
		{areStrictEqual, 2},
		{isInRange, 3},
		{isLengthInRange, 3},
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

func TestExport(t *testing.T) {
	functions := Export()

	expectedFunctions := []string{
		"type", "isNull", "isDefined", "isEmpty", "isNotEmpty",
		"isBool", "isNumber", "isString", "isList", "isMap", "isArray", "isObject",
		"isInteger", "isFloat", "isPositive", "isNegative", "isZero", "isEven", "isOdd",
		"isNan", "isInfinite", "isFinite",
		"isNumericString", "isAlpha", "isAlphanumeric", "isDigit", "isLower", "isUpper", "isWhitespace",
		"isEmail", "isUrl", "isIpAddress", "isUUID", "isJSON", "isBase64", "isHex",
		"hasLength", "isInRange", "isLengthInRange",
		"canConvertToNumber", "canConvertToString", "canConvertToBool",
		"areEqual", "areStrictEqual",
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

// Test Edge Cases
func TestEdgeCases(t *testing.T) {
	t.Run("base64_empty_string", func(t *testing.T) {
		_, fn := isBase64()
		result, err := fn([]lang.Value{lang.StringValue("")})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// Empty string should be valid base64 and have length 0 (multiple of 4)
		if !bool(result.(lang.BoolValue)) {
			t.Error("empty string should be valid base64")
		}
	})

	t.Run("deep_equal_nested", func(t *testing.T) {
		_, fn := areEqual()

		// Test nested structures
		nested1 := lang.MapValue{
			"outer": lang.MapValue{
				"inner": lang.ListValue{lang.NumberValue(1), lang.NumberValue(2)},
			},
		}
		nested2 := lang.MapValue{
			"outer": lang.MapValue{
				"inner": lang.ListValue{lang.NumberValue(1), lang.NumberValue(2)},
			},
		}
		nested3 := lang.MapValue{
			"outer": lang.MapValue{
				"inner": lang.ListValue{lang.NumberValue(1), lang.NumberValue(3)},
			},
		}

		result1, err := fn([]lang.Value{nested1, nested2})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !bool(result1.(lang.BoolValue)) {
			t.Error("deeply equal nested structures should be equal")
		}

		result2, err := fn([]lang.Value{nested1, nested3})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if bool(result2.(lang.BoolValue)) {
			t.Error("deeply unequal nested structures should not be equal")
		}
	})
}

// Benchmark Tests
func BenchmarkTypeOf(b *testing.B) {
	_, fn := typeOf()
	input := []lang.Value{lang.StringValue("test")}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fn(input)
	}
}

func BenchmarkIsEmail(b *testing.B) {
	_, fn := isEmail()
	input := []lang.Value{lang.StringValue("test@example.com")}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fn(input)
	}
}

func BenchmarkDeepEqual(b *testing.B) {
	_, fn := areEqual()

	list1 := lang.ListValue{lang.NumberValue(1), lang.NumberValue(2), lang.NumberValue(3)}
	list2 := lang.ListValue{lang.NumberValue(1), lang.NumberValue(2), lang.NumberValue(3)}
	input := []lang.Value{list1, list2}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fn(input)
	}
}
