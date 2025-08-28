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

// List/Array Functions
// These functions provide comprehensive list and array manipulation capabilities

// Basic List Operations
func listLength(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("list_length: expected 1 argument")
	}
	if list, ok := args[0].(lang.ListValue); ok {
		return lang.NumberValue(float64(len(list))), nil
	}
	return nil, fmt.Errorf("list_length: expected list, got %T", args[0])
}

func listIsEmpty(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("list_is_empty: expected 1 argument")
	}
	if list, ok := args[0].(lang.ListValue); ok {
		return lang.BoolValue(len(list) == 0), nil
	}
	return nil, fmt.Errorf("list_is_empty: expected list, got %T", args[0])
}

func listGet(args []lang.Value) (lang.Value, error) {
	if len(args) < 2 || len(args) > 3 {
		return nil, errors.New("list_get: expected 2 or 3 arguments (list, index, default?)")
	}
	list, ok := args[0].(lang.ListValue)
	if !ok {
		return nil, fmt.Errorf("list_get: expected list, got %T", args[0])
	}

	indexNum, err := lib.ToNumber(args[1])
	if err != nil {
		return nil, fmt.Errorf("list_get: index %w", err)
	}
	index := int(indexNum)

	// Handle negative indices
	if index < 0 {
		index = len(list) + index
	}

	if index < 0 || index >= len(list) {
		// Return default value if provided
		if len(args) == 3 {
			return args[2], nil
		}
		return nil, fmt.Errorf("list_get: index %d out of bounds (list length: %d)", index, len(list))
	}

	return list[index], nil
}

func listSet(args []lang.Value) (lang.Value, error) {
	if len(args) != 3 {
		return nil, errors.New("list_set: expected 3 arguments (list, index, value)")
	}
	list, ok := args[0].(lang.ListValue)
	if !ok {
		return nil, fmt.Errorf("list_set: expected list, got %T", args[0])
	}

	indexNum, err := lib.ToNumber(args[1])
	if err != nil {
		return nil, fmt.Errorf("list_set: index %w", err)
	}
	index := int(indexNum)
	value := args[2]

	// Handle negative indices
	if index < 0 {
		index = len(list) + index
	}

	if index < 0 || index >= len(list) {
		return nil, fmt.Errorf("list_set: index %d out of bounds (list length: %d)", index, len(list))
	}

	// Create a copy and modify
	result := make(lang.ListValue, len(list))
	copy(result, list)
	result[index] = value

	return result, nil
}

// List Construction
func listAppend(args []lang.Value) (lang.Value, error) {
	if len(args) < 2 {
		return nil, errors.New("append: expected at least 2 arguments (list, ...values)")
	}
	list, ok := args[0].(lang.ListValue)
	if !ok {
		return nil, fmt.Errorf("append: expected list, got %T", args[0])
	}

	result := make(lang.ListValue, len(list))
	copy(result, list)

	for i := 1; i < len(args); i++ {
		result = append(result, args[i])
	}

	return result, nil
}

func listPrepend(args []lang.Value) (lang.Value, error) {
	if len(args) < 2 {
		return nil, errors.New("prepend: expected at least 2 arguments (list, ...values)")
	}
	list, ok := args[0].(lang.ListValue)
	if !ok {
		return nil, fmt.Errorf("prepend: expected list, got %T", args[0])
	}

	// Prepend all values
	var result lang.ListValue
	for i := 1; i < len(args); i++ {
		result = append(result, args[i])
	}
	result = append(result, list...)

	return result, nil
}

