package string

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/vedadiyan/exql/lang"
	"github.com/vedadiyan/exql/lib"
)

func length() (string, func([]lang.Value) (lang.Value, error)) {
	name := "len"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, lib.StringError(name, args[0])
		}
		return lang.NumberValue(float64(utf8.RuneCountInString(string(str)))), nil
	}
	return name, fn
}

func size() (string, func([]lang.Value) (lang.Value, error)) {
	name := "size"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, lib.StringError(name, args[0])
		}
		return lang.NumberValue(float64(len(string(str)))), nil
	}
	return name, fn
}

func voncat() (string, func([]lang.Value) (lang.Value, error)) {
	name := "concat"
	fn := func(args []lang.Value) (lang.Value, error) {
		var result string
		for i, arg := range args {
			str, err := lib.ToString(arg)
			if err != nil {
				return nil, fmt.Errorf("%s: argument %d %w", name, i, err)
			}
			result += string(str)
		}
		return lang.StringValue(result), nil
	}
	return name, fn
}

func repeat() (string, func([]lang.Value) (lang.Value, error)) {
	name := "repeat"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 2 {
			return nil, lib.ArgumentError(name, 2)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: string %w", name, err)
		}
		countVal, err := lib.ToNumber(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: count %w", name, err)
		}
		count := int(countVal)
		if count < 0 {
			count = 0
		}
		return lang.StringValue(strings.Repeat(string(str), count)), nil
	}
	return name, fn
}

func reverse() (string, func([]lang.Value) (lang.Value, error)) {
	name := "reverse"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, lib.StringError(name, args[0])
		}
		runes := []rune(string(str))
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		return lang.StringValue(string(runes)), nil
	}
	return name, fn
}

func toUpper() (string, func([]lang.Value) (lang.Value, error)) {
	name := "upper"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, lib.StringError(name, args[0])
		}
		return lang.StringValue(strings.ToUpper(string(str))), nil
	}
	return name, fn
}

func toLower() (string, func([]lang.Value) (lang.Value, error)) {
	name := "lower"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, lib.StringError(name, args[0])
		}
		return lang.StringValue(strings.ToLower(string(str))), nil
	}
	return name, fn
}

func title() (string, func([]lang.Value) (lang.Value, error)) {
	name := "title"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, lib.StringError(name, args[0])
		}
		return lang.StringValue(strings.Title(string(str))), nil
	}
	return name, fn
}

func capitalize() (string, func([]lang.Value) (lang.Value, error)) {
	name := "capitalize"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, lib.StringError(name, args[0])
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
	return name, fn
}

func swapCase() (string, func([]lang.Value) (lang.Value, error)) {
	name := "swap_case"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, lib.StringError(name, args[0])
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
	return name, fn
}

func trim() (string, func([]lang.Value) (lang.Value, error)) {
	name := "trim"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) < 1 || len(args) > 2 {
			return nil, fmt.Errorf("%s: expected 1 or 2 arguments", name)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, lib.StringError(name, args[0])
		}
		if len(args) == 1 {
			return lang.StringValue(strings.TrimSpace(string(str))), nil
		}
		cutset, err := lib.ToString(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: cutset %w", name, err)
		}
		return lang.StringValue(strings.Trim(string(str), string(cutset))), nil
	}
	return name, fn
}

func trimLeft() (string, func([]lang.Value) (lang.Value, error)) {
	name := "trim_left"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) < 1 || len(args) > 2 {
			return nil, fmt.Errorf("%s: expected 1 or 2 arguments", name)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, lib.StringError(name, args[0])
		}
		if len(args) == 1 {
			return lang.StringValue(strings.TrimLeftFunc(string(str), unicode.IsSpace)), nil
		}
		cutset, err := lib.ToString(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: cutset %w", name, err)
		}
		return lang.StringValue(strings.TrimLeft(string(str), string(cutset))), nil
	}
	return name, fn
}

func trimRight() (string, func([]lang.Value) (lang.Value, error)) {
	name := "trim_right"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) < 1 || len(args) > 2 {
			return nil, fmt.Errorf("%s: expected 1 or 2 arguments", name)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, lib.StringError(name, args[0])
		}
		if len(args) == 1 {
			return lang.StringValue(strings.TrimRightFunc(string(str), unicode.IsSpace)), nil
		}
		cutset, err := lib.ToString(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: cutset %w", name, err)
		}
		return lang.StringValue(strings.TrimRight(string(str), string(cutset))), nil
	}
	return name, fn
}

