package string

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/vedadiyan/exql/lang"
	"github.com/vedadiyan/exql/lib"
)

// String Functions
// These functions provide comprehensive string manipulation and processing capabilities

// Basic String Operations
func stringLen(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("len: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("len: %w", err)
	}
	return lang.NumberValue(float64(utf8.RuneCountInString(string(str)))), nil
}

func stringSize(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("size: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("size: %w", err)
	}
	return lang.NumberValue(float64(len(string(str)))), nil // Byte count
}

func stringConcat(args []lang.Value) (lang.Value, error) {
	result := ""
	for i, arg := range args {
		str, err := lib.ToString(arg)
		if err != nil {
			return nil, fmt.Errorf("concat: argument %d %w", i, err)
		}
		result += string(str)
	}
	return lang.StringValue(result), nil
}

func stringRepeat(args []lang.Value) (lang.Value, error) {
	if len(args) != 2 {
		return nil, errors.New("repeat: expected 2 arguments")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("repeat: string %w", err)
	}
	countVal, err := lib.ToNumber(args[1])
	if err != nil {
		return nil, fmt.Errorf("repeat: count %w", err)
	}
	count := int(countVal)
	if count < 0 {
		count = 0
	}
	return lang.StringValue(strings.Repeat(string(str), count)), nil
}

func stringReverse(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("reverse: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("reverse: %w", err)
	}
	runes := []rune(string(str))
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return lang.StringValue(string(runes)), nil
}

// Case Conversion
func stringUpper(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("upper: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("upper: %w", err)
	}
	return lang.StringValue(strings.ToUpper(string(str))), nil
}

func stringLower(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("lower: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("lower: %w", err)
	}
	return lang.StringValue(strings.ToLower(string(str))), nil
}

func stringTitle(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("title: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("title: %w", err)
	}
	return lang.StringValue(strings.Title(string(str))), nil
}

