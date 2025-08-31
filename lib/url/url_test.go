package url

import (
	"net/url"
	"testing"

	"github.com/vedadiyan/exql/lang"
)

func TestArgumentErrors(t *testing.T) {
	functions := []struct {
		fn           func() (string, lang.Function)
		expectedArgs int
	}{
		{parse, 1},
		{urlEncode, 1},
		{urlDecode, 1},
		{urlHost, 1},
		{port, 1},
		{path, 1},
		{query, 1},
		{queryParam, 2},
		{fragment, 1},
		{scheme, 1},
		{user, 1},
		{build, 1},
		{isAbsolute, 1},
		{pathSegments, 1},
		{queryString, 1},
		{clean, 1},
	}

	for _, tf := range functions {
		name, fn := tf.fn()
		t.Run(name+"_wrong_args", func(t *testing.T) {
			var args []lang.Value
			if tf.expectedArgs == 0 {
				args = []lang.Value{lang.NumberValue(1)} // Add extra arg for 0-arg functions
			} else {
				args = make([]lang.Value, tf.expectedArgs-1) // One less than expected
				for i := range args {
					args[i] = lang.StringValue("test")
				}
			}
			_, err := fn(args)
			if err == nil {
				t.Errorf("Expected error for wrong argument count in %s", name)
			}
		})
	}

	// Test join function separately as it expects at least 2 arguments
	t.Run("join_wrong_args", func(t *testing.T) {
		_, fn := join()
		_, err := fn([]lang.Value{lang.StringValue("http://example.com")}) // Only 1 arg, needs at least 2
		if err == nil {
			t.Error("Expected error for join with only 1 argument")
		}
	})
}

func TestExport(t *testing.T) {
	functions := Export()

	expectedFunctions := []string{
		"parse", "encode", "decode", "host", "port", "path", "query", "query_param",
		"fragment", "scheme", "user", "build", "join", "is_absolute", "path_segments",
		"query_string", "clean",
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

func TestParse(t *testing.T) {
	_, fn := parse()

	tests := []struct {
		name     string
		input    string
		expected map[string]string
	}{
		{
			"complete URL",
			"https://user@example.com:8080/path/to/resource?param1=value1&param2=value2#section1",
			map[string]string{
				"scheme":   "https",
				"host":     "example.com",
				"port":     "8080",
				"path":     "/path/to/resource",
				"query":    "param1=value1&param2=value2",
				"fragment": "section1",
				"user":     "user",
			},
		},
		{
			"simple URL",
			"http://example.com",
			map[string]string{
				"scheme":   "http",
				"host":     "example.com",
				"port":     "",
				"path":     "",
				"query":    "",
				"fragment": "",
				"user":     "",
			},
		},
		{
			"URL with path only",
			"https://example.com/path",
			map[string]string{
				"scheme":   "https",
				"host":     "example.com",
				"port":     "",
				"path":     "/path",
				"query":    "",
				"fragment": "",
				"user":     "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.StringValue(tt.input)})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			resultMap, ok := result.(lang.MapValue)
			if !ok {
				t.Fatal("expected MapValue result")
			}
			for key, expectedValue := range tt.expected {
				if actual, exists := resultMap[key]; exists {
					if string(actual.(lang.StringValue)) != expectedValue {
						t.Errorf("key %s: expected %s, got %s", key, expectedValue, string(actual.(lang.StringValue)))
					}
				} else {
					t.Errorf("key %s not found in result", key)
				}
			}
		})
	}

	// Test invalid URL
	t.Run("invalid URL", func(t *testing.T) {
		_, err := fn([]lang.Value{lang.StringValue("://invalid")})
		if err == nil {
			t.Error("expected error for invalid URL")
		}
	})
}

func TestUrlEncode(t *testing.T) {
	_, fn := urlEncode()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple string", "hello world", "hello+world"},
		{"special chars", "hello@world.com", "hello%40world.com"},
		{"already encoded", "hello%20world", "hello%2520world"},
		{"empty string", "", ""},
		{"symbols", "test&param=value", "test%26param%3Dvalue"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.StringValue(tt.input)})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if string(result.(lang.StringValue)) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, string(result.(lang.StringValue)))
			}
		})
	}
}

