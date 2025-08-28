package json

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/vedadiyan/exql/lang"
	"github.com/vedadiyan/exql/lib"
)

// JSON Processing Functions
// These functions help parse, manipulate, and extract data from JSON

func jsonParse(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("json_parse: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("json_parse: %w", err)
	}

	var result interface{}
	err = json.Unmarshal([]byte(string(str)), &result)
	if err != nil {
		return nil, fmt.Errorf("json_parse: invalid JSON: %w", err)
	}
	return convertJSONToValue(result), nil
}

func jsonString(args []lang.Value) (lang.Value, error) {
	if len(args) < 1 || len(args) > 2 {
		return nil, errors.New("json_string: expected 1 or 2 arguments (data, pretty?)")
	}

	indent := ""
	if len(args) == 2 {
		pretty, err := lib.ToBool(args[1])
		if err != nil {
			return nil, fmt.Errorf("json_string: pretty flag %w", err)
		}
		if pretty {
			indent = "  " // Pretty print with 2 spaces
		}
	}

	var jsonBytes []byte
	var err error

	if indent != "" {
		jsonBytes, err = json.MarshalIndent(convertValueToJSON(args[0]), "", indent)
	} else {
		jsonBytes, err = json.Marshal(convertValueToJSON(args[0]))
	}

	if err != nil {
		return nil, fmt.Errorf("json_string: marshalling failed: %w", err)
	}
	return lang.StringValue(string(jsonBytes)), nil
}

func jsonValid(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("json_valid: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("json_valid: %w", err)
	}

	var result interface{}
	err = json.Unmarshal([]byte(string(str)), &result)
	return lang.BoolValue(err == nil), nil
}

func jsonGet(args []lang.Value) (lang.Value, error) {
	if len(args) != 2 {
		return nil, errors.New("json_get: expected 2 arguments (data, path)")
	}

	// First argument can be JSON string or already parsed data
	var data interface{}
	if jsonStr, ok := args[0].(lang.StringValue); ok {
		err := json.Unmarshal([]byte(string(jsonStr)), &data)
		if err != nil {
			return nil, fmt.Errorf("json_get: invalid JSON: %w", err)
		}
	} else {
		data = convertValueToJSON(args[0])
	}

	pathStr, err := lib.ToString(args[1])
	if err != nil {
		return nil, fmt.Errorf("json_get: path %w", err)
	}
	path := string(pathStr)

	result := jsonGetByPath(data, path)
	return result, nil
}

func jsonSet(args []lang.Value) (lang.Value, error) {
	if len(args) != 3 {
		return nil, errors.New("json_set: expected 3 arguments (data, path, value)")
	}

	// First argument can be JSON string or already parsed data
	var data interface{}
	if jsonStr, ok := args[0].(lang.StringValue); ok {
		err := json.Unmarshal([]byte(string(jsonStr)), &data)
		if err != nil {
			return nil, fmt.Errorf("json_set: invalid JSON: %w", err)
		}
	} else {
		data = convertValueToJSON(args[0])
	}

	pathStr, err := lib.ToString(args[1])
	if err != nil {
		return nil, fmt.Errorf("json_set: path %w", err)
	}
	path := string(pathStr)
	value := convertValueToJSON(args[2])

	result := jsonSetByPath(data, path, value)
	return convertJSONToValue(result), nil
}

func jsonDelete(args []lang.Value) (lang.Value, error) {
	if len(args) != 2 {
		return nil, errors.New("json_delete: expected 2 arguments (data, path)")
	}

	// First argument can be JSON string or already parsed data
	var data interface{}
	if jsonStr, ok := args[0].(lang.StringValue); ok {
		err := json.Unmarshal([]byte(string(jsonStr)), &data)
		if err != nil {
			return nil, fmt.Errorf("json_delete: invalid JSON: %w", err)
		}
	} else {
		data = convertValueToJSON(args[0])
	}

	pathStr, err := lib.ToString(args[1])
	if err != nil {
		return nil, fmt.Errorf("json_delete: path %w", err)
	}
	path := string(pathStr)

	result := jsonDeleteByPath(data, path)
	return convertJSONToValue(result), nil
}

func jsonHas(args []lang.Value) (lang.Value, error) {
	if len(args) != 2 {
		return nil, errors.New("json_has: expected 2 arguments (data, path)")
	}

	result, err := jsonGet(args)
	if err != nil {
		return nil, fmt.Errorf("json_has: %w", err)
	}
	return lang.BoolValue(result != nil), nil
}

func jsonKeys(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("json_keys: expected 1 argument")
	}

	// First argument can be JSON string or already parsed data
	var data interface{}
	if jsonStr, ok := args[0].(lang.StringValue); ok {
		err := json.Unmarshal([]byte(string(jsonStr)), &data)
		if err != nil {
			return nil, fmt.Errorf("json_keys: invalid JSON: %w", err)
		}
	} else {
		data = convertValueToJSON(args[0])
	}

	switch v := data.(type) {
	case map[string]interface{}:
		keys := make(lang.ListValue, 0, len(v))
		for key := range v {
			keys = append(keys, lang.StringValue(key))
		}
		return keys, nil
	default:
		return nil, fmt.Errorf("json_keys: expected object, got %T", data)
	}
}

