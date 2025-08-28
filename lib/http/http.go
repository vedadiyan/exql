package http

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/vedadiyan/exql/lang"
	"github.com/vedadiyan/exql/lib"
)

func getHeaderValue(headers lang.MapValue, headerName string) string {
	for key, value := range headers {
		if strings.EqualFold(key, headerName) {
			if str, err := lib.ToString(value); err == nil {
				return string(str)
			}
		}
	}
	return ""
}

func parseCookies(cookieHeader string) lang.MapValue {
	result := lang.MapValue{}
	cookies := strings.Split(cookieHeader, ";")
	for _, cookie := range cookies {
		parts := strings.SplitN(strings.TrimSpace(cookie), "=", 2)
		if len(parts) == 2 {
			result[parts[0]] = lang.StringValue(parts[1])
		}
	}
	return result
}

func header() (string, lang.Function) {
	name := "header"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 2 {
			return nil, lib.ArgumentError(name, 2)
		}
		contextMap, ok := args[0].(lang.MapValue)
		if !ok {
			return nil, lib.ContextError(name, args[0])
		}
		headerName, err := lib.ToString(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: header name %w", name, err)
		}

		if headers, exists := contextMap["headers"]; exists {
			if headersMap, ok := headers.(lang.MapValue); ok {
				for key, value := range headersMap {
					if strings.EqualFold(key, string(headerName)) {
						return value, nil
					}
				}
			}
		}
		return nil, nil
	}
	return name, fn
}

func headers() (string, lang.Function) {
	name := "headers"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		contextMap, ok := args[0].(lang.MapValue)
		if !ok {
			return nil, lib.ContextError(name, args[0])
		}
		if headers, exists := contextMap["headers"]; exists {
			if headersMap, ok := headers.(lang.MapValue); ok {
				return headersMap, nil
			}
		}
		return lang.MapValue{}, nil
	}
	return name, fn
}

func method() (string, lang.Function) {
	name := "method"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		contextMap, ok := args[0].(lang.MapValue)
		if !ok {
			return nil, lib.ContextError(name, args[0])
		}
		if method, exists := contextMap["method"]; exists {
			return method, nil
		}
		return lang.StringValue(""), nil
	}
	return name, fn
}

func path() (string, lang.Function) {
	name := "path"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		contextMap, ok := args[0].(lang.MapValue)
		if !ok {
			return nil, lib.ContextError(name, args[0])
		}
		if path, exists := contextMap["path"]; exists {
			return path, nil
		}
		return lang.StringValue(""), nil
	}
	return name, fn
}

func query() (string, lang.Function) {
	name := "query"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		contextMap, ok := args[0].(lang.MapValue)
		if !ok {
			return nil, lib.ContextError(name, args[0])
		}
		if query, exists := contextMap["query"]; exists {
			if queryMap, ok := query.(lang.MapValue); ok {
				return queryMap, nil
			}
		}
		return lang.MapValue{}, nil
	}
	return name, fn
}

func queryParam() (string, lang.Function) {
	name := "query_param"
	_, Query := query()
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 2 {
			return nil, lib.ArgumentError(name, 2)
		}
		queryParams, err := Query(args[:1])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		paramName, err := lib.ToString(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: param name %w", name, err)
		}

		if queryMap, ok := queryParams.(lang.MapValue); ok {
			return queryMap[string(paramName)], nil
		}
		return nil, nil
	}
	return name, fn
}

func body() (string, lang.Function) {
	name := "body"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		contextMap, ok := args[0].(lang.MapValue)
		if !ok {
			return nil, lib.ContextError(name, args[0])
		}
		if body, exists := contextMap["body"]; exists {
			return body, nil
		}
		return lang.StringValue(""), nil
	}
	return name, fn
}

func status() (string, lang.Function) {
	name := "status"
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		contextMap, ok := args[0].(lang.MapValue)
		if !ok {
			return nil, lib.ContextError(name, args[0])
		}
		if status, exists := contextMap["status"]; exists {
			return status, nil
		}
		return lang.NumberValue(0), nil
	}
	return name, fn
}

func ip() (string, lang.Function) {
	name := "ip"
	_, Headers := headers()
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		contextMap, ok := args[0].(lang.MapValue)
		if !ok {
			return nil, lib.ContextError(name, args[0])
		}

		ipFields := []string{"remote_ip", "client_ip", "x_forwarded_for", "x_real_ip", "ip"}

		for _, field := range ipFields {
			if ip, exists := contextMap[field]; exists {
				ipStr, err := lib.ToString(ip)
				if err == nil && string(ipStr) != "" {
					if strings.Contains(string(ipStr), ",") {
						return lang.StringValue(strings.TrimSpace(strings.Split(string(ipStr), ",")[0])), nil
					}
					return lang.StringValue(string(ipStr)), nil
				}
			}
		}
		headers, err := Headers(args)
		if err != nil {
			return lang.StringValue(""), nil
		}
		if headersMap, ok := headers.(lang.MapValue); ok {
			xForwardedFor := getHeaderValue(headersMap, "X-Forwarded-For")
			if xForwardedFor != "" {
				return lang.StringValue(strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])), nil
			}
			xRealIP := getHeaderValue(headersMap, "X-Real-IP")
			if xRealIP != "" {
				return lang.StringValue(xRealIP), nil
			}
		}
		return lang.StringValue(""), nil
	}
	return name, fn
}

