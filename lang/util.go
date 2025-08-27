package lang

type DefaultContext struct {
	values map[string]Value
	funcs  map[string]Function
}

type DefaultContextOption func(*DefaultContext)

func NewDefaultContext(opts ...DefaultContextOption) *DefaultContext {
	out := new(DefaultContext)

	for _, opt := range opts {
		opt(out)
	}

	if out.values == nil {
		out.values = make(map[string]Value)
	}
	if out.funcs == nil {
		out.funcs = make(map[string]Function)
	}
	return out
}

func (c *DefaultContext) SetVariable(name string, value Value) {
	c.values[name] = value
}
func (c *DefaultContext) SetFunction(name string, function Function) {
	c.funcs[name] = function
}

func (c *DefaultContext) GetVariable(name string) Value {
	return c.values[name]
}
func (c *DefaultContext) GetFunction(name string) Function {
	return c.funcs[name]
}
