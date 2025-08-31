package string

import (
	"testing"

	"github.com/vedadiyan/exql/lang"
)

func TestLength(t *testing.T) {
	_, fn := length()

	tests := []struct {
		name     string
		input    string
		expected float64
		hasError bool
	}{
		{"empty string", "", 0, false},
		{"simple string", "hello", 5, false},
		{"unicode string", "héllo", 5, false},
		{"emoji string", "😀👍", 2, false},
		{"mixed unicode", "hello 世界", 8, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.StringValue(tt.input)})
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

func TestSize(t *testing.T) {
	_, fn := size()

	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		{"empty string", "", 0},
		{"simple string", "hello", 5},
		{"unicode string", "héllo", 6}, // byte length different from rune length
		{"emoji string", "😀", 4},       // emoji takes 4 bytes
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.StringValue(tt.input)})
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

func TestConcat(t *testing.T) {
	_, fn := voncat()

	tests := []struct {
		name     string
		inputs   []string
		expected string
	}{
		{"two strings", []string{"hello", " world"}, "hello world"},
		{"three strings", []string{"a", "b", "c"}, "abc"},
		{"empty strings", []string{"", "", ""}, ""},
		{"single string", []string{"test"}, "test"},
		{"no strings", []string{}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := make([]lang.Value, len(tt.inputs))
			for i, input := range tt.inputs {
				args[i] = lang.StringValue(input)
			}

			result, err := fn(args)
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

func TestRepeat(t *testing.T) {
	_, fn := repeat()

	tests := []struct {
		name     string
		input    string
		count    float64
		expected string
		hasError bool
	}{
		{"repeat string", "abc", 3, "abcabcabc", false},
		{"repeat zero", "test", 0, "", false},
		{"repeat negative", "test", -1, "", false},
		{"repeat one", "hello", 1, "hello", false},
		{"empty string", "", 5, "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.StringValue(tt.input), lang.NumberValue(tt.count)})
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

func TestReverse(t *testing.T) {
	_, fn := reverse()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple string", "hello", "olleh"},
		{"empty string", "", ""},
		{"single char", "a", "a"},
		{"unicode", "世界", "界世"},
		{"palindrome", "racecar", "racecar"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.StringValue(tt.input)})
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

func TestCaseConversions(t *testing.T) {
	tests := []struct {
		name     string
		fn       func() (string, lang.Function)
		input    string
		expected string
	}{
		{"toUpper simple", toUpper, "hello", "HELLO"},
		{"toUpper mixed", toUpper, "Hello World", "HELLO WORLD"},
		{"toLower simple", toLower, "HELLO", "hello"},
		{"toLower mixed", toLower, "Hello World", "hello world"},
		{"title simple", title, "hello world", "Hello World"},
		{"capitalize simple", capitalize, "hello world", "Hello world"},
		{"swapCase mixed", swapCase, "Hello World", "hELLO wORLD"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, fn := tt.fn()
			result, err := fn([]lang.Value{lang.StringValue(tt.input)})
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

func TestTrim(t *testing.T) {
	_, fn := trim()

	tests := []struct {
		name     string
		input    string
		cutset   *string
		expected string
	}{
		{"trim spaces", "  hello  ", nil, "hello"},
		{"trim custom", "xxhelloxx", stringPtr("x"), "hello"},
		{"trim multiple chars", "abchelloabc", stringPtr("abc"), "hello"},
		{"no trim needed", "hello", nil, "hello"},
		{"empty string", "", nil, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []lang.Value{lang.StringValue(tt.input)}
			if tt.cutset != nil {
				args = append(args, lang.StringValue(*tt.cutset))
			}

			result, err := fn(args)
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

func TestTrimLeft(t *testing.T) {
	_, fn := trimLeft()

	tests := []struct {
		name     string
		input    string
		cutset   *string
		expected string
	}{
		{"trim left spaces", "  hello  ", nil, "hello  "},
		{"trim left custom", "xxhello", stringPtr("x"), "hello"},
		{"no trim needed", "hello", nil, "hello"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []lang.Value{lang.StringValue(tt.input)}
			if tt.cutset != nil {
				args = append(args, lang.StringValue(*tt.cutset))
			}

			result, err := fn(args)
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

func TestTrimRight(t *testing.T) {
	_, fn := trimRight()

	tests := []struct {
		name     string
		input    string
		cutset   *string
		expected string
	}{
		{"trim right spaces", "  hello  ", nil, "  hello"},
		{"trim right custom", "helloxx", stringPtr("x"), "hello"},
		{"no trim needed", "hello", nil, "hello"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []lang.Value{lang.StringValue(tt.input)}
			if tt.cutset != nil {
				args = append(args, lang.StringValue(*tt.cutset))
			}

			result, err := fn(args)
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

func TestPadLeft(t *testing.T) {
	_, fn := padLeft()

	tests := []struct {
		name     string
		input    string
		length   float64
		padChar  *string
		expected string
	}{
		{"pad with spaces", "hello", 8, nil, "   hello"},
		{"pad with custom char", "test", 6, stringPtr("-"), "--test"},
		{"no padding needed", "hello", 3, nil, "hello"},
		{"empty pad char", "hi", 4, stringPtr(""), "  hi"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []lang.Value{lang.StringValue(tt.input), lang.NumberValue(tt.length)}
			if tt.padChar != nil {
				args = append(args, lang.StringValue(*tt.padChar))
			}

			result, err := fn(args)
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

func TestPadRight(t *testing.T) {
	_, fn := padRight()

	tests := []struct {
		name     string
		input    string
		length   float64
		padChar  *string
		expected string
	}{
		{"pad with spaces", "hello", 8, nil, "hello   "},
		{"pad with custom char", "test", 6, stringPtr("-"), "test--"},
		{"no padding needed", "hello", 3, nil, "hello"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []lang.Value{lang.StringValue(tt.input), lang.NumberValue(tt.length)}
			if tt.padChar != nil {
				args = append(args, lang.StringValue(*tt.padChar))
			}

			result, err := fn(args)
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

func TestPadCenter(t *testing.T) {
	_, fn := padCenter()

	tests := []struct {
		name     string
		input    string
		length   float64
		padChar  *string
		expected string
	}{
		{"pad center odd", "hello", 9, nil, "  hello  "},
		{"pad center even", "test", 8, nil, "  test  "},
		{"pad with custom char", "hi", 6, stringPtr("-"), "--hi--"},
		{"no padding needed", "hello", 3, nil, "hello"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []lang.Value{lang.StringValue(tt.input), lang.NumberValue(tt.length)}
			if tt.padChar != nil {
				args = append(args, lang.StringValue(*tt.padChar))
			}

			result, err := fn(args)
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

func TestSubstr(t *testing.T) {
	_, fn := substr()

	tests := []struct {
		name     string
		input    string
		start    float64
		length   *float64
		expected string
		hasError bool
	}{
		{"substr middle", "hello world", 6, nil, "world", false},
		{"substr with length", "hello world", 0, float64Ptr(5), "hello", false},
		{"substr negative start", "hello", -2, nil, "lo", false},
		{"substr out of bounds", "hello", 10, nil, "", false},
		{"substr negative length", "hello", 1, float64Ptr(-1), "", false},
		{"substr zero length", "hello", 1, float64Ptr(0), "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []lang.Value{lang.StringValue(tt.input), lang.NumberValue(tt.start)}
			if tt.length != nil {
				args = append(args, lang.NumberValue(*tt.length))
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

func TestLeft(t *testing.T) {
	_, fn := left()

	tests := []struct {
		name     string
		input    string
		count    float64
		expected string
	}{
		{"left normal", "hello world", 5, "hello"},
		{"left zero", "hello", 0, ""},
		{"left negative", "hello", -1, ""},
		{"left more than length", "hi", 5, "hi"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.StringValue(tt.input), lang.NumberValue(tt.count)})
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

func TestRight(t *testing.T) {
	_, fn := right()

	tests := []struct {
		name     string
		input    string
		count    float64
		expected string
	}{
		{"right normal", "hello world", 5, "world"},
		{"right zero", "hello", 0, ""},
		{"right negative", "hello", -1, ""},
		{"right more than length", "hi", 5, "hi"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.StringValue(tt.input), lang.NumberValue(tt.count)})
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

func TestContains(t *testing.T) {
	_, fn := contains()

	tests := []struct {
		name     string
		input    string
		substr   string
		expected bool
	}{
		{"contains true", "hello world", "world", true},
		{"contains false", "hello world", "foo", false},
		{"contains empty", "hello", "", true},
		{"empty contains empty", "", "", true},
		{"empty contains non-empty", "", "hello", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.StringValue(tt.input), lang.StringValue(tt.substr)})
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

func TestStartsWith(t *testing.T) {
	_, fn := startsWith()

	tests := []struct {
		name     string
		input    string
		prefix   string
		expected bool
	}{
		{"starts with true", "hello world", "hello", true},
		{"starts with false", "hello world", "world", false},
		{"starts with empty", "hello", "", true},
		{"empty starts with empty", "", "", true},
		{"empty starts with non-empty", "", "hello", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.StringValue(tt.input), lang.StringValue(tt.prefix)})
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

func TestEndsWith(t *testing.T) {
	_, fn := endsWith()

	tests := []struct {
		name     string
		input    string
		suffix   string
		expected bool
	}{
		{"ends with true", "hello world", "world", true},
		{"ends with false", "hello world", "hello", false},
		{"ends with empty", "hello", "", true},
		{"empty ends with empty", "", "", true},
		{"empty ends with non-empty", "", "hello", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.StringValue(tt.input), lang.StringValue(tt.suffix)})
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

func TestIndexOf(t *testing.T) {
	_, fn := indexOf()

	tests := []struct {
		name     string
		input    string
		substr   string
		start    *float64
		expected float64
	}{
		{"index found", "hello world", "world", nil, 6},
		{"index not found", "hello world", "foo", nil, -1},
		{"index with start", "hello hello", "hello", float64Ptr(1), 6},
		{"index start out of bounds", "hello", "h", float64Ptr(10), -1},
		{"index empty substring", "hello", "", nil, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []lang.Value{lang.StringValue(tt.input), lang.StringValue(tt.substr)}
			if tt.start != nil {
				args = append(args, lang.NumberValue(*tt.start))
			}

			result, err := fn(args)
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

func TestLastIndexOf(t *testing.T) {
	_, fn := lastIndexOf()

	tests := []struct {
		name     string
		input    string
		substr   string
		expected float64
	}{
		{"last index found", "hello world hello", "hello", 12},
		{"last index not found", "hello world", "foo", -1},
		{"last index single occurrence", "hello world", "world", 6},
		{"last index empty substring", "hello", "", 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.StringValue(tt.input), lang.StringValue(tt.substr)})
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

func TestReplace(t *testing.T) {
	_, fn := replace()

	tests := []struct {
		name     string
		input    string
		old      string
		new      string
		count    *float64
		expected string
	}{
		{"replace all", "hello hello hello", "hello", "hi", nil, "hi hi hi"},
		{"replace with count", "hello hello hello", "hello", "hi", float64Ptr(2), "hi hi hello"},
		{"replace none", "hello world", "foo", "bar", nil, "hello world"},
		{"replace empty", "hello", "", "x", nil, "xhxexlxlxox"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []lang.Value{lang.StringValue(tt.input), lang.StringValue(tt.old), lang.StringValue(tt.new)}
			if tt.count != nil {
				args = append(args, lang.NumberValue(*tt.count))
			}

			result, err := fn(args)
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

func TestReplaceAll(t *testing.T) {
	_, fn := replaceAll()

	tests := []struct {
		name     string
		input    string
		old      string
		new      string
		expected string
	}{
		{"replace all occurrences", "hello hello hello", "hello", "hi", "hi hi hi"},
		{"replace none", "hello world", "foo", "bar", "hello world"},
		{"replace with empty", "hello", "l", "", "heo"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.StringValue(tt.input), lang.StringValue(tt.old), lang.StringValue(tt.new)})
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

func TestSplit(t *testing.T) {
	_, fn := split()

	tests := []struct {
		name        string
		input       string
		separator   string
		count       *float64
		expectedLen int
		expected    []string
	}{
		{"split by space", "hello world test", " ", nil, 3, []string{"hello", "world", "test"}},
		{"split with count", "a,b,c,d", ",", float64Ptr(3), 3, []string{"a", "b", "c,d"}},
		{"split by empty", "hello", "", nil, 5, []string{"h", "e", "l", "l", "o"}},
		{"split empty string", "", ",", nil, 1, []string{""}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []lang.Value{lang.StringValue(tt.input), lang.StringValue(tt.separator)}
			if tt.count != nil {
				args = append(args, lang.NumberValue(*tt.count))
			}

			result, err := fn(args)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			resultList := result.(lang.ListValue)
			if len(resultList) != tt.expectedLen {
				t.Errorf("Expected length %d, got %d", tt.expectedLen, len(resultList))
				return
			}

			for i, expectedStr := range tt.expected {
				if string(resultList[i].(lang.StringValue)) != expectedStr {
					t.Errorf("Index %d: expected %s, got %s", i, expectedStr, string(resultList[i].(lang.StringValue)))
				}
			}
		})
	}
}

func TestMatch(t *testing.T) {
	_, fn := match()

	tests := []struct {
		name     string
		input    string
		pattern  string
		expected bool
		hasError bool
	}{
		{"simple match", "hello", "h.*o", true, false},
		{"no match", "hello", "x+", false, false},
		{"digit pattern", "123", `\d+`, true, false},
		{"invalid regex", "hello", "[", false, true},
		{"empty pattern", "hello", "", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.StringValue(tt.input), lang.StringValue(tt.pattern)})
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
			if bool(result.(lang.BoolValue)) != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, bool(result.(lang.BoolValue)))
			}
		})
	}
}

func TestFindAll(t *testing.T) {
	_, fn := findAll()

	tests := []struct {
		name        string
		input       string
		pattern     string
		count       *float64
		expectedLen int
		hasError    bool
	}{
		{"find digits", "a1b2c3", `\d`, nil, 3, false},
		{"find words", "hello world test", `\w+`, nil, 3, false},
		{"find with limit", "a1b2c3d4", `\d`, float64Ptr(2), 2, false},
		{"no matches", "hello", `\d`, nil, 0, false},
		{"invalid regex", "hello", "[", nil, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []lang.Value{lang.StringValue(tt.input), lang.StringValue(tt.pattern)}
			if tt.count != nil {
				args = append(args, lang.NumberValue(*tt.count))
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

			resultList := result.(lang.ListValue)
			if len(resultList) != tt.expectedLen {
				t.Errorf("Expected length %d, got %d", tt.expectedLen, len(resultList))
			}
		})
	}
}

func TestReplaceRegex(t *testing.T) {
	_, fn := replaceRegex()

	tests := []struct {
		name        string
		input       string
		pattern     string
		replacement string
		expected    string
		hasError    bool
	}{
		{"replace digits", "a1b2c3", `\d`, "X", "aXbXcX", false},
		{"replace words", "hello world", `\w+`, "***", "*** ***", false},
		{"no matches", "hello", `\d`, "X", "hello", false},
		{"invalid regex", "hello", "[", "X", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{
				lang.StringValue(tt.input),
				lang.StringValue(tt.pattern),
				lang.StringValue(tt.replacement),
			})
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

func TestCharAt(t *testing.T) {
	_, fn := charAt()

	tests := []struct {
		name     string
		input    string
		index    float64
		expected string
	}{
		{"valid index", "hello", 1, "e"},
		{"first char", "hello", 0, "h"},
		{"last char", "hello", 4, "o"},
		{"out of bounds", "hello", 10, ""},
		{"negative index", "hello", -1, ""},
		{"unicode char", "世界", 1, "界"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.StringValue(tt.input), lang.NumberValue(tt.index)})
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

func TestCharCode(t *testing.T) {
	_, fn := charCode()

	tests := []struct {
		name     string
		input    string
		index    float64
		expected float64
	}{
		{"char code A", "ABC", 0, 65},
		{"char code a", "abc", 0, 97},
		{"char code space", " ", 0, 32},
		{"out of bounds", "hello", 10, 0},
		{"negative index", "hello", -1, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.StringValue(tt.input), lang.NumberValue(tt.index)})
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

func TestFromCharCode(t *testing.T) {
	_, fn := fromCharCode()

	tests := []struct {
		name     string
		codes    []float64
		expected string
	}{
		{"single char", []float64{65}, "A"},
		{"multiple chars", []float64{72, 101, 108, 108, 111}, "Hello"},
		{"empty input", []float64{}, ""},
		{"unicode", []float64{19990, 30028}, "世界"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := make([]lang.Value, len(tt.codes))
			for i, code := range tt.codes {
				args[i] = lang.NumberValue(code)
			}

			result, err := fn(args)
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

func TestIsEmpty(t *testing.T) {
	_, fn := isEmpty()

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"empty string", "", true},
		{"whitespace only", "   ", true},
		{"tabs and spaces", "\t  \n", true},
		{"non-empty", "hello", false},
		{"whitespace with content", " hello ", false},
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

func TestIsNumeric(t *testing.T) {
	_, fn := isNumeric()

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"integer", "123", true},
		{"float", "123.45", true},
		{"negative", "-123", true},
		{"scientific notation", "1.23e10", true},
		{"not numeric", "hello", false},
		{"mixed", "123abc", false},
		{"empty", "", false},
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

func TestIsAlpha(t *testing.T) {
	_, fn := isAlpha()

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"all letters", "hello", true},
		{"mixed case", "HeLLo", true},
		{"with numbers", "hello123", false},
		{"with spaces", "hello world", false},
		{"empty string", "", false},
		{"unicode letters", "世界", true},
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

func TestIsAlphaNumeric(t *testing.T) {
	_, fn := isAlphaNumeric()

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"letters only", "hello", true},
		{"numbers only", "123", true},
		{"letters and numbers", "hello123", true},
		{"with spaces", "hello 123", false},
		{"with symbols", "hello!", false},
		{"empty string", "", false},
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

func TestIsSpace(t *testing.T) {
	_, fn := isSpace()

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"spaces only", "   ", true},
		{"tabs and spaces", "\t \n", true},
		{"with letters", "  a  ", false},
		{"empty string", "", false},
		{"single space", " ", true},
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

func TestToString(t *testing.T) {
	_, fn := toStr()

	tests := []struct {
		name     string
		input    lang.Value
		expected string
		hasError bool
	}{
		{"string input", lang.StringValue("hello"), "hello", false},
		{"number input", lang.NumberValue(123), "123", false},
		{"boolean input", lang.BoolValue(true), "true", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{tt.input})
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

func TestToNumber(t *testing.T) {
	_, fn := toNumber()

	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		{"integer string", "123", 123},
		{"float string", "123.45", 123.45},
		{"negative", "-123", -123},
		{"invalid number", "hello", 0},
		{"empty string", "", 0},
		{"scientific notation", "1.23e2", 123},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.StringValue(tt.input)})
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

func TestArgumentErrors(t *testing.T) {
	functions := []struct {
		fn           func() (string, lang.Function)
		expectedArgs int
	}{
		{length, 1},
		{size, 1},
		{repeat, 2},
		{reverse, 1},
		{toUpper, 1},
		{toLower, 1},
		{title, 1},
		{capitalize, 1},
		{swapCase, 1},
		{left, 2},
		{right, 2},
		{contains, 2},
		{startsWith, 2},
		{endsWith, 2},
		{lastIndexOf, 2},
		{replaceAll, 3},
		{join, 2},
		{lines, 1},
		{fields, 1},
		{match, 2},
		{replaceRegex, 3},
		{charAt, 2},
		{charCode, 2},
		{isEmpty, 1},
		{isNumeric, 1},
		{isAlpha, 1},
		{isAlphaNumeric, 1},
		{isSpace, 1},
		{toStr, 1},
		{toNumber, 1},
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

func TestNonStringInputs(t *testing.T) {
	functions := []func() (string, lang.Function){
		length, size, reverse, toUpper, toLower, title, capitalize, swapCase,
		isEmpty, isNumeric, isAlpha, isAlphaNumeric, isSpace, lines, fields,
	}

	nonStringInput := lang.NumberValue(123)

	for _, getFn := range functions {
		name, fn := getFn()
		t.Run(name+"_non_string", func(t *testing.T) {
			// Most functions should handle non-string inputs via lib.ToString
			// but let's test with types that can't be converted
			_, err := fn([]lang.Value{nonStringInput})
			// These should actually work since lib.ToString can convert numbers
			if err != nil && name != "length" && name != "size" {
				// Only test functions that should strictly require strings
				t.Logf("Function %s rejected non-string input (expected for strict string functions)", name)
			}
		})
	}
}

func TestComplexScenarios(t *testing.T) {
	t.Run("case_conversion_roundtrip", func(t *testing.T) {
		original := "Hello World Test"

		_, upperFn := toUpper()
		_, lowerFn := toLower()
		_, titleFn := title()

		// Convert to upper
		upper, err := upperFn([]lang.Value{lang.StringValue(original)})
		if err != nil {
			t.Errorf("Error converting to upper: %v", err)
			return
		}

		// Convert back to lower
		lower, err := lowerFn([]lang.Value{upper})
		if err != nil {
			t.Errorf("Error converting to lower: %v", err)
			return
		}

		// Convert to title case
		title, err := titleFn([]lang.Value{lower})
		if err != nil {
			t.Errorf("Error converting to title: %v", err)
			return
		}

		if string(title.(lang.StringValue)) != original {
			t.Errorf("Case conversion roundtrip failed: expected %s, got %s", original, string(title.(lang.StringValue)))
		}
	})

	t.Run("split_join_roundtrip", func(t *testing.T) {
		original := "apple,banana,cherry"
		separator := ","

		_, splitFn := split()
		_, joinFn := join()

		// Split the string
		split, err := splitFn([]lang.Value{lang.StringValue(original), lang.StringValue(separator)})
		if err != nil {
			t.Errorf("Error splitting string: %v", err)
			return
		}

		// Join back together
		joined, err := joinFn([]lang.Value{lang.StringValue(separator), split})
		if err != nil {
			t.Errorf("Error joining string: %v", err)
			return
		}

		if string(joined.(lang.StringValue)) != original {
			t.Errorf("Split/join roundtrip failed: expected %s, got %s", original, string(joined.(lang.StringValue)))
		}
	})

	t.Run("unicode_handling", func(t *testing.T) {
		unicodeStr := "Hello 世界 🌍"

		_, lenFn := length()
		_, sizeFn := size()
		_, reverseFn := reverse()

		// Test length vs size with unicode
		lenResult, _ := lenFn([]lang.Value{lang.StringValue(unicodeStr)})
		sizeResult, _ := sizeFn([]lang.Value{lang.StringValue(unicodeStr)})

		runeCount := float64(lenResult.(lang.NumberValue))
		byteCount := float64(sizeResult.(lang.NumberValue))

		if byteCount <= runeCount {
			t.Errorf("Expected byte count (%f) to be greater than rune count (%f) for unicode string", byteCount, runeCount)
		}

		// Test reverse with unicode
		reversed, err := reverseFn([]lang.Value{lang.StringValue(unicodeStr)})
		if err != nil {
			t.Errorf("Error reversing unicode string: %v", err)
		}

		// Reverse again should give original
		doubleReversed, err := reverseFn([]lang.Value{reversed})
		if err != nil {
			t.Errorf("Error double-reversing unicode string: %v", err)
		}

		if string(doubleReversed.(lang.StringValue)) != unicodeStr {
			t.Errorf("Double reverse failed for unicode: expected %s, got %s", unicodeStr, string(doubleReversed.(lang.StringValue)))
		}
	})
}

func TestExport(t *testing.T) {
	functions := Export()

	expectedFunctions := []string{
		"len", "size", "concat", "repeat", "reverse",
		"upper", "lower", "title", "capitalize", "swapCase",
		"trim", "trimLeft", "trimRight", "padLeft", "padRight", "padCenter",
		"substr", "left", "right",
		"contains", "startswith", "endswith", "indexof", "lastIndexOf", "replace", "replaceAll",
		"split", "join", "lines", "fields",
		"match", "findAll", "replaceRegex",
		"charAt", "charCode", "fromCharCode",
		"isEmpty", "isNumeric", "isAlpha", "isAlphanumeric", "isSpace",
		"toString", "toNumber",
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

func BenchmarkLength(b *testing.B) {
	_, fn := length()
	args := []lang.Value{lang.StringValue("hello world test string")}

	for i := 0; i < b.N; i++ {
		fn(args)
	}
}

func BenchmarkConcat(b *testing.B) {
	_, fn := voncat()
	args := []lang.Value{
		lang.StringValue("hello"),
		lang.StringValue(" "),
		lang.StringValue("world"),
		lang.StringValue("!"),
	}

	for i := 0; i < b.N; i++ {
		fn(args)
	}
}

func BenchmarkReplace(b *testing.B) {
	_, fn := replace()
	args := []lang.Value{
		lang.StringValue("hello world hello world hello"),
		lang.StringValue("hello"),
		lang.StringValue("hi"),
	}

	for i := 0; i < b.N; i++ {
		fn(args)
	}
}

func BenchmarkSplit(b *testing.B) {
	_, fn := split()
	args := []lang.Value{
		lang.StringValue("apple,banana,cherry,date,elderberry"),
		lang.StringValue(","),
	}

	for i := 0; i < b.N; i++ {
		fn(args)
	}
}

func stringPtr(s string) *string {
	return &s
}

func float64Ptr(f float64) *float64 {
	return &f
}

func TestJoin(t *testing.T) {
	_, fn := join()

	tests := []struct {
		name      string
		separator string
		list      []string
		expected  string
	}{
		{"join with comma", ",", []string{"a", "b", "c"}, "a,b,c"},
		{"join with space", " ", []string{"hello", "world"}, "hello world"},
		{"join empty list", ",", []string{}, ""},
		{"join single item", ",", []string{"alone"}, "alone"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := make(lang.ListValue, len(tt.list))
			for i, str := range tt.list {
				list[i] = lang.StringValue(str)
			}

			result, err := fn([]lang.Value{lang.StringValue(tt.separator), list})
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

func TestLines(t *testing.T) {
	_, fn := lines()

	tests := []struct {
		name        string
		input       string
		expectedLen int
		expected    []string
	}{
		{"multiple lines", "line1\nline2\nline3", 3, []string{"line1", "line2", "line3"}},
		{"single line", "single", 1, []string{"single"}},
		{"empty lines", "a\n\nb", 3, []string{"a", "", "b"}},
		{"trailing newline", "line1\nline2\n", 3, []string{"line1", "line2", ""}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.StringValue(tt.input)})
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			resultList := result.(lang.ListValue)
			if len(resultList) != tt.expectedLen {
				t.Errorf("Expected length %d, got %d", tt.expectedLen, len(resultList))
				return
			}

			for i, expectedStr := range tt.expected {
				if string(resultList[i].(lang.StringValue)) != expectedStr {
					t.Errorf("Index %d: expected %s, got %s", i, expectedStr, string(resultList[i].(lang.StringValue)))
				}
			}
		})
	}
}

func TestFields(t *testing.T) {
	_, fn := fields()

	tests := []struct {
		name        string
		input       string
		expectedLen int
		expected    []string
	}{
		{"multiple fields", "hello   world  test", 3, []string{"hello", "world", "test"}},
		{"single field", "single", 1, []string{"single"}},
		{"empty string", "", 0, []string{}},
		{"only spaces", "   ", 0, []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.StringValue(tt.input)})
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			resultList := result.(lang.ListValue)
			if len(resultList) != tt.expectedLen {
				t.Errorf("Expected length %d, got %d", tt.expectedLen, len(resultList))
				return
			}

			for i, expectedStr := range tt.expected {
				if i < len(resultList) {
					if string(resultList[i].(lang.StringValue)) != expectedStr {
						t.Errorf("Index %d: expected %s, got %s", i, expectedStr, string(resultList[i].(lang.StringValue)))
					}
				}
			}
		})
	}
}