func padLeft() (string, func([]lang.Value) (lang.Value, error)) {
	name := "pad_left"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) < 2 || len(args) > 3 {
			return nil, fmt.Errorf("%s: expected 2 or 3 arguments", name)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, lib.StringError(name, args[0])
		}
		totalLenVal, err := lib.ToNumber(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: total length %w", name, err)
		}
		totalLen := int(totalLenVal)
		padChar := " "
		if len(args) == 3 {
			padCharVal, err := lib.ToString(args[2])
			if err != nil {
				return nil, fmt.Errorf("%s: pad character %w", name, err)
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
	return name, fn
}

func padRight() (string, func([]lang.Value) (lang.Value, error)) {
	name := "pad_right"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) < 2 || len(args) > 3 {
			return nil, fmt.Errorf("%s: expected 2 or 3 arguments", name)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, lib.StringError(name, args[0])
		}
		totalLenVal, err := lib.ToNumber(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: total length %w", name, err)
		}
		totalLen := int(totalLenVal)
		padChar := " "
		if len(args) == 3 {
			padCharVal, err := lib.ToString(args[2])
			if err != nil {
				return nil, fmt.Errorf("%s: pad character %w", name, err)
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
	return name, fn
}

func padCenter() (string, func([]lang.Value) (lang.Value, error)) {
	name := "pad_center"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) < 2 || len(args) > 3 {
			return nil, fmt.Errorf("%s: expected 2 or 3 arguments", name)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, lib.StringError(name, args[0])
		}
		totalLenVal, err := lib.ToNumber(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: total length %w", name, err)
		}
		totalLen := int(totalLenVal)
		padChar := " "
		if len(args) == 3 {
			padCharVal, err := lib.ToString(args[2])
			if err != nil {
				return nil, fmt.Errorf("%s: pad character %w", name, err)
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
	return name, fn
}

func substr() (string, func([]lang.Value) (lang.Value, error)) {
	name := "substr"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) < 2 || len(args) > 3 {
			return nil, fmt.Errorf("%s: expected 2 or 3 arguments", name)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, lib.StringError(name, args[0])
		}
		startVal, err := lib.ToNumber(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: start %w", name, err)
		}
		start := int(startVal)

		runes := []rune(string(str))
		strLen := len(runes)

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
			return nil, fmt.Errorf("%s: length %w", name, err)
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
	return name, fn
}

func left() (string, func([]lang.Value) (lang.Value, error)) {
	name := "left"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 2 {
			return nil, lib.ArgumentError(name, 2)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, lib.StringError(name, args[0])
		}
		countVal, err := lib.ToNumber(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: count %w", name, err)
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
	return name, fn
}

func right() (string, func([]lang.Value) (lang.Value, error)) {
	name := "right"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 2 {
			return nil, lib.ArgumentError(name, 2)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, lib.StringError(name, args[0])
		}
		countVal, err := lib.ToNumber(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: count %w", name, err)
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
	return name, fn
}

func contains() (string, func([]lang.Value) (lang.Value, error)) {
	name := "contains"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 2 {
			return nil, lib.ArgumentError(name, 2)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, lib.StringError(name, args[0])
		}
		substr, err := lib.ToString(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: substring %w", name, err)
		}
		return lang.BoolValue(strings.Contains(string(str), string(substr))), nil
	}
	return name, fn
}

func startsWith() (string, func([]lang.Value) (lang.Value, error)) {
	name := "startswith"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 2 {
			return nil, lib.ArgumentError(name, 2)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, lib.StringError(name, args[0])
		}
		prefix, err := lib.ToString(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: prefix %w", name, err)
		}
		return lang.BoolValue(strings.HasPrefix(string(str), string(prefix))), nil
	}
	return name, fn
}

func endsWith() (string, func([]lang.Value) (lang.Value, error)) {
	name := "endswith"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 2 {
			return nil, lib.ArgumentError(name, 2)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, lib.StringError(name, args[0])
		}
		suffix, err := lib.ToString(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: suffix %w", name, err)
		}
		return lang.BoolValue(strings.HasSuffix(string(str), string(suffix))), nil
	}
	return name, fn
}

func indexOf() (string, func([]lang.Value) (lang.Value, error)) {
	name := "indexof"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) < 2 || len(args) > 3 {
			return nil, fmt.Errorf("%s: expected 2 or 3 arguments", name)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, lib.StringError(name, args[0])
		}
		substr, err := lib.ToString(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: substring %w", name, err)
		}
		start := 0

		if len(args) == 3 {
			startVal, err := lib.ToNumber(args[2])
			if err != nil {
				return nil, fmt.Errorf("%s: start %w", name, err)
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
	return name, fn
}

func lastIndexOf() (string, func([]lang.Value) (lang.Value, error)) {
	name := "last_index_of"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 2 {
			return nil, lib.ArgumentError(name, 2)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, lib.StringError(name, args[0])
		}
		substr, err := lib.ToString(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: substring %w", name, err)
		}

		s := string(str)
		sub := string(substr)
		index := strings.LastIndex(s, sub)
		if index == -1 {
			return lang.NumberValue(-1), nil
		}

		return lang.NumberValue(float64(utf8.RuneCountInString(s[:index]))), nil
	}
	return name, fn
}

func replace() (string, func([]lang.Value) (lang.Value, error)) {
	name := "replace"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) < 3 || len(args) > 4 {
			return nil, fmt.Errorf("%s: expected 3 or 4 arguments", name)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, lib.StringError(name, args[0])
		}
		old, err := lib.ToString(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: old %w", name, err)
		}
		new, err := lib.ToString(args[2])
		if err != nil {
			return nil, fmt.Errorf("%s: new %w", name, err)
		}
		n := -1

		if len(args) == 4 {
			nVal, err := lib.ToNumber(args[3])
			if err != nil {
				return nil, fmt.Errorf("%s: count %w", name, err)
			}
			n = int(nVal)
		}

		return lang.StringValue(strings.Replace(string(str), string(old), string(new), n)), nil
	}
	return name, fn
}