func userAgent() (string, lang.Function) {
	name := "user_agent"
	_, Headers := headers()
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		headers, err := Headers(args)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		if headersMap, ok := headers.(lang.MapValue); ok {
			return lang.StringValue(getHeaderValue(headersMap, "User-Agent")), nil
		}
		return lang.StringValue(""), nil
	}
	return name, fn
}

func contentType() (string, lang.Function) {
	name := "content_type"
	_, Headers := headers()
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		headers, err := Headers(args)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		if headersMap, ok := headers.(lang.MapValue); ok {
			ct := getHeaderValue(headersMap, "Content-Type")
			if idx := strings.Index(ct, ";"); idx != -1 {
				ct = strings.TrimSpace(ct[:idx])
			}
			return lang.StringValue(ct), nil
		}
		return lang.StringValue(""), nil
	}
	return name, fn
}

func contentLength() (string, lang.Function) {
	name := "content_length"
	_, Headers := headers()
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		headers, err := Headers(args)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		if headersMap, ok := headers.(lang.MapValue); ok {
			cl := getHeaderValue(headersMap, "Content-Length")
			if cl == "" {
				return lang.NumberValue(0), nil
			}
			if length, err := strconv.ParseFloat(cl, 64); err == nil {
				return lang.NumberValue(length), nil
			}
		}
		return lang.NumberValue(0), nil
	}
	return name, fn
}

func host() (string, lang.Function) {
	name := "host"
	_, Headers := headers()
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		headers, err := Headers(args)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		if headersMap, ok := headers.(lang.MapValue); ok {
			return lang.StringValue(getHeaderValue(headersMap, "Host")), nil
		}
		return lang.StringValue(""), nil
	}
	return name, fn
}

func scheme() (string, lang.Function) {
	name := "scheme"
	_, Headers := headers()
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		contextMap, ok := args[0].(lang.MapValue)
		if !ok {
			return nil, lib.ContextError(name, args[0])
		}

		if scheme, exists := contextMap["scheme"]; exists {
			return scheme, nil
		}

		headers, err := Headers(args)
		if err == nil {
			if headersMap, ok := headers.(lang.MapValue); ok {
				proto := getHeaderValue(headersMap, "X-Forwarded-Proto")
				if proto != "" {
					return lang.StringValue(proto), nil
				}
			}
		}
		return lang.StringValue("https"), nil
	}
	return name, fn
}

func port() (string, lang.Function) {
	name := "port"
	_, Scheme := scheme()
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		contextMap, ok := args[0].(lang.MapValue)
		if !ok {
			return nil, lib.ContextError(name, args[0])
		}
		if port, exists := contextMap["port"]; exists {
			return port, nil
		}
		scheme, err := Scheme(args)
		if err == nil {
			if schemeStr, err := lib.ToString(scheme); err == nil && string(schemeStr) == "https" {
				return lang.NumberValue(443), nil
			}
		}
		return lang.NumberValue(80), nil
	}
	return name, fn
}

func cookies() (string, lang.Function) {
	name := "cookies"
	_, Headers := headers()
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		headers, err := Headers(args)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		if headersMap, ok := headers.(lang.MapValue); ok {
			cookieHeader := getHeaderValue(headersMap, "Cookie")
			if cookieHeader != "" {
				return parseCookies(cookieHeader), nil
			}
		}
		return lang.MapValue{}, nil
	}
	return name, fn
}

func cookie() (string, lang.Function) {
	name := "cookie"
	_, Cookies := cookies()
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 2 {
			return nil, lib.ArgumentError(name, 2)
		}
		cookies, err := Cookies(args[:1])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		cookieName, err := lib.ToString(args[1])
		if err != nil {
			return nil, fmt.Errorf("%s: cookie name %w", name, err)
		}

		if cookiesMap, ok := cookies.(lang.MapValue); ok {
			return cookiesMap[string(cookieName)], nil
		}
		return nil, nil
	}
	return name, fn
}

func referer() (string, lang.Function) {
	name := "referer"
	_, Headers := headers()
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		headers, err := Headers(args)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		if headersMap, ok := headers.(lang.MapValue); ok {
			return lang.StringValue(getHeaderValue(headersMap, "Referer")), nil
		}
		return lang.StringValue(""), nil
	}
	return name, fn
}

func authorization() (string, lang.Function) {
	name := "authorization"
	_, Headers := headers()
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		headers, err := Headers(args)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		if headersMap, ok := headers.(lang.MapValue); ok {
			return lang.StringValue(getHeaderValue(headersMap, "Authorization")), nil
		}
		return lang.StringValue(""), nil
	}
	return name, fn
}

func accept() (string, lang.Function) {
	name := "accept"
	_, Headers := headers()
	fn := func(args []lang.Value) (lang.Value, error) {
		if len(args) != 1 {
			return nil, lib.ArgumentError(name, 1)
		}
		headers, err := Headers(args)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		if headersMap, ok := headers.(lang.MapValue); ok {
			return lang.StringValue(getHeaderValue(headersMap, "Accept")), nil
		}
		return lang.StringValue(""), nil
	}
	return name, fn
}

var HttpFunctions = []func() (string, lang.Function){
	header,
	headers,
	method,
	path,
	query,
	queryParam,
	body,
	status,
	ip,
	userAgent,
	contentType,
	contentLength,
	host,
	scheme,
	port,
	cookies,
	cookie,
	referer,
	authorization,
	accept,
}
