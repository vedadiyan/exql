package maps

import (
	"reflect"
	"testing"

	"github.com/vedadiyan/exql/lang"
)

func createTestMap() lang.MapValue {
	return lang.MapValue{
		"name":    lang.StringValue("John"),
		"age":     lang.NumberValue(30),
		"active":  lang.BoolValue(true),
		"city":    lang.StringValue("New York"),
		"score":   lang.NumberValue(85.5),
		"enabled": lang.BoolValue(false),
	}
}

func TestKeys(t *testing.T) {
	_, fn := keys()
	testMap := createTestMap()

	result, err := fn([]lang.Value{testMap})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	keysList, ok := result.(lang.ListValue)
	if !ok {
		t.Errorf("Expected ListValue, got %T", result)
		return
	}

	if len(keysList) != len(testMap) {
		t.Errorf("Expected %d keys, got %d", len(testMap), len(keysList))
	}

	// Keys should be sorted
	expectedKeys := []string{"active", "age", "city", "enabled", "name", "score"}
	for i, keyVal := range keysList {
		if i < len(expectedKeys) {
			if string(keyVal.(lang.StringValue)) != expectedKeys[i] {
				t.Errorf("Key %d: expected %s, got %s", i, expectedKeys[i], string(keyVal.(lang.StringValue)))
			}
		}
	}
}

func TestValues(t *testing.T) {
	_, fn := values()
	testMap := createTestMap()

	result, err := fn([]lang.Value{testMap})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	valuesList, ok := result.(lang.ListValue)
	if !ok {
		t.Errorf("Expected ListValue, got %T", result)
		return
	}

	if len(valuesList) != len(testMap) {
		t.Errorf("Expected %d values, got %d", len(testMap), len(valuesList))
	}
}

