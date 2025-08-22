package lang

import (
	"fmt"
	"strconv"
)

type yyLex struct {
	input  string
	pos    int
	result ExprNode
	error  error
}

func (l *yyLex) Lex(lval *yySymType) int {
	// Skip whitespace
	for l.pos < len(l.input) && (isWhitespace(l.input[l.pos])) {
		l.pos++
	}

	if l.pos >= len(l.input) {
		return 0 // EOF
	}

	// Check for keywords and operators
	if matched, newPos := l.matchKeyword("and"); matched {
		l.pos = newPos
		return AND
	}
	if matched, newPos := l.matchKeyword("or"); matched {
		l.pos = newPos
		return OR
	}
	if matched, newPos := l.matchKeyword("not"); matched {
		l.pos = newPos
		return NOT
	}
	if matched, newPos := l.matchKeyword("in"); matched {
		l.pos = newPos
		return IN
	}
	if matched, newPos := l.matchKeyword("true"); matched {
		l.pos = newPos
		lval.boolean = true
		return BOOLEAN
	}
	if matched, newPos := l.matchKeyword("false"); matched {
		l.pos = newPos
		lval.boolean = false
		return BOOLEAN
	}

	// Two-character operators
	if l.pos+1 < len(l.input) {
		twoChar := l.input[l.pos : l.pos+2]
		switch twoChar {
		case "==":
			l.pos += 2
			return EQ
		case "!=":
			l.pos += 2
			return NE
		case "<=":
			l.pos += 2
			return LE
		case ">=":
			l.pos += 2
			return GE
		}
	}

	ch := l.input[l.pos]

	switch ch {
	case '(':
		l.pos++
		return LPAREN
	case ')':
		l.pos++
		return RPAREN
	case '[':
		l.pos++
		return LBRACKET
	case ']':
		l.pos++
		return RBRACKET
	case '.':
		l.pos++
		return DOT
	case ',':
		l.pos++
		return COMMA
	case '+':
		l.pos++
		return int(ch)
	case '-':
		l.pos++
		return int(ch)
	case '*':
		l.pos++
		return int(ch)
	case '/':
		l.pos++
		return int(ch)
	case '=':
		l.pos++
		return EQ
	case '<':
		l.pos++
		return LT
	case '>':
		l.pos++
		return GT
	case '\'':
		if token, newPos := l.readSingleQuotedString(lval); token != 0 {
			l.pos = newPos
			return token
		}
	case '"':
		if token, newPos := l.readDoubleQuotedString(lval); token != 0 {
			l.pos = newPos
			return token
		}
	default:
		if ch >= '0' && ch <= '9' {
			if token, newPos := l.readNumber(lval); token != 0 {
				l.pos = newPos
				return token
			}
		}
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_' {
			if token, newPos := l.readIdentifier(lval); token != 0 {
				l.pos = newPos
				return token
			}
		}
	}
	return 0
}

func (l *yyLex) matchKeyword(keyword string) (bool, int) {
	if l.pos+len(keyword) > len(l.input) {
		return false, l.pos
	}

	if l.input[l.pos:l.pos+len(keyword)] == keyword {
		// Check that it's not part of a larger identifier
		nextPos := l.pos + len(keyword)
		if nextPos < len(l.input) {
			nextChar := l.input[nextPos]
			if (nextChar >= 'a' && nextChar <= 'z') ||
				(nextChar >= 'A' && nextChar <= 'Z') ||
				(nextChar >= '0' && nextChar <= '9') ||
				nextChar == '_' {
				return false, l.pos
			}
		}
		return true, l.pos + len(keyword)
	}
	return false, l.pos
}

func (l *yyLex) readSingleQuotedString(lval *yySymType) (int, int) {
	start := l.pos + 1 // Skip opening quote
	pos := start
	for pos < len(l.input) && l.input[pos] != '\'' {
		pos++
	}
	if pos < len(l.input) {
		pos++ // Skip closing quote
		lval.str = l.input[start : pos-1]
		return STRING, pos
	}
	return 0, l.pos // Error - unclosed string
}

func (l *yyLex) readDoubleQuotedString(lval *yySymType) (int, int) {
	start := l.pos + 1 // Skip opening quote
	pos := start
	for pos < len(l.input) && l.input[pos] != '"' {
		pos++
	}
	if pos < len(l.input) {
		pos++ // Skip closing quote
		lval.str = l.input[start : pos-1]
		return DSTRING, pos
	}
	return 0, l.pos // Error - unclosed string
}

func (l *yyLex) readNumber(lval *yySymType) (int, int) {
	start := l.pos
	pos := l.pos
	for pos < len(l.input) && ((l.input[pos] >= '0' && l.input[pos] <= '9') || l.input[pos] == '.') {
		pos++
	}
	if pos > start {
		num, _ := strconv.ParseFloat(l.input[start:pos], 64)
		lval.num = num
		return NUMBER, pos
	}
	return 0, l.pos
}

func (l *yyLex) readIdentifier(lval *yySymType) (int, int) {
	start := l.pos
	pos := l.pos
	for pos < len(l.input) {
		ch := l.input[pos]
		if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '_') {
			break
		}
		pos++
	}
	if pos > start {
		lval.str = l.input[start:pos]
		return IDENTIFIER, pos
	}
	return 0, l.pos
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func (l *yyLex) Error(s string) {
	// Find the current token or problematic area
	start := l.pos
	end := l.pos

	// If we're at EOF, back up to show last token
	if l.pos >= len(l.input) {
		start = len(l.input) - 1
		if start < 0 {
			start = 0
		}
	}

	// Extend backwards to find token start
	for start > 0 && !isWhitespace(l.input[start-1]) {
		start--
	}

	// Extend forwards to find token end
	for end < len(l.input) && !isWhitespace(l.input[end]) {
		end++
	}

	// Extract the problematic token
	var token string
	if start < len(l.input) && end <= len(l.input) && start < end {
		token = l.input[start:end]
	} else {
		token = "<EOF>"
	}

	// Create context around the error
	contextStart := 0
	if l.pos > 20 {
		contextStart = l.pos - 20
	}
	contextEnd := len(l.input)
	if l.pos+20 < len(l.input) {
		contextEnd = l.pos + 20
	}

	context := l.input[contextStart:contextEnd]

	// Calculate relative position in context for pointer
	relativePos := l.pos - contextStart
	if relativePos < 0 {
		relativePos = 0
	}
	if relativePos > len(context) {
		relativePos = len(context)
	}

	// Create pointer string
	pointer := ""
	for i := 0; i < relativePos-len(token); i++ {
		if context[i] == '\t' {
			pointer += "\t"
		} else {
			pointer += " "
		}
	}
	pointer += "^"

	l.error = fmt.Errorf("%s near token '%s' at position %d\n%s\n%s",
		s, token, l.pos-len(token), context, pointer)
}

func ParseExpression(input string) (ExprNode, error) {
	yyErrorVerbose = true
	lexer := &yyLex{input: input}
	yyParse(lexer)
	if lexer.error != nil {
		return nil, lexer.error
	}
	return lexer.result, nil
}
