package exql

import (
	"github.com/vedadiyan/exql/lang"
)

func Parse(expr string) (lang.ExprNode, error) {

	return lang.ParseExpression(expr)
}

func Eval(expr string, context lang.Context) (lang.Value, error) {
	result, err := Parse(expr)
	if err != nil {
		return nil, err
	}
	return result.Evaluate(context)
}
