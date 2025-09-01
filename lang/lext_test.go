package lang

import (
	"strings"
	"testing"
)

const (
	EOF = 0
)

func TestKeywordTokens(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"and keyword", "and", AND},
		{"or keyword", "or", OR},
		{"not keyword", "not", NOT},
		{"in keyword", "in", IN},
		{"true keyword", "true", BOOLEAN},
		{"false keyword", "false", BOOLEAN},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := &yyLex{input: tt.input}
			var lval yySymType
			token := lexer.Lex(&lval)

			if token != tt.expected {
				t.Errorf("expected token %d, got %d", tt.expected, token)
			}

			// Test boolean values specifically
			if tt.input == "true" && !lval.boolean {
				t.Error("expected true boolean value")
			}
			if tt.input == "false" && lval.boolean {
				t.Error("expected false boolean value")
			}
		})
	}
}

func TestOperatorTokens(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"equality", "==", EQ},
		{"inequality", "!=", NE},
		{"less than or equal", "<=", LE},
		{"greater than or equal", ">=", GE},
		{"single equals", "=", EQ},
		{"less than", "<", LT},
		{"greater than", ">", GT},
		{"plus", "+", int('+')},
		{"minus", "-", int('-')},
		{"multiply", "*", int('*')},
		{"divide", "/", int('/')},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := &yyLex{input: tt.input}
			var lval yySymType
			token := lexer.Lex(&lval)

			if token != tt.expected {
				t.Errorf("expected token %d, got %d", tt.expected, token)
			}
		})
	}
}

func TestPunctuationTokens(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"left paren", "(", LPAREN},
		{"right paren", ")", RPAREN},
		{"left bracket", "[", LBRACKET},
		{"right bracket", "]", RBRACKET},
		{"dot", ".", DOT},
		{"comma", ",", COMMA},
		{"question mark", "?", QMARK},
		{"colon", ":", COLON},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := &yyLex{input: tt.input}
			var lval yySymType
			token := lexer.Lex(&lval)

			if token != tt.expected {
				t.Errorf("expected token %d, got %d", tt.expected, token)
			}
		})
	}
}