func TestUrlDecode(t *testing.T) {
	_, fn := urlDecode()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"encoded spaces", "hello+world", "hello world"},
		{"encoded symbols", "hello%40world.com", "hello@world.com"},
		{"encoded query", "test%26param%3Dvalue", "test&param=value"},
		{"already decoded", "hello world", "hello world"},
		{"empty string", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.StringValue(tt.input)})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if string(result.(lang.StringValue)) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, string(result.(lang.StringValue)))
			}
		})
	}

	// Test invalid encoded string
	t.Run("invalid encoding", func(t *testing.T) {
		_, err := fn([]lang.Value{lang.StringValue("hello%ZZ")})
		if err == nil {
			t.Error("expected error for invalid URL encoding")
		}
	})
}

func TestUrlHost(t *testing.T) {
	_, fn := urlHost()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple host", "http://example.com", "example.com"},
		{"host with port", "https://example.com:8080", "example.com"},
		{"localhost", "http://localhost:3000", "localhost"},
		{"IP address", "http://192.168.1.1:8080", "192.168.1.1"},
		{"no port", "https://api.example.com", "api.example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.StringValue(tt.input)})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if string(result.(lang.StringValue)) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, string(result.(lang.StringValue)))
			}
		})
	}
}

func TestPort(t *testing.T) {
	_, fn := port()

	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		{"explicit port", "http://example.com:8080", 8080},
		{"https default", "https://example.com", 443},
		{"http default", "http://example.com", 80},
		{"ftp default", "ftp://example.com", 21},
		{"ssh default", "ssh://example.com", 22},
		{"unknown scheme", "custom://example.com", 0},
		{"localhost with port", "http://localhost:3000", 3000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.StringValue(tt.input)})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if float64(result.(lang.NumberValue)) != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, float64(result.(lang.NumberValue)))
			}
		})
	}

	// Test invalid port
	t.Run("invalid port", func(t *testing.T) {
		_, err := fn([]lang.Value{lang.StringValue("http://example.com:abc")})
		if err == nil {
			t.Error("expected error for invalid port")
		}
	})
}

func TestPath(t *testing.T) {
	_, fn := path()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple path", "http://example.com/path/to/resource", "/path/to/resource"},
		{"root path", "http://example.com/", "/"},
		{"no path", "http://example.com", ""},
		{"complex path", "https://api.example.com/v1/users/123", "/v1/users/123"},
		{"path with query", "http://example.com/search?q=test", "/search"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.StringValue(tt.input)})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if string(result.(lang.StringValue)) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, string(result.(lang.StringValue)))
			}
		})
	}
}

func TestQuery(t *testing.T) {
	_, fn := query()

	tests := []struct {
		name     string
		input    string
		expected map[string]interface{}
	}{
		{
			"single param",
			"http://example.com?param=value",
			map[string]interface{}{"param": "value"},
		},
		{
			"multiple params",
			"http://example.com?param1=value1&param2=value2",
			map[string]interface{}{"param1": "value1", "param2": "value2"},
		},
		{
			"no query",
			"http://example.com",
			map[string]interface{}{},
		},
		{
			"array param",
			"http://example.com?tags=tag1&tags=tag2",
			map[string]interface{}{"tags": []string{"tag1", "tag2"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.StringValue(tt.input)})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			resultMap, ok := result.(lang.MapValue)
			if !ok {
				t.Fatal("expected MapValue result")
			}

			if len(resultMap) != len(tt.expected) {
				t.Errorf("expected %d params, got %d", len(tt.expected), len(resultMap))
			}

			for key, expectedValue := range tt.expected {
				actual, exists := resultMap[key]
				if !exists {
					t.Errorf("key %s not found in result", key)
					continue
				}

				switch exp := expectedValue.(type) {
				case string:
					if string(actual.(lang.StringValue)) != exp {
						t.Errorf("key %s: expected %s, got %s", key, exp, string(actual.(lang.StringValue)))
					}
				case []string:
					actualList, ok := actual.(lang.ListValue)
					if !ok {
						t.Errorf("key %s: expected list, got %T", key, actual)
						continue
					}
					if len(actualList) != len(exp) {
						t.Errorf("key %s: expected %d items, got %d", key, len(exp), len(actualList))
						continue
					}
					for i, expItem := range exp {
						if string(actualList[i].(lang.StringValue)) != expItem {
							t.Errorf("key %s[%d]: expected %s, got %s", key, i, expItem, string(actualList[i].(lang.StringValue)))
						}
					}
				}
			}
		})
	}
}