func TestSize(t *testing.T) {
	_, fn := size()

	tests := []struct {
		name     string
		input    lang.Value
		expected float64
		hasError bool
	}{
		{"non-empty map", createTestMap(), 6, false},
		{"empty map", lang.MapValue{}, 0, false},
		{"non-map input", lang.StringValue("test"), 0, true},
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
		{"empty map", lang.MapValue{}, true, false},
		{"non-empty map", createTestMap(), false, false},
		{"non-map input", lang.StringValue("test"), false, true},
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

func TestHas(t *testing.T) {
	_, fn := has()
	testMap := createTestMap()

	tests := []struct {
		name     string
		key      string
		expected bool
		hasError bool
	}{
		{"existing key", "name", true, false},
		{"non-existing key", "nonexistent", false, false},
		{"numeric key conversion", "age", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{testMap, lang.StringValue(tt.key)})
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
	testMap := createTestMap()

	tests := []struct {
		name     string
		key      string
		def      lang.Value
		expected lang.Value
		hasError bool
	}{
		{"existing key", "name", nil, lang.StringValue("John"), false},
		{"non-existing key with default", "nonexistent", lang.StringValue("default"), lang.StringValue("default"), false},
		{"non-existing key without default", "nonexistent", nil, nil, false},
		{"numeric value", "age", nil, lang.NumberValue(30), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var args []lang.Value
			args = []lang.Value{testMap, lang.StringValue(tt.key)}
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
	testMap := createTestMap()

	tests := []struct {
		name     string
		key      string
		value    lang.Value
		hasError bool
	}{
		{"set existing key", "name", lang.StringValue("Jane"), false},
		{"set new key", "email", lang.StringValue("jane@example.com"), false},
		{"set numeric value", "salary", lang.NumberValue(50000), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{testMap, lang.StringValue(tt.key), tt.value})
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

			resultMap := result.(lang.MapValue)
			if resultMap[tt.key] != tt.value {
				t.Errorf("Expected key %s to have value %v, got %v", tt.key, tt.value, resultMap[tt.key])
			}

			// Original map should be unchanged
			if len(resultMap) <= len(testMap) && tt.key == "email" {
				t.Errorf("Expected new map to have more keys when adding new key")
			}
		})
	}
}

func TestRemove(t *testing.T) {
	_, fn := remove()
	testMap := createTestMap()
	originalSize := len(testMap)

	tests := []struct {
		name     string
		key      string
		hasError bool
	}{
		{"remove existing key", "name", false},
		{"remove non-existing key", "nonexistent", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{testMap, lang.StringValue(tt.key)})
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

			resultMap := result.(lang.MapValue)
			if _, exists := resultMap[tt.key]; exists && tt.key == "name" {
				t.Errorf("Key %s should have been removed", tt.key)
			}

			if tt.key == "name" && len(resultMap) != originalSize-1 {
				t.Errorf("Expected map size to decrease by 1, got size %d", len(resultMap))
			}
		})
	}
}

func TestMerge(t *testing.T) {
	_, fn := merge()

	map1 := lang.MapValue{
		"a": lang.NumberValue(1),
		"b": lang.StringValue("test"),
	}
	map2 := lang.MapValue{
		"b": lang.StringValue("overwritten"),
		"c": lang.BoolValue(true),
	}
	map3 := lang.MapValue{
		"d": lang.NumberValue(4),
	}

	tests := []struct {
		name        string
		maps        []lang.MapValue
		expectedLen int
		hasError    bool
	}{
		{"merge two maps", []lang.MapValue{map1, map2}, 3, false},
		{"merge three maps", []lang.MapValue{map1, map2, map3}, 4, false},
		{"single map", []lang.MapValue{map1}, 2, false},
		{"no arguments", []lang.MapValue{}, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := make([]lang.Value, len(tt.maps))
			for i, m := range tt.maps {
				args[i] = m
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

			resultMap := result.(lang.MapValue)
			if len(resultMap) != tt.expectedLen {
				t.Errorf("Expected length %d, got %d", tt.expectedLen, len(resultMap))
			}

			// Test specific merge behavior
			if tt.name == "merge two maps" {
				if resultMap["b"] != lang.StringValue("overwritten") {
					t.Errorf("Expected b to be overwritten")
				}
			}
		})
	}
}

func TestMergeDeep(t *testing.T) {
	_, fn := mergeDeep()

	nestedMap1 := lang.MapValue{
		"user": lang.MapValue{
			"name": lang.StringValue("John"),
			"age":  lang.NumberValue(30),
		},
		"config": lang.MapValue{
			"theme": lang.StringValue("dark"),
		},
	}

	nestedMap2 := lang.MapValue{
		"user": lang.MapValue{
			"email": lang.StringValue("john@example.com"),
			"age":   lang.NumberValue(31), // This should overwrite
		},
		"config": lang.MapValue{
			"language": lang.StringValue("en"),
		},
	}

	result, err := fn([]lang.Value{nestedMap1, nestedMap2})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	resultMap := result.(lang.MapValue)
	userMap := resultMap["user"].(lang.MapValue)

	// Check that both name and email exist
	if userMap["name"] != lang.StringValue("John") {
		t.Errorf("Expected name to be preserved")
	}
	if userMap["email"] != lang.StringValue("john@example.com") {
		t.Errorf("Expected email to be merged")
	}
	if userMap["age"] != lang.NumberValue(31) {
		t.Errorf("Expected age to be overwritten")
	}

	configMap := resultMap["config"].(lang.MapValue)
	if configMap["theme"] != lang.StringValue("dark") {
		t.Errorf("Expected theme to be preserved")
	}
	if configMap["language"] != lang.StringValue("en") {
		t.Errorf("Expected language to be merged")
	}
}

func TestInvert(t *testing.T) {
	_, fn := invert()
	testMap := lang.MapValue{
		"a": lang.StringValue("1"),
		"b": lang.StringValue("2"),
		"c": lang.NumberValue(3),
	}

	result, err := fn([]lang.Value{testMap})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	resultMap := result.(lang.MapValue)
	if len(resultMap) != len(testMap) {
		t.Errorf("Expected same length, got %d vs %d", len(resultMap), len(testMap))
	}

	if resultMap["1"] != lang.StringValue("a") {
		t.Errorf("Expected inverted mapping")
	}
}

func TestFilter(t *testing.T) {
	_, fn := filter()
	testMap := lang.MapValue{
		"name":     lang.StringValue("John"),
		"empty":    lang.StringValue(""),
		"age":      lang.NumberValue(30),
		"zero":     lang.NumberValue(0),
		"active":   lang.BoolValue(true),
		"inactive": lang.BoolValue(false),
		"nothing":  nil,
	}

	result, err := fn([]lang.Value{testMap})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	resultMap := result.(lang.MapValue)
	if len(resultMap) >= len(testMap) {
		t.Errorf("Expected filtered map to be smaller")
	}

	// Non-null values should remain
	if resultMap["name"] != lang.StringValue("John") {
		t.Errorf("Expected non-null values to remain")
	}
	if resultMap["age"] != lang.NumberValue(30) {
		t.Errorf("Expected non-zero numbers to remain")
	}
	if resultMap["active"] != lang.BoolValue(true) {
		t.Errorf("Expected true boolean to remain")
	}
}

func TestFilterKeys(t *testing.T) {
	_, fn := filterKeys()
	testMap := createTestMap()
	keysToKeep := lang.ListValue{lang.StringValue("name"), lang.StringValue("age")}

	result, err := fn([]lang.Value{testMap, keysToKeep})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	resultMap := result.(lang.MapValue)
	if len(resultMap) != 2 {
		t.Errorf("Expected 2 keys, got %d", len(resultMap))
	}

	if resultMap["name"] != lang.StringValue("John") {
		t.Errorf("Expected name to be kept")
	}
	if resultMap["age"] != lang.NumberValue(30) {
		t.Errorf("Expected age to be kept")
	}
	if _, exists := resultMap["city"]; exists {
		t.Errorf("Expected city to be filtered out")
	}
}

func TestOmitKeys(t *testing.T) {
	_, fn := omitKeys()
	testMap := createTestMap()
	keysToOmit := lang.ListValue{lang.StringValue("name"), lang.StringValue("age")}

	result, err := fn([]lang.Value{testMap, keysToOmit})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	resultMap := result.(lang.MapValue)
	expectedLen := len(testMap) - 2
	if len(resultMap) != expectedLen {
		t.Errorf("Expected %d keys, got %d", expectedLen, len(resultMap))
	}

	if _, exists := resultMap["name"]; exists {
		t.Errorf("Expected name to be omitted")
	}
	if _, exists := resultMap["age"]; exists {
		t.Errorf("Expected age to be omitted")
	}
	if resultMap["city"] != lang.StringValue("New York") {
		t.Errorf("Expected city to be kept")
	}
}

func TestRename(t *testing.T) {
	_, fn := rename()
	testMap := lang.MapValue{
		"name": lang.StringValue("John"),
		"age":  lang.NumberValue(30),
	}
	renameMap := lang.MapValue{
		"name": lang.StringValue("full_name"),
		"age":  lang.StringValue("years_old"),
	}

	result, err := fn([]lang.Value{testMap, renameMap})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	resultMap := result.(lang.MapValue)
	if len(resultMap) != len(testMap) {
		t.Errorf("Expected same length after rename")
	}

	if resultMap["full_name"] != lang.StringValue("John") {
		t.Errorf("Expected renamed key full_name")
	}
	if resultMap["years_old"] != lang.NumberValue(30) {
		t.Errorf("Expected renamed key years_old")
	}
	if _, exists := resultMap["name"]; exists {
		t.Errorf("Expected original key name to be removed")
	}
}

func TestToList(t *testing.T) {
	_, fn := toList()
	testMap := lang.MapValue{
		"a": lang.NumberValue(1),
		"b": lang.StringValue("test"),
	}

	result, err := fn([]lang.Value{testMap})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	resultList := result.(lang.ListValue)
	if len(resultList) != len(testMap) {
		t.Errorf("Expected %d pairs, got %d", len(testMap), len(resultList))
	}

	// Check first pair (should be sorted by key)
	firstPair := resultList[0].(lang.ListValue)
	if len(firstPair) != 2 {
		t.Errorf("Expected pair to have 2 elements")
	}
}

func TestFromList(t *testing.T) {
	_, fn := fromList()
	testList := lang.ListValue{
		lang.ListValue{lang.StringValue("name"), lang.StringValue("John")},
		lang.ListValue{lang.StringValue("age"), lang.NumberValue(30)},
	}

	result, err := fn([]lang.Value{testList})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	resultMap := result.(lang.MapValue)
	if len(resultMap) != 2 {
		t.Errorf("Expected 2 keys, got %d", len(resultMap))
	}

	if resultMap["name"] != lang.StringValue("John") {
		t.Errorf("Expected name to be John")
	}
	if resultMap["age"] != lang.NumberValue(30) {
		t.Errorf("Expected age to be 30")
	}
}

func TestToQueryString(t *testing.T) {
	_, fn := toQueryString()
	testMap := lang.MapValue{
		"name": lang.StringValue("John Doe"),
		"age":  lang.NumberValue(30),
		"tags": lang.ListValue{lang.StringValue("dev"), lang.StringValue("golang")},
	}

	result, err := fn([]lang.Value{testMap})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	queryStr := string(result.(lang.StringValue))
	if queryStr == "" {
		t.Errorf("Expected non-empty query string")
	}

	// Should contain the parameters (order may vary due to sorting)
	if !contains(queryStr, "age=30") {
		t.Errorf("Expected age parameter in query string")
	}
	if !contains(queryStr, "name=John Doe") {
		t.Errorf("Expected name parameter in query string")
	}
}

func TestFromQueryString(t *testing.T) {
	_, fn := fromQueryString()

	tests := []struct {
		name        string
		queryString string
		expectedLen int
	}{
		{"simple query", "name=John&age=30", 2},
		{"array values", "tag=a&tag=b&tag=c", 1}, // Should create array for duplicate keys
		{"empty string", "", 0},
		{"single param", "name=John", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.StringValue(tt.queryString)})
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			resultMap := result.(lang.MapValue)
			if len(resultMap) != tt.expectedLen {
				t.Errorf("Expected %d keys, got %d", tt.expectedLen, len(resultMap))
			}

			if tt.name == "array values" {
				tags := resultMap["tag"].(lang.ListValue)
				if len(tags) != 3 {
					t.Errorf("Expected 3 tag values, got %d", len(tags))
				}
			}
		})
	}
}

