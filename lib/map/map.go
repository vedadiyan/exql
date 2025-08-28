package maps

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/vedadiyan/exql/lang"
	"github.com/vedadiyan/exql/lib"
)

// Map/Object Functions
// These functions provide comprehensive map and object manipulation capabilities

// Basic Map Operations
func mapKeys(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("keys: expected 1 argument")
	}
	m, ok := args[0].(lang.MapValue)
	if !ok {
		return nil, fmt.Errorf("keys: expected map, got %T", args[0])
	}

	keys := make(lang.ListValue, 0, len(m))
	for key := range m {
		keys = append(keys, lang.StringValue(key))
	}

	// Sort keys for consistent ordering
	sort.Slice(keys, func(i, j int) bool {
		return string(keys[i].(lang.StringValue)) < string(keys[j].(lang.StringValue))
	})

	return keys, nil
}

func mapValues(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("values: expected 1 argument")
	}
	m, ok := args[0].(lang.MapValue)
	if !ok {
		return nil, fmt.Errorf("values: expected map, got %T", args[0])
	}

	// Get keys in sorted order for consistent results
	keys, err := mapKeys(args)
	if err != nil {
		return nil, fmt.Errorf("values: %w", err)
	}
	keysList := keys.(lang.ListValue)

	values := make(lang.ListValue, 0, len(m))
	for _, keyVal := range keysList {
		key := string(keyVal.(lang.StringValue))
		values = append(values, m[key])
	}

	return values, nil
}

func mapSize(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("map_size: expected 1 argument")
	}
	m, ok := args[0].(lang.MapValue)
	if !ok {
		return nil, fmt.Errorf("map_size: expected map, got %T", args[0])
	}
	return lang.NumberValue(float64(len(m))), nil
}

func mapIsEmpty(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("map_is_empty: expected 1 argument")
	}
	m, ok := args[0].(lang.MapValue)
	if !ok {
		return nil, fmt.Errorf("map_is_empty: expected map, got %T", args[0])
	}
	return lang.BoolValue(len(m) == 0), nil
}

func mapHas(args []lang.Value) (lang.Value, error) {
	if len(args) != 2 {
		return nil, errors.New("has: expected 2 arguments (map, key)")
	}
	m, ok := args[0].(lang.MapValue)
	if !ok {
		return nil, fmt.Errorf("has: expected map, got %T", args[0])
	}

	keyStr, err := lib.ToString(args[1])
	if err != nil {
		return nil, fmt.Errorf("has: key %w", err)
	}
	key := string(keyStr)
	_, exists := m[key]
	return lang.BoolValue(exists), nil
}

func mapGet(args []lang.Value) (lang.Value, error) {
	if len(args) < 2 || len(args) > 3 {
		return nil, errors.New("map_get: expected 2 or 3 arguments (map, key, default?)")
	}
	m, ok := args[0].(lang.MapValue)
	if !ok {
		return nil, fmt.Errorf("map_get: expected map, got %T", args[0])
	}

	keyStr, err := lib.ToString(args[1])
	if err != nil {
		return nil, fmt.Errorf("map_get: key %w", err)
	}
	key := string(keyStr)

	value, exists := m[key]
	if !exists {
		if len(args) == 3 {
			return args[2], nil // Return default value
		}
		return nil, nil
	}

	return value, nil
}

func mapSet(args []lang.Value) (lang.Value, error) {
	if len(args) != 3 {
		return nil, errors.New("map_set: expected 3 arguments (map, key, value)")
	}
	m, ok := args[0].(lang.MapValue)
	if !ok {
		return nil, fmt.Errorf("map_set: expected map, got %T", args[0])
	}

	keyStr, err := lib.ToString(args[1])
	if err != nil {
		return nil, fmt.Errorf("map_set: key %w", err)
	}
	key := string(keyStr)
	value := args[2]

	// Create a new map with the updated value
	result := make(lang.MapValue, len(m)+1)
	for k, v := range m {
		result[k] = v
	}
	result[key] = value

	return result, nil
}