func TestQueryParam(t *testing.T) {
	_, fn := queryParam()

	tests := []struct {
		name     string
		url      string
		param    string
		expected interface{}
	}{
		{"existing param", "http://example.com?param=value", "param", "value"},
		{"missing param", "http://example.com?other=value", "param", nil},
		{"array param", "http://example.com?tags=tag1&tags=tag2", "tags", []string{"tag1", "tag2"}},
		{"empty query", "http://example.com", "param", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.StringValue(tt.url), lang.StringValue(tt.param)})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.expected == nil {
				if result != nil {
					t.Errorf("expected nil, got %v", result)
				}
			} else {
				switch exp := tt.expected.(type) {
				case string:
					if string(result.(lang.StringValue)) != exp {
						t.Errorf("expected %s, got %s", exp, string(result.(lang.StringValue)))
					}
				case []string:
					actualList, ok := result.(lang.ListValue)
					if !ok {
						t.Errorf("expected list, got %T", result)
						break
					}
					if len(actualList) != len(exp) {
						t.Errorf("expected %d items, got %d", len(exp), len(actualList))
						break
					}
					for i, expItem := range exp {
						if string(actualList[i].(lang.StringValue)) != expItem {
							t.Errorf("item %d: expected %s, got %s", i, expItem, string(actualList[i].(lang.StringValue)))
						}
					}
				}
			}
		})
	}
}

func TestFragment(t *testing.T) {
	_, fn := fragment()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"with fragment", "http://example.com#section1", "section1"},
		{"no fragment", "http://example.com", ""},
		{"fragment with query", "http://example.com?param=value#section1", "section1"},
		{"empty fragment", "http://example.com#", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.StringValue(tt.input)})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if string(result.(lang.StringValue)) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, string(result.(lang.StringValue)))
			}
		})
	}
}

func TestScheme(t *testing.T) {
	_, fn := scheme()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"http", "http://example.com", "http"},
		{"https", "https://example.com", "https"},
		{"ftp", "ftp://example.com", "ftp"},
		{"custom", "custom://example.com", "custom"},
		{"no scheme", "//example.com", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.StringValue(tt.input)})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if string(result.(lang.StringValue)) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, string(result.(lang.StringValue)))
			}
		})
	}
}

func TestUser(t *testing.T) {
	_, fn := user()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"with user", "http://user@example.com", "user"},
		{"no user", "http://example.com", ""},
		{"user with password", "http://user:pass@example.com", "user"},
		{"empty user", "http://@example.com", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.StringValue(tt.input)})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if string(result.(lang.StringValue)) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, string(result.(lang.StringValue)))
			}
		})
	}
}

