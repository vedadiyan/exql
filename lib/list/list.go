package list

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"strconv"

	"github.com/vedadiyan/exql/lang"
	"github.com/vedadiyan/exql/lib"
)

func length() (string, lang.Function) {
	name := "length"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		if list, ok := args[0].(lang.ListValue); ok {
			return lang.NumberValue(float64(len(list))), nil
		}
		return nil, lib.ListError(name, args[0])
	}
	return name, fn
}

func isEmpty() (string, lang.Function) {
	name := "isEmpty"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		if list, ok := args[0].(lang.ListValue); ok {
			return lang.BoolValue(len(list) == 0), nil
		}
		return nil, lib.ListError(name, args[0])
	}
	return name, fn
}

func get() (string, lang.Function) {
	name := "get"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) < 2 || len(args) > 3 {
			return nil, errors.New("get: expected 2 or 3 arguments (list, index, default?)")
		}
		list, ok := args[0].(lang.ListValue)
		if !ok {
			return nil, lib.ListError(name, args[0])
		}
		indexNum, err := lib.ToNumber(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: index %w", name, err)
		}
		index := int(indexNum)
		if index < 0 {
			index = len(list) + index
		}
		if index < 0 || index >= len(list) {
			if len(args) == 3 {
				return args[2], nil
			}
			return nil, fmt.Errorf("%s: index %d out of bounds (list length: %d)", name, index, len(list))
		}
		return list[index], nil
	}
	return name, fn
}

func set() (string, lang.Function) {
	name := "set"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 3 {
			return nil, lib.ArgumentError(name, 3)
		}
		list, ok := args[0].(lang.ListValue)
		if !ok {
			return nil, lib.ListError(name, args[0])
		}
		indexNum, err := lib.ToNumber(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: index %w", name, err)
		}
		index := int(indexNum)
		value := args[2]
		if index < 0 {
			index = len(list) + index
		}
		if index < 0 || index >= len(list) {
			return nil, fmt.Errorf("%s: index %d out of bounds (list length: %d)", name, index, len(list))
		}
		result := make(lang.ListValue, len(list))
		copy(result, list)
		result[index] = value
		return result, nil
	}
	return name, fn
}

func aappend() (string, lang.Function) {
	name := "append"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) < 2 {
			return nil, errors.New("append: expected at least 2 arguments (list, ...values)")
		}
		list, ok := args[0].(lang.ListValue)
		if !ok {
			return nil, lib.ListError(name, args[0])
		}
		result := make(lang.ListValue, len(list))
		copy(result, list)
		for i := 1; i < len(args); i++ {
			result = append(result, args[i])
		}
		return result, nil
	}
	return name, fn
}

func prepend() (string, lang.Function) {
	name := "prepend"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) < 2 {
			return nil, errors.New("prepend: expected at least 2 arguments (list, ...values)")
		}
		list, ok := args[0].(lang.ListValue)
		if !ok {
			return nil, lib.ListError(name, args[0])
		}
		var result lang.ListValue
		for i := 1; i < len(args); i++ {
			result = append(result, args[i])
		}
		result = append(result, list...)
		return result, nil
	}
	return name, fn
}

func insert() (string, lang.Function) {
	name := "insert"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 3 {
			return nil, lib.ArgumentError(name, 3)
		}
		list, ok := args[0].(lang.ListValue)
		if !ok {
			return nil, lib.ListError(name, args[0])
		}
		indexNum, err := lib.ToNumber(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: index %w", name, err)
		}
		index := int(indexNum)
		value := args[2]
		if index < 0 {
			index = len(list) + index + 1
		}
		if index < 0 {
			index = 0
		}
		if index > len(list) {
			index = len(list)
		}
		result := make(lang.ListValue, 0, len(list)+1)
		result = append(result, list[:index]...)
		result = append(result, value)
		result = append(result, list[index:]...)
		return result, nil
	}
	return name, fn
}

