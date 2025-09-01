package exql

import (
	"github.com/vedadiyan/exql/lang"
	"github.com/vedadiyan/exql/lib/crypt"
	"github.com/vedadiyan/exql/lib/http"
	"github.com/vedadiyan/exql/lib/ip"
	"github.com/vedadiyan/exql/lib/json"
	"github.com/vedadiyan/exql/lib/list"
	maps "github.com/vedadiyan/exql/lib/map"
	str "github.com/vedadiyan/exql/lib/string"
	"github.com/vedadiyan/exql/lib/time"
	"github.com/vedadiyan/exql/lib/url"
	"github.com/vedadiyan/exql/lib/util"
)

type DefaultContext struct {
	values map[string]lang.Value
	funcs  map[string]lang.Function
}

type DefaultContextOption func(*DefaultContext)

func WithFunctions(funcs map[string]lang.Function) DefaultContextOption {
	return func(dc *DefaultContext) {
		dc.funcs = funcs
	}
}

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

func Exports() map[string]lang.Value {
	return map[string]lang.Value{
		"crypt":  NewDefaultContext(WithFunctions(crypt.Export())),
		"http":   NewDefaultContext(WithFunctions(http.Export())),
		"ip":     NewDefaultContext(WithFunctions(ip.Export())),
		"json":   NewDefaultContext(WithFunctions(json.Export())),
		"list":   NewDefaultContext(WithFunctions(list.Export())),
		"map":    NewDefaultContext(WithFunctions(maps.Export())),
		"string": NewDefaultContext(WithFunctions(str.Export())),
		"time":   NewDefaultContext(WithFunctions(time.Export())),
		"url":    NewDefaultContext(WithFunctions(url.Export())),
		"util":   NewDefaultContext(WithFunctions(util.Export())),
	}
}
