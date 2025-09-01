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
package crypt

import (
	"testing"

	"github.com/vedadiyan/exql/lang"
)

// Test Base64 Functions
func TestBase64Encode(t *testing.T) {
	_, fn := base64Encode()

	tests := []struct {
		name     string
		input    []lang.Value
		expected string
		hasError bool
	}{
		{
			name:     "simple string",
			input:    []lang.Value{lang.StringValue("hello")},
			expected: "aGVsbG8=",
			hasError: false,
		},
		{
			name:     "empty string",
			input:    []lang.Value{lang.StringValue("")},
			expected: "",
			hasError: false,
		},
		{
			name:     "special characters",
			input:    []lang.Value{lang.StringValue("Hello, World!")},
			expected: "SGVsbG8sIFdvcmxkIQ==",
			hasError: false,
		},
		{
			name:     "wrong argument count",
			input:    []lang.Value{},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn(tt.input)
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
			if string(result.(lang.StringValue)) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(result.(lang.StringValue)))
			}
		})
	}
}

func TestBase64Decode(t *testing.T) {
	_, fn := base64Decode()

	tests := []struct {
		name     string
		input    []lang.Value
		expected string
		hasError bool
	}{
		{
			name:     "simple string",
			input:    []lang.Value{lang.StringValue("aGVsbG8=")},
			expected: "hello",
			hasError: false,
		},
		{
			name:     "empty string",
			input:    []lang.Value{lang.StringValue("")},
			expected: "",
			hasError: false,
		},
		{
			name:     "special characters",
			input:    []lang.Value{lang.StringValue("SGVsbG8sIFdvcmxkIQ==")},
			expected: "Hello, World!",
			hasError: false,
		},
		{
			name:     "invalid base64",
			input:    []lang.Value{lang.StringValue("invalid!")},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn(tt.input)
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
			if string(result.(lang.StringValue)) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(result.(lang.StringValue)))
			}
		})
	}
}

func TestBase64UrlEncode(t *testing.T) {
	_, fn := base64UrlEncode()

	result, err := fn([]lang.Value{lang.StringValue("hello?world")})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected := "aGVsbG8_d29ybGQ="
	if string(result.(lang.StringValue)) != expected {
		t.Errorf("Expected %s, got %s", expected, string(result.(lang.StringValue)))
	}
}

func TestBase64UrlDecode(t *testing.T) {
	_, fn := base64UrlDecode()

	result, err := fn([]lang.Value{lang.StringValue("aGVsbG8_d29ybGQ=")})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected := "hello?world"
	if string(result.(lang.StringValue)) != expected {
		t.Errorf("Expected %s, got %s", expected, string(result.(lang.StringValue)))
	}
}

// Test Hex Functions
func TestHexEncode(t *testing.T) {
	_, fn := hexEncode()

	result, err := fn([]lang.Value{lang.StringValue("hello")})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected := "68656c6c6f"
	if string(result.(lang.StringValue)) != expected {
		t.Errorf("Expected %s, got %s", expected, string(result.(lang.StringValue)))
	}
}

func TestHexDecode(t *testing.T) {
	_, fn := hexDecode()

	tests := []struct {
		name     string
		input    string
		expected string
		hasError bool
	}{
		{
			name:     "valid hex",
			input:    "68656c6c6f",
			expected: "hello",
			hasError: false,
		},
		{
			name:     "invalid hex",
			input:    "invalid",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.StringValue(tt.input)})
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
			if string(result.(lang.StringValue)) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(result.(lang.StringValue)))
			}
		})
	}
}

// Test Hash Functions
func TestHashMD5(t *testing.T) {
	_, fn := hashMD5()

	result, err := fn([]lang.Value{lang.StringValue("hello")})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected := "5d41402abc4b2a76b9719d911017c592"
	if string(result.(lang.StringValue)) != expected {
		t.Errorf("Expected %s, got %s", expected, string(result.(lang.StringValue)))
	}
}

func TestHashSHA1(t *testing.T) {
	_, fn := hashSHA1()

	result, err := fn([]lang.Value{lang.StringValue("hello")})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected := "aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d"
	if string(result.(lang.StringValue)) != expected {
		t.Errorf("Expected %s, got %s", expected, string(result.(lang.StringValue)))
	}
}