func stringCapitalize(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("capitalize: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("capitalize: %w", err)
	}
	s := string(str)
	if len(s) == 0 {
		return lang.StringValue(""), nil
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	for i := 1; i < len(runes); i++ {
		runes[i] = unicode.ToLower(runes[i])
	}
	return lang.StringValue(string(runes)), nil
}

func stringSwapCase(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("swap_case: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("swap_case: %w", err)
	}
	runes := []rune(string(str))
	for i, r := range runes {
		if unicode.IsUpper(r) {
			runes[i] = unicode.ToLower(r)
		} else if unicode.IsLower(r) {
			runes[i] = unicode.ToUpper(r)
		}
	}
	return lang.StringValue(string(runes)), nil
}

// Trimming and Padding
func stringTrim(args []lang.Value) (lang.Value, error) {
	if len(args) < 1 || len(args) > 2 {
		return nil, errors.New("trim: expected 1 or 2 arguments")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("trim: string %w", err)
	}
	if len(args) == 1 {
		return lang.StringValue(strings.TrimSpace(string(str))), nil
	}
	cutset, err := lib.ToString(args[1])
	if err != nil {
		return nil, fmt.Errorf("trim: cutset %w", err)
	}
	return lang.StringValue(strings.Trim(string(str), string(cutset))), nil
}

func stringTrimLeft(args []lang.Value) (lang.Value, error) {
	if len(args) < 1 || len(args) > 2 {
		return nil, errors.New("trim_left: expected 1 or 2 arguments")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("trim_left: string %w", err)
	}
	if len(args) == 1 {
		return lang.StringValue(strings.TrimLeftFunc(string(str), unicode.IsSpace)), nil
	}
	cutset, err := lib.ToString(args[1])
	if err != nil {
		return nil, fmt.Errorf("trim_left: cutset %w", err)
	}
	return lang.StringValue(strings.TrimLeft(string(str), string(cutset))), nil
}

func stringTrimRight(args []lang.Value) (lang.Value, error) {
	if len(args) < 1 || len(args) > 2 {
		return nil, errors.New("trim_right: expected 1 or 2 arguments")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("trim_right: string %w", err)
	}
	if len(args) == 1 {
		return lang.StringValue(strings.TrimRightFunc(string(str), unicode.IsSpace)), nil
	}
	cutset, err := lib.ToString(args[1])
	if err != nil {
		return nil, fmt.Errorf("trim_right: cutset %w", err)
	}
	return lang.StringValue(strings.TrimRight(string(str), string(cutset))), nil
}

func stringPadLeft(args []lang.Value) (lang.Value, error) {
	if len(args) < 2 || len(args) > 3 {
		return nil, errors.New("pad_left: expected 2 or 3 arguments")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("pad_left: string %w", err)
	}
	totalLenVal, err := lib.ToNumber(args[1])
	if err != nil {
		return nil, fmt.Errorf("pad_left: total length %w", err)
	}
	totalLen := int(totalLenVal)
	padChar := " "
	if len(args) == 3 {
		padCharVal, err := lib.ToString(args[2])
		if err != nil {
			return nil, fmt.Errorf("pad_left: pad character %w", err)
		}
		padChar = string(padCharVal)
		if len(padChar) == 0 {
			padChar = " "
		}
	}

	s := string(str)
	currentLen := utf8.RuneCountInString(s)
	if totalLen <= currentLen {
		return lang.StringValue(s), nil
	}

	padLen := totalLen - currentLen
	padding := strings.Repeat(padChar, padLen)
	return lang.StringValue(padding + s), nil
}

func stringPadRight(args []lang.Value) (lang.Value, error) {
	if len(args) < 2 || len(args) > 3 {
		return nil, errors.New("pad_right: expected 2 or 3 arguments")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("pad_right: string %w", err)
	}
	totalLenVal, err := lib.ToNumber(args[1])
	if err != nil {
		return nil, fmt.Errorf("pad_right: total length %w", err)
	}
	totalLen := int(totalLenVal)
	padChar := " "
	if len(args) == 3 {
		padCharVal, err := lib.ToString(args[2])
		if err != nil {
			return nil, fmt.Errorf("pad_right: pad character %w", err)
		}
		padChar = string(padCharVal)
		if len(padChar) == 0 {
			padChar = " "
		}
	}

	s := string(str)
	currentLen := utf8.RuneCountInString(s)
	if totalLen <= currentLen {
		return lang.StringValue(s), nil
	}

	padLen := totalLen - currentLen
	padding := strings.Repeat(padChar, padLen)
	return lang.StringValue(s + padding), nil
}

func stringPadCenter(args []lang.Value) (lang.Value, error) {
	if len(args) < 2 || len(args) > 3 {
		return nil, errors.New("pad_center: expected 2 or 3 arguments")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("pad_center: string %w", err)
	}
	totalLenVal, err := lib.ToNumber(args[1])
	if err != nil {
		return nil, fmt.Errorf("pad_center: total length %w", err)
	}
	totalLen := int(totalLenVal)
	padChar := " "
	if len(args) == 3 {
		padCharVal, err := lib.ToString(args[2])
		if err != nil {
			return nil, fmt.Errorf("pad_center: pad character %w", err)
		}
		padChar = string(padCharVal)
		if len(padChar) == 0 {
			padChar = " "
		}
	}

	s := string(str)
	currentLen := utf8.RuneCountInString(s)
	if totalLen <= currentLen {
		return lang.StringValue(s), nil
	}

	padLen := totalLen - currentLen
	leftPad := padLen / 2
	rightPad := padLen - leftPad

	leftPadding := strings.Repeat(padChar, leftPad)
	rightPadding := strings.Repeat(padChar, rightPad)

	return lang.StringValue(leftPadding + s + rightPadding), nil
}

// Substring Operations
func stringSubstr(args []lang.Value) (lang.Value, error) {
	if len(args) < 2 || len(args) > 3 {
		return nil, errors.New("substr: expected 2 or 3 arguments")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("substr: string %w", err)
	}
	startVal, err := lib.ToNumber(args[1])
	if err != nil {
		return nil, fmt.Errorf("substr: start %w", err)
	}
	start := int(startVal)

	runes := []rune(string(str))
	strLen := len(runes)

	// Handle negative start index
	if start < 0 {
		start = strLen + start
	}
	if start < 0 || start >= strLen {
		return lang.StringValue(""), nil
	}

	if len(args) == 2 {
		return lang.StringValue(string(runes[start:])), nil
	}

	lengthVal, err := lib.ToNumber(args[2])
	if err != nil {
		return nil, fmt.Errorf("substr: length %w", err)
	}
	length := int(lengthVal)
	if length < 0 {
		return lang.StringValue(""), nil
	}

	end := start + length
	if end > strLen {
		end = strLen
	}

	return lang.StringValue(string(runes[start:end])), nil
}

func stringLeft(args []lang.Value) (lang.Value, error) {
	if len(args) != 2 {
		return nil, errors.New("left: expected 2 arguments")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("left: string %w", err)
	}
	countVal, err := lib.ToNumber(args[1])
	if err != nil {
		return nil, fmt.Errorf("left: count %w", err)
	}
	count := int(countVal)

	if count <= 0 {
		return lang.StringValue(""), nil
	}

	runes := []rune(string(str))
	if count >= len(runes) {
		return lang.StringValue(string(str)), nil
	}

	return lang.StringValue(string(runes[:count])), nil
}

func stringRight(args []lang.Value) (lang.Value, error) {
	if len(args) != 2 {
		return nil, errors.New("right: expected 2 arguments")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("right: string %w", err)
	}
	countVal, err := lib.ToNumber(args[1])
	if err != nil {
		return nil, fmt.Errorf("right: count %w", err)
	}
	count := int(countVal)

	if count <= 0 {
		return lang.StringValue(""), nil
	}

	runes := []rune(string(str))
	if count >= len(runes) {
		return lang.StringValue(string(str)), nil
	}

	return lang.StringValue(string(runes[len(runes)-count:])), nil
}

// Search and Replace
func stringContains(args []lang.Value) (lang.Value, error) {
	if len(args) != 2 {
		return nil, errors.New("contains: expected 2 arguments")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("contains: string %w", err)
	}
	substr, err := lib.ToString(args[1])
	if err != nil {
		return nil, fmt.Errorf("contains: substring %w", err)
	}
	return lang.BoolValue(strings.Contains(string(str), string(substr))), nil
}

func stringStartsWith(args []lang.Value) (lang.Value, error) {
	if len(args) != 2 {
		return nil, errors.New("startswith: expected 2 arguments")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("startswith: string %w", err)
	}
	prefix, err := lib.ToString(args[1])
	if err != nil {
		return nil, fmt.Errorf("startswith: prefix %w", err)
	}
	return lang.BoolValue(strings.HasPrefix(string(str), string(prefix))), nil
}

func stringEndsWith(args []lang.Value) (lang.Value, error) {
	if len(args) != 2 {
		return nil, errors.New("endswith: expected 2 arguments")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("endswith: string %w", err)
	}
	suffix, err := lib.ToString(args[1])
	if err != nil {
		return nil, fmt.Errorf("endswith: suffix %w", err)
	}
	return lang.BoolValue(strings.HasSuffix(string(str), string(suffix))), nil
}

func stringIndexOf(args []lang.Value) (lang.Value, error) {
	if len(args) < 2 || len(args) > 3 {
		return nil, errors.New("indexof: expected 2 or 3 arguments")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("indexof: string %w", err)
	}
	substr, err := lib.ToString(args[1])
	if err != nil {
		return nil, fmt.Errorf("indexof: substring %w", err)
	}
	start := 0

	if len(args) == 3 {
		startVal, err := lib.ToNumber(args[2])
		if err != nil {
			return nil, fmt.Errorf("indexof: start %w", err)
		}
		start = int(startVal)
		if start < 0 {
			start = 0
		}
	}

	runes := []rune(string(str))
	if start >= len(runes) {
		return lang.NumberValue(-1), nil
	}

	searchStr := string(runes[start:])
	index := strings.Index(searchStr, string(substr))
	if index == -1 {
		return lang.NumberValue(-1), nil
	}

	return lang.NumberValue(float64(start + utf8.RuneCountInString(searchStr[:index]))), nil
}

func stringLastIndexOf(args []lang.Value) (lang.Value, error) {
	if len(args) != 2 {
		return nil, errors.New("lastindexof: expected 2 arguments")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("lastindexof: string %w", err)
	}
	substr, err := lib.ToString(args[1])
	if err != nil {
		return nil, fmt.Errorf("lastindexof: substring %w", err)
	}

	s := string(str)
	sub := string(substr)
	index := strings.LastIndex(s, sub)
	if index == -1 {
		return lang.NumberValue(-1), nil
	}

	return lang.NumberValue(float64(utf8.RuneCountInString(s[:index]))), nil
}

func stringReplace(args []lang.Value) (lang.Value, error) {
	if len(args) < 3 || len(args) > 4 {
		return nil, errors.New("replace: expected 3 or 4 arguments")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("replace: string %w", err)
	}
	old, err := lib.ToString(args[1])
	if err != nil {
		return nil, fmt.Errorf("replace: old %w", err)
	}
	new, err := lib.ToString(args[2])
	if err != nil {
		return nil, fmt.Errorf("replace: new %w", err)
	}
	n := -1

	if len(args) == 4 {
		nVal, err := lib.ToNumber(args[3])
		if err != nil {
			return nil, fmt.Errorf("replace: count %w", err)
		}
		n = int(nVal)
	}

	return lang.StringValue(strings.Replace(string(str), string(old), string(new), n)), nil
}

func stringReplaceAll(args []lang.Value) (lang.Value, error) {
	if len(args) != 3 {
		return nil, errors.New("replace_all: expected 3 arguments")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("replace_all: string %w", err)
	}
	old, err := lib.ToString(args[1])
	if err != nil {
		return nil, fmt.Errorf("replace_all: old %w", err)
	}
	new, err := lib.ToString(args[2])
	if err != nil {
		return nil, fmt.Errorf("replace_all: new %w", err)
	}
	return lang.StringValue(strings.ReplaceAll(string(str), string(old), string(new))), nil
}

// String Splitting and Joining
func stringSplit(args []lang.Value) (lang.Value, error) {
	if len(args) < 2 || len(args) > 3 {
		return nil, errors.New("split: expected 2 or 3 arguments")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("split: string %w", err)
	}
	sep, err := lib.ToString(args[1])
	if err != nil {
		return nil, fmt.Errorf("split: separator %w", err)
	}
	n := -1

	if len(args) == 3 {
		nVal, err := lib.ToNumber(args[2])
		if err != nil {
			return nil, fmt.Errorf("split: count %w", err)
		}
		n = int(nVal)
	}

	var parts []string
	if n == -1 {
		parts = strings.Split(string(str), string(sep))
	} else {
		parts = strings.SplitN(string(str), string(sep), n)
	}

	result := make(lang.ListValue, len(parts))
	for i, part := range parts {
		result[i] = lang.StringValue(part)
	}
	return result, nil
}

func stringJoin(args []lang.Value) (lang.Value, error) {
	if len(args) != 2 {
		return nil, errors.New("join: expected 2 arguments")
	}
	sep, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("join: separator %w", err)
	}
	list, ok := args[1].(lang.ListValue)
	if !ok {
		return nil, errors.New("join: second argument must be a list")
	}

	parts := make([]string, len(list))
	for i, item := range list {
		str, err := lib.ToString(item)
		if err != nil {
			return nil, fmt.Errorf("join: list item %d %w", i, err)
		}
		parts[i] = string(str)
	}
	return lang.StringValue(strings.Join(parts, string(sep))), nil
}