func listInsert(args []lang.Value) (lang.Value, error) {
	if len(args) != 3 {
		return nil, errors.New("list_insert: expected 3 arguments (list, index, value)")
	}
	list, ok := args[0].(lang.ListValue)
	if !ok {
		return nil, fmt.Errorf("list_insert: expected list, got %T", args[0])
	}

	indexNum, err := lib.ToNumber(args[1])
	if err != nil {
		return nil, fmt.Errorf("list_insert: index %w", err)
	}
	index := int(indexNum)
	value := args[2]

	// Handle negative indices
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

func listRemove(args []lang.Value) (lang.Value, error) {
	if len(args) != 2 {
		return nil, errors.New("list_remove: expected 2 arguments (list, index)")
	}
	list, ok := args[0].(lang.ListValue)
	if !ok {
		return nil, fmt.Errorf("list_remove: expected list, got %T", args[0])
	}

	indexNum, err := lib.ToNumber(args[1])
	if err != nil {
		return nil, fmt.Errorf("list_remove: index %w", err)
	}
	index := int(indexNum)

	// Handle negative indices
	if index < 0 {
		index = len(list) + index
	}

	if index < 0 || index >= len(list) {
		return nil, fmt.Errorf("list_remove: index %d out of bounds (list length: %d)", index, len(list))
	}

	result := make(lang.ListValue, 0, len(list)-1)
	result = append(result, list[:index]...)
	result = append(result, list[index+1:]...)

	return result, nil
}

func listConcat(args []lang.Value) (lang.Value, error) {
	var result lang.ListValue

	for i, arg := range args {
		if list, ok := arg.(lang.ListValue); ok {
			result = append(result, list...)
		} else {
			return nil, fmt.Errorf("list_concat: argument %d expected list, got %T", i+1, arg)
		}
	}

	return result, nil
}

// List Access
func listFirst(args []lang.Value) (lang.Value, error) {
	if len(args) < 1 || len(args) > 2 {
		return nil, errors.New("first: expected 1 or 2 arguments (list, default?)")
	}
	list, ok := args[0].(lang.ListValue)
	if !ok {
		return nil, fmt.Errorf("first: expected list, got %T", args[0])
	}
	if len(list) == 0 {
		if len(args) == 2 {
			return args[1], nil // Default value
		}
		return nil, errors.New("first: list is empty")
	}
	return list[0], nil
}

func listLast(args []lang.Value) (lang.Value, error) {
	if len(args) < 1 || len(args) > 2 {
		return nil, errors.New("last: expected 1 or 2 arguments (list, default?)")
	}
	list, ok := args[0].(lang.ListValue)
	if !ok {
		return nil, fmt.Errorf("last: expected list, got %T", args[0])
	}
	if len(list) == 0 {
		if len(args) == 2 {
			return args[1], nil // Default value
		}
		return nil, errors.New("last: list is empty")
	}
	return list[len(list)-1], nil
}

func listHead(args []lang.Value) (lang.Value, error) {
	return listFirst(args)
}

func listTail(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("tail: expected 1 argument")
	}
	list, ok := args[0].(lang.ListValue)
	if !ok {
		return nil, fmt.Errorf("tail: expected list, got %T", args[0])
	}
	if len(list) <= 1 {
		return lang.ListValue{}, nil
	}

	result := make(lang.ListValue, len(list)-1)
	copy(result, list[1:])
	return result, nil
}

func listRest(args []lang.Value) (lang.Value, error) {
	return listTail(args)
}

func listInit(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("list_init: expected 1 argument")
	}
	list, ok := args[0].(lang.ListValue)
	if !ok {
		return nil, fmt.Errorf("list_init: expected list, got %T", args[0])
	}
	if len(list) <= 1 {
		return lang.ListValue{}, nil
	}

	result := make(lang.ListValue, len(list)-1)
	copy(result, list[:len(list)-1])
	return result, nil
}