func replaceAll() (string, func([]lang.Value) (lang.Value, error)) {
	name := "replace_all"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 3 {
			return nil, lib.ArgumentError(name, 3)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, lib.StringError(name, args[0])
		}
		old, err := lib.ToString(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: old %w", name, err)
		}
		new, err := lib.ToString(args[2])
		if err != nil {
			return nil, fmt.Errorf("%s: new %w", name, err)
		}
		return lang.StringValue(strings.ReplaceAll(string(str), string(old), string(new))), nil
	}
	return name, fn
}

func split() (string, func([]lang.Value) (lang.Value, error)) {
	name := "split"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) < 2 || len(args) > 3 {
			return nil, fmt.Errorf("%s: expected 2 or 3 arguments", name)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, lib.StringError(name, args[0])
		}
		sep, err := lib.ToString(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: separator %w", name, err)
		}
		n := -1

		if len(args) == 3 {
			nVal, err := lib.ToNumber(args[2])
			if err != nil {
				return nil, fmt.Errorf("%s: count %w", name, err)
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
	return name, fn
}

func join() (string, func([]lang.Value) (lang.Value, error)) {
	name := "join"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 2 {
			return nil, lib.ArgumentError(name, 2)
		}
		sep, err := lib.ToString(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: separator %w", name, err)
		}
		list, ok := args[1].(lang.ListValue)
		if !ok {
			return nil, lib.ListError(name, args[1])
		}

		parts := make([]string, len(list))
		for i, item := range list {
			str, err := lib.ToString(item)
			if err != nil {
				return nil, fmt.Errorf("%s: list item %d %w", name, i, err)
			}
			parts[i] = string(str)
		}
		return lang.StringValue(strings.Join(parts, string(sep))), nil
	}
	return name, fn
}

func lines() (string, func([]lang.Value) (lang.Value, error)) {
	name := "lines"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, lib.StringError(name, args[0])
		}
		lines := strings.Split(string(str), "\n")
		result := make(lang.ListValue, len(lines))
		for i, line := range lines {
			result[i] = lang.StringValue(line)
		}
		return result, nil
	}
	return name, fn
}

func fields() (string, func([]lang.Value) (lang.Value, error)) {
	name := "fields"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, lib.StringError(name, args[0])
		}
		fields := strings.Fields(string(str))
		result := make(lang.ListValue, len(fields))
		for i, field := range fields {
			result[i] = lang.StringValue(field)
		}
		return result, nil
	}
	return name, fn
}

func match() (string, func([]lang.Value) (lang.Value, error)) {
	name := "match"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 2 {
			return nil, lib.ArgumentError(name, 2)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, lib.StringError(name, args[0])
		}
		pattern, err := lib.ToString(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: pattern %w", name, err)
		}

		matched, err := regexp.MatchString(string(pattern), string(str))
		if err != nil {
			return nil, fmt.Errorf("%s: invalid regex pattern: %w", name, err)
		}
		return lang.BoolValue(matched), nil
	}
	return name, fn
}

func findAll() (string, func([]lang.Value) (lang.Value, error)) {
	name := "find_all"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) < 2 || len(args) > 3 {
			return nil, fmt.Errorf("%s: expected 2 or 3 arguments", name)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, lib.StringError(name, args[0])
		}
		pattern, err := lib.ToString(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: pattern %w", name, err)
		}
		n := -1

		if len(args) == 3 {
			nVal, err := lib.ToNumber(args[2])
			if err != nil {
				return nil, fmt.Errorf("%s: count %w", name, err)
			}
			n = int(nVal)
		}

		re, err := regexp.Compile(string(pattern))
		if err != nil {
			return nil, fmt.Errorf("%s: invalid regex pattern: %w", name, err)
		}

		matches := re.FindAllString(string(str), n)
		result := make(lang.ListValue, len(matches))
		for i, match := range matches {
			result[i] = lang.StringValue(match)
		}
		return result, nil
	}
	return name, fn
}