func TestStringTokens(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
		value    string
	}{
		{"single quoted string", "'hello'", STRING, "hello"},
		{"double quoted string", `"world"`, DSTRING, "world"},
		{"empty single quoted", "''", STRING, ""},
		{"empty double quoted", `""`, DSTRING, ""},
		{"string with spaces", "'hello world'", STRING, "hello world"},
		{"string with escape", `'hello\'s'`, STRING, "hello\\'s"},
		{"double string with escape", `"say \"hello\""`, DSTRING, `say \"hello\"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := &yyLex{input: tt.input}
			var lval yySymType
			token := lexer.Lex(&lval)

			if token != tt.expected {
				t.Errorf("expected token %d, got %d", tt.expected, token)
			}
			if lval.str != tt.value {
				t.Errorf("expected string value %q, got %q", tt.value, lval.str)
			}
		})
	}
}

func TestNumberTokens(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		{"integer", "42", 42},
		{"zero", "0", 0},
		{"decimal", "3.14", 3.14},
		{"leading zero decimal", "0.5", 0.5},
		{"large number", "12345", 12345},
		{"decimal with trailing zeros", "1.00", 1.00},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := &yyLex{input: tt.input}
			var lval yySymType
			token := lexer.Lex(&lval)

			if token != NUMBER {
				t.Errorf("expected NUMBER token, got %d", token)
			}
			if lval.num != tt.expected {
				t.Errorf("expected number value %f, got %f", tt.expected, lval.num)
			}
		})
	}
}

func TestIdentifierTokens(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple identifier", "hello", "hello"},
		{"with underscore", "hello_world", "hello_world"},
		{"with numbers", "var123", "var123"},
		{"starts with underscore", "_private", "_private"},
		{"mixed case", "myVar", "myVar"},
		{"single char", "x", "x"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := &yyLex{input: tt.input}
			var lval yySymType
			token := lexer.Lex(&lval)

			if token != IDENTIFIER {
				t.Errorf("expected IDENTIFIER token, got %d", token)
			}
			if lval.str != tt.expected {
				t.Errorf("expected identifier %q, got %q", tt.expected, lval.str)
			}
		})
	}
}

func TestWhitespaceHandling(t *testing.T) {
	tests := []struct {
		name  string
		input string
		token int
	}{
		{"spaces before identifier", "   hello", IDENTIFIER},
		{"tabs before number", "\t\t42", NUMBER},
		{"newlines before string", "\n\n'test'", STRING},
		{"mixed whitespace", " \t\n\r and", AND},
		{"only whitespace", "   ", EOF},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := &yyLex{input: tt.input}
			var lval yySymType
			token := lexer.Lex(&lval)

			if token != tt.token {
				t.Errorf("expected token %d, got %d", tt.token, token)
			}
		})
	}
}

func TestKeywordBoundaries(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"and as part of identifier", "android", IDENTIFIER},
		{"or as part of identifier", "order", IDENTIFIER},
		{"not as part of identifier", "notation", IDENTIFIER},
		{"in as part of identifier", "input", IDENTIFIER},
		{"true as part of identifier", "truthy", IDENTIFIER},
		{"false as part of identifier", "falsely", IDENTIFIER},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := &yyLex{input: tt.input}
			var lval yySymType
			token := lexer.Lex(&lval)

			if token != tt.expected {
				t.Errorf("expected token %d, got %d", tt.expected, token)
			}
		})
	}
}

func TestMultipleTokens(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []int
	}{
		{"simple expression", "x + 1", []int{IDENTIFIER, int('+'), NUMBER, EOF}},
		{"boolean expression", "true and false", []int{BOOLEAN, AND, BOOLEAN, EOF}},
		{"comparison", "x == 42", []int{IDENTIFIER, EQ, NUMBER, EOF}},
		{"function call", "func(arg)", []int{IDENTIFIER, LPAREN, IDENTIFIER, RPAREN, EOF}},
		{"array access", "arr[0]", []int{IDENTIFIER, LBRACKET, NUMBER, RBRACKET, EOF}},
		{"ternary operator", "x ? y : z", []int{IDENTIFIER, QMARK, IDENTIFIER, COLON, IDENTIFIER, EOF}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := &yyLex{input: tt.input}
			var lval yySymType

			for i, expectedToken := range tt.expected {
				token := lexer.Lex(&lval)
				if token != expectedToken {
					t.Errorf("token %d: expected %d, got %d", i, expectedToken, token)
					break
				}
			}
		})
	}
}

func TestMatchKeyword(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		keyword  string
		position int
		expected bool
		newPos   int
	}{
		{"exact match", "and", "and", 0, true, 3},
		{"keyword at start", "and something", "and", 0, true, 3},
		{"keyword not at boundary", "android", "and", 0, false, 0},
		{"keyword at end", "test and", "and", 5, true, 8},
		{"no match", "hello", "and", 0, false, 0},
		{"partial match", "an", "and", 0, false, 0},
		{"beyond input", "test", "testing", 0, false, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := &yyLex{input: tt.input, pos: tt.position}
			matched, newPos := lexer.matchKeyword(tt.keyword)

			if matched != tt.expected {
				t.Errorf("expected match %v, got %v", tt.expected, matched)
			}
			if newPos != tt.newPos {
				t.Errorf("expected new position %d, got %d", tt.newPos, newPos)
			}
		})
	}
}

func TestReadSingleQuotedString(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		position  int
		expectPos int
		expectVal string
		expectTok int
	}{
		{"simple string", "'hello'", 0, 7, "hello", STRING},
		{"empty string", "''", 0, 2, "", STRING},
		{"string with space", "'hello world'", 0, 13, "hello world", STRING},
		{"unclosed string", "'hello", 0, 0, "", 0},
		{"escaped quote", "'don\\'t'", 0, 8, "don\\'t", STRING},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := &yyLex{input: tt.input, pos: tt.position}
			var lval yySymType
			token, newPos := lexer.readSingleQuotedString(&lval)

			if token != tt.expectTok {
				t.Errorf("expected token %d, got %d", tt.expectTok, token)
			}
			if newPos != tt.expectPos {
				t.Errorf("expected position %d, got %d", tt.expectPos, newPos)
			}
			if token != 0 && lval.str != tt.expectVal {
				t.Errorf("expected value %q, got %q", tt.expectVal, lval.str)
			}
		})
	}
}

func TestReadDoubleQuotedString(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		position  int
		expectPos int
		expectVal string
		expectTok int
	}{
		{"simple string", `"hello"`, 0, 7, "hello", DSTRING},
		{"empty string", `""`, 0, 2, "", DSTRING},
		{"string with space", `"hello world"`, 0, 13, "hello world", DSTRING},
		{"unclosed string", `"hello`, 0, 0, "", 0},
		// {"escaped quote", `"say \"hi\""`, 0, 11, `say \"hi\"`, DSTRING},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := &yyLex{input: tt.input, pos: tt.position}
			var lval yySymType
			token, newPos := lexer.readDoubleQuotedString(&lval)

			if token != tt.expectTok {
				t.Errorf("expected token %d, got %d", tt.expectTok, token)
			}
			if newPos != tt.expectPos {
				t.Errorf("expected position %d, got %d", tt.expectPos, newPos)
			}
			if token != 0 && lval.str != tt.expectVal {
				t.Errorf("expected value %q, got %q", tt.expectVal, lval.str)
			}
		})
	}
}

