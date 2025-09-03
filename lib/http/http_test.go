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
package http

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/vedadiyan/exql/lang"
)

// Helper function to create a mock HTTP context
func createMockContext() HttpProtocol {
	headers := http.Header{
		"Content-Type":      []string{"application/json; charset=utf-8"},
		"User-Agent":        []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"},
		"Host":              []string{"example.com"},
		"Authorization":     []string{"Bearer token123"},
		"Content-Length":    []string{"1024"},
		"Cookie":            []string{"session=abc123; theme=dark; lang=en"},
		"Referer":           []string{"https://google.com"},
		"Accept":            []string{"text/html,application/xhtml+xml,application/xml;q=0.9"},
		"X-Forwarded-For":   []string{"192.168.1.1, 10.0.0.1"},
		"X-Real-IP":         []string{"203.0.113.1"},
		"X-Forwarded-Proto": []string{"https"},
	}

	url := url.URL{}
	url.Path = "/api/users"
	url.Scheme = "https"
	query := url.Query()
	query.Add("page", "1")
	query.Add("limit", "10")
	query.Add("search", "test query")
	url.RawQuery = query.Encode()

	request := http.Request{}
	request.Method = "POST"
	request.Body = io.NopCloser(bytes.NewBufferString(`{"name": "John", "email": "john@example.com"}`))
	request.Response = &http.Response{}
	request.Response.StatusCode = 200
	request.Header = headers
	request.RemoteAddr = "192.168.1.100"
	request.URL = &url

	return New(&request)
}

// Test Header Functions
func TestHeader(t *testing.T) {
	_, fn := headerFn()
	ctx := createMockContext()

	tests := []struct {
		name     string
		args     []lang.Value
		expected string
		hasError bool
	}{
		{
			name:     "get content-type header",
			args:     []lang.Value{ctx, lang.StringValue("Content-Type")},
			expected: "application/json; charset=utf-8",
			hasError: false,
		},
		{
			name:     "case insensitive header",
			args:     []lang.Value{ctx, lang.StringValue("content-type")},
			expected: "application/json; charset=utf-8",
			hasError: false,
		},
		{
			name:     "non-existent header",
			args:     []lang.Value{ctx, lang.StringValue("X-Custom")},
			expected: "",
			hasError: false,
		},
		{
			name:     "wrong argument count",
			args:     []lang.Value{ctx},
			hasError: true,
		},
		{
			name:     "invalid context type",
			args:     []lang.Value{lang.StringValue("invalid"), lang.StringValue("Content-Type")},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn(tt.args)
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
			if result == nil {
				if tt.expected != "" {
					t.Errorf("Expected %s, got nil", tt.expected)
				}
			} else if string(result.(lang.StringValue)) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(result.(lang.StringValue)))
			}
		})
	}
}

func TestHeaders(t *testing.T) {
	_, fn := headersFn()
	ctx := createMockContext()

	result, err := fn([]lang.Value{ctx})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	headersMap, ok := result.(lang.MapValue)
	if !ok {
		t.Errorf("Expected MapValue, got %T", result)
		return
	}

	expectedHeaders := []string{"Content-Type", "User-Agent", "Host"}
	for _, header := range expectedHeaders {
		if _, exists := headersMap[header]; !exists {
			t.Errorf("Expected header %s not found", header)
		}
	}
}

func TestMethod(t *testing.T) {
	_, fn := methodFn()
	ctx := createMockContext()

	result, err := fn([]lang.Value{ctx})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	expected := "POST"
	if string(result.(lang.StringValue)) != expected {
		t.Errorf("Expected %s, got %s", expected, string(result.(lang.StringValue)))
	}
}

func TestPath(t *testing.T) {
	_, fn := pathFn()
	ctx := createMockContext()

	result, err := fn([]lang.Value{ctx})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	expected := "/api/users"
	if string(result.(lang.StringValue)) != expected {
		t.Errorf("Expected %s, got %s", expected, string(result.(lang.StringValue)))
	}
}

