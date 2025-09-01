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
	EachValue   Value
	ExprNode    interface {
		Evaluate(ctx Context) (Value, error)
	}
	Context interface {
		GetVariable(name string) Value
		GetFunction(name string) Function
	}
	Function     func(args []Value) (Value, error)
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
		Namespace ExprNode
		Name      string
		Args      []ExprNode
	}
	ListNode struct {
		Elements []ExprNode
	}
	EachNode  struct{}
	RangeNode struct {
		Begin Value
		End   Value
	}
)

func (n *BinaryOpNode) Evaluate(ctx Context) (Value, error) {
	left, err := n.Left.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	right, err := n.Right.Evaluate(ctx)
	if err != nil {
		return nil, err
	}

	switch n.Operator {
	case "and":
		return BoolValue(ToBool(left) && ToBool(right)), nil
	case "or":
		return BoolValue(ToBool(left) || ToBool(right)), nil
	case "=", "==":
		return BoolValue(equal(left, right)), nil
	case "!=":
		return BoolValue(!equal(left, right)), nil
	case "<":
		return BoolValue(compare(left, right) < 0), nil
	case "<=":
		return BoolValue(compare(left, right) <= 0), nil
	case ">":
		return BoolValue(compare(left, right) > 0), nil
	case ">=":
		return BoolValue(compare(left, right) >= 0), nil
	case "in":
		return BoolValue(contains(right, left)), nil
	case "not in":
		return BoolValue(!contains(right, left)), nil
	case "+":
		return NumberValue(ToNumber(left) + ToNumber(right)), nil
	case "-":
		return NumberValue(ToNumber(left) - ToNumber(right)), nil
	case "*":
		return NumberValue(ToNumber(left) * ToNumber(right)), nil
	case "/":
		return NumberValue(ToNumber(left) / ToNumber(right)), nil
	}
	return nil, fmt.Errorf("expectation failed: %s not supported", n.Operator)
}

func (n *UnaryOpNode) Evaluate(ctx Context) (Value, error) {
	operand, err := n.Operand.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	switch n.Operator {
	case "not":
		return BoolValue(!ToBool(operand)), nil
	case "-":
		return NumberValue(-ToNumber(operand)), nil
	}
	return nil, fmt.Errorf("expectation failed: %s not supported", n.Operator)
}

func (n *LiteralNode) Evaluate(ctx Context) (Value, error) {
	return n.Value, nil
}

func (n *VariableNode) Evaluate(ctx Context) (Value, error) {
	return ctx.GetVariable(n.Name), nil
}

func (n *FieldAccessNode) Evaluate(ctx Context) (Value, error) {
	obj, err := n.Object.Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	return n.evaluate(obj), nil
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

func (n *IndexAccessNode) Evaluate(ctx Context) (Value, error) {
	obj, err := n.Object.Evaluate(ctx)
	if err != nil {
		return nil, err
	}

	index, err := n.Index.Evaluate(ctx)
	if err != nil {
		return nil, err
	}

	switch obj := obj.(type) {
	case MapValue:
		{
			if strIndex, ok := index.(StringValue); ok {
				expr := new(FieldAccessNode)
				expr.Field = string(strIndex)
				expr.Object = n.Object
				return expr.Evaluate(ctx)
			}
			return nil, fmt.Errorf("expectation failed: %T not supported", index)
		}
	case ListValue:
		{
			switch index := index.(type) {
			case NumberValue:
				{
					idx := int(index)
					if idx >= 0 && idx < len(obj) {
						return obj[idx], nil
					}
					return nil, fmt.Errorf("expectation failed: index %d is out of range", idx)
				}
			case StringValue:
				{
					expr := new(FieldAccessNode)
					expr.Field = string(index)
					expr.Object = n.Object
					return expr.Evaluate(ctx)
				}
			case BoolValue:
				{
					return nil, fmt.Errorf("expectation failed: %T not supported", index)
				}
			case EachValue:
				{
					return obj, nil
				}
			default:
				{
					return nil, fmt.Errorf("expectation failed: %T not supported", index)
				}
			}
		}
	default:
		{
			return nil, fmt.Errorf("expectation failed: %T not supported", obj)
		}
	}
}

func (n *FunctionCallNode) Evaluate(ctx Context) (Value, error) {
	namespace := ctx
	if n.Namespace != nil {
		value, err := n.Namespace.Evaluate(ctx)
		if err != nil {
			return nil, err
		}
		n, ok := value.(Context)
		if !ok {
			return nil, fmt.Errorf("unexpected identifier %v", value)
		}
		namespace = n
	}
	fn := namespace.GetFunction(n.Name)
	if fn == nil {
		return BoolValue(false), nil
	}

	args := make([]Value, len(n.Args))
	for i, arg := range n.Args {
		val, err := arg.Evaluate(ctx)
		if err != nil {
			return nil, err
		}
		args[i] = val
	}

	return fn(args)
}

func (n *ListNode) Evaluate(ctx Context) (Value, error) {
	elements := make([]Value, len(n.Elements))
	for i, elem := range n.Elements {
		val, err := elem.Evaluate(ctx)
		if err != nil {
			return nil, err
		}
		elements[i] = val
	}
	return ListValue(elements), nil
}

func (n *EachNode) Evaluate(ctx Context) (Value, error) {
	return EachValue(0), nil
}

func (n *RangeNode) Evaluate(ctx Context) (Value, error) {
	return EachValue(0), nil
}

func ToBool(v Value) bool {
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

func ToNumber(v Value) float64 {
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
	return fmt.Sprintf("%T %v", a, a) == fmt.Sprintf("%T %v", b, b)
}

func compare(a, b Value) int {
	aNum, bNum := ToNumber(a), ToNumber(b)
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
