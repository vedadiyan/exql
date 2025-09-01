package json

import (
	"reflect"
	"testing"

	"github.com/vedadiyan/exql/lang"
)

func TestParse(t *testing.T) {
	_, fn := parse()

	tests := []struct {
		name     string
		input    string
		expected interface{}
		hasError bool
	}{
		{
			name:     "simple object",
			input:    `{"name": "John", "age": 30}`,
			expected: lang.MapValue{"name": lang.StringValue("John"), "age": lang.NumberValue(30)},
			hasError: false,
		},
		{
			name:     "array",
			input:    `[1, 2, 3]`,
			expected: lang.ListValue{lang.NumberValue(1), lang.NumberValue(2), lang.NumberValue(3)},
			hasError: false,
		},
		{
			name:     "string",
			input:    `"hello"`,
			expected: lang.StringValue("hello"),
			hasError: false,
		},
		{
			name:     "number",
			input:    `42`,
			expected: lang.NumberValue(42),
			hasError: false,
		},
		{
			name:     "boolean",
			input:    `true`,
			expected: lang.BoolValue(true),
			hasError: false,
		},
		{
			name:     "null",
			input:    `null`,
			expected: nil,
			hasError: false,
		},
		{
			name:     "invalid json",
			input:    `{"invalid": }`,
			hasError: true,
		},
		{
			name:     "wrong args",
			input:    "",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var args []lang.Value
			if tt.name == "wrong args" {
				args = []lang.Value{}
			} else {
				args = []lang.Value{lang.StringValue(tt.input)}
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

			switch expected := tt.expected.(type) {
			case nil:
				if result != nil {
					t.Errorf("Expected nil, got %v", result)
				}
			case lang.MapValue:
				resultMap, ok := result.(lang.MapValue)
				if !ok {
					t.Errorf("Expected MapValue, got %T", result)
					return
				}
				if len(resultMap) != len(expected) {
					t.Errorf("Map length mismatch: expected %d, got %d", len(expected), len(resultMap))
					return
				}
				for key, expectedVal := range expected {
					if actualVal, exists := resultMap[key]; !exists {
						t.Errorf("Missing key %s", key)
					} else if actualVal != expectedVal {
						t.Errorf("Key %s: expected %v, got %v", key, expectedVal, actualVal)
					}
				}
			case lang.ListValue:
				resultList, ok := result.(lang.ListValue)
				if !ok {
					t.Errorf("Expected ListValue, got %T", result)
					return
				}
				if len(resultList) != len(expected) {
					t.Errorf("List length mismatch: expected %d, got %d", len(expected), len(resultList))
					return
				}
				for i, expectedVal := range expected {
					if resultList[i] != expectedVal {
						t.Errorf("Index %d: expected %v, got %v", i, expectedVal, resultList[i])
					}
				}
			default:
				if result != expected {
					t.Errorf("Expected %v, got %v", expected, result)
				}
			}
		})
	}
}

func TestSstring(t *testing.T) {
	_, fn := sstring()

	tests := []struct {
		name     string
		data     lang.Value
		pretty   *bool
		expected string
		hasError bool
	}{
		{
			name:     "simple object",
			data:     lang.MapValue{"name": lang.StringValue("John")},
			expected: `{"name":"John"}`,
			hasError: false,
		},
		{
			name:     "pretty printed",
			data:     lang.MapValue{"name": lang.StringValue("John")},
			pretty:   boolPtr(true),
			expected: "{\n  \"name\": \"John\"\n}",
			hasError: false,
		},
		{
			name:     "array",
			data:     lang.ListValue{lang.NumberValue(1), lang.NumberValue(2)},
			expected: `[1,2]`,
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var args []lang.Value
			args = []lang.Value{tt.data}
			if tt.pretty != nil {
				args = append(args, lang.BoolValue(*tt.pretty))
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

			if string(result.(lang.StringValue)) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(result.(lang.StringValue)))
			}
		})
	}
}

func TestValid(t *testing.T) {
	_, fn := valid()

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid object", `{"name": "John"}`, true},
		{"valid array", `[1, 2, 3]`, true},
		{"valid string", `"hello"`, true},
		{"invalid json", `{"invalid": }`, false},
		{"empty string", ``, false},
		{"malformed", `{name: "John"}`, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.StringValue(tt.input)})
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

