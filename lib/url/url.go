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
package url

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/vedadiyan/exql/lang"
	"github.com/vedadiyan/exql/lib"
)

func parse() (string, lang.Function) {
	name := "parse"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, lib.StringError(name, args[0])
		}
		u, err := url.Parse(string(str))
		if err != nil {
			return nil, fmt.Errorf("%s: invalid URL string: %w", name, err)
		}

		host := u.Host
		port := ""
		if colonIdx := strings.LastIndex(host, ":"); colonIdx != -1 {
			port = host[colonIdx+1:]
			host = host[:colonIdx]
		}

		return lang.MapValue{
			"scheme":   lang.StringValue(u.Scheme),
			"host":     lang.StringValue(host),
			"port":     lang.StringValue(port),
			"path":     lang.StringValue(u.Path),
			"query":    lang.StringValue(u.RawQuery),
			"fragment": lang.StringValue(u.Fragment),
			"user":     lang.StringValue(getUserInfo(u.User)),
		}, nil
	}
	return name, fn
}

func urlEncode() (string, lang.Function) {
	name := "encode"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, lib.StringError(name, args[0])
		}
		return lang.StringValue(url.QueryEscape(string(str))), nil
	}
	return name, fn
}

func urlDecode() (string, lang.Function) {
	name := "decode"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, lib.StringError(name, args[0])
		}
		decoded, err := url.QueryUnescape(string(str))
		if err != nil {
			return nil, fmt.Errorf("%s: invalid URL string: %w", name, err)
		}
		return lang.StringValue(decoded), nil
	}
	return name, fn
}

func urlHost() (string, lang.Function) {
	name := "host"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, lib.StringError(name, args[0])
		}
		u, err := url.Parse(string(str))
		if err != nil {
			return nil, fmt.Errorf("%s: invalid URL string: %w", name, err)
		}
		host := u.Host
		if colonIdx := strings.LastIndex(host, ":"); colonIdx != -1 {
			host = host[:colonIdx]
		}
		return lang.StringValue(host), nil
	}
	return name, fn
}

func port() (string, lang.Function) {
	name := "port"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, lib.StringError(name, args[0])
		}
		u, err := url.Parse(string(str))
		if err != nil {
			return nil, fmt.Errorf("%s: invalid URL string: %w", name, err)
		}

		if colonIdx := strings.LastIndex(u.Host, ":"); colonIdx != -1 {
			portStr := u.Host[colonIdx+1:]
			if port, err := strconv.Atoi(portStr); err == nil {
				return lang.NumberValue(float64(port)), nil
			}
			return nil, fmt.Errorf("%s: invalid port number in URL: '%s'", name, portStr)
		}

		switch strings.ToLower(u.Scheme) {
		case "https":
			return lang.NumberValue(443), nil
		case "http":
			return lang.NumberValue(80), nil
		case "ftp":
			return lang.NumberValue(21), nil
		case "ssh":
			return lang.NumberValue(22), nil
		default:
			return lang.NumberValue(0), nil
		}
	}
	return name, fn
}

func path() (string, lang.Function) {
	name := "path"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, lib.StringError(name, args[0])
		}
		u, err := url.Parse(string(str))
		if err != nil {
			return nil, fmt.Errorf("%s: invalid URL string: %w", name, err)
		}
		return lang.StringValue(u.Path), nil
	}
	return name, fn
}

func query() (string, lang.Function) {
	name := "query"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, lib.StringError(name, args[0])
		}
		u, err := url.Parse(string(str))
		if err != nil {
			return nil, fmt.Errorf("%s: invalid URL string: %w", name, err)
		}

		result := lang.MapValue{}
		for k, v := range u.Query() {
			if len(v) == 1 {
				result[k] = lang.StringValue(v[0])
			} else {
				list := make(lang.ListValue, len(v))
				for i, val := range v {
					list[i] = lang.StringValue(val)
				}
				result[k] = list
			}
		}
		return result, nil
	}
	return name, fn
}

func queryParam() (string, lang.Function) {
	name := "query_param"
	_, query := query()
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 2 {
			return nil, lib.ArgumentError(name, 2)
		}
		queryParams, err := query(args[:1])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		paramName, err := lib.ToString(args[1])
		if err != nil {
			return nil, lib.StringError(name, args[1])
		}

		if queryMap, ok := queryParams.(lang.MapValue); ok {
			return queryMap[string(paramName)], nil
		}
		return nil, fmt.Errorf("%s: internal error parsing URL query", name)
	}
	return name, fn
}