func TestReadNumber(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		position  int
		expectPos int
		expectVal float64
		expectTok int
	}{
		{"integer", "42", 0, 2, 42, NUMBER},
		{"decimal", "3.14", 0, 4, 3.14, NUMBER},
		{"zero", "0", 0, 1, 0, NUMBER},
		{"leading zero", "0.5", 0, 3, 0.5, NUMBER},
		{"number at position", "abc123", 3, 6, 123, NUMBER},
		{"not a number", "abc", 0, 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := &yyLex{input: tt.input, pos: tt.position}
			var lval yySymType
			token, newPos := lexer.readNumber(&lval)

			if token != tt.expectTok {
				t.Errorf("expected token %d, got %d", tt.expectTok, token)
			}
			if newPos != tt.expectPos {
				t.Errorf("expected position %d, got %d", tt.expectPos, newPos)
			}
			if token != 0 && lval.num != tt.expectVal {
				t.Errorf("expected value %f, got %f", tt.expectVal, lval.num)
			}
		})
	}
}

func TestReadIdentifier(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		position  int
		expectPos int
		expectVal string
		expectTok int
	}{
		{"simple identifier", "hello", 0, 5, "hello", IDENTIFIER},
		{"with underscore", "hello_world", 0, 11, "hello_world", IDENTIFIER},
		{"with numbers", "var123", 0, 6, "var123", IDENTIFIER},
		{"starts with underscore", "_private", 0, 8, "_private", IDENTIFIER},
		{"identifier at position", "123abc", 3, 6, "abc", IDENTIFIER},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := &yyLex{input: tt.input, pos: tt.position}
			var lval yySymType
			token, newPos := lexer.readIdentifier(&lval)

			if token != tt.expectTok {
				t.Errorf("expected token %d, got %d", tt.expectTok, token)
			}
			if newPos != tt.expectPos {
				t.Errorf("expected position %d, got %d", tt.expectPos, newPos)
			}
			if token != 0 && lval.str != tt.expectVal {
				t.Errorf("expected value %q, got %q", tt.expectVal, lval.str)
			}
		})
	}
}

