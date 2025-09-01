package exql

import "testing"

func TestParser(t *testing.T) {
	ctx := NewDefaultContext()
	ctx.values = Exports()
	Eval("url.parse('test')", ctx)
}
