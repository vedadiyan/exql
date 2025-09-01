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
package exql

import (
	"testing"

	"github.com/vedadiyan/exql/lang"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		expectErr  bool
	}{
		// Valid expressions
		{"number", "42", false},
		{"string", "'hello'", false},
		{"boolean", "true", false},
		{"variable", "x", false},
		{"arithmetic", "3 + 5", false},
		{"comparison", "x > 10", false},
		{"logical", "true and false", false},
		{"function call", "func(1, 2)", false},
		{"field access", "obj.field", false},
		{"index access", "arr[0]", false},
		{"complex expression", "(x + 5) * 2 > threshold and active", false},
		{"list literal", "[1, 2, 3]", false},
		{"nested access", "users[0].name", false},

		// Invalid expressions
		{"syntax error", "3 + +", true},
		{"unclosed paren", "(3 + 5", true},
		{"unclosed string", "'hello", true},
		{"invalid operator", "3 // 5", true},
		{"empty expression", "", true},
		{"invalid identifier", "123abc", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, err := Parse(tt.expression)

			if tt.expectErr {
				if err == nil {
					t.Error("expected parsing error but got none")
				}
				if ast != nil {
					t.Error("expected nil AST on parse error")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected parsing error: %v", err)
			}

			if ast == nil {
				t.Error("expected non-nil AST for valid expression")
			}
		})
	}
}

func TestEval(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		setup      func(*DefaultContext)
		expected   lang.Value
		expectErr  bool
	}{
		{
			"literal number",
			"42",
			nil,
			lang.NumberValue(42),
			false,
		},
		{
			"literal string",
			"'hello world'",
			nil,
			lang.StringValue("hello world"),
			false,
		},
		{
			"literal boolean true",
			"true",
			nil,
			lang.BoolValue(true),
			false,
		},
		{
			"arithmetic expression",
			"3 + 5 * 2",
			nil,
			lang.NumberValue(13),
			false,
		},
		{
			"comparison true",
			"10 > 5",
			nil,
			lang.BoolValue(true),
			false,
		},
		{
			"logical and",
			"true and false",
			nil,
			lang.BoolValue(false),
			false,
		},
		{
			"variable access",
			"x",
			func(ctx *DefaultContext) {
				ctx.SetVariable("x", lang.NumberValue(100))
			},
			lang.NumberValue(100),
			false,
		},
		{
			"variable in expression",
			"x + 20",
			func(ctx *DefaultContext) {
				ctx.SetVariable("x", lang.NumberValue(30))
			},
			lang.NumberValue(50),
			false,
		},
		{
			"field access",
			"user.name",
			func(ctx *DefaultContext) {
				ctx.SetVariable("user", lang.MapValue{
					"name": lang.StringValue("John"),
					"age":  lang.NumberValue(25),
				})
			},
			lang.StringValue("John"),
			false,
		},
		{
			"index access",
			"items[1]",
			func(ctx *DefaultContext) {
				ctx.SetVariable("items", lang.ListValue{
					lang.StringValue("first"),
					lang.StringValue("second"),
					lang.StringValue("third"),
				})
			},
			lang.StringValue("second"),
			false,
		},
		{
			"function call",
			"add(10, 20)",
			func(ctx *DefaultContext) {
				ctx.SetFunction("add", func(args []lang.Value) (lang.Value, error) {
					a := toNumber(args[0])
					b := toNumber(args[1])
					return lang.NumberValue(a + b), nil
				})
			},
			lang.NumberValue(30),
			false,
		},
		{
			"undefined variable",
			"undefined_var",
			nil,
			nil,
			false,
		},
		{
			"undefined function returns false",
			"undefined_func(1, 2)",
			nil,
			lang.BoolValue(false),
			false,
		},
		{
			"complex expression",
			"(user.age > 18) and user.active",
			func(ctx *DefaultContext) {
				ctx.SetVariable("user", lang.MapValue{
					"age":    lang.NumberValue(25),
					"active": lang.BoolValue(true),
				})
			},
			lang.BoolValue(true),
			false,
		},
		{
			"list creation",
			"[1, 'hello', true]",
			nil,
			lang.ListValue{
				lang.NumberValue(1),
				lang.StringValue("hello"),
				lang.BoolValue(true),
			},
			false,
		},
		{
			"membership test in",
			"'apple' in fruits",
			func(ctx *DefaultContext) {
				ctx.SetVariable("fruits", lang.ListValue{
					lang.StringValue("apple"),
					lang.StringValue("banana"),
					lang.StringValue("cherry"),
				})
			},
			lang.BoolValue(true),
			false,
		},
		{
			"membership test not in",
			"'grape' not in fruits",
			func(ctx *DefaultContext) {
				ctx.SetVariable("fruits", lang.ListValue{
					lang.StringValue("apple"),
					lang.StringValue("banana"),
				})
			},
			lang.BoolValue(true),
			false,
		},
		{
			"parse error",
			"3 + +",
			nil,
			nil,
			true,
		},
		{
			"evaluation error",
			"arr[10]",
			func(ctx *DefaultContext) {
				ctx.SetVariable("arr", lang.ListValue{
					lang.StringValue("only"),
					lang.StringValue("two"),
				})
			},
			nil,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := NewDefaultContext()
			if tt.setup != nil {
				tt.setup(ctx)
			}

			result, err := Eval(tt.expression, ctx)

			if tt.expectErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !valueEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestParseAndEvalSeparately(t *testing.T) {
	ctx := NewDefaultContext()
	ctx.SetVariable("x", lang.NumberValue(10))

	// Parse once
	ast, err := Parse("x * 2 + 5")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	// Evaluate multiple times with different contexts
	result1, err := ast.Evaluate(ctx)
	if err != nil {
		t.Fatalf("evaluation error: %v", err)
	}

	if !valueEqual(result1, lang.NumberValue(25)) {
		t.Errorf("expected 25, got %v", result1)
	}

	// Change variable and re-evaluate
	ctx.SetVariable("x", lang.NumberValue(3))
	result2, err := ast.Evaluate(ctx)
	if err != nil {
		t.Fatalf("evaluation error: %v", err)
	}

	if !valueEqual(result2, lang.NumberValue(11)) {
		t.Errorf("expected 11, got %v", result2)
	}
}

func TestBuiltInLibraries(t *testing.T) {
	ctx := NewDefaultContext(WithBuiltInLibrary())

	// Test that built-in libraries are available
	libraries := []string{"string", "util", "time", "json", "list", "map", "url", "http", "crypt", "ip"}

	for _, lib := range libraries {
		t.Run("library_"+lib, func(t *testing.T) {
			result, err := Eval(lib, ctx)
			if err != nil {
				t.Fatalf("error accessing library %s: %v", lib, err)
			}

			if result == nil {
				t.Errorf("library %s should not be nil", lib)
			}

			// Should be a DefaultContext
			if _, ok := result.(*DefaultContext); !ok {
				t.Errorf("library %s should be a DefaultContext, got %T", lib, result)
			}
		})
	}
}

func TestRealWorldScenario(t *testing.T) {
	// Simulate a real-world user authorization scenario
	ctx := NewDefaultContext()

	// Set up user data
	ctx.SetVariable("user", lang.MapValue{
		"id":     lang.NumberValue(123),
		"name":   lang.StringValue("Alice"),
		"role":   lang.StringValue("admin"),
		"active": lang.BoolValue(true),
		"permissions": lang.ListValue{
			lang.StringValue("read"),
			lang.StringValue("write"),
			lang.StringValue("delete"),
		},
	})

	ctx.SetVariable("resource", lang.MapValue{
		"type":     lang.StringValue("document"),
		"owner_id": lang.NumberValue(123),
		"public":   lang.BoolValue(false),
	})

	// Test various authorization rules
	authTests := []struct {
		name     string
		rule     string
		expected bool
		scenario string
	}{
		{
			"admin access",
			"user.role == 'admin'",
			true,
			"Admin users should have access",
		},
		{
			"owner access",
			"user.id == resource.owner_id",
			true,
			"Resource owners should have access",
		},
		{
			"active user requirement",
			"user.active",
			true,
			"Only active users should have access",
		},
		{
			"write permission check",
			"'write' in user.permissions",
			true,
			"User should have write permission",
		},
		{
			"complex authorization",
			"user.active and (user.role == 'admin' or user.id == resource.owner_id) and 'write' in user.permissions",
			true,
			"Complex rule should pass for admin owner with write permission",
		},
		{
			"public resource access",
			"resource.public or user.active",
			true,
			"Should have access to public resources or if user is active",
		},
	}

	for _, tt := range authTests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Eval(tt.rule, ctx)
			if err != nil {
				t.Fatalf("error in scenario '%s': %v", tt.scenario, err)
			}

			boolResult, ok := result.(lang.BoolValue)
			if !ok {
				t.Fatalf("expected boolean result for authorization rule, got %T", result)
			}

			if bool(boolResult) != tt.expected {
				t.Errorf("scenario '%s' failed: expected %v, got %v", tt.scenario, tt.expected, bool(boolResult))
			}
		})
	}
}