func fragment() (string, lang.Function) {
	name := "fragment"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, lib.StringError(name, args[0])
		}
		u, err := url.Parse(string(str))
		if err != nil {
			return nil, fmt.Errorf("%s: invalid URL string: %w", name, err)
		}
		return lang.StringValue(u.Fragment), nil
	}
	return name, fn
}

func scheme() (string, lang.Function) {
	name := "scheme"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, lib.StringError(name, args[0])
		}
		u, err := url.Parse(string(str))
		if err != nil {
			return nil, fmt.Errorf("%s: invalid URL string: %w", name, err)
		}
		return lang.StringValue(u.Scheme), nil
	}
	return name, fn
}

func user() (string, lang.Function) {
	name := "user"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, lib.StringError(name, args[0])
		}
		u, err := url.Parse(string(str))
		if err != nil {
			return nil, fmt.Errorf("%s: invalid URL string: %w", name, err)
		}
		return lang.StringValue(getUserInfo(u.User)), nil
	}
	return name, fn
}

func build() (string, lang.Function) {
	name := "build"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}

		parts, ok := args[0].(lang.MapValue)
		if !ok {
			return nil, lib.MapError(name, args[0])
		}

		var u url.URL

		if scheme, exists := parts["scheme"]; exists {
			s, err := lib.ToString(scheme)
			if err != nil {
				return nil, fmt.Errorf("%s: invalid scheme value: %w", name, err)
			}
			u.Scheme = string(s)
		}

		if host, exists := parts["host"]; exists {
			hostStr, err := lib.ToString(host)
			if err != nil {
				return nil, fmt.Errorf("%s: invalid host value: %w", name, err)
			}
			h := string(hostStr)
			if port, exists := parts["port"]; exists {
				portStr, err := lib.ToString(port)
				if err != nil {
					return nil, fmt.Errorf("%s: invalid port value: %w", name, err)
				}
				p := string(portStr)
				if p != "" && p != "0" {
					h = h + ":" + p
				}
			}
			u.Host = h
		}

		if path, exists := parts["path"]; exists {
			p, err := lib.ToString(path)
			if err != nil {
				return nil, fmt.Errorf("%s: invalid path value: %w", name, err)
			}
			u.Path = string(p)
		}

		if query, exists := parts["query"]; exists {
			q, err := lib.ToString(query)
			if err != nil {
				return nil, fmt.Errorf("%s: invalid query value: %w", name, err)
			}
			u.RawQuery = string(q)
		}

		if fragment, exists := parts["fragment"]; exists {
			f, err := lib.ToString(fragment)
			if err != nil {
				return nil, fmt.Errorf("%s: invalid fragment value: %w", name, err)
			}
			u.Fragment = string(f)
		}

		if user, exists := parts["user"]; exists {
			userStr, err := lib.ToString(user)
			if err != nil {
				return nil, fmt.Errorf("%s: invalid user value: %w", name, err)
			}
			uStr := string(userStr)
			if uStr != "" {
				if password, exists := parts["password"]; exists {
					passStr, err := lib.ToString(password)
					if err != nil {
						return nil, fmt.Errorf("%s: invalid password value: %w", name, err)
					}
					u.User = url.UserPassword(uStr, string(passStr))
				} else {
					u.User = url.User(uStr)
				}
			}
		}

		return lang.StringValue(u.String()), nil
	}
	return name, fn
}

func join() (string, lang.Function) {
	name := "join"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) < 2 {
			return nil, fmt.Errorf("%s: expected at least 2 arguments", name)
		}

		base, err := lib.ToString(args[0])
		if err != nil {
			return nil, fmt.Errorf("%s: base URL %w", name, err)
		}
		u, err := url.Parse(string(base))
		if err != nil {
			return nil, fmt.Errorf("%s: invalid base URL string: %w", name, err)
		}

		for i := 1; i < len(args); i++ {
			pathSegment, err := lib.ToString(args[i])
			if err != nil {
				return nil, fmt.Errorf("%s: path segment %d %w", name, i, err)
			}
			if string(pathSegment) == "" {
				continue
			}
			u = u.JoinPath(string(pathSegment))
		}

		return lang.StringValue(u.String()), nil
	}
	return name, fn
}

