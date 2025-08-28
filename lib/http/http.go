package http

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/vedadiyan/exql/lang"
	"github.com/vedadiyan/exql/lib"
)

// HTTP/Request Inspection Functions for Sidecar Proxy
// These functions expect a context (usually the first argument) containing request/response data

func httpHeader(args []lang.Value) (lang.Value, error) {
	if len(args) != 2 {
		return nil, errors.New("header: expected 2 arguments (context, header_name)")
	}
	// Extract header from context
	contextMap, ok := args[0].(lang.MapValue)
	if !ok {
		return nil, fmt.Errorf("header: first argument must be map, got %T", args[0])
	}
	headerName, err := lib.ToString(args[1])
	if err != nil {
		return nil, fmt.Errorf("header: header name %w", err)
	}

	if headers, exists := contextMap["headers"]; exists {
		if headersMap, ok := headers.(lang.MapValue); ok {
			// Case-insensitive header lookup
			for key, value := range headersMap {
				if strings.EqualFold(key, string(headerName)) {
					return value, nil
				}
			}
		}
	}
	return nil, nil
}

func httpHeaders(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("headers: expected 1 argument (context)")
	}
	contextMap, ok := args[0].(lang.MapValue)
	if !ok {
		return nil, fmt.Errorf("headers: argument must be map, got %T", args[0])
	}
	if headers, exists := contextMap["headers"]; exists {
		if headersMap, ok := headers.(lang.MapValue); ok {
			return headersMap, nil
		}
	}
	return lang.MapValue{}, nil
}

func httpMethod(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("method: expected 1 argument (context)")
	}
	contextMap, ok := args[0].(lang.MapValue)
	if !ok {
		return nil, fmt.Errorf("method: argument must be map, got %T", args[0])
	}
	if method, exists := contextMap["method"]; exists {
		return method, nil
	}
	return lang.StringValue(""), nil
}

func httpPath(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("path: expected 1 argument (context)")
	}
	contextMap, ok := args[0].(lang.MapValue)
	if !ok {
		return nil, fmt.Errorf("path: argument must be map, got %T", args[0])
	}
	if path, exists := contextMap["path"]; exists {
		return path, nil
	}
	return lang.StringValue(""), nil
}

func httpQuery(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("query: expected 1 argument (context)")
	}
	contextMap, ok := args[0].(lang.MapValue)
	if !ok {
		return nil, fmt.Errorf("query: argument must be map, got %T", args[0])
	}
	if query, exists := contextMap["query"]; exists {
		if queryMap, ok := query.(lang.MapValue); ok {
			return queryMap, nil
		}
	}
	return lang.MapValue{}, nil
}

func httpQueryParam(args []lang.Value) (lang.Value, error) {
	if len(args) != 2 {
		return nil, errors.New("query_param: expected 2 arguments (context, param_name)")
	}
	queryParams, err := httpQuery(args[:1])
	if err != nil {
		return nil, fmt.Errorf("query_param: %w", err)
	}
	paramName, err := lib.ToString(args[1])
	if err != nil {
		return nil, fmt.Errorf("query_param: param name %w", err)
	}

	if queryMap, ok := queryParams.(lang.MapValue); ok {
		return queryMap[string(paramName)], nil
	}
	return nil, nil
}

func httpBody(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("body: expected 1 argument (context)")
	}
	contextMap, ok := args[0].(lang.MapValue)
	if !ok {
		return nil, fmt.Errorf("body: argument must be map, got %T", args[0])
	}
	if body, exists := contextMap["body"]; exists {
		return body, nil
	}
	return lang.StringValue(""), nil
}

func httpStatus(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("status: expected 1 argument (context)")
	}
	contextMap, ok := args[0].(lang.MapValue)
	if !ok {
		return nil, fmt.Errorf("status: argument must be map, got %T", args[0])
	}
	if status, exists := contextMap["status"]; exists {
		return status, nil
	}
	return lang.NumberValue(0), nil
}

func httpIP(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("ip: expected 1 argument (context)")
	}
	contextMap, ok := args[0].(lang.MapValue)
	if !ok {
		return nil, fmt.Errorf("ip: argument must be map, got %T", args[0])
	}

	// Try multiple IP fields in order of preference
	ipFields := []string{"remote_ip", "client_ip", "x_forwarded_for", "x_real_ip", "ip"}

	for _, field := range ipFields {
		if ip, exists := contextMap[field]; exists {
			ipStr, err := lib.ToString(ip)
			if err == nil && string(ipStr) != "" {
				// If it's X-Forwarded-For, take the first IP
				if strings.Contains(string(ipStr), ",") {
					return lang.StringValue(strings.TrimSpace(strings.Split(string(ipStr), ",")[0])), nil
				}
				return lang.StringValue(string(ipStr)), nil
			}
		}
	}

	// Fallback to headers
	headers, err := httpHeaders(args)
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

func httpUserAgent(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("user_agent: expected 1 argument (context)")
	}
	headers, err := httpHeaders(args)
	if err != nil {
		return nil, fmt.Errorf("user_agent: %w", err)
	}
	if headersMap, ok := headers.(lang.MapValue); ok {
		return lang.StringValue(getHeaderValue(headersMap, "User-Agent")), nil
	}
	return lang.StringValue(""), nil
}