func stringLines(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("lines: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("lines: %w", err)
	}
	lines := strings.Split(string(str), "\n")
	result := make(lang.ListValue, len(lines))
	for i, line := range lines {
		result[i] = lang.StringValue(line)
	}
	return result, nil
}

func stringFields(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("fields: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("fields: %w", err)
	}
	fields := strings.Fields(string(str))
	result := make(lang.ListValue, len(fields))
	for i, field := range fields {
		result[i] = lang.StringValue(field)
	}
	return result, nil
}

// Regular Expressions
func stringMatch(args []lang.Value) (lang.Value, error) {
	if len(args) != 2 {
		return nil, errors.New("match: expected 2 arguments")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("match: string %w", err)
	}
	pattern, err := lib.ToString(args[1])
	if err != nil {
		return nil, fmt.Errorf("match: pattern %w", err)
	}

	matched, err := regexp.MatchString(string(pattern), string(str))
	if err != nil {
		return nil, fmt.Errorf("match: invalid regex pattern: %w", err)
	}
	return lang.BoolValue(matched), nil
}

func stringFindAll(args []lang.Value) (lang.Value, error) {
	if len(args) < 2 || len(args) > 3 {
		return nil, errors.New("find_all: expected 2 or 3 arguments")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("find_all: string %w", err)
	}
	pattern, err := lib.ToString(args[1])
	if err != nil {
		return nil, fmt.Errorf("find_all: pattern %w", err)
	}
	n := -1

	if len(args) == 3 {
		nVal, err := lib.ToNumber(args[2])
		if err != nil {
			return nil, fmt.Errorf("find_all: count %w", err)
		}
		n = int(nVal)
	}

	re, err := regexp.Compile(string(pattern))
	if err != nil {
		return nil, fmt.Errorf("find_all: invalid regex pattern: %w", err)
	}

	matches := re.FindAllString(string(str), n)
	result := make(lang.ListValue, len(matches))
	for i, match := range matches {
		result[i] = lang.StringValue(match)
	}
	return result, nil
}

