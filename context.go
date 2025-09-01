package exql

import "github.com/vedadiyan/exql/lang"

type DefaultContext struct {
	values map[string]lang.Value
	funcs  map[string]lang.Function
}

type DefaultContextOption func(*DefaultContext)

func NewDefaultContext(opts ...DefaultContextOption) *DefaultContext {
	out := new(DefaultContext)

	for _, opt := range opts {
		opt(out)
	}

	if out.values == nil {
		out.values = make(map[string]lang.Value)
	}
	if out.funcs == nil {
		out.funcs = make(map[string]lang.Function)
	}
	return out
}

func (c *DefaultContext) SetVariable(name string, value lang.Value) {
	c.values[name] = value
}
func (c *DefaultContext) SetFunction(name string, function lang.Function) {
	c.funcs[name] = function
}

func (c *DefaultContext) GetVariable(name string) lang.Value {
	return c.values[name]
}
func (c *DefaultContext) GetFunction(name string) lang.Function {
	return c.funcs[name]
}
