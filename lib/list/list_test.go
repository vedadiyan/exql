package list

import (
	"math/rand"
	"testing"
	"time"

	"github.com/vedadiyan/exql/lang"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestLength(t *testing.T) {
	_, fn := length()

	tests := []struct {
		name     string
		input    lang.Value
		expected float64
		hasError bool
	}{
		{"empty list", lang.ListValue{}, 0, false},
		{"single item", lang.ListValue{lang.NumberValue(1)}, 1, false},
		{"multiple items", lang.ListValue{lang.NumberValue(1), lang.StringValue("test"), lang.BoolValue(true)}, 3, false},
		{"non-list input", lang.StringValue("test"), 0, true},
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
			if float64(result.(lang.NumberValue)) != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, float64(result.(lang.NumberValue)))
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
		hasError bool
	}{
		{"empty list", lang.ListValue{}, true, false},
		{"non-empty list", lang.ListValue{lang.NumberValue(1)}, false, false},
		{"non-list input", lang.StringValue("test"), false, true},
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
			if bool(result.(lang.BoolValue)) != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, bool(result.(lang.BoolValue)))
			}
		})
	}
}

func TestGet(t *testing.T) {
	_, fn := get()
	testList := lang.ListValue{lang.NumberValue(10), lang.StringValue("test"), lang.BoolValue(true)}

	tests := []struct {
		name     string
		list     lang.ListValue
		index    float64
		def      lang.Value
		expected lang.Value
		hasError bool
	}{
		{"valid index", testList, 0, nil, lang.NumberValue(10), false},
		{"negative index", testList, -1, nil, lang.BoolValue(true), false},
		{"out of bounds with default", testList, 5, lang.StringValue("default"), lang.StringValue("default"), false},
		{"out of bounds without default", testList, 5, nil, nil, true},
		{"negative out of bounds", testList, -5, nil, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var args []lang.Value
			args = []lang.Value{tt.list, lang.NumberValue(tt.index)}
			if tt.def != nil {
				args = append(args, tt.def)
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
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestSet(t *testing.T) {
	_, fn := set()
	testList := lang.ListValue{lang.NumberValue(10), lang.StringValue("test"), lang.BoolValue(true)}

	tests := []struct {
		name     string
		index    float64
		value    lang.Value
		hasError bool
	}{
		{"valid index", 1, lang.StringValue("modified"), false},
		{"negative index", -1, lang.BoolValue(false), false},
		{"out of bounds", 5, lang.StringValue("new"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{testList, lang.NumberValue(tt.index), tt.value})
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
			if len(resultList) != len(testList) {
				t.Errorf("Result list length changed unexpectedly")
			}
		})
	}
}

func TestAppend(t *testing.T) {
	_, fn := aappend()
	testList := lang.ListValue{lang.NumberValue(1), lang.NumberValue(2)}

	tests := []struct {
		name        string
		values      []lang.Value
		expectedLen int
		hasError    bool
	}{
		{"single value", []lang.Value{lang.NumberValue(3)}, 3, false},
		{"multiple values", []lang.Value{lang.NumberValue(3), lang.StringValue("test")}, 4, false},
		{"no values", []lang.Value{}, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []lang.Value{testList}
			args = append(args, tt.values...)

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

func TestPrepend(t *testing.T) {
	_, fn := prepend()
	testList := lang.ListValue{lang.NumberValue(1), lang.NumberValue(2)}

	result, err := fn([]lang.Value{testList, lang.StringValue("first"), lang.NumberValue(0)})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}
	resultList := result.(lang.ListValue)
	if len(resultList) != 4 {
		t.Errorf("Expected length 4, got %d", len(resultList))
	}
	if resultList[0] != lang.StringValue("first") {
		t.Errorf("Expected first element to be 'first', got %v", resultList[0])
	}
}

func TestInsert(t *testing.T) {
	_, fn := insert()
	testList := lang.ListValue{lang.NumberValue(1), lang.NumberValue(3)}

	tests := []struct {
		name     string
		index    float64
		value    lang.Value
		hasError bool
	}{
		{"middle insert", 1, lang.NumberValue(2), false},
		{"beginning insert", 0, lang.NumberValue(0), false},
		{"end insert", 2, lang.NumberValue(4), false},
		{"negative index", -1, lang.NumberValue(-1), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{testList, lang.NumberValue(tt.index), tt.value})
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
			if len(resultList) != len(testList)+1 {
				t.Errorf("Expected length %d, got %d", len(testList)+1, len(resultList))
			}
		})
	}
}

func TestRemove(t *testing.T) {
	_, fn := remove()
	testList := lang.ListValue{lang.NumberValue(1), lang.NumberValue(2), lang.NumberValue(3)}

	tests := []struct {
		name     string
		index    float64
		hasError bool
	}{
		{"valid index", 1, false},
		{"first element", 0, false},
		{"last element", 2, false},
		{"negative index", -1, false},
		{"out of bounds", 5, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{testList, lang.NumberValue(tt.index)})
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
			if len(resultList) != len(testList)-1 {
				t.Errorf("Expected length %d, got %d", len(testList)-1, len(resultList))
			}
		})
	}
}

func TestConcat(t *testing.T) {
	_, fn := concat()

	tests := []struct {
		name        string
		lists       []lang.ListValue
		expectedLen int
		hasError    bool
	}{
		{
			"two lists",
			[]lang.ListValue{
				{lang.NumberValue(1), lang.NumberValue(2)},
				{lang.NumberValue(3), lang.NumberValue(4)},
			},
			4,
			false,
		},
		{
			"empty lists",
			[]lang.ListValue{{}, {}},
			0,
			false,
		},
		{
			"single list",
			[]lang.ListValue{{lang.NumberValue(1)}},
			1,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := make([]lang.Value, len(tt.lists))
			for i, list := range tt.lists {
				args[i] = list
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

func TestFirst(t *testing.T) {
	_, fn := first()

	tests := []struct {
		name     string
		list     lang.ListValue
		def      lang.Value
		expected lang.Value
		hasError bool
	}{
		{"non-empty list", lang.ListValue{lang.NumberValue(10), lang.NumberValue(20)}, nil, lang.NumberValue(10), false},
		{"empty list with default", lang.ListValue{}, lang.StringValue("default"), lang.StringValue("default"), false},
		{"empty list no default", lang.ListValue{}, nil, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var args []lang.Value
			args = []lang.Value{tt.list}
			if tt.def != nil {
				args = append(args, tt.def)
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
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestLast(t *testing.T) {
	_, fn := last()

	tests := []struct {
		name     string
		list     lang.ListValue
		def      lang.Value
		expected lang.Value
		hasError bool
	}{
		{"non-empty list", lang.ListValue{lang.NumberValue(10), lang.NumberValue(20)}, nil, lang.NumberValue(20), false},
		{"empty list with default", lang.ListValue{}, lang.StringValue("default"), lang.StringValue("default"), false},
		{"empty list no default", lang.ListValue{}, nil, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var args []lang.Value
			args = []lang.Value{tt.list}
			if tt.def != nil {
				args = append(args, tt.def)
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
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestTail(t *testing.T) {
	_, fn := tail()

	tests := []struct {
		name        string
		list        lang.ListValue
		expectedLen int
	}{
		{"multiple items", lang.ListValue{lang.NumberValue(1), lang.NumberValue(2), lang.NumberValue(3)}, 2},
		{"single item", lang.ListValue{lang.NumberValue(1)}, 0},
		{"empty list", lang.ListValue{}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{tt.list})
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

func TestInit(t *testing.T) {
	_, fn := iinit()

	tests := []struct {
		name        string
		list        lang.ListValue
		expectedLen int
	}{
		{"multiple items", lang.ListValue{lang.NumberValue(1), lang.NumberValue(2), lang.NumberValue(3)}, 2},
		{"single item", lang.ListValue{lang.NumberValue(1)}, 0},
		{"empty list", lang.ListValue{}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{tt.list})
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

func TestSlice(t *testing.T) {
	_, fn := slice()
	testList := lang.ListValue{lang.NumberValue(0), lang.NumberValue(1), lang.NumberValue(2), lang.NumberValue(3), lang.NumberValue(4)}

	tests := []struct {
		name        string
		start       float64
		end         *float64
		expectedLen int
		hasError    bool
	}{
		{"normal slice", 1, float64Ptr(3), 2, false},
		{"slice to end", 2, nil, 3, false},
		{"negative start", -3, nil, 3, false},
		{"negative end", 1, float64Ptr(-1), 3, false},
		{"start > end", 3, float64Ptr(1), 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var args []lang.Value
			args = []lang.Value{testList, lang.NumberValue(tt.start)}
			if tt.end != nil {
				args = append(args, lang.NumberValue(*tt.end))
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

func TestTake(t *testing.T) {
	_, fn := take()
	testList := lang.ListValue{lang.NumberValue(1), lang.NumberValue(2), lang.NumberValue(3), lang.NumberValue(4)}

	tests := []struct {
		name        string
		count       float64
		expectedLen int
		hasError    bool
	}{
		{"take some", 2, 2, false},
		{"take all", 4, 4, false},
		{"take more than available", 10, 4, false},
		{"take zero", 0, 0, false},
		{"negative count", -1, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{testList, lang.NumberValue(tt.count)})
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

func TestDrop(t *testing.T) {
	_, fn := drop()
	testList := lang.ListValue{lang.NumberValue(1), lang.NumberValue(2), lang.NumberValue(3), lang.NumberValue(4)}

	tests := []struct {
		name        string
		count       float64
		expectedLen int
		hasError    bool
	}{
		{"drop some", 2, 2, false},
		{"drop all", 4, 0, false},
		{"drop more than available", 10, 0, false},
		{"drop zero", 0, 4, false},
		{"negative count", -1, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{testList, lang.NumberValue(tt.count)})
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

func TestReverse(t *testing.T) {
	_, fn := reverse()
	testList := lang.ListValue{lang.NumberValue(1), lang.NumberValue(2), lang.NumberValue(3)}

	result, err := fn([]lang.Value{testList})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}
	resultList := result.(lang.ListValue)
	if len(resultList) != len(testList) {
		t.Errorf("Expected length %d, got %d", len(testList), len(resultList))
	}
	if resultList[0] != lang.NumberValue(3) || resultList[2] != lang.NumberValue(1) {
		t.Errorf("List not properly reversed")
	}
}

func TestSort(t *testing.T) {
	_, fn := ssort()
	testList := lang.ListValue{lang.NumberValue(3), lang.NumberValue(1), lang.NumberValue(2)}

	result, err := fn([]lang.Value{testList})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}
	resultList := result.(lang.ListValue)
	if len(resultList) != len(testList) {
		t.Errorf("Expected length %d, got %d", len(testList), len(resultList))
	}
	if resultList[0] != lang.NumberValue(1) || resultList[2] != lang.NumberValue(3) {
		t.Errorf("List not properly sorted")
	}
}

func TestSortDesc(t *testing.T) {
	_, fn := sortDesc()
	testList := lang.ListValue{lang.NumberValue(1), lang.NumberValue(3), lang.NumberValue(2)}

	result, err := fn([]lang.Value{testList})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}
	resultList := result.(lang.ListValue)
	if len(resultList) != len(testList) {
		t.Errorf("Expected length %d, got %d", len(testList), len(resultList))
	}
	if resultList[0] != lang.NumberValue(3) || resultList[2] != lang.NumberValue(1) {
		t.Errorf("List not properly sorted descending")
	}
}

func TestShuffle(t *testing.T) {
	_, fn := shuffle()
	testList := lang.ListValue{lang.NumberValue(1), lang.NumberValue(2), lang.NumberValue(3), lang.NumberValue(4), lang.NumberValue(5)}

	result, err := fn([]lang.Value{testList})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}
	resultList := result.(lang.ListValue)
	if len(resultList) != len(testList) {
		t.Errorf("Expected length %d, got %d", len(testList), len(resultList))
	}
}

func TestUnique(t *testing.T) {
	_, fn := unique()
	testList := lang.ListValue{lang.NumberValue(1), lang.NumberValue(2), lang.NumberValue(1), lang.NumberValue(3), lang.NumberValue(2)}

	result, err := fn([]lang.Value{testList})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}
	resultList := result.(lang.ListValue)
	if len(resultList) != 3 {
		t.Errorf("Expected length 3, got %d", len(resultList))
	}
}

func TestFlatten(t *testing.T) {
	_, fn := flatten()

	tests := []struct {
		name        string
		list        lang.ListValue
		depth       *float64
		expectedLen int
		hasError    bool
	}{
		{
			"simple flatten",
			lang.ListValue{lang.NumberValue(1), lang.ListValue{lang.NumberValue(2), lang.NumberValue(3)}, lang.NumberValue(4)},
			nil,
			4,
			false,
		},
		{
			"depth 0",
			lang.ListValue{lang.NumberValue(1), lang.ListValue{lang.NumberValue(2)}},
			float64Ptr(0),
			2,
			false,
		},
		{
			"negative depth",
			lang.ListValue{lang.NumberValue(1)},
			float64Ptr(-1),
			0,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var args []lang.Value
			args = []lang.Value{tt.list}
			if tt.depth != nil {
				args = append(args, lang.NumberValue(*tt.depth))
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

func TestContains(t *testing.T) {
	_, fn := contains()
	testList := lang.ListValue{lang.NumberValue(1), lang.StringValue("test"), lang.BoolValue(true)}

	tests := []struct {
		name     string
		value    lang.Value
		expected bool
	}{
		{"contains number", lang.NumberValue(1), true},
		{"contains string", lang.StringValue("test"), true},
		{"contains bool", lang.BoolValue(true), true},
		{"not contains", lang.NumberValue(99), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{testList, tt.value})
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
	testList := lang.ListValue{lang.NumberValue(1), lang.StringValue("test"), lang.NumberValue(1), lang.BoolValue(true)}

	tests := []struct {
		name     string
		value    lang.Value
		start    *float64
		expected float64
	}{
		{"first occurrence", lang.NumberValue(1), nil, 0},
		{"string value", lang.StringValue("test"), nil, 1},
		{"with start index", lang.NumberValue(1), float64Ptr(1), 2},
		{"not found", lang.NumberValue(99), nil, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var args []lang.Value
			args = []lang.Value{testList, tt.value}
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
	testList := lang.ListValue{lang.NumberValue(1), lang.StringValue("test"), lang.NumberValue(1), lang.BoolValue(true)}

	tests := []struct {
		name     string
		value    lang.Value
		expected float64
	}{
		{"last occurrence", lang.NumberValue(1), 2},
		{"single occurrence", lang.StringValue("test"), 1},
		{"not found", lang.NumberValue(99), -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{testList, tt.value})
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

func TestCount(t *testing.T) {
	_, fn := count()
	testList := lang.ListValue{lang.NumberValue(1), lang.NumberValue(1), lang.StringValue("test"), lang.NumberValue(1)}

	tests := []struct {
		name     string
		value    lang.Value
		expected float64
	}{
		{"multiple occurrences", lang.NumberValue(1), 3},
		{"single occurrence", lang.StringValue("test"), 1},
		{"no occurrences", lang.NumberValue(99), 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{testList, tt.value})
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

func TestRange(t *testing.T) {
	_, fn := rrange()

	tests := []struct {
		name        string
		args        []float64
		expectedLen int
		hasError    bool
	}{
		{"single arg", []float64{5}, 5, false},
		{"start and end", []float64{2, 5}, 3, false},
		{"with step", []float64{0, 10, 2}, 5, false},
		{"negative step", []float64{10, 0, -2}, 5, false},
		{"zero step", []float64{0, 5, 0}, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := make([]lang.Value, len(tt.args))
			for i, arg := range tt.args {
				args[i] = lang.NumberValue(arg)
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

func TestRepeat(t *testing.T) {
	_, fn := repeat()

	tests := []struct {
		name        string
		value       lang.Value
		count       float64
		expectedLen int
		hasError    bool
	}{
		{"repeat number", lang.NumberValue(42), 3, 3, false},
		{"repeat string", lang.StringValue("test"), 2, 2, false},
		{"repeat zero", lang.NumberValue(1), 0, 0, false},
		{"negative count", lang.NumberValue(1), -1, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{tt.value, lang.NumberValue(tt.count)})
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
			if tt.expectedLen > 0 && resultList[0] != tt.value {
				t.Errorf("Expected all elements to be %v, got %v", tt.value, resultList[0])
			}
		})
	}
}

func TestZip(t *testing.T) {
	_, fn := zip()

	tests := []struct {
		name        string
		lists       []lang.ListValue
		expectedLen int
		hasError    bool
	}{
		{
			"two equal lists",
			[]lang.ListValue{
				{lang.NumberValue(1), lang.NumberValue(2)},
				{lang.StringValue("a"), lang.StringValue("b")},
			},
			2,
			false,
		},
		{
			"different length lists",
			[]lang.ListValue{
				{lang.NumberValue(1), lang.NumberValue(2), lang.NumberValue(3)},
				{lang.StringValue("a"), lang.StringValue("b")},
			},
			2,
			false,
		},
		{
			"empty lists",
			[]lang.ListValue{{}, {}},
			0,
			false,
		},
		{
			"no arguments",
			[]lang.ListValue{},
			0,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := make([]lang.Value, len(tt.lists))
			for i, list := range tt.lists {
				args[i] = list
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
			if tt.expectedLen > 0 {
				firstTuple, ok := resultList[0].(lang.ListValue)
				if !ok || len(firstTuple) != len(tt.lists) {
					t.Errorf("Expected tuples of length %d", len(tt.lists))
				}
			}
		})
	}
}

func TestFilter(t *testing.T) {
	_, fn := filter()
	testList := lang.ListValue{
		lang.NumberValue(1),
		nil,
		lang.StringValue("test"),
		lang.StringValue(""),
		lang.BoolValue(true),
		lang.BoolValue(false),
		lang.NumberValue(0),
	}

	result, err := fn([]lang.Value{testList})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}
	resultList := result.(lang.ListValue)
	if len(resultList) >= len(testList) {
		t.Errorf("Expected filtered list to be shorter than original")
	}
}

func TestMap(t *testing.T) {
	_, fn := mmap()
	testList := lang.ListValue{lang.NumberValue(1), lang.NumberValue(2), lang.NumberValue(3)}

	result, err := fn([]lang.Value{testList})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}
	resultList := result.(lang.ListValue)
	if len(resultList) != len(testList) {
		t.Errorf("Expected same length, got %d vs %d", len(resultList), len(testList))
	}
}

func TestHelperFunctions(t *testing.T) {
	t.Run("compareValues", func(t *testing.T) {
		tests := []struct {
			a        lang.Value
			b        lang.Value
			expected int
		}{
			{lang.NumberValue(1), lang.NumberValue(2), -1},
			{lang.NumberValue(2), lang.NumberValue(1), 1},
			{lang.NumberValue(1), lang.NumberValue(1), 0},
			{lang.StringValue("a"), lang.StringValue("b"), -1},
			{lang.StringValue("b"), lang.StringValue("a"), 1},
			{lang.StringValue("a"), lang.StringValue("a"), 0},
		}

		for _, tt := range tests {
			result := compareValues(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("compareValues(%v, %v) = %d, expected %d", tt.a, tt.b, result, tt.expected)
			}
		}
	})

	t.Run("equalValues", func(t *testing.T) {
		tests := []struct {
			a        lang.Value
			b        lang.Value
			expected bool
		}{
			{lang.NumberValue(1), lang.NumberValue(1), true},
			{lang.NumberValue(1), lang.NumberValue(2), false},
			{lang.StringValue("test"), lang.StringValue("test"), true},
			{lang.StringValue("test"), lang.StringValue("other"), false},
			{lang.BoolValue(true), lang.BoolValue(true), true},
			{lang.BoolValue(true), lang.BoolValue(false), false},
		}

		for _, tt := range tests {
			result := equalValues(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("equalValues(%v, %v) = %v, expected %v", tt.a, tt.b, result, tt.expected)
			}
		}
	})

	t.Run("valueToString", func(t *testing.T) {
		tests := []struct {
			value    lang.Value
			expected string
		}{
			{lang.StringValue("test"), "test"},
			{lang.NumberValue(42), "42"},
			{lang.BoolValue(true), "true"},
			{lang.BoolValue(false), "false"},
			{nil, ""},
		}

		for _, tt := range tests {
			result := valueToString(tt.value)
			if result != tt.expected {
				t.Errorf("valueToString(%v) = %s, expected %s", tt.value, result, tt.expected)
			}
		}
	})

	t.Run("isNullValue", func(t *testing.T) {
		tests := []struct {
			value    lang.Value
			expected bool
		}{
			{lang.StringValue(""), true},
			{lang.StringValue("test"), false},
			{lang.NumberValue(0), true},
			{lang.NumberValue(42), false},
			{lang.BoolValue(false), true},
			{lang.BoolValue(true), false},
			{nil, true},
		}

		for _, tt := range tests {
			result := isNullValue(tt.value)
			if result != tt.expected {
				t.Errorf("isNullValue(%v) = %v, expected %v", tt.value, result, tt.expected)
			}
		}
	})
}

func TestArgumentErrors(t *testing.T) {
	functions := []struct {
		fn           func() (string, lang.Function)
		expectedArgs int
	}{
		{length, 1},
		{isEmpty, 1},
		{set, 3},
		{insert, 3},
		{remove, 2},
		{contains, 2},
		{indexOf, 2},
		{lastIndexOf, 2},
		{count, 2},
		{take, 2},
		{drop, 2},
		{reverse, 1},
		{ssort, 1},
		{sortDesc, 1},
		{shuffle, 1},
		{unique, 1},
		{tail, 1},
		{iinit, 1},
		{repeat, 2},
		{filter, 1},
		{mmap, 1},
	}

	for _, tf := range functions {
		name, fn := tf.fn()
		t.Run(name+"_wrong_args", func(t *testing.T) {
			args := make([]lang.Value, tf.expectedArgs-1)
			for i := range args {
				args[i] = lang.ListValue{lang.NumberValue(1)}
			}
			_, err := fn(args)
			if err == nil {
				t.Errorf("Expected error for wrong argument count in %s", name)
			}
		})
	}
}

func TestNonListInputs(t *testing.T) {
	functions := []func() (string, lang.Function){
		length, isEmpty, set, aappend, prepend, insert, remove,
		first, last, head, tail, rest, iinit, take, drop,
		reverse, ssort, sortDesc, shuffle, unique, flatten,
		contains, indexOf, lastIndexOf, count, filter, mmap,
	}

	nonListInput := lang.StringValue("not a list")

	for _, getFn := range functions {
		name, fn := getFn()
		t.Run(name+"_non_list", func(t *testing.T) {
			args := []lang.Value{nonListInput}
			if name == "set" || name == "insert" {
				args = append(args, lang.NumberValue(0), lang.StringValue("value"))
			} else if name == "append" || name == "prepend" {
				args = append(args, lang.StringValue("value"))
			} else if name == "contains" || name == "indexOf" || name == "lastIndexOf" || name == "count" {
				args = append(args, lang.StringValue("value"))
			} else if name == "take" || name == "drop" {
				args = append(args, lang.NumberValue(1))
			} else if name == "remove" {
				args = append(args, lang.NumberValue(0))
			}

			_, err := fn(args)
			if err == nil {
				t.Errorf("Expected error for non-list input in %s", name)
			}
		})
	}
}

func TestExport(t *testing.T) {
	functions := Export()

	expectedFunctions := []string{
		"length", "isEmpty", "get", "set", "append", "prepend", "insert",
		"remove", "concat", "first", "last", "head", "tail", "rest", "init",
		"slice", "take", "drop", "reverse", "sort", "sortDesc", "shuffle",
		"unique", "flatten", "contains", "indexOf", "lastIndexOf", "count",
		"range", "repeat", "zip", "filter", "map",
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
	testList := lang.ListValue{lang.NumberValue(1), lang.NumberValue(2), lang.NumberValue(3)}
	args := []lang.Value{testList}

	for i := 0; i < b.N; i++ {
		fn(args)
	}
}

func BenchmarkAppend(b *testing.B) {
	_, fn := aappend()
	testList := lang.ListValue{lang.NumberValue(1), lang.NumberValue(2)}
	args := []lang.Value{testList, lang.NumberValue(3)}

	for i := 0; i < b.N; i++ {
		fn(args)
	}
}

func BenchmarkSort(b *testing.B) {
	_, fn := ssort()
	testList := lang.ListValue{lang.NumberValue(3), lang.NumberValue(1), lang.NumberValue(4), lang.NumberValue(2)}
	args := []lang.Value{testList}

	for i := 0; i < b.N; i++ {
		fn(args)
	}
}

func BenchmarkContains(b *testing.B) {
	_, fn := contains()
	testList := make(lang.ListValue, 1000)
	for i := range testList {
		testList[i] = lang.NumberValue(float64(i))
	}
	args := []lang.Value{testList, lang.NumberValue(500)}

	for i := 0; i < b.N; i++ {
		fn(args)
	}
}

func float64Ptr(f float64) *float64 {
	return &f
}