func stringReplaceRegex(args []lang.Value) (lang.Value, error) {
	if len(args) != 3 {
		return nil, errors.New("replace_regex: expected 3 arguments")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("replace_regex: string %w", err)
	}
	pattern, err := lib.ToString(args[1])
	if err != nil {
		return nil, fmt.Errorf("replace_regex: pattern %w", err)
	}
	replacement, err := lib.ToString(args[2])
	if err != nil {
		return nil, fmt.Errorf("replace_regex: replacement %w", err)
	}

	re, err := regexp.Compile(string(pattern))
	if err != nil {
		return nil, fmt.Errorf("replace_regex: invalid regex pattern: %w", err)
	}

	return lang.StringValue(re.ReplaceAllString(string(str), string(replacement))), nil
}

// Character Operations
func stringCharAt(args []lang.Value) (lang.Value, error) {
	if len(args) != 2 {
		return nil, errors.New("char_at: expected 2 arguments")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("char_at: string %w", err)
	}
	indexVal, err := lib.ToNumber(args[1])
	if err != nil {
		return nil, fmt.Errorf("char_at: index %w", err)
	}
	index := int(indexVal)

	runes := []rune(string(str))
	if index < 0 || index >= len(runes) {
		return lang.StringValue(""), nil
	}

	return lang.StringValue(string(runes[index])), nil
}