// Helper functions for testing
func valueEqual(a, b lang.Value) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	switch va := a.(type) {
	case lang.BoolValue:
		if vb, ok := b.(lang.BoolValue); ok {
			return bool(va) == bool(vb)
		}
	case lang.NumberValue:
		if vb, ok := b.(lang.NumberValue); ok {
			return float64(va) == float64(vb)
		}
	case lang.StringValue:
		if vb, ok := b.(lang.StringValue); ok {
			return string(va) == string(vb)
		}
	case lang.ListValue:
		if vb, ok := b.(lang.ListValue); ok {
			if len(va) != len(vb) {
				return false
			}
			for i := range va {
				if !valueEqual(va[i], vb[i]) {
					return false
				}
			}
			return true
		}
	case lang.MapValue:
		if vb, ok := b.(lang.MapValue); ok {
			if len(va) != len(vb) {
				return false
			}
			for key, valueA := range va {
				if valueB, exists := vb[key]; exists {
					if !valueEqual(valueA, valueB) {
						return false
					}
				} else {
					return false
				}
			}
			return true
		}
	}
	return false
}

func toNumber(v lang.Value) float64 {
	switch val := v.(type) {
	case lang.NumberValue:
		return float64(val)
	case lang.StringValue:
		// Simplified conversion for testing
		return 0
	case lang.BoolValue:
		if val {
			return 1
		}
		return 0
	}
	return 0
}