func jsonValues(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("json_values: expected 1 argument")
	}

	// First argument can be JSON string or already parsed data
	var data interface{}
	if jsonStr, ok := args[0].(lang.StringValue); ok {
		err := json.Unmarshal([]byte(string(jsonStr)), &data)
		if err != nil {
			return nil, fmt.Errorf("json_values: invalid JSON: %w", err)
		}
	} else {
		data = convertValueToJSON(args[0])
	}

	switch v := data.(type) {
	case map[string]interface{}:
		values := make(lang.ListValue, 0, len(v))
		for _, value := range v {
			values = append(values, convertJSONToValue(value))
		}
		return values, nil
	case []interface{}:
		values := make(lang.ListValue, len(v))
		for i, value := range v {
			values[i] = convertJSONToValue(value)
		}
		return values, nil
	default:
		return nil, fmt.Errorf("json_values: expected object or array, got %T", data)
	}
}

func jsonLength(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("json_length: expected 1 argument")
	}

	// First argument can be JSON string or already parsed data
	var data interface{}
	if jsonStr, ok := args[0].(lang.StringValue); ok {
		err := json.Unmarshal([]byte(string(jsonStr)), &data)
		if err != nil {
			return nil, fmt.Errorf("json_length: invalid JSON: %w", err)
		}
	} else {
		data = convertValueToJSON(args[0])
	}

	switch v := data.(type) {
	case map[string]interface{}:
		return lang.NumberValue(float64(len(v))), nil
	case []interface{}:
		return lang.NumberValue(float64(len(v))), nil
	case string:
		return lang.NumberValue(float64(len(v))), nil
	default:
		return nil, fmt.Errorf("json_length: cannot get length of %T", data)
	}
}

func jsonMerge(args []lang.Value) (lang.Value, error) {
	if len(args) < 2 {
		return nil, errors.New("json_merge: expected at least 2 arguments")
	}

	result := make(map[string]interface{})

	for i, arg := range args {
		var data interface{}
		if jsonStr, ok := arg.(lang.StringValue); ok {
			err := json.Unmarshal([]byte(string(jsonStr)), &data)
			if err != nil {
				return nil, fmt.Errorf("json_merge: argument %d invalid JSON: %w", i+1, err)
			}
		} else {
			data = convertValueToJSON(arg)
		}

		if obj, ok := data.(map[string]interface{}); ok {
			for key, value := range obj {
				result[key] = value
			}
		} else {
			return nil, fmt.Errorf("json_merge: argument %d is not an object", i+1)
		}
	}

	return convertJSONToValue(result), nil
}

func jsonType(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("json_type: expected 1 argument")
	}

	// First argument can be JSON string or already parsed data
	var data interface{}
	if jsonStr, ok := args[0].(lang.StringValue); ok {
		err := json.Unmarshal([]byte(string(jsonStr)), &data)
		if err != nil {
			return lang.StringValue("invalid"), nil // Return "invalid" instead of error for type check
		}
	} else {
		data = convertValueToJSON(args[0])
	}

	switch data.(type) {
	case nil:
		return lang.StringValue("null"), nil
	case bool:
		return lang.StringValue("boolean"), nil
	case float64:
		return lang.StringValue("number"), nil
	case string:
		return lang.StringValue("string"), nil
	case []interface{}:
		return lang.StringValue("array"), nil
	case map[string]interface{}:
		return lang.StringValue("object"), nil
	default:
		return lang.StringValue("unknown"), nil
	}
}

// Helper function to get value by JSON path (simple dot notation)
func jsonGetByPath(data interface{}, path string) lang.Value {
	if path == "" {
		return convertJSONToValue(data)
	}

	parts := strings.Split(path, ".")
	current := data

	for _, part := range parts {
		switch v := current.(type) {
		case map[string]interface{}:
			current = v[part]
		case []interface{}:
			// Handle array index
			if idx, err := strconv.Atoi(part); err == nil && idx >= 0 && idx < len(v) {
				current = v[idx]
			} else {
				return nil
			}
		default:
			return nil
		}

		if current == nil {
			return nil
		}
	}

	return convertJSONToValue(current)
}

// Helper function to set value by JSON path
func jsonSetByPath(data interface{}, path string, value interface{}) interface{} {
	if path == "" {
		return value
	}

	parts := strings.Split(path, ".")
	return jsonSetByPathParts(data, parts, value)
}