func TestGet(t *testing.T) {
	_, fn := get()

	jsonData := `{"user": {"name": "John", "age": 30}, "items": [1, 2, 3]}`
	objectData := lang.MapValue{
		"user": lang.MapValue{
			"name": lang.StringValue("John"),
			"age":  lang.NumberValue(30),
		},
		"items": lang.ListValue{lang.NumberValue(1), lang.NumberValue(2), lang.NumberValue(3)},
	}

	tests := []struct {
		name     string
		data     lang.Value
		path     string
		expected lang.Value
		hasError bool
	}{
		{
			name:     "get nested object property",
			data:     lang.StringValue(jsonData),
			path:     "user.name",
			expected: lang.StringValue("John"),
			hasError: false,
		},
		{
			name:     "get array element",
			data:     lang.StringValue(jsonData),
			path:     "items.1",
			expected: lang.NumberValue(2),
			hasError: false,
		},
		{
			name:     "get root property",
			data:     objectData,
			path:     "user",
			expected: lang.MapValue{"name": lang.StringValue("John"), "age": lang.NumberValue(30)},
			hasError: false,
		},
		{
			name:     "non-existent path",
			data:     lang.StringValue(jsonData),
			path:     "nonexistent",
			expected: nil,
			hasError: false,
		},
		{
			name:     "empty path",
			data:     lang.StringValue(`"hello"`),
			path:     "",
			expected: lang.StringValue("hello"),
			hasError: false,
		},
		{
			name:     "invalid json",
			data:     lang.StringValue(`{"invalid": }`),
			path:     "test",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{tt.data, lang.StringValue(tt.path)})
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

			if !reflect.DeepEqual(result, tt.expected) {
				if tt.expected == nil {
					if result != nil {
						t.Errorf("Expected nil, got %v", result)
					}
				} else {
					t.Errorf("Expected %v, got %v", tt.expected, result)
				}
			}
		})
	}
}

func TestSet(t *testing.T) {
	_, fn := set()

	tests := []struct {
		name     string
		data     string
		path     string
		value    lang.Value
		hasError bool
	}{
		{
			name:     "set existing property",
			data:     `{"name": "John"}`,
			path:     "name",
			value:    lang.StringValue("Jane"),
			hasError: false,
		},
		{
			name:     "set new property",
			data:     `{"name": "John"}`,
			path:     "age",
			value:    lang.NumberValue(25),
			hasError: false,
		},
		{
			name:     "set nested property",
			data:     `{"user": {"name": "John"}}`,
			path:     "user.age",
			value:    lang.NumberValue(30),
			hasError: false,
		},
		{
			name:     "set array element",
			data:     `{"items": [1, 2, 3]}`,
			path:     "items.1",
			value:    lang.NumberValue(5),
			hasError: false,
		},
		{
			name:     "invalid json",
			data:     `{"invalid": }`,
			path:     "test",
			value:    lang.StringValue("value"),
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.StringValue(tt.data), lang.StringValue(tt.path), tt.value})
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

			if result == nil {
				t.Errorf("Expected non-nil result")
			}
		})
	}
}

func TestRemove(t *testing.T) {
	_, fn := remove()

	tests := []struct {
		name     string
		data     string
		path     string
		hasError bool
	}{
		{
			name:     "remove existing property",
			data:     `{"name": "John", "age": 30}`,
			path:     "name",
			hasError: false,
		},
		{
			name:     "remove array element",
			data:     `{"items": [1, 2, 3]}`,
			path:     "items.1",
			hasError: false,
		},
		{
			name:     "remove non-existent property",
			data:     `{"name": "John"}`,
			path:     "nonexistent",
			hasError: false,
		},
		{
			name:     "invalid json",
			data:     `{"invalid": }`,
			path:     "test",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.StringValue(tt.data), lang.StringValue(tt.path)})
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

			if result == nil {
				t.Errorf("Expected non-nil result")
			}
		})
	}
}

