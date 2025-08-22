package main

import (
	"fmt"

	"github.com/vedadiyan/exql/lang"
)

type Context struct {
	vars  map[string]lang.Value
	funcs map[string]lang.Function
}

func main() {
	input := "role = 'admin' and not age <= 18"
	context := new(Context)
	context.vars = make(map[string]lang.Value)
	context.vars["role"] = "admin"
	context.vars["age"] = lang.NumberValue(20)
	context.funcs = make(map[string]lang.Function)

	expr, err := lang.ParseExpression(input)
	if err != nil {
		panic(err)
	}

	fmt.Println(expr.Evaluate(context))
}

func (c *Context) GetVariable(name string) lang.Value {
	return c.vars[name]
}
func (c *Context) GetFunction(name string) lang.Function {
	return c.funcs[name]
}