func remove() (string, lang.Function) {
	name := "remove"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 2 {
			return nil, lib.ArgumentError(name, 2)
		}
		list, ok := args[0].(lang.ListValue)
		if !ok {
			return nil, lib.ListError(name, args[0])
		}
		indexNum, err := lib.ToNumber(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: index %w", name, err)
		}
		index := int(indexNum)
		if index < 0 {
			index = len(list) + index
		}
		if index < 0 || index >= len(list) {
			return nil, fmt.Errorf("%s: index %d out of bounds (list length: %d)", name, index, len(list))
		}
		result := make(lang.ListValue, 0, len(list)-1)
		result = append(result, list[:index]...)
		result = append(result, list[index+1:]...)
		return result, nil
	}
	return name, fn
}

func concat() (string, lang.Function) {
	name := "concat"
	fn := func(args []lang.Value) (lang.Value, error) {
		var result lang.ListValue
		for i, arg := range args {
			if list, ok := arg.(lang.ListValue); ok {
				result = append(result, list...)
			} else {
				return nil, fmt.Errorf("%s: argument %d expected list, got %T", name, i+1, arg)
			}
		}
		return result, nil
	}
	return name, fn
}

func first() (string, lang.Function) {
	name := "first"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) < 1 || len(args) > 2 {
			return nil, errors.New("first: expected 1 or 2 arguments (list, default?)")
		}
		list, ok := args[0].(lang.ListValue)
		if !ok {
			return nil, lib.ListError(name, args[0])
		}
		if len(list) == 0 {
			if len(args) == 2 {
				return args[1], nil
			}
			return nil, errors.New("first: list is empty")
		}
		return list[0], nil
	}
	return name, fn
}

func last() (string, lang.Function) {
	name := "last"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) < 1 || len(args) > 2 {
			return nil, errors.New("last: expected 1 or 2 arguments (list, default?)")
		}
		list, ok := args[0].(lang.ListValue)
		if !ok {
			return nil, lib.ListError(name, args[0])
		}
		if len(list) == 0 {
			if len(args) == 2 {
				return args[1], nil
			}
			return nil, errors.New("last: list is empty")
		}
		return list[len(list)-1], nil
	}
	return name, fn
}

func head() (string, lang.Function) {
	name := "head"
	_, First := first()
	fn := func(args []lang.Value) (lang.Value, error) {
		return First(args)
	}
	return name, fn
}

func tail() (string, lang.Function) {
	name := "tail"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		list, ok := args[0].(lang.ListValue)
		if !ok {
			return nil, lib.ListError(name, args[0])
		}
		if len(list) <= 1 {
			return lang.ListValue{}, nil
		}
		result := make(lang.ListValue, len(list)-1)
		copy(result, list[1:])
		return result, nil
	}
	return name, fn
}

func rest() (string, lang.Function) {
	name := "rest"
	_, Tail := tail()
	fn := func(args []lang.Value) (lang.Value, error) {
		return Tail(args)
	}
	return name, fn
}

func iinit() (string, lang.Function) {
	name := "init"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		list, ok := args[0].(lang.ListValue)
		if !ok {
			return nil, lib.ListError(name, args[0])
		}
		if len(list) <= 1 {
			return lang.ListValue{}, nil
		}
		result := make(lang.ListValue, len(list)-1)
		copy(result, list[:len(list)-1])
		return result, nil
	}
	return name, fn
}

func slice() (string, lang.Function) {
	name := "slice"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) < 2 || len(args) > 3 {
			return nil, errors.New("slice: expected 2 or 3 arguments (list, start, end?)")
		}
		list, ok := args[0].(lang.ListValue)
		if !ok {
			return nil, lib.ListError(name, args[0])
		}
		startNum, err := lib.ToNumber(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: start %w", name, err)
		}
		start := int(startNum)
		end := len(list)
		if len(args) == 3 {
			endNum, err := lib.ToNumber(args[2])
			if err != nil {
				return nil, fmt.Errorf("%s: end %w", name, err)
			}
			end = int(endNum)
		}
		if start < 0 {
			start = len(list) + start
		}
		if end < 0 {
			end = len(list) + end
		}
		if start < 0 {
			start = 0
		}
		if end > len(list) {
			end = len(list)
		}
		if start > end {
			return lang.ListValue{}, nil
		}
		result := make(lang.ListValue, end-start)
		copy(result, list[start:end])
		return result, nil
	}
	return name, fn
}

