package url

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/vedadiyan/exql/lang"
	"github.com/vedadiyan/exql/lib"
)

// URL/URI Processing Functions
// These functions help parse, manipulate, and extract information from URLs

func urlParse(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("url_parse: expected 1 argument")
	}
	urlStr, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("url_parse: %w", err)
	}
	u, err := url.Parse(string(urlStr))
	if err != nil {
		return nil, fmt.Errorf("url_parse: invalid URL string: %w", err)
	}

	// Parse port from host
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

func urlEncode(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("url_encode: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("url_encode: %w", err)
	}
	return lang.StringValue(url.QueryEscape(string(str))), nil
}

func urlDecode(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("url_decode: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("url_decode: %w", err)
	}
	decoded, err := url.QueryUnescape(string(str))
	if err != nil {
		return nil, fmt.Errorf("url_decode: invalid URL string: %w", err)
	}
	return lang.StringValue(decoded), nil
}

func urlHost(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("url_host: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("url_host: %w", err)
	}
	u, err := url.Parse(string(str))
	if err != nil {
		return nil, fmt.Errorf("url_host: invalid URL string: %w", err)
	}
	// Remove port if present
	host := u.Host
	if colonIdx := strings.LastIndex(host, ":"); colonIdx != -1 {
		host = host[:colonIdx]
	}
	return lang.StringValue(host), nil
}

func urlPort(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("url_port: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("url_port: %w", err)
	}
	u, err := url.Parse(string(str))
	if err != nil {
		return nil, fmt.Errorf("url_port: invalid URL string: %w", err)
	}

	// Extract port from host
	if colonIdx := strings.LastIndex(u.Host, ":"); colonIdx != -1 {
		portStr := u.Host[colonIdx+1:]
		if port, err := strconv.Atoi(portStr); err == nil {
			return lang.NumberValue(float64(port)), nil
		}
		// If port is not a number, we return an error instead of 0
		return nil, fmt.Errorf("url_port: invalid port number in URL: '%s'", portStr)
	}

	// Return default port based on scheme
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

func urlPath(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("url_path: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("url_path: %w", err)
	}
	u, err := url.Parse(string(str))
	if err != nil {
		return nil, fmt.Errorf("url_path: invalid URL string: %w", err)
	}
	return lang.StringValue(u.Path), nil
}

func urlQuery(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("url_query: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("url_query: %w", err)
	}
	u, err := url.Parse(string(str))
	if err != nil {
		return nil, fmt.Errorf("url_query: invalid URL string: %w", err)
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

func urlQueryParam(args []lang.Value) (lang.Value, error) {
	if len(args) != 2 {
		return nil, errors.New("url_query_param: expected 2 arguments")
	}
	queryParams, err := urlQuery(args[:1])
	if err != nil {
		return nil, fmt.Errorf("url_query_param: %w", err)
	}
	paramName, err := lib.ToString(args[1])
	if err != nil {
		return nil, fmt.Errorf("url_query_param: %w", err)
	}

	if queryMap, ok := queryParams.(lang.MapValue); ok {
		return queryMap[string(paramName)], nil
	}
	return nil, errors.New("url_query_param: internal error parsing URL query")
}

func urlFragment(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("url_fragment: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("url_fragment: %w", err)
	}
	u, err := url.Parse(string(str))
	if err != nil {
		return nil, fmt.Errorf("url_fragment: invalid URL string: %w", err)
	}
	return lang.StringValue(u.Fragment), nil
}

func urlScheme(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("url_scheme: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("url_scheme: %w", err)
	}
	u, err := url.Parse(string(str))
	if err != nil {
		return nil, fmt.Errorf("url_scheme: invalid URL string: %w", err)
	}
	return lang.StringValue(u.Scheme), nil
}

func urlUser(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("url_user: expected 1 argument")
	}
	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("url_user: %w", err)
	}
	u, err := url.Parse(string(str))
	if err != nil {
		return nil, fmt.Errorf("url_user: invalid URL string: %w", err)
	}
	return lang.StringValue(getUserInfo(u.User)), nil
}

func urlBuild(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("url_build: expected 1 argument")
	}

	parts, ok := args[0].(lang.MapValue)
	if !ok {
		return nil, fmt.Errorf("url_build: expected a map, got %T", args[0])
	}

	var u url.URL

	if scheme, exists := parts["scheme"]; exists {
		s, err := lib.ToString(scheme)
		if err != nil {
			return nil, fmt.Errorf("url_build: invalid scheme value: %w", err)
		}
		u.Scheme = string(s)
	}

	if host, exists := parts["host"]; exists {
		hostStr, err := lib.ToString(host)
		if err != nil {
			return nil, fmt.Errorf("url_build: invalid host value: %w", err)
		}
		h := string(hostStr)
		if port, exists := parts["port"]; exists {
			portStr, err := lib.ToString(port)
			if err != nil {
				return nil, fmt.Errorf("url_build: invalid port value: %w", err)
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
			return nil, fmt.Errorf("url_build: invalid path value: %w", err)
		}
		u.Path = string(p)
	}

	if query, exists := parts["query"]; exists {
		q, err := lib.ToString(query)
		if err != nil {
			return nil, fmt.Errorf("url_build: invalid query value: %w", err)
		}
		u.RawQuery = string(q)
	}

	if fragment, exists := parts["fragment"]; exists {
		f, err := lib.ToString(fragment)
		if err != nil {
			return nil, fmt.Errorf("url_build: invalid fragment value: %w", err)
		}
		u.Fragment = string(f)
	}

	if user, exists := parts["user"]; exists {
		userStr, err := lib.ToString(user)
		if err != nil {
			return nil, fmt.Errorf("url_build: invalid user value: %w", err)
		}
		uStr := string(userStr)
		if uStr != "" {
			if password, exists := parts["password"]; exists {
				passStr, err := lib.ToString(password)
				if err != nil {
					return nil, fmt.Errorf("url_build: invalid password value: %w", err)
				}
				u.User = url.UserPassword(uStr, string(passStr))
			} else {
				u.User = url.User(uStr)
			}
		}
	}

	return lang.StringValue(u.String()), nil
}

func urlJoin(args []lang.Value) (lang.Value, error) {
	if len(args) < 2 {
		return nil, errors.New("url_join: expected at least 2 arguments")
	}

	base, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("url_join: base URL %w", err)
	}
	u, err := url.Parse(string(base))
	if err != nil {
		return nil, fmt.Errorf("url_join: invalid base URL string: %w", err)
	}

	for i := 1; i < len(args); i++ {
		pathSegment, err := lib.ToString(args[i])
		if err != nil {
			return nil, fmt.Errorf("url_join: path segment %d %w", i, err)
		}
		if string(pathSegment) == "" {
			continue
		}

		ref, err := url.Parse(string(pathSegment))
		if err != nil {
			return nil, fmt.Errorf("url_join: invalid path segment: %w", err)
		}

		u = u.ResolveReference(ref)
	}

	return lang.StringValue(u.String()), nil
}