func mapDelete(args []lang.Value) (lang.Value, error) {
	if len(args) != 2 {
		return nil, errors.New("map_delete: expected 2 arguments (map, key)")
	}
	m, ok := args[0].(lang.MapValue)
	if !ok {
		return nil, fmt.Errorf("map_delete: expected map, got %T", args[0])
	}

	keyStr, err := lib.ToString(args[1])
	if err != nil {
		return nil, fmt.Errorf("map_delete: key %w", err)
	}
	key := string(keyStr)

	// Create a new map without the key
	result := make(lang.MapValue, len(m))
	for k, v := range m {
		if k != key {
			result[k] = v
		}
	}

	return result, nil
}

func mapMerge(args []lang.Value) (lang.Value, error) {
	if len(args) == 0 {
		return nil, errors.New("merge: expected at least 1 argument")
	}

	result := make(lang.MapValue)

	for i, arg := range args {
		if m, ok := arg.(lang.MapValue); ok {
			for key, value := range m {
				result[key] = value
			}
		} else {
			return nil, fmt.Errorf("merge: argument %d expected map, got %T", i+1, arg)
		}
	}

	return result, nil
}

func mapMergeDeep(args []lang.Value) (lang.Value, error) {
	if len(args) == 0 {
		return nil, errors.New("merge_deep: expected at least 1 argument")
	}

	result := make(lang.MapValue)
	if first, ok := args[0].(lang.MapValue); ok {
		// Deep copy first map
		for k, v := range first {
			result[k] = deepCopyValue(v)
		}
	} else {
		return nil, fmt.Errorf("merge_deep: argument 1 expected map, got %T", args[0])
	}

	for i := 1; i < len(args); i++ {
		if m, ok := args[i].(lang.MapValue); ok {
			result = deepMergeMaps(result, m)
		} else {
			return nil, fmt.Errorf("merge_deep: argument %d expected map, got %T", i+1, args[i])
		}
	}

	return result, nil
}

func deepMergeMaps(dest, src lang.MapValue) lang.MapValue {
	for key, srcValue := range src {
		if destValue, exists := dest[key]; exists {
			// If both values are maps, merge them recursively
			if destMap, destIsMap := destValue.(lang.MapValue); destIsMap {
				if srcMap, srcIsMap := srcValue.(lang.MapValue); srcIsMap {
					dest[key] = deepMergeMaps(destMap, srcMap)
					continue
				}
			}
		}
		// Otherwise, just set the value
		dest[key] = deepCopyValue(srcValue)
	}
	return dest
}

func deepCopyValue(v lang.Value) lang.Value {
	switch val := v.(type) {
	case lang.MapValue:
		result := make(lang.MapValue, len(val))
		for k, v := range val {
			result[k] = deepCopyValue(v)
		}
		return result
	case lang.ListValue:
		result := make(lang.ListValue, len(val))
		for i, v := range val {
			result[i] = deepCopyValue(v)
		}
		return result
	default:
		return v // Primitive values can be copied directly
	}
}

// Map Transformation
func mapInvert(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("map_invert: expected 1 argument")
	}
	m, ok := args[0].(lang.MapValue)
	if !ok {
		return nil, fmt.Errorf("map_invert: expected map, got %T", args[0])
	}

	result := make(lang.MapValue, len(m))
	for key, value := range m {
		valueStr := valueToStringForMap(value)
		result[valueStr] = lang.StringValue(key)
	}

	return result, nil
}

func mapFilter(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("map_filter: expected 1 argument")
	}
	m, ok := args[0].(lang.MapValue)
	if !ok {
		return nil, fmt.Errorf("map_filter: expected map, got %T", args[0])
	}

	// Filter out nil/null values
	result := make(lang.MapValue)
	for key, value := range m {
		if value != nil && !isNullValueForMap(value) {
			result[key] = value
		}
	}

	return result, nil
}