func jsonSetByPathParts(data interface{}, parts []string, value interface{}) interface{} {
	if len(parts) == 0 {
		return value
	}

	key := parts[0]

	if len(parts) == 1 {
		// Last part - set the value
		switch v := data.(type) {
		case map[string]interface{}:
			if v == nil {
				v = make(map[string]interface{})
			}
			v[key] = value
			return v
		case []interface{}:
			if idx, err := strconv.Atoi(key); err == nil && idx >= 0 {
				// Extend array if needed
				for len(v) <= idx {
					v = append(v, nil)
				}
				v[idx] = value
				return v
			}
		default:
			// Create new object
			newObj := make(map[string]interface{})
			newObj[key] = value
			return newObj
		}
	}

	// Recursive case
	switch v := data.(type) {
	case map[string]interface{}:
		if v == nil {
			v = make(map[string]interface{})
		}
		if v[key] == nil {
			// Determine if next part is array index
			nextKey := parts[1]
			if _, err := strconv.Atoi(nextKey); err == nil {
				v[key] = make([]interface{}, 0)
			} else {
				v[key] = make(map[string]interface{})
			}
		}
		v[key] = jsonSetByPathParts(v[key], parts[1:], value)
		return v
	case []interface{}:
		if idx, err := strconv.Atoi(key); err == nil && idx >= 0 {
			// Extend array if needed
			for len(v) <= idx {
				v = append(v, nil)
			}
			if v[idx] == nil {
				// Determine if next part is array index
				nextKey := parts[1]
				if _, err := strconv.Atoi(nextKey); err == nil {
					v[idx] = make([]interface{}, 0)
				} else {
					v[idx] = make(map[string]interface{})
				}
			}
			v[idx] = jsonSetByPathParts(v[idx], parts[1:], value)
			return v
		}
	default:
		// Create new structure
		newObj := make(map[string]interface{})
		// Determine if next part is array index
		if len(parts) > 1 {
			nextKey := parts[1]
			if _, err := strconv.Atoi(nextKey); err == nil {
				newObj[key] = make([]interface{}, 0)
			} else {
				newObj[key] = make(map[string]interface{})
			}
			newObj[key] = jsonSetByPathParts(newObj[key], parts[1:], value)
		} else {
			newObj[key] = value
		}
		return newObj
	}

	return data
}

// Helper function to delete value by JSON path
func jsonDeleteByPath(data interface{}, path string) interface{} {
	if path == "" {
		return nil
	}

	parts := strings.Split(path, ".")
	if len(parts) == 1 {
		// Single key - delete directly
		key := parts[0]
		switch v := data.(type) {
		case map[string]interface{}:
			delete(v, key)
			return v
		case []interface{}:
			if idx, err := strconv.Atoi(key); err == nil && idx >= 0 && idx < len(v) {
				// Remove element from array
				return append(v[:idx], v[idx+1:]...)
			}
		}
		return data
	}

	// Navigate to parent and delete the last key
	parentPath := strings.Join(parts[:len(parts)-1], ".")
	parent := jsonGetByPath(data, parentPath)
	if parent == nil {
		return data
	}

	lastKey := parts[len(parts)-1]
	parentJSON := convertValueToJSON(parent)

	switch v := parentJSON.(type) {
	case map[string]interface{}:
		delete(v, lastKey)
	case []interface{}:
		if idx, err := strconv.Atoi(lastKey); err == nil && idx >= 0 && idx < len(v) {
			parentJSON = append(v[:idx], v[idx+1:]...)
		}
	}

	// Set the modified parent back
	return jsonSetByPath(data, parentPath, parentJSON)
}

// Convert Go JSON types to our Value types
func convertJSONToValue(v interface{}) lang.Value {
	switch val := v.(type) {
	case nil:
		return nil
	case bool:
		return lang.BoolValue(val)
	case float64:
		return lang.NumberValue(val)
	case string:
		return lang.StringValue(val)
	case []interface{}:
		result := make(lang.ListValue, len(val))
		for i, item := range val {
			result[i] = convertJSONToValue(item)
		}
		return result
	case map[string]interface{}:
		result := lang.MapValue{}
		for k, item := range val {
			result[k] = convertJSONToValue(item)
		}
		return result
	default:
		return nil
	}
}

// Convert our Value types to Go JSON types
func convertValueToJSON(v lang.Value) interface{} {
	switch val := v.(type) {
	case nil:
		return nil
	case lang.BoolValue:
		return bool(val)
	case lang.NumberValue:
		return float64(val)
	case lang.StringValue:
		return string(val)
	case lang.ListValue:
		result := make([]interface{}, len(val))
		for i, item := range val {
			result[i] = convertValueToJSON(item)
		}
		return result
	case lang.MapValue:
		result := map[string]interface{}{}
		for k, item := range val {
			result[k] = convertValueToJSON(item)
		}
		return result
	default:
		return nil
	}
}

// Functions that would be in the BuiltinFunctions map:
var JsonFunctions = map[string]func([]lang.Value) (lang.Value, error){
	"json_parse":  jsonParse,
	"json_string": jsonString,
	"json_valid":  jsonValid,
	"json_get":    jsonGet,
	"json_set":    jsonSet,
	"json_delete": jsonDelete,
	"json_has":    jsonHas,
	"json_keys":   jsonKeys,
	"json_values": jsonValues,
	"json_length": jsonLength,
	"json_merge":  jsonMerge,
	"json_type":   jsonType,
}