func replaceRegex() (string, func([]lang.Value) (lang.Value, error)) {
	name := "replace_regex"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 3 {
			return nil, lib.ArgumentError(name, 3)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, lib.StringError(name, args[0])
		}
		pattern, err := lib.ToString(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: pattern %w", name, err)
		}
		replacement, err := lib.ToString(args[2])
		if err != nil {
			return nil, fmt.Errorf("%s: replacement %w", name, err)
		}

		re, err := regexp.Compile(string(pattern))
		if err != nil {
			return nil, fmt.Errorf("%s: invalid regex pattern: %w", name, err)
		}

		return lang.StringValue(re.ReplaceAllString(string(str), string(replacement))), nil
	}
	return name, fn
}

func charAt() (string, func([]lang.Value) (lang.Value, error)) {
	name := "char_at"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 2 {
			return nil, lib.ArgumentError(name, 2)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, lib.StringError(name, args[0])
		}
		indexVal, err := lib.ToNumber(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: index %w", name, err)
		}
		index := int(indexVal)

		runes := []rune(string(str))
		if index < 0 || index >= len(runes) {
			return lang.StringValue(""), nil
		}

		return lang.StringValue(string(runes[index])), nil
	}
	return name, fn
}

func charCode() (string, func([]lang.Value) (lang.Value, error)) {
	name := "char_code"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 2 {
			return nil, lib.ArgumentError(name, 2)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, lib.StringError(name, args[0])
		}
		indexVal, err := lib.ToNumber(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: index %w", name, err)
		}
		index := int(indexVal)

		runes := []rune(string(str))
		if index < 0 || index >= len(runes) {
			return lang.NumberValue(0), nil
		}

		return lang.NumberValue(float64(runes[index])), nil
	}
	return name, fn
}

func fromCharCode() (string, func([]lang.Value) (lang.Value, error)) {
	name := "from_char_code"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) == 0 {
			return lang.StringValue(""), nil
		}

		var runes []rune
		for i, arg := range args {
			codeVal, err := lib.ToNumber(arg)
			if err != nil {
				return nil, fmt.Errorf("%s: argument %d %w", name, i, err)
			}
			code := int(codeVal)
			runes = append(runes, rune(code))
		}

		return lang.StringValue(string(runes)), nil
	}
	return name, fn
}

func isEmpty() (string, func([]lang.Value) (lang.Value, error)) {
	name := "is_empty"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, lib.StringError(name, args[0])
		}
		return lang.BoolValue(len(strings.TrimSpace(string(str))) == 0), nil
	}
	return name, fn
}

func isNumeric() (string, func([]lang.Value) (lang.Value, error)) {
	name := "is_numeric"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, lib.StringError(name, args[0])
		}
		_, parseErr := strconv.ParseFloat(string(str), 64)
		return lang.BoolValue(parseErr == nil), nil
	}
	return name, fn
}

func isAlpha() (string, func([]lang.Value) (lang.Value, error)) {
	name := "is_alpha"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, lib.StringError(name, args[0])
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
	return name, fn
}

func isAlphaNumeric() (string, func([]lang.Value) (lang.Value, error)) {
	name := "is_alphanumeric"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, lib.StringError(name, args[0])
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
	return name, fn
}

func isSpace() (string, func([]lang.Value) (lang.Value, error)) {
	name := "is_space"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, lib.StringError(name, args[0])
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
	return name, fn
}

func toStr() (string, func([]lang.Value) (lang.Value, error)) {
	name := "tostring"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, lib.StringError(name, args[0])
		}
		return str, nil
	}
	return name, fn
}

func toNumber() (string, func([]lang.Value) (lang.Value, error)) {
	name := "tonumber"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, lib.StringError(name, args[0])
		}
		if f, parseErr := strconv.ParseFloat(string(str), 64); parseErr == nil {
			return lang.NumberValue(f), nil
		}
		return lang.NumberValue(0), nil
	}
	return name, fn
}

var StringFunctions = []func() (string, func([]lang.Value) (lang.Value, error)){

	length,
	size,
	voncat,
	repeat,
	reverse,

	toUpper,
	toLower,
	title,
	capitalize,
	swapCase,

	trim,
	trimLeft,
	trimRight,
	padLeft,
	padRight,
	padCenter,

	substr,
	left,
	right,

	contains,
	startsWith,
	endsWith,
	indexOf,
	lastIndexOf,
	replace,
	replaceAll,

	split,
	join,
	lines,
	fields,

	match,
	findAll,
	replaceRegex,

	charAt,
	charCode,
	fromCharCode,

	isEmpty,
	isNumeric,
	isAlpha,
	isAlphaNumeric,
	isSpace,

	toStr,
	toNumber,
}
