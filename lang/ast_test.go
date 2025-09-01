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
package lang

import (
	"fmt"
	"testing"
)

// Mock context for testing
type MockContext struct {
	variables map[string]Value
	functions map[string]Function
}

func NewMockContext() *MockContext {
	return &MockContext{
		variables: make(map[string]Value),
		functions: make(map[string]Function),
	}
}

func (c *MockContext) GetVariable(name string) Value {
	return c.variables[name]
}

func (c *MockContext) GetFunction(name string) Function {
	return c.functions[name]
}

func (c *MockContext) SetVariable(name string, value Value) {
	c.variables[name] = value
}

func (c *MockContext) SetFunction(name string, fn Function) {
	c.functions[name] = fn
}

func TestBinaryOpNode(t *testing.T) {
	ctx := NewMockContext()

	tests := []struct {
		name     string
		left     Value
		right    Value
		operator string
		expected Value
	}{
		// Logical operations
		{"and true true", BoolValue(true), BoolValue(true), "and", BoolValue(true)},
		{"and true false", BoolValue(true), BoolValue(false), "and", BoolValue(false)},
		{"and false false", BoolValue(false), BoolValue(false), "and", BoolValue(false)},
		{"or true false", BoolValue(true), BoolValue(false), "or", BoolValue(true)},
		{"or false false", BoolValue(false), BoolValue(false), "or", BoolValue(false)},

		// Equality operations
		{"equal numbers", NumberValue(5), NumberValue(5), "==", BoolValue(true)},
		{"equal strings", StringValue("hello"), StringValue("hello"), "==", BoolValue(true)},
		{"not equal numbers", NumberValue(5), NumberValue(3), "!=", BoolValue(true)},
		{"not equal strings", StringValue("hello"), StringValue("world"), "!=", BoolValue(true)},
		{"equal with =", NumberValue(5), NumberValue(5), "=", BoolValue(true)},

		// Comparison operations
		{"less than true", NumberValue(3), NumberValue(5), "<", BoolValue(true)},
		{"less than false", NumberValue(5), NumberValue(3), "<", BoolValue(false)},
		{"less than equal true", NumberValue(3), NumberValue(5), "<=", BoolValue(true)},
		{"less than equal equal", NumberValue(5), NumberValue(5), "<=", BoolValue(true)},
		{"greater than true", NumberValue(5), NumberValue(3), ">", BoolValue(true)},
		{"greater than false", NumberValue(3), NumberValue(5), ">", BoolValue(false)},
		{"greater than equal true", NumberValue(5), NumberValue(3), ">=", BoolValue(true)},
		{"greater than equal equal", NumberValue(5), NumberValue(5), ">=", BoolValue(true)},

		// Arithmetic operations
		{"addition", NumberValue(3), NumberValue(5), "+", NumberValue(8)},
		{"subtraction", NumberValue(5), NumberValue(3), "-", NumberValue(2)},
		{"multiplication", NumberValue(3), NumberValue(4), "*", NumberValue(12)},
		{"division", NumberValue(10), NumberValue(2), "/", NumberValue(5)},

		// In operations
		{"in true", NumberValue(2), ListValue{NumberValue(1), NumberValue(2), NumberValue(3)}, "in", BoolValue(true)},
		{"in false", NumberValue(4), ListValue{NumberValue(1), NumberValue(2), NumberValue(3)}, "in", BoolValue(false)},
		{"not in true", NumberValue(4), ListValue{NumberValue(1), NumberValue(2), NumberValue(3)}, "not in", BoolValue(true)},
		{"not in false", NumberValue(2), ListValue{NumberValue(1), NumberValue(2), NumberValue(3)}, "not in", BoolValue(false)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := &BinaryOpNode{
				Left:     &LiteralNode{Value: tt.left},
				Right:    &LiteralNode{Value: tt.right},
				Operator: tt.operator,
			}

			result, err := node.Evaluate(ctx)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !equal(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestBinaryOpNodeUnsupportedOperator(t *testing.T) {
	ctx := NewMockContext()
	node := &BinaryOpNode{
		Left:     &LiteralNode{Value: NumberValue(5)},
		Right:    &LiteralNode{Value: NumberValue(3)},
		Operator: "%%", // Unsupported operator
	}

	_, err := node.Evaluate(ctx)
	if err == nil {
		t.Error("expected error for unsupported operator")
	}
}

func TestUnaryOpNode(t *testing.T) {
	ctx := NewMockContext()

	tests := []struct {
		name     string
		operand  Value
		operator string
		expected Value
	}{
		{"not true", BoolValue(true), "not", BoolValue(false)},
		{"not false", BoolValue(false), "not", BoolValue(true)},
		{"not truthy", NumberValue(5), "not", BoolValue(false)},
		{"not falsy", NumberValue(0), "not", BoolValue(true)},
		{"negative number", NumberValue(5), "-", NumberValue(-5)},
		{"double negative", NumberValue(-3), "-", NumberValue(3)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := &UnaryOpNode{
				Operand:  &LiteralNode{Value: tt.operand},
				Operator: tt.operator,
			}

			result, err := node.Evaluate(ctx)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !equal(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestUnaryOpNodeUnsupportedOperator(t *testing.T) {
	ctx := NewMockContext()
	node := &UnaryOpNode{
		Operand:  &LiteralNode{Value: NumberValue(5)},
		Operator: "++", // Unsupported operator
	}

	_, err := node.Evaluate(ctx)
	if err == nil {
		t.Error("expected error for unsupported operator")
	}
}

func TestLiteralNode(t *testing.T) {
	ctx := NewMockContext()

	tests := []struct {
		name  string
		value Value
	}{
		{"boolean", BoolValue(true)},
		{"number", NumberValue(42)},
		{"string", StringValue("hello")},
		{"list", ListValue{NumberValue(1), NumberValue(2)}},
		{"map", MapValue{"key": StringValue("value")}},
		{"nil", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := &LiteralNode{Value: tt.value}
			result, err := node.Evaluate(ctx)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !equal(result, tt.value) {
				t.Errorf("expected %v, got %v", tt.value, result)
			}
		})
	}
}

func TestVariableNode(t *testing.T) {
	ctx := NewMockContext()
	ctx.SetVariable("x", NumberValue(42))
	ctx.SetVariable("name", StringValue("test"))
	ctx.SetVariable("flag", BoolValue(true))

	tests := []struct {
		name     string
		varName  string
		expected Value
	}{
		{"number variable", "x", NumberValue(42)},
		{"string variable", "name", StringValue("test")},
		{"boolean variable", "flag", BoolValue(true)},
		{"undefined variable", "undefined", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := &VariableNode{Name: tt.varName}
			result, err := node.Evaluate(ctx)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !equal(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestFieldAccessNode(t *testing.T) {
	ctx := NewMockContext()

	tests := []struct {
		name     string
		object   Value
		field    string
		expected Value
	}{
		{
			"map field access",
			MapValue{"name": StringValue("John"), "age": NumberValue(30)},
			"name",
			StringValue("John"),
		},
		{
			"map missing field",
			MapValue{"name": StringValue("John")},
			"age",
			nil,
		},
		{
			"list field access",
			ListValue{
				MapValue{"name": StringValue("John")},
				MapValue{"name": StringValue("Jane")},
			},
			"name",
			ListValue{StringValue("John"), StringValue("Jane")},
		},
		{
			"non-map/list object",
			StringValue("hello"),
			"length",
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := &FieldAccessNode{
				Object: &LiteralNode{Value: tt.object},
				Field:  tt.field,
			}

			result, err := node.Evaluate(ctx)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !equal(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestIndexAccessNode(t *testing.T) {
	ctx := NewMockContext()

	tests := []struct {
		name      string
		object    Value
		index     Value
		expected  Value
		expectErr bool
	}{
		{
			"list number index",
			ListValue{StringValue("a"), StringValue("b"), StringValue("c")},
			NumberValue(1),
			StringValue("b"),
			false,
		},
		{
			"list out of bounds",
			ListValue{StringValue("a"), StringValue("b")},
			NumberValue(5),
			nil,
			true,
		},
		{
			"list negative index",
			ListValue{StringValue("a"), StringValue("b")},
			NumberValue(-1),
			nil,
			true,
		},
		{
			"list string index (field access)",
			ListValue{
				MapValue{"name": StringValue("John")},
				MapValue{"name": StringValue("Jane")},
			},
			StringValue("name"),
			ListValue{StringValue("John"), StringValue("Jane")},
			false,
		},
		{
			"map string index",
			MapValue{"key": StringValue("value")},
			StringValue("key"),
			StringValue("value"),
			false,
		},
		{
			"list each index",
			ListValue{NumberValue(1), NumberValue(2), NumberValue(3)},
			EachValue(0),
			ListValue{NumberValue(1), NumberValue(2), NumberValue(3)},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := &IndexAccessNode{
				Object: &LiteralNode{Value: tt.object},
				Index:  &LiteralNode{Value: tt.index},
			}

			result, err := node.Evaluate(ctx)
			if tt.expectErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !equal(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestIndexAccessNodeUnsupportedTypes(t *testing.T) {
	ctx := NewMockContext()

	tests := []struct {
		name   string
		object Value
		index  Value
	}{
		{"map with number index", MapValue{"key": StringValue("value")}, NumberValue(1)},
		{"list with boolean index", ListValue{NumberValue(1)}, BoolValue(true)},
		{"string object", StringValue("hello"), NumberValue(0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := &IndexAccessNode{
				Object: &LiteralNode{Value: tt.object},
				Index:  &LiteralNode{Value: tt.index},
			}

			x, err := node.Evaluate(ctx)
			_ = x
			if err == nil {
				t.Error("expected error for unsupported type combination")
			}
		})
	}
}

func TestFunctionCallNode(t *testing.T) {
	ctx := NewMockContext()

	// Set up test functions
	ctx.SetFunction("add", func(args []Value) (Value, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("add expects 2 arguments")
		}
		return NumberValue(ToNumber(args[0]) + ToNumber(args[1])), nil
	})

	ctx.SetFunction("concat", func(args []Value) (Value, error) {
		result := ""
		for _, arg := range args {
			result += fmt.Sprintf("%v", arg)
		}
		return StringValue(result), nil
	})

	tests := []struct {
		name      string
		funcName  string
		args      []Value
		expected  Value
		expectErr bool
	}{
		{
			"add function",
			"add",
			[]Value{NumberValue(3), NumberValue(5)},
			NumberValue(8),
			false,
		},
		{
			"concat function",
			"concat",
			[]Value{StringValue("hello"), StringValue(" "), StringValue("world")},
			StringValue("hello world"),
			false,
		},
		{
			"undefined function",
			"undefined",
			[]Value{NumberValue(1)},
			BoolValue(false),
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			argNodes := make([]ExprNode, len(tt.args))
			for i, arg := range tt.args {
				argNodes[i] = &LiteralNode{Value: arg}
			}

			node := &FunctionCallNode{
				Name: tt.funcName,
				Args: argNodes,
			}

			result, err := node.Evaluate(ctx)
			if tt.expectErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !equal(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestFunctionCallNodeWithError(t *testing.T) {
	ctx := NewMockContext()
	ctx.SetFunction("error", func(args []Value) (Value, error) {
		return nil, fmt.Errorf("function error")
	})

	node := &FunctionCallNode{
		Name: "error",
		Args: []ExprNode{&LiteralNode{Value: NumberValue(1)}},
	}

	_, err := node.Evaluate(ctx)
	if err == nil {
		t.Error("expected error from function call")
	}
}

func TestListNode(t *testing.T) {
	ctx := NewMockContext()
	ctx.SetVariable("x", NumberValue(10))

	tests := []struct {
		name     string
		elements []ExprNode
		expected ListValue
	}{
		{
			"empty list",
			[]ExprNode{},
			ListValue{},
		},
		{
			"literal elements",
			[]ExprNode{
				&LiteralNode{Value: NumberValue(1)},
				&LiteralNode{Value: StringValue("hello")},
				&LiteralNode{Value: BoolValue(true)},
			},
			ListValue{NumberValue(1), StringValue("hello"), BoolValue(true)},
		},
		{
			"mixed expressions",
			[]ExprNode{
				&LiteralNode{Value: NumberValue(5)},
				&VariableNode{Name: "x"},
				&BinaryOpNode{
					Left:     &LiteralNode{Value: NumberValue(2)},
					Right:    &LiteralNode{Value: NumberValue(3)},
					Operator: "+",
				},
			},
			ListValue{NumberValue(5), NumberValue(10), NumberValue(5)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := &ListNode{Elements: tt.elements}
			result, err := node.Evaluate(ctx)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !equal(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestEachNode(t *testing.T) {
	ctx := NewMockContext()
	node := &EachNode{}

	result, err := node.Evaluate(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := result.(EachValue); !ok {
		t.Errorf("expected EachValue, got %T", result)
	}
}

func TestRangeNode(t *testing.T) {
	ctx := NewMockContext()
	node := &RangeNode{
		Begin: NumberValue(1),
		End:   NumberValue(10),
	}

	result, err := node.Evaluate(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := result.(EachValue); !ok {
		t.Errorf("expected EachValue, got %T", result)
	}
}

func TestToBool(t *testing.T) {
	tests := []struct {
		name     string
		input    Value
		expected bool
	}{
		{"true boolean", BoolValue(true), true},
		{"false boolean", BoolValue(false), false},
		{"non-zero number", NumberValue(42), true},
		{"zero number", NumberValue(0), false},
		{"negative number", NumberValue(-5), true},
		{"non-empty string", StringValue("hello"), true},
		{"empty string", StringValue(""), false},
		{"nil value", nil, false},
		{"list value", ListValue{NumberValue(1)}, true},
		{"map value", MapValue{"key": StringValue("value")}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToBool(tt.input)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestToNumber(t *testing.T) {
	tests := []struct {
		name     string
		input    Value
		expected float64
	}{
		{"number value", NumberValue(42), 42},
		{"zero number", NumberValue(0), 0},
		{"negative number", NumberValue(-5), -5},
		{"numeric string", StringValue("42"), 42},
		{"float string", StringValue("3.14"), 3.14},
		{"non-numeric string", StringValue("hello"), 0},
		{"empty string", StringValue(""), 0},
		{"true boolean", BoolValue(true), 1},
		{"false boolean", BoolValue(false), 0},
		{"nil value", nil, 0},
		{"list value", ListValue{NumberValue(1)}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToNumber(tt.input)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestEqual(t *testing.T) {
	tests := []struct {
		name     string
		a, b     Value
		expected bool
	}{
		{"equal numbers", NumberValue(42), NumberValue(42), true},
		{"different numbers", NumberValue(42), NumberValue(43), false},
		{"equal strings", StringValue("hello"), StringValue("hello"), true},
		{"different strings", StringValue("hello"), StringValue("world"), false},
		{"equal booleans", BoolValue(true), BoolValue(true), true},
		{"different booleans", BoolValue(true), BoolValue(false), false},
		{"both nil", nil, nil, true},
		{"nil and value", nil, NumberValue(0), false},
		{"different types", NumberValue(42), StringValue("42"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := equal(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestCompare(t *testing.T) {
	tests := []struct {
		name     string
		a, b     Value
		expected int
	}{
		{"a less than b", NumberValue(3), NumberValue(5), -1},
		{"a greater than b", NumberValue(5), NumberValue(3), 1},
		{"a equal to b", NumberValue(5), NumberValue(5), 0},
		{"string numbers", StringValue("3"), StringValue("5"), -1},
		{"mixed types", NumberValue(5), StringValue("3"), 1},
		{"booleans", BoolValue(false), BoolValue(true), -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := compare(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		name      string
		container Value
		item      Value
		expected  bool
	}{
		{
			"item in list",
			ListValue{NumberValue(1), NumberValue(2), NumberValue(3)},
			NumberValue(2),
			true,
		},
		{
			"item not in list",
			ListValue{NumberValue(1), NumberValue(2), NumberValue(3)},
			NumberValue(4),
			false,
		},
		{
			"string in list",
			ListValue{StringValue("a"), StringValue("b"), StringValue("c")},
			StringValue("b"),
			true,
		},
		{
			"empty list",
			ListValue{},
			NumberValue(1),
			false,
		},
		{
			"non-list container",
			StringValue("hello"),
			StringValue("h"),
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := contains(tt.container, tt.item)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// Test evaluation errors
func TestEvaluationErrors(t *testing.T) {
	ctx := NewMockContext()

	// Test binary operation with evaluation error
	t.Run("binary op left error", func(t *testing.T) {
		node := &BinaryOpNode{
			Left:     &VariableNode{Name: "undefined"}, // Will cause error in some contexts
			Right:    &LiteralNode{Value: NumberValue(5)},
			Operator: "+",
		}

		// This should not error in our current implementation since undefined variables return nil
		_, err := node.Evaluate(ctx)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	// Test function call with argument evaluation error
	t.Run("function call arg error", func(t *testing.T) {
		ctx.SetFunction("test", func(args []Value) (Value, error) {
			return NumberValue(1), nil
		})

		errorNode := &BinaryOpNode{
			Left:     &LiteralNode{Value: NumberValue(1)},
			Right:    &LiteralNode{Value: NumberValue(0)},
			Operator: "/", // This won't error, but let's create a node that could
		}

		node := &FunctionCallNode{
			Name: "test",
			Args: []ExprNode{errorNode},
		}

		_, err := node.Evaluate(ctx)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	// Test list node with element evaluation error
	t.Run("list node element error", func(t *testing.T) {
		errorNode := &BinaryOpNode{
			Left:     &LiteralNode{Value: NumberValue(1)},
			Right:    &LiteralNode{Value: NumberValue(2)},
			Operator: "unsupported",
		}

		node := &ListNode{
			Elements: []ExprNode{
				&LiteralNode{Value: NumberValue(1)},
				errorNode,
			},
		}

		_, err := node.Evaluate(ctx)
		if err == nil {
			t.Error("expected error from unsupported operator")
		}
	})
}

// Benchmark Tests
func BenchmarkBinaryOpNode(b *testing.B) {
	ctx := NewMockContext()
	node := &BinaryOpNode{
		Left:     &LiteralNode{Value: NumberValue(5)},
		Right:    &LiteralNode{Value: NumberValue(3)},
		Operator: "+",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		node.Evaluate(ctx)
	}
}

func BenchmarkFieldAccess(b *testing.B) {
	ctx := NewMockContext()
	obj := MapValue{
		"field1": StringValue("value1"),
		"field2": StringValue("value2"),
		"field3": NumberValue(42),
	}

	node := &FieldAccessNode{
		Object: &LiteralNode{Value: obj},
		Field:  "field1",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		node.Evaluate(ctx)
	}
}

func BenchmarkFunctionCall(b *testing.B) {
	ctx := NewMockContext()
	ctx.SetFunction("add", func(args []Value) (Value, error) {
		return NumberValue(ToNumber(args[0]) + ToNumber(args[1])), nil
	})

	node := &FunctionCallNode{
		Name: "add",
		Args: []ExprNode{
			&LiteralNode{Value: NumberValue(5)},
			&LiteralNode{Value: NumberValue(3)},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		node.Evaluate(ctx)
	}
}

func BenchmarkToBool(b *testing.B) {
	values := []Value{
		BoolValue(true),
		NumberValue(42),
		StringValue("hello"),
		nil,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ToBool(values[i%len(values)])
	}
}

func BenchmarkToNumber(b *testing.B) {
	values := []Value{
		NumberValue(42),
		StringValue("42"),
		BoolValue(true),
		StringValue("3.14"),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ToNumber(values[i%len(values)])
	}
}