func TestHas(t *testing.T) {
	_, fn := has()

	jsonData := `{"user": {"name": "John"}, "items": [1, 2, 3]}`

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{"existing property", "user.name", true},
		{"existing array index", "items.1", true},
		{"non-existent property", "nonexistent", false},
		{"invalid array index", "items.10", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.StringValue(jsonData), lang.StringValue(tt.path)})
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

func TestKeys(t *testing.T) {
	_, fn := keys()

	tests := []struct {
		name     string
		data     string
		expected []string
		hasError bool
	}{
		{
			name:     "object keys",
			data:     `{"name": "John", "age": 30, "city": "NYC"}`,
			expected: []string{"name", "age", "city"},
			hasError: false,
		},
		{
			name:     "empty object",
			data:     `{}`,
			expected: []string{},
			hasError: false,
		},
		{
			name:     "array input",
			data:     `[1, 2, 3]`,
			hasError: true,
		},
		{
			name:     "invalid json",
			data:     `{"invalid": }`,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.StringValue(tt.data)})
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

			keys, ok := result.(lang.ListValue)
			if !ok {
				t.Errorf("Expected ListValue, got %T", result)
				return
			}

			if len(keys) != len(tt.expected) {
				t.Errorf("Expected %d keys, got %d", len(tt.expected), len(keys))
				return
			}

			keyMap := make(map[string]bool)
			for _, key := range keys {
				keyMap[string(key.(lang.StringValue))] = true
			}

			for _, expectedKey := range tt.expected {
				if !keyMap[expectedKey] {
					t.Errorf("Missing expected key: %s", expectedKey)
				}
			}
		})
	}
}

func TestValues(t *testing.T) {
	_, fn := values()

	tests := []struct {
		name     string
		data     string
		hasError bool
	}{
		{
			name:     "object values",
			data:     `{"name": "John", "age": 30}`,
			hasError: false,
		},
		{
			name:     "array values",
			data:     `[1, 2, 3]`,
			hasError: false,
		},
		{
			name:     "string input",
			data:     `"hello"`,
			hasError: true,
		},
		{
			name:     "invalid json",
			data:     `{"invalid": }`,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.StringValue(tt.data)})
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

			if _, ok := result.(lang.ListValue); !ok {
				t.Errorf("Expected ListValue, got %T", result)
			}
		})
	}
}

func TestLength(t *testing.T) {
	_, fn := length()

	tests := []struct {
		name     string
		data     string
		expected float64
		hasError bool
	}{
		{
			name:     "object length",
			data:     `{"name": "John", "age": 30}`,
			expected: 2,
			hasError: false,
		},
		{
			name:     "array length",
			data:     `[1, 2, 3, 4, 5]`,
			expected: 5,
			hasError: false,
		},
		{
			name:     "string length",
			data:     `"hello"`,
			expected: 5,
			hasError: false,
		},
		{
			name:     "number input",
			data:     `42`,
			hasError: true,
		},
		{
			name:     "invalid json",
			data:     `{"invalid": }`,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.StringValue(tt.data)})
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

			if float64(result.(lang.NumberValue)) != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, float64(result.(lang.NumberValue)))
			}
		})
	}
}

func TestMerge(t *testing.T) {
	_, fn := merge()

	tests := []struct {
		name     string
		inputs   []string
		hasError bool
	}{
		{
			name:     "merge two objects",
			inputs:   []string{`{"a": 1}`, `{"b": 2}`},
			hasError: false,
		},
		{
			name:     "merge multiple objects",
			inputs:   []string{`{"a": 1}`, `{"b": 2}`, `{"c": 3}`},
			hasError: false,
		},
		{
			name:     "overwrite properties",
			inputs:   []string{`{"a": 1, "b": 2}`, `{"a": 3}`},
			hasError: false,
		},
		{
			name:     "single argument",
			inputs:   []string{`{"a": 1}`},
			hasError: true,
		},
		{
			name:     "non-object input",
			inputs:   []string{`{"a": 1}`, `[1, 2, 3]`},
			hasError: true,
		},
		{
			name:     "invalid json",
			inputs:   []string{`{"a": 1}`, `{"invalid": }`},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := make([]lang.Value, len(tt.inputs))
			for i, input := range tt.inputs {
				args[i] = lang.StringValue(input)
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

			if _, ok := result.(lang.MapValue); !ok {
				t.Errorf("Expected MapValue, got %T", result)
			}
		})
	}
}