func mapFilterKeys(args []lang.Value) (lang.Value, error) {
	if len(args) != 2 {
		return nil, errors.New("map_filter_keys: expected 2 arguments (map, keys_list)")
	}
	m, ok := args[0].(lang.MapValue)
	if !ok {
		return nil, fmt.Errorf("map_filter_keys: expected map, got %T", args[0])
	}

	keysToKeep, ok := args[1].(lang.ListValue)
	if !ok {
		return nil, fmt.Errorf("map_filter_keys: expected list for keys, got %T", args[1])
	}

	// Convert keys to strings for comparison
	keySet := make(map[string]bool)
	for i, keyVal := range keysToKeep {
		keyStr, err := lib.ToString(keyVal)
		if err != nil {
			return nil, fmt.Errorf("map_filter_keys: key %d %w", i, err)
		}
		key := string(keyStr)
		keySet[key] = true
	}

	result := make(lang.MapValue)
	for key, value := range m {
		if keySet[key] {
			result[key] = value
		}
	}

	return result, nil
}

func mapOmitKeys(args []lang.Value) (lang.Value, error) {
	if len(args) != 2 {
		return nil, errors.New("map_omit_keys: expected 2 arguments (map, keys_list)")
	}
	m, ok := args[0].(lang.MapValue)
	if !ok {
		return nil, fmt.Errorf("map_omit_keys: expected map, got %T", args[0])
	}

	keysToOmit, ok := args[1].(lang.ListValue)
	if !ok {
		return nil, fmt.Errorf("map_omit_keys: expected list for keys, got %T", args[1])
	}

	// Convert keys to strings for comparison
	keySet := make(map[string]bool)
	for i, keyVal := range keysToOmit {
		keyStr, err := lib.ToString(keyVal)
		if err != nil {
			return nil, fmt.Errorf("map_omit_keys: key %d %w", i, err)
		}
		key := string(keyStr)
		keySet[key] = true
	}

	result := make(lang.MapValue)
	for key, value := range m {
		if !keySet[key] {
			result[key] = value
		}
	}

	return result, nil
}

func mapRename(args []lang.Value) (lang.Value, error) {
	if len(args) != 2 {
		return nil, errors.New("map_rename: expected 2 arguments (map, rename_map)")
	}
	m, ok := args[0].(lang.MapValue)
	if !ok {
		return nil, fmt.Errorf("map_rename: expected map, got %T", args[0])
	}

	renameMap, ok := args[1].(lang.MapValue)
	if !ok {
		return nil, fmt.Errorf("map_rename: expected map for rename mapping, got %T", args[1])
	}

	result := make(lang.MapValue, len(m))

	for key, value := range m {
		newKey := key
		if newKeyVal, exists := renameMap[key]; exists {
			newKeyStr, err := lib.ToString(newKeyVal)
			if err != nil {
				return nil, fmt.Errorf("map_rename: new key for '%s' %w", key, err)
			}
			newKey = string(newKeyStr)
		}
		result[newKey] = value
	}

	return result, nil
}

// Map Conversion
func mapToList(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("map_to_list: expected 1 argument")
	}
	m, ok := args[0].(lang.MapValue)
	if !ok {
		return nil, fmt.Errorf("map_to_list: expected map, got %T", args[0])
	}

	// Get keys in sorted order for consistent results
	keys, err := mapKeys(args)
	if err != nil {
		return nil, fmt.Errorf("map_to_list: %w", err)
	}
	keysList := keys.(lang.ListValue)

	result := make(lang.ListValue, 0, len(m))
	for _, keyVal := range keysList {
		key := string(keyVal.(lang.StringValue))
		pair := lang.ListValue{lang.StringValue(key), m[key]}
		result = append(result, pair)
	}

	return result, nil
}

func mapFromList(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("map_from_list: expected 1 argument")
	}
	list, ok := args[0].(lang.ListValue)
	if !ok {
		return nil, fmt.Errorf("map_from_list: expected list, got %T", args[0])
	}

	result := make(lang.MapValue)

	for i, item := range list {
		if pair, ok := item.(lang.ListValue); ok {
			if len(pair) < 2 {
				return nil, fmt.Errorf("map_from_list: pair %d must have at least 2 elements, got %d", i, len(pair))
			}
			keyStr, err := lib.ToString(pair[0])
			if err != nil {
				return nil, fmt.Errorf("map_from_list: pair %d key %w", i, err)
			}
			key := string(keyStr)
			value := pair[1]
			result[key] = value
		} else {
			return nil, fmt.Errorf("map_from_list: item %d expected list (key-value pair), got %T", i, item)
		}
	}

	return result, nil
}