func TestBuild(t *testing.T) {
	_, fn := build()

	tests := []struct {
		name     string
		input    lang.MapValue
		expected string
	}{
		{
			"complete URL",
			lang.MapValue{
				"scheme":   lang.StringValue("https"),
				"host":     lang.StringValue("example.com"),
				"port":     lang.StringValue("8080"),
				"path":     lang.StringValue("/path"),
				"query":    lang.StringValue("param=value"),
				"fragment": lang.StringValue("section1"),
			},
			"https://example.com:8080/path?param=value#section1",
		},
		{
			"minimal URL",
			lang.MapValue{
				"scheme": lang.StringValue("http"),
				"host":   lang.StringValue("example.com"),
			},
			"http://example.com",
		},
		{
			"with user",
			lang.MapValue{
				"scheme": lang.StringValue("https"),
				"host":   lang.StringValue("example.com"),
				"user":   lang.StringValue("user"),
			},
			"https://user@example.com",
		},
		{
			"with user and password",
			lang.MapValue{
				"scheme":   lang.StringValue("https"),
				"host":     lang.StringValue("example.com"),
				"user":     lang.StringValue("user"),
				"password": lang.StringValue("pass"),
			},
			"https://user:pass@example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{tt.input})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if string(result.(lang.StringValue)) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, string(result.(lang.StringValue)))
			}
		})
	}

	// Test non-map input
	t.Run("non-map input", func(t *testing.T) {
		_, err := fn([]lang.Value{lang.StringValue("not a map")})
		if err == nil {
			t.Error("expected error for non-map input")
		}
	})
}

func TestJoin(t *testing.T) {
	_, fn := join()

	tests := []struct {
		name     string
		base     string
		paths    []string
		expected string
	}{
		{
			"simple join",
			"http://example.com",
			[]string{"path", "to", "resource"},
			"http://example.com/path/to/resource",
		},
		{
			"with existing path",
			"http://example.com/api",
			[]string{"v1", "users"},
			"http://example.com/api/v1/users",
		},
		{
			"empty segments",
			"http://example.com",
			[]string{"", "path", "", "resource"},
			"http://example.com/path/resource",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := make([]lang.Value, len(tt.paths)+1)
			args[0] = lang.StringValue(tt.base)
			for i, path := range tt.paths {
				args[i+1] = lang.StringValue(path)
			}

			result, err := fn(args)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if string(result.(lang.StringValue)) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, string(result.(lang.StringValue)))
			}
		})
	}
}

func TestIsAbsolute(t *testing.T) {
	_, fn := isAbsolute()

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"absolute http", "http://example.com", true},
		{"absolute https", "https://example.com", true},
		{"absolute ftp", "ftp://example.com", true},
		{"relative path", "/path/to/resource", false},
		{"relative current", "./resource", false},
		{"relative parent", "../resource", false},
		{"protocol relative", "//example.com", false},
		{"query only", "?param=value", false},
		{"fragment only", "#section", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.StringValue(tt.input)})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if bool(result.(lang.BoolValue)) != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, bool(result.(lang.BoolValue)))
			}
		})
	}
}

func TestPathSegments(t *testing.T) {
	_, fn := pathSegments()

	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{"simple path", "/path/to/resource", []string{"path", "to", "resource"}},
		{"URL with path", "http://example.com/api/v1/users", []string{"api", "v1", "users"}},
		{"root path", "/", []string{}},
		{"no path", "http://example.com", []string{}},
		{"encoded segments", "/path/hello%20world", []string{"path", "hello world"}},
		{"trailing slash", "/path/to/", []string{"path", "to"}},
		{"relative path", "path/to/resource", []string{"path", "to", "resource"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.StringValue(tt.input)})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			resultList, ok := result.(lang.ListValue)
			if !ok {
				t.Fatal("expected ListValue result")
			}

			if len(resultList) != len(tt.expected) {
				t.Errorf("expected %d segments, got %d", len(tt.expected), len(resultList))
			}

			for i, expected := range tt.expected {
				if i < len(resultList) {
					if string(resultList[i].(lang.StringValue)) != expected {
						t.Errorf("segment %d: expected %s, got %s", i, expected, string(resultList[i].(lang.StringValue)))
					}
				}
			}
		})
	}
}