func urlIsAbsolute(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("url_is_absolute: expected 1 argument")
	}

	str, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("url_is_absolute: %w", err)
	}

	u, err := url.Parse(string(str))
	if err != nil {
		return nil, fmt.Errorf("url_is_absolute: invalid URL string: %w", err)
	}

	return lang.BoolValue(u.IsAbs()), nil
}

func urlPathSegments(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("url_path_segments: expected 1 argument")
	}

	path, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("url_path_segments: %w", err)
	}
	pathStr := string(path)

	// If it's a full URL, extract just the path
	if strings.Contains(pathStr, "://") {
		u, err := url.Parse(pathStr)
		if err != nil {
			return nil, fmt.Errorf("url_path_segments: invalid URL string: %w", err)
		}
		pathStr = u.Path
	}

	// Clean and split path
	pathStr = strings.Trim(pathStr, "/")
	if pathStr == "" {
		return lang.ListValue{}, nil
	}

	segments := strings.Split(pathStr, "/")
	result := make(lang.ListValue, len(segments))
	for i, segment := range segments {
		decoded, err := url.QueryUnescape(segment)
		if err != nil {
			// If unescaping fails, we return an error instead of continuing with the unescaped string.
			return nil, fmt.Errorf("url_path_segments: invalid path segment '%s': %w", segment, err)
		}
		result[i] = lang.StringValue(decoded)
	}

	return result, nil
}

func urlQueryString(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("url_query_string: expected 1 argument")
	}

	params, ok := args[0].(lang.MapValue)
	if !ok {
		return nil, fmt.Errorf("url_query_string: expected a map, got %T", args[0])
	}

	values := url.Values{}
	for key, value := range params {
		switch v := value.(type) {
		case lang.ListValue:
			for _, item := range v {
				itemStr, err := lib.ToString(item)
				if err != nil {
					return nil, fmt.Errorf("url_query_string: list element %w", err)
				}
				values.Add(key, string(itemStr))
			}
		default:
			valStr, err := lib.ToString(value)
			if err != nil {
				return nil, fmt.Errorf("url_query_string: map value for key '%s' %w", key, err)
			}
			values.Set(key, string(valStr))
		}
	}

	return lang.StringValue(values.Encode()), nil
}

func urlClean(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("url_clean: expected 1 argument")
	}

	urlStr, err := lib.ToString(args[0])
	if err != nil {
		return nil, fmt.Errorf("url_clean: %w", err)
	}
	u, err := url.Parse(string(urlStr))
	if err != nil {
		return nil, fmt.Errorf("url_clean: invalid URL string: %w", err)
	}

	// Clean the path
	u.Path = cleanURLPath(u.Path)

	return lang.StringValue(u.String()), nil
}

// Helper function to get user info as string
func getUserInfo(userInfo *url.Userinfo) string {
	if userInfo == nil {
		return ""
	}
	return userInfo.Username()
}

// Helper function to clean URL paths (similar to path.Clean but for URLs)
func cleanURLPath(path string) string {
	if path == "" {
		return "/"
	}

	// Split path into segments
	segments := strings.Split(path, "/")
	cleaned := make([]string, 0, len(segments))

	for _, segment := range segments {
		switch segment {
		case "", ".":
			// Skip empty and current directory
			continue
		case "..":
			// Parent directory - remove last segment if possible
			if len(cleaned) > 0 && cleaned[len(cleaned)-1] != ".." {
				cleaned = cleaned[:len(cleaned)-1]
			} else if !strings.HasPrefix(path, "/") {
				// Only keep ".." for relative paths
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

var UrlFunctions = map[string]func(args []lang.Value) (lang.Value, error){
	"url_parse":         urlParse,
	"url_encode":        urlEncode,
	"url_decode":        urlDecode,
	"url_host":          urlHost,
	"url_port":          urlPort,
	"url_path":          urlPath,
	"url_query":         urlQuery,
	"url_query_param":   urlQueryParam,
	"url_fragment":      urlFragment,
	"url_scheme":        urlScheme,
	"url_user":          urlUser,
	"url_build":         urlBuild,
	"url_join":          urlJoin,
	"url_is_absolute":   urlIsAbsolute,
	"url_path_segments": urlPathSegments,
	"url_query_string":  urlQueryString,
	"url_clean":         urlClean,
}