func take() (string, lang.Function) {
	name := "take"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 2 {
			return nil, lib.ArgumentError(name, 2)
		}
		list, ok := args[0].(lang.ListValue)
		if !ok {
			return nil, lib.ListError(name, args[0])
		}
		countNum, err := lib.ToNumber(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: count %w", name, err)
		}
		count := int(countNum)
		if count < 0 {
			return nil, errors.New("take: count must be non-negative")
		}
		if count == 0 {
			return lang.ListValue{}, nil
		}
		if count >= len(list) {
			result := make(lang.ListValue, len(list))
			copy(result, list)
			return result, nil
		}
		result := make(lang.ListValue, count)
		copy(result, list[:count])
		return result, nil
	}
	return name, fn
}

func drop() (string, lang.Function) {
	name := "drop"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 2 {
			return nil, lib.ArgumentError(name, 2)
		}
		list, ok := args[0].(lang.ListValue)
		if !ok {
			return nil, lib.ListError(name, args[0])
		}
		countNum, err := lib.ToNumber(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: count %w", name, err)
		}
		count := int(countNum)
		if count < 0 {
			return nil, errors.New("drop: count must be non-negative")
		}
		if count == 0 {
			result := make(lang.ListValue, len(list))
			copy(result, list)
			return result, nil
		}
		if count >= len(list) {
			return lang.ListValue{}, nil
		}
		result := make(lang.ListValue, len(list)-count)
		copy(result, list[count:])
		return result, nil
	}
	return name, fn
}

func reverse() (string, lang.Function) {
	name := "reverse"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		list, ok := args[0].(lang.ListValue)
		if !ok {
			return nil, lib.ListError(name, args[0])
		}
		result := make(lang.ListValue, len(list))
		for i, v := range list {
			result[len(list)-1-i] = v
		}
		return result, nil
	}
	return name, fn
}

func ssort() (string, lang.Function) {
	name := "sort"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		list, ok := args[0].(lang.ListValue)
		if !ok {
			return nil, lib.ListError(name, args[0])
		}
		result := make(lang.ListValue, len(list))
		copy(result, list)
		sort.Slice(result, func(i, j int) bool {
			return compareValues(result[i], result[j]) < 0
		})
		return result, nil
	}
	return name, fn
}

func sortDesc() (string, lang.Function) {
	name := "sortDesc"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		list, ok := args[0].(lang.ListValue)
		if !ok {
			return nil, lib.ListError(name, args[0])
		}
		result := make(lang.ListValue, len(list))
		copy(result, list)
		sort.Slice(result, func(i, j int) bool {
			return compareValues(result[i], result[j]) > 0
		})
		return result, nil
	}
	return name, fn
}

func shuffle() (string, lang.Function) {
	name := "shuffle"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		list, ok := args[0].(lang.ListValue)
		if !ok {
			return nil, lib.ListError(name, args[0])
		}
		result := make(lang.ListValue, len(list))
		copy(result, list)
		for i := len(result) - 1; i > 0; i-- {
			j := rand.Intn(i + 1)
			result[i], result[j] = result[j], result[i]
		}
		return result, nil
	}
	return name, fn
}

func unique() (string, lang.Function) {
	name := "unique"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		list, ok := args[0].(lang.ListValue)
		if !ok {
			return nil, lib.ListError(name, args[0])
		}
		seen := make(map[string]bool)
		var result lang.ListValue
		for _, item := range list {
			key := valueToString(item)
			if !seen[key] {
				seen[key] = true
				result = append(result, item)
			}
		}
		return result, nil
	}
	return name, fn
}

func flatten() (string, lang.Function) {
	name := "flatten"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) < 1 || len(args) > 2 {
			return nil, errors.New("flatten: expected 1 or 2 arguments (list, depth?)")
		}
		list, ok := args[0].(lang.ListValue)
		if !ok {
			return nil, lib.ListError(name, args[0])
		}
		depth := 1
		if len(args) == 2 {
			depthNum, err := lib.ToNumber(args[1])
			if err != nil {
				return nil, fmt.Errorf("%s: depth %w", name, err)
			}
			depth = int(depthNum)
			if depth < 0 {
				return nil, errors.New("flatten: depth must be non-negative")
			}
		}
		return flattenList(list, depth), nil
	}
	return name, fn
}