func TestQueryString(t *testing.T) {
	_, fn := queryString()

	tests := []struct {
		name     string
		input    lang.MapValue
		expected string
	}{
		{
			"single param",
			lang.MapValue{"param": lang.StringValue("value")},
			"param=value",
		},
		{
			"multiple params",
			lang.MapValue{
				"param1": lang.StringValue("value1"),
				"param2": lang.StringValue("value2"),
			},
			"param1=value1&param2=value2",
		},
		{
			"empty map",
			lang.MapValue{},
			"",
		},
		{
			"array param",
			lang.MapValue{
				"tags": lang.ListValue{
					lang.StringValue("tag1"),
					lang.StringValue("tag2"),
				},
			},
			"tags=tag1&tags=tag2",
		},
		{
			"mixed types",
			lang.MapValue{
				"name":   lang.StringValue("test"),
				"count":  lang.NumberValue(42),
				"active": lang.BoolValue(true),
			},
			"active=true&count=42&name=test", // Note: url.Values.Encode() sorts keys
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{tt.input})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if string(result.(lang.StringValue)) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, string(result.(lang.StringValue)))
			}
		})
	}

	// Test non-map input
	t.Run("non-map input", func(t *testing.T) {
		_, err := fn([]lang.Value{lang.StringValue("not a map")})
		if err == nil {
			t.Error("expected error for non-map input")
		}
	})
}

func TestClean(t *testing.T) {
	_, fn := clean()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple path", "http://example.com/path/to/resource", "http://example.com/path/to/resource"},
		{"with dots", "http://example.com/path/./to/../resource", "http://example.com/path/resource"},
		{"multiple dots", "http://example.com/path/../../../resource", "http://example.com/resource"},
		{"trailing dots", "http://example.com/path/to/.", "http://example.com/path/to"},
		{"root with dots", "http://example.com/../", "http://example.com/"},
		{"complex cleanup", "http://example.com/a/b/../c/./d", "http://example.com/a/c/d"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.StringValue(tt.input)})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if string(result.(lang.StringValue)) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, string(result.(lang.StringValue)))
			}
		})
	}
}

func TestCleanPath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty path", "", "/"},
		{"simple path", "/path/to/resource", "/path/to/resource"},
		{"with current dir", "/path/./to/resource", "/path/to/resource"},
		{"with parent dir", "/path/to/../resource", "/path/resource"},
		{"multiple parent dirs", "/path/to/../../resource", "/resource"},
		{"root parent", "/../resource", "/resource"},
		{"relative path", "path/to/resource", "path/to/resource"},
		{"relative with dots", "path/../resource", "resource"},
		{"only dots", "./", "/"},
		{"complex case", "/a/b/../c/./d/../e", "/a/c/e"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cleanpath(tt.input)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestGetUserInfo(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"nil userinfo", "", ""},
		{"with user", "user", "user"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.input == "" {
				result := getUserInfo(nil)
				if result != tt.expected {
					t.Errorf("expected %s, got %s", tt.expected, result)
				}
			} else {
				// For testing with actual userinfo, we need to create a URL and parse it
				testURL := "http://" + tt.input + "@example.com"
				u, err := url.Parse(testURL)
				if err != nil {
					t.Fatalf("failed to parse test URL: %v", err)
				}
				result := getUserInfo(u.User)
				if result != tt.expected {
					t.Errorf("expected %s, got %s", tt.expected, result)
				}
			}
		})
	}
}

// Benchmark Tests
func BenchmarkParse(b *testing.B) {
	_, fn := parse()
	input := []lang.Value{lang.StringValue("https://user@example.com:8080/path?param=value#fragment")}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fn(input)
	}
}

func BenchmarkUrlEncode(b *testing.B) {
	_, fn := urlEncode()
	input := []lang.Value{lang.StringValue("hello world & special chars @#$%")}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fn(input)
	}
}

func BenchmarkJoin(b *testing.B) {
	_, fn := join()
	input := []lang.Value{
		lang.StringValue("http://example.com"),
		lang.StringValue("api"),
		lang.StringValue("v1"),
		lang.StringValue("users"),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fn(input)
	}
}