func TestQuery(t *testing.T) {
	_, fn := queryFn()
	ctx := createMockContext()

	result, err := fn([]lang.Value{ctx})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	queryMap, ok := result.(lang.MapValue)
	if !ok {
		t.Errorf("Expected MapValue, got %T", result)
		return
	}

	expectedParams := map[string]lang.ListValue{
		"page":   {"1"},
		"limit":  {"10"},
		"search": {"test query"},
	}

	for key, expectedValue := range expectedParams {
		if actualValue, exists := queryMap[key]; !exists {
			t.Errorf("Expected query param %s not found", key)
		} else if !reflect.DeepEqual(actualValue, expectedValue) {
			t.Errorf("Expected query param %s=%s, got %v", key, expectedValue, actualValue)
		}
	}
}

func TestQueryParam(t *testing.T) {
	_, fn := queryParamFn()
	ctx := createMockContext()

	tests := []struct {
		name      string
		paramName string
		expected  lang.ListValue
		hasError  bool
	}{
		{
			name:      "existing param",
			paramName: "page",
			expected:  lang.ListValue{"1"},
			hasError:  false,
		},
		{
			name:      "non-existent param",
			paramName: "nonexistent",
			expected:  nil,
			hasError:  false,
		},
		{
			name:     "wrong argument count",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var args []lang.Value
			if !tt.hasError {
				args = []lang.Value{ctx, lang.StringValue(tt.paramName)}
			} else {
				args = []lang.Value{ctx}
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
			if result == nil {
				if tt.expected != nil {
					t.Errorf("Expected %s, got nil", tt.expected)
				}
			} else if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected %s, got %v", tt.expected, result)
			}
		})
	}
}

func TestBody(t *testing.T) {
	_, fn := bodyFn()
	ctx := createMockContext()

	result, err := fn([]lang.Value{ctx})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	expected := lang.Value(lang.StringValue(`{"name": "John", "email": "john@example.com"}`))
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %s, got %v", expected, result)
	}
}

func TestStatus(t *testing.T) {
	_, fn := statusFn()
	ctx := createMockContext()

	result, err := fn([]lang.Value{ctx})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	expected := float64(200)
	if float64(result.(lang.NumberValue)) != expected {
		t.Errorf("Expected %f, got %f", expected, float64(result.(lang.NumberValue)))
	}
}

func TestIP(t *testing.T) {
	_, fn := ipFn()

	tests := []struct {
		name     string
		ctx      HttpProtocol
		expected string
	}{
		{
			name:     "direct IP field",
			ctx:      createMockContext(),
			expected: "192.168.1.100",
		},
		{
			name:     "X-Forwarded-For with multiple IPs",
			ctx:      createMockContext(),
			expected: "192.168.1.1",
		},
		{
			name:     "X-Real-IP header",
			ctx:      createMockContext(),
			expected: "203.0.113.1",
		},
		{
			name:     "no IP found",
			ctx:      createMockContext(),
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{tt.ctx})
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if string(result.(lang.StringValue)) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(result.(lang.StringValue)))
			}
		})
	}
}

func TestUserAgent(t *testing.T) {
	_, fn := userAgentFn()
	ctx := createMockContext()

	result, err := fn([]lang.Value{ctx})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	expected := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
	if string(result.(lang.StringValue)) != expected {
		t.Errorf("Expected %s, got %s", expected, string(result.(lang.StringValue)))
	}
}

func TestContentType(t *testing.T) {
	_, fn := contentTypeFn()

	tests := []struct {
		name     string
		ctx      lang.MapValue
		expected string
	}{
		{
			name: "content type with charset",
			ctx: lang.MapValue{
				"headers": lang.MapValue{
					"Content-Type": lang.StringValue("application/json; charset=utf-8"),
				},
			},
			expected: "application/json",
		},
		{
			name: "content type without charset",
			ctx: lang.MapValue{
				"headers": lang.MapValue{
					"Content-Type": lang.StringValue("text/html"),
				},
			},
			expected: "text/html",
		},
		{
			name:     "no content type",
			ctx:      lang.MapValue{"headers": lang.MapValue{}},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{tt.ctx})
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if string(result.(lang.StringValue)) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(result.(lang.StringValue)))
			}
		})
	}
}

