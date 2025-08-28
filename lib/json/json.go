package json

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/vedadiyan/exql/lang"
	"github.com/vedadiyan/exql/lib"
)

func parse() (string, lang.Function) {
	name := "json_parse"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		var result interface{}
		err = json.Unmarshal([]byte(string(str)), &result)
		if err != nil {
			return nil, fmt.Errorf("%s: invalid JSON: %w", name, err)
		}
		return convertJSONToValue(result), nil
	}
	return name, fn
}

func sstring() (string, lang.Function) {
	name := "json_string"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) < 1 || len(args) > 2 {
			return nil, fmt.Errorf("%s: expected 1 or 2 arguments (data, pretty?)", name)
		}
		indent := ""
		if len(args) == 2 {
			pretty, err := lib.ToBool(args[1])
			if err != nil {
				return nil, fmt.Errorf("%s: pretty flag %w", name, err)
			}
			if pretty {
				indent = "  "
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
			return nil, fmt.Errorf("%s: marshalling failed: %w", name, err)
		}
		return lang.StringValue(string(jsonBytes)), nil
	}
	return name, fn
}

func valid() (string, lang.Function) {
	name := "json_valid"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		var result interface{}
		err = json.Unmarshal([]byte(string(str)), &result)
		return lang.BoolValue(err == nil), nil
	}
	return name, fn
}

func get() (string, lang.Function) {
	name := "json_get"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 2 {
			return nil, lib.ArgumentError(name, 2)
		}
		var data interface{}
		if jsonStr, ok := args[0].(lang.StringValue); ok {
			err := json.Unmarshal([]byte(string(jsonStr)), &data)
			if err != nil {
				return nil, fmt.Errorf("%s: invalid JSON: %w", name, err)
			}
		} else {
			data = convertValueToJSON(args[0])
		}
		pathStr, err := lib.ToString(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: path %w", name, err)
		}
		path := string(pathStr)
		result := jsonGetByPath(data, path)
		return result, nil
	}
	return name, fn
}

func set() (string, lang.Function) {
	name := "json_set"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 3 {
			return nil, lib.ArgumentError(name, 3)
		}
		var data interface{}
		if jsonStr, ok := args[0].(lang.StringValue); ok {
			err := json.Unmarshal([]byte(string(jsonStr)), &data)
			if err != nil {
				return nil, fmt.Errorf("%s: invalid JSON: %w", name, err)
			}
		} else {
			data = convertValueToJSON(args[0])
		}
		pathStr, err := lib.ToString(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: path %w", name, err)
		}
		path := string(pathStr)
		value := convertValueToJSON(args[2])
		result := jsonSetByPath(data, path, value)
		return convertJSONToValue(result), nil
	}
	return name, fn
}

func remove() (string, lang.Function) {
	name := "json_delete"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 2 {
			return nil, lib.ArgumentError(name, 2)
		}
		var data interface{}
		if jsonStr, ok := args[0].(lang.StringValue); ok {
			err := json.Unmarshal([]byte(string(jsonStr)), &data)
			if err != nil {
				return nil, fmt.Errorf("%s: invalid JSON: %w", name, err)
			}
		} else {
			data = convertValueToJSON(args[0])
		}
		pathStr, err := lib.ToString(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: path %w", name, err)
		}
		path := string(pathStr)
		result := jsonDeleteByPath(data, path)
		return convertJSONToValue(result), nil
	}
	return name, fn
}

func has() (string, lang.Function) {
	name := "json_has"
	_, Get := get()
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 2 {
			return nil, lib.ArgumentError(name, 2)
		}
		result, err := Get(args)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		return lang.BoolValue(result != nil), nil
	}
	return name, fn
}

func keys() (string, lang.Function) {
	name := "json_keys"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		var data interface{}
		if jsonStr, ok := args[0].(lang.StringValue); ok {
			err := json.Unmarshal([]byte(string(jsonStr)), &data)
			if err != nil {
				return nil, fmt.Errorf("%s: invalid JSON: %w", name, err)
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
			return nil, fmt.Errorf("%s: expected object, got %T", name, data)
		}
	}
	return name, fn
}

func values() (string, lang.Function) {
	name := "json_values"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		var data interface{}
		if jsonStr, ok := args[0].(lang.StringValue); ok {
			err := json.Unmarshal([]byte(string(jsonStr)), &data)
			if err != nil {
				return nil, fmt.Errorf("%s: invalid JSON: %w", name, err)
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
			return nil, fmt.Errorf("%s: expected object or array, got %T", name, data)
		}
	}
	return name, fn
}

func length() (string, lang.Function) {
	name := "json_length"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		var data interface{}
		if jsonStr, ok := args[0].(lang.StringValue); ok {
			err := json.Unmarshal([]byte(string(jsonStr)), &data)
			if err != nil {
				return nil, fmt.Errorf("%s: invalid JSON: %w", name, err)
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
			return nil, fmt.Errorf("%s: cannot get length of %T", name, data)
		}
	}
	return name, fn
}

func merge() (string, lang.Function) {
	name := "json_merge"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("%s: expected at least 2 arguments", name)
		}
		result := make(map[string]interface{})
		for i, arg := range args {
			var data interface{}
			if jsonStr, ok := arg.(lang.StringValue); ok {
				err := json.Unmarshal([]byte(string(jsonStr)), &data)
				if err != nil {
					return nil, fmt.Errorf("%s: argument %d invalid JSON: %w", name, i+1, err)
				}
			} else {
				data = convertValueToJSON(arg)
			}
			if obj, ok := data.(map[string]interface{}); ok {
				for key, value := range obj {
					result[key] = value
				}
			} else {
				return nil, fmt.Errorf("%s: argument %d is not an object", name, i+1)
			}
		}
		return convertJSONToValue(result), nil
	}
	return name, fn
}

func ttype() (string, lang.Function) {
	name := "json_type"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		var data interface{}
		if jsonStr, ok := args[0].(lang.StringValue); ok {
			err := json.Unmarshal([]byte(string(jsonStr)), &data)
			if err != nil {
				return lang.StringValue("invalid"), nil
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
	return name, fn
}

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
		switch v := data.(type) {
		case map[string]interface{}:
			if v == nil {
				v = make(map[string]interface{})
			}
			v[key] = value
			return v
		case []interface{}:
			if idx, err := strconv.Atoi(key); err == nil && idx >= 0 {
				for len(v) <= idx {
					v = append(v, nil)
				}
				v[idx] = value
				return v
			}
		default:
			newObj := make(map[string]interface{})
			newObj[key] = value
			return newObj
		}
	}
	switch v := data.(type) {
	case map[string]interface{}:
		if v == nil {
			v = make(map[string]interface{})
		}
		if v[key] == nil {
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
			for len(v) <= idx {
				v = append(v, nil)
			}
			if v[idx] == nil {
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
		newObj := make(map[string]interface{})
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

func jsonDeleteByPath(data interface{}, path string) interface{} {
	if path == "" {
		return nil
	}
	parts := strings.Split(path, ".")
	if len(parts) == 1 {
		key := parts[0]
		switch v := data.(type) {
		case map[string]interface{}:
			delete(v, key)
			return v
		case []interface{}:
			if idx, err := strconv.Atoi(key); err == nil && idx >= 0 && idx < len(v) {
				return append(v[:idx], v[idx+1:]...)
			}
		}
		return data
	}
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
	return jsonSetByPath(data, parentPath, parentJSON)
}

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

var JsonFunctions = []func() (string, lang.Function){
	parse,
	sstring,
	valid,
	get,
	set,
	remove,
	has,
	keys,
	values,
	length,
	merge,
	ttype,
}