func listSlice(args []lang.Value) (lang.Value, error) {
	if len(args) < 2 || len(args) > 3 {
		return nil, errors.New("slice: expected 2 or 3 arguments (list, start, end?)")
	}
	list, ok := args[0].(lang.ListValue)
	if !ok {
		return nil, fmt.Errorf("slice: expected list, got %T", args[0])
	}

	startNum, err := lib.ToNumber(args[1])
	if err != nil {
		return nil, fmt.Errorf("slice: start %w", err)
	}
	start := int(startNum)
	end := len(list)

	if len(args) == 3 {
		endNum, err := lib.ToNumber(args[2])
		if err != nil {
			return nil, fmt.Errorf("slice: end %w", err)
		}
		end = int(endNum)
	}

	// Handle negative indices
	if start < 0 {
		start = len(list) + start
	}
	if end < 0 {
		end = len(list) + end
	}

	// Clamp to bounds
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

func listTake(args []lang.Value) (lang.Value, error) {
	if len(args) != 2 {
		return nil, errors.New("take: expected 2 arguments (list, count)")
	}
	list, ok := args[0].(lang.ListValue)
	if !ok {
		return nil, fmt.Errorf("take: expected list, got %T", args[0])
	}

	countNum, err := lib.ToNumber(args[1])
	if err != nil {
		return nil, fmt.Errorf("take: count %w", err)
	}
	count := int(countNum)

	if count < 0 {
		return nil, errors.New("take: count must be non-negative")
	}
	if count == 0 {
		return lang.ListValue{}, nil
	}
	if count >= len(list) {
		// Return copy of entire list
		result := make(lang.ListValue, len(list))
		copy(result, list)
		return result, nil
	}

	result := make(lang.ListValue, count)
	copy(result, list[:count])
	return result, nil
}

func listDrop(args []lang.Value) (lang.Value, error) {
	if len(args) != 2 {
		return nil, errors.New("drop: expected 2 arguments (list, count)")
	}
	list, ok := args[0].(lang.ListValue)
	if !ok {
		return nil, fmt.Errorf("drop: expected list, got %T", args[0])
	}

	countNum, err := lib.ToNumber(args[1])
	if err != nil {
		return nil, fmt.Errorf("drop: count %w", err)
	}
	count := int(countNum)

	if count < 0 {
		return nil, errors.New("drop: count must be non-negative")
	}
	if count == 0 {
		// Return copy of entire list
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

// List Transformation
func listReverse(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("reverse: expected 1 argument")
	}
	list, ok := args[0].(lang.ListValue)
	if !ok {
		return nil, fmt.Errorf("reverse: expected list, got %T", args[0])
	}

	result := make(lang.ListValue, len(list))
	for i, v := range list {
		result[len(list)-1-i] = v
	}
	return result, nil
}

func listSort(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("sort: expected 1 argument")
	}
	list, ok := args[0].(lang.ListValue)
	if !ok {
		return nil, fmt.Errorf("sort: expected list, got %T", args[0])
	}

	result := make(lang.ListValue, len(list))
	copy(result, list)

	sort.Slice(result, func(i, j int) bool {
		return compareValues(result[i], result[j]) < 0
	})

	return result, nil
}

func listSortDesc(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("sort_desc: expected 1 argument")
	}
	list, ok := args[0].(lang.ListValue)
	if !ok {
		return nil, fmt.Errorf("sort_desc: expected list, got %T", args[0])
	}

	result := make(lang.ListValue, len(list))
	copy(result, list)

	sort.Slice(result, func(i, j int) bool {
		return compareValues(result[i], result[j]) > 0
	})

	return result, nil
}

func listShuffle(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("shuffle: expected 1 argument")
	}
	list, ok := args[0].(lang.ListValue)
	if !ok {
		return nil, fmt.Errorf("shuffle: expected list, got %T", args[0])
	}

	result := make(lang.ListValue, len(list))
	copy(result, list)

	for i := len(result) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		result[i], result[j] = result[j], result[i]
	}

	return result, nil
}