func mapToQueryString(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("map_to_query_string: expected 1 argument")
	}
	m, ok := args[0].(lang.MapValue)
	if !ok {
		return nil, fmt.Errorf("map_to_query_string: expected map, got %T", args[0])
	}

	var parts []string

	// Get keys in sorted order for consistent results
	keys, err := mapKeys(args)
	if err != nil {
		return nil, fmt.Errorf("map_to_query_string: %w", err)
	}
	keysList := keys.(lang.ListValue)

	for _, keyVal := range keysList {
		key := string(keyVal.(lang.StringValue))
		value := m[key]

		if list, ok := value.(lang.ListValue); ok {
			// Handle array values
			for _, item := range list {
				itemStr, err := lib.ToString(item)
				if err != nil {
					return nil, fmt.Errorf("map_to_query_string: array value %w", err)
				}
				parts = append(parts, key+"="+string(itemStr))
			}
		} else {
			valueStr, err := lib.ToString(value)
			if err != nil {
				return nil, fmt.Errorf("map_to_query_string: value for key '%s' %w", key, err)
			}
			parts = append(parts, key+"="+string(valueStr))
		}
	}

	return lang.StringValue(strings.Join(parts, "&")), nil
}

func mapFromQueryString(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("map_from_query_string: expected 1 argument")
	}
	queryStr, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("map_from_query_string: %w", err)
	}

	result := make(lang.MapValue)

	if string(queryStr) == "" {
		return result, nil
	}

	pairs := strings.Split(string(queryStr), "&")
	for _, pair := range pairs {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			key := parts[0]
			value := parts[1]

			// Check if key already exists
			if existing, exists := result[key]; exists {
				// Convert to list if not already
				if list, ok := existing.(lang.ListValue); ok {
					result[key] = append(list, lang.StringValue(value))
				} else {
					result[key] = lang.ListValue{existing, lang.StringValue(value)}
				}
			} else {
				result[key] = lang.StringValue(value)
			}
		}
	}

	return result, nil
}

// Map Path Operations (dot notation)
func mapGetPath(args []lang.Value) (lang.Value, error) {
	if len(args) < 2 || len(args) > 3 {
		return nil, errors.New("map_get_path: expected 2 or 3 arguments (map, path, default?)")
	}
	m, ok := args[0].(lang.MapValue)
	if !ok {
		return nil, fmt.Errorf("map_get_path: expected map, got %T", args[0])
	}

	pathStr, err := lib.ToString(args[1])
	if err != nil {
		return nil, fmt.Errorf("map_get_path: path %w", err)
	}
	path := string(pathStr)
	current := lang.Value(m)

	if path == "" {
		return current, nil
	}

	parts := strings.Split(path, ".")
	for _, part := range parts {
		if currentMap, ok := current.(lang.MapValue); ok {
			if value, exists := currentMap[part]; exists {
				current = value
			} else {
				if len(args) == 3 {
					return args[2], nil // Return default value
				}
				return nil, nil
			}
		} else {
			if len(args) == 3 {
				return args[2], nil // Return default value
			}
			return nil, fmt.Errorf("map_get_path: path '%s' invalid at part '%s' (not a map)", path, part)
		}
	}

	return current, nil
}

func mapSetPath(args []lang.Value) (lang.Value, error) {
	if len(args) != 3 {
		return nil, errors.New("map_set_path: expected 3 arguments (map, path, value)")
	}
	m, ok := args[0].(lang.MapValue)
	if !ok {
		return nil, fmt.Errorf("map_set_path: expected map, got %T", args[0])
	}

	pathStr, err := lib.ToString(args[1])
	if err != nil {
		return nil, fmt.Errorf("map_set_path: path %w", err)
	}
	path := string(pathStr)
	value := args[2]

	if path == "" {
		return nil, errors.New("map_set_path: path cannot be empty")
	}

	// Deep copy the original map
	result := deepCopyValue(m).(lang.MapValue)

	parts := strings.Split(path, ".")
	current := result

	// Navigate to the parent of the target
	for i := 0; i < len(parts)-1; i++ {
		part := parts[i]
		if next, exists := current[part]; exists {
			if nextMap, ok := next.(lang.MapValue); ok {
				current = nextMap
			} else {
				// Replace with a new map
				newMap := make(lang.MapValue)
				current[part] = newMap
				current = newMap
			}
		} else {
			// Create new map
			newMap := make(lang.MapValue)
			current[part] = newMap
			current = newMap
		}
	}

	// Set the final value
	current[parts[len(parts)-1]] = value

	return result, nil
}