func httpContentType(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("content_type: expected 1 argument (context)")
	}
	headers, err := httpHeaders(args)
	if err != nil {
		return nil, fmt.Errorf("content_type: %w", err)
	}
	if headersMap, ok := headers.(lang.MapValue); ok {
		ct := getHeaderValue(headersMap, "Content-Type")
		// Extract just the media type, not the parameters
		if idx := strings.Index(ct, ";"); idx != -1 {
			ct = strings.TrimSpace(ct[:idx])
		}
		return lang.StringValue(ct), nil
	}
	return lang.StringValue(""), nil
}

func httpContentLength(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("content_length: expected 1 argument (context)")
	}
	headers, err := httpHeaders(args)
	if err != nil {
		return nil, fmt.Errorf("content_length: %w", err)
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

func httpHost(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("host: expected 1 argument (context)")
	}
	headers, err := httpHeaders(args)
	if err != nil {
		return nil, fmt.Errorf("host: %w", err)
	}
	if headersMap, ok := headers.(lang.MapValue); ok {
		return lang.StringValue(getHeaderValue(headersMap, "Host")), nil
	}
	return lang.StringValue(""), nil
}

func httpScheme(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("scheme: expected 1 argument (context)")
	}
	contextMap, ok := args[0].(lang.MapValue)
	if !ok {
		return nil, fmt.Errorf("scheme: argument must be map, got %T", args[0])
	}

	if scheme, exists := contextMap["scheme"]; exists {
		return scheme, nil
	}

	// Check X-Forwarded-Proto header
	headers, err := httpHeaders(args)
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

func httpPort(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("port: expected 1 argument (context)")
	}
	contextMap, ok := args[0].(lang.MapValue)
	if !ok {
		return nil, fmt.Errorf("port: argument must be map, got %T", args[0])
	}
	if port, exists := contextMap["port"]; exists {
		return port, nil
	}

	// Infer from scheme if not provided
	scheme, err := httpScheme(args)
	if err == nil {
		if schemeStr, err := lib.ToString(scheme); err == nil && string(schemeStr) == "https" {
			return lang.NumberValue(443), nil
		}
	}
	return lang.NumberValue(80), nil
}

func httpCookies(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("cookies: expected 1 argument (context)")
	}
	headers, err := httpHeaders(args)
	if err != nil {
		return nil, fmt.Errorf("cookies: %w", err)
	}
	if headersMap, ok := headers.(lang.MapValue); ok {
		cookieHeader := getHeaderValue(headersMap, "Cookie")
		if cookieHeader != "" {
			return parseCookies(cookieHeader), nil
		}
	}
	return lang.MapValue{}, nil
}

func httpCookie(args []lang.Value) (lang.Value, error) {
	if len(args) != 2 {
		return nil, errors.New("cookie: expected 2 arguments (context, cookie_name)")
	}
	cookies, err := httpCookies(args[:1])
	if err != nil {
		return nil, fmt.Errorf("cookie: %w", err)
	}
	cookieName, err := lib.ToString(args[1])
	if err != nil {
		return nil, fmt.Errorf("cookie: cookie name %w", err)
	}

	if cookiesMap, ok := cookies.(lang.MapValue); ok {
		return cookiesMap[string(cookieName)], nil
	}
	return nil, nil
}

func httpReferer(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("referer: expected 1 argument (context)")
	}
	headers, err := httpHeaders(args)
	if err != nil {
		return nil, fmt.Errorf("referer: %w", err)
	}
	if headersMap, ok := headers.(lang.MapValue); ok {
		return lang.StringValue(getHeaderValue(headersMap, "Referer")), nil
	}
	return lang.StringValue(""), nil
}

func httpAuthorization(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("authorization: expected 1 argument (context)")
	}
	headers, err := httpHeaders(args)
	if err != nil {
		return nil, fmt.Errorf("authorization: %w", err)
	}
	if headersMap, ok := headers.(lang.MapValue); ok {
		return lang.StringValue(getHeaderValue(headersMap, "Authorization")), nil
	}
	return lang.StringValue(""), nil
}

func httpAccept(args []lang.Value) (lang.Value, error) {
	if len(args) != 1 {
		return nil, errors.New("accept: expected 1 argument (context)")
	}
	headers, err := httpHeaders(args)
	if err != nil {
		return nil, fmt.Errorf("accept: %w", err)
	}
	if headersMap, ok := headers.(lang.MapValue); ok {
		return lang.StringValue(getHeaderValue(headersMap, "Accept")), nil
	}
	return lang.StringValue(""), nil
}

// Helper function for case-insensitive header lookup
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

// Parse cookie header into a map
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

// Functions that would be in the BuiltinFunctions map:
var HttpFunctions = map[string]func([]lang.Value) (lang.Value, error){
	"header":         httpHeader,
	"headers":        httpHeaders,
	"method":         httpMethod,
	"path":           httpPath,
	"query":          httpQuery,
	"query_param":    httpQueryParam,
	"body":           httpBody,
	"status":         httpStatus,
	"ip":             httpIP,
	"user_agent":     httpUserAgent,
	"content_type":   httpContentType,
	"content_length": httpContentLength,
	"host":           httpHost,
	"scheme":         httpScheme,
	"port":           httpPort,
	"cookies":        httpCookies,
	"cookie":         httpCookie,
	"referer":        httpReferer,
	"authorization":  httpAuthorization,
	"accept":         httpAccept,
}