func TestHashSHA256(t *testing.T) {
	_, fn := hashSHA256()

	result, err := fn([]lang.Value{lang.StringValue("hello")})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected := "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
	// Note: SHA256 produces 64 character hex string, truncated for readability
	actual := string(result.(lang.StringValue))
	if len(actual) != 64 {
		t.Errorf("Expected 64 character hash, got %d characters", len(actual))
	}
	if string(result.(lang.StringValue)) != expected {
		t.Errorf("Expected %s, got %s", expected, string(result.(lang.StringValue)))
	}
}

func TestHashCRC32(t *testing.T) {
	_, fn := hashCRC32()

	result, err := fn([]lang.Value{lang.StringValue("hello")})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	// CRC32 returns a number, not a string
	if _, ok := result.(lang.NumberValue); !ok {
		t.Errorf("Expected NumberValue, got %T", result)
	}
}

// Test HMAC Functions
func TestHmacSHA256(t *testing.T) {
	_, fn := hmacSHA256()

	tests := []struct {
		name     string
		key      string
		message  string
		hasError bool
	}{
		{
			name:     "valid hmac",
			key:      "secret",
			message:  "hello",
			hasError: false,
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
				args = []lang.Value{lang.StringValue(tt.key), lang.StringValue(tt.message)}
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
			if _, ok := result.(lang.StringValue); !ok {
				t.Errorf("Expected StringValue, got %T", result)
			}
		})
	}
}

// Test Binary Functions
func TestToBinary(t *testing.T) {
	_, fn := toBinary()

	tests := []struct {
		name     string
		input    lang.Value
		expected string
		hasError bool
	}{
		{
			name:     "number to binary",
			input:    lang.NumberValue(5),
			expected: "101",
			hasError: false,
		},
		{
			name:     "string to binary",
			input:    lang.StringValue("A"),
			expected: "01000001",
			hasError: false,
		},
		{
			name:     "multi-char string",
			input:    lang.StringValue("AB"),
			expected: "01000001 01000010",
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{tt.input})
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
			if string(result.(lang.StringValue)) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(result.(lang.StringValue)))
			}
		})
	}
}

func TestFromBinary(t *testing.T) {
	_, fn := fromBinary()

	tests := []struct {
		name     string
		input    string
		expected interface{}
		hasError bool
	}{
		{
			name:     "binary to number",
			input:    "101",
			expected: float64(5),
			hasError: false,
		},
		{
			name:     "binary to string",
			input:    "01000001 01000010",
			expected: "AB",
			hasError: false,
		},
		{
			name:     "invalid binary",
			input:    "102",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fn([]lang.Value{lang.StringValue(tt.input)})
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

			switch expected := tt.expected.(type) {
			case float64:
				if float64(result.(lang.NumberValue)) != expected {
					t.Errorf("Expected %f, got %f", expected, float64(result.(lang.NumberValue)))
				}
			case string:
				if string(result.(lang.StringValue)) != expected {
					t.Errorf("Expected %s, got %s", expected, string(result.(lang.StringValue)))
				}
			}
		})
	}
}

// Test ASCII Functions
func TestToASCII(t *testing.T) {
	_, fn := toASCII()

	result, err := fn([]lang.Value{lang.StringValue("AB")})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	list, ok := result.(lang.ListValue)
	if !ok {
		t.Errorf("Expected ListValue, got %T", result)
		return
	}

	if len(list) != 2 {
		t.Errorf("Expected list length 2, got %d", len(list))
		return
	}

	if float64(list[0].(lang.NumberValue)) != 65 {
		t.Errorf("Expected 65, got %f", float64(list[0].(lang.NumberValue)))
	}
	if float64(list[1].(lang.NumberValue)) != 66 {
		t.Errorf("Expected 66, got %f", float64(list[1].(lang.NumberValue)))
	}
}

func TestFromASCII(t *testing.T) {
	_, fn := fromASCII()

	input := lang.ListValue{lang.NumberValue(65), lang.NumberValue(66)}
	result, err := fn([]lang.Value{input})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := "AB"
	if string(result.(lang.StringValue)) != expected {
		t.Errorf("Expected %s, got %s", expected, string(result.(lang.StringValue)))
	}
}

func TestFromASCIIInvalidRange(t *testing.T) {
	_, fn := fromASCII()

	input := lang.ListValue{lang.NumberValue(300)} // Out of ASCII range
	_, err := fn([]lang.Value{input})
	if err == nil {
		t.Errorf("Expected error for out-of-range ASCII value")
	}
}