func mapHasPath(args []lang.Value) (lang.Value, error) {
	if len(args) != 2 {
		return nil, errors.New("map_has_path: expected 2 arguments (map, path)")
	}
	m, ok := args[0].(lang.MapValue)
	if !ok {
		return nil, fmt.Errorf("map_has_path: expected map, got %T", args[0])
	}

	pathStr, err := lib.ToString(args[1])
	if err != nil {
		return nil, fmt.Errorf("map_has_path: path %w", err)
	}
	path := string(pathStr)
	current := lang.Value(m)

	if path == "" {
		return lang.BoolValue(true), nil
	}

	parts := strings.Split(path, ".")
	for _, part := range parts {
		if currentMap, ok := current.(lang.MapValue); ok {
			if value, exists := currentMap[part]; exists {
				current = value
			} else {
				return lang.BoolValue(false), nil
			}
		} else {
			return lang.BoolValue(false), nil
		}
	}

	return lang.BoolValue(true), nil
}

func mapDeletePath(args []lang.Value) (lang.Value, error) {
	if len(args) != 2 {
		return nil, errors.New("map_delete_path: expected 2 arguments (map, path)")
	}
	m, ok := args[0].(lang.MapValue)
	if !ok {
		return nil, fmt.Errorf("map_delete_path: expected map, got %T", args[0])
	}

	pathStr, err := lib.ToString(args[1])
	if err != nil {
		return nil, fmt.Errorf("map_delete_path: path %w", err)
	}
	path := string(pathStr)

	if path == "" {
		return nil, errors.New("map_delete_path: path cannot be empty")
	}

	parts := strings.Split(path, ".")

	if len(parts) == 1 {
		// Simple key deletion
		return mapDelete(args)
	}

	// Deep copy the original map
	result := deepCopyValue(m).(lang.MapValue)

	// Navigate to the parent of the target
	current := result
	for i := 0; i < len(parts)-1; i++ {
		part := parts[i]
		if next, exists := current[part]; exists {
			if nextMap, ok := next.(lang.MapValue); ok {
				current = nextMap
			} else {
				// Path doesn't exist
				return result, nil
			}
		} else {
			// Path doesn't exist
			return result, nil
		}
	}

	// Delete the final key
	delete(current, parts[len(parts)-1])

	return result, nil
}

// Helper functions
func valueToStringForMap(v lang.Value) string {
	switch val := v.(type) {
	case lang.StringValue:
		return string(val)
	case lang.NumberValue:
		if val == lang.NumberValue(float64(int64(val))) {
			return strconv.FormatInt(int64(val), 10)
		}
		str, _ := lib.ToString(v)
		return string(str)
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

func isNullValueForMap(v lang.Value) bool {
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
var MapFunctions = map[string]func([]lang.Value) (lang.Value, error){
	// Basic operations
	"keys":         mapKeys,
	"values":       mapValues,
	"map_size":     mapSize,
	"map_is_empty": mapIsEmpty,
	"has":          mapHas,
	"map_get":      mapGet,
	"map_set":      mapSet,
	"map_delete":   mapDelete,

	// Merging
	"merge":      mapMerge,
	"merge_deep": mapMergeDeep,

	// Transformation
	"map_invert":      mapInvert,
	"map_filter":      mapFilter,
	"map_filter_keys": mapFilterKeys,
	"map_omit_keys":   mapOmitKeys,
	"map_rename":      mapRename,

	// Conversion
	"map_to_list":           mapToList,
	"map_from_list":         mapFromList,
	"map_to_query_string":   mapToQueryString,
	"map_from_query_string": mapFromQueryString,

	// Path operations
	"map_get_path":    mapGetPath,
	"map_set_path":    mapSetPath,
	"map_has_path":    mapHasPath,
	"map_delete_path": mapDeletePath,
}
