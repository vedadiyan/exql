/*
 * Copyright 2025 Pouya Vedadiyan
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package maps

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/vedadiyan/exql/lang"
	"github.com/vedadiyan/exql/lib"
)

// Basic Map Operations
func keys() (string, lang.Function) {
	name := "keys"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		m, ok := args[0].(lang.MapValue)
		if !ok {
			return nil, lib.MapError(name, args[0])
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
	return name, fn
}

func values() (string, lang.Function) {
	name := "values"
	_, Keys := keys()
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		m, ok := args[0].(lang.MapValue)
		if !ok {
			return nil, lib.MapError(name, args[0])
		}

		// Get keys in sorted order for consistent results
		keys, err := Keys(args)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		keysList := keys.(lang.ListValue)

		values := make(lang.ListValue, 0, len(m))
		for _, keyVal := range keysList {
			key := string(keyVal.(lang.StringValue))
			values = append(values, m[key])
		}

		return values, nil
	}
	return name, fn
}

func size() (string, lang.Function) {
	name := "size"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		m, ok := args[0].(lang.MapValue)
		if !ok {
			return nil, lib.MapError(name, args[0])
		}
		return lang.NumberValue(float64(len(m))), nil
	}
	return name, fn
}

func isEmpty() (string, lang.Function) {
	name := "isEmpty"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		m, ok := args[0].(lang.MapValue)
		if !ok {
			return nil, lib.MapError(name, args[0])
		}
		return lang.BoolValue(len(m) == 0), nil
	}
	return name, fn
}

func has() (string, lang.Function) {
	name := "has"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("%s: expected 2 arguments (map, key)", name)
		}
		m, ok := args[0].(lang.MapValue)
		if !ok {
			return nil, lib.MapError(name, args[0])
		}

		keyStr, err := lib.ToString(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: key %w", name, err)
		}
		key := string(keyStr)
		_, exists := m[key]
		return lang.BoolValue(exists), nil
	}
	return name, fn
}

func get() (string, lang.Function) {
	name := "get"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) < 2 || len(args) > 3 {
			return nil, fmt.Errorf("%s: expected 2 or 3 arguments (map, key, default?)", name)
		}
		m, ok := args[0].(lang.MapValue)
		if !ok {
			return nil, lib.MapError(name, args[0])
		}

		keyStr, err := lib.ToString(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: key %w", name, err)
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
	return name, fn
}

func set() (string, lang.Function) {
	name := "set"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 3 {
			return nil, lib.ArgumentError(name, 3)
		}
		m, ok := args[0].(lang.MapValue)
		if !ok {
			return nil, lib.MapError(name, args[0])
		}

		keyStr, err := lib.ToString(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: key %w", name, err)
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
	return name, fn
}

func remove() (string, lang.Function) {
	name := "delete"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 2 {
			return nil, lib.ArgumentError(name, 2)
		}
		m, ok := args[0].(lang.MapValue)
		if !ok {
			return nil, lib.MapError(name, args[0])
		}

		keyStr, err := lib.ToString(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: key %w", name, err)
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
	return name, fn
}

func merge() (string, lang.Function) {
	name := "merge"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) == 0 {
			return nil, fmt.Errorf("%s: expected at least 1 argument", name)
		}

		result := make(lang.MapValue)

		for i, arg := range args {
			if m, ok := arg.(lang.MapValue); ok {
				for key, value := range m {
					result[key] = value
				}
			} else {
				return nil, fmt.Errorf("%s: argument %d expected map, got %T", name, i+1, arg)
			}
		}

		return result, nil
	}
	return name, fn
}

func mergeDeep() (string, lang.Function) {
	name := "mergeDeep"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) == 0 {
			return nil, fmt.Errorf("%s: expected at least 1 argument", name)
		}

		result := make(lang.MapValue)
		if first, ok := args[0].(lang.MapValue); ok {
			// Deep copy first map
			for k, v := range first {
				result[k] = deepCopyValue(v)
			}
		} else {
			return nil, fmt.Errorf("%s: argument 1 expected map, got %T", name, args[0])
		}

		for i := 1; i < len(args); i++ {
			if m, ok := args[i].(lang.MapValue); ok {
				result = deepMergeMaps(result, m)
			} else {
				return nil, fmt.Errorf("%s: argument %d expected map, got %T", name, i+1, args[i])
			}
		}

		return result, nil
	}
	return name, fn
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
func invert() (string, lang.Function) {
	name := "invert"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		m, ok := args[0].(lang.MapValue)
		if !ok {
			return nil, lib.MapError(name, args[0])
		}

		result := make(lang.MapValue, len(m))
		for key, value := range m {
			valueStr := valueToStringForMap(value)
			result[valueStr] = lang.StringValue(key)
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
		m, ok := args[0].(lang.MapValue)
		if !ok {
			return nil, lib.MapError(name, args[0])
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
	return name, fn
}

func filterKeys() (string, lang.Function) {
	name := "filterKeys"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 2 {
			return nil, lib.ArgumentError(name, 2)
		}
		m, ok := args[0].(lang.MapValue)
		if !ok {
			return nil, lib.MapError(name, args[0])
		}

		keysToKeep, ok := args[1].(lang.ListValue)
		if !ok {
			return nil, fmt.Errorf("%s: expected list for keys, got %T", name, args[1])
		}

		// Convert keys to strings for comparison
		keySet := make(map[string]bool)
		for i, keyVal := range keysToKeep {
			keyStr, err := lib.ToString(keyVal)
			if err != nil {
				return nil, fmt.Errorf("%s: key %d %w", name, i, err)
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
	return name, fn
}

func omitKeys() (string, lang.Function) {
	name := "omitKeys"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 2 {
			return nil, lib.ArgumentError(name, 2)
		}
		m, ok := args[0].(lang.MapValue)
		if !ok {
			return nil, lib.MapError(name, args[0])
		}

		keysToOmit, ok := args[1].(lang.ListValue)
		if !ok {
			return nil, fmt.Errorf("%s: expected list for keys, got %T", name, args[1])
		}

		// Convert keys to strings for comparison
		keySet := make(map[string]bool)
		for i, keyVal := range keysToOmit {
			keyStr, err := lib.ToString(keyVal)
			if err != nil {
				return nil, fmt.Errorf("%s: key %d %w", name, i, err)
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
	return name, fn
}

func rename() (string, lang.Function) {
	name := "rename"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 2 {
			return nil, lib.ArgumentError(name, 2)
		}
		m, ok := args[0].(lang.MapValue)
		if !ok {
			return nil, lib.MapError(name, args[0])
		}

		renameMap, ok := args[1].(lang.MapValue)
		if !ok {
			return nil, fmt.Errorf("%s: expected map for rename mapping, got %T", name, args[1])
		}

		result := make(lang.MapValue, len(m))

		for key, value := range m {
			newKey := key
			if newKeyVal, exists := renameMap[key]; exists {
				newKeyStr, err := lib.ToString(newKeyVal)
				if err != nil {
					return nil, fmt.Errorf("%s: new key for '%s' %w", name, key, err)
				}
				newKey = string(newKeyStr)
			}
			result[newKey] = value
		}

		return result, nil
	}
	return name, fn
}

// Map Conversion
func toList() (string, lang.Function) {
	name := "toList"
	_, Keys := keys()
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		m, ok := args[0].(lang.MapValue)
		if !ok {
			return nil, lib.MapError(name, args[0])
		}

		// Get keys in sorted order for consistent results
		keys, err := Keys(args)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
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
	return name, fn
}

func fromList() (string, lang.Function) {
	name := "fromList"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		list, ok := args[0].(lang.ListValue)
		if !ok {
			return nil, fmt.Errorf("%s: expected list, got %T", name, args[0])
		}

		result := make(lang.MapValue)

		for i, item := range list {
			if pair, ok := item.(lang.ListValue); ok {
				if len(pair) < 2 {
					return nil, fmt.Errorf("%s: pair %d must have at least 2 elements, got %d", name, i, len(pair))
				}
				keyStr, err := lib.ToString(pair[0])
				if err != nil {
					return nil, fmt.Errorf("%s: pair %d key %w", name, i, err)
				}
				key := string(keyStr)
				value := pair[1]
				result[key] = value
			} else {
				return nil, fmt.Errorf("%s: item %d expected list (key-value pair), got %T", name, i, item)
			}
		}

		return result, nil
	}
	return name, fn
}

func toQueryString() (string, lang.Function) {
	name := "toQueryString"
	_, Keys := keys()
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		m, ok := args[0].(lang.MapValue)
		if !ok {
			return nil, lib.MapError(name, args[0])
		}

		var parts []string

		// Get keys in sorted order for consistent results
		keys, err := Keys(args)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
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
						return nil, fmt.Errorf("%s: array value %w", name, err)
					}
					parts = append(parts, key+"="+string(itemStr))
				}
			} else {
				valueStr, err := lib.ToString(value)
				if err != nil {
					return nil, fmt.Errorf("%s: value for key '%s' %w", name, key, err)
				}
				parts = append(parts, key+"="+string(valueStr))
			}
		}

		return lang.StringValue(strings.Join(parts, "&")), nil
	}
	return name, fn
}

func fromQueryString() (string, lang.Function) {
	name := "fromQueryString"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		queryStr, err := lib.ToString(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
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
	return name, fn
}

// Map Path Operations (dot notation)
func getPath() (string, lang.Function) {
	name := "getPath"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) < 2 || len(args) > 3 {
			return nil, fmt.Errorf("%s: expected 2 or 3 arguments (map, path, default?)", name)
		}
		m, ok := args[0].(lang.MapValue)
		if !ok {
			return nil, lib.MapError(name, args[0])
		}

		pathStr, err := lib.ToString(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: path %w", name, err)
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
				return nil, fmt.Errorf("%s: path '%s' invalid at part '%s' (not a map)", name, path, part)
			}
		}

		return current, nil
	}
	return name, fn
}

func setPath() (string, lang.Function) {
	name := "setPath"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 3 {
			return nil, lib.ArgumentError(name, 3)
		}
		m, ok := args[0].(lang.MapValue)
		if !ok {
			return nil, lib.MapError(name, args[0])
		}

		pathStr, err := lib.ToString(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: path %w", name, err)
		}
		path := string(pathStr)
		value := args[2]

		if path == "" {
			return nil, fmt.Errorf("%s: path cannot be empty", name)
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
	return name, fn
}

func hasPath() (string, lang.Function) {
	name := "hasPath"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 2 {
			return nil, lib.ArgumentError(name, 2)
		}
		m, ok := args[0].(lang.MapValue)
		if !ok {
			return nil, lib.MapError(name, args[0])
		}

		pathStr, err := lib.ToString(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: path %w", name, err)
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
	return name, fn
}

func deletePath() (string, lang.Function) {
	name := "deletePath"
	_, Delete := remove()
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 2 {
			return nil, lib.ArgumentError(name, 2)
		}
		m, ok := args[0].(lang.MapValue)
		if !ok {
			return nil, lib.MapError(name, args[0])
		}

		pathStr, err := lib.ToString(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: path %w", name, err)
		}
		path := string(pathStr)

		if path == "" {
			return nil, fmt.Errorf("%s: path cannot be empty", name)
		}

		parts := strings.Split(path, ".")

		if len(parts) == 1 {
			// Simple key deletion
			return Delete(args)
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
	return name, fn
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

var mapFunctions = []func() (string, lang.Function){
	// Basic operations
	keys,
	values,
	size,
	isEmpty,
	has,
	get,
	set,
	remove,

	// Merging
	merge,
	mergeDeep,

	// Transformation
	invert,
	filter,
	filterKeys,
	omitKeys,
	rename,

	// Conversion
	toList,
	fromList,
	toQueryString,
	fromQueryString,

	// Path operations
	getPath,
	setPath,
	hasPath,
	deletePath,
}

func Export() map[string]lang.Function {
	out := make(map[string]lang.Function)
	for _, value := range mapFunctions {
		name, fn := value()
		out[name] = fn
	}
	return out
}
