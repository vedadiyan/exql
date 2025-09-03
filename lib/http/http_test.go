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
func mockRequest() HttpProtocol {
	headers := http.Header{
		"Content-Type":      []string{"application/json; charset=utf-8"},
		"User-Agent":        []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"},
		"Host":              []string{"example.com"},
		"Authorization":     []string{"Bearer token123"},
		"Content-Length":    []string{"1024"},
		"Cookie":            []string{"session=abc123; path=/", "theme=dark; path=/", "lang=en; path=/"},
		"Referer":           []string{"https://google.com"},
		"Accept":            []string{"text/html,application/xhtml+xml,application/xml;q=0.9"},
		"X-Forwarded-For":   []string{"192.168.1.1, 10.0.0.1"},
		"X-Real-IP":         []string{"203.0.113.1"},
		"X-Forwarded-Proto": []string{"https"},
	}

	trailers := http.Header{
		"Test-1": []string{"trailer-value"},
		"Test-2": []string{"custom-trailer"},
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
	request.ContentLength = 1024
	request.Host = "example.com"
	request.Proto = "HTTP/1.1"
	request.ProtoMajor = 1
	request.ProtoMinor = 1
	request.TransferEncoding = []string{"chunked", "gzip"}
	request.Body = io.NopCloser(bytes.NewBufferString(`{"name": "John", "email": "john@example.com"}`))
	request.Response = &http.Response{
		StatusCode: 200,
		Trailer:    trailers,
	}
	request.Header = headers
	request.Trailer = trailers
	request.RemoteAddr = "192.168.1.100"
	request.URL = &url
	request.TransferEncoding = []string{"chunked", "gzip"}

	return New(&request)
}

func mockResponse() HttpProtocol {
	headers := http.Header{
		"Content-Type":      []string{"application/json; charset=utf-8"},
		"User-Agent":        []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"},
		"Host":              []string{"example.com"},
		"Authorization":     []string{"Bearer token123"},
		"Content-Length":    []string{"1024"},
		"Set-Cookie":        []string{"session=abc123; path=/", "theme=dark; path=/", "lang=en; path=/"},
		"Referer":           []string{"https://google.com"},
		"Accept":            []string{"text/html,application/xhtml+xml,application/xml;q=0.9"},
		"X-Forwarded-For":   []string{"192.168.1.1, 10.0.0.1"},
		"X-Real-IP":         []string{"203.0.113.1"},
		"X-Forwarded-Proto": []string{"https"},
		"X-Trailer-One":     []string{"trailer1"},
		"X-Trailer-Two":     []string{"trailer2", "trailer3"},
	}

	url := url.URL{}
	url.Path = "/api/users"
	url.Scheme = "https"
	url.Host = "example.com"

	query := url.Query()
	query.Add("page", "1")
	query.Add("limit", "10")
	query.Add("search", "test query")
	url.RawQuery = query.Encode()

	request := http.Response{}
	request.ContentLength = 1024
	request.Body = io.NopCloser(bytes.NewBufferString(`{"name": "John", "email": "john@example.com"}`))
	request.StatusCode = 200
	request.Header = headers
	request.Request = &http.Request{URL: &url}

	return New(&request)
}

// Helper function to create a mock HTTP context with specific headers
func createMockContextWithHeaders(headers http.Header) HttpProtocol {
	url := url.URL{}
	url.Path = "/api/users"
	url.Scheme = "https"

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

// Helper function to create a mock HTTP context with specific scheme
func createMockContextWithScheme(scheme string) HttpProtocol {
	headers := http.Header{}
	if scheme != "" {
		headers.Set("X-Forwarded-Proto", scheme)
	}

	url := url.URL{}
	url.Path = "/api/users"
	if scheme != "" {
		url.Scheme = scheme
	}

	request := http.Request{}
	request.Method = "POST"
	request.Header = headers
	request.URL = &url

	return New(&request)
}

// Helper function to create a mock HTTP context with specific port
func createMockContextWithPort(port string, scheme string) HttpProtocol {
	headers := http.Header{}
	if port != "" {
		headers.Set("Host", "example.com:"+port)
	} else {
		headers.Set("Host", "example.com")
	}

	url := url.URL{}
	url.Path = "/api/users"
	url.Scheme = scheme
	if port != "" {
		url.Host = "example.com:" + port
	} else {
		url.Host = "example.com"
	}

	request := http.Request{}
	request.Method = "POST"
	request.Header = headers
	request.URL = &url

	return New(&request)
}

// Test Header Functions
func TestHeader(t *testing.T) {
	_, fn := headerFn()
	ctx := mockRequest()

	tests := []struct {
		name     string
		args     []lang.Value
		expected lang.Value
		hasError bool
	}{
		{
			name:     "get content-type header",
			args:     []lang.Value{ctx, lang.StringValue("Content-Type")},
			expected: lang.StringValue("application/json; charset=utf-8"),
			hasError: false,
		},
		{
			name:     "case insensitive header",
			args:     []lang.Value{ctx, lang.StringValue("content-type")},
			expected: lang.StringValue("application/json; charset=utf-8"),
			hasError: false,
		},
		{
			name:     "non-existent header",
			args:     []lang.Value{ctx, lang.StringValue("X-Custom")},
			expected: lang.StringValue(""),
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
				if tt.expected != nil {
					t.Errorf("Expected %v, got nil", tt.expected)
				}
			} else if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestHeaders(t *testing.T) {
	_, fn := headersFn()
	ctx := mockRequest()

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

func TestTrailer(t *testing.T) {
	_, fn := trailerFn()
	ctx := mockRequest()

	tests := []struct {
		name        string
		trailerName string
		expected    lang.Value
		hasError    bool
	}{
		{
			name:        "existing trailer",
			trailerName: "Test-1",
			expected:    lang.StringValue("trailer-value"),
			hasError:    false,
		},
		{
			name:        "non-existent trailer",
			trailerName: "X-Missing",
			expected:    lang.StringValue(""),
			hasError:    false,
		},
		{
			name:     "wrong argument count",
			hasError: true,
		},
		{
			name:        "invalid context type",
			trailerName: "X-Test",
			hasError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var args []lang.Value
			if tt.name == "invalid context type" {
				args = []lang.Value{lang.StringValue("invalid"), lang.StringValue(tt.trailerName)}
			} else if !tt.hasError {
				args = []lang.Value{ctx, lang.StringValue(tt.trailerName)}
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
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestTrailers(t *testing.T) {
	_, fn := trailersFn()
	ctx := mockRequest()

	result, err := fn([]lang.Value{ctx})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	trailersMap, ok := result.(lang.MapValue)
	if !ok {
		t.Errorf("Expected MapValue, got %T", result)
		return
	}

	expectedTrailers := []string{"Test-1", "Test-2"}
	for _, trailer := range expectedTrailers {
		if _, exists := trailersMap[trailer]; !exists {
			t.Errorf("Expected trailer %s not found", trailer)
		}
	}
}

func TestRouteValues(t *testing.T) {
	_, fn := routeValuesFn()
	ctx := mockRequest()

	result, err := fn([]lang.Value{ctx})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	routeMap, ok := result.(lang.MapValue)
	if !ok {
		t.Errorf("Expected MapValue, got %T", result)
		return
	}

	if routeMap == nil {
		t.Errorf("Expected non-nil MapValue")
	}
}

func TestPattern(t *testing.T) {
	_, fn := patternFn()
	ctx := mockRequest()

	result, err := fn([]lang.Value{ctx})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	expected := lang.StringValue("")
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestProto(t *testing.T) {
	_, fn := protoFn()
	ctx := mockRequest()

	result, err := fn([]lang.Value{ctx})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	expected := lang.StringValue("HTTP/1.1")
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestProtoMajor(t *testing.T) {
	_, fn := protoMajorFn()
	ctx := mockRequest()

	result, err := fn([]lang.Value{ctx})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	expected := lang.NumberValue(1)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestProtoMinor(t *testing.T) {
	_, fn := protoMinorFn()
	ctx := mockRequest()

	result, err := fn([]lang.Value{ctx})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	expected := lang.NumberValue(1)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestTransferEncoding(t *testing.T) {
	_, fn := transferEncodingFn()
	ctx := mockRequest()

	result, err := fn([]lang.Value{ctx})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	expected := lang.ListValue{
		lang.StringValue("chunked"),
		lang.StringValue("gzip"),
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestUrl(t *testing.T) {
	_, fn := urlFn()
	ctx := mockResponse()

	result, err := fn([]lang.Value{ctx})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	urlMap, ok := result.(lang.MapValue)
	if !ok {
		t.Errorf("Expected MapValue, got %T", result)
		return
	}

	expectedFields := []string{"scheme", "host", "path", "rawQuery", "fragment", "user"}
	for _, field := range expectedFields {
		if _, exists := urlMap[field]; !exists {
			t.Errorf("Expected URL field %s not found", field)
		}
	}

	// Check user field structure
	userValue, exists := urlMap["user"]
	if !exists {
		t.Errorf("Expected 'user' field in URL")
		return
	}

	userMap, ok := userValue.(lang.MapValue)
	if !ok {
		t.Errorf("Expected MapValue for user field, got %T", userValue)
		return
	}

	if _, exists := userMap["username"]; !exists {
		t.Errorf("Expected 'username' in user field")
	}
	if _, exists := userMap["password"]; !exists {
		t.Errorf("Expected 'password' in user field")
	}
}

func TestMethod(t *testing.T) {
	_, fn := methodFn()
	ctx := mockRequest()

	result, err := fn([]lang.Value{ctx})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	expected := lang.StringValue("POST")
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestPath(t *testing.T) {
	_, fn := pathFn()
	ctx := mockRequest()

	result, err := fn([]lang.Value{ctx})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	expected := lang.StringValue("/api/users")
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestQuery(t *testing.T) {
	_, fn := queryFn()
	ctx := mockRequest()

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
			t.Errorf("Expected query param %s=%v, got %v", key, expectedValue, actualValue)
		}
	}
}

func TestQueryParam(t *testing.T) {
	_, fn := queryParamFn()
	ctx := mockRequest()

	tests := []struct {
		name      string
		paramName string
		expected  lang.Value
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
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestBody(t *testing.T) {
	_, fn := bodyFn()
	ctx := mockRequest()

	result, err := fn([]lang.Value{ctx})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	expected := lang.StringValue(`{"name": "John", "email": "john@example.com"}`)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestStatus(t *testing.T) {
	_, fn := statusFn()
	ctx := mockRequest()

	result, err := fn([]lang.Value{ctx})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	expected := lang.NumberValue(200)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestIP(t *testing.T) {
	_, fn := ipFn()

	tests := []struct {
		name     string
		ctx      HttpProtocol
		expected lang.Value
	}{
		{
			name:     "X-Forwarded-For with multiple IPs",
			ctx:      mockRequest(),
			expected: lang.StringValue("192.168.1.1"), // First IP from X-Forwarded-For should take precedence
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{tt.ctx})
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestUserAgent(t *testing.T) {
	_, fn := userAgentFn()
	ctx := mockRequest()

	result, err := fn([]lang.Value{ctx})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	expected := lang.StringValue("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestContentType(t *testing.T) {
	_, fn := contentTypeFn()

	tests := []struct {
		name     string
		ctx      HttpProtocol
		expected lang.Value
	}{
		{
			name:     "content type with charset",
			ctx:      mockRequest(),
			expected: lang.StringValue("application/json; charset=utf-8"),
		},
		{
			name: "content type without charset",
			ctx: createMockContextWithHeaders(http.Header{
				"Content-Type": []string{"text/html"},
			}),
			expected: lang.StringValue("text/html"),
		},
		{
			name:     "no content type",
			ctx:      createMockContextWithHeaders(http.Header{}),
			expected: lang.StringValue(""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{tt.ctx})
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestContentLength(t *testing.T) {
	_, fn := contentLengthFn()

	tests := []struct {
		name     string
		ctx      HttpProtocol
		expected lang.Value
	}{
		{
			name:     "valid content length",
			ctx:      mockRequest(),
			expected: lang.NumberValue(1024),
		},
		{
			name: "invalid content length",
			ctx: createMockContextWithHeaders(http.Header{
				"Content-Length": []string{"invalid"},
			}),
			expected: lang.NumberValue(0),
		},
		{
			name:     "no content length",
			ctx:      createMockContextWithHeaders(http.Header{}),
			expected: lang.NumberValue(0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{tt.ctx})
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestHost(t *testing.T) {
	_, fn := hostFn()
	ctx := mockRequest()

	result, err := fn([]lang.Value{ctx})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	expected := lang.StringValue("example.com")
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestScheme(t *testing.T) {
	_, fn := schemeFn()

	tests := []struct {
		name     string
		ctx      HttpProtocol
		expected lang.Value
	}{
		{
			name:     "https scheme from URL",
			ctx:      mockRequest(),
			expected: lang.StringValue("https"),
		},
		{
			name:     "X-Forwarded-Proto header",
			ctx:      createMockContextWithScheme("http"),
			expected: lang.StringValue("http"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{tt.ctx})
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestPort(t *testing.T) {
	_, fn := portFn()

	tests := []struct {
		name     string
		ctx      HttpProtocol
		expected lang.Value
	}{
		{
			name:     "explicit port 8080",
			ctx:      createMockContextWithPort("8080", "http"),
			expected: lang.StringValue("8080"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{tt.ctx})
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestCookies(t *testing.T) {
	_, fn := cookiesFn()
	ctx := mockResponse()

	result, err := fn([]lang.Value{ctx})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	// The result should be a ListValue containing MapValues with cookie properties
	resultList, ok := result.(lang.ListValue)
	if !ok {
		t.Errorf("Expected ListValue, got %T", result)
		return
	}

	// Check that we have the expected number of cookies
	expectedCookieCount := 3 // session, theme, lang
	if len(resultList) != expectedCookieCount {
		t.Errorf("Expected %d cookies, got %d", expectedCookieCount, len(resultList))
		return
	}

	// Check that each element is a MapValue and contains expected cookie names
	expectedCookieNames := map[string]bool{
		"session": false,
		"theme":   false,
		"lang":    false,
	}

	for _, cookieValue := range resultList {
		cookieMap, ok := cookieValue.(lang.MapValue)
		if !ok {
			t.Errorf("Expected MapValue for cookie, got %T", cookieValue)
			continue
		}

		nameValue, exists := cookieMap["name"]
		if !exists {
			t.Errorf("Cookie missing 'name' field")
			continue
		}

		cookieName := string(nameValue.(lang.StringValue))
		if _, expected := expectedCookieNames[cookieName]; expected {
			expectedCookieNames[cookieName] = true

			// Check that required fields exist
			if _, exists := cookieMap["value"]; !exists {
				t.Errorf("Cookie %s missing 'value' field", cookieName)
			}
			if _, exists := cookieMap["domain"]; !exists {
				t.Errorf("Cookie %s missing 'domain' field", cookieName)
			}
		}
	}

	// Verify all expected cookies were found
	for cookieName, found := range expectedCookieNames {
		if !found {
			t.Errorf("Expected cookie %s not found", cookieName)
		}
	}
}

func TestCookie(t *testing.T) {
	_, fn := cookieFn()
	ctx := mockResponse()

	tests := []struct {
		name       string
		cookieName string
		expected   lang.Value
		hasError   bool
	}{
		{
			name:       "existing cookie session",
			cookieName: "session",
			expected: lang.MapValue{
				"domain":      lang.StringValue(""),
				"expires":     lang.StringValue("0001-01-01 00:00:00 +0000 UTC"),
				"httpOnly":    lang.BoolValue(false),
				"maxAge":      lang.NumberValue(0),
				"name":        lang.StringValue("session"),
				"partitioned": lang.BoolValue(false),
				"path":        lang.StringValue("/"),
				"quoted":      lang.BoolValue(false),
				"raw":         lang.StringValue("session=abc123; path=/"),
				"rawExpires":  lang.StringValue(""),
				"sameSite":    lang.NumberValue(0),
				"secure":      lang.BoolValue(false),
				"unparsed":    lang.ListValue{},
				"value":       lang.StringValue("abc123"),
			},
			hasError: false,
		},
		{
			name:       "existing cookie theme",
			cookieName: "theme",
			expected: lang.MapValue{
				"domain":      lang.StringValue(""),
				"expires":     lang.StringValue("0001-01-01 00:00:00 +0000 UTC"),
				"httpOnly":    lang.BoolValue(false),
				"maxAge":      lang.NumberValue(0),
				"name":        lang.StringValue("theme"),
				"partitioned": lang.BoolValue(false),
				"path":        lang.StringValue("/"),
				"quoted":      lang.BoolValue(false),
				"raw":         lang.StringValue("theme=dark; path=/"),
				"rawExpires":  lang.StringValue(""),
				"sameSite":    lang.NumberValue(0),
				"secure":      lang.BoolValue(false),
				"unparsed":    lang.ListValue{},
				"value":       lang.StringValue("dark"),
			},
			hasError: false,
		},
		{
			name:       "non-existent cookie",
			cookieName: "nonexistent",
			expected:   nil,
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
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestReferer(t *testing.T) {
	_, fn := refererFn()
	ctx := mockRequest()

	result, err := fn([]lang.Value{ctx})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	expected := lang.StringValue("https://google.com")
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestAuthorization(t *testing.T) {
	_, fn := authorizationFn()
	ctx := mockRequest()

	result, err := fn([]lang.Value{ctx})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	expected := lang.StringValue("Bearer token123")
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestAccept(t *testing.T) {
	_, fn := acceptFn()
	ctx := mockRequest()

	result, err := fn([]lang.Value{ctx})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	expected := lang.StringValue("text/html,application/xhtml+xml,application/xml;q=0.9")
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

// Test Error Cases
func TestFunctionsWithInvalidContext(t *testing.T) {
	functions := []func() (string, lang.Function){
		methodFn,
		pathFn,
		bodyFn,
		statusFn,
		trailerFn,
		patternFn,
		protoFn,
		protoMajorFn,
		protoMinorFn,
		transferEncodingFn,
	}

	invalidContext := lang.StringValue("invalid")

	for _, tf := range functions {
		name, fn := tf()
		t.Run(name+"_invalid_context", func(t *testing.T) {
			var args []lang.Value
			if name == "trailer" {
				args = []lang.Value{invalidContext, lang.StringValue("X-Test")}
			} else {
				args = []lang.Value{invalidContext}
			}
			_, err := fn(args)
			if err == nil {
				t.Errorf("Expected error for invalid context in %s", name)
			}
		})
	}
}

func TestFunctionsWithMissingData(t *testing.T) {
	// Create empty HTTP contexts with minimal request data
	emptyUrl := &url.URL{}
	emptyRequest := &http.Request{
		Method: "",
		URL:    emptyUrl,
		Body:   io.NopCloser(bytes.NewBufferString("")),
		Response: &http.Response{
			StatusCode: 0,
		},
		Header:           http.Header{},
		Proto:            "",
		ProtoMajor:       0,
		ProtoMinor:       0,
		TransferEncoding: []string{},
	}
	emptyCtx := New(emptyRequest)

	tests := []func() (string, lang.Function){
		methodFn,
		pathFn,
		bodyFn,
		statusFn,
		protoFn,
		protoMajorFn,
		protoMinorFn,
	}

	expected := []lang.Value{
		lang.StringValue(""),
		lang.StringValue(""),
		lang.StringValue(""),
		lang.NumberValue(0),
		lang.StringValue(""),
		lang.NumberValue(0),
		lang.NumberValue(0),
	}

	for it, tt := range tests {
		name, fn := tt()
		t.Run(name+"_missing_data", func(t *testing.T) {
			result, err := fn([]lang.Value{emptyCtx})
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if !reflect.DeepEqual(result, expected[it]) {
				t.Errorf("Expected %v, got %v", expected[it], result)
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
		"authorization", "accept", "trailer", "trailers", "routeValues",
		"pattern", "proto", "protoMajor", "protoMinor", "transferEncoding", "url",
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
	ctx := mockRequest()
	args := []lang.Value{ctx, lang.StringValue("Content-Type")}

	for i := 0; i < b.N; i++ {
		fn(args)
	}
}

func BenchmarkTrailer(b *testing.B) {
	_, fn := trailerFn()
	ctx := mockRequest()
	args := []lang.Value{ctx, lang.StringValue("X-Trailer-Test")}

	for i := 0; i < b.N; i++ {
		fn(args)
	}
}

func BenchmarkIP(b *testing.B) {
	_, fn := ipFn()
	ctx := mockRequest()
	args := []lang.Value{ctx}

	for i := 0; i < b.N; i++ {
		fn(args)
	}
}

func BenchmarkCookies(b *testing.B) {
	_, fn := cookiesFn()
	ctx := mockResponse()
	args := []lang.Value{ctx}

	for i := 0; i < b.N; i++ {
		fn(args)
	}
}

func BenchmarkTransferEncoding(b *testing.B) {
	_, fn := transferEncodingFn()
	ctx := mockRequest()
	args := []lang.Value{ctx}

	for i := 0; i < b.N; i++ {
		fn(args)
	}
}

func BenchmarkUrl(b *testing.B) {
	_, fn := urlFn()
	ctx := mockResponse()
	args := []lang.Value{ctx}

	for i := 0; i < b.N; i++ {
		fn(args)
	}
}