func isAbsolute() (string, lang.Function) {
	name := "is_absolute"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}

		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, lib.StringError(name, args[0])
		}

		u, err := url.Parse(string(str))
		if err != nil {
			return nil, fmt.Errorf("%s: invalid URL string: %w", name, err)
		}

		return lang.BoolValue(u.IsAbs()), nil
	}
	return name, fn
}

func pathSegments() (string, lang.Function) {
	name := "path_segments"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}

		path, err := lib.ToString(args[0])
		if err != nil {
			return nil, lib.StringError(name, args[0])
		}
		pathStr := string(path)

		if strings.Contains(pathStr, "://") {
			u, err := url.Parse(pathStr)
			if err != nil {
				return nil, fmt.Errorf("%s: invalid URL string: %w", name, err)
			}
			pathStr = u.Path
		}

		pathStr = strings.Trim(pathStr, "/")
		if pathStr == "" {
			return lang.ListValue{}, nil
		}

		segments := strings.Split(pathStr, "/")
		result := make(lang.ListValue, len(segments))
		for i, segment := range segments {
			decoded, err := url.QueryUnescape(segment)
			if err != nil {
				return nil, fmt.Errorf("%s: invalid path segment '%s': %w", name, segment, err)
			}
			result[i] = lang.StringValue(decoded)
		}

		return result, nil
	}
	return name, fn
}

func queryString() (string, lang.Function) {
	name := "query_string"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}

		params, ok := args[0].(lang.MapValue)
		if !ok {
			return nil, lib.MapError(name, args[0])
		}

		values := url.Values{}
		for key, value := range params {
			switch v := value.(type) {
			case lang.ListValue:
				for _, item := range v {
					itemStr, err := lib.ToString(item)
					if err != nil {
						return nil, fmt.Errorf("%s: list element %w", name, err)
					}
					values.Add(key, string(itemStr))
				}
			default:
				valStr, err := lib.ToString(value)
				if err != nil {
					return nil, fmt.Errorf("%s: map value for key '%s' %w", name, key, err)
				}
				values.Set(key, string(valStr))
			}
		}

		return lang.StringValue(values.Encode()), nil
	}
	return name, fn
}

func clean() (string, lang.Function) {
	name := "clean"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}

		str, err := lib.ToString(args[0])
		if err != nil {
			return nil, lib.StringError(name, args[0])
		}
		u, err := url.Parse(string(str))
		if err != nil {
			return nil, fmt.Errorf("%s: invalid URL string: %w", name, err)
		}

		u.Path = cleanpath(u.Path)

		return lang.StringValue(u.String()), nil
	}
	return name, fn
}

func getUserInfo(userInfo *url.Userinfo) string {
	if userInfo == nil {
		return ""
	}
	return userInfo.Username()
}

func cleanpath(path string) string {
	if path == "" {
		return "/"
	}

	segments := strings.Split(path, "/")
	cleaned := make([]string, 0, len(segments))

	for _, segment := range segments {
		switch segment {
		case "", ".":
			continue
		case "..":
			if len(cleaned) > 0 && cleaned[len(cleaned)-1] != ".." {
				cleaned = cleaned[:len(cleaned)-1]
			} else if !strings.HasPrefix(path, "/") {
				cleaned = append(cleaned, segment)
			}
		default:
			cleaned = append(cleaned, segment)
		}
	}

	result := strings.Join(cleaned, "/")
	if strings.HasPrefix(path, "/") && !strings.HasPrefix(result, "/") {
		result = "/" + result
	}
	if result == "" {
		result = "/"
	}

	return result
}

var functions = []func() (string, lang.Function){
	parse,
	urlEncode,
	urlDecode,
	urlHost,
	port,
	path,
	query,
	queryParam,
	fragment,
	scheme,
	user,
	build,
	join,
	isAbsolute,
	pathSegments,
	queryString,
	clean,
}

func Export() map[string]lang.Function {
	out := make(map[string]lang.Function)
	for _, value := range functions {
		name, fn := value()
		out[name] = fn
	}
	return out
}