func stringCharCode(args []lang.Value) (lang.Value, error) {
	if len(args) != 2 {
		return nil, errors.New("char_code: expected 2 arguments")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("char_code: string %w", err)
	}
	indexVal, err := lib.ToNumber(args[1])
	if err != nil {
		return nil, fmt.Errorf("char_code: index %w", err)
	}
	index := int(indexVal)

	runes := []rune(string(str))
	if index < 0 || index >= len(runes) {
		return lang.NumberValue(0), nil
	}

	return lang.NumberValue(float64(runes[index])), nil
}

func stringFromCharCode(args []lang.Value) (lang.Value, error) {
	if len(args) == 0 {
		return lang.StringValue(""), nil
	}

	var runes []rune
	for i, arg := range args {
		codeVal, err := lib.ToNumber(arg)
		if err != nil {
			return nil, fmt.Errorf("from_char_code: argument %d %w", i, err)
		}
		code := int(codeVal)
		runes = append(runes, rune(code))
	}

	return lang.StringValue(string(runes)), nil
}

// String Validation
func stringIsEmpty(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("is_empty: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("is_empty: %w", err)
	}
	return lang.BoolValue(len(strings.TrimSpace(string(str))) == 0), nil
}

func stringIsNumeric(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("is_numeric: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("is_numeric: %w", err)
	}
	_, parseErr := strconv.ParseFloat(string(str), 64)
	return lang.BoolValue(parseErr == nil), nil
}