func TestTtype(t *testing.T) {
	_, fn := ttype()

	tests := []struct {
		name     string
		data     string
		expected string
	}{
		{"null", `null`, "null"},
		{"boolean true", `true`, "boolean"},
		{"boolean false", `false`, "boolean"},
		{"number", `42`, "number"},
		{"string", `"hello"`, "string"},
		{"array", `[1, 2, 3]`, "array"},
		{"object", `{"name": "John"}`, "object"},
		{"invalid json", `{"invalid": }`, "invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.StringValue(tt.data)})
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if string(result.(lang.StringValue)) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(result.(lang.StringValue)))
			}
		})
	}
}

func TestConvertJSONToValue(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected lang.Value
	}{
		{"nil", nil, nil},
		{"bool", true, lang.BoolValue(true)},
		{"float64", 42.0, lang.NumberValue(42.0)},
		{"string", "hello", lang.StringValue("hello")},
		{"slice", []interface{}{1.0, 2.0}, lang.ListValue{lang.NumberValue(1), lang.NumberValue(2)}},
		{"map", map[string]interface{}{"key": "value"}, lang.MapValue{"key": lang.StringValue("value")}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertJSONToValue(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestConvertValueToJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    lang.Value
		expected interface{}
	}{
		{"nil", nil, nil},
		{"bool", lang.BoolValue(true), true},
		{"number", lang.NumberValue(42), 42.0},
		{"string", lang.StringValue("hello"), "hello"},
		{"list", lang.ListValue{lang.NumberValue(1), lang.NumberValue(2)}, []interface{}{1.0, 2.0}},
		{"map", lang.MapValue{"key": lang.StringValue("value")}, map[string]interface{}{"key": "value"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertValueToJSON(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestJsonGetByPath(t *testing.T) {
	data := map[string]interface{}{
		"user": map[string]interface{}{
			"name": "John",
			"age":  30.0,
		},
		"items": []interface{}{1.0, 2.0, 3.0},
	}

	tests := []struct {
		name     string
		path     string
		expected lang.Value
	}{
		{"empty path", "", convertJSONToValue(data)},
		{"simple property", "user", convertJSONToValue(data["user"])},
		{"nested property", "user.name", lang.StringValue("John")},
		{"array index", "items.1", lang.NumberValue(2)},
		{"non-existent", "nonexistent", nil},
		{"invalid array index", "items.10", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := jsonGetByPath(data, tt.path)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestArgumentErrors(t *testing.T) {
	functions := []struct {
		fn           func() (string, lang.Function)
		expectedArgs int
	}{
		{parse, 1},
		{get, 2},
		{set, 3},
		{remove, 2},
		{has, 2},
		{keys, 1},
		{values, 1},
		{length, 1},
		{ttype, 1},
	}

	for _, tf := range functions {
		name, fn := tf.fn()
		t.Run(name+"_wrong_args", func(t *testing.T) {
			args := make([]lang.Value, tf.expectedArgs-1)
			for i := range args {
				args[i] = lang.StringValue("test")
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
		"parse", "string", "valid", "get", "set", "delete", "has",
		"keys", "values", "length", "merge", "type",
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

func BenchmarkParse(b *testing.B) {
	_, fn := parse()
	jsonStr := `{"name": "John", "age": 30, "items": [1, 2, 3]}`
	args := []lang.Value{lang.StringValue(jsonStr)}

	for i := 0; i < b.N; i++ {
		fn(args)
	}
}

func BenchmarkGet(b *testing.B) {
	_, fn := get()
	jsonStr := `{"user": {"name": "John", "details": {"age": 30}}}`
	args := []lang.Value{lang.StringValue(jsonStr), lang.StringValue("user.details.age")}

	for i := 0; i < b.N; i++ {
		fn(args)
	}
}

func boolPtr(b bool) *bool {
	return &b
}