func TestIsWhitespace(t *testing.T) {
	tests := []struct {
		name     string
		char     byte
		expected bool
	}{
		{"space", ' ', true},
		{"tab", '\t', true},
		{"newline", '\n', true},
		{"carriage return", '\r', true},
		{"letter", 'a', false},
		{"digit", '5', false},
		{"punctuation", '.', false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isWhitespace(tt.char)
			if result != tt.expected {
				t.Errorf("expected %v for char %q, got %v", tt.expected, tt.char, result)
			}
		})
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		position int
		message  string
	}{
		{"error at beginning", "invalid", 0, "syntax error"},
		{"error in middle", "hello invalid world", 6, "syntax error"},
		{"error at end", "hello world", 11, "syntax error"},
		{"empty input", "", 0, "syntax error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := &yyLex{input: tt.input, pos: tt.position}
			lexer.Error(tt.message)

			if lexer.error == nil {
				t.Error("expected error to be set")
			}

			errorStr := lexer.error.Error()
			if !strings.Contains(errorStr, tt.message) {
				t.Errorf("expected error to contain %q, got %q", tt.message, errorStr)
			}

			// Check that error includes position information
			if !strings.Contains(errorStr, "position") {
				t.Errorf("expected error to include position information, got %q", errorStr)
			}
		})
	}
}

func TestErrorContext(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		position int
		message  string
	}{
		{"short input", "bad", 2, "error"},
		{"long input", "this is a very long input with an error in the middle somewhere", 30, "error"},
		{"error at start", "error at beginning", 0, "error"},
		{"error at end", "input with error", 16, "error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := &yyLex{input: tt.input, pos: tt.position}
			lexer.Error(tt.message)

			if lexer.error == nil {
				t.Fatal("expected error to be set")
			}

			errorStr := lexer.error.Error()

			// Should contain the error message
			if !strings.Contains(errorStr, tt.message) {
				t.Errorf("expected error to contain message %q", tt.message)
			}

			// Should contain position
			if !strings.Contains(errorStr, "position") {
				t.Errorf("expected error to contain position info")
			}

			// Should contain context (part of input)
			lines := strings.Split(errorStr, "\n")
			if len(lines) < 2 {
				t.Errorf("expected multi-line error with context")
			}
		})
	}
}

func TestUnknownCharacter(t *testing.T) {
	tests := []struct {
		name  string
		input string
		char  byte
	}{
		{"hash", "#", '#'},
		{"at symbol", "@", '@'},
		{"dollar", "$", '$'},
		{"percent", "%", '%'},
		{"ampersand", "&", '&'},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := &yyLex{input: tt.input}
			var lval yySymType
			token := lexer.Lex(&lval)

			// Unknown characters should return 0 (EOF/error)
			if token != 0 {
				t.Errorf("expected token 0 for unknown character, got %d", token)
			}
		})
	}
}

// Benchmark Tests
func BenchmarkLexKeyword(b *testing.B) {
	lexer := &yyLex{input: "and or not"}
	var lval yySymType

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lexer.pos = 0
		lexer.Lex(&lval)
		lexer.Lex(&lval)
		lexer.Lex(&lval)
	}
}

func BenchmarkLexNumber(b *testing.B) {
	lexer := &yyLex{input: "123.456"}
	var lval yySymType

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lexer.pos = 0
		lexer.Lex(&lval)
	}
}

func BenchmarkLexString(b *testing.B) {
	lexer := &yyLex{input: `"hello world"`}
	var lval yySymType

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lexer.pos = 0
		lexer.Lex(&lval)
	}
}

func BenchmarkLexIdentifier(b *testing.B) {
	lexer := &yyLex{input: "identifier"}
	var lval yySymType

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lexer.pos = 0
		lexer.Lex(&lval)
	}
}

func BenchmarkLexComplexExpression(b *testing.B) {
	lexer := &yyLex{input: "func(arg1, 'string', 42) and not condition"}
	var lval yySymType

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lexer.pos = 0
		for {
			token := lexer.Lex(&lval)
			if token == 0 {
				break
			}
		}
	}
}