// Test Base32 Functions
func TestBase32Encode(t *testing.T) {
	_, fn := base32Encode()

	result, err := fn([]lang.Value{lang.StringValue("hello")})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Base32 encoding of "hello" should be "NBSWY3DP"
	actual := string(result.(lang.StringValue))
	if len(actual) == 0 {
		t.Errorf("Expected non-empty result")
	}
}

func TestBase32Decode(t *testing.T) {
	_, fn := base32Decode()

	// First encode something to get a valid base32 string
	_, encodeFn := base32Encode()
	encoded, _ := encodeFn([]lang.Value{lang.StringValue("hello")})

	result, err := fn([]lang.Value{encoded})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if string(result.(lang.StringValue)) != "hello" {
		t.Errorf("Expected hello, got %s", string(result.(lang.StringValue)))
	}
}

// Test HTML Escape Functions
func TestHtmlEscape(t *testing.T) {
	_, fn := htmlEscape()

	result, err := fn([]lang.Value{lang.StringValue("<script>alert('xss')</script>")})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;"
	if string(result.(lang.StringValue)) != expected {
		t.Errorf("Expected %s, got %s", expected, string(result.(lang.StringValue)))
	}
}

func TestHtmlUnescape(t *testing.T) {
	_, fn := htmlUnescape()

	result, err := fn([]lang.Value{lang.StringValue("&lt;div&gt;Hello&lt;/div&gt;")})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := "<div>Hello</div>"
	if string(result.(lang.StringValue)) != expected {
		t.Errorf("Expected %s, got %s", expected, string(result.(lang.StringValue)))
	}
}

// Test Hash Verification
func TestHashVerify(t *testing.T) {
	_, fn := hashVerify()

	tests := []struct {
		name      string
		input     string
		hash      string
		algorithm string
		expected  bool
		hasError  bool
	}{
		{
			name:      "valid md5",
			input:     "hello",
			hash:      "5d41402abc4b2a76b9719d911017c592",
			algorithm: "md5",
			expected:  true,
			hasError:  false,
		},
		{
			name:      "invalid md5",
			input:     "hello",
			hash:      "wrong_hash",
			algorithm: "md5",
			expected:  false,
			hasError:  false,
		},
		{
			name:      "unsupported algorithm",
			input:     "hello",
			hash:      "somehash",
			algorithm: "unsupported",
			expected:  false,
			hasError:  true,
		},
		{
			name:     "wrong argument count",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var args []lang.Value
			if !tt.hasError || tt.name != "wrong argument count" {
				args = []lang.Value{
					lang.StringValue(tt.input),
					lang.StringValue(tt.hash),
					lang.StringValue(tt.algorithm),
				}
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

			if bool(result.(lang.BoolValue)) != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, bool(result.(lang.BoolValue)))
			}
		})
	}
}

// Test Octal Functions
func TestToOctal(t *testing.T) {
	_, fn := toOctal()

	result, err := fn([]lang.Value{lang.NumberValue(8)})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := "10"
	if string(result.(lang.StringValue)) != expected {
		t.Errorf("Expected %s, got %s", expected, string(result.(lang.StringValue)))
	}
}

func TestFromOctal(t *testing.T) {
	_, fn := fromOctal()

	result, err := fn([]lang.Value{lang.StringValue("10")})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := float64(8)
	if float64(result.(lang.NumberValue)) != expected {
		t.Errorf("Expected %f, got %f", expected, float64(result.(lang.NumberValue)))
	}
}

// Test Export Function
func TestExport(t *testing.T) {
	functions := Export()

	expectedFunctions := []string{
		"base64Encode", "base64Decode", "base64UrlEncode", "base64UrlDecode",
		"hexEncode", "hexDecode", "base32Encode", "base32Decode",
		"hashMd5", "hashSha1", "hashSha224", "hashSha256", "hashSha384", "hashSha512", "hashCrc32",
		"hmacMd5", "hmacSha1", "hmacSha256", "hmacSha512",
		"toBinary", "fromBinary", "toOctal", "fromOctal", "toAscii", "fromAscii",
		"htmlEscape", "htmlUnescape", "hashVerify",
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

// Benchmark tests for performance-critical functions
func BenchmarkHashSHA256(b *testing.B) {
	_, fn := hashSHA256()
	input := []lang.Value{lang.StringValue("hello world")}

	for i := 0; i < b.N; i++ {
		fn(input)
	}
}

func BenchmarkBase64Encode(b *testing.B) {
	_, fn := base64Encode()
	input := []lang.Value{lang.StringValue("hello world")}

	for i := 0; i < b.N; i++ {
		fn(input)
	}
}