func listUnique(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("unique: expected 1 argument")
	}
	list, ok := args[0].(lang.ListValue)
	if !ok {
		return nil, fmt.Errorf("unique: expected list, got %T", args[0])
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

func listFlatten(args []lang.Value) (lang.Value, error) {
	if len(args) < 1 || len(args) > 2 {
		return nil, errors.New("flatten: expected 1 or 2 arguments (list, depth?)")
	}
	list, ok := args[0].(lang.ListValue)
	if !ok {
		return nil, fmt.Errorf("flatten: expected list, got %T", args[0])
	}

	depth := 1
	if len(args) == 2 {
		depthNum, err := lib.ToNumber(args[1])
		if err != nil {
			return nil, fmt.Errorf("flatten: depth %w", err)
		}
		depth = int(depthNum)
		if depth < 0 {
			return nil, errors.New("flatten: depth must be non-negative")
		}
	}

	return flattenList(list, depth), nil
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

// List Search and Testing
func listContains(args []lang.Value) (lang.Value, error) {
	if len(args) != 2 {
		return nil, errors.New("list_contains: expected 2 arguments (list, value)")
	}
	list, ok := args[0].(lang.ListValue)
	if !ok {
		return nil, fmt.Errorf("list_contains: expected list, got %T", args[0])
	}

	searchValue := args[1]
	for _, item := range list {
		if equalValues(item, searchValue) {
			return lang.BoolValue(true), nil
		}
	}

	return lang.BoolValue(false), nil
}

func listIndexOf(args []lang.Value) (lang.Value, error) {
	if len(args) < 2 || len(args) > 3 {
		return nil, errors.New("list_index_of: expected 2 or 3 arguments (list, value, start?)")
	}
	list, ok := args[0].(lang.ListValue)
	if !ok {
		return nil, fmt.Errorf("list_index_of: expected list, got %T", args[0])
	}

	searchValue := args[1]
	start := 0

	if len(args) == 3 {
		startNum, err := lib.ToNumber(args[2])
		if err != nil {
			return nil, fmt.Errorf("list_index_of: start %w", err)
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

func listLastIndexOf(args []lang.Value) (lang.Value, error) {
	if len(args) != 2 {
		return nil, errors.New("list_last_index_of: expected 2 arguments (list, value)")
	}
	list, ok := args[0].(lang.ListValue)
	if !ok {
		return nil, fmt.Errorf("list_last_index_of: expected list, got %T", args[0])
	}

	searchValue := args[1]
	for i := len(list) - 1; i >= 0; i-- {
		if equalValues(list[i], searchValue) {
			return lang.NumberValue(float64(i)), nil
		}
	}

	return lang.NumberValue(-1), nil
}

func listCount(args []lang.Value) (lang.Value, error) {
	if len(args) != 2 {
		return nil, errors.New("list_count: expected 2 arguments (list, value)")
	}
	list, ok := args[0].(lang.ListValue)
	if !ok {
		return nil, fmt.Errorf("list_count: expected list, got %T", args[0])
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

// List Generation
func listRange(args []lang.Value) (lang.Value, error) {
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

func listRepeat(args []lang.Value) (lang.Value, error) {
	if len(args) != 2 {
		return nil, errors.New("list_repeat: expected 2 arguments (value, count)")
	}

	value := args[0]
	countNum, err := lib.ToNumber(args[1])
	if err != nil {
		return nil, fmt.Errorf("list_repeat: count %w", err)
	}
	count := int(countNum)

	if count < 0 {
		return nil, errors.New("list_repeat: count must be non-negative")
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

func listZip(args []lang.Value) (lang.Value, error) {
	if len(args) == 0 {
		return nil, errors.New("zip: expected at least 1 argument")
	}

	// Convert all args to lists
	var lists []lang.ListValue
	minLen := -1

	for i, arg := range args {
		if list, ok := arg.(lang.ListValue); ok {
			lists = append(lists, list)
			if minLen == -1 || len(list) < minLen {
				minLen = len(list)
			}
		} else {
			return nil, fmt.Errorf("zip: argument %d expected list, got %T", i+1, arg)
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

// List Filtering and Mapping (simplified versions without function arguments)
func listFilter(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("filter: expected 1 argument")
	}
	list, ok := args[0].(lang.ListValue)
	if !ok {
		return nil, fmt.Errorf("filter: expected list, got %T", args[0])
	}

	var result lang.ListValue
	for _, item := range list {
		if item != nil && !isNullValue(item) {
			result = append(result, item)
		}
	}

	return result, nil
}

func listMap(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("map: expected 1 argument")
	}
	list, ok := args[0].(lang.ListValue)
	if !ok {
		return nil, fmt.Errorf("map: expected list, got %T", args[0])
	}

	// For now, just return a copy of the original list (would need function support)
	result := make(lang.ListValue, len(list))
	copy(result, list)
	return result, nil
}

// Helper functions
func compareValues(a, b lang.Value) int {
	aNum, aErr := lib.ToNumber(a)
	bNum, bErr := lib.ToNumber(b)

	// If both can be converted to numbers, compare numerically
	if aErr == nil && bErr == nil {
		if aNum < bNum {
			return -1
		}
		if aNum > bNum {
			return 1
		}
		return 0
	}

	// Fallback to string comparison
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

// Functions that would be in the BuiltinFunctions map:
var ListFunctions = map[string]func([]lang.Value) (lang.Value, error){
	// Basic operations
	"list_length":   listLength,
	"list_is_empty": listIsEmpty,
	"list_get":      listGet,
	"list_set":      listSet,

	// Construction
	"append":      listAppend,
	"prepend":     listPrepend,
	"list_insert": listInsert,
	"list_remove": listRemove,
	"list_concat": listConcat,

	// Access
	"first":     listFirst,
	"last":      listLast,
	"head":      listHead,
	"tail":      listTail,
	"rest":      listRest,
	"list_init": listInit,
	"slice":     listSlice,
	"take":      listTake,
	"drop":      listDrop,

	// Transformation
	"reverse":   listReverse,
	"sort":      listSort,
	"sort_desc": listSortDesc,
	"shuffle":   listShuffle,
	"unique":    listUnique,
	"flatten":   listFlatten,

	// Search and testing
	"list_contains":      listContains,
	"list_index_of":      listIndexOf,
	"list_last_index_of": listLastIndexOf,
	"list_count":         listCount,

	// Generation
	"range":       listRange,
	"list_repeat": listRepeat,
	"zip":         listZip,

	// Filtering and mapping (simplified)
	"filter": listFilter,
	"map":    listMap,
}