func TestGetPath(t *testing.T) {
	_, fn := getPath()
	nestedMap := lang.MapValue{
		"user": lang.MapValue{
			"profile": lang.MapValue{
				"name": lang.StringValue("John"),
			},
		},
		"settings": lang.MapValue{
			"theme": lang.StringValue("dark"),
		},
	}

	tests := []struct {
		name     string
		path     string
		def      lang.Value
		expected lang.Value
		hasError bool
	}{
		{"nested path", "user.profile.name", nil, lang.StringValue("John"), false},
		{"simple path", "settings.theme", nil, lang.StringValue("dark"), false},
		{"non-existent path with default", "user.email", lang.StringValue("default"), lang.StringValue("default"), false},
		{"non-existent path without default", "user.email", nil, nil, false},
		{"empty path", "", nil, nestedMap, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var args []lang.Value
			args = []lang.Value{nestedMap, lang.StringValue(tt.path)}
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
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestSetPath(t *testing.T) {
	_, fn := setPath()
	testMap := lang.MapValue{
		"user": lang.MapValue{
			"name": lang.StringValue("John"),
		},
	}

	tests := []struct {
		name     string
		path     string
		value    lang.Value
		hasError bool
	}{
		{"set nested existing", "user.name", lang.StringValue("Jane"), false},
		{"set new nested", "user.email", lang.StringValue("jane@example.com"), false},
		{"create new path", "settings.theme", lang.StringValue("dark"), false},
		{"empty path", "", lang.StringValue("value"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{testMap, lang.StringValue(tt.path), tt.value})
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

			// Verify the path was set correctly
			_, getPathFn := getPath()
			getValue, err := getPathFn([]lang.Value{result, lang.StringValue(tt.path)})
			if err != nil {
				t.Errorf("Error getting set value: %v", err)
				return
			}
			if getValue != tt.value {
				t.Errorf("Expected set value %v, got %v", tt.value, getValue)
			}
		})
	}
}

func TestHasPath(t *testing.T) {
	_, fn := hasPath()
	nestedMap := lang.MapValue{
		"user": lang.MapValue{
			"profile": lang.MapValue{
				"name": lang.StringValue("John"),
			},
		},
	}

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{"existing nested path", "user.profile.name", true},
		{"existing partial path", "user.profile", true},
		{"non-existing path", "user.email", false},
		{"non-existing nested path", "user.profile.age", false},
		{"empty path", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{nestedMap, lang.StringValue(tt.path)})
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

func TestDeletePath(t *testing.T) {
	_, fn := deletePath()
	nestedMap := lang.MapValue{
		"user": lang.MapValue{
			"name":  lang.StringValue("John"),
			"email": lang.StringValue("john@example.com"),
		},
		"settings": lang.MapValue{
			"theme": lang.StringValue("dark"),
		},
	}

	tests := []struct {
		name     string
		path     string
		hasError bool
	}{
		{"delete nested path", "user.email", false},
		{"delete simple path", "settings", false},
		{"delete non-existent path", "user.age", false},
		{"empty path", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{nestedMap, lang.StringValue(tt.path)})
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

			// Verify the path was deleted
			_, hasPathFn := hasPath()
			hasResult, err := hasPathFn([]lang.Value{result, lang.StringValue(tt.path)})
			if err != nil {
				t.Errorf("Error checking deleted path: %v", err)
				return
			}
			if tt.path != "" && bool(hasResult.(lang.BoolValue)) {
				t.Errorf("Expected path %s to be deleted", tt.path)
			}
		})
	}
}

func TestHelperFunctions(t *testing.T) {
	t.Run("valueToStringForMap", func(t *testing.T) {
		tests := []struct {
			input    lang.Value
			expected string
		}{
			{lang.StringValue("test"), "test"},
			{lang.NumberValue(42), "42"},
			{lang.NumberValue(42.5), "42.5"},
			{lang.BoolValue(true), "true"},
			{lang.BoolValue(false), "false"},
			{nil, ""},
		}

		for _, tt := range tests {
			result := valueToStringForMap(tt.input)
			if result != tt.expected {
				t.Errorf("valueToStringForMap(%v) = %s, expected %s", tt.input, result, tt.expected)
			}
		}
	})

	t.Run("isNullValueForMap", func(t *testing.T) {
		tests := []struct {
			input    lang.Value
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
			result := isNullValueForMap(tt.input)
			if result != tt.expected {
				t.Errorf("isNullValueForMap(%v) = %v, expected %v", tt.input, result, tt.expected)
			}
		}
	})

	t.Run("deepCopyValue", func(t *testing.T) {
		originalMap := lang.MapValue{
			"nested": lang.MapValue{
				"key": lang.StringValue("value"),
			},
			"list": lang.ListValue{lang.NumberValue(1), lang.NumberValue(2)},
		}

		copied := deepCopyValue(originalMap)
		copiedMap := copied.(lang.MapValue)

		// Modify the copy
		nestedMap := copiedMap["nested"].(lang.MapValue)
		nestedMap["key"] = lang.StringValue("modified")

		// Original should be unchanged
		originalNested := originalMap["nested"].(lang.MapValue)
		if originalNested["key"] != lang.StringValue("value") {
			t.Errorf("Original map was modified during deep copy")
		}
	})
}

func TestArgumentErrors(t *testing.T) {
	functions := []struct {
		fn           func() (string, lang.Function)
		expectedArgs int
	}{
		{keys, 1},
		{values, 1},
		{size, 1},
		{isEmpty, 1},
		{set, 3},
		{remove, 2},
		{invert, 1},
		{filter, 1},
		{filterKeys, 2},
		{omitKeys, 2},
		{rename, 2},
		{toList, 1},
		{fromList, 1},
		{toQueryString, 1},
		{fromQueryString, 1},
		{setPath, 3},
		{hasPath, 2},
		{deletePath, 2},
	}

	for _, tf := range functions {
		name, fn := tf.fn()
		t.Run(name+"_wrong_args", func(t *testing.T) {
			args := make([]lang.Value, tf.expectedArgs-1)
			for i := range args {
				args[i] = createTestMap()
			}
			_, err := fn(args)
			if err == nil {
				t.Errorf("Expected error for wrong argument count in %s", name)
			}
		})
	}
}

func TestNonMapInputs(t *testing.T) {
	functions := []func() (string, lang.Function){
		keys, values, size, isEmpty, set, remove, invert, filter,
		toList, toQueryString,
	}

	nonMapInput := lang.StringValue("not a map")

	for _, getFn := range functions {
		name, fn := getFn()
		t.Run(name+"_non_map", func(t *testing.T) {
			args := []lang.Value{nonMapInput}
			if name == "set" {
				args = append(args, lang.StringValue("key"), lang.StringValue("value"))
			} else if name == "remove" {
				args = append(args, lang.StringValue("key"))
			}

			_, err := fn(args)
			if err == nil {
				t.Errorf("Expected error for non-map input in %s", name)
			}
		})
	}
}

func TestComplexScenarios(t *testing.T) {
	t.Run("round_trip_conversion", func(t *testing.T) {
		originalMap := createTestMap()

		// Convert to list and back
		_, toListFn := toList()
		_, fromListFn := fromList()

		listResult, err := toListFn([]lang.Value{originalMap})
		if err != nil {
			t.Errorf("Error converting to list: %v", err)
			return
		}

		mapResult, err := fromListFn([]lang.Value{listResult})
		if err != nil {
			t.Errorf("Error converting from list: %v", err)
			return
		}

		resultMap := mapResult.(lang.MapValue)
		if len(resultMap) != len(originalMap) {
			t.Errorf("Round trip conversion lost data")
		}

		for key, value := range originalMap {
			if resultMap[key] != value {
				t.Errorf("Value mismatch for key %s", key)
			}
		}
	})

	t.Run("deep_path_operations", func(t *testing.T) {
		// Create a deeply nested structure
		deepMap := lang.MapValue{
			"level1": lang.MapValue{
				"level2": lang.MapValue{
					"level3": lang.MapValue{
						"value": lang.StringValue("deep"),
					},
				},
			},
		}

		// Test deep get/set/has/delete operations
		_, getPathFn := getPath()
		_, setPathFn := setPath()
		_, hasPathFn := hasPath()
		_, deletePathFn := deletePath()

		// Test deep get
		result, err := getPathFn([]lang.Value{deepMap, lang.StringValue("level1.level2.level3.value")})
		if err != nil {
			t.Errorf("Error getting deep path: %v", err)
		}
		if result != lang.StringValue("deep") {
			t.Errorf("Expected 'deep', got %v", result)
		}

		// Test deep set
		modified, err := setPathFn([]lang.Value{deepMap, lang.StringValue("level1.level2.level3.new"), lang.StringValue("added")})
		if err != nil {
			t.Errorf("Error setting deep path: %v", err)
		}

		// Verify the new value exists
		newResult, err := getPathFn([]lang.Value{modified, lang.StringValue("level1.level2.level3.new")})
		if err != nil {
			t.Errorf("Error getting newly set value: %v", err)
		}
		if newResult != lang.StringValue("added") {
			t.Errorf("Expected 'added', got %v", newResult)
		}

		// Test deep has
		hasResult, err := hasPathFn([]lang.Value{modified, lang.StringValue("level1.level2.level3.new")})
		if err != nil {
			t.Errorf("Error checking deep path: %v", err)
		}
		if !bool(hasResult.(lang.BoolValue)) {
			t.Errorf("Expected path to exist")
		}

		// Test deep delete
		deleted, err := deletePathFn([]lang.Value{modified, lang.StringValue("level1.level2.level3.new")})
		if err != nil {
			t.Errorf("Error deleting deep path: %v", err)
		}

		// Verify the value no longer exists
		hasAfterDelete, err := hasPathFn([]lang.Value{deleted, lang.StringValue("level1.level2.level3.new")})
		if err != nil {
			t.Errorf("Error checking deleted path: %v", err)
		}
		if bool(hasAfterDelete.(lang.BoolValue)) {
			t.Errorf("Expected path to be deleted")
		}
	})
}

func TestEdgeCases(t *testing.T) {
	t.Run("empty_map_operations", func(t *testing.T) {
		emptyMap := lang.MapValue{}

		// Test operations on empty map
		_, keysFn := keys()
		result, err := keysFn([]lang.Value{emptyMap})
		if err != nil {
			t.Errorf("Error getting keys from empty map: %v", err)
		}
		if len(result.(lang.ListValue)) != 0 {
			t.Errorf("Expected empty keys list")
		}

		// Test filter on empty map
		_, filterFn := filter()
		filtered, err := filterFn([]lang.Value{emptyMap})
		if err != nil {
			t.Errorf("Error filtering empty map: %v", err)
		}
		if len(filtered.(lang.MapValue)) != 0 {
			t.Errorf("Expected empty filtered map")
		}
	})

	t.Run("special_characters_in_keys", func(t *testing.T) {
		specialMap := lang.MapValue{
			"key.with.dots":        lang.StringValue("value1"),
			"key with spaces":      lang.StringValue("value2"),
			"key-with-dashes":      lang.StringValue("value3"),
			"key_with_underscores": lang.StringValue("value4"),
		}

		_, keysFn := keys()
		result, err := keysFn([]lang.Value{specialMap})
		if err != nil {
			t.Errorf("Error getting keys with special characters: %v", err)
		}

		keys := result.(lang.ListValue)
		if len(keys) != len(specialMap) {
			t.Errorf("Expected %d keys, got %d", len(specialMap), len(keys))
		}
	})
}

func TestExport(t *testing.T) {
	functions := Export()

	expectedFunctions := []string{
		"keys", "values", "size", "isEmpty", "has", "get", "set", "delete",
		"merge", "mergeDeep", "invert", "filter", "filterKeys", "omitKeys",
		"rename", "toList", "fromList", "toQueryString", "fromQueryString",
		"getPath", "setPath", "hasPath", "deletePath",
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

func BenchmarkKeys(b *testing.B) {
	_, fn := keys()
	testMap := createTestMap()
	args := []lang.Value{testMap}

	for i := 0; i < b.N; i++ {
		fn(args)
	}
}

func BenchmarkGet(b *testing.B) {
	_, fn := get()
	testMap := createTestMap()
	args := []lang.Value{testMap, lang.StringValue("name")}

	for i := 0; i < b.N; i++ {
		fn(args)
	}
}

func BenchmarkSet(b *testing.B) {
	_, fn := set()
	testMap := createTestMap()
	args := []lang.Value{testMap, lang.StringValue("newKey"), lang.StringValue("newValue")}

	for i := 0; i < b.N; i++ {
		fn(args)
	}
}

func BenchmarkMerge(b *testing.B) {
	_, fn := merge()
	map1 := createTestMap()
	map2 := lang.MapValue{"extra": lang.StringValue("value")}
	args := []lang.Value{map1, map2}

	for i := 0; i < b.N; i++ {
		fn(args)
	}
}

// Helper function for string contains check
func contains(s, substr string) bool {
	return len(substr) <= len(s) && (len(substr) == 0 || s[len(s)-len(substr):] == substr ||
		s[:len(substr)] == substr || containsMiddle(s, substr))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
