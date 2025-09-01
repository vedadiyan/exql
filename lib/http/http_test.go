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
	"testing"

	"github.com/vedadiyan/exql/lang"
)

// Helper function to create a mock HTTP context
func createMockContext() lang.MapValue {
	headers := lang.MapValue{
		"Content-Type":      lang.StringValue("application/json; charset=utf-8"),
		"User-Agent":        lang.StringValue("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"),
		"Host":              lang.StringValue("example.com"),
		"Authorization":     lang.StringValue("Bearer token123"),
		"Content-Length":    lang.StringValue("1024"),
		"Cookie":            lang.StringValue("session=abc123; theme=dark; lang=en"),
		"Referer":           lang.StringValue("https://google.com"),
		"Accept":            lang.StringValue("text/html,application/xhtml+xml,application/xml;q=0.9"),
		"X-Forwarded-For":   lang.StringValue("192.168.1.1, 10.0.0.1"),
		"X-Real-IP":         lang.StringValue("203.0.113.1"),
		"X-Forwarded-Proto": lang.StringValue("https"),
	}

	query := lang.MapValue{
		"page":   lang.StringValue("1"),
		"limit":  lang.StringValue("10"),
		"search": lang.StringValue("test query"),
	}

	return lang.MapValue{
		"method":  lang.StringValue("POST"),
		"path":    lang.StringValue("/api/users"),
		"body":    lang.StringValue(`{"name": "John", "email": "john@example.com"}`),
		"status":  lang.NumberValue(200),
		"scheme":  lang.StringValue("https"),
		"port":    lang.NumberValue(443),
		"headers": headers,
		"query":   query,
		"ip":      lang.StringValue("192.168.1.100"),
	}
}

// Test Helper Functions
func TestGetHeaderValue(t *testing.T) {
	headers := lang.MapValue{
		"Content-Type": lang.StringValue("application/json"),
		"user-agent":   lang.StringValue("test-agent"),
		"ACCEPT":       lang.StringValue("text/html"),
	}

	tests := []struct {
		name       string
		headerName string
		expected   string
	}{
		{"exact match", "Content-Type", "application/json"},
		{"case insensitive lowercase", "content-type", "application/json"},
		{"case insensitive uppercase", "USER-AGENT", "test-agent"},
		{"mixed case", "Accept", "text/html"},
		{"not found", "X-Custom-Header", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getHeaderValue(headers, tt.headerName)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestParseCookies(t *testing.T) {
	tests := []struct {
		name         string
		cookieHeader string
		expected     map[string]string
	}{
		{
			name:         "single cookie",
			cookieHeader: "session=abc123",
			expected:     map[string]string{"session": "abc123"},
		},
		{
			name:         "multiple cookies",
			cookieHeader: "session=abc123; theme=dark; lang=en",
			expected:     map[string]string{"session": "abc123", "theme": "dark", "lang": "en"},
		},
		{
			name:         "empty cookie header",
			cookieHeader: "",
			expected:     map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseCookies(tt.cookieHeader)
			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d cookies, got %d", len(tt.expected), len(result))
			}
			for key, expectedValue := range tt.expected {
				if actualValue, exists := result[key]; !exists {
					t.Errorf("Expected cookie %s not found", key)
				} else if string(actualValue.(lang.StringValue)) != expectedValue {
					t.Errorf("Expected cookie %s=%s, got %s", key, expectedValue, string(actualValue.(lang.StringValue)))
				}
			}
		})
	}
}

// Test Header Functions
func TestHeader(t *testing.T) {
	_, fn := header()
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
	_, fn := headers()
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
	_, fn := method()
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
	_, fn := path()
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
	_, fn := query()
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

	expectedParams := map[string]string{
		"page":   "1",
		"limit":  "10",
		"search": "test query",
	}

	for key, expectedValue := range expectedParams {
		if actualValue, exists := queryMap[key]; !exists {
			t.Errorf("Expected query param %s not found", key)
		} else if string(actualValue.(lang.StringValue)) != expectedValue {
			t.Errorf("Expected query param %s=%s, got %s", key, expectedValue, string(actualValue.(lang.StringValue)))
		}
	}
}

func TestQueryParam(t *testing.T) {
	_, fn := queryParam()
	ctx := createMockContext()

	tests := []struct {
		name      string
		paramName string
		expected  string
		hasError  bool
	}{
		{
			name:      "existing param",
			paramName: "page",
			expected:  "1",
			hasError:  false,
		},
		{
			name:      "non-existent param",
			paramName: "nonexistent",
			expected:  "",
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
				if tt.expected != "" {
					t.Errorf("Expected %s, got nil", tt.expected)
				}
			} else if string(result.(lang.StringValue)) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(result.(lang.StringValue)))
			}
		})
	}
}

func TestBody(t *testing.T) {
	_, fn := body()
	ctx := createMockContext()

	result, err := fn([]lang.Value{ctx})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	expected := `{"name": "John", "email": "john@example.com"}`
	if string(result.(lang.StringValue)) != expected {
		t.Errorf("Expected %s, got %s", expected, string(result.(lang.StringValue)))
	}
}

func TestStatus(t *testing.T) {
	_, fn := status()
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
	_, fn := ip()

	tests := []struct {
		name     string
		ctx      lang.MapValue
		expected string
	}{
		{
			name:     "direct IP field",
			ctx:      lang.MapValue{"ip": lang.StringValue("192.168.1.100")},
			expected: "192.168.1.100",
		},
		{
			name:     "X-Forwarded-For with multiple IPs",
			ctx:      lang.MapValue{"headers": lang.MapValue{"X-Forwarded-For": lang.StringValue("192.168.1.1, 10.0.0.1")}},
			expected: "192.168.1.1",
		},
		{
			name:     "X-Real-IP header",
			ctx:      lang.MapValue{"headers": lang.MapValue{"X-Real-IP": lang.StringValue("203.0.113.1")}},
			expected: "203.0.113.1",
		},
		{
			name:     "no IP found",
			ctx:      lang.MapValue{},
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
	_, fn := userAgent()
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
	_, fn := contentType()

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
	_, fn := contentLength()

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
	_, fn := host()
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
	_, fn := scheme()

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
	_, fn := port()

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
	_, fn := cookies()
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
	_, fn := cookie()
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
	_, fn := referer()
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
	_, fn := authorization()
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
	_, fn := accept()
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
		method,
		path,
		body,
		status,
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
		method,
		path,
		body,
		status,
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
	_, fn := header()
	ctx := createMockContext()
	args := []lang.Value{ctx, lang.StringValue("Content-Type")}

	for i := 0; i < b.N; i++ {
		fn(args)
	}
}

func BenchmarkIP(b *testing.B) {
	_, fn := ip()
	ctx := createMockContext()
	args := []lang.Value{ctx}

	for i := 0; i < b.N; i++ {
		fn(args)
	}
}

func BenchmarkCookies(b *testing.B) {
	_, fn := cookies()
	ctx := createMockContext()
	args := []lang.Value{ctx}

	for i := 0; i < b.N; i++ {
		fn(args)
	}
}