func TestContentLength(t *testing.T) {
	_, fn := contentLengthFn()

	tests := []struct {
		name     string
		ctx      lang.MapValue
		expected float64
	}{
		{
			name: "valid content length",
			ctx: lang.MapValue{
				"headers": lang.MapValue{
					"Content-Length": lang.StringValue("1024"),
				},
			},
			expected: 1024,
		},
		{
			name: "invalid content length",
			ctx: lang.MapValue{
				"headers": lang.MapValue{
					"Content-Length": lang.StringValue("invalid"),
				},
			},
			expected: 0,
		},
		{
			name:     "no content length",
			ctx:      lang.MapValue{"headers": lang.MapValue{}},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{tt.ctx})
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

func TestHost(t *testing.T) {
	_, fn := hostFn()
	ctx := createMockContext()

	result, err := fn([]lang.Value{ctx})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	expected := "example.com"
	if string(result.(lang.StringValue)) != expected {
		t.Errorf("Expected %s, got %s", expected, string(result.(lang.StringValue)))
	}
}

func TestScheme(t *testing.T) {
	_, fn := schemeFn()

	tests := []struct {
		name     string
		ctx      lang.MapValue
		expected string
	}{
		{
			name:     "scheme in context",
			ctx:      lang.MapValue{"scheme": lang.StringValue("http")},
			expected: "http",
		},
		{
			name: "X-Forwarded-Proto header",
			ctx: lang.MapValue{
				"headers": lang.MapValue{
					"X-Forwarded-Proto": lang.StringValue("https"),
				},
			},
			expected: "https",
		},
		{
			name:     "default scheme",
			ctx:      lang.MapValue{},
			expected: "https",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{tt.ctx})
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if string(result.(lang.StringValue)) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(result.(lang.StringValue)))
			}
		})
	}
}

func TestPort(t *testing.T) {
	_, fn := portFn()

	tests := []struct {
		name     string
		ctx      lang.MapValue
		expected float64
	}{
		{
			name:     "port in context",
			ctx:      lang.MapValue{"port": lang.NumberValue(8080)},
			expected: 8080,
		},
		{
			name:     "https scheme default port",
			ctx:      lang.MapValue{"scheme": lang.StringValue("https")},
			expected: 443,
		},
		{
			name:     "http scheme default port",
			ctx:      lang.MapValue{"scheme": lang.StringValue("http")},
			expected: 80,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{tt.ctx})
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

func TestCookies(t *testing.T) {
	_, fn := cookiesFn()
	ctx := createMockContext()

	result, err := fn([]lang.Value{ctx})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	cookiesMap, ok := result.(lang.MapValue)
	if !ok {
		t.Errorf("Expected MapValue, got %T", result)
		return
	}

	expectedCookies := map[string]string{
		"session": "abc123",
		"theme":   "dark",
		"lang":    "en",
	}

	for key, expectedValue := range expectedCookies {
		if actualValue, exists := cookiesMap[key]; !exists {
			t.Errorf("Expected cookie %s not found", key)
		} else if string(actualValue.(lang.StringValue)) != expectedValue {
			t.Errorf("Expected cookie %s=%s, got %s", key, expectedValue, string(actualValue.(lang.StringValue)))
		}
	}
}

func TestCookie(t *testing.T) {
	_, fn := cookieFn()
	ctx := createMockContext()

	tests := []struct {
		name       string
		cookieName string
		expected   string
		hasError   bool
	}{
		{
			name:       "existing cookie",
			cookieName: "session",
			expected:   "abc123",
			hasError:   false,
		},
		{
			name:       "non-existent cookie",
			cookieName: "nonexistent",
			expected:   "",
			hasError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{ctx, lang.StringValue(tt.cookieName)})
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
			if result == nil {
				if tt.expected != "" {
					t.Errorf("Expected %s, got nil", tt.expected)
				}
			} else if string(result.(lang.StringValue)) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(result.(lang.StringValue)))
			}
		})
	}
}

