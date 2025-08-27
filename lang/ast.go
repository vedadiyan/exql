package lang

import (
	"fmt"
	"strconv"
)

type (
	Value       interface{}
	BoolValue   bool
	StringValue string
	NumberValue float64
	ListValue   []Value
	MapValue    map[string]Value
	ExprNode    interface {
		Evaluate(ctx Context) Value
	}
	Context interface {
		GetVariable(name string) Value
		GetFunction(name string) Function
	}
	Function interface {
		Call(args []Value) Value
	}
	BinaryOpNode struct {
		Left, Right ExprNode
		Operator    string
	}
	UnaryOpNode struct {
		Operand  ExprNode
		Operator string
	}
	LiteralNode struct {
		Value Value
	}
	VariableNode struct {
		Name string
	}
	FieldAccessNode struct {
		Object ExprNode
		Field  string
	}
	IndexAccessNode struct {
		Object ExprNode
		Index  ExprNode
	}
	FunctionCallNode struct {
		Name string
		Args []ExprNode
	}
	ListNode struct {
		Elements []ExprNode
	}
	EachNode  struct{}
	EachValue Value
)

func (n *BinaryOpNode) Evaluate(ctx Context) Value {
	left := n.Left.Evaluate(ctx)
	right := n.Right.Evaluate(ctx)

	switch n.Operator {
	case "and":
		return BoolValue(toBool(left) && toBool(right))
	case "or":
		return BoolValue(toBool(left) || toBool(right))
	case "=", "==":
		return BoolValue(equal(left, right))
	case "!=":
		return BoolValue(!equal(left, right))
	case "<":
		return BoolValue(compare(left, right) < 0)
	case "<=":
		return BoolValue(compare(left, right) <= 0)
	case ">":
		return BoolValue(compare(left, right) > 0)
	case ">=":
		return BoolValue(compare(left, right) >= 0)
	case "in":
		return BoolValue(contains(right, left))
	case "not in":
		return BoolValue(!contains(right, left))
	case "+":
		return NumberValue(toNumber(left) + toNumber(right))
	case "-":
		return NumberValue(toNumber(left) - toNumber(right))
	case "*":
		return NumberValue(toNumber(left) * toNumber(right))
	case "/":
		return NumberValue(toNumber(left) / toNumber(right))
	}
	return BoolValue(false)
}

func (n *UnaryOpNode) Evaluate(ctx Context) Value {
	operand := n.Operand.Evaluate(ctx)
	switch n.Operator {
	case "not":
		return BoolValue(!toBool(operand))
	case "-":
		return NumberValue(-toNumber(operand))
	}
	return operand
}

func (n *LiteralNode) Evaluate(ctx Context) Value {
	return n.Value
}

func (n *VariableNode) Evaluate(ctx Context) Value {
	return ctx.GetVariable(n.Name)
}

func (n *FieldAccessNode) Evaluate(ctx Context) Value {
	obj := n.Object.Evaluate(ctx)
	return n.evaluate(obj)
}

func (n *FieldAccessNode) evaluate(obj Value) Value {
	switch obj := obj.(type) {
	case ListValue:
		{
			values := make(ListValue, 0)
			for _, i := range obj {
				values = append(values, n.evaluate(i))
			}
			return values
		}
	case MapValue:
		{
			return obj[n.Field]
		}
	default:
		{
			return nil
		}
	}
}

func (n *IndexAccessNode) Evaluate(ctx Context) Value {
	obj := n.Object.Evaluate(ctx)
	index := n.Index.Evaluate(ctx)

	switch obj := obj.(type) {
	case MapValue:
		{
			if strIndex, ok := index.(StringValue); ok {
				expr := new(FieldAccessNode)
				expr.Field = string(strIndex)
				expr.Object = n.Object
				return expr.Evaluate(ctx)
			}
			return nil
		}
	case ListValue:
		{
			switch index := index.(type) {
			case NumberValue:
				{
					idx := int(index)
					if idx >= 0 && idx < len(obj) {
						return obj[idx]
					}
					return nil
				}
			case StringValue:
				{
					expr := new(FieldAccessNode)
					expr.Field = string(index)
					expr.Object = n.Object
					return expr.Evaluate(ctx)
				}
			case EachValue:
				{
					return obj
				}
			default:
				{
					return nil
				}
			}
		}
	default:
		{
			return nil
		}
	}
}

func (n *FunctionCallNode) Evaluate(ctx Context) Value {
	fn := ctx.GetFunction(n.Name)
	if fn == nil {
		return BoolValue(false)
	}

	args := make([]Value, len(n.Args))
	for i, arg := range n.Args {
		args[i] = arg.Evaluate(ctx)
	}

	return fn.Call(args)
}

func (n *ListNode) Evaluate(ctx Context) Value {
	elements := make([]Value, len(n.Elements))
	for i, elem := range n.Elements {
		elements[i] = elem.Evaluate(ctx)
	}
	return ListValue(elements)
}

func (n *EachNode) Evaluate(ctx Context) Value {
	return EachValue(0)
}

func toBool(v Value) bool {
	switch val := v.(type) {
	case BoolValue:
		return bool(val)
	case NumberValue:
		return val != 0
	case StringValue:
		return val != ""
	default:
		return v != nil
	}
}

func toNumber(v Value) float64 {
	switch val := v.(type) {
	case NumberValue:
		return float64(val)
	case StringValue:
		if f, err := strconv.ParseFloat(string(val), 64); err == nil {
			return f
		}
	case BoolValue:
		if val {
			return 1
		}
		return 0
	}
	return 0
}

func equal(a, b Value) bool {
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}

func compare(a, b Value) int {
	aNum, bNum := toNumber(a), toNumber(b)
	if aNum < bNum {
		return -1
	}
	if aNum > bNum {
		return 1
	}
	return 0
}

func contains(container, item Value) bool {
	if listVal, ok := container.(ListValue); ok {
		for _, elem := range listVal {
			if equal(elem, item) {
				return true
			}
		}
	}
	return false
}