func flattenList(list lang.ListValue, depth int) lang.ListValue {
	if depth <= 0 {
		return list
	}
	var result lang.ListValue
	for _, item := range list {
		if subList, ok := item.(lang.ListValue); ok {
			flattened := flattenList(subList, depth-1)
			result = append(result, flattened...)
		} else {
			result = append(result, item)
		}
	}
	return result
}

func contains() (string, lang.Function) {
	name := "contains"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 2 {
			return nil, lib.ArgumentError(name, 2)
		}
		list, ok := args[0].(lang.ListValue)
		if !ok {
			return nil, lib.ListError(name, args[0])
		}
		searchValue := args[1]
		for _, item := range list {
			if equalValues(item, searchValue) {
				return lang.BoolValue(true), nil
			}
		}
		return lang.BoolValue(false), nil
	}
	return name, fn
}

func indexOf() (string, lang.Function) {
	name := "indexOf"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) < 2 || len(args) > 3 {
			return nil, fmt.Errorf("%s: expected 2 or 3 arguments (list, value, start?)", name)
		}
		list, ok := args[0].(lang.ListValue)
		if !ok {
			return nil, lib.ListError(name, args[0])
		}
		searchValue := args[1]
		start := 0
		if len(args) == 3 {
			startNum, err := lib.ToNumber(args[2])
			if err != nil {
				return nil, fmt.Errorf("%s: start %w", name, err)
			}
			start = int(startNum)
			if start < 0 {
				start = 0
			}
		}
		for i := start; i < len(list); i++ {
			if equalValues(list[i], searchValue) {
				return lang.NumberValue(float64(i)), nil
			}
		}
		return lang.NumberValue(-1), nil
	}
	return name, fn
}

func lastIndexOf() (string, lang.Function) {
	name := "lastIndexOf"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 2 {
			return nil, lib.ArgumentError(name, 2)
		}
		list, ok := args[0].(lang.ListValue)
		if !ok {
			return nil, lib.ListError(name, args[0])
		}
		searchValue := args[1]
		for i := len(list) - 1; i >= 0; i-- {
			if equalValues(list[i], searchValue) {
				return lang.NumberValue(float64(i)), nil
			}
		}
		return lang.NumberValue(-1), nil
	}
	return name, fn
}

func count() (string, lang.Function) {
	name := "count"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 2 {
			return nil, lib.ArgumentError(name, 2)
		}
		list, ok := args[0].(lang.ListValue)
		if !ok {
			return nil, lib.ListError(name, args[0])
		}
		searchValue := args[1]
		count := 0
		for _, item := range list {
			if equalValues(item, searchValue) {
				count++
			}
		}
		return lang.NumberValue(float64(count)), nil
	}
	return name, fn
}

func rrange() (string, lang.Function) {
	name := "range"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) < 1 || len(args) > 3 {
			return nil, errors.New("range: expected 1, 2, or 3 arguments (end) or (start, end) or (start, end, step)")
		}
		var start, end, step float64
		var err error
		if len(args) == 1 {
			start = 0
			if end, err = lib.ToNumber(args[0]); err != nil {
				return nil, fmt.Errorf("range: end %w", err)
			}
			step = 1
		} else if len(args) == 2 {
			if start, err = lib.ToNumber(args[0]); err != nil {
				return nil, fmt.Errorf("range: start %w", err)
			}
			if end, err = lib.ToNumber(args[1]); err != nil {
				return nil, fmt.Errorf("range: end %w", err)
			}
			step = 1
		} else {
			if start, err = lib.ToNumber(args[0]); err != nil {
				return nil, fmt.Errorf("range: start %w", err)
			}
			if end, err = lib.ToNumber(args[1]); err != nil {
				return nil, fmt.Errorf("range: end %w", err)
			}
			if step, err = lib.ToNumber(args[2]); err != nil {
				return nil, fmt.Errorf("range: step %w", err)
			}
		}
		if step == 0 {
			return nil, errors.New("range: step cannot be zero")
		}
		var result lang.ListValue
		if step > 0 {
			for i := start; i < end; i += step {
				result = append(result, lang.NumberValue(i))
			}
		} else {
			for i := start; i > end; i += step {
				result = append(result, lang.NumberValue(i))
			}
		}
		return result, nil
	}
	return name, fn
}