func stringIsAlpha(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("is_alpha: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("is_alpha: %w", err)
	}
	s := string(str)
	if len(s) == 0 {
		return lang.BoolValue(false), nil
	}

	for _, r := range s {
		if !unicode.IsLetter(r) {
			return lang.BoolValue(false), nil
		}
	}
	return lang.BoolValue(true), nil
}

func stringIsAlphaNumeric(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("is_alphanumeric: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("is_alphanumeric: %w", err)
	}
	s := string(str)
	if len(s) == 0 {
		return lang.BoolValue(false), nil
	}

	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return lang.BoolValue(false), nil
		}
	}
	return lang.BoolValue(true), nil
}

func stringIsSpace(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("is_space: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("is_space: %w", err)
	}
	s := string(str)
	if len(s) == 0 {
		return lang.BoolValue(false), nil
	}

	for _, r := range s {
		if !unicode.IsSpace(r) {
			return lang.BoolValue(false), nil
		}
	}
	return lang.BoolValue(true), nil
}

// Type Conversion
func toStr(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("tostring: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("tostring: %w", err)
	}
	return str, nil
}

func stringToNumber(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("tonumber: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("tonumber: %w", err)
	}
	if f, parseErr := strconv.ParseFloat(string(str), 64); parseErr == nil {
		return lang.NumberValue(f), nil
	}
	return lang.NumberValue(0), nil
}

// Functions that would be in the BuiltinFunctions map:
var StringFunctions = map[string]func([]lang.Value) (lang.Value, error){
	// Basic operations
	"len":     stringLen,
	"size":    stringSize,
	"concat":  stringConcat,
	"repeat":  stringRepeat,
	"reverse": stringReverse,

	// Case conversion
	"upper":      stringUpper,
	"lower":      stringLower,
	"title":      stringTitle,
	"capitalize": stringCapitalize,
	"swap_case":  stringSwapCase,

	// Trimming and padding
	"trim":       stringTrim,
	"trim_left":  stringTrimLeft,
	"trim_right": stringTrimRight,
	"pad_left":   stringPadLeft,
	"pad_right":  stringPadRight,
	"pad_center": stringPadCenter,

	// Substring operations
	"substr": stringSubstr,
	"left":   stringLeft,
	"right":  stringRight,

	// Search and replace
	"contains":    stringContains,
	"startswith":  stringStartsWith,
	"endswith":    stringEndsWith,
	"indexof":     stringIndexOf,
	"lastindexof": stringLastIndexOf,
	"replace":     stringReplace,
	"replace_all": stringReplaceAll,

	// Splitting and joining
	"split":  stringSplit,
	"join":   stringJoin,
	"lines":  stringLines,
	"fields": stringFields,

	// Regular expressions
	"match":         stringMatch,
	"find_all":      stringFindAll,
	"replace_regex": stringReplaceRegex,

	// Character operations
	"char_at":        stringCharAt,
	"char_code":      stringCharCode,
	"from_char_code": stringFromCharCode,

	// Validation
	"is_empty":        stringIsEmpty,
	"is_numeric":      stringIsNumeric,
	"is_alpha":        stringIsAlpha,
	"is_alphanumeric": stringIsAlphaNumeric,
	"is_space":        stringIsSpace,

	// Type conversion
	"tostring": toStr,
	"tonumber": stringToNumber,
}