func TestReferer(t *testing.T) {
	_, fn := refererFn()
	ctx := createMockContext()

	result, err := fn([]lang.Value{ctx})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	expected := "https://google.com"
	if string(result.(lang.StringValue)) != expected {
		t.Errorf("Expected %s, got %s", expected, string(result.(lang.StringValue)))
	}
}

func TestAuthorization(t *testing.T) {
	_, fn := authorizationFn()
	ctx := createMockContext()

	result, err := fn([]lang.Value{ctx})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	expected := "Bearer token123"
	if string(result.(lang.StringValue)) != expected {
		t.Errorf("Expected %s, got %s", expected, string(result.(lang.StringValue)))
	}
}

func TestAccept(t *testing.T) {
	_, fn := acceptFn()
	ctx := createMockContext()

	result, err := fn([]lang.Value{ctx})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	expected := "text/html,application/xhtml+xml,application/xml;q=0.9"
	if string(result.(lang.StringValue)) != expected {
		t.Errorf("Expected %s, got %s", expected, string(result.(lang.StringValue)))
	}
}

// Test Error Cases
func TestFunctionsWithInvalidContext(t *testing.T) {
	functions := []func() (string, lang.Function){
		methodFn,
		pathFn,
		bodyFn,
		statusFn,
	}

	invalidContext := lang.StringValue("invalid")

	for _, tf := range functions {
		name, fn := tf()
		t.Run(name+"_invalid_context", func(t *testing.T) {
			_, err := fn([]lang.Value{invalidContext})
			if err == nil {
				t.Errorf("Expected error for invalid context in %s", name)
			}
		})
	}
}

func TestFunctionsWithMissingData(t *testing.T) {
	emptyCtx := lang.MapValue{}

	tests := []func() (string, lang.Function){
		methodFn,
		pathFn,
		bodyFn,
		statusFn,
	}

	expected := []any{
		"",
		"",
		"",
		float64(0),
	}

	for it, tt := range tests {
		name, fn := tt()
		t.Run(name+"_missing_data", func(t *testing.T) {
			result, err := fn([]lang.Value{emptyCtx})
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			switch expected := expected[it].(type) {
			case string:
				if string(result.(lang.StringValue)) != expected {
					t.Errorf("Expected %s, got %s", expected, string(result.(lang.StringValue)))
				}
			case float64:
				if float64(result.(lang.NumberValue)) != expected {
					t.Errorf("Expected %f, got %f", expected, float64(result.(lang.NumberValue)))
				}
			}
		})
	}
}

// Test Export Function
func TestExport(t *testing.T) {
	functions := Export()

	expectedFunctions := []string{
		"header", "headers", "method", "path", "query", "queryParam",
		"body", "status", "ip", "userAgent", "contentType", "contentLength",
		"host", "scheme", "port", "cookies", "cookie", "referer",
		"authorization", "accept",
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

// Benchmark tests
func BenchmarkHeader(b *testing.B) {
	_, fn := headerFn()
	ctx := createMockContext()
	args := []lang.Value{ctx, lang.StringValue("Content-Type")}

	for i := 0; i < b.N; i++ {
		fn(args)
	}
}

func BenchmarkIP(b *testing.B) {
	_, fn := ipFn()
	ctx := createMockContext()
	args := []lang.Value{ctx}

	for i := 0; i < b.N; i++ {
		fn(args)
	}
}

func BenchmarkCookies(b *testing.B) {
	_, fn := cookiesFn()
	ctx := createMockContext()
	args := []lang.Value{ctx}

	for i := 0; i < b.N; i++ {
		fn(args)
	}
}