func repeat() (string, lang.Function) {
	name := "repeat"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 2 {
			return nil, lib.ArgumentError(name, 2)
		}
		value := args[0]
		countNum, err := lib.ToNumber(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: count %w", name, err)
		}
		count := int(countNum)
		if count < 0 {
			return nil, errors.New("repeat: count must be non-negative")
		}
		if count == 0 {
			return lang.ListValue{}, nil
		}
		result := make(lang.ListValue, count)
		for i := 0; i < count; i++ {
			result[i] = value
		}
		return result, nil
	}
	return name, fn
}

func zip() (string, lang.Function) {
	name := "zip"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) == 0 {
			return nil, errors.New("zip: expected at least 1 argument")
		}
		var lists []lang.ListValue
		minLen := -1
		for i, arg := range args {
			if list, ok := arg.(lang.ListValue); ok {
				lists = append(lists, list)
				if minLen == -1 || len(list) < minLen {
					minLen = len(list)
				}
			} else {
				return nil, fmt.Errorf("%s: argument %d expected list, got %T", name, i+1, arg)
			}
		}
		if minLen <= 0 {
			return lang.ListValue{}, nil
		}
		result := make(lang.ListValue, minLen)
		for i := 0; i < minLen; i++ {
			tuple := make(lang.ListValue, len(lists))
			for j, list := range lists {
				tuple[j] = list[i]
			}
			result[i] = tuple
		}
		return result, nil
	}
	return name, fn
}

func filter() (string, lang.Function) {
	name := "filter"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		list, ok := args[0].(lang.ListValue)
		if !ok {
			return nil, lib.ListError(name, args[0])
		}
		var result lang.ListValue
		for _, item := range list {
			if item != nil && !isNullValue(item) {
				result = append(result, item)
			}
		}
		return result, nil
	}
	return name, fn
}

func mmap() (string, lang.Function) {
	name := "map"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		list, ok := args[0].(lang.ListValue)
		if !ok {
			return nil, lib.ListError(name, args[0])
		}
		result := make(lang.ListValue, len(list))
		copy(result, list)
		return result, nil
	}
	return name, fn
}

func compareValues(a, b lang.Value) int {
	aNum, aErr := lib.ToNumber(a)
	bNum, bErr := lib.ToNumber(b)
	if aErr == nil && bErr == nil {
		if aNum < bNum {
			return -1
		}
		if aNum > bNum {
			return 1
		}
		return 0
	}
	aStr := valueToString(a)
	bStr := valueToString(b)
	if aStr < bStr {
		return -1
	}
	if aStr > bStr {
		return 1
	}
	return 0
}

func equalValues(a, b lang.Value) bool {
	return valueToString(a) == valueToString(b)
}

func valueToString(v lang.Value) string {
	switch val := v.(type) {
	case lang.StringValue:
		return string(val)
	case lang.NumberValue:
		return strconv.FormatFloat(float64(val), 'g', -1, 64)
	case lang.BoolValue:
		if val {
			return "true"
		}
		return "false"
	case nil:
		return ""
	default:
		return fmt.Sprintf("%v", v)
	}
}

func isNullValue(v lang.Value) bool {
	switch val := v.(type) {
	case lang.StringValue:
		return string(val) == ""
	case lang.NumberValue:
		return float64(val) == 0
	case lang.BoolValue:
		return !bool(val)
	case nil:
		return true
	default:
		return false
	}
}

var listFunctions = []func() (string, lang.Function){
	length,
	isEmpty,
	get,
	set,
	aappend,
	prepend,
	insert,
	remove,
	concat,
	first,
	last,
	head,
	tail,
	rest,
	iinit,
	slice,
	take,
	drop,
	reverse,
	ssort,
	sortDesc,
	shuffle,
	unique,
	flatten,
	contains,
	indexOf,
	lastIndexOf,
	count,
	rrange,
	repeat,
	zip,
	filter,
	mmap,
}

func Export() map[string]lang.Function {
	out := make(map[string]lang.Function)
	for _, value := range listFunctions {
		name, fn := value()
		out[name] = fn
	}
	return out
}
